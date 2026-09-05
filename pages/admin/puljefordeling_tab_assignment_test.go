package admin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
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

func TestPuljefordelingIndex_AssignmentPickerOffersAllManualAssignmentTypes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en åpen puljefordeling med en billettholder.",
		When:  "Når administratoren åpner siden med tildelingsdialogen.",
		Then:  "Så skal dialogen tilby spilleder, spiller og førstevalg som separate handlinger.",
	})

	// Given
	expectedLabels := []string{"Legg til som spilleder", "Legg til som spiller", "Legg til som førstevalg"}
	expectedActions := []string{
		"/admin/api/puljefordeling/assign/gm",
		"/admin/api/puljefordeling/assign",
		"/admin/api/puljefordeling/assign/first-choice",
	}
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assignment_picker_actions")
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `
		INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1, 'Kari', 'Nordmann', 0, '', 0, 1)
	`)

	// When
	doc := templtest.Render(t, puljefordelingIndex(db, logger, fredag, nil))
	buttons := doc.Find("#puljefordeling-assign-dialog .billettholder-search-buttons button")

	// Then
	if buttons.Length() != len(expectedLabels) {
		t.Fatalf("assignment action count mismatch\nexpected: %d\nactual:   %d", len(expectedLabels), buttons.Length())
	}
	for index, expectedLabel := range expectedLabels {
		button := buttons.Eq(index)
		if actualLabel := strings.TrimSpace(button.Text()); actualLabel != expectedLabel {
			t.Fatalf("assignment label mismatch at index %d\nexpected: %q\nactual:   %q", index, expectedLabel, actualLabel)
		}
		if actualAction := button.AttrOr("data-on:click", ""); !strings.Contains(actualAction, expectedActions[index]) {
			t.Fatalf("assignment action mismatch at index %d\nexpected action to contain: %q\nactual: %q", index, expectedActions[index], actualAction)
		}
	}
}

func TestPuljefordelingGMRoute_AssignsAndRendersRemovableGM(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et arrangement uten spilleder i en åpen pulje.",
		When:  "Når administratoren tildeler en billettholder som spilleder.",
		Then:  "Så skal spillederen lagres og vises med en fjernhandling.",
	})

	// Given
	expectedAssignments := 1
	expectedRemovePath := "/admin/api/puljefordeling/FredagKveld/evA/1/gm"
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_gm_route")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, fredag)
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)

	// When
	rec := postAssignmentSignals(t, router, "/api/puljefordeling/assign/gm", 1, "evA", string(fredag))
	actualAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
		  AND role = ? AND source = ?
	`, fredag, models.EventPlayerRoleGM, models.EventPlayerSourceManual)
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag, nil))
	removeAction := doc.Find(".pulje-remove-gm").AttrOr("data-on:click", "")

	// Then
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d (%s)", rec.Code, rec.Body.String())
	}
	if actualAssignments != expectedAssignments {
		t.Fatalf("GM assignment count mismatch\nexpected: %d\nactual:   %d", expectedAssignments, actualAssignments)
	}
	if !strings.Contains(removeAction, expectedRemovePath) {
		t.Fatalf("GM remove action mismatch\nexpected action to contain: %q\nactual: %q", expectedRemovePath, removeAction)
	}
}

func TestPuljefordelingFirstChoiceRoute_AssignsAndRendersDistinctRemoval(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder og et arrangement i en åpen pulje.",
		When:  "Når administratoren tildeler arrangementet som førstevalg.",
		Then:  "Så skal førstevalget lagres og vises med sin egen fjernhandling.",
	})

	// Given
	expectedAssignments := 1
	expectedInterestLevel := models.InterestLevelHigh
	expectedRemovePath := "/admin/api/puljefordeling/FredagKveld/evA/1/first-choice"
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_assign_first_choice_route")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, fredag)
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)

	// When
	rec := postAssignmentSignals(t, router, "/api/puljefordeling/assign/first-choice", 1, "evA", string(fredag))
	actualAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
		  AND role = ? AND source = ?
	`, fredag, models.EventPlayerRolePlayer, models.EventPlayerSourceManual)
	var actualInterestLevel models.InterestLevel
	if err := db.QueryRow(`
		SELECT interest_level FROM interests
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
	`, fredag).Scan(&actualInterestLevel); err != nil {
		t.Fatalf("read first-choice interest: %v", err)
	}
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag, nil))
	removeAction := doc.Find(".pulje-remove-first-choice").AttrOr("data-on:click", "")

	// Then
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d (%s)", rec.Code, rec.Body.String())
	}
	if actualAssignments != expectedAssignments {
		t.Fatalf("first-choice assignment count mismatch\nexpected: %d\nactual:   %d", expectedAssignments, actualAssignments)
	}
	if actualInterestLevel != expectedInterestLevel {
		t.Fatalf("interest level mismatch\nexpected: %q\nactual:   %q", expectedInterestLevel, actualInterestLevel)
	}
	if !strings.Contains(removeAction, expectedRemovePath) {
		t.Fatalf("first-choice remove action mismatch\nexpected action to contain: %q\nactual: %q", expectedRemovePath, removeAction)
	}
}

func TestPuljefordelingRemoveGMRoute_DeletesManualGM(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en manuelt tildelt spilleder i en åpen pulje.",
		When:  "Når administratoren bruker spillederens fjernhandling.",
		Then:  "Så skal den manuelle spilledertildelingen fjernes.",
	})

	// Given
	expectedAssignments := 0
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_remove_gm_route")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, fredag)
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evA', ?, 1, ?, ?)
	`, fredag, models.EventPlayerRoleGM, models.EventPlayerSourceManual)

	// When
	req := httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/FredagKveld/evA/1/gm", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	actualAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1 AND role = ?
	`, fredag, models.EventPlayerRoleGM)

	// Then
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d (%s)", rec.Code, rec.Body.String())
	}
	if actualAssignments != expectedAssignments {
		t.Fatalf("GM assignment count mismatch\nexpected: %d\nactual:   %d", expectedAssignments, actualAssignments)
	}
}

func TestPuljefordelingRemoveFirstChoiceRoute_DeletesSeatAndInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et administrativt tildelt førstevalg i en åpen pulje.",
		When:  "Når administratoren bruker førstevalgets fjernhandling.",
		Then:  "Så skal både spillerplassen og den samsvarende interessen fjernes.",
	})

	// Given
	expectedAssignments := 0
	expectedInterests := 0
	db, logger := testutil.CreateTestDBAndLogger(t, "puljefordeling_remove_first_choice_route")
	router := chi.NewRouter()
	puljefordelingRoute(router, db, &live.Manager{}, logger, nil)
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, fredag)
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)
	if err := puljefordeling.AddFirstChoiceSeat(db, fredag, "evA", 1); err != nil {
		t.Fatalf("seed first-choice assignment: %v", err)
	}

	// When
	req := httptest.NewRequest(http.MethodDelete, "/api/puljefordeling/FredagKveld/evA/1/first-choice", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	actualAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
	`, fredag)
	actualInterests := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM interests
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
	`, fredag)

	// Then
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d (%s)", rec.Code, rec.Body.String())
	}
	if actualAssignments != expectedAssignments {
		t.Fatalf("first-choice assignment count mismatch\nexpected: %d\nactual:   %d", expectedAssignments, actualAssignments)
	}
	if actualInterests != expectedInterests {
		t.Fatalf("first-choice interest count mismatch\nexpected: %d\nactual:   %d", expectedInterests, actualInterests)
	}
}

func TestAddFirstChoiceThenEmulate_PinsAddedPlayer(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "puljefordeling_add_then_pin")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", models.PuljeStatusOpen, "2026-01-01 18:00")
	testutil.MustExec(t, db, `INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES ('evA','Alpha','','','','','',4)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES ('evA',?,1)`, string(fredag))
	testutil.MustExec(t, db, `INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		VALUES (1,'Kari','Nordmann',0,'',0,1)`)

	// Add Kari through the real picker add path (the + button's endpoint).
	if err := puljefordeling.AddFirstChoiceSeat(db, fredag, "evA", 1); err != nil {
		t.Fatalf("AddFirstChoiceSeat: %v", err)
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

func postAssignmentSignals(t *testing.T, router http.Handler, path string, bhID int, eventID, pulje string) *httptest.ResponseRecorder {
	t.Helper()
	body := fmt.Sprintf(`{"assignmentBillettholderId":%d,"assignmentEventId":%q,"assignmentPuljeId":%q}`, bhID, eventID, pulje)
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func postAssignSignals(t *testing.T, router http.Handler, bhID int, eventID, pulje string) *httptest.ResponseRecorder {
	t.Helper()
	return postAssignmentSignals(t, router, "/api/puljefordeling/assign", bhID, eventID, pulje)
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
