package myprofile

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/checkIn"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupMyProfileRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, db *sql.DB, logger *slog.Logger) error {

	router.Route("/my-profile", func(ticketRouter chi.Router) {
		ticketRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			if err := layouts.Base(
				"Min profil",
				userctx.GetUserRequestInfo(ctx),
				myProfile(),
			).Render(ctx, w); err != nil {
				logger.Error("Failed to render profile page", "error", err)
			}

		})

		ticketRouter.Get("/my-tickets", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()

			tickets, ticketsErr := checkIn.GetTicketsFromCheckIn(logger, userctx.GetUserRequestInfo(ctx).Email)

			if ticketsErr != nil {
				logger.Error("Failed to get tickets from check-in", "ticketsErr", ticketsErr)

				if err := layouts.Base(
					"Ingen billetter",
					userctx.GetUserRequestInfo(ctx),
					noTickets(),
				).Render(ctx, w); err != nil {
					logger.Error("Failed to render no tickets page", "error", err)
				}
				return
			}

			if len(tickets) == 0 {
				if err := layouts.Base(
					"Ingen billetter",
					userctx.GetUserRequestInfo(ctx),
					noTickets(),
				).Render(ctx, w); err != nil {
					logger.Error("Failed to render no tickets page", "error", err)
				}
				return
			}

			if err := layouts.Base(
				"Mine billetter",
				userctx.GetUserRequestInfo(ctx),
				myTickets(),
			).Render(ctx, w); err != nil {
				logger.Error("Failed to render tickets page", "error", err)
			}

		})
	})

	return nil
}
