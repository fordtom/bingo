package commands

import "github.com/bwmarrin/discordgo"

// SetActiveGame returns the set_active_game subcommand definition
func SetActiveGame() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "set_active_game",
		Description: "Set which game is currently active",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "game_id",
				Description: "ID of the game to make active",
				Required:    true,
			},
		},
	}
}

// HandleSetActiveGame processes the set_active_game command
func HandleSetActiveGame(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, db interface{}) {
	// TODO: Implement
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "set_active_game command not yet implemented",
		},
	})
}
