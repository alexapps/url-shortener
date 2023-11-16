package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CREATE_TABLE_Q = `
	CREATE TABLE IF NOT EXISTS URL(
        id INTEGER PRIMARY KEY,
        alias TEXT NOT NULL UNIQUE,
        url TEXT NOT NULL);
     CREATE INDEX IF NOT EXISTS idx_alias ON url(alias)
	`
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// used for errors or log messages as function path name
	const op = "storage.sqlie.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s / %w", op, err)
	}

	//TODO: implement migration
	stmt, err := db.Prepare(CREATE_TABLE_Q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}
