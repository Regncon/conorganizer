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
