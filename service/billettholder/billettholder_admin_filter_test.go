package billettholderService

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestGetBillettholdereWithFilters_WhenFilteringWithoutFirstChoiceAndGM_ReturnsPublishedGMWithoutPublishedFirstChoice(t *testing.T) {
	// Given expected billettholdere who are GM/DM in a published pulje and have not received a published first choice,
	// when the admin billettholder query applies both filters,
	// then only billettholdere matching both published-event filters are returned.

	// Given
	expectedBillettholderIDs := []int{3}
	filters := BillettholderFilters{
		WithoutFirstChoice: true,
		GMOrDM:             true,
	}

	db := testutil.CreateTestDB(t, "billettholder_admin_filters")
	seedBillettholderFilterLookups(t, db)
	seedBillettholderFilterBillettholdere(t, db)
	seedBillettholderFilterPuljer(t, db)
	seedBillettholderFilterEvents(t, db)
	seedBillettholderFilterEventPuljer(t, db)
	seedBillettholderFilterInterests(t, db)
	seedBillettholderFilterAssignments(t, db)

	// When
	actualBillettholdere, err := GetBillettholdereWithFilters("", db, filters)

	// Then
	if err != nil {
		t.Fatalf("expected filtered billettholdere to load: %v", err)
	}

	actualBillettholderIDs := collectBillettholderIDs(actualBillettholdere)
	if !reflect.DeepEqual(expectedBillettholderIDs, actualBillettholderIDs) {
		t.Fatalf("filtered billettholder IDs mismatch\nexpected: %v\nactual:   %v", expectedBillettholderIDs, actualBillettholderIDs)
	}
}

func collectBillettholderIDs(billettholdere []models.Billettholder) []int {
	ids := make([]int, 0, len(billettholdere))
	for _, billettholder := range billettholdere {
		ids = append(ids, billettholder.ID)
	}
	return ids
}

func seedBillettholderFilterLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderFilterTest(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?)`, models.EventStatusApproved)
	mustExecBillettholderFilterTest(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExecBillettholderFilterTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExecBillettholderFilterTest(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	mustExecBillettholderFilterTest(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?)`, models.InterestLevelHigh)
	mustExecBillettholderFilterTest(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?)`, models.PuljeStatusOpen)
}

func seedBillettholderFilterBillettholdere(t *testing.T, db *sql.DB) {
	t.Helper()

	for id := 1; id <= 5; id++ {
		mustExecBillettholderFilterTest(t, db, `
			INSERT INTO billettholdere (
				id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
			) VALUES (?, ?, 'Participant', 1, 'Ticket', 1, ?, ?)
		`, id, billettholderFilterFirstName(id), 1000+id, 2000+id)
	}
}

func billettholderFilterFirstName(id int) string {
	switch id {
	case 1:
		return "PublishedFirstChoice"
	case 2:
		return "UnpublishedFirstChoice"
	case 3:
		return "PublishedGMWithoutFirstChoice"
	case 4:
		return "UnpublishedGMWithoutFirstChoice"
	case 5:
		return "PublishedGMWithFirstChoice"
	default:
		return "Participant"
	}
}

func seedBillettholderFilterPuljer(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderFilterTest(t, db, `
		INSERT INTO puljer (
			id, name, status, start_at, end_at
		) VALUES
			(?, 'Fredag kveld', ?, '2026-10-09T18:00:00Z', '2026-10-09T23:00:00Z')
	`, models.PuljeFredagKveld, models.PuljeStatusOpen)
}

func seedBillettholderFilterEvents(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderFilterTest(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type,
			age_group, event_runtime, host_name, email, phone_number,
			max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES
			('published-first-choice', 'Published First Choice', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('unpublished-first-choice', 'Unpublished First Choice', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('published-gm', 'Published GM', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('unpublished-gm', 'Unpublished GM', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('published-gm-with-first-choice', 'Published GM With First Choice', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('published-first-choice-for-gm', 'Published First Choice For GM', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?)
	`,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
	)
}

func seedBillettholderFilterEventPuljer(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderFilterTest(t, db, `
		INSERT INTO relation_event_puljer (
			event_id, pulje_id, is_in_pulje, is_published
		) VALUES
			('published-first-choice', ?, 1, 1),
			('unpublished-first-choice', ?, 1, 0),
			('published-gm', ?, 1, 1),
			('unpublished-gm', ?, 1, 0),
			('published-gm-with-first-choice', ?, 1, 1),
			('published-first-choice-for-gm', ?, 1, 1)
	`,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
	)
}

func seedBillettholderFilterInterests(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderFilterTest(t, db, `
		INSERT INTO interests (
			billettholder_id, event_id, pulje_id, interest_level
		) VALUES
			(1, 'published-first-choice', ?, ?),
			(2, 'unpublished-first-choice', ?, ?),
			(5, 'published-first-choice-for-gm', ?, ?)
	`,
		models.PuljeFredagKveld, models.InterestLevelHigh,
		models.PuljeFredagKveld, models.InterestLevelHigh,
		models.PuljeFredagKveld, models.InterestLevelHigh,
	)
}

func seedBillettholderFilterAssignments(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderFilterTest(t, db, `
		INSERT INTO relation_events_players (
			event_id, pulje_id, billettholder_id, role
		) VALUES
			('published-first-choice', ?, 1, ?),
			('unpublished-first-choice', ?, 2, ?),
			('published-gm', ?, 3, ?),
			('unpublished-gm', ?, 4, ?),
			('published-gm-with-first-choice', ?, 5, ?),
			('published-first-choice-for-gm', ?, 5, ?)
	`,
		models.PuljeFredagKveld, models.EventPlayerRolePlayer,
		models.PuljeFredagKveld, models.EventPlayerRolePlayer,
		models.PuljeFredagKveld, models.EventPlayerRoleGM,
		models.PuljeFredagKveld, models.EventPlayerRoleGM,
		models.PuljeFredagKveld, models.EventPlayerRoleGM,
		models.PuljeFredagKveld, models.EventPlayerRolePlayer,
	)
}

func mustExecBillettholderFilterTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
