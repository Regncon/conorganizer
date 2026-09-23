package admin

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	datastar "github.com/starfederation/datastar-go/datastar"
)

type puljeInteresseSignaler struct {
	BillettholderID int                  `json:"assignmentBillettholderId"`
	EventID         string               `json:"assignmentEventId"`
	PuljeID         string               `json:"assignmentPuljeId"`
	Level           models.InterestLevel `json:"assignmentInterestLevel"`
	LukkDialog      bool                 `json:"assignmentCloseDialog"`
}

func puljeInteresseHandler(db *sql.DB, liveManager *live.Manager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var signaler puljeInteresseSignaler
		if err := datastar.ReadSignals(r, &signaler); err != nil {
			http.Error(w, "Ugyldige interessedata", http.StatusBadRequest)
			return
		}
		puljeID, ok := models.ParsePulje(signaler.PuljeID)
		if !ok {
			http.Error(w, "Ugyldig pulje", http.StatusBadRequest)
			return
		}
		if err := puljefordeling.SetInteresse(db, puljeID, signaler.EventID, signaler.BillettholderID, signaler.Level); err != nil {
			switch {
			case errors.Is(err, puljefordeling.ErrPuljeCompleted):
				http.Error(w, "Puljen er publisert og kan ikke endres", http.StatusConflict)
			case errors.Is(err, sql.ErrNoRows):
				http.Error(w, "Fant ikke billettholderen", http.StatusNotFound)
			case errors.Is(err, puljefordeling.ErrUgyldigTildeling):
				http.Error(w, "Ugyldig interesse", http.StatusBadRequest)
			default:
				logger.Error(err.Error(), "pulje_id", puljeID, "event_id", signaler.EventID, "billettholder_id", signaler.BillettholderID)
				http.Error(w, "Kunne ikke oppdatere interesse", http.StatusInternalServerError)
			}
			return
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
			logger.Error(err.Error(), "pulje_id", puljeID, "event_id", signaler.EventID, "billettholder_id", signaler.BillettholderID)
		}
		if signaler.LukkDialog {
			sse := datastar.NewSSE(w, r)
			if err := sse.MarshalAndPatchSignals(map[string]bool{"assignmentActionCompleted": true}); err != nil {
				logger.Error(err.Error(), "pulje_id", puljeID, "event_id", signaler.EventID, "billettholder_id", signaler.BillettholderID)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
