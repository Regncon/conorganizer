package service

import (
	"database/sql"
	"fmt"
	"github.com/Regncon/conorganizer/models"
	_ "github.com/mattn/go-sqlite3"
)

type EventService interface {
	AddEvent(name, description string) (int64, error)
	GetEvents() ([]models.Event, error)
	InitDB(dbPath string) (*sql.DB, error)
}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	// Verify the connection is working
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	if err = createEventsTable(db); err != nil {
		return nil, fmt.Errorf("failed to create events table: %w", err)
	}

	return db, nil
}

func createEventsTable(db *sql.DB) error {
	tableCreationQuery := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL
	)`

	_, err := db.Exec(tableCreationQuery)
	if err != nil {
		return fmt.Errorf("failed to create events table: %w", err)
	}

	return nil
}
