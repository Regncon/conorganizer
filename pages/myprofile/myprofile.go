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

	router.Get("/my-profile", func(w http.ResponseWriter, r *http.Request) {
		my_profile("Min profil").Render(r.Context(), w)
		return
	})

	return nil
}
