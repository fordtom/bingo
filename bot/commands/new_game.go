package commands

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/fordtom/bingo/db"
)

// NewGame returns the new_game subcommand definition
func NewGame() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "new_game",
		Description: "Create a complete bingo game with events and player boards",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "title",
				Description: "Name for the bingo game",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "grid_size",
				Description: "Size of the grid (2-10, typically 3, 4, or 5)",
				Required:    true,
				MinValue:    floatPtr(2),
				MaxValue:    10,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "player_ids",
				Description: "Discord users who will participate (@player1 @player2 ...)",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionAttachment,
				Name:        "events_csv",
				Description: "CSV file containing event descriptions",
				Required:    true,
			},
		},
	}
}

// HandleNewGame processes the new_game command
func HandleNewGame(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption, database *db.DB) {
	ctx := context.Background()

	// Parse options
	var title string
	var gridSize int
	var playerIDsStr string
	var attachmentID string

	for _, opt := range options {
		switch opt.Name {
		case "title":
			title = opt.StringValue()
		case "grid_size":
			gridSize = int(opt.IntValue())
		case "player_ids":
			playerIDsStr = opt.StringValue()
		case "events_csv":
			attachmentID = opt.Value.(string)
		}
	}

	// Parse player IDs from mentions
	playerIDs := parseMentionsToIDs(playerIDsStr)
	if len(playerIDs) == 0 {
		respondError(s, i, "No valid player mentions found. Use @username format.")
		return
	}

	// Fetch and parse CSV
	attachment := i.ApplicationCommandData().Resolved.Attachments[attachmentID]
	events, err := fetchAndParseCSV(attachment.URL)
	if err != nil {
		respondError(s, i, "Error parsing CSV: "+err.Error())
		return
	}

	// Validate event count
	minEvents := gridSize * gridSize
	if len(events) < minEvents {
		respondError(s, i, fmt.Sprintf("CSV must contain at least %d events for a %dx%d grid (found %d).", minEvents, gridSize, gridSize, len(events)))
		return
	}

	// Create game
	gameID, err := database.CreateGame(ctx, title, gridSize)
	if err != nil {
		respondError(s, i, "Error creating game: "+err.Error())
		return
	}

	// Set as active if no active game exists
	activeGame, err := database.GetActiveGame(ctx)
	if err == nil && activeGame == nil {
		database.SetActiveGame(ctx, gameID)
	}

	// Create game data in transaction
	err = database.WithTx(ctx, func(tx *sql.Tx) error {
		// Create events
		eventIDs := make([]int64, len(events))
		for i, desc := range events {
			displayID := i + 1
			eventID, err := database.CreateEvent(ctx, tx, gameID, displayID, desc)
			if err != nil {
				return fmt.Errorf("error creating event: %w", err)
			}
			eventIDs[i] = eventID
		}

		// Distribute events to boards
		boardAssignments := distributeEvents(eventIDs, playerIDs, gridSize)

		// Create boards and squares
		for playerID, assignments := range boardAssignments {
			boardID, err := database.CreateBoard(ctx, tx, gameID, playerID, gridSize)
			if err != nil {
				return fmt.Errorf("error creating board: %w", err)
			}

			squares := make([]db.BoardSquare, 0, len(assignments))
			for _, assign := range assignments {
				squares = append(squares, db.BoardSquare{
					BoardID: boardID,
					Row:     assign.row,
					Column:  assign.col,
					EventID: assign.eventID,
				})
			}

			if err := database.CreateBoardSquares(ctx, tx, boardID, squares); err != nil {
				return fmt.Errorf("error creating board squares: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		respondError(s, i, err.Error())
		return
	}

	titleText := fmt.Sprintf("Game Created: #%d — %s", gameID, title)
	msg := fmt.Sprintf("%dx%d grid | %d events | %d players", gridSize, gridSize, len(events), len(playerIDs))
	respondEmbed(s, i, titleText, msg, colorSuccess, false)
}

// fetchAndParseCSV fetches a CSV file from URL and parses event descriptions
func fetchAndParseCSV(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch CSV: %s", resp.Status)
	}
	defer resp.Body.Close()

	reader := csv.NewReader(resp.Body)
	reader.TrimLeadingSpace = true

	var events []string
	firstRow := true

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(record) == 0 || strings.TrimSpace(record[0]) == "" {
			continue
		}

		// Skip header if first cell is "description"
		if firstRow && strings.ToLower(strings.TrimSpace(record[0])) == "description" {
			firstRow = false
			continue
		}
		firstRow = false

		events = append(events, strings.TrimSpace(record[0]))
	}

	return events, nil
}

type assignment struct {
	eventID int64
	row     int
	col     int
}

// distributeEvents implements the balanced distribution algorithm
func distributeEvents(eventIDs []int64, playerIDs []int64, gridSize int) map[int64][]assignment {
	cellsPerBoard := gridSize * gridSize
	totalCells := cellsPerBoard * len(playerIDs)

	// If more events than total cells, shuffle to avoid front-of-file bias
	if len(eventIDs) > totalCells {
		shuffled := make([]int64, len(eventIDs))
		copy(shuffled, eventIDs)
		rand.Shuffle(len(shuffled), func(i, j int) {
			shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
		})
		eventIDs = shuffled
	}

	// Initialize boards
	boards := make(map[int64][]assignment)
	for _, pid := range playerIDs {
		boards[pid] = make([]assignment, 0, cellsPerBoard)
	}

	// Track which events each board has
	boardHasEvent := make(map[int64]map[int64]bool)
	for _, pid := range playerIDs {
		boardHasEvent[pid] = make(map[int64]bool)
	}

	// Round-robin pass: ensure each event used at least once
	playerIdx := 0
	for _, eventID := range eventIDs {
		placed := false
		attempts := 0
		for attempts < len(playerIDs) {
			pid := playerIDs[playerIdx]
			if len(boards[pid]) < cellsPerBoard && !boardHasEvent[pid][eventID] {
				row := len(boards[pid]) / gridSize
				col := len(boards[pid]) % gridSize
				boards[pid] = append(boards[pid], assignment{eventID, row, col})
				boardHasEvent[pid][eventID] = true
				placed = true
				playerIdx = (playerIdx + 1) % len(playerIDs)
				break
			}
			playerIdx = (playerIdx + 1) % len(playerIDs)
			attempts++
		}
		if !placed {
			break // All boards full
		}
	}

	// Fill remaining cells with random events (no duplicates per board)
	for _, pid := range playerIDs {
		for len(boards[pid]) < cellsPerBoard {
			// Pick random event not yet on this board
			eventID := eventIDs[rand.Intn(len(eventIDs))]
			if !boardHasEvent[pid][eventID] {
				row := len(boards[pid]) / gridSize
				col := len(boards[pid]) % gridSize
				boards[pid] = append(boards[pid], assignment{eventID, row, col})
				boardHasEvent[pid][eventID] = true
			}
		}
	}

	return boards
}
