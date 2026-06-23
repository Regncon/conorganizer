package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

// Given a Completed (published) pulje with interested players,
// When RerunPuljeAssignments re-solves it,
// Then it writes fresh solver seats and leaves the status Completed (no unpublish).
func TestRerunPuljeAssignments_ReSolvesAndPreservesCompletedStatus(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_rerun_service_completed")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedParticipant(t, db, 1, "Anna", "A")
	seedInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)
	setPuljeStatus(t, db, fredag, models.PuljeStatusCompleted)

	if err := RerunPuljeAssignments(db, fredag, logger); err != nil {
		t.Fatalf("RerunPuljeAssignments: %v", err)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(fredag)).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != string(models.PuljeStatusCompleted) {
		t.Errorf("expected status to stay Completed after rerun, got %s", status)
	}

	var solverSeats int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND source = 'solver'`,
		string(fredag),
	).Scan(&solverSeats); err != nil {
		t.Fatalf("count solver seats: %v", err)
	}
	if solverSeats == 0 {
		t.Errorf("expected rerun to write solver seats, found none")
	}
}

// Given a Locked pulje, When rerun runs, Then the status stays Locked.
func TestRerunPuljeAssignments_PreservesLockedStatus(t *testing.T) {
	db, logger := testutil.CreateTestDBAndLogger(t, "test_rerun_service_locked")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedParticipant(t, db, 1, "Anna", "A")
	seedInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)
	setPuljeStatus(t, db, fredag, models.PuljeStatusLocked)

	if err := RerunPuljeAssignments(db, fredag, logger); err != nil {
		t.Fatalf("RerunPuljeAssignments: %v", err)
	}

	var status string
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(fredag)).Scan(&status); err != nil {
		t.Fatalf("read status: %v", err)
	}
	if status != string(models.PuljeStatusLocked) {
		t.Errorf("expected status to stay Locked after rerun, got %s", status)
	}
}
