package puljefordeling

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

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

func TestAddFirstChoiceSeat_AssignsPlayerAndRecordsHighInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billetthelder uten interesse for et arrangement.",
		When:  "Når en administrator tildeler arrangementet som førstevalg.",
		Then:  "Så skal en manuell spillerplass og veldig interessert lagres sammen.",
	})

	// Given
	expectedPlayerAssignments := 1
	expectedInterestLevel := models.InterestLevelHigh
	db, _ := testutil.CreateTestDBAndLogger(t, "add_first_choice_seat")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")

	// When
	err := AddFirstChoiceSeat(db, fredag, "evA", 1)
	actualPlayerAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
		  AND role = ? AND source = ?
	`, fredag, models.EventPlayerRolePlayer, models.EventPlayerSourceManual)
	var actualInterestLevel models.InterestLevel
	if queryErr := db.QueryRow(`
		SELECT interest_level FROM interests
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
	`, fredag).Scan(&actualInterestLevel); queryErr != nil {
		t.Fatalf("read first-choice interest: %v", queryErr)
	}

	// Then
	if err != nil {
		t.Fatalf("expected first-choice assignment to succeed: %v", err)
	}
	if actualPlayerAssignments != expectedPlayerAssignments {
		t.Fatalf("player assignment count mismatch\nexpected: %d\nactual:   %d", expectedPlayerAssignments, actualPlayerAssignments)
	}
	if actualInterestLevel != expectedInterestLevel {
		t.Fatalf("interest level mismatch\nexpected: %q\nactual:   %q", expectedInterestLevel, actualInterestLevel)
	}
}

func TestRemoveFirstChoiceSeat_RemovesMatchingSeatAndInterestOnly(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en administrativt tildelt førstevalg og en annen interesse i samme pulje.",
		When:  "Når administratoren fjerner førstevalget.",
		Then:  "Så skal førstevalgets spillerplass og interesse fjernes mens den andre interessen beholdes.",
	})

	// Given
	expectedPlayerAssignments := 0
	expectedMatchingInterests := 0
	expectedOtherInterests := 1
	db, _ := testutil.CreateTestDBAndLogger(t, "remove_first_choice_seat")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedEvent(t, db, "evB", "Bravo", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evB", fredag, models.InterestLevelMedium)
	if err := AddFirstChoiceSeat(db, fredag, "evA", 1); err != nil {
		t.Fatalf("seed first-choice assignment: %v", err)
	}

	// When
	err := RemoveFirstChoiceSeat(db, fredag, "evA", 1)
	actualPlayerAssignments := testutil.QueryInt(t, db, `
		SELECT COUNT(*) FROM relation_events_players
		WHERE event_id = 'evA' AND pulje_id = ? AND billettholder_id = 1
	`, fredag)
	actualMatchingInterests := interestCount(t, db, "evA", fredag, 1)
	actualOtherInterests := interestCount(t, db, "evB", fredag, 1)

	// Then
	if err != nil {
		t.Fatalf("expected first-choice removal to succeed: %v", err)
	}
	if actualPlayerAssignments != expectedPlayerAssignments {
		t.Fatalf("player assignment count mismatch\nexpected: %d\nactual:   %d", expectedPlayerAssignments, actualPlayerAssignments)
	}
	if actualMatchingInterests != expectedMatchingInterests {
		t.Fatalf("matching interest count mismatch\nexpected: %d\nactual:   %d", expectedMatchingInterests, actualMatchingInterests)
	}
	if actualOtherInterests != expectedOtherInterests {
		t.Fatalf("other interest count mismatch\nexpected: %d\nactual:   %d", expectedOtherInterests, actualOtherInterests)
	}
}

func TestRemoveFirstChoiceSeat_LeavesOrdinaryHighInterestWithoutSeat(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et vanlig veldig interessert-valg uten en manuell spillerplass.",
		When:  "Når en foreldet fjernhandling for førstevalg blir sendt.",
		Then:  "Så skal den vanlige interessen beholdes.",
	})

	// Given
	expectedInterests := 1
	db, _ := testutil.CreateTestDBAndLogger(t, "remove_stale_first_choice")
	const fredag = models.PuljeFredagKveld
	seedPulje(t, db, fredag, "Fredag Kveld", "2026-09-04T18:00:00Z")
	seedEvent(t, db, "evA", "Alpha", 4, fredag)
	seedParticipant(t, db, 1, "Kari", "Nordmann")
	seedInterest(t, db, 1, "evA", fredag, models.InterestLevelHigh)

	// When
	err := RemoveFirstChoiceSeat(db, fredag, "evA", 1)
	actualInterests := interestCount(t, db, "evA", fredag, 1)

	// Then
	if err != nil {
		t.Fatalf("expected stale first-choice removal to be a no-op: %v", err)
	}
	if actualInterests != expectedInterests {
		t.Fatalf("interest count mismatch\nexpected: %d\nactual:   %d", expectedInterests, actualInterests)
	}
}
