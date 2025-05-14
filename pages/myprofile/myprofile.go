package myprofile

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupMyProfileRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, db *sql.DB, logger *slog.Logger) error {

	router.Route("/my-profile", func(ticketRouter chi.Router) {
		ticketRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			myProfile().Render(r.Context(), w)
			return
		})

		ticketRouter.Get("/my-tickets", func(w http.ResponseWriter, r *http.Request) {
			noTicket().Render(r.Context(), w)
		})
	})

	return nil
}
