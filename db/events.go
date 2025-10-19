package db

import (
	"context"
	"database/sql"
)

// CreateEvent creates a single event for a game with a specific display_id
func (db *DB) CreateEvent(ctx context.Context, tx *sql.Tx, gameID int64, displayID int, description string) (int64, error) {
	result, err := tx.ExecContext(ctx,
		"INSERT INTO events (game_id, display_id, description, status) VALUES (?, ?, ?, 'OPEN')",
		gameID, displayID, description,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// CreateEvents bulk creates events for a game with sequential display_ids starting at 1
func (db *DB) CreateEvents(ctx context.Context, tx *sql.Tx, gameID int64, descriptions []string) error {
	stmt, err := tx.PrepareContext(ctx,
		"INSERT INTO events (game_id, display_id, description, status) VALUES (?, ?, ?, 'OPEN')",
	)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i, description := range descriptions {
		displayID := i + 1 // 1-indexed for users
		if _, err := stmt.ExecContext(ctx, gameID, displayID, description); err != nil {
			return err
		}
	}
	return nil
}

// GetEventByDisplayID retrieves an event by its user-facing display_id within a game
func (db *DB) GetEventByDisplayID(ctx context.Context, gameID int64, displayID int) (*Event, error) {
	var event Event
	err := db.conn.QueryRowContext(ctx,
		"SELECT event_id, game_id, display_id, description, status FROM events WHERE game_id = ? AND display_id = ?",
		gameID, displayID,
	).Scan(&event.ID, &event.GameID, &event.DisplayID, &event.Description, &event.Status)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &event, nil
}

// GetGameEvents retrieves all events for a game, ordered by display_id
func (db *DB) GetGameEvents(ctx context.Context, gameID int64) ([]Event, error) {
	rows, err := db.conn.QueryContext(ctx,
		"SELECT event_id, game_id, display_id, description, status FROM events WHERE game_id = ? ORDER BY display_id",
		gameID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.GameID, &event.DisplayID, &event.Description, &event.Status); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, rows.Err()
}

// UpdateEventStatus changes an event's status to 'OPEN' or 'CLOSED'
func (db *DB) UpdateEventStatus(ctx context.Context, eventID int64, status string) error {
	_, err := db.conn.ExecContext(ctx,
		"UPDATE events SET status = ? WHERE event_id = ?",
		status, eventID,
	)
	return err
}
