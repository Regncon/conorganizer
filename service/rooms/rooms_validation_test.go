package rooms

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestValidateRooms_WhenPublicNotesExceedsMaxLength_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a room with public notes longer than 1000 characters.",
		When:  "When the room is validated.",
		Then:  "Then validation rejects the public notes field.",
	})

	// Given
	expectedError := true
	room := roomFixture("Hakkebakken", "101", 1)
	room.PublicNotes = strings.Repeat("a", maxRoomNotesLength+1)

	// When
	errors := ValidateRooms(room)
	actualError := errors.HasError(models.RoomErrorPublicNotes)

	// Then
	if expectedError != actualError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestValidateRooms_WhenAdminNotesExceedsMaxLength_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a room with admin notes longer than 1000 characters.",
		When:  "When the room is validated.",
		Then:  "Then validation rejects the admin notes field.",
	})

	// Given
	expectedError := true
	room := roomFixture("Hakkebakken", "101", 1)
	room.AdminNotes = strings.Repeat("a", maxRoomNotesLength+1)

	// When
	errors := ValidateRooms(room)
	actualError := errors.HasError(models.RoomErrorAdminNotes)

	// Then
	if expectedError != actualError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestValidateRooms_WhenNotesAreExactlyMaxLength_ReturnsNoError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a room with public and admin notes exactly 1000 characters long.",
		When:  "When the room is validated.",
		Then:  "Then validation accepts both notes fields.",
	})

	// Given
	expectedError := false
	room := roomFixture("Hakkebakken", "101", 1)
	room.PublicNotes = strings.Repeat("a", maxRoomNotesLength)
	room.AdminNotes = strings.Repeat("a", maxRoomNotesLength)

	// When
	errors := ValidateRooms(room)
	actualPublicNotesError := errors.HasError(models.RoomErrorPublicNotes)
	actualAdminNotesError := errors.HasError(models.RoomErrorAdminNotes)

	// Then
	if expectedError != actualPublicNotesError {
		t.Fatalf("public notes error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualPublicNotesError)
	}
	if expectedError != actualAdminNotesError {
		t.Fatalf("admin notes error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualAdminNotesError)
	}
}
