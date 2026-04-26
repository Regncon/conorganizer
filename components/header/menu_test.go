package header

import (
	"bytes"
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/service/requestctx"
)

// Gitt at brukeren ikke er innlogget, når hovednavigasjonen vises,
// så skal brukeren bare få interne navigasjonslenker til forsiden og innlogging.
func TestMenu_AnonymousUserOnlyReceivesPublicInternalNavigation(t *testing.T) {
	// Given
	expectedInternalHrefs := []string{"/", "/auth"}
	userInfo := requestctx.UserRequestInfo{}

	// When
	doc := renderMenu(t, userInfo)
	actualInternalHrefs := collectUniqueInternalMenuHrefs(doc)

	// Then
	assertSameHrefs(t, expectedInternalHrefs, actualInternalHrefs)
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
	doc := renderMenu(t, userInfo)
	actualInternalHrefs := collectUniqueInternalMenuHrefs(doc)

	// Then
	assertSameHrefs(t, expectedInternalHrefs, actualInternalHrefs)
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
	doc := renderMenu(t, userInfo)
	actualInternalHrefs := collectUniqueInternalMenuHrefs(doc)

	// Then
	assertSameHrefs(t, expectedInternalHrefs, actualInternalHrefs)
}

func renderMenu(t *testing.T, userInfo requestctx.UserRequestInfo) *goquery.Document {
	t.Helper()

	var html bytes.Buffer
	if err := Menu(userInfo).Render(context.Background(), &html); err != nil {
		t.Fatalf("render menu: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(&html)
	if err != nil {
		t.Fatalf("parse menu html: %v", err)
	}

	return doc
}

func collectUniqueInternalMenuHrefs(doc *goquery.Document) []string {
	uniqueHrefs := make(map[string]struct{})

	doc.Find("a[href]").Each(func(_ int, link *goquery.Selection) {
		href, exists := link.Attr("href")
		if !exists {
			return
		}

		href = strings.TrimSpace(href)
		if href == "" || !strings.HasPrefix(href, "/") || strings.HasPrefix(href, "//") {
			return
		}

		uniqueHrefs[href] = struct{}{}
	})

	hrefs := make([]string, 0, len(uniqueHrefs))
	for href := range uniqueHrefs {
		hrefs = append(hrefs, href)
	}

	slices.Sort(hrefs)
	return hrefs
}

func assertSameHrefs(t *testing.T, expectedHrefs []string, actualHrefs []string) {
	t.Helper()

	expectedHrefs = sortedCopy(expectedHrefs)
	actualHrefs = sortedCopy(actualHrefs)

	if !slices.Equal(expectedHrefs, actualHrefs) {
		t.Fatalf("expected internal menu hrefs %v, got %v", expectedHrefs, actualHrefs)
	}
}

func sortedCopy(values []string) []string {
	copiedValues := append([]string(nil), values...)
	slices.Sort(copiedValues)
	return copiedValues
}
