package billettholderadmin

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/Regncon/conorganizer/models"
	billettholderService "github.com/Regncon/conorganizer/service/billettholder"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestBillettholderAdminPage_DeclaresOpenDialogSignalWithKeyBeforeModifier(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt billettholder-adminsiden.",
		When:  "Når siden rendres.",
		Then:  "Så skal billettholderInterestOpenDialogId få standardverdien med data-signals:navn__ifmissing.",
	})

	// Given
	expectedValue := "''"
	db, logger := testutil.CreateTestDBAndLogger(t, "billettholder_admin_page_signals")

	// When
	doc := templtest.Render(t, BillettholderAdminPage(db, logger))

	// Then
	actual, ok := doc.Find("#billettholder-admin-container").Attr("data-signals:billettholder-interest-open-dialog-id__ifmissing")
	if !ok || actual != expectedValue {
		t.Fatalf("open dialog default = %q (present %t), want %q", actual, ok, expectedValue)
	}
	if names := templtest.SignalsAttrsWithModifierBeforeKey(doc); len(names) > 0 {
		t.Fatalf("expected no data-signals__modifier:key attributes, got %v", names)
	}
}

func TestBillettholderInterestDialog_DeclaresFirstPuljeAsActivePuljeDefault(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder med interesser i to puljer.",
		When:  "Når interessedialogen rendres.",
		Then:  "Så skal den første puljen være standardverdi for billettholderInterestActivePulje i data-signals__ifmissing.",
	})

	// Given
	expected := map[string]models.Pulje{"billettholderInterestActivePulje": models.PuljeLordagMorgen}
	billettholder := models.Billettholder{ID: 42, FirstName: "Kari", LastName: "Nordmann"}
	sections := []billettholderInterestPuljeSection{
		{PuljeID: models.PuljeLordagMorgen, Name: "Lørdag morgen"},
		{PuljeID: models.PuljeLordagKveld, Name: "Lørdag kveld"},
	}

	// When
	doc := templtest.Render(t, billettholderInterestDialog(billettholder, sections, "billettholder-interests-42", 0, 0))

	// Then
	raw, ok := doc.Find("#billettholder-interests-42").Attr("data-signals__ifmissing")
	if !ok {
		t.Fatal("expected the dialog to have a data-signals__ifmissing attribute")
	}
	var actual map[string]models.Pulje
	if err := json.Unmarshal([]byte(raw), &actual); err != nil {
		t.Fatalf("data-signals__ifmissing is not a JSON object: %q: %v", raw, err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("active pulje default = %v, want %v", actual, expected)
	}
	if names := templtest.SignalsAttrsWithModifierBeforeKey(doc); len(names) > 0 {
		t.Fatalf("expected no data-signals__modifier:key attributes, got %v", names)
	}
}

func TestBillettholderAdminOverview_DeclaresActiveFiltersAsJSONSignals(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at filteret Uten førstevalg er på og GM/DM er av.",
		When:  "Når oversikten rendres.",
		Then:  "Så skal data-signals inneholde filterWithoutFirstChoice og filterGmOrDm som JSON.",
	})

	// Given
	expected := map[string]any{"filterWithoutFirstChoice": true, "filterGmOrDm": false}
	db, logger := testutil.CreateTestDBAndLogger(t, "billettholder_admin_overview_signals")
	filters := billettholderService.BillettholderFilters{WithoutFirstChoice: true}

	// When
	doc := templtest.Render(t, billettholderAdminOverview(db, logger, filters))

	// Then
	actual := templtest.SignalsObject(t, doc.Find(".billettholder-admin-overview"), "data-signals")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("overview filter signals = %v, want %v", actual, expected)
	}
}
