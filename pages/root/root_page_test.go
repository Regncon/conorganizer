package root

import (
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestRootPageContent_RendersHomeBreadcrumb(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren åpner forsiden.",
		When:  "Når forsiden vises.",
		Then:  "Så skal brødsmulestien vise Hjem som gjeldende side.",
	})

	// Given
	expectedBreadcrumb := []string{"Hjem"}

	db := createRootPageTestDB(t)
	setProgramPublishing(t, db, false)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualBreadcrumb := templtest.CollectTexts(doc, ".breadcrumb-end")

	// Then
	if !slices.Equal(expectedBreadcrumb, actualBreadcrumb) {
		t.Fatalf("breadcrumb mismatch\nexpected: %v\nactual:   %v", expectedBreadcrumb, actualBreadcrumb)
	}
}

func TestRootPageContent_RendersSubmitEventCallToAction(t *testing.T) {
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
	expectedImageSrc := "/static/awesome-dragon-generated.png"
	expectedImageAlt := "Sent inn arrangement"

	db := createRootPageTestDB(t)
	setProgramPublishing(t, db, false)

	// When
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualText := strings.Join(templtest.CollectTexts(doc, ".call-to-action"), " ")
	actualHref, actualHrefExists := doc.Find(".call-to-action a").Attr("href")
	actualImageSrc, actualImageSrcExists := doc.Find(".call-to-action img.call-to-action-avatar").Attr("src")
	actualImageAlt, actualImageAltExists := doc.Find(".call-to-action img.call-to-action-avatar").Attr("alt")

	// Then
	for _, expectedTextPart := range expectedTextParts {
		assertTextContains(t, actualText, expectedTextPart)
	}
	if !actualHrefExists || actualHref != expectedHref {
		t.Fatalf("CTA href mismatch\nexpected: %q\nactual:   %q", expectedHref, actualHref)
	}
	if !actualImageSrcExists || actualImageSrc != expectedImageSrc {
		t.Fatalf("CTA image src mismatch\nexpected: %q\nactual:   %q", expectedImageSrc, actualImageSrc)
	}
	if !actualImageAltExists || actualImageAlt != expectedImageAlt {
		t.Fatalf("CTA image alt mismatch\nexpected: %q\nactual:   %q", expectedImageAlt, actualImageAlt)
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
	doc := templtest.Render(t, rootPageContent(db, false, nil))
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
	doc := templtest.Render(t, rootPageContent(db, false, nil))
	actualText := rootPageText(doc)

	// Then
	assertTextContains(t, actualText, expectedTextPart)
	for _, unexpectedTextPart := range unexpectedTextParts {
		assertTextDoesNotContain(t, actualText, unexpectedTextPart)
	}
}
