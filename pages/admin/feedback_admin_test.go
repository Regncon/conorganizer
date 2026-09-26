package admin

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/feedback"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestFeedbackAdminPage_RendersEntriesNewestFirst(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt tilbakemeldinger lagret på ulike tidspunkt.",
		When:  "Når admin-listen rendres.",
		Then:  "Så vises den nyeste tilbakemeldingen først.",
	})

	// Given
	expectedMessages := []string{"nyest", "midten", "eldst"}
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_admin_newest_first")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "midten", "2026-09-20T12:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "nyest", "2026-09-21T12:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "eldst", "2026-09-19T12:00:00.000Z")
	entries := mustListFeedback(t, db, "")

	// When
	doc := templtest.Render(t, feedbackAdminPage(entries, ""))
	actualMessages := templtest.CollectTexts(doc, ".feedback-entry-message")

	// Then
	if len(actualMessages) != len(expectedMessages) {
		t.Fatalf("expected %d messages, got %d: %v", len(expectedMessages), len(actualMessages), actualMessages)
	}
	for i, expectedMessage := range expectedMessages {
		if actualMessages[i] != expectedMessage {
			t.Fatalf("message order mismatch at %d\nexpected: %v\nactual:   %v", i, expectedMessages, actualMessages)
		}
	}
}

func TestFeedbackAdminPage_FiltersByCategoryAndMarksActiveTab(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt tilbakemeldinger i flere kategorier.",
		When:  "Når admin-listen filtreres på én kategori via ?om=.",
		Then:  "Så vises kun den kategoriens tilbakemeldinger, og fanen markeres som aktiv.",
	})

	// Given
	expectedMessage := "om festivalen"
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_admin_filter")
	insertFeedbackAdminEntry(t, db, feedback.CategoryWebsite, "om nettsiden", "2026-09-20T12:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryConvention, expectedMessage, "2026-09-20T13:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "om noe annet", "2026-09-20T14:00:00.000Z")

	activeCategory := resolveFeedbackCategoryFilter("festivalen")
	entries := mustListFeedback(t, db, activeCategory)

	// When
	doc := templtest.Render(t, feedbackAdminPage(entries, activeCategory))
	actualMessages := templtest.CollectTexts(doc, ".feedback-entry-message")
	activeTabText := templtest.CollectTexts(doc, `a[aria-current="page"]`)

	// Then
	if len(actualMessages) != 1 || actualMessages[0] != expectedMessage {
		t.Fatalf("expected only %q, got %v", expectedMessage, actualMessages)
	}
	if len(activeTabText) != 1 || activeTabText[0] != "Festivalen" {
		t.Fatalf("expected only the Festivalen tab to be marked active, got %v", activeTabText)
	}
}

func TestResolveFeedbackCategoryFilter_FallsBackToAllForEmptyOrUnknownSlug(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tom eller ukjent ?om=-verdi.",
		When:  "Når filteret slås opp.",
		Then:  "Så vises alle kategorier, altså ingen filtrering.",
	})

	// Given
	unresolvableSlugs := []string{"", "unknown", "Nettsiden"}

	for _, slug := range unresolvableSlugs {
		// When
		actual := resolveFeedbackCategoryFilter(slug)

		// Then
		if actual != "" {
			t.Fatalf("expected slug %q to resolve to no filter, got %q", slug, actual)
		}
	}
}

func TestFeedbackAdminPage_RendersEmptyState(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at det ikke finnes noen tilbakemeldinger.",
		When:  "Når admin-listen rendres.",
		Then:  "Så vises en tom-tilstand i stedet for en liste.",
	})

	// Given
	expectedEmptyText := "Ingen tilbakemeldinger å vise ennå."

	// When
	doc := templtest.Render(t, feedbackAdminPage(nil, ""))
	actualEmptyText := templtest.CollectTexts(doc, ".feedback-empty")
	actualEntryCount := doc.Find(".feedback-entry").Length()

	// Then
	if len(actualEmptyText) != 1 || actualEmptyText[0] != expectedEmptyText {
		t.Fatalf("expected empty state text %q, got %v", expectedEmptyText, actualEmptyText)
	}
	if actualEntryCount != 0 {
		t.Fatalf("expected no rendered entries, got %d", actualEntryCount)
	}
}

func TestFormatFeedbackTime_ShowsOsloLocalTime(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tilbakemelding lagret med UTC-tidspunkt om sommeren.",
		When:  "Når tidspunktet formateres for admin-listen.",
		Then:  "Så vises det i norsk tid (Europe/Oslo).",
	})

	// Given
	expected := "20.09.2026 14:00"
	createdAt := models.NewDBDateTime(time.Date(2026, time.September, 20, 12, 0, 0, 0, time.UTC))

	// When
	actual := formatFeedbackTime(createdAt)

	// Then
	if actual != expected {
		t.Fatalf("expected Oslo time %q, got %q", expected, actual)
	}
}

func TestFeedbackAdminRoute_NonAdminIsForbiddenAndSeesNoFeedback(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som ikke er admin, og lagrede tilbakemeldinger.",
		When:  "Når brukeren åpner /admin/tilbakemeldinger/ gjennom admin-kravet.",
		Then:  "Så skal tilgang nektes, og ingen tilbakemeldinger vises.",
	})

	// Given
	expectedStatus := http.StatusForbidden
	secretMessage := "hemmelig tilbakemelding"
	db, logger := testutil.CreateTestDBAndLogger(t, "feedback_admin_non_admin")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, secretMessage, "2026-09-20T12:00:00.000Z")
	router := chi.NewRouter()
	adminRouter := router.With(
		userctx.UserMiddleware(logger, db),
		authctx.RequireAdmin(logger, authctx.WithForbiddenHandler(userctx.AdminForbiddenHandler(db, logger))),
	)
	adminRouter.Route("/admin", func(r chi.Router) {
		feedbackAdminRoute(r, db, logger)
	})
	request := httptest.NewRequest(http.MethodGet, "/admin/tilbakemeldinger/", nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), "feedback-non-admin", "non-admin@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
	if strings.Contains(recorder.Body.String(), secretMessage) {
		t.Fatalf("expected no feedback in the forbidden response, got: %s", recorder.Body.String())
	}
}

func insertFeedbackAdminEntry(t *testing.T, db *sql.DB, category feedback.Category, message, createdAt string) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT INTO feedback (category, message, created_at) VALUES (?, ?, ?)`, string(category), message, createdAt)
}

func mustListFeedback(t *testing.T, db *sql.DB, category feedback.Category) []feedback.Entry {
	t.Helper()
	entries, err := feedback.List(db, category)
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	return entries
}
