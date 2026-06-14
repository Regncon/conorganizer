package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestCheckWritableDirectorySucceedsAndDoesNotLeaveTempFile(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a writable image directory.",
		When:  "When the startup check probes it.",
		Then:  "Then the check succeeds and removes its temporary file.",
	})

	// Given
	expectedTempFilePrefix := ".conorganizer-write-check-"

	dir := t.TempDir()

	// When
	if err := CheckWritableDirectory(dir); err != nil {
		t.Fatalf("expected writable directory check to succeed: %v", err)
	}

	// Then
	assertNoDirectoryEntryWithPrefix(t, dir, expectedTempFilePrefix)
}

func TestCheckWritableDirectoryFailsForMissingDirectory(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the image directory does not exist.",
		When:  "When the startup check probes it.",
		Then:  "Then the check fails with a missing directory error.",
	})

	// Given
	expectedErrorText := "does not exist"

	dir := filepath.Join(t.TempDir(), "missing")

	// When
	err := CheckWritableDirectory(dir)

	// Then
	assertErrorContains(t, err, expectedErrorText)
}

func TestCheckWritableDirectoryFailsForFilePath(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the configured image path points to a file.",
		When:  "When the startup check probes it.",
		Then:  "Then the check fails before attempting a write probe.",
	})

	// Given
	expectedErrorText := "is not a directory"

	filePath := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(filePath, []byte("not a directory"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	// When
	err := CheckWritableDirectory(filePath)

	// Then
	assertErrorContains(t, err, expectedErrorText)
}

func assertNoDirectoryEntryWithPrefix(t *testing.T, dir string, prefix string) {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read temp dir: %v", err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), prefix) {
			t.Fatalf("write check left temp file behind: %s", entry.Name())
		}
	}
}
