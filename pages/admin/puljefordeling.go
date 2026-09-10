package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

var errPuljeNotFound = errors.New("pulje not found")

func getPuljer(db *sql.DB) ([]models.PuljeRow, error) {
	const query = `
		SELECT id, name, status, start_at, end_at
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
}

// getBillettholderAge returns the participant's display name and whether they
// are over 18. The assign endpoint needs both: the age decides whether a pin
// into an 18+ game must be confirmed, the name goes into that confirmation.
func getBillettholderAge(db *sql.DB, billettholderID int) (string, bool, error) {
	var firstName, lastName string
	var isOver18 bool
	if err := db.QueryRow(
		`SELECT first_name, last_name, is_over_18 FROM billettholdere WHERE id = ?`,
		billettholderID,
	).Scan(&firstName, &lastName, &isOver18); err != nil {
		return "", false, fmt.Errorf("get billettholder %d: %w", billettholderID, err)
	}
	return strings.TrimSpace(firstName + " " + lastName), isOver18, nil
}

// getEventAgeGroup returns the event's title and age group, so the assign
// endpoint can name the game in the age warning it sends back.
func getEventAgeGroup(db *sql.DB, eventID string) (string, models.AgeGroup, error) {
	var title string
	var ageGroup models.AgeGroup
	if err := db.QueryRow(`SELECT title, age_group FROM events WHERE id = ?`, eventID).Scan(&title, &ageGroup); err != nil {
		return "", "", fmt.Errorf("get event %s: %w", eventID, err)
	}
	return title, ageGroup, nil
}

// adultsOnlyWarning returns the confirmation an admin must accept before placing
// a participant under 18 in an 18+ game, or "" when the placement raises no age
// conflict. The solver never makes such a placement; an admin may, but only
// deliberately, so every endpoint that seats someone asks this first. Callers map
// a wrapped sql.ErrNoRows to 404.
func adultsOnlyWarning(db *sql.DB, billettholderID int, eventID string) (string, error) {
	playerName, isOver18, err := getBillettholderAge(db, billettholderID)
	if err != nil {
		return "", err
	}
	eventTitle, ageGroup, err := getEventAgeGroup(db, eventID)
	if err != nil {
		return "", err
	}
	if isOver18 || ageGroup != models.AgeGroupAdultsOnly {
		return "", nil
	}
	return fmt.Sprintf("%s er under 18 år, og «%s» er 18+. Vil du plassere likevel?", playerName, eventTitle), nil
}
