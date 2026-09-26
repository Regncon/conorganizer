package feedback

import (
	"database/sql"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/feedback"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
)

func TestFeedbackRoute_PostWithOnlyWentWellStoresIt(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som bare fyller ut 'Hva fungerte bra?'.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal den lagres, og skjemaet skal erstattes med en takkemelding.",
	})

	// Given
	expectedWentWell := "Påmeldingen var enkel å finne."
	db, router := setupFeedbackRouteTest(t, "feedback_post_went_well_only")

	// When
	recorder := postFeedbackSignals(t, router, feedbackSubmission{
		Category: string(feedback.CategoryWebsite),
		WentWell: expectedWentWell,
	})

	// Then
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	elements := feedbackPatchedElements(t, recorder.Body.String())
	if !strings.Contains(elements, "Takk for tilbakemeldingen!") {
		t.Fatalf("expected the form to be replaced by a thank-you message, got: %s", elements)
	}
	if !strings.Contains(elements, `id="feedback-form"`) {
		t.Fatalf("expected the thank-you fragment to reuse the #feedback-form id, got: %s", elements)
	}
	entries := mustListFeedback(t, db, "")
	if len(entries) != 1 || entries[0].WentWell != expectedWentWell || entries[0].CouldImprove != "" || entries[0].Category != feedback.CategoryWebsite {
		t.Fatalf("expected exactly one stored entry matching the submission, got: %+v", entries)
	}
}

func TestFeedbackRoute_PostWithTopicsStoresThem(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som velger tema i tillegg til tekst.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal de valgte temaene lagres sammen med tilbakemeldingen.",
	})

	// Given
	expectedTopics := []feedback.Topic{feedback.TopicSignup, feedback.TopicMobile}
	db, router := setupFeedbackRouteTest(t, "feedback_post_with_topics")

	// When
	recorder := postFeedbackSignals(t, router, feedbackSubmission{
		Category: string(feedback.CategoryWebsite),
		WentWell: "Registreringen gikk fint.",
		Topics:   []string{string(feedback.TopicSignup), string(feedback.TopicMobile)},
	})

	// Then
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d\nbody: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	entries := mustListFeedback(t, db, "")
	if len(entries) != 1 {
		t.Fatalf("expected exactly one stored entry, got: %+v", entries)
	}
	if len(entries[0].Topics) != len(expectedTopics) {
		t.Fatalf("expected topics %v, got %v", expectedTopics, entries[0].Topics)
	}
	for i, topic := range expectedTopics {
		if entries[0].Topics[i] != topic {
			t.Fatalf("expected topics %v, got %v", expectedTopics, entries[0].Topics)
		}
	}
}

func TestFeedbackRoute_PostSuccessClearsEarlierValidationErrors(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en bruker som tidligere fikk en valideringsfeil.",
		When:  "Når en gyldig tilbakemelding sendes inn.",
		Then:  "Så skal feilmeldingene tømmes, slik at de ikke dukker opp igjen i et nytt skjema.",
	})

	// Given
	expectedErrors := map[string]string{
		feedbackCategorySignal:     "",
		feedbackWentWellSignal:     "",
		feedbackCouldImproveSignal: "",
		feedbackTopicsSignal:       "",
	}
	_, router := setupFeedbackRouteTest(t, "feedback_post_clears_errors")

	// When
	recorder := postFeedbackSignals(t, router, feedbackSubmission{
		Category: string(feedback.CategoryConvention),
		WentWell: "Fin festival.",
	})

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", expectedErrors)
}

func TestFeedbackRoute_ResetClearsValidationErrors(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en bruker som har sendt inn og ber om et nytt skjema.",
		When:  "Når skjemaet tilbakestilles.",
		Then:  "Så skal feilmeldingene tømmes sammen med det nye skjemaet.",
	})

	// Given
	expectedErrors := map[string]string{
		feedbackCategorySignal:     "",
		feedbackWentWellSignal:     "",
		feedbackCouldImproveSignal: "",
		feedbackTopicsSignal:       "",
	}
	_, router := setupFeedbackRouteTest(t, "feedback_reset_clears_errors")

	// When
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/tilbakemelding/reset", nil))

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", expectedErrors)
}

func TestFeedbackRoute_PostWithBothTextsEmptyShowsRequiredErrorAndStoresNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som lar begge tekstfeltene stå tomme.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal en norsk feilmelding vises, og ingenting skal lagres.",
	})

	// Given
	expectedError := "Fyll ut minst ett av feltene."
	db, router := setupFeedbackRouteTest(t, "feedback_post_both_empty")

	// When
	recorder := postFeedbackSignals(t, router, feedbackSubmission{
		Category:     string(feedback.CategoryOther),
		WentWell:     "   ",
		CouldImprove: "",
	})

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", map[string]string{
		feedbackWentWellSignal: expectedError,
	})
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM feedback`); count != 0 {
		t.Fatalf("expected no stored feedback, got %d", count)
	}
}

func TestFeedbackRoute_PostUnknownCategoryShowsErrorAndStoresNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innsending uten en gyldig kategori.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal en norsk feilmelding vises ved kategorivalget, og ingenting skal lagres.",
	})

	// Given
	expectedError := "Velg hva tilbakemeldingen gjelder."
	db, router := setupFeedbackRouteTest(t, "feedback_post_bad_category")

	// When
	recorder := postFeedbackSignals(t, router, feedbackSubmission{
		Category: "",
		WentWell: "Noe fornuftig å si.",
	})

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", map[string]string{
		feedbackCategorySignal: expectedError,
	})
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM feedback`); count != 0 {
		t.Fatalf("expected no stored feedback, got %d", count)
	}
}

func TestFeedbackRoute_PostTooLongTextShowsErrorAndStoresNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tekst som er lengre enn maks antall tegn.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal en norsk feilmelding vises ved feltet, og ingenting skal lagres.",
	})

	// Given
	tooLong := strings.Repeat("a", feedback.MaxTextLength+1)
	db, router := setupFeedbackRouteTest(t, "feedback_post_too_long")

	// When
	recorder := postFeedbackSignals(t, router, feedbackSubmission{
		Category: string(feedback.CategoryWebsite),
		WentWell: tooLong,
	})

	// Then
	body := recorder.Body.String()
	if !strings.Contains(body, "for lang") {
		t.Fatalf("expected a 'too long' validation message, got: %s", body)
	}
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM feedback`); count != 0 {
		t.Fatalf("expected no stored feedback, got %d", count)
	}
}

func TestFeedbackRoute_PostWithEmailShowsPersonalInfoErrorAndStoresNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tilbakemelding som inneholder en e-postadresse.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal en norsk feilmelding om personopplysninger vises, og ingenting skal lagres.",
	})

	// Given
	expectedError := "Det ser ut som du har skrevet en e-postadresse eller et telefonnummer. Fjern det før du sender."
	db, router := setupFeedbackRouteTest(t, "feedback_post_email")

	// When
	recorder := postFeedbackSignals(t, router, feedbackSubmission{
		Category: string(feedback.CategoryWebsite),
		WentWell: "Kontakt meg på kari@example.no for detaljer.",
	})

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", map[string]string{
		feedbackWentWellSignal: expectedError,
	})
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM feedback`); count != 0 {
		t.Fatalf("expected no stored feedback, got %d", count)
	}
}

func TestFeedbackRoute_PostWithPhoneNumberShowsPersonalInfoErrorAndStoresNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tilbakemelding som inneholder et telefonnummer.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal en norsk feilmelding om personopplysninger vises, og ingenting skal lagres.",
	})

	// Given
	expectedError := "Det ser ut som du har skrevet en e-postadresse eller et telefonnummer. Fjern det før du sender."
	db, router := setupFeedbackRouteTest(t, "feedback_post_phone")

	// When
	recorder := postFeedbackSignals(t, router, feedbackSubmission{
		Category:     string(feedback.CategoryWebsite),
		CouldImprove: "Ring meg på 12345678 gjerne.",
	})

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", map[string]string{
		feedbackCouldImproveSignal: expectedError,
	})
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM feedback`); count != 0 {
		t.Fatalf("expected no stored feedback, got %d", count)
	}
}

func TestFeedbackRoute_UnauthenticatedGetShowsLoginWithNesteLink(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en ikke-innlogget bruker som åpner tilbakemeldingssiden med en forhåndsvalgt kategori.",
		When:  "Når siden hentes.",
		Then:  "Så skal innloggingssiden vises, med en retur-lenke tilbake til akkurat den siden.",
	})

	// Given
	expectedLoginHref := "/auth?neste=" + url.QueryEscape("/tilbakemelding?om=festivalen")
	_, router := setupFeedbackRouteTest(t, "feedback_get_unauthenticated")

	// When
	request := httptest.NewRequest(http.MethodGet, "/tilbakemelding?om=festivalen", nil)
	request.Header.Set("X-Test-Anonymous", "true")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), expectedLoginHref) {
		t.Fatalf("expected login page to contain return link %q\nbody: %s", expectedLoginHref, recorder.Body.String())
	}
}

func TestFeedbackRoute_UnauthenticatedPostIsRejectedAndStoresNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en ikke-innlogget bruker med gyldige skjemadata.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal forespørselen avvises, og ingenting skal lagres.",
	})

	// Given
	expectedStatus := http.StatusUnauthorized
	db, router := setupFeedbackRouteTest(t, "feedback_post_unauthenticated")
	body, err := json.Marshal(feedbackSubmission{
		Category: string(feedback.CategoryWebsite),
		WentWell: "Skal ikke lagres.",
	})
	if err != nil {
		t.Fatalf("failed to marshal Datastar signals: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/tilbakemelding", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Test-Anonymous", "true")
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM feedback`); count != 0 {
		t.Fatalf("expected no stored feedback, got %d", count)
	}
}

func TestFeedbackRoute_GetWithoutOrUnknownOmPreselectsWebsite(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som åpner tilbakemeldingssiden uten ?om= eller med en ukjent verdi.",
		When:  "Når siden hentes.",
		Then:  "Så skal 'Nettsiden' være forhåndsvalgt.",
	})

	_, router := setupFeedbackRouteTest(t, "feedback_get_default_category")
	for _, target := range []string{"/tilbakemelding", "/tilbakemelding?om=ukjent"} {
		// Given
		expectedPressed := []string{"Nettsiden"}

		// When
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, target, nil))

		// Then
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected status %d for %s, got %d", http.StatusOK, target, recorder.Code)
		}
		assertPressedCategories(t, recorder.Body.String(), expectedPressed)
	}
}

func TestFeedbackRoute_ResetWithOmPreselectsThatCategory(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en bruker som sendte en tilbakemelding om festivalen.",
		When:  "Når skjemaet tilbakestilles med ?om=festivalen.",
		Then:  "Så skal det nye skjemaet ha 'Festivalen' forhåndsvalgt.",
	})

	// Given
	expectedPressed := []string{"Festivalen"}
	_, router := setupFeedbackRouteTest(t, "feedback_reset_keeps_category")

	// When
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/tilbakemelding/reset?om=festivalen", nil))

	// Then
	assertPressedCategories(t, feedbackPatchedElements(t, recorder.Body.String()), expectedPressed)
}

func assertPressedCategories(t *testing.T, html string, expected []string) {
	t.Helper()
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("failed to parse HTML: %v", err)
	}
	pressed := []string{}
	doc.Find(`.feedback-category-button[aria-pressed="true"]`).Each(func(_ int, button *goquery.Selection) {
		pressed = append(pressed, strings.TrimSpace(button.Text()))
	})
	if !slices.Equal(pressed, expected) {
		t.Fatalf("expected pressed categories %v, got %v", expected, pressed)
	}
}

func setupFeedbackRouteTest(t *testing.T, dbName string) (*sql.DB, chi.Router) {
	t.Helper()

	db, logger := testutil.CreateTestDBAndLogger(t, dbName)
	router := chi.NewRouter()
	isLoggedInRouter := router.With(withFeedbackTestUser, userctx.UserMiddleware(logger, db))
	if err := SetupFeedbackRoute(isLoggedInRouter, db, logger); err != nil {
		t.Fatalf("expected feedback route setup to succeed: %v", err)
	}
	return db, router
}

// withFeedbackTestUser simulates a logged-in user for every request except
// the unauthenticated GET test, which deliberately skips it by calling the
// router directly without adding this context.
func withFeedbackTestUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Test-Anonymous") == "true" {
			next.ServeHTTP(w, r)
			return
		}
		ctx := authctx.WithUserToken(r.Context(), "feedback-route-user", "feedback-route-user@example.com")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func postFeedbackSignals(t *testing.T, router http.Handler, submission feedbackSubmission) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(submission)
	if err != nil {
		t.Fatalf("failed to marshal Datastar signals: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, "/tilbakemelding", strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func mustListFeedback(t *testing.T, db *sql.DB, category feedback.Category) []feedback.Entry {
	t.Helper()
	entries, err := feedback.List(db, category)
	if err != nil {
		t.Fatalf("expected feedback list to succeed: %v", err)
	}
	return entries
}

func feedbackPatchedElements(t *testing.T, body string) string {
	t.Helper()
	var elements []string
	for line := range strings.SplitSeq(body, "\n") {
		if payload, ok := strings.CutPrefix(line, "data: elements "); ok {
			elements = append(elements, payload)
		}
	}
	return strings.Join(elements, "\n")
}

// assertFeedbackSignal checks that the feedbackErrors patch contains the
// expected messages. Keys not mentioned in expected are still required to be
// present and cleared (empty), since Patch always sends every initialized key.
func assertFeedbackSignal(t *testing.T, body string, signalKey string, expected map[string]string) {
	t.Helper()
	for line := range strings.SplitSeq(body, "\n") {
		payload, ok := strings.CutPrefix(line, "data: signals ")
		if !ok {
			continue
		}
		patch := map[string]map[string]string{}
		if err := json.Unmarshal([]byte(payload), &patch); err != nil {
			t.Fatalf("failed to unmarshal Datastar signal patch %q: %v", payload, err)
		}
		values, ok := patch[signalKey]
		if !ok {
			continue
		}
		merged := map[string]string{}
		maps.Copy(merged, values)
		for _, key := range []string{feedbackCategorySignal, feedbackWentWellSignal, feedbackCouldImproveSignal, feedbackTopicsSignal} {
			if _, wanted := expected[key]; !wanted {
				expected[key] = ""
			}
		}
		for key, expectedValue := range expected {
			if merged[key] != expectedValue {
				t.Fatalf("signal %s.%s mismatch\nexpected: %q\nactual:   %q\nbody: %s", signalKey, key, expectedValue, merged[key], body)
			}
		}
		return
	}
	t.Fatalf("expected a patch for signal %q\nbody: %s", signalKey, body)
}
