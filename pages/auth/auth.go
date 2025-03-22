package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/supertokens/supertokens-golang/recipe/session"
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

		authRouter.Get("/test", session.VerifySession(nil, func(w http.ResponseWriter, r *http.Request) {
			sessionContainer := session.GetSessionFromRequestContext(r.Context())

			userID := sessionContainer.GetUserID()

			fmt.Println(userID)
			w.Write([]byte("test"))
		}))

		authRouter.Get("/signin", func(w http.ResponseWriter, r *http.Request) {

		})

		// eventForm.SetupExampleInlineValidation(db, newRouter, logger)

	})
	return nil
}
