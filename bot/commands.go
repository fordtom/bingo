package bot

import (
	"github.com/bwmarrin/discordgo"
)

// commands defines the commands
func commands() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "bg",
			Description: "Bingo game commands",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "create",
					Description: "Create a new bingo game",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "name",
							Description: "The name of the bingo game",
							Required:    true,
						},
					},
				},
			},
		},
	}
}
