package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// InitTestDBFrom copies the schema (no data) from templateDBPath into a new DB at testDBPath.
// It recreates tables, indexes, triggers and views, and copies PRAGMA user_version.
// testDBPath will be overwritten if it already exists.
func InitTestDBFrom(templateDBPath, testDBPath string) (*sql.DB, error) {
	if templateDBPath == "" {
		return nil, errors.New("templateDBPath is required")
	}
	if testDBPath == "" {
		return nil, errors.New("testDBPath is required")
	}
	if _, err := os.Stat(templateDBPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("template database not found: %s", templateDBPath)
	}

	// Ensure destination dir exists
	if dir := filepath.Dir(testDBPath); dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return nil, fmt.Errorf("destination directory does not exist: %s", dir)
		}
	}

	// Start clean
	if _, err := os.Stat(testDBPath); err == nil {
		_ = os.Remove(testDBPath)
	}

	src, err := sql.Open("sqlite", templateDBPath)
	if err != nil {
		return nil, fmt.Errorf("open template: %w", err)
	}
	defer src.Close()

	// Pull schema objects from the template DB
	var tables, views, indexes, triggers []string
	rows, err := src.Query(`
		SELECT type, sql
		FROM sqlite_master
		WHERE sql IS NOT NULL
		  AND name NOT LIKE 'sqlite_%'
		ORDER BY
		  CASE type
		    WHEN 'table'  THEN 1
		    WHEN 'view'   THEN 2
		    WHEN 'index'  THEN 3
		    WHEN 'trigger'THEN 4
		    ELSE 5
		  END, name`)
	if err != nil {
		return nil, fmt.Errorf("read schema from template: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var typ, stmt string
		if err := rows.Scan(&typ, &stmt); err != nil {
			return nil, fmt.Errorf("scan schema row: %w", err)
		}
		switch typ {
		case "table":
			tables = append(tables, stmt)
		case "view":
			views = append(views, stmt)
		case "index":
			indexes = append(indexes, stmt)
		case "trigger":
			triggers = append(triggers, stmt)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate schema: %w", err)
	}

	// Copy PRAGMA user_version (optional, handy)
	var userVersion int
	_ = src.QueryRow(`PRAGMA user_version;`).Scan(&userVersion)

	dst, err := sql.Open("sqlite", testDBPath)
	if err != nil {
		return nil, fmt.Errorf("open destination: %w", err)
	}

	tx, err := dst.Begin()
	if err != nil {
		dst.Close()
		return nil, fmt.Errorf("begin tx: %w", err)
	}
	// Safer while creating objects in arbitrary order
	if _, err := tx.Exec(`PRAGMA foreign_keys=OFF;`); err != nil {
		tx.Rollback()
		dst.Close()
		return nil, fmt.Errorf("disable foreign_keys: %w", err)
	}

	execMany := func(stmts []string) error {
		for _, s := range stmts {
			if _, err := tx.Exec(s); err != nil {
				return fmt.Errorf("exec schema statement failed: %w\nSQL: %s", err, s)
			}
		}
		return nil
	}

	if err := execMany(tables); err != nil {
		tx.Rollback()
		dst.Close()
		return nil, err
	}
	if err := execMany(views); err != nil {
		tx.Rollback()
		dst.Close()
		return nil, err
	}
	if err := execMany(indexes); err != nil {
		tx.Rollback()
		dst.Close()
		return nil, err
	}
	if err := execMany(triggers); err != nil {
		tx.Rollback()
		dst.Close()
		return nil, err
	}

	if _, err := tx.Exec(fmt.Sprintf(`PRAGMA user_version=%d;`, userVersion)); err != nil {
		tx.Rollback()
		dst.Close()
		return nil, fmt.Errorf("set user_version: %w", err)
	}
	if _, err := tx.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		tx.Rollback()
		dst.Close()
		return nil, fmt.Errorf("enable foreign_keys: %w", err)
	}
	if err := tx.Commit(); err != nil {
		dst.Close()
		return nil, fmt.Errorf("commit schema tx: %w", err)
	}

	if err := dst.Ping(); err != nil {
		dst.Close()
		return nil, fmt.Errorf("ping new test DB: %w", err)
	}

	return dst, nil
}
