package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

type Gateway struct {
	config *Config

	token string

	conn *websocket.Conn

	heartbeatInterval time.Duration

	events eventHandler
}

func New(token string, opts ...ConfigOpt) *Gateway {
	conf := DefaultConfig()
	for _, opt := range opts {
		opt(conf)
	}

	return &Gateway{
		config: conf,
		token:  token,
		events: make(eventHandler),
	}
}

func (gtw *Gateway) Open() {
	gtw.config.Logger.Info("Openning gateway...")
	go gtw.connect()
}

func (gtw *Gateway) HandleEvent(eventName string, handler func(any)) {
	gtw.events.Append(eventName, handler)
}

func (gtw *Gateway) connect() {
	defer gtw.config.Logger.Info("Connection finished.")

	conn, res, err := gtw.config.Dialer.Dial(gtw.config.URL, nil)
	if err != nil {
		gtw.config.Logger.Error("handshake failed with:", slog.Int("status", res.StatusCode))
		log.Fatal("dial error:", err)
	}
	defer conn.Close()

	gtw.conn = conn

outer:
	for {
		typ, b, err := gtw.conn.ReadMessage()
		if err != nil {
			gtw.config.Logger.Error("Read message error:", slog.Any("err", err))
			break outer
		}

		switch typ {
		case websocket.CloseMessage:
			gtw.config.Logger.Debug("Closing connection...")
			break outer

		case websocket.TextMessage:
			msg := map[string]any{}
			if err := json.Unmarshal(b, &msg); err != nil {
				gtw.config.Logger.Error("read message: %w", slog.Any("err", err))
				break outer
			}

			gtw.handlePayload(msg)

		default:
			gtw.config.Logger.Debug("Conn message:", slog.Int("typ", typ), slog.String("b", string(b)))
		}
	}
}

func (gtw *Gateway) handlePayload(payload map[string]any) {
	opcode, data, seq, t := payload["op"], payload["d"], payload["s"], payload["t"]
	if v, ok := seq.(int); ok {
		gtw.config.LastSequence = &v
	}

	var eventName string
	if v, ok := t.(string); ok {
		eventName = v
	}

	switch op := int(opcode.(float64)); op {
	case 10: // Hello
		m := map[string]any{}
		if v, ok := data.(map[string]any); ok {
			m = v
		}

		gtw.heartbeatInterval = time.Duration(m["heartbeat_interval"].(float64)) * time.Millisecond
		gtw.identify()
		go gtw.heartbeat()

	case 11: // Heartbeat ACK
		gtw.config.Logger.Info("Heartbeat acknowledged")

	case 0: // Dispatch event
		go gtw.handleEvent(eventName, data)

	default:
		gtw.config.Logger.Debug("Unhandled payload:", slog.Any("payload", payload))
	}
}

func (gtw *Gateway) handleEvent(eventName string, data any) {
	gtw.config.Logger.Info("Event received:", slog.Any("eventName", eventName))

	switch eventName {
	case "READY":
		gtw.config.Logger.Info("Gateway connection is ready!")
		gtw.events.Call(eventName, data)

	case "MESSAGE_CREATE":
		gtw.events.Call(eventName, data)

	default:
		fmt.Println("\nUnhandled event:", eventName, data)
		gtw.events.Call(eventName, data)
	}
}

func (gtw *Gateway) identify() {
	gtw.config.Logger.Info("Identifying...")

	identify := Identify{}
	identify.Op = 2
	identify.D = IdentifyData{
		Token:          gtw.token,
		Intents:        gtw.config.Intents,
		Compress:       gtw.config.Compress,
		LargeThreshold: gtw.config.LargeThreshold,
		Properties: IdentifyConnectionProperties{
			OS:      gtw.config.OS,
			Browser: gtw.config.Browser,
			Device:  gtw.config.Device,
		},
	}

	if err := gtw.conn.WriteJSON(identify); err != nil {
		panic(err)
	}
}

func (gtw *Gateway) heartbeat() {
	gtw.config.Logger.Info("Starting heartbeat...")

	heartbeatTicker := time.NewTicker(gtw.heartbeatInterval)
	defer heartbeatTicker.Stop()

	for range heartbeatTicker.C {
		heartb := Heartbeat{}
		heartb.Op = 1
		if gtw.config.LastSequence != nil {
			heartb.D = *gtw.config.LastSequence
		}

		if err := gtw.conn.WriteJSON(heartb); err != nil {
			panic(err)
		}

		gtw.config.Logger.Debug("heartbeat sent")
	}
}
