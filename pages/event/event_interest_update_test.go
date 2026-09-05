package event

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestUpdateInterest_WhenOpenRegistrationEvent_StoresFallbackInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt eit arrangement med open påmelding og ein billettheldar utan tildeling i pulja.",
		When:  "Når billettheldaren vel litt interessert som reserveval.",
		Then:  "Så skal reserveinteressa lagrast på vanleg måte.",
	})

	// Given
	expectedInterest := models.InterestLevelLow
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelMedium)
	mustExecEventInterestTest(t, db, `UPDATE events SET is_open_registration = 1 WHERE id = ?`, fixture.eventID)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		expectedInterest,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err != nil {
		t.Fatalf("expected fallback interest to succeed: %v", err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}

func TestUpdateInterest_WhenParticipationRuleBlocksChange_KeepsExistingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt ein billettheldar som ikkje kan endre vanleg interesse i pulja.",
		When:  "Når billettheldaren forsøker å endre interessa.",
		Then:  "Så skal forsøket avvisast og eksisterande interesse beholdast.",
	})

	testCases := []struct {
		name        string
		expectedErr error
		choice      models.InterestLevel
		mutate      func(t *testing.T, fixture eventInterestUpdateFixture, db *sql.DB)
	}{
		{
			name:        "high interest on open registration event",
			expectedErr: errInterestHighForOpenRegistration,
			choice:      models.InterestLevelHigh,
			mutate: func(t *testing.T, fixture eventInterestUpdateFixture, db *sql.DB) {
				mustExecEventInterestTest(t, db, `UPDATE events SET is_open_registration = 1 WHERE id = ?`, fixture.eventID)
			},
		},
		{
			name:        "underage billettholder on adults-only event",
			expectedErr: errInterestAdultsOnly,
			choice:      models.InterestLevelLow,
			mutate: func(t *testing.T, fixture eventInterestUpdateFixture, db *sql.DB) {
				mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupAdultsOnly)
				mustExecEventInterestTest(t, db, `UPDATE events SET age_group = ? WHERE id = ?`, models.AgeGroupAdultsOnly, fixture.eventID)
				mustExecEventInterestTest(t, db, `UPDATE billettholdere SET is_over_18 = 0 WHERE id = ?`, fixture.billettholderID)
			},
		},
		{
			name:        "player assignment in same pulje",
			expectedErr: errInterestAssignedInPulje,
			choice:      models.InterestLevelLow,
			mutate: func(t *testing.T, fixture eventInterestUpdateFixture, db *sql.DB) {
				seedEventInterestBlockingEvent(t, db, "assigned-event", "Assigned Event", fixture.puljeID)
				mustExecEventInterestTest(t, db, `
					INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
					VALUES (?, ?, ?, ?, ?)
				`, "assigned-event", fixture.puljeID, fixture.billettholderID, models.EventPlayerRolePlayer, models.EventPlayerSourceManual)
			},
		},
		{
			name:        "gamemaster assignment in same pulje",
			expectedErr: errInterestGamemasterInPulje,
			choice:      models.InterestLevelLow,
			mutate: func(t *testing.T, fixture eventInterestUpdateFixture, db *sql.DB) {
				seedEventInterestBlockingEvent(t, db, "gm-event", "GM Event", fixture.puljeID)
				mustExecEventInterestTest(t, db, `
					INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
					VALUES (?, ?, ?, ?, ?)
				`, "gm-event", fixture.puljeID, fixture.billettholderID, models.EventPlayerRoleGM, models.EventPlayerSourceManual)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Given
			expectedInterest := models.InterestLevelMedium
			db := createEventInterestTestDB(t)
			fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, expectedInterest)
			testCase.mutate(t, fixture, db)

			// When
			err := updateInterest(
				fixture.userExternalID,
				fixture.billettholderID,
				fixture.eventID,
				testCase.choice,
				string(fixture.puljeID),
				db,
			)
			actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

			// Then
			if !errors.Is(err, testCase.expectedErr) {
				t.Fatalf("error mismatch\nexpected: %v\nactual:   %v", testCase.expectedErr, err)
			}
			if actualInterest != expectedInterest {
				t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
			}
		})
	}
}

func TestUpdateInterest_WhenPuljeIsOpen_UpdatesInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder har meldt interesse i en åpen pulje.",
		When:  "Når interessen endres.",
		Then:  "Så skal den nye interessen lagres.",
	})

	// Given
	expectedInterest := models.InterestLevelLow

	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		expectedInterest,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err != nil {
		t.Fatalf("expected open pulje interest update to succeed: %v", err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}

func TestUpdateInterest_WhenProgramPublishingIsOff_RejectsInterestChangeAndKeepsExistingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at publisering av program er skrudd av.",
		When:  "Når interessen forsøkes endret.",
		Then:  "Så skal endringen avvises og eksisterende interesse beholdes.",
	})

	// Given
	expectedInterest := models.InterestLevelHigh
	expectedErrorText := "program"

	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, expectedInterest)
	setEventInterestProgramPublishing(t, db, false)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		models.InterestLevelLow,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err == nil {
		t.Errorf("expected unpublished program to reject interest update")
	} else if !strings.Contains(strings.ToLower(err.Error()), expectedErrorText) {
		t.Errorf("error mismatch\nexpected to contain: %q\nactual:              %v", expectedErrorText, err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}

func TestUpdateInterest_WhenEventIsNotPublishedInPulje_RejectsInterestChangeAndKeepsExistingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at arrangementet ikke er publisert i puljen.",
		When:  "Når interessen forsøkes endret.",
		Then:  "Så skal endringen avvises og eksisterende interesse beholdes.",
	})

	// Given
	expectedInterest := models.InterestLevelHigh
	expectedErrorText := "published"

	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, expectedInterest)
	mustExecEventInterestTest(t, db, `
		UPDATE relation_event_puljer
		SET is_published = 0
		WHERE event_id = ? AND pulje_id = ?
	`, fixture.eventID, fixture.puljeID)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		models.InterestLevelLow,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err == nil {
		t.Errorf("expected unpublished event pulje relation to reject interest update")
	} else if !strings.Contains(strings.ToLower(err.Error()), expectedErrorText) {
		t.Errorf("error mismatch\nexpected to contain: %q\nactual:              %v", expectedErrorText, err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}

func TestUpdateInterest_WhenPuljeIsLocked_RejectsInterestChangeAndKeepsExistingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder allerede har meldt interesse i en låst pulje.",
		When:  "Når interessen forsøkes endret.",
		Then:  "Så skal endringen avvises og eksisterende interesse beholdes.",
	})

	// Given
	expectedInterest := models.InterestLevelHigh
	expectedErrorText := "locked"

	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusLocked, expectedInterest)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		models.InterestLevelLow,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err == nil {
		t.Errorf("expected locked pulje to reject interest update")
	} else if !strings.Contains(strings.ToLower(err.Error()), expectedErrorText) {
		t.Errorf("error mismatch\nexpected to contain: %q\nactual:              %v", expectedErrorText, err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}

func TestUpdateInterest_WhenPuljeIsCompleted_RejectsInterestChangeAndKeepsExistingInterest(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder allerede har meldt interesse i en fullført pulje.",
		When:  "Når interessen forsøkes endret.",
		Then:  "Så skal endringen avvises og eksisterende interesse beholdes.",
	})

	// Given
	expectedInterest := models.InterestLevelHigh
	expectedErrorText := "completed"

	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusCompleted, expectedInterest)

	// When
	err := updateInterest(
		fixture.userExternalID,
		fixture.billettholderID,
		fixture.eventID,
		models.InterestLevelLow,
		string(fixture.puljeID),
		db,
	)
	actualInterest := getEventInterestTestInterest(t, db, fixture.eventID, fixture.billettholderID, fixture.puljeID)

	// Then
	if err == nil {
		t.Errorf("expected completed pulje to reject interest update")
	} else if !strings.Contains(strings.ToLower(err.Error()), expectedErrorText) {
		t.Errorf("error mismatch\nexpected to contain: %q\nactual:              %v", expectedErrorText, err)
	}
	if actualInterest != expectedInterest {
		t.Fatalf("interest level mismatch\nexpected: %s\nactual:   %s", expectedInterest, actualInterest)
	}
}
