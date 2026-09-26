// Package feedback wires up the logged-in "Gi tilbakemelding" form at
// /tilbakemelding. The admin list at /admin/tilbakemeldinger/ lives in
// pages/admin.
package feedback

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/Regncon/conorganizer/components/errorfeedback"
	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/feedback"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

const (
	feedbackCategorySignal = "feedbackCategory"
	feedbackMessageSignal  = "feedbackMessage"
	feedbackFormElementID  = "feedback-form"
)

// feedbackSubmission is the Datastar signal payload the feedback form posts.
type feedbackSubmission struct {
	Category string `json:"feedbackCategory"`
	Message  string `json:"feedbackMessage"`
}

// SetupFeedbackRoute registers the logged-in feedback form. router is the
// logged-in router (the same tier as /profile).
func SetupFeedbackRoute(router chi.Router, db *sql.DB, logger *slog.Logger) error {
	logger = logger.With("component", "feedback")

	router.Route("/tilbakemelding", func(feedbackRouter chi.Router) {
		feedbackRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			renderFeedbackPage(w, r, db, logger, preselectedCategory(r))
		})

		feedbackRouter.Get("/reset", func(w http.ResponseWriter, r *http.Request) {
			sse := datastar.NewSSE(w, r)
			if err := sse.PatchElementTempl(feedbackFormWrapper("")); err != nil {
				logger.Error(fmt.Errorf("failed to patch reset feedback form: %w", err).Error())
			}
			// The error signals use __ifmissing, so the new form would show
			// stale errors unless they are cleared explicitly.
			if err := newFeedbackErrors().Patch(sse); err != nil {
				logger.Error(fmt.Errorf("failed to clear feedback errors on reset: %w", err).Error())
			}
		})

		feedbackRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
			submitFeedback(w, r, db, logger)
		})
	})

	return nil
}

func preselectedCategory(r *http.Request) feedback.Category {
	category, _ := feedback.CategoryFromSlug(r.URL.Query().Get("om"))
	return category
}

func renderFeedbackPage(w http.ResponseWriter, r *http.Request, db *sql.DB, logger *slog.Logger, preselected feedback.Category) {
	ctx := r.Context()
	user := userctx.GetUserRequestInfo(ctx)
	if err := layouts.Base(
		"Gi tilbakemelding",
		user,
		db,
		logger,
		feedbackFormPage(preselected),
	).Render(ctx, w); err != nil {
		logger.Error(fmt.Errorf("failed to render feedback page: %w", err).Error(), "user_id", user.Id)
	}
}

func submitFeedback(w http.ResponseWriter, r *http.Request, db *sql.DB, logger *slog.Logger) {
	submission := &feedbackSubmission{}
	if err := datastar.ReadSignals(r, submission); err != nil {
		http.Error(w, "Klarte ikke å lese skjemadata", http.StatusBadRequest)
		return
	}

	category := feedback.Category(submission.Category)
	message := strings.TrimSpace(submission.Message)

	feedbackErrors := newFeedbackErrors()
	if !category.Valid() {
		feedbackErrors.Set(feedbackCategorySignal, "Velg hva tilbakemeldingen gjelder.")
	}
	if message == "" {
		feedbackErrors.Set(feedbackMessageSignal, "Skriv en tilbakemelding før du sender inn.")
	} else if length := utf8.RuneCountInString(message); length > feedback.MaxMessageLength {
		feedbackErrors.Set(feedbackMessageSignal, fmt.Sprintf("Tilbakemeldingen er for lang (%d av maks %d tegn).", length, feedback.MaxMessageLength))
	}

	sse := datastar.NewSSE(w, r)
	if feedbackErrors.HasErrors() {
		if err := feedbackErrors.Patch(sse); err != nil {
			logger.Error(fmt.Errorf("failed to patch feedback validation errors: %w", err).Error())
		}
		return
	}

	if err := feedback.Submit(db, category, submission.Message); err != nil {
		logger.Error(fmt.Errorf("failed to store feedback: %w", err).Error())
		feedbackErrors.Set(feedbackMessageSignal, "Klarte ikke å lagre tilbakemeldingen. Prøv igjen.")
		if patchErr := feedbackErrors.Patch(sse); patchErr != nil {
			logger.Error(fmt.Errorf("failed to patch feedback storage error: %w", patchErr).Error())
		}
		return
	}

	// Clear any errors from earlier attempts so they do not reappear when
	// the user resets the form to send another one.
	if err := feedbackErrors.Patch(sse); err != nil {
		logger.Error(fmt.Errorf("failed to clear feedback validation errors: %w", err).Error())
	}
	if err := sse.PatchElementTempl(feedbackThankYou()); err != nil {
		logger.Error(fmt.Errorf("failed to patch feedback thank-you: %w", err).Error())
	}
}

func newFeedbackErrors() *errorfeedback.FeedbackErrors {
	return errorfeedback.New(feedbackCategorySignal, feedbackMessageSignal)
}
