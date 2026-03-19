package userctx

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/go-chi/chi/v5/middleware"
)

func UserMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	userLogger := logger.With("component", "user")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestID := middleware.GetReqID(ctx)
			userInfo := GetUserRequestInfo(ctx)
			if userInfo.IsLoggedIn {
				userLogger.Debug("User is logged in", "request_id", requestID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if !userInfo.IsLoggedIn {
				userLogger.Warn("User is not logged in", "request_id", requestID, "path", r.URL.Path)
				w.WriteHeader(http.StatusUnauthorized)
				if err := layouts.Base("Unauthorized", requestctx.UserRequestInfo{}, Unauthorized()).Render(r.Context(), w); err != nil {
					userLogger.Error("Failed to render unauthorized page", "error", err, "request_id", requestID)
				}
				return
			}

		})
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

func GetIdFromUserIdInDb(userId string, db *sql.DB, logger *slog.Logger) (string, error) {
	var userDbId string
	userQuery := "SELECT id FROM users WHERE user_id = ?"
	userRow := db.QueryRow(userQuery, userId)
	if userRowErr := userRow.Scan(&userDbId); userRowErr != nil {
		logger.With("component", "user").Error("Failed to find user", "user_id", userId, "error", userRowErr)
		return "", userRowErr
	}
	return userDbId, nil
}

func GetIdFromUserIdInDbFromContext(ctx context.Context, db *sql.DB, logger *slog.Logger) (string, error) {
	userId := GetUserRequestInfo(ctx).Id
	return GetIdFromUserIdInDb(userId, db, logger)
}
