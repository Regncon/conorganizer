package feedback

import (
	"database/sql"
	"encoding/json"
	"maps"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/feedback"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
)

func TestFeedbackRoute_PostStoresFeedbackAndReturnsThankYou(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker med en gyldig kategori og melding.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal den lagres, og skjemaet skal erstattes med en takkemelding.",
	})

	// Given
	expectedMessage := "Påmeldingen var enkel å finne."
	db, router := setupFeedbackRouteTest(t, "feedback_post_success")

	// When
	recorder := postFeedbackSignals(t, router, map[string]string{
		"feedbackCategory": string(feedback.CategoryWebsite),
		"feedbackMessage":  expectedMessage,
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
	if len(entries) != 1 || entries[0].Message != expectedMessage || entries[0].Category != feedback.CategoryWebsite {
		t.Fatalf("expected exactly one stored entry matching the submission, got: %+v", entries)
	}
}

func TestFeedbackRoute_PostSuccessClearsEarlierValidationErrors(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en bruker som tidligere fikk en valideringsfeil.",
		When:  "Når en gyldig tilbakemelding sendes inn.",
		Then:  "Så skal feilmeldingene tømmes, slik at de ikke dukker opp igjen i et nytt skjema.",
	})

	// Given
	expectedErrors := map[string]string{"feedbackCategory": "", "feedbackMessage": ""}
	_, router := setupFeedbackRouteTest(t, "feedback_post_clears_errors")

	// When
	recorder := postFeedbackSignals(t, router, map[string]string{
		"feedbackCategory": string(feedback.CategoryConvention),
		"feedbackMessage":  "Fin festival.",
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
	expectedErrors := map[string]string{"feedbackCategory": "", "feedbackMessage": ""}
	_, router := setupFeedbackRouteTest(t, "feedback_reset_clears_errors")

	// When
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/tilbakemelding/reset", nil))

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", expectedErrors)
}

func TestFeedbackRoute_PostEmptyMessageShowsErrorAndStoresNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som sender inn en tom melding.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal en norsk feilmelding vises ved feltet, og ingenting skal lagres.",
	})

	// Given
	expectedError := "Skriv en tilbakemelding før du sender inn."
	db, router := setupFeedbackRouteTest(t, "feedback_post_empty")

	// When
	recorder := postFeedbackSignals(t, router, map[string]string{
		"feedbackCategory": string(feedback.CategoryOther),
		"feedbackMessage":  "   ",
	})

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", map[string]string{
		"feedbackCategory": "",
		"feedbackMessage":  expectedError,
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
	recorder := postFeedbackSignals(t, router, map[string]string{
		"feedbackCategory": "",
		"feedbackMessage":  "Noe fornuftig å si.",
	})

	// Then
	assertFeedbackSignal(t, recorder.Body.String(), "feedbackErrors", map[string]string{
		"feedbackCategory": expectedError,
		"feedbackMessage":  "",
	})
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM feedback`); count != 0 {
		t.Fatalf("expected no stored feedback, got %d", count)
	}
}

func TestFeedbackRoute_PostTooLongMessageShowsErrorAndStoresNothing(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en melding som er lengre enn maks antall tegn.",
		When:  "Når tilbakemeldingen sendes inn.",
		Then:  "Så skal en norsk feilmelding vises ved feltet, og ingenting skal lagres.",
	})

	// Given
	tooLong := strings.Repeat("a", feedback.MaxMessageLength+1)
	db, router := setupFeedbackRouteTest(t, "feedback_post_too_long")

	// When
	recorder := postFeedbackSignals(t, router, map[string]string{
		"feedbackCategory": string(feedback.CategoryWebsite),
		"feedbackMessage":  tooLong,
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
	body, err := json.Marshal(map[string]string{
		"feedbackCategory": string(feedback.CategoryWebsite),
		"feedbackMessage":  "Skal ikke lagres.",
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

func postFeedbackSignals(t *testing.T, router http.Handler, signals map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(signals)
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
		if values, ok := patch[signalKey]; ok {
			merged := map[string]string{}
			maps.Copy(merged, values)
			for key, expectedValue := range expected {
				if merged[key] != expectedValue {
					t.Fatalf("signal %s.%s mismatch\nexpected: %q\nactual:   %q\nbody: %s", signalKey, key, expectedValue, merged[key], body)
				}
			}
			return
		}
	}
	t.Fatalf("expected a patch for signal %q\nbody: %s", signalKey, body)
}
