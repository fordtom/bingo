package commands

import "github.com/bwmarrin/discordgo"

const Prefix = "bg"

// All returns all command definitions assembled into the /{Prefix} command
func All() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        Prefix,
			Description: "Bingo game commands",
			Options: []*discordgo.ApplicationCommandOption{
				NewGame(),
				DeleteGame(),
				SetActiveGame(),
				ListGames(),
				ListEvents(),
				ViewBoard(),
				Vote(),
				Help(),
			},
		},
	}
}

// floatPtr returns a pointer to a float64 value (helper for MinValue)
func floatPtr(f float64) *float64 {
	return &f
}
