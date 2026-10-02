package root

import (
	"strings"
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

func TestRootPageContent_RendersBannerMoreEvents(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren åpner forsiden.",
		When:  "Når innsendingseksjonen vises.",
		Then:  "Så skal den gi en tydelig inngang til å sende inn arrangement.",
	})

	// Given
	expectedTextParts := []string{
		"Vil du arrangere noe under Regncon?",
		"Send inn arrangement",
	}
	expectedHref := "/profile"
	expectedImageSrc := "/static/calltoactiondragon1.webp"
	expectedImageAlt := "Send inn arrangement"

	db := createRootPageTestDB(t)
	setProgramPublishing(t, db, false)

	// When
	doc := templtest.Render(t, rootPageContent(db, nil))
	actualText := strings.Join(templtest.CollectTexts(doc, ".banner"), " ")
	actualHref, actualHrefExists := doc.Find(".banner a").Attr("href")
	actualImageSrc, actualImageSrcExists := doc.Find(".banner img.banner-avatar").Attr("src")
	actualImageAlt, actualImageAltExists := doc.Find(".banner img.banner-avatar").Attr("alt")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		assertTextContains(t, actualText, expectedTextPart)
	}
	if !actualHrefExists || actualHref != expectedHref {
		t.Fatalf("banner href mismatch\nexpected: %q\nactual:   %q", expectedHref, actualHref)
	}
	if !actualImageSrcExists || actualImageSrc != expectedImageSrc {
		t.Fatalf("banner image src mismatch\nexpected: %q\nactual:   %q", expectedImageSrc, actualImageSrc)
	}
	if !actualImageAltExists || actualImageAlt != expectedImageAlt {
		t.Fatalf("banner image alt mismatch\nexpected: %q\nactual:   %q", expectedImageAlt, actualImageAlt)
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
