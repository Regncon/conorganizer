package utils

import (
	"database/sql"
	"fmt"

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
