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
	invalidRoomConcurrency.MaxConcurrentGames = -1

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
