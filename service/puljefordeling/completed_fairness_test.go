package puljefordeling

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestEmulateSeatings_CompletedPuljeUsesActualWinnerForLaterFairness(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "A completed fredag seated Anna, while Bob now runs a game on sondag.", When: "Later seating is emulated.", Then: "Anna keeps her actual seat and Bob wins the contested lordag seat."})
	// Given
	expectedFredag := []string{"Anna A"}
	expectedLordag := []string{"Bob B"}
	db, _ := testutil.CreateTestDBAndLogger(t, "completed_fairness")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedPulje(t, db, models.PuljeLordagMorgen, "Lordag Morgen", "2026-09-05T10:00:00Z")
	seedPulje(t, db, models.PuljeSondagMorgen, "Sondag Morgen", "2026-09-06T10:00:00Z")
	seedEvent(t, db, "fredag", "Fredag", 1, models.PuljeFredagKveld)
	seedEvent(t, db, "lordag", "Lordag", 1, models.PuljeLordagMorgen)
	seedEvent(t, db, "sondag", "Sondag", 1, models.PuljeSondagMorgen)
	seedParticipant(t, db, 1, "Anna", "A")
	seedParticipant(t, db, 2, "Bob", "B")
	for _, id := range []int{1, 2} {
		seedInterest(t, db, id, "fredag", models.PuljeFredagKveld, models.InterestLevelHigh)
		seedInterest(t, db, id, "lordag", models.PuljeLordagMorgen, models.InterestLevelHigh)
	}
	seedSolverSeat(t, db, "fredag", models.PuljeFredagKveld, 1)
	completePulje(t, db, models.PuljeFredagKveld)
	seedGM(t, db, "sondag", models.PuljeSondagMorgen, 2)
	// When
	em, err := EmulateSeatings(db)
	// Then
	if err != nil {
		t.Fatal(err)
	}
	fredag, _ := findEvent(em.Puljer[0], "fredag")
	lordag, _ := findEvent(em.Puljer[1], "lordag")
	if got := playerNames(fredag.AssignedPlayers); !slices.Equal(got, expectedFredag) {
		t.Errorf("completed fredag: want %v, got %v", expectedFredag, got)
	}
	if got := playerNames(lordag.AssignedPlayers); !slices.Equal(got, expectedLordag) {
		t.Errorf("later fairness: want %v, got %v", expectedLordag, got)
	}
	if em.SatisfiedTotal != 2 {
		t.Errorf("want both participants satisfied, got %d", em.SatisfiedTotal)
	}
}

func TestEmulateSeatings_CompletedPuljePreservesManualSeatsWithoutFillingVacancies(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "A completed pulje has a manual seat without interest and a vacant seat.", When: "Seating is emulated.", Then: "Only the persisted manual seat is displayed."})
	// Given
	expectedNames := []string{"Anna A"}
	db, _ := testutil.CreateTestDBAndLogger(t, "completed_manual")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "fredag", "Fredag", 2, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Anna", "A")
	seedParticipant(t, db, 2, "Bob", "B")
	seedManualSeat(t, db, "fredag", models.PuljeFredagKveld, 1)
	seedInterest(t, db, 2, "fredag", models.PuljeFredagKveld, models.InterestLevelHigh)
	completePulje(t, db, models.PuljeFredagKveld)
	// When
	em, err := EmulateSeatings(db)
	// Then
	if err != nil {
		t.Fatal(err)
	}
	ev, _ := findEvent(em.Puljer[0], "fredag")
	if got := playerNames(ev.AssignedPlayers); !slices.Equal(got, expectedNames) {
		t.Fatalf("want persisted seats %v, got %v", expectedNames, got)
	}
	if !ev.AssignedPlayers[0].Manual {
		t.Error("persisted manual seat lost its manual flag")
	}
	if em.SatisfiedTotal != 0 {
		t.Errorf("seat without interest must not satisfy a preference, got %d", em.SatisfiedTotal)
	}
}

func TestEmulateSeatings_CompletedEmptyPuljeDoesNotInventSeats(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "A completed pulje has interested participants but no persisted seats.", When: "Seating is emulated.", Then: "The completed pulje remains empty and nobody is satisfied."})
	// Given
	expectedUnassigned := []string{"Anna A"}
	db, _ := testutil.CreateTestDBAndLogger(t, "completed_empty")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "fredag", "Fredag", 1, models.PuljeFredagKveld)
	seedParticipant(t, db, 1, "Anna", "A")
	seedInterest(t, db, 1, "fredag", models.PuljeFredagKveld, models.InterestLevelHigh)
	completePulje(t, db, models.PuljeFredagKveld)
	// When
	em, err := EmulateSeatings(db)
	// Then
	if err != nil {
		t.Fatal(err)
	}
	ev, _ := findEvent(em.Puljer[0], "fredag")
	if len(ev.AssignedPlayers) != 0 {
		t.Errorf("completed empty pulje gained seats: %v", ev.AssignedPlayers)
	}
	if got := em.Puljer[0].Unassigned; !slices.Equal(got, expectedUnassigned) {
		t.Errorf("want unassigned %v, got %v", expectedUnassigned, got)
	}
	if em.SatisfiedTotal != 0 {
		t.Errorf("empty completed pulje must not satisfy anyone, got %d", em.SatisfiedTotal)
	}
}

func completePulje(t *testing.T, db *sql.DB, pulje models.Pulje) {
	t.Helper()
	if _, err := db.Exec(`UPDATE puljer SET status = ? WHERE id = ?`, models.PuljeStatusCompleted, pulje); err != nil {
		t.Fatal(err)
	}
}

func seedSolverSeat(t *testing.T, db *sql.DB, eventID string, pulje models.Pulje, id int) {
	t.Helper()
	if _, err := db.Exec(`INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES (?, ?, ?, ?, ?)`, eventID, pulje, id, models.EventPlayerRolePlayer, SourceSolver); err != nil {
		t.Fatal(err)
	}
}
