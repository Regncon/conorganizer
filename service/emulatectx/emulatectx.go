// Package emulatectx lets an admin browse the site read-only as a specific
// non-admin player ("view as"). A middleware runs after authctx.AuthMiddleware:
// when the real user is an admin and an emulate_target cookie is present, it
// replaces the request identity with a synthesized token for the target player
// (no Admin role) and stashes the real admin token so the banner and exit flow
// keep working. Navigating to any /admin route auto-exits emulation.
package emulatectx

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/descope/go-sdk/descope"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/Regncon/conorganizer/service/authctx"
)

// CookieName is the cookie holding the emulated target's users.id.
const CookieName = "emulate_target"

type realAdminTokenKey struct{}

// EmulateMiddleware swaps the request identity to the cookie's target player when
// the real user is an admin. It must be mounted after authctx.AuthMiddleware.
func EmulateMiddleware(db *sql.DB, logger *slog.Logger) func(http.Handler) http.Handler {
	logger = logger.With("component", "emulate")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())

			cookie, cookieErr := r.Cookie(CookieName)
			if cookieErr != nil || cookie.Value == "" {
				next.ServeHTTP(w, r)
				return
			}

			// Security gate: only a real admin can ever emulate. A non-admin
			// presenting the cookie is ignored and the cookie cleared.
			if !authctx.GetAdminFromUserToken(r.Context()) {
				clearCookie(w, r)
				next.ServeHTTP(w, r)
				return
			}

			// Auto-exit: any /admin navigation ends emulation and restores the
			// admin's own view.
			if isAdminPath(r.URL.Path) {
				clearCookie(w, r)
				next.ServeHTTP(w, r)
				return
			}

			// Read-only: block writes while emulating.
			if !isSafeMethod(r.Method) {
				http.Error(w, "Du ser siden som en spiller (skrivebeskyttet). Avslutt emuleringen for å gjøre endringer.", http.StatusForbidden)
				return
			}

			external, email, loadErr := loadTargetUser(db, cookie.Value)
			if loadErr != nil {
				logger.Warn(fmt.Errorf("emulation target lookup failed: %w", loadErr).Error(), "request_id", requestID, "target", cookie.Value)
				clearCookie(w, r)
				next.ServeHTTP(w, r)
				return
			}

			realToken, _ := authctx.GetUserTokenFromContext(r.Context())
			ctx := context.WithValue(r.Context(), realAdminTokenKey{}, realToken)
			ctx = authctx.WithUserToken(ctx, syntheticToken(external, email))
			logger.Debug("emulating player", "request_id", requestID, "target", cookie.Value)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RealAdminTokenFromContext returns the stashed real admin token when the request
// is an admin emulating a player.
func RealAdminTokenFromContext(ctx context.Context) (*descope.Token, bool) {
	token, ok := ctx.Value(realAdminTokenKey{}).(*descope.Token)
	if !ok || token == nil {
		return nil, false
	}
	return token, true
}

// IsEmulating reports whether the current request is an admin emulating a player.
func IsEmulating(ctx context.Context) bool {
	_, ok := RealAdminTokenFromContext(ctx)
	return ok
}

func syntheticToken(externalID, email string) *descope.Token {
	return &descope.Token{
		ID: externalID,
		Claims: map[string]any{
			"email": email,
			"roles": []any{},
		},
	}
}

func loadTargetUser(db *sql.DB, userID string) (externalID, email string, err error) {
	const query = `SELECT external_id, email FROM users WHERE id = ?`
	row := db.QueryRow(query, userID)
	if scanErr := row.Scan(&externalID, &email); scanErr != nil {
		return "", "", fmt.Errorf("load target user %q: %w", userID, scanErr)
	}
	return externalID, email, nil
}

func isAdminPath(path string) bool {
	return path == "/admin" || strings.HasPrefix(path, "/admin/")
}

func isSafeMethod(method string) bool {
	switch method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

func clearCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
}
