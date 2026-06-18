package event

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEventInterestPanel_WhenScheduledWarningHasFired_RendersWarningState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at den planlagte varselmeldingen for en åpen pulje har blitt sendt.",
		When:  "Når interessepanelet rendres på nytt.",
		Then:  "Så skal billettholderen se varselstatus ved knappen.",
	})

	// Given
	expectedHelperVisible := true
	expectedHelperClass := "pulje-interest-state--warning"
	expectedMessagePart := "låses snart"
	expectedExternalLinkIconVisible := true

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
	doc := templtest.Render(t, EventInterestPanel(true, puljer, string(models.PuljeFredagKveld), true))
	helper := doc.Find(".event-interest-helper")
	actualHelperVisible := helper.Length() > 0
	actualMessage := strings.Join(strings.Fields(helper.Text()), " ")
	actualHasExpectedClass := helper.HasClass(expectedHelperClass)
	actualExternalLinkIconVisible := doc.Find(`a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon`).Length() > 0

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
	if actualExternalLinkIconVisible != expectedExternalLinkIconVisible {
		t.Fatalf("external link icon visibility mismatch\nexpected: %v\nactual:   %v", expectedExternalLinkIconVisible, actualExternalLinkIconVisible)
	}
}

func TestEventInterestPanel_WhenCurrentTimeIsBeforeWarningThreshold_RendersNoWarningState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en åpen pulje ikke nærmer seg låsing.",
		When:  "Når interessepanelet rendres.",
		Then:  "Så skal ingen låseadvarsel vises ved knappen.",
	})

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
	doc := templtest.Render(t, EventInterestPanel(true, puljer, string(models.PuljeFredagKveld), true))
	actualHelperVisible := doc.Find(".event-interest-helper").Length() > 0

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
}

func TestEventInterestPanel_WhenScheduledUrgentWarningHasFired_RendersUrgentWarningState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at den planlagte hastevarselmeldingen for en åpen pulje har blitt sendt.",
		When:  "Når interessepanelet rendres.",
		Then:  "Så skal statusen for den valgte puljen vises ved knappen.",
	})

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
	doc := templtest.Render(t, EventInterestPanel(true, puljer, string(models.PuljeLordagMorgen), true))
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

func TestEventInterestPanel_WhenDifferentPuljeHasCompletedStatus_RendersNoStatus(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en annen pulje enn den valgte er fullført.",
		When:  "Når interessepanelet rendres med valgt pulje i query-parameteren.",
		Then:  "Så skal statusen for den andre puljen ikke vises.",
	})

	// Given
	expectedHelperVisible := false

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusCompleted,
			now.Add(-1*time.Hour),
		),
		buildEventInterestTestPulje(
			models.PuljeLordagMorgen,
			"Lørdag morgen",
			models.PuljeStatusOpen,
			now.Add(4*time.Hour),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, string(models.PuljeLordagMorgen), true))
	actualHelperVisible := doc.Find(".event-interest-helper").Length() > 0

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
}

func TestEventInterestPanel_WhenSelectedPuljeIsCompleted_RendersCompletedStatus(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at den valgte puljen er fullført.",
		When:  "Når interessepanelet rendres med valgt pulje i query-parameteren.",
		Then:  "Så skal fullførtstatusen vises ved knappen.",
	})

	// Given
	expectedHelperVisible := true
	expectedHelperClass := "pulje-interest-state--completed"
	expectedMessagePart := "Puljefordelingen er klar"

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusCompleted,
			now.Add(-1*time.Hour),
		),
		buildEventInterestTestPulje(
			models.PuljeLordagMorgen,
			"Lørdag morgen",
			models.PuljeStatusOpen,
			now.Add(4*time.Hour),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, string(models.PuljeFredagKveld), true))
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

func TestEventInterestPanel_WhenPuljeQueryIsMissing_RendersNoStatus(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en pulje har status, men URL-en ikke har puljeparameter.",
		When:  "Når interessepanelet rendres.",
		Then:  "Så skal ingen puljestatus vises.",
	})

	// Given
	expectedHelperVisible := false

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusCompleted,
			now.Add(-1*time.Hour),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, "", true))
	actualHelperVisible := doc.Find(".event-interest-helper").Length() > 0

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
}

func TestEventInterestPanel_WhenPuljeQueryIsInvalid_RendersNoStatus(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en pulje har status, men URL-en har en ugyldig puljeparameter.",
		When:  "Når interessepanelet rendres.",
		Then:  "Så skal ingen puljestatus vises.",
	})

	// Given
	expectedHelperVisible := false

	now := time.Now()
	puljer := []models.PuljeRow{
		buildEventInterestTestPulje(
			models.PuljeFredagKveld,
			"Fredag kveld",
			models.PuljeStatusCompleted,
			now.Add(-1*time.Hour),
		),
	}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, "fredag_kveld", true))
	actualHelperVisible := doc.Find(".event-interest-helper").Length() > 0

	// Then
	if actualHelperVisible != expectedHelperVisible {
		t.Fatalf("helper visibility mismatch\nexpected: %v\nactual:   %v", expectedHelperVisible, actualHelperVisible)
	}
}

func TestEventInterestPanel_WhenInterestIsUnavailableForTicketHolder_RendersUnavailableMessage(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at interessevalg ikke er åpnet for arrangementet og brukeren har billett.",
		When:  "Når interessepanelet rendres.",
		Then:  "Så skal panelet vise en melding i stedet for knappen for å melde interesse.",
	})

	// Given
	expectedMessages := []string{"Interessevalg er ikke åpnet for dette arrangementet ennå."}
	expectedInterestButtonVisible := false

	puljer := []models.PuljeRow{}

	// When
	doc := templtest.Render(t, EventInterestPanel(true, puljer, "", false))
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
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at interessevalg ikke er åpnet for arrangementet og brukeren ikke har billett.",
		When:  "Når interessepanelet rendres.",
		Then:  "Så skal brukeren fortsatt se lenken for å hente billett.",
	})

	// Given
	expectedHrefs := []string{"/profile/tickets"}

	puljer := []models.PuljeRow{}

	// When
	doc := templtest.Render(t, EventInterestPanel(false, puljer, "", false))
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}
