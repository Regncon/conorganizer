package rooms

import (
	"strings"
	"unicode/utf8"

	"github.com/Regncon/conorganizer/models"
)

// ValidateRooms validates that required entries in `room` is valid and returns all encutered errors
func ValidateRooms(room models.Room) models.RoomFormErrors {
	errors := models.RoomFormErrors{}

	if room.Name != "" && strings.TrimSpace(room.Name) == "" {
		errors.AddError(models.RoomErrorRoomNumber, "Rom namn kan ikkje berre innehalde mellomrom")
	}

	if utf8.RuneCountInString(room.Name) > 50 {
		errors.AddError(
			models.RoomErrorName,
			"Namn kan ikkje vere lengre enn 50 teikn",
		)
	}

	if strings.TrimSpace(room.RoomNumber) == "" {
		errors.AddError(models.RoomErrorRoomNumber, "Rom nummer er påkrevd")
	}

	if utf8.RuneCountInString(room.RoomNumber) > 10 {
		errors.AddError(
			models.RoomErrorRoomNumber,
			"Rom nummer kan ikkje vere lengre enn 10 teikn",
		)
	}

	return errors
}
