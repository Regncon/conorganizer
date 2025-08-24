package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Regncon/conorganizer/backup-service/config"
	"github.com/Regncon/conorganizer/backup-service/services"
	"github.com/Regncon/conorganizer/backup-service/utils"
	"github.com/Regncon/conorganizer/backup-service/web"
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
	logger.Info("All startup tasks completed, now running server")

	// Run once on startup
	/* backupService.Manual()
	backupService.Hourly()
	backupService.Daily()
	backupService.Weekly()
	backupService.Yearly() */

	// Create new context for dashboard to avoid polution
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Create dashboard web server
	server := &http.Server{
		Addr:    ":8080",
		Handler: web.NewRouter(ctx, logger, db),
	}

	go func() {
		// Start server
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Failed to start webserver", "errorr", err)
			os.Exit(1)
		}
	}()

	// Wait for Ctrl+C
	<-ctx.Done()
	logger.Info("Shutting down...")

	// Gracefully shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown server gracefully", "err", err)
	} else {
		logger.Info("Server shutdown complete")
	}
}
