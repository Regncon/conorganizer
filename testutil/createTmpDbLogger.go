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

	uniqueName := name + "_" + t.Name() + "_" + uuid.New().String() + ".db"

	databaseTestsDir := filepath.Join(projectRoot, "database", "tests")
	if err := os.MkdirAll(databaseTestsDir, 0o755); err != nil {
		t.Fatalf("failed to create test db directory: %v", err)
	}

	testDBPath := filepath.Join(databaseTestsDir, uniqueName)
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
		t.Fatalf("unable to get caller info")
	}

	testutilDir := filepath.Dir(thisFilePath)
	projectRoot := filepath.Dir(testutilDir)

	return projectRoot
}
