package admin

import (
	"reflect"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljefordelingIndex_DeclaresSignalDefaultsWithKeyBeforeModifier(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt puljefordelingssiden med tildelingsdialogen.",
		When:  "Når siden rendres.",
		Then:  "Så skal signal-standardverdiene skrives som data-signals:navn__ifmissing, slik at Datastar leser navnet.",
	})

	// Given
	expectedAttrs := []struct {
		selector string
		name     string
		value    string
	}{
		{"#puljefordeling-page", "data-signals:assignment-picker-title__ifmissing", "'Administrer deltaker'"},
		{"#puljefordeling-page", "data-signals:assignment-billettholder-id__ifmissing", "0"},
		{"#puljefordeling-page", "data-signals:age-warning-role__ifmissing", "'Player'"},
		{"#puljefordeling-page", "data-signals:_puljefordeling-assign-dialog-open__ifmissing", "false"},
		{"#tildeling-dialog", "data-signals:assignment-role__ifmissing", "'Player'"},
		{"#tildeling-dialog", "data-signals:tildeling-open__ifmissing", "false"},
	}
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_signal_defaults")

	// When
	doc := templtest.Render(t, puljefordelingIndex(db, logger, models.PuljeFredagKveld, nil))

	// Then
	for _, expected := range expectedAttrs {
		actual, ok := doc.Find(expected.selector).Attr(expected.name)
		if !ok || actual != expected.value {
			t.Fatalf("%s %s = %q (present %t), want %q", expected.selector, expected.name, actual, ok, expected.value)
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
