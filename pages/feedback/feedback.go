// Package feedback wires up the logged-in "Gi tilbakemelding" form at
// /tilbakemelding. The admin list at /admin/tilbakemeldinger/ lives in
// pages/admin.
package feedback

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

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
	feedbackTopicsSignal   = "feedbackTopics"
	feedbackFormElementID  = "feedback-form"
)

// feedbackSubmission is the Datastar signal payload the feedback form posts.
type feedbackSubmission struct {
	Category string   `json:"feedbackCategory"`
	Message  string   `json:"feedbackMessage"`
	Topics   []string `json:"feedbackTopics"`
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
			if err := sse.PatchElementTempl(feedbackFormWrapper(preselectedCategory(r))); err != nil {
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

// preselectedCategory maps ?om=<slug> to a category. Some category is always
// selected: an unrecognized or missing slug falls back to the website
// category, since the form always needs one preselected.
func preselectedCategory(r *http.Request) feedback.Category {
	if category, ok := feedback.CategoryFromSlug(r.URL.Query().Get("om")); ok {
		return category
	}
	return feedback.CategoryWebsite
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
	payload := &feedbackSubmission{}
	if err := datastar.ReadSignals(r, payload); err != nil {
		http.Error(w, "Klarte ikke å lese skjemadata", http.StatusBadRequest)
		return
	}

	submission := feedback.Submission{
		Category: feedback.Category(payload.Category),
		Message:  payload.Message,
		Topics:   toTopics(payload.Topics),
	}

	sse := datastar.NewSSE(w, r)
	feedbackErrors := newFeedbackErrors()

	if fieldErrors := submission.Validate(); len(fieldErrors) > 0 {
		for _, fieldError := range fieldErrors {
			feedbackErrors.Set(feedbackSignalForField(fieldError.Field), feedbackFieldErrorMessage(fieldError))
		}
		if err := feedbackErrors.Patch(sse); err != nil {
			logger.Error(fmt.Errorf("failed to patch feedback validation errors: %w", err).Error())
		}
		return
	}

	if err := feedback.Submit(db, submission); err != nil {
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
	if err := sse.PatchElementTempl(feedbackThankYou(submission.Category)); err != nil {
		logger.Error(fmt.Errorf("failed to patch feedback thank-you: %w", err).Error())
	}
}

func toTopics(values []string) []feedback.Topic {
	topics := make([]feedback.Topic, len(values))
	for i, value := range values {
		topics[i] = feedback.Topic(value)
	}
	return topics
}

// feedbackSignalForField maps a service field name to the frontend signal
// name the matching error message is shown under.
func feedbackSignalForField(field string) string {
	switch field {
	case feedback.FieldCategory:
		return feedbackCategorySignal
	case feedback.FieldMessage:
		return feedbackMessageSignal
	case feedback.FieldTopics:
		return feedbackTopicsSignal
	}
	return feedbackMessageSignal
}

// feedbackFieldErrorMessage is the Norwegian message shown for one
// FieldError, matching the same rule the frontend already checked live.
func feedbackFieldErrorMessage(fieldError feedback.FieldError) string {
	switch fieldError.Code {
	case feedback.CodeRequired:
		return "Skriv en tilbakemelding."
	case feedback.CodeTooLong:
		return fmt.Sprintf("Teksten er for lang (maks %d tegn).", feedback.MaxTextLength)
	case feedback.CodePersonalInfo:
		return "Det ser ut som du har skrevet en e-postadresse eller et telefonnummer. Fjern det før du sender."
	case feedback.CodeInvalid:
		if fieldError.Field == feedback.FieldTopics {
			return "Velg emner som passer til valgt kategori."
		}
		return "Velg hva tilbakemeldingen gjelder."
	}
	return "Klarte ikke å validere skjemaet. Prøv igjen."
}

func newFeedbackErrors() *errorfeedback.FeedbackErrors {
	return errorfeedback.New(feedbackCategorySignal, feedbackMessageSignal, feedbackTopicsSignal)
}
