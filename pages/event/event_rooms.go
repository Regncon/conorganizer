package event

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/program"
	"github.com/Regncon/conorganizer/service/rooms"
)

type eventRoom struct {
	Pulje       models.PuljeRow
	ProgramTime string
	Room        models.Room
	// MapPath is empty when the room is hidden, or has no map.
	MapPath string
}

// getEventRooms returns the room assignments for an event's active puljer.
// Nobody, admins included, sees the room for a pulje whose room assignment is not published yet:
// the room is cleared to its zero value so the schedule entry renders as unassigned.
func getEventRooms(db *sql.DB, eventID string, programPublished bool) ([]eventRoom, error) {
	if !programPublished {
		return nil, nil
	}

	const query = `
		SELECT p.id, p.name, p.start_at, p.end_at, p.rooms_published, ep.program_time,
			COALESCE(r.id, 0), COALESCE(r.name, ''), COALESCE(r.room_number, ''), COALESCE(r.floor, 0), COALESCE(r.public_notes, '')
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
		var roomsPublished bool
		if err := rows.Scan(&assignment.Pulje.ID, &assignment.Pulje.Name, &assignment.Pulje.StartAt, &assignment.Pulje.EndAt, &roomsPublished, &assignment.ProgramTime,
			&assignment.Room.ID, &assignment.Room.Name, &assignment.Room.RoomNumber, &assignment.Room.Floor, &assignment.Room.PublicNotes); err != nil {
			return nil, fmt.Errorf("scan room for event %s: %w", eventID, err)
		}
		if roomsPublished {
			assignment.MapPath, _ = rooms.MapPathForRoom(assignment.Room.RoomNumber)
		} else {
			assignment.Room = models.Room{}
		}
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
	Puljer  []program.PuljeTime
	// Schedule is when the event is in this room, one line each.
	Schedule []string
}

// eventSchedule is when the event happens across all its puljer, one line each.
func eventSchedule(assignments []eventRoom, isProgramEvent bool) ([]string, error) {
	puljeTimes := make([]program.PuljeTime, 0, len(assignments))
	for _, assignment := range assignments {
		puljeTimes = append(puljeTimes, program.PuljeTime{Pulje: assignment.Pulje, ProgramTime: assignment.ProgramTime})
	}
	return program.EventSchedule(puljeTimes, isProgramEvent)
}

// Input and output follow the first occurrence of each room in the schedule.
// Cleared (hidden) rooms all share ID 0 and are grouped together as unassigned, same as any other pulje without a room.
func groupEventRooms(assignments []eventRoom, isProgramEvent bool) ([]eventRoomGroup, error) {
	var groups []eventRoomGroup
	groupIndexByRoomID := make(map[int]int)
	for _, assignment := range assignments {
		groupIndex, exists := groupIndexByRoomID[assignment.Room.ID]
		if !exists {
			groupIndex = len(groups)
			groupIndexByRoomID[assignment.Room.ID] = groupIndex
			groups = append(groups, eventRoomGroup{Room: assignment.Room, MapPath: assignment.MapPath})
		}
		groups[groupIndex].Puljer = append(groups[groupIndex].Puljer, program.PuljeTime{Pulje: assignment.Pulje, ProgramTime: assignment.ProgramTime})
	}
	for index := range groups {
		schedule, err := program.EventSchedule(groups[index].Puljer, isProgramEvent)
		if err != nil {
			return nil, err
		}
		groups[index].Schedule = schedule
	}
	return groups, nil
}
