package newevent

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const eventNotFoundMessage = "Fant ikke arrangementet."

func RequireEventOwnerOrAdmin(db *sql.DB, logger *slog.Logger) func(http.Handler) http.Handler {
	logger = logger.With("component", "new_event")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			userInfo := userctx.GetUserRequestInfo(ctx)
			if userInfo.IsAdmin {
				next.ServeHTTP(w, r)
				return
			}

			eventID := chi.URLParam(r, "id")
			isOwner, err := eventIsOwnedByUser(db, eventID, userInfo.Id)
			if err != nil {
				logger.Error(err.Error(), "request_id", middleware.GetReqID(ctx), "event_id", eventID, "user_id", userInfo.Id)
				http.Error(w, "Klarte ikke å sjekke tilgangen til arrangementet", http.StatusInternalServerError)
				return
			}
			if !isOwner {
				logger.Warn("User is neither owner nor admin of event", "request_id", middleware.GetReqID(ctx), "path", r.URL.Path, "event_id", eventID, "user_id", userInfo.Id)
				http.Error(w, eventNotFoundMessage, http.StatusNotFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func eventIsOwnedByUser(db *sql.DB, eventID string, externalUserID string) (bool, error) {
	if eventID == "" || externalUserID == "" {
		return false, nil
	}

	var isOwner bool
	if err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM events
			JOIN users ON users.id = events.user_id
			WHERE events.id = ? AND users.external_id = ?
		)
	`, eventID, externalUserID).Scan(&isOwner); err != nil {
		return false, fmt.Errorf("check owner of event %q: %w", eventID, err)
	}

	return isOwner, nil
}
