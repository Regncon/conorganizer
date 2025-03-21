package auth

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetupAuthRoute(router chi.Router, logger *slog.Logger) error {
	// eventLayoutRoute(router, db, err)
	// eventForm.NewEventLayoutRoute(router, db, err)

	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			login().Render(r.Context(), w)
		})

		authRouter.Get("/logout", func(w http.ResponseWriter, r *http.Request) {

			http.Redirect(w, r, "/", http.StatusSeeOther)
		})
		// eventForm.SetupExampleInlineValidation(db, newRouter, logger)

	})
	return nil
}
