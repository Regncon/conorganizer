package root

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
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
