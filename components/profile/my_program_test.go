package profilecomponent

import (
	"database/sql"
	"log/slog"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestGetAllEventsForUser_WhenPlayerAssignmentIsInOpenPulje_ReturnsInterestsInstead(t *testing.T) {
	// Given a player assignment in an open pulje,
	// when the profile program data is loaded,
	// then the player result is hidden and the user's interests are returned.

	// Given
	expectedEventTitles := []string{}
	expectedInterestNames := []string{"Open Wish Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusOpen)
	insertProfileProgramPublishedEvent(t, db, "open-assigned-event", "Open Assigned Event")
	insertProfileProgramPublishedEvent(t, db, "open-wish-event", "Open Wish Event")
	insertProfileProgramPlayer(t, db, "open-assigned-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "open-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	events, eventsErr := GetAllEventsForUser(userInfo, billettholderID, db, logger)
	interests, interestsErr := getAllInterestsForUser(userInfo, billettholderID, db, logger)

	// Then
	if eventsErr != nil {
		t.Fatalf("expected event query to succeed: %v", eventsErr)
	}
	if interestsErr != nil {
		t.Fatalf("expected interest query to succeed: %v", interestsErr)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramInterestNames(t, expectedInterestNames, interests)
}

func TestGetAllEventsForUser_WhenPlayerAssignmentIsInLockedPulje_ReturnsInterestsInstead(t *testing.T) {
	// Given a player assignment in a locked pulje,
	// when the profile program data is loaded,
	// then the player result is hidden and the user's interests are returned.

	// Given
	expectedEventTitles := []string{}
	expectedInterestNames := []string{"Locked Wish Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked)
	insertProfileProgramPublishedEvent(t, db, "locked-assigned-event", "Locked Assigned Event")
	insertProfileProgramPublishedEvent(t, db, "locked-wish-event", "Locked Wish Event")
	insertProfileProgramPlayer(t, db, "locked-assigned-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "locked-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	events, eventsErr := GetAllEventsForUser(userInfo, billettholderID, db, logger)
	interests, interestsErr := getAllInterestsForUser(userInfo, billettholderID, db, logger)

	// Then
	if eventsErr != nil {
		t.Fatalf("expected event query to succeed: %v", eventsErr)
	}
	if interestsErr != nil {
		t.Fatalf("expected interest query to succeed: %v", interestsErr)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramInterestNames(t, expectedInterestNames, interests)
}

func TestGetAllEventsForUser_WhenPlayerAssignmentIsInCompletedPulje_ReturnsPlayerResult(t *testing.T) {
	// Given a player assignment in a completed pulje,
	// when the profile program data is loaded,
	// then the assigned event is returned as the user's program.

	// Given
	expectedEventTitles := []string{"Completed Assigned Event"}
	expectedInterestNames := []string{}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "completed-assigned-event", "Completed Assigned Event")
	insertProfileProgramPublishedEvent(t, db, "completed-wish-event", "Completed Wish Event")
	insertProfileProgramPlayer(t, db, "completed-assigned-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "completed-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	events, eventsErr := GetAllEventsForUser(userInfo, billettholderID, db, logger)
	interests, interestsErr := getAllInterestsForUser(userInfo, billettholderID, db, logger)

	// Then
	if eventsErr != nil {
		t.Fatalf("expected event query to succeed: %v", eventsErr)
	}
	if interestsErr != nil {
		t.Fatalf("expected interest query to succeed: %v", interestsErr)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramInterestNames(t, expectedInterestNames, interests)
}

func TestGetAllEventsForUser_WhenGMEventIsInOpenPulje_ReturnsGMEvent(t *testing.T) {
	// Given a GM assignment in an open pulje,
	// when the profile program data is loaded,
	// then the GM event is returned.

	// Given
	expectedEventTitles := []string{"Open GM Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusOpen)
	insertProfileProgramPublishedEvent(t, db, "open-gm-event", "Open GM Event")
	insertProfileProgramPlayer(t, db, "open-gm-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)

	// When
	events, err := GetAllEventsForUser(userInfo, billettholderID, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected event query to succeed: %v", err)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramEventsAreGM(t, events)
}

func TestGetAllEventsForUser_WhenGMEventIsInLockedPulje_ReturnsGMEvent(t *testing.T) {
	// Given a GM assignment in a locked pulje,
	// when the profile program data is loaded,
	// then the GM event is returned.

	// Given
	expectedEventTitles := []string{"Locked GM Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked)
	insertProfileProgramPublishedEvent(t, db, "locked-gm-event", "Locked GM Event")
	insertProfileProgramPlayer(t, db, "locked-gm-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)

	// When
	events, err := GetAllEventsForUser(userInfo, billettholderID, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected event query to succeed: %v", err)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramEventsAreGM(t, events)
}

func TestGetAllEventsForUser_WhenGMEventIsInCompletedPulje_ReturnsGMEvent(t *testing.T) {
	// Given a GM assignment in a completed pulje,
	// when the profile program data is loaded,
	// then the GM event is returned.

	// Given
	expectedEventTitles := []string{"Completed GM Event"}

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "completed-gm-event", "Completed GM Event")
	insertProfileProgramPlayer(t, db, "completed-gm-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)

	// When
	events, err := GetAllEventsForUser(userInfo, billettholderID, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected event query to succeed: %v", err)
	}
	assertProfileProgramEventTitles(t, expectedEventTitles, events)
	assertProfileProgramEventsAreGM(t, events)
}

func TestMyProgram_WhenPuljeIsNotCompleted_RendersInterestsAndHidesPlayerResult(t *testing.T) {
	// Given a player assignment and a wish in an open pulje,
	// when Mitt festivalprogram is rendered,
	// then the visible HTML shows the wish and hides the player allocation.

	// Given
	expectedVisibleText := "Visible Wish Event"
	hiddenVisibleText := "Hidden Player Result"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusOpen)
	insertProfileProgramPublishedEvent(t, db, "hidden-player-result", hiddenVisibleText)
	insertProfileProgramPublishedEvent(t, db, "visible-wish-event", expectedVisibleText)
	insertProfileProgramPlayer(t, db, "hidden-player-result", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "visible-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if !strings.Contains(actualText, expectedVisibleText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedVisibleText, actualText)
	}
	if !strings.Contains(actualText, models.InterestLevelHigh.Label()) {
		t.Fatalf("expected rendered profile program to contain interest level %q\nactual text: %s", models.InterestLevelHigh.Label(), actualText)
	}
	if strings.Contains(actualText, hiddenVisibleText) {
		t.Fatalf("expected rendered profile program to hide %q\nactual text: %s", hiddenVisibleText, actualText)
	}
}

func TestMyProgram_WhenPuljeIsCompleted_RendersPlayerResult(t *testing.T) {
	// Given a player assignment in a completed pulje,
	// when Mitt festivalprogram is rendered,
	// then the visible HTML shows what the user is playing.

	// Given
	expectedVisibleText := "Completed Player Result"
	hiddenVisibleText := "Completed Wish Hidden By Result"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusCompleted)
	insertProfileProgramPublishedEvent(t, db, "completed-player-result", expectedVisibleText)
	insertProfileProgramPublishedEvent(t, db, "completed-wish-event", hiddenVisibleText)
	insertProfileProgramPlayer(t, db, "completed-player-result", models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer)
	insertProfileProgramInterest(t, db, "completed-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if !strings.Contains(actualText, expectedVisibleText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedVisibleText, actualText)
	}
	if strings.Contains(actualText, hiddenVisibleText) {
		t.Fatalf("expected rendered profile program to hide %q\nactual text: %s", hiddenVisibleText, actualText)
	}
}

func TestMyProgram_WhenGMEventIsInNotCompletedPulje_RendersGMEventOverInterests(t *testing.T) {
	// Given a GM assignment and a wish in a locked pulje,
	// when Mitt festivalprogram is rendered,
	// then the visible HTML shows the GM event instead of the interests.

	// Given
	expectedVisibleText := "Locked GM Event"
	hiddenVisibleText := "Locked Wish Hidden By GM"

	db, logger := createProfileProgramTestDB(t)
	userInfo, billettholderID := seedProfileProgramUser(t, db)
	insertProfileProgramPulje(t, db, models.PuljeFredagKveld, models.PuljeStatusLocked)
	insertProfileProgramPublishedEvent(t, db, "locked-gm-event", expectedVisibleText)
	insertProfileProgramPublishedEvent(t, db, "locked-wish-event", hiddenVisibleText)
	insertProfileProgramPlayer(t, db, "locked-gm-event", models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM)
	insertProfileProgramInterest(t, db, "locked-wish-event", models.PuljeFredagKveld, billettholderID, models.InterestLevelHigh)

	// When
	doc := templtest.Render(t, MyProgram(userInfo, billettholderID, db, logger, nil))
	actualText := profileProgramVisibleText(doc)

	// Then
	if !strings.Contains(actualText, expectedVisibleText) {
		t.Fatalf("expected rendered profile program to contain %q\nactual text: %s", expectedVisibleText, actualText)
	}
	if strings.Contains(actualText, hiddenVisibleText) {
		t.Fatalf("expected rendered profile program to hide %q\nactual text: %s", hiddenVisibleText, actualText)
	}
}

func createProfileProgramTestDB(t *testing.T) (*sql.DB, *slog.Logger) {
	t.Helper()

	db, logger := testutil.CreateTestDBAndLogger(t, "profile_program")
	seedProfileProgramLookups(t, db)

	return db, logger
}

func seedProfileProgramLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, status := range []models.EventStatus{
		models.EventStatusAnnounced,
	} {
		mustExecProfileProgramTest(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}

	mustExecProfileProgramTest(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeOther)
	mustExecProfileProgramTest(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupDefault)
	mustExecProfileProgramTest(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeNormal)

	for _, status := range []models.PuljeStatus{
		models.PuljeStatusOpen,
		models.PuljeStatusLocked,
		models.PuljeStatusCompleted,
	} {
		mustExecProfileProgramTest(t, db, `INSERT INTO pulje_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}

	for _, interestLevel := range []models.InterestLevel{
		models.InterestLevelHigh,
		models.InterestLevelMedium,
		models.InterestLevelLow,
	} {
		mustExecProfileProgramTest(t, db, `INSERT INTO interest_levels(interest_level) VALUES (?) ON CONFLICT(interest_level) DO NOTHING`, interestLevel)
	}
}

func seedProfileProgramUser(t *testing.T, db *sql.DB) (requestctx.UserRequestInfo, int) {
	t.Helper()

	const externalID = "profile-program-user"
	const email = "profile-program-user@example.com"
	const billettholderID = 1001

	mustExecProfileProgramTest(t, db, `
		INSERT INTO users(id, external_id, email)
		VALUES(?, ?, ?)
	`, 501, externalID, email)
	mustExecProfileProgramTest(t, db, `
		INSERT INTO billettholdere(
			id,
			first_name,
			last_name,
			ticket_type_id,
			ticket_type,
			order_id,
			ticket_id
		)
		VALUES(?, ?, ?, ?, ?, ?, ?)
	`, billettholderID, "Profile", "Program", 1, "Festivalpass", 100, 200)
	mustExecProfileProgramTest(t, db, `
		INSERT INTO relation_billettholdere_users(billettholder_id, user_id)
		VALUES(?, ?)
	`, billettholderID, 501)

	return requestctx.UserRequestInfo{
		IsLoggedIn: true,
		Id:         externalID,
		Email:      email,
	}, billettholderID
}

func insertProfileProgramPulje(t *testing.T, db *sql.DB, puljeID models.Pulje, status models.PuljeStatus) {
	t.Helper()

	mustExecProfileProgramTest(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES(?, ?, ?, ?, ?)
	`, puljeID, "Fredag kveld", status, "2026-10-09T18:00:00Z", "2026-10-09T23:00:00Z")
}

func insertProfileProgramPublishedEvent(t *testing.T, db *sql.DB, eventID string, title string) {
	t.Helper()

	mustExecProfileProgramTest(t, db, `
		INSERT INTO events(
			id,
			title,
			intro,
			description,
			host_name,
			email,
			phone_number,
			max_players,
			status
		)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, eventID, title, "Intro for "+title, "Description for "+title, "Host", "host@example.com", "12345678", 4, models.EventStatusAnnounced)
	mustExecProfileProgramTest(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published)
		VALUES(?, ?, 1, 1)
	`, eventID, models.PuljeFredagKveld)
}

func insertProfileProgramPlayer(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, billettholderID int, role models.EventPlayerRole) {
	t.Helper()

	mustExecProfileProgramTest(t, db, `
		INSERT INTO relation_events_players(event_id, pulje_id, billettholder_id, role)
		VALUES(?, ?, ?, ?)
	`, eventID, puljeID, billettholderID, role)
}

func insertProfileProgramInterest(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, billettholderID int, interestLevel models.InterestLevel) {
	t.Helper()

	mustExecProfileProgramTest(t, db, `
		INSERT INTO interests(billettholder_id, event_id, pulje_id, interest_level)
		VALUES(?, ?, ?, ?)
	`, billettholderID, eventID, puljeID, interestLevel)
}

func assertProfileProgramEventTitles(t *testing.T, expectedTitles []string, events []UserEvent) {
	t.Helper()

	actualTitles := make([]string, 0, len(events))
	for _, event := range events {
		actualTitles = append(actualTitles, event.Title)
	}

	if !slices.Equal(expectedTitles, actualTitles) {
		t.Fatalf("event titles mismatch\nexpected: %v\nactual:   %v", expectedTitles, actualTitles)
	}
}

func assertProfileProgramInterestNames(t *testing.T, expectedNames []string, interests []UserInterest) {
	t.Helper()

	actualNames := make([]string, 0, len(interests))
	for _, interest := range interests {
		actualNames = append(actualNames, interest.EventName)
	}

	if !slices.Equal(expectedNames, actualNames) {
		t.Fatalf("interest names mismatch\nexpected: %v\nactual:   %v", expectedNames, actualNames)
	}
}

func assertProfileProgramEventsAreGM(t *testing.T, events []UserEvent) {
	t.Helper()

	for _, event := range events {
		if !event.IsGM {
			t.Fatalf("expected event %q to be marked as GM", event.Title)
		}
	}
}

func profileProgramVisibleText(doc *goquery.Document) string {
	return strings.Join(strings.Fields(doc.Text()), " ")
}

func mustExecProfileProgramTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("failed to execute query: %v\nquery: %s", err, query)
	}
}
