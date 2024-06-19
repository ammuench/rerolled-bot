package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/ammuench/rerolled-bot/internal/db"
	"github.com/ammuench/rerolled-bot/internal/discord"
)

const (
	discordBotEnvKey = "DISCORD_BOT_KEY"
	// TODO: REMOVE THIS AFTER DEV
	guildID = 1246302013860483142
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

	tursoDBInstance, err := db.InitDB()
	if err != nil {
		log.Fatal("Error initializing turso db: ", err)
	}

	fmt.Println("Connected to turso db...")

	discordBot, err := discordgo.New("Bot " + discordBotKey)
	if err != nil {
		log.Fatal("Error creating discord session: ", err)
	}

	err = discordBot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	registeredCmds, err := discord.InitializeCommands(discordBot)
	if err != nil {
		log.Fatal("Error creating discord slash commands: ", err)
	}

	fmt.Println("Rerolled-Bot is now ~*O N L I N E*~")
	fmt.Println("Press CTRL-C to exit.")

	// Register close signal
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	discord.TeardownAllCommands(discordBot, registeredCmds)
	discordBot.Close()
	tursoDBInstance.Close()
}
