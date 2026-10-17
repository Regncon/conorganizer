package event_components

import (
	"reflect"
	"testing"

	ticketholder "github.com/Regncon/conorganizer/components/ticket_holder"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestInterestContent_DeclaresSavedInterestAsJSONSignals(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder som har valgt Veldig interessert.",
		When:  "Når interesseinnholdet rendres.",
		Then:  "Så skal data-signals inneholde selectedInterestLevel og currentInterestLevelChoice som JSON.",
	})

	// Given
	expected := map[string]any{
		"selectedInterestLevel":      string(models.InterestLevelHigh),
		"currentInterestLevelChoice": ticketholder.InterestChoicePending,
	}

	// When
	doc := templtest.Render(t, InterestContent("drager-og-fangehull", "Drager og Fangehull", InterestNoticeState{CanChooseInterest: true}, models.InterestLevelHigh))

	// Then
	actual := templtest.SignalsObject(t, doc.Find("#interest-content"), "data-signals")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("interest content signals = %v, want %v", actual, expected)
	}
}
