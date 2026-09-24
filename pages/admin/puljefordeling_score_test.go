package admin

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	smodel "github.com/Regncon/conorganizer/service/puljefordeling/solver/model"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljeEventBox_ShowsSolverScoreWithCalculation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en deltaker som fordelingen ga plass med bonus for manglende førstevalg og spillederrolle.",
		When:  "Når arrangementskortet rendres.",
		Then:  "Så skal algoritmepoengene vises, med hver del av utregningen i tooltipen.",
	})

	// Given
	expectedTotal := "900"
	expectedLines := []string{
		"Veldig interessert, mangler førstevalg", "+800",
		"Ikke fått førstevalg i 2 puljer", "+40",
		"Spilleder i helgen", "+60",
		"Sum", "900",
	}
	score := smodel.ScoreBreakdown{Score: 5, Band: 800, Misses: 2, MissBonus: 40, DMBump: 60, Total: 900}
	ev := puljefordeling.EmulatedEvent{EventID: "ev1", Title: "Drager", Capacity: 4, AssignedPlayers: []puljefordeling.AssignedPlayer{
		{BillettholderID: 1, Name: "Kari Nordmann", Level: models.InterestLevelHigh, Score: &score},
	}}

	// When
	doc := templtest.Render(t, puljeEventBox(models.PuljeFredagKveld, ev, false, nil))

	// Then
	chip := doc.Find(".pulje-players .pulje-score")
	if got := strings.TrimSpace(chip.Find(".pulje-score-value").Text()); got != expectedTotal {
		t.Fatalf("expected score chip %q, got %q", expectedTotal, got)
	}
	details := chip.Find(".pulje-score-details")
	if id, _ := details.Attr("id"); id == "" || chip.AttrOr("aria-describedby", "") != id {
		t.Fatalf("expected the chip to be described by its tooltip, got id %q", id)
	}
	actualLines := templtest.CollectTexts(doc, ".pulje-score-details dt, .pulje-score-details dd")
	if strings.Join(actualLines, "|") != strings.Join(expectedLines, "|") {
		t.Fatalf("expected calculation %v, got %v", expectedLines, actualLines)
	}
}

func TestPuljeEventBox_PinnedPlayerShowsNoSolverScore(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en manuelt plassert deltaker uten algoritmepoeng.",
		When:  "Når arrangementskortet rendres.",
		Then:  "Så skal det ikke vises noen poengbrikke.",
	})

	// Given
	ev := puljefordeling.EmulatedEvent{EventID: "ev1", Title: "Drager", Capacity: 4, AssignedPlayers: []puljefordeling.AssignedPlayer{
		{BillettholderID: 1, Name: "Kari Nordmann", Level: models.InterestLevelHigh, Manual: true},
	}}

	// When
	doc := templtest.Render(t, puljeEventBox(models.PuljeFredagKveld, ev, false, nil))

	// Then
	if n := doc.Find(".pulje-players .pulje-score").Length(); n != 0 {
		t.Fatalf("expected no score chip for a pinned player, got %d", n)
	}
}
