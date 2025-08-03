package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/pierrec/lz4/v4"
)

// RotateFile is responsible for deleting old files
func RotateFile(directory, interval string, keep int) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	var targets []os.DirEntry
	for _, f := range files {
		if strings.HasPrefix(f.Name(), interval) {
			targets = append(targets, f)
		}
	}

	if len(targets) <= keep {
		return nil
	}

	sort.Slice(targets, func(i, j int) bool {
		fi, _ := targets[i].Info()
		fj, _ := targets[j].Info()
		return fi.ModTime().Before(fj.ModTime())
	})

	for i := 0; i < len(targets)-keep; i++ {
		os.Remove(filepath.Join(directory, targets[i].Name()))
	}

	return nil
}

// DecompressLZ4 takes a path to a .lz4-compressed snapshot file,
// decompresses it, and returns the path to a temporary .db file.
func DecompressLZ4(lz4Path string) (string, error) {
	// Open compressed input file
	inFile, err := os.Open(lz4Path)
	if err != nil {
		return "", fmt.Errorf("failed to open LZ4 file: %w", err)
	}
	defer inFile.Close()

	// Create temp output file for decompressed .db
	outFile, err := os.CreateTemp("", "*.db")
	if err != nil {
		return "", fmt.Errorf("failed to create temp DB file: %w", err)
	}
	defer outFile.Close()

	// Create LZ4 reader and decompress to file
	lz4Reader := lz4.NewReader(inFile)
	if _, err := io.Copy(outFile, lz4Reader); err != nil {
		return "", fmt.Errorf("failed to decompress LZ4 file: %w", err)
	}

	return outFile.Name(), nil
}

// RotateBackups moves the validated .db file to the backup directory,
// renaming it using the given timestamp. It returns the final path.
func RotateBackups(dbPath string, backupPath string, retention int) (string, error) {
	// Step 1: Ensure directory exists
	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Step 2: Scan for .db files
	entries, err := os.ReadDir(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to read backup directory: %w", err)
	}

	type fileEntry struct {
		name string
		time time.Time
	}
	var dbFiles []fileEntry

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".db") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		dbFiles = append(dbFiles, fileEntry{
			name: entry.Name(),
			time: info.ModTime(),
		})
	}

	// Step 3: Delete oldest if limit reached
	if len(dbFiles) >= retention {
		sort.Slice(dbFiles, func(i, j int) bool {
			return dbFiles[i].time.Before(dbFiles[j].time)
		})

		oldest := filepath.Join(backupPath, dbFiles[0].name)
		if err := os.Remove(oldest); err != nil {
			return "", fmt.Errorf("failed to delete oldest backup file: %w", err)
		}
	}

	// Step 4: Move & rename validated db file
	timestamp := time.Now().In(time.Local).Format("2006-01-02T15-04")
	destName := fmt.Sprintf("%s.db", timestamp)
	destPath := filepath.Join(backupPath, destName)

	if err := os.Rename(dbPath, destPath); err != nil {
		return "", fmt.Errorf("failed to move backup file: %w", err)
	}

	return destPath, nil
}
