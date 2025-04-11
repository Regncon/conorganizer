package auth

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/pages/auth/redirect"
	"github.com/Regncon/conorganizer/service"
	"github.com/go-chi/chi/v5"
)

func SetupAuthRoute(router chi.Router, logger *slog.Logger) error {
	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			loginForm().Render(r.Context(), w)
		})

		authRouter.Group(func(protectedRoute chi.Router) {
			protectedRoute.Use(service.AuthMiddleware(logger))

			protectedRoute.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				userToken, err := service.GetUserTokenFromContext(r.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				isAdmin := service.GetAdminFromUserToken(r.Context())
				layouts.Base("Is logged in test").Render(r.Context(), w)
				w.Write(fmt.Appendf(nil, "Test successful! Authenticated as: %v, and is admin: %v", userToken.Claims["email"], isAdmin))
			})

		})

		authRouter.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     service.SessionCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     service.RefreshCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			redirectUrl := "/"
			redirect.Redirect(redirectUrl, "Logging you out").Render(r.Context(), w)
		})
	})

	return nil
}
