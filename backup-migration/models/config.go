package models

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Secrets struct {
	AWS_ENDPOINT_URL_S3   string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	AWS_REGION            string
	BUCKET_NAME           string
	DB_PREFIX             string
}

type Configs struct {
	S3_secrets           Secrets
	S3_secrets_isMissing bool
	S3_prefix_old        string
	S3_prefix_new        string
	S3_status            string
	S3_process_status    string
	DB_path              string
	DB_process_status    string
	DB_validation_status string
}

func Config(logger *slog.Logger) *Configs {
	// Create new config
	cfg := &Configs{}

	// Check for envrionment variables
	err := godotenv.Load()
	if err != nil {
		logger.Error("Missing .env file")
		cfg.S3_secrets_isMissing = true
		return cfg
	}

	// Load variables from environment
	secrets := Secrets{
		AWS_ENDPOINT_URL_S3:   os.Getenv("AWS_ENDPOINT_URL_S3"),
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWS_REGION:            os.Getenv("AWS_REGION"),
		BUCKET_NAME:           os.Getenv("BUCKET_NAME"),
		DB_PREFIX:             os.Getenv("DB_PREFIX"),
	}

	// Update config
	cfg.S3_secrets_isMissing = false
	cfg.S3_secrets = secrets

	return cfg
}

func (c *Configs) Update() {
	// accept keyvalue as arg to update a field on &configs{}
}

func (c *Configs) CheckMissing() {
	// accept keyvalue as arg to update a field on &configs{}
}
