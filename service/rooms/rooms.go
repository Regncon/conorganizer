package rooms

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

// CreateRoom creates a room, on success updates `ID` with its entry ID
func CreateRoom(db *sql.DB, data models.Room) (*models.Room, models.RoomFormErrors) {
	errors := ValidateRooms(data)
	if errors.HasErrors() {
		return nil, errors
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

	result, err := db.Exec(query,
		data.Name,
		data.RoomNumber,
		data.Floor,
		data.MaxConcurrentGames,
		data.Notes,
		data.IsDisabled,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed: rooms.room_number") {
			errors.AddError(models.RoomErrorRoomNumber, "room number must be unique")
		} else {
			errors.AddError(models.RoomError, fmt.Sprintf("failed to create room: %s", err.Error()))
		}

		return nil, errors
	}

	id, err := result.LastInsertId()
	if err != nil {
		errors.AddError(models.RoomError, fmt.Sprintf("failed to get created room id: %s", err.Error()))
		return nil, errors
	}

	data.ID = int(id)
	return &data, nil
}

// DeleteRoom removes a room given an ID, since pragma is enabled the change will cascade
func DeleteRoom(db *sql.DB, roomID int) error {
	query := `DELETE FROM rooms WHERE id = ?`

	result, err := db.Exec(query, roomID)
	if err != nil {
		return fmt.Errorf("delete room %d: %w", roomID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("room %d not found", roomID)
	}

	return nil
}

// UpdateRoom updates a room based on its ID with new Room type data
func UpdateRoom(db *sql.DB, data models.Room) (*models.Room, models.RoomFormErrors) {
	errors := ValidateRooms(data)
	if errors.HasErrors() {
		return nil, errors
	}

	if data.ID < 1 {
		errors.AddError(models.RoomError, "room ID is required and must be a valid positive number")
		return nil, errors
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
		query,
		data.Name,
		data.RoomNumber,
		data.Floor,
		data.MaxConcurrentGames,
		data.Notes,
		data.IsDisabled,
		data.ID,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.RoomNumber,
		&updated.Floor,
		&updated.MaxConcurrentGames,
		&updated.Notes,
		&updated.IsDisabled,
	)

	if err != nil {
		errors.AddError(models.RoomError, fmt.Sprintf("failed to update room: %s", err.Error()))
		return nil, errors
	}

	return &updated, nil
}

// UpdateRoom updates a room based on its ID with partial new information
func UpdateRoomPartial(db *sql.DB, data models.RoomInput) (*models.Room, models.RoomFormErrors) {
	// Init error handling and check for ID before continuing
	var errors models.RoomFormErrors
	if data.ID < 1 {
		errors.AddError(models.RoomError, "room ID is required and must be a valid positive number")
		return nil, errors
	}

	// Set up params based on partial data
	setParts := []string{}
	args := []any{}

	if data.Name != nil || strings.TrimSpace(*data.Name) != "" {
		setParts = append(setParts, "name = ?")
		args = append(args, *data.Name)
	}

	if data.RoomNumber != nil {
		if strings.TrimSpace(*data.RoomNumber) == "" {
			errors.AddError(models.RoomErrorRoomNumber, "room number cannot be empty")
			return nil, errors
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
			errors.AddError(models.RoomErrorRoomNumber, "max concurrent games must be greater than 0")
			return nil, errors
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

	// Check if any data was being updated
	if len(args) == 0 {
		errors.AddError(models.RoomError, "UpdateRoomPartial called without any updated data")
	}

	// Check if errors exists before running database update
	if errors.HasErrors() {
		return nil, errors
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

	// Add ID to constructed args
	args = append(args, data.ID)

	// Update and return updated data
	var updated models.Room
	err := db.QueryRow(query, args...).Scan(
		&updated.ID,
		&updated.Name,
		&updated.RoomNumber,
		&updated.Floor,
		&updated.MaxConcurrentGames,
		&updated.Notes,
		&updated.IsDisabled,
	)

	if err != nil {
		errors.AddError(models.RoomError, fmt.Sprintf("failed to update room: %s", err.Error()))
		return nil, errors
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
func GetAllRooms(db *sql.DB) ([]models.Room, error) {
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

	var rooms []models.Room

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

		rooms = append(rooms, room)
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
            r.floor,
            r.max_concurrent_games,
            r.is_disabled,
            r.notes,

            e.id,
            e.title,
            e.max_players
        FROM puljer p
        CROSS JOIN rooms r
        LEFT JOIN relation_event_puljer rep
            ON rep.pulje_id = p.id
            AND rep.room_id = r.id
            AND rep.is_in_pulje = 1
        LEFT JOIN events e
            ON e.id = rep.event_id
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
			&row.Floor,
			&row.MaxConcurrentGames,
			&row.IsDisabled,
			&row.RoomNotes,

			&row.EventID,
			&row.EventTitle,
			&row.EventMaxPlayers,
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
				Floor:              row.Floor,
				MaxConcurrentGames: row.MaxConcurrentGames,
				Notes:              row.RoomNotes,
				IsDisabled:         row.IsDisabled,
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
					EventID:    row.EventID.String,
					Title:      row.EventTitle.String,
					MaxPlayers: int(row.EventMaxPlayers.Int32),
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
			"error assigning room to relation_event_puljer: %w",
			err,
		)
	}

	return result, nil
}
