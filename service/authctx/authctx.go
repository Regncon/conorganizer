package authctx

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/descope/go-sdk/descope/client"
)

const (
	SessionCookieName = "session_token"
	RefreshCookieName = "refresh_token"
)

type authctxKey string

const (
	ctxSessionError authctxKey = "sessionError"
	ctxUserToken    authctxKey = "userToken"
)

const ProjectID = "P2ufzqahlYUHDIprVXtkuCx8MH5C" // TODO: get from env

func AuthMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			descopeClient, descopeClientError := client.NewWithConfig(&client.Config{ProjectID: ProjectID})
			if descopeClientError != nil {
				logger.Error("Failed to create Descope client", slog.String("projectID", ProjectID), "descopeClientError", descopeClientError)
				ctx := context.WithValue(r.Context(), ctxSessionError, descopeClientError)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxSessionError, nil)
			ctx = context.WithValue(ctx, ctxUserToken, nil)

			sessionCookie, sessionCookieError := r.Cookie(SessionCookieName)
			refreshCookie, refreshCookieError := r.Cookie(RefreshCookieName)

			if sessionCookieError == nil && refreshCookieError == nil {
				_, userToken, validateTokenError := descopeClient.Auth.ValidateAndRefreshSessionWithTokens(
					r.Context(), sessionCookie.Value, refreshCookie.Value)
				if validateTokenError == nil && userToken != nil {
					ctx = context.WithValue(ctx, ctxUserToken, userToken)

					http.SetCookie(w, &http.Cookie{
						Name:     SessionCookieName,
						Value:    userToken.JWT,
						Path:     "/",
						Expires:  time.Now().AddDate(1, 0, 0),
						HttpOnly: true,
						// Secure:   true,
						// SameSite: http.SameSiteStrictMode,
						Secure:   false,
						SameSite: http.SameSiteLaxMode,
					})

					logger.Info("Successfully validated and refreshed session", "email", userToken.Claims["email"])
				}

			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
