package profilepage

import (
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestProfilePage_RendersBreadcrumbAndBillettholderSelectionMetadata(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at Min Side vises med en valgt billettinnehaver.",
		When:  "Når profilsiden rendres.",
		Then:  "Så skal brødsmulestien og metadata for gyldig billettvalg være tilgjengelig.",
	})

	// Given
	expectedBreadcrumb := []string{"Min Side"}
	expectedSelectedID := "22"
	expectedValidIDs := "[11,22]"
	expectedInitPath := "/profile/api"

	db, logger := testutil.CreateTestDBAndLogger(t, "profile_page")
	user := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		Id:         "profile-page-user",
		Email:      "profile-page-user@example.com",
	}

	// When
	doc := templtest.Render(t, ProfilePage(user, nil, nil, 22, []int{11, 22}, db, logger, nil))
	actualBreadcrumb := templtest.CollectTexts(doc, ".breadcrumb-end")
	profileContainer := doc.Find(".profile-container")
	actualSelectedID, actualSelectedIDExists := profileContainer.Attr("data-profile-selected-billettholder-id")
	actualValidIDs, actualValidIDsExists := profileContainer.Attr("data-profile-valid-billettholder-ids")
	actualInit, actualInitExists := profileContainer.Attr("data-init")

	// Then
	if !slices.Equal(expectedBreadcrumb, actualBreadcrumb) {
		t.Fatalf("breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedBreadcrumb, actualBreadcrumb)
	}
	if !actualSelectedIDExists || actualSelectedID != expectedSelectedID {
		t.Fatalf("selected billettholder metadata mismatch\nexpected: %q\nactual:   %q", expectedSelectedID, actualSelectedID)
	}
	if !actualValidIDsExists || actualValidIDs != expectedValidIDs {
		t.Fatalf("valid billettholder metadata mismatch\nexpected: %q\nactual:   %q", expectedValidIDs, actualValidIDs)
	}
	if !actualInitExists || !strings.Contains(actualInit, expectedInitPath) {
		t.Fatalf("profile live init mismatch\nexpected data-init to contain: %q\nactual:                        %q", expectedInitPath, actualInit)
	}
}
