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

const puljeTildelingsURL = "/admin/api/puljefordeling/assign"

type puljeTildelingssignaler struct {
	BillettholderID int                    `json:"assignmentBillettholderId"`
	EventID         string                 `json:"assignmentEventId"`
	PuljeID         string                 `json:"assignmentPuljeId"`
	Role            models.EventPlayerRole `json:"assignmentRole"`
	FraLeggTil      bool                   `json:"assignmentFromAddMenu"`
	FraEventID      string                 `json:"assignmentFromEventId"`
	FraManuellPlass bool                   `json:"assignmentFromManualSeat"`
	Bekreftelse     string                 `json:"assignmentConfirmation"`
	AlderBekreftet  bool                   `json:"assignmentAgeConfirmed"`
	LukkDialog      bool                   `json:"assignmentCloseDialog"`
}

func puljeTildelingsHandler(db *sql.DB, liveManager *live.Manager, logger *slog.Logger) http.HandlerFunc {
	logger = logger.With("component", "admin_tildelinger")
	return func(w http.ResponseWriter, r *http.Request) {
		var signaler puljeTildelingssignaler
		if err := datastar.ReadSignals(r, &signaler); err != nil {
			http.Error(w, "Ugyldige tildelingsdata", http.StatusBadRequest)
			return
		}
		pulje, ok := models.ParsePulje(signaler.PuljeID)
		if !ok {
			http.Error(w, "Ugyldig pulje", http.StatusBadRequest)
			return
		}
		role := signaler.Role
		if role == "" {
			role = models.EventPlayerRolePlayer
		}
		valg := puljefordeling.Tildelingsvalg{
			PuljeID: pulje, EventID: signaler.EventID, BillettholderID: signaler.BillettholderID,
			Role: role, FraLeggTil: signaler.FraLeggTil,
			Bekreftelse: signaler.Bekreftelse, AlderBekreftet: signaler.AlderBekreftet,
			FraEventID: signaler.FraEventID, FraManuellPlass: signaler.FraManuellPlass,
		}
		varsel, err := puljefordeling.TildelBillettholder(db, valg)
		if err != nil {
			tildelingsfeil(w, logger, err)
			return
		}
		if varsel != nil {
			sendTildelingsvarsel(w, r, logger, *varsel, valg)
			return
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketEvents, live.BucketInterests, live.BucketRooms); err != nil {
			logger.Error(err.Error(), "pulje_id", pulje, "event_id", signaler.EventID, "billettholder_id", signaler.BillettholderID)
		}
		if signaler.LukkDialog {
			sse := datastar.NewSSE(w, r)
			if err := sse.MarshalAndPatchSignals(map[string]bool{"assignmentActionCompleted": true}); err != nil {
				logger.Error(err.Error(), "pulje_id", pulje, "event_id", signaler.EventID, "billettholder_id", signaler.BillettholderID)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func tildelingsfeil(w http.ResponseWriter, logger *slog.Logger, err error) {
	switch {
	case errors.Is(err, puljefordeling.ErrPuljeCompleted):
		http.Error(w, "Puljen er publisert og kan ikke endres", http.StatusConflict)
	case errors.Is(err, sql.ErrNoRows):
		http.Error(w, "Fant ikke billettholder, arrangement eller pulje", http.StatusNotFound)
	case errors.Is(err, puljefordeling.ErrUgyldigTildeling):
		http.Error(w, "Ugyldig tildeling", http.StatusBadRequest)
	default:
		logger.Error(err.Error())
		http.Error(w, "Kunne ikke oppdatere tildelingen", http.StatusInternalServerError)
	}
}

func sendTildelingsvarsel(w http.ResponseWriter, r *http.Request, logger *slog.Logger, varsel puljefordeling.Tildelingsvarsel, valg puljefordeling.Tildelingsvalg) {
	sse := datastar.NewSSE(w, r)
	if len(varsel.Tildelinger) == 0 && varsel.Aldersvarsel != "" && varsel.Kapasitetsvarsel == "" {
		if err := sse.MarshalAndPatchSignals(map[string]any{
			"ageWarningText": varsel.Aldersvarsel, "ageWarningBillettholderId": valg.BillettholderID,
			"ageWarningEventId": valg.EventID, "ageWarningPuljeId": string(valg.PuljeID),
			"ageWarningIsPlayer": valg.Role == models.EventPlayerRolePlayer, "ageWarningIsGm": valg.Role == models.EventPlayerRoleGM,
			"ageWarningRole": string(valg.Role), "ageWarningFromAddMenu": valg.FraLeggTil,
			"ageWarningFromEventId": valg.FraEventID, "ageWarningFromManualSeat": valg.FraManuellPlass,
		}); err != nil {
			logger.Error(err.Error())
		}
		return
	}
	if err := sse.PatchElementTempl(puljeTildelingsDialogInnhold(varsel, valg)); err != nil {
		logger.Error(err.Error())
		return
	}
	if err := sse.MarshalAndPatchSignals(map[string]any{"tildelingOpen": true}); err != nil {
		logger.Error(err.Error())
	}
}
