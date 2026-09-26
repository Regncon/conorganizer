package admin

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestUpdatePuljeStatus_RejectsLockingWithoutClosingWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en åpen pulje uten aktivt stengevarsel.",
		When:  "Når admin prøver å lukke puljefordelingen.",
		Then:  "Så skal endringen avvises fordi steg 1 ikke er gjort.",
	})

	// Given
	expectedStatus := models.PuljeStatusOpen
	db, _ := testutil.CreateTestDBAndLogger(t, "puljefordeling_step_lock_requires_warning")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", expectedStatus, "2026-01-01 18:00")

	// When
	err := updatePuljeStatus(db, models.PuljeFredagKveld, models.PuljeStatusLocked)

	// Then
	if !errors.Is(err, errPuljeStepOrder) {
		t.Fatalf("expected step order error, got %v", err)
	}
	assertPuljeStatus(t, db, models.PuljeFredagKveld, expectedStatus)
}

func TestUpdatePuljeStatus_RejectsPublishingOpenPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en åpen pulje med aktivt stengevarsel.",
		When:  "Når admin prøver å publisere puljefordelingen.",
		Then:  "Så skal endringen avvises fordi steg 2 ikke er gjort.",
	})

	// Given
	expectedStatus := models.PuljeStatusOpen
	db, _ := testutil.CreateTestDBAndLogger(t, "puljefordeling_step_publish_requires_lock")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", expectedStatus, "2026-01-01 18:00")
	testutil.MustExec(t, db, `UPDATE puljer SET closing_warning_active = TRUE WHERE id = ?`, string(models.PuljeFredagKveld))

	// When
	err := updatePuljeStatus(db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if !errors.Is(err, errPuljeStepOrder) {
		t.Fatalf("expected step order error, got %v", err)
	}
	assertPuljeStatus(t, db, models.PuljeFredagKveld, expectedStatus)
}

func TestUpdatePuljeStatus_RejectsReopeningPublishedPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en publisert pulje.",
		When:  "Når admin prøver å åpne puljefordelingen igjen.",
		Then:  "Så skal endringen avvises fordi steg 3 må fjernes først.",
	})

	// Given
	expectedStatus := models.PuljeStatusCompleted
	db, _ := testutil.CreateTestDBAndLogger(t, "puljefordeling_step_reopen_requires_unpublish")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", expectedStatus, "2026-01-01 18:00")

	// When
	err := updatePuljeStatus(db, models.PuljeFredagKveld, models.PuljeStatusOpen)

	// Then
	if !errors.Is(err, errPuljeStepOrder) {
		t.Fatalf("expected step order error, got %v", err)
	}
	assertPuljeStatus(t, db, models.PuljeFredagKveld, expectedStatus)
}

func TestPuljeStatusToggles_OpenWithoutWarningOnlyEnablesFirstStep(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en åpen pulje uten stengevarsel.",
		When:  "Når status-stegene rendres.",
		Then:  "Så skal bare steg 1 kunne krysses av.",
	})

	// Given
	expected := []stepState{{enabled: true}, {}, {}}
	row := models.PuljeRow{ID: models.PuljeFredagKveld, Name: "Fredag Kveld", Status: models.PuljeStatusOpen}

	// When
	doc := templtest.Render(t, puljeStatusToggles(row, false))

	// Then
	assertStepStates(t, doc, expected)
}

func TestPuljeStatusToggles_ActiveWarningEnablesLockStep(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en åpen pulje med aktivt stengevarsel.",
		When:  "Når status-stegene rendres.",
		Then:  "Så skal steg 1 være gjort og steg 2 kunne krysses av, men ikke steg 3.",
	})

	// Given
	expected := []stepState{{checked: true, enabled: true}, {enabled: true}, {}}
	row := models.PuljeRow{ID: models.PuljeFredagKveld, Name: "Fredag Kveld", Status: models.PuljeStatusOpen, ClosingWarningActive: true}

	// When
	doc := templtest.Render(t, puljeStatusToggles(row, false))

	// Then
	assertStepStates(t, doc, expected)
}

func TestPuljeStatusToggles_LockedPuljeEnablesUnlockAndPublish(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lukket pulje (stengevarselet er nullstilt ved lukking).",
		When:  "Når status-stegene rendres.",
		Then:  "Så skal steg 1 vises som gjort og låst, og steg 2 og 3 kunne endres.",
	})

	// Given
	expected := []stepState{{checked: true}, {checked: true, enabled: true}, {enabled: true}}
	row := models.PuljeRow{ID: models.PuljeFredagKveld, Name: "Fredag Kveld", Status: models.PuljeStatusLocked}

	// When
	doc := templtest.Render(t, puljeStatusToggles(row, false))

	// Then
	assertStepStates(t, doc, expected)
}

func TestPuljeStatusToggles_CompletedPuljeOnlyEnablesUnpublish(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en publisert pulje.",
		When:  "Når status-stegene rendres.",
		Then:  "Så skal alle steg være gjort, og bare steg 3 kunne fjernes.",
	})

	// Given
	expected := []stepState{{checked: true}, {checked: true}, {checked: true, enabled: true}}
	row := models.PuljeRow{ID: models.PuljeFredagKveld, Name: "Fredag Kveld", Status: models.PuljeStatusCompleted}

	// When
	doc := templtest.Render(t, puljeStatusToggles(row, false))

	// Then
	assertStepStates(t, doc, expected)
}

type stepState struct {
	checked bool
	enabled bool
}

func assertStepStates(t *testing.T, doc *goquery.Document, expected []stepState) {
	t.Helper()
	inputs := doc.Find(".pulje-status-control input[type=checkbox]")
	if inputs.Length() != len(expected) {
		t.Fatalf("expected %d steps, got %d", len(expected), inputs.Length())
	}
	inputs.Each(func(i int, input *goquery.Selection) {
		actual := stepState{checked: input.Is("[checked]"), enabled: !input.Is("[disabled]")}
		if actual != expected[i] {
			t.Errorf("step %d: expected %+v, got %+v", i+1, expected[i], actual)
		}
	})
}

func assertPuljeStatus(t *testing.T, db *sql.DB, puljeID models.Pulje, expected models.PuljeStatus) {
	t.Helper()
	var actual models.PuljeStatus
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(puljeID)).Scan(&actual); err != nil {
		t.Fatalf("load pulje status: %v", err)
	}
	if actual != expected {
		t.Fatalf("expected pulje status %s, got %s", expected, actual)
	}
}
