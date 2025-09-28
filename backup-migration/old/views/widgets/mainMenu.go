package widgets

import (
	"backup-migration/services"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func MainMenu(reg *services.Registry, window fyne.Window) fyne.CanvasObject {
	// Load required files into app
	selectLabel := widget.NewLabel("Select files")
	selectDb := NewFileOpener(window, "db", ".db", reg.AppState.DB.Path)
	selectEnv := NewFileOpener(window, "env", ".env", reg.AppState.ENV.Path)
	selectContainer := container.New(layout.NewVBoxLayout(), selectLabel, selectEnv, selectDb)

	// Select new database prefix for litestream sync
	prefixLabel := widget.NewLabel("Choose new prefix")
	prefixInput := widget.NewEntry()
	prefixInput.SetPlaceHolder("S3 loading failed")
	prefixContainer := container.New(layout.NewVBoxLayout(), prefixLabel, prefixInput)
	prefixContainer.Hide()

	// Run migration
	migrationLabel := widget.NewLabel("Migration")
	validateDb := widget.NewButtonWithIcon("validate database", theme.SearchReplaceIcon(), func() {
		log.Print("validating .db file")
	})
	uploadDb := widget.NewButtonWithIcon("Start migration", theme.UploadIcon(), func() {

	})
	migrationContainer := container.New(layout.NewVBoxLayout(), migrationLabel, validateDb, uploadDb)

	// Combine inputs into a menu container
	menuContainer := container.New(layout.NewVBoxLayout(), selectContainer, prefixContainer, migrationContainer)
	line := canvas.NewLine(color.White)
	border := container.NewBorder(nil, nil, nil, line, menuContainer)
	padding := container.NewPadded(border)

	return padding
}
