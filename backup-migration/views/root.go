package views

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Regncon/conorganizer/backup-migration/services"
	"github.com/Regncon/conorganizer/backup-migration/views/widgets"
)

func NewRootView(ctx context.Context, reg *services.Registry) fyne.CanvasObject {
	// Main menu
	menu := widgets.RootMenu(ctx, reg)
	menuScroll := container.NewVScroll(menu)
	menuScroll.SetMinSize(fyne.NewSize(260, 0))

	// Main content
	mainLabel := widget.NewLabel("Logs")
	mainLogs := widget.NewTextGrid()
	mainContainer := container.NewVBox(mainLabel, mainLogs)

	// Combining main and menu
	spacer := layout.NewSpacer()
	return container.NewBorder(nil, nil, menuScroll, spacer, mainContainer)
}
