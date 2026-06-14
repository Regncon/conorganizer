package root

import (
	"database/sql"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func createRootPageTestDB(t *testing.T) *sql.DB {
	t.Helper()

	return testutil.CreateTestDB(t, "root_page")
}

func seedRootPageLookups(t *testing.T, db *sql.DB) {
	t.Helper()

	for _, status := range []models.EventStatus{
		models.EventStatusDraft,
		models.EventStatusSubmitted,
		models.EventStatusApproved,
		models.EventStatusArchived,
		models.EventStatusAnnounced,
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

func collectRootPageHrefs(doc *goquery.Document, selector string) []string {
	hrefs := make([]string, 0)
	doc.Find(selector).Each(func(_ int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if exists {
			hrefs = append(hrefs, href)
		}
	})
	slices.Sort(hrefs)
	return hrefs
}

func rootPageText(doc *goquery.Document) string {
	return strings.Join(strings.Fields(doc.Text()), " ")
}

func assertTextContains(t *testing.T, actualText string, expectedTextPart string) {
	t.Helper()

	if !strings.Contains(actualText, expectedTextPart) {
		t.Fatalf("text mismatch\nexpected text to contain: %q\nactual text:              %q", expectedTextPart, actualText)
	}
}

func assertTextDoesNotContain(t *testing.T, actualText string, unexpectedTextPart string) {
	t.Helper()

	if strings.Contains(actualText, unexpectedTextPart) {
		t.Fatalf("text mismatch\nexpected text not to contain: %q\nactual text:                  %q", unexpectedTextPart, actualText)
	}
}
