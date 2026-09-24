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

func TestPuljeEventBox_ParticipantInterestCanBeChanged(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en manuelt plassert deltaker uten interesse på arrangementet.",
		When:  "Når arrangementskortet rendres i en åpen pulje.",
		Then:  "Så skal interessefjesset åpne en meny som setter interessen for akkurat dette arrangementet.",
	})

	// Given
	expectedLevels := []string{"Veldig interessert", "Middels interessert", "Litt interessert", "Ikkje interessert"}
	ev := puljefordeling.EmulatedEvent{EventID: "ev1", Title: "Drager", Capacity: 4, AssignedPlayers: []puljefordeling.AssignedPlayer{
		{BillettholderID: 7, Name: "Kari Nordmann", Manual: true},
	}}

	// When
	doc := templtest.Render(t, puljeEventBox(models.PuljeFredagKveld, ev, false, nil))

	// Then
	trigger := doc.Find(".pulje-players button.pulje-interest-trigger")
	menuID := trigger.AttrOr("popovertarget", "")
	menu := doc.Find("#" + menuID + "[popover]")
	if menuID == "" || menu.Length() != 1 {
		t.Fatalf("expected the interest face to open a popover menu, got target %q", menuID)
	}
	actualLevels := templtest.CollectTexts(doc, "#"+menuID+" button")
	if len(actualLevels) != len(expectedLevels) {
		t.Fatalf("expected interest options %v, got %v", expectedLevels, actualLevels)
	}
	for i, expected := range expectedLevels {
		if !strings.Contains(actualLevels[i], expected) {
			t.Fatalf("expected option %d to be %q, got %q", i+1, expected, actualLevels[i])
		}
	}
	action := menu.Find("button").First().AttrOr("data-on:click", "")
	for _, part := range []string{"$assignmentBillettholderId = 7", `$assignmentEventId = "ev1"`, `$assignmentInterestLevel = "Veldig interessert"`, "@put('/admin/api/puljefordeling/interest')"} {
		if !strings.Contains(action, part) {
			t.Fatalf("expected interest action to contain %q, got %q", part, action)
		}
	}
}

func TestPuljeEventBox_PublishedParticipantInterestIsReadOnly(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en deltaker i en publisert pulje.",
		When:  "Når arrangementskortet rendres.",
		Then:  "Så skal interessen bare vises, uten meny for å endre den.",
	})

	// Given
	ev := puljefordeling.EmulatedEvent{EventID: "ev1", Title: "Drager", Capacity: 4, AssignedPlayers: []puljefordeling.AssignedPlayer{
		{BillettholderID: 7, Name: "Kari Nordmann", Level: models.InterestLevelHigh},
	}}

	// When
	doc := templtest.Render(t, puljeEventBox(models.PuljeFredagKveld, ev, true, nil))

	// Then
	if n := doc.Find(".pulje-interest-trigger, [popover]").Length(); n != 0 {
		t.Fatalf("expected no interest menu in a published pulje, got %d elements", n)
	}
	if got := strings.TrimSpace(doc.Find(".pulje-players .pulje-emoji").Text()); got != models.InterestLevelHigh.Emoji() {
		t.Fatalf("expected the interest face to still show, got %q", got)
	}
}
