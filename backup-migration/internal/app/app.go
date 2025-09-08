package app

import (
	"backup-migration/services"
	"backup-migration/views"
	"context"
	"log/slog"

	"fyne.io/fyne/v2/app"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	// Create a new dependency container for sharing services
	reg := services.NewRegistry()

	// Setup GUI with Fyne
	fyneApp := app.New()
	window := fyneApp.NewWindow("RegnCon - Database Migration Tool™")

	// Load root layout / content
	root := views.NewRoot(window, reg, logger)
	window.SetContent(root)

	// Graceful shutdown if ctx is cancelled
	go func() {
		<-ctx.Done()
		fyneApp.Quit()
	}()

	// Start GUI
	window.ShowAndRun()

	return nil
}
