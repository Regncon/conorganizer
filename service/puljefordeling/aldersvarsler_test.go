package puljefordeling

import (
	"reflect"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestFinnAldersvarsler_SpillerUnder18PaaVoksenarrangementVarsles(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et 18+-arrangement der solveren har satt en spiller under 18, og en GM under 18 leder det.",
		When:  "Når aldersvarsler beregnes for forhåndsvisningen.",
		Then:  "Så gir arrangementet ett varsel som bare nevner spilleren, ikke GM-en.",
	})

	// Given
	expected := []Aldersvarsel{{
		EventID:    "ev18",
		EventTitle: "Voksenspill",
		Spillere: []Aldersvarselspiller{
			{BillettholderID: 2, Navn: "Kari Nordmann"},
		},
	}}
	pulje := EmulatedPulje{PuljeID: models.PuljeFredagKveld, Events: []EmulatedEvent{{
		EventID:  "ev18",
		Title:    "Voksenspill",
		AgeGroup: models.AgeGroupAdultsOnly,
		AssignedPlayers: []AssignedPlayer{
			{BillettholderID: 1, Name: "Ola Voksen", IsOver18: true},
			{BillettholderID: 2, Name: "Kari Nordmann", IsOver18: false},
		},
		AssignedGMs: []AssignedGM{{BillettholderID: 3, Name: "Gunnar GM", IsOver18: false}},
	}}}

	// When
	varsler := FinnAldersvarsler(pulje)

	// Then
	if !reflect.DeepEqual(varsler, expected) {
		t.Errorf("aldersvarsler = %+v\nvil ha        %+v", varsler, expected)
	}
}

func TestFinnAldersvarsler_VoksneOgVanligeArrangementerVarslesIkke(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt voksne på et 18+-arrangement og en spiller under 18 på et arrangement uten aldersgrense.",
		When:  "Når aldersvarsler beregnes for forhåndsvisningen.",
		Then:  "Så gis det ingen aldersvarsler.",
	})

	// Given
	expectedWarnings := 0
	pulje := EmulatedPulje{PuljeID: models.PuljeFredagKveld, Events: []EmulatedEvent{
		{
			EventID:         "ev18",
			Title:           "Voksenspill",
			AgeGroup:        models.AgeGroupAdultsOnly,
			AssignedPlayers: []AssignedPlayer{{BillettholderID: 1, Name: "Ola Voksen", IsOver18: true}},
			AssignedGMs:     []AssignedGM{{BillettholderID: 3, Name: "Gunnar GM", IsOver18: true}},
		},
		{
			EventID:         "evAlle",
			Title:           "For alle",
			AgeGroup:        models.AgeGroupDefault,
			AssignedPlayers: []AssignedPlayer{{BillettholderID: 2, Name: "Kari Nordmann", IsOver18: false}},
		},
	}}

	// When
	varsler := FinnAldersvarsler(pulje)

	// Then
	if len(varsler) != expectedWarnings {
		t.Errorf("antall aldersvarsler = %d, vil ha %d: %+v", len(varsler), expectedWarnings, varsler)
	}
}
