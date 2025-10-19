package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

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
func HandleSetActiveGame(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.DB) {
	ctx := context.Background()
	gameID := options[0].IntValue()

	// Verify game exists
	game, err := database.GetGame(ctx, gameID)
	if err != nil {
		respondError(s, i, "Error checking game: "+err.Error())
		return
	}
	if game == nil {
		respondError(s, i, fmt.Sprintf("Game #%d not found.", gameID))
		return
	}

	// Set as active
	if err := database.SetActiveGame(ctx, gameID); err != nil {
		respondError(s, i, "Error setting active game: "+err.Error())
		return
	}

	titleText := "Active Game Set"
	desc := fmt.Sprintf("✓ Game #%d (**%s**) is now active.", gameID, game.Title)
	respondEmbed(s, i, titleText, desc, colorSuccess, false)
}
