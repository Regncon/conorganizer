package billettholderadmin

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestGetBillettholderInterestSectionsByBillettholderID_ReturnsAssignedAndInterestedEventsGroupedByPulje(t *testing.T) {
	// Given expected pulje sections with assigned rows separated from interest rows and pulje publication state,
	// when the real billettholder interest loader reads the database,
	// then each requested billettholder gets assigned rows first and interest rows grouped by level with the correct pulje-specific publication status.

	// Given
	expectedBillettholderID := 42
	expectedSections := []expectedBillettholderInterestSection{
		{
			PuljeID: models.PuljeFredagKveld,
			Name:    "Fredag kveld",
			Assigned: []expectedBillettholderInterestRow{
				{
					EventID:       "assigned-gm-without-interest",
					EventTitle:    "Assigned GM Without Interest",
					EventStatus:   models.EventStatusPublished,
					IsPublished:   true,
					InterestLevel: models.InterestLevelNone,
					AssignedRole:  models.EventPlayerRoleGM,
				},
				{
					EventID:       "assigned-player-first-choice",
					EventTitle:    "Assigned Player First Choice",
					EventStatus:   models.EventStatusPublished,
					IsPublished:   true,
					InterestLevel: models.InterestLevelHigh,
					AssignedRole:  models.EventPlayerRolePlayer,
				},
			},
			High: []expectedBillettholderInterestRow{
				{
					EventID:       "high-interest",
					EventTitle:    "High Interest",
					EventStatus:   models.EventStatusPublished,
					IsPublished:   false,
					InterestLevel: models.InterestLevelHigh,
				},
			},
			Medium: []expectedBillettholderInterestRow{
				{
					EventID:       "medium-interest",
					EventTitle:    "Medium Interest",
					EventStatus:   models.EventStatusApproved,
					IsPublished:   true,
					InterestLevel: models.InterestLevelMedium,
				},
			},
			Low: []expectedBillettholderInterestRow{
				{
					EventID:       "low-interest",
					EventTitle:    "Low Interest",
					EventStatus:   models.EventStatusDraft,
					IsPublished:   false,
					InterestLevel: models.InterestLevelLow,
				},
			},
		},
		{
			PuljeID: models.PuljeLordagMorgen,
			Name:    "Lørdag morgen",
			Assigned: []expectedBillettholderInterestRow{
				{
					EventID:       "saturday-assigned-player",
					EventTitle:    "Saturday Assigned Player",
					EventStatus:   models.EventStatusPublished,
					IsPublished:   true,
					InterestLevel: models.InterestLevelNone,
					AssignedRole:  models.EventPlayerRolePlayer,
				},
			},
			High: []expectedBillettholderInterestRow{
				{
					EventID:       "saturday-high-interest",
					EventTitle:    "Saturday High Interest",
					EventStatus:   models.EventStatusApproved,
					IsPublished:   true,
					InterestLevel: models.InterestLevelHigh,
				},
			},
		},
	}

	db, _, createDBErr := testutil.CreateTemporaryDBAndLogger("billettholder_interests", t)
	if createDBErr != nil {
		t.Fatalf("failed to create test database: %v", createDBErr)
	}
	defer db.Close()

	seedBillettholderInterestLookups(t, db)
	seedBillettholderInterestBillettholdere(t, db, expectedBillettholderID)
	seedBillettholderInterestPuljer(t, db)
	seedBillettholderInterestEvents(t, db)
	seedBillettholderInterestEventPuljer(t, db)
	seedBillettholderInterestRows(t, db, expectedBillettholderID)
	seedBillettholderInterestAssignments(t, db, expectedBillettholderID)

	// When
	actualSectionsByBillettholderID, err := getBillettholderInterestSectionsByBillettholderID(db, []int{expectedBillettholderID})

	// Then
	if err != nil {
		t.Fatalf("expected billettholder interest sections to load: %v", err)
	}

	actualSections, ok := actualSectionsByBillettholderID[expectedBillettholderID]
	if !ok {
		t.Fatalf("expected sections for billettholder ID %d", expectedBillettholderID)
	}

	assertBillettholderInterestSections(t, expectedSections, actualSections)
}

type expectedBillettholderInterestSection struct {
	PuljeID  models.Pulje
	Name     string
	Assigned []expectedBillettholderInterestRow
	High     []expectedBillettholderInterestRow
	Medium   []expectedBillettholderInterestRow
	Low      []expectedBillettholderInterestRow
}

type expectedBillettholderInterestRow struct {
	EventID       string
	EventTitle    string
	EventStatus   models.EventStatus
	IsPublished   bool
	InterestLevel models.InterestLevel
	AssignedRole  models.EventPlayerRole
}

func assertBillettholderInterestSections(
	t *testing.T,
	expected []expectedBillettholderInterestSection,
	actual []billettholderInterestPuljeSection,
) {
	t.Helper()

	normalizedActual := make([]expectedBillettholderInterestSection, 0, len(actual))
	for _, section := range actual {
		normalizedActual = append(normalizedActual, expectedBillettholderInterestSection{
			PuljeID:  section.PuljeID,
			Name:     section.Name,
			Assigned: normalizeBillettholderInterestRows(section.Assigned),
			High:     normalizeBillettholderInterestRows(section.High),
			Medium:   normalizeBillettholderInterestRows(section.Medium),
			Low:      normalizeBillettholderInterestRows(section.Low),
		})
	}

	if !reflect.DeepEqual(expected, normalizedActual) {
		t.Fatalf("billettholder interest sections mismatch\nexpected: %#v\nactual:   %#v", expected, normalizedActual)
	}
}

func normalizeBillettholderInterestRows(rows []billettholderInterestEventRow) []expectedBillettholderInterestRow {
	if len(rows) == 0 {
		return nil
	}

	normalized := make([]expectedBillettholderInterestRow, 0, len(rows))
	for _, row := range rows {
		normalized = append(normalized, expectedBillettholderInterestRow{
			EventID:       row.EventID,
			EventTitle:    row.EventTitle,
			EventStatus:   row.EventStatus,
			IsPublished:   row.IsPublished,
			InterestLevel: row.InterestLevel,
			AssignedRole:  row.AssignedRole,
		})
	}
	return normalized
}

func seedBillettholderInterestLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderInterestTest(t, db, `INSERT OR IGNORE INTO event_statuses(status) VALUES (?), (?), (?)`, models.EventStatusPublished, models.EventStatusApproved, models.EventStatusDraft)
	mustExecBillettholderInterestTest(t, db, `INSERT OR IGNORE INTO events_types(event_type) VALUES (?)`, models.EventTypeOther)
	mustExecBillettholderInterestTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupDefault)
	mustExecBillettholderInterestTest(t, db, `INSERT OR IGNORE INTO event_runtimes(runtime) VALUES (?)`, models.RunTimeNormal)
	mustExecBillettholderInterestTest(t, db, `INSERT OR IGNORE INTO interest_levels(interest_level) VALUES (?), (?), (?)`, models.InterestLevelHigh, models.InterestLevelMedium, models.InterestLevelLow)
	mustExecBillettholderInterestTest(t, db, `INSERT OR IGNORE INTO pulje_statuses(status) VALUES (?)`, models.PuljeStatusPublished)
}

func seedBillettholderInterestBillettholdere(t *testing.T, db *sql.DB, billettholderID int) {
	t.Helper()

	mustExecBillettholderInterestTest(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES
			(?, 'Test', 'Billettholder', 1, 'Ticket', 1, 1001, 2001),
			(?, 'Other', 'Billettholder', 1, 'Ticket', 1, 1002, 2002)
	`, billettholderID, billettholderID+1)
}

func seedBillettholderInterestPuljer(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderInterestTest(t, db, `
		INSERT INTO puljer (
			id, name, status, start_at, end_at
		) VALUES
			(?, 'Fredag kveld', ?, '2026-10-09T18:00:00Z', '2026-10-09T23:00:00Z'),
			(?, 'Lørdag morgen', ?, '2026-10-10T10:00:00Z', '2026-10-10T14:00:00Z')
	`, models.PuljeFredagKveld, models.PuljeStatusPublished, models.PuljeLordagMorgen, models.PuljeStatusPublished)
}

func seedBillettholderInterestEvents(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderInterestTest(t, db, `
		INSERT INTO events (
			id, title, intro, description, system, event_type,
			age_group, event_runtime, host_name, email, phone_number,
			max_players, beginner_friendly, can_be_run_in_english,
			status
		) VALUES
			('assigned-gm-without-interest', 'Assigned GM Without Interest', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('assigned-player-first-choice', 'Assigned Player First Choice', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('high-interest', 'High Interest', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('medium-interest', 'Medium Interest', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('low-interest', 'Low Interest', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('saturday-assigned-player', 'Saturday Assigned Player', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('saturday-high-interest', 'Saturday High Interest', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?),
			('other-billettholder-interest', 'Other Billettholder Interest', 'intro', 'description', '', ?, ?, ?, 'Host', 'host@example.com', '11111111', 4, 1, 1, ?)
	`,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusPublished,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusPublished,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusPublished,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusDraft,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusPublished,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusApproved,
		models.EventTypeOther, models.AgeGroupDefault, models.RunTimeNormal, models.EventStatusPublished,
	)
}

func seedBillettholderInterestEventPuljer(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExecBillettholderInterestTest(t, db, `
		INSERT INTO relation_event_puljer (
			event_id, pulje_id, is_in_pulje, is_published
		) VALUES
			('assigned-gm-without-interest', ?, 1, 1),
			('assigned-player-first-choice', ?, 1, 1),
			('high-interest', ?, 1, 0),
			('medium-interest', ?, 1, 1),
			('low-interest', ?, 1, 0),
			('saturday-assigned-player', ?, 1, 1),
			('saturday-high-interest', ?, 1, 0),
			('saturday-high-interest', ?, 1, 1),
			('other-billettholder-interest', ?, 1, 1)
	`,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeFredagKveld,
		models.PuljeLordagMorgen,
		models.PuljeFredagKveld,
		models.PuljeLordagMorgen,
		models.PuljeFredagKveld,
	)
}

func seedBillettholderInterestRows(t *testing.T, db *sql.DB, billettholderID int) {
	t.Helper()

	mustExecBillettholderInterestTest(t, db, `
		INSERT INTO interests (
			billettholder_id, event_id, pulje_id, interest_level
		) VALUES
			(?, 'assigned-player-first-choice', ?, ?),
			(?, 'high-interest', ?, ?),
			(?, 'medium-interest', ?, ?),
			(?, 'low-interest', ?, ?),
			(?, 'saturday-high-interest', ?, ?),
			(?, 'other-billettholder-interest', ?, ?)
	`,
		billettholderID, models.PuljeFredagKveld, models.InterestLevelHigh,
		billettholderID, models.PuljeFredagKveld, models.InterestLevelHigh,
		billettholderID, models.PuljeFredagKveld, models.InterestLevelMedium,
		billettholderID, models.PuljeFredagKveld, models.InterestLevelLow,
		billettholderID, models.PuljeLordagMorgen, models.InterestLevelHigh,
		billettholderID+1, models.PuljeFredagKveld, models.InterestLevelHigh,
	)
}

func seedBillettholderInterestAssignments(t *testing.T, db *sql.DB, billettholderID int) {
	t.Helper()

	mustExecBillettholderInterestTest(t, db, `
		INSERT INTO relation_events_players (
			event_id, pulje_id, billettholder_id, role
		) VALUES
			('assigned-gm-without-interest', ?, ?, ?),
			('assigned-player-first-choice', ?, ?, ?),
			('saturday-assigned-player', ?, ?, ?)
	`,
		models.PuljeFredagKveld, billettholderID, models.EventPlayerRoleGM,
		models.PuljeFredagKveld, billettholderID, models.EventPlayerRolePlayer,
		models.PuljeLordagMorgen, billettholderID, models.EventPlayerRolePlayer,
	)
}

func mustExecBillettholderInterestTest(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("exec failed: %v\nquery:\n%s", err, query)
	}
}
