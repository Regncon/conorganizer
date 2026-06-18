package event

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func createEventVisibilityTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db := testutil.CreateTestDB(t, "event_visibility")
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?), (?), (?), (?)`, models.EventStatusApproved, models.EventStatusAnnounced, models.EventStatusDraft, models.EventStatusArchived)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?)`, models.PuljeStatusOpen)
	mustExecEventVisibilityTest(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
	setEventVisibilityProgramPublishing(t, db, true)

	return db
}

func seedEventVisibilityUser(t *testing.T, db *sql.DB, userID int, externalID string) {
	t.Helper()

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO users (id, external_id, email, is_admin)
		VALUES (?, ?, ?, 0)
	`, userID, externalID, externalID+"@example.com")
}

func seedEventVisibilityEvent(t *testing.T, db *sql.DB, eventID string, title string, status models.EventStatus, userID sql.NullInt64) {
	t.Helper()

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type,
			age_group, event_runtime, host_name, user_id, email, phone_number,
			max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES (?, ?, 'intro', 'description', '', ?, ?, ?, 'Host', ?, 'host@example.com', '11111111', 4, 1, 1, ?)
	`, eventID, title, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, userID, status)
}

func seedEventVisibilityPulje(t *testing.T, db *sql.DB, puljeID models.Pulje) {
	t.Helper()

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at)
		VALUES (?, 'Fredag kveld', ?, '2026-10-09T18:30:00+02:00', '2026-10-09T23:00:00+02:00')
	`, puljeID, models.PuljeStatusOpen)
}

func seedEventVisibilityEventPulje(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, published bool) {
	t.Helper()

	isPublished := 0
	if published {
		isPublished = 1
	}

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published)
		VALUES (?, ?, 1, ?)
	`, eventID, puljeID, isPublished)
}

func setEventVisibilityProgramPublishing(t *testing.T, db *sql.DB, programPublished bool) {
	t.Helper()

	isPublished := 0
	if programPublished {
		isPublished = 1
	}

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO program_publishing_state(id, is_published)
		VALUES (1, ?)
		ON CONFLICT(id) DO UPDATE SET is_published = excluded.is_published
	`, isPublished)
}

func mustExecEventVisibilityTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
