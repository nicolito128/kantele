package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nicolito128/kantele/gateway"
	"github.com/nicolito128/kantele/rest"
)

func main() {
	// loading env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("BOT_TOKEN")

	// using the kantele gateway
	gtw := gateway.New(token, gateway.WithIntents(gateway.Intents(33280)))

	// using a rest client for HTTP requests to discord api
	restClient := rest.New(token)

	// throw an event when "MESSAGE_CREATE" is received
	gtw.HandleEvent("MESSAGE_CREATE", func(data any) {
		message, ok := data.(map[string]any)
		if !ok {
			panic("message should be parsed")
		}

		content := message["content"].(string)

		if strings.HasPrefix(content, "!ping") {
			channelId := message["channel_id"].(string)
			restClient.Post(fmt.Sprintf("/channels/%s/messages", channelId), struct {
				Content string `json:"content"`
			}{
				Content: "pong!",
			})
		}
	})

	// throw an event when "READY" is received
	gtw.HandleEvent("READY", func(a any) {
		res, err := restClient.Get("/users/@me")
		if err != nil {
			log.Fatal("rest error: ", err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal("response body read error: ", err)
		}

		fmt.Println(string(b))
	})

	// open the gateway connection
	gtw.Open()

	// hold the main goroutine running
	<-make(chan struct{})
}
