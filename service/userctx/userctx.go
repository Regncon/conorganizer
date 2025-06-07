package userctx

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/descope/go-sdk/descope/client"
)

type userctxKey struct{}

var userContextKey = userctxKey{}

type UserRequestInfo struct {
	IsLoggedIn bool
	Id         string
	Email      string
}

func IsLoggedInMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isLoggedIn := false

			sessionCookie, sessionCookieError := r.Cookie(authctx.SessionCookieName)
			if sessionCookieError == nil {
				descopeClient, descopeClientError := client.NewWithConfig(&client.Config{ProjectID: authctx.ProjectID})
				if descopeClientError == nil {
					_, userToken, validateTokenError := descopeClient.Auth.ValidateSessionWithToken(r.Context(), sessionCookie.Value)
					if validateTokenError == nil && userToken != nil {
						isLoggedIn = true
					}
				}
			}

			userInfo := UserRequestInfo{
				IsLoggedIn: isLoggedIn,
			}

			ctx := context.WithValue(r.Context(), userContextKey, userInfo)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserRequestInfo(ctx context.Context) UserRequestInfo {
	userInfo, ok := ctx.Value(userContextKey).(UserRequestInfo)

	if !ok {
		return UserRequestInfo{
			IsLoggedIn: false,
		}
	}

	userId, _ := authctx.GetUserIDFromToken(ctx)
	email, _ := authctx.GetEmailFromToken(ctx)

	return UserRequestInfo{
		IsLoggedIn: userInfo.IsLoggedIn,
		Id:         userId,
		Email:      email,
	}
}
