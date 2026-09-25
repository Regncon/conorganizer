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

// puljeLagringsgrense is how many changes the save dialog lists one by one;
// above it the dialog only gives the count.
const puljeLagringsgrense = 20

var errPuljeUlagret = errors.New("pulje has unsaved distribution changes")

// puljeLagringsforhandsvisningHandler opens the save dialog, listing what saving
// the distribution would change.
func puljeLagringsforhandsvisningHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pulje, ok := models.ParsePulje(chi.URLParam(r, "pulje"))
		if !ok {
			http.Error(w, "Invalid pulje ID", http.StatusBadRequest)
			return
		}
		status, err := puljefordeling.Lagringsstatus(db, pulje)
		if err != nil {
			logger.Error(err.Error(), "pulje_id", pulje)
			http.Error(w, "Failed to compare distribution", http.StatusInternalServerError)
			return
		}
		sendTildelingsdialog(w, r, logger, puljeLagringsDialogInnhold(pulje, status))
	}
}

// puljeLagringsHandler saves the distribution once the admin has confirmed the
// changes they were shown. If the changes differ by now, the dialog is shown
// again with the current changes instead.
func puljeLagringsHandler(db *sql.DB, liveManager *live.Manager, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pulje, ok := models.ParsePulje(chi.URLParam(r, "pulje"))
		if !ok {
			http.Error(w, "Invalid pulje ID", http.StatusBadRequest)
			return
		}
		var signaler struct {
			Bekreftelse string `json:"lagreBekreftelse"`
		}
		if err := datastar.ReadSignals(r, &signaler); err != nil {
			http.Error(w, "Invalid signals", http.StatusBadRequest)
			return
		}
		status, err := puljefordeling.LagreBekreftetFordeling(db, pulje, signaler.Bekreftelse)
		if err != nil {
			if errors.Is(err, puljefordeling.ErrPuljeCompleted) {
				http.Error(w, "Pulje is published; changes are not allowed", http.StatusConflict)
				return
			}
			logger.Error(err.Error(), "pulje_id", pulje)
			http.Error(w, "Failed to commit distribution", http.StatusInternalServerError)
			return
		}
		if status != nil {
			sendTildelingsdialog(w, r, logger, puljeLagringsDialogInnhold(pulje, *status))
			return
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketEvents, live.BucketRooms); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast commit: %w", err).Error(), "pulje_id", pulje)
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// puljeLagringsURL is the endpoint that saves the confirmed distribution.
func puljeLagringsURL(pulje models.Pulje) string {
	return fmt.Sprintf("/admin/api/puljefordeling/%s/commit", pulje)
}
