package banner

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestBannerMoreEvents_RendersSubmissionPrompt(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at banneret for innsending av arrangement vises.",
		When:  "Når BannerMoreEvents rendres direkte.",
		Then:  "Så skal det vise teksten og lenken for innsending.",
	})

	// Given
	expectedTextParts := []string{
		"Vil du arrangere noe under Regncon?",
		"Send inn arrangement",
	}
	expectedHref := "/profile"
	expectedImageSrc := "/static/calltoactiondragon1.webp"
	expectedImageAlt := "Send inn arrangement"

	// When
	doc := templtest.Render(t, BannerMoreEvents())
	actualText := strings.Join(templtest.CollectTexts(doc, ".banner"), " ")
	actualHref, actualHrefExists := doc.Find(".banner a").Attr("href")
	actualImageSrc, actualImageSrcExists := doc.Find(".banner img.banner-avatar").Attr("src")
	actualImageAlt, actualImageAltExists := doc.Find(".banner img.banner-avatar").Attr("alt")

	// Then
	assertBannerTextContains(t, actualText, expectedTextParts)
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

func TestBannerDuringFestival_RendersFestivalInformation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at banneret for festivaldeltakere vises.",
		When:  "Når BannerDuringFestival rendres direkte.",
		Then:  "Så skal det forklare hvordan programoversikten og puljene fungerer.",
	})

	// Given
	expectedTextParts := []string{
		"Festivalen er i gang!",
		"Arrangementene i programoversikten kan du som regel sjekke ut når det passer deg",
		"Arrangementene i puljene har begrensede plasser, som fordeles tilfeldig",
	}
	expectedImageSrc := "/static/calltoactiondragon1.webp"
	expectedImageAlt := "Sjekk ut programmet!"

	// When
	doc := templtest.Render(t, BannerDuringFestival())
	actualText := strings.Join(templtest.CollectTexts(doc, ".banner"), " ")
	actualImageSrc, actualImageSrcExists := doc.Find(".banner img.banner-avatar").Attr("src")
	actualImageAlt, actualImageAltExists := doc.Find(".banner img.banner-avatar").Attr("alt")

	// Then
	assertBannerTextContains(t, actualText, expectedTextParts)
	if doc.Find(".banner-text-small").Length() != 2 {
		t.Fatalf("banner paragraph count mismatch\nexpected: 2\nactual:   %d", doc.Find(".banner-text-small").Length())
	}
	if !actualImageSrcExists || actualImageSrc != expectedImageSrc {
		t.Fatalf("banner image src mismatch\nexpected: %q\nactual:   %q", expectedImageSrc, actualImageSrc)
	}
	if !actualImageAltExists || actualImageAlt != expectedImageAlt {
		t.Fatalf("banner image alt mismatch\nexpected: %q\nactual:   %q", expectedImageAlt, actualImageAlt)
	}
}

func assertBannerTextContains(t *testing.T, actualText string, expectedTextParts []string) {
	t.Helper()
	for _, expectedTextPart := range expectedTextParts {
		if !strings.Contains(actualText, expectedTextPart) {
			t.Fatalf("banner text missing %q\nactual: %q", expectedTextPart, actualText)
		}
	}
}
