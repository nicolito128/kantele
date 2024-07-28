package gateway

type PayloadEventMessage interface {
	eventMessage()
}

// Ref: https://discord.com/developers/docs/topics/gateway-events#payload-structure
type Payload[DTyp any] struct {
	// Opcode
	Op int `json:"op"`
	// Event data
	D DTyp `json:"d,omitempty"`
	// Sequence number of event used for resuming sessions and heartbeating
	S int `json:"s,omitempty"`
	// Event name
	T string `json:"t,omitempty"`
}

func (Payload[DTyp]) eventMessage() {}

// Ref: https://discord.com/developers/docs/topics/gateway-events#hello
type Hello struct {
	Payload[HelloData]
}

type HelloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

// Ref: https://discord.com/developers/docs/topics/gateway#sending-heartbeats
type Heartbeat struct {
	Payload[int]
}

// Ref: https://discord.com/developers/docs/topics/gateway-events#identify
type Identify struct {
	Payload[IdentifyData]
}

type IdentifyData struct {
	Token          string                       `json:"token"`
	Properties     IdentifyConnectionProperties `json:"properties"`
	Compress       bool                         `json:"compress,omitempty"`
	LargeThreshold int                          `json:"large_threshold,omitempty"`
	Shard          *[2]int                      `json:"shard,omitempty"`
	Intents        Intents                      `json:"intents"`
	//! IMPLEMENT  Presence       PresenceUpdate         `json:"presence,omitempty"`
}

// Ref: https://discord.com/developers/docs/topics/gateway-events#identify-identify-connection-properties
type IdentifyConnectionProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}
