package views

import (
	"backup-migration/services"
	"context"
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type action struct {
	Title string
	Do    func()
}

func NewRoot(ctx context.Context, reg *services.Registry, w fyne.Window) fyne.CanvasObject {
	console, consoleObj := NewConsoleWidget(1000)

	// Wire a UI logger that writes into the console widget.
	reg.UILogger = slog.New(slog.NewTextHandler(console, nil))
	reg.ConsoleWriter = console

	actions := make([]action, 0, 8)

	// Always-available actions
	actions = append(actions,
		action{"Say hello", func() { reg.UILogger.Info("Hello from action") }},
		action{"Clear console", func() { console.Clear() }},
	)

	// S3 actions
	if reg.S3 != nil {
		actions = append(actions,
			action{"S3: List latest snapshots (10)", func() {
				go func() {
					snaps, err := reg.S3.Browse(context.Background(), "", 10)
					if err != nil {
						reg.UILogger.Error("browse snapshots failed", "error", err.Error())
						return
					}
					fmt.Fprintln(reg.ConsoleWriter, "GEN\tSNAP#\tSIZE\tLASTMOD\tKEY")
					for _, s := range snaps {
						fmt.Fprintf(reg.ConsoleWriter, "%s\t%d\t%d\t%s\t%s\n",
							s.Generation, s.SnapshotNum, s.Size, s.LastModified.Format("2006-01-02 15:04:05"), s.Key)
					}
				}()
			}},
		)
	} else {
		actions = append(actions,
			action{"S3 not configured — open .env", func() {
				fmt.Fprintln(reg.ConsoleWriter, "S3 not configured. Ensure .env has AWS_* and BUCKET_NAME, then restart.")
			}},
		)
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
