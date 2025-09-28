package views

import (
	"backup-migration/services"
	"backup-migration/views/widgets"
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

func NewRoot(ctx context.Context, reg *services.Registry, w fyne.Window) fyne.CanvasObject {
	menuBox := widgets.MainMenu(reg, w)
	outputView := OutputView(reg, w)

	content := container.New(layout.NewHBoxLayout(), menuBox, outputView)

	return content
}
