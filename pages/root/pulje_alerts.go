package root

import (
	"time"

	"github.com/Regncon/conorganizer/components/icons"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/program"
)

const puljeAlertHideAfterStart = time.Hour

type puljeAlert struct {
	PuljeName string
	Message   string
	Class     string
	Icon      icons.IconType
	LinkHref  string
}

func collectPuljeAlerts(blocks []program.PuljeBlock, now time.Time) []puljeAlert {
	showPuljeName := len(blocks) > 1
	alerts := make([]puljeAlert, 0, len(blocks))
	for _, block := range blocks {
		if now.After(block.Pulje.StartAt.TimeOrZero().Add(puljeAlertHideAfterStart)) {
			continue
		}
		alert, ok := puljeAlertFor(block.Pulje)
		if !ok {
			continue
		}
		if showPuljeName {
			alert.PuljeName = block.Pulje.Name
		}
		if len(block.RaffleEvents()) == 0 {
			alert.LinkHref = ""
		}
		alerts = append(alerts, alert)
	}
	return alerts
}

func puljeAlertFor(pulje models.PuljeRow) (puljeAlert, bool) {
	switch pulje.Status {
	case models.PuljeStatusOpen:
		if pulje.ClosingWarningActive {
			return puljeAlert{
				Message:  "Interessevalget stenger snart! Hvis du vil endre valgene dine for kommende pulje, gjør det nå.",
				Class:    "is-closing",
				Icon:     icons.WarningOutline,
				LinkHref: "#" + puljeAnchorID(pulje.ID),
			}, true
		}
	case models.PuljeStatusLocked:
		return puljeAlert{
			Message: "Interessevalg for kommende pulje er nå låst og kan ikke endres. Vi jobber med å fordele spillere og publiserer resultatet snart!",
			Class:   "is-locked",
			Icon:    icons.ClockLock,
		}, true
	case models.PuljeStatusCompleted:
		return puljeAlert{
			Message: "Puljefordelingen er klar! Se hva du fikk på profilen din.",
			Class:   "is-completed",
			Icon:    icons.ProgressComplete,
		}, true
	}
	return puljeAlert{}, false
}

func puljeAnchorID(puljeID models.Pulje) string {
	return "pulje-" + string(puljeID)
}
