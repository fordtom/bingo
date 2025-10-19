package db

import (
	"context"
)

// CreateVote records a vote for an event by a user
func (db *DB) CreateVote(ctx context.Context, eventID, userID int64) error {
	_, err := db.conn.ExecContext(ctx,
		"INSERT INTO votes (event_id, user_id) VALUES (?, ?)",
		eventID, userID,
	)
	return err
}

// GetVoteCount returns the number of votes for an event
func (db *DB) GetVoteCount(ctx context.Context, eventID int64) (int, error) {
	var count int
	err := db.conn.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM votes WHERE event_id = ?",
		eventID,
	).Scan(&count)
	return count, err
}

// HasUserVoted checks if a user has already voted on an event
func (db *DB) HasUserVoted(ctx context.Context, eventID, userID int64) (bool, error) {
	var count int
	err := db.conn.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM votes WHERE event_id = ? AND user_id = ?",
		eventID, userID,
	).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetEventVoters returns all user IDs who voted on an event
func (db *DB) GetEventVoters(ctx context.Context, eventID int64) ([]int64, error) {
	rows, err := db.conn.QueryContext(ctx,
		"SELECT user_id FROM votes WHERE event_id = ? ORDER BY voted_at",
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var voters []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		voters = append(voters, userID)
	}
	return voters, rows.Err()
}
