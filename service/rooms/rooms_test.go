package rooms

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestCreateRoom_CreatesRoomWithGeneratedID(t *testing.T) {
	// Given valid room input,
	// when the room is created,
	// then the persisted room is returned with a generated ID.

	// Given
	expectedRoom := models.Room{
		ID:                 1,
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er eit gyldig rom",
		IsDisabled:         false,
	}
	inputRoom := expectedRoom
	inputRoom.ID = 0

	db := createRoomsTestDB(t)

	// When
	actualRoom, err := CreateRoom(db, inputRoom)

	// Then
	if err != nil {
		t.Fatalf("expected room creation to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestCreateRoom_WhenCalledRepeatedly_AutoIncrementsID(t *testing.T) {
	// Given two valid room inputs,
	// when both rooms are created,
	// then their IDs are allocated in insert order.

	// Given
	expectedIDs := []int{1, 2}
	db := createRoomsTestDB(t)
	firstRoom := roomFixture("Hakkebakken", "101", 1)
	secondRoom := roomFixture("Tangerud", "201", 2)

	// When
	firstCreated, firstErr := CreateRoom(db, firstRoom)
	secondCreated, secondErr := CreateRoom(db, secondRoom)

	// Then
	if firstErr != nil {
		t.Fatalf("expected first room creation to succeed: %v", firstErr)
	}
	if secondErr != nil {
		t.Fatalf("expected second room creation to succeed: %v", secondErr)
	}
	actualIDs := []int{firstCreated.ID, secondCreated.ID}
	if !slices.Equal(expectedIDs, actualIDs) {
		t.Fatalf("room IDs mismatch\nexpected: %v\nactual:   %v", expectedIDs, actualIDs)
	}
}

func TestCreateRoom_WhenRoomNumberDoesNotMatchFloor_ReturnsError(t *testing.T) {
	// Given room input where the room number starts with another floor,
	// when the room is created,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	invalidRoom := roomFixture("Tangerud", "203", 3)

	// When
	_, err := CreateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestCreateRoom_WhenMaxConcurrentGamesIsInvalid_ReturnsError(t *testing.T) {
	// Given room input without capacity for any events,
	// when the room is created,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	invalidRoom := roomFixture("Hakkebakken", "101", 1)
	invalidRoom.MaxConcurrentGames = 0

	// When
	_, err := CreateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestDeleteRoom_RemovesOnlyTargetRoom(t *testing.T) {
	// Given two stored rooms,
	// when one room is deleted,
	// then only the other room remains.

	// Given
	expectedRemainingRoomIDs := []int{2}
	db := createRoomsTestDB(t)
	roomToDelete := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	insertRoom(t, db, roomFixture("Tangerud", "201", 2))

	// When
	err := DeleteRoom(db, roomToDelete.ID)

	// Then
	if err != nil {
		t.Fatalf("expected room deletion to succeed: %v", err)
	}
	actualRemainingRoomIDs := queryRoomIDs(t, db)
	if !slices.Equal(expectedRemainingRoomIDs, actualRemainingRoomIDs) {
		t.Fatalf("remaining room IDs mismatch\nexpected: %v\nactual:   %v", expectedRemainingRoomIDs, actualRemainingRoomIDs)
	}
}

func TestDeleteRoom_WhenRoomDoesNotExist_ReturnsError(t *testing.T) {
	// Given an empty room table,
	// when a missing room is deleted,
	// then the caller receives an error.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	err := DeleteRoom(db, 999)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoom_UpdatesAllFieldsWithoutChangingID(t *testing.T) {
	// Given an existing room and replacement room data,
	// when the room is updated,
	// then every mutable field changes while the ID stays the same.

	// Given
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	expectedRoom := models.Room{
		ID:                 existingRoom.ID,
		Name:               "Tangerud",
		RoomNumber:         "209",
		Floor:              2,
		MaxConcurrentGames: 3,
		Notes:              "Dette er ei oppdatert note",
		IsDisabled:         true,
	}

	// When
	actualRoom, err := UpdateRoom(db, expectedRoom)

	// Then
	if err != nil {
		t.Fatalf("expected room update to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestUpdateRoom_WhenRoomNumberIsEmpty_ReturnsError(t *testing.T) {
	// Given replacement room data with an empty room number,
	// when the room is updated,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	invalidRoom := existingRoom
	invalidRoom.RoomNumber = ""

	// When
	_, err := UpdateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoom_WhenRoomNumberDoesNotMatchFloor_ReturnsError(t *testing.T) {
	// Given replacement room data whose number starts with another floor,
	// when the room is updated,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	invalidRoom := existingRoom
	invalidRoom.RoomNumber = "203"
	invalidRoom.Floor = 3

	// When
	_, err := UpdateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoom_WhenMaxConcurrentGamesIsInvalid_ReturnsError(t *testing.T) {
	// Given replacement room data with invalid event capacity,
	// when the room is updated,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	invalidRoom := existingRoom
	invalidRoom.MaxConcurrentGames = 0

	// When
	_, err := UpdateRoom(db, invalidRoom)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_UpdatesProvidedFields(t *testing.T) {
	// Given an existing room and partial input for every mutable field,
	// when the partial update runs,
	// then the returned room contains all supplied values.

	// Given
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	expectedRoom := models.Room{
		ID:                 existingRoom.ID,
		Name:               "Tangerud",
		RoomNumber:         "303",
		Floor:              3,
		MaxConcurrentGames: 3,
		Notes:              "",
		IsDisabled:         true,
	}
	input := models.RoomInput{
		ID:                 existingRoom.ID,
		Name:               ptr(expectedRoom.Name),
		RoomNumber:         ptr(expectedRoom.RoomNumber),
		Floor:              ptr(expectedRoom.Floor),
		MaxConcurrentGames: ptr(expectedRoom.MaxConcurrentGames),
		Notes:              ptr(expectedRoom.Notes),
		IsDisabled:         ptr(expectedRoom.IsDisabled),
	}

	// When
	actualRoom, err := UpdateRoomPartial(db, input)

	// Then
	if err != nil {
		t.Fatalf("expected partial room update to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestUpdateRoomPartial_WhenOnlyNameIsProvided_LeavesOtherFieldsUnchanged(t *testing.T) {
	// Given an existing room and partial input with only a new name,
	// when the partial update runs,
	// then only the name changes.

	// Given
	db := createRoomsTestDB(t)
	existingRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	expectedRoom := existingRoom
	expectedRoom.Name = "Tangerud"
	input := models.RoomInput{
		ID:   existingRoom.ID,
		Name: ptr(expectedRoom.Name),
	}

	// When
	actualRoom, err := UpdateRoomPartial(db, input)

	// Then
	if err != nil {
		t.Fatalf("expected partial room update to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestUpdateRoomPartial_WhenIDIsMissing_ReturnsError(t *testing.T) {
	// Given partial room input without a room ID,
	// when the partial update runs,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_WhenNameIsEmpty_ReturnsError(t *testing.T) {
	// Given partial room input with an empty name,
	// when the partial update runs,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{ID: 1, Name: ptr("")})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_WhenRoomNumberIsEmpty_ReturnsError(t *testing.T) {
	// Given partial room input with an empty room number,
	// when the partial update runs,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{ID: 1, RoomNumber: ptr("")})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_WhenMaxConcurrentGamesIsInvalid_ReturnsError(t *testing.T) {
	// Given partial room input with invalid event capacity,
	// when the partial update runs,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{ID: 1, MaxConcurrentGames: ptr(0)})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestUpdateRoomPartial_WhenNoFieldsAreProvided_ReturnsError(t *testing.T) {
	// Given partial room input with an ID but no updated fields,
	// when the partial update runs,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := UpdateRoomPartial(db, models.RoomInput{ID: 1})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestGetRoomByID_ReturnsStoredRoom(t *testing.T) {
	// Given a stored room,
	// when it is fetched by ID,
	// then the matching room is returned.

	// Given
	db := createRoomsTestDB(t)
	expectedRoom := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))

	// When
	actualRoom, err := GetRoomByID(db, expectedRoom.ID)

	// Then
	if err != nil {
		t.Fatalf("expected room lookup to succeed: %v", err)
	}
	assertRoomMatches(t, expectedRoom, *actualRoom)
}

func TestGetRoomByID_WhenIDIsInvalid_ReturnsError(t *testing.T) {
	// Given an invalid room ID,
	// when it is fetched,
	// then validation rejects it.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := GetRoomByID(db, 0)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestGetRoomByID_WhenRoomDoesNotExist_ReturnsError(t *testing.T) {
	// Given a positive room ID with no stored room,
	// when it is fetched,
	// then the caller receives an error.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := GetRoomByID(db, 999)
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func TestGetAllRooms_ReturnsRoomsOrderedByFloorAndNumber(t *testing.T) {
	// Given rooms inserted outside display order,
	// when all rooms are listed,
	// then rooms are ordered by floor and room number.

	// Given
	expectedRoomNumbers := []string{"101", "102", "201"}
	db := createRoomsTestDB(t)
	insertRoom(t, db, roomFixture("Tangerud", "201", 2))
	insertRoom(t, db, roomFixture("Brumms hus", "102", 1))
	insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))

	// When
	actualRooms, err := GetAllRooms(db)

	// Then
	if err != nil {
		t.Fatalf("expected room listing to succeed: %v", err)
	}
	actualRoomNumbers := roomNumbers(actualRooms)
	if !slices.Equal(expectedRoomNumbers, actualRoomNumbers) {
		t.Fatalf("room order mismatch\nexpected: %v\nactual:   %v", expectedRoomNumbers, actualRoomNumbers)
	}
}

func TestGetAllRoomStatusesByPulje_ReturnsRoomsAndPuljeAssignments(t *testing.T) {
	// Given rooms, puljer, and event assignments across puljer,
	// when room statuses are listed,
	// then every pulje has every room with only its assigned events.

	// Given
	db := createRoomsTestDB(t)
	seedRoomEventLookups(t, db)
	roomOne := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	roomTwo := insertRoom(t, db, roomFixture("Tangerud", "201", 2))
	fridayPulje := insertPulje(t, db, models.Pulje("Friday"), "Fredag kveld")
	saturdayPulje := insertPulje(t, db, models.Pulje("Saturday"), "Laurdag")
	alphaEvent := insertEvent(t, db, "alpha-event", "Alpha Event", 5)
	betaEvent := insertEvent(t, db, "beta-event", "Beta Event", 4)
	gammaEvent := insertEvent(t, db, "gamma-event", "Gamma Event", 6)
	insertEventPulje(t, db, alphaEvent, fridayPulje, sql.NullInt64{Int64: int64(roomOne.ID), Valid: true})
	insertEventPulje(t, db, betaEvent, fridayPulje, sql.NullInt64{Int64: int64(roomOne.ID), Valid: true})
	insertEventPulje(t, db, gammaEvent, saturdayPulje, sql.NullInt64{Int64: int64(roomTwo.ID), Valid: true})
	expectedAssignments := map[models.Pulje]map[int64][]string{
		fridayPulje: {
			int64(roomOne.ID): []string{"Alpha Event", "Beta Event"},
			int64(roomTwo.ID): []string{},
		},
		saturdayPulje: {
			int64(roomOne.ID): []string{},
			int64(roomTwo.ID): []string{"Gamma Event"},
		},
	}

	// When
	actualStatuses, err := GetAllRoomStatusesByPulje(db, fridayPulje)

	// Then
	if err != nil {
		t.Fatalf("expected room status listing to succeed: %v", err)
	}
	assertRoomStatusAssignments(t, expectedAssignments, actualStatuses)
}

func TestAssignRoomToRelationEventPuljer_AssignsRoomToEventPulje(t *testing.T) {
	// Given an event pulje relation without a room,
	// when a room is assigned to the event,
	// then the relation stores that room.

	// Given
	db := createRoomsTestDB(t)
	seedRoomEventLookups(t, db)
	room := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	puljeID := insertPulje(t, db, models.Pulje("Friday"), "Fredag kveld")
	eventID := insertEvent(t, db, "alpha-event", "Alpha Event", 5)
	insertEventPulje(t, db, eventID, puljeID, sql.NullInt64{})
	expectedEventPulje := models.EventPulje{
		EventID:     eventID,
		PuljeID:     puljeID,
		IsInPulje:   true,
		IsPublished: false,
		RoomID:      sql.NullInt64{Int64: int64(room.ID), Valid: true},
	}

	// When
	actualEventPulje, err := AssignRoomToRelationEventPuljer(db, int64(room.ID), eventID)

	// Then
	if err != nil {
		t.Fatalf("expected room assignment to succeed: %v", err)
	}
	if actualEventPulje != expectedEventPulje {
		t.Fatalf("event pulje mismatch\nexpected: %+v\nactual:   %+v", expectedEventPulje, actualEventPulje)
	}
}

func TestAssignRoomToRelationEventPuljer_WhenRelationDoesNotExist_ReturnsError(t *testing.T) {
	// Given no event pulje relation for an event,
	// when a room is assigned to that event,
	// then the caller receives an error.

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := AssignRoomToRelationEventPuljer(db, 1, "missing-event")
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

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
