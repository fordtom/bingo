package commands

import "github.com/bwmarrin/discordgo"

// NewGame returns the new_game subcommand definition
func NewGame() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "new_game",
		Description: "Create a complete bingo game with events and player boards",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "Name for the bingo game",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "grid_size",
				Description: "Size of the grid (2-10, typically 3, 4, or 5)",
				Required:    true,
				MinValue:    floatPtr(2),
				MaxValue:    10,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_ids",
				Description: "Discord users who will participate (@player1 @player2 ...)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "events_csv",
				Description: "CSV file containing event descriptions",
				Required:    true,
			},
		},
	}
}

// HandleNewGame processes the new_game command
func HandleNewGame(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	// TODO: Implement
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "new_game command not yet implemented",
		},
	})
}
