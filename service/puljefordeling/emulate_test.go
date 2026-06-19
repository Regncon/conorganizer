package puljefordeling

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func seedPulje(t *testing.T, db *sql.DB, id models.Pulje, name, startAt string) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO puljer (id, name, status, start_at, end_at) VALUES (?, ?, 'Open', ?, ?)`,
		string(id), name, startAt, startAt,
	)
	if err != nil {
		t.Fatalf("seed pulje %s: %v", id, err)
	}
}

func seedEvent(t *testing.T, db *sql.DB, id, title string, maxPlayers int, pulje models.Pulje) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		 VALUES (?, ?, '', '', '', '', '', ?)`,
		id, title, maxPlayers,
	)
	if err != nil {
		t.Fatalf("seed event %s: %v", id, err)
	}
	_, err = db.Exec(
		`INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje) VALUES (?, ?, 1)`,
		id, string(pulje),
	)
	if err != nil {
		t.Fatalf("place event %s in %s: %v", id, pulje, err)
	}
}

func seedParticipant(t *testing.T, db *sql.DB, id int, first, last string) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, order_id, ticket_id)
		 VALUES (?, ?, ?, 0, '', 0, ?)`,
		id, first, last, id,
	)
	if err != nil {
		t.Fatalf("seed participant %d: %v", id, err)
	}
}

func seedInterest(t *testing.T, db *sql.DB, bhID int, eventID string, pulje models.Pulje, level models.InterestLevel) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level) VALUES (?, ?, ?, ?)`,
		bhID, eventID, string(pulje), string(level),
	)
	if err != nil {
		t.Fatalf("seed interest bh=%d ev=%s: %v", bhID, eventID, err)
	}
}

func seedGM(t *testing.T, db *sql.DB, eventID string, pulje models.Pulje, bhID int) {
	t.Helper()
	_, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role) VALUES (?, ?, ?, 'GM')`,
		eventID, string(pulje), bhID,
	)
	if err != nil {
		t.Fatalf("seed GM bh=%d ev=%s: %v", bhID, eventID, err)
	}
}

func findEvent(p EmulatedPulje, eventID string) (EmulatedEvent, bool) {
	for _, e := range p.Events {
		if e.EventID == eventID {
			return e, true
		}
	}
	return EmulatedEvent{}, false
}

func playerNames(aps []AssignedPlayer) []string {
	out := make([]string, len(aps))
	for i, ap := range aps {
		out[i] = ap.Name
	}
	return out
}

func TestEmulateSeatings(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")

	// Event A: capacity 2. Event B: capacity 5 (fallback). GM event G run by p99.
	seedEvent(t, db, "evA", "Alpha", 2, fredag)
	seedEvent(t, db, "evB", "Bravo", 5, fredag)
	seedEvent(t, db, "evG", "Gamma", 4, fredag)

	// Participants.
	for _, p := range []struct {
		id          int
		first, last string
	}{
		{1, "Anna", "A"},
		{2, "Bob", "B"},
		{3, "Cara", "C"},
		{99, "Game", "Master"},
	} {
		seedParticipant(t, db, p.id, p.first, p.last)
	}

	// p99 is the GM of evG in this pulje.
	seedGM(t, db, "evG", fredag, 99)

	// All three regulars want evA highly (capacity only 2 → one loses),
	// and all also have a fallback interest in evB.
	for _, bh := range []int{1, 2, 3} {
		seedInterest(t, db, bh, "evA", fredag, models.InterestLevelHigh)
		seedInterest(t, db, bh, "evB", fredag, models.InterestLevelLow)
	}
	// The GM also expresses interest in evA, but must NOT be seated (busy).
	seedInterest(t, db, 99, "evA", fredag, models.InterestLevelHigh)

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	if em.Year != 2026 {
		t.Errorf("year: want 2026, got %d", em.Year)
	}
	if em.PlayerCount != 4 {
		t.Errorf("player count: want 4 (all with interest), got %d", em.PlayerCount)
	}
	if len(em.Puljer) != 1 {
		t.Fatalf("want 1 pulje, got %d", len(em.Puljer))
	}

	p := em.Puljer[0]

	evA, ok := findEvent(p, "evA")
	if !ok {
		t.Fatal("evA missing from result")
	}
	// Capacity respected.
	if len(evA.AssignedPlayers) != 2 {
		t.Errorf("evA capacity 2: want 2 seated, got %v", evA.AssignedPlayers)
	}
	// GM must never be seated as a participant in their own slot.
	if slices.Contains(playerNames(evA.AssignedPlayers), "Game Master") {
		t.Errorf("GM should not be seated as a player, got %v", playerNames(evA.AssignedPlayers))
	}
	// Everyone seated in evA wanted it highly → 🔥 level surfaced.
	for _, ap := range evA.AssignedPlayers {
		if ap.Level != models.InterestLevelHigh {
			t.Errorf("evA seat %q: want level High, got %q", ap.Name, ap.Level)
		}
	}

	// GM name surfaced on the event they run.
	evG, ok := findEvent(p, "evG")
	if !ok {
		t.Fatal("evG missing from result")
	}
	if evG.GMName != "Game Master" {
		t.Errorf("evG GM: want \"Game Master\", got %q", evG.GMName)
	}

	// The regular who lost the evA lottery should fall back to evB, not be
	// left unassigned (everyone had an evB interest).
	if len(p.Unassigned) != 0 {
		t.Errorf("want nobody unassigned (evB is a fallback), got %v", p.Unassigned)
	}
	evB, _ := findEvent(p, "evB")
	if len(evB.AssignedPlayers) != 1 {
		t.Errorf("evB should seat the 1 evA loser, got %v", evB.AssignedPlayers)
	}
	// The evB occupant simply lost the evA lottery and took their fallback — they
	// were never tentatively seated in evA, so no reverse-edge bump occurred and
	// they must NOT be flagged as moved. (Moved is reserved for players the solver
	// actively relocated off a higher-scoring seat; see the solver's bump test.)
	if len(evB.AssignedPlayers) == 1 {
		if got := evB.AssignedPlayers[0]; got.Moved {
			t.Errorf("evB occupant %q only lost contention and should not be marked Moved", got.Name)
		}
	}

	// Two regulars got their top-choice (evA score 5) → newly satisfied.
	if p.NewlySatisfied != 2 {
		t.Errorf("newly satisfied: want 2, got %d", p.NewlySatisfied)
	}
	if em.SatisfiedTotal != 2 {
		t.Errorf("satisfied total: want 2, got %d", em.SatisfiedTotal)
	}
}
