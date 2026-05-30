package service

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const (
	defaultSQLiteBusyTimeoutMillis = 5000
	defaultSQLiteSynchronous       = "NORMAL"
)

var defaultRequiredSQLiteTables = []string{
	"users",
	"events",
	"billettholdere",
	"puljer",
}

type SQLiteConfig struct {
	BusyTimeoutMillis int
	Synchronous       string
	RequireWAL        bool
	MaxOpenConns      int
	MaxIdleConns      int
	RequiredTables    []string
}

func DefaultSQLiteConfig() SQLiteConfig {
	return SQLiteConfig{
		BusyTimeoutMillis: defaultSQLiteBusyTimeoutMillis,
		Synchronous:       defaultSQLiteSynchronous,
		RequireWAL:        true,
		MaxOpenConns:      1,
		MaxIdleConns:      1,
		RequiredTables:    append([]string(nil), defaultRequiredSQLiteTables...),
	}
}

func InitDB(databaseFileName string) (*sql.DB, error) {
	return InitDBWithConfig(databaseFileName, DefaultSQLiteConfig())
}

func InitDBWithConfig(databaseFileName string, config SQLiteConfig) (*sql.DB, error) {
	if strings.TrimSpace(databaseFileName) == "" {
		return nil, fmt.Errorf("database path is required")
	}

	config = normalizeSQLiteConfig(config)

	dir := filepath.Dir(databaseFileName)
	if !isMemorySQLiteDatabase(databaseFileName) && dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return nil, fmt.Errorf("directory path does not exist: %s", dir)
		} else if err != nil {
			return nil, fmt.Errorf("could not access database directory %q: %w", dir, err)
		}
	}

	if !isMemorySQLiteDatabase(databaseFileName) {
		if _, err := os.Stat(databaseFileName); os.IsNotExist(err) {
			return nil, fmt.Errorf("database file does not exist: %s", databaseFileName)
		} else if err != nil {
			return nil, fmt.Errorf("could not access database file %q: %w", databaseFileName, err)
		}
	}

	db, err := sql.Open("sqlite", sqliteDSN(databaseFileName, config))
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	closeOnError := true
	defer func() {
		if closeOnError {
			_ = db.Close()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping DB: %w", err)
	}

	if err := verifySQLiteConfiguration(ctx, db, config); err != nil {
		return nil, err
	}

	if err := verifyRequiredSQLiteTables(ctx, db, config.RequiredTables); err != nil {
		return nil, err
	}

	closeOnError = false
	return db, nil
}

func normalizeSQLiteConfig(config SQLiteConfig) SQLiteConfig {
	if config.BusyTimeoutMillis <= 0 {
		config.BusyTimeoutMillis = defaultSQLiteBusyTimeoutMillis
	}
	if strings.TrimSpace(config.Synchronous) == "" {
		config.Synchronous = defaultSQLiteSynchronous
	}
	config.Synchronous = strings.ToUpper(strings.TrimSpace(config.Synchronous))
	if config.MaxOpenConns <= 0 {
		config.MaxOpenConns = 1
	}
	if config.MaxIdleConns <= 0 || config.MaxIdleConns > config.MaxOpenConns {
		config.MaxIdleConns = config.MaxOpenConns
	}
	if config.RequiredTables == nil {
		config.RequiredTables = append([]string(nil), defaultRequiredSQLiteTables...)
	}
	return config
}

func sqliteDSN(databaseFileName string, config SQLiteConfig) string {
	values := url.Values{}
	if config.RequireWAL {
		values.Add("_pragma", "journal_mode(WAL)")
	}
	values.Add("_pragma", "foreign_keys(ON)")
	values.Add("_pragma", fmt.Sprintf("busy_timeout(%d)", config.BusyTimeoutMillis))
	values.Add("_pragma", fmt.Sprintf("synchronous(%s)", config.Synchronous))

	if filepath.IsAbs(databaseFileName) {
		uri := url.URL{
			Scheme:   "file",
			Path:     databaseFileName,
			RawQuery: values.Encode(),
		}
		return uri.String()
	}

	uri := url.URL{
		Scheme:   "file",
		Opaque:   databaseFileName,
		RawQuery: values.Encode(),
	}
	return uri.String()
}

func verifySQLiteConfiguration(ctx context.Context, db *sql.DB, config SQLiteConfig) error {
	var foreignKeys int
	if err := db.QueryRowContext(ctx, "PRAGMA foreign_keys;").Scan(&foreignKeys); err != nil {
		return fmt.Errorf("verify SQLite foreign_keys pragma: %w", err)
	}
	if foreignKeys != 1 {
		return fmt.Errorf("SQLite foreign_keys pragma is disabled")
	}

	if config.RequireWAL {
		var journalMode string
		if err := db.QueryRowContext(ctx, "PRAGMA journal_mode;").Scan(&journalMode); err != nil {
			return fmt.Errorf("verify SQLite journal_mode pragma: %w", err)
		}
		if !strings.EqualFold(journalMode, "wal") {
			return fmt.Errorf("SQLite journal_mode is %q, expected WAL", journalMode)
		}
	}

	var busyTimeoutMillis int
	if err := db.QueryRowContext(ctx, "PRAGMA busy_timeout;").Scan(&busyTimeoutMillis); err != nil {
		return fmt.Errorf("verify SQLite busy_timeout pragma: %w", err)
	}
	if busyTimeoutMillis < config.BusyTimeoutMillis {
		return fmt.Errorf("SQLite busy_timeout is %dms, expected at least %dms", busyTimeoutMillis, config.BusyTimeoutMillis)
	}

	expectedSynchronous, err := sqliteSynchronousValue(config.Synchronous)
	if err != nil {
		return err
	}
	var synchronous int
	if err := db.QueryRowContext(ctx, "PRAGMA synchronous;").Scan(&synchronous); err != nil {
		return fmt.Errorf("verify SQLite synchronous pragma: %w", err)
	}
	if synchronous != expectedSynchronous {
		return fmt.Errorf("SQLite synchronous is %d, expected %d (%s)", synchronous, expectedSynchronous, config.Synchronous)
	}

	return nil
}

func sqliteSynchronousValue(synchronous string) (int, error) {
	switch strings.ToUpper(strings.TrimSpace(synchronous)) {
	case "OFF":
		return 0, nil
	case "NORMAL":
		return 1, nil
	case "FULL":
		return 2, nil
	case "EXTRA":
		return 3, nil
	default:
		return 0, fmt.Errorf("unsupported SQLite synchronous mode %q", synchronous)
	}
}

func verifyRequiredSQLiteTables(ctx context.Context, db *sql.DB, requiredTables []string) error {
	for _, tableName := range requiredTables {
		var exists int
		if err := db.QueryRowContext(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM sqlite_schema
				WHERE type = 'table'
				  AND name = ?
			);
		`, tableName).Scan(&exists); err != nil {
			return fmt.Errorf("verify required SQLite table %q: %w", tableName, err)
		}
		if exists != 1 {
			return fmt.Errorf("required SQLite table %q is missing", tableName)
		}
	}
	return nil
}

func isMemorySQLiteDatabase(databaseFileName string) bool {
	return databaseFileName == ":memory:" || strings.Contains(databaseFileName, "mode=memory")
}
