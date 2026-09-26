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

func TestFeedbackFormPage_DefaultsToWebsiteCategoryWithoutOmParam(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en innlogget bruker åpner tilbakemeldingssiden uten ?om=.",
		When:  "Når siden rendres.",
		Then:  "Så skal 'Nettsiden' være forhåndsvalgt, og bare de andre skal være utrykket.",
	})

	// Given
	expectedPressedLabel := "Nettsiden"
	expectedLabels := []string{"Nettsiden", "Festivalen", "Noe annet"}

	// When
	doc := templtest.Render(t, feedbackFormPage(feedback.CategoryWebsite))
	buttons := doc.Find(".feedback-category-button")

	// Then
	if buttons.Length() != len(expectedLabels) {
		t.Fatalf("expected %d category buttons, got %d", len(expectedLabels), buttons.Length())
	}
	pressedCount := 0
	for i, expectedLabel := range expectedLabels {
		button := buttons.Eq(i)
		actualLabel := strings.TrimSpace(button.Text())
		if actualLabel != expectedLabel {
			t.Fatalf("expected category button %d to be %q, got %q", i, expectedLabel, actualLabel)
		}
		pressed, _ := button.Attr("aria-pressed")
		if pressed == "true" {
			pressedCount++
			if actualLabel != expectedPressedLabel {
				t.Fatalf("expected %q to be the only pressed category, got %q pressed", expectedPressedLabel, actualLabel)
			}
		}
	}
	if pressedCount != 1 {
		t.Fatalf("expected exactly one pressed category button, got %d", pressedCount)
	}
	if doc.Find("textarea#feedback-went-well").AttrOr("maxlength", "") != strconv.Itoa(feedback.MaxTextLength) {
		t.Fatalf("expected 'went well' textarea maxlength %d", feedback.MaxTextLength)
	}
	if doc.Find("textarea#feedback-could-improve").AttrOr("maxlength", "") != strconv.Itoa(feedback.MaxTextLength) {
		t.Fatalf("expected 'could improve' textarea maxlength %d", feedback.MaxTextLength)
	}
	if templtest.CollectTexts(doc, "button[type=submit]")[0] != "Send tilbakemelding" {
		t.Fatalf("expected submit button labelled 'Send tilbakemelding'")
	}
}

func TestFeedbackFormPage_RendersSmallAnonymityHint(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at tilbakemeldingssiden rendres.",
		When:  "Når siden vises.",
		Then:  "Så skal den vise en kort, rolig anonymitetslinje, ikke en stor advarselsboks.",
	})

	// Given
	expectedHintPart := "Tilbakemeldingen er anonym"

	// When
	doc := templtest.Render(t, feedbackFormPage(feedback.CategoryWebsite))
	actualHint := strings.Join(templtest.CollectTexts(doc, ".feedback-privacy-hint"), " ")

	// Then
	if !strings.Contains(actualHint, expectedHintPart) {
		t.Fatalf("expected anonymity hint to contain %q, got %q", expectedHintPart, actualHint)
	}
	if doc.Find(".feedback-privacy-hint .inline-icon").Length() != 1 {
		t.Fatalf("expected the anonymity hint to include a small icon")
	}
}

func TestFeedbackFormPage_PreselectsCategoryFromOmParamAndShowsOnlyItsTopics(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren kom fra en QR-kode med ?om=festivalen.",
		When:  "Når tilbakemeldingssiden rendres med den forhåndsvalgte kategorien.",
		Then:  "Så skal 'Festivalen' vises som trykket, og bare temaene for den kategorien skal vises.",
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

	visibleGroupSelector := `.feedback-topics[data-show="$feedbackCategory === 'convention'"]`
	convChips := templtest.CollectTexts(doc, visibleGroupSelector+" .feedback-topic-chip")
	if len(convChips) != len(feedback.TopicsFor(feedback.CategoryConvention)) {
		t.Fatalf("expected %d convention topic chips, got %v", len(feedback.TopicsFor(feedback.CategoryConvention)), convChips)
	}
	for _, topic := range feedback.TopicsFor(feedback.CategoryConvention) {
		if !strings.Contains(strings.Join(convChips, "|"), topic.Label()) {
			t.Fatalf("expected convention topic chip %q, got %v", topic.Label(), convChips)
		}
	}
	websiteChipsSelector := `.feedback-topics[data-show="$feedbackCategory === 'website'"] .feedback-topic-chip`
	if doc.Find(websiteChipsSelector).Length() != len(feedback.TopicsFor(feedback.CategoryWebsite)) {
		t.Fatalf("expected website topic group to still render its own chips (shown/hidden client-side)")
	}
}

func TestFeedbackFormPage_ExposesPersonalInfoPatternValuesForTheScript(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at tilbakemeldingssiden rendres.",
		When:  "Når skjemaet leter etter personopplysninger i nettleseren.",
		Then:  "Så skal skjemaet eksponere de samme mønstrene som backend bruker.",
	})

	// Given
	expectedEmailMarker := feedback.EmailMarker
	expectedDigitPattern := feedback.DigitRunPattern

	// When
	doc := templtest.Render(t, feedbackFormPage(feedback.CategoryWebsite))
	form := doc.Find("form").First()

	// Then
	if got := form.AttrOr("data-personal-info-email-marker", ""); got != expectedEmailMarker {
		t.Fatalf("expected email marker %q, got %q", expectedEmailMarker, got)
	}
	if got := form.AttrOr("data-personal-info-digit-pattern", ""); got != expectedDigitPattern {
		t.Fatalf("expected digit pattern %q, got %q", expectedDigitPattern, got)
	}
	if got := form.AttrOr("data-personal-info-date-pattern", ""); got != feedback.DateOrTimeRangePattern {
		t.Fatalf("expected date pattern %q, got %q", feedback.DateOrTimeRangePattern, got)
	}
	if got := form.AttrOr("data-personal-info-date-replacement", ""); got != feedback.DateOrTimeRangeReplacement {
		t.Fatalf("expected date replacement %q, got %q", feedback.DateOrTimeRangeReplacement, got)
	}
	digitCounts := form.AttrOr("data-personal-info-digit-counts", "")
	for _, count := range feedback.BannedDigitCounts {
		if !strings.Contains(digitCounts, strconv.Itoa(count)) {
			t.Fatalf("expected banned digit counts %q to contain %d", digitCounts, count)
		}
	}
	if doc.Find("script").Length() == 0 {
		t.Fatalf("expected an inline script implementing the personal info check")
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
	doc := templtest.Render(t, feedbackThankYou(feedback.CategoryWebsite))

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

func TestFeedbackFormPage_HidesTopicsQuestionWhenCategoryHasNoTopics(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at 'Noe annet' er forhåndsvalgt, som ikke har noen temaer.",
		When:  "Når tilbakemeldingssiden rendres.",
		Then:  "Så skal spørsmålet om tema være skjult, og bare vises for kategorier som har temaer.",
	})

	// Given
	expectedShowExpr := "['website', 'convention'].includes($feedbackCategory)"

	// When
	otherDoc := templtest.Render(t, feedbackFormPage(feedback.CategoryOther))
	websiteDoc := templtest.Render(t, feedbackFormPage(feedback.CategoryWebsite))

	// Then
	otherTopics := otherDoc.Find("fieldset.feedback-topics-field")
	if got := otherTopics.AttrOr("data-show", ""); got != expectedShowExpr {
		t.Fatalf("expected topics fieldset data-show %q, got %q", expectedShowExpr, got)
	}
	if got := otherTopics.AttrOr("style", ""); got != "display: none" {
		t.Fatalf("expected topics fieldset to start hidden for 'Noe annet', got style %q", got)
	}
	if _, hasStyle := websiteDoc.Find("fieldset.feedback-topics-field").Attr("style"); hasStyle {
		t.Fatalf("expected topics fieldset to start visible for 'Nettsiden'")
	}
}

func TestFeedbackFormPage_DescribesTextFieldsAndSubmitButtonForAssistiveTech(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at tilbakemeldingssiden rendres.",
		When:  "Når en skjermleser leser tekstfeltene og send-knappen.",
		Then:  "Så skal feltene peke på feilmelding og tegnteller, og personvernhintet og knappen skal kobles til når hintet vises.",
	})

	for _, fieldKey := range []string{"went-well", "could-improve"} {
		// Given
		expectedDescribedBy := "feedback-" + fieldKey + "-error feedback-" + fieldKey + "-counter"
		expectedHintID := "feedback-" + fieldKey + "-personal-info"

		// When
		doc := templtest.Render(t, feedbackFormPage(feedback.CategoryWebsite))
		textarea := doc.Find("textarea#feedback-" + fieldKey)

		// Then
		if got := textarea.AttrOr("aria-describedby", ""); got != expectedDescribedBy {
			t.Fatalf("expected textarea aria-describedby %q, got %q", expectedDescribedBy, got)
		}
		for _, id := range strings.Fields(expectedDescribedBy + " " + expectedHintID) {
			if doc.Find("#"+id).Length() != 1 {
				t.Fatalf("expected exactly one element with id %q", id)
			}
		}
		if !strings.Contains(textarea.AttrOr("data-attr:aria-describedby", ""), expectedHintID) {
			t.Fatalf("expected textarea to reference %q while the hint is shown", expectedHintID)
		}
		if !strings.Contains(doc.Find("button[type=submit]").AttrOr("data-attr:aria-describedby", ""), expectedHintID) {
			t.Fatalf("expected submit button to reference %q while it is disabled by that hint", expectedHintID)
		}
	}
}

func TestFeedbackThankYou_ResetKeepsSubmittedCategory(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren nettopp sendte en tilbakemelding om festivalen.",
		When:  "Når takke-fragmentet rendres.",
		Then:  "Så skal knappen for en ny tilbakemelding hente et skjema med festivalen forhåndsvalgt.",
	})

	// Given
	expectedClick := "@get('/tilbakemelding/reset?om=festivalen')"

	// When
	doc := templtest.Render(t, feedbackThankYou(feedback.CategoryConvention))

	// Then
	if got := doc.Find("button").AttrOr("data-on:click", ""); got != expectedClick {
		t.Fatalf("expected reset click %q, got %q", expectedClick, got)
	}
}
