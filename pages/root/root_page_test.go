package root

import (
	"database/sql"
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestRootPageContent_WhenProgramPublishingIsOff_HidesScrollnav(t *testing.T) {
	// Gitt at publisering av program er skrudd av,
	// når forsiden vises,
	// så skal puljefilteret skjules.

	// Given
	expectedScrollnavVisible := false

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, false)
	insertRootPagePulje(t, db)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualScrollnavVisible := templtest.HasSelector(doc, ".program-scrollnav-container")

	// Then
	if actualScrollnavVisible != expectedScrollnavVisible {
		t.Fatalf("scrollnav visibility mismatch\nexpected: %v\nactual:   %v", expectedScrollnavVisible, actualScrollnavVisible)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOff_OnlyShowsApprovedEvents(t *testing.T) {
	// Gitt at publisering av program er skrudd av,
	// når forsiden vises,
	// så skal den flate arrangementslisten bare vise godkjente arrangementer.

	// Given
	expectedTitles := []string{"Alpha Approved", "Beta Approved"}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, false)
	insertRootPageEvent(t, db, "draft-event", "Draft Event", models.EventStatusDraft)
	insertRootPageEvent(t, db, "submitted-event", "Submitted Event", models.EventStatusSubmitted)
	insertRootPageEvent(t, db, "beta-approved", "Beta Approved", models.EventStatusApproved)
	insertRootPageEvent(t, db, "alpha-approved", "Alpha Approved", models.EventStatusApproved)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualTitles := templtest.CollectTexts(doc, ".event-card-title")

	// Then
	if !slices.Equal(expectedTitles, actualTitles) {
		t.Fatalf("event titles mismatch\nexpected: %v\nactual:   %v", expectedTitles, actualTitles)
	}
}

func createRootPageTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, _, err := testutil.CreateTemporaryDBAndLogger("root_page", t)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Fatalf("failed to close test database: %v", err)
		}
	})

	return db
}

func seedRootPageLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, status := range []models.EventStatus{
		models.EventStatusDraft,
		models.EventStatusSubmitted,
		models.EventStatusApproved,
		models.EventStatusArchived,
	} {
		mustExec(t, db, `INSERT INTO event_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}

	mustExec(t, db, `INSERT INTO events_types(event_type) VALUES (?) ON CONFLICT(event_type) DO NOTHING`, models.EventTypeOther)
	mustExec(t, db, `INSERT INTO age_groups(age_group) VALUES (?) ON CONFLICT(age_group) DO NOTHING`, models.AgeGroupDefault)
	mustExec(t, db, `INSERT INTO event_runtimes(runtime) VALUES (?) ON CONFLICT(runtime) DO NOTHING`, models.RunTimeNormal)

	for _, status := range []string{"not_published", "published", "open", "locked", "completed"} {
		mustExec(t, db, `INSERT INTO pulje_statuses(status) VALUES (?) ON CONFLICT(status) DO NOTHING`, status)
	}
}

func setProgramPublishing(t *testing.T, db *sql.DB, isPublished bool) {
	t.Helper()

	value := 0
	if isPublished {
		value = 1
	}

	mustExec(t, db, `
		INSERT INTO program_publishing_state(id, is_published)
		VALUES(1, ?)
		ON CONFLICT(id) DO UPDATE SET is_published = excluded.is_published
	`, value)
}

func insertRootPagePulje(t *testing.T, db *sql.DB) {
	t.Helper()

	mustExec(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES(?, ?, ?, ?, ?)
	`, models.PuljeFredagKveld, "Fredag kveld", puljeStatusForRootPageTest(t, db), "2026-10-09T18:00:00Z", "2026-10-09T23:00:00Z")
}

func insertRootPageEvent(t *testing.T, db *sql.DB, id string, title string, status models.EventStatus) {
	t.Helper()

	mustExec(t, db, `
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
	`, id, title, "Intro", "Description", "Host", "host@example.com", "12345678", 4, status)
}

func puljeStatusForRootPageTest(t *testing.T, db *sql.DB) string {
	t.Helper()

	var schema string
	if err := db.QueryRow(`
		SELECT sql
		FROM sqlite_schema
		WHERE type = 'table' AND name = 'puljer'
	`).Scan(&schema); err != nil {
		t.Fatalf("failed to read puljer schema: %v", err)
	}

	if strings.Contains(schema, "'open'") {
		return string(models.PuljeStatusOpen)
	}

	return string(models.PuljeStatusPublished)
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("failed to execute query: %v\nquery: %s", err, query)
	}
}
