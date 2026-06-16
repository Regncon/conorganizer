package rooms

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestCreateRoom_CreatesRoomWithGeneratedID(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given valid room input.",
		When:  "When the room is created.",
		Then:  "Then the persisted room is returned with a generated ID.",
	})

	// Given
	expectedRoom := models.Room{
		ID:                 1,
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er eit gyldig rom",
		IsDisabled:         false,
	}
	inputRoom := expectedRoom
	inputRoom.ID = 0

	db := createRoomsTestDB(t)

	// When
	actualRoom, err := CreateRoom(db, inputRoom)

	// Then
	if err != nil {
		t.Fatalf("expected room creation to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestCreateRoom_WhenCalledRepeatedly_AutoIncrementsID(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given two valid room inputs.",
		When:  "When both rooms are created.",
		Then:  "Then their IDs are allocated in insert order.",
	})

	// Given
	expectedIDs := []int{1, 2}
	db := createRoomsTestDB(t)
	firstRoom := roomFixture("Hakkebakken", "101", 1)
	secondRoom := roomFixture("Tangerud", "201", 2)

	// When
	firstCreated, firstErr := CreateRoom(db, firstRoom)
	secondCreated, secondErr := CreateRoom(db, secondRoom)

	// Then
	if firstErr.HasErrors() {
		t.Fatalf("expected first room creation to succeed: %v", firstErr)
	}
	if secondErr.HasErrors() {
		t.Fatalf("expected second room creation to succeed: %v", secondErr)
	}
	actualIDs := []int{firstCreated.ID, secondCreated.ID}
	if !slices.Equal(expectedIDs, actualIDs) {
		t.Fatalf("room IDs mismatch\nexpected: %v\nactual:   %v", expectedIDs, actualIDs)
	}
}

func TestCreateRoom_WhenMaxConcurrentGamesIsInvalid_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given room input without capacity for any events.",
		When:  "When the room is created.",
		Then:  "Then validation rejects it.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	invalidRoom := roomFixture("Hakkebakken", "101", 1)
	invalidRoom.MaxConcurrentGames = -1

	// When
	_, err := CreateRoom(db, invalidRoom)
	actualError := err.HasError(models.RoomErrorMaxConcurrent)

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestDeleteRoom_RemovesOnlyTargetRoom(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given two stored rooms.",
		When:  "When one room is deleted.",
		Then:  "Then only the other room remains.",
	})

	// Given
	expectedRemainingRoomIDs := []int{2}
	db := createRoomsTestDB(t)
	roomToDelete := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	insertRoom(t, db, roomFixture("Tangerud", "201", 2))

	// When
	err := DeleteRoom(db, roomToDelete.ID)

	// Then
	if err != nil {
		t.Fatalf("expected room deletion to succeed: %v", err)
	}
	actualRemainingRoomIDs := queryRoomIDs(t, db)
	if !slices.Equal(expectedRemainingRoomIDs, actualRemainingRoomIDs) {
		t.Fatalf("remaining room IDs mismatch\nexpected: %v\nactual:   %v", expectedRemainingRoomIDs, actualRemainingRoomIDs)
	}
}

func TestDeleteRoom_WhenRoomDoesNotExist_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an empty room table.",
		When:  "When a missing room is deleted.",
		Then:  "Then the caller receives an error.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	err := DeleteRoom(db, 999)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestGetRoomByID_ReturnsStoredRoom(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a stored room.",
		When:  "When it is fetched by ID.",
		Then:  "Then the matching room is returned.",
	})

	// Given
	db := createRoomsTestDB(t)
	expectedRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))

	// When
	actualRoom, err := GetRoomByID(db, expectedRoom.ID)

	// Then
	if err != nil {
		t.Fatalf("expected room lookup to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestGetRoomByID_WhenIDIsInvalid_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an invalid room ID.",
		When:  "When it is fetched.",
		Then:  "Then validation rejects it.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := GetRoomByID(db, 0)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestGetRoomByID_WhenRoomDoesNotExist_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a positive room ID with no stored room.",
		When:  "When it is fetched.",
		Then:  "Then the caller receives an error.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := GetRoomByID(db, 999)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestGetAllRooms_ReturnsRoomsOrderedByFloorAndNumber(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given rooms inserted outside display order.",
		When:  "When all rooms are listed.",
		Then:  "Then rooms are ordered by floor and room number.",
	})

	// Given
	expectedRoomNumbers := []string{"101", "102", "201"}
	db := createRoomsTestDB(t)
	insertRoom(t, db, roomFixture("Tangerud", "201", 2))
	insertRoom(t, db, roomFixture("Brumms hus", "102", 1))
	insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))

	// When
	actualRooms, err := GetAllRooms(db)

	// Then
	if err != nil {
		t.Fatalf("expected room listing to succeed: %v", err)
	}
	actualRoomNumbers := roomNumbers(actualRooms)
	if !slices.Equal(expectedRoomNumbers, actualRoomNumbers) {
		t.Fatalf("room order mismatch\nexpected: %v\nactual:   %v", expectedRoomNumbers, actualRoomNumbers)
	}
}
