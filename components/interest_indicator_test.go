package components

import (
	"slices"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestInterestIndicator_FillsEachHeartByInterestLevel(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at tre billettholdere har meldt veldig, middels og litt interesse på arrangementet.",
		When:  "Når interesseindikatoren vises.",
		Then:  "Så skal hjertene fylles helt, 65 prosent og 45 prosent.",
	})

	// Given
	expectedHeartStyles := []string{"--interest-fill: 100%;", "--interest-fill: 65%;", "--interest-fill: 45%;"}

	interests := []BillettholderInterest{
		{BillettholderID: 1, BillettholderName: "Anna Aas", InterestLevel: models.InterestLevelHigh, IsSelected: true},
		{BillettholderID: 2, BillettholderName: "Bjørn Berg", InterestLevel: models.InterestLevelMedium},
		{BillettholderID: 3, BillettholderName: "Cecilie Carlsen", InterestLevel: models.InterestLevelLow},
	}

	// When
	doc := templtest.Render(t, InterestIndicator(interests))

	// Then
	actualHeartStyles := make([]string, 0)
	doc.Find(".interest-indicator-heart").Each(func(_ int, heart *goquery.Selection) {
		actualHeartStyles = append(actualHeartStyles, heart.AttrOr("style", ""))
	})
	if !slices.Equal(expectedHeartStyles, actualHeartStyles) {
		t.Fatalf("interest heart styles mismatch\nexpected: %v\nactual:   %v", expectedHeartStyles, actualHeartStyles)
	}
}
