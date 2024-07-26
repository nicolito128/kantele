package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/nicolito128/kantele/gateway"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("BOT_TOKEN")

	client := gateway.New(token)

	client.HandleEvent("MESSAGE_CREATE", func(data any) {
		message, ok := data.(map[string]any)
		if !ok {
			panic("message should be parsed")
		}

		content := message["content"].(string)

		if strings.HasPrefix(content, "!ping") {
			channelId := message["channel_id"].(string)
			client.Rest().Post(fmt.Sprintf("/channels/%s/messages", channelId), struct {
				Content string `json:"content"`
			}{
				Content: "pong!",
			})
		}
	})

	client.HandleEvent("READY", func(a any) {
		res, err := client.Rest().Get("/users/@me")
		if err != nil {
			log.Fatal("rest error: ", err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal("response body read error: ", err)
		}

		fmt.Println(string(b))
	})

	// Open connection
	client.Open()

	<-make(chan struct{})
}
