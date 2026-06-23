package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	pf "github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

var errPuljeNotFound = errors.New("pulje not found")

func isValidPuljeStatus(status models.PuljeStatus) bool {
	switch status {
	case models.PuljeStatusOpen, models.PuljeStatusLocked, models.PuljeStatusCompleted:
		return true
	default:
		return false
	}
}

func puljeIsLocked(status models.PuljeStatus) bool {
	return status == models.PuljeStatusLocked || status == models.PuljeStatusCompleted
}

func puljeIsCompleted(status models.PuljeStatus) bool {
	return status == models.PuljeStatusCompleted
}

func puljeStatusUpdateAction(
	pulje models.PuljeRow,
	message string,
	checkedStatus models.PuljeStatus,
	uncheckedStatus models.PuljeStatus,
) string {
	return fmt.Sprintf(
		"if (!confirm(%q)) { evt.preventDefault(); } else { $puljeStatus = evt.currentTarget.checked ? %q : %q; @put('/admin/api/puljer/%s/status') }",
		message,
		string(checkedStatus),
		string(uncheckedStatus),
		pulje.ID,
	)
}

func updatePuljeStatus(db *sql.DB, puljeID models.Pulje, status models.PuljeStatus) error {
	const query = `UPDATE puljer SET status = ? WHERE id = ?`

	result, err := db.Exec(query, status, puljeID)
	if err != nil {
		return fmt.Errorf("update pulje %s status to %s: %w", puljeID, status, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected for pulje %s status update: %w", puljeID, err)
	}

	if rowsAffected == 0 {
		return errPuljeNotFound
	}

	return nil
}

func puljefordelingStatusRoute(router chi.Router, db *sql.DB, liveManager *live.Manager, logger *slog.Logger) {
	logger = logger.With("component", "admin_puljefordeling")

	router.Put("/api/puljer/{puljeId}/status", func(w http.ResponseWriter, r *http.Request) {
		puljeID, ok := models.ParsePulje(chi.URLParam(r, "puljeId"))
		if !ok {
			http.Error(w, "Invalid pulje ID", http.StatusBadRequest)
			return
		}

		type Store struct {
			PuljeStatus models.PuljeStatus `json:"puljeStatus"`
		}

		store := &Store{}
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !isValidPuljeStatus(store.PuljeStatus) {
			http.Error(w, "Invalid pulje status", http.StatusBadRequest)
			return
		}

		// Read the current status to decide the transition side-effects.
		var currentRaw string
		if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, string(puljeID)).Scan(&currentRaw); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "Pulje not found", http.StatusNotFound)
				return
			}
			logger.Error(fmt.Errorf("read current pulje status: %w", err).Error(), "pulje_id", puljeID)
			http.Error(w, "Failed to read pulje status", http.StatusInternalServerError)
			return
		}
		current := models.PuljeStatus(currentRaw)
		next := store.PuljeStatus

		var transitionErr error
		switch {
		case current == models.PuljeStatusOpen && next == models.PuljeStatusLocked:
			// Lock: commit the distribution and set status (atomic, inside the service).
			transitionErr = pf.CommitPuljeAssignments(db, puljeID, logger)
		case current == models.PuljeStatusLocked && next == models.PuljeStatusOpen:
			// Unlock: drop solver seats and reopen (atomic, inside the service).
			transitionErr = pf.RevertPuljeAssignments(db, puljeID)
		default:
			// Locked↔Completed and any no-op: status only.
			transitionErr = updatePuljeStatus(db, puljeID, next)
		}
		if transitionErr != nil {
			if errors.Is(transitionErr, errPuljeNotFound) {
				http.Error(w, "Pulje not found", http.StatusNotFound)
				return
			}
			logger.Error(transitionErr.Error(), "pulje_id", puljeID, "pulje_status", next)
			http.Error(w, "Failed to update pulje status", http.StatusInternalServerError)
			return
		}

		// Live broadcasts are best-effort: the status transition is already
		// committed, so a failed broadcast is logged but does not fail the request.
		if err := liveManager.Broadcast(r.Context(), live.BucketEvents); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast pulje status update: %w", err).Error(), "pulje_id", puljeID, "pulje_status", next)
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast interests update: %w", err).Error(), "pulje_id", puljeID, "pulje_status", next)
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
