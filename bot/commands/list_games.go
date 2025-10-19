package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

// ListGames returns the list_games subcommand definition
func ListGames() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "list_games",
		Description: "List all bingo games with their details and statistics",
	}
}

// HandleListGames processes the list_games command
func HandleListGames(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.DB) {
	ctx := context.Background()

	games, err := database.ListGames(ctx)
	if err != nil {
		respondError(s, i, "Error fetching games: "+err.Error())
		return
	}

	if len(games) == 0 {
		respondSuccess(s, i, "No games found. Create one with `/bg new_game`.")
		return
	}

	var lines []string
	title := "Bingo Games"

	for _, game := range games {
		open, closed, err := database.GetEventCounts(ctx, game.ID)
		if err != nil {
			respondError(s, i, "Error fetching event counts: "+err.Error())
			return
		}

		playerCount, err := database.GetPlayerCountForGame(ctx, game.ID)
		if err != nil {
			respondError(s, i, "Error fetching player count: "+err.Error())
			return
		}

		activeMarker := ""
		if game.IsActive {
			activeMarker = " **(active)**"
		}

		line := fmt.Sprintf("**#%d** %s%s\n  %dx%d grid | %d open, %d closed | %d players",
			game.ID, game.Title, activeMarker, game.GridSize, game.GridSize, open, closed, playerCount)
		lines = append(lines, line)
	}

	desc := strings.Join(lines, "\n")
	respondEmbed(s, i, title, desc, colorInfo, false)
}
