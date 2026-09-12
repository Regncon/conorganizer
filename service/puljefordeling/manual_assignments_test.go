package puljefordeling

import (
	"errors"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestAddManualSeatFromInterest_PreservesInterestAndItsScoringMeaning(t *testing.T) {
	tests := []struct {
		name              string
		level             models.InterestLevel
		expectedSatisfied int
	}{
		{
			name:              "very interested counts as first choice",
			level:             models.InterestLevelHigh,
			expectedSatisfied: 1,
		},
		{
			name:              "ordinary interest remains an ordinary assignment",
			level:             models.InterestLevelMedium,
			expectedSatisfied: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			db, _ := testutil.CreateTestDBAndLogger(t, "add_manual_seat_from_interest")
			const fredag = models.PuljeFredagKveld
			seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
			seedEvent(t, db, "evA", "Alpha", 4, fredag)
			seedParticipant(t, db, 1, "Kari", "Nordmann")
			seedInterest(t, db, 1, "evA", fredag, tt.level)

			// When
			err := AddManualSeatFromInterest(db, fredag, "evA", 1)
			em, emulateErr := EmulateSeatings(db)

			// Then
			if err != nil {
				t.Fatalf("expected interest-backed assignment to succeed: %v", err)
			}
			if emulateErr != nil {
				t.Fatalf("emulate seating: %v", emulateErr)
			}
			if got := manualSeatCount(t, db, "evA", fredag, 1); got != 1 {
				t.Fatalf("manual seat count mismatch\nexpected: 1\nactual:   %d", got)
			}
			if got := interestCount(t, db, "evA", fredag, 1); got != 1 {
				t.Fatalf("the matching interest should remain\nexpected: 1\nactual:   %d", got)
			}
			if got := em.Puljer[0].NewlySatisfied; got != tt.expectedSatisfied {
				t.Fatalf("newly satisfied count mismatch\nexpected: %d\nactual:   %d", tt.expectedSatisfied, got)
			}

		})
	}
}

func TestAddManualSeatFromInterest_RejectsStaleInterest(t *testing.T) {
	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "add_manual_seat_from_missing_interest")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")

	// When
	err := AddManualSeatFromInterest(db, fredag, "evA", 1)

	// Then
	if !errors.Is(err, ErrInterestNotFound) {
		t.Fatalf("expected ErrInterestNotFound, got %v", err)
	}
	if got := manualSeatCount(t, db, "evA", fredag, 1); got != 0 {
		t.Fatalf("a stale interest must not create a seat, found %d", got)
	}
}

func TestAddManualGM_AssignsGMWithoutChangingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billetthelder med interesse for et arrangement i en pulje.",
		When:  "Når en administrator tildeler billetthelderen som spilleder.",
		Then:  "Så skal spilledertildelingen lagres manuelt og interessen beholdes.",
	})

	// Given
	expectedGMAssignments := 1
	expectedInterests := 1
	db, _ := testutil.CreateTestDBAndLogger(t, "add_manual_gm")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelMedium)

	// When
	err := AddManualGM(db, fredag, "evA", 1)
	actualGMAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
		  AND role = ? AND source = ?
	`, fredag, models.EventPlayerRoleGM, models.EventPlayerSourceManual)
	actualInterests := interestCount(t, db, "evA", fredag, 1)

	// Then
	if err != nil {
		t.Fatalf("expected GM assignment to succeed: %v", err)
	}
	if actualGMAssignments != expectedGMAssignments {
		t.Fatalf("GM assignment count mismatch\nexpected: %d\nactual:   %d", expectedGMAssignments, actualGMAssignments)
	}
	if actualInterests != expectedInterests {
		t.Fatalf("interest count mismatch\nexpected: %d\nactual:   %d", expectedInterests, actualInterests)
	}
}

func TestRemoveManualGM_RemovesOnlyMatchingGM(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en manuelt tildelt spilleder med en interesse i et annet arrangement.",
		When:  "Når administratoren fjerner spilledertildelingen.",
		Then:  "Så skal bare den samsvarende spilledertildelingen fjernes.",
	})

	// Given
	expectedGMAssignments := 0
	expectedOtherInterests := 1
	db, _ := testutil.CreateTestDBAndLogger(t, "remove_manual_gm")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedEvent(t, db, "evB", "Bravo", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evB", fredag, models.InterestLevelMedium)
	testutil.MustExec(t, db, `
		INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
		VALUES ('evA', ?, 1, ?, ?)
	`, fredag, models.EventPlayerRoleGM, models.EventPlayerSourceManual)

	// When
	err := RemoveManualGM(db, fredag, "evA", 1)
	actualGMAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1 AND role = ?
	`, fredag, models.EventPlayerRoleGM)
	actualOtherInterests := interestCount(t, db, "evB", fredag, 1)

	// Then
	if err != nil {
		t.Fatalf("expected GM removal to succeed: %v", err)
	}
	if actualGMAssignments != expectedGMAssignments {
		t.Fatalf("GM assignment count mismatch\nexpected: %d\nactual:   %d", expectedGMAssignments, actualGMAssignments)
	}
	if actualOtherInterests != expectedOtherInterests {
		t.Fatalf("other interest count mismatch\nexpected: %d\nactual:   %d", expectedOtherInterests, actualOtherInterests)
	}
}
