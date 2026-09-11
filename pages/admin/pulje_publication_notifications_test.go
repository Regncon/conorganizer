package admin

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
)

func TestPuljeStatusRoute_PublishingQueuesParticipantNotification(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en deltaker med en plass og aktiverte varsler.",
		When:  "Når administrator publiserer puljen.",
		Then:  "Så publiseres resultatet og ett varsel legges i kø.",
	})
	// Given
	expectedJobs := 1
	db, logger := testutil.CreateTestDBAndLogger(t, "publish_notification_route")
	seedPublicationNotification(t, db)
	router := chi.NewRouter()
	puljefordelingStatusRoute(router, db, &live.Manager{}, logger)

	// When
	rec := postPuljePublicationStatus(router, "Completed")

	// Then
	if rec.Code != http.StatusNoContent {
		t.Errorf("publish returned %d: %s", rec.Code, rec.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs`); got != expectedJobs {
		t.Errorf("queued notifications = %d, want %d", got, expectedJobs)
	}
	var status string
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, models.PuljeFredagKveld).Scan(&status); err != nil {
		t.Fatal(err)
	}
	if status != string(models.PuljeStatusCompleted) {
		t.Errorf("pulje status = %q", status)
	}
}

func TestPuljeStatusRoute_RepeatedPublicationDoesNotQueueDuplicate(t *testing.T) {
	// Given
	expectedJobs := 1
	db, logger := testutil.CreateTestDBAndLogger(t, "repeat_publish_notification_route")
	seedPublicationNotification(t, db)
	router := chi.NewRouter()
	puljefordelingStatusRoute(router, db, &live.Manager{}, logger)
	if rec := postPuljePublicationStatus(router, "Completed"); rec.Code != http.StatusNoContent {
		t.Errorf("first publish: %d %s", rec.Code, rec.Body.String())
	}

	// When
	rec := postPuljePublicationStatus(router, "Completed")

	// Then
	if rec.Code != http.StatusNoContent {
		t.Errorf("repeat publish: %d %s", rec.Code, rec.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs`); got != expectedJobs {
		t.Errorf("queued notifications = %d, want %d", got, expectedJobs)
	}
}

func TestProgramPublishingRoute_PublishingPreviouslyCompletedPuljeQueuesNotification(t *testing.T) {
	// Given
	expectedJobs := 1
	db, logger := testutil.CreateTestDBAndLogger(t, "publish_completed_program_route")
	seedPublicationNotification(t, db)
	testutil.MustExec(t, db, `UPDATE program_publishing_state SET is_published = 0`)
	testutil.MustExec(t, db, `UPDATE puljer SET status = 'Completed'`)
	router := chi.NewRouter()
	programPublishingRoute(router, db, &live.Manager{}, logger)
	req := httptest.NewRequest(http.MethodPut, "/api/program-publishing", strings.NewReader(`{"programPublished":true}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	// When
	router.ServeHTTP(rec, req)

	// Then
	if rec.Code != http.StatusNoContent {
		t.Errorf("program publish returned %d: %s", rec.Code, rec.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM web_push_jobs`); got != expectedJobs {
		t.Errorf("queued notifications = %d, want %d", got, expectedJobs)
	}
}

func seedPublicationNotification(t *testing.T, db *sql.DB) {
	t.Helper()
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusLocked, "2026-09-04 18:00")
	seedTabEventWithInterest(t, db, "evA", "Alpha", models.PuljeFredagKveld)
	testutil.MustExec(t, db, `UPDATE events SET status = 'Annonsert' WHERE id = 'evA'`)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET is_published = 1 WHERE event_id = 'evA'`)
	testutil.MustExec(t, db, `INSERT INTO program_publishing_state (id, is_published) VALUES (1, 1) ON CONFLICT(id) DO UPDATE SET is_published = 1`)
	testutil.MustExec(t, db, `INSERT INTO users (id, external_id, email) VALUES (1, 'publication-user', 'publication@example.test')`)
	testutil.MustExec(t, db, `INSERT INTO relation_billettholdere_users (billettholder_id, user_id) VALUES (1, 1)`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES ('evA', ?, 1, 'Player', 'solver')`, models.PuljeFredagKveld)
	testutil.MustExec(t, db, `INSERT INTO web_push_subscriptions (user_id, endpoint, p256dh, auth) VALUES (1, 'https://fcm.googleapis.com/fcm/send/test-route', 'fixture-key', 'fixture-auth')`)
}

func postPuljePublicationStatus(router http.Handler, status string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, "/api/puljer/FredagKveld/status", strings.NewReader(`{"puljeStatus":"`+status+`"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}
