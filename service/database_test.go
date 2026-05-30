package service

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "modernc.org/sqlite"
)

func TestInitDBAppliesProductionSQLiteSettings(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "events.db")
	createMinimalProductionDB(t, dbPath)

	db, err := InitDB(dbPath)
	if err != nil {
		t.Fatalf("InitDB() error = %v", err)
	}
	defer db.Close()

	var foreignKeys int
	if err := db.QueryRow("PRAGMA foreign_keys;").Scan(&foreignKeys); err != nil {
		t.Fatalf("query foreign_keys pragma: %v", err)
	}
	if foreignKeys != 1 {
		t.Fatalf("foreign_keys = %d, want 1", foreignKeys)
	}

	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("query journal_mode pragma: %v", err)
	}
	if !strings.EqualFold(journalMode, "wal") {
		t.Fatalf("journal_mode = %q, want WAL", journalMode)
	}

	var busyTimeoutMillis int
	if err := db.QueryRow("PRAGMA busy_timeout;").Scan(&busyTimeoutMillis); err != nil {
		t.Fatalf("query busy_timeout pragma: %v", err)
	}
	if busyTimeoutMillis < defaultSQLiteBusyTimeoutMillis {
		t.Fatalf("busy_timeout = %d, want at least %d", busyTimeoutMillis, defaultSQLiteBusyTimeoutMillis)
	}

	var synchronous int
	if err := db.QueryRow("PRAGMA synchronous;").Scan(&synchronous); err != nil {
		t.Fatalf("query synchronous pragma: %v", err)
	}
	if synchronous != 1 {
		t.Fatalf("synchronous = %d, want 1 for NORMAL", synchronous)
	}

	if _, err := db.Exec("INSERT INTO child(parent_id) VALUES (999);"); err == nil {
		t.Fatalf("insert violating foreign key succeeded, want failure")
	}

	stats := db.Stats()
	if stats.MaxOpenConnections != 1 {
		t.Fatalf("MaxOpenConnections = %d, want 1", stats.MaxOpenConnections)
	}
}

func TestInitDBSupportsRelativeDatabasePath(t *testing.T) {
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "database"), 0o755); err != nil {
		t.Fatalf("create database dir: %v", err)
	}
	dbPath := filepath.Join(dir, "database", "events.db")
	createMinimalProductionDB(t, dbPath)
	t.Chdir(dir)

	db, err := InitDB(filepath.Join("database", "events.db"))
	if err != nil {
		t.Fatalf("InitDB() with relative path error = %v", err)
	}
	defer db.Close()

	var journalMode string
	if err := db.QueryRow("PRAGMA journal_mode;").Scan(&journalMode); err != nil {
		t.Fatalf("query journal_mode pragma: %v", err)
	}
	if !strings.EqualFold(journalMode, "wal") {
		t.Fatalf("journal_mode = %q, want WAL", journalMode)
	}
}

func TestInitDBFailsWhenRequiredTableMissing(t *testing.T) {
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

	_, err := InitDB(dbPath)
	if err == nil {
		t.Fatalf("InitDB() error = nil, want missing table error")
	}
	if !strings.Contains(err.Error(), `required SQLite table "events" is missing`) {
		t.Fatalf("InitDB() error = %v, want missing events table", err)
	}
}

func TestInitDBFailsWhenDatabaseFileMissing(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "missing.db")

	_, err := InitDB(dbPath)
	if err == nil {
		t.Fatalf("InitDB() error = nil, want missing database file error")
	}
	if !strings.Contains(err.Error(), "database file does not exist") {
		t.Fatalf("InitDB() error = %v, want missing database file error", err)
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
