package header

import (
	"testing"

	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestMenu_AnonymousUserOnlyReceivesPublicNavigation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren ikke er innlogget.",
		When:  "Når hovednavigasjonen vises.",
		Then:  "Så skal brukeren bare få navigasjonslenker til forsiden og innlogging.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "test_room_services")
	expectedHrefs := []string{"/", "/auth"}
	userInfo := requestctx.UserRequestInfo{}

	// When
	doc := templtest.Render(t, Menu(userInfo, db, logger))
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestMenu_LoggedInUserOnlyReceivesUserNavigation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren er innlogget uten adminrettigheter.",
		When:  "Når hovednavigasjonen vises.",
		Then:  "Så skal brukeren bare få navigasjonslenker til forsiden, egen profil, utlogging og vanlige spørsmål.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "test_room_services")
	expectedHrefs := []string{
		"/",
		"/profile",
		"/tilbakemelding",
		"/auth/logout",
		"https://www.regncon.no/vanlege-sporsmal/",
	}
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		IsAdmin:    false,
	}

	// When
	doc := templtest.Render(t, Menu(userInfo, db, logger))
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualExternalLinkIconVisible := doc.Find(`a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon`).Length() > 0

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	if !actualExternalLinkIconVisible {
		t.Fatalf("expected external FAQ link to include external link icon")
	}
}

func TestMenu_LoggedInUserSeesFeedbackLinkOnlyInBurgerMenus(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren er innlogget.",
		When:  "Når hovednavigasjonen vises.",
		Then:  "Så skal lenken til å gi tilbakemelding finnes i burgermenyen på skrivebord og mobil, men ikke i toppmenyen.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "test_room_services")
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
	}

	// When
	doc := templtest.Render(t, Menu(userInfo, db, logger))
	feedbackLinkInDesktopDropdown := doc.Find(`#main-menu-user-details .dropdown-panel a[href="/tilbakemelding"]`).Length() > 0
	feedbackLinkInMobileDialog := doc.Find(`#modal-hamburger-phone a[href="/tilbakemelding"]`).Length() > 0
	feedbackLinkInMainMenuButtons := doc.Find(`.main-menu-buttons a[href="/tilbakemelding"]`).Length() > 0

	// Then
	if !feedbackLinkInDesktopDropdown {
		t.Fatalf("expected feedback link in desktop dropdown menu")
	}
	if !feedbackLinkInMobileDialog {
		t.Fatalf("expected feedback link in mobile menu dialog")
	}
	if feedbackLinkInMainMenuButtons {
		t.Fatalf("expected no feedback link in the top menu bar")
	}
}

func TestMenu_AnonymousUserDoesNotSeeFeedbackLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren ikke er innlogget.",
		When:  "Når hovednavigasjonen vises.",
		Then:  "Så skal lenken til å gi tilbakemelding ikke vises.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "test_room_services")
	userInfo := requestctx.UserRequestInfo{}

	// When
	doc := templtest.Render(t, Menu(userInfo, db, logger))
	feedbackLinkVisible := doc.Find(`a[href="/tilbakemelding"]`).Length() > 0

	// Then
	if feedbackLinkVisible {
		t.Fatalf("expected anonymous menu to not contain the feedback link")
	}
}

func TestMenu_AdminUserReceivesUserAndAdminNavigation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren er admin.",
		When:  "Når hovednavigasjonen vises.",
		Then:  "Så skal brukeren få navigasjonslenker til forsiden, egen profil, utlogging, adminområdene og vanlige spørsmål.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "test_room_services")
	expectedHrefs := []string{
		"/",
		"/profile",
		"/tilbakemelding",
		"/auth/logout",
		"/admin",
		"/admin/billettholder/",
		"/admin/approval/",
		"https://www.regncon.no/vanlege-sporsmal/",
	}
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		IsAdmin:    true,
	}

	// When
	doc := templtest.Render(t, Menu(userInfo, db, logger))
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualExternalLinkIconVisible := doc.Find(`a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon`).Length() > 0

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	if !actualExternalLinkIconVisible {
		t.Fatalf("expected external FAQ link to include external link icon")
	}
}
