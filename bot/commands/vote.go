package commands

import (
	"context"
	"fmt"

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
	userID := int64(parseUserID(i.Member.User.ID))

	// Parse options
	displayID := int(options[0].IntValue())
	var gameID int64

	// Get game (either specified or active)
	if len(options) > 1 && options[1].Name == "game_id" {
		gameID = options[1].IntValue()
	} else {
		game, err := database.GetActiveGame(ctx)
		if err != nil {
			respondError(s, i, "Error fetching active game: "+err.Error())
			return
		}
		if game == nil {
			respondError(s, i, "No active game found. Please specify a game_id or set an active game.")
			return
		}
		gameID = game.ID
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

	// TODO: Determine threshold (could be based on number of players with boards)
	// For now, using a simple threshold of 2 votes as example
	threshold := 2

	response := fmt.Sprintf("✓ Voted for event #%d: **%s**\nCurrent votes: %d", displayID, event.Description, voteCount)

	// Close event if threshold reached
	if voteCount >= threshold {
		if err := database.UpdateEventStatus(ctx, event.ID, "CLOSED"); err != nil {
			respondError(s, i, "Vote recorded, but error closing event: "+err.Error())
			return
		}
		response += fmt.Sprintf("\n🎉 Event has been marked as occurred (reached %d votes)!", threshold)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

// Helper to parse Discord snowflake ID to int64
func parseUserID(snowflake string) uint64 {
	var id uint64
	fmt.Sscanf(snowflake, "%d", &id)
	return id
}

// Helper to respond with error message
func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
