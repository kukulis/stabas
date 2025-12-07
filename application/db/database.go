package db

import (
	"database/sql"
	"fmt"
	"sync"

	_ "modernc.org/sqlite"
)

type Database struct {
	db   *sql.DB
	path string
	mu   sync.Mutex
}

func NewDatabase(path string) *Database {
	return &Database{
		path: path,
	}
}

func (d *Database) GetDB() (*sql.DB, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.db != nil {
		return d.db, nil
	}

	db, err := sql.Open("sqlite", d.path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	d.db = db
	return d.db, nil
}

func (d *Database) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.db != nil {
		return d.db.Close()
	}
	return nil
}
