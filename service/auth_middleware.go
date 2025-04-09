package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Regncon/conorganizer/pages/auth/redirect"
	"github.com/descope/go-sdk/descope"
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

func GetUserTokenFromContext(ctx context.Context) (*descope.Token, error) {
	if errVal := ctx.Value(ctxSessionError); errVal != nil {
		return nil, fmt.Errorf("authentication error: %v", errVal)
	}

	userToken, ok := ctx.Value(ctxUserToken).(*descope.Token)
	if !ok || userToken == nil {
		return nil, fmt.Errorf("user token not found")
	}

	return userToken, nil
}

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
				queryParam := "error=not_logged_in"
				loginURL := "/auth?" + queryParam

				// w.Header().Set("Content-Type", "text/html")
				// w.WriteHeader(http.StatusOK)
				// fmt.Fprintf(w, `
				//     <html>
				//         <head>
				//             <title>Redirecting...</title>
				//             <meta http-equiv="refresh" content="5;url=%s" />
				//         </head>
				//         <body>
				//             <p>You are not logged in. Redirecting to the login page in 5 seconds...</p>
				//             <a href="%s">Click here if you are not redirected.</a>
				//         </body>
				//     </html>
				// `, loginURL, loginURL)
				redirect.Redirect(loginURL).Render(ctx, w)
				return
			}

			refreshCookie, refreshCookieError := r.Cookie(RefreshCookieName)
			if refreshCookieError != nil {
				logger.Error("No refresh cookie found", "refreshCookieError", refreshCookieError)
				ctx := context.WithValue(r.Context(), ctxSessionError, refreshCookieError)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			_, userToken, validateTokenError := descopeClient.Auth.ValidateAndRefreshSessionWithTokens(
				r.Context(), sessionCookie.Value, refreshCookie.Value)
			if validateTokenError != nil {
				logger.Error("Failed to validate/refresh session", "validateTokenError", validateTokenError)
				ctx := context.WithValue(r.Context(), ctxSessionError, validateTokenError)
				next.ServeHTTP(w, r.WithContext(ctx))
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
