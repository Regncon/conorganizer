package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

var errPuljeNotFound = errors.New("pulje not found")
var errPuljeStepOrder = errors.New("pulje status steps must be done in order")

func getPuljer(db *sql.DB) ([]models.PuljeRow, error) {
	const query = `
		SELECT id, name, status, closing_warning_active, start_at, end_at
		FROM puljer
		ORDER BY start_at ASC
	`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query puljer: %w", err)
	}
	defer rows.Close()

	puljer := make([]models.PuljeRow, 0)
	for rows.Next() {
		var pulje models.PuljeRow
		if err := rows.Scan(
			&pulje.ID,
			&pulje.Name,
			&pulje.Status,
			&pulje.ClosingWarningActive,
			&pulje.StartAt,
			&pulje.EndAt,
		); err != nil {
			return nil, fmt.Errorf("scan pulje row: %w", err)
		}
		puljer = append(puljer, pulje)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate pulje rows: %w", err)
	}

	return puljer, nil
}

func getPuljeStatus(db *sql.DB, pulje models.Pulje) (models.PuljeStatus, error) {
	var status models.PuljeStatus
	if err := db.QueryRow(`SELECT status FROM puljer WHERE id = ?`, pulje).Scan(&status); err != nil {
		return "", fmt.Errorf("get pulje %s status: %w", pulje, err)
	}
	return status, nil
}

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

// puljeStatusStepAllowed enforces the three puljefordeling steps in order:
// 1. closing warning, 2. locked, 3. published. A step can only be ticked when
// the previous one is done, and only unticked while the next one is not.
// Locking clears the closing warning, so a non-open pulje counts step 1 as done.
func puljeStatusStepAllowed(current models.PuljeRow, next models.PuljeStatus) bool {
	switch next {
	case models.PuljeStatusOpen:
		return current.Status != models.PuljeStatusCompleted
	case models.PuljeStatusLocked:
		return current.Status != models.PuljeStatusOpen || current.ClosingWarningActive
	case models.PuljeStatusCompleted:
		return current.Status != models.PuljeStatusOpen
	default:
		return false
	}
}

func puljeClosingWarningStepDone(pulje models.PuljeRow) bool {
	return pulje.ClosingWarningActive || pulje.Status != models.PuljeStatusOpen
}

func updatePuljeStatus(db *sql.DB, puljeID models.Pulje, status models.PuljeStatus) error {
	current := models.PuljeRow{ID: puljeID}
	err := db.QueryRow(
		`SELECT status, closing_warning_active FROM puljer WHERE id = ?`,
		puljeID,
	).Scan(&current.Status, &current.ClosingWarningActive)
	if errors.Is(err, sql.ErrNoRows) {
		return errPuljeNotFound
	}
	if err != nil {
		return fmt.Errorf("load pulje %s status: %w", puljeID, err)
	}
	if !puljeStatusStepAllowed(current, status) {
		return errPuljeStepOrder
	}
	if status == models.PuljeStatusCompleted && current.Status != models.PuljeStatusCompleted {
		saveStatus, err := puljefordeling.LoadSaveStatus(db, puljeID)
		if err != nil {
			return fmt.Errorf("check unsaved distribution for %s: %w", puljeID, err)
		}
		if saveStatus.HasChanges() {
			return errPuljeUnsaved
		}
	}

	const query = `
		UPDATE puljer
		SET status = ?, closing_warning_active = CASE WHEN ? = 'Open' THEN closing_warning_active ELSE FALSE END
		WHERE id = ?
	`

	result, err := db.Exec(query, status, status, puljeID)
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

func updatePuljeClosingWarning(db *sql.DB, puljeID models.Pulje, active bool) error {
	result, err := db.Exec(
		`UPDATE puljer
		 SET closing_warning_active = CASE WHEN ? AND status = 'Open' THEN TRUE ELSE FALSE END
		 WHERE id = ?`,
		active,
		puljeID,
	)
	if err != nil {
		return fmt.Errorf("update pulje %s closing warning: %w", puljeID, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected for pulje %s closing warning update: %w", puljeID, err)
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

		if err := updatePuljeStatus(db, puljeID, store.PuljeStatus); err != nil {
			if errors.Is(err, errPuljeNotFound) {
				http.Error(w, "Pulje not found", http.StatusNotFound)
				return
			}
			if errors.Is(err, errPuljeStepOrder) {
				http.Error(w, "Forrige steg må være fullført først", http.StatusConflict)
				return
			}
			if errors.Is(err, errPuljeUnsaved) {
				http.Error(w, "Lagre fordelingen før den publiseres", http.StatusConflict)
				return
			}
			logger.Error(err.Error(), "pulje_id", puljeID, "pulje_status", store.PuljeStatus)
			http.Error(w, "Failed to update pulje status", http.StatusInternalServerError)
			return
		}

		if err := liveManager.Broadcast(r.Context(), live.BucketEvents); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast pulje status update: %w", err).Error(), "pulje_id", puljeID, "pulje_status", store.PuljeStatus)
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})

	router.Put("/api/puljer/{puljeId}/closing-warning", func(w http.ResponseWriter, r *http.Request) {
		puljeID, ok := models.ParsePulje(chi.URLParam(r, "puljeId"))
		if !ok {
			http.Error(w, "Invalid pulje ID", http.StatusBadRequest)
			return
		}

		type Store struct {
			ClosingWarningActive bool `json:"closingWarningActive"`
		}
		store := &Store{}
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := updatePuljeClosingWarning(db, puljeID, store.ClosingWarningActive); err != nil {
			if errors.Is(err, errPuljeNotFound) {
				http.Error(w, "Pulje not found", http.StatusNotFound)
				return
			}
			logger.Error(err.Error(), "pulje_id", puljeID)
			http.Error(w, "Failed to update pulje closing warning", http.StatusInternalServerError)
			return
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketEvents); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast pulje closing warning update: %w", err).Error(), "pulje_id", puljeID)
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
