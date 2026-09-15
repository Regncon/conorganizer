package event

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/rooms"
)

type eventRoom struct {
	PuljeID   models.Pulje
	PuljeName string
	Room      models.Room
	MapPath   string
}

func getEventRooms(db *sql.DB, eventID string, programPublished bool) ([]eventRoom, error) {
	if !programPublished {
		return nil, nil
	}

	const query = `
		SELECT p.id, p.name, r.id, r.name, r.room_number, r.floor
		FROM relation_event_puljer ep
		JOIN puljer p ON p.id = ep.pulje_id
		JOIN rooms r ON r.id = ep.room_id
		WHERE ep.event_id = ?
			AND ep.is_in_pulje = 1
		ORDER BY p.start_at, p.id
	`
	rows, err := db.Query(query, eventID)
	if err != nil {
		return nil, fmt.Errorf("query rooms for event %s: %w", eventID, err)
	}
	defer rows.Close()

	var assignments []eventRoom
	for rows.Next() {
		var assignment eventRoom
		if err := rows.Scan(&assignment.PuljeID, &assignment.PuljeName,
			&assignment.Room.ID, &assignment.Room.Name, &assignment.Room.RoomNumber, &assignment.Room.Floor); err != nil {
			return nil, fmt.Errorf("scan room for event %s: %w", eventID, err)
		}
		assignment.MapPath, _ = rooms.MapPathForRoom(assignment.Room.RoomNumber)
		assignments = append(assignments, assignment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rooms for event %s: %w", eventID, err)
	}
	return assignments, nil
}

func eventRoomDialogID(eventID string, puljeID models.Pulje) string {
	return fmt.Sprintf("event-room-map-%s-%s", eventID, puljeID)
}
