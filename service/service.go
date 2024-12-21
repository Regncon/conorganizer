package service

import (
	"database/sql"
	"fmt"
	"github.com/Regncon/conorganizer/models"
	_ "github.com/mattn/go-sqlite3"
)

// EventService defines the interface for event operations.
type EventService interface {
	AddEvent(name, description string) (int64, error)
	GetEvents() ([]models.Event, error)
	InitDB(dbPath string) (*sql.DB, error)
}

// DatabaseService is a concrete implementation of EventService.
type DatabaseService struct {
	db *sql.DB
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

	return db, nil
}

// NewDatabaseService creates a new DatabaseService instance.
func NewDatabaseService(db *sql.DB) *DatabaseService {
	return &DatabaseService{db: db}
}

// AddEvent adds a new event to the database.
func (s *DatabaseService) AddEvent(name, description string) (int64, error) {
	query := "INSERT INTO events (name, description) VALUES (?, ?)"
	result, err := s.db.Exec(query, name, description)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetEvents retrieves all events from the database as a list of models.Event.
func (s *DatabaseService) GetEvents() ([]models.Event, error) {
	query := "SELECT id, name, description, start_time, end_time FROM events"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(&event.ID, &event.Name, &event.Description); err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}
