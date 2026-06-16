package profilecomponent

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestMyEvents_WhenUserHasNoEvents_RendersCreateEventEntry(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren ikke har egne arrangementer.",
		When:  "Når Mine arrangementer vises.",
		Then:  "Så skal brukeren fortsatt se en tydelig inngang til å sende inn arrangement.",
	})

	// Given
	expectedFormAction := "/profile/api/create"
	expectedFormMethod := "post"
	expectedButtonText := []string{"Send inn arrangement"}

	// When
	doc := templtest.Render(t, MyEvents(nil))
	actualFormAction, actualFormActionExists := doc.Find(`form.submit-event-message`).Attr("action")
	actualFormMethod, actualFormMethodExists := doc.Find(`form.submit-event-message`).Attr("method")
	actualButtonText := templtest.CollectTexts(doc, `form.submit-event-message button[type="submit"]`)

	// Then
	if !actualFormActionExists || actualFormAction != expectedFormAction {
		t.Fatalf("create form action mismatch\nexpected: %q\nactual:   %q", expectedFormAction, actualFormAction)
	}
	if !actualFormMethodExists || actualFormMethod != expectedFormMethod {
		t.Fatalf("create form method mismatch\nexpected: %q\nactual:   %q", expectedFormMethod, actualFormMethod)
	}
	if !slices.Equal(expectedButtonText, actualButtonText) {
		t.Fatalf("create button text mismatch\nexpected: %v\nactual:   %v", expectedButtonText, actualButtonText)
	}
}

func TestMyEvents_RendersStatusAwareEventLinks(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren har arrangementer i ulike statuser.",
		When:  "Når Mine arrangementer vises.",
		Then:  "Så skal kladder åpne redigering og innsendte/godkjente/annonserte arrangementer åpne arrangementsvisningen.",
	})

	// Given
	expectedHrefs := []string{
		"/event/announced-event",
		"/event/approved-event",
		"/event/submitted-event",
		"/profile/new/draft-event",
	}
	events := []models.EventCardModel{
		{Id: "draft-event", Title: "Draft Event", Status: models.EventStatusDraft},
		{Id: "submitted-event", Title: "Submitted Event", Status: models.EventStatusSubmitted},
		{Id: "approved-event", Title: "Approved Event", Status: models.EventStatusApproved},
		{Id: "announced-event", Title: "Announced Event", Status: models.EventStatusAnnounced},
	}

	// When
	doc := templtest.Render(t, MyEvents(events))
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestMyEvents_WhenEventTitleIsMissing_RendersFallbackTitle(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at et av brukerens arrangementer mangler tittel.",
		When:  "Når Mine arrangementer vises.",
		Then:  "Så skal kortet vise en forståelig fallback uten å miste lenken til arrangementet.",
	})

	// Given
	expectedTitle := []string{"Mangler navn"}
	expectedHref := "/profile/new/untitled-event"
	events := []models.EventCardModel{
		{Id: "untitled-event", Status: models.EventStatusDraft},
	}

	// When
	doc := templtest.Render(t, MyEvents(events))
	actualTitle := templtest.CollectTexts(doc, ".profile-event-bar-title")
	actualHref, actualHrefExists := doc.Find(".profile-event-bar").Attr("href")

	// Then
	if !slices.Equal(expectedTitle, actualTitle) {
		t.Fatalf("fallback title mismatch\nexpected: %v\nactual:   %v", expectedTitle, actualTitle)
	}
	if !actualHrefExists || actualHref != expectedHref {
		t.Fatalf("event href mismatch\nexpected: %q\nactual:   %q", expectedHref, actualHref)
	}
}
