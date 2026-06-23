package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestEmulatePulje_ReturnsOnlyRequestedPulje(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pulje_single")

	const fredag = models.PuljeFredagKveld
	const lordag = models.PuljeLordagMorgen
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedPulje(t, db, lordag, "Lørdag Morgen", "2026-09-05T10:00:00Z")
	seedEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedEvent(t, db, "evL", "Lørdagsspill", 4, lordag)
	seedParticipant(t, db, 1, "Anna", "A")
	seedInterest(t, db, 1, "evF", fredag, models.InterestLevelHigh)
	seedInterest(t, db, 1, "evL", lordag, models.InterestLevelHigh)

	pulje, meta, err := EmulatePulje(db, fredag)
	if err != nil {
		t.Fatalf("EmulatePulje: %v", err)
	}
	if pulje.PuljeID != fredag {
		t.Errorf("expected pulje %s, got %s", fredag, pulje.PuljeID)
	}
	if _, ok := findEvent(pulje, "evF"); !ok {
		t.Errorf("expected fredag event evF in result")
	}
	if _, ok := findEvent(pulje, "evL"); ok {
		t.Errorf("lørdag event evL must NOT appear in the fredag pulje result")
	}
	if meta.PlayerCount != 1 {
		t.Errorf("expected PlayerCount 1, got %d", meta.PlayerCount)
	}
}

func TestEmulatePulje_UnknownPuljeErrors(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pulje_unknown")
	seedPulje(t, db, models.PuljeFredagKveld, "Fredag Kveld", "2026-09-04T18:00:00Z")

	if _, _, err := EmulatePulje(db, models.Pulje("DoesNotExist")); err == nil {
		t.Fatal("expected error for unknown pulje, got nil")
	}
}

func TestEmulatePulje_MatchesFullEmulationForLastPulje(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pulje_equiv")

	const fredag = models.PuljeFredagKveld
	const lordag = models.PuljeLordagMorgen
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedPulje(t, db, lordag, "Lørdag Morgen", "2026-09-05T10:00:00Z")
	seedEvent(t, db, "evL", "Lørdagsspill", 4, lordag)
	seedParticipant(t, db, 1, "Anna", "A")
	seedInterest(t, db, 1, "evL", lordag, models.InterestLevelHigh)

	scoped, _, err := EmulatePulje(db, lordag)
	if err != nil {
		t.Fatalf("EmulatePulje: %v", err)
	}
	full, err := EmulateSeatings(db)
	if err != nil {
		t.Fatalf("EmulateSeatings: %v", err)
	}
	// lørdag is the last pulje chronologically -> index 1 in the full result.
	want := full.Puljer[1]
	if scoped.NewlySatisfied != want.NewlySatisfied || scoped.TotalScore != want.TotalScore {
		t.Errorf("scoped=%+v want=%+v", scoped, want)
	}
}

func TestEmulatePulje_AssignedPlayerHasBillettholderID(t *testing.T) {
	db, _ := testutil.CreateTestDBAndLogger(t, "test_emulate_pulje_bhid")

	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evF", "Fredagsspill", 4, fredag)
	seedParticipant(t, db, 7, "Anna", "A")
	seedInterest(t, db, 7, "evF", fredag, models.InterestLevelHigh)

	pulje, _, err := EmulatePulje(db, fredag)
	if err != nil {
		t.Fatalf("EmulatePulje: %v", err)
	}
	ev, ok := findEvent(pulje, "evF")
	if !ok {
		t.Fatal("evF missing")
	}
	if len(ev.AssignedPlayers) != 1 {
		t.Fatalf("expected 1 assigned player, got %d", len(ev.AssignedPlayers))
	}
	if ev.AssignedPlayers[0].BillettholderID != 7 {
		t.Errorf("expected BillettholderID 7, got %d", ev.AssignedPlayers[0].BillettholderID)
	}
}
