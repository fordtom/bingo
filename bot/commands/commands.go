package commands

import "github.com/bwmarrin/discordgo"

// All returns all command definitions assembled into the /bg command
func All() []*discordgo.ApplicationCommand {
	return []*discordgo.ApplicationCommand{
		{
			Name:        "bg",
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
