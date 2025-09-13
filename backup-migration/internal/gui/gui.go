package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/Regncon/conorganizer/backup-migration/services"
	"github.com/Regncon/conorganizer/backup-migration/views"
)

func NewFyneApp(ctx context.Context, reg *services.Registry) {
	// Init fyne
	reg.App = app.NewWithID("RegnconMigrationTool")

	// Create application window
	reg.Window = reg.App.NewWindow("RegnCon - Database Migration Tool™")
	reg.Window.Resize(fyne.NewSize(800, 400))

	// Draw root view
	rootView := views.NewRootView(ctx, reg)
	reg.Window.SetContent(rootView)

	// Graceful shutdown if ctx is cancelled
	go func() {
		<-ctx.Done()
		reg.App.Quit()
	}()

	// Start GUI
	reg.Window.ShowAndRun()
}
