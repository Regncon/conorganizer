package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/keyvalue"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupEventRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, db *sql.DB, logger *slog.Logger, eventImageDir *string) error {
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

	resetMVC := func(mvc *root.TodoMVC) {
		mvc.Mode = root.TodoViewModeAll
		mvc.Todos = []*root.Todo{
			{Text: "Learn a backend language", Completed: true},
			{Text: "Learn Datastar", Completed: false},
			{Text: "Create Hypermedia", Completed: false},
			{Text: "???", Completed: false},
			{Text: "Profit", Completed: false},
		}
		mvc.EditingIdx = -1
	}

	mvcSession := func(w http.ResponseWriter, r *http.Request) (string, *root.TodoMVC, error) {
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
			resetMVC(mvc)

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

	//TODO FIX THIS SO WE SE THE ROUTER AND PAS IT IN (hard to find if we do this)
	eventLayoutRoute(router, db, logger, eventImageDir, err)

	router.Route("/event/api", func(eventApiRouter chi.Router) {
		eventApiRouter.Route("/{idx}", func(eventIdRouter chi.Router) {
			eventIdRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				eventID := chi.URLParam(r, "idx")
				sessionID, mvc, err := mvcSession(w, r)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				sse := datastar.NewSSE(w, r)

				// Watch for updates
				ctx := r.Context()
				watcher, err := kv.Watch(ctx, sessionID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer watcher.Stop()

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
						isAdmin := authctx.GetAdminFromUserToken(ctx)
						c := event_page(eventID, isAdmin, logger, db, eventImageDir, r)
						if err := sse.PatchElementTempl(c); err != nil {
							sse.ConsoleError(err)
							return
						}
					}
				}
			})

			eventIdRouter.Route("/interest", func(eventInterest chi.Router) {
				eventInterest.Route("/update", func(updateInterestRouter chi.Router) {

					updateInterestRouter.Put("/{interestLevel}", func(w http.ResponseWriter, r *http.Request) {
						type Put struct {
							BillettHolderId int    `json:"billettHolderId"`
							PuljeId         string `json:"puljeId"`
						}
						signals := &Put{}

						if readSignalErr := datastar.ReadSignals(r, signals); readSignalErr != nil {
							logger.Error("Failed to read signals", "error", readSignalErr)
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
							return
						}
						fmt.Printf("Signals: %+v\n", signals)
						ctx := r.Context()
						userInfo := userctx.GetUserRequestInfo(ctx)
						// billettholderId, err := strconv.Atoi(chi.URLParam(r, "billettholderId"))
						// if err != nil {
						// 	logger.Error("Failed to convert billettholderId to int", "error", err)
						// 	http.Error(w, "Failed to convert billettholderId to int", http.StatusBadRequest)
						// 	return
						// }

						eventId := chi.URLParam(r, "idx")
						// convert interestLevel string to InterestLevels struct
						var interestLevel InterestLevels

						switch chi.URLParam(r, "interestLevel") {
						case "high":
							interestLevel.High = "high"
						case "medium":
							interestLevel.Medium = "medium"
						case "low":
							interestLevel.Low = "low"
						case "none":
							interestLevel.None = "none"
						}
						sessionId, _, mvcErr := mvcSession(w, r)

						if mvcErr != nil {
							http.Error(w, mvcErr.Error(), http.StatusInternalServerError)
							return
						}

						updateInterest(userInfo.Id, signals.BillettHolderId, eventId, interestLevel, signals.PuljeId, db, logger)

						logger.Info(fmt.Sprintf("%d", signals.BillettHolderId), eventId, signals.PuljeId, sessionId, userInfo, fmt.Sprintf("%+v", interestLevel), "ASDFG")

						if err := keyvalue.BroadcastUpdate(kv, r); err != nil {
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

type InterestLevels struct {
	Low    string `json:"low"`
	Medium string `json:"medium"`
	High   string `json:"high"`
	None   string `json:"none"`
}

func ConvertInterestLevelToDbInterestLevel(interest InterestLevels) string {
	switch {
	case interest.High != "":
		return "Veldig interessert"
	case interest.Medium != "":
		return "Middels interessert"
	case interest.Low != "":
		return "Litt interessert"
	case interest.None != "":
		return "Ikke interessert"
	default:
		return "Ikke interessert"
	}
}

func updateInterest(
	userId string,
	billettholderId int,
	eventID string,
	interest InterestLevels,
	puljeId string,
	db *sql.DB,
	logger *slog.Logger,
) error {
	puljeQuery := `SELECT EXISTS (SELECT * FROM event_puljer WHERE event_id = $1 AND pulje_id = $2 AND is_active = 1 AND is_published = 1)`
	_, puljerErr := db.Query(puljeQuery, eventID, puljeId)
	if puljerErr != nil {
		logger.Error("failed to check if pulje exists", "error", puljerErr)
		return puljerErr
	}

	userHasAccessToBillettHolderIdQuery := `
        SELECT EXISTS
            (SELECT *
                FROM billettholdere_users [BU]
                JOIN users [U] ON [BU].user_id = [U].id
                WHERE [BU].billettholder_id = $1 AND [U].user_id = $2)`
	_, userHasAccessErr := db.Query(userHasAccessToBillettHolderIdQuery, billettholderId, userId)

	if userHasAccessErr != nil {
		logger.Error("failed to check if user has access to billettholder", "error", userHasAccessErr)
		return userHasAccessErr
	}

	if interest.None != "" {
		dropQuery := `DELETE FROM interests WHERE event_id = $1 AND pulje_id = $2 AND billettholder_id = $3`
		dropRows, dropErr := db.Exec(dropQuery, eventID, puljeId, billettholderId)
		if dropErr != nil {
			logger.Error("failed to drop interest", "error", dropErr)
			return dropErr
		}

		_, dropAffectedErr := dropRows.RowsAffected()
		if dropAffectedErr != nil {
			logger.Error("failed to get affected rows", "error", dropAffectedErr)
			return dropAffectedErr
		}

		logger.Info("Interest dropped successfully")

		return fmt.Errorf("interest dropped successfully")
	}

	updateQuery := `
                INSERT INTO interests (billettholder_id, event_id, pulje_id, interest_level)
                VALUES (?, ?, ?, ?)
                ON CONFLICT(billettholder_id, pulje_id, event_id) DO UPDATE SET
                    interest_level = excluded.interest_level
            `
	updateRows, updateAffectedErr := db.Exec(updateQuery, billettholderId, eventID, puljeId, ConvertInterestLevelToDbInterestLevel(interest))
	if updateAffectedErr != nil {
		logger.Error("failed to update interest", "error", updateAffectedErr)
		return updateAffectedErr
	}

	updateAffected, updateAffectedErr := updateRows.RowsAffected()
	if updateAffectedErr != nil {
		logger.Error("failed to get affected rows", "error", updateAffectedErr)
		return updateAffectedErr
	}

	if updateAffected == 0 {
		logger.Info("no rows were updated")
		return nil
	}

	logger.Info("Interest updated successfully")

	return nil
}
