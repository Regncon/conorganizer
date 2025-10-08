package service

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(databaseFileName string) (*sql.DB, error) {
	dir := filepath.Dir(databaseFileName)
	if dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return nil, fmt.Errorf("directory path does not exist: %s", dir)
		}
	}

	if _, err := os.Stat(databaseFileName); os.IsNotExist(err) {
		db, err := sql.Open("sqlite", databaseFileName)
		if err != nil {
			return nil, fmt.Errorf("failed to open DB: %w", err)
		}
		if err = db.Ping(); err != nil {
			return nil, fmt.Errorf("failed to ping DB: %w", err)
		}

		return db, nil
	}

	db, err := sql.Open("sqlite", databaseFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	return db, nil
}
