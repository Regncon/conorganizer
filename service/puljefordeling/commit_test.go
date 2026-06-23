package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestCommitPuljeAssignments_WritesSolverSeatsPreservesManual(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_commit")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)

	seedParticipant(t, db, 1, "Anna", "A")
	seedParticipant(t, db, 2, "Bob", "B")
	seedParticipant(t, db, 3, "Kid", "K")

	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 2, "evA", fredag, models.InterestLevelHigh)
	// Kid is manually pinned, no interest.
	seedAssignment(t, db, "evA", fredag, 3, "manual")

	if err := CommitPuljeAssignments(db, fredag, logger); err != nil {
		t.Fatalf("CommitPuljeAssignments: %v", err)
	}

	// Status flipped to Locked.
	var status string
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(fredag)).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != string(models.PuljeStatusLocked) {
		t.Errorf("status: want Locked, got %s", status)
	}

	// Anna & Bob written as solver; Kid preserved as manual.
	rows, err := db.Query(`SELECT billettholder_id, source FROM relation_events_players WHERE pulje_id = ? AND role = 'Player' ORDER BY billettholder_id`, string(fredag))
	if err != nil {
		t.Fatalf("query seats: %v", err)
	}
	defer rows.Close()
	type seat struct {
		bh     int
		source string
	}
	var seats []seat
	for rows.Next() {
		var s seat
		if err := rows.Scan(&s.bh, &s.source); err != nil {
			t.Fatalf("scan: %v", err)
		}
		seats = append(seats, s)
	}
	if len(seats) != 3 {
		t.Fatalf("want 3 Player rows, got %d (%v)", len(seats), seats)
	}
	bySource := map[int]string{}
	for _, s := range seats {
		bySource[s.bh] = s.source
	}
	if bySource[1] != "solver" || bySource[2] != "solver" {
		t.Errorf("Anna/Bob should be solver-written, got %v", bySource)
	}
	if bySource[3] != "manual" {
		t.Errorf("Kid should remain manual, got %q", bySource[3])
	}
}

func TestRevertPuljeAssignments_RemovesOnlySolverSeats(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_revert")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Anna", "A")
	seedParticipant(t, db, 2, "Kid", "K")
	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)
	seedAssignment(t, db, "evA", fredag, 2, "manual")

	if err := CommitPuljeAssignments(db, fredag, logger); err != nil {
		t.Fatalf("commit: %v", err)
	}
	if err := RevertPuljeAssignments(db, fredag); err != nil {
		t.Fatalf("revert: %v", err)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(fredag)).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != string(models.PuljeStatusOpen) {
		t.Errorf("status after revert: want Open, got %s", status)
	}

	var solverCount, manualCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND role='Player' AND source='solver'`, string(fredag)).Scan(&solverCount); err != nil {
		t.Fatalf("count solver: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND role='Player' AND source='manual'`, string(fredag)).Scan(&manualCount); err != nil {
		t.Fatalf("count manual: %v", err)
	}
	if solverCount != 0 {
		t.Errorf("solver seats should be gone, got %d", solverCount)
	}
	if manualCount != 1 {
		t.Errorf("manual seat must survive revert, got %d", manualCount)
	}
}

