package event

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/rooms"
)

type eventRoom struct {
	Pulje   models.PuljeRow
	Room    models.Room
	MapPath string
}

func getEventRooms(db *sql.DB, eventID string, programPublished bool) ([]eventRoom, error) {
	if !programPublished {
		return nil, nil
	}

	const query = `
		SELECT p.id, p.name, p.start_at, p.end_at,
			COALESCE(r.id, 0), COALESCE(r.name, ''), COALESCE(r.room_number, ''), COALESCE(r.floor, 0)
		FROM relation_event_puljer ep
		JOIN puljer p ON p.id = ep.pulje_id
		LEFT JOIN rooms r ON r.id = ep.room_id
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
		if err := rows.Scan(&assignment.Pulje.ID, &assignment.Pulje.Name, &assignment.Pulje.StartAt, &assignment.Pulje.EndAt,
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

func roomMapAlt(room models.Room) string {
	if room.Floor == 0 {
		return fmt.Sprintf("Kart med vei fra heisene og trappen til %s i %d. etasje", room.Name, room.Floor)
	}
	return fmt.Sprintf("Kart med vei fra trappegangen til %s, rom %s i %d. etasje", room.Name, room.RoomNumber, room.Floor)
}

type eventRoomGroup struct {
	Room    models.Room
	MapPath string
	Puljer  []models.PuljeRow
}

// Input and output follow the first occurrence of each room in the schedule.
func groupEventRooms(assignments []eventRoom) []eventRoomGroup {
	var groups []eventRoomGroup
	groupIndexByRoomID := make(map[int]int)
	for _, assignment := range assignments {
		groupIndex, exists := groupIndexByRoomID[assignment.Room.ID]
		if !exists {
			groupIndex = len(groups)
			groupIndexByRoomID[assignment.Room.ID] = groupIndex
			groups = append(groups, eventRoomGroup{Room: assignment.Room, MapPath: assignment.MapPath})
		}
		groups[groupIndex].Puljer = append(groups[groupIndex].Puljer, assignment.Pulje)
	}
	return groups
}
