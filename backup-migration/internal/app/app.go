package app

import (
	"backup-migration/services"
	"backup-migration/views"
	"context"
	"log/slog"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
)

func Run(ctx context.Context, logger *slog.Logger) error {
	// Setup GUI with Fyne
	fyneApp := app.NewWithID("FyneTymes")
	window := fyneApp.NewWindow("RegnCon - Database Migration Tool™")

	// Create a new dependency container for sharing services
	reg := services.NewRegistry(logger)
	reg.S3.GetExistingPrefixes(ctx)

	// Load root layout / content
	root := views.NewRoot(ctx, reg, window)
	window.SetContent(root)

	// Trigger refresh on content change
	reg.AppState.OnChange.AddListener(binding.NewDataListener(func() {
		root.Refresh()
	}))

	// Graceful shutdown if ctx is cancelled
	go func() {
		<-ctx.Done()
		fyneApp.Quit()
	}()

	// Start GUI
	window.ShowAndRun()

	return nil
}
