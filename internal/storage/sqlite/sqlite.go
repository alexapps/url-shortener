package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexapps/url-shortener/internal/storage"
	"github.com/mattn/go-sqlite3"
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

	INSERT_INTO_Q = `INSERT INTO url(url, alias) VALUES(?, ?)`
	GET_URL_Q     = `SELECT url FROM url WHERE alias = ?`
	DELETE_URL_Q  = `DELETE FROM url WHERE alias = ?`
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// used for errors or log messages as function path name
	const op = "storage.sqlie.New"
	//TODO: implement db migration
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s / %w", op, err)
	}

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

// SaveURL save URL. Returns an id of the last inserted row.
// Be cearful, the id could not be supported in some db
func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const op = "storage.sqlie.SaveURL"

	stmt, err := s.db.Prepare(INSERT_INTO_Q)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.sqlie.GetURL"

	stmt, err := s.db.Prepare(GET_URL_Q)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}
		return "", fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.sqlie.DeleteURL"

	stmt, err := s.db.Prepare(DELETE_URL_Q)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return err
}
