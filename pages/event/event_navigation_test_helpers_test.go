package event

import (
	"database/sql"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func collectEventNavigationHrefs(doc *goquery.Document) []string {
	uniqueHrefs := map[string]struct{}{}
	doc.Find(".breadcrumb-prev-next-buttons a[href], .previous-next-buttons a[href]").Each(func(_ int, selection *goquery.Selection) {
		href, exists := selection.Attr("href")
		if !exists || !strings.HasPrefix(href, "/event/") {
			return
		}
		uniqueHrefs[href] = struct{}{}
	})

	hrefs := make([]string, 0, len(uniqueHrefs))
	for href := range uniqueHrefs {
		hrefs = append(hrefs, href)
	}
	slices.Sort(hrefs)
	return hrefs
}

func collectEventNavigationText(doc *goquery.Document) string {
	return strings.Join(templtest.CollectTexts(doc, ".breadcrumb-prev-next-buttons, .previous-next-buttons"), " ")
}

func assertNoDisabledEventNavigationPlaceholders(t *testing.T, doc *goquery.Document) {
	t.Helper()

	disabledPlaceholders := doc.Find(`.previous-next-buttons [aria-disabled="true"], .breadcrumb-prev-next-buttons div.btn`)
	if disabledPlaceholders.Length() > 0 {
		t.Fatalf("expected no disabled previous/next placeholders, found %d", disabledPlaceholders.Length())
	}
}

func assertEventNavigationNotRendered(t *testing.T, doc *goquery.Document) {
	t.Helper()

	if templtest.HasSelector(doc, ".breadcrumb-prev-next-buttons") {
		t.Fatalf("expected breadcrumb previous/next navigation not to render")
	}
	if templtest.HasSelector(doc, ".previous-next-buttons") {
		t.Fatalf("expected event previous/next navigation not to render")
	}
}

func seedEventNavigationPulje(t *testing.T, db *sql.DB, puljeID models.Pulje, name string, startAt string, endAt string) {
	t.Helper()

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO puljer (id, name, status, start_at, end_at)
		VALUES (?, ?, ?, ?, ?)
	`, puljeID, name, models.PuljeStatusOpen, startAt, endAt)
}

func seedEventNavigationEventPulje(t *testing.T, db *sql.DB, eventID string, puljeID models.Pulje) {
	t.Helper()

	mustExecEventVisibilityTest(t, db, `
		INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, is_published)
		VALUES (?, ?, 1, 1)
	`, eventID, puljeID)
}
