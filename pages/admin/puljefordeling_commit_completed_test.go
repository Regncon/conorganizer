package admin

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
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
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/puljefordeling/FredagKveld/commit", nil))

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

func TestPuljefordelingTabContent_SaveDisabledOnlyWhenCompleted(t *testing.T) {
	for _, status := range []models.PuljeStatus{models.PuljeStatusOpen, models.PuljeStatusLocked, models.PuljeStatusCompleted} {
		t.Run(string(status), func(t *testing.T) {
			// Given
			expectedDisabled := status == models.PuljeStatusCompleted
			db, logger := testutil.CreateTestDBAndLogger(t, "commit_button_"+string(status))
			seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", status, "2026-09-04 18:00")

			// When
			doc := templtest.Render(t, PuljefordelingTabContent(db, logger, models.PuljeFredagKveld, nil))

			// Then
			button := doc.Find("button.pulje-commit")
			if button.Length() != 1 {
				t.Fatalf("expected one save button, got %d", button.Length())
			}
			if _, disabled := button.Attr("disabled"); disabled != expectedDisabled {
				t.Errorf("save disabled = %t, want %t", disabled, expectedDisabled)
			}
		})
	}
}
