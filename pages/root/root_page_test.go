package root

import (
	"database/sql"
	"slices"
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

func TestRootPageContent_WhenProgramPublishingIsOn_ShowsScrollnav(t *testing.T) {
	// Gitt at publisering av program er skrudd på,
	// når forsiden vises,
	// så skal puljefilteret vises.

	// Given
	expectedScrollnavVisible := true

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualScrollnavVisible := templtest.HasSelector(doc, ".program-scrollnav-container")

	// Then
	if actualScrollnavVisible != expectedScrollnavVisible {
		t.Fatalf("scrollnav visibility mismatch\nexpected: %v\nactual:   %v", expectedScrollnavVisible, actualScrollnavVisible)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOn_OnlyShowsPublishedPuljeEvents(t *testing.T) {
	// Gitt at publisering av program er skrudd på,
	// når forsiden vises,
	// så skal puljevisningen bare vise godkjente arrangementer som er publisert i en pulje.

	// Given
	expectedTitles := []string{"Published Approved"}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)

	insertRootPageEvent(t, db, "published-approved", "Published Approved", models.EventStatusApproved)
	insertRootPageEventPulje(t, db, "published-approved", models.PuljeFredagKveld, true)

	insertRootPageEvent(t, db, "unpublished-approved", "Unpublished Approved", models.EventStatusApproved)
	insertRootPageEventPulje(t, db, "unpublished-approved", models.PuljeFredagKveld, false)

	insertRootPageEvent(t, db, "unrelated-approved", "Unrelated Approved", models.EventStatusApproved)

	insertRootPageEvent(t, db, "published-submitted", "Published Submitted", models.EventStatusSubmitted)
	insertRootPageEventPulje(t, db, "published-submitted", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualTitles := templtest.CollectTexts(doc, ".event-card-title")

	// Then
	if !slices.Equal(expectedTitles, actualTitles) {
		t.Fatalf("event titles mismatch\nexpected: %v\nactual:   %v", expectedTitles, actualTitles)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOn_RendersPuljeSectionsInTimeOrder(t *testing.T) {
	// Gitt at publisering av program er skrudd på,
	// når forsiden vises,
	// så skal arrangementene grupperes i puljer sortert etter starttidspunkt.

	// Given
	expectedPuljeHeadings := []string{
		"Fredag kveld (18:00 - 23:00)",
		"Lordag morgen (10:00 - 14:00)",
	}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePuljeWithDetails(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-10-09T18:00:00Z", "2026-10-09T23:00:00Z")
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagMorgen, "Lordag morgen", "2026-10-10T10:00:00Z", "2026-10-10T14:00:00Z")

	insertRootPageEvent(t, db, "lordag-event", "Lordag Event", models.EventStatusApproved)
	insertRootPageEventPulje(t, db, "lordag-event", models.PuljeLordagMorgen, true)

	insertRootPageEvent(t, db, "fredag-event", "Fredag Event", models.EventStatusApproved)
	insertRootPageEventPulje(t, db, "fredag-event", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualPuljeHeadings := templtest.CollectTexts(doc, ".pulje-heading")

	// Then
	if !slices.Equal(expectedPuljeHeadings, actualPuljeHeadings) {
		t.Fatalf("pulje headings mismatch\nexpected: %v\nactual:   %v", expectedPuljeHeadings, actualPuljeHeadings)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOn_SortsEventsAlphabeticallyWithinPulje(t *testing.T) {
	// Gitt at publisering av program er skrudd på,
	// når forsiden vises,
	// så skal arrangementene sorteres alfabetisk innenfor hver pulje.

	// Given
	expectedTitles := []string{"Alpha Event", "Beta Event"}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)

	insertRootPageEvent(t, db, "beta-event", "Beta Event", models.EventStatusApproved)
	insertRootPageEventPulje(t, db, "beta-event", models.PuljeFredagKveld, true)

	insertRootPageEvent(t, db, "alpha-event", "Alpha Event", models.EventStatusApproved)
	insertRootPageEventPulje(t, db, "alpha-event", models.PuljeFredagKveld, true)

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

	for _, status := range []models.PuljeStatus{
		models.PuljeStatusOpen,
		models.PuljeStatusLocked,
		models.PuljeStatusCompleted,
	} {
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

	insertRootPagePuljeWithDetails(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-10-09T18:00:00Z", "2026-10-09T23:00:00Z")
}

func insertRootPagePuljeWithDetails(t *testing.T, db *sql.DB, puljeID models.Pulje, name string, startAt string, endAt string) {
	t.Helper()

	mustExec(t, db, `
		INSERT INTO puljer(id, name, status, start_at, end_at)
		VALUES(?, ?, ?, ?, ?)
	`, puljeID, name, models.PuljeStatusOpen, startAt, endAt)
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

func insertRootPageEventPulje(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje, isPublished bool) {
	t.Helper()

	published := 0
	if isPublished {
		published = 1
	}

	mustExec(t, db, `
		INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published)
		VALUES(?, ?, 1, ?)
	`, eventID, puljeID, published)
}

func mustExec(t *testing.T, db *sql.DB, query string, args ...any) {
	t.Helper()

	if _, err := db.Exec(query, args...); err != nil {
		t.Fatalf("failed to execute query: %v\nquery: %s", err, query)
	}
}
