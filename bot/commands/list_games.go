package commands

import "github.com/bwmarrin/discordgo"

// ListGames returns the list_games subcommand definition
func ListGames() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "list_games",
		Description: "List all bingo games with their details and statistics",
	}
}

// HandleListGames processes the list_games command
func HandleListGames(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	// TODO: Implement
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "list_games command not yet implemented",
		},
	})
}
