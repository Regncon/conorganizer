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
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5"
)

// SetupPuljefordelingEmulateRoute wires the read-only emulate page.
// Authorization is handled by the admin router this is mounted on.
func SetupPuljefordelingEmulateRoute(router chi.Router, db *sql.DB, logger *slog.Logger) {
	logger = logger.With("component", "puljefordeling_emulate")

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		userInfo := userctx.GetUserRequestInfo(r.Context())

		em, emErr := puljefordeling.EmulateSeatings(db)
		if emErr != nil {
			logger.Error(fmt.Errorf("emulation failed: %w", emErr).Error(), "user_id", userInfo.Id)
		}

		if err := layouts.Base(
			"Puljefordeling – Emulering",
			userInfo,
			emulatePage(em, emErr),
		).Render(r.Context(), w); err != nil {
			logger.Error(fmt.Errorf("failed to render emulate page: %w", err).Error(), "user_id", userInfo.Id)
		}
	})
}
