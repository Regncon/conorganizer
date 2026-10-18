package admin

import (
	"reflect"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljefordelingIndex_DeclaresSignalDefaultsAsIfMissingJSON(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt puljefordelingssiden med tildelingsdialogen.",
		When:  "Når siden rendres.",
		Then:  "Så skal signal-standardverdiene ligge som JSON i data-signals__ifmissing, slik at Datastar leser navnene.",
	})

	// Given
	expectedSignals := []struct {
		selector string
		name     string
		value    any
	}{
		{"#puljefordeling-page", "assignmentPickerTitle", "Administrer deltaker"},
		{"#puljefordeling-page", "assignmentBillettholderId", float64(0)},
		{"#puljefordeling-page", "ageWarningRole", string(models.EventPlayerRolePlayer)},
		{"#puljefordeling-page", "_puljefordelingAssignDialogOpen", false},
		{"#tildeling-dialog", "assignmentRole", string(models.EventPlayerRolePlayer)},
		{"#tildeling-dialog", "tildelingOpen", false},
	}
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_signal_defaults")

	// When
	doc := templtest.Render(t, puljefordelingIndex(db, logger, models.PuljeFredagKveld, nil))

	// Then
	for _, expected := range expectedSignals {
		signals := templtest.SignalsObject(t, doc.Find(expected.selector), "data-signals__ifmissing")
		if actual, ok := signals[expected.name]; !ok || actual != expected.value {
			t.Fatalf("%s %s = %v (present %t), want %v", expected.selector, expected.name, actual, ok, expected.value)
		}
	}
	if names := templtest.SignalsAttrsWithModifierBeforeKey(doc); len(names) > 0 {
		t.Fatalf("expected no data-signals__modifier:key attributes, got %v", names)
	}
}

func TestPuljefordelingTabContent_DeclaresPuljeStatusAsJSONSignal(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en åpen pulje.",
		When:  "Når puljefordeling-fanen rendres.",
		Then:  "Så skal data-signals inneholde puljeStatus som JSON.",
	})

	// Given
	expected := map[string]any{"puljeStatus": string(models.PuljeStatusOpen)}
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_tab_status_signal")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")

	// When
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld, nil))

	// Then
	actual := templtest.SignalsObject(t, doc.Find("#puljefordeling-tab"), "data-signals")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("pulje status signals = %v, want %v", actual, expected)
	}
}

func TestPubliserProgram_DeclaresPublishedStateAsJSONSignal(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programmet er publisert.",
		When:  "Når publiseringskontrollen rendres.",
		Then:  "Så skal data-signals inneholde programPublished som JSON.",
	})

	// Given
	expected := map[string]any{"programPublished": true}
	db := testutil.CreateTestDB(t, "publiser_program_signal")
	testutil.MustExec(t, db, `INSERT OR REPLACE INTO program_publishing_state (id, is_published) VALUES (1, 1)`)

	// When
	doc := templtest.Render(t, publiserProgram(db))

	// Then
	actual := templtest.SignalsObject(t, doc.Find(".publiser-program"), "data-signals")
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("program published signals = %v, want %v", actual, expected)
	}
}
