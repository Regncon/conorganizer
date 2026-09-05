package admin

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/components/formsubmission"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestPuljefordelingRemoveManualSeatRoute_DeletesPin(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_remove_route")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evA',?,1,'Player','manual')`, string(fredag))

	req := httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/FredagKveld/evA/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204 No Content, got %d (%s)", rec.Code, rec.Body.String())
	}

	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE event_id='evA' AND billettholder_id=1 AND source='manual'`,
	).Scan(&n); err != nil {
		t.Fatalf("count manual seats: %v", err)
	}
	if n != 0 {
		t.Fatalf("manual seat should be deleted, still found %d", n)
	}
}

func TestPuljefordelingTabContent_RendersAddPickerAndManualRemove(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_tab_interactive")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	// Kari is a manual pin → her tile must offer a × that removes her seat.
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evA',?,1,'Player','manual')`, string(fredag))

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag, nil))
	html, err := doc.Html()
	if err != nil {
		t.Fatalf("render html: %v", err)
	}

	// The × on Kari's manual tile deletes her manual seat.
	if !strings.Contains(html, "/admin/api/puljefordeling/FredagKveld/evA/1") {
		t.Errorf("manual tile should contain the remove URL")
	}

	// The + button opens the dialog scoped to this event (attribute decoded).
	addClick := doc.Find(".pulje-add").AttrOr("data-on:click", "")
	if !strings.Contains(addClick, "$assignmentEventId = 'evA'") {
		t.Errorf("+ button should set assignmentEventId to the event; got %q", addClick)
	}
}

// The picker is a modal dialog. If it renders inside the SSE-updated section
// (#puljefordeling-tab), every add re-renders the section and orphans the open
// modal's backdrop, locking the page. It must live in the stable outer wrapper.
func TestPuljefordelingIndex_DialogRendersOutsideLiveRegion(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_dialog_placement")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")

	doc := templtest.Render(t, puljefordelingIndex(db, logger, fredag, nil))

	if got := doc.Find("#puljefordeling-assign-dialog").Length(); got != 1 {
		t.Fatalf("expected exactly one assign dialog, got %d", got)
	}
	if got := doc.Find("#puljefordeling-tab #puljefordeling-assign-dialog").Length(); got != 0 {
		t.Errorf("assign dialog must NOT be inside the live #puljefordeling-tab section (orphans the modal backdrop on SSE re-render)")
	}
}

func TestAddFirstChoiceThenEmulate_PinsAddedPlayer(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_add_then_pin")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)

	// Add Kari through the real picker add path (the + button's endpoint).
	if err := formsubmission.AddPlayersFirstChoice(1, "evA", string(fredag), db, logger); err != nil {
		t.Fatalf("AddPlayersFirstChoice: %v", err)
	}

	// A subsequent emulation must pin her into evA, marked as a manual placement.
	em, err := puljefordeling.EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}
	var evA puljefordeling.EmulatedEvent
	for _, p := range em.Puljer {
		if p.PuljeID == fredag {
			for _, e := range p.Events {
				if e.EventID == "evA" {
					evA = e
				}
			}
		}
	}
	names := make([]string, len(evA.AssignedPlayers))
	for i, ap := range evA.AssignedPlayers {
		names[i] = ap.Name
	}
	if !slices.Contains(names, "Kari Nordmann") {
		t.Fatalf("added player should be pinned into evA, got %v", names)
	}
	for _, ap := range evA.AssignedPlayers {
		if ap.Name == "Kari Nordmann" && !ap.Manual {
			t.Errorf("added player should be marked as a manual placement")
		}
	}
}

func postAssignSignals(t *testing.T, router http.Handler, bhID int, eventID, pulje string) *httptest.ResponseRecorder {
	t.Helper()
	body := fmt.Sprintf(`{"assignmentBillettholderId":%d,"assignmentEventId":%q,"assignmentPuljeId":%q}`, bhID, eventID, pulje)
	req := httptest.NewRequest(http.MethodPost, "/api/puljefordeling/assign", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func TestPuljefordelingCommitRoute_PersistsSolverPicks(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_commit_route")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	testutil.MustExec(t, db, `INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level) VALUES (1,'evA',?,?)`,
		string(fredag), string(models.InterestLevelHigh))

	req := httptest.NewRequest(http.MethodPost, "/api/puljefordeling/FredagKveld/commit", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d (%s)", rec.Code, rec.Body.String())
	}

	var source string
	if err := db.QueryRow(
		`SELECT source FROM relation_events_players WHERE event_id='evA' AND billettholder_id=1 AND role='Player'`,
	).Scan(&source); err != nil {
		t.Fatalf("expected a committed seat after commit: %v", err)
	}
	if source != "solver" {
		t.Errorf("committed solver pick should be source='solver', got %q", source)
	}
}

func TestPuljefordelingAssignRoute_PinsWithoutCreatingInterest(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_route")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)

	rec := postAssignSignals(t, router, 1, "evA", "FredagKveld")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204, got %d (%s)", rec.Code, rec.Body.String())
	}

	var seats, interests int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE event_id='evA' AND billettholder_id=1 AND source='manual'`).Scan(&seats); err != nil {
		t.Fatalf("count seats: %v", err)
	}
	if seats != 1 {
		t.Fatalf("want 1 manual seat, got %d", seats)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM interests WHERE event_id='evA' AND billettholder_id=1`).Scan(&interests); err != nil {
		t.Fatalf("count interests: %v", err)
	}
	if interests != 0 {
		t.Fatalf("manual pin must not create an interest, found %d", interests)
	}
}

func TestPuljefordelingAssignRoute_RejectsWhenPublished(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_published")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusCompleted, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)

	rec := postAssignSignals(t, router, 1, "evA", "FredagKveld")
	if rec.Code != http.StatusConflict {
		t.Fatalf("assigning into a published pulje should be 409, got %d", rec.Code)
	}
	var seats int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE event_id='evA' AND billettholder_id=1`).Scan(&seats); err != nil {
		t.Fatalf("count seats: %v", err)
	}
	if seats != 0 {
		t.Fatalf("no seat should be created for a published pulje, found %d", seats)
	}
}

func TestPuljefordelingRemoveManualSeatRoute_RejectsWhenPublished(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_remove_published")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusCompleted, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evA',?,1,'Player','manual')`, string(fredag))

	req := httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/FredagKveld/evA/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("removing from a published pulje should be rejected with 409, got %d", rec.Code)
	}
	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE event_id='evA' AND billettholder_id=1 AND source='manual'`,
	).Scan(&n); err != nil {
		t.Fatalf("count manual seats: %v", err)
	}
	if n != 1 {
		t.Fatalf("manual seat must survive a rejected remove, found %d", n)
	}
}

func TestPuljefordelingRemoveManualSeatRoute_RejectsInvalidPulje(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_remove_route_invalid")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/NotAPulje/evA/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid pulje should be rejected with 400, got %d", rec.Code)
	}
}

// postAssignSignalsConfirmed posts the assign signals with the age confirmation
// flag the age-warning dialog sets when the admin insists on the placement.
func postAssignSignalsConfirmed(t *testing.T, router http.Handler, bhID int, eventID, pulje string) *httptest.ResponseRecorder {
	t.Helper()
	body := fmt.Sprintf(
		`{"assignmentBillettholderId":%d,"assignmentEventId":%q,"assignmentPuljeId":%q,"assignmentAgeConfirmed":true}`,
		bhID, eventID, pulje,
	)
	req := httptest.NewRequest(http.MethodPost, "/api/puljefordeling/assign", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

// seedAssignFixture sets up one open pulje, one event and one participant. The
// participant is a minor unless over18 is set (is_over_18 defaults to 0).
func seedAssignFixture(t *testing.T, db *sql.DB, pulje models.Pulje, ageGroup models.AgeGroup, over18 bool) {
	t.Helper()
	seedTabPulje(t, db, pulje, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players, age_group)
		VALUES ('evA','Voksenspel','','','','','',4,?)`, string(ageGroup))
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(pulje))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id, is_over_18)
		VALUES (1,'Kari','Nordmann',0,'',0,1,?)`, boolToInt(over18))
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func countManualSeats(t *testing.T, db *sql.DB) int {
	t.Helper()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE event_id='evA' AND billettholder_id=1 AND source='manual'`).Scan(&n); err != nil {
		t.Fatalf("count manual seats: %v", err)
	}
	return n
}

// An admin dropping a minor into an 18+ game must be asked first: the endpoint
// writes nothing and answers with the signals that open the confirm dialog.
func TestPuljefordelingAssignRoute_MinorInAdultsOnlyAsksForConfirmation(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_minor_warns")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)

	rec := postAssignSignals(t, router, 1, "evA", "FredagKveld")
	body := rec.Body.String()

	if !strings.Contains(body, "datastar-patch-signals") {
		t.Fatalf("expected a datastar signal patch, got %d: %s", rec.Code, body)
	}
	if !strings.Contains(body, "ageWarningText") {
		t.Errorf("expected the ageWarningText signal in the response, got: %s", body)
	}
	if !strings.Contains(body, "er under 18") {
		t.Errorf("expected a Norwegian under-18 warning in the response, got: %s", body)
	}
	if n := countManualSeats(t, db); n != 0 {
		t.Fatalf("an unconfirmed placement must not write a seat, found %d", n)
	}
}

// Once the admin confirms, the pin is their deliberate override and goes through.
func TestPuljefordelingAssignRoute_MinorInAdultsOnlyPinsWhenConfirmed(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_minor_confirmed")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, false)

	rec := postAssignSignalsConfirmed(t, router, 1, "evA", "FredagKveld")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204 after confirmation, got %d (%s)", rec.Code, rec.Body.String())
	}
	if n := countManualSeats(t, db); n != 1 {
		t.Fatalf("confirmed placement should create one manual seat, found %d", n)
	}
}

// An adult in an 18+ game is not an age conflict: unchanged, no dialog.
func TestPuljefordelingAssignRoute_AdultInAdultsOnlyNeedsNoConfirmation(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_adult")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupAdultsOnly, true)

	rec := postAssignSignals(t, router, 1, "evA", "FredagKveld")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204 for an adult, got %d (%s)", rec.Code, rec.Body.String())
	}
	if n := countManualSeats(t, db); n != 1 {
		t.Fatalf("adult placement should create one manual seat, found %d", n)
	}
}

// A minor in an ordinary game is not an age conflict either.
func TestPuljefordelingAssignRoute_MinorInDefaultEventNeedsNoConfirmation(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_minor_default")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupDefault, false)

	rec := postAssignSignals(t, router, 1, "evA", "FredagKveld")
	if rec.Code != http.StatusNoContent {
		t.Fatalf("want 204 for a minor in a Default event, got %d (%s)", rec.Code, rec.Body.String())
	}
	if n := countManualSeats(t, db); n != 1 {
		t.Fatalf("placement should create one manual seat, found %d", n)
	}
}

func TestPuljefordelingAssignRoute_UnknownBillettholderIs404(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_unknown_bh")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupDefault, false)

	rec := postAssignSignals(t, router, 999, "evA", "FredagKveld")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown billettholder should be 404, got %d (%s)", rec.Code, rec.Body.String())
	}
}

func TestPuljefordelingAssignRoute_UnknownEventIs404(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_unknown_event")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)

	seedAssignFixture(t, db, models.PuljeFredagKveld, models.AgeGroupDefault, false)

	rec := postAssignSignals(t, router, 1, "ukjent", "FredagKveld")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown event should be 404, got %d (%s)", rec.Code, rec.Body.String())
	}
}

// Like the picker dialog, the age warning dialog must sit outside the live
// region — an SSE re-render of the section would orphan its open backdrop.
func TestPuljefordelingIndex_AgeDialogRendersOutsideLiveRegion(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_age_dialog_placement")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")

	doc := templtest.Render(t, puljefordelingIndex(db, logger, fredag, nil))

	if got := doc.Find("#puljefordeling-age-dialog").Length(); got != 1 {
		t.Fatalf("expected exactly one age warning dialog, got %d", got)
	}
	if got := doc.Find("#puljefordeling-tab #puljefordeling-age-dialog").Length(); got != 0 {
		t.Errorf("age dialog must NOT be inside the live #puljefordeling-tab section")
	}

	html, err := doc.Html()
	if err != nil {
		t.Fatalf("render html: %v", err)
	}
	if !strings.Contains(html, "assignmentAgeConfirmed") {
		t.Errorf("the age dialog's confirm button should set assignmentAgeConfirmed")
	}
	if !strings.Contains(html, "/admin/api/puljefordeling/assign") {
		t.Errorf("the age dialog's confirm button should re-post to the assign endpoint")
	}
}

// datastar 1.0.x compiles every expression with a plain Function(), so `await`
// is a SyntaxError and the handler silently never runs. The confirm button must
// dispatch the request and reset the confirmation flag synchronously.
func TestPuljefordelingIndex_AgeDialogConfirmResetsFlagSynchronously(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_age_dialog_no_await")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")

	doc := templtest.Render(t, puljefordelingIndex(db, logger, models.PuljeFredagKveld, nil))
	dialog := doc.Find("#puljefordeling-age-dialog")
	if dialog.Length() != 1 {
		t.Fatalf("expected one age dialog, got %d", dialog.Length())
	}

	confirm := dialog.Find("button").FilterFunction(func(_ int, s *goquery.Selection) bool {
		return strings.Contains(s.AttrOr("data-on:click", ""), "assignmentAgeConfirmed = true")
	})
	if confirm.Length() != 1 {
		t.Fatalf("expected exactly one confirm button setting assignmentAgeConfirmed, got %d", confirm.Length())
	}
	action := confirm.AttrOr("data-on:click", "")
	if strings.Contains(action, "await") {
		t.Errorf("confirm action must not use await (datastar compiles with Function, not AsyncFunction); got %q", action)
	}
	if strings.Contains(action, ".finally(") || strings.Contains(action, ".then(") {
		t.Errorf("confirm action must reset the flag synchronously, not in a promise callback; got %q", action)
	}
	// The flag is page-global and sent with every request, so it must be reset
	// on the statement right after the action call (datastar serialises the
	// payload synchronously); resetting later would pre-confirm any placement
	// fired while the re-post is in flight.
	post := strings.Index(action, "@post(")
	reset := strings.Index(action, "$assignmentAgeConfirmed = false")
	if post < 0 || reset < 0 || reset < post {
		t.Errorf("confirm action should reset assignmentAgeConfirmed right after dispatching the request; got %q", action)
	}
}

// closedBy="any" lets Esc or a backdrop click dismiss the dialog. If that path
// left $ageWarningText set, an identical warning for the same person and game
// would patch an unchanged value and the data-effect would never re-open it.
func TestPuljefordelingIndex_AgeDialogClearsWarningOnClose(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_age_dialog_on_close")
	seedTabPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")

	doc := templtest.Render(t, puljefordelingIndex(db, logger, models.PuljeFredagKveld, nil))
	onClose := doc.Find("#puljefordeling-age-dialog").AttrOr("data-on:close", "")
	if !strings.Contains(onClose, "$ageWarningText = ''") {
		t.Errorf("age dialog should clear $ageWarningText on close; got %q", onClose)
	}
}
