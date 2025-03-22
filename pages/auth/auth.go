package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/service"
	"github.com/descope/go-sdk/descope"
	"github.com/go-chi/chi/v5"
)

func SetupAuthRoute(router chi.Router, logger *slog.Logger) error {

	// extract from request authorization header. The above sample code sends the the session token in authorization header.
	// sessionToken := r.Header.Get("Authorization")

	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			loginForm().Render(r.Context(), w)
		})

		authRouter.Get("/test", service.ValidateSession(logger, func(w http.ResponseWriter, r *http.Request) {

			if errVal := r.Context().Value(service.CtxSessionError); errVal != nil {
				http.Error(w, fmt.Sprintf("Session error: %v", errVal), http.StatusUnauthorized)
				return
			}

			userToken, ok := r.Context().Value(service.CtxUserToken).(*descope.Token)
			fmt.Printf("userToken: %v\n", userToken)
			fmt.Printf("ok: %v\n", ok)
			fmt.Printf("userToken.Claims: %v\n", userToken.Claims["email"])
			if !ok || userToken == nil {
				http.Error(w, "User token not found", http.StatusUnauthorized)
				return
			}

			w.Write([]byte(fmt.Sprintf("Test, authorized: true, email: %v", userToken.Claims["email"])))
		}))

	})

	return nil
}
