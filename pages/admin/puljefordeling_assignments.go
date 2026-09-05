package admin

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Regncon/conorganizer/models"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

type puljeAssignmentTarget struct {
	BillettholderID int
	EventID         string
	PuljeID         models.Pulje
}

func readPuljeAssignmentTarget(r *http.Request) (puljeAssignmentTarget, error) {
	var signals struct {
		BillettholderID int    `json:"assignmentBillettholderId"`
		EventID         string `json:"assignmentEventId"`
		PuljeID         string `json:"assignmentPuljeId"`
	}
	if err := datastar.ReadSignals(r, &signals); err != nil {
		return puljeAssignmentTarget{}, fmt.Errorf("read assignment signals: %w", err)
	}
	puljeID, ok := models.ParsePulje(signals.PuljeID)
	if !ok {
		return puljeAssignmentTarget{}, fmt.Errorf("invalid pulje ID %q", signals.PuljeID)
	}
	if signals.EventID == "" {
		return puljeAssignmentTarget{}, fmt.Errorf("event ID is required")
	}
	if signals.BillettholderID <= 0 {
		return puljeAssignmentTarget{}, fmt.Errorf("billettholder ID must be greater than zero")
	}
	return puljeAssignmentTarget{
		BillettholderID: signals.BillettholderID,
		EventID:         signals.EventID,
		PuljeID:         puljeID,
	}, nil
}

func readPuljeRemovalTarget(r *http.Request) (puljeAssignmentTarget, error) {
	puljeValue := chi.URLParam(r, "pulje")
	puljeID, ok := models.ParsePulje(puljeValue)
	if !ok {
		return puljeAssignmentTarget{}, fmt.Errorf("invalid pulje ID %q", puljeValue)
	}
	eventID := chi.URLParam(r, "event")
	if eventID == "" {
		return puljeAssignmentTarget{}, fmt.Errorf("event ID is required")
	}
	billettholderID, err := strconv.Atoi(chi.URLParam(r, "billettholderId"))
	if err != nil || billettholderID <= 0 {
		return puljeAssignmentTarget{}, fmt.Errorf("invalid billettholder ID %q", chi.URLParam(r, "billettholderId"))
	}
	return puljeAssignmentTarget{
		BillettholderID: billettholderID,
		EventID:         eventID,
		PuljeID:         puljeID,
	}, nil
}

func puljeAllowsAdminAssignmentChanges(w http.ResponseWriter, db *sql.DB, logger *slog.Logger, pulje models.Pulje) bool {
	status, err := getPuljeStatus(db, pulje)
	if err != nil {
		logger.Error(err.Error(), "pulje_id", pulje)
		http.Error(w, "Failed to read pulje status", http.StatusInternalServerError)
		return false
	}
	if puljeIsCompleted(status) {
		http.Error(w, "Pulje is published; changes are not allowed", http.StatusConflict)
		return false
	}
	return true
}
