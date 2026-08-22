package event

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
)

const eventNotAnnouncedMessage = "Dette arrangementet er ikke annonsert ennå. Kom tilbake senere, så får du se hva som venter."
const eventArchivedMessage = "Dette arrangementet er ikke tilgjengelig lenger."

type eventHiddenReason string

const (
	eventHiddenReasonUnannounced  eventHiddenReason = "unannounced"
	eventHiddenReasonArchivedGone eventHiddenReason = "archived_gone"
)

type eventViewDecision struct {
	CanView                  bool
	ShowUnannouncedWarning   bool
	ShowArchivedWarning      bool
	HiddenReason             eventHiddenReason
	HiddenResponseStatusCode int
}

func decideEventView(event *models.Event, userInfo requestctx.UserRequestInfo, db *sql.DB) (eventViewDecision, error) {
	if event == nil {
		return eventViewDecision{
			CanView:                  false,
			HiddenReason:             eventHiddenReasonUnannounced,
			HiddenResponseStatusCode: http.StatusNotFound,
		}, nil
	}
	if event.Status == models.EventStatusAnnounced {
		return eventViewDecision{CanView: true}, nil
	}
	if userInfo.IsAdmin {
		return eventPrivilegedViewDecision(event), nil
	}

	isOwner, err := eventOwnerMatchesUser(event, userInfo, db)
	if err != nil {
		return eventViewDecision{}, err
	}
	if isOwner {
		return eventPrivilegedViewDecision(event), nil
	}

	return eventHiddenViewDecision(event), nil
}

func eventPrivilegedViewDecision(event *models.Event) eventViewDecision {
	if event.Status == models.EventStatusArchived {
		return eventViewDecision{
			CanView:             true,
			ShowArchivedWarning: true,
		}
	}

	return eventViewDecision{
		CanView:                true,
		ShowUnannouncedWarning: true,
	}
}

func eventHiddenViewDecision(event *models.Event) eventViewDecision {
	if event.Status == models.EventStatusArchived {
		return eventViewDecision{
			CanView:                  false,
			HiddenReason:             eventHiddenReasonArchivedGone,
			HiddenResponseStatusCode: http.StatusGone,
		}
	}

	return eventViewDecision{
		CanView:                  false,
		HiddenReason:             eventHiddenReasonUnannounced,
		HiddenResponseStatusCode: http.StatusOK,
	}
}

func eventOwnerMatchesUser(event *models.Event, userInfo requestctx.UserRequestInfo, db *sql.DB) (bool, error) {
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
