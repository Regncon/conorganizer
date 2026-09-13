package rooms

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	if got := doc.Find(`.room .room-event[draggable="true"]`).AttrOr("data-on:dragstart", ""); got != "$draggedEventId = 'assigned-room-event'" {
		t.Fatalf("assigned room event card should set the dragged signal on dragstart\nactual: %s", got)
	}
	if got := doc.Find(`.event-list .room-event[draggable="true"]`).AttrOr("data-on:dragstart", ""); got != "$draggedEventId = 'missing-room-event'" {
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

func TestRoomCard_ResolvesMapByRoomNumber(t *testing.T) {
	for _, tc := range []struct {
		number string
		path   string
	}{
		{number: "705", path: "/static/rooms/terminus-7-etasje-705.svg"},
		{number: "716"},
		{number: "../705"},
	} {
		t.Run(tc.number, func(t *testing.T) {
			// Given
			room := models.Room{ID: 42, Name: "Et annet navn", RoomNumber: tc.number, Floor: 7}

			// When
			doc := templtest.Render(t, roomCard(room))

			// Then
			if tc.path == "" {
				if doc.Find("img, a.room-map").Length() != 0 || doc.Find(".room-map-missing").Length() != 1 {
					t.Fatal("expected an unavailable-map message without a broken image or link")
				}
			} else {
				if got := doc.Find(".room-map img").AttrOr("src", ""); got != tc.path {
					t.Fatalf("expected map %q, got %q", tc.path, got)
				}
				if got := doc.Find("a.room-map").AttrOr("href", ""); got != tc.path {
					t.Fatalf("expected full-size link %q, got %q", tc.path, got)
				}
			}
			if !strings.Contains(doc.Find("button.room-edit").AttrOr("data-on:click", ""), "/admin/rooms/api/42") {
				t.Fatal("editing must still use the database room ID")
			}
			if _, exists := doc.Find(".room").Attr("data-on:click"); exists {
				t.Fatal("opening a map must not also open the edit dialog")
			}
		})
	}
}

func roomPageFloorIDs(floorGroups []FloorGroup) []int {
	floors := make([]int, 0, len(floorGroups))
	for _, floorGroup := range floorGroups {
		floors = append(floors, floorGroup.Floor)
	}
	return floors
}

func TestRoomsAssignmentMap_UsesDatabaseIDsAndCurrentPuljeAssignments(t *testing.T) {
	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "assignment_map")
	seedRoomsPageLookups(t, db)
	room := createRoomsPageRoom(t, db, "Amalie", "705", 7)
	createRoomsPageRoom(t, db, "Annet rom", "201", 2)
	insertRoomsPagePulje(t, db, models.PuljeFredagKveld)
	insertRoomsPagePulje(t, db, models.PuljeLordagKveld)
	insertRoomsPageEvent(t, db, "friday-event", "Fredag", 4)
	insertRoomsPageEvent(t, db, "saturday-event", "Lørdag", 4)
	insertRoomsPageEventPulje(t, db, "friday-event", models.PuljeFredagKveld, room.ID)
	insertRoomsPageEventPulje(t, db, "saturday-event", models.PuljeLordagKveld, room.ID)

	// When
	doc := templtest.Render(t, RoomsAssignmentPageContent(db, logger, models.PuljeFredagKveld, nil))
	var mappedRooms []models.RoomByPulje
	err := json.Unmarshal([]byte(doc.Find("room-map").AttrOr("rooms", "")), &mappedRooms)

	// Then
	if err != nil {
		t.Fatal(err)
	}
	if len(mappedRooms) != 1 || mappedRooms[0].ID != int64(room.ID) || mappedRooms[0].RoomNumber != "705" {
		t.Fatalf("expected the seventh-floor room with its database ID, got %+v", mappedRooms)
	}
	if len(mappedRooms[0].AssignedEventsID) != 1 || mappedRooms[0].AssignedEventsID[0].EventID != "friday-event" {
		t.Fatalf("expected only Friday assignments, got %+v", mappedRooms[0].AssignedEventsID)
	}
	if doc.Find(`#assignment-dialog button[data-event-id="friday-event"]`).Length() != 1 || doc.Find(`#assignment-dialog button[data-event-id="saturday-event"]`).Length() != 0 {
		t.Fatal("map event picker must contain only the current pulje's events")
	}
	if doc.Find(".rooms-container .room").Length() != 1 || doc.Find("room-map .room").Length() != 1 {
		t.Fatal("rooms outside the map must remain in the assignment list")
	}
	if doc.Find("#map-event").Length() != 0 {
		t.Fatal("map must use per-room add buttons, not the old dropdown")
	}
	if doc.Find(`room-map .room-add-event[aria-label="Legg til arrangement i rom 705"]`).Length() != 1 {
		t.Fatal("mapped room must have an accessible add button")
	}
	if got := doc.Find("room-map .room-event-count").Text(); got != "1" {
		t.Fatalf("expected one assigned event in the room counter, got %q", got)
	}
}

func TestRoomEventCard_IsSharedAcrossAssignedAndUnassignedLocations(t *testing.T) {
	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "shared_room_cards")
	seedRoomsPageLookups(t, db)
	mapped := createRoomsPageRoom(t, db, "Amalie", "705", 7)
	unmapped := createRoomsPageRoom(t, db, "Utenfor kartet", "716", 7)
	otherFloor := createRoomsPageRoom(t, db, "Annen etasje", "201", 2)
	insertRoomsPagePulje(t, db, models.PuljeFredagKveld)
	for _, roomID := range []int{0, mapped.ID, unmapped.ID, otherFloor.ID} {
		id := fmt.Sprintf("event-%d", roomID)
		insertRoomsPageEvent(t, db, id, "Et arrangement", 4)
		insertRoomsPageEventPulje(t, db, id, models.PuljeFredagKveld, roomID)
	}

	// When
	doc := templtest.Render(t, RoomsAssignmentPageContent(db, logger, models.PuljeFredagKveld, nil))

	// Then
	for _, roomID := range []int{0, mapped.ID, unmapped.ID, otherFloor.ID} {
		card := doc.Find(fmt.Sprintf(`.room-event[data-event-id="event-%d"]`, roomID))
		if card.Length() != 1 || card.Find("img").AttrOr("src", "") != "/static/placeholder_banner.svg" || !strings.Contains(card.Text(), "Et arrangement") {
			t.Fatalf("event in room %d should render once with an image and title", roomID)
		}
		remove := card.Find(".room-event-remove")
		if roomID == 0 {
			if remove.Length() != 0 {
				t.Fatal("unassigned event must not offer removal")
			}
		} else if got := remove.AttrOr("data-on:click", ""); got != fmt.Sprintf("@delete('/admin/rooms/api/assignment/FredagKveld/event-%d/%d')", roomID, roomID) {
			t.Fatalf("removal must target the current pulje and room, got %q", got)
		}
	}
}

func TestRoomAssignmentPicker_ExistsWhenPuljeHasNoEvents(t *testing.T) {
	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "empty_room_picker")
	createRoomsPageRoom(t, db, "Amalie", "705", 7)

	// When
	doc := templtest.Render(t, RoomsAssignmentPageContent(db, logger, models.PuljeFredagKveld, nil))

	// Then
	if doc.Find("#assignment-dialog").Length() != 1 || !strings.Contains(doc.Find("#assignment-dialog").Text(), "Ingen arrangementer i denne puljen.") {
		t.Fatal("add button must open a dialog with an empty state, even without events")
	}
}

func createRoomsPageRoom(t *testing.T, db *sql.DB, name string, roomNumber string, floor int) models.Room {
	t.Helper()

	room, err := roomService.CreateRoom(db, models.Room{
		Name:       name,
		RoomNumber: roomNumber,
		Floor:      floor,
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
