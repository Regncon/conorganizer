package rooms

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestUpdateRoom_UpdatesAllFieldsWithoutChangingID(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an existing room and replacement room data.",
		When:  "When the room is updated.",
		Then:  "Then every mutable field changes while the ID stays the same.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given replacement room data with an empty room number.",
		When:  "When the room is updated.",
		Then:  "Then validation rejects it.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given replacement room data whose number starts with another floor.",
		When:  "When the room is updated.",
		Then:  "Then validation rejects it.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Given replacement room data with invalid event capacity.",
		When:  "When the room is updated.",
		Then:  "Then validation rejects it.",
	})

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
