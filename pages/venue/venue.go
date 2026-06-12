package venue

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5"
)

func SetupVenueRoute(router chi.Router, liveManager *live.Manager, db *sql.DB, eventImageDir *string, logger *slog.Logger) error {
	router.Route("/venue", func(venueRouter chi.Router) {
		venueRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			userInfo := userctx.GetUserRequestInfo(r.Context())
			if err := layouts.Base(
				"Oversikt over rom",
				userInfo,
				venuePage(),
			).Render(ctx, w); err != nil {
				logger.Error(fmt.Errorf("failed to render venue page: %w", err).Error(), "user_id", userInfo.Id)
			}
		})
	})
	return nil
}
