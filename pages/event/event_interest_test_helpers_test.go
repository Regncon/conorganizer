package event

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

type eventInterestUpdateFixture struct {
	userExternalID  string
	billettholderID int
	eventID         string
	puljeID         models.Pulje
}

func createEventInterestTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db := testutil.CreateTestDB(t, "event_interest")
	seedEventInterestLookups(t, db)
	return db
}

func seedEventInterestLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?), (?)`, models.EventStatusApproved, models.EventStatusAnnounced)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?), (?), (?)`, models.PuljeStatusOpen, models.PuljeStatusLocked, models.PuljeStatusCompleted)
	setEventInterestProgramPublishing(t, db, true)
}

func seedEventInterestUpdateFixture(
	t *testing.T,
	db *sql.DB,
	puljeStatus models.PuljeStatus,
	existingInterest models.InterestLevel,
) eventInterestUpdateFixture {
	t.Helper()

	fixture := eventInterestUpdateFixture{
		userExternalID:  "event-interest-user",
		billettholderID: 901,
		eventID:         "event-interest-event",
		puljeID:         models.PuljeFredagKveld,
	}

	mustExecEventInterestTest(t, db, `
		INSERT INTO users (id, external_id, email, is_admin)
		VALUES (501, ?, 'event-interest-user@example.com', 0)
	`, fixture.userExternalID)
	mustExecEventInterestTest(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, 'Event', 'Interest', 1, 'Ticket', 1, 7001, 8001)
	`, fixture.billettholderID)
	mustExecEventInterestTest(t, db, `
		INSERT INTO relation_billettholdere_users (billettholder_id, user_id)
		VALUES (?, 501)
	`, fixture.billettholderID)
	mustExecEventInterestTest(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at)
		VALUES (?, 'Fredag kveld', ?, '2026-10-09T18:30:00+02:00', '2026-10-09T23:00:00+02:00')
	`, fixture.puljeID, puljeStatus)
	mustExecEventInterestTest(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type,
			age_group, event_runtime, host_name, email, phone_number,
			max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES (?, 'Interest Event', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?)
	`, fixture.eventID, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusAnnounced)
	mustExecEventInterestTest(t, db, `
		INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published)
		VALUES (?, ?, 1, 1)
	`, fixture.eventID, fixture.puljeID)
	mustExecEventInterestTest(t, db, `
		INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
		VALUES (?, ?, ?, ?)
	`, fixture.billettholderID, fixture.eventID, fixture.puljeID, existingInterest)

	return fixture
}

func getEventInterestTestInterest(t *testing.T, db *sql.DB, eventID string, billettholderID int, puljeID models.Pulje) models.InterestLevel {
	t.Helper()

	interest, err := getSelectedInterest(eventID, billettholderID, string(puljeID), db)
	if err != nil {
		t.Fatalf("failed to get selected interest: %v", err)
	}
	return interest
}

func setEventInterestProgramPublishing(t *testing.T, db *sql.DB, programPublished bool) {
	t.Helper()

	isPublished := 0
	if programPublished {
		isPublished = 1
	}

	mustExecEventInterestTest(t, db, `
		INSERT INTO program_publishing_state(id, is_published)
		VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET is_published = excluded.is_published
	`, isPublished)
}

func buildEventInterestTestPulje(id models.Pulje, name string, status models.PuljeStatus, startAt time.Time) models.PuljeRow {
	return models.PuljeRow{
		ID:      id,
		Name:    name,
		Status:  status,
		StartAt: models.NewDBDateTime(startAt),
		EndAt:   models.NewDBDateTime(startAt.Add(4 * time.Hour)),
	}
}

func mustExecEventInterestTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
