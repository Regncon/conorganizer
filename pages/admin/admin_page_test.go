package admin

import (
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

// TestAdminPage_RendersBreadcrumb is a small smoke check that the admin
// landing page keeps the unchanged Hjem / Admin breadcrumb. The rest of the
// page's contract (sections, cards, routes, assets, publish round trip) is
// covered in the deeper admin_landing_integration_test.go.
func TestAdminPage_RendersBreadcrumb(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at adminforsiden er lastet inn.",
		When:  "Når admininnholdet rendres.",
		Then:  "Så skal brødsmulesti til Admin være synlig.",
	})

	// Given
	expectedBreadcrumb := []string{"Admin"}
	db := testutil.CreateTestDB(t, "admin_page")
	testutil.MustExec(t, db, `
		INSERT INTO program_publishing_state(id, is_published)
		VALUES(1, 0)
		ON CONFLICT(id) DO UPDATE SET is_published = excluded.is_published
	`)

	// When
	doc := templtest.Render(t, adminPage(db))
	actualBreadcrumb := templtest.CollectTexts(doc, ".breadcrumb-end")

	// Then
	if !slices.Equal(expectedBreadcrumb, actualBreadcrumb) {
		t.Fatalf("breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedBreadcrumb, actualBreadcrumb)
	}
}
