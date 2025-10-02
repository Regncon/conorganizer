package profilepage

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	ticketspage "github.com/Regncon/conorganizer/pages/profile/tickets"
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
				"Min profil side",
				user,
				ProfilePage(user, db, logger),
			).Render(ctx, w)
		})

		profileRouter.Get("/tickets", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			var user = userctx.GetUserRequestInfo(ctx)

			// todo: enable error checking
			tickets, _ := checkIn.GetTicketsFromCheckIn(logger, user.Email)

			layouts.Base(
				"Mine biletter",
				user,
				ticketspage.ProfileTicketsPage(tickets),
			).Render(ctx, w)
		})
	})

	return nil
}
