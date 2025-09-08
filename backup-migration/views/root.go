package views

import (
	"backup-migration/services"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func NewRoot(w fyne.Window, reg *services.Registry, baseLogger *slog.Logger) fyne.CanvasObject {
	console, consoleObj := NewConsoleWidget(1000)

	// Wire a UI logger that writes into the console widget.
	reg.UILogger = slog.New(slog.NewTextHandler(console, nil))
	reg.ConsoleWriter = console

	// Example action menu; replace with real service calls later.
	actions := []struct {
		Title string
		Do    func()
	}{
		{"Say hello", func() { reg.UILogger.Info("Hello from action") }},
		{"Clear console", func() { console.Clear() }},
	}

	menu := widget.NewList(
		func() int { return len(actions) },
		func() fyne.CanvasObject { return widget.NewButton("action", nil) },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			btn := o.(*widget.Button)
			btn.SetText(actions[i].Title)
			btn.OnTapped = actions[i].Do
		},
	)

	split := container.NewHSplit(menu, consoleObj)
	split.Offset = 0.22
	return split
}
