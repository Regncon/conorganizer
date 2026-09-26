package admin

import (
	"database/sql"
	"encoding/json"
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
	expectedTexts := []string{"nyest", "midten", "eldst"}
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_admin_newest_first")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "midten", "", "2026-09-20T12:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "nyest", "", "2026-09-21T12:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "eldst", "", "2026-09-19T12:00:00.000Z")
	entries := mustListFeedback(t, db, "")

	// When
	doc := templtest.Render(t, feedbackAdminPage(entries, ""))
	actualTexts := templtest.CollectTexts(doc, ".feedback-entry-went-well")

	// Then
	if len(actualTexts) != len(expectedTexts) {
		t.Fatalf("expected %d entries, got %d: %v", len(expectedTexts), len(actualTexts), actualTexts)
	}
	for i, expectedText := range expectedTexts {
		if actualTexts[i] != expectedText {
			t.Fatalf("entry order mismatch at %d\nexpected: %v\nactual:   %v", i, expectedTexts, actualTexts)
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
	expectedText := "om festivalen"
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_admin_filter")
	insertFeedbackAdminEntry(t, db, feedback.CategoryWebsite, "om nettsiden", "", "2026-09-20T12:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryConvention, expectedText, "", "2026-09-20T13:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "om noe annet", "", "2026-09-20T14:00:00.000Z")

	activeCategory := resolveFeedbackCategoryFilter("festivalen")
	entries := mustListFeedback(t, db, activeCategory)

	// When
	doc := templtest.Render(t, feedbackAdminPage(entries, activeCategory))
	actualTexts := templtest.CollectTexts(doc, ".feedback-entry-went-well")
	activeTabText := templtest.CollectTexts(doc, `a[aria-current="page"]`)

	// Then
	if len(actualTexts) != 1 || actualTexts[0] != expectedText {
		t.Fatalf("expected only %q, got %v", expectedText, actualTexts)
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

func TestFeedbackAdminPage_ShowsBothTextsAndTopicsWhenPresent(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tilbakemelding med begge tekstfelt og emner.",
		When:  "Når admin-listen rendres.",
		Then:  "Så vises begge tekstene under sine overskrifter, og emnene som chips.",
	})

	// Given
	expectedWentWell := "Påmeldingen var enkel."
	expectedCouldImprove := "Programmet var tregt på mobil."
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_admin_texts_and_topics")
	insertFeedbackAdminEntryWithTopics(t, db, feedback.CategoryWebsite, expectedWentWell, expectedCouldImprove,
		[]feedback.Topic{feedback.TopicSignup, feedback.TopicMobile}, "2026-09-20T12:00:00.000Z")
	entries := mustListFeedback(t, db, "")

	// When
	doc := templtest.Render(t, feedbackAdminPage(entries, ""))
	actualHeadings := templtest.CollectTexts(doc, ".feedback-entry-heading")
	actualWentWell := templtest.CollectTexts(doc, ".feedback-entry-went-well")
	actualCouldImprove := templtest.CollectTexts(doc, ".feedback-entry-could-improve")
	actualTopics := templtest.CollectTexts(doc, ".feedback-entry-topic")

	// Then
	expectedHeadings := []string{"Fungerte bra", "Kan bli bedre"}
	if len(actualHeadings) != len(expectedHeadings) {
		t.Fatalf("expected headings %v, got %v", expectedHeadings, actualHeadings)
	}
	for i, expectedHeading := range expectedHeadings {
		if actualHeadings[i] != expectedHeading {
			t.Fatalf("expected heading %q at %d, got %q", expectedHeading, i, actualHeadings[i])
		}
	}
	if len(actualWentWell) != 1 || actualWentWell[0] != expectedWentWell {
		t.Fatalf("expected went-well text %q, got %v", expectedWentWell, actualWentWell)
	}
	if len(actualCouldImprove) != 1 || actualCouldImprove[0] != expectedCouldImprove {
		t.Fatalf("expected could-improve text %q, got %v", expectedCouldImprove, actualCouldImprove)
	}
	expectedTopics := []string{"Påmelding", "Mobil"}
	if len(actualTopics) != len(expectedTopics) {
		t.Fatalf("expected topics %v, got %v", expectedTopics, actualTopics)
	}
	for i, expectedTopic := range expectedTopics {
		if actualTopics[i] != expectedTopic {
			t.Fatalf("expected topic %q at %d, got %q", expectedTopic, i, actualTopics[i])
		}
	}
}

func TestFeedbackAdminPage_OmitsEmptyTextSectionAndTopicList(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tilbakemelding med kun ett tekstfelt utfylt og ingen emner.",
		When:  "Når admin-listen rendres.",
		Then:  "Så vises kun den utfylte tekstseksjonen, og ingen emneliste.",
	})

	// Given
	expectedWentWell := "Alt fungerte bra."
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_admin_omits_empty")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, expectedWentWell, "", "2026-09-20T12:00:00.000Z")
	entries := mustListFeedback(t, db, "")

	// When
	doc := templtest.Render(t, feedbackAdminPage(entries, ""))
	actualHeadings := templtest.CollectTexts(doc, ".feedback-entry-heading")
	actualWentWell := templtest.CollectTexts(doc, ".feedback-entry-went-well")
	actualCouldImprove := templtest.CollectTexts(doc, ".feedback-entry-could-improve")
	actualTopicListCount := doc.Find(".feedback-entry-topics").Length()

	// Then
	if len(actualHeadings) != 1 || actualHeadings[0] != "Fungerte bra" {
		t.Fatalf("expected only the went-well heading, got %v", actualHeadings)
	}
	if len(actualWentWell) != 1 || actualWentWell[0] != expectedWentWell {
		t.Fatalf("expected went-well text %q, got %v", expectedWentWell, actualWentWell)
	}
	if len(actualCouldImprove) != 0 {
		t.Fatalf("expected no could-improve section, got %v", actualCouldImprove)
	}
	if actualTopicListCount != 0 {
		t.Fatalf("expected no topic list, got %d", actualTopicListCount)
	}
}

func TestFeedbackAdminPage_PreservesLineBreaksInText(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tilbakemelding med linjeskift i teksten.",
		When:  "Når admin-listen rendres.",
		Then:  "Så beholdes linjeskiftene i den viste teksten.",
	})

	// Given
	expectedWentWell := "Første linje\nAndre linje"
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_admin_line_breaks")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, expectedWentWell, "", "2026-09-20T12:00:00.000Z")
	entries := mustListFeedback(t, db, "")

	// When
	doc := templtest.Render(t, feedbackAdminPage(entries, ""))
	actualWentWell := doc.Find(".feedback-entry-went-well").First().Text()

	// Then
	if actualWentWell != expectedWentWell {
		t.Fatalf("expected line breaks to be preserved, expected %q, got %q", expectedWentWell, actualWentWell)
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
	secretText := "hemmelig tilbakemelding"
	db, logger := testutil.CreateTestDBAndLogger(t, "feedback_admin_non_admin")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, secretText, "", "2026-09-20T12:00:00.000Z")
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
	if strings.Contains(recorder.Body.String(), secretText) {
		t.Fatalf("expected no feedback in the forbidden response, got: %s", recorder.Body.String())
	}
}

func insertFeedbackAdminEntry(t *testing.T, db *sql.DB, category feedback.Category, wentWell, couldImprove, createdAt string) {
	t.Helper()
	insertFeedbackAdminEntryWithTopics(t, db, category, wentWell, couldImprove, nil, createdAt)
}

func insertFeedbackAdminEntryWithTopics(t *testing.T, db *sql.DB, category feedback.Category, wentWell, couldImprove string, topics []feedback.Topic, createdAt string) {
	t.Helper()
	if topics == nil {
		topics = []feedback.Topic{}
	}
	topicsJSON, err := json.Marshal(topics)
	if err != nil {
		t.Fatalf("marshal topics: %v", err)
	}
	testutil.MustExec(t, db,
		`INSERT INTO feedback (category, went_well, could_improve, topics, created_at) VALUES (?, ?, ?, ?, ?)`,
		string(category), wentWell, couldImprove, string(topicsJSON), createdAt,
	)
}

func mustListFeedback(t *testing.T, db *sql.DB, category feedback.Category) []feedback.Entry {
	t.Helper()
	entries, err := feedback.List(db, category)
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	return entries
}
