package puljefordeling

import (
	"sort"

	"github.com/Regncon/conorganizer/models"
)

type Aldersvarseldeltaker struct {
	BillettholderID int
	Navn            string
	Role            models.EventPlayerRole
}

// Aldersvarsel names everyone under 18 who plays in or runs one 18+ event.
type Aldersvarsel struct {
	EventID    string
	EventTitle string
	Deltakere  []Aldersvarseldeltaker
}

// FinnAldersvarsler finds participants under 18 on 18+ events in one displayed
// pulje. The solver does not enforce the age limit (the board handles it
// manually), so the warning covers solver seats as well as manual pins and GMs.
func FinnAldersvarsler(pulje EmulatedPulje) []Aldersvarsel {
	var varsler []Aldersvarsel
	for _, event := range pulje.Events {
		if event.AgeGroup != models.AgeGroupAdultsOnly {
			continue
		}
		varsel := Aldersvarsel{EventID: event.EventID, EventTitle: event.Title}
		for _, gm := range event.AssignedGMs {
			if !gm.IsOver18 {
				varsel.Deltakere = append(varsel.Deltakere, Aldersvarseldeltaker{BillettholderID: gm.BillettholderID, Navn: gm.Name, Role: models.EventPlayerRoleGM})
			}
		}
		for _, player := range event.AssignedPlayers {
			if player.BillettholderID > 0 && !player.IsOver18 {
				varsel.Deltakere = append(varsel.Deltakere, Aldersvarseldeltaker{BillettholderID: player.BillettholderID, Navn: player.Name, Role: models.EventPlayerRolePlayer})
			}
		}
		if len(varsel.Deltakere) == 0 {
			continue
		}
		sort.SliceStable(varsel.Deltakere, func(i, j int) bool {
			return varsel.Deltakere[i].Navn < varsel.Deltakere[j].Navn
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
