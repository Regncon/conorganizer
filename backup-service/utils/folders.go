package utils

import (
	"fmt"
	"os"
)

// CreateRequiredFolders ensures all necessary folders exist with proper permissions.
func CreateRequiredFolders() error {
	requiredDirs := []string{
		"/data/regncon/tmp",
		"/data/regncon/logs",
		"/data/regncon/backup/hourly",
		"/data/regncon/backup/daily",
		"/data/regncon/backup/weekly",
		"/data/regncon/backup/yearly",
		"/data/regncon/backup/manually",
	}

	for _, dir := range requiredDirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}
