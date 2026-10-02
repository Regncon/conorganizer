package event_components

import (
	"strings"
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
	doc := templtest.Render(t, ProgramPuljeInterests(interests, "2026-10-03", models.PuljeFredagKveld))
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
	doc := templtest.Render(t, ProgramPuljeInterests(nil, date, models.PuljeFredagKveld))
	actualHref, exists := doc.Find(".pulje-interests-description a").Attr("href")

	// Then
	if !exists {
		t.Fatal("expected see-arrangements link to have an href")
	}
	if want := "/?date=2026-10-03#pulje-FredagKveld"; actualHref != want {
		t.Fatalf("see-arrangements href = %q, want %q", actualHref, want)
	}
}

func TestProgramPuljeInterests_SaysHowManyAndInvitesToFindMore(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a billettholder with three interests in a pulje.",
		When:  "When the profile program interest section is rendered.",
		Then:  "Then it says they only have three and invites them to find more arrangementer.",
	})

	// Given
	expectedSummary := "Du har bare meldt interesse for 3 arrangement(er) i denne puljen."
	expectedLink := "Finn flere arrangementer"
	interests := []Interest{
		{EventID: "a", EventName: "A", InterestLevel: models.InterestLevelHigh},
		{EventID: "b", EventName: "B", InterestLevel: models.InterestLevelHigh},
		{EventID: "c", EventName: "C", InterestLevel: models.InterestLevelLow},
	}

	// When
	doc := templtest.Render(t, ProgramPuljeInterests(interests, "2026-10-03", models.PuljeFredagKveld))

	// Then
	if got := doc.Find(".pulje-interests-description p").Text(); !strings.Contains(got, expectedSummary) {
		t.Fatalf("summary = %q, want it to contain %q", got, expectedSummary)
	}
	if got := strings.TrimSpace(doc.Find(".pulje-interests-description a").Text()); got != expectedLink {
		t.Fatalf("link text = %q, want %q", got, expectedLink)
	}
}

func TestProgramPuljeInterests_WhenMoreThanThree_DoesNotSayOnly(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a billettholder with four interests in a pulje.",
		When:  "When the profile program interest section is rendered.",
		Then:  "Then it says how many they have without calling it only.",
	})

	// Given
	expectedSummary := "Du har meldt interesse for 4 arrangement(er) i denne puljen."
	interests := []Interest{
		{EventID: "a", EventName: "A", InterestLevel: models.InterestLevelHigh},
		{EventID: "b", EventName: "B", InterestLevel: models.InterestLevelHigh},
		{EventID: "c", EventName: "C", InterestLevel: models.InterestLevelMedium},
		{EventID: "d", EventName: "D", InterestLevel: models.InterestLevelLow},
	}

	// When
	doc := templtest.Render(t, ProgramPuljeInterests(interests, "2026-10-03", models.PuljeFredagKveld))

	// Then
	got := doc.Find(".pulje-interests-description p").Text()
	if !strings.Contains(got, expectedSummary) {
		t.Fatalf("summary = %q, want it to contain %q", got, expectedSummary)
	}
	if strings.Contains(got, "bare") {
		t.Fatalf("summary = %q, should not say bare with more than three interests", got)
	}
}
