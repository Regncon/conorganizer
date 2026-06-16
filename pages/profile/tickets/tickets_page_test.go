package profileticketspage

import (
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestProfileTicketsPageContent_WhenUserHasNoTicketHolders_RendersEmptyStateAndFetchAction(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren ikke har billettinnehavere knyttet til profilen.",
		When:  "Når billettsiden rendres.",
		Then:  "Så skal siden vise brødsmulesti, tomtilstand og handling for å hente billetter.",
	})

	// Given
	expectedBreadcrumb := []string{"Billetter"}
	expectedTextParts := []string{
		"Mine billetter",
		"Ingen billetter funnet",
		"Hent billetter",
		"Kjøp billetter",
	}
	db, logger := testutil.CreateTestDBAndLogger(t, "profile_tickets_page")

	// When
	doc := templtest.Render(t, profileTicketsPageContent("profile-tickets-user", db, logger))
	actualBreadcrumb := templtest.CollectTexts(doc, ".breadcrumb-end")
	actualText := strings.Join(templtest.CollectTexts(doc, "#profile-tickets-page-wrapper, section"), " ")
	actualFetchButtonVisible := templtest.HasSelector(doc, `button[data-on\:click="@post('/profile/tickets/api/get-tickets')"]`)
	actualBuyHref, actualBuyHrefExists := doc.Find(`a[href="https://event.checkin.no/221572/regncon-xxxiv-2026"]`).Attr("href")
	actualPlaceholderImageVisible := templtest.HasSelector(doc, `img[src="/static/sobbingtemp.png"]`)

	// Then
	if !slices.Equal(expectedBreadcrumb, actualBreadcrumb) {
		t.Fatalf("breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedBreadcrumb, actualBreadcrumb)
	}
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("expected ticket page text to contain %q\nactual text: %s", expectedTextPart, actualText)
		}
	}
	if !actualFetchButtonVisible {
		t.Fatalf("expected fetch tickets button to be visible")
	}
	if !actualBuyHrefExists || actualBuyHref != checkInTicketsURL {
		t.Fatalf("buy ticket href mismatch\nexpected: %q\nactual:   %q", checkInTicketsURL, actualBuyHref)
	}
	if actualPlaceholderImageVisible {
		t.Fatalf("expected placeholder image to be removed from empty state")
	}
}

func TestBillettholderCard_RendersEmailKindsAndOnlyAllowsManualDelete(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettinnehaver har billett-, tilknyttet og manuell e-post.",
		When:  "Når billettkortet rendres.",
		Then:  "Så skal e-postene vises og bare den manuelle adressen ha slettehandling.",
	})

	// Given
	expectedTextParts := []string{
		"Festivalpass",
		"Kari Nordmann",
		"ticket@example.com",
		"associated@example.com",
		"manual@example.com",
		"Over 18",
	}
	expectedDeleteAction := "@post('/profile/tickets/api/delete-email/42/3/')"
	billettholder := models.Billettholder{
		ID:         42,
		FirstName:  "Kari",
		LastName:   "Nordmann",
		TicketType: "Festivalpass",
		IsOver18:   true,
		Emails: []models.BillettholderEmail{
			{ID: 1, Email: "ticket@example.com", Kind: models.BillettholderEmailKindTicket},
			{ID: 2, Email: "associated@example.com", Kind: models.BillettholderEmailKindAssociated},
			{ID: 3, Email: "manual@example.com", Kind: models.BillettholderEmailKindManual},
		},
	}

	// When
	doc := templtest.Render(t, billettholderCard(billettholder))
	actualText := strings.Join(templtest.CollectTexts(doc, ".card"), " ")
	deleteButtons := doc.Find(`button[title="Slett epostadresse"]`)
	actualDeleteAction, actualDeleteActionExists := deleteButtons.Attr("data-on:click")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("expected ticket card text to contain %q\nactual text: %s", expectedTextPart, actualText)
		}
	}
	if deleteButtons.Length() != 1 {
		t.Fatalf("delete button count mismatch\nexpected: 1\nactual:   %d", deleteButtons.Length())
	}
	if !actualDeleteActionExists || actualDeleteAction != expectedDeleteAction {
		t.Fatalf("delete action mismatch\nexpected: %q\nactual:   %q", expectedDeleteAction, actualDeleteAction)
	}
}
