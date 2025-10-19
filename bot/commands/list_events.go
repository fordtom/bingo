package commands

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

// ListEvents returns the list_events subcommand definition
func ListEvents() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "list_events",
		Description: "List all events for a game with their status and vote counts",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "game_id",
				Description: "ID of the game to list events for (uses active game if not provided)",
				Required:    false,
			},
		},
	}
}

// HandleListEvents processes the list_events command
func HandleListEvents(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.DB) {
	ctx := context.Background()

	gameID, err := getGameIDOrActive(ctx, database, options, "game_id")
	if err != nil {
		respondError(s, i, err.Error())
		return
	}

	// Get game info
	game, err := database.GetGame(ctx, gameID)
	if err != nil {
		respondError(s, i, "Error fetching game: "+err.Error())
		return
	}

	// Get events
	events, err := database.GetGameEvents(ctx, gameID)
	if err != nil {
		respondError(s, i, "Error fetching events: "+err.Error())
		return
	}

	if len(events) == 0 {
		respondSuccess(s, i, fmt.Sprintf("No events found for game #%d (**%s**).", gameID, game.Title))
		return
	}

	// Calculate threshold
	playerCount, err := database.GetPlayerCountForGame(ctx, gameID)
	if err != nil {
		respondError(s, i, "Error fetching player count: "+err.Error())
		return
	}

	threshold := playerCount
	if playerCount > 3 {
		threshold = int(math.Ceil(0.6 * float64(playerCount)))
	}

	var lines []string
	title := fmt.Sprintf("Events for Game #%d: %s", gameID, game.Title)

	for _, event := range events {
		if event.Status == "CLOSED" {
			lines = append(lines, fmt.Sprintf("**#%d** %s ✅", event.DisplayID, event.Description))
		} else {
			voteCount, err := database.GetVoteCount(ctx, event.ID)
			if err != nil {
				respondError(s, i, "Error fetching vote count: "+err.Error())
				return
			}
			lines = append(lines, fmt.Sprintf("**#%d** %s (%d/%d votes)", event.DisplayID, event.Description, voteCount, threshold))
		}
	}

	desc := strings.Join(lines, "\n")
	respondEmbed(s, i, title, desc, colorInfo, false)
}
