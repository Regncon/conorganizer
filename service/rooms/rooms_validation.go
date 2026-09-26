package rooms

import (
	"strings"
	"unicode/utf8"

	"github.com/Regncon/conorganizer/models"
)

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

	return errors
}
