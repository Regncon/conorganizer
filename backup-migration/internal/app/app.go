package app

import (
	"backup-migration/views"
	"context"
	"log/slog"

	fyneApp "fyne.io/fyne/v2/app"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	// Setup GUI with Fyne
	app := fyneApp.New()
	fyne, err := views.Run(ctx, app)
	if err != nil {
		logger.Error("Failed to start Fyne GUI", "Reason", err)
	}

	// Graceful shutdown if ctx is cancelled
	go func() {
		<-ctx.Done()
		app.Quit()
	}()

	// Start GUI
	fyne.ShowAndRun()

	return nil
}
