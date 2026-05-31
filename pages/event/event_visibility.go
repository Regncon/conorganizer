package event

import (
	"database/sql"
	"fmt"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
)

const eventNotAnnouncedMessage = "Dette arrangementet er ikke annonsert ennå. Kom tilbake senere, så får du se hva som venter."

func canViewEvent(event *models.Event, userInfo requestctx.UserRequestInfo, db *sql.DB) (bool, error) {
	if event == nil {
		return false, nil
	}
	if event.Status == models.EventStatusAnnounced {
		return true, nil
	}
	if userInfo.IsAdmin {
		return true, nil
	}
	if !userInfo.IsLoggedIn || userInfo.Id == "" || !event.UserID.Valid {
		return false, nil
	}

	var isOwner bool
	if err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE id = ? AND external_id = ?
		)
	`, event.UserID.Int64, userInfo.Id).Scan(&isOwner); err != nil {
		return false, fmt.Errorf("check event owner visibility for event %q: %w", event.ID, err)
	}

	return isOwner, nil
}
