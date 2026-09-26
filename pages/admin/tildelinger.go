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

// readAssignmentSignals reads the assignment signals shared by the preview and
// the assign endpoints.
func readAssignmentSignals(w http.ResponseWriter, r *http.Request) (puljeTildelingssignaler, puljefordeling.Tildelingsvalg, bool) {
	var signals puljeTildelingssignaler
	if err := datastar.ReadSignals(r, &signals); err != nil {
		http.Error(w, "Ugyldige tildelingsdata", http.StatusBadRequest)
		return signals, puljefordeling.Tildelingsvalg{}, false
	}
	pulje, ok := models.ParsePulje(signals.PuljeID)
	if !ok {
		http.Error(w, "Ugyldig pulje", http.StatusBadRequest)
		return signals, puljefordeling.Tildelingsvalg{}, false
	}
	role := signals.Role
	if role == "" {
		role = models.EventPlayerRolePlayer
	}
	return signals, puljefordeling.Tildelingsvalg{
		PuljeID: pulje, EventID: signals.EventID, BillettholderID: signals.BillettholderID,
		Role: role, FraLeggTil: signals.FraLeggTil,
		Bekreftelse: signals.Bekreftelse, AlderBekreftet: signals.AlderBekreftet,
		FraEventID: signals.FraEventID, FraManuellPlass: signals.FraManuellPlass,
	}, true
}

// puljeAssignPreviewHandler opens the "Er du sikker?" dialog for a manual
// assignment, listing who else in the pulje would change seats. Nothing is saved.
func puljeAssignPreviewHandler(db *sql.DB, logger *slog.Logger) http.HandlerFunc {
	logger = logger.With("component", "admin_tildelinger")
	return func(w http.ResponseWriter, r *http.Request) {
		_, assignment, ok := readAssignmentSignals(w, r)
		if !ok {
			return
		}
		assignment.Bekreftelse = ""
		warning, err := puljefordeling.PreviewAssignment(db, assignment)
		if err != nil {
			tildelingsfeil(w, logger, err)
			return
		}
		if warning == nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		sendAssignmentDialog(w, r, logger, puljeTildelingsDialogInnhold(*warning, assignment))
	}
}

// puljeRemovalPreviewHandler opens the "Er du sikker?" dialog for
// removing a manual Player seat or a GM, listing who else would change seats.
func puljeRemovalPreviewHandler(db *sql.DB, logger *slog.Logger, role models.EventPlayerRole) http.HandlerFunc {
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
		preview, err := puljefordeling.PreviewRemoval(db, pulje, eventID, billettholderID, role)
		if err != nil {
			tildelingsfeil(w, logger, err)
			return
		}
		sendAssignmentDialog(w, r, logger, puljeRemovalDialogContent(preview, puljeRemovalURL(pulje, eventID, billettholderID, role)))
	}
}

// puljeRemovalURL is the endpoint that removes the seat or GM once confirmed.
func puljeRemovalURL(pulje models.Pulje, eventID string, billettholderID int, role models.EventPlayerRole) string {
	url := fmt.Sprintf("/admin/api/puljefordeling/%s/%s/%d", pulje, eventID, billettholderID)
	if role == models.EventPlayerRoleGM {
		url += "/gm"
	}
	return url
}

func sendAssignmentDialog(w http.ResponseWriter, r *http.Request, logger *slog.Logger, content templ.Component) {
	sse := datastar.NewSSE(w, r)
	if err := sse.PatchElementTempl(content); err != nil {
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
		signaler, valg, ok := readAssignmentSignals(w, r)
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
