package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/descope/go-sdk/descope/client"
)

// Definer egne typer for context-nøkler for å unngå kollisjoner
type ctxKey string

const (
	CtxSessionError ctxKey = "sessionError"
	CtxUserToken    ctxKey = "userToken"
)

func ValidateSession(logger *slog.Logger, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		projectID := "P2ufzqahlYUHDIprVXtkuCx8MH5C"
		descopeClient, validateAndRefreshError := client.NewWithConfig(&client.Config{ProjectID: projectID})
		if validateAndRefreshError != nil {
			logger.Error("Failed to create Descope client", slog.String("projectID", projectID), validateAndRefreshError)
			ctx := context.WithValue(r.Context(), CtxSessionError, validateAndRefreshError)
			next(w, r.WithContext(ctx))
			return
		}
		logger.Info("Descope client created successfully")

		cookieSession, validateAndRefreshError := r.Cookie("session_token")
		if validateAndRefreshError != nil {
			ctx := context.WithValue(r.Context(), CtxSessionError, validateAndRefreshError)
			next(w, r.WithContext(ctx))
			return
		}

		cookieRefresh, validateAndRefreshError := r.Cookie("refresh_token")
		if validateAndRefreshError != nil {
			ctx := context.WithValue(r.Context(), CtxSessionError, validateAndRefreshError)
			next(w, r.WithContext(ctx))
			return
		}
		fmt.Printf("cookieSession: %v\n", cookieSession.Value)
		fmt.Printf("cookieRefresh: %v\n", cookieRefresh.Value)
		_, userToken, validateAndRefreshError := descopeClient.Auth.ValidateAndRefreshSessionWithTokens(
			r.Context(), cookieSession.Value, cookieRefresh.Value)
		if validateAndRefreshError != nil {
			logger.Error("Could not validate user session", "validateAndRefreshError", validateAndRefreshError)
			ctx := context.WithValue(r.Context(), CtxSessionError, validateAndRefreshError)
			next(w, r.WithContext(ctx))
			return
		}

		cookie := &http.Cookie{
			Name:    "session_token",
			Value:   userToken.JWT,
			Path:    "/",
			Expires: time.Unix(userToken.Expiration, 0),
		}
		http.SetCookie(w, cookie)
		logger.Info("Successfully validated user session", slog.Any("email", userToken.Claims["email"]))

		ctx := context.WithValue(r.Context(), CtxUserToken, userToken)
		next(w, r.WithContext(ctx))
	}
}
