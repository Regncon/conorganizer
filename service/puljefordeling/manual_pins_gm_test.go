package puljefordeling

import (
	"errors"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestAddManualSeat_RejectsGMInSamePuljeWithoutChangingSeats(t *testing.T) {
	for _, gmEvent := range []string{"evA", "evB"} {
		t.Run(gmEvent, func(t *testing.T) {
			bdd.Behavior(t, bdd.BDD{Given: "A billettholder is GM in the pulje and has an existing player seat.", When: "An admin tries to pin them into an event.", Then: "The pin is rejected and existing seats are preserved."})
			// Given
			const expectedSeats = 2
			db, _ := testutil.CreateTestDBAndLogger(t, "pin_gm_"+gmEvent)
			pulje := models.PuljeFredagKveld
			seedPulje(t, db, pulje, "Fredag Kveld", "2026-01-01 18:00")
			seedEvent(t, db, "evA", "A", 1, pulje)
			seedEvent(t, db, "evB", "B", 1, pulje)
			seedEvent(t, db, "evC", "C", 1, pulje)
			seedParticipant(t, db, 1, "Kari", "Nordmann")
			seedGM(t, db, gmEvent, pulje, 1)
			seedManualSeat(t, db, "evC", pulje, 1)
			// When
			err := AddManualSeat(db, pulje, "evA", 1)
			// Then
			if !errors.Is(err, ErrGMInPulje) {
				t.Error("expected GM pin to be rejected")
			}
			if got := manualSeatCount(t, db, "evA", pulje, 1); got != 0 {
				t.Errorf("conflicting pin persisted: %d", got)
			}
			if got := manualSeatCount(t, db, "evC", pulje, 1); got != 1 {
				t.Errorf("prior seat changed: %d", got)
			}
			var seats int
			if err := db.QueryRow(`SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id = 1`).Scan(&seats); err != nil {
				t.Fatal(err)
			}
			if seats != expectedSeats {
				t.Errorf("want %d unchanged seats, got %d", expectedSeats, seats)
			}
		})
	}
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
