package views

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Regncon/conorganizer/backup-migration/services"
)

func NewRootView(ctx context.Context, reg *services.Registry, window fyne.Window) fyne.CanvasObject {
	return container.New(layout.NewCenterLayout(), widget.NewLabel("asdasd"))
}
