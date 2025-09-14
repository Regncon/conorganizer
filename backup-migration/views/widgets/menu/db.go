package menu

import (
	"context"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Regncon/conorganizer/backup-migration/services"
)

func MenuWidgetDB(ctx context.Context, reg *services.Registry, isConnected binding.Bool, isPrefixValid binding.Bool) fyne.CanvasObject {
	// Databinds
	isValidPath := binding.NewBool()
	isValidated := binding.NewBool()
	isUploaded := binding.NewBool()

	// Activity status
	isWorking := widget.NewActivity()

	// Labels
	menuLabel := widget.NewLabel("Database migration")

	// Input/output
	pathInput := widget.NewEntry()
	pathInput.PlaceHolder = "Choose local .db"

	// Buttons
	fileOpenerBtn := widget.NewButton("Browse", func() {
		fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			// Error opening window?
			if err != nil {
				dialog.ShowError(err, reg.Window)
			}
			// User canceled dialog
			if reader == nil {
				return
			}
			// Only allow .db files
			if reader.URI().Extension() != ".db" {
				pathInput.SetText("")
				dialog.ShowError(errors.New("you can only select .db files for upload! :> "), reg.Window)
				return
			}

			defer reader.Close()

			// Set file path
			uri := reader.URI()
			pathInput.SetText(uri.Path())

		}, reg.Window)
		fileDialog.Show()
	})

	validateBtn := widget.NewButton("Validate", func() {
		isWorking.Show()
		isWorking.Start()

		go func() {
			// do something
			fyne.Do(func() {
				isValidated.Set(true)
				isWorking.Stop()
				isWorking.Hide()
			})
		}()
	})

	uploadBtn := widget.NewButton("Upload", func() {
		isWorking.Show()
		isWorking.Start()

		go func() {
			// do something
			fyne.Do(func() {
				isUploaded.Set(true)
				isWorking.Stop()
				isWorking.Hide()
			})
		}()
	})

	// containers
	labelContainer := container.NewBorder(nil, nil, menuLabel, isWorking)
	fileOpenContainer := container.NewBorder(nil, nil, nil, fileOpenerBtn, pathInput)

	// Initial gui states
	validateBtn.Disable()
	uploadBtn.Disable()

	// Watchers
	isValidPath.AddListener(binding.NewDataListener(func() {
		if val, _ := isValidPath.Get(); val {
			validateBtn.Enable()
		} else {
			validateBtn.Disable()
			uploadBtn.Disable()
		}

	}))
	isValidated.AddListener(binding.NewDataListener(func() {
		validated, _ := isValidated.Get()
		prefixed, _ := isPrefixValid.Get()

		fmt.Println("validation changed")

		if validated && prefixed {
			uploadBtn.Enable()
		} else {
			uploadBtn.Disable()
		}

	}))
	isPrefixValid.AddListener(binding.NewDataListener(func() {
		validated, _ := isValidated.Get()
		prefixed, _ := isPrefixValid.Get()

		fmt.Println("prefix changed")

		if validated && prefixed {
			uploadBtn.Enable()
		} else {
			uploadBtn.Disable()
		}

	}))
	isConnected.AddListener(binding.NewDataListener(func() {
		if val, _ := isConnected.Get(); val {
			uploadBtn.Show()
		} else {
			uploadBtn.Hide()
		}
	}))

	// Validators
	validateExtension := validation.NewRegexp(`(?i)\.db$`, "must end with .db")
	pathInput.Validator = func(s string) error {
		// Check extension
		if err := validateExtension(s); err != nil {
			isValidPath.Set(false)
			return err
		}
		isValidPath.Set(true)
		return nil
	}

	return container.NewVBox(
		labelContainer,
		fileOpenContainer,
		validateBtn,
		uploadBtn,
	)
}
