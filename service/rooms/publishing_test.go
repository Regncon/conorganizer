package rooms

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestSetRoomsPublished_RoundTripsThroughGet(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a pulje with an unpublished room assignment.",
		When:  "When the room assignment is published.",
		Then:  "Then reading it back reports it as published.",
	})

	// Given
	expectedPublished := true
	db := createRoomsTestDB(t)
	seedRoomEventLookups(t, db)
	pulje := insertPulje(t, db, models.Pulje("Friday"), "Fredag kveld")

	// When
	err := SetRoomsPublished(db, pulje, expectedPublished)
	if err != nil {
		t.Fatalf("expected setting rooms published to succeed: %v", err)
	}
	actualPublished, err := RoomsPublished(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("expected reading rooms published to succeed: %v", err)
	}
	if actualPublished != expectedPublished {
		t.Fatalf("published mismatch\nexpected: %v\nactual:   %v", expectedPublished, actualPublished)
	}
}

func TestSetRoomsPublished_CanUnpublishAgain(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a pulje with a published room assignment.",
		When:  "When the room assignment is unpublished.",
		Then:  "Then reading it back reports it as unpublished.",
	})

	// Given
	expectedPublished := false
	db := createRoomsTestDB(t)
	seedRoomEventLookups(t, db)
	pulje := insertPulje(t, db, models.Pulje("Friday"), "Fredag kveld")
	if err := SetRoomsPublished(db, pulje, true); err != nil {
		t.Fatalf("expected setting rooms published to succeed: %v", err)
	}

	// When
	err := SetRoomsPublished(db, pulje, expectedPublished)
	if err != nil {
		t.Fatalf("expected setting rooms published to succeed: %v", err)
	}
	actualPublished, err := RoomsPublished(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("expected reading rooms published to succeed: %v", err)
	}
	if actualPublished != expectedPublished {
		t.Fatalf("published mismatch\nexpected: %v\nactual:   %v", expectedPublished, actualPublished)
	}
}

func TestRoomsPublished_DefaultsToUnpublished(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a newly created pulje.",
		When:  "When its rooms published state is read.",
		Then:  "Then it defaults to unpublished.",
	})

	// Given
	expectedPublished := false
	db := createRoomsTestDB(t)
	seedRoomEventLookups(t, db)
	pulje := insertPulje(t, db, models.Pulje("Friday"), "Fredag kveld")

	// When
	actualPublished, err := RoomsPublished(db, pulje)

	// Then
	if err != nil {
		t.Fatalf("expected reading rooms published to succeed: %v", err)
	}
	if actualPublished != expectedPublished {
		t.Fatalf("published mismatch\nexpected: %v\nactual:   %v", expectedPublished, actualPublished)
	}
}

func TestRoomsPublished_WhenPuljeDoesNotExist_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given no pulje with the given ID.",
		When:  "When its rooms published state is read.",
		Then:  "Then the caller receives an error.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := RoomsPublished(db, models.Pulje("missing-pulje"))
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestSetRoomsPublished_WhenPuljeDoesNotExist_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given no pulje with the given ID.",
		When:  "When its rooms published state is set.",
		Then:  "Then the caller receives an error.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	err := SetRoomsPublished(db, models.Pulje("missing-pulje"), true)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}
