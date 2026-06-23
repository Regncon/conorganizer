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

func TestPuljefordelingTabContent_RerunButtonOnlyWhenLocked(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_tab_rerun_button")
	const fredag = models.PuljeFredagKveld
	seedTabPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedTabEvent(t, db, "evF", "Fredagsspill", 4, fredag)

	if templtest.Render(t, PuljefordelingTabContent(db, logger, fredag)).
		Find(".puljefordeling-tab-rerun").Length() != 0 {
		t.Errorf("open pulje must not show the rerun button")
	}

	if _, err := db.Exec(`UPDATE puljer SET status = 'Locked' WHERE id = ?`, string(fredag)); err != nil {
		t.Fatalf("lock pulje: %v", err)
	}
	if templtest.Render(t, PuljefordelingTabContent(db, logger, fredag)).
		Find(".puljefordeling-tab-rerun").Length() == 0 {
		t.Errorf("locked pulje must show the rerun button")
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
