package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestCommitDistribution_PersistsSolverPicksAsSolverSource(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "commit_solver")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)

	if err := CommitDistribution(db, fredag); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}

	var source string
	if err := db.QueryRow(
		`SELECT source FROM relation_events_players WHERE event_id='evA' AND pulje_id=? AND billettholder_id=1 AND role='Player'`,
		string(fredag),
	).Scan(&source); err != nil {
		t.Fatalf("expected a committed solver seat: %v", err)
	}
	if source != "solver" {
		t.Errorf("solver pick should be committed as source='solver', got %q", source)
	}
}

func TestCommitDistribution_KeepsManualSeatsManual(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "commit_manual")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 2, "Pin", "Ned")
	// Manual pin, no interest.
	if _, err := db.Exec(
		`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		 VALUES ('evA', ?, 2, 'Player', 'manual')`, string(fredag),
	); err != nil {
		t.Fatalf("seed manual seat: %v", err)
	}

	if err := CommitDistribution(db, fredag); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}

	var source string
	if err := db.QueryRow(
		`SELECT source FROM relation_events_players WHERE event_id='evA' AND pulje_id=? AND billettholder_id=2`,
		string(fredag),
	).Scan(&source); err != nil {
		t.Fatalf("manual seat should survive commit: %v", err)
	}
	if source != "manual" {
		t.Errorf("manual seat must stay source='manual', got %q", source)
	}
}

func TestCommitDistribution_PreservesRegistrationAndDoesNotAddOrdinarySeat(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en påmeldt billetthelder med en ordinær interesse i samme pulje.",
		When:  "Når puljefordelingen lagres.",
		Then:  "Så skal påmeldingskilden beholdes uten at billetthelderen får en solverplass i puljen.",
	})

	db, _ := testutil.CreateTestDBAndLogger(t, "commit_registration")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evOpen", "Open", 100, fredag)
	seedEvent(t, db, "evOrdinary", "Ordinary", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedRegistration(t, db, 1, "evOpen", fredag)
	seedInterest(t, db, 1, "evOrdinary", fredag, models.InterestLevelHigh)

	if err := CommitDistribution(db, fredag); err != nil {
		t.Fatalf("CommitDistribution: %v", err)
	}

	registrationCount := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_events_players
		WHERE event_id = 'evOpen' AND pulje_id = ? AND billettholder_id = 1
		  AND role = ? AND source = ?
	`, fredag, models.EventPlayerRolePlayer, models.EventPlayerSourceRegistration)
	if registrationCount != 1 {
		t.Fatalf("registration should keep its source, got %d matching rows", registrationCount)
	}

	solverSeatCount := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_events_players
		WHERE pulje_id = ? AND billettholder_id = 1
		  AND role = ? AND source = ?
	`, fredag, models.EventPlayerRolePlayer, models.EventPlayerSourceSolver)
	if solverSeatCount != 0 {
		t.Fatalf("registered participant must not receive an ordinary solver seat in the same pulje, got %d", solverSeatCount)
	}
}
