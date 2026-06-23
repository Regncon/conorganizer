package admin

import (
	"database/sql"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestPuljefordelingTabContent_RendersPuljeAndControls(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_puljefordeling_tab_content")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedTabParticipant(t, db, 1, "Anna", "A")
	seedTabInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))

	if doc.Find("h3:contains('Fredagsspill')").Length() == 0 {
		html, _ := doc.Html()
		t.Errorf("expected event title in rendered tab, got:\n%s", html)
	}
	if doc.Find("span:contains('Puljefordeling lukket')").Length() == 0 {
		t.Errorf("expected lock control in rendered tab")
	}
	dataInit := doc.Find("#puljefordeling-tab").AttrOr("data-init", "")
	if !strings.Contains(dataInit, "/admin/puljefordeling/api/FredagKveld") {
		t.Errorf("expected data-init SSE url, got: %q", dataInit)
	}
}

func seedTabPulje(t *testing.T, db *sql.DB, id models.Pulje, name, startAt string) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, ?, 'Open', ?, ?)`,
		string(id), name, startAt, startAt,
	); err != nil {
		t.Fatalf("seed pulje %s: %v", id, err)
	}
}

func seedTabEvent(t *testing.T, db *sql.DB, id, title string, maxPlayers int, pulje models.Pulje) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES (?, ?, '', '', '', '', '', ?)`,
		id, title, maxPlayers,
	); err != nil {
		t.Fatalf("seed event %s: %v", id, err)
	}
	if _, err := db.Exec(
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES (?, ?, 1)`,
		id, string(pulje),
	); err != nil {
		t.Fatalf("place event %s in %s: %v", id, pulje, err)
	}
}

func seedTabParticipant(t *testing.T, db *sql.DB, id int, first, last string) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		 VALUES (?, ?, ?, 0, '', 0, ?)`,
		id, first, last, id,
	); err != nil {
		t.Fatalf("seed participant %d: %v", id, err)
	}
}

func seedTabInterest(t *testing.T, db *sql.DB, bhID int, eventID string, pulje models.Pulje, level models.InterestLevel) {
	t.Helper()
	if _, err := db.Exec(
		`INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level) VALUES (?, ?, ?, ?)`,
		bhID, eventID, string(pulje), string(level),
	); err != nil {
		t.Fatalf("seed interest bh=%d ev=%s: %v", bhID, eventID, err)
	}
}

func TestPuljefordelingTabContent_RerunButtonStates(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_tab_rerun_button")
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)

	// Open pulje: the button is present but disabled and reads "Auto fordeling".
	openBtn := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag)).Find(".puljefordeling-tab-rerun")
	if openBtn.Length() == 0 {
		t.Fatal("open pulje must still show the (disabled) auto-fordeling button")
	}
	if _, disabled := openBtn.Attr("disabled"); !disabled {
		t.Errorf("open pulje button must be disabled")
	}
	if strings.TrimSpace(openBtn.Text()) != "Auto fordeling" {
		t.Errorf("open pulje button must read 'Auto fordeling', got %q", openBtn.Text())
	}

	// Locked pulje: the button is enabled and reads "Rerun fordeling".
	if _, err := db.Exec(`UPDATE puljer SET status = 'Locked' WHERE id = ?`, string(fredag)); err != nil {
		t.Fatalf("lock pulje: %v", err)
	}
	lockedBtn := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag)).Find(".puljefordeling-tab-rerun")
	if _, disabled := lockedBtn.Attr("disabled"); disabled {
		t.Errorf("locked pulje button must be enabled")
	}
	if strings.TrimSpace(lockedBtn.Text()) != "Rerun fordeling" {
		t.Errorf("locked pulje button must read 'Rerun fordeling', got %q", lockedBtn.Text())
	}
}

func TestPuljefordelingTabContent_RemoveGatedByState(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_tab_remove_gating")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedTabParticipant(t, db, 1, "Anna", "A")
	seedTabInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)

	// Open pulje: Anna is solver-placed (not pinned) -> no remove button.
	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))
	if doc.Find(".puljefordeling-tab-remove").Length() != 0 {
		t.Errorf("open pulje: solver-placed player must not have a remove button")
	}
	if doc.Find(".puljefordeling-tab-add").Length() == 0 {
		t.Errorf("expected an add (+) button on the game card")
	}

	// Lock the pulje: every seated player becomes removable.
	if _, err := db.Exec(`UPDATE puljer SET status = 'Locked' WHERE id = ?`, string(fredag)); err != nil {
		t.Fatalf("lock pulje: %v", err)
	}
	// Persist Anna as a committed seat so the locked view shows her.
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES ('evF', ?, 1, 'Player', 'solver')`, string(fredag),
	); err != nil {
		t.Fatalf("seed committed seat: %v", err)
	}
	docLocked := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))
	if docLocked.Find(".puljefordeling-tab-remove").Length() == 0 {
		t.Errorf("locked pulje: seated player must have a remove button")
	}
}

func TestPuljefordelingTabContent_RendersDragDropWiring(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_tab_drag_drop_wiring")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedTabParticipant(t, db, 1, "Anna", "A")
	seedTabInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)

	doc := templtest.Render(t, PuljefordelingTabContent(db, logger, fredag))

	eventCard := doc.Find(`.puljefordeling-tab-event[data-dnd-accept="pulje-player"]`)
	if eventCard.Length() == 0 {
		t.Fatal("expected puljefordeling event card to be a player drop target")
	}
	if got := eventCard.AttrOr("data-dnd-drop-url-template", ""); got != "/admin/puljefordeling/api/FredagKveld/move/evF/{id}" {
		t.Fatalf("drop url template mismatch\nexpected: %s\nactual:   %s", "/admin/puljefordeling/api/FredagKveld/move/evF/{id}", got)
	}

	player := doc.Find(`.puljefordeling-tab-players li[draggable="true"][data-dnd-kind="pulje-player"][data-dnd-id="1"]`)
	if player.Length() == 0 {
		t.Fatal("expected assigned player to be draggable")
	}
}

func TestMovePuljePlayerAssignment_MovesPlayerOnlyAssignment(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_move_pulje_player_assignment")

	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "old-event", "Old Event", 4, fredag)
	seedTabEvent(t, db, "new-event", "New Event", 4, fredag)
	seedTabEvent(t, db, "gm-event", "GM Event", 4, fredag)
	seedTabParticipant(t, db, 1, "Anna", "A")
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES ('old-event', ?, 1, 'Player', 'solver')`, string(fredag),
	); err != nil {
		t.Fatalf("seed old player assignment: %v", err)
	}
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES ('gm-event', ?, 1, 'GM', 'manual')`, string(fredag),
	); err != nil {
		t.Fatalf("seed gm assignment: %v", err)
	}

	if err := movePuljePlayerAssignment(db, fredag, "new-event", 1); err != nil {
		t.Fatalf("move player assignment: %v", err)
	}

	var playerEvent string
	var playerSource string
	if err := db.QueryRow(
		`SELECT event_id, source FROM relation_events_players WHERE pulje_id = ? AND billettholder_id = 1 AND role = 'Player'`,
		string(fredag),
	).Scan(&playerEvent, &playerSource); err != nil {
		t.Fatalf("read moved player assignment: %v", err)
	}
	if playerEvent != "new-event" || playerSource != models.EventPlayerSourceManual {
		t.Fatalf("moved player mismatch\nexpected event/source: new-event/%s\nactual event/source:   %s/%s", models.EventPlayerSourceManual, playerEvent, playerSource)
	}

	var gmCount int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE event_id = 'gm-event' AND pulje_id = ? AND billettholder_id = 1 AND role = 'GM'`,
		string(fredag),
	).Scan(&gmCount); err != nil {
		t.Fatalf("count gm assignment: %v", err)
	}
	if gmCount != 1 {
		t.Fatalf("expected GM assignment to be preserved, got %d", gmCount)
	}
}
