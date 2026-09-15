package admin

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/components/formsubmission"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	datastar "github.com/starfederation/datastar-go/datastar"
)

type tildelingssignaler struct {
	BillettholderID int                    `json:"assignmentBillettholderId"`
	EventID         string                 `json:"assignmentEventId"`
	PuljeID         string                 `json:"assignmentPuljeId"`
	Role            models.EventPlayerRole `json:"assignmentRole"`
	FraLeggTil      bool                   `json:"assignmentFromAddMenu"`
	Bekreftelse     string                 `json:"assignmentConfirmation"`
	AlderBekreftet  bool                   `json:"assignmentAgeConfirmed"`
	Fjern           bool                   `json:"assignmentRemove"`
	IsPlayer        bool                   `json:"assignmentIsPlayer"`
	IsGM            bool                   `json:"assignmentIsGm"`
}

type tildelingsrute struct {
	URL        string
	Role       models.EventPlayerRole
	FraLeggTil bool
	Forstevalg bool
	Oppdater   bool
}

func tildelingsHandler(db *sql.DB, liveManager *live.Manager, logger *slog.Logger, rute tildelingsrute) http.HandlerFunc {
	logger = logger.With("component", "admin_tildelinger")
	return func(w http.ResponseWriter, r *http.Request) {
		var signaler tildelingssignaler
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
		if rute.Role != "" {
			role = rute.Role
		} else if role == "" {
			role = models.EventPlayerRolePlayer
			if signaler.IsGM {
				role = models.EventPlayerRoleGM
			}
			if rute.Oppdater && !signaler.IsPlayer && !signaler.IsGM {
				signaler.Fjern = true
			}
		}
		if signaler.Fjern && rute.Oppdater {
			if err := puljefordeling.FjernTildeling(db, pulje, signaler.EventID, signaler.BillettholderID, role); err != nil {
				tildelingsfeil(w, logger, err)
				return
			}
		} else {
			valg := puljefordeling.Tildelingsvalg{
				PuljeID: pulje, EventID: signaler.EventID, BillettholderID: signaler.BillettholderID,
				Role: role, FraLeggTil: rute.FraLeggTil || signaler.FraLeggTil,
				Bekreftelse: signaler.Bekreftelse, AlderBekreftet: signaler.AlderBekreftet,
				Forstevalg: rute.Forstevalg,
			}
			varsel, err := puljefordeling.TildelBillettholder(db, valg)
			if err != nil {
				tildelingsfeil(w, logger, err)
				return
			}
			if varsel != nil {
				sendTildelingsvarsel(w, r, logger, *varsel, valg, rute.URL)
				return
			}
		}
		if err := liveManager.Broadcast(r.Context(), live.BucketEvents, live.BucketInterests, live.BucketRooms); err != nil {
			logger.Error(err.Error(), "pulje_id", pulje, "event_id", signaler.EventID, "billettholder_id", signaler.BillettholderID)
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

func sendTildelingsvarsel(w http.ResponseWriter, r *http.Request, logger *slog.Logger, varsel puljefordeling.Tildelingsvarsel, valg puljefordeling.Tildelingsvalg, retryURL string) {
	sse := datastar.NewSSE(w, r)
	if len(varsel.Tildelinger) == 0 && varsel.Aldersvarsel != "" && varsel.Kapasitetsvarsel == "" {
		if err := sse.MarshalAndPatchSignals(map[string]any{
			"ageWarningText": varsel.Aldersvarsel, "ageWarningBillettholderId": valg.BillettholderID,
			"ageWarningEventId": valg.EventID, "ageWarningPuljeId": string(valg.PuljeID),
			"ageWarningIsPlayer": valg.Role == models.EventPlayerRolePlayer, "ageWarningIsGm": valg.Role == models.EventPlayerRoleGM,
			"ageWarningRole": string(valg.Role), "ageWarningFromAddMenu": valg.FraLeggTil,
			"ageWarningMethod": r.Method, "ageWarningUrl": retryURL,
		}); err != nil {
			logger.Error(err.Error())
		}
		return
	}
	if err := sse.PatchElementTempl(formsubmission.TildelingsDialogInnhold(varsel, valg, r.Method, retryURL)); err != nil {
		logger.Error(err.Error())
		return
	}
	if err := sse.MarshalAndPatchSignals(map[string]any{"tildelingOpen": true}); err != nil {
		logger.Error(err.Error())
	}
}
