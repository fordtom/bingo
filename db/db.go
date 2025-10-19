package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./bingo.db")
	if err != nil {
		return nil, err
	}

	// SQLite-specific optimizations for concurrency
	db.SetMaxOpenConns(1)                 // SQLite benefits from single connection
	db.Exec("PRAGMA busy_timeout = 5000") // Wait up to 5s for locks
	db.Exec("PRAGMA journal_mode = WAL")  // Better concurrent read performance

	return db, nil
}
