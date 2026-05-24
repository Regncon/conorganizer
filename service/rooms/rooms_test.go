package rooms

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestCreateRoom(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	var validRoom = models.Room{
		ID:                 0,
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}

	var invalidRoomNumber = validRoom
	invalidRoomNumber.Floor = 2

	var invalidRoomConcurrency = validRoom
	invalidRoomConcurrency.MaxConcurrentGames = 0

	// When
	// - create room is called once
	createRoomResult, err := CreateRoom(db, validRoom)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// - create room is called again, tesing auto increment
	createRoomResult2, err := CreateRoom(db, validRoom)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// - create room is called with an invalid room number based on the floor being its prefix
	_, errInvalidRoomNumber := CreateRoom(db, invalidRoomNumber)

	// - create room is called with an invalid room number based on the floor being its prefix
	_, errInvalidConcurrency := CreateRoom(db, invalidRoomConcurrency)

	// Then
	if validRoom.ID == createRoomResult.ID {
		t.Fatalf("Room ID did not update correctly after insert\nexpected: %d\nrecieved: %d", validRoom.ID, createRoomResult.ID)
	}

	validRoom.ID = createRoomResult.ID
	if validRoom != *createRoomResult {
		t.Fatalf("createRoomResult did not match happyRoom\nexpected: \t%v\nrecieved: \t%v", validRoom, createRoomResult)
	}

	var expectedID = createRoomResult.ID + 1
	if createRoomResult2.ID != expectedID {
		t.Fatalf("createRoom did not auto increment ID\nexpected: %d\nrecieved: %d", expectedID, createRoomResult2.ID)
	}

	if errInvalidRoomNumber == nil {
		t.Fatal("expected error when creating room with invalid room number, but it was allowed")
	}

	if errInvalidConcurrency == nil {
		t.Fatal("expected error when creating room with 0 or less max games, but it was allowed")
	}
}

func TestUpdateRoom(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	var originalRoom = models.Room{
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}

	resultCreateRoom, err := CreateRoom(db, originalRoom)
	if err != nil {
		t.Fatalf("unexpected error creating original room: %v", err)
	}

	var updatedRoom = models.Room{
		ID:                 resultCreateRoom.ID,
		Name:               "Tangerud",
		RoomNumber:         "209",
		Floor:              2,
		MaxConcurrentGames: 3,
		Notes:              "Dette er en oppdatert note",
		IsDisabled:         true,
	}

	var updatedRoomInvalidRoomNumber = *resultCreateRoom
	updatedRoomInvalidRoomNumber.RoomNumber = ""

	var updatedRoomInvalidRoomConcurrency = *resultCreateRoom
	updatedRoomInvalidRoomConcurrency.MaxConcurrentGames = -1

	var updatedRoomInvalidRoomNumberFloor = *resultCreateRoom
	updatedRoomInvalidRoomNumberFloor.Floor = 3
	updatedRoomInvalidRoomNumberFloor.RoomNumber = "203"

	// When
	resultUpdateRoomValid, err := UpdateRoom(db, updatedRoom)
	if err != nil {
		t.Fatalf("unexpected error when updating valid room: %v", err)
	}

	_, errInvalidRoomNumber := UpdateRoom(db, updatedRoomInvalidRoomNumber)
	_, errInvalidRoomConcurrency := UpdateRoom(db, updatedRoomInvalidRoomConcurrency)
	_, errInvalidRoomNumberFloor := UpdateRoom(db, updatedRoomInvalidRoomNumberFloor)

	// Then
	if resultCreateRoom.ID != resultUpdateRoomValid.ID {
		t.Fatalf("UpdateRoom caused room ID to change\nexpected: \t%d\nrecieved: \t%d", resultCreateRoom.ID, resultUpdateRoomValid.ID)
	}
	if resultCreateRoom.Name == resultUpdateRoomValid.Name {
		t.Fatalf("UpdateRoom did not update name correctly\nexpected: \t%s\nrecieved: \t%s", resultCreateRoom.Name, resultUpdateRoomValid.Name)
	}
	if resultCreateRoom.RoomNumber == resultUpdateRoomValid.RoomNumber {
		t.Fatalf("UpdateRoom did not update room number correctly\nexpected: \t%s\nrecieved: \t%s", resultCreateRoom.RoomNumber, resultUpdateRoomValid.RoomNumber)
	}
	if resultCreateRoom.Floor == resultUpdateRoomValid.Floor {
		t.Fatalf("UpdateRoom did not update floor number correctly\nexpected: \t%d\nrecieved: \t%d", resultCreateRoom.Floor, resultUpdateRoomValid.Floor)
	}
	if resultCreateRoom.Notes == resultUpdateRoomValid.Notes {
		t.Fatalf("UpdateRoom did not update notes correctly\nexpected: \t%s\nrecieved: \t%s", resultCreateRoom.Notes, resultUpdateRoomValid.Notes)
	}
	if resultCreateRoom.MaxConcurrentGames == resultUpdateRoomValid.MaxConcurrentGames {
		t.Fatalf("UpdateRoom did not update max concurrent games correctly\nexpected: \t%d\nrecieved: \t%d", resultCreateRoom.MaxConcurrentGames, resultUpdateRoomValid.MaxConcurrentGames)
	}

	if errInvalidRoomNumber == nil {
		t.Fatalf("UpdateRoom allowed updating invalid room number")
	}
	if errInvalidRoomNumberFloor == nil {
		t.Fatalf("UpdateRoom allowed updating with an invalid room number and floor combination")
	}
	if errInvalidRoomConcurrency == nil {
		t.Fatalf("UpdateRoom allowed updating invalid max concurrent games")
	}
}

func TestUpdateRoomPartial(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	var originalRoom = models.Room{
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}
	var originalRoomPartial = models.Room{
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}

	createRoomResult, err := CreateRoom(db, originalRoom)
	if err != nil {
		t.Fatalf("unexpected error creating original room: %v", err)
	}
	createRoomPartialResult, err := CreateRoom(db, originalRoomPartial)
	if err != nil {
		t.Fatalf("unexpected error creating original room: %v", err)
	}

	// When
	var updatedName string = "Tangerud"
	var updatedRoomNumber string = "303"
	var updatedFloor int = 3
	var updatedConcurrent int = 3
	var updatedNotes string = ""
	var updatedDisables bool = true

	var updatedRoom = models.RoomInput{
		ID:                 createRoomResult.ID,
		Name:               &updatedName,
		RoomNumber:         &updatedRoomNumber,
		Floor:              &updatedFloor,
		MaxConcurrentGames: &updatedConcurrent,
		Notes:              &updatedNotes,
		IsDisabled:         &updatedDisables,
	}
	updatedRoomResult, err := UpdateRoomPartial(db, updatedRoom)
	if err != nil {
		t.Fatalf("unexpected error when updating valid room: %v", err)
	}
	partialUpdatedRoomResult, err := UpdateRoomPartial(db, models.RoomInput{ID: createRoomPartialResult.ID, Name: &updatedName})
	if err != nil {
		t.Fatalf("unexpected error when partially updating valid room: %v", err)
	}

	var invalidRoomNumber string = ""
	var invalidName string = ""
	var invalidConcurrent int = -1

	_, errInvalidID := UpdateRoomPartial(db, models.RoomInput{})
	_, errInvalidName := UpdateRoomPartial(db, models.RoomInput{ID: 1, Name: &invalidName})
	_, errInvalidRoomNumber := UpdateRoomPartial(db, models.RoomInput{ID: 1, RoomNumber: &invalidRoomNumber})
	_, errInvalidConcurrency := UpdateRoomPartial(db, models.RoomInput{ID: 1, MaxConcurrentGames: &invalidConcurrent})

	// Then
	if originalRoom.Name != createRoomResult.Name {
		t.Errorf("Original name was different to what create room returned")
	}
	if createRoomResult.Name == updatedRoomResult.Name {
		t.Errorf("Room name persisted after updateRoom was called successfully\nexpected: \t%s\nrecieved: \t%s", createRoomResult.Name, updatedRoomResult.Name)
	}
	if createRoomPartialResult.Name == partialUpdatedRoomResult.Name {
		t.Errorf("Room name persisted after updateRoom was called successfully\nexpected: \t%s\nrecieved: \t%s", createRoomPartialResult.Name, partialUpdatedRoomResult.Name)
	}

	if errInvalidID == nil {
		t.Errorf("UpdateRoom allowed update when ID was omited")
	}
	if errInvalidName == nil {
		t.Errorf("UpdateRoom allowed update when name was an empty string")
	}
	if errInvalidConcurrency == nil {
		t.Errorf("UpdateRoom allowed update when max concurrent games was less than 1")
	}
	if errInvalidRoomNumber == nil {
		t.Errorf("UpdateRoom allowed update when room number was an empty string")
	}
}
