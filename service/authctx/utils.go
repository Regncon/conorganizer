package authctx

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/descope/go-sdk/descope"
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

func GetAdminFromUserToken(ctx context.Context) bool {
	userToken, err := GetUserTokenFromContext(ctx)
	if err != nil {
		return false
	}

	rolesClaim, ok := userToken.Claims["roles"]
	if !ok {
		return false
	}

	roles, ok := rolesClaim.([]any)
	if !ok {
		return false
	}

	for _, role := range roles {
		if roleStr, ok := role.(string); ok && roleStr == "Admin" {
			return true
		}
	}

	return false
}

func RequireAdmin(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			isAdmin := GetAdminFromUserToken(r.Context())
			if !isAdmin {
				logger.Warn("User is not an admin")
				http.Error(w, "You are not an admin", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUserIDFromToken(ctx context.Context) (string, error) {
	userToken, err := GetUserTokenFromContext(ctx)
	if err != nil {
		return "", err
	}

	return userToken.ID, nil
}
