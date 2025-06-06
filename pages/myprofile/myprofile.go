package myprofile

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupMyProfileRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, db *sql.DB, logger *slog.Logger) error {

	router.Route("/my-profile", func(ticketRouter chi.Router) {
		ticketRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			layouts.Base(
				"Min profil",
				userctx.GetUserRequestInfo(ctx).IsLoggedIn,
				myProfile(),
			).Render(ctx, w)
			return
		})

		ticketRouter.Get("/my-tickets", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			//Todo get tickets and change what page to render
			layouts.Base(
				"Mine billetter",
				userctx.GetUserRequestInfo(ctx).IsLoggedIn,
				myTickets(),
			).Render(ctx, w)

			layouts.Base(
				"Ingen billetter",
				userctx.GetUserRequestInfo(ctx).IsLoggedIn,
				noTickets(),
			).Render(ctx, w)
		})
	})

	return nil
}
