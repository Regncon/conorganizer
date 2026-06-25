package admin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

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
	puljefordelingRoute(router, db, &live.Manager{}, logger)

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

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))
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

	doc := templtest.Render(t, puljefordelingIndex(db, logger, fredag))

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
	puljefordelingRoute(router, db, &live.Manager{}, logger)

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
	puljefordelingRoute(router, db, &live.Manager{}, logger)

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
	puljefordelingRoute(router, db, &live.Manager{}, logger)

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
	puljefordelingRoute(router, db, &live.Manager{}, logger)

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
	puljefordelingRoute(router, db, &live.Manager{}, logger)

	req := httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/NotAPulje/evA/1", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid pulje should be rejected with 400, got %d", rec.Code)
	}
}
