package service

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestInitDBAppliesProductionSQLiteSettings(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a database with the required production schema.",
		When:  "When the application opens it.",
		Then:  "Then SQLite is configured with the required production settings.",
	})

	// Given
	expectedForeignKeys := 1
	expectedJournalMode := "wal"
	expectedBusyTimeoutMillis := defaultSQLiteBusyTimeoutMillis
	expectedSynchronous := 1
	expectedMaxOpenConnections := 1

	dbPath := filepath.Join(t.TempDir(), "events.db")
	createMinimalProductionDB(t, dbPath)

	// When
	db, err := InitDB(dbPath)

	// Then
	if err != nil {
		t.Fatalf("expected database initialization to succeed: %v", err)
	}
	defer db.Close()

	assertSQLitePragmaInt(t, db, "foreign_keys", expectedForeignKeys)
	assertSQLitePragmaString(t, db, "journal_mode", expectedJournalMode)
	assertSQLitePragmaAtLeast(t, db, "busy_timeout", expectedBusyTimeoutMillis)
	assertSQLitePragmaInt(t, db, "synchronous", expectedSynchronous)
	assertForeignKeysRejectInvalidChild(t, db)
	assertMaxOpenConnections(t, db, expectedMaxOpenConnections)
}

func TestInitDBSupportsRelativeDatabasePath(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a relative database path with the required schema.",
		When:  "When the application opens it.",
		Then:  "Then SQLite is configured the same way as for absolute paths.",
	})

	// Given
	expectedJournalMode := "wal"

	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "database"), 0o755); err != nil {
		t.Fatalf("create database dir: %v", err)
	}
	dbPath := filepath.Join(dir, "database", "events.db")
	createMinimalProductionDB(t, dbPath)
	t.Chdir(dir)

	// When
	db, err := InitDB(filepath.Join("database", "events.db"))

	// Then
	if err != nil {
		t.Fatalf("expected relative database initialization to succeed: %v", err)
	}
	defer db.Close()

	assertSQLitePragmaString(t, db, "journal_mode", expectedJournalMode)
}

func TestInitDBFailsWhenRequiredTableMissing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a database missing a required core table.",
		When:  "When the application opens it.",
		Then:  "Then initialization fails without modifying the schema.",
	})

	// Given
	expectedErrorText := `required SQLite table "events" is missing`

	dbPath := filepath.Join(t.TempDir(), "events.db")
	db := openSQLiteForTest(t, dbPath)
	if _, err := db.Exec(`
		CREATE TABLE users(id INTEGER PRIMARY KEY);
		CREATE TABLE billettholdere(id INTEGER PRIMARY KEY);
		CREATE TABLE puljer(id TEXT PRIMARY KEY);
	`); err != nil {
		t.Fatalf("create incomplete schema: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close setup db: %v", err)
	}

	// When
	_, err := InitDB(dbPath)

	// Then
	assertErrorContains(t, err, expectedErrorText)
}

func TestInitDBFailsWhenDatabaseFileMissing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the configured database file does not exist.",
		When:  "When the application opens it.",
		Then:  "Then initialization fails instead of creating an empty database.",
	})

	// Given
	expectedErrorText := "database file does not exist"

	dbPath := filepath.Join(t.TempDir(), "missing.db")

	// When
	_, err := InitDB(dbPath)

	// Then
	assertErrorContains(t, err, expectedErrorText)
}

func assertSQLitePragmaInt(t *testing.T, db *sql.DB, pragmaName string, expected int) {
	t.Helper()

	var actual int
	if err := db.QueryRow("PRAGMA " + pragmaName + ";").Scan(&actual); err != nil {
		t.Fatalf("query %s pragma: %v", pragmaName, err)
	}
	if actual != expected {
		t.Fatalf("%s pragma mismatch\nexpected: %d\nactual:   %d", pragmaName, expected, actual)
	}
}

func assertSQLitePragmaAtLeast(t *testing.T, db *sql.DB, pragmaName string, expectedMinimum int) {
	t.Helper()

	var actual int
	if err := db.QueryRow("PRAGMA " + pragmaName + ";").Scan(&actual); err != nil {
		t.Fatalf("query %s pragma: %v", pragmaName, err)
	}
	if actual < expectedMinimum {
		t.Fatalf("%s pragma too low\nexpected at least: %d\nactual:            %d", pragmaName, expectedMinimum, actual)
	}
}

func assertSQLitePragmaString(t *testing.T, db *sql.DB, pragmaName string, expected string) {
	t.Helper()

	var actual string
	if err := db.QueryRow("PRAGMA " + pragmaName + ";").Scan(&actual); err != nil {
		t.Fatalf("query %s pragma: %v", pragmaName, err)
	}
	if !strings.EqualFold(actual, expected) {
		t.Fatalf("%s pragma mismatch\nexpected: %q\nactual:   %q", pragmaName, expected, actual)
	}
}

func assertForeignKeysRejectInvalidChild(t *testing.T, db *sql.DB) {
	t.Helper()

	if _, err := db.Exec("INSERT INTO child(parent_id) VALUES (999);"); err == nil {
		t.Fatalf("expected foreign key enforcement to reject invalid child row")
	}
}

func assertMaxOpenConnections(t *testing.T, db *sql.DB, expected int) {
	t.Helper()

	actual := db.Stats().MaxOpenConnections
	if actual != expected {
		t.Fatalf("max open connections mismatch\nexpected: %d\nactual:   %d", expected, actual)
	}
}

func assertErrorContains(t *testing.T, err error, expectedText string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error containing %q, got nil", expectedText)
	}
	if !strings.Contains(err.Error(), expectedText) {
		t.Fatalf("error mismatch\nexpected to contain: %q\nactual:              %v", expectedText, err)
	}
}

func createMinimalProductionDB(t *testing.T, dbPath string) {
	t.Helper()

	db := openSQLiteForTest(t, dbPath)
	defer db.Close()

	if _, err := db.Exec(`
		CREATE TABLE users(id INTEGER PRIMARY KEY);
		CREATE TABLE events(id INTEGER PRIMARY KEY);
		CREATE TABLE billettholdere(id INTEGER PRIMARY KEY);
		CREATE TABLE puljer(id TEXT PRIMARY KEY);
		CREATE TABLE parent(id INTEGER PRIMARY KEY);
		CREATE TABLE child(parent_id INTEGER REFERENCES parent(id));
	`); err != nil {
		t.Fatalf("create minimal production schema: %v", err)
	}
}

func openSQLiteForTest(t *testing.T, dbPath string) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatalf("open sqlite database: %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		t.Fatalf("ping sqlite database: %v", err)
	}
	return db
}
