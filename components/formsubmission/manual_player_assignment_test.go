package formsubmission

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestUpdatePlayerStatus_WhenAssigningPlayer_RemovesOnlyMatchingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billetthelder med interesser i to arrangementer i samme pulje.",
		When:  "Når administratoren tildeler billetthelderen som spiller i ett arrangement.",
		Then:  "Så skal spillertildelingen lagres og bare den samsvarende interessen fjernes.",
	})

	// Given
	expectedAssignments := 1
	expectedMatchingInterests := 0
	expectedOtherInterests := 1
	db := testutil.CreateTestDB(t, "approval_manual_player_assignment")
	seedManualPlayerAssignmentFixture(t, db)

	// When
	err := UpdatePlayerStatus("event-a", string(models.PuljeFredagKveld), 1, true, false, db)
	actualAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'event-a' AND pulje_id = ? AND billettholder_id = 1
			AND role = ? AND source = ?
	`, models.PuljeFredagKveld, models.EventPlayerRolePlayer, models.EventPlayerSourceManual)
	actualMatchingInterests := manualPlayerInterestCount(t, db, "event-a")
	actualOtherInterests := manualPlayerInterestCount(t, db, "event-b")

	// Then
	if err != nil {
		t.Fatalf("expected manual player assignment to succeed: %v", err)
	}
	if actualAssignments != expectedAssignments {
		t.Fatalf("assignment count mismatch\nexpected: %d\nactual:   %d", expectedAssignments, actualAssignments)
	}
	if actualMatchingInterests != expectedMatchingInterests {
		t.Fatalf("matching interest count mismatch\nexpected: %d\nactual:   %d", expectedMatchingInterests, actualMatchingInterests)
	}
	if actualOtherInterests != expectedOtherInterests {
		t.Fatalf("other interest count mismatch\nexpected: %d\nactual:   %d", expectedOtherInterests, actualOtherInterests)
	}
}

func seedManualPlayerAssignmentFixture(t *testing.T, db *sql.DB) {
	t.Helper()

	testutil.MustExec(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium)
	testutil.MustExec(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at)
		VALUES (?, 'Fredag kveld', ?, '2026-09-04T18:00:00Z', '2026-09-04T23:00:00Z')
	`, models.PuljeFredagKveld, models.PuljeStatusOpen)
	testutil.MustExec(t, db, `
		INSERT INTO events (id, title, intro, description, host_name, email, phone_number, max_players)
		VALUES
			('event-a', 'Alpha', '', '', '', '', '', 4),
			('event-b', 'Bravo', '', '', '', '', '', 4)
	`)
	testutil.MustExec(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (1, 'Kari', 'Nordmann', 1, 'Ticket', 1, 1001, 2001)
	`)
	testutil.MustExec(t, db, `
		INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
		VALUES
			(1, 'event-a', ?, ?),
			(1, 'event-b', ?, ?)
	`, models.PuljeFredagKveld, models.InterestLevelHigh, models.PuljeFredagKveld, models.InterestLevelMedium)
}

func manualPlayerInterestCount(t *testing.T, db *sql.DB, eventID string) int {
	t.Helper()
	return testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM interests
		WHERE billettholder_id = 1 AND event_id = ? AND pulje_id = ?
	`, eventID, models.PuljeFredagKveld)
}
