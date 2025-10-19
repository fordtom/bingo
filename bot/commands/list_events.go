package commands

import "github.com/bwmarrin/discordgo"

// ListEvents returns the list_events subcommand definition
func ListEvents() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "list_events",
		Description: "List all events for a game with their status and vote counts",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "game_id",
				Description: "ID of the game to list events for (uses active game if not provided)",
				Required:    false,
			},
		},
	}
}

// HandleListEvents processes the list_events command
func HandleListEvents(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, db interface{}) {
	// TODO: Implement
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "list_events command not yet implemented",
		},
	})
}
