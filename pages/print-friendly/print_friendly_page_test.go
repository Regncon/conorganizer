package printfriendly

import (
	"database/sql"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func createPrintPageTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db := testutil.CreateTestDB(t, "print_friendly_page")
	for _, status := range []models.EventStatus{
		models.EventStatusDraft,
		models.EventStatusAnnounced,
	} {
		testutil.MustExec(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}
	for _, status := range []models.PuljeStatus{
		models.PuljeStatusOpen,
		models.PuljeStatusLocked,
		models.PuljeStatusCompleted,
	} {
		testutil.MustExec(t, db, `INSERT INTO pulje_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}
	return db
}

// insertPrintPulje seeds a pulje with its room assignment already published,
// since most tests using this helper assert on visible room information.
// Tests covering the unpublished state override it with an explicit UPDATE.
func insertPrintPulje(t *testing.T, db *sql.DB, id models.Pulje, name, start, end string) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT INTO puljer(id, name, status, start_at, end_at, rooms_published) VALUES (?, ?, ?, ?, ?, 1)`, id, name, models.PuljeStatusOpen, start, end)
}

func insertPrintEvent(t *testing.T, db *sql.DB, id, title string, inPuljefordeling bool) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO events(id, title, intro, description, host_name, email, phone_number, max_players, status, is_in_puljefordeling)
		VALUES (?, ?, 'Intro', 'Description', 'Host', 'host@example.com', '12345678', 4, ?, ?)
	`, id, title, models.EventStatusAnnounced, inPuljefordeling)
}

func assignPrintEvent(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, roomID int) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published, room_id)
		VALUES (?, ?, 1, 1, ?)
	`, eventID, puljeID, roomID)
}

func insertPrintRoom(t *testing.T, db *sql.DB, id int, name, number string) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT INTO rooms(id, name, room_number, floor, max_concurrent_games) VALUES (?, ?, ?, 7, 1)`, id, name, number)
}

func TestPrintFriendlyPage_OrdersDaysAndDeduplicatesProgramEventsWithinEachDay(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A program event runs in two Saturday puljer and again on Sunday, alongside raffle events.",
		When:  "The printable program is rendered.",
		Then:  "Each day shows its program event once before raffle events in pulje order, with maps for that day's rooms.",
	})

	// Given
	expectedDayHeadings := []string{"Lørdag 10.10", "Søndag 11.10"}
	expectedTitles := []string{"Shared Program", "Morning Raffle", "Evening Raffle", "Shared Program", "Sunday Raffle"}
	expectedSaturdayMaps := []string{"/static/rooms/terminus-7-etasje-705.svg", "/static/rooms/terminus-7-etasje-710.svg"}
	expectedSundayMaps := []string{"/static/rooms/terminus-7-etasje-705.svg"}
	db := createPrintPageTestDB(t)
	insertPrintPulje(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	insertPrintPulje(t, db, models.PuljeLordagKveld, "Lørdag kveld", "2026-10-10T18:00:00Z", "2026-10-10T23:00:00Z")
	insertPrintPulje(t, db, models.PuljeSondagMorgen, "Søndag morgen", "2026-10-11T10:00:00Z", "2026-10-11T15:00:00Z")
	insertPrintRoom(t, db, 705, "Morning room", "705")
	insertPrintRoom(t, db, 710, "Evening room", "710")
	insertPrintEvent(t, db, "shared-program", "Shared Program", false)
	insertPrintEvent(t, db, "morning-raffle", "Morning Raffle", true)
	insertPrintEvent(t, db, "evening-raffle", "Evening Raffle", true)
	insertPrintEvent(t, db, "sunday-raffle", "Sunday Raffle", true)
	assignPrintEvent(t, db, "shared-program", models.PuljeLordagMorgen, 705)
	assignPrintEvent(t, db, "shared-program", models.PuljeLordagKveld, 710)
	assignPrintEvent(t, db, "shared-program", models.PuljeSondagMorgen, 705)
	assignPrintEvent(t, db, "morning-raffle", models.PuljeLordagMorgen, 705)
	assignPrintEvent(t, db, "evening-raffle", models.PuljeLordagKveld, 710)
	assignPrintEvent(t, db, "sunday-raffle", models.PuljeSondagMorgen, 705)

	// When
	doc := templtest.Render(t, printFriendlyPage(db, nil, testutil.NewTestLogger()))
	actualDayHeadings := make([]string, 0)
	doc.Find("section.print-day").Each(func(_ int, section *goquery.Selection) {
		actualDayHeadings = append(actualDayHeadings, strings.Join(strings.Fields(section.Find(".print-sheet-day").First().Text()), " "))
	})
	actualTitles := templtest.CollectTexts(doc, "article.print-event-sheet .event-header .title")
	saturdayMaps := printSheetMapSources(doc.Find("section.print-day").Eq(0).Find("article.print-event-sheet").First())
	sundayMaps := printSheetMapSources(doc.Find("section.print-day").Eq(1).Find("article.print-event-sheet").First())

	// Then
	if !slices.Equal(actualDayHeadings, expectedDayHeadings) {
		t.Fatalf("day headings = %v, want %v", actualDayHeadings, expectedDayHeadings)
	}
	if !slices.Equal(actualTitles, expectedTitles) {
		t.Fatalf("printed event order = %v, want %v", actualTitles, expectedTitles)
	}
	if !slices.Equal(saturdayMaps, expectedSaturdayMaps) {
		t.Fatalf("Saturday program maps = %v, want %v", saturdayMaps, expectedSaturdayMaps)
	}
	if !slices.Equal(sundayMaps, expectedSundayMaps) {
		t.Fatalf("Sunday program maps = %v, want %v", sundayMaps, expectedSundayMaps)
	}
}

func TestPrintFriendlyPage_ShowsOneMapForAProgramEventUsingTheSameRoomTwice(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A program event runs in two puljer on the same day in one room.",
		When:  "The printable program is rendered.",
		Then:  "Its single sheet points to the description for the time and prints the room map once.",
	})

	// Given
	expectedMap := "/static/rooms/terminus-7-etasje-710.svg"
	db := createPrintPageTestDB(t)
	insertPrintPulje(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	insertPrintPulje(t, db, models.PuljeLordagKveld, "Lørdag kveld", "2026-10-10T18:00:00Z", "2026-10-10T23:00:00Z")
	insertPrintRoom(t, db, 710, "Shared room", "710")
	insertPrintEvent(t, db, "shared-program", "Shared Program", false)
	assignPrintEvent(t, db, "shared-program", models.PuljeLordagMorgen, 710)
	assignPrintEvent(t, db, "shared-program", models.PuljeLordagKveld, 710)

	// When
	doc := templtest.Render(t, printFriendlyPage(db, nil, testutil.NewTestLogger()))
	sheet := doc.Find("article.print-event-sheet")
	actualMaps := printSheetMapSources(sheet)
	actualTimes := strings.Join(strings.Fields(sheet.Find(".print-room").Text()), " ")

	// Then
	if sheet.Length() != 1 {
		t.Fatalf("printed sheets = %d, want 1", sheet.Length())
	}
	if !slices.Equal(actualMaps, []string{expectedMap}) {
		t.Fatalf("room maps = %v, want one %q", actualMaps, expectedMap)
	}
	if !strings.Contains(actualTimes, "Se beskrivelse for tidspunkt") {
		t.Fatalf("printed times = %q, want it to point to the description", actualTimes)
	}
	if strings.Contains(actualTimes, "10:00 - 15:00") || strings.Contains(actualTimes, "18:00 - 23:00") {
		t.Fatalf("printed times = %q, want no pulje times for a programarrangement", actualTimes)
	}
}

func TestPrintFriendlyPage_ShowsPuljeTimeForPuljeEvent(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An event in puljefordelingen runs in one pulje.",
		When:  "The printable program is rendered.",
		Then:  "Its sheet shows the pulje name and time.",
	})

	// Given
	expectedTime := "Lørdag morgen · 10:00 - 15:00"
	db := createPrintPageTestDB(t)
	insertPrintPulje(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	insertPrintRoom(t, db, 705, "Morning room", "705")
	insertPrintEvent(t, db, "morning-raffle", "Morning Raffle", true)
	assignPrintEvent(t, db, "morning-raffle", models.PuljeLordagMorgen, 705)

	// When
	doc := templtest.Render(t, printFriendlyPage(db, nil, testutil.NewTestLogger()))
	actualTimes := strings.Join(strings.Fields(doc.Find("article.print-event-sheet .print-room").Text()), " ")

	// Then
	if !strings.Contains(actualTimes, expectedTime) {
		t.Fatalf("printed times = %q, want %q", actualTimes, expectedTime)
	}
	if strings.Contains(actualTimes, "Se beskrivelse for tidspunkt") {
		t.Fatalf("printed times = %q, pulje event should not point to the description", actualTimes)
	}
}

func TestPrintFriendlyPage_ShowsNoRoomForUnpublishedPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A print event assigned to a room in a pulje whose rooms are not published.",
		When:  "Anyone, admins included, renders the printable program.",
		Then:  "The sheet shows no room information for that occurrence.",
	})

	// Given
	db := createPrintPageTestDB(t)
	insertPrintPulje(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	testutil.MustExec(t, db, `UPDATE puljer SET rooms_published = 0 WHERE id = ?`, models.PuljeLordagMorgen)
	insertPrintRoom(t, db, 705, "Morning room", "705")
	insertPrintEvent(t, db, "shared-program", "Shared Program", false)
	assignPrintEvent(t, db, "shared-program", models.PuljeLordagMorgen, 705)

	// When
	doc := templtest.Render(t, printFriendlyPage(db, nil, testutil.NewTestLogger()))
	sheet := doc.Find("article.print-event-sheet")

	// Then
	if strings.Contains(sheet.Text(), "Morning room") {
		t.Fatal("room name should be hidden for an unpublished pulje")
	}
	if !strings.Contains(sheet.Find(".print-room").Text(), "Rom ikke tildelt") {
		t.Fatal("expected the occurrence to render as having no room assigned")
	}
	if sheet.Find(".print-room-map").Length() != 0 {
		t.Fatal("map should be hidden for an unpublished pulje")
	}
}

func printSheetMapSources(sheet *goquery.Selection) []string {
	paths := make([]string, 0)
	sheet.Find(".print-room-map img").Each(func(_ int, image *goquery.Selection) {
		paths = append(paths, image.AttrOr("src", ""))
	})
	return paths
}
