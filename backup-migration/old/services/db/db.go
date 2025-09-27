package db

import (
	"backup-migration/models"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

type Client struct {
	db   *sql.DB
	path string
}

func NewDBClient(state *models.AppState, logger *slog.Logger) (*Client, error) {
	client := &Client{}
	dbPath, _ := state.DB.Path.Get()

	fmt.Println("opening db")

	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		fmt.Printf("directory path does not exist: %s", dir)
		return nil, fmt.Errorf("directory path does not exist: %s", dir)
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		db, err := sql.Open("sqlite", dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open DB: %w", err)
		}
		if err = db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping DB: %w", err)
		}
		client.db = db
		client.path = dir
	}

	return client, nil
}

func (c *Client) Open() {
	fmt.Println("database opened")
}

func (c *Client) Close() {
	c.db.Close()
	fmt.Println("database closed")
}

func (c *Client) Validated() {
	fmt.Println("database closed")
}
