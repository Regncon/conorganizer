package profilepage

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
)

func TestSelectedBillettholderIDFromRequest_WhenQueryIDBelongsToUser_ReturnsQueryID(t *testing.T) {
	// Gitt at URL-en peker på en billettinnehaver brukeren eier,
	// når Min Side velger aktiv billettinnehaver,
	// så skal URL-valget brukes.

	// Given
	expectedSelectedID := 202
	user := profileSelectionUser("owner@example.com")
	billettholdere := []models.Billettholder{
		profileSelectionBillettholder(101, "owner@example.com"),
		profileSelectionBillettholder(expectedSelectedID, "other@example.com"),
	}
	request := profileSelectionRequest(t, "/profile?b_id=202")

	// When
	actualSelectedID := selectedBillettholderIDFromRequest(request, user, billettholdere, testutil.NewTestLogger())

	// Then
	if actualSelectedID != expectedSelectedID {
		t.Fatalf("selected billettholder mismatch\nexpected: %d\nactual:   %d", expectedSelectedID, actualSelectedID)
	}
}

func TestSelectedBillettholderIDFromRequest_WhenQueryIDIsInvalid_UsesEmailMatch(t *testing.T) {
	// Gitt at URL-en har en ugyldig billettinnehaver-ID,
	// når Min Side velger aktiv billettinnehaver,
	// så skal brukerens e-postmatch brukes som fallback.

	// Given
	expectedSelectedID := 101
	user := profileSelectionUser("owner@example.com")
	billettholdere := []models.Billettholder{
		profileSelectionBillettholder(expectedSelectedID, "owner@example.com"),
		profileSelectionBillettholder(202, "other@example.com"),
	}
	request := profileSelectionRequest(t, "/profile?b_id=ikke-et-tall")

	// When
	actualSelectedID := selectedBillettholderIDFromRequest(request, user, billettholdere, testutil.NewTestLogger())

	// Then
	if actualSelectedID != expectedSelectedID {
		t.Fatalf("selected billettholder mismatch\nexpected: %d\nactual:   %d", expectedSelectedID, actualSelectedID)
	}
}

func TestSelectedBillettholderIDFromRequest_WhenQueryIDIsNotRelatedToUser_UsesEmailMatch(t *testing.T) {
	// Gitt at URL-en peker på en billettinnehaver uten brukerrelasjon,
	// når Min Side velger aktiv billettinnehaver,
	// så skal brukerens e-postmatch brukes som fallback.

	// Given
	expectedSelectedID := 101
	user := profileSelectionUser("owner@example.com")
	billettholdere := []models.Billettholder{
		profileSelectionBillettholder(expectedSelectedID, "owner@example.com"),
		profileSelectionBillettholder(202, "other@example.com"),
	}
	request := profileSelectionRequest(t, "/profile?b_id=999")

	// When
	actualSelectedID := selectedBillettholderIDFromRequest(request, user, billettholdere, testutil.NewTestLogger())

	// Then
	if actualSelectedID != expectedSelectedID {
		t.Fatalf("selected billettholder mismatch\nexpected: %d\nactual:   %d", expectedSelectedID, actualSelectedID)
	}
}

func TestSelectedBillettholderIDFromRequest_WhenNoQueryAndEmailMatches_ReturnsEmailMatchedHolder(t *testing.T) {
	// Gitt at URL-en ikke velger billettinnehaver og brukerens e-post matcher en billettinnehaver,
	// når Min Side velger aktiv billettinnehaver,
	// så skal e-postmatchen brukes.

	// Given
	expectedSelectedID := 202
	user := profileSelectionUser("owner@example.com")
	billettholdere := []models.Billettholder{
		profileSelectionBillettholder(101, "first@example.com"),
		profileSelectionBillettholder(expectedSelectedID, "OWNER@example.com"),
	}
	request := profileSelectionRequest(t, "/profile")

	// When
	actualSelectedID := selectedBillettholderIDFromRequest(request, user, billettholdere, testutil.NewTestLogger())

	// Then
	if actualSelectedID != expectedSelectedID {
		t.Fatalf("selected billettholder mismatch\nexpected: %d\nactual:   %d", expectedSelectedID, actualSelectedID)
	}
}

func TestSelectedBillettholderIDFromRequest_WhenNoQueryAndNoEmailMatch_ReturnsFirstHolder(t *testing.T) {
	// Gitt at URL-en ikke velger billettinnehaver og ingen e-post matcher brukeren,
	// når Min Side velger aktiv billettinnehaver,
	// så skal første tilgjengelige billettinnehaver brukes.

	// Given
	expectedSelectedID := 101
	user := profileSelectionUser("owner@example.com")
	billettholdere := []models.Billettholder{
		profileSelectionBillettholder(expectedSelectedID, "first@example.com"),
		profileSelectionBillettholder(202, "second@example.com"),
	}
	request := profileSelectionRequest(t, "/profile")

	// When
	actualSelectedID := selectedBillettholderIDFromRequest(request, user, billettholdere, testutil.NewTestLogger())

	// Then
	if actualSelectedID != expectedSelectedID {
		t.Fatalf("selected billettholder mismatch\nexpected: %d\nactual:   %d", expectedSelectedID, actualSelectedID)
	}
}

func TestSelectedBillettholderIDFromRequest_WhenUserHasNoTicketHolders_ReturnsZero(t *testing.T) {
	// Gitt at brukeren ikke har billettinnehavere,
	// når Min Side velger aktiv billettinnehaver,
	// så skal ingen billettinnehaver være valgt.

	// Given
	expectedSelectedID := 0
	user := profileSelectionUser("owner@example.com")
	request := profileSelectionRequest(t, "/profile")

	// When
	actualSelectedID := selectedBillettholderIDFromRequest(request, user, nil, testutil.NewTestLogger())

	// Then
	if actualSelectedID != expectedSelectedID {
		t.Fatalf("selected billettholder mismatch\nexpected: %d\nactual:   %d", expectedSelectedID, actualSelectedID)
	}
}

func profileSelectionUser(email string) requestctx.UserRequestInfo {
	return requestctx.UserRequestInfo{
		IsLoggedIn: true,
		Id:         "profile-selection-user",
		Email:      email,
	}
}

func profileSelectionBillettholder(id int, emails ...string) models.Billettholder {
	billettholderEmails := make([]models.BillettholderEmail, 0, len(emails))
	for _, email := range emails {
		billettholderEmails = append(billettholderEmails, models.BillettholderEmail{Email: email})
	}

	return models.Billettholder{
		ID:     id,
		Emails: billettholderEmails,
	}
}

func profileSelectionRequest(t *testing.T, target string) *http.Request {
	t.Helper()

	return httptest.NewRequest(http.MethodGet, target, nil)
}
