package edit_form

import (
	"context"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestEditFormIndex_HasNoPlayerAssignmentDialogs(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt redigeringssiden for et arrangement.",
		When:  "Når siden blir rendret.",
		Then:  "Så inneholder den ikke lenger spillertildeling.",
	})

	db := testutil.CreateTestDB(t, "edit-form-index-age-dialog")
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEditFormNavigationLookups(t, db)
	seedEditFormNavigationEvent(t, db, "submitted-event", "Submitted Event", models.EventStatusSubmitted, "2026-01-02T10:00:00Z")

	doc := templtest.Render(t, editFormIndex("submitted-event", context.Background(), db, nil, logger))

	if got := doc.Find("#approval-age-dialog, #tildeling-dialog, .who-is-interested").Length(); got != 0 {
		t.Fatalf("expected no player assignment UI on the event form, got %d elements", got)
	}
}
