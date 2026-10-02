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
	doc := templtest.Render(t, ProgramPuljeInterests(interests, "2026-10-03"))
	collapse := doc.Find(".pulje-interests-collapse")
	actualPreserveAttr, actualPreserveAttrExists := collapse.Attr("data-preserve-attr")

	// Then
	if !actualPreserveAttrExists || actualPreserveAttr != "open" {
		t.Fatalf("interest dropdown preserve attr mismatch\nexpected: %q\nactual:   %q", "open", actualPreserveAttr)
	}
}

func TestProgramPuljeInterests_LinksToProgramDate(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the profile program interest section is rendered for a pulje day.",
		When:  "When the user chooses to see arrangements.",
		Then:  "Then the link selects that program date on the root page.",
	})

	// Given
	date := "2026-10-03"

	// When
	doc := templtest.Render(t, ProgramPuljeInterests(nil, date))
	actualHref, exists := doc.Find(".pulje-interests-description a").Attr("href")

	// Then
	if !exists {
		t.Fatal("expected see-arrangements link to have an href")
	}
	if want := "/?date=2026-10-03"; actualHref != want {
		t.Fatalf("see-arrangements href = %q, want %q", actualHref, want)
	}
}
