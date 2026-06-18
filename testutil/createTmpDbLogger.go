package testutil

import (
	"database/sql"
	"log/slog"
	"path/filepath"
	"testing"

	"github.com/Regncon/conorganizer/service"
)

func CreateTestDB(t testing.TB, name string) *sql.DB {
	t.Helper()

	db, err := initTestDB(t, name)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close test database: %v", err)
		}
	})

	return db
}

func CreateTestDBAndLogger(t testing.TB, name string) (*sql.DB, *slog.Logger) {
	t.Helper()

	return CreateTestDB(t, name), NewTestLogger()
}

func initTestDB(t testing.TB, name string) (*sql.DB, error) {
	t.Helper()

	return service.InitTestDBFrom(filepath.Join(t.TempDir(), name+".db"))
}

func NewTestLogger() *slog.Logger {
	return NewSlogAdapter(&StubLogger{})
}

func MustExec(t testing.TB, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("failed to execute query: %v\nquery: %s", err, query)
	}
}

func QueryInt(t testing.TB, db *sql.DB, query string, args ...any) int {
	t.Helper()

	var value int
	if err := db.QueryRow(query, args...).Scan(&value); err != nil {
		t.Fatalf("failed to query integer value: %v\nquery: %s", err, query)
	}
	return value
}
