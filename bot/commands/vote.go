package commands

import (
	"context"
	"fmt"
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

// Vote returns the vote subcommand definition
func Vote() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "vote",
		Description: "Vote that an event has occurred",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "event_id",
				Description: "ID of the event to vote for",
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

// HandleVote processes the vote command
func HandleVote(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.DB) {
	ctx := context.Background()
	userID := parseUserID(i.Member.User.ID)

	// Parse options
	displayID := int(options[0].IntValue())

	// Get game ID (either specified or active)
	gameID, err := getGameIDOrActive(ctx, database, options, "game_id")
	if err != nil {
		respondError(s, i, err.Error())
		return
	}

	// Look up event by display_id
	event, err := database.GetEventByDisplayID(ctx, gameID, displayID)
	if err != nil {
		respondError(s, i, "Error fetching event: "+err.Error())
		return
	}
	if event == nil {
		respondError(s, i, fmt.Sprintf("Event #%d not found in the current game.", displayID))
		return
	}

	// Check if event is already closed
	if event.Status == "CLOSED" {
		respondError(s, i, fmt.Sprintf("Event #%d has already been marked as occurred.", displayID))
		return
	}

	// Check if user already voted
	hasVoted, err := database.HasUserVoted(ctx, event.ID, userID)
	if err != nil {
		respondError(s, i, "Error checking vote status: "+err.Error())
		return
	}
	if hasVoted {
		respondError(s, i, fmt.Sprintf("You have already voted for event #%d.", displayID))
		return
	}

	// Record the vote
	if err := database.CreateVote(ctx, event.ID, userID); err != nil {
		respondError(s, i, "Error recording vote: "+err.Error())
		return
	}

	// Get updated vote count
	voteCount, err := database.GetVoteCount(ctx, event.ID)
	if err != nil {
		respondError(s, i, "Vote recorded, but error checking vote count: "+err.Error())
		return
	}

	// Calculate threshold based on player count
	playerCount, err := database.GetPlayerCountForGame(ctx, gameID)
	if err != nil {
		respondError(s, i, "Vote recorded, but error fetching player count: "+err.Error())
		return
	}

	threshold := playerCount
	if playerCount > 3 {
		threshold = int(math.Ceil(0.6 * float64(playerCount)))
	}

	response := fmt.Sprintf("✓ Voted for event #%d: **%s**\nCurrent votes: %d/%d", displayID, event.Description, voteCount, threshold)

	// Close event if threshold reached
	if voteCount >= threshold {
		if err := database.UpdateEventStatus(ctx, event.ID, "CLOSED"); err != nil {
			respondError(s, i, "Vote recorded, but error closing event: "+err.Error())
			return
		}
		response += "\n🎉 Event has been marked as occurred!"

		// Check for winners
		winners, err := checkWinners(ctx, database, gameID)
		if err != nil {
			response += "\n(Error checking for winners: " + err.Error() + ")"
		} else if len(winners) > 0 {
			response += "\n\n🏆 **BINGO!** Winners: "
			for idx, winnerID := range winners {
				if idx > 0 {
					response += ", "
				}
				response += fmt.Sprintf("<@%d>", winnerID)
			}
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// checkWinners checks all boards for bingo (row, column, or diagonal)
func checkWinners(ctx context.Context, database *db.DB, gameID int64) ([]int64, error) {
	playerIDs, err := database.GetGamePlayerIDs(ctx, gameID)
	if err != nil {
		return nil, err
	}

	var winners []int64
	for _, playerID := range playerIDs {
		if hasWon, err := checkPlayerWin(ctx, database, gameID, playerID); err != nil {
			return nil, err
		} else if hasWon {
			winners = append(winners, playerID)
		}
	}

	return winners, nil
}

// checkPlayerWin checks if a single player has won
func checkPlayerWin(ctx context.Context, database *db.DB, gameID, playerID int64) (bool, error) {
	board, squares, err := database.GetUserBoard(ctx, gameID, playerID)
	if err != nil {
		return false, err
	}
	if board == nil {
		return false, nil
	}

	// Build grid
	gridSize := board.GridSize
	grid := make([][]bool, gridSize)
	for i := range grid {
		grid[i] = make([]bool, gridSize)
	}
	for _, sq := range squares {
		grid[sq.Row][sq.Column] = (sq.EventStatus == "CLOSED")
	}

	// Check rows
	for row := 0; row < gridSize; row++ {
		allClosed := true
		for col := 0; col < gridSize; col++ {
			if !grid[row][col] {
				allClosed = false
				break
			}
		}
		if allClosed {
			return true, nil
		}
	}

	// Check columns
	for col := 0; col < gridSize; col++ {
		allClosed := true
		for row := 0; row < gridSize; row++ {
			if !grid[row][col] {
				allClosed = false
				break
			}
		}
		if allClosed {
			return true, nil
		}
	}

	// Check diagonal (top-left to bottom-right)
	allClosed := true
	for i := 0; i < gridSize; i++ {
		if !grid[i][i] {
			allClosed = false
			break
		}
	}
	if allClosed {
		return true, nil
	}

	// Check diagonal (top-right to bottom-left)
	allClosed = true
	for i := 0; i < gridSize; i++ {
		if !grid[i][gridSize-1-i] {
			allClosed = false
			break
		}
	}
	if allClosed {
		return true, nil
	}

	return false, nil
}
