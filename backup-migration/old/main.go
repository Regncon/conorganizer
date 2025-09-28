package main

import (
	"backup-migration/internal/app"
	"context"
	"log"
	"log/slog"
	"os"
)

func main() {
	// Set up logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Run app entrypoint
	if err := app.Run(context.Background(), logger); err != nil {
		log.Fatal(err)
	}
}
