package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nicolito128/kantele/discord"
	"github.com/nicolito128/kantele/rest"
)

const (
	GatewayURL = "wss://gateway.discord.gg/?v=10&encoding=json"
)

type Gateway struct {
	token             string
	sequence          int
	conn              *websocket.Conn
	restClient        *rest.Client
	heartbeatInterval time.Duration
}

func New(token string) *Gateway {
	return &Gateway{
		token:      token,
		restClient: rest.NewClient(token),
	}
}

func (gtw *Gateway) Open() {
	go gtw.connect()
}

func (gtw *Gateway) Client() *rest.Client {
	return gtw.restClient
}

func (gtw *Gateway) connect() {
	fmt.Println("Openning gateway connection...")

	conn, res, err := websocket.DefaultDialer.Dial(GatewayURL, nil)
	if err != nil {
		log.Printf("handshake failed with status %d", res.StatusCode)
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	gtw.conn = conn

outer:
	for {
		typ, b, err := gtw.conn.ReadMessage()
		if err != nil {
			log.Fatal("read message:", err)
		}

		switch typ {
		case websocket.CloseMessage:
			log.Println("Closing connection...")
			break outer

		case websocket.TextMessage:
			msg := map[string]any{}
			if err := json.Unmarshal(b, &msg); err != nil {
				panic(err)
			}

			gtw.handlePayload(gtw.token, msg)

		default:
			log.Println("\nconn message:", typ, string(b))
		}
	}

	log.Println("Connection finished!")
}

func (gtw *Gateway) handlePayload(token string, payload map[string]any) {
	opcode, data, seq, t := payload["op"], payload["d"], payload["s"], payload["t"]
	if v, ok := seq.(int); ok {
		gtw.sequence = v
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
		gtw.identify(token)
		go gtw.heartbeat()

	case 11: // Heartbeat ACK
		log.Println("Heartbeat acknowledged")

	case 0: // Dispatch event
		go gtw.handleEvent(eventName, data)

	default:
		log.Println("Unhandled payload:", payload)
	}
}

func (gtw *Gateway) handleEvent(eventName string, data any) {
	fmt.Println("\nEvent received:", eventName, data)

	switch eventName {
	case "READY":
		log.Println("\nBot is ready!")

	case "MESSAGE_CREATE":
		message, ok := data.(map[string]any)
		if !ok {
			panic("message shuld be parsed")
		}

		content := message["content"].(string)

		if strings.HasPrefix(content, "!ping") {
			channelId := message["channel_id"].(string)
			gtw.restClient.Post(fmt.Sprintf("/channels/%s/messages", channelId), struct {
				Content string `json:"content"`
			}{
				Content: "pong!",
			})
		}

	default:
		fmt.Println("\nUnhandled event:", eventName, data)
	}
}

func (gtw *Gateway) identify(token string) {
	log.Println("Identifying...")

	idPayload := discord.Identify{
		Payload: discord.Payload[discord.IdentifyData]{
			Op: 2,
			D: discord.IdentifyData{
				Token:   token,
				Intents: 33280,
				Properties: discord.IdentifyProperties{
					OS:      "win32",
					Browser: "kantele",
					Device:  "kantele",
				},
			},
		},
	}

	if err := gtw.conn.WriteJSON(idPayload); err != nil {
		panic(err)
	}
}

func (gtw *Gateway) heartbeat() {
	fmt.Println("Starting heartbeat...")

	heartbeatTicker := time.NewTicker(gtw.heartbeatInterval)
	defer heartbeatTicker.Stop()

	for range heartbeatTicker.C {
		heartb := discord.Heartbeat{
			Payload: discord.Payload[int]{
				Op: 1,
				D:  gtw.sequence,
			},
		}

		if err := gtw.conn.WriteJSON(heartb); err != nil {
			panic(err)
		}
	}
}
