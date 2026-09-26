package event

import (
	"database/sql"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func createEventRoomTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db := createEventVisibilityTestDB(t)
	seedEventVisibilityEvent(t, db, "room-event", "Room event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "room-event", models.PuljeFredagKveld, true)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?), (?)`, models.PuljeStatusCompleted, models.PuljeStatusLocked)
	testutil.MustExec(t, db, `INSERT INTO rooms(id, name, room_number, floor, max_concurrent_games) VALUES (42, 'Amalie Hansen', '705', 7, 1)`)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = 42 WHERE event_id = 'room-event'`)
	return db
}

func TestEventRoomMap_InfoDeskShowsGroundFloorRoutes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An event is assigned to Info Desk, room 004 on the ground floor.",
		When:  "The event page displays the room map.",
		Then:  "The map points to room 004 and describes routes from both the lifts and stairs.",
	})

	// Given
	expectedMapPath := "/static/rooms/terminus-0-etasje-004.svg"
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE rooms SET name = 'Info Desk', room_number = '004', floor = 0 WHERE id = 42`)

	// When
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
	mapImage := doc.Find(".event-room-dialog img")

	// Then
	if got := mapImage.AttrOr("src", ""); got != expectedMapPath {
		t.Fatalf("map = %q, want %q", got, expectedMapPath)
	}
	if got := mapImage.AttrOr("alt", ""); !strings.Contains(got, "heisene og trappen") {
		t.Fatalf("map alt text does not describe both starting points: %q", got)
	}
}

func TestEventRoomVisibility(t *testing.T) {
	tests := []struct {
		name     string
		update   string
		wantRoom bool
		wantMap  bool
	}{
		{name: "open allocation", wantRoom: true, wantMap: true},
		{name: "unpublished program", update: `UPDATE program_publishing_state SET is_published = 0`},
		{name: "completed allocation", update: `UPDATE puljer SET status = 'Completed'`, wantRoom: true, wantMap: true},
		{name: "locked allocation", update: `UPDATE puljer SET status = 'Locked'`, wantRoom: true, wantMap: true},
		{name: "legacy unpublished occurrence flag is ignored", update: `UPDATE relation_event_puljer SET is_published = 0`, wantRoom: true, wantMap: true},
		{name: "removed occurrence", update: `UPDATE relation_event_puljer SET is_in_pulje = 0`},
		{name: "unassigned room", update: `UPDATE relation_event_puljer SET room_id = NULL`},
		{name: "deleted room", update: `DELETE FROM rooms WHERE id = 42`},
		{name: "unknown map", update: `UPDATE rooms SET room_number = '999'`, wantRoom: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := createEventRoomTestDB(t)
			if test.update != "" {
				testutil.MustExec(t, db, test.update)
			}
			request := httptest.NewRequest("GET", "/event/room-event", nil)
			doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
			if got := doc.Find(".event-room-list .event-room-name").Length() > 0; got != test.wantRoom {
				t.Fatalf("room visible = %v, want %v", got, test.wantRoom)
			}
			if got := doc.Find(".event-room-button").Length() > 0; got != test.wantMap {
				t.Fatalf("map button visible = %v, want %v", got, test.wantMap)
			}
			if got := doc.Find(".event-room-dialog img").Length() > 0; got != test.wantMap {
				t.Fatalf("map rendered = %v, want %v", got, test.wantMap)
			}
			if !test.wantRoom && strings.Contains(doc.Find("#event-room-locations").Text(), "Amalie Hansen") {
				t.Fatal("hidden room name leaked into markup")
			}
		})
	}
}

func TestEventRoomsUseEachPuljeAssignment(t *testing.T) {
	db := createEventRoomTestDB(t)
	seedEventVisibilityPulje(t, db, models.PuljeLordagKveld)
	seedEventVisibilityEventPulje(t, db, "room-event", models.PuljeLordagKveld, false)
	testutil.MustExec(t, db, `UPDATE puljer SET name = 'Lørdag kveld', status = ?, start_at = '2026-10-10T18:30:00+02:00' WHERE id = ?`, models.PuljeStatusCompleted, models.PuljeLordagKveld)
	testutil.MustExec(t, db, `INSERT INTO rooms(id, name, room_number, floor, max_concurrent_games) VALUES (43, 'Lørdagsrommet', '710', 7, 1)`)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = 43 WHERE pulje_id = ?`, models.PuljeLordagKveld)
	seedEventVisibilityEvent(t, db, "other-event", "Other event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityEventPulje(t, db, "other-event", models.PuljeFredagKveld, true)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = 43 WHERE event_id = 'other-event'`)

	assignments, err := getEventRooms(db, "room-event", true)
	if err != nil {
		t.Fatal(err)
	}
	if len(assignments) != 2 {
		t.Fatalf("got %d assignments, want 2", len(assignments))
	}
	for i, want := range []struct {
		pulje  models.Pulje
		roomID int
		number string
	}{
		{models.PuljeFredagKveld, 42, "705"},
		{models.PuljeLordagKveld, 43, "710"},
	} {
		got := assignments[i]
		if got.Pulje.ID != want.pulje || got.Room.ID != want.roomID || !strings.HasSuffix(got.MapPath, "-"+want.number+".svg") {
			t.Fatalf("assignment %d: %+v", i, got)
		}
	}
	request := httptest.NewRequest("GET", "/event/room-event?pulje=LordagKveld", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
	for _, assignment := range assignments {
		id := eventRoomDialogID("room-event", assignment.Pulje.ID)
		if got := doc.Find("#"+id+" img").AttrOr("src", ""); got != assignment.MapPath {
			t.Fatalf("map = %q, want %q", got, assignment.MapPath)
		}
		if !strings.Contains(doc.Find("#"+id+"-button").Text(), assignment.Pulje.Name) {
			t.Fatal("button is missing its pulje name")
		}
	}
	// A reassignment updates the map without changing the dialog's identity.
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = 43 WHERE event_id = 'room-event' AND pulje_id = ?`, models.PuljeFredagKveld)
	doc = templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
	if got := doc.Find("#"+eventRoomDialogID("room-event", models.PuljeFredagKveld)+" img").AttrOr("src", ""); got != assignments[1].MapPath {
		t.Fatalf("reassigned map = %q", got)
	}
}

func TestEventRoomButtonsGroupSchedulesByRoom(t *testing.T) {
	tests := []struct {
		name           string
		secondRoomID   any
		roomNumber     string
		roomName       string
		wantButtons    int
		wantUnassigned int
	}{
		{name: "same room", secondRoomID: 42, wantButtons: 1},
		{name: "different rooms", secondRoomID: 43, roomNumber: "710", roomName: "Lucie Wolf", wantButtons: 2},
		{name: "same name but different rooms", secondRoomID: 43, roomNumber: "710", roomName: "Amalie Hansen", wantButtons: 2},
		{name: "unassigned occurrence", secondRoomID: nil, wantButtons: 1, wantUnassigned: 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := createEventRoomTestDB(t)
			seedEventVisibilityPulje(t, db, models.PuljeLordagMorgen)
			seedEventVisibilityEventPulje(t, db, "room-event", models.PuljeLordagMorgen, true)
			testutil.MustExec(t, db, `UPDATE puljer SET name = 'Lørdag morgen', start_at = '2026-10-10T10:00:00+02:00', end_at = '2026-10-10T15:00:00+02:00' WHERE id = ?`, models.PuljeLordagMorgen)
			if test.roomNumber != "" {
				testutil.MustExec(t, db, `INSERT INTO rooms(id, name, room_number, floor, max_concurrent_games) VALUES (43, ?, ?, 7, 1)`, test.roomName, test.roomNumber)
			}
			testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = ? WHERE pulje_id = ?`, test.secondRoomID, models.PuljeLordagMorgen)
			request := httptest.NewRequest("GET", "/event/room-event", nil)
			doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
			if got := doc.Find(".event-room-button").Length(); got != test.wantButtons {
				t.Fatalf("got %d buttons, want %d", got, test.wantButtons)
			}
			if got := doc.Find(".event-schedule-unassigned").Length(); got != test.wantUnassigned {
				t.Fatalf("got %d unassigned schedule rows, want %d", got, test.wantUnassigned)
			}
			list := doc.Find(".event-room-list")
			for _, schedule := range []string{"Fredag kveld · 18:30 - 23:00", "Lørdag morgen · 10:00 - 15:00"} {
				if strings.Count(list.Text(), schedule) != 1 {
					t.Fatalf("schedule should appear exactly once in the list: %s", schedule)
				}
			}
			if strings.Contains(doc.Text(), "Pulje(r)") {
				t.Fatal("redundant pulje section is still rendered")
			}
			if test.name == "same room" {
				button := doc.Find(".event-room-button")
				if button.Find(".event-room-schedule > span").Length() != 2 {
					t.Fatal("shared room button must include both times")
				}
				if strings.Index(button.Text(), "Fredag kveld") > strings.Index(button.Text(), "Lørdag morgen") {
					t.Fatal("puljer are not chronological")
				}
				id := button.AttrOr("aria-controls", "")
				if doc.Find("#"+id+" .event-room-schedule > span").Length() != 2 {
					t.Fatal("shared room modal must include both times")
				}
			}
		})
	}
}

func TestEventScheduleIsHiddenBeforeProgramPublication(t *testing.T) {
	db := createEventRoomTestDB(t)
	setEventVisibilityProgramPublishing(t, db, false)
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
	for _, text := range []string{"Pulje(r)", "Fredag kveld", "18:30", "Amalie Hansen"} {
		if strings.Contains(doc.Text(), text) {
			t.Fatalf("unpublished schedule leaked into page: %s", text)
		}
	}
	if doc.Find(".event-room-schedule, .event-room-button, .event-room-dialog").Length() != 0 {
		t.Fatal("unpublished schedule or map is rendered")
	}
}

func TestUnassignedEventKeepsItsPublishedSchedule(t *testing.T) {
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = NULL`)
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
	if got := doc.Find(".event-schedule-unassigned").Text(); !strings.Contains(got, "Fredag kveld · 18:30 - 23:00") {
		t.Fatalf("missing unassigned event schedule: %q", got)
	}
	if doc.Find(".event-room-button").Length() != 0 {
		t.Fatal("unassigned event has a map button")
	}
}
