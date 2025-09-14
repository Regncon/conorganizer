package menu

import (
	"context"
	"errors"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/Regncon/conorganizer/backup-migration/services"
)

func MenuWidgetPrefix(ctx context.Context, reg *services.Registry, isConnected binding.Bool, prefixList binding.StringList, isValidPrefix binding.Bool) fyne.CanvasObject {
	// Labels
	newPrefixLabel := widget.NewLabel("Set new prefix")
	oldPrefixLabel := widget.NewLabel("Existing prefixes")

	// Input/output
	prefixInput := widget.NewEntry()
	prefixInput.PlaceHolder = "Prefix must be unique"

	// Table
	existingTable := widget.NewListWithData(prefixList,
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(di binding.DataItem, co fyne.CanvasObject) {
			co.(*widget.Label).Bind(di.(binding.String))
		},
	)

	// containers
	oldPrefixContainer := container.NewBorder(oldPrefixLabel, nil, nil, nil, existingTable)
	newPrefixContainer := container.NewBorder(newPrefixLabel, nil, nil, nil, prefixInput)
	menuContainer := container.NewBorder(oldPrefixContainer, nil, nil, nil, newPrefixContainer)

	// Watchers
	isConnected.AddListener(binding.NewDataListener(func() {
		if val, _ := isConnected.Get(); val {
			menuContainer.Show()
		} else {
			menuContainer.Hide()
		}
	}))
	prefixList.AddListener(binding.NewDataListener(func() {
		val, _ := prefixList.Get()

		// todo, fix cell height

		if len(val) > 0 {
			oldPrefixContainer.Show()
		} else {
			oldPrefixContainer.Hide()
		}
	}))

	// Validators
	limitPrefixChars := validation.NewRegexp(`(?i)^[a-z0-9-]+$`, "Only alphanumerical and dashes allowed")
	prefixInput.Validator = func(s string) error {
		// Only accept aA09-
		if err := limitPrefixChars(s); err != nil {
			isValidPrefix.Set(false)
			return err
		}

		// Check if prefix already exists
		existingPrefixesVal, _ := prefixList.Get()
		for _, value := range existingPrefixesVal {
			if strings.EqualFold(s, value) {
				isValidPrefix.Set(false)
				return errors.New("prefix already exists in S3")
			}
		}

		isValidPrefix.Set(true)
		return nil
	}

	// Initial gui state
	oldPrefixContainer.Hide()

	return menuContainer
}
