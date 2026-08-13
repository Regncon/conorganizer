package rooms

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestGetAllRoomStatusesByPulje_ReturnsRoomsAndPuljeAssignments(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given rooms, puljer, and event assignments across puljer.",
		When:  "When room statuses are listed.",
		Then:  "Then every pulje has every room with only its assigned events.",
	})

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

func TestAssignRoomToEventPulje_AssignsRoomToEventPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given an event pulje relation without a room.",
		When:  "When a room is assigned to the event.",
		Then:  "Then the relation stores that room.",
	})

	// Given
	db := createRoomsTestDB(t)
	seedRoomEventLookups(t, db)
	room := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	puljeID := insertPulje(t, db, models.Pulje("Friday"), "Fredag kveld")
	eventID := insertEvent(t, db, "alpha-event", "Alpha Event", 5)
	key := models.EventPuljeKey{EventID: eventID, PuljeID: puljeID}
	insertEventPulje(t, db, eventID, puljeID, sql.NullInt64{})
	expectedEventPulje := models.EventPulje{
		EventID:     eventID,
		PuljeID:     puljeID,
		IsInPulje:   true,
		IsPublished: false,
		RoomID:      sql.NullInt64{Int64: int64(room.ID), Valid: true},
	}

	// When
	actualEventPulje, err := AssignRoomToEventPulje(db, int64(room.ID), key)

	// Then
	if err != nil {
		t.Fatalf("expected room assignment to succeed: %v", err)
	}
	if actualEventPulje != expectedEventPulje {
		t.Fatalf("event pulje mismatch\nexpected: %+v\nactual:   %+v", expectedEventPulje, actualEventPulje)
	}
}

func TestAssignRoomToEventPulje_AssignsOnlyTargetPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given the same event is in two puljer.",
		When:  "When a room is assigned to the event in one pulje.",
		Then:  "Then the other event-pulje relation is left unchanged.",
	})

	// Given
	db := createRoomsTestDB(t)
	seedRoomEventLookups(t, db)
	room := insertRoom(t, db, roomFixture("Hakkebakken", "101", 1))
	fridayPulje := insertPulje(t, db, models.Pulje("Friday"), "Fredag kveld")
	saturdayPulje := insertPulje(t, db, models.Pulje("Saturday"), "Laurdag")
	eventID := insertEvent(t, db, "alpha-event", "Alpha Event", 5)
	targetKey := models.EventPuljeKey{EventID: eventID, PuljeID: fridayPulje}
	otherKey := models.EventPuljeKey{EventID: eventID, PuljeID: saturdayPulje}
	insertEventPulje(t, db, eventID, fridayPulje, sql.NullInt64{})
	insertEventPulje(t, db, eventID, saturdayPulje, sql.NullInt64{})

	// When
	if _, err := AssignRoomToEventPulje(db, int64(room.ID), targetKey); err != nil {
		t.Fatalf("expected room assignment to succeed: %v", err)
	}
	actualTargetRoomID := queryEventPuljeRoomID(t, db, targetKey)
	actualOtherRoomID := queryEventPuljeRoomID(t, db, otherKey)

	// Then
	if actualTargetRoomID != (sql.NullInt64{Int64: int64(room.ID), Valid: true}) {
		t.Fatalf("target room id mismatch\nexpected: %+v\nactual:   %+v", sql.NullInt64{Int64: int64(room.ID), Valid: true}, actualTargetRoomID)
	}
	if actualOtherRoomID.Valid {
		t.Fatalf("other pulje should not have a room assigned\nactual: %+v", actualOtherRoomID)
	}
}

func TestAssignRoomToEventPulje_WhenRelationDoesNotExist_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given no event pulje relation for an event.",
		When:  "When a room is assigned to that event.",
		Then:  "Then the caller receives an error.",
	})

	// Given
	expectedError := true
	db := createRoomsTestDB(t)

	// When
	_, err := AssignRoomToEventPulje(db, 1, models.EventPuljeKey{
		EventID: "missing-event",
		PuljeID: models.Pulje("missing-pulje"),
	})
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func queryEventPuljeRoomID(t testing.TB, db *sql.DB, key models.EventPuljeKey) sql.NullInt64 {
	t.Helper()

	var roomID sql.NullInt64
	err := db.QueryRow(`
		SELECT room_id
		FROM relation_event_puljer
		WHERE event_id = ? AND pulje_id = ?
	`, key.EventID, key.PuljeID).Scan(&roomID)
	if err != nil {
		t.Fatalf("failed to query event-pulje room id: %v", err)
	}
	return roomID
}
