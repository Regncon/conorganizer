package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Regncon/conorganizer/components/event_components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	eventservice "github.com/Regncon/conorganizer/service/eventService"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/program"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func patchInterestErrorSignal(sse *datastar.ServerSentEventGenerator, errorMessage string) error {
	signalJSON, err := json.Marshal(map[string]string{
		"interestErrorMessage": errorMessage,
	})
	if err != nil {
		return fmt.Errorf("marshal interest error signal: %w", err)
	}
	if err := sse.PatchSignals(signalJSON); err != nil {
		return fmt.Errorf("patch interest error signal: %w", err)
	}
	return nil
}

func interestErrorMessageFromError(err error) string {
	if err == nil {
		return ""
	}
	if strings.Contains(err.Error(), "does not have access") {
		return "Du har ikke tilgang til å endre interessen til denne billettholderen. Kontakt styret."
	}
	if strings.Contains(err.Error(), "is not active for event") {
		return "Denne puljen er ikke tilgjengelig for dette arrangementet."
	}
	if strings.Contains(err.Error(), "is locked for event") {
		return "Puljen er låst. Du kan ikke melde eller endre interesse lenger mens vi fordeler spillere."
	}
	if strings.Contains(err.Error(), "is completed for event") {
		return "Puljefordelingen er klar. Gå til profilen din for å se hva du fikk."
	}
	if strings.Contains(err.Error(), "program is not published") {
		return "Interessevalget er ikke åpnet ennå."
	}
	return "Det oppstod en feil da interessen skulle lagres. Prøv igjen, eller kontakt styret dersom feilen fortsetter."
}

func SetupEventRoute(router chi.Router, liveManager *live.Manager, db *sql.DB, logger *slog.Logger, eventImageDir *string) error {
	logger = logger.With("component", "event")

	//TODO FIX THIS SO WE SE THE ROUTER AND PAS IT IN (hard to find if we do this)
	eventLayoutRoute(router, db, logger, eventImageDir, nil)

	router.Route("/event/api", func(eventApiRouter chi.Router) {
		eventApiRouter.Route("/{idx}", func(eventIdRouter chi.Router) {
			eventIdRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				eventID := chi.URLParam(r, "idx")
				liveManager.Stream(w, r, live.Page{
					Buckets: []live.Bucket{live.BucketEvents, live.BucketInterests, live.BucketRooms},
					Render: func(ctx context.Context, r *http.Request) templ.Component {
						isAdmin := authctx.GetAdminFromUserToken(ctx)
						return event_page(eventID, isAdmin, logger, db, eventImageDir, r)
					},
				})
			})

			eventIdRouter.Route("/interest", func(eventInterest chi.Router) {
				eventInterest.Put("/selected-interest", selectedInterestHandler(db, logger, eventImageDir))
				eventInterest.Route("/update", func(updateInterestRouter chi.Router) {
					updateInterestRouter.Put("/interest", interestUpdateHandler(liveManager, db, logger))
				})
			})
		})
	})

	return nil
}

func hasValidInterestChoice(interest models.InterestLevel) bool {
	return interest.Valid()
}

func getSelectedInterest(eventId string, billettholderId int, puljeId string, db *sql.DB) (models.InterestLevel, error) {
	query := `SELECT interest_level FROM interests WHERE event_id = $1 AND billettholder_id = $2 AND pulje_id = $3`
	var interestLevel models.InterestLevel
	err := db.QueryRow(query, eventId, billettholderId, puljeId).Scan(&interestLevel)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.InterestLevelNone, nil
		}
		return models.InterestLevelNone, fmt.Errorf("failed to get selected interest: %w", err)
	}
	return interestLevel, nil
}

func updateInterest(
	userId string,
	billettholderId int,
	eventID string,
	currentInterestLevelChoice models.InterestLevel,
	puljeId string,
	db *sql.DB,
) error {

	if eventID == "" {
		return fmt.Errorf("event id is required")
	}
	if billettholderId <= 0 {
		return fmt.Errorf("billettholder id is required")
	}
	if puljeId == "" {
		return fmt.Errorf("pulje id is required")
	}
	if !hasValidInterestChoice(currentInterestLevelChoice) {
		return fmt.Errorf("interest level is required")
	}

	programPublished, programPublishedErr := program.IsPublished(db)
	if programPublishedErr != nil {
		return fmt.Errorf("failed to check program publishing state: %w", programPublishedErr)
	}
	if !programPublished {
		return fmt.Errorf("program is not published")
	}

	puljeQuery := `
		SELECT p.status
		FROM relation_event_puljer ep
		JOIN puljer p ON p.id = ep.pulje_id
		JOIN events e ON e.id = ep.event_id
		WHERE ep.event_id = $1
			AND ep.pulje_id = $2
			AND ep.is_in_pulje = 1
			AND e.is_in_puljefordeling = 1
			AND e.status = $3
	`
	var puljeStatus models.PuljeStatus
	puljerErr := db.QueryRow(puljeQuery, eventID, puljeId, models.EventStatusAnnounced).Scan(&puljeStatus)
	if puljerErr != nil {
		if puljerErr == sql.ErrNoRows {
			return fmt.Errorf("pulje %s is not active for event %s", puljeId, eventID)
		}
		return fmt.Errorf("failed to check if pulje %s exists for event %s: %w", puljeId, eventID, puljerErr)
	}
	if puljeStatus == models.PuljeStatusLocked {
		return fmt.Errorf("pulje %s is locked for event %s", puljeId, eventID)
	}
	if puljeStatus == models.PuljeStatusCompleted {
		return fmt.Errorf("pulje %s is completed for event %s", puljeId, eventID)
	}

	userHasAccessToBillettHolderIdQuery := `
        SELECT EXISTS
            (SELECT 1
                FROM relation_billettholdere_users [BU]
                JOIN users [U] ON [BU].user_id = [U].id
                WHERE [BU].billettholder_id = $1 AND [U].external_id = $2)`
	var userHasAccess bool
	userHasAccessErr := db.QueryRow(userHasAccessToBillettHolderIdQuery, billettholderId, userId).Scan(&userHasAccess)

	if userHasAccessErr != nil {
		return fmt.Errorf("failed to check if user %s has access to billettholder %d: %w", userId, billettholderId, userHasAccessErr)
	}
	if !userHasAccess {
		return fmt.Errorf("user %s does not have access to this billettholder interest", userId)
	}

	if currentInterestLevelChoice == models.InterestLevelNone {
		dropQuery := `DELETE FROM interests WHERE event_id = $1 AND pulje_id = $2 AND billettholder_id = $3`
		dropRows, dropErr := db.Exec(dropQuery, eventID, puljeId, billettholderId)
		if dropErr != nil {
			return fmt.Errorf("failed to drop interest for event %s, pulje %s, billettholder %d: %w", eventID, puljeId, billettholderId, dropErr)
		}
		_, dropAffectedErr := dropRows.RowsAffected()
		if dropAffectedErr != nil {
			return fmt.Errorf("failed to get affected rows when dropping interest for event %s, pulje %s, billettholder %d: %w", eventID, puljeId, billettholderId, dropAffectedErr)
		}

		return nil
	}

	updateQuery := `
                INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
                VALUES (?, ?, ?, ?)
                ON CONFLICT(billettholder_id, pulje_id, event_id) DO UPDATE SET
                    interest_level = excluded.interest_level
            `
	updateRows, updateErr := db.Exec(updateQuery, billettholderId, eventID, puljeId, currentInterestLevelChoice)
	if updateErr != nil {
		return fmt.Errorf("failed to update interest for event %s, pulje %s, billettholder %d: %w", eventID, puljeId, billettholderId, updateErr)
	}

	updateAffected, updateAffectedErr := updateRows.RowsAffected()
	if updateAffectedErr != nil {
		return fmt.Errorf("failed to get affected rows when updating interest for event %s, pulje %s, billettholder %d: %w", eventID, puljeId, billettholderId, updateAffectedErr)
	}

	if updateAffected == 0 {
		return nil
	}

	return nil
}

type interestSignals struct {
	BillettHolderId int    `json:"billettHolderId"`
	PuljeId         string `json:"puljeId"`
}

type interestUpdateSignals struct {
	interestSignals
	CurrentInterestLevelChoice models.InterestLevel `json:"currentInterestLevelChoice"`
}

func interestUpdateHandler(liveManager *live.Manager, db *sql.DB, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		signals := &interestUpdateSignals{}

		if readSignalErr := datastar.ReadSignals(r, signals); readSignalErr != nil {
			logger.Error(fmt.Errorf("failed to read event interest signals: %w", readSignalErr).Error())
			http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		userInfo := userctx.GetUserRequestInfo(ctx)
		sse := datastar.NewSSE(w, r)

		eventId := chi.URLParam(r, "idx")
		if eventId == "" {
			logger.Error("Rejected interest update: missing event id", "user_id", userInfo.Id, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
			if err := patchInterestErrorSignal(sse, "Mangler arrangement."); err != nil {
				logger.Error(err.Error(), "user_id", userInfo.Id, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
			}
			return
		}
		if signals.BillettHolderId <= 0 {
			logger.Error("Rejected interest update: missing billettholder id", "event_id", eventId, "user_id", userInfo.Id, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
			if err := patchInterestErrorSignal(sse, "Velg billettholder f\u00f8r du melder interesse."); err != nil {
				logger.Error(err.Error(), "event_id", eventId, "user_id", userInfo.Id, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
			}
			return
		}
		if signals.PuljeId == "" {
			logger.Error("Rejected interest update: missing pulje id", "event_id", eventId, "user_id", userInfo.Id, "billettholder_id", signals.BillettHolderId)
			if err := patchInterestErrorSignal(sse, "Velg pulje f\u00f8r du melder interesse."); err != nil {
				logger.Error(err.Error(), "event_id", eventId, "user_id", userInfo.Id, "billettholder_id", signals.BillettHolderId)
			}
			return
		}

		if err := updateInterest(userInfo.Id, signals.BillettHolderId, eventId, signals.CurrentInterestLevelChoice, signals.PuljeId, db); err != nil {
			logger.Error(
				err.Error(),
				"event_id", eventId,
				"user_id", userInfo.Id,
				"pulje_id", signals.PuljeId,
				"billettholder_id", signals.BillettHolderId,
			)
			if patchErr := patchInterestErrorSignal(sse, interestErrorMessageFromError(err)); patchErr != nil {
				logger.Error(patchErr.Error(), "event_id", eventId, "user_id", userInfo.Id, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
			}
			return
		}

		if err := patchInterestErrorSignal(sse, ""); err != nil {
			logger.Error(err.Error(), "event_id", eventId, "user_id", userInfo.Id, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
		}

		logger.Debug("Interest update request handled",
			"event_id", eventId,
			"pulje_id", signals.PuljeId,
			"user_id", userInfo.Id,
			"billettholder_id", signals.BillettHolderId,
		)

		if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
			logger.Error(fmt.Errorf("failed to broadcast interest update: %w", err).Error(), "event_id", eventId, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
			http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
			return
		}
	}
}

func selectedInterestHandler(db *sql.DB, logger *slog.Logger, eventImageDir *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		eventID := chi.URLParam(r, "idx")
		var signals interestSignals
		if err := datastar.ReadSignals(r, &signals); err != nil {
			http.Error(w, "Ugyldig billettholdervalg.", http.StatusBadRequest)
			return
		}
		if signals.BillettHolderId <= 0 || signals.PuljeId == "" {
			http.Error(w, "Velg billettholder og pulje.", http.StatusBadRequest)
			return
		}

		userInfo := userctx.GetUserRequestInfo(r.Context())
		logger := logger.With("event_id", eventID, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
		// Billettholdere are associated by email, which is also what GetTicketHolders
		// uses to build the picker. Checking the same relation here keeps the endpoint
		// from rejecting a billettholder the picker just offered.
		var hasAccess bool
		err := db.QueryRowContext(r.Context(), `
			SELECT EXISTS (
				SELECT 1 FROM relation_billettholder_emails
				WHERE billettholder_id = ?1 AND email = ?2
			)
		`, signals.BillettHolderId, userInfo.Email).Scan(&hasAccess)
		if err != nil {
			logger.Error(fmt.Errorf("failed to check billettholder access: %w", err).Error())
			http.Error(w, "Kunne ikke hente billettholder.", http.StatusInternalServerError)
			return
		}
		if !hasAccess {
			http.Error(w, "Du har ikke tilgang til denne billettholderen.", http.StatusForbidden)
			return
		}

		event, err := eventservice.GetEventById(eventID, db)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Kunne ikke hente arrangementet.", http.StatusInternalServerError)
			return
		}
		viewDecision, err := decideEventView(event, userInfo, db)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Kunne ikke sjekke tilgang til arrangementet.", http.StatusInternalServerError)
			return
		}
		if !viewDecision.CanView {
			http.Error(w, "Arrangementet er ikke tilgjengelig.", viewDecision.HiddenResponseStatusCode)
			return
		}

		notices, err := event_components.LoadInterestNoticeState(signals.BillettHolderId, signals.PuljeId, event.AgeGroup, eventImageDir, db)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Kunne ikke hente interessevarsler.", http.StatusInternalServerError)
			return
		}
		interest, err := getSelectedInterest(eventID, signals.BillettHolderId, signals.PuljeId, db)
		if err != nil {
			logger.Error(err.Error())
			http.Error(w, "Kunne ikke hente valgt interesse.", http.StatusInternalServerError)
			return
		}
		sse := datastar.NewSSE(w, r)
		if err := sse.PatchElementTempl(
			event_components.InterestContent(eventID, event.Title, notices, interest),
			datastar.WithModeReplace(),
		); err != nil {
			logger.Error(fmt.Errorf("failed to patch selected interest content: %w", err).Error())
		}
	}
}
