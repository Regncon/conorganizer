package admin

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestAdminPage_RendersBreadcrumbAndAdminAreaLinks(t *testing.T) {
	// Gitt at adminforsiden er lastet inn,
	// når admininnholdet rendres,
	// så skal brødsmulesti og lenker til de viktigste adminområdene være synlige.

	// Given
	expectedBreadcrumb := []string{"Admin"}
	expectedHrefs := []string{
		"/admin/approval/",
		"/admin/billettholder/",
		"/admin/rooms/",
	}
	db := testutil.CreateTestDB(t, "admin_page")
	testutil.MustExec(t, db, `
		INSERT INTO program_publishing_state(id, is_published)
		VALUES(1, 0)
		ON CONFLICT(id) DO UPDATE SET is_published = excluded.is_published
	`)

	// When
	doc := templtest.Render(t, adminPage(db))
	actualBreadcrumb := templtest.CollectTexts(doc, ".breadcrumb-end")
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	if !slices.Equal(expectedBreadcrumb, actualBreadcrumb) {
		t.Fatalf("breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedBreadcrumb, actualBreadcrumb)
	}
	for _, expectedHref := range expectedHrefs {
		if !slices.Contains(actualHrefs, expectedHref) {
			t.Fatalf("expected admin page href %q in %v", expectedHref, actualHrefs)
		}
	}
}
