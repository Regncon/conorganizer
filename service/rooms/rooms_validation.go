package rooms

import (
	"strings"

	"github.com/Regncon/conorganizer/models"
)

// ValidateRooms validates that required entries in `room` is valid and returns all encutered errors
func ValidateRooms(room models.Room) models.RoomFormErrors {
	errors := models.RoomFormErrors{}

	if strings.TrimSpace(room.RoomNumber) == "" {
		errors.AddError(models.RoomErrorRoomNumber, "Rom nummer er påkrevd")
	}

	if room.MaxConcurrentGames < 1 {
		errors.AddError(models.RoomErrorMaxConcurrent, "Maks samtidige spill må være minst 1")
	}

	return errors
}
