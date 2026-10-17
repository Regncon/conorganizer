package ticketholder

import (
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestTicketHolderPicker_DeclaresBillettholderSignalWithKeyBeforeModifier(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som er knyttet til en billettholder.",
		When:  "Når billettholdervelgeren rendres.",
		Then:  "Så skal billettHolderId få standardverdien med data-signals:billett-holder-id__ifmissing.",
	})

	// Given
	expectedValue := "''"
	db := createGetTicketHoldersTestDB(t)
	insertGetTicketHoldersBillettholder(t, db, 801, "kari.nordmann@example.com")
	linkGetTicketHoldersBillettholder(t, db, 801)

	// When
	doc := templtest.Render(t, TicketHolderPicker(getTicketHoldersTestUserInfo(), "drager-og-fangehull", db, testutil.NewTestLogger()))

	// Then
	actual, ok := doc.Find(".ticket-holder-container").Attr("data-signals:billett-holder-id__ifmissing")
	if !ok || actual != expectedValue {
		t.Fatalf("billett-holder-id default = %q (present %t), want %q", actual, ok, expectedValue)
	}
	if names := templtest.SignalsAttrsWithModifierBeforeKey(doc); len(names) > 0 {
		t.Fatalf("expected no data-signals__modifier:key attributes, got %v", names)
	}
}
