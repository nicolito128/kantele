package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

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

	// Open connection
	client.Open()

	// Rest example
	go func() {
		time.Sleep(2 * time.Second)

		res, err := client.Client().Get("/users/@me")
		if err != nil {
			log.Fatal("rest error: ", err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal("response body read error: ", err)
		}

		fmt.Println(string(b))
	}()

	<-make(chan struct{})
}
