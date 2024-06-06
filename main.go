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

	"github.com/ammuench/rerolled-bot/internal/db"
)

const (
	discordBotEnvKey  = "DISCORD_BOT_KEY"
	supabaseURLEnvKey = "SUPABASE_URL"
	supabaseKeyEnvKey = "SUPABASE_KEY"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Unable to load .env file")
	}

	discordBotKey, dbKeyExists := os.LookupEnv(discordBotEnvKey)
	if !dbKeyExists {
		log.Fatal("No Discord Bot Key in .env file")
	}

	tursoDbInstance, err := db.InitDB()
	if err != nil {
		log.Fatal("Error initializing turso db: ", err)
	}

	fmt.Println("Connected to turso db...")

	discordBot, err := discordgo.New("Bot " + discordBotKey)
	if err != nil {
		log.Fatal("Error creating discord session: ", err)
	}

	fmt.Println("Rerolled-Bot is started...")

	discordBot.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == "!ping" {
			s.ChannelMessageSend(m.ChannelID, "!pong @ "+time.Now().String())
		}
	})

	discordBot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = discordBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Rerolled-Bot is now ~*running*~.")
	fmt.Println("Press CTRL-C to exit.")

	// Register Handlers
	// discordBot.ApplicationCommandCreate()

	// Register close signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discordBot.Close()
	tursoDbInstance.Close()
}
