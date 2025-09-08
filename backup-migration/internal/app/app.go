package app

import (
	"backup-migration/services"
	"backup-migration/views"
	"context"

	"fyne.io/fyne/v2/app"
)

func Run(ctx context.Context, reg *services.Registry) error {
	// Setup GUI with Fyne
	fyneApp := app.New()
	window := fyneApp.NewWindow("RegnCon - Database Migration Tool™")

	// Load root layout / content
	root := views.NewRoot(ctx, reg, window)
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
