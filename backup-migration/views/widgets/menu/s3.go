package menu

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/Regncon/conorganizer/backup-migration/services"
	"github.com/Regncon/conorganizer/backup-migration/utils"
)

func MenuWidgetS3(ctx context.Context, reg *services.Registry, isConnected binding.Bool, prefixList binding.StringList) fyne.CanvasObject {
	// Activity status
	isWorking := widget.NewActivity()

	// Data binds
	connectedError := binding.NewString()
	lastModified := binding.NewString()
	downloadPath := binding.NewString()

	// Data binds initial values
	isConnected.Set(false)

	// Labels
	sectionLabel := widget.NewLabel("S3 Storage")

	// Status text
	lastModifiedLabel := widget.NewLabelWithData(lastModified)

	// Input/output
	downloadPathEntry := widget.NewEntryWithData(downloadPath)

	// Buttons
	connectBtn := widget.NewButton("Connect to S3", func() {
		isWorking.Show()
		isWorking.Start()

		go func() {
			err := reg.S3.Connect(&reg.Config.S3)
			if err != nil {
				isConnected.Set(false)
				connectedError.Set(fmt.Sprint(err))
				isWorking.Hide()
				isWorking.Stop()
				return
			}

			isConnected.Set(true)

			res, err := reg.S3.ListExistingPrefixes(&reg.Config)
			if err != nil {
				connectedError.Set(fmt.Sprint(err))
			}

			err = prefixList.Set(*res)
			if err != nil {
				connectedError.Set(fmt.Sprint(err))
			}

			fyne.Do(func() {
				isWorking.Hide()
				isWorking.Stop()
			})
		}()
	})
	latestBtn := widget.NewButton("Check for latest", func() {
		isWorking.Show()
		isWorking.Start()

		go func() {
			obj, err := reg.S3.GetLatestBackup(&reg.Config)
			if err != nil {
				fmt.Println("Attempted to start S3 client without cfg")
				isWorking.Stop()
				isWorking.Hide()
			}
			fyne.Do(func() {
				lastModifiedLabel.Show()
				lastModified.Set(utils.TimeAgo(obj.LastModified))

				isWorking.Stop()
				isWorking.Hide()
			})
		}()
	})
	downloadBtn := widget.NewButton("Download latest", func() {
		isWorking.Show()
		isWorking.Start()
		go func() {
			latest, err := reg.S3.GetLatestBackup(&reg.Config)
			if err != nil {
				fmt.Println(err)
				isWorking.Stop()
				isWorking.Hide()
			}
			path, err := reg.S3.Download(&reg.Config, latest.Key)
			if err != nil {
				fmt.Println(err)
				isWorking.Stop()
				isWorking.Hide()
			}

			fyne.Do(func() {
				downloadPath.Set(*path)
				downloadPathEntry.Show()
				isWorking.Stop()
				isWorking.Hide()
			})
		}()
	})

	// containers
	labelGroup := container.NewBorder(nil, nil, sectionLabel, isWorking)
	latestGroup := container.NewBorder(nil, nil, lastModifiedLabel, nil, latestBtn)

	// Initial gui states
	lastModifiedLabel.Hide()
	downloadPathEntry.Hide()

	// Watchers
	isConnected.AddListener(binding.NewDataListener(func() {
		status, _ := isConnected.Get()

		fyne.Do(func() {
			if status {
				latestBtn.Enable()
				downloadBtn.Enable()
			} else {
				latestBtn.Disable()
				downloadBtn.Disable()
			}
		})
	}))

	return container.NewVBox(labelGroup, connectBtn, latestGroup, downloadBtn, downloadPathEntry)
}
