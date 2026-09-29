package puljefordeling

import (
	"sort"

	"github.com/Regncon/conorganizer/models"
)

type Aldersvarselspiller struct {
	BillettholderID int
	Navn            string
}

// Aldersvarsel names every player under 18 seated on one 18+ event.
type Aldersvarsel struct {
	EventID    string
	EventTitle string
	Spillere   []Aldersvarselspiller
}

// FinnAldersvarsler finds players under 18 on 18+ events in one displayed
// pulje. The solver does not enforce the age limit (the board handles it
// manually), so the warning covers solver seats as well as manual pins.
// GMs are not included.
func FinnAldersvarsler(pulje EmulatedPulje) []Aldersvarsel {
	var varsler []Aldersvarsel
	for _, event := range pulje.Events {
		if event.AgeGroup != models.AgeGroupAdultsOnly {
			continue
		}
		varsel := Aldersvarsel{EventID: event.EventID, EventTitle: event.Title}
		for _, player := range event.AssignedPlayers {
			if player.BillettholderID > 0 && !player.IsOver18 {
				varsel.Spillere = append(varsel.Spillere, Aldersvarselspiller{BillettholderID: player.BillettholderID, Navn: player.Name})
			}
		}
		if len(varsel.Spillere) == 0 {
			continue
		}
		sort.SliceStable(varsel.Spillere, func(i, j int) bool {
			return varsel.Spillere[i].Navn < varsel.Spillere[j].Navn
		})
		varsler = append(varsler, varsel)
	}
	sort.Slice(varsler, func(i, j int) bool {
		if varsler[i].EventTitle != varsler[j].EventTitle {
			return varsler[i].EventTitle < varsler[j].EventTitle
		}
		return varsler[i].EventID < varsler[j].EventID
	})
	return varsler
}
