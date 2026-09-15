package admin

import (
	"database/sql"
	"log/slog"

	"github.com/Regncon/conorganizer/components/formsubmission"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/go-chi/chi/v5"
)

// approvalEventPlayersRoute wires the approval page's player/GM assignment
// endpoints. It lives apart from the rest of the admin router so the placement
// rules it enforces can be exercised without mounting the whole admin tree.
func approvalEventPlayersRoute(router chi.Router, db *sql.DB, liveManager *live.Manager, logger *slog.Logger) {
	router.Route("/event-players", func(router chi.Router) {
		router.Post("/post/add_first_choice", tildelingsHandler(db, liveManager, logger, tildelingsrute{
			URL: formsubmission.AddFirstChoiceURL, Role: models.EventPlayerRolePlayer, FraLeggTil: true, Forstevalg: true,
		}))
		router.Post("/post/add_gm", tildelingsHandler(db, liveManager, logger, tildelingsrute{
			URL: formsubmission.AddGMURL, Role: models.EventPlayerRoleGM, FraLeggTil: true,
		}))
		router.Put("/update_status", tildelingsHandler(db, liveManager, logger, tildelingsrute{
			URL: formsubmission.UpdatePlayerStatusURL, FraLeggTil: true, Oppdater: true,
		}))
	})
}
