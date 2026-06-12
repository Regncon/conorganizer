package root

import (
	"database/sql"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestRootPageContent_RendersHomeBreadcrumb(t *testing.T) {
	// Gitt at brukeren åpner forsiden,
	// når forsiden vises,
	// så skal brødsmulestien vise Hjem som gjeldende side.

	// Given
	expectedBreadcrumb := []string{"Hjem"}

	db := createRootPageTestDB(t)
	setProgramPublishing(t, db, false)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualBreadcrumb := templtest.CollectTexts(doc, ".breadcrumb-end")

	// Then
	if !slices.Equal(expectedBreadcrumb, actualBreadcrumb) {
		t.Fatalf("breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedBreadcrumb, actualBreadcrumb)
	}
}

func TestRootPageContent_RendersSubmitEventCallToAction(t *testing.T) {
	// Gitt at brukeren åpner forsiden,
	// når innsendingseksjonen vises,
	// så skal den gi en tydelig inngang til å sende inn arrangement.

	// Given
	expectedTextParts := []string{
		"Vil du arrangere noe under Regncon?",
		"Send inn arrangement",
	}
	expectedHref := "/profile"
	expectedImageSrc := "/static/call-to-action-avatar.webp"
	expectedImageAlt := "Sent inn arrangement"

	db := createRootPageTestDB(t)
	setProgramPublishing(t, db, false)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualText := strings.Join(templtest.CollectTexts(doc, ".call-to-action"), " ")
	actualHref, actualHrefExists := doc.Find(".call-to-action a").Attr("href")
	actualImageSrc, actualImageSrcExists := doc.Find(".call-to-action img.call-to-action-avatar").Attr("src")
	actualImageAlt, actualImageAltExists := doc.Find(".call-to-action img.call-to-action-avatar").Attr("alt")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		assertTextContains(t, actualText, expectedTextPart)
	}
	if !actualHrefExists || actualHref != expectedHref {
		t.Fatalf("CTA href mismatch\nexpected: %q\nactual:   %q", expectedHref, actualHref)
	}
	if !actualImageSrcExists || actualImageSrc != expectedImageSrc {
		t.Fatalf("CTA image src mismatch\nexpected: %q\nactual:   %q", expectedImageSrc, actualImageSrc)
	}
	if !actualImageAltExists || actualImageAlt != expectedImageAlt {
		t.Fatalf("CTA image alt mismatch\nexpected: %q\nactual:   %q", expectedImageAlt, actualImageAlt)
	}
}

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

func TestRootPageContent_WhenProgramPublishingIsOff_OnlyShowsAnnouncedEvents(t *testing.T) {
	// Gitt at publisering av program er skrudd av,
	// når forsiden vises,
	// så skal den flate arrangementslisten bare vise annonserte arrangementer.

	// Given
	expectedTitles := []string{"Alpha Announced", "Beta Announced"}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, false)
	insertRootPageEvent(t, db, "draft-event", "Draft Event", models.EventStatusDraft)
	insertRootPageEvent(t, db, "submitted-event", "Submitted Event", models.EventStatusSubmitted)
	insertRootPageEvent(t, db, "approved-event", "Approved Event", models.EventStatusApproved)
	insertRootPageEvent(t, db, "beta-announced", "Beta Announced", models.EventStatusAnnounced)
	insertRootPageEvent(t, db, "alpha-announced", "Alpha Announced", models.EventStatusAnnounced)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualTitles := templtest.CollectTexts(doc, ".event-card-title")

	// Then
	if !slices.Equal(expectedTitles, actualTitles) {
		t.Fatalf("event titles mismatch\nexpected: %v\nactual:   %v", expectedTitles, actualTitles)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOff_RendersEventLinksWithoutPulje(t *testing.T) {
	// Gitt at programmet ikke er publisert,
	// når annonserte arrangementer vises på forsiden,
	// så skal arrangementskortene lenke direkte til arrangementssidene uten puljekontekst.

	// Given
	expectedHrefs := []string{"/event/alpha-announced", "/event/beta-announced"}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, false)
	insertRootPageEvent(t, db, "beta-announced", "Beta Announced", models.EventStatusAnnounced)
	insertRootPageEvent(t, db, "alpha-announced", "Alpha Announced", models.EventStatusAnnounced)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualHrefs := collectRootPageHrefs(doc, ".event-card-container")

	// Then
	if !slices.Equal(expectedHrefs, actualHrefs) {
		t.Fatalf("event card hrefs mismatch\nexpected: %v\nactual:   %v", expectedHrefs, actualHrefs)
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

func TestRootPageContent_WhenProgramPublishingIsOn_OnlyShowsAnnouncedPublishedPuljeEvents(t *testing.T) {
	// Gitt at publisering av program er skrudd på,
	// når forsiden vises,
	// så skal puljevisningen bare vise annonserte arrangementer som er publisert i en pulje.

	// Given
	expectedTitles := []string{"Published Announced"}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)

	insertRootPageEvent(t, db, "published-announced", "Published Announced", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "published-announced", models.PuljeFredagKveld, true)

	insertRootPageEvent(t, db, "unpublished-announced", "Unpublished Announced", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "unpublished-announced", models.PuljeFredagKveld, false)

	insertRootPageEvent(t, db, "unrelated-approved", "Unrelated Approved", models.EventStatusApproved)
	insertRootPageEvent(t, db, "published-approved", "Published Approved", models.EventStatusApproved)
	insertRootPageEventPulje(t, db, "published-approved", models.PuljeFredagKveld, true)

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

func TestRootPageContent_WhenProgramPublishingIsOn_RendersEventLinksWithPulje(t *testing.T) {
	// Gitt at programmet er publisert,
	// når publiserte puljearrangementer vises på forsiden,
	// så skal arrangementskortene lenke til arrangementssiden med valgt puljekontekst.

	// Given
	expectedHrefs := []string{"/event/alpha-event?pulje=FredagKveld"}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)
	insertRootPageEvent(t, db, "alpha-event", "Alpha Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "alpha-event", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualHrefs := collectRootPageHrefs(doc, ".event-card-container")

	// Then
	if !slices.Equal(expectedHrefs, actualHrefs) {
		t.Fatalf("event card hrefs mismatch\nexpected: %v\nactual:   %v", expectedHrefs, actualHrefs)
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

	insertRootPageEvent(t, db, "lordag-event", "Lordag Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "lordag-event", models.PuljeLordagMorgen, true)

	insertRootPageEvent(t, db, "fredag-event", "Fredag Event", models.EventStatusAnnounced)
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

	insertRootPageEvent(t, db, "beta-event", "Beta Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "beta-event", models.PuljeFredagKveld, true)

	insertRootPageEvent(t, db, "alpha-event", "Alpha Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "alpha-event", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualTitles := templtest.CollectTexts(doc, ".event-card-title")

	// Then
	if !slices.Equal(expectedTitles, actualTitles) {
		t.Fatalf("event titles mismatch\nexpected: %v\nactual:   %v", expectedTitles, actualTitles)
	}
}

func TestRootPageContent_WhenProgramPublishingStateCannotLoad_RendersFriendlyError(t *testing.T) {
	// Gitt at forsiden ikke kan lese publiseringsstatus,
	// når forsiden vises,
	// så skal brukeren se en vennlig feil uten tekniske detaljer.

	// Given
	expectedTextPart := rootPageLoadErrorMessage
	unexpectedTextParts := []string{
		"Error fetching",
		"program_publishing_state",
		"query program publishing state",
		"no such table",
	}

	db := createRootPageTestDB(t)
	mustExec(t, db, `DROP TABLE program_publishing_state`)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualText := rootPageText(doc)

	// Then
	assertTextContains(t, actualText, expectedTextPart)
	for _, unexpectedTextPart := range unexpectedTextParts {
		assertTextDoesNotContain(t, actualText, unexpectedTextPart)
	}
}

func TestRootPageContent_WhenEventsCannotLoad_RendersFriendlyError(t *testing.T) {
	// Gitt at forsiden ikke kan lese arrangementslisten,
	// når forsiden vises,
	// så skal brukeren se en vennlig feil uten tekniske detaljer.

	// Given
	expectedTextPart := rootEventsLoadErrorMessage
	unexpectedTextParts := []string{
		"Error fetching",
		"query announced events",
		"no such table",
	}

	db := createRootPageTestDB(t)
	setProgramPublishing(t, db, false)
	mustExec(t, db, `DROP TABLE events`)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualText := rootPageText(doc)

	// Then
	assertTextContains(t, actualText, expectedTextPart)
	for _, unexpectedTextPart := range unexpectedTextParts {
		assertTextDoesNotContain(t, actualText, unexpectedTextPart)
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
