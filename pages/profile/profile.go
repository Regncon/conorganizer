package profile

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	mytickets "github.com/Regncon/conorganizer/pages/profile/myTickets"
	"github.com/Regncon/conorganizer/service/checkIn"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupProfileRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, db *sql.DB, logger *slog.Logger) error {
	router.Route("/profile", func(profileRouter chi.Router) {
		profileRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			var user = userctx.GetUserRequestInfo(ctx)

			layouts.Base(
				"Min profil",
				userctx.GetUserRequestInfo(ctx),
				profilePage(user, db, logger),
			).Render(ctx, w)

		})

		profileRouter.Get("/events", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			var user = userctx.GetUserRequestInfo(ctx)

			layouts.Base(
				"Mine events",
				userctx.GetUserRequestInfo(ctx),
				profilePage(user, db, logger),
			).Render(ctx, w)
		})

		profileRouter.Get("/tickets", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			tickets, ticketsErr := checkIn.GetTicketsFromCheckIn(logger, userctx.GetUserRequestInfo(ctx).Email)

			if ticketsErr != nil {
				logger.Error("Failed to get tickets from check-in", "ticketsErr", ticketsErr)

				layouts.Base(
					"Ingen billetter",
					userctx.GetUserRequestInfo(ctx),
					mytickets.MyTicketsPage(tickets),
				).Render(ctx, w)
				return
			}
		})
	})

	return nil
}
