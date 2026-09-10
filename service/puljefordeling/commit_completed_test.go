package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestCommitDistribution_RejectsCompletedPuljeWithoutChangingSeats(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en publisert pulje med lagrede plasser og endrede interesser.",
		When:  "Når fordelingen forsøkes lagret på nytt.",
		Then:  "Så avvises lagringen og plassene bevares.",
	})
	// Given
	expectedEvent := "evA"
	db, _ := testutil.CreateTestDBAndLogger(t, "commit_completed")
	const pulje = models.PuljeFredagKveld
	seedPulje(t, db, pulje, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, pulje)
	seedEvent(t, db, "evB", "Beta", 4, pulje)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evB", pulje, models.InterestLevelHigh)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES ('evA', ?, 1, 'Player', 'solver')`, pulje)
	testutil.MustExec(t, db, `UPDATE puljer SET status = 'Completed' WHERE id = ?`, pulje)

	// When
	err := CommitDistribution(db, pulje)

	// Then
	if err == nil {
		t.Error("expected completed pulje commit to be rejected")
	}
	var eventID string
	if err := db.QueryRow(`SELECT event_id FROM relation_events_players WHERE billettholder_id = 1 AND pulje_id = ?`, pulje).Scan(&eventID); err != nil {
		t.Fatal(err)
	}
	if eventID != expectedEvent {
		t.Errorf("published seat changed: got %s, want %s", eventID, expectedEvent)
	}
}

func TestCommitDistribution_AllowsSavingAfterUnpublishing(t *testing.T) {
	// Given
	expectedEvent := "evB"
	db, _ := testutil.CreateTestDBAndLogger(t, "commit_unpublished")
	const pulje = models.PuljeFredagKveld
	seedPulje(t, db, pulje, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evB", "Beta", 4, pulje)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evB", pulje, models.InterestLevelHigh)
	testutil.MustExec(t, db, `UPDATE puljer SET status = 'Completed' WHERE id = ?`, pulje)
	testutil.MustExec(t, db, `UPDATE puljer SET status = 'Locked' WHERE id = ?`, pulje)

	// When
	err := CommitDistribution(db, pulje)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	var eventID string
	if err := db.QueryRow(`SELECT event_id FROM relation_events_players WHERE billettholder_id = 1 AND pulje_id = ?`, pulje).Scan(&eventID); err != nil {
		t.Fatal(err)
	}
	if eventID != expectedEvent {
		t.Errorf("got seat %s, want %s", eventID, expectedEvent)
	}
}
