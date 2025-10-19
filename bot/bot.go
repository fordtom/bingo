package bot

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/bot/commands"
)

type Bot struct {
	session   *discordgo.Session
	channelID string
	userID    string
}

func Setup(s *discordgo.Session, channelID string) (*Bot, func(), error) {
	bot := &Bot{
		session:   s,
		channelID: channelID,
		userID:    s.State.User.ID,
	}

	// Register handlers
	s.AddHandler(bot.handleMessageCreate)
	s.AddHandler(bot.handleInteractionCreate)

	// Register slash commands with Discord
	registeredCommands, err := bot.registerCommands()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		bot.cleanupCommands(registeredCommands)
	}

	return bot, cleanup, nil
}

// handleMessageCreate handles read text messages
func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.ChannelID == b.channelID {
		if m.Content == "bingo" {
			s.ChannelMessageSend(m.ChannelID, "thats me boss")
		}
	}
}

// handleInteractionCreate routes slash commands
func (b *Bot) handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	if data.Name != "bg" {
		return
	}

	subCmd := data.Options[0]

	switch subCmd.Name {
	case "new_game":
		commands.HandleNewGame(s, i, subCmd.Options)
	case "delete_game":
		commands.HandleDeleteGame(s, i, subCmd.Options)
	case "set_active_game":
		commands.HandleSetActiveGame(s, i, subCmd.Options)
	case "list_games":
		commands.HandleListGames(s, i, subCmd.Options)
	case "list_events":
		commands.HandleListEvents(s, i, subCmd.Options)
	case "view_board":
		commands.HandleViewBoard(s, i, subCmd.Options)
	case "vote":
		commands.HandleVote(s, i, subCmd.Options)
	case "help":
		commands.HandleHelp(s, i, subCmd.Options)
	}
}

// registerCommands registers commands
func (b *Bot) registerCommands() ([]*discordgo.ApplicationCommand, error) {
	commandDefinitions := commands.All()
	registeredCommands := make([]*discordgo.ApplicationCommand, 0, len(commandDefinitions))

	for _, cmd := range commandDefinitions {
		registered, err := b.session.ApplicationCommandCreate(b.userID, "", cmd)
		if err != nil {
			log.Printf("Error creating command %s: %v", cmd.Name, err)
			return nil, err
		}
		registeredCommands = append(registeredCommands, registered)
		log.Printf("Command %s registered", cmd.Name)
	}
	return registeredCommands, nil
}

// cleanupCommands removes registered commands
func (b *Bot) cleanupCommands(commands []*discordgo.ApplicationCommand) {
	for _, command := range commands {
		err := b.session.ApplicationCommandDelete(b.userID, "", command.ID)
		if err != nil {
			log.Printf("Error deleting command %s: %v", command.Name, err)
		}
	}
}
