package components

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/components/icons"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestBillettholderInterest_PicksTheHeartIconForEachInterestLevel(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt interessenivåene veldig, middels, litt og ingen interesse.",
		When:  "Når hjertet for hvert nivå velges.",
		Then:  "Så skal hjertet være fullt, fylt til middels, fylt til litt og tomt.",
	})

	// Given
	expectedIcons := []icons.IconType{icons.HeartFilled, icons.HeartFilledMedium, icons.HeartFilledLow, icons.HeartUnfilled}
	levels := []models.InterestLevel{models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow, models.InterestLevelNone}

	// When
	actualIcons := make([]icons.IconType, 0, len(levels))
	for _, level := range levels {
		actualIcons = append(actualIcons, BillettholderInterest{InterestLevel: level}.HeartIcon())
	}

	// Then
	if !slices.Equal(expectedIcons, actualIcons) {
		t.Fatalf("heart icons mismatch\nexpected: %v\nactual:   %v", expectedIcons, actualIcons)
	}
}

func TestInterestIndicator_WritesYouForTheAccountsOwnBillettholder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at kontoens egen billettholder og en annen billettholder har meldt interesse.",
		When:  "Når interesseindikatoren vises.",
		Then:  "Så skal verktøytipset skrive «Du» for kontoens egen billettholder og fornavnet for den andre.",
	})

	// Given
	expectedTooltip := "Du er veldig interessert\nBjørn er interessert"

	interests := []BillettholderInterest{
		{BillettholderID: 1, BillettholderName: "Anna Aas", InterestLevel: models.InterestLevelHigh, IsSelected: true, IsOwn: true},
		{BillettholderID: 2, BillettholderName: "Bjørn Berg", InterestLevel: models.InterestLevelMedium},
	}

	// When
	doc := templtest.Render(t, InterestIndicator(interests))

	// Then
	if actual := doc.Find(".interest-indicator").AttrOr("data-tippy-content", ""); actual != expectedTooltip {
		t.Fatalf("interest tooltip mismatch\nexpected: %q\nactual:   %q", expectedTooltip, actual)
	}
}
