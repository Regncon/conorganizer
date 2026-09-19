package root

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestRootPageContent_WhenProgramPublishingIsOn_ShowsProgramDaySelector(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at publisering av program er skrudd på.",
		When:  "Når forsiden vises.",
		Then:  "Så skal dagvelgeren vises.",
	})

	// Given
	expectedDaySelectorVisible := true

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePulje(t, db)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualDaySelectorVisible := templtest.HasSelector(doc, ".program-day-selector-container")

	// Then
	if actualDaySelectorVisible != expectedDaySelectorVisible {
		t.Fatalf("day selector visibility mismatch\nexpected: %v\nactual:   %v", expectedDaySelectorVisible, actualDaySelectorVisible)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOn_ShowsProgramDaysWithActiveDay(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programmet er publisert og puljene dekker fredag, lørdag og søndag.",
		When:  "Når forsiden vises med lørdag valgt.",
		Then:  "Så skal dagene vises med riktig aktiv dag.",
	})

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePuljeWithDetails(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-10-02T18:00:00Z", "2026-10-02T23:00:00Z")
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-03T10:00:00Z", "2026-10-03T15:00:00Z")
	insertRootPagePuljeWithDetails(t, db, models.PuljeSondagMorgen, "Søndag morgen", "2026-10-04T10:00:00Z", "2026-10-04T15:00:00Z")
	insertRootPageEvent(t, db, "friday-event", "Friday Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "friday-event", models.PuljeFredagKveld, true)
	insertRootPageEvent(t, db, "saturday-event", "Saturday Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "saturday-event", models.PuljeLordagMorgen, true)
	insertRootPageEvent(t, db, "sunday-event", "Sunday Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "sunday-event", models.PuljeSondagMorgen, true)

	doc := templtest.Render(t, rootPageContentForDate(db, false, nil, "2026-10-03"))

	if actual := templtest.CollectTexts(doc, ".program-day-selector .btn"); !slices.Equal(
		[]string{"Fredag 2.10", "Lørdag 3.10", "Søndag 4.10"}, actual,
	) {
		t.Fatalf("program day labels mismatch\nexpected: %v\nactual:   %v", []string{"Fredag 2.10", "Lørdag 3.10", "Søndag 4.10"}, actual)
	}

	if actual := doc.Find(".program-day-selector .is-active").Text(); actual != "Lørdag 3.10" {
		t.Fatalf("active program day mismatch\nexpected: %q\nactual:   %q", "Lørdag 3.10", actual)
	}

	if actual := templtest.CollectTexts(doc, ".event-card-title"); !slices.Equal([]string{"Saturday Event"}, actual) {
		t.Fatalf("selected day event titles mismatch\nexpected: %v\nactual:   %v", []string{"Saturday Event"}, actual)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOn_OnlyShowsAnnouncedPublishedPuljeEvents(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at publisering av program er skrudd på.",
		When:  "Når forsiden vises.",
		Then:  "Så skal puljevisningen bare vise annonserte arrangementer som er publisert i en pulje.",
	})

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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at programmet er publisert.",
		When:  "Når publiserte puljearrangementer vises på forsiden.",
		Then:  "Så skal arrangementskortene lenke til arrangementssiden med valgt puljekontekst.",
	})

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

func TestRootPageContent_WhenProgramPublishingIsOn_RendersSelectedDatePuljeSectionsInTimeOrder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at publisering av program er skrudd på.",
		When:  "Når forsiden vises.",
		Then:  "Så skal den valgte dagen vise puljene sortert etter starttidspunkt.",
	})

	// Given
	expectedPuljeHeadings := []string{
		"Fredag kveld (18:00 - 23:00)",
		"Lordag morgen (20:00 - 22:00)",
	}

	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePuljeWithDetails(t, db, models.PuljeFredagKveld, "Fredag kveld", "2026-10-09T18:00:00Z", "2026-10-09T23:00:00Z")
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagMorgen, "Lordag morgen", "2026-10-09T20:00:00Z", "2026-10-09T22:00:00Z")

	insertRootPageEvent(t, db, "lordag-event", "Lordag Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "lordag-event", models.PuljeLordagMorgen, true)

	insertRootPageEvent(t, db, "fredag-event", "Fredag Event", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "fredag-event", models.PuljeFredagKveld, true)

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, false, nil, "2026-10-09"))
	actualPuljeHeadings := templtest.CollectTexts(doc, ".pulje-heading")

	// Then
	if !slices.Equal(expectedPuljeHeadings, actualPuljeHeadings) {
		t.Fatalf("pulje headings mismatch\nexpected: %v\nactual:   %v", expectedPuljeHeadings, actualPuljeHeadings)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOn_SortsEventsAlphabeticallyWithinPulje(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at publisering av program er skrudd på.",
		When:  "Når forsiden vises.",
		Then:  "Så skal arrangementene sorteres alfabetisk innenfor hver pulje.",
	})

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
