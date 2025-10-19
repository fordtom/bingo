package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/bot"
	"github.com/joho/godotenv"
)

func loadEnv() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token, ok := os.LookupEnv("DISCORD_TOKEN")
	if !ok {
		log.Fatal("DISCORD_TOKEN is not set")
	}

	channelID, ok := os.LookupEnv("CHANNEL_ID")
	if !ok {
		log.Fatal("CHANNEL_ID is not set")
	}

	return token, channelID
}

func main() {

	discordToken, channelID := loadEnv()

	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatal("Error creating bot: ", err)
	}

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = session.Open()
	if err != nil {
		log.Fatal("Error opening session: ", err)
	}
	defer session.Close()

	_, cleanup, err := bot.Setup(session, channelID)
	if err != nil {
		log.Fatal("Error setting up bot: ", err)
	}
	defer cleanup()

	log.Println("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("Bot is shutting down")
}
