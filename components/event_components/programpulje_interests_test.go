package event_components

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestProgramPuljeInterests_PreservesOpenStateAcrossLiveUpdates(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the profile program interest dropdown is rendered.",
		When:  "When Datastar morphs it during a live update.",
		Then:  "Then the details open attribute is preserved.",
	})

	// Given
	interests := []Interest{
		{
			EventID:       "interest-event",
			EventName:     "Interest Event",
			InterestLevel: models.InterestLevelHigh,
		},
	}

	// When
	doc := templtest.Render(t, ProgramPuljeInterests(interests))
	collapse := doc.Find(".pulje-interests-collapse")
	actualPreserveAttr, actualPreserveAttrExists := collapse.Attr("data-preserve-attr")

	// Then
	if !actualPreserveAttrExists || actualPreserveAttr != "open" {
		t.Fatalf("interest dropdown preserve attr mismatch\nexpected: %q\nactual:   %q", "open", actualPreserveAttr)
	}
}
