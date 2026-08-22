package authctx

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestForbidden_RendersAdminAccessMessageAndHomeLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en innlogget bruker mangler adminrolle.",
		When:  "Når tilgangsfeilsiden vises.",
		Then:  "Så skal brukeren få en tydelig forklaring og lenke til arrangementslisten.",
	})

	// Given
	expectedHrefs := []string{"/"}
	expectedTextParts := []string{
		"Du har ikke tilgang",
		"Du er logget inn, men denne siden krever administratortilgang.",
		"Gå til arrangementslisten",
	}

	// When
	doc := templtest.Render(t, Forbidden())
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualText := strings.Join(templtest.CollectTexts(doc, ".access-denied"), " ")

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("forbidden page text mismatch\nexpected text to contain: %q\nactual text:              %q", expectedTextPart, actualText)
		}
	}
}
