package db

import (
	"context"
	"database/sql"
)

// CreateBoard creates a board for a user in a game
func (db *DB) CreateBoard(ctx context.Context, tx *sql.Tx, gameID, userID int64, gridSize int) (int64, error) {
	result, err := tx.ExecContext(ctx,
		"INSERT INTO boards (game_id, user_id, grid_size) VALUES (?, ?, ?)",
		gameID, userID, gridSize,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetUserBoard retrieves a user's board with all squares populated with event details
func (db *DB) GetUserBoard(ctx context.Context, gameID, userID int64) (*Board, []BoardSquareWithEvent, error) {
	// First get the board
	var board Board
	err := db.conn.QueryRowContext(ctx,
		"SELECT board_id, game_id, user_id, grid_size FROM boards WHERE game_id = ? AND user_id = ?",
		gameID, userID,
	).Scan(&board.ID, &board.GameID, &board.UserID, &board.GridSize)
	if err == sql.ErrNoRows {
		return nil, nil, nil
	}
	if err != nil {
		return nil, nil, err
	}

	// Then get all squares with event details
	rows, err := db.conn.QueryContext(ctx,
		`SELECT bs.board_id, bs.row, bs.column, bs.event_id, e.description, e.status
		 FROM board_squares bs
		 JOIN events e ON bs.event_id = e.event_id
		 WHERE bs.board_id = ?
		 ORDER BY bs.row, bs.column`,
		board.ID,
	)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var squares []BoardSquareWithEvent
	for rows.Next() {
		var square BoardSquareWithEvent
		if err := rows.Scan(
			&square.BoardID, &square.Row, &square.Column, &square.EventID,
			&square.EventDescription, &square.EventStatus,
		); err != nil {
			return nil, nil, err
		}
		squares = append(squares, square)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	return &board, squares, nil
}

// CreateBoardSquare creates a single board square
func (db *DB) CreateBoardSquare(ctx context.Context, tx *sql.Tx, boardID int64, row, col int, eventID int64) error {
	_, err := tx.ExecContext(ctx,
		"INSERT INTO board_squares (board_id, row, column, event_id) VALUES (?, ?, ?, ?)",
		boardID, row, col, eventID,
	)
	return err
}

// CreateBoardSquares bulk creates squares for a board
func (db *DB) CreateBoardSquares(ctx context.Context, tx *sql.Tx, boardID int64, squares []BoardSquare) error {
	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO board_squares (board_id, row, column, event_id) VALUES (?, ?, ?, ?)",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, square := range squares {
		if _, err := stmt.ExecContext(ctx, boardID, square.Row, square.Column, square.EventID); err != nil {
			return err
		}
	}
	return nil
}
