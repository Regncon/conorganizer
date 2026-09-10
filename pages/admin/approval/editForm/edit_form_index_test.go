package edit_form

import (
	"context"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

// The live stream morphs everything inside #edit-form-container. A <dialog> in
// that subtree loses its `open` attribute on the next patch without a `close`
// event, which would leave the age warning wedged shut. It must sit outside.
func TestEditFormIndex_AgeDialogRendersOutsideLiveRegion(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt redigeringssida for eit arrangement.",
		When:  "Når sida blir rendra.",
		Then:  "Så ligg aldersgrense-dialogen utanfor det live-oppdaterte området.",
	})

	db := testutil.CreateTestDB(t, "edit-form-index-age-dialog")
	logger := testutil.NewSlogAdapter(&testutil.StubLogger{})
	seedEditFormNavigationLookups(t, db)
	seedEditFormNavigationEvent(t, db, "submitted-event", "Submitted Event", models.EventStatusSubmitted, "2026-01-02T10:00:00Z")

	doc := templtest.Render(t, editFormIndex("submitted-event", context.Background(), db, nil, logger))

	if got := doc.Find("#approval-age-dialog").Length(); got != 1 {
		t.Fatalf("expected exactly one age warning dialog on the page, got %d", got)
	}
	if got := doc.Find("#edit-form-container #approval-age-dialog").Length(); got != 0 {
		t.Errorf("age dialog must NOT be inside the live #edit-form-container")
	}
}
