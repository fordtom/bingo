package commands

import "github.com/bwmarrin/discordgo"

// DeleteGame returns the delete_game subcommand definition
func DeleteGame() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "delete_game",
		Description: "Delete a game and all associated data",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "game_id",
				Description: "ID of the game to delete",
				Required:    true,
			},
		},
	}
}

// HandleDeleteGame processes the delete_game command
func HandleDeleteGame(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	// TODO: Implement
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "delete_game command not yet implemented",
		},
	})
}
