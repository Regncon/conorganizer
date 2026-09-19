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

	// Given
	expectedLabels := []string{"Fredag 2.10", "Lørdag 3.10", "Søndag 4.10"}
	expectedActiveDay := "Lørdag 3.10"
	expectedEventTitles := []string{"Saturday Event"}
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

	// When
	doc := templtest.Render(t, rootPageContentForDate(db, false, nil, "2026-10-03"))
	actualLabels := templtest.CollectTexts(doc, ".program-day-selector .btn")
	actualActiveDay := doc.Find(".program-day-selector .is-active").Text()
	actualEventTitles := templtest.CollectTexts(doc, ".event-card-title")

	// Then
	if !slices.Equal(expectedLabels, actualLabels) {
		t.Fatalf("program day labels mismatch\nexpected: %v\nactual:   %v", expectedLabels, actualLabels)
	}
	if actualActiveDay != expectedActiveDay {
		t.Fatalf("active program day mismatch\nexpected: %q\nactual:   %q", expectedActiveDay, actualActiveDay)
	}
	if !slices.Equal(expectedEventTitles, actualEventTitles) {
		t.Fatalf("selected day event titles mismatch\nexpected: %v\nactual:   %v", expectedEventTitles, actualEventTitles)
	}
}

func TestRootPageContent_WhenProgramPublishingIsOn_ShowsAnnouncedActivePuljeEventsRegardlessOfLegacyPublishedFlag(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at publisering av program er skrudd på.",
		When:  "Når forsiden vises.",
		Then:  "Så skal puljevisningen vise alle annonserte arrangementer som er med i en pulje.",
	})

	// Given
	expectedTitles := []string{"Published Announced", "Unpublished Announced"}

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
	expectedHrefs := []string{"/event/alpha-event?date=2026-10-09&pulje=FredagKveld"}

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

func TestRootPageContent_WhenProgramAndRaffleEventsAreMixed_DeduplicatesProgramAndRendersSections(t *testing.T) {
	db := createRootPageTestDB(t)
	seedRootPageLookups(t, db)
	setProgramPublishing(t, db, true)
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagMorgen, "Lørdag morgen", "2026-10-10T10:00:00Z", "2026-10-10T15:00:00Z")
	insertRootPagePuljeWithDetails(t, db, models.PuljeLordagKveld, "Lørdag kveld", "2026-10-10T18:00:00Z", "2026-10-10T23:00:00Z")

	for _, event := range []struct {
		id    string
		title string
	}{
		{id: "program-alpha", title: "Alpha Program"},
		{id: "program-beta", title: "Beta Program"},
		{id: "program-gamma", title: "Gamma Program"},
	} {
		insertRootPageEvent(t, db, event.id, event.title, models.EventStatusAnnounced)
		mustExec(t, db, `UPDATE events SET system = 'Hidden System', host_name = 'Hidden Host' WHERE id = ?`, event.id)
	}
	insertRootPageEventPulje(t, db, "program-alpha", models.PuljeLordagMorgen, true)
	insertRootPageEventPulje(t, db, "program-alpha", models.PuljeLordagKveld, true)
	insertRootPageEventPulje(t, db, "program-beta", models.PuljeLordagMorgen, true)
	insertRootPageEventPulje(t, db, "program-gamma", models.PuljeLordagKveld, true)
	for _, eventID := range []string{"program-alpha", "program-beta", "program-gamma"} {
		setRootPageEventInPuljefordeling(t, db, eventID, false)
	}

	insertRootPageEvent(t, db, "raffle-morning", "Morning Raffle", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "raffle-morning", models.PuljeLordagMorgen, true)
	insertRootPageEvent(t, db, "raffle-evening", "Evening Raffle", models.EventStatusAnnounced)
	insertRootPageEventPulje(t, db, "raffle-evening", models.PuljeLordagKveld, true)

	doc := templtest.Render(t, rootPageContentForDate(db, false, nil, "2026-10-10"))
	programTitles := templtest.CollectTexts(doc, ".program-event-card .event-card-title")
	if !slices.Equal(programTitles, []string{"Alpha Program", "Beta Program", "Gamma Program"}) {
		t.Fatalf("program event titles mismatch: %v", programTitles)
	}
	if got := doc.Find(".program-event-card").Length(); got != 3 {
		t.Fatalf("program event card count = %d, want 3", got)
	}
	if got := doc.Find(".program-event-card .event-card-subtitle, .program-event-card .event-card-footer-gamemaster").Length(); got != 0 {
		t.Fatalf("program cards rendered %d system/GM elements", got)
	}
	if got := doc.Find(".program-event-card .event-card-main-body .event-card-description").Length(); got != 0 {
		t.Fatalf("program cards rendered %d descriptions in the body", got)
	}
	if got := doc.Find(".program-event-card .event-card-footer-description").Length(); got != 3 {
		t.Fatalf("program cards rendered %d descriptions in the footer, want 3", got)
	}
	if got := doc.Find(".raffle-eventcard-grid .raffle-event-card").Length(); got != 2 {
		t.Fatalf("raffle event card count = %d, want 2", got)
	}
	if got := doc.Find(".program-event-row").Length(); got != 2 {
		t.Fatalf("program row count = %d, want 2", got)
	}
	if got := doc.Find(".program-event-row.reversed").Length(); got != 1 {
		t.Fatalf("reversed program row count = %d, want 1", got)
	}

	alphaHref := doc.Find(`.program-event-card`).First().AttrOr("href", "")
	if alphaHref != "/event/program-alpha?date=2026-10-10&pulje=LordagMorgen" {
		t.Fatalf("program event href = %q", alphaHref)
	}
	if headings := templtest.CollectTexts(doc, ".pulje-heading"); !slices.Equal(headings, []string{"Programoversikt", "Lørdag morgen (10:00 - 15:00)", "Lørdag kveld (18:00 - 23:00)"}) {
		t.Fatalf("section headings mismatch: %v", headings)
	}
}
