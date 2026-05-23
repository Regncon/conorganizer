package profilepage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	profilecomponent "github.com/Regncon/conorganizer/components/profile"
	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/models"
	profileticketspage "github.com/Regncon/conorganizer/pages/profile/tickets"
	"github.com/Regncon/conorganizer/pages/root"
	billettholderService "github.com/Regncon/conorganizer/service/billettholder"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupProfileRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, db *sql.DB, eventImageDir *string, logger *slog.Logger) error {
	nc, err := ns.Client()
	if err != nil {
		return fmt.Errorf("error creating nats client: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("error creating jetstream client: %w", err)
	}

	kv, err := js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "events",
		Description: "Regncon Event Store",
		Compression: true,
		TTL:         time.Hour,
		MaxBytes:    16 * 1024 * 1024,
	})
	if err != nil {
		return fmt.Errorf("error creating key value: %w", err)
	}

	profileSession := func(w http.ResponseWriter, r *http.Request) (string, *root.TodoMVC, error) {
		ctx := r.Context()
		sessionID, err := upsertSessionID(store, r, w)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get session id: %w", err)
		}

		mvc := &root.TodoMVC{}
		if entry, err := kv.Get(ctx, sessionID); err != nil {
			if err != jetstream.ErrKeyNotFound {
				return "", nil, fmt.Errorf("failed to get key value: %w", err)
			}

			if err := saveMVC(ctx, mvc, sessionID, kv); err != nil {
				return "", nil, fmt.Errorf("failed to save mvc: %w", err)
			}
		} else {
			if err := json.Unmarshal(entry.Value(), mvc); err != nil {
				return "", nil, fmt.Errorf("failed to unmarshal mvc: %w", err)
			}
		}
		return sessionID, mvc, nil
	}

	var profileTicketsErr error
	router.Route("/profile", func(profileRouter chi.Router) {
		profileRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			requestLogger := logger.With("component", "profile")
			ctx := r.Context()
			user := userctx.GetUserRequestInfo(ctx)
			events := GetEventsByExternalID(user.Id, db, requestLogger)
			billettholdere, err := billettholderService.GetBillettholdere(user.Id, db)
			if err != nil {
				requestLogger.Error(err.Error(), "user_id", user.Id)
			}

			selectedBillettholderID := selectedBillettholderIDFromRequest(r, user, billettholdere, requestLogger)
			validBillettholderIDs := billettholderIDs(billettholdere)

			tickets := make([]profilecomponent.TicketHolder, 0, len(billettholdere))
			for _, billettholder := range billettholdere {
				email := ""
				if len(billettholder.Emails) > 0 {
					email = billettholder.Emails[0].Email
				}
				tickets = append(tickets, profilecomponent.TicketHolder{
					Name:   strings.TrimSpace(billettholder.FirstName + " " + billettholder.LastName),
					Ticket: billettholder.TicketType,
					Email:  email,
				})
			}

			if err := layouts.Base(
				"Min profil side",
				user,
				ProfilePage(user, events, tickets, selectedBillettholderID, validBillettholderIDs, db, requestLogger, eventImageDir),
			).Render(ctx, w); err != nil {
				requestLogger.Error(fmt.Errorf("failed to render profile page: %w", err).Error(), "user_id", user.Id)
			}
		})

		profileRouter.Route("/api", func(profileApiRouter chi.Router) {
			profileApiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				requestLogger := logger.With("component", "profile")
				sessionID, mvc, err := profileSession(w, r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				sse := datastar.NewSSE(w, r)
				ctx := r.Context()
				user := userctx.GetUserRequestInfo(ctx)
				watcher, err := kv.Watch(ctx, sessionID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer func() {
					if err := watcher.Stop(); err != nil {
						requestLogger.Error(fmt.Errorf("failed to stop profile watcher: %w", err).Error())
					}
				}()

				renderProfileMainColumn := func() error {
					events := GetEventsByExternalID(user.Id, db, requestLogger)
					billettholdere, err := billettholderService.GetBillettholdere(user.Id, db)
					if err != nil {
						requestLogger.Error(err.Error(), "user_id", user.Id)
					}
					selectedBillettholderID := selectedBillettholderIDFromRequest(r, user, billettholdere, requestLogger)
					return sse.PatchElementTempl(ProfileMainColumn(user, events, selectedBillettholderID, db, requestLogger, eventImageDir))
				}

				if err := renderProfileMainColumn(); err != nil {
					_ = sse.ConsoleError(err)
					return
				}

				for {
					select {
					case <-ctx.Done():
						return
					case entry := <-watcher.Updates():
						if entry == nil {
							continue
						}
						if err := json.Unmarshal(entry.Value(), mvc); err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
							return
						}
						if err := renderProfileMainColumn(); err != nil {
							_ = sse.ConsoleError(err)
							return
						}
					}
				}
			})
		})

		profileTicketsErr = profileticketspage.ProfileTicketsRoute(profileRouter, store, ns, db, logger)
	})

	if profileTicketsErr != nil {
		return fmt.Errorf("error setting up profile tickets route: %w", profileTicketsErr)
	}

	return nil
}

func GetEventsByExternalID(externalID string, db *sql.DB, logger *slog.Logger) []models.EventCardModel {
	logger = logger.With("component", "profile")
	var events []models.EventCardModel

	// Get events where event created id is the same as user
	userID, err := userctx.GetUserIDFromExternalID(externalID, db, logger)
	if err != nil {
		logger.Error(fmt.Errorf("failed to resolve user_id for external_id %q: %w", externalID, err).Error())
		return events
	}

	// Query for events created by user
	eventsQuery := "SELECT id, title, intro, status, system, host_name, beginner_friendly, event_type, age_group, event_runtime, can_be_run_in_english FROM events WHERE user_id = ?"
	rows, eventsQueryErr := db.Query(eventsQuery, userID)
	if eventsQueryErr != nil {
		logger.Error(fmt.Errorf("failed to query events for external_id %q: %w", externalID, eventsQueryErr).Error())
		return events
	}
	defer rows.Close()

	// Validate database query return
	for rows.Next() {
		var event models.EventCardModel
		if scanErr := rows.Scan(&event.Id, &event.Title, &event.Intro, &event.Status, &event.System, &event.HostName, &event.BeginnerFriendly, &event.EventType, &event.AgeGroup, &event.Runtime, &event.CanBeRunInEnglish); scanErr != nil {
			logger.Error(fmt.Errorf("failed to scan event row for external_id %q: %w", externalID, scanErr).Error())
			return events
		}
		events = append(events, event)
	}

	return events
}

func selectedBillettholderIDFromRequest(r *http.Request, user requestctx.UserRequestInfo, billettholdere []models.Billettholder, logger *slog.Logger) int {
	rawID := strings.TrimSpace(r.URL.Query().Get("b_id"))
	if rawID == "" {
		return defaultSelectedBillettholderID(user, billettholdere, logger)
	}

	billettholderID, err := strconv.Atoi(rawID)
	if err != nil {
		logger.Debug("Ignoring invalid profile billettholder id", "user_id", user.Id, "b_id", rawID)
		return defaultSelectedBillettholderID(user, billettholdere, logger)
	}

	if hasBillettholderID(billettholdere, billettholderID) {
		return billettholderID
	}

	logger.Debug("Ignoring profile billettholder id without user relation", "user_id", user.Id, "b_id", billettholderID)
	return defaultSelectedBillettholderID(user, billettholdere, logger)
}

func hasBillettholderID(billettholdere []models.Billettholder, billettholderID int) bool {
	for _, billettholder := range billettholdere {
		if billettholder.ID == billettholderID {
			return true
		}
	}
	return false
}

func billettholderIDs(billettholdere []models.Billettholder) []int {
	ids := make([]int, 0, len(billettholdere))
	for _, billettholder := range billettholdere {
		ids = append(ids, billettholder.ID)
	}
	return ids
}

func defaultSelectedBillettholderID(user requestctx.UserRequestInfo, billettholdere []models.Billettholder, logger *slog.Logger) int {
	for _, billettholder := range billettholdere {
		for _, email := range billettholder.Emails {
			if strings.EqualFold(email.Email, user.Email) {
				return billettholder.ID
			}
		}
	}

	if len(billettholdere) == 0 {
		logger.Debug("No billettholder available for profile selection", "user_id", user.Id)
		return 0
	}
	return billettholdere[0].ID
}

func saveMVC(ctx context.Context, mvc *root.TodoMVC, sessionID string, kv jetstream.KeyValue) error {
	b, err := json.Marshal(mvc)
	if err != nil {
		return fmt.Errorf("failed to marshal mvc: %w", err)
	}
	if _, err := kv.Put(ctx, sessionID, b); err != nil {
		return fmt.Errorf("failed to put key value: %w", err)
	}
	return nil
}

func upsertSessionID(store sessions.Store, r *http.Request, w http.ResponseWriter) (string, error) {
	sess, err := store.Get(r, "connections")
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}
	id, ok := sess.Values["id"].(string)
	if !ok {
		id = toolbelt.NextEncodedID()
		sess.Values["id"] = id
		if err := sess.Save(r, w); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}
	}
	return id, nil
}
