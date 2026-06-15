package userctx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/go-chi/chi/v5/middleware"
)

func UserMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	logger = logger.With("component", "user")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestID := middleware.GetReqID(ctx)
			userInfo := GetUserRequestInfo(ctx)
			if userInfo.IsLoggedIn {
				logger.Debug("User is logged in", "request_id", requestID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if !userInfo.IsLoggedIn {
				logger.Warn("User is not logged in", "request_id", requestID, "path", r.URL.Path)
				w.WriteHeader(http.StatusUnauthorized)
				if err := layouts.Base("Logg inn", requestctx.UserRequestInfo{}, Unauthenticated()).Render(r.Context(), w); err != nil {
					logger.Error(fmt.Errorf("failed to render unauthenticated page: %w", err).Error(), "request_id", requestID)
				}
				return
			}

		})
	}
}

func AdminForbiddenHandler(logger *slog.Logger) http.HandlerFunc {
	logger = logger.With("component", "user")
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetReqID(r.Context())
		userInfo := GetUserRequestInfo(r.Context())

		w.WriteHeader(http.StatusForbidden)
		if err := layouts.Base("Ingen tilgang", userInfo, authctx.Forbidden()).Render(r.Context(), w); err != nil {
			logger.Error(fmt.Errorf("failed to render forbidden page: %w", err).Error(), "request_id", requestID)
		}
	}
}

func GetUserRequestInfo(ctx context.Context) requestctx.UserRequestInfo {

	userToken, userTokenErr := authctx.GetUserTokenFromContext(ctx)
	userId, _ := authctx.GetUserIDFromToken(ctx)
	email, _ := authctx.GetEmailFromToken(ctx)
	isAdmin := authctx.GetAdminFromUserToken(ctx)

	return requestctx.UserRequestInfo{
		IsLoggedIn: userTokenErr == nil || userToken != nil,
		Id:         userId,
		Email:      email,
		IsAdmin:    isAdmin,
	}
}

func GetUserIDFromExternalID(externalID string, db *sql.DB, logger *slog.Logger) (string, error) {
	var userID string
	userQuery := "SELECT id FROM users WHERE external_id = ?"
	userRow := db.QueryRow(userQuery, externalID)
	if userRowErr := userRow.Scan(&userID); userRowErr != nil {
		return "", fmt.Errorf("failed to find user external_id %q: %w", externalID, userRowErr)
	}
	return userID, nil
}

func GetUserIDFromExternalIDFromContext(ctx context.Context, db *sql.DB, logger *slog.Logger) (string, error) {
	externalID := GetUserRequestInfo(ctx).Id
	return GetUserIDFromExternalID(externalID, db, logger)
}
