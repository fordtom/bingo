package commands

import "github.com/bwmarrin/discordgo"

// Help returns the help subcommand definition
func Help() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "help",
		Description: "Display help information about all bot commands",
	}
}

// HandleHelp processes the help command
func HandleHelp(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	// TODO: Implement
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "help command not yet implemented",
		},
	})
}
