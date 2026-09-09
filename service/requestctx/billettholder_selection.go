package requestctx

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const SelectedBillettholderCookieName = "selectedBillettholderId"

type selectedBillettholderIDKey struct{}

// SelectedBillettholderID returns a presentation hint from the selection cookie.
// Callers must validate that the ID belongs to an associated billettholder.
func SelectedBillettholderID(ctx context.Context) int {
	id, _ := ctx.Value(selectedBillettholderIDKey{}).(int)
	return id
}

func BillettholderSelectionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		selectedID := 0
		if cookie, err := r.Cookie(SelectedBillettholderCookieName); err == nil {
			if id, err := strconv.Atoi(cookie.Value); err == nil && id > 0 {
				selectedID = id
			}
		}
		ctx := context.WithValue(r.Context(), selectedBillettholderIDKey{}, selectedID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ClearBillettholderSelectionCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     SelectedBillettholderCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   requestIsSecure(r),
		SameSite: http.SameSiteLaxMode,
	})
}

// Match authctx's handling of TLS termination at the application's proxy.
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
