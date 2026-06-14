package event

import (
	"database/sql"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEventPageContent_WhenProgramPublishingIsOff_RendersPreviousNextForMiddleAnnouncedEvent(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en midtre annonsert arrangementsside før programmet er publisert.",
		When:  "Når siden rendres.",
		Then:  "Så vises forrige/neste med eksisterende etiketter og lenker uten pulje.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt første og siste annonserte arrangement før programmet er publisert.",
		When:  "Når sidene rendres.",
		Then:  "Så rendres ikke manglende forrige/neste-side som deaktivert knapp.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt samme arrangement publisert i to puljer.",
		When:  "Når arrangementssiden rendres med ulike puljeverdier.",
		Then:  "Så følger forrige/neste-lenkene den konkrete forsidelisteforekomsten.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt et publisert arrangement i en pulje.",
		When:  "Når siden rendres uten eller med feil puljeverdi.",
		Then:  "Så rendres ikke forrige/neste-komponentene.",
	})

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
