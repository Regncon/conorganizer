package authctx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/descope/go-sdk/descope"
	"github.com/descope/go-sdk/descope/client"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	SessionCookieName = "session_token"
	RefreshCookieName = "refresh_token"
)

type sessionErrorKey string
type userTokenKey string

const (
	ctxSessionError sessionErrorKey = "sessionError"
	ctxUserToken    userTokenKey    = "userToken"
)

const authCookieMaxAgeSeconds = 365 * 24 * 60 * 60

type SessionValidator interface {
	ValidateSessionWithToken(ctx context.Context, sessionToken string) (bool, *descope.Token, error)
	RefreshSessionWithToken(ctx context.Context, refreshToken string) (bool, *descope.Token, error)
}

func NewSessionValidatorFromEnv() (SessionValidator, error) {
	projectID := os.Getenv("DESCOPE_PROJECT_ID")
	descopeClient, err := client.NewWithConfig(&client.Config{ProjectID: projectID})
	if err != nil {
		return nil, fmt.Errorf("failed to create Descope client for project %q: %w", projectID, err)
	}
	return descopeClient.Auth, nil
}

func SetAuthCookies(w http.ResponseWriter, r *http.Request, sessionToken, refreshToken string) {
	SetSessionCookie(w, r, sessionToken)
	if refreshToken != "" {
		http.SetCookie(w, authCookie(r, RefreshCookieName, refreshToken, authCookieMaxAgeSeconds))
	}
}

func SetSessionCookie(w http.ResponseWriter, r *http.Request, sessionToken string) {
	if sessionToken == "" {
		return
	}
	http.SetCookie(w, authCookie(r, SessionCookieName, sessionToken, authCookieMaxAgeSeconds))
}

func ClearAuthCookies(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, expiredAuthCookie(r, SessionCookieName))
	http.SetCookie(w, expiredAuthCookie(r, RefreshCookieName))
}

func authCookie(r *http.Request, name, value string, maxAge int) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().AddDate(1, 0, 0),
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   requestIsSecure(r),
		SameSite: http.SameSiteLaxMode,
	}
}

func expiredAuthCookie(r *http.Request, name string) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   requestIsSecure(r),
		SameSite: http.SameSiteLaxMode,
	}
}

func requestIsSecure(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	for proto := range strings.SplitSeq(r.Header.Get("X-Forwarded-Proto"), ",") {
		if strings.EqualFold(strings.TrimSpace(proto), "https") {
			return true
		}
	}
	return strings.Contains(strings.ToLower(r.Header.Get("Forwarded")), "proto=https")
}

func AuthMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	logger = logger.With("component", "auth")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
			sessionValidator, descopeClientError := NewSessionValidatorFromEnv()
			if descopeClientError != nil {
				logger.Error(descopeClientError.Error(), "request_id", requestID)
				ctx := context.WithValue(r.Context(), ctxSessionError, descopeClientError)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, ctxSessionError, nil)
			ctx = context.WithValue(ctx, ctxUserToken, nil)

			userToken, refreshed, validateTokenError := validateRequestSession(r.Context(), sessionValidator, r)
			if validateTokenError == nil && userToken != nil {
				ctx = context.WithValue(ctx, ctxUserToken, userToken)
				if refreshed && userToken.JWT != "" {
					SetSessionCookie(w, r, userToken.JWT)
					logger.Debug("Successfully refreshed session", "request_id", requestID)
				}
			}
			if validateTokenError != nil {
				logger.Warn(fmt.Errorf("failed to validate session: %w", validateTokenError).Error(), "request_id", requestID)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func validateRequestSession(ctx context.Context, sessionValidator SessionValidator, r *http.Request) (*descope.Token, bool, error) {
	sessionCookie, sessionCookieError := r.Cookie(SessionCookieName)
	refreshCookie, refreshCookieError := r.Cookie(RefreshCookieName)

	var sessionErr error
	if sessionCookieError == nil && sessionCookie.Value != "" {
		userOK, userToken, err := sessionValidator.ValidateSessionWithToken(ctx, sessionCookie.Value)
		if err == nil && userOK && userToken != nil {
			return userToken, false, nil
		}
		if err == nil {
			err = fmt.Errorf("session token rejected")
		}
		sessionErr = err
	}

	if refreshCookieError == nil && refreshCookie.Value != "" {
		userOK, userToken, err := sessionValidator.RefreshSessionWithToken(ctx, refreshCookie.Value)
		if err == nil && userOK && userToken != nil {
			return userToken, true, nil
		}
		if err == nil {
			err = fmt.Errorf("refresh token rejected")
		}
		if sessionErr != nil {
			return nil, false, fmt.Errorf("session validation failed: %w; refresh failed: %w", sessionErr, err)
		}
		return nil, false, err
	}

	return nil, false, sessionErr
}
