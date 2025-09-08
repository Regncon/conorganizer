package views

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Run(_ context.Context, app fyne.App) (fyne.Window, error) {
	// Create a new Fyne ui window
	window := app.NewWindow("RegnCon - Database Migration Tool™")

	// Add some content
	hello := widget.NewLabel("Hellow Fyne!")
	window.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	return window, nil
}
