package db

import (
	"context"
	"database/sql"
	_ "embed"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/init_schema.sql
var initSchemaSQL string

// DB wraps the database connection and provides data access methods
type DB struct {
	conn *sql.DB
}

// EventStatus represents the status of an event
type EventStatus string

const (
	EventStatusOpen   EventStatus = "OPEN"
	EventStatusClosed EventStatus = "CLOSED"
)

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
	path := os.Getenv("DB_PATH")
	if path == "" {
		path = "./bingo.db"
	}
	conn, err := sql.Open("sqlite3", path)
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

	// Ensure schema exists
	if err := ensureSchema(conn); err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

// ensureSchema creates tables if they don't exist
func ensureSchema(conn *sql.DB) error {
	// Check if tables already exist by looking for the games table
	var count int
	err := conn.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='games'").Scan(&count)
	if err != nil {
		return err
	}

	// If tables don't exist, create them
	if count == 0 {
		_, err := conn.Exec(initSchemaSQL)
		return err
	}

	return nil
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
