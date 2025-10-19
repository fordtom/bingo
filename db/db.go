package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// DB wraps the database connection and provides data access methods
type DB struct {
	conn *sql.DB
}

// Domain types
type Game struct {
	ID       int64
	Title    string
	IsActive bool
	GridSize int
}

type Event struct {
	ID          int64
	GameID      int64
	DisplayID   int
	Description string
	Status      string
}

type Board struct {
	ID       int64
	GameID   int64
	UserID   int64
	GridSize int
}

type BoardSquare struct {
	BoardID int64
	Row     int
	Column  int
	EventID int64
}

type BoardSquareWithEvent struct {
	BoardSquare
	EventDescription string
	EventStatus      string
}

type Vote struct {
	EventID int64
	UserID  int64
	VotedAt time.Time
}

// InitDB creates and configures the database connection
func InitDB() (*DB, error) {
	conn, err := sql.Open("sqlite3", "./bingo.db")
	if err != nil {
		return nil, err
	}

	// SQLite-specific optimizations for concurrency
	conn.SetMaxOpenConns(1)
	if _, err := conn.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		return nil, err
	}
	if _, err := conn.Exec("PRAGMA journal_mode = WAL"); err != nil {
		return nil, err
	}
	if _, err := conn.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// BeginTx starts a transaction
func (db *DB) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return db.conn.BeginTx(ctx, nil)
}

// WithTx runs a function in a transaction with automatic rollback/commit
func (db *DB) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
