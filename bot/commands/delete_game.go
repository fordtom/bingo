package commands

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

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
func HandleDeleteGame(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.DB) {
	ctx := context.Background()

	gameID, ok := getIntOption(options, "game_id")
	if !ok {
		respondError(s, i, "Missing required game_id option.")
		return
	}

	// Check if game exists and if it's active
	game, err := database.GetGame(ctx, gameID)
	if err != nil {
		respondError(s, i, "Error checking game: "+err.Error())
		return
	}
	if game == nil {
		respondError(s, i, fmt.Sprintf("Game #%d not found.", gameID))
		return
	}

	wasActive := game.IsActive

	// Delete the game and all associated data
	if err := database.DeleteGameCascade(ctx, gameID); err != nil {
		respondError(s, i, "Error deleting game: "+err.Error())
		return
	}

	response := fmt.Sprintf("✓ Game #%d (**%s**) deleted.", gameID, game.Title)

	// If we deleted the active game, set a new one
	if wasActive {
		games, err := database.ListGames(ctx)
		if err == nil && len(games) > 0 {
			// Set lowest ID game as active
			minID := games[0].ID
			for _, g := range games {
				if g.ID < minID {
					minID = g.ID
				}
			}
			if err := database.SetActiveGame(ctx, minID); err == nil {
				response += fmt.Sprintf("\nGame #%d is now active.", minID)
			}
		} else {
			response += "\nNo remaining games."
		}
	}

	titleText := "Game Deleted"
	respondEmbed(s, i, titleText, response, colorSuccess, false)
}
