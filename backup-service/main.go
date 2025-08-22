package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Regncon/conorganizer/backup-service/config"
	"github.com/Regncon/conorganizer/backup-service/services"
	"github.com/Regncon/conorganizer/backup-service/utils"
)

func main() {
	// Set timezone
	os.Setenv("TZ", "Europe/Oslo")

	// Set up logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Creating required folders
	if err := utils.CreateRequiredFolders(); err != nil {
		logger.Error("Could not create required folders", "error", err)
		os.Exit(1)
	}

	// Load config with params and secrets
	cfg := config.Load(logger)

	// Initialize database
	db, err := services.InitDB()
	if err != nil {
		logger.Error("Could not initialize DB", "error", err)
	}
	defer db.Close()

	// Start S3 client
	s3Client, err := services.NewS3Client(cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize S3 client", "error", err)
	}

	// Define backup service
	backupService := services.NewBackupService(cfg, db, s3Client, logger)

	// Start scheduler
	err = services.StartScheduler(backupService, logger)
	if err != nil {
		logger.Error("Failed to start backup scheduler", "error", err)
	}

	// Start dashboard web server
	// todo

	// Block forever, with graceful shutdown support
	logger.Info("All startup tasks completed, now running server")
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Info("Shutdown signal received. Exiting.")
}
