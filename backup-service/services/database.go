package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB(databaseFileName string, sqlFileName string) (*sql.DB, error) {
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

		if err = initializeDatabase(db, sqlFileName); err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
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

func initializeDatabase(db *sql.DB, filename string) error {
	sqlContent, err := loadSQLFile(filename)
	if err != nil {
		return fmt.Errorf("error loading SQL file: %w", err)
	}

	_, err = db.Exec(sqlContent)
	if err != nil {
		return fmt.Errorf("failed to execute SQL commands: %w", err)
	}

	return nil
}

func loadSQLFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	return string(data), nil
}
