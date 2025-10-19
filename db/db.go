package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./bingo.db") // "sqlite3" string refers to the registered driver
	if err != nil {
		return nil, err
	}
	return db, nil
}
