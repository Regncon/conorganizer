package authctx

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Regncon/conorganizer/components/redirect"
	"github.com/descope/go-sdk/descope/client"
)

const (
	SessionCookieName = "session_token"
	RefreshCookieName = "refresh_token"
)

type ctxKey string

const (
	ctxSessionError ctxKey = "sessionError"
	ctxUserToken    ctxKey = "userToken"
)

func AuthMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//todo: get projectID from env?
			projectID := "P2ufzqahlYUHDIprVXtkuCx8MH5C"
			descopeClient, descopeClientError := client.NewWithConfig(&client.Config{ProjectID: projectID})
			if descopeClientError != nil {
				logger.Error("Failed to create Descope client", slog.String("projectID", projectID), "descopeClientError", descopeClientError)
				ctx := context.WithValue(r.Context(), ctxSessionError, descopeClientError)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Get session and refresh tokens from cookies
			sessionCookie, sessionCookieError := r.Cookie(SessionCookieName)
			if sessionCookieError != nil {
				logger.Error("No session cookie found", "sessionCookieError", sessionCookieError)
				ctx := context.WithValue(r.Context(), ctxSessionError, sessionCookieError)
				redirectUrl := "/auth"
				redirect.Redirect(redirectUrl, "Redirecting to login").Render(ctx, w)
				return
			}

			refreshCookie, refreshCookieError := r.Cookie(RefreshCookieName)
			if refreshCookieError != nil {
				logger.Error("No refresh cookie found", "refreshCookieError", refreshCookieError)
				ctx := context.WithValue(r.Context(), ctxSessionError, refreshCookieError)
				redirectUrl := "/auth"
				redirect.Redirect(redirectUrl, "Redirecting to login").Render(ctx, w)
				return
			}

			_, userToken, validateTokenError := descopeClient.Auth.ValidateAndRefreshSessionWithTokens(
				r.Context(), sessionCookie.Value, refreshCookie.Value)
			if validateTokenError != nil {
				logger.Error("Failed to validate/refresh session", "validateTokenError", validateTokenError)
				ctx := context.WithValue(r.Context(), ctxSessionError, validateTokenError)
				redirectUrl := "/auth"
				redirect.Redirect(redirectUrl, "Redirecting to login").Render(ctx, w)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     SessionCookieName,
				Value:    userToken.JWT,
				Path:     "/",
				Expires:  time.Now().AddDate(1, 0, 0),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			logger.Info("Successfully validated and refreshed session", "email", userToken.Claims["email"])

			ctx := context.WithValue(r.Context(), ctxUserToken, userToken)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
