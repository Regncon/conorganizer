package puljefordeling

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestEmulateSeatings_AllGMsExcludedFromTheirPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an event has two GMs who also express player interests", When: "the pulje is emulated", Then: "both GMs remain excluded from player seats"})
	// Given
	expectedPlayers := []string{"Player Regular"}
	db, _ := testutil.CreateTestDBAndLogger(t, "all_gms_excluded")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "run", "Run", 4, fredag)
	seedEvent(t, db, "play", "Play", 4, fredag)
	for id, name := range map[int]string{1: "First", 2: "Second", 3: "Player"} {
		seedParticipant(t, db, id, name, "Regular")
		seedInterest(t, db, id, "play", fredag, models.InterestLevelHigh)
	}
	seedGM(t, db, "run", fredag, 1)
	seedGM(t, db, "run", fredag, 2)

	// When
	em, err := EmulateSeatings(db)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	ev, _ := findEvent(em.Puljer[0], "play")
	if got := playerNames(ev.AssignedPlayers); !slices.Equal(got, expectedPlayers) {
		t.Fatalf("players: got %v, want %v; all GMs are busy", got, expectedPlayers)
	}
	if len(em.Puljer[0].Unassigned) != 0 {
		t.Fatalf("busy GMs listed as unassigned: %v", em.Puljer[0].Unassigned)
	}
}

func TestEmulateSeatings_EachGMReceivesWeekendPriority(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "two GMs share an event and later contend for a seat", When: "the weekend is emulated", Then: "each GM receives priority over an equally interested regular player"})
	for _, gmID := range []int{1, 2} {
		t.Run(fmt.Sprintf("GM%d", gmID), func(t *testing.T) {
			// Given
			expectedWinner := gmID
			db, _ := testutil.CreateTestDBAndLogger(t, "all_gms_priority")
			const fredag, lordag = models.PuljeFredagKveld, models.PuljeLordagMorgen
			seedPulje(t, db, fredag, "Fredag", "2027-09-04T18:00:00Z")
			seedPulje(t, db, lordag, "Lørdag", "2027-09-05T10:00:00Z")
			seedEvent(t, db, "run", "Run", 4, fredag)
			seedEvent(t, db, "play", "Play", 1, lordag)
			for id := 1; id <= 3; id++ {
				seedParticipant(t, db, id, fmt.Sprint(id), "Person")
			}
			seedGM(t, db, "run", fredag, 1)
			seedGM(t, db, "run", fredag, 2)
			seedInterest(t, db, gmID, "play", lordag, models.InterestLevelHigh)
			seedInterest(t, db, 3, "play", lordag, models.InterestLevelHigh)

			// When
			em, err := EmulateSeatings(db)

			// Then
			if err != nil {
				t.Fatal(err)
			}
			ev, _ := findEvent(em.Puljer[1], "play")
			if len(ev.AssignedPlayers) != 1 || ev.AssignedPlayers[0].BillettholderID != expectedWinner {
				t.Fatalf("winner: got %+v, want GM %d", ev.AssignedPlayers, expectedWinner)
			}
			if !ev.AssignedPlayers[0].IsDM {
				t.Fatal("GM winner missing weekend GM flag")
			}
		})
	}
}

func TestEmulateSeatings_DisplaysAllGMsAndPreservesMinorWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{Given: "an event has an adult GM and a minor GM", When: "the emulation is displayed", Then: "both names appear and the minor warning remains available"})
	// Given
	expectedNames := "Alice GM, Bob GM"
	db, _ := testutil.CreateTestDBAndLogger(t, "all_gms_display")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "run", "Run", 4, fredag)
	seedParticipant(t, db, 1, "Alice", "GM")
	seedParticipant(t, db, 2, "Bob", "GM")
	markParticipantOver18(t, db, 2)
	seedGM(t, db, "run", fredag, 1)
	seedGM(t, db, "run", fredag, 2)

	// When
	em, err := EmulateSeatings(db)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	ev, _ := findEvent(em.Puljer[0], "run")
	if ev.GMName != expectedNames {
		t.Errorf("GM names: got %q, want %q", ev.GMName, expectedNames)
	}
	if ev.GMIsOver18 {
		t.Error("a minor GM must keep the age warning visible")
	}
}
