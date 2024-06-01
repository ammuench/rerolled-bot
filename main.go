package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
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

	discordBot, err := discordgo.New("Bot " + discordBotKey)
	if err != nil {
		log.Fatal("Error creating discord session: ", err)
	}

	fmt.Println("Rerolled-Bot is started...")

	discordBot.AddHandler(func (s* discordgo.Session, m* discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "!ping" {
			s.ChannelMessageSend(m.ChannelID, "!pong @ " + time.Now().String())
		}

	})

	discordBot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = discordBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Rerolled-Bot is now ~*running*~. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discordBot.Close()
}
