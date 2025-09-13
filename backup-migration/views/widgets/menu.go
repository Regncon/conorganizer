package widgets

import (
	"context"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Regncon/conorganizer/backup-migration/services"
	"github.com/Regncon/conorganizer/backup-migration/utils"
)

func RootMenu(ctx context.Context, reg *services.Registry) fyne.CanvasObject {
	// Menu label
	menuLabel := widget.NewLabel("Select files for migration")

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

	// DB section
	dbPath := widget.NewEntry()
	dbPath.PlaceHolder = "Choose .db to upload it"
	dbButton := widget.NewButtonWithIcon(".db", theme.FolderOpenIcon(), func() {
		err := reg.DB.Load(reg.Config)
		if err != nil {
			dbPath.Text = fmt.Sprint(err)
		}
		dbPath.Refresh()
	})
	dbButton.Resize(fyne.NewSize(50, 0))
	dbContainer := container.NewBorder(nil, nil, nil, dbButton, dbPath)

	// Prefix section
	prefixLabel := widget.NewLabel("Database prefix")
	prefix := widget.NewEntry()
	prefix.PlaceHolder = reg.Config.S3.Prefix
	prefixContainer := container.NewVBox(prefixLabel, prefix)

	// S3 section
	s3Activity := widget.NewActivity()
	s3Label := widget.NewLabel("S3 Storage")
	s3LabelGroup := container.NewBorder(nil, nil, s3Label, s3Activity)
	s3ConnectButton := widget.NewButton("Connect", func() {
		s3Activity.Show()
		s3Activity.Start()

		go func() {
			err := reg.S3.Connect(&reg.Config)
			if err != nil {
				fmt.Println("Attempted to start S3 client without cfg")
			}
			fyne.Do(func() {
				s3Activity.Stop()
				s3Activity.Hide()
				s3Label.Text = "S3 Storage connected"
				s3Label.Refresh()
			})
		}()
	})
	s3Latest := canvas.NewText("", color.White)
	s3Latest.Hide()
	s3LatestButton := widget.NewButton("Check for latest", func() {
		s3Activity.Show()
		s3Activity.Start()

		go func() {
			obj, err := reg.S3.GetLatestBackup(&reg.Config)
			if err != nil {
				fmt.Println("Attempted to start S3 client without cfg")
			}
			fyne.Do(func() {
				s3Activity.Stop()
				s3Activity.Hide()

				s3Latest.Text = utils.TimeAgo(obj.LastModified)
				s3Latest.Show()
			})
		}()
	})
	s3Container := container.NewVBox(s3LabelGroup, s3ConnectButton, s3LatestButton, s3Latest)

	// Migration section
	migrationLabel := widget.NewLabel("Migrations")
	migrationContainer := container.NewVBox(migrationLabel)

	// Combine menu items into one container
	content := container.NewVBox(menuLabel, envContainer, dbContainer, prefixContainer, s3Container, migrationContainer)
	return content
}
