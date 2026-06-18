package billettholderadmin

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestBillettholderCard_RendersDetailsAndOnlyAllowsManualEmailDelete(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder har billett-, tilknyttet og manuell e-post.",
		When:  "Når adminkortet rendres.",
		Then:  "Så skal detaljene vises og bare manuell e-post kunne slettes.",
	})

	// Given
	expectedTextParts := []string{
		"Bestilling:",
		"9001",
		"Festivalpass",
		"Kari Nordmann",
		"ticket@example.com",
		"associated@example.com",
		"manual@example.com",
		"Under 18",
	}
	expectedDeleteAction := "@post('/admin/billettholder/api/delete-email/42/3/')"
	billettholder := models.Billettholder{
		ID:         42,
		FirstName:  "Kari",
		LastName:   "Nordmann",
		TicketType: "Festivalpass",
		OrderID:    9001,
		IsOver18:   false,
		Emails: []models.BillettholderEmail{
			{ID: 1, Email: "ticket@example.com", Kind: models.BillettholderEmailKindTicket},
			{ID: 2, Email: "associated@example.com", Kind: models.BillettholderEmailKindAssociated},
			{ID: 3, Email: "manual@example.com", Kind: models.BillettholderEmailKindManual},
		},
	}

	// When
	doc := templtest.Render(t, billettholderCard(billettholder, "", nil))
	actualText := strings.Join(templtest.CollectTexts(doc, ".billettholder-card"), " ")
	deleteButtons := doc.Find(`button[title="Slett epostadresse"]`)
	actualDeleteAction, actualDeleteActionExists := deleteButtons.Attr("data-on:click")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("expected billettholder card text to contain %q\nactual text: %s", expectedTextPart, actualText)
		}
	}
	if deleteButtons.Length() != 1 {
		t.Fatalf("delete button count mismatch\nexpected: 1\nactual:   %d", deleteButtons.Length())
	}
	if !actualDeleteActionExists || actualDeleteAction != expectedDeleteAction {
		t.Fatalf("delete action mismatch\nexpected: %q\nactual:   %q", expectedDeleteAction, actualDeleteAction)
	}
}

func TestBillettholderCard_UsesDatastarStateForInterestDialogOpenState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at et billettholderkort har en interessedialog.",
		When:  "Når adminkortet rendres.",
		Then:  "Så skal åpning styres av Datastar-state og dialogens native open-attributt bevares under live patches.",
	})

	// Given
	billettholder := models.Billettholder{
		ID:        42,
		FirstName: "Kari",
		LastName:  "Nordmann",
	}
	expectedDialogID := "billettholder-interests-42"
	expectedOpenAction := fmt.Sprintf(
		"$billettholderInterestActivePulje = '%s'; $billettholderInterestOpenDialogId = '%s'",
		models.PuljeFredagKveld,
		expectedDialogID,
	)

	// When
	doc := templtest.Render(t, billettholderCard(billettholder, "", nil))
	openButton := doc.Find(".billettholder-interest-open")
	dialog := doc.Find("#" + expectedDialogID)
	actualOpenAction, actualOpenActionExists := openButton.Attr("data-on:click")
	actualPreserveAttr, actualPreserveAttrExists := dialog.Attr("data-preserve-attr")
	actualEffect, actualEffectExists := dialog.Attr("data-effect")
	actualCloseAction, actualCloseActionExists := dialog.Attr("data-on:close")

	// Then
	if !actualOpenActionExists || actualOpenAction != expectedOpenAction {
		t.Fatalf("interest dialog open action mismatch\nexpected: %q\nactual:   %q", expectedOpenAction, actualOpenAction)
	}
	if !actualPreserveAttrExists || actualPreserveAttr != "open" {
		t.Fatalf("dialog preserve attr mismatch\nexpected: %q\nactual:   %q", "open", actualPreserveAttr)
	}
	for _, expectedPart := range []string{
		"$billettholderInterestOpenDialogId == '" + expectedDialogID + "'",
		"el.showModal()",
		"el.close()",
	} {
		if !actualEffectExists || !strings.Contains(actualEffect, expectedPart) {
			t.Fatalf("expected dialog effect to contain %q\nactual effect: %q", expectedPart, actualEffect)
		}
	}
	if !actualCloseActionExists || !strings.Contains(actualCloseAction, "$billettholderInterestOpenDialogId = ''") {
		t.Fatalf("expected close action to reset open dialog id\nactual close action: %q", actualCloseAction)
	}
}

func TestHighlightSearchTerm_EscapesTextBeforeHighlighting(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at søket matcher tekst som også inneholder HTML-tegn.",
		When:  "Når teksten markeres.",
		Then:  "Så skal originalteksten escapes før markeringen settes inn.",
	})

	// Given
	expectedTextPart := "&lt;"
	unexpectedTextPart := "<script>"

	// When
	actualHTML := highlightSearchTerm("<script>", "s")

	// Then
	if !strings.Contains(actualHTML, expectedTextPart) {
		t.Fatalf("expected highlighted HTML to contain escaped text %q\nactual HTML: %s", expectedTextPart, actualHTML)
	}
	if strings.Contains(actualHTML, unexpectedTextPart) {
		t.Fatalf("expected highlighted HTML not to contain raw script tag\nactual HTML: %s", actualHTML)
	}
}
