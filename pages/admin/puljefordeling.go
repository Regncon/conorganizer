package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

var errPuljeNotFound = errors.New("pulje not found")

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

func updatePuljeStatus(db *sql.DB, puljeID models.Pulje, status models.PuljeStatus) error {
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
