package formsubmission

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestGetAssigneesForEvent_AggregatesPlayerAndGMFlags(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given one billettholder assigned as both Player and GM to the same event and pulje.",
		When:  "When assignees are fetched for the event.",
		Then:  "Then one assignee row reports both role flags.",
	})

	// Given
	const expectedAssignees = 1
	db, logger := testutil.CreateTestDBAndLogger(t, "assignees-dual-role")
	seedBaseTables(t, db)
	seedBillettholdere(t, db, []billettholderFixture{{id: 20, firstName: "Dual", lastName: "Role"}})
	seedInterests(t, db, []interestFixture{{billettholderID: 20, eventID: eventE1, puljeID: puljeP1, interestLevel: models.InterestLevelHigh}})
	seedRoleAssignment(t, db, 20, eventE1, puljeP1, models.EventPlayerRolePlayer, "manual")
	seedRoleAssignment(t, db, 20, eventE1, puljeP1, models.EventPlayerRoleGM, "manual")

	// When
	assignees, err := GetAssigneesForEvent(eventE1, db, logger)

	// Then
	if err != nil {
		t.Fatalf("get assignees: %v", err)
	}
	if len(assignees) != expectedAssignees {
		t.Fatalf("assignee count = %d, want %d", len(assignees), expectedAssignees)
	}
	if !assignees[0].IsPlayer || !assignees[0].IsGamemaster {
		t.Fatalf("role flags = (Player=%t, GM=%t), want both true", assignees[0].IsPlayer, assignees[0].IsGamemaster)
	}
}

func TestGetInterestsForEvent_UnassignedInterestHasNoRoleFlags(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given one unassigned billettholder interested in an event.",
		When:  "When interests are fetched for the event.",
		Then:  "Then the interest reports both Player and GM flags as false.",
	})

	// Given
	const expectedInterests = 1
	db, logger := testutil.CreateTestDBAndLogger(t, "interests-no-role")
	seedBaseTables(t, db)
	seedBillettholdere(t, db, []billettholderFixture{{id: 21, firstName: "No", lastName: "Role"}})
	seedInterests(t, db, []interestFixture{{billettholderID: 21, eventID: eventE1, puljeID: puljeP1, interestLevel: models.InterestLevelMedium}})

	// When
	interests, err := GetInterestsForEvent(eventE1, db, logger)

	// Then
	if err != nil {
		t.Fatalf("get interests: %v", err)
	}
	if len(interests) != expectedInterests {
		t.Fatalf("interest count = %d, want %d", len(interests), expectedInterests)
	}
	if interests[0].IsPlayer || interests[0].IsGamemaster {
		t.Fatalf("role flags = (Player=%t, GM=%t), want both false", interests[0].IsPlayer, interests[0].IsGamemaster)
	}
}

func seedRoleAssignment(t *testing.T, db *sql.DB, billettholderID int, eventID, puljeID string, role models.EventPlayerRole, source string) {
	t.Helper()
	if _, err := db.Exec(`
		INSERT INTO relation_events_players (billettholder_id, event_id, pulje_id, role, source)
		VALUES (?, ?, ?, ?, ?)`, billettholderID, eventID, puljeID, role, source); err != nil {
		t.Fatalf("seed %s assignment: %v", role, err)
	}
}
