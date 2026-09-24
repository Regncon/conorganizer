package profilecomponent

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestLoggedInAccount_RendersAccountEmail(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren er logget inn med en e-postadresse.",
		When:  "Når Kontoadministrasjon vises.",
		Then:  "Så skal brukeren se hvilken konto de er innlogget som.",
	})

	// Given
	expectedLabel := []string{"Innlogget som"}
	expectedEmail := []string{"ola.nordmann@example.no"}

	// When
	doc := templtest.Render(t, LoggedInAccount("ola.nordmann@example.no"))
	actualLabel := templtest.CollectTexts(doc, ".logged-in-account-label")
	actualEmail := templtest.CollectTexts(doc, ".logged-in-account-email")

	// Then
	if !slices.Equal(expectedLabel, actualLabel) {
		t.Fatalf("label mismatch\nexpected: %v\nactual:   %v", expectedLabel, actualLabel)
	}
	if !slices.Equal(expectedEmail, actualEmail) {
		t.Fatalf("email mismatch\nexpected: %v\nactual:   %v", expectedEmail, actualEmail)
	}
}
