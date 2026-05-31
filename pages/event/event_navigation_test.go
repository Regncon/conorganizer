package event

import (
	"database/sql"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEventPageContent_WhenProgramPublishingIsOff_RendersPreviousNextForMiddleAnnouncedEvent(t *testing.T) {
	// Gitt en midtre annonsert arrangementsside før programmet er publisert,
	// når siden rendres,
	// så vises forrige/neste med eksisterende etiketter og lenker uten pulje.

	// Given
	expectedHrefs := []string{"/event/alpha-announced", "/event/delta-announced"}
	expectedText := []string{"Forrige arrangement", "Neste arrangement"}

	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	setEventVisibilityProgramPublishing(t, db, false)
	seedEventVisibilityEvent(t, db, "alpha-announced", "Alpha Announced", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityEvent(t, db, "beta-announced", "Beta Announced", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityEvent(t, db, "delta-announced", "Delta Announced", models.EventStatusAnnounced, sql.NullInt64{})
	request := httptest.NewRequest("GET", "/event/beta-announced", nil)

	// When
	doc := templtest.Render(t, event_page_content("beta-announced", false, logger, db, nil, request))
	actualHrefs := collectEventNavigationHrefs(doc)
	actualText := collectEventNavigationText(doc)

	// Then
	if !slices.Equal(expectedHrefs, actualHrefs) {
		t.Fatalf("event navigation hrefs mismatch\nexpected: %v\nactual:   %v", expectedHrefs, actualHrefs)
	}
	for _, text := range expectedText {
		if !strings.Contains(actualText, text) {
			t.Fatalf("expected event navigation text to contain %q, got %q", text, actualText)
		}
	}
	assertNoDisabledEventNavigationPlaceholders(t, doc)
}

func TestEventPageContent_WhenProgramPublishingIsOff_DoesNotRenderMissingEdgeButtons(t *testing.T) {
	// Gitt første og siste annonserte arrangement før programmet er publisert,
	// når sidene rendres,
	// så rendres ikke manglende forrige/neste-side som deaktivert knapp.

	// Given
	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	setEventVisibilityProgramPublishing(t, db, false)
	seedEventVisibilityEvent(t, db, "alpha-announced", "Alpha Announced", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityEvent(t, db, "beta-announced", "Beta Announced", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventVisibilityEvent(t, db, "delta-announced", "Delta Announced", models.EventStatusAnnounced, sql.NullInt64{})

	cases := []struct {
		name          string
		eventID       string
		expectedHrefs []string
	}{
		{
			name:          "first has next only",
			eventID:       "alpha-announced",
			expectedHrefs: []string{"/event/beta-announced"},
		},
		{
			name:          "last has previous only",
			eventID:       "delta-announced",
			expectedHrefs: []string{"/event/beta-announced"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// When
			request := httptest.NewRequest("GET", "/event/"+tc.eventID, nil)
			doc := templtest.Render(t, event_page_content(tc.eventID, false, logger, db, nil, request))
			actualHrefs := collectEventNavigationHrefs(doc)

			// Then
			if !slices.Equal(tc.expectedHrefs, actualHrefs) {
				t.Fatalf("event navigation hrefs mismatch\nexpected: %v\nactual:   %v", tc.expectedHrefs, actualHrefs)
			}
			assertNoDisabledEventNavigationPlaceholders(t, doc)
		})
	}
}

func TestEventPageContent_WhenProgramPublishingIsOn_RendersPuljeSpecificPreviousNextLinks(t *testing.T) {
	// Gitt samme arrangement publisert i to puljer,
	// når arrangementssiden rendres med ulike puljeverdier,
	// så følger forrige/neste-lenkene den konkrete forsidelisteforekomsten.

	// Given
	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	setEventVisibilityProgramPublishing(t, db, true)
	seedEventNavigationPulje(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-10-09T18:00:00Z", "2026-10-09T23:00:00Z")
	seedEventNavigationPulje(t, db, models.PuljeLordagMorgen, "Lordag morgen", "2026-10-10T10:00:00Z", "2026-10-10T14:00:00Z")
	seedEventVisibilityEvent(t, db, "alpha-fredag", "Alpha Fredag", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventNavigationEventPulje(t, db, "alpha-fredag", models.PuljeFredagKveld)
	seedEventVisibilityEvent(t, db, "shared-event", "Shared Event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventNavigationEventPulje(t, db, "shared-event", models.PuljeFredagKveld)
	seedEventNavigationEventPulje(t, db, "shared-event", models.PuljeLordagMorgen)
	seedEventVisibilityEvent(t, db, "zeta-fredag", "Zeta Fredag", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventNavigationEventPulje(t, db, "zeta-fredag", models.PuljeFredagKveld)
	seedEventVisibilityEvent(t, db, "lima-lordag", "Lima Lordag", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventNavigationEventPulje(t, db, "lima-lordag", models.PuljeLordagMorgen)
	seedEventVisibilityEvent(t, db, "zulu-lordag", "Zulu Lordag", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventNavigationEventPulje(t, db, "zulu-lordag", models.PuljeLordagMorgen)

	cases := []struct {
		name          string
		path          string
		expectedHrefs []string
	}{
		{
			name: "fredag occurrence",
			path: "/event/shared-event?pulje=FredagKveld",
			expectedHrefs: []string{
				"/event/alpha-fredag?pulje=FredagKveld",
				"/event/zeta-fredag?pulje=FredagKveld",
			},
		},
		{
			name: "lordag occurrence",
			path: "/event/shared-event?pulje=LordagMorgen",
			expectedHrefs: []string{
				"/event/lima-lordag?pulje=LordagMorgen",
				"/event/zulu-lordag?pulje=LordagMorgen",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// When
			request := httptest.NewRequest("GET", tc.path, nil)
			doc := templtest.Render(t, event_page_content("shared-event", false, logger, db, nil, request))
			actualHrefs := collectEventNavigationHrefs(doc)

			// Then
			if !slices.Equal(tc.expectedHrefs, actualHrefs) {
				t.Fatalf("event navigation hrefs mismatch\nexpected: %v\nactual:   %v", tc.expectedHrefs, actualHrefs)
			}
			assertNoDisabledEventNavigationPlaceholders(t, doc)
		})
	}
}

func TestEventPageContent_WhenProgramPublishingIsOn_DoesNotRenderNavigationForMissingOrWrongPulje(t *testing.T) {
	// Gitt et publisert arrangement i en pulje,
	// når siden rendres uten eller med feil puljeverdi,
	// så rendres ikke forrige/neste-komponentene.

	// Given
	db := createEventVisibilityTestDB(t)
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	setEventVisibilityProgramPublishing(t, db, true)
	seedEventNavigationPulje(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-10-09T18:00:00Z", "2026-10-09T23:00:00Z")
	seedEventVisibilityEvent(t, db, "alpha-fredag", "Alpha Fredag", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventNavigationEventPulje(t, db, "alpha-fredag", models.PuljeFredagKveld)
	seedEventVisibilityEvent(t, db, "shared-event", "Shared Event", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventNavigationEventPulje(t, db, "shared-event", models.PuljeFredagKveld)
	seedEventVisibilityEvent(t, db, "zeta-fredag", "Zeta Fredag", models.EventStatusAnnounced, sql.NullInt64{})
	seedEventNavigationEventPulje(t, db, "zeta-fredag", models.PuljeFredagKveld)

	cases := []struct {
		name string
		path string
	}{
		{name: "missing pulje", path: "/event/shared-event"},
		{name: "invalid pulje", path: "/event/shared-event?pulje=Nope"},
		{name: "wrong pulje", path: "/event/shared-event?pulje=LordagMorgen"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// When
			request := httptest.NewRequest("GET", tc.path, nil)
			doc := templtest.Render(t, event_page_content("shared-event", false, logger, db, nil, request))

			// Then
			assertEventNavigationNotRendered(t, doc)
		})
	}
}

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
