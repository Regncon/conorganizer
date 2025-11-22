package service

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "modernc.org/sqlite"
)

func InitTestDBFrom(testDBPath string) (*sql.DB, error) {
	if testDBPath == "" {
		return nil, fmt.Errorf("testDBPath is required")
	}

	rootDir, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	schemaPath := filepath.Join(rootDir, "schema.sql")

	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("read schema file %q: %w", schemaPath, err)
	}
	if len(schemaBytes) == 0 {
		return nil, fmt.Errorf("schema file %q is empty", schemaPath)
	}

	db, err := sql.Open("sqlite", testDBPath)
	if err != nil {
		return nil, fmt.Errorf("open destination DB: %w", err)
	}

	tx, err := db.Begin()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("begin transaction: %w", err)
	}

	if _, err := tx.Exec(`PRAGMA foreign_keys=OFF;`); err != nil {
		tx.Rollback()
		db.Close()
		return nil, fmt.Errorf("disable foreign_keys: %w", err)
	}

	if _, err := tx.Exec(string(schemaBytes)); err != nil {
		tx.Rollback()
		db.Close()
		return nil, fmt.Errorf("execute schema.sql: %w", err)
	}

	if _, err := tx.Exec(`PRAGMA foreign_keys=ON;`); err != nil {
		tx.Rollback()
		db.Close()
		return nil, fmt.Errorf("enable foreign_keys: %w", err)
	}

	if err := tx.Commit(); err != nil {
		db.Close()
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return db, nil
}

func findProjectRoot() (string, error) {
	currentDir, err := getThisFileDir()
	if err != nil {
		return "", err
	}

	for {
		goMod := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goMod); err == nil {
			return currentDir, nil
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("could not find go.mod when searching from %s", currentDir)
		}

		currentDir = parent
	}
}

func getThisFileDir() (string, error) {
	_, thisFilePath, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("could not determine caller file path")
	}

	return filepath.Dir(thisFilePath), nil
}
