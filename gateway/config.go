package gateway

import (
	"log/slog"
	"runtime"

	"github.com/gorilla/websocket"
)

type ConfigOpt func(*Config)

func DefaultConfig() *Config {
	return &Config{
		Logger:         slog.Default(),
		Dialer:         websocket.DefaultDialer,
		URL:            "wss://gateway.discord.gg/?v=10&encoding=json",
		LargeThreshold: 50,
		Intents:        IntentsNone,
		OS:             runtime.GOOS,
		Browser:        "kantele",
		Device:         "kantele",
	}
}

type Config struct {
	// Logger is the Logger of the Gateway. Defaults to slog.Default().
	Logger *slog.Logger
	// Dialer is the websocket.Dialer of the Gateway. Defaults to websocket.DefaultDialer.
	Dialer *websocket.Dialer
	// LargeThreshold is the threshold for the Gateway. Defaults to 50
	// Ref: https://discord.com/developers/docs/topics/gateway-events#identify-identify-structure.
	LargeThreshold int
	// Intents is the Intents for the Gateway. Defaults to IntentsNone.
	Intents Intents
	// Whether this connection supports compression of packets. Defaults to false.
	Compress bool
	// URL is the URL of the Gateway. Defaults to fetch from Discord.
	URL string
	// LastSequence is the last sequence received by the Gateway.
	LastSequence *int
	// OS is the operating system it should send on login. Defaults to runtime.GOOS.
	OS string
	// Browser is the browser it should send on login. Defaults to "kantele".
	Browser string
	// Device is the device it should send on login. Defaults to "kantele".
	Device string
}

func WithIntents(intents Intents) ConfigOpt {
	return func(c *Config) {
		c.Intents = intents
	}
}
