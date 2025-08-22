package config

import (
	"log/slog"
	"os"

	"github.com/Regncon/conorganizer/backup-service/models"
	"github.com/joho/godotenv"
)

// Load Loads the required system environment variables
func Load(logger *slog.Logger) models.Config {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	return models.Config{
		AWS_ENDPOINT_URL_S3:   os.Getenv("AWS_ENDPOINT_URL_S3"),
		AWS_ACCESS_KEY_ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWS_SECRET_ACCESS_KEY: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AWS_REGION:            os.Getenv("AWS_REGION"),
		BUCKET_NAME:           os.Getenv("BUCKET_NAME"),
		DB_PREFIX:             os.Getenv("DB_PREFIX"),
		DISCORD_WEBHOOK_URL:   os.Getenv("DISCORD_WEBHOOK_URL"),
		DISCORD_SECRET_KEY:    os.Getenv("DISCORD_SECRET_KEY"),
	}
}
