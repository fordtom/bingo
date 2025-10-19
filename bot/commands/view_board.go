package commands

import "github.com/bwmarrin/discordgo"

// ViewBoard returns the view_board subcommand definition
func ViewBoard() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "view_board",
		Description: "Display a player's bingo board",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The Discord user whose board to display",
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

// HandleViewBoard processes the view_board command
func HandleViewBoard(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, db interface{}) {
	// TODO: Implement
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "view_board command not yet implemented",
		},
	})
}
