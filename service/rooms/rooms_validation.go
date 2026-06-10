package rooms

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

// ValidateRooms validates that required entries in `room` is valid and returns all encutered errors
func ValidateRooms(room models.Room) models.RoomFormErrors {
	errors := models.RoomFormErrors{}

	if strings.TrimSpace(room.RoomNumber) == "" {
		errors.RoomNumber = "Romnummer er påkrevd"
	}

	if !strings.HasPrefix(room.RoomNumber, strconv.Itoa(room.Floor)) {
		errors.RoomNumber = "Romnummer må starte med etasje som første tall"
	}

	if room.MaxConcurrentGames < 1 {
		errors.MaxConcurrentGames = "Maks samtidige spill må være minst 1"
	}

	return errors
}

// ValidateRoomsByPulje is used for validating that a snapshot of rooms, based on a pulje, is valid (eg. assigned vs max concurrent events per pulje)
func ValidateRoomsByPulje(db *sql.DB, puljeID string) {}

// ValidateDisabledRoomsCascade Validates that room disabled status has cascaded and no orphans exist in `relation_event_puljer`
func ValidateDisabledRoomsCascade(db *sql.DB) {}
