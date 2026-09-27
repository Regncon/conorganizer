package root

import (
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/program"
)

type puljeAlert struct {
	PuljeName string
	Message   string
	Class     string
}

func collectPuljeAlerts(blocks []program.PuljeBlock) []puljeAlert {
	showPuljeName := len(blocks) > 1
	alerts := make([]puljeAlert, 0, len(blocks))
	for _, block := range blocks {
		alert, ok := puljeAlertFor(block.Pulje)
		if !ok {
			continue
		}
		if showPuljeName {
			alert.PuljeName = block.Pulje.Name
		}
		alerts = append(alerts, alert)
	}
	return alerts
}

func puljeAlertFor(pulje models.PuljeRow) (puljeAlert, bool) {
	switch pulje.Status {
	case models.PuljeStatusOpen:
		if pulje.ClosingWarningActive {
			return puljeAlert{Message: "puljen er i ferd med å bli låst", Class: "is-closing"}, true
		}
	case models.PuljeStatusLocked:
		return puljeAlert{Message: "puljen er låst", Class: "is-locked"}, true
	case models.PuljeStatusCompleted:
		return puljeAlert{Message: "puljen er lukket", Class: "is-completed"}, true
	}
	return puljeAlert{}, false
}
