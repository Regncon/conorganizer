package header

import (
	"testing"

	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestMenu_AnonymousUserOnlyReceivesPublicNavigation(t *testing.T) {
	// Gitt at brukeren ikke er innlogget,
	// når hovednavigasjonen vises,
	// så skal brukeren bare få navigasjonslenker til forsiden og innlogging.

	// Given
    db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()
	expectedHrefs := []string{"/", "/auth"}
	userInfo := requestctx.UserRequestInfo{}

	// When
	doc := templtest.Render(t, Menu(userInfo, db))
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestMenu_LoggedInUserOnlyReceivesUserNavigation(t *testing.T) {
	// Gitt at brukeren er innlogget uten adminrettigheter,
	// når hovednavigasjonen vises,
	// så skal brukeren bare få navigasjonslenker til forsiden, egen profil, utlogging og vanlege spørsmål.

	// Given
    db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()
	expectedHrefs := []string{
		"/",
		"/profile",
		"/auth/logout",
		"https://www.regncon.no/vanlege-sporsmal/",
	}
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		IsAdmin:    false,
	}

	// When
	doc := templtest.Render(t, Menu(userInfo, db))
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestMenu_AdminUserReceivesUserAndAdminNavigation(t *testing.T) {
	// Gitt at brukeren er admin,
	// når hovednavigasjonen vises,
	// så skal brukeren få navigasjonslenker til forsiden, egen profil, utlogging, adminområdene og vanlege spørsmål.

	// Given
    db, _, err := testutil.CreateTemporaryDBAndLogger("test_room_services", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()
	expectedHrefs := []string{
		"/",
		"/profile",
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
	doc := templtest.Render(t, Menu(userInfo, db))
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}
