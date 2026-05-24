package rooms

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

// CreateRoom creates a room, on success updates `ID` with its entry ID
func CreateRoom(db *sql.DB, data models.Room) (*models.Room, error) {
	if strings.TrimSpace(data.RoomNumber) == "" {
		return nil, fmt.Errorf("Room number is required")
	}

	// We can disable this check if we want
	if !strings.HasPrefix(data.RoomNumber, fmt.Sprintf("%d", data.Floor)) {
		return nil, fmt.Errorf(
			"Room number must start with the floor number, eg: %dxx, got: %s",
			data.Floor,
			data.RoomNumber,
		)
	}

	if data.MaxConcurrentGames < 1 {
		return nil, fmt.Errorf("Max concurrent events must be greater than 0, got: %d", data.MaxConcurrentGames)
	}

	query := `
        INSERT INTO rooms (
            name,
			room_number,
			floor,
			max_concurrent_games,
			notes,
			is_disabled
        )
        VALUES (?, ?, ?, ?, ?, ?)
    `

	result, err := db.Exec(query, data.Name, data.RoomNumber, data.Floor, data.MaxConcurrentGames, data.Notes, data.IsDisabled)
	if err != nil {
		return nil, fmt.Errorf("Failed to create room: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("Failed to get created room id: %w", err)
	}

	data.ID = int(id)
	return &data, nil
}

// UpdateRoom updates a room based on its ID
func UpdateRoom(db *sql.DB, data models.Room) {}

// GetRoomByID returns a room pointer based on a roomID
func GetRoomByID(db *sql.DB, roomID int) {}

// GetAllRooms returns a list of all rooms stored in DB
func GetAllRooms(db *sql.DB) {}

// GetAllRoomStatusesByPuljeID Generates a list of all rooms, but unique to a pulje
func GetAllRoomStatusesByPuljeID(db *sql.DB, puljeID string) {
	// This function needs to return a detailed overview of available rooms, where
	// assigned events are limited to pulje

	// Should this include complete events from event puljer, just event puljer id or just convert this to a number?
}

// SetRelationEventPuljeRoom assigns a room to an event in `relation_event_puljer`
func AssignRoomToRelationEventPuljer(db *sql.DB, roomID int, relationEventPuljeID string) {
	// This function will assign a room by id to an event in relation_event_puljer
	// Validate that the room does not exceed max events based on pulje?

	// Move this function to event pujer as parent...
}

func GetRelationEventPuljerByRoomIDAndPuljeID(db *sql.DB, roomID int, puljeID string) {
	// This functino will return all events assigned to a room limited by the pulje
}
