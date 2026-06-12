package rooms

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
)

func TestCreateRoom_CreatesRoomWithGeneratedID(t *testing.T) {
	// Given valid room input,
	// when the room is created,
	// then the persisted room is returned with a generated ID.

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
	// Given two valid room inputs,
	// when both rooms are created,
	// then their IDs are allocated in insert order.

	// Given
	expectedIDs := []int{1, 2}
	db := createRoomsTestDB(t)
	firstRoom := roomFixture("Hakkebakken", "101", 1)
	secondRoom := roomFixture("Tangerud", "201", 2)

	// When
	firstCreated, firstErr := CreateRoom(db, firstRoom)
	secondCreated, secondErr := CreateRoom(db, secondRoom)

	// Then
	if firstErr != nil {
		t.Fatalf("expected first room creation to succeed: %v", firstErr)
	}
	if secondErr != nil {
		t.Fatalf("expected second room creation to succeed: %v", secondErr)
	}
	actualIDs := []int{firstCreated.ID, secondCreated.ID}
	if !slices.Equal(expectedIDs, actualIDs) {
		t.Fatalf("room IDs mismatch\nexpected: %v\nactual:   %v", expectedIDs, actualIDs)
	}
}

func TestCreateRoom_WhenRoomNumberDoesNotMatchFloor_ReturnsError(t *testing.T) {
	// Given room input where the room number starts with another floor,
	// when the room is created,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	invalidRoom := roomFixture("Tangerud", "203", 3)

	// When
	_, err := CreateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestCreateRoom_WhenMaxConcurrentGamesIsInvalid_ReturnsError(t *testing.T) {
	// Given room input without capacity for any events,
	// when the room is created,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	invalidRoom := roomFixture("Hakkebakken", "101", 1)
	invalidRoom.MaxConcurrentGames = 0

	// When
	_, err := CreateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestDeleteRoom_RemovesOnlyTargetRoom(t *testing.T) {
	// Given two stored rooms,
	// when one room is deleted,
	// then only the other room remains.

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
	// Given an empty room table,
	// when a missing room is deleted,
	// then the caller receives an error.

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
	// Given a stored room,
	// when it is fetched by ID,
	// then the matching room is returned.

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
	// Given an invalid room ID,
	// when it is fetched,
	// then validation rejects it.

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
	// Given a positive room ID with no stored room,
	// when it is fetched,
	// then the caller receives an error.

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
	// Given rooms inserted outside display order,
	// when all rooms are listed,
	// then rooms are ordered by floor and room number.

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
