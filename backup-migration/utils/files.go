package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(data []byte, path string, name string) (*string, error) {
	if name == "" {
		return nil, fmt.Errorf("empty file name")
	}

	// Check if path is valid, try to create if not
	fileInfo, err := os.Stat(path)
	if err != nil {
		// Try to create dir
		if os.IsNotExist(err) {
			err := os.MkdirAll(path, 0755)
			if err != nil {
				return nil, fmt.Errorf("unable to create dir db backup: %w", err)
			}
		} else {
			return nil, fmt.Errorf("not sure, ill fix later: %w", err)
		}
	} else if !fileInfo.IsDir() {
		fmt.Println(path)
		return nil, fmt.Errorf("path is not a dir %w", err)
	}

	// Create a tmp file while transfering bytes
	tmp, err := os.CreateTemp(path, name)
	if err != nil {
		return nil, fmt.Errorf("create temp file: %w", err)
	}

	// Cleanup tmp file if stuff breaks
	defer func() {
		_ = tmp.Close()
		_ = os.Remove(tmp.Name())
	}()

	// Write bytes to file
	if _, err := tmp.Write(data); err != nil {
		return nil, fmt.Errorf("error writing to tmmp file: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		return nil, fmt.Errorf("error syncing tmp file: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return nil, fmt.Errorf("error closing tmp file: %w", err)
	}

	// Rename tmp file to its propper name and extension
	fullPath := filepath.Join(path, name)
	if err := os.Rename(tmp.Name(), fullPath); err != nil {
		return nil, fmt.Errorf("failed to rename tmp file: %w", err)
	}
	return &fullPath, nil
}
