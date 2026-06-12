package rooms

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
)

func TestUpdateRoom_UpdatesAllFieldsWithoutChangingID(t *testing.T) {
	// Given an existing room and replacement room data,
	// when the room is updated,
	// then every mutable field changes while the ID stays the same.

	// Given
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	expectedRoom := models.Room{
		ID:                 existingRoom.ID,
		Name:               "Tangerud",
		RoomNumber:         "209",
		Floor:              2,
		MaxConcurrentGames: 3,
		Notes:              "Dette er ei oppdatert note",
		IsDisabled:         true,
	}

	// When
	actualRoom, err := UpdateRoom(db, expectedRoom)

	// Then
	if err != nil {
		t.Fatalf("expected room update to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestUpdateRoom_WhenRoomNumberIsEmpty_ReturnsError(t *testing.T) {
	// Given replacement room data with an empty room number,
	// when the room is updated,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	invalidRoom := existingRoom
	invalidRoom.RoomNumber = ""

	// When
	_, err := UpdateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoom_WhenRoomNumberDoesNotMatchFloor_ReturnsError(t *testing.T) {
	// Given replacement room data whose number starts with another floor,
	// when the room is updated,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	invalidRoom := existingRoom
	invalidRoom.RoomNumber = "203"
	invalidRoom.Floor = 3

	// When
	_, err := UpdateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoom_WhenMaxConcurrentGamesIsInvalid_ReturnsError(t *testing.T) {
	// Given replacement room data with invalid event capacity,
	// when the room is updated,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	invalidRoom := existingRoom
	invalidRoom.MaxConcurrentGames = 0

	// When
	_, err := UpdateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}
