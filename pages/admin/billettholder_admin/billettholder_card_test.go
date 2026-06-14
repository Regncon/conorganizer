package billettholderadmin

import (
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
