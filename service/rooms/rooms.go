package rooms

import (
	"database/sql"
	"errors"
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
func GetRoomByID(db *sql.DB, roomID int) (*models.Room, error) {
	if roomID < 1 {
		return nil, fmt.Errorf("invalid room ID")
	}

	query := `
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
	`

	var room models.Room

	err := db.QueryRow(query, roomID).Scan(
		&room.ID,
		&room.Name,
		&room.RoomNumber,
		&room.Floor,
		&room.MaxConcurrentGames,
		&room.Notes,
		&room.IsDisabled,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("room with ID %d was not found", roomID)
		}

		return nil, fmt.Errorf("failed to get room with ID: %w", err)
	}

	return &room, nil
}

// GetAllRooms returns a list of all rooms stored in DB
func GetAllRooms(db *sql.DB) ([]*models.Room, error) {
	query := `
		SELECT
			id,
			name,
			room_number,
			floor,
			max_concurrent_games,
			notes,
			is_disabled
		FROM rooms
		ORDER BY floor ASC, room_number ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms: %w", err)
	}
	defer rows.Close()

	var rooms []*models.Room

	for rows.Next() {

		var room models.Room

		err := rows.Scan(
			&room.ID,
			&room.Name,
			&room.RoomNumber,
			&room.Floor,
			&room.MaxConcurrentGames,
			&room.Notes,
			&room.IsDisabled,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan room: %w", err)
		}

		rooms = append(rooms, &room)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating rooms: %w", err)
	}

	return rooms, nil
}

// GetAllRoomStatusesByPulje Generates a list of all rooms, but unique to a pulje
func GetAllRoomStatusesByPulje(db *sql.DB, pulje models.Pulje) (models.RoomStatusByPulje, error) {
	// This function needs to return a detailed overview of available rooms, where
	// assigned events are limited to pulje
	query := `
        SELECT
            p.id,

            r.id,
            r.name,
            r.room_number,
            r.max_concurrent_games,
            r.notes,

            e.id,
            e.title
        FROM puljer p
        CROSS JOIN rooms r
        LEFT JOIN relation_event_puljer rep
            ON rep.pulje_id = p.id
            AND rep.room_id = r.id
            AND rep.is_in_pulje = 1
        LEFT JOIN events e
            ON e.id = rep.event_id
        WHERE
            r.is_disabled = 0
        ORDER BY
            p.id,
            r.floor,
            r.room_number,
            e.title;
        `

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get rooms statuses: %w", err)
	}
	defer rows.Close()

	result := make(models.RoomStatusByPulje)
	for rows.Next() {
		var row models.RoomStatusRow

		err := rows.Scan(
			&row.PuljeID,

			&row.RoomID,
			&row.RoomName,
			&row.RoomNumber,
			&row.MaxConcurrentGames,
			&row.RoomNotes,

			&row.EventID,
			&row.EventTitle,
		)
		if err != nil {
			return nil, fmt.Errorf("scan room status row: %w", err)
		}

		// Ensure pulje exists
		if _, exists := result[row.PuljeID]; !exists {
			result[row.PuljeID] = make(map[int64]models.RoomByPulje)
		}

		// Create room if it doesn't exist yet
		room, exists := result[row.PuljeID][row.RoomID]
		if !exists {
			room = models.RoomByPulje{
				ID:                 row.RoomID,
				Name:               row.RoomName,
				RoomNumber:         row.RoomNumber,
				MaxConcurrentGames: row.MaxConcurrentGames,
				Notes:              row.RoomNotes,
				AssignedEventsID:   []models.RoomEventPuljeSummary{},
			}
		}

		// Add event if assigned
		if row.EventID.Valid {
			room.AssignedEventsID = append(
				room.AssignedEventsID,
				models.RoomEventPuljeSummary{
					EventPuljeID: fmt.Sprintf(
						"%s:%s",
						row.EventID.String,
						row.PuljeID,
					),
					EventID: row.EventID.String,
					Title:   row.EventTitle.String,
				},
			)
		}

		result[row.PuljeID][row.RoomID] = room
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate room status rows: %w", err)
	}

	return result, nil
}

// SetRelationEventPuljeRoom assigns a room to an event in `relation_event_puljer`
func AssignRoomToRelationEventPuljer(db *sql.DB, roomID int64, relationEventPuljeID string) (models.EventPulje, error) {
	query := `
		UPDATE relation_event_puljer
		SET room_id = ?
		WHERE event_id = ?
		RETURNING
			event_id,
			pulje_id,
			is_in_pulje,
			is_published,
			room_id
	`

	var result models.EventPulje
	err := db.QueryRow(query, roomID, relationEventPuljeID).Scan(
		&result.EventID,
		&result.PuljeID,
		&result.IsInPulje,
		&result.IsPublished,
		&result.RoomID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.EventPulje{}, fmt.Errorf(
				"no relation_event_puljer found for event_id=%s",
				relationEventPuljeID,
			)
		}
		return models.EventPulje{}, fmt.Errorf(
			"Error assigning room to relation_event_puljer: %w",
			err,
		)
	}

	return result, nil
}

func GetRelationEventPuljerByRoomIDAndPuljeID(db *sql.DB, roomID int, puljeID string) {
	// This functino will return all events assigned to a room limited by the pulje
}
