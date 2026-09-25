package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
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

// lesTildelingssignaler reads the assignment signals shared by the preview and
// the assign endpoints.
func lesTildelingssignaler(w http.ResponseWriter, r *http.Request) (puljeTildelingssignaler, puljefordeling.Tildelingsvalg, bool) {
	var signaler puljeTildelingssignaler
	if err := datastar.ReadSignals(r, &signaler); err != nil {
		http.Error(w, "Ugyldige tildelingsdata", http.StatusBadRequest)
		return signaler, puljefordeling.Tildelingsvalg{}, false
	}
	pulje, ok := models.ParsePulje(signaler.PuljeID)
	if !ok {
		http.Error(w, "Ugyldig pulje", http.StatusBadRequest)
		return signaler, puljefordeling.Tildelingsvalg{}, false
	}
	role := signaler.Role
	if role == "" {
		role = models.EventPlayerRolePlayer
	}
	return signaler, puljefordeling.Tildelingsvalg{
		PuljeID: pulje, EventID: signaler.EventID, BillettholderID: signaler.BillettholderID,
		Role: role, FraLeggTil: signaler.FraLeggTil,
		Bekreftelse: signaler.Bekreftelse, AlderBekreftet: signaler.AlderBekreftet,
		FraEventID: signaler.FraEventID, FraManuellPlass: signaler.FraManuellPlass,
	}, true
}

// puljeForhandsvisningsHandler opens the "Er du sikker?" dialog for a manual
// assignment, listing who else in the pulje would change seats. Nothing is saved.
func puljeForhandsvisningsHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc {
	logger = logger.With("component", "admin_tildelinger")
	return func(w http.ResponseWriter, r *http.Request) {
		_, valg, ok := lesTildelingssignaler(w, r)
		if !ok {
			return
		}
		valg.Bekreftelse = ""
		varsel, err := puljefordeling.ForhandsvisTildeling(db, valg)
		if err != nil {
			tildelingsfeil(w, logger, err)
			return
		}
		if varsel == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		sendTildelingsdialog(w, r, logger, puljeTildelingsDialogInnhold(*varsel, valg))
	}
}

// puljeFjerningsforhandsvisningHandler opens the "Er du sikker?" dialog for
// removing a manual Player seat or a GM, listing who else would change seats.
func puljeFjerningsforhandsvisningHandler(db *sql.DB, logger *slog.Logger, role models.EventPlayerRole) http.HandlerFunc {
	logger = logger.With("component", "admin_tildelinger")
	return func(w http.ResponseWriter, r *http.Request) {
		pulje, ok := models.ParsePulje(chi.URLParam(r, "pulje"))
		if !ok {
			http.Error(w, "Ugyldig pulje", http.StatusBadRequest)
			return
		}
		billettholderID, err := strconv.Atoi(chi.URLParam(r, "billettholderId"))
		if err != nil || billettholderID <= 0 {
			http.Error(w, "Ugyldig billettholder", http.StatusBadRequest)
			return
		}
		eventID := chi.URLParam(r, "event")
		varsel, err := puljefordeling.ForhandsvisFjerning(db, pulje, eventID, billettholderID, role)
		if err != nil {
			tildelingsfeil(w, logger, err)
			return
		}
		sendTildelingsdialog(w, r, logger, puljeFjerningsDialogInnhold(varsel, puljeFjerningsURL(pulje, eventID, billettholderID, role)))
	}
}

// puljeFjerningsURL is the endpoint that removes the seat or GM once confirmed.
func puljeFjerningsURL(pulje models.Pulje, eventID string, billettholderID int, role models.EventPlayerRole) string {
	url := fmt.Sprintf("/admin/api/puljefordeling/%s/%s/%d", pulje, eventID, billettholderID)
	if role == models.EventPlayerRoleGM {
		url += "/gm"
	}
	return url
}

func sendTildelingsdialog(w http.ResponseWriter, r *http.Request, logger *slog.Logger, innhold templ.Component) {
	sse := datastar.NewSSE(w, r)
	if err := sse.PatchElementTempl(innhold); err != nil {
		logger.Error(err.Error())
		return
	}
	if err := sse.MarshalAndPatchSignals(map[string]any{"tildelingOpen": true}); err != nil {
		logger.Error(err.Error())
	}
}

func puljeTildelingsHandler(db *sql.DB, liveManager *live.Manager, logger *slog.Logger) http.HandlerFunc {
	logger = logger.With("component", "admin_tildelinger")
	return func(w http.ResponseWriter, r *http.Request) {
		signaler, valg, ok := lesTildelingssignaler(w, r)
		if !ok {
			return
		}
		pulje := valg.PuljeID
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
