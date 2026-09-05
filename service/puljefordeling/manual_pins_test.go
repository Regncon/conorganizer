package puljefordeling

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestEmulateSeatings_ManualPinSeatsPlayerWithoutInterest(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "emulate_manual_pin")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")

	// Kari expressed NO interest in evA. A manual seat (source='manual') must
	// still pin her into evA when the distribution is emulated.
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES (?, ?, ?, 'Player', 'manual')`,
		"evA", string(fredag), 1,
	); err != nil {
		t.Fatalf("seed manual seat: %v", err)
	}

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	evA, ok := findEvent(em.Puljer[0], "evA")
	if !ok {
		t.Fatal("evA missing from result")
	}
	if !slices.Contains(playerNames(evA.AssignedPlayers), "Kari Nordmann") {
		t.Fatalf("manual pin should seat Kari in evA, got %v", playerNames(evA.AssignedPlayers))
	}
	for _, ap := range evA.AssignedPlayers {
		if ap.Name == "Kari Nordmann" {
			if !ap.Manual {
				t.Errorf("Kari was added manually and should be marked Manual")
			}
			if ap.BillettholderID != 1 {
				t.Errorf("assigned player should carry the billettholder id for removal: want 1, got %d", ap.BillettholderID)
			}
		}
	}
}

func manualSeatCount(t *testing.T, db interface {
	QueryRow(string, ...any) *sql.Row
}, eventID string, pulje models.Pulje, bhID int) int {
	t.Helper()
	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players
		 WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ? AND source = 'manual' AND role = 'Player'`,
		eventID, string(pulje), bhID,
	).Scan(&n); err != nil {
		t.Fatalf("count manual seats: %v", err)
	}
	return n
}

func interestCount(t *testing.T, db interface {
	QueryRow(string, ...any) *sql.Row
}, eventID string, pulje models.Pulje, bhID int) int {
	t.Helper()
	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM interests
		 WHERE event_id = ? AND pulje_id = ? AND billettholder_id = ?`,
		eventID, string(pulje), bhID,
	).Scan(&n); err != nil {
		t.Fatalf("count interests: %v", err)
	}
	return n
}

func TestAddManualSeat_CreatesPinAndRemovesOnlyMatchingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billetthelder med interesser i to arrangementer i samme pulje.",
		When:  "Når en administrator tildeler ett av arrangementene manuelt.",
		Then:  "Så skal tildelingen lagres og bare den samsvarende interessen fjernes.",
	})

	// Given
	expectedManualSeats := 1
	expectedMatchingInterests := 0
	expectedOtherInterests := 1
	db, _ := testutil.CreateTestDBAndLogger(t, "add_manual_seat")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedEvent(t, db, "evB", "Bravo", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 1, "evB", fredag, models.InterestLevelMedium)

	// When
	if err := AddManualSeat(db, fredag, "evA", 1); err != nil {
		t.Fatalf("AddManualSeat: %v", err)
	}
	actualManualSeats := manualSeatCount(t, db, "evA", fredag, 1)
	actualMatchingInterests := interestCount(t, db, "evA", fredag, 1)
	actualOtherInterests := interestCount(t, db, "evB", fredag, 1)

	// Then
	if actualManualSeats != expectedManualSeats {
		t.Fatalf("manual seat count mismatch\nexpected: %d\nactual:   %d", expectedManualSeats, actualManualSeats)
	}
	if actualMatchingInterests != expectedMatchingInterests {
		t.Fatalf("matching interest count mismatch\nexpected: %d\nactual:   %d", expectedMatchingInterests, actualMatchingInterests)
	}
	if actualOtherInterests != expectedOtherInterests {
		t.Fatalf("other interest count mismatch\nexpected: %d\nactual:   %d", expectedOtherInterests, actualOtherInterests)
	}
}

func TestAddManualSeat_MovePreservesOpenRegistration(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billetthelder med en åpen påmelding og en ordinær manuell tildeling i samme pulje.",
		When:  "Når administratoren flytter den ordinære tildelingen til et annet arrangement.",
		Then:  "Så skal påmeldingen beholdes ved siden av den nye manuelle tildelingen.",
	})

	// Given
	expectedTotalPlayerSeats := 2
	expectedOldManualSeats := 0
	expectedNewManualSeats := 1
	expectedRegistrations := 1
	db, _ := testutil.CreateTestDBAndLogger(t, "move_manual_seat_preserves_registration")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedEvent(t, db, "evB", "Bravo", 4, fredag)
	seedEvent(t, db, "evOpen", "Open", 100, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evOpen', ?, 1, ?, ?)
	`, fredag, models.EventPlayerRolePlayer, models.EventPlayerSourceRegistration)
	if err := AddManualSeat(db, fredag, "evA", 1); err != nil {
		t.Fatalf("AddManualSeat evA: %v", err)
	}

	// When
	if err := AddManualSeat(db, fredag, "evB", 1); err != nil {
		t.Fatalf("AddManualSeat evB: %v", err)
	}
	actualTotalPlayerSeats := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE pulje_id = ? AND billettholder_id = 1 AND role = ?
	`, fredag, models.EventPlayerRolePlayer)
	actualOldManualSeats := manualSeatCount(t, db, "evA", fredag, 1)
	actualNewManualSeats := manualSeatCount(t, db, "evB", fredag, 1)
	actualRegistrations := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evOpen' AND pulje_id = ? AND billettholder_id = 1 AND source = ?
	`, fredag, models.EventPlayerSourceRegistration)

	// Then
	if actualTotalPlayerSeats != expectedTotalPlayerSeats {
		t.Fatalf("total player seat count mismatch\nexpected: %d\nactual:   %d", expectedTotalPlayerSeats, actualTotalPlayerSeats)
	}
	if actualOldManualSeats != expectedOldManualSeats {
		t.Fatalf("old manual seat count mismatch\nexpected: %d\nactual:   %d", expectedOldManualSeats, actualOldManualSeats)
	}
	if actualNewManualSeats != expectedNewManualSeats {
		t.Fatalf("new manual seat count mismatch\nexpected: %d\nactual:   %d", expectedNewManualSeats, actualNewManualSeats)
	}
	if actualRegistrations != expectedRegistrations {
		t.Fatalf("registration count mismatch\nexpected: %d\nactual:   %d", expectedRegistrations, actualRegistrations)
	}
}

func TestAddManualSeat_MoveClearsPreviousEventSeat(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "add_manual_seat_move")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedEvent(t, db, "evB", "Bravo", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")

	// Pin into A, then move to B. The participant must end up in exactly one event.
	if err := AddManualSeat(db, fredag, "evA", 1); err != nil {
		t.Fatalf("AddManualSeat A: %v", err)
	}
	if err := AddManualSeat(db, fredag, "evB", 1); err != nil {
		t.Fatalf("AddManualSeat B: %v", err)
	}

	var count int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id=? AND billettholder_id=1 AND role='Player'`,
		string(fredag),
	).Scan(&count); err != nil {
		t.Fatalf("count seats: %v", err)
	}
	if count != 1 {
		t.Fatalf("moving a pin must leave exactly one seat, got %d", count)
	}

	var eventID string
	if err := db.QueryRow(
		`SELECT event_id FROM relation_events_players WHERE pulje_id=? AND billettholder_id=1`,
		string(fredag),
	).Scan(&eventID); err != nil {
		t.Fatalf("read seat: %v", err)
	}
	if eventID != "evB" {
		t.Errorf("after move the single seat should be in evB, got %q", eventID)
	}
}

func TestRemoveManualSeat_DeletesManualPlayerSeat(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "remove_manual_seat")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES (?, ?, ?, 'Player', 'manual')`,
		"evA", string(fredag), 1,
	); err != nil {
		t.Fatalf("seed manual seat: %v", err)
	}

	if err := RemoveManualSeat(db, fredag, "evA", 1); err != nil {
		t.Fatalf("RemoveManualSeat: %v", err)
	}

	if got := manualSeatCount(t, db, "evA", fredag, 1); got != 0 {
		t.Fatalf("manual seat should be deleted, still found %d", got)
	}
}

func TestRemoveManualSeat_LeavesSolverSeatIntact(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "remove_manual_seat_guard")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	// A solver-sourced seat must not be removable through the manual-remove path.
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES (?, ?, ?, 'Player', 'solver')`,
		"evA", string(fredag), 1,
	); err != nil {
		t.Fatalf("seed solver seat: %v", err)
	}

	if err := RemoveManualSeat(db, fredag, "evA", 1); err != nil {
		t.Fatalf("RemoveManualSeat: %v", err)
	}

	var n int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE source = 'solver'`,
	).Scan(&n); err != nil {
		t.Fatalf("count solver seats: %v", err)
	}
	if n != 1 {
		t.Fatalf("solver seat must remain, found %d", n)
	}
}

func TestEmulateSeatings_PinnedPlayerCountsAsSatisfiedAndScored(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "emulate_manual_pin_satisfied")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	// One seat in evA; two participants want it as their top choice.
	seedEvent(t, db, "evA", "Alpha", 1, fredag)
	seedParticipant(t, db, 1, "Pin", "Ned")
	seedParticipant(t, db, 2, "Loser", "Out")
	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 2, "evA", fredag, models.InterestLevelHigh)

	// Pin the player who DID want evA (top choice) into the single seat.
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES (?, ?, ?, 'Player', 'manual')`,
		"evA", string(fredag), 1,
	); err != nil {
		t.Fatalf("seed manual seat: %v", err)
	}

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	p := em.Puljer[0]
	// The pinned player wanted evA as a top choice, so pinning them there must
	// count as a newly-satisfied top choice and contribute their score.
	if p.NewlySatisfied != 1 {
		t.Errorf("pinned top-choice player should be newly satisfied: want 1, got %d", p.NewlySatisfied)
	}
	if em.SatisfiedTotal != 1 {
		t.Errorf("satisfied total: want 1, got %d", em.SatisfiedTotal)
	}
	if p.TotalScore < int(models.InterestLevelHigh.Score()) {
		t.Errorf("pinned top-choice player should contribute their score; total score too low: %d", p.TotalScore)
	}

	evA, _ := findEvent(p, "evA")
	for _, ap := range evA.AssignedPlayers {
		if ap.Name == "Pin Ned" {
			if !ap.Manual {
				t.Errorf("pinned player should be marked Manual")
			}
			if ap.Level != models.InterestLevelHigh {
				t.Errorf("pinned player's surfaced interest level: want High, got %q", ap.Level)
			}
		}
	}
}

func TestEmulateSeatings_ManualPinReservesSeatUnderContention(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "emulate_manual_pin_contention")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	// Single seat in evA; two highly-interested regulars compete for it.
	seedEvent(t, db, "evA", "Alpha", 1, fredag)
	seedParticipant(t, db, 1, "Pin", "Ned")
	seedParticipant(t, db, 2, "Eager", "One")
	seedParticipant(t, db, 3, "Eager", "Two")
	seedInterest(t, db, 2, "evA", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 3, "evA", fredag, models.InterestLevelHigh)

	// Admin pins Pin Ned into the single evA seat; the solver must honour it and
	// leave the eager regulars without that seat.
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES (?, ?, ?, 'Player', 'manual')`,
		"evA", string(fredag), 1,
	); err != nil {
		t.Fatalf("seed manual seat: %v", err)
	}

	em, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}

	evA, ok := findEvent(em.Puljer[0], "evA")
	if !ok {
		t.Fatal("evA missing from result")
	}
	names := playerNames(evA.AssignedPlayers)
	if !slices.Contains(names, "Pin Ned") {
		t.Fatalf("pinned player must keep the contested seat, got %v", names)
	}
	if len(evA.AssignedPlayers) != 1 {
		t.Fatalf("evA capacity is 1 and a pin fills it; want 1 seated, got %v", names)
	}
}
