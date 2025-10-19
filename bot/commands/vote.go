package commands

import "github.com/bwmarrin/discordgo"

// Vote returns the vote subcommand definition
func Vote() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "vote",
		Description: "Vote that an event has occurred",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "event_id",
				Description: "ID of the event to vote for",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "game_id",
				Description: "ID of the game (uses active game if not provided)",
				Required:    false,
			},
		},
	}
}

// HandleVote processes the vote command
func HandleVote(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	// TODO: Implement
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "vote command not yet implemented",
		},
	})
}
