package event

import (
	"database/sql"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

const eventRoomPendingText = "Stedet for arrangementet kommer!"

func createEventRoomTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db := createEventVisibilityTestDB(t)
	seedEventVisibilityEvent(t, db, "room-event", "Room event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityPulje(t, db, models.PuljeFredagKveld)
	seedEventVisibilityEventPulje(t, db, "room-event", models.PuljeFredagKveld, true)
	testutil.MustExec(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?), (?)`, models.PuljeStatusCompleted, models.PuljeStatusLocked)
	testutil.MustExec(t, db, `INSERT INTO rooms(id, name, room_number, floor, max_concurrent_games, public_notes, admin_notes) VALUES (42, 'Amalie Hansen', '705', 7, 1, '', '')`)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = 42 WHERE event_id = 'room-event'`)
	return db
}

func TestEventRoomList_ShowsPublicNoteInListAndDialog(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A published room assignment whose room has a public note.",
		When:  "The event page displays the room.",
		Then:  "The public note is shown both in the room list and in the room's map dialog.",
	})

	// Given
	expectedNote := "Ta av deg skoene før du går inn."
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE rooms SET public_notes = ?, admin_notes = 'Kun for admin' WHERE id = 42`, expectedNote)

	// When
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))

	// Then
	if got := doc.Find(".event-room-list .event-room-note").Text(); !strings.Contains(got, expectedNote) {
		t.Fatalf("list note = %q, want it to contain %q", got, expectedNote)
	}
	if got := doc.Find(".event-room-dialog .event-room-note").Text(); !strings.Contains(got, expectedNote) {
		t.Fatalf("dialog note = %q, want it to contain %q", got, expectedNote)
	}
}

func TestEventRoomList_NeverRendersAdminNotes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A published room assignment whose room has an admin-only note.",
		When:  "The event page displays the room.",
		Then:  "The admin note text never appears anywhere on the page.",
	})

	// Given
	adminOnlyNote := "Hemmelig admin-notat"
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE rooms SET public_notes = '', admin_notes = ? WHERE id = 42`, adminOnlyNote)

	// When
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))

	// Then
	if strings.Contains(doc.Text(), adminOnlyNote) {
		t.Fatal("admin note leaked onto the event page")
	}
}

func TestEventRoomList_RendersPublicNoteMarkupAsText(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A published room assignment whose public note contains HTML markup.",
		When:  "The event page displays the room.",
		Then:  "The markup is shown as plain text and never becomes an element.",
	})

	// Given
	expectedNote := `<img src=x onerror="alert(1)">`
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE rooms SET public_notes = ? WHERE id = 42`, expectedNote)

	// When
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))

	// Then
	if got := doc.Find(".event-room-list .event-room-note").Text(); !strings.Contains(got, expectedNote) {
		t.Fatalf("list note = %q, want it to contain %q as text", got, expectedNote)
	}
	if doc.Find(".event-room-note img").Length() != 0 {
		t.Fatal("public note markup was rendered as HTML")
	}
}

func TestEventRoomList_NoNoteElementWhenRoomHasNoPublicNotes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A published room assignment whose room has no public notes.",
		When:  "The event page displays the room.",
		Then:  "No note element is rendered for that room.",
	})

	// Given
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE rooms SET public_notes = '' WHERE id = 42`)

	// When
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))

	// Then
	if doc.Find(".event-room-note").Length() != 0 {
		t.Fatal("note element rendered for a room without public notes")
	}
}

func TestEventRoomVisibility_NonAdminSeesTimeOnlyWhenPuljeRoomsUnpublished(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A room assignment for a pulje whose rooms are not published.",
		When:  "A non-admin views the event page.",
		Then:  "The schedule shows the time but hides the room name and map.",
	})

	// Given
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE puljer SET rooms_published = 0 WHERE id = ?`, models.PuljeFredagKveld)

	// When
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))

	// Then
	if doc.Find(".event-room-name").Length() != 0 {
		t.Fatal("room name should be hidden from non-admins for an unpublished pulje")
	}
	if doc.Find(".event-room-button, .event-room-dialog img").Length() != 0 {
		t.Fatal("map should be hidden from non-admins for an unpublished pulje")
	}
	if strings.Contains(doc.Text(), "Amalie Hansen") {
		t.Fatal("hidden room name leaked into markup")
	}
	if got := doc.Find(".event-room-list .event-room-schedule").Text(); !strings.Contains(got, "Fredag kveld · 18:30 - 23:00") {
		t.Fatalf("expected the time to remain visible, got %q", got)
	}
	if got := doc.Find(".event-room-pending").Text(); got != eventRoomPendingText {
		t.Fatalf("expected the pending room text, got %q", got)
	}
}

func TestEventRoomVisibility_NonAdminSeesRoomWhenPuljeRoomsPublished(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A room assignment for a pulje whose rooms are published.",
		When:  "A non-admin views the event page.",
		Then:  "The schedule shows the room name and map.",
	})

	// Given
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE puljer SET rooms_published = 1 WHERE id = ?`, models.PuljeFredagKveld)

	// When
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))

	// Then
	if got := doc.Find(".event-room-name").Text(); got != "Amalie Hansen" {
		t.Fatalf("expected the room name to be visible, got %q", got)
	}
	if doc.Find(".event-room-button").Length() != 1 {
		t.Fatal("expected the map button to be visible")
	}
	if doc.Find(".event-room-pending").Length() != 0 {
		t.Fatal("pending room text shown although the room is published")
	}
}

func TestEventRoomVisibility_AdminSeesTimeOnlyWhenPuljeRoomsUnpublished(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A room assignment for a pulje whose rooms are not published.",
		When:  "An admin views the event page.",
		Then:  "The admin, like everyone else, sees the time but not the room name or map.",
	})

	// Given
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE puljer SET rooms_published = 0 WHERE id = ?`, models.PuljeFredagKveld)

	// When
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", true, testutil.NewTestLogger(), db, nil, request))

	// Then
	if strings.Contains(doc.Text(), "Amalie Hansen") {
		t.Fatal("unpublished room name should be hidden from admins too")
	}
	if doc.Find(".event-room-button, .event-room-dialog img").Length() != 0 {
		t.Fatal("unpublished room map should be hidden from admins too")
	}
	if got := doc.Find(".event-room-list .event-room-schedule").Text(); !strings.Contains(got, "Fredag kveld · 18:30 - 23:00") {
		t.Fatalf("expected the time to remain visible, got %q", got)
	}
	if got := doc.Find(".event-room-pending").Text(); got != eventRoomPendingText {
		t.Fatalf("expected the pending room text, got %q", got)
	}
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
	if got := mapImage.AttrOr("alt", ""); strings.Contains(got, "004") {
		t.Fatalf("ground-floor map alt text should use the room name without a number: %q", got)
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
	bdd.Behavior(t, bdd.BDD{
		Given: "An event whose puljer are in different rooms, each with a map.",
		When:  "The event page displays the rooms.",
		Then:  "Each pulje's map button sits in its own room's row, and its map dialog shows that pulje's room.",
	})

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
		if got := doc.Find("#" + id + "-button").Closest(".event-room-row").Find(".event-room-name").Text(); got != assignment.Room.Name {
			t.Fatalf("map button for %s is in the row for %q, want %q", assignment.Pulje.Name, got, assignment.Room.Name)
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
	bdd.Behavior(t, bdd.BDD{
		Given: "An event with two puljer in the same room, different rooms, or one pulje without a room.",
		When:  "The event page displays the schedule.",
		Then:  "Each time is listed once, each room gets one row, and the pending text stays hidden because the event has a room.",
	})

	tests := []struct {
		name         string
		secondRoomID any
		roomNumber   string
		roomName     string
		wantButtons  int
		wantRooms    []string
	}{
		{name: "same room", secondRoomID: 42, wantButtons: 1, wantRooms: []string{"Amalie Hansen"}},
		{name: "different rooms", secondRoomID: 43, roomNumber: "710", roomName: "Lucie Wolf", wantButtons: 2, wantRooms: []string{"Amalie Hansen", "Lucie Wolf"}},
		{name: "same name but different rooms", secondRoomID: 43, roomNumber: "710", roomName: "Amalie Hansen", wantButtons: 2, wantRooms: []string{"Amalie Hansen", "Amalie Hansen"}},
		{name: "unassigned occurrence", secondRoomID: nil, wantButtons: 1, wantRooms: []string{"Amalie Hansen"}},
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
			if got := templtest.CollectTexts(doc, ".event-room-row .event-room-name"); !slices.Equal(got, test.wantRooms) {
				t.Fatalf("rooms = %q, want %q", got, test.wantRooms)
			}
			if doc.Find(".event-room-pending").Length() != 0 {
				t.Fatal("pending room text shown although the event has a room")
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
				box := doc.Find(".event-room-list")
				if box.Find(".event-room-schedule > span").Length() != 2 {
					t.Fatal("shared room box must include both times")
				}
				if strings.Index(box.Text(), "Fredag kveld") > strings.Index(box.Text(), "Lørdag morgen") {
					t.Fatal("puljer are not chronological")
				}
				id := box.Find(".event-room-button").AttrOr("aria-controls", "")
				if doc.Find("#"+id+" .event-room-schedule > span").Length() != 2 {
					t.Fatal("shared room modal must include both times")
				}
			}
		})
	}
}

func TestEventScheduleIsHiddenBeforeProgramPublication(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An event with a room assignment while the program is not published.",
		When:  "The event page is displayed.",
		Then:  "No schedule, room, placeholder text or map is rendered.",
	})

	db := createEventRoomTestDB(t)
	setEventVisibilityProgramPublishing(t, db, false)
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
	for _, text := range []string{"Pulje(r)", "Fredag kveld", "18:30", "Amalie Hansen", "Sted og tidspunkt", eventRoomPendingText} {
		if strings.Contains(doc.Text(), text) {
			t.Fatalf("unpublished schedule leaked into page: %s", text)
		}
	}
	if doc.Find(".event-room-schedule, .event-room-button, .event-room-dialog").Length() != 0 {
		t.Fatal("unpublished schedule or map is rendered")
	}
}

func TestUnassignedEventKeepsItsPublishedSchedule(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A published event whose pulje has no room yet.",
		When:  "The event page is displayed.",
		Then:  "The time is shown with the pending room text and no map button.",
	})

	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = NULL`)
	request := httptest.NewRequest("GET", "/event/room-event", nil)
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))
	if got := doc.Find(".event-room-list .event-room-schedule").Text(); !strings.Contains(got, "Fredag kveld · 18:30 - 23:00") {
		t.Fatalf("missing unassigned event schedule: %q", got)
	}
	if got := doc.Find(".event-room-pending").Text(); got != eventRoomPendingText {
		t.Fatalf("expected the pending room text, got %q", got)
	}
	if doc.Find(".event-room-button").Length() != 0 {
		t.Fatal("unassigned event has a map button")
	}
}

func TestEventSchedule_ProgramEventShowsItsTimePerDayInsteadOfPuljeTimes(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et programarrangement fredag kveld, og hele lørdagen i samme rom med tidspunktet «Hele dagen» i begge lørdagspuljene.",
		When:  "Når arrangementssiden vises.",
		Then:  "Så viser tidsplanen tidspunktet per dag, én gang for lørdag, og ikke puljetidene.",
	})

	// Given
	expectedSchedule := []string{"Fredag 9.10 · 18:00–20:00", "Lørdag 10.10 · Hele dagen"}
	db := createEventRoomTestDB(t)
	testutil.MustExec(t, db, `UPDATE events SET is_in_puljefordeling = 0 WHERE id = 'room-event'`)
	for _, pulje := range []struct {
		id      models.Pulje
		name    string
		startAt string
	}{
		{models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00+02:00"},
		{models.PuljeLordagKveld, "Lørdag kveld", "2026-10-10T18:00:00+02:00"},
	} {
		seedEventVisibilityPulje(t, db, pulje.id)
		testutil.MustExec(t, db, `UPDATE puljer SET name = ?, start_at = ? WHERE id = ?`, pulje.name, pulje.startAt, pulje.id)
		seedEventVisibilityEventPulje(t, db, "room-event", pulje.id, true)
	}
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET room_id = 42, program_time = 'Hele dagen' WHERE event_id = 'room-event'`)
	testutil.MustExec(t, db, `UPDATE relation_event_puljer SET program_time = '18:00–20:00' WHERE event_id = 'room-event' AND pulje_id = ?`, models.PuljeFredagKveld)
	request := httptest.NewRequest("GET", "/event/room-event", nil)

	// When
	doc := templtest.Render(t, event_page_content("room-event", false, testutil.NewTestLogger(), db, nil, request))

	// Then
	schedule := templtest.CollectTexts(doc, ".event-room-list > .event-room-schedule > span")
	if !slices.Equal(schedule, expectedSchedule) {
		t.Fatalf("program event schedule = %q, want %q", schedule, expectedSchedule)
	}
	if strings.Contains(doc.Find(".event-room-list").Text(), "18:30 - 23:00") {
		t.Fatal("program event schedule still shows the pulje time")
	}
}
