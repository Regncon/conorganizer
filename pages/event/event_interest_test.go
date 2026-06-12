package event

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEventInterestPanel_WhenScheduledWarningHasFired_RendersWarningState(t *testing.T) {
	// Gitt at den planlagte varselmeldingen for en åpen pulje har blitt sendt,
	// når interessepanelet rendres på nytt,
	// så skal billettholderen se varselstatus ved knappen.

	// Given
	expectedHelperVisible := true
	expectedHelperClass := "pulje-interest-state--warning"
	expectedMessagePart := "låses snart"

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusOpen,
			now.Add(2*time.Hour),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, true))
	helper := doc.Find(".event-interest-helper")
	actualHelperVisible := helper.Length() > 0
	actualMessage := strings.Join(strings.Fields(helper.Text()), " ")
	actualHasExpectedClass := helper.HasClass(expectedHelperClass)

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
	if !actualHasExpectedClass {
		t.Fatalf("helper class mismatch\nexpected helper to have class: %s", expectedHelperClass)
	}
	if !strings.Contains(actualMessage, expectedMessagePart) {
		t.Fatalf("helper message mismatch\nexpected to contain: %q\nactual:              %q", expectedMessagePart, actualMessage)
	}
}

func TestEventInterestPanel_WhenCurrentTimeIsBeforeWarningThreshold_RendersNoWarningState(t *testing.T) {
	// Gitt at en åpen pulje ikke nærmer seg låsing,
	// når interessepanelet rendres,
	// så skal ingen låseadvarsel vises ved knappen.

	// Given
	expectedHelperVisible := false

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusOpen,
			now.Add(4*time.Hour),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, true))
	actualHelperVisible := doc.Find(".event-interest-helper").Length() > 0

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
}

func TestEventInterestPanel_WhenScheduledUrgentWarningHasFired_RendersUrgentWarningState(t *testing.T) {
	// Gitt at den planlagte hastevarselmeldingen for en åpen pulje har blitt sendt,
	// når interessepanelet rendres,
	// så skal den mest presserende puljemeldingen vises ved knappen.

	// Given
	expectedHelperVisible := true
	expectedHelperClass := "pulje-interest-state--urgent-warning"
	expectedMessagePart := "låses straks"

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusLocked,
			now.Add(-1*time.Hour),
		),
		buildEventInterestTestPulje(
			models.PuljeLordagMorgen,
			"Lørdag morgen",
			models.PuljeStatusOpen,
			now.Add(45*time.Minute),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, true))
	helper := doc.Find(".event-interest-helper")
	actualHelperVisible := helper.Length() > 0
	actualMessage := strings.Join(strings.Fields(helper.Text()), " ")
	actualHasExpectedClass := helper.HasClass(expectedHelperClass)

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
	if !actualHasExpectedClass {
		t.Fatalf("helper class mismatch\nexpected helper to have class: %s", expectedHelperClass)
	}
	if !strings.Contains(actualMessage, expectedMessagePart) {
		t.Fatalf("helper message mismatch\nexpected to contain: %q\nactual:              %q", expectedMessagePart, actualMessage)
	}
}

func TestEventInterestPanel_WhenInterestIsUnavailableForTicketHolder_RendersUnavailableMessage(t *testing.T) {
	// Gitt at interessevalg ikke er åpnet for arrangementet og brukeren har billett,
	// når interessepanelet rendres,
	// så skal panelet vise en melding i stedet for knappen for å melde interesse.

	// Given
	expectedMessages := []string{"Interessevalg er ikke åpnet for dette arrangementet ennå."}
	expectedInterestButtonVisible := false

	puljer := []models.PuljeRow{}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, false))
	actualMessages := templtest.CollectTexts(doc, ".event-interest-unavailable-message")
	actualInterestButtonVisible := doc.Find(".event-interest-open-button").Length() > 0

	// Then
	if !slices.Equal(expectedMessages, actualMessages) {
		t.Fatalf("unavailable message mismatch\nexpected: %v\nactual:   %v", expectedMessages, actualMessages)
	}
	if actualInterestButtonVisible != expectedInterestButtonVisible {
		t.Fatalf("interest button visibility mismatch\nexpected: %v\nactual:   %v", expectedInterestButtonVisible, actualInterestButtonVisible)
	}
}

func TestEventInterestPanel_WhenInterestIsUnavailableAndUserHasNoTicket_RendersTicketCTA(t *testing.T) {
	// Gitt at interessevalg ikke er åpnet for arrangementet og brukeren ikke har billett,
	// når interessepanelet rendres,
	// så skal brukeren fortsatt se lenken for å hente billett.

	// Given
	expectedHrefs := []string{"/profile/tickets"}

	puljer := []models.PuljeRow{}

	// When
	doc := templtest.Render(t, EventInterestPanel(false, puljer, false))
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}
