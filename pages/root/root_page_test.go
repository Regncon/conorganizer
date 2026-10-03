package root

import (
	"testing"

	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

var userInfo = requestctx.UserRequestInfo{}

func TestRootPageContent_DoesNotRenderBreadcrumb(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren åpner forsiden.",
		When:  "Når forsiden vises.",
		Then:  "Så skal det ikke vises noen brødsmulesti.",
	})

	// Given
	expectedBreadcrumbVisible := false

	db := createRootPageTestDB(t)
	setProgramPublishing(t, db, false)

	// When
	doc := templtest.Render(t, rootPageContent(db, nil))
	actualBreadcrumbVisible := templtest.HasSelector(doc, ".breadcrumb-container")

	// Then
	if actualBreadcrumbVisible != expectedBreadcrumbVisible {
		t.Fatalf("breadcrumb visibility mismatch\nexpected: %v\nactual:   %v", expectedBreadcrumbVisible, actualBreadcrumbVisible)
	}
}

func TestRootPageContent_WhenProgramPublishingStateCannotLoad_RendersFriendlyError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at forsiden ikke kan lese publiseringsstatus.",
		When:  "Når forsiden vises.",
		Then:  "Så skal brukeren se en vennlig feil uten tekniske detaljer.",
	})

	// Given
	expectedTextPart := rootPageLoadErrorMessage
	unexpectedTextParts := []string{
		"Error fetching",
		"program_publishing_state",
		"query program publishing state",
		"no such table",
	}

	db := createRootPageTestDB(t)
	mustExec(t, db, `DROP TABLE program_publishing_state`)

	// When
	doc := templtest.Render(t, rootPageContent(db, nil))
	actualText := rootPageText(doc)

	// Then
	assertTextContains(t, actualText, expectedTextPart)
	for _, unexpectedTextPart := range unexpectedTextParts {
		assertTextDoesNotContain(t, actualText, unexpectedTextPart)
	}
}

func TestRootPageContent_WhenEventsCannotLoad_RendersFriendlyError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at forsiden ikke kan lese arrangementslisten.",
		When:  "Når forsiden vises.",
		Then:  "Så skal brukeren se en vennlig feil uten tekniske detaljer.",
	})

	// Given
	expectedTextPart := rootEventsLoadErrorMessage
	unexpectedTextParts := []string{
		"Error fetching",
		"query announced events",
		"no such table",
	}

	db := createRootPageTestDB(t)
	setProgramPublishing(t, db, false)
	mustExec(t, db, `DROP TABLE events`)

	// When
	doc := templtest.Render(t, rootPageContent(db, nil))
	actualText := rootPageText(doc)

	// Then
	assertTextContains(t, actualText, expectedTextPart)
	for _, unexpectedTextPart := range unexpectedTextParts {
		assertTextDoesNotContain(t, actualText, unexpectedTextPart)
	}
}
