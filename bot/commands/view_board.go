package commands

import (
	"bytes"
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

// ViewBoard returns the view_board subcommand definition
func ViewBoard() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "view_board",
		Description: "Display a player's bingo board",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "The Discord user whose board to display",
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

// HandleViewBoard processes the view_board command
func HandleViewBoard(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.DB) {
	ctx := context.Background()

	// Parse user
	var userID int64
	var userSnowflake string
	for _, opt := range options {
		if opt.Name == "user" {
			userSnowflake = opt.UserValue(s).ID
			userID = parseUserID(userSnowflake)
			break
		}
	}

	// Get game ID
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

	// Get board
	board, squares, err := database.GetUserBoard(ctx, gameID, userID)
	if err != nil {
		respondError(s, i, "Error fetching board: "+err.Error())
		return
	}
	if board == nil {
		respondError(s, i, fmt.Sprintf("No board found for <@%d> in game #%d.", userID, gameID))
		return
	}

	// Build grid from squares
	gridSize := board.GridSize
	grid := make([][]db.BoardSquareWithEvent, gridSize)
	for i := range grid {
		grid[i] = make([]db.BoardSquareWithEvent, gridSize)
	}
	for _, sq := range squares {
		grid[sq.Row][sq.Column] = sq
	}

	// Generate board image
	imageBytes, err := GenerateBoardImage(grid, gridSize)
	if err != nil {
		respondError(s, i, "Error generating board image: "+err.Error())
		return
	}

	// Create filename and message
	displayName := userDisplayName(s, i.GuildID, userSnowflake)
	filename := fmt.Sprintf("board_game%d_user%d.png", gameID, userID)
	message := fmt.Sprintf("Board for %s — Game #%d: %s", displayName, gameID, game.Title)

	// Respond to interaction first
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		respondError(s, i, "Error responding: "+err.Error())
		return
	}

	// Send image as follow-up
	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Files: []*discordgo.File{
			{
				Name:        filename,
				ContentType: "image/png",
				Reader:      bytes.NewReader(imageBytes),
			},
		},
	})
	if err != nil {
		// Log error but don't fail - message was already sent
		fmt.Printf("Error sending board image: %v\n", err)
	}
}
