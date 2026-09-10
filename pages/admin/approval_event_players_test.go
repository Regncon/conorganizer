package admin

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/go-chi/chi/v5"
)

// approvalRouterFor mounts only the approval page's player/GM endpoints, which
// is all these tests need to exercise the placement rules.
func approvalRouterFor(t *testing.T, db *sql.DB) http.Handler {
	t.Helper()
	router := chi.NewRouter()
	approvalEventPlayersRoute(router, db, &live.Manager{}, testutil.NewTestLogger())
	return router
}

// postApprovalSignals posts the approval page's assignment signals. Extra
// signals (role flags, the age confirmation) are spliced in as raw JSON so each
// test spells out exactly what the browser would have sent.
func postApprovalSignals(t *testing.T, router http.Handler, method, path string, bhID int, eventID, pulje, extra string) *httptest.ResponseRecorder {
	t.Helper()
	body := fmt.Sprintf(
		`{"assignmentBillettholderId":%d,"assignmentEventId":%q,"assignmentPuljeId":%q%s}`,
		bhID, eventID, pulje, extra,
	)
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func countEventPlayerRows(t *testing.T, db *sql.DB) int {
	t.Helper()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE event_id='evA' AND billettholder_id=1`).Scan(&n); err != nil {
		t.Fatalf("count event player rows: %v", err)
	}
	return n
}

func countInterestRows(t *testing.T, db *sql.DB) int {
	t.Helper()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM interests WHERE event_id='evA' AND billettholder_id=1`).Scan(&n); err != nil {
		t.Fatalf("count interest rows: %v", err)
	}
	return n
}

const (
	approvalFirstChoicePath = "/event-players/post/add_first_choice"
	approvalAddGMPath       = "/event-players/post/add_gm"
	approvalUpdatePath      = "/event-players/update_status"
)

// The approval page seats people into the very same events the puljefordeling
// tab does, so it has to ask the same question before putting a minor in an 18+
// game: nothing is written until the admin confirms.
func TestApprovalFirstChoice_MinorInAdultsOnlyAsksForConfirmation(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_first_choice_minor_warns")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)
	router := approvalRouterFor(t, db)

	rec := postApprovalSignals(t, router, http.MethodPost, approvalFirstChoicePath, 1, "evA", "FredagKveld", "")
	body := rec.Body.String()

	if !strings.Contains(body, "datastar-patch-signals") {
		t.Fatalf("expected a datastar signal patch, got %d: %s", rec.Code, body)
	}
	if !strings.Contains(body, "ageWarningText") {
		t.Errorf("expected the ageWarningText signal in the response, got: %s", body)
	}
	if !strings.Contains(body, "er under 18") {
		t.Errorf("expected a Norwegian under-18 warning, got: %s", body)
	}
	if !strings.Contains(body, approvalFirstChoicePath) {
		t.Errorf("the warning should carry the endpoint to retry, got: %s", body)
	}
	if n := countEventPlayerRows(t, db); n != 0 {
		t.Fatalf("an unconfirmed placement must not seat anyone, found %d rows", n)
	}
	if n := countInterestRows(t, db); n != 0 {
		t.Fatalf("an unconfirmed placement must not write an interest either, found %d rows", n)
	}
}

// Once the admin has seen the warning and insists, the placement is a deliberate
// override and goes through unchanged.
func TestApprovalFirstChoice_MinorInAdultsOnlySeatsWhenConfirmed(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_first_choice_minor_confirmed")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)
	router := approvalRouterFor(t, db)

	rec := postApprovalSignals(t, router, http.MethodPost, approvalFirstChoicePath, 1, "evA", "FredagKveld", `,"assignmentAgeConfirmed":true`)
	if strings.Contains(rec.Body.String(), "ageWarningText") {
		t.Fatalf("a confirmed placement must not warn again: %s", rec.Body.String())
	}
	if n := countEventPlayerRows(t, db); n != 1 {
		t.Fatalf("confirmed placement should seat the participant, found %d rows", n)
	}
	if n := countInterestRows(t, db); n != 1 {
		t.Fatalf("confirmed first choice should still write the interest, found %d rows", n)
	}
}

func TestApprovalFirstChoice_AdultInAdultsOnlyNeedsNoConfirmation(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_first_choice_adult")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, true)
	router := approvalRouterFor(t, db)

	postApprovalSignals(t, router, http.MethodPost, approvalFirstChoicePath, 1, "evA", "FredagKveld", "")
	if n := countEventPlayerRows(t, db); n != 1 {
		t.Fatalf("an adult should be seated straight away, found %d rows", n)
	}
}

func TestApprovalFirstChoice_MinorInDefaultEventNeedsNoConfirmation(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_first_choice_minor_default")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupDefault, false)
	router := approvalRouterFor(t, db)

	postApprovalSignals(t, router, http.MethodPost, approvalFirstChoicePath, 1, "evA", "FredagKveld", "")
	if n := countEventPlayerRows(t, db); n != 1 {
		t.Fatalf("a minor in an ordinary game should be seated straight away, found %d rows", n)
	}
}

// "Legg til som spelar" on an interest row goes through update_status; it needs
// the same gate as the first-choice button.
func TestApprovalUpdateStatus_MinorInAdultsOnlyAsksForConfirmation(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_update_status_minor_warns")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)
	router := approvalRouterFor(t, db)

	rec := postApprovalSignals(t, router, http.MethodPut, approvalUpdatePath, 1, "evA", "FredagKveld", `,"assignmentIsPlayer":true,"assignmentIsGm":false`)
	if !strings.Contains(rec.Body.String(), "ageWarningText") {
		t.Fatalf("expected an age warning, got %d: %s", rec.Code, rec.Body.String())
	}
	if n := countEventPlayerRows(t, db); n != 0 {
		t.Fatalf("an unconfirmed placement must not seat anyone, found %d rows", n)
	}
}

func TestApprovalUpdateStatus_MinorInAdultsOnlySeatsWhenConfirmed(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_update_status_minor_confirmed")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)
	router := approvalRouterFor(t, db)

	postApprovalSignals(t, router, http.MethodPut, approvalUpdatePath, 1, "evA", "FredagKveld", `,"assignmentIsPlayer":true,"assignmentIsGm":false,"assignmentAgeConfirmed":true`)
	if n := countEventPlayerRows(t, db); n != 1 {
		t.Fatalf("confirmed placement should seat the participant, found %d rows", n)
	}
}

// Taking someone back out of an 18+ game resolves the conflict rather than
// creating one, so it must never be blocked by the warning.
func TestApprovalUpdateStatus_RemovalNeedsNoConfirmation(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_update_status_removal")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role) VALUES ('evA', 'FredagKveld', 1, 'Player')`)
	router := approvalRouterFor(t, db)

	rec := postApprovalSignals(t, router, http.MethodPut, approvalUpdatePath, 1, "evA", "FredagKveld", `,"assignmentIsPlayer":false,"assignmentIsGm":false`)
	if strings.Contains(rec.Body.String(), "ageWarningText") {
		t.Fatalf("removing a seat must not ask for age confirmation: %s", rec.Body.String())
	}
	if n := countEventPlayerRows(t, db); n != 0 {
		t.Fatalf("the seat should be gone, found %d rows", n)
	}
}

// A speleleiar under 18 in an 18+ game is the same deliberate override, so the
// GM button asks too.
func TestApprovalAddGM_MinorInAdultsOnlyAsksForConfirmation(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_add_gm_minor_warns")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)
	router := approvalRouterFor(t, db)

	rec := postApprovalSignals(t, router, http.MethodPost, approvalAddGMPath, 1, "evA", "FredagKveld", "")
	if !strings.Contains(rec.Body.String(), "ageWarningText") {
		t.Fatalf("expected an age warning, got %d: %s", rec.Code, rec.Body.String())
	}
	if n := countEventPlayerRows(t, db); n != 0 {
		t.Fatalf("an unconfirmed GM placement must not write a row, found %d", n)
	}
}

func TestApprovalAddGM_MinorInAdultsOnlySeatsWhenConfirmed(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "approval_add_gm_minor_confirmed")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)
	router := approvalRouterFor(t, db)

	postApprovalSignals(t, router, http.MethodPost, approvalAddGMPath, 1, "evA", "FredagKveld", `,"assignmentAgeConfirmed":true`)
	if n := countEventPlayerRows(t, db); n != 1 {
		t.Fatalf("confirmed GM placement should write a row, found %d", n)
	}
}
