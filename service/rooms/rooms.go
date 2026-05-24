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

	// We can disable this check if we want, if so, remember to update test
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

// UpdateRoom updates a room based on its ID with new Room type data
func UpdateRoom(db *sql.DB, data models.Room) (*models.Room, error) {
	if strings.TrimSpace(data.RoomNumber) == "" {
		return nil, fmt.Errorf("Room number is required")
	}

	// We can disable this check if we want, if so, remember to update test
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
		UPDATE rooms
		SET
			name = ?,
			room_number = ?,
			floor = ?,
			max_concurrent_games = ?,
			notes = ?,
			is_disabled = ?
		WHERE id = ?
		RETURNING
			id,
			name,
			room_number,
			floor,
			max_concurrent_games,
			notes,
			is_disabled
	`

	var updated models.Room

	err := db.QueryRow(
		query, data.Name, data.RoomNumber, data.Floor, data.MaxConcurrentGames, data.Notes, data.IsDisabled, data.ID,
	).Scan(
		&updated.ID, &updated.Name, &updated.RoomNumber, &updated.Floor, &updated.MaxConcurrentGames, &updated.Notes, &updated.IsDisabled,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	return &updated, nil
}

// UpdateRoom updates a room based on its ID with partial new information
func UpdateRoomPartial(db *sql.DB, data models.RoomInput) (*models.Room, error) {
	if data.ID < 1 {
		return nil, fmt.Errorf("Room ID is required and must be a valid positive number")
	}

	// Set up params based on partial data
	setParts := []string{}
	args := []any{}

	if data.Name != nil {
		if strings.TrimSpace(*data.Name) == "" {
			return nil, fmt.Errorf("room name cannot be empty")
		}

		setParts = append(setParts, "name = ?")
		args = append(args, *data.Name)
	}

	if data.RoomNumber != nil {
		if strings.TrimSpace(*data.RoomNumber) == "" {
			return nil, fmt.Errorf("room number cannot be empty")
		}

		setParts = append(setParts, "room_number = ?")
		args = append(args, *data.RoomNumber)
	}

	if data.Floor != nil {
		setParts = append(setParts, "floor = ?")
		args = append(args, *data.Floor)
	}

	if data.MaxConcurrentGames != nil {
		if *data.MaxConcurrentGames < 1 {
			return nil, fmt.Errorf("max concurrent games must be greater than 0")
		}

		setParts = append(setParts, "max_concurrent_games = ?")
		args = append(args, *data.MaxConcurrentGames)
	}

	if data.Notes != nil {
		setParts = append(setParts, "notes = ?")
		args = append(args, *data.Notes)
	}

	if data.IsDisabled != nil {
		setParts = append(setParts, "is_disabled = ?")
		args = append(args, *data.IsDisabled)
	}

	if len(args) == 0 {
		return nil, fmt.Errorf("Update room called without any updated data")
	}

	// Construct query based on partial data
	query := fmt.Sprintf(`
		UPDATE rooms
		SET %s
		WHERE id = ?
        RETURNING
			id,
			name,
			room_number,
			floor,
			max_concurrent_games,
			notes,
			is_disabled;
	`, strings.Join(setParts, ", "))

	// Construct args based on partial data
	args = append(args, data.ID)

	_, err := db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	var updated models.Room

	err = db.QueryRow(`
		SELECT
			id,
			name,
			room_number,
			floor,
			max_concurrent_games,
			notes,
			is_disabled
		FROM rooms
		WHERE id = ?
	`, data.ID).Scan(
		&updated.ID,
		&updated.Name,
		&updated.RoomNumber,
		&updated.Floor,
		&updated.MaxConcurrentGames,
		&updated.Notes,
		&updated.IsDisabled,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated room: %w", err)
	}

	return &updated, nil
}

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
