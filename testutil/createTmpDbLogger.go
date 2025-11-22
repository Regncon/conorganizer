package testutil

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Regncon/conorganizer/service"
	"github.com/google/uuid"
)

func CreateTemporaryDBAndLogger(name string, t *testing.T) (*sql.DB, *slog.Logger, error) {
	t.Helper()

	projectRoot := getProjectRoot(t)

	uniqueDatabaseName := name + "_" + t.Name() + "_" + uuid.New().String() + ".db"

	databaseTestsDirPath := filepath.Join(projectRoot, "database", "tests")
	if err := os.MkdirAll(databaseTestsDirPath, 0o755); err != nil {
		t.Fatalf("failed to create database tests directory: %v", err)
	}

	testDBPath := filepath.Join(databaseTestsDirPath, uniqueDatabaseName)
	seedDBPath := filepath.Join(projectRoot, "database", "events.db")

	db, err := service.InitTestDBFrom(seedDBPath, testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	sl := &StubLogger{}
	slogger := NewSlogAdapter(sl)

	return db, slogger, nil
}

func getProjectRoot(t *testing.T) string {
	t.Helper()

	_, thisFilePath, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("unable to determine caller information")
	}

	currentDir := filepath.Dir(thisFilePath)

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		_, err := os.Stat(goModPath)
		if err == nil {
			return currentDir
		}

		parentDir := filepath.Dir(currentDir)
		if parentDir == currentDir {
			t.Fatalf("could not find go.mod when walking up from %s", thisFilePath)
		}

		currentDir = parentDir
	}
}
