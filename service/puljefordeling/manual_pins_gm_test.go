package puljefordeling

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestAddManualSeat_AllowsGMInSamePuljeWithoutChangingGMSeat(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "a billettholder is GM in the pulje and has an existing Player seat", When: "an admin moves the Player pin", Then: "the Player pin moves and the GM assignment remains"})

	// Given
	const expectedGMAssignments = 1
	const expectedPlayerAssignments = 1
	db, _ := testutil.CreateTestDBAndLogger(t, "pin_gm_same_pulje")
	pulje := models.PuljeFredagKveld
	seedPulje(t, db, pulje, "Fredag Kveld", "2026-01-01 18:00")
	seedEvent(t, db, "evA", "A", 1, pulje)
	seedEvent(t, db, "evB", "B", 1, pulje)
	seedEvent(t, db, "evC", "C", 1, pulje)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedGM(t, db, "evB", pulje, 1)
	seedManualSeat(t, db, "evC", pulje, 1)

	// When
	err := AddManualSeat(db, pulje, "evA", 1)

	// Then
	if err != nil {
		t.Fatalf("expected same-pulje GM to be manually placed as Player: %v", err)
	}
	if got := assignmentCountByRole(t, db, pulje, 1, models.EventPlayerRoleGM); got != expectedGMAssignments {
		t.Fatalf("GM assignments: got %d, want %d", got, expectedGMAssignments)
	}
	if got := assignmentCountByRole(t, db, pulje, 1, models.EventPlayerRolePlayer); got != expectedPlayerAssignments {
		t.Fatalf("Player assignments: got %d, want %d", got, expectedPlayerAssignments)
	}
	if got := manualSeatCount(t, db, "evA", pulje, 1); got != 1 {
		t.Fatalf("target Player pin: got %d, want 1", got)
	}
}

func assignmentCountByRole(t *testing.T, db *sql.DB, pulje models.Pulje, billettholderID int, role models.EventPlayerRole) int {
	t.Helper()
	var count int
	if err := db.QueryRow(
		`SELECT COUNT(*) FROM relation_events_players WHERE pulje_id = ? AND billettholder_id = ? AND role = ?`,
		pulje, billettholderID, role,
	).Scan(&count); err != nil {
		t.Fatalf("count %s assignments: %v", role, err)
	}
	return count
}

func TestAddManualSeat_AllowsGMInDifferentPulje(t *testing.T) {
	// Given
	const expectedPins = 1
	db, _ := testutil.CreateTestDBAndLogger(t, "pin_gm_other_pulje")
	fredag := models.PuljeFredagKveld
	lordag := models.PuljeLordagMorgen
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-01-01 18:00")
	seedPulje(t, db, lordag, "Lordag Morgen", "2026-01-02 10:00")
	seedEvent(t, db, "evA", "A", 1, fredag)
	seedEvent(t, db, "evB", "B", 1, lordag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedGM(t, db, "evB", lordag, 1)
	// When
	err := AddManualSeat(db, fredag, "evA", 1)
	// Then
	if err != nil {
		t.Fatal(err)
	}
	if got := manualSeatCount(t, db, "evA", fredag, 1); got != expectedPins {
		t.Errorf("want %d pin, got %d", expectedPins, got)
	}
}
