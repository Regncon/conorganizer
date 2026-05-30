package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckWritableDirectorySucceedsAndDoesNotLeaveTempFile(t *testing.T) {
	dir := t.TempDir()

	if err := CheckWritableDirectory(dir); err != nil {
		t.Fatalf("CheckWritableDirectory() error = %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read temp dir: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".conorganizer-write-check-") {
			t.Fatalf("write check left temp file behind: %s", entry.Name())
		}
	}
}

func TestCheckWritableDirectoryFailsForMissingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "missing")

	err := CheckWritableDirectory(dir)
	if err == nil {
		t.Fatalf("CheckWritableDirectory() error = nil, want missing directory error")
	}
	if !strings.Contains(err.Error(), "does not exist") {
		t.Fatalf("CheckWritableDirectory() error = %v, want missing directory error", err)
	}
}

func TestCheckWritableDirectoryFailsForFilePath(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(filePath, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	err := CheckWritableDirectory(filePath)
	if err == nil {
		t.Fatalf("CheckWritableDirectory() error = nil, want not directory error")
	}
	if !strings.Contains(err.Error(), "is not a directory") {
		t.Fatalf("CheckWritableDirectory() error = %v, want not directory error", err)
	}
}
