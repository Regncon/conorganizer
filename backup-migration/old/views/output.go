package views

import (
	"backup-migration/services"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func OutputView(reg *services.Registry, window fyne.Window) fyne.CanvasObject {
	text1 := canvas.NewText("Database not set", color.White)
	text2 := canvas.NewText("Env not set", color.White)
	text3 := canvas.NewText("Just hang on tight...", color.White)
	textContent := container.New(layout.NewVBoxLayout(), text1, text2, layout.NewSpacer(), text3)

	tabs := container.NewAppTabs(
		container.NewTabItem("Status", textContent),
		container.NewTabItem("Environment", widget.NewLabel("Hello")),
		container.NewTabItem("S3 Storage", widget.NewLabel("World!")),
	)

	// Create container
	tabs.SetTabLocation(container.TabLocationTop)
	content := container.NewVBox(tabs)

	// Refresh on state change
	reg.AppState.OnChange.AddListener(binding.NewDataListener(func() {
		dbPath, _ := reg.AppState.DB.Path.Get()
		if dbPath != "" {
			text1.Text = "Database path: " + dbPath
		}

		envPath, _ := reg.AppState.ENV.Path.Get()
		if envPath != "" {
			text2.Text = "Env path: " + envPath
		}

	}))

	return content
}
