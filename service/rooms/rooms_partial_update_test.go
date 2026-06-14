package rooms

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestUpdateRoomPartial_UpdatesProvidedFields(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing room and partial input for every mutable field.",
		When:  "When the partial update runs.",
		Then:  "Then the returned room contains all supplied values.",
	})

	// Given
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	expectedRoom := models.Room{
		ID:                 existingRoom.ID,
		Name:               "Tangerud",
		RoomNumber:         "303",
		Floor:              3,
		MaxConcurrentGames: 3,
		Notes:              "",
		IsDisabled:         true,
	}
	input := models.RoomInput{
		ID:                 existingRoom.ID,
		Name:               ptr(expectedRoom.Name),
		RoomNumber:         ptr(expectedRoom.RoomNumber),
		Floor:              ptr(expectedRoom.Floor),
		MaxConcurrentGames: ptr(expectedRoom.MaxConcurrentGames),
		Notes:              ptr(expectedRoom.Notes),
		IsDisabled:         ptr(expectedRoom.IsDisabled),
	}

	// When
	actualRoom, err := UpdateRoomPartial(db, input)

	// Then
	if err != nil {
		t.Fatalf("expected partial room update to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestUpdateRoomPartial_WhenOnlyNameIsProvided_LeavesOtherFieldsUnchanged(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing room and partial input with only a new name.",
		When:  "When the partial update runs.",
		Then:  "Then only the name changes.",
	})

	// Given
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	expectedRoom := existingRoom
	expectedRoom.Name = "Tangerud"
	input := models.RoomInput{
		ID:   existingRoom.ID,
		Name: ptr(expectedRoom.Name),
	}

	// When
	actualRoom, err := UpdateRoomPartial(db, input)

	// Then
	if err != nil {
		t.Fatalf("expected partial room update to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestUpdateRoomPartial_WhenIDIsMissing_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given partial room input without a room ID.",
		When:  "When the partial update runs.",
		Then:  "Then validation rejects it.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_WhenNameIsEmpty_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given partial room input with an empty name.",
		When:  "When the partial update runs.",
		Then:  "Then validation rejects it.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{ID: 1, Name: ptr("")})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_WhenRoomNumberIsEmpty_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given partial room input with an empty room number.",
		When:  "When the partial update runs.",
		Then:  "Then validation rejects it.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{ID: 1, RoomNumber: ptr("")})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_WhenMaxConcurrentGamesIsInvalid_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given partial room input with invalid event capacity.",
		When:  "When the partial update runs.",
		Then:  "Then validation rejects it.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{ID: 1, MaxConcurrentGames: ptr(0)})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_WhenNoFieldsAreProvided_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given partial room input with an ID but no updated fields.",
		When:  "When the partial update runs.",
		Then:  "Then validation rejects it.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{ID: 1})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}
