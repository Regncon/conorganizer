package event_components

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestInterestContent_WhenUserCanChoose_ExplainsThereIsNoLimitPerPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a billettholder who can choose an interest.",
		When:  "When the interest picker is rendered.",
		Then:  "Then it says they can choose as many arrangementer as they like, and share interest levels.",
	})

	// Given
	expectedHint := "Du kan melde interesse for så mange arrangementer du vil i hver pulje."

	// When
	doc := templtest.Render(t, InterestContent("e1", "Spill", InterestNoticeState{CanChooseInterest: true}, models.InterestLevelNone))

	// Then
	if got := doc.Find(".interest-limit-hint").Text(); !strings.Contains(got, expectedHint) {
		t.Fatalf("interest hint = %q, want it to contain %q", got, expectedHint)
	}
}

func TestInterestContent_WhenUserCannotChoose_HidesLimitHint(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a billettholder who is already assigned in the pulje.",
		When:  "When the interest content is rendered.",
		Then:  "Then there is no hint about choosing more arrangementer.",
	})

	// Given
	state := InterestNoticeState{ShowAssignedEvent: true}

	// When
	doc := templtest.Render(t, InterestContent("e1", "Spill", state, models.InterestLevelNone))

	// Then
	if doc.Find(".interest-limit-hint").Length() != 0 {
		t.Fatal("expected no interest hint when the billettholder cannot choose")
	}
}
