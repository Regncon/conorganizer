package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/Regncon/conorganizer/backup-migration/config"
)

type DBClient struct {
	DB *sql.DB
}

func NewDBClient() *DBClient {
	return &DBClient{}
}

func (c *DBClient) Load(cfg config.Config) error {
	fmt.Println(cfg.DB.Path)
	if cfg.DB.Path == "" {
		return errors.New("DB Load must be called with a valid cfg path")
	}

	// Check if file exists
	fileInfo, err := os.Stat(cfg.DB.Path)
	if err != nil {
		return fmt.Errorf("DB Load: %s: %w", cfg.DB.Path, err)
	}
	if fileInfo.IsDir() {
		return fmt.Errorf("DB Load: %s is a dir", cfg.DB.Path)
	}

	// Open database
	db, err := sql.Open("sqlite", cfg.DB.Path)
	if err != nil {
		return fmt.Errorf("failed to open DB: %w", err)
	}

	// Close existing connections before assigning
	if c.DB != nil {
		c.DB.Close()
	}

	// Assign new db
	c.DB = db
	fmt.Println("New database loaded")

	return nil
}

func (c *DBClient) Close() error {
	if c.DB == nil {
		return fmt.Errorf("DB Close: not connected anyway /shrug")
	}
	err := c.DB.Close()
	c.DB = nil
	return err
}

func (c *DBClient) Validate() error {
	if c.DB == nil {
		return errors.New("you've somehow managed to validate on a db file that hsn't been loaded, impressive")
	}

	// Limit validation runtime in case of stck
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if db file is locked
	_, _ = c.DB.ExecContext(ctx, "PRAGMA busy_timeout = 3000")

	// Checks
	if err := c.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("validate: ping failed: %w", err)
	}
	var one int
	if err := c.DB.QueryRowContext(ctx, "SELECT 1").Scan(&one); err != nil {
		return fmt.Errorf("validate: trivial query failed: %w", err)
	}
	var cnt int
	if err := c.DB.QueryRowContext(ctx, "SELECT count(*) FROM sqlite_master").Scan(&cnt); err != nil {
		return fmt.Errorf("validate: reading sqlite_master failed: %w", err)
	}
	msgs, err := quickCheck(ctx, c.DB)
	if err != nil {
		return fmt.Errorf("validate: quick_check failed: %w", err)
	}
	if len(msgs) > 0 {
		return fmt.Errorf("validate: integrity errors: %s", strings.Join(msgs, "; "))
	}

	fmt.Println("Validation passed woho")

	return nil
}

// todo, move to util|tests
func quickCheck(ctx context.Context, db *sql.DB) ([]string, error) {
	rows, err := db.QueryContext(ctx, "PRAGMA quick_check")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return nil, err
		}
		if strings.ToLower(strings.TrimSpace(s)) != "ok" {
			msgs = append(msgs, s)
		}
	}
	return msgs, rows.Err()
}
