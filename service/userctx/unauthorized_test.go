package userctx

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestUnauthorized_RendersClearLoginAndHomeLinks(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en bruker ikke har tilgang.",
		When:  "Når tilgangsfeilsiden vises.",
		Then:  "Så skal brukeren få tydelige veier til innlogging og forsiden.",
	})

	// Given
	expectedHrefs := []string{"/", "/auth"}
	expectedTextParts := []string{
		"Du har ikkje tilgang",
		"Logg inn",
		"Gå til arrangement lista",
	}

	// When
	doc := templtest.Render(t, Unauthorized())
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualText := strings.Join(templtest.CollectTexts(doc, ".unauthorized"), " ")

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("unauthorized page text mismatch\nexpected text to contain: %q\nactual text:              %q", expectedTextPart, actualText)
		}
	}
}
