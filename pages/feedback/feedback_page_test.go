package feedback

import (
	"strconv"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/service/feedback"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestFeedbackFormPage_RendersPrivacyHintAndCategories(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en innlogget bruker åpner tilbakemeldingssiden uten forhåndsvalgt kategori.",
		When:  "Når siden rendres.",
		Then:  "Så skal den vise et tydelig personvernhint og alle tre kategoriene, ingen forhåndsvalgt.",
	})

	// Given
	expectedHintPart := "Ikke skriv personopplysninger"
	expectedLabels := []string{"Nettsiden", "Festivalen", "Noe annet"}

	// When
	doc := templtest.Render(t, feedbackFormPage(""))
	actualHint := strings.Join(templtest.CollectTexts(doc, ".feedback-privacy-hint"), " ")
	buttons := doc.Find(".feedback-category-button")

	// Then
	if !strings.Contains(actualHint, expectedHintPart) {
		t.Fatalf("expected privacy hint to contain %q, got %q", expectedHintPart, actualHint)
	}
	if buttons.Length() != len(expectedLabels) {
		t.Fatalf("expected %d category buttons, got %d", len(expectedLabels), buttons.Length())
	}
	for i, expectedLabel := range expectedLabels {
		button := buttons.Eq(i)
		actualLabel := strings.TrimSpace(button.Text())
		if actualLabel != expectedLabel {
			t.Fatalf("expected category button %d to be %q, got %q", i, expectedLabel, actualLabel)
		}
		if pressed, exists := button.Attr("aria-pressed"); !exists || pressed != "false" {
			t.Fatalf("expected category button %q to start unpressed, got %q (exists=%v)", expectedLabel, pressed, exists)
		}
	}
	if doc.Find("textarea#feedback-message").AttrOr("maxlength", "") != strconv.Itoa(feedback.MaxMessageLength) {
		t.Fatalf("expected textarea maxlength %d", feedback.MaxMessageLength)
	}
	if templtest.CollectTexts(doc, "button[type=submit]")[0] != "Send tilbakemelding" {
		t.Fatalf("expected submit button labelled 'Send tilbakemelding'")
	}
}

func TestFeedbackFormPage_PreselectsCategoryFromOmParam(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren kom fra en QR-kode med ?om=festivalen.",
		When:  "Når tilbakemeldingssiden rendres med den forhåndsvalgte kategorien.",
		Then:  "Så skal 'Festivalen' vises som trykket, og de andre som ikke trykket.",
	})

	// Given
	preselected := feedback.CategoryConvention

	// When
	doc := templtest.Render(t, feedbackFormPage(preselected))

	// Then
	pressedTexts := []string{}
	doc.Find(`.feedback-category-button[aria-pressed="true"]`).Each(func(_ int, button *goquery.Selection) {
		pressedTexts = append(pressedTexts, strings.TrimSpace(button.Text()))
	})
	if len(pressedTexts) != 1 || pressedTexts[0] != "Festivalen" {
		t.Fatalf("expected only 'Festivalen' to be pressed, got %v", pressedTexts)
	}
	selectedClassCount := doc.Find(".feedback-category-button.selected").Length()
	if selectedClassCount != 1 {
		t.Fatalf("expected exactly one selected category button, got %d", selectedClassCount)
	}
}

func TestFeedbackThankYou_RendersThankYouMessageAndResetButton(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en tilbakemelding nettopp ble lagret.",
		When:  "Når takke-fragmentet rendres.",
		Then:  "Så skal det vise en takkemelding og en knapp for å sende en ny.",
	})

	// Given
	expectedHeading := "Takk for tilbakemeldingen!"

	// When
	doc := templtest.Render(t, feedbackThankYou())

	// Then
	if headings := templtest.CollectTexts(doc, "h2"); len(headings) == 0 || headings[0] != expectedHeading {
		t.Fatalf("expected heading %q", expectedHeading)
	}
	if doc.Find(`#feedback-form`).Length() != 1 {
		t.Fatalf("expected thank-you fragment to reuse the #feedback-form id so it can replace the form")
	}
	if templtest.CollectTexts(doc, "button")[0] != "Send en tilbakemelding til" {
		t.Fatalf("expected a button to send another feedback")
	}
}
