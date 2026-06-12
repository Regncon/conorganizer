package rooms

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
)

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
