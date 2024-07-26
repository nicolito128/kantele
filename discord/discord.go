package discord

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

type Hello struct {
	Payload[HelloData]
}

type HelloData struct {
	HeartbeatInterval int `json:"heartbeat_interval"`
}

type Heartbeat struct {
	Payload[int]
}

type Identify struct {
	Payload[IdentifyData]
}

type IdentifyData struct {
	Token      string             `json:"token"`
	Intents    int                `json:"intents"`
	Properties IdentifyProperties `json:"properties"`
}

type IdentifyProperties struct {
	OS      string `json:"os"`
	Browser string `json:"browser"`
	Device  string `json:"device"`
}
