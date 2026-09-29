package rooms

import (
	"strings"
	"unicode/utf8"

	"github.com/Regncon/conorganizer/models"
)

// maxRoomNotesLength is the maximum number of characters allowed in public and admin notes.
const maxRoomNotesLength = 1000

// ValidateRooms validates that required entries in `room` is valid and returns all encutered errors
func ValidateRooms(room models.Room) models.RoomFormErrors {
	errors := models.RoomFormErrors{}
	errors.ResetErrors()

	if room.Name != "" && strings.TrimSpace(room.Name) == "" {
		errors.AddError(models.RoomErrorName, "Romnavn kan ikke bare inneholde mellomrom")
	}

	if utf8.RuneCountInString(room.Name) > 50 {
		errors.AddError(
			models.RoomErrorName,
			"Navn kan ikke være lengre enn 50 tegn",
		)
	}

	if strings.TrimSpace(room.RoomNumber) == "" {
		errors.AddError(models.RoomErrorRoomNumber, "Romnummer er påkrevd")
	}

	if utf8.RuneCountInString(room.RoomNumber) > 10 {
		errors.AddError(
			models.RoomErrorRoomNumber,
			"Romnummer kan ikke være lengre enn 10 tegn",
		)
	}

	if utf8.RuneCountInString(room.PublicNotes) > maxRoomNotesLength {
		errors.AddError(
			models.RoomErrorPublicNotes,
			"Offentlige notater kan ikke være lengre enn 1000 tegn",
		)
	}

	if utf8.RuneCountInString(room.AdminNotes) > maxRoomNotesLength {
		errors.AddError(
			models.RoomErrorAdminNotes,
			"Admin-notater kan ikke være lengre enn 1000 tegn",
		)
	}

	return errors
}
