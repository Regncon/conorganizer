package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go/jetstream"
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
		return "Du har ikkje tilgang til å endre interessa til denne billettheldaren. Kontakt styret."
	}
	if strings.Contains(err.Error(), "is not active and published for event") {
		return "Denne pulja er ikkje tilgjengeleg for dette arrangementet."
	}
	if strings.Contains(err.Error(), "is locked for event") {
		return "Pulja er låst. Du kan ikkje melde eller endre interesse lenger medan vi fordeler spelarar."
	}
	if strings.Contains(err.Error(), "is completed for event") {
		return "Puljefordelinga er klar. Gå til profilen din for å sjå kva du fekk."
	}
	if strings.Contains(err.Error(), "program is not published") {
		return "Interessevalget er ikke åpnet ennå."
	}
	return "Det oppstod ein feil då interessa skulle lagrast. Prøv igjen, eller kontakt styret dersom feilen held fram."
}

func SetupEventRoute(router chi.Router, ns *embeddednats.Server, liveManager *live.Manager, db *sql.DB, logger *slog.Logger, eventImageDir *string) error {
	logger = logger.With("component", "event")
	nc, err := ns.Client()
	if err != nil {
		return fmt.Errorf("error creating nats client: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("error creating jetstream client: %w", err)
	}

	if err := setupPuljeScheduledBroadcasts(context.Background(), js, liveManager, db, logger); err != nil {
		return fmt.Errorf("error setting up pulje scheduled broadcasts: %w", err)
	}

	//TODO FIX THIS SO WE SE THE ROUTER AND PAS IT IN (hard to find if we do this)
	eventLayoutRoute(router, db, logger, eventImageDir, err)

	router.Route("/event/api", func(eventApiRouter chi.Router) {
		eventApiRouter.Route("/{idx}", func(eventIdRouter chi.Router) {
			eventIdRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				eventID := chi.URLParam(r, "idx")
				liveManager.Stream(w, r, live.Page{
					Buckets: []live.Bucket{live.BucketEvents, live.BucketInterests},
					Render: func(ctx context.Context, r *http.Request) templ.Component {
						isAdmin := authctx.GetAdminFromUserToken(ctx)
						return event_page(eventID, isAdmin, logger, db, eventImageDir, r)
					},
				})
			})

			eventIdRouter.Route("/interest", func(eventInterest chi.Router) {

				eventInterest.Put("/selected-interest", func(w http.ResponseWriter, r *http.Request) {
					eventId := chi.URLParam(r, "idx")
					type Signals struct {
						BillettHolderId int    `json:"billettHolderId"`
						PuljeId         string `json:"puljeId"`
					}
					signals := &Signals{}
					if readSignalErr := datastar.ReadSignals(r, signals); readSignalErr != nil {
						logger.Error(fmt.Errorf("failed to read event interest signals: %w", readSignalErr).Error())
						http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
						return
					}

					interest, err := getSelectedInterest(eventId, signals.BillettHolderId, signals.PuljeId, db)
					if err != nil {
						logger.Error(fmt.Errorf("failed to get selected interest: %w", err).Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}

					sse := datastar.NewSSE(w, r)
					signalJSON := fmt.Appendf(nil, `{"selectedInterestLevel": %q, "currentInterestLevelChoice": "Pending choice"}`, interest)
					if err := sse.PatchSignals(signalJSON); err != nil {
						logger.Error(fmt.Errorf("failed to patch selected interest signal: %w", err).Error(), "event_id", eventId, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId, "selectedInterestLevel", interest)
					}

				})

				eventInterest.Route("/update", func(updateInterestRouter chi.Router) {

					updateInterestRouter.Put("/interest", func(w http.ResponseWriter, r *http.Request) {
						type Put struct {
							BillettHolderId            int                  `json:"billettHolderId"`
							PuljeId                    string               `json:"puljeId"`
							CurrentInterestLevelChoice models.InterestLevel `json:"currentInterestLevelChoice"`
						}
						signals := &Put{}

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
							if err := patchInterestErrorSignal(sse, "Vel billetthelder f\u00f8r du melder interesse."); err != nil {
								logger.Error(err.Error(), "event_id", eventId, "user_id", userInfo.Id, "pulje_id", signals.PuljeId, "billettholder_id", signals.BillettHolderId)
							}
							return
						}
						if signals.PuljeId == "" {
							logger.Error("Rejected interest update: missing pulje id", "event_id", eventId, "user_id", userInfo.Id, "billettholder_id", signals.BillettHolderId)
							if err := patchInterestErrorSignal(sse, "Vel pulje f\u00f8r du melder interesse."); err != nil {
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
					})

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

	programPublished, programPublishedErr := getProgramPublished(db)
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
			AND ep.is_published = 1
			AND e.status = $3
	`
	var puljeStatus models.PuljeStatus
	puljerErr := db.QueryRow(puljeQuery, eventID, puljeId, models.EventStatusAnnounced).Scan(&puljeStatus)
	if puljerErr != nil {
		if puljerErr == sql.ErrNoRows {
			return fmt.Errorf("pulje %s is not active and published for event %s", puljeId, eventID)
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
