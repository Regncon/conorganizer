package billettholderadmin

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/billettholder"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestBillettholderAdminOverview_RendersCountsFromPublishedPuljeAssignments(t *testing.T) {
	// Given expected overview counts where unpublished first-choice and GM assignments do not count,
	// when the billettholder admin overview component renders,
	// then the visible counts reflect only published-pulje first-choice and GM/DM status.

	// Given
	expectedTexts := []string{
		"Totalt: 5",
		"Uten førstevalg: 4",
		"GM/DM: 1",
	}

	db, logger, createDBErr := testutil.CreateTemporaryDBAndLogger("billettholder_admin_overview", t)
	if createDBErr != nil {
		t.Fatalf("failed to create test database: %v", createDBErr)
	}
	defer db.Close()

	seedBillettholderOverviewLookups(t, db)
	seedBillettholderOverviewBillettholdere(t, db)
	seedBillettholderOverviewPuljer(t, db)
	seedBillettholderOverviewEvents(t, db)
	seedBillettholderOverviewEventPuljer(t, db)
	seedBillettholderOverviewInterests(t, db)
	seedBillettholderOverviewAssignments(t, db)

	// When
	doc := templtest.Render(t, billettholderAdminOverview(db, logger, billettholderService.BillettholderFilters{}))
	actualTexts := templtest.CollectTexts(doc, ".billettholder-admin-overview-count")

	// Then
	if !reflect.DeepEqual(expectedTexts, actualTexts) {
		t.Fatalf("overview count text mismatch\nexpected: %v\nactual:   %v", expectedTexts, actualTexts)
	}
}

func seedBillettholderOverviewLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderOverviewTest(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?)`, models.EventStatusApproved)
	mustExecBillettholderOverviewTest(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExecBillettholderOverviewTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExecBillettholderOverviewTest(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	mustExecBillettholderOverviewTest(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?)`, models.InterestLevelHigh)
	mustExecBillettholderOverviewTest(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?)`, models.PuljeStatusPublished)
}

func seedBillettholderOverviewBillettholdere(t *testing.T, db *sql.DB) {
	t.Helper()

	for id := 1; id <= 5; id++ {
		mustExecBillettholderOverviewTest(t, db, `
			INSERT INTO billettholdere (
				id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
			) VALUES (?, 'Overview', 'Participant', 1, 'Ticket', 1, ?, ?)
		`, id, 3000+id, 4000+id)
	}
}

func seedBillettholderOverviewPuljer(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderOverviewTest(t, db, `
		INSERT INTO puljer (
			id, name, status, start_at, end_at
		) VALUES
			(?, 'Fredag kveld', ?, '2026-10-09T18:00:00Z', '2026-10-09T23:00:00Z')
	`, models.PuljeFredagKveld, models.PuljeStatusPublished)
}

func seedBillettholderOverviewEvents(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderOverviewTest(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type,
			age_group, event_runtime, host_name, email, phone_number,
			max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES
			('overview-published-first-choice', 'Published First Choice', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('overview-unpublished-first-choice', 'Unpublished First Choice', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('overview-published-gm', 'Published GM', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('overview-unpublished-gm', 'Unpublished GM', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?)
	`,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
	)
}

func seedBillettholderOverviewEventPuljer(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderOverviewTest(t, db, `
		INSERT INTO relation_event_puljer (
			event_id, pulje_id, is_in_pulje, is_published
		) VALUES
			('overview-published-first-choice', ?, 1, 1),
			('overview-unpublished-first-choice', ?, 1, 0),
			('overview-published-gm', ?, 1, 1),
			('overview-unpublished-gm', ?, 1, 0)
	`,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
	)
}

func seedBillettholderOverviewInterests(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderOverviewTest(t, db, `
		INSERT INTO interests (
			billettholder_id, event_id, pulje_id, interest_level
		) VALUES
			(1, 'overview-published-first-choice', ?, ?),
			(2, 'overview-unpublished-first-choice', ?, ?)
	`,
		models.PuljeFredagKveld, models.InterestLevelHigh,
		models.PuljeFredagKveld, models.InterestLevelHigh,
	)
}

func seedBillettholderOverviewAssignments(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderOverviewTest(t, db, `
		INSERT INTO relation_events_players (
			event_id, pulje_id, billettholder_id, role
		) VALUES
			('overview-published-first-choice', ?, 1, ?),
			('overview-unpublished-first-choice', ?, 2, ?),
			('overview-published-gm', ?, 3, ?),
			('overview-unpublished-gm', ?, 4, ?)
	`,
		models.PuljeFredagKveld, models.EventPlayerRolePlayer,
		models.PuljeFredagKveld, models.EventPlayerRolePlayer,
		models.PuljeFredagKveld, models.EventPlayerRoleGM,
		models.PuljeFredagKveld, models.EventPlayerRoleGM,
	)
}

func mustExecBillettholderOverviewTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
