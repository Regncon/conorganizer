package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

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
	c.DB = db

	return nil
}

func (c *DBClient) Close(cfg config.Config) error {
	if c.DB == nil {
		return fmt.Errorf("DB Close: not connected anyway /shrug")
	}
	err := c.DB.Close()
	c.DB = nil
	return err
}
