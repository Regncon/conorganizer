package admin

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/feedback"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/qrtest"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestFeedbackAdminPage_LinksToPrintableQRCode(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt admin-listen over tilbakemeldinger.",
		When:  "Når siden rendres.",
		Then:  "Så finnes en knapp som åpner utskriftsvennlige QR-koder til tilbakemeldingsskjemaet.",
	})

	// Given
	expectedHref := "/admin/tilbakemeldinger/qr"

	// When
	doc := templtest.Render(t, feedbackAdminPage(nil, ""))

	// Then
	link := doc.Find(`a[href="` + expectedHref + `"]`)
	if link.Length() != 1 || !strings.Contains(link.Text(), "QR-kode") {
		t.Fatalf("expected one QR-kode button linking to %s, got %d", expectedHref, link.Length())
	}
}

func TestFeedbackQRRoute_ScansToTheFeedbackFormOnThePublicHost(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en administrator bak proxyen på regncon.no (https).",
		When:  "Når QR-siden for tilbakemeldinger åpnes uten valgt kategori.",
		Then:  "Så skannes QR-koden til https://regncon.no/tilbakemelding, og samme URL står under koden.",
	})

	// Given
	expectedURL := "https://regncon.no/tilbakemelding"

	// When
	doc := getFeedbackQRPage(t, "/admin/tilbakemeldinger/qr", "regncon.no", "https")

	// Then
	assertFeedbackQR(t, doc, expectedURL)
}

func TestFeedbackQRRoute_PreselectsTheChosenCategory(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at administratoren velger kategorien Festivalen på QR-siden.",
		When:  "Når QR-siden åpnes med ?om=festivalen.",
		Then:  "Så åpner QR-koden skjemaet med Festivalen forhåndsvalgt, og valget er markert.",
	})

	// Given
	expectedURL := "https://regncon.no/tilbakemelding?om=festivalen"

	// When
	doc := getFeedbackQRPage(t, "/admin/tilbakemeldinger/qr?om=festivalen", "regncon.no", "https")

	// Then
	assertFeedbackQR(t, doc, expectedURL)
	current := doc.Find(`.feedback-qr-choices a[aria-current="page"]`)
	if current.Length() != 1 || strings.TrimSpace(current.Text()) != feedback.CategoryConvention.Label() {
		t.Fatalf("expected %q to be marked as chosen, got %q", feedback.CategoryConvention.Label(), current.Text())
	}
}

func TestFeedbackQRRoute_UnknownCategoryFallsBackToThePlainForm(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en ukjent kategori i URL-en.",
		When:  "Når QR-siden åpnes med ?om=tull.",
		Then:  "Så peker QR-koden til skjemaet uten kategori.",
	})

	// Given
	expectedURL := "https://regncon.no/tilbakemelding"

	// When
	doc := getFeedbackQRPage(t, "/admin/tilbakemeldinger/qr?om=tull", "regncon.no", "https")

	// Then
	assertFeedbackQR(t, doc, expectedURL)
}

func TestFeedbackQRRoute_UsesPlainHTTPWithoutAProxy(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lokal server uten TLS og uten proxy-header.",
		When:  "Når QR-siden åpnes.",
		Then:  "Så bruker QR-koden http og vertsnavnet fra forespørselen.",
	})

	// Given
	expectedURL := "http://localhost:7331/tilbakemelding"

	// When
	doc := getFeedbackQRPage(t, "/admin/tilbakemeldinger/qr", "localhost:7331", "")

	// Then
	assertFeedbackQR(t, doc, expectedURL)
}

func TestFeedbackQRRoute_NonAdminIsForbidden(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som ikke er admin.",
		When:  "Når brukeren åpner QR-siden gjennom admin-kravet.",
		Then:  "Så nektes tilgang.",
	})

	// Given
	expectedStatus := http.StatusForbidden
	db, logger := testutil.CreateTestDBAndLogger(t, "feedback_qr_non_admin")
	router := chi.NewRouter()
	adminRouter := router.With(
		userctx.UserMiddleware(logger, db),
		authctx.RequireAdmin(logger, authctx.WithForbiddenHandler(userctx.AdminForbiddenHandler(db, logger))),
	)
	adminRouter.Route("/admin", func(r chi.Router) {
		feedbackAdminRoute(r, db, logger)
	})
	request := httptest.NewRequest(http.MethodGet, "/admin/tilbakemeldinger/qr", nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), "feedback-non-admin", "non-admin@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
}

func getFeedbackQRPage(t *testing.T, target, host, forwardedProto string) *goquery.Document {
	t.Helper()
	db, logger := testutil.CreateTestDBAndLogger(t, "feedback_qr")
	router := chi.NewRouter()
	router.Route("/admin", func(r chi.Router) {
		feedbackAdminRoute(r, db, logger)
	})
	request := httptest.NewRequest(http.MethodGet, target, nil)
	request.Host = host
	if forwardedProto != "" {
		request.Header.Set("X-Forwarded-Proto", forwardedProto)
	}
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	doc, err := goquery.NewDocumentFromReader(recorder.Body)
	if err != nil {
		t.Fatalf("parse QR page: %v", err)
	}
	return doc
}

func assertFeedbackQR(t *testing.T, doc *goquery.Document, expectedURL string) {
	t.Helper()
	svg, err := goquery.OuterHtml(doc.Find(".feedback-qr-code svg"))
	if err != nil || svg == "" {
		t.Fatalf("expected a QR code SVG on the page: %v", err)
	}
	if got := qrtest.ScanSVG(t, svg); got != expectedURL {
		t.Fatalf("QR code scans to %q, want %q", got, expectedURL)
	}
	if got := strings.TrimSpace(doc.Find(".feedback-qr-url").Text()); got != expectedURL {
		t.Fatalf("printed URL = %q, want %q", got, expectedURL)
	}
}
