package utils

import (
	"database/sql"
	"fmt"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func ValidateSnapshot(dbPath string) error {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("failed to open SQLite file: %w", err)
	}
	defer db.Close()

	var result string
	err = db.QueryRow("PRAGMA integrity_check;").Scan(&result)
	if err != nil {
		return fmt.Errorf("integrity check query failed: %w", err)
	}

	if result != "ok" {
		return fmt.Errorf("integrity check failed: %s", result)
	}

	return nil
}

func ValidateJournals(tmpDir string) (string, error) {
	var tmpPath = filepath.Join(tmpDir, "events.db")
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)", tmpPath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()

	var dbPath = filepath.Join(tmpDir, "events-validated.db")

	if _, err = db.Exec(`PRAGMA wal_checkpoint(FULL);`); err != nil {
		return "", fmt.Errorf("wal checkpoint failed: %w", err)
	}
	if _, err = db.Exec(`VACUUM INTO ?;`, dbPath); err != nil {
		return "", fmt.Errorf("vacuum into failed: %w", err)
	}
	return dbPath, nil
}

func ValidateEvents(dbPath string) (int64, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return 0, fmt.Errorf("failed to open SQLite file: %w", err)
	}
	defer db.Close()

	var result int64
	if err := db.QueryRow(`SELECT COUNT(*) FROM events;`).Scan(&result); err != nil {
		return 0, fmt.Errorf("count events: %w", err)
	}

	return result, nil
}
