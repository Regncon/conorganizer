package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/components/formsubmission"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

// approvalEventPlayersRoute wires the approval page's player/GM assignment
// endpoints. It lives apart from the rest of the admin router so the placement
// rules it enforces can be exercised without mounting the whole admin tree.
func approvalEventPlayersRoute(router chi.Router, db *sql.DB, liveManager *live.Manager, baseLogger *slog.Logger) {
	logger := baseLogger.With("component", "admin")

	router.Route("/event-players", func(eventPlayersRouter chi.Router) {
		eventPlayersRouter.Post("/post/add_first_choice", func(w http.ResponseWriter, r *http.Request) {
			type Store struct {
				BillettholderId int    `json:"assignmentBillettholderId"`
				EventId         string `json:"assignmentEventId"`
				PuljeId         string `json:"assignmentPuljeId"`
				// Set only by the age warning dialog, when the admin has seen the
				// warning and still wants the placement.
				AgeConfirmed bool `json:"assignmentAgeConfirmed"`
			}

			store := &Store{}

			if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
				http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
				return
			}
			if store.BillettholderId <= 0 {
				logger.Error(fmt.Errorf("invalid billettholder id for add first choice (event_id=%s, pulje_id=%s): invalid assignmentBillettholderId %d: must be greater than 0", store.EventId, store.PuljeId, store.BillettholderId).Error())
				http.Error(w, fmt.Errorf("invalid assignmentBillettholderId %d: must be greater than 0", store.BillettholderId).Error(), http.StatusNotFound)
				return
			}

			if approvalAgeGate(w, r, db, logger, ageGateRequest{
				billettholderID: store.BillettholderId,
				eventID:         store.EventId,
				puljeID:         store.PuljeId,
				isPlayer:        true,
				confirmed:       store.AgeConfirmed,
				retryMethod:     http.MethodPost,
				retryURL:        formsubmission.AddFirstChoiceURL,
			}) {
				return
			}

			var addFirstChoiceErr = formsubmission.AddPlayersFirstChoice(
				store.BillettholderId,
				store.EventId,
				store.PuljeId,
				db,
				baseLogger,
			)
			if addFirstChoiceErr != nil {
				logger.Error(fmt.Errorf("failed to add player as first choice: %w", addFirstChoiceErr).Error())
				http.Error(w, addFirstChoiceErr.Error(), http.StatusInternalServerError)
				return
			}
			logger.Info(
				"Successfully added player as first choice",
				"event_id", store.EventId,
				"pulje_id", store.PuljeId,
				"billettholder_id", store.BillettholderId,
			)
			if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
				logger.Error(fmt.Errorf("failed to broadcast add first choice update: %w", err).Error())
				http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
				return
			}

		})
		eventPlayersRouter.Post("/post/add_gm", func(w http.ResponseWriter, r *http.Request) {

			type Store struct {
				BillettholderId int    `json:"assignmentBillettholderId"`
				EventId         string `json:"assignmentEventId"`
				PuljeId         string `json:"assignmentPuljeId"`
				AgeConfirmed    bool   `json:"assignmentAgeConfirmed"`
			}
			store := &Store{}

			if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
				http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
				return
			}
			if store.BillettholderId <= 0 {
				logger.Error(fmt.Errorf("invalid billettholder id for add GM (event_id=%s, pulje_id=%s): invalid assignmentBillettholderId %d: must be greater than 0", store.EventId, store.PuljeId, store.BillettholderId).Error())
				http.Error(w, fmt.Errorf("invalid assignmentBillettholderId %d: must be greater than 0", store.BillettholderId).Error(), http.StatusNotFound)
				return
			}

			if approvalAgeGate(w, r, db, logger, ageGateRequest{
				billettholderID: store.BillettholderId,
				eventID:         store.EventId,
				puljeID:         store.PuljeId,
				isGM:            true,
				confirmed:       store.AgeConfirmed,
				retryMethod:     http.MethodPost,
				retryURL:        formsubmission.AddGMURL,
			}) {
				return
			}

			var updatePlayerStatusErr = formsubmission.UpdatePlayerStatus(
				store.EventId,
				store.PuljeId,
				store.BillettholderId,
				false,
				true,
				db,
				baseLogger,
			)
			if updatePlayerStatusErr != nil {
				logger.Error(fmt.Errorf("failed to add player as GM: %w", updatePlayerStatusErr).Error())
				http.Error(w, updatePlayerStatusErr.Error(), http.StatusInternalServerError)
				return
			}
			logger.Info(
				"Successfully Added player as GM",
				"event_id", store.EventId,
				"pulje_id", store.PuljeId,
				"billettholder_id", store.BillettholderId,
				"role", models.EventPlayerRoleGM,
			)
			if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
				logger.Error(fmt.Errorf("failed to broadcast add GM update: %w", err).Error())
				http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
				return
			}
		})
		eventPlayersRouter.Put("/update_status", func(w http.ResponseWriter, r *http.Request) {
			type Store struct {
				BillettholderId int    `json:"assignmentBillettholderId"`
				EventId         string `json:"assignmentEventId"`
				PuljeId         string `json:"assignmentPuljeId"`
				IsPlayer        bool   `json:"assignmentIsPlayer"`
				IsGm            bool   `json:"assignmentIsGm"`
				AgeConfirmed    bool   `json:"assignmentAgeConfirmed"`
			}
			store := &Store{}

			if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
				http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
				return
			}

			// Only a placement can break the age rule; taking someone back out
			// resolves it, so removals never have to be confirmed.
			if store.IsPlayer || store.IsGm {
				if approvalAgeGate(w, r, db, logger, ageGateRequest{
					billettholderID: store.BillettholderId,
					eventID:         store.EventId,
					puljeID:         store.PuljeId,
					isPlayer:        store.IsPlayer,
					isGM:            store.IsGm,
					confirmed:       store.AgeConfirmed,
					retryMethod:     http.MethodPut,
					retryURL:        formsubmission.UpdatePlayerStatusURL,
				}) {
					return
				}
			}

			var updatePlayerStatusErr = formsubmission.UpdatePlayerStatus(
				store.EventId,
				store.PuljeId,
				store.BillettholderId,
				store.IsPlayer,
				store.IsGm,
				db,
				baseLogger,
			)
			if updatePlayerStatusErr != nil {
				http.Error(w, updatePlayerStatusErr.Error(), http.StatusInternalServerError)
				return
			}
			logger.Info(
				"Successfully updated player status",
				"event_id", store.EventId,
				"pulje_id", store.PuljeId,
				"billettholder_id", store.BillettholderId,
				"assignment_is_player", store.IsPlayer,
				"assignment_is_gm", store.IsGm,
			)
			if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
				logger.Error(fmt.Errorf("failed to broadcast player status update: %w", err).Error())
				http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
				return
			}
		})
	})
}

// ageGateRequest is the placement the admin is about to make, as far as the age
// rule cares: who goes into which game, in which role, and where the
// confirmation should be sent if the admin insists.
type ageGateRequest struct {
	billettholderID int
	eventID         string
	puljeID         string
	isPlayer        bool
	isGM            bool
	confirmed       bool
	retryMethod     string
	retryURL        string
}

// approvalAgeGate answers with the confirmation dialog when the admin is about to
// place someone under 18 in an 18+ game, and reports true when the caller must
// stop. The billettholder picker only ever sends an id, so the age can only be
// checked here; the retry url and method travel with the warning because one
// dialog serves every placement button on the page.
func approvalAgeGate(w http.ResponseWriter, r *http.Request, db *sql.DB, logger *slog.Logger, req ageGateRequest) bool {
	if req.confirmed {
		return false
	}

	warning, err := adultsOnlyWarning(db, req.billettholderID, req.eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Billettholder or event not found", http.StatusNotFound)
			return true
		}
		logger.Error(err.Error(), "event_id", req.eventID, "billettholder_id", req.billettholderID)
		http.Error(w, "Failed to read age restriction", http.StatusInternalServerError)
		return true
	}
	if warning == "" {
		return false
	}

	sse := datastar.NewSSE(w, r)
	if err := sse.MarshalAndPatchSignals(map[string]any{
		"ageWarningText":            warning,
		"ageWarningBillettholderId": req.billettholderID,
		"ageWarningEventId":         req.eventID,
		"ageWarningPuljeId":         req.puljeID,
		"ageWarningIsPlayer":        req.isPlayer,
		"ageWarningIsGm":            req.isGM,
		"ageWarningMethod":          req.retryMethod,
		"ageWarningUrl":             req.retryURL,
	}); err != nil {
		logger.Error(fmt.Errorf("failed to patch age warning signals: %w", err).Error(), "event_id", req.eventID, "billettholder_id", req.billettholderID)
	}
	return true
}
