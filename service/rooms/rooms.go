package rooms

import (
	"database/sql"

	"github.com/Regncon/conorganizer/models"
)

// CreateRoom creates a room
func CreateRoom(db *sql.DB, data models.Room)

// UpdateRoom updates a room based on its ID
func UpdateRoom(db *sql.DB, data models.Room)

// GetRoomByID returns a room pointer based on a roomID
func GetRoomByID(db *sql.DB, roomID string)

// GetAllRooms returns a list of all rooms stored in DB
func GetAllRooms(db *sql.DB)

// GetAllRoomStatusesByPuljeID Generates a list of all rooms, but unique to a pulje
func GetAllRoomStatusesByPuljeID(db *sql.DB, puljeID string) {
	// This function needs to return a detailed overview of available rooms, where
	// assigned events are limited to pulje
}

// SetRelationEventPuljeRoom assigns a room to an event in event_puljer_relation table
func AssignRoomToRelationEventPuljer(db *sql.DB, roomID string, relationEventPuljeID string) {
	// This function will assign a room by id to an event in event_puljer_relation
	// Validate that the room does not exceed max events based on pulje?
}
