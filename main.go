package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const discordBotKeyRef = "DISCORD_BOT_KEY"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load .env file")
	}

	discordBotKey, dbKeyExists := os.LookupEnv(discordBotKeyRef)
	if !dbKeyExists {
		log.Fatal("No Discord Bot Key in .env file")
	}

	fmt.Printf("Discord bot started with key %v", discordBotKey)

}
