package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func CheckWritableDirectory(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("directory path is empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %q does not exist: %w", path, err)
		}
		return fmt.Errorf("could not access directory %q: %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path %q is not a directory", path)
	}

	tempFile, err := os.CreateTemp(path, ".conorganizer-write-check-*")
	if err != nil {
		return fmt.Errorf("directory %q is not writable: %w", path, err)
	}

	tempFileName := tempFile.Name()
	closeErr := tempFile.Close()
	removeErr := os.Remove(tempFileName)
	if err := errors.Join(closeErr, removeErr); err != nil {
		return fmt.Errorf("clean up directory write check file %q: %w", tempFileName, err)
	}

	return nil
}
