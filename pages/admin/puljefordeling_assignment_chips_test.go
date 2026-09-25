package admin

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljeAssignmentInterestRow_ExplainsStatusIconsWithText(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder under 18 som allerede har plass i puljen og har hatt førstevalg.",
		When:  "Når raden i hurtigtildelingen rendres.",
		Then:  "Så skal hvert ikon ha en synlig tekst og en forklaring.",
	})

	// Given
	expected := map[string]string{
		"is-assigned":      "Har plass i puljen",
		"is-under-18":      "Under 18",
		"has-first-choice": "Har hatt førstevalg",
	}
	candidate := puljeAssignmentInterestCandidate{BillettholderID: 1, FirstName: "Ola", LastName: "Nordmann", IsAssignedInPulje: true, HasPreviousFirstChoice: true}

	// When
	doc := templtest.Render(t, puljeAssignmentInterestRow(candidate))

	// Then
	for class, label := range expected {
		chip := doc.Find(".pulje-assignment-chip." + class)
		if got := strings.TrimSpace(chip.Text()); got != label {
			t.Errorf("chip %s: expected label %q, got %q", class, label, got)
		}
		if chip.AttrOr("data-tippy-content", "") == "" || chip.AttrOr("aria-label", "") == "" {
			t.Errorf("chip %s: expected a tooltip and an accessible description", class)
		}
	}
}
