package userctx

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/service/authctx"
)

type userctxKey struct{}

var userContextKey = userctxKey{}

type UserRequestInfo struct {
	IsLoggedIn bool
}

func IsLoggedInMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isLoggedIn := false

			sessionCookie, sessionCookieError := r.Cookie(authctx.SessionCookieName)
			if sessionCookieError == nil {
				descopeClient, descopeClientError := authctx.GetDescopeClient()
				if descopeClientError == nil {
					_, userToken, validateTokenError := descopeClient.Auth.ValidateSessionWithToken(r.Context(), sessionCookie.Value)
					if validateTokenError == nil && userToken != nil {
						isLoggedIn = true
					}
				}
			}

			ctx := context.WithValue(r.Context(), userContextKey, isLoggedIn)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserRequestInfo(ctx context.Context) UserRequestInfo {
	if ctx == nil {
		return UserRequestInfo{IsLoggedIn: false}
	}
	isLoggedIn, ok := ctx.Value(userContextKey).(bool)
	if !ok {
		return UserRequestInfo{IsLoggedIn: false}
	}
	return UserRequestInfo{
		IsLoggedIn: isLoggedIn,
	}
}
