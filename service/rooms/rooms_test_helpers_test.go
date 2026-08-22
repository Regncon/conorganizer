package rooms

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func createRoomsTestDB(t testing.TB) *sql.DB {
	t.Helper()

	return testutil.CreateTestDB(t, "rooms")
}

func roomFixture(name string, roomNumber string, floor int) models.Room {
	return models.Room{
		Name:               name,
		RoomNumber:         roomNumber,
		Floor:              floor,
		MaxConcurrentGames: 2,
		Notes:              "Romnotat",
		IsDisabled:         false,
	}
}

func insertRoom(t testing.TB, db *sql.DB, input models.Room) models.Room {
	t.Helper()

	var room models.Room
	err := db.QueryRow(`
		INSERT INTO rooms (
			name,
			room_number,
			floor,
			max_concurrent_games,
			notes,
			is_disabled
		)
		VALUES (?, ?, ?, ?, ?, ?)
		RETURNING
			id,
			name,
			room_number,
			floor,
			max_concurrent_games,
			notes,
			is_disabled
	`, input.Name, input.RoomNumber, input.Floor, input.MaxConcurrentGames, input.Notes, input.IsDisabled).Scan(
		&room.ID,
		&room.Name,
		&room.RoomNumber,
		&room.Floor,
		&room.MaxConcurrentGames,
		&room.Notes,
		&room.IsDisabled,
	)
	if err != nil {
		t.Fatalf("failed to insert room: %v", err)
	}

	return room
}

func queryRoomIDs(t testing.TB, db *sql.DB) []int {
	t.Helper()

	rows, err := db.Query(`SELECT id FROM rooms ORDER BY id`)
	if err != nil {
		t.Fatalf("failed to query room IDs: %v", err)
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("failed to scan room ID: %v", err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("failed to iterate room IDs: %v", err)
	}

	return ids
}

func seedRoomEventLookups(t testing.TB, db *sql.DB) {
	t.Helper()

	testutil.MustExec(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, models.EventStatusAnnounced)
	testutil.MustExec(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeBoardGame)
	testutil.MustExec(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupAdultsOnly)
	testutil.MustExec(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeLongRunning)
	testutil.MustExec(t, db, `INSERT INTO pulje_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, models.PuljeStatusOpen)
}

func insertPulje(t testing.TB, db *sql.DB, id models.Pulje, name string) models.Pulje {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at)
		VALUES (?, ?, ?, ?, ?)
	`, id, name, models.PuljeStatusOpen, "2025-10-03", "2025-10-03")

	return id
}

func insertEvent(t testing.TB, db *sql.DB, eventID string, title string, maxPlayers int) string {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO events (
			id,
			title,
			intro,
			description,
			system,
			event_type,
			age_group,
			event_runtime,
			host_name,
			email,
			phone_number,
			max_players,
			beginner_friendly,
			can_be_run_in_english,
			status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, eventID, title, "Intro", "Description", "System", models.EventTypeBoardGame, models.AgeGroupAdultsOnly, models.RunTimeLongRunning, "Host", "host@example.com", "12345678", maxPlayers, true, true, models.EventStatusAnnounced)

	return eventID
}

func insertEventPulje(t testing.TB, db *sql.DB, eventID string, puljeID models.Pulje, roomID sql.NullInt64) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO relation_event_puljer (event_id, pulje_id, room_id)
		VALUES (?, ?, ?)
	`, eventID, puljeID, roomID)
}

func roomNumbers(rooms []models.Room) []string {
	numbers := make([]string, 0, len(rooms))
	for _, room := range rooms {
		numbers = append(numbers, room.RoomNumber)
	}
	return numbers
}

func assertRoomMatches(t testing.TB, expected models.Room, actual models.Room) {
	t.Helper()

	if expected != actual {
		t.Fatalf("room mismatch\nexpected: %+v\nactual:   %+v", expected, actual)
	}
}

func assertRoomStatusAssignments(t testing.TB, expected map[models.Pulje]map[int64][]string, actual models.RoomStatusByPulje) {
	t.Helper()

	for expectedPulje, expectedRooms := range expected {
		actualRooms, exists := actual[expectedPulje]
		if !exists {
			t.Fatalf("expected pulje %s to exist in room statuses", expectedPulje)
		}

		for expectedRoomID, expectedTitles := range expectedRooms {
			actualRoom, exists := actualRooms[expectedRoomID]
			if !exists {
				t.Fatalf("expected room %d to exist in pulje %s", expectedRoomID, expectedPulje)
			}

			actualTitles := roomStatusEventTitles(actualRoom)
			if !slices.Equal(expectedTitles, actualTitles) {
				t.Fatalf(
					"assigned event titles mismatch for pulje %s room %d\nexpected: %v\nactual:   %v",
					expectedPulje,
					expectedRoomID,
					expectedTitles,
					actualTitles,
				)
			}
		}
	}
}

func roomStatusEventTitles(room models.RoomByPulje) []string {
	titles := make([]string, 0, len(room.AssignedEventsID))
	for _, assignedEvent := range room.AssignedEventsID {
		titles = append(titles, assignedEvent.Title)
	}
	return titles
}

func ptr[T any](value T) *T {
	return &value
}
