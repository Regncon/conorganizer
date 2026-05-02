package header

import (
	"testing"

	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

// Gitt at brukeren ikke er innlogget, når hovednavigasjonen vises,
// så skal brukeren bare få interne navigasjonslenker til forsiden og innlogging.
func TestMenu_AnonymousUserOnlyReceivesPublicInternalNavigation(t *testing.T) {
	// Given
	expectedInternalHrefs := []string{"/", "/auth"}
	userInfo := requestctx.UserRequestInfo{}

	// When
	doc := templtest.Render(t, Menu(userInfo))
	actualInternalHrefs := templtest.CollectUniqueInternalHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedInternalHrefs, actualInternalHrefs)
}

// Gitt at brukeren er innlogget uten adminrettigheter, når hovednavigasjonen vises,
// så skal brukeren bare få interne navigasjonslenker til forsiden, egen profil og utlogging.
func TestMenu_LoggedInUserOnlyReceivesUserInternalNavigation(t *testing.T) {
	// Given
	expectedInternalHrefs := []string{"/", "/profile", "/auth/logout"}
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		IsAdmin:    false,
	}

	// When
	doc := templtest.Render(t, Menu(userInfo))
	actualInternalHrefs := templtest.CollectUniqueInternalHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedInternalHrefs, actualInternalHrefs)
}

// Gitt at brukeren er admin, når hovednavigasjonen vises,
// så skal brukeren få interne navigasjonslenker til forsiden, egen profil, utlogging og adminområdene.
func TestMenu_AdminUserReceivesUserAndAdminInternalNavigation(t *testing.T) {
	// Given
	expectedInternalHrefs := []string{
		"/",
		"/profile",
		"/auth/logout",
		"/admin",
		"/admin/billettholder/",
		"/admin/approval/",
	}
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		IsAdmin:    true,
	}

	// When
	doc := templtest.Render(t, Menu(userInfo))
	actualInternalHrefs := templtest.CollectUniqueInternalHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedInternalHrefs, actualInternalHrefs)
}
