package main

import (
	"backup-migration/internal/app"
	"backup-migration/services"
	"context"
	"log"
	"log/slog"
	"os"
)

func main() {
	// Set up logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Create a new dependency container for sharing services
	reg := services.NewRegistry(logger)

	// Run app entrypoint
	if err := app.Run(context.Background(), reg); err != nil {
		log.Fatal(err)
	}
}
