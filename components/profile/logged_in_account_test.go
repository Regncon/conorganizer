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
		When:  "Når kontoboksen vises på Min Side.",
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

func TestLoggedInAccount_RendersResetPasswordAndLogoutLinks(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren er logget inn.",
		When:  "Når kontoboksen vises på Min Side.",
		Then:  "Så skal brukeren kunne resette passord og logge ut fra boksen.",
	})

	// Given
	expectedResetPasswordText := []string{"Reset passord"}
	expectedLogoutText := []string{"Logg ut"}

	// When
	doc := templtest.Render(t, LoggedInAccount("ola.nordmann@example.no"))
	actualResetPasswordText := templtest.CollectTexts(doc, `a[href="/profile/descope-profile"]`)
	actualLogoutText := templtest.CollectTexts(doc, `a[href="/auth/logout"]`)

	// Then
	if !slices.Equal(expectedResetPasswordText, actualResetPasswordText) {
		t.Fatalf("reset password link mismatch\nexpected: %v\nactual:   %v", expectedResetPasswordText, actualResetPasswordText)
	}
	if !slices.Equal(expectedLogoutText, actualLogoutText) {
		t.Fatalf("logout link mismatch\nexpected: %v\nactual:   %v", expectedLogoutText, actualLogoutText)
	}
}
