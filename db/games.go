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
