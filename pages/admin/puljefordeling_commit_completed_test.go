package admin

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestPuljefordelingCommitRoute_RejectsCompletedPulje(t *testing.T) {
	// Given
	expectedStatus := http.StatusConflict
	db, logger := testutil.CreateTestDBAndLogger(t, "commit_completed_route")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusCompleted, "2026-09-04 18:00")
	seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	// When
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/puljefordeling/FredagKveld/commit", strings.NewReader(`{"saveConfirmation":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(rec, req)

	// Then
	if rec.Code != expectedStatus {
		t.Errorf("got status %d, want %d: %s", rec.Code, expectedStatus, rec.Body.String())
	}
	var seats int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players`).Scan(&seats); err != nil {
		t.Fatal(err)
	}
	if seats != 0 {
		t.Errorf("completed pulje acquired %d new seats", seats)
	}
}

func TestPuljefordelingTabContent_SaveEnabledOnlyWithUnsavedChanges(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt puljer som er åpne, lukket, publisert eller nettopp lagret.",
		When:  "Når puljefordeling-fanen rendres.",
		Then:  "Så skal «Lagre fordeling» bare være aktiv når det finnes ulagrede endringer.",
	})

	cases := []struct {
		name            string
		status          models.PuljeStatus
		saved           bool
		expectedEnabled bool
	}{
		{"Open with unsaved changes", models.PuljeStatusOpen, false, true},
		{"Locked with unsaved changes", models.PuljeStatusLocked, false, true},
		{"Open and saved", models.PuljeStatusOpen, true, false},
		{"Completed", models.PuljeStatusCompleted, false, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Given
			db, logger := testutil.CreateTestDBAndLogger(t, "commit_button_"+strings.ReplaceAll(c.name, " ", "_"))
			seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", c.status, "2026-09-04 18:00")
			seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)
			if c.saved {
				if err := puljefordeling.CommitDistribution(db, models.PuljeFredagKveld); err != nil {
					t.Fatalf("CommitDistribution: %v", err)
				}
			}

			// When
			doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld, nil))

			// Then
			button := doc.Find("button.pulje-commit")
			if button.Length() != 1 {
				t.Fatalf("expected one save button, got %d", button.Length())
			}
			if _, disabled := button.Attr("disabled"); disabled == c.expectedEnabled {
				t.Errorf("save enabled = %t, want %t", !disabled, c.expectedEnabled)
			}
		})
	}
}

func TestPuljefordelingSavePreview_SaysItIsTheFirstSave(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der fordelingen aldri er lagret.",
		When:  "Når admin trykker «Lagre fordeling».",
		Then:  "Så skal dialogen si at dette er første lagring, og bekreftelsen lagre plassene.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "commit_first_save")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-09-04 18:00")
	seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	// When
	preview := postAssignmentSignals(t, router, http.MethodPost, "/api/puljefordeling/FredagKveld/commit/preview", 0, "", "", "")

	// Then
	for _, part := range []string{"Er du sikker?", "Første lagring av fordelingen", "1 deltaker får lagret plass"} {
		if !strings.Contains(preview.Body.String(), part) {
			t.Errorf("expected the dialog to contain %q: %s", part, preview.Body.String())
		}
	}
	confirmed := postAssignmentSignals(t, router, http.MethodPost, "/api/puljefordeling/FredagKveld/commit", 0, "", "", `,"saveConfirmation":"`+confirmationFromResponse(t, preview)+`"`)
	if confirmed.Code != http.StatusNoContent {
		t.Fatalf("expected the confirmed save to succeed, got %d: %s", confirmed.Code, confirmed.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE source = 'solver'`); got != 1 {
		t.Fatalf("expected one saved solver seat, got %d", got)
	}
}

func TestPuljefordelingSave_StaleConfirmationShowsCurrentChanges(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje med ulagrede endringer.",
		When:  "Når lagring bekreftes med en utdatert bekreftelse.",
		Then:  "Så skal ingenting lagres, og dialogen vises på nytt med dagens endringer.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "commit_stale")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-09-04 18:00")
	seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	// When
	rec := postAssignmentSignals(t, router, http.MethodPost, "/api/puljefordeling/FredagKveld/commit", 0, "", "", `,"saveConfirmation":"utdatert"`)

	// Then
	if !strings.Contains(rec.Body.String(), "Første lagring av fordelingen") {
		t.Fatalf("expected the dialog again, got %d: %s", rec.Code, rec.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players`); got != 0 {
		t.Fatalf("stale confirmation saved %d seats", got)
	}
}

func TestPuljeStatusRoute_PublishingRequiresSavedDistribution(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lukket pulje med ulagrede endringer i fordelingen.",
		When:  "Når admin prøver å publisere puljen.",
		Then:  "Så skal publiseringen avvises til fordelingen er lagret.",
	})

	// Given
	expectedStatus := models.PuljeStatusLocked
	db, _ := testutil.CreateTestDBAndLogger(t, "publish_requires_save")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", expectedStatus, "2026-09-04 18:00")
	seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)

	// When
	err := updatePuljeStatus(db, models.PuljeFredagKveld, models.PuljeStatusCompleted)

	// Then
	if !errors.Is(err, errPuljeUnsaved) {
		t.Fatalf("expected publishing to require a saved distribution, got %v", err)
	}
	assertPuljeStatus(t, db, models.PuljeFredagKveld, expectedStatus)
	if err := puljefordeling.CommitDistribution(db, models.PuljeFredagKveld); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}
	if err := updatePuljeStatus(db, models.PuljeFredagKveld, models.PuljeStatusCompleted); err != nil {
		t.Fatalf("expected publishing to succeed once saved: %v", err)
	}
}
