package testutil

import (
	"database/sql"
	"log/slog"
	"testing"

	"github.com/Regncon/conorganizer/service"
	"github.com/google/uuid"
)

// CreateTemporaryDBAndLogger creates a temporary named database and slogger that you can use during tests
func CreateTemporaryDBAndLogger(name string, t *testing.T) (*sql.DB, *slog.Logger, error) {
	uniqueDatabaseName := name + "_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/tests/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom("../../database/events.db", testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	sl := &StubLogger{}
	slogger := NewSlogAdapter(sl)

	return db, slogger, nil
}
