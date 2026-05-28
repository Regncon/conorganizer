package rooms

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestCreateRoom(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	var validRoom = models.Room{
		ID:                 0,
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}

	var invalidRoomNumber = validRoom
	invalidRoomNumber.Floor = 2

	var invalidRoomConcurrency = validRoom
	invalidRoomConcurrency.MaxConcurrentGames = 0

	// When
	// - create room is called once
	createRoomResult, err := CreateRoom(db, validRoom)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// - create room is called again, tesing auto increment
	createRoomResult2, err := CreateRoom(db, validRoom)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// - create room is called with an invalid room number based on the floor being its prefix
	_, errInvalidRoomNumber := CreateRoom(db, invalidRoomNumber)

	// - create room is called with an invalid room number based on the floor being its prefix
	_, errInvalidConcurrency := CreateRoom(db, invalidRoomConcurrency)

	// Then
	if validRoom.ID == createRoomResult.ID {
		t.Fatalf("Room ID did not update correctly after insert\nexpected: %d\nrecieved: %d", validRoom.ID, createRoomResult.ID)
	}

	validRoom.ID = createRoomResult.ID
	if validRoom != *createRoomResult {
		t.Fatalf("createRoomResult did not match happyRoom\nexpected: \t%v\nrecieved: \t%v", validRoom, createRoomResult)
	}

	var expectedID = createRoomResult.ID + 1
	if createRoomResult2.ID != expectedID {
		t.Fatalf("createRoom did not auto increment ID\nexpected: %d\nrecieved: %d", expectedID, createRoomResult2.ID)
	}

	if errInvalidRoomNumber == nil {
		t.Fatal("expected error when creating room with invalid room number, but it was allowed")
	}

	if errInvalidConcurrency == nil {
		t.Fatal("expected error when creating room with 0 or less max games, but it was allowed")
	}
}

func TestUpdateRoom(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	var originalRoom = models.Room{
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}

	resultCreateRoom, err := CreateRoom(db, originalRoom)
	if err != nil {
		t.Fatalf("unexpected error creating original room: %v", err)
	}

	var updatedRoom = models.Room{
		ID:                 resultCreateRoom.ID,
		Name:               "Tangerud",
		RoomNumber:         "209",
		Floor:              2,
		MaxConcurrentGames: 3,
		Notes:              "Dette er en oppdatert note",
		IsDisabled:         true,
	}

	var updatedRoomInvalidRoomNumber = *resultCreateRoom
	updatedRoomInvalidRoomNumber.RoomNumber = ""

	var updatedRoomInvalidRoomConcurrency = *resultCreateRoom
	updatedRoomInvalidRoomConcurrency.MaxConcurrentGames = -1

	var updatedRoomInvalidRoomNumberFloor = *resultCreateRoom
	updatedRoomInvalidRoomNumberFloor.Floor = 3
	updatedRoomInvalidRoomNumberFloor.RoomNumber = "203"

	// When
	resultUpdateRoomValid, err := UpdateRoom(db, updatedRoom)
	if err != nil {
		t.Fatalf("unexpected error when updating valid room: %v", err)
	}

	_, errInvalidRoomNumber := UpdateRoom(db, updatedRoomInvalidRoomNumber)
	_, errInvalidRoomConcurrency := UpdateRoom(db, updatedRoomInvalidRoomConcurrency)
	_, errInvalidRoomNumberFloor := UpdateRoom(db, updatedRoomInvalidRoomNumberFloor)

	// Then
	if resultCreateRoom.ID != resultUpdateRoomValid.ID {
		t.Fatalf("UpdateRoom caused room ID to change\nexpected: \t%d\nrecieved: \t%d", resultCreateRoom.ID, resultUpdateRoomValid.ID)
	}
	if resultCreateRoom.Name == resultUpdateRoomValid.Name {
		t.Fatalf("UpdateRoom did not update name correctly\nexpected: \t%s\nrecieved: \t%s", resultCreateRoom.Name, resultUpdateRoomValid.Name)
	}
	if resultCreateRoom.RoomNumber == resultUpdateRoomValid.RoomNumber {
		t.Fatalf("UpdateRoom did not update room number correctly\nexpected: \t%s\nrecieved: \t%s", resultCreateRoom.RoomNumber, resultUpdateRoomValid.RoomNumber)
	}
	if resultCreateRoom.Floor == resultUpdateRoomValid.Floor {
		t.Fatalf("UpdateRoom did not update floor number correctly\nexpected: \t%d\nrecieved: \t%d", resultCreateRoom.Floor, resultUpdateRoomValid.Floor)
	}
	if resultCreateRoom.Notes == resultUpdateRoomValid.Notes {
		t.Fatalf("UpdateRoom did not update notes correctly\nexpected: \t%s\nrecieved: \t%s", resultCreateRoom.Notes, resultUpdateRoomValid.Notes)
	}
	if resultCreateRoom.MaxConcurrentGames == resultUpdateRoomValid.MaxConcurrentGames {
		t.Fatalf("UpdateRoom did not update max concurrent games correctly\nexpected: \t%d\nrecieved: \t%d", resultCreateRoom.MaxConcurrentGames, resultUpdateRoomValid.MaxConcurrentGames)
	}

	if errInvalidRoomNumber == nil {
		t.Fatalf("UpdateRoom allowed updating invalid room number")
	}
	if errInvalidRoomNumberFloor == nil {
		t.Fatalf("UpdateRoom allowed updating with an invalid room number and floor combination")
	}
	if errInvalidRoomConcurrency == nil {
		t.Fatalf("UpdateRoom allowed updating invalid max concurrent games")
	}
}

func TestUpdateRoomPartial(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	var originalRoom = models.Room{
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}
	var originalRoomPartial = models.Room{
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}

	createRoomResult, err := CreateRoom(db, originalRoom)
	if err != nil {
		t.Fatalf("unexpected error creating original room: %v", err)
	}
	createRoomPartialResult, err := CreateRoom(db, originalRoomPartial)
	if err != nil {
		t.Fatalf("unexpected error creating original room: %v", err)
	}

	// When
	var updatedName string = "Tangerud"
	var updatedRoomNumber string = "303"
	var updatedFloor int = 3
	var updatedConcurrent int = 3
	var updatedNotes string = ""
	var updatedDisables bool = true

	var updatedRoom = models.RoomInput{
		ID:                 createRoomResult.ID,
		Name:               &updatedName,
		RoomNumber:         &updatedRoomNumber,
		Floor:              &updatedFloor,
		MaxConcurrentGames: &updatedConcurrent,
		Notes:              &updatedNotes,
		IsDisabled:         &updatedDisables,
	}
	updatedRoomResult, err := UpdateRoomPartial(db, updatedRoom)
	if err != nil {
		t.Fatalf("unexpected error when updating valid room: %v", err)
	}
	partialUpdatedRoomResult, err := UpdateRoomPartial(db, models.RoomInput{ID: createRoomPartialResult.ID, Name: &updatedName})
	if err != nil {
		t.Fatalf("unexpected error when partially updating valid room: %v", err)
	}

	var invalidRoomNumber string = ""
	var invalidName string = ""
	var invalidConcurrent int = -1

	_, errInvalidID := UpdateRoomPartial(db, models.RoomInput{})
	_, errInvalidName := UpdateRoomPartial(db, models.RoomInput{ID: 1, Name: &invalidName})
	_, errInvalidRoomNumber := UpdateRoomPartial(db, models.RoomInput{ID: 1, RoomNumber: &invalidRoomNumber})
	_, errInvalidConcurrency := UpdateRoomPartial(db, models.RoomInput{ID: 1, MaxConcurrentGames: &invalidConcurrent})

	// Then
	if originalRoom.Name != createRoomResult.Name {
		t.Errorf("Original name was different to what create room returned")
	}
	if createRoomResult.Name == updatedRoomResult.Name {
		t.Errorf("Room name persisted after updateRoom was called successfully\nexpected: \t%s\nrecieved: \t%s", createRoomResult.Name, updatedRoomResult.Name)
	}
	if createRoomPartialResult.Name == partialUpdatedRoomResult.Name {
		t.Errorf("Room name persisted after updateRoom was called successfully\nexpected: \t%s\nrecieved: \t%s", createRoomPartialResult.Name, partialUpdatedRoomResult.Name)
	}

	if errInvalidID == nil {
		t.Errorf("UpdateRoom allowed update when ID was omited")
	}
	if errInvalidName == nil {
		t.Errorf("UpdateRoom allowed update when name was an empty string")
	}
	if errInvalidConcurrency == nil {
		t.Errorf("UpdateRoom allowed update when max concurrent games was less than 1")
	}
	if errInvalidRoomNumber == nil {
		t.Errorf("UpdateRoom allowed update when room number was an empty string")
	}
}

func TestGetRoomByID(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	var validRoom = models.Room{
		Name:               "Hakkebakken",
		RoomNumber:         "101",
		Floor:              1,
		MaxConcurrentGames: 2,
		Notes:              "Dette er et gyldig rom",
		IsDisabled:         false,
	}

	createdRoom, err := CreateRoom(db, validRoom)
	if err != nil {
		t.Fatalf("unexpected error creating room: %v", err)
	}

	// When
	roomResult, err := GetRoomByID(db, createdRoom.ID)
	if err != nil {
		t.Fatalf("unexpected error getting room by ID: %v", err)
	}

	_, errInvalidID := GetRoomByID(db, 0)
	_, errMissingRoom := GetRoomByID(db, -1)

	// Then
	if *roomResult != *createdRoom {
		t.Fatalf(
			"GetRoomByID did not return expected room\nexpected:\t%v\nrecieved:\t%v",
			*createdRoom,
			*roomResult,
		)
	}

	if errInvalidID == nil {
		t.Fatal("expected error when getting room with invalid ID, but it was allowed")
	}
	if errMissingRoom == nil {
		t.Fatal("expected error when getting non-existing room, but it was allowed")
	}
}

func TestGetAllRooms(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	var validRooms = []models.Room{
		{
			Name:               "Hakkebakken",
			RoomNumber:         "101",
			Floor:              1,
			MaxConcurrentGames: 2,
			Notes:              "Room 1",
			IsDisabled:         false,
		}, {
			Name:               "Tangerud",
			RoomNumber:         "201",
			Floor:              2,
			MaxConcurrentGames: 4,
			Notes:              "Second room",
			IsDisabled:         false,
		}, {
			Name:               "Hundremeter Skogen",
			RoomNumber:         "301",
			Floor:              3,
			MaxConcurrentGames: 1,
			Notes:              "Second room",
			IsDisabled:         false,
		},
	}

	var createdRooms []models.Room
	for _, room := range validRooms {
		createdRoom, err := CreateRoom(db, room)
		if err != nil {
			t.Fatalf("unexpected error creating room: %v", err)
		}

		createdRooms = append(createdRooms, *createdRoom)
	}

	// When
	resultRooms, err := GetAllRooms(db)
	if err != nil {
		t.Fatalf("unexpected error getting all rooms: %v", err)
	}

	// Then
	// - ensure all rooms were returned
	if len(resultRooms) != 3 {
		t.Fatalf(
			"expected 3 rooms, recieved: %d",
			len(resultRooms),
		)
	}

	// - ensure ordering is correct
	if resultRooms[0].ID != createdRooms[0].ID {
		t.Fatalf(
			"expected first room ID %d, recieved %d",
			createdRooms[0].ID,
			resultRooms[0].ID,
		)
	}
	if resultRooms[len(resultRooms)-1].ID != createdRooms[len(createdRooms)-1].ID {
		t.Fatalf(
			"expected last room ID %d, recieved %d",
			createdRooms[len(createdRooms)-1].ID,
			resultRooms[len(resultRooms)-1].ID,
		)
	}
}

func TestGetAllRoomStatusesByPulje(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services_event_puljer", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	// Seed databases with required data for relations
	rooms := insertRooms(t, db)
	puljer := insertPuljer(t, db)
	events := insertEvents(t, db)

	query := `
        INSERT INTO relation_event_puljer (
            event_id,
            pulje_id,
            room_id
        )
        VALUES (?, ?, ?)
        RETURNING
            event_id,
            pulje_id,
            is_in_pulje,
            is_published,
            room_id
    `

	var eventPuljerSource = []models.EventPulje{
		{
			EventID: events[0],
			PuljeID: models.Pulje(puljer[0]),
			RoomID:  sql.NullInt64{Int64: int64(rooms[0]), Valid: true},
		}, {
			EventID: events[1],
			PuljeID: models.Pulje(puljer[0]),
			RoomID:  sql.NullInt64{Int64: int64(rooms[0]), Valid: true},
		}, {
			EventID: events[2],
			PuljeID: models.Pulje(puljer[1]),
			RoomID:  sql.NullInt64{Int64: int64(rooms[1]), Valid: true},
		}, {
			EventID: events[3],
			PuljeID: models.Pulje(puljer[1]),
			RoomID:  sql.NullInt64{Int64: int64(rooms[2]), Valid: true},
		}, {
			EventID: events[4],
			PuljeID: models.Pulje(puljer[2]),
			RoomID:  sql.NullInt64{},
		},
	}

	var createdEventPuljer []models.EventPulje
	for _, eventSource := range eventPuljerSource {
		var createdEvent models.EventPulje

		err := db.QueryRow(
			query,
			eventSource.EventID,
			eventSource.PuljeID,
			eventSource.RoomID,
		).Scan(
			&createdEvent.EventID,
			&createdEvent.PuljeID,
			&createdEvent.IsInPulje,
			&createdEvent.IsPublished,
			&createdEvent.RoomID,
		)

		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}
		createdEventPuljer = append(createdEventPuljer, createdEvent)
	}

	var expectedRoomStatuses = make(models.RoomStatusByPulje)
	for _, puljeID := range puljer {
		pulje := models.Pulje(puljeID)

		if expectedRoomStatuses[pulje] == nil {
			expectedRoomStatuses[pulje] = make(map[int64]models.RoomByPulje)
		}

		for _, roomID := range rooms {
			expectedRoomStatuses[pulje][int64(roomID)] = models.RoomByPulje{
				ID:                 roomID,
				Name:               "",
				RoomNumber:         "",
				MaxConcurrentGames: 0,
				Notes:              "",
				AssignedEventsID:   []models.RoomEventPuljeSummary{},
			}
		}
	}
	for _, eventPulje := range createdEventPuljer {
		if !eventPulje.IsInPulje || !eventPulje.RoomID.Valid {
			continue
		}

		pulje := eventPulje.PuljeID
		roomID := eventPulje.RoomID.Int64

		room := expectedRoomStatuses[pulje][roomID]

		room.AssignedEventsID = append(room.AssignedEventsID, models.RoomEventPuljeSummary{
			EventID: eventPulje.EventID,
		})

		expectedRoomStatuses[pulje][roomID] = room
	}

	// When
	result, err := GetAllRoomStatusesByPulje(db, models.PuljeFredagKveld)
	if err != nil {
		t.Fatalf("Unexpected arror: %v", err)
	}

	// Then
	for pulje, rooms := range expectedRoomStatuses {
		for roomID, expectedRoom := range rooms {

			actualRoom := result[pulje][roomID]

			if len(actualRoom.AssignedEventsID) != len(expectedRoom.AssignedEventsID) {
				t.Errorf(
					"pulje=%s room=%d expected %d events, got %d",
					pulje,
					roomID,
					len(expectedRoom.AssignedEventsID),
					len(actualRoom.AssignedEventsID),
				)
			}
		}
	}
}

func TestAssignRoomToRelationEventPuljer(t *testing.T) {
	// Given
	db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services_assignment", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	// Seed databases with required data for relations
	rooms := insertRooms(t, db)
	puljer := insertPuljer(t, db)
	events := insertEvents(t, db)

	// Simplified relatinal insert
	eventID := events[0]
	puljeID := puljer[0]

	_, err = db.Exec(`
		INSERT INTO relation_event_puljer (event_id, pulje_id)
		VALUES (?, ?)
	`, eventID, puljeID)
	if err != nil {
		t.Fatalf("failed to insert relation: %v", err)
	}

	// When
	result, err := AssignRoomToRelationEventPuljer(db, rooms[0], eventID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Thenn
	if result.EventID != eventID {
		t.Fatalf("Event ID did not match between input and ouptut\nexpected:\t%s\nrecieved:\t%s", eventID, result.EventID)
	}
	if result.PuljeID != models.Pulje(puljeID) {
		t.Fatalf("Pulje ID did not match between input and ouptut\nexpected:\t%s\nrecieved:\t%s", models.Pulje(puljeID), result.PuljeID)
	}

	if !result.RoomID.Valid {
		t.Fatalf("expected room_id to be set")
	}

	if result.RoomID.Int64 != int64(rooms[0]) {
		t.Fatalf(
			"Room number did not update\nexpected:\t%d\nrecieved:\t%d",
			rooms[0],
			result.RoomID.Int64,
		)
	}
}

func insertPuljer(t *testing.T, db *sql.DB) []string {
	t.Helper()

	puljerQuery := `
        INSERT INTO puljer (
			id, name, status, start_at, end_at
		) VALUES
			('Friday', 'Fredag kveld', 'Open', '2025-10-03', '2025-10-03'),
			('SaturdayMorning', 'Lørdag morgen', 'Open', '2025-10-04', '2025-10-04'),
			('SaturdayEvening', 'Lørdag kveld', 'Open', '2025-10-04', '2025-10-04'),
			('Sunday', 'Søndag morgen', 'Open', '2025-10-05', '2025-10-05')
        RETURNING id
	`
	rows, err := db.Query(puljerQuery)
	if err != nil {
		t.Fatalf("failed to insert puljer: %v", err)
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("failed to scan pulje id: %v", err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("row iteration failed: %v", err)
	}

	return ids
}

func insertRooms(t *testing.T, db *sql.DB) []int64 {
	t.Helper()

	rooms := []models.Room{
		{
			Name:               "Hundremeterskogen",
			RoomNumber:         "101",
			Floor:              1,
			MaxConcurrentGames: 2,
			Notes:              "Rom inspirert av skogen der Ole Brumm bor",
			IsDisabled:         false,
		},
		{
			Name:               "Brumms Hus",
			RoomNumber:         "102",
			Floor:              1,
			MaxConcurrentGames: 3,
			Notes:              "Koselig rom med plass til flere spill",
			IsDisabled:         false,
		},
		{
			Name:               "Tigerguttens Hjorne",
			RoomNumber:         "201",
			Floor:              2,
			MaxConcurrentGames: 2,
			Notes:              "Aktivt rom for mindre grupper",
			IsDisabled:         false,
		},
		{
			Name:               "Nasse Noffs Sti",
			RoomNumber:         "202",
			Floor:              2,
			MaxConcurrentGames: 1,
			Notes:              "Lite og stille rom",
			IsDisabled:         false,
		},
		{
			Name:               "Ugles Topp",
			RoomNumber:         "301",
			Floor:              3,
			MaxConcurrentGames: 4,
			Notes:              "Stort rom egnet for parallelle aktiviteter, men er inaktivt",
			IsDisabled:         true,
		},
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
            RETURNING id
        `
	var roomIDs []int64
	for _, room := range rooms {
		var roomID int64
		err := db.QueryRow(query,
			room.Name,
			room.RoomNumber,
			room.Floor,
			room.MaxConcurrentGames,
			room.Notes,
			room.IsDisabled,
		).Scan(&roomID)
		if err != nil {
			t.Fatalf("Failed to create room: %v", err)
		}
		roomIDs = append(roomIDs, roomID)
	}

	return roomIDs
}

func insertEvents(t *testing.T, db *sql.DB) []string {
	events := []models.Event{
		{
			Title:             "Mysteriet i Hundremeterskogen",
			Intro:             "Et rolig mysterium for nye spillere.",
			Description:       "Spillerne må finne ut hvorfor honningkrukkene til Ole Brumm forsvinner om natten.",
			System:            "Call of Cthulhu",
			EventType:         models.EventTypeBoardGame,
			AgeGroup:          models.AgeGroupAdultsOnly,
			Runtime:           models.RunTimeLongRunning,
			HostName:          "Kristoffer",
			Email:             "brumm@example.com",
			PhoneNumber:       "90000001",
			MaxPlayers:        5,
			BeginnerFriendly:  true,
			CanBeRunInEnglish: true,
			Notes:             "Passer godt for førstegangsspillere",
			Status:            "Publisert",
		},
		{
			Title:             "Tigerguttens Turnering",
			Intro:             "En energisk konkurranse med raske utfordringer.",
			Description:       "Deltakerne konkurrerer i kreative oppgaver og samarbeid under press.",
			System:            "Dungeons & Dragons 5e",
			EventType:         models.EventTypeBoardGame,
			AgeGroup:          models.AgeGroupAdultsOnly,
			Runtime:           models.RunTimeLongRunning,
			HostName:          "Ole",
			Email:             "tiger@example.com",
			PhoneNumber:       "90000002",
			MaxPlayers:        6,
			BeginnerFriendly:  true,
			CanBeRunInEnglish: false,
			Notes:             "",
			Status:            "Kladd",
		},
		{
			Title:             "Nasse Noffs Mørke Skog",
			Intro:             "Et skrekkeventyr i dype skoger.",
			Description:       "Noe beveger seg mellom trærne, og spillerne må overleve natten.",
			System:            "Vaesen",
			EventType:         models.EventTypeBoardGame,
			AgeGroup:          models.AgeGroupAdultsOnly,
			Runtime:           models.RunTimeLongRunning,
			HostName:          "Anne",
			Email:             "nasse@example.com",
			PhoneNumber:       "90000003",
			MaxPlayers:        4,
			BeginnerFriendly:  false,
			CanBeRunInEnglish: true,
			Notes:             "Inneholder skrekkelementer",
			Status:            "Publisert",
		},
		{
			Title:             "Ugles Kunnskapsprove",
			Intro:             "Quiz og strategi i kombinasjon.",
			Description:       "Spillerne må samarbeide for å løse gåter og vinne over Ugle.",
			System:            "Custom",
			EventType:         models.EventTypeBoardGame,
			AgeGroup:          models.AgeGroupAdultsOnly,
			Runtime:           models.RunTimeLongRunning,
			HostName:          "Mari",
			Email:             "ugle@example.com",
			PhoneNumber:       "90000004",
			MaxPlayers:        8,
			BeginnerFriendly:  true,
			CanBeRunInEnglish: true,
			Notes:             "",
			Status:            "Publisert",
		},
		{
			Title:             "Kengus Eventyrreise",
			Intro:             "Et familievennlig fantasy-eventyr.",
			Description:       "Bli med Kengu og Ro på en reise gjennom magiske landskap.",
			System:            "Pathfinder 2e",
			EventType:         models.EventTypeBoardGame,
			AgeGroup:          models.AgeGroupAdultsOnly,
			Runtime:           models.RunTimeLongRunning,
			HostName:          "Sindre",
			Email:             "kengu@example.com",
			PhoneNumber:       "90000005",
			MaxPlayers:        5,
			BeginnerFriendly:  true,
			CanBeRunInEnglish: false,
			Notes:             "Familievennlig innhold",
			Status:            "Godkjent",
		},
	}

	query := `
        INSERT INTO events (
            title,
            intro,
            description,
            system,
            host_name,
            user_id,
            created_by_id,
            updated_by_id,
            email,
            phone_number,
            max_players,
            age_group,
            event_runtime,
            beginner_friendly,
            can_be_run_in_english,
            status
        ) VALUES (
            ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
        ) RETURNING id`

	var eventIDs []string
	for _, event := range events {
		var eventID string
		err := db.QueryRow(query,
			event.Title,
			event.Intro,
			event.Description,
			event.System,
			event.HostName,
			event.UserID,
			event.CreatedByID,
			event.UpdatedByID,
			event.Email,
			event.PhoneNumber,
			event.MaxPlayers,
			event.AgeGroup,
			event.Runtime,
			event.BeginnerFriendly,
			event.CanBeRunInEnglish,
			event.Status,
		).Scan(&eventID)
		if err != nil {
			t.Fatalf("Failed to create event: %v", err)
		}
		eventIDs = append(eventIDs, eventID)
	}

	return eventIDs
}
