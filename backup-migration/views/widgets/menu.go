package widgets

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Regncon/conorganizer/backup-migration/services"
	"github.com/Regncon/conorganizer/backup-migration/views/widgets/menu"
)

func RootMenu(ctx context.Context, reg *services.Registry) fyne.CanvasObject {
	// Menu label
	menuLabel := widget.NewLabel("Select config file")

	// Env section
	envPath := widget.NewEntry()
	envPath.PlaceHolder = reg.Config.EnvPath
	envButton := widget.NewButtonWithIcon(".env", theme.FolderOpenIcon(), func() {
		err := reg.Config.LoadEnv(".env")
		if err != nil {
			envPath.Text = fmt.Sprint(err)
		}

		envPath.Text = reg.Config.EnvPath
		envPath.Refresh()
	})
	envContainer := container.NewBorder(nil, nil, nil, envButton, envPath)

	// S3 section
	// todo move bindings to reg?
	prefixList := binding.NewStringList()
	isConnected := binding.NewBool()
	isConnected.Set(false)
	s3menu := menu.MenuWidgetS3(ctx, reg, isConnected, prefixList)

	// Database section
	dbMenu := menu.MenuWidgetDB(ctx, reg, isConnected)

	// Prefix section
	prefix := menu.MenuWidgetPrefix(ctx, reg, isConnected, prefixList)

	// Application cta
	exitButton := widget.NewButton("Exit", func() {
		reg.App.Quit()
	})
	bottomMenuContainer := container.NewVBox(exitButton)

	// Spacer
	menuSpacer := layout.NewSpacer()

	// Combine menu items into one container
	content := container.NewVBox(
		menuLabel,
		envContainer,
		s3menu,
		prefix,
		dbMenu,
		menuSpacer,
		bottomMenuContainer,
		menuSpacer,
		menuSpacer,
	)
	return content
}
