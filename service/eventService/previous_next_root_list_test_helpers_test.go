package eventservice

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

type expectedPreviousNext struct {
	previousURL   string
	previousTitle string
	nextURL       string
	nextTitle     string
}

func assertPreviousNextMatches(t *testing.T, expected expectedPreviousNext, actual components.PreviousNext) {
	t.Helper()

	if actual.PreviousUrl != expected.previousURL {
		t.Fatalf("previous URL mismatch\nexpected: %q\nactual:   %q", expected.previousURL, actual.PreviousUrl)
	}
	if actual.PreviousTitle != expected.previousTitle {
		t.Fatalf("previous title mismatch\nexpected: %q\nactual:   %q", expected.previousTitle, actual.PreviousTitle)
	}
	if actual.NextUrl != expected.nextURL {
		t.Fatalf("next URL mismatch\nexpected: %q\nactual:   %q", expected.nextURL, actual.NextUrl)
	}
	if actual.NextTitle != expected.nextTitle {
		t.Fatalf("next title mismatch\nexpected: %q\nactual:   %q", expected.nextTitle, actual.NextTitle)
	}
}

func collectPreviousNextRootListEventIDs(events []models.EventCardModel) []string {
	ids := make([]string, 0, len(events))
	for _, event := range events {
		ids = append(ids, event.Id)
	}
	return ids
}

func createPreviousNextRootListTestDB(t *testing.T) *sql.DB {
	t.Helper()

	return testutil.CreateTestDB(t, "previous-next-root-list")
}

func seedPreviousNextRootListLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, status := range []models.EventStatus{
		models.EventStatusDraft,
		models.EventStatusSubmitted,
		models.EventStatusApproved,
		models.EventStatusArchived,
		models.EventStatusAnnounced,
	} {
		mustExecPreviousNextRootList(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}

	mustExecPreviousNextRootList(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeOther)
	mustExecPreviousNextRootList(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupDefault)
	mustExecPreviousNextRootList(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeNormal)
	mustExecPreviousNextRootList(t, db, `INSERT INTO pulje_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, models.PuljeStatusOpen)
}

func seedPreviousNextRootListPulje(t *testing.T, db *sql.DB, puljeID models.Pulje, name string, startAt string, endAt string) {
	t.Helper()

	mustExecPreviousNextRootList(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES(?, ?, ?, ?, ?)
	`, puljeID, name, models.PuljeStatusOpen, startAt, endAt)
}

func seedPreviousNextRootListEvent(t *testing.T, db *sql.DB, id string, title string, status models.EventStatus) {
	t.Helper()

	mustExecPreviousNextRootList(t, db, `
		INSERT INTO events(
			id,
			title,
			intro,
			description,
			event_type,
			age_group,
			event_runtime,
			host_name,
			email,
			phone_number,
			max_players,
			beginner_friendly,
			can_be_run_in_english,
			status
		)
		VALUES(?, ?, 'Intro', 'Description', ?, ?, ?, 'Host', 'host@example.com', '12345678', 4, 1, 1, ?)
	`, id, title, models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, status)
}

func seedPreviousNextRootListEventPulje(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, isInPulje bool, isPublished bool) {
	t.Helper()

	inPulje := 0
	if isInPulje {
		inPulje = 1
	}
	published := 0
	if isPublished {
		published = 1
	}

	mustExecPreviousNextRootList(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published)
		VALUES(?, ?, ?, ?)
	`, eventID, puljeID, inPulje, published)
}

func mustExecPreviousNextRootList(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
