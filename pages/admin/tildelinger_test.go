package admin

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestTildeling_AddMenuListsAllAssignmentsBeforePinning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Kari is GM for three arrangementer in one pulje.", When: "An admin adds her as player through Legg til.", Then: "All GM assignments are shown before any player seat is written."})
	// Given
	expectedTitles := []string{"Arrangement X", "Arrangement Y", "Arrangement Z"}
	db, router := tildelingsFixture(t)
	for _, eventID := range []string{"evA", "evB", "evC"} {
		testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES (?, 'FredagKveld', 1, 'GM')`, eventID)
	}
	// When
	rec := postTildeling(t, router, "evB", "Player", true, "")
	// Then
	if rec.Code != http.StatusOK {
		t.Fatalf("expected confirmation, got %d: %s", rec.Code, rec.Body.String())
	}
	for _, title := range expectedTitles {
		if !strings.Contains(rec.Body.String(), title) {
			t.Errorf("confirmation missing %q: %s", title, rec.Body.String())
		}
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE role='Player'`); got != 0 {
		t.Fatalf("unconfirmed add wrote %d player seats", got)
	}
	confirmationFromResponse(t, rec)
}

func TestTildeling_ConfirmedAddRetainsEveryGMAndPlayerPin(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Kari is GM for X, Y and Z and pinned as player on Y.", When: "An admin confirms adding her as player on Z.", Then: "All three GM assignments and both player pins remain."})
	// Given
	const expectedGMs = 3
	const expectedPlayers = 2
	db, router := tildelingsFixture(t)
	for _, eventID := range []string{"evA", "evB", "evC"} {
		testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES (?, 'FredagKveld', 1, 'GM')`, eventID)
	}
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evB','FredagKveld',1,'Player','manual')`)
	warning := postTildeling(t, router, "evC", "Player", true, "")
	// When
	rec := postTildeling(t, router, "evC", "Player", true, confirmationFromResponse(t, warning))
	// Then
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected confirmed add, got %d: %s", rec.Code, rec.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE role='GM'`); got != expectedGMs {
		t.Fatalf("want %d GM assignments, got %d", expectedGMs, got)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE role='Player'`); got != expectedPlayers {
		t.Fatalf("want %d player pins, got %d", expectedPlayers, got)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE event_id='evC' AND role='Player' AND source='manual'`); got != 1 {
		t.Fatal("player pin was not added to Z")
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE event_id='evB' AND role='Player' AND source='manual'`); got != 1 {
		t.Fatal("adding a player pin removed the existing pin on Y")
	}
}

func TestTildeling_DragCannotConfirmGMOverlap(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Kari is GM in the selected pulje.", When: "Her player tile is dragged to another arrangement.", Then: "The popup directs the admin to Legg til without offering an override."})
	// Given
	const expectedPins = 0
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('evA','FredagKveld',1,'GM')`)
	// When
	rec := postTildeling(t, router, "evB", "Player", false, "")
	// Then
	if !strings.Contains(rec.Body.String(), "Legg til") || strings.Contains(rec.Body.String(), "data-confirmation=") {
		t.Fatalf("drag should direct admin to add-menu without override: %s", rec.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE role='Player'`); got != expectedPins {
		t.Fatalf("drag wrote %d unexpected pins", got)
	}
}

func TestTildeling_DialogSurvivesDistributionRefresh(t *testing.T) {
	// Given
	const selector = "#tildeling-dialog"
	db, _ := tildelingsFixture(t)
	// When
	doc := templtest.Render(t, puljefordelingIndex(db, testutil.NewTestLogger(), models.PuljeFredagKveld, nil))
	// Then
	if doc.Find(selector).Length() != 1 || doc.Find("#puljefordeling-tab "+selector).Length() != 0 {
		t.Fatal("assignment confirmation must live once outside the refreshed distribution")
	}
}

func TestTildeling_AddGMRetainsPlayerOnSameArrangement(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Kari is a player on X and GM on Y.", When: "An admin confirms adding her as GM on X through the approval add-menu.", Then: "Her player seat and both GM assignments remain."})
	// Given
	const expectedRoles = 3
	db, _ := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('evA','FredagKveld',1,'Player'), ('evB','FredagKveld',1,'GM')`)
	router := approvalRouterFor(t, db)
	warning := postApprovalSignals(t, router, http.MethodPost, approvalAddGMPath, 1, "evA", "FredagKveld", "")
	confirmation := confirmationFromResponse(t, warning)
	// When
	rec := postApprovalSignals(t, router, http.MethodPost, approvalAddGMPath, 1, "evA", "FredagKveld", `,"assignmentConfirmation":"`+confirmation+`"`)
	// Then
	if rec.Code != http.StatusNoContent {
		t.Fatalf("GM confirmation failed: %d %s", rec.Code, rec.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1`); got != expectedRoles {
		t.Fatalf("want %d roles, got %d", expectedRoles, got)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE event_id='evA' AND role='Player'`); got != 1 {
		t.Fatal("adding GM removed the player seat")
	}
}

func TestTildeling_RemovingPlayerRetainsGMOnSameArrangement(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Kari has both GM and Player roles on X.", When: "An admin removes her Player role.", Then: "The GM role remains."})
	// Given
	const expectedGMs = 1
	db, _ := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('evA','FredagKveld',1,'Player'), ('evA','FredagKveld',1,'GM')`)
	router := approvalRouterFor(t, db)
	// When
	rec := postApprovalSignals(t, router, http.MethodPut, approvalUpdatePath, 1, "evA", "FredagKveld", `,"assignmentRole":"Player","assignmentRemove":true`)
	// Then
	if rec.Code != http.StatusNoContent {
		t.Fatalf("role removal failed: %d %s", rec.Code, rec.Body.String())
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE role='GM'`); got != expectedGMs {
		t.Fatalf("want %d GM, got %d", expectedGMs, got)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE role='Player'`); got != 0 {
		t.Fatalf("player role remained: %d", got)
	}
}

func TestTildeling_ChangedAssignmentsRefreshConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "An admin has a confirmation for Kari's existing assignment.", When: "A second assignment is added before the confirmation is submitted.", Then: "The new assignment is displayed and the player pin waits for fresh confirmation."})
	// Given
	const expectedPins = 0
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('evA','FredagKveld',1,'GM')`)
	warning := postTildeling(t, router, "evB", "Player", true, "")
	oldConfirmation := confirmationFromResponse(t, warning)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('evC','FredagKveld',1,'GM')`)
	// When
	rec := postTildeling(t, router, "evB", "Player", true, oldConfirmation)
	// Then
	if got := confirmationFromResponse(t, rec); got == oldConfirmation {
		t.Fatal("confirmation did not change with assignments")
	}
	if !strings.Contains(rec.Body.String(), "Arrangement Z") {
		t.Fatal("refreshed confirmation omitted the new GM assignment")
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE role='Player'`); got != expectedPins {
		t.Fatalf("stale confirmation wrote %d pins", got)
	}
}

func TestTildeling_AgeAndGMWarningsShareOneConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "Kari is under 18 and GM in the pulje.", When: "An admin adds her as player to an adults-only arrangement.", Then: "One confirmation includes her GM assignment and the age warning."})
	// Given
	const expectedPins = 0
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `UPDATE billettholdere SET is_over_18=0 WHERE id=1`)
	testutil.MustExec(t, db, `UPDATE events SET age_group=? WHERE id='evB'`, models.AgeGroupAdultsOnly)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role) VALUES ('evA','FredagKveld',1,'GM')`)
	// When
	rec := postTildeling(t, router, "evB", "Player", true, "")
	// Then
	for _, expected := range []string{"tildeling-dialog-innhold", "Arrangement X", "under 18"} {
		if !strings.Contains(rec.Body.String(), expected) {
			t.Errorf("missing %q: %s", expected, rec.Body.String())
		}
	}
	if strings.Contains(rec.Body.String(), "ageWarningText") {
		t.Fatal("combined warning should not open a second age dialog")
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE role='Player'`); got != expectedPins {
		t.Fatalf("warning wrote %d pins", got)
	}
}

func tildelingsFixture(t *testing.T) (*sql.DB, http.Handler) {
	t.Helper()
	db, logger := testutil.CreateTestDBAndLogger(t, "tildelinger")
	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupDefault, true)
	testutil.MustExec(t, db, `UPDATE events SET title='Arrangement X' WHERE id='evA'`)
	for eventID, title := range map[string]string{"evB": "Arrangement Y", "evC": "Arrangement Z"} {
		testutil.MustExec(t, db, `INSERT INTO events(id,title,intro,description,host_name,email,phone_number,max_players,is_in_puljefordeling) VALUES (?,?,'','','','','',4,1)`, eventID, title)
		testutil.MustExec(t, db, `INSERT INTO relation_event_puljer(event_id,pulje_id,is_in_pulje) VALUES (?,'FredagKveld',1)`, eventID)
	}
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)
	return db, router
}

func postTildeling(t *testing.T, router http.Handler, eventID, role string, fromAddMenu bool, confirmation string) *httptest.ResponseRecorder {
	t.Helper()
	data, err := json.Marshal(map[string]any{"assignmentBillettholderId": 1, "assignmentEventId": eventID, "assignmentPuljeId": "FredagKveld", "assignmentRole": role, "assignmentFromAddMenu": fromAddMenu, "assignmentConfirmation": confirmation})
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/puljefordeling/assign", strings.NewReader(string(data)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func confirmationFromResponse(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	match := regexp.MustCompile(`data-confirmation="([a-f0-9]+)"`).FindStringSubmatch(rec.Body.String())
	if len(match) != 2 {
		t.Fatalf("missing confirmation token in %d: %s", rec.Code, rec.Body.String())
	}
	return match[1]
}
