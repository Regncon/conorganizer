package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

// Config Application configs and secrets
type Config struct {
	AWS_ENDPOINT_URL_S3   string
	AWS_ACCESS_KEY_ID     string
	AWS_SECRET_ACCESS_KEY string
	BUCKET_NAME           string
	DISCORD_WEBHOOK_URL   string
	DISCORD_SECRET_KEY    string
}

// Load Loads the required system environment variables
func Load(logger *slog.Logger) Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	return Config{
		AWS_ENDPOINT_URL_S3:   os.Getenv("AWS_ENDPOINT_URL_S3"),
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		BUCKET_NAME:           os.Getenv("BUCKET_NAME"),
		DISCORD_WEBHOOK_URL:   os.Getenv("DISCORD_WEBHOOK_URL"),
		DISCORD_SECRET_KEY:    os.Getenv("DISCORD_SECRET_KEY"),
	}
}
