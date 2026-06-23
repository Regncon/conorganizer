package admin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

func setupRerunRouter(db *sql.DB, logger *slog.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Route("/admin/puljefordeling", func(pf chi.Router) {
		SetupPuljefordelingTabRoute(pf, db, &live.Manager{}, logger)
	})
	return r
}

func TestRerun_LockedPuljeReSolves(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_rerun_locked")
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedTabParticipant(t, db, 1, "Anna", "A")
	seedTabInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)
	if _, err := db.Exec(`UPDATE puljer SET status = 'Locked' WHERE id = ?`, string(fredag)); err != nil {
		t.Fatalf("lock pulje: %v", err)
	}

	router := setupRerunRouter(db, logger)
	req := httptest.NewRequest(http.MethodPut, "/admin/puljefordeling/api/FredagKveld/rerun", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND source = 'solver'`,
		string(fredag),
	).Scan(&n); err != nil {
		t.Fatalf("count solver seats: %v", err)
	}
	if n == 0 {
		t.Errorf("expected rerun to commit solver seats, found none")
	}
}

func TestRerun_OpenPuljeConflicts(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_rerun_open")
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")

	router := setupRerunRouter(db, logger)
	req := httptest.NewRequest(http.MethodPut, "/admin/puljefordeling/api/FredagKveld/rerun", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Errorf("expected 409 for open pulje, got %d", rec.Code)
	}
}

func TestRerun_InvalidPuljeBadRequest(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_rerun_invalid")
	router := setupRerunRouter(db, logger)
	req := httptest.NewRequest(http.MethodPut, "/admin/puljefordeling/api/Nope/rerun", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid pulje, got %d", rec.Code)
	}
}
