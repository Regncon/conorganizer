package addbillettholder

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/checkIn"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestTicketCard_WhenTicketCanBeConverted_RendersConvertAction(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en check-in-billett ikke er middag og ikke er konvertert.",
		When:  "Når billettkortet rendres.",
		Then:  "Så skal admin se handlingen for å konvertere billetten til deltaker.",
	})

	// Given
	expectedTextParts := []string{
		"Bestilling:",
		"9001",
		"Adult",
		"OlaNordmann",
		"ola@example.com",
		"Over 18",
		"Konverter billett til deltager",
	}
	ticket := addBillettholderTestTicket(1)

	// When
	doc := templtest.Render(t, ticketCard(ticket, false, ""))
	actualText := strings.Join(templtest.CollectTexts(doc, ".card"), " ")
	actualConvertButtonVisible := templtest.HasSelector(doc, "button.btn--outline")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("expected ticket card text to contain %q\nactual text: %s", expectedTextPart, actualText)
		}
	}
	if !actualConvertButtonVisible {
		t.Fatalf("expected convert button to be visible")
	}
}

func TestTicketCard_WhenTicketIsDinner_RendersDinnerWarningWithoutConvertAction(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en check-in-billett er en middagsbillett.",
		When:  "Når billettkortet rendres.",
		Then:  "Så skal kortet markere middag og ikke tilby vanlig deltakerkonvertering.",
	})

	// Given
	expectedTextPart := "Dette er en middagsbillett"
	unexpectedTextPart := "Konverter billett til deltager"
	ticket := addBillettholderTestTicket(checkIn.TicketTypeMiddag)

	// When
	doc := templtest.Render(t, ticketCard(ticket, false, ""))
	actualText := strings.Join(templtest.CollectTexts(doc, ".card"), " ")

	// Then
	if !strings.Contains(actualText, expectedTextPart) {
		t.Fatalf("expected ticket card text to contain %q\nactual text: %s", expectedTextPart, actualText)
	}
	if strings.Contains(actualText, unexpectedTextPart) {
		t.Fatalf("expected ticket card text not to contain %q\nactual text: %s", unexpectedTextPart, actualText)
	}
}

func TestTicketCard_WhenTicketIsAlreadyBillettholder_RendersConvertedWarningWithoutConvertAction(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en check-in-billett allerede er konvertert til billettholder.",
		When:  "Når billettkortet rendres.",
		Then:  "Så skal kortet markere dette og ikke tilby ny konvertering.",
	})

	// Given
	expectedTextPart := "allerede konvertert til en billettholder"
	unexpectedTextPart := "Konverter billett til deltager"
	ticket := addBillettholderTestTicket(1)

	// When
	doc := templtest.Render(t, ticketCard(ticket, true, ""))
	actualText := strings.Join(templtest.CollectTexts(doc, ".card"), " ")

	// Then
	if !strings.Contains(actualText, expectedTextPart) {
		t.Fatalf("expected ticket card text to contain %q\nactual text: %s", expectedTextPart, actualText)
	}
	if strings.Contains(actualText, unexpectedTextPart) {
		t.Fatalf("expected ticket card text not to contain %q\nactual text: %s", unexpectedTextPart, actualText)
	}
}

func TestIsBillettholder_WhenTicketIDExists_ReturnsTrue(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en check-in-billett allerede finnes som billettholder.",
		When:  "Når siden sjekker konverteringsstatus.",
		Then:  "Så skal billetten behandles som konvertert.",
	})

	// Given
	expectedResult := true
	billettholdere := []models.Billettholder{
		{TicketID: 1001},
		{TicketID: 1002},
	}

	// When
	actualResult := isBillettholder(1002, billettholdere)

	// Then
	if actualResult != expectedResult {
		t.Fatalf("isBillettholder mismatch\nexpected: %v\nactual:   %v", expectedResult, actualResult)
	}
}

func addBillettholderTestTicket(ticketTypeID int) checkIn.CheckInTicket {
	return checkIn.CheckInTicket{
		ID:        1001,
		OrderID:   9001,
		TypeId:    ticketTypeID,
		Type:      "Adult",
		FirstName: "Ola",
		LastName:  "Nordmann",
		Email:     "ola@example.com",
		IsOver18:  true,
	}
}
