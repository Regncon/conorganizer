package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestEmulateSeatings_SeatedPlayerCarriesSolverScore(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder med Veldig interessert på et arrangement med plass.",
		When:  "Når puljefordelingen emuleres.",
		Then:  "Så skal deltakeren ha algoritmepoengene fordelingen ga plassen.",
	})

	// Given
	expectedTotal := 810 // unsatisfied top choice (800) + never seated (10)
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_player_score")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "ev1", "Drager", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "ev1", fredag, models.InterestLevelHigh)

	// When
	em, err := EmulateSeatings(db)

	// Then
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}
	ev, _ := findEvent(em.Puljer[0], "ev1")
	player, ok := findAssigned(ev.AssignedPlayers, 1)
	if !ok {
		t.Fatal("expected Kari to be seated")
	}
	if player.Score == nil || player.Score.Total != expectedTotal {
		t.Fatalf("expected solver score %d, got %+v", expectedTotal, player.Score)
	}
}

func TestEmulateSeatings_PinnedPlayerHasNoSolverScore(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder som er manuelt plassert på et arrangement.",
		When:  "Når puljefordelingen emuleres.",
		Then:  "Så skal deltakeren ikke ha algoritmepoeng, fordi fordelingen ikke vurderte plassen.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pinned_no_score")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "ev1", "Drager", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "ev1", fredag, models.InterestLevelHigh)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source) VALUES ('ev1', ?, 1, 'Player', 'manual')`, string(fredag))

	// When
	em, err := EmulateSeatings(db)

	// Then
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}
	ev, _ := findEvent(em.Puljer[0], "ev1")
	player, ok := findAssigned(ev.AssignedPlayers, 1)
	if !ok || !player.Manual {
		t.Fatalf("expected Kari to be pinned, got %+v", player)
	}
	if player.Score != nil {
		t.Fatalf("expected no solver score for a pinned player, got %+v", player.Score)
	}
}
