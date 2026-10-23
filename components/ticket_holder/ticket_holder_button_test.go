package ticketholder

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/a-h/templ"
)

func TestTicketHolderButton_RendersItsStyleOnceForSeveralBillettholdere(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt flere billettholdere som vises som knapper i samme liste.",
		When:  "Når knappene rendres i samme Render-kall.",
		Then:  "Så skal style-elementet til knappen bare rendres én gang.",
	})

	// Given
	expectedStyleCount := 1
	billettholdere := []BillettHolder{
		{Id: 1, Name: "Anders Andersen", Color: ColorFromName("Anders Andersen")},
		{Id: 2, Name: "Berit Bertelsen", Color: ColorFromName("Berit Bertelsen")},
		{Id: 3, Name: "Carl Carlsen", Color: ColorFromName("Carl Carlsen")},
	}
	buttons := make([]templ.Component, 0, len(billettholdere))
	for _, billettholder := range billettholdere {
		buttons = append(buttons, TicketHolderButton(billettholder))
	}

	// When
	doc := templtest.Render(t, templ.Join(buttons...))

	// Then
	if got := countStylesContaining(doc, ".ticket-holder-picker"); got != expectedStyleCount {
		t.Fatalf("ticket holder button style count = %d, want %d", got, expectedStyleCount)
	}
	if got := doc.Find(".ticket-holder-picker").Length(); got != len(billettholdere) {
		t.Fatalf("rendered button count = %d, want %d", got, len(billettholdere))
	}
}

func countStylesContaining(doc *goquery.Document, selector string) int {
	return doc.Find("style").FilterFunction(func(_ int, style *goquery.Selection) bool {
		return strings.Contains(style.Text(), selector)
	}).Length()
}
