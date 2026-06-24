package rooms

import (
	"database/sql"
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	roomService "github.com/Regncon/conorganizer/service/rooms"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestGetRoomsByFloor_ReturnsFloorsInDescendingOrder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt rom i flere etasjer.",
		When:  "Når romoversikten grupperer rom etter etasje.",
		Then:  "Så skal etasjene vises ovenfra og ned.",
	})

	// Given
	expectedFloors := []int{3, 2, 1}
	db, logger := testutil.CreateTestDBAndLogger(t, "rooms_page_floors")
	createRoomsPageRoom(t, db, "Hakkebakken", "101", 1)
	createRoomsPageRoom(t, db, "Tangerud", "201", 2)
	createRoomsPageRoom(t, db, "Topprommet", "301", 3)

	// When
	floorGroups := getRoomsByFloor(db, logger)
	actualFloors := roomPageFloorIDs(floorGroups)

	// Then
	if !slices.Equal(expectedFloors, actualFloors) {
		t.Fatalf("floor order mismatch\nexpected: %v\nactual:   %v", expectedFloors, actualFloors)
	}
}

func TestRoomsPageContent_RendersRoomDetailsAndCreateAction(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at rom er registrert.",
		When:  "Når romoversikten rendres.",
		Then:  "Så skal romdetaljer og handling for nytt rom være synlige.",
	})

	// Given
	expectedTextParts := []string{
		"201",
		"Tangerud",
		"Maks aktive spill",
		"2",
		"Ja",
		"Ligg til nytt rom",
	}
	db, logger := testutil.CreateTestDBAndLogger(t, "rooms_page_content")
	createRoomsPageRoom(t, db, "Tangerud", "201", 2)

	// When
	doc := templtest.Render(t, RoomsPageContent(db, logger))
	actualText := strings.Join(templtest.CollectTexts(doc, "#room-administration"), " ")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("expected room page text to contain %q\nactual text: %s", expectedTextPart, actualText)
		}
	}
}

func TestRoomsAssignmentPageContent_RendersMissingRoomEventsAndAssignedRooms(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje med ett arrangement uten rom og ett arrangement med rom.",
		When:  "Når romfordelingen rendres.",
		Then:  "Så skal manglende rom og tildelt rom vises hver for seg.",
	})

	// Given
	expectedTextParts := []string{
		"1 Eventer i pulje uten tildelt rom",
		"Missing Room Event",
		"Romfordelig for FredagKveld",
		"Assigned Room Event",
		"201",
	}
	db, logger := testutil.CreateTestDBAndLogger(t, "rooms_assignment_page")
	seedRoomsPageLookups(t, db)
	room := createRoomsPageRoom(t, db, "Tangerud", "201", 2)
	insertRoomsPagePulje(t, db, models.PuljeFredagKveld)
	insertRoomsPageEvent(t, db, "missing-room-event", "Missing Room Event", 4)
	insertRoomsPageEvent(t, db, "assigned-room-event", "Assigned Room Event", 5)
	insertRoomsPageEventPulje(t, db, "missing-room-event", models.PuljeFredagKveld, 0)
	insertRoomsPageEventPulje(t, db, "assigned-room-event", models.PuljeFredagKveld, room.ID)

	// When
	doc := templtest.Render(t, RoomsAssignmentPageContent(db, logger, models.PuljeFredagKveld, nil))
	actualText := strings.Join(templtest.CollectTexts(doc, "#room-assignment"), " ")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("expected room assignment text to contain %q\nactual text: %s", expectedTextPart, actualText)
		}
	}

	// Drag-and-drop is wired with Datastar: drop targets @post via the
	// $draggedEventId signal set on dragstart — no custom fetch or headers.
	roomDropTarget := doc.Find(".room")
	if roomDropTarget.Length() == 0 {
		t.Fatal("expected a room drop target")
	}
	wantDrop := "@post('/admin/rooms/api/assignment/FredagKveld/' + $draggedEventId + '/1')"
	if got := roomDropTarget.AttrOr("data-on:drop__prevent", ""); !strings.Contains(got, wantDrop) {
		t.Fatalf("room drop handler mismatch\nexpected to contain: %s\nactual:              %s", wantDrop, got)
	}
	if got := doc.Find(`.room-event[draggable="true"]`).AttrOr("data-on:dragstart", ""); got != "$draggedEventId = 'assigned-room-event'" {
		t.Fatalf("assigned room event card should set the dragged signal on dragstart\nactual: %s", got)
	}
	if got := doc.Find(`.event-list a[draggable="true"]`).AttrOr("data-on:dragstart", ""); got != "$draggedEventId = 'missing-room-event'" {
		t.Fatalf("missing-room event link should set the dragged signal on dragstart\nactual: %s", got)
	}
}

func TestCalculatePopulation_CountsPlayersAndGMs(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt arrangementer med maks antall spillere.",
		When:  "Når romfordelingen beregner omtrentlig personbelastning.",
		Then:  "Så skal den telle spillere pluss en GM per arrangement.",
	})

	// Given
	expectedPopulation := 11
	events := []models.RoomEventPuljeSummary{
		{MaxPlayers: 4},
		{MaxPlayers: 5},
	}

	// When
	actualPopulation := calculatePopulation(events)

	// Then
	if actualPopulation != expectedPopulation {
		t.Fatalf("population mismatch\nexpected: %d\nactual:   %d", expectedPopulation, actualPopulation)
	}
}

func roomPageFloorIDs(floorGroups []FloorGroup) []int {
	floors := make([]int, 0, len(floorGroups))
	for _, floorGroup := range floorGroups {
		floors = append(floors, floorGroup.Floor)
	}
	return floors
}

func createRoomsPageRoom(t *testing.T, db *sql.DB, name string, roomNumber string, floor int) models.Room {
	t.Helper()

	room, err := roomService.CreateRoom(db, models.Room{
		Name:               name,
		RoomNumber:         roomNumber,
		Floor:              floor,
		MaxConcurrentGames: 2,
	})
	if err.HasErrors() {
		t.Fatalf("failed to create room: %v", err)
	}
	return *room
}

func seedRoomsPageLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	testutil.MustExec(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, models.EventStatusAnnounced)
	testutil.MustExec(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeOther)
	testutil.MustExec(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupDefault)
	testutil.MustExec(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeNormal)
	testutil.MustExec(t, db, `INSERT INTO pulje_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, models.PuljeStatusOpen)
}

func insertRoomsPagePulje(t *testing.T, db *sql.DB, puljeID models.Pulje) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES(?, ?, ?, '2026-10-09T18:00:00Z', '2026-10-09T23:00:00Z')
	`, puljeID, string(puljeID), models.PuljeStatusOpen)
}

func insertRoomsPageEvent(t *testing.T, db *sql.DB, eventID string, title string, maxPlayers int) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO events(
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
			status
		)
		VALUES(?, ?, 'Intro', 'Description', 'System', ?, ?, ?, 'Host', 'host@example.com', '12345678', ?, ?)
	`, eventID, title, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, maxPlayers, models.EventStatusAnnounced)
}

func insertRoomsPageEventPulje(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, roomID int) {
	t.Helper()

	if roomID == 0 {
		testutil.MustExec(t, db, `
			INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published, room_id)
			VALUES(?, ?, 1, 1, NULL)
		`, eventID, puljeID)
		return
	}

	testutil.MustExec(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published, room_id)
		VALUES(?, ?, 1, 1, ?)
	`, eventID, puljeID, roomID)
}
