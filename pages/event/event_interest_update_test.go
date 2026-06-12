package event

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
)

func TestUpdateInterest_WhenPuljeIsOpen_UpdatesInterest(t *testing.T) {
	// Gitt at en billettholder har meldt interesse i en åpen pulje,
	// når interessen endres,
	// så skal den nye interessen lagres.

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
	// Gitt at publisering av program er skrudd av,
	// når interessen forsøkes endret,
	// så skal endringen avvises og eksisterende interesse beholdes.

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
	// Gitt at arrangementet ikke er publisert i puljen,
	// når interessen forsøkes endret,
	// så skal endringen avvises og eksisterende interesse beholdes.

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
	// Gitt at en billettholder allerede har meldt interesse i en låst pulje,
	// når interessen forsøkes endret,
	// så skal endringen avvises og eksisterende interesse beholdes.

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
	// Gitt at en billettholder allerede har meldt interesse i en fullført pulje,
	// når interessen forsøkes endret,
	// så skal endringen avvises og eksisterende interesse beholdes.

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
