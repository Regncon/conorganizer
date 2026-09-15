package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestRelationEventsPlayers_AllowsIndependentRolesForSameAssignment(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a billettholder assigned as GM to an event in a pulje.",
		When:  "When the same billettholder is also assigned as Player to that event.",
		Then:  "Then both roles are stored independently and either role can be removed alone.",
	})

	// Given
	const expectedRoleRows = 2
	db, _ := testutil.CreateTestDBAndLogger(t, "test_schema_independent_roles")
	seedPulje(t, db, models.PuljeFredagKveld, "Pulje 1", "2026-09-15T09:00:00Z")
	seedEvent(t, db, "event-1", "Event 1", 4, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Ada", "Lovelace")
	if _, err := db.Exec(`
		INSERT INTO relation_events_players (billettholder_id, event_id, pulje_id, role)
		VALUES (1, 'event-1', ?, 'GM')`, models.PuljeFredagKveld); err != nil {
		t.Fatalf("seed GM assignment: %v", err)
	}

	// When
	_, err := db.Exec(`
		INSERT INTO relation_events_players (billettholder_id, event_id, pulje_id, role)
		VALUES (1, 'event-1', ?, 'Player')`, models.PuljeFredagKveld)

	// Then
	if err != nil {
		t.Fatalf("insert Player beside GM: %v", err)
	}
	var roleRows int
	if err := db.QueryRow(`
		SELECT COUNT(*) FROM relation_events_players
		WHERE billettholder_id = 1 AND event_id = 'event-1' AND pulje_id = ?`, models.PuljeFredagKveld).Scan(&roleRows); err != nil {
		t.Fatalf("count role rows: %v", err)
	}
	if roleRows != expectedRoleRows {
		t.Fatalf("role row count = %d, want %d", roleRows, expectedRoleRows)
	}
	if _, err := db.Exec(`
		INSERT INTO relation_events_players (billettholder_id, event_id, pulje_id, role)
		VALUES (1, 'event-1', ?, 'Player')`, models.PuljeFredagKveld); err == nil {
		t.Fatal("duplicate Player role insert succeeded, want primary-key violation")
	}
	if _, err := db.Exec(`
		DELETE FROM relation_events_players
		WHERE billettholder_id = 1 AND event_id = 'event-1' AND pulje_id = ? AND role = 'Player'`, models.PuljeFredagKveld); err != nil {
		t.Fatalf("delete Player role: %v", err)
	}
	var remainingRole models.EventPlayerRole
	if err := db.QueryRow(`
		SELECT role FROM relation_events_players
		WHERE billettholder_id = 1 AND event_id = 'event-1' AND pulje_id = ?`, models.PuljeFredagKveld).Scan(&remainingRole); err != nil {
		t.Fatalf("read remaining role: %v", err)
	}
	if remainingRole != models.EventPlayerRoleGM {
		t.Fatalf("remaining role = %q, want %q", remainingRole, models.EventPlayerRoleGM)
	}
}
