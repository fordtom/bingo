package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/bot"
	"github.com/fordtom/bingo/db"
	"github.com/joho/godotenv"
)

// initLogger configures the global logger to write to stdout and a file.
// It returns a cleanup function that must be deferred to close the file.
func initLogger() (func(), error) {
	path := os.Getenv("LOG_FILE")
	if path == "" {
		path = "bingo.log"
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, f))
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix("bingo ")
	return func() { _ = f.Close() }, nil
}

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
	cleanup, err := initLogger()
	if err != nil {
		log.Fatalf("log setup failed: %v", err)
	}
	defer cleanup()

	discordToken, channelID := loadEnv()

	// Initialize database
	database, err := db.InitDB()
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	defer database.Close()

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

	var botCleanup func()
	_, botCleanup, err = bot.Setup(session, channelID, database)
	if err != nil {
		log.Fatal("Error setting up bot: ", err)
	}
	defer botCleanup()

	log.Println("Bot is running")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("Bot is shutting down")
}
