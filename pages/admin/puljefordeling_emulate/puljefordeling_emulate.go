// Package puljefordeling_emulate hosts the admin page for previewing how
// participants would be distributed across events in each pulje. It is a
// read-only "what-if" tool — it never writes assignments to the database.
package puljefordeling_emulate

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

// runStore holds the per-pulje late-boost toggles posted from the page.
type runStore struct {
	BoostFredagKveld  bool `json:"boostFredagKveld"`
	BoostLordagMorgen bool `json:"boostLordagMorgen"`
	BoostLordagKveld  bool `json:"boostLordagKveld"`
	BoostSondagMorgen bool `json:"boostSondagMorgen"`
}

func (s runStore) boostMap() map[models.Pulje]bool {
	return map[models.Pulje]bool{
		models.PuljeFredagKveld:  s.BoostFredagKveld,
		models.PuljeLordagMorgen: s.BoostLordagMorgen,
		models.PuljeLordagKveld:  s.BoostLordagKveld,
		models.PuljeSondagMorgen: s.BoostSondagMorgen,
	}
}

// boostSignal returns the Datastar signal name used to toggle late boost for a
// pulje, matching the json tags on runStore.
func boostSignal(pulje models.Pulje) string {
	switch pulje {
	case models.PuljeFredagKveld:
		return "boostFredagKveld"
	case models.PuljeLordagMorgen:
		return "boostLordagMorgen"
	case models.PuljeLordagKveld:
		return "boostLordagKveld"
	case models.PuljeSondagMorgen:
		return "boostSondagMorgen"
	default:
		return ""
	}
}

// SetupPuljefordelingEmulateRoute wires the emulate page and its re-run
// endpoint. Authorization is handled by the admin router this is mounted on.
func SetupPuljefordelingEmulateRoute(router chi.Router, db *sql.DB, logger *slog.Logger) {
	logger = logger.With("component", "puljefordeling_emulate")

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userInfo := userctx.GetUserRequestInfo(r.Context())

		em, emErr := puljefordeling.EmulateSeatings(db, nil)
		if emErr != nil {
			logger.Error(fmt.Errorf("initial emulation failed: %w", emErr).Error(), "user_id", userInfo.Id)
		}

		if err := layouts.Base(
			"Puljefordeling – Emulering",
			userInfo,
			emulatePage(em, emErr),
		).Render(r.Context(), w); err != nil {
			logger.Error(fmt.Errorf("failed to render emulate page: %w", err).Error(), "user_id", userInfo.Id)
		}
	})

	router.Post("/api/run", func(w http.ResponseWriter, r *http.Request) {
		store := &runStore{}
		if err := datastar.ReadSignals(r, store); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sse := datastar.NewSSE(w, r)

		em, emErr := puljefordeling.EmulateSeatings(db, store.boostMap())
		if emErr != nil {
			logger.Error(fmt.Errorf("emulation re-run failed: %w", emErr).Error())
		}

		if err := sse.PatchElementTempl(emulationResults(em, emErr)); err != nil {
			_ = sse.ConsoleError(err)
			return
		}
	})
}
