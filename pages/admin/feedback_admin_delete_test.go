package admin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/feedback"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestFeedbackAdminPage_DeleteButtonOpensAreYouSureDialog(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en lagret tilbakemelding i admin-listen, filtrert på kategori.",
		When:  "Når admin-listen rendres.",
		Then:  "Så har tilbakemeldingen en slett-knapp som åpner en «Er du sikker?»-dialog med avbryt og bekreft.",
	})

	// Given
	db, _ := testutil.CreateTestDBAndLogger(t, "feedback_admin_delete_button")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "rydd meg bort", "2026-09-20T12:00:00.000Z")
	entry := mustListFeedback(t, db, "")[0]
	expectedDialogID := fmt.Sprintf("feedback-delete-dialog-%d", entry.ID)
	expectedDeleteAction := fmt.Sprintf("@delete('/admin/tilbakemeldinger/%d?om=%s')", entry.ID, feedback.CategoryOther.Slug())

	// When
	doc := templtest.Render(t, feedbackAdminPage([]feedback.Entry{entry}, feedback.CategoryOther))

	// Then
	deleteButton := doc.Find(".feedback-entry .feedback-delete-button")
	if deleteButton.Length() != 1 {
		t.Fatalf("expected one delete button, got %d", deleteButton.Length())
	}
	if deleteButton.AttrOr("commandfor", "") != expectedDialogID || deleteButton.AttrOr("command", "") != "show-modal" {
		t.Fatalf("expected the delete button to open %q, got commandfor=%q command=%q", expectedDialogID, deleteButton.AttrOr("commandfor", ""), deleteButton.AttrOr("command", ""))
	}
	if deleteButton.Find("svg").Length() != 1 {
		t.Fatal("expected the delete button to show the trash icon")
	}
	dialog := doc.Find("dialog#" + expectedDialogID)
	if dialog.Length() != 1 {
		t.Fatalf("expected a confirmation dialog %q", expectedDialogID)
	}
	if !strings.Contains(dialog.Text(), "Er du sikker") {
		t.Fatalf("expected the dialog to ask if the admin is sure, got %q", dialog.Text())
	}
	cancel := dialog.Find(`button[command="close"]`)
	if cancel.AttrOr("commandfor", "") != expectedDialogID || cancel.Find("svg").Length() != 1 {
		t.Fatal("expected a cancel button with an icon that closes the dialog")
	}
	confirm := dialog.Find(".feedback-delete-confirm")
	if confirm.AttrOr("data-on:click", "") != expectedDeleteAction || confirm.Find("svg").Length() != 1 {
		t.Fatalf("expected a confirm button with an icon running %q, got %q", expectedDeleteAction, confirm.AttrOr("data-on:click", ""))
	}
}

func TestFeedbackAdminRoute_AdminDeletesFeedbackAndGetsUpdatedList(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt to lagrede tilbakemeldinger.",
		When:  "Når en administrator bekrefter sletting av den ene.",
		Then:  "Så slettes den, og svaret er den oppdaterte listen uten den slettede tilbakemeldingen.",
	})

	// Given
	expectedRemaining := "beholdes"
	deletedMessage := "utdatert klage om kaffen"
	db, logger := testutil.CreateTestDBAndLogger(t, "feedback_admin_delete_route")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, deletedMessage, "2026-09-20T12:00:00.000Z")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, expectedRemaining, "2026-09-21T12:00:00.000Z")
	target := feedbackEntryWithMessage(t, mustListFeedback(t, db, ""), deletedMessage)
	router := chi.NewRouter()
	router.Route("/admin", func(r chi.Router) {
		feedbackAdminRoute(r, db, logger)
	})
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/tilbakemeldinger/%d", target.ID), nil)
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", recorder.Code, recorder.Body.String())
	}
	remaining := mustListFeedback(t, db, "")
	if len(remaining) != 1 || remaining[0].Message != expectedRemaining {
		t.Fatalf("expected only %q to remain, got %+v", expectedRemaining, remaining)
	}
	body := recorder.Body.String()
	if !strings.Contains(body, "datastar-patch-elements") || !strings.Contains(body, expectedRemaining) || !strings.Contains(body, "1 tilbakemelding") {
		t.Fatalf("expected the response to patch in the updated list, got: %s", body)
	}
	if strings.Contains(body, deletedMessage) {
		t.Fatalf("expected the deleted feedback to be gone from the updated list, got: %s", body)
	}
}

func TestFeedbackAdminRoute_DeletingUnknownFeedbackIsNotFound(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en tilbakemelding som ikke finnes.",
		When:  "Når en administrator prøver å slette den.",
		Then:  "Så svarer serveren 404.",
	})

	// Given
	expectedStatus := http.StatusNotFound
	db, logger := testutil.CreateTestDBAndLogger(t, "feedback_admin_delete_unknown")
	router := chi.NewRouter()
	router.Route("/admin", func(r chi.Router) {
		feedbackAdminRoute(r, db, logger)
	})
	request := httptest.NewRequest(http.MethodDelete, "/admin/tilbakemeldinger/999", nil)
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
}

func TestFeedbackAdminRoute_NonAdminCannotDeleteFeedback(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en innlogget bruker som ikke er admin, og en lagret tilbakemelding.",
		When:  "Når brukeren sender en slett-forespørsel gjennom admin-kravet.",
		Then:  "Så nektes tilgang, og tilbakemeldingen blir liggende.",
	})

	// Given
	expectedStatus := http.StatusForbidden
	db, logger := testutil.CreateTestDBAndLogger(t, "feedback_admin_delete_non_admin")
	insertFeedbackAdminEntry(t, db, feedback.CategoryOther, "blir liggende", "2026-09-20T12:00:00.000Z")
	target := mustListFeedback(t, db, "")[0]
	router := chi.NewRouter()
	adminRouter := router.With(
		userctx.UserMiddleware(logger, db),
		authctx.RequireAdmin(logger, authctx.WithForbiddenHandler(userctx.AdminForbiddenHandler(db, logger))),
	)
	adminRouter.Route("/admin", func(r chi.Router) {
		feedbackAdminRoute(r, db, logger)
	})
	request := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/admin/tilbakemeldinger/%d", target.ID), nil)
	request = request.WithContext(authctx.WithUserToken(request.Context(), "feedback-non-admin", "non-admin@example.com"))
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, recorder.Code)
	}
	if got := len(mustListFeedback(t, db, "")); got != 1 {
		t.Fatalf("expected the feedback to remain, found %d entries", got)
	}
}

func feedbackEntryWithMessage(t *testing.T, entries []feedback.Entry, message string) feedback.Entry {
	t.Helper()
	for _, entry := range entries {
		if entry.Message == message {
			return entry
		}
	}
	t.Fatalf("no feedback with message %q", message)
	return feedback.Entry{}
}
