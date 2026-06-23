package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

func mustQueryRow(t *testing.T, db *sql.DB, query string, args ...any) *sql.Row {
	t.Helper()
	return db.QueryRow(query, args...)
}

func putPuljeStatus(t *testing.T, router http.Handler, pulje models.Pulje, status models.PuljeStatus) {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"puljeStatus": string(status)})
	req := httptest.NewRequest(http.MethodPut, "/api/puljer/"+string(pulje)+"/status", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
}

func TestPuljeStatusHandler_LockCommitsUnlockReverts(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_admin_lock")
	router := chi.NewRouter()
	puljefordelingStatusRoute(router, db, &live.Manager{}, logger)

	const fredag = models.PuljeFredagKveld
	testutil.MustExec(t, db, `INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, 'F', 'Open', ?, ?)`,
		string(fredag), "2026-09-04T18:00:00Z", "2026-09-04T22:00:00Z")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA', ?, 1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id) VALUES (1,'Anna','A',0,'',0,1)`)
	testutil.MustExec(t, db, `INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level) VALUES (1,'evA',?, ?)`,
		string(fredag), string(models.InterestLevelHigh))

	// Lock → commit writes a solver seat and status becomes Locked.
	putPuljeStatus(t, router, fredag, models.PuljeStatusLocked)

	var solverCount int
	mustQueryRow(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE pulje_id=? AND source='solver'`, string(fredag)).Scan(&solverCount)
	if solverCount != 1 {
		t.Errorf("after lock: want 1 solver seat, got %d", solverCount)
	}
	var status string
	mustQueryRow(t, db, `SELECT status FROM puljer WHERE id=?`, string(fredag)).Scan(&status)
	if status != string(models.PuljeStatusLocked) {
		t.Errorf("after lock: want status Locked, got %s", status)
	}

	// Unlock → revert removes solver seats and status becomes Open.
	putPuljeStatus(t, router, fredag, models.PuljeStatusOpen)

	mustQueryRow(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE pulje_id=? AND source='solver'`, string(fredag)).Scan(&solverCount)
	if solverCount != 0 {
		t.Errorf("after unlock: want 0 solver seats, got %d", solverCount)
	}
	mustQueryRow(t, db, `SELECT status FROM puljer WHERE id=?`, string(fredag)).Scan(&status)
	if status != string(models.PuljeStatusOpen) {
		t.Errorf("after unlock: want status Open, got %s", status)
	}
}
