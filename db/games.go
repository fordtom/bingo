package db

import (
	"context"
	"database/sql"
)

// CreateGame creates a new game and returns its ID
func (db *DB) CreateGame(ctx context.Context, title string, gridSize int) (int64, error) {
	result, err := db.conn.ExecContext(ctx,
		"INSERT INTO games (title, grid_size) VALUES (?, ?)",
		title, gridSize,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetGame retrieves a specific game by ID
func (db *DB) GetGame(ctx context.Context, gameID int64) (*Game, error) {
	var game Game
	err := db.conn.QueryRowContext(ctx,
		"SELECT game_id, title, is_active, grid_size FROM games WHERE game_id = ?",
		gameID,
	).Scan(&game.ID, &game.Title, &game.IsActive, &game.GridSize)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &game, nil
}

// GetActiveGame retrieves the currently active game (returns nil if none)
func (db *DB) GetActiveGame(ctx context.Context) (*Game, error) {
	var game Game
	err := db.conn.QueryRowContext(ctx,
		"SELECT game_id, title, is_active, grid_size FROM games WHERE is_active = 1",
	).Scan(&game.ID, &game.Title, &game.IsActive, &game.GridSize)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &game, nil
}

// ListGames retrieves all games
func (db *DB) ListGames(ctx context.Context) ([]Game, error) {
	rows, err := db.conn.QueryContext(ctx,
		"SELECT game_id, title, is_active, grid_size FROM games ORDER BY game_id DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var games []Game
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Title, &game.IsActive, &game.GridSize); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, rows.Err()
}

// SetActiveGame sets a game as active and unsets all others
func (db *DB) SetActiveGame(ctx context.Context, gameID int64) error {
	return db.WithTx(ctx, func(tx *sql.Tx) error {
		// Unset all active games
		if _, err := tx.ExecContext(ctx, "UPDATE games SET is_active = 0 WHERE is_active = 1"); err != nil {
			return err
		}
		// Set the target game as active
		if _, err := tx.ExecContext(ctx, "UPDATE games SET is_active = 1 WHERE game_id = ?", gameID); err != nil {
			return err
		}
		return nil
	})
}

// DeleteGame removes a game
func (db *DB) DeleteGame(ctx context.Context, gameID int64) error {
	_, err := db.conn.ExecContext(ctx, "DELETE FROM games WHERE game_id = ?", gameID)
	return err
}

// DeleteGameCascade removes a game and all associated data in proper order
func (db *DB) DeleteGameCascade(ctx context.Context, gameID int64) error {
	return db.WithTx(ctx, func(tx *sql.Tx) error {
		// Delete boards (cascades to board_squares via FK)
		if _, err := tx.ExecContext(ctx, "DELETE FROM boards WHERE game_id = ?", gameID); err != nil {
			return err
		}
		// Delete votes for events in this game
		if _, err := tx.ExecContext(ctx, "DELETE FROM votes WHERE event_id IN (SELECT event_id FROM events WHERE game_id = ?)", gameID); err != nil {
			return err
		}
		// Delete events
		if _, err := tx.ExecContext(ctx, "DELETE FROM events WHERE game_id = ?", gameID); err != nil {
			return err
		}
		// Delete game
		if _, err := tx.ExecContext(ctx, "DELETE FROM games WHERE game_id = ?", gameID); err != nil {
			return err
		}
		return nil
	})
}

// GetPlayerCountForGame returns the number of players (boards) for a game
func (db *DB) GetPlayerCountForGame(ctx context.Context, gameID int64) (int, error) {
	var count int
	err := db.conn.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM boards WHERE game_id = ?",
		gameID,
	).Scan(&count)
	return count, err
}

// GetGamePlayerIDs returns all player user IDs for a game
func (db *DB) GetGamePlayerIDs(ctx context.Context, gameID int64) ([]int64, error) {
	rows, err := db.conn.QueryContext(ctx,
		"SELECT user_id FROM boards WHERE game_id = ? ORDER BY user_id",
		gameID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playerIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		playerIDs = append(playerIDs, userID)
	}
	return playerIDs, rows.Err()
}

// GetEventCounts returns open and closed event counts for a game
func (db *DB) GetEventCounts(ctx context.Context, gameID int64) (open, closed int, err error) {
	rows, err := db.conn.QueryContext(ctx,
		"SELECT status, COUNT(*) FROM events WHERE game_id = ? GROUP BY status",
		gameID,
	)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return 0, 0, err
		}
		if status == "OPEN" {
			open = count
		} else if status == "CLOSED" {
			closed = count
		}
	}
	return open, closed, rows.Err()
}
