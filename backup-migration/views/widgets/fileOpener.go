package widgets

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewFileOpener(window fyne.Window, label string, ext string, bind binding.String) fyne.CanvasObject {
	filePath := widget.NewEntry()
	filePath.SetPlaceHolder("No file selected")
	filePath.Resize(fyne.NewSize(200, 10))

	button := widget.NewButtonWithIcon(label, theme.FolderOpenIcon(), func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, window)
				return
			}
			if reader == nil {
				return
			}
			defer reader.Close()

			// File info
			uri := reader.URI()
			fileExt := uri.Extension()

			// Check extension
			if fileExt != ext {
				dialog.ShowInformation("Note", "Selected file type is not supported", window)
				return
			}

			// Update app state
			filePath.SetText(uri.Path())
			bind.Set(filePath.Text)
			fmt.Printf("%s loaded: %s \n", ext, uri.Path())

		}, window)
		fileDialog.Show()
	})

	container := container.New(layout.NewHBoxLayout(), filePath, button)

	return container
}
