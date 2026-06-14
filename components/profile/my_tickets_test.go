package profilecomponent

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestMyTickets_RendersTicketHolderSummaryAndTicketsLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren har billettinnehavere knyttet til profilen.",
		When:  "Når billettseksjonen vises på Min Side.",
		Then:  "Så skal navn, billettype, e-post og lenken til billettsiden være synlig.",
	})

	// Given
	expectedTextParts := []string{
		"Ola Nordmann",
		"Festivalpass",
		"ola@example.com",
		"Mine billettar",
	}
	expectedHref := "/profile/tickets"
	tickets := []TicketHolder{
		{Name: "Ola Nordmann", Ticket: "Festivalpass", Email: "ola@example.com"},
	}

	// When
	doc := templtest.Render(t, MyTickets(tickets))
	actualText := strings.Join(templtest.CollectTexts(doc, ".surface-pane"), " ")
	actualHref, actualHrefExists := doc.Find(`a[href="/profile/tickets"]`).Attr("href")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("expected ticket section text to contain %q\nactual text: %s", expectedTextPart, actualText)
		}
	}
	if !actualHrefExists || actualHref != expectedHref {
		t.Fatalf("ticket page href mismatch\nexpected: %q\nactual:   %q", expectedHref, actualHref)
	}
}

func TestMyTickets_WhenUserHasNoTicketHolders_RendersTicketsLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren ikke har billettinnehavere knyttet til profilen.",
		When:  "Når billettseksjonen vises på Min Side.",
		Then:  "Så skal brukeren fortsatt kunne gå videre til billettsiden.",
	})

	// Given
	expectedHref := "/profile/tickets"
	expectedLinkText := []string{"Mine billettar"}

	// When
	doc := templtest.Render(t, MyTickets(nil))
	actualLinkText := templtest.CollectTexts(doc, `a[href="/profile/tickets"]`)
	actualHref, actualHrefExists := doc.Find(`a[href="/profile/tickets"]`).Attr("href")

	// Then
	if !actualHrefExists || actualHref != expectedHref {
		t.Fatalf("ticket page href mismatch\nexpected: %q\nactual:   %q", expectedHref, actualHref)
	}
	if strings.Join(actualLinkText, " ") != strings.Join(expectedLinkText, " ") {
		t.Fatalf("ticket page link text mismatch\nexpected: %v\nactual:   %v", expectedLinkText, actualLinkText)
	}
}
