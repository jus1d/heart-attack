package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	db *sql.DB
}

func New(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &DB{db: db}, nil
}

func (d *DB) Migrate(ctx context.Context) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			telegram_id INTEGER UNIQUE NOT NULL,
			username TEXT,
			first_name TEXT,
			last_name TEXT,
			is_bot INTEGER NOT NULL DEFAULT 0,
			language_code TEXT,
			is_premium INTEGER NOT NULL DEFAULT 0,
			sex TEXT,
			about TEXT NOT NULL DEFAULT '',
			state TEXT NOT NULL DEFAULT 'start',
			time_ranges TEXT NOT NULL DEFAULT '000000',
			is_admin INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS pairs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			dill_id INTEGER NOT NULL REFERENCES users(id),
			doe_id INTEGER NOT NULL REFERENCES users(id),
			score REAL NOT NULL,
			time_intersection TEXT NOT NULL,
			is_fullmatch INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS places (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			description TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS meetings (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			pair_id INTEGER NOT NULL REFERENCES pairs(id),
			place_id INTEGER NOT NULL REFERENCES places(id),
			time TEXT NOT NULL,
			dill_confirmed INTEGER NOT NULL DEFAULT 0,
			doe_confirmed INTEGER NOT NULL DEFAULT 0,
			dill_cancelled INTEGER NOT NULL DEFAULT 0,
			doe_cancelled INTEGER NOT NULL DEFAULT 0
		)`,
	}

	for _, q := range queries {
		if _, err := d.db.ExecContext(ctx, q); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}

	return nil
}

func (d *DB) Close() error {
	return d.db.Close()
}
