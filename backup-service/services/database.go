package services

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	dbPath := "/data/regncon/logs/logs.db"
	sqlInitPath := "/usr/local/share/regncon/initialize.sql"

	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
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

		if err = initializeDatabase(db, sqlInitPath); err != nil {
			return nil, fmt.Errorf("failed to initialize database: %w", err)
		}

		return db, nil
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	// Run migration
	if err := MigrateBackupLogsTable(db); err != nil {
		return nil, fmt.Errorf("failed to apply DB migration: %w", err)
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

func MigrateBackupLogsTable(db *sql.DB) error {
	rows, err := db.Query(`PRAGMA table_info(backup_logs)`)
	if err != nil {
		return fmt.Errorf("failed to inspect table: %w", err)
	}
	defer rows.Close()

	columns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue sql.NullString

		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return fmt.Errorf("failed to scan column info: %w", err)
		}
		columns[name] = true
	}

	// Add missing columns if not found
	if !columns["stage"] {
		_, err := db.Exec(`ALTER TABLE backup_logs ADD COLUMN stage TEXT NOT NULL DEFAULT 'pending' CHECK (stage IN ('pending', 'downloading', 'decompressing', 'validating', 'moving', 'completed'))`)
		if err != nil {
			return fmt.Errorf("failed to add column 'stage': %w", err)
		}
	}

	if !columns["file_path"] {
		_, err := db.Exec(`ALTER TABLE backup_logs ADD COLUMN file_path TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add column 'file_path': %w", err)
		}
	}

	if !columns["file_size"] {
		_, err := db.Exec(`ALTER TABLE backup_logs ADD COLUMN file_size INTEGER NOT NULL DEFAULT 0`)
		if err != nil {
			return fmt.Errorf("failed to add column 'file_size': %w", err)
		}
	}

	if !columns["db_prefix"] {
		_, err := db.Exec(`ALTER TABLE backup_logs ADD COLUMN db_prefix TEXT NOT NULL DEFAULT ''`)
		if err != nil {
			return fmt.Errorf("failed to add column 'db_prefix': %w", err)
		}
	}

	if !columns["events"] {
		_, err := db.Exec(`ALTER TABLE backup_logs ADD COLUMN events INTEGER NOT NULL DEFAULT 0`)
		if err != nil {
			return fmt.Errorf("failed to add column 'events': %w", err)
		}
	}

	return nil
}
