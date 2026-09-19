package eventservice

import (
	"context"
	"net/http/httptest"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestGetPreviousNextForRootEventList_WhenProgramIsNotPublished_UsesAnnouncedAlphabeticalRootList(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt noen annonserte og interne arrangementer.",
		When:  "Når forrige/neste hentes før programmet er publisert.",
		Then:  "Så brukes den flate annonserte forsidelisten i alfabetisk rekkefølge.",
	})

	// Given
	expectedAnnouncedOrder := []string{"alpha-announced", "beta-announced", "delta-announced"}
	expectedMiddle := expectedPreviousNext{
		previousURL:   "/event/alpha-announced",
		previousTitle: "Alpha Announced",
		nextURL:       "/event/delta-announced",
		nextTitle:     "Delta Announced",
	}
	expectedFirst := expectedPreviousNext{
		nextURL:   "/event/beta-announced",
		nextTitle: "Beta Announced",
	}
	expectedLast := expectedPreviousNext{
		previousURL:   "/event/beta-announced",
		previousTitle: "Beta Announced",
	}

	ctx := context.Background()
	imgDir := ""
	db := createPreviousNextRootListTestDB(t)
	seedPreviousNextRootListLookups(t, db)
	seedPreviousNextRootListEvent(t, db, "draft-event", "Draft Event", models.EventStatusDraft)
	seedPreviousNextRootListEvent(t, db, "submitted-event", "Submitted Event", models.EventStatusSubmitted)
	seedPreviousNextRootListEvent(t, db, "approved-event", "Approved Event", models.EventStatusApproved)
	seedPreviousNextRootListEvent(t, db, "archived-event", "Archived Event", models.EventStatusArchived)
	seedPreviousNextRootListEvent(t, db, "beta-announced", "Beta Announced", models.EventStatusAnnounced)
	seedPreviousNextRootListEvent(t, db, "delta-announced", "Delta Announced", models.EventStatusAnnounced)
	seedPreviousNextRootListEvent(t, db, "alpha-announced", "Alpha Announced", models.EventStatusAnnounced)

	request := httptest.NewRequest("GET", "/event/beta-announced?pulje=FredagKveld", nil)

	// When
	announcedEvents, err := root.GetAnnouncedEventsAlphabetically(db)
	if err != nil {
		t.Fatalf("expected announced root events query to succeed: %v", err)
	}
	actualAnnouncedOrder := collectPreviousNextRootListEventIDs(announcedEvents)

	middle, err := GetPreviousNextForRootEventList(ctx, db, "beta-announced", false, request, &imgDir)
	if err != nil {
		t.Fatalf("expected middle previous/next query to succeed: %v", err)
	}
	first, err := GetPreviousNextForRootEventList(ctx, db, "alpha-announced", false, request, &imgDir)
	if err != nil {
		t.Fatalf("expected first previous/next query to succeed: %v", err)
	}
	last, err := GetPreviousNextForRootEventList(ctx, db, "delta-announced", false, request, &imgDir)
	if err != nil {
		t.Fatalf("expected last previous/next query to succeed: %v", err)
	}
	approved, err := GetPreviousNextForRootEventList(ctx, db, "approved-event", false, request, &imgDir)
	if err != nil {
		t.Fatalf("expected approved previous/next query to succeed: %v", err)
	}
	submitted, err := GetPreviousNextForRootEventList(ctx, db, "submitted-event", false, request, &imgDir)
	if err != nil {
		t.Fatalf("expected submitted previous/next query to succeed: %v", err)
	}
	draft, err := GetPreviousNextForRootEventList(ctx, db, "draft-event", false, request, &imgDir)
	if err != nil {
		t.Fatalf("expected draft previous/next query to succeed: %v", err)
	}
	archived, err := GetPreviousNextForRootEventList(ctx, db, "archived-event", false, request, &imgDir)
	if err != nil {
		t.Fatalf("expected archived previous/next query to succeed: %v", err)
	}

	// Then
	if !slices.Equal(expectedAnnouncedOrder, actualAnnouncedOrder) {
		t.Fatalf("announced root event order mismatch\nexpected: %v\nactual:   %v", expectedAnnouncedOrder, actualAnnouncedOrder)
	}
	assertPreviousNextMatches(t, expectedMiddle, middle)
	assertPreviousNextMatches(t, expectedFirst, first)
	assertPreviousNextMatches(t, expectedLast, last)
	assertPreviousNextMatches(t, expectedPreviousNext{}, approved)
	assertPreviousNextMatches(t, expectedPreviousNext{}, submitted)
	assertPreviousNextMatches(t, expectedPreviousNext{}, draft)
	assertPreviousNextMatches(t, expectedPreviousNext{}, archived)
}

func TestGetPreviousNextForRootEventList_WhenProgramIsPublished_IgnoresLegacyPuljePublishedFlag(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt publiserte, upubliserte og interne puljerader.",
		When:  "Når forrige/neste hentes etter at programmet er publisert.",
		Then:  "Så brukes annonserte rader som er med i puljen, uavhengig av det gamle publiseringsflagget.",
	})

	// Given
	ctx := context.Background()
	imgDir := ""
	db := createPreviousNextRootListTestDB(t)
	seedPreviousNextRootListLookups(t, db)
	seedPreviousNextRootListPulje(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-10-09T18:00:00Z", "2026-10-09T23:00:00Z")
	seedPreviousNextRootListPulje(t, db, models.PuljeLordagMorgen, "Lordag morgen", "2026-10-10T10:00:00Z", "2026-10-10T14:00:00Z")
	seedPreviousNextRootListPulje(t, db, models.PuljeSondagMorgen, "Sondag morgen", "2026-10-11T10:00:00Z", "2026-10-11T14:00:00Z")

	seedPreviousNextRootListEvent(t, db, "alpha-fredag", "Alpha Fredag", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "alpha-fredag", models.PuljeFredagKveld, true, true)
	seedPreviousNextRootListEvent(t, db, "shared-event", "Shared Event", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "shared-event", models.PuljeFredagKveld, true, true)
	seedPreviousNextRootListEventPulje(t, db, "shared-event", models.PuljeLordagMorgen, true, true)
	seedPreviousNextRootListEvent(t, db, "zeta-fredag", "Zeta Fredag", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "zeta-fredag", models.PuljeFredagKveld, true, true)
	seedPreviousNextRootListEvent(t, db, "lima-lordag", "Lima Lordag", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "lima-lordag", models.PuljeLordagMorgen, true, true)
	seedPreviousNextRootListEvent(t, db, "zulu-lordag", "Zulu Lordag", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "zulu-lordag", models.PuljeLordagMorgen, true, true)

	seedPreviousNextRootListEvent(t, db, "not-in-pulje", "Beta Not In Pulje", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "not-in-pulje", models.PuljeFredagKveld, false, true)
	seedPreviousNextRootListEvent(t, db, "unpublished-pulje", "Beta Unpublished", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "unpublished-pulje", models.PuljeFredagKveld, true, false)
	seedPreviousNextRootListEvent(t, db, "approved-pulje", "Beta Approved", models.EventStatusApproved)
	seedPreviousNextRootListEventPulje(t, db, "approved-pulje", models.PuljeFredagKveld, true, true)

	cases := []struct {
		name      string
		currentID string
		path      string
		expected  expectedPreviousNext
	}{
		{
			name:      "same event in fredag occurrence uses fredag neighbors",
			currentID: "shared-event",
			path:      "/event/shared-event?pulje=FredagKveld",
			expected: expectedPreviousNext{
				previousURL:   "/event/unpublished-pulje?date=2026-10-09&pulje=FredagKveld",
				previousTitle: "Beta Unpublished",
				nextURL:       "/event/zeta-fredag?date=2026-10-09&pulje=FredagKveld",
				nextTitle:     "Zeta Fredag",
			},
		},
		{
			name:      "same event in lordag occurrence uses lordag neighbors",
			currentID: "shared-event",
			path:      "/event/shared-event?pulje=LordagMorgen",
			expected: expectedPreviousNext{
				previousURL:   "/event/lima-lordag?date=2026-10-10&pulje=LordagMorgen",
				previousTitle: "Lima Lordag",
				nextURL:       "/event/zulu-lordag?date=2026-10-10&pulje=LordagMorgen",
				nextTitle:     "Zulu Lordag",
			},
		},
		{
			name:      "is_in_pulje zero row is excluded",
			currentID: "not-in-pulje",
			path:      "/event/not-in-pulje?pulje=FredagKveld",
			expected:  expectedPreviousNext{},
		},
		{
			name:      "legacy unpublished pulje row remains included",
			currentID: "unpublished-pulje",
			path:      "/event/unpublished-pulje?pulje=FredagKveld",
			expected: expectedPreviousNext{
				previousURL:   "/event/alpha-fredag?date=2026-10-09&pulje=FredagKveld",
				previousTitle: "Alpha Fredag",
				nextURL:       "/event/shared-event?date=2026-10-09&pulje=FredagKveld",
				nextTitle:     "Shared Event",
			},
		},
		{
			name:      "non announced event row is excluded",
			currentID: "approved-pulje",
			path:      "/event/approved-pulje?pulje=FredagKveld",
			expected:  expectedPreviousNext{},
		},
		{
			name:      "missing pulje query returns no navigation",
			currentID: "shared-event",
			path:      "/event/shared-event",
			expected:  expectedPreviousNext{},
		},
		{
			name:      "invalid pulje query returns no navigation",
			currentID: "shared-event",
			path:      "/event/shared-event?pulje=Nope",
			expected:  expectedPreviousNext{},
		},
		{
			name:      "valid but wrong pulje query returns no navigation",
			currentID: "shared-event",
			path:      "/event/shared-event?pulje=SondagMorgen",
			expected:  expectedPreviousNext{},
		},
		{
			name:      "first root occurrence has no previous",
			currentID: "alpha-fredag",
			path:      "/event/alpha-fredag?pulje=FredagKveld",
			expected: expectedPreviousNext{
				nextURL:   "/event/unpublished-pulje?date=2026-10-09&pulje=FredagKveld",
				nextTitle: "Beta Unpublished",
			},
		},
		{
			name:      "last root occurrence has no next",
			currentID: "zulu-lordag",
			path:      "/event/zulu-lordag?pulje=LordagMorgen",
			expected: expectedPreviousNext{
				previousURL:   "/event/shared-event?date=2026-10-10&pulje=LordagMorgen",
				previousTitle: "Shared Event",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// When
			request := httptest.NewRequest("GET", tc.path, nil)
			actual, err := GetPreviousNextForRootEventList(ctx, db, tc.currentID, true, request, &imgDir)

			// Then
			if err != nil {
				t.Fatalf("expected previous/next query to succeed: %v", err)
			}
			assertPreviousNextMatches(t, tc.expected, actual)
		})
	}
}

func TestGetPreviousNextForRootEventList_WhenProgramAndRaffleEventsShareADay_UsesRenderedOrder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en dag har programarrangementer og rafflearrangementer i flere puljer.",
		When:  "Når forrige og neste arrangement hentes for hvert arrangement.",
		Then:  "Så skal navigasjonen følge den rendrerte rekkefølgen uten å gå rundt til starten.",
	})

	// Given
	expectedCases := []struct {
		id    string
		pulje models.Pulje
		want  expectedPreviousNext
	}{
		{
			id: "program-alpha", pulje: models.PuljeLordagMorgen,
			want: expectedPreviousNext{nextURL: "/event/program-beta?date=2026-10-10&pulje=LordagMorgen", nextTitle: "Beta Program"},
		},
		{
			id: "program-beta", pulje: models.PuljeLordagMorgen,
			want: expectedPreviousNext{
				previousURL: "/event/program-alpha?date=2026-10-10&pulje=LordagMorgen", previousTitle: "Alpha Program",
				nextURL: "/event/morning-raffle?date=2026-10-10&pulje=LordagMorgen", nextTitle: "Morning Raffle",
			},
		},
		{
			id: "morning-raffle", pulje: models.PuljeLordagMorgen,
			want: expectedPreviousNext{
				previousURL: "/event/program-beta?date=2026-10-10&pulje=LordagMorgen", previousTitle: "Beta Program",
				nextURL: "/event/evening-raffle?date=2026-10-10&pulje=LordagKveld", nextTitle: "Evening Raffle",
			},
		},
		{
			id: "evening-raffle", pulje: models.PuljeLordagKveld,
			want: expectedPreviousNext{previousURL: "/event/morning-raffle?date=2026-10-10&pulje=LordagMorgen", previousTitle: "Morning Raffle"},
		},
	}

	db := createPreviousNextRootListTestDB(t)
	seedPreviousNextRootListLookups(t, db)
	seedPreviousNextRootListPulje(t, db, models.PuljeLordagMorgen, "Lordag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	seedPreviousNextRootListPulje(t, db, models.PuljeLordagKveld, "Lordag kveld", "2026-10-10T18:00:00Z", "2026-10-10T23:00:00Z")

	seedPreviousNextRootListEvent(t, db, "program-alpha", "Alpha Program", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "program-alpha", models.PuljeLordagMorgen, true, true)
	seedPreviousNextRootListEventPulje(t, db, "program-alpha", models.PuljeLordagKveld, true, true)
	setPreviousNextRootListEventInPuljefordeling(t, db, "program-alpha", false)
	seedPreviousNextRootListEvent(t, db, "program-beta", "Beta Program", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "program-beta", models.PuljeLordagMorgen, true, true)
	setPreviousNextRootListEventInPuljefordeling(t, db, "program-beta", false)
	seedPreviousNextRootListEvent(t, db, "morning-raffle", "Morning Raffle", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "morning-raffle", models.PuljeLordagMorgen, true, true)
	seedPreviousNextRootListEvent(t, db, "evening-raffle", "Evening Raffle", models.EventStatusAnnounced)
	seedPreviousNextRootListEventPulje(t, db, "evening-raffle", models.PuljeLordagKveld, true, true)

	// When
	for _, tc := range expectedCases {
		t.Run(tc.id, func(t *testing.T) {
			// When
			request := httptest.NewRequest("GET", "/event/"+tc.id+"?date=2026-10-10&pulje="+string(tc.pulje), nil)
			actual, err := GetPreviousNextForRootEventList(context.Background(), db, tc.id, true, request, nil)

			// Then
			if err != nil {
				t.Fatalf("get previous/next: %v", err)
			}
			assertPreviousNextMatches(t, tc.want, actual)
		})
	}
}
