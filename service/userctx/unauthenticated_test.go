package userctx

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestUnauthenticated_RendersClearLoginAndHomeLinks(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en bruker ikke er logget inn.",
		When:  "Når innloggingsfeilsiden vises.",
		Then:  "Så skal brukeren få tydelige veier til innlogging og arrangementslisten.",
	})

	// Given
	expectedHrefs := []string{"/", "/auth"}
	expectedTextParts := []string{
		"Du må logge inn",
		"Logg inn for å se denne siden.",
		"Logg inn",
		"Gå til arrangementslisten",
	}

	// When
	doc := templtest.Render(t, Unauthenticated("/auth"))
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualText := strings.Join(templtest.CollectTexts(doc, ".access-denied"), " ")

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("unauthenticated page text mismatch\nexpected text to contain: %q\nactual text:              %q", expectedTextPart, actualText)
		}
	}
}

func TestUnauthenticated_LoginLinkCarriesGivenHref(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren ble avvist fra en bestemt side.",
		When:  "Når innloggingsfeilsiden vises med en retur-lenke for den siden.",
		Then:  "Så skal 'Logg inn'-lenken peke dit, slik at brukeren kommer tilbake etter innlogging.",
	})

	// Given
	expectedHref := "/auth?neste=%2Ftilbakemelding%3Fom%3Dfestivalen"

	// When
	doc := templtest.Render(t, Unauthenticated(expectedHref))
	actualHref, exists := doc.Find(".btn-login").Attr("href")

	// Then
	if !exists || actualHref != expectedHref {
		t.Fatalf("expected login link href %q, got %q (exists=%v)", expectedHref, actualHref, exists)
	}
}
