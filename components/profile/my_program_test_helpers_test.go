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
)

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
