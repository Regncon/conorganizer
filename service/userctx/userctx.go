package userctx

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/requestctx"
)

func UserMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userInfo := GetUserRequestInfo(ctx)
			if userInfo.IsLoggedIn {
				logger.Info("User is logged in")
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if !userInfo.IsLoggedIn {
				logger.Warn("User is not logged in")
				w.WriteHeader(http.StatusUnauthorized)
				layouts.Base("Unauthorized", requestctx.UserRequestInfo{}, Unauthorized()).Render(r.Context(), w)
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
		logger.Error("Failed to find user", "user_id", userId, "error", userRowErr)
		return "", userRowErr
	}
	return userDbId, nil
}
