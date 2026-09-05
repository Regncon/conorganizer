package ticketholder

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestTicketHolderInterestPicker_WhenEventHasOpenRegistration_ReplacesHighInterestWithRegistration(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt eit arrangement med open påmelding.",
		When:  "Når interesseveljaren blir vist.",
		Then:  "Så skal direkte på- og avmelding erstatte valet veldig interessert.",
	})

	// Given
	expectedRegistrationButtons := 1
	expectedHighInterestButtons := 0
	expectedDeregistrationButtons := 1
	expectedEndpoint := "/event/api/open-event/registration"

	// When
	doc := templtest.Render(t, TicketHolderInterestPicker("open-event", true))
	actualRegistrationButtons := doc.Find(".interest-register").Length()
	actualHighInterestButtons := doc.Find(".interest-high").Length()
	actualDeregistrationButtons := doc.Find(".interest-deregister").Length()
	actualRegistrationAction := doc.Find(".interest-register").AttrOr("data-on:click", "")
	actualDeregistrationAction := doc.Find(".interest-deregister").AttrOr("data-on:click", "")

	// Then
	if actualRegistrationButtons != expectedRegistrationButtons {
		t.Fatalf("registration button count mismatch\nexpected: %d\nactual:   %d", expectedRegistrationButtons, actualRegistrationButtons)
	}
	if actualHighInterestButtons != expectedHighInterestButtons {
		t.Fatalf("high-interest button count mismatch\nexpected: %d\nactual:   %d", expectedHighInterestButtons, actualHighInterestButtons)
	}
	if actualDeregistrationButtons != expectedDeregistrationButtons {
		t.Fatalf("deregistration button count mismatch\nexpected: %d\nactual:   %d", expectedDeregistrationButtons, actualDeregistrationButtons)
	}
	if !strings.Contains(actualRegistrationAction, expectedEndpoint) || !strings.Contains(actualRegistrationAction, "$isRegistered = true") {
		t.Fatalf("registration action mismatch: %q", actualRegistrationAction)
	}
	if !strings.Contains(actualDeregistrationAction, expectedEndpoint) || !strings.Contains(actualDeregistrationAction, "$isRegistered = false") {
		t.Fatalf("deregistration action mismatch: %q", actualDeregistrationAction)
	}
}

func TestTicketHolderInterestPicker_WhenEventUsesOrdinarySelection_KeepsHighInterestAndAllowsManualAssignmentOptOut(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt eit vanleg arrangement utan open påmelding.",
		When:  "Når interesseveljaren blir vist.",
		Then:  "Så skal veldig interessert vere tilgjengeleg og manuell tildeling kunne meldast av.",
	})

	// Given
	expectedRegistrationButtons := 0
	expectedHighInterestButtons := 1
	expectedDeregistrationButtons := 1

	// When
	doc := templtest.Render(t, TicketHolderInterestPicker("ordinary-event", false))
	actualRegistrationButtons := doc.Find(".interest-register").Length()
	actualHighInterestButtons := doc.Find(".interest-high").Length()
	actualDeregistrationButtons := doc.Find(".interest-deregister").Length()
	actualHighInterestAction := doc.Find(".interest-high").AttrOr("data-on:click", "")

	// Then
	if actualRegistrationButtons != expectedRegistrationButtons {
		t.Fatalf("registration button count mismatch\nexpected: %d\nactual:   %d", expectedRegistrationButtons, actualRegistrationButtons)
	}
	if actualHighInterestButtons != expectedHighInterestButtons {
		t.Fatalf("high-interest button count mismatch\nexpected: %d\nactual:   %d", expectedHighInterestButtons, actualHighInterestButtons)
	}
	if actualDeregistrationButtons != expectedDeregistrationButtons {
		t.Fatalf("deregistration button count mismatch\nexpected: %d\nactual:   %d", expectedDeregistrationButtons, actualDeregistrationButtons)
	}
	if !strings.Contains(actualHighInterestAction, "/event/api/ordinary-event/interest/update/interest") {
		t.Fatalf("high-interest action mismatch: %q", actualHighInterestAction)
	}
}
