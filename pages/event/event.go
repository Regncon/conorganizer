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
	billettholderService "github.com/Regncon/conorganizer/service/billettholder"
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
						c := event_page(eventID, isAdmin, logger, db, eventImageDir)
						if err := sse.PatchElementTempl(c); err != nil {
							sse.ConsoleError(err)
							return
						}
					}
				}
			})
		})

		eventApiRouter.Route("/{id}", func(eventIdRouter chi.Router) {
			eventIdRouter.Route("/new", func(eventNew chi.Router) {
				eventNew.Route("/interest", func(eventInterest chi.Router) {
					eventInterest.Put("/update", func(w http.ResponseWriter, r *http.Request) {

						type Put struct {
							InterestLevel string `json:"interest_level"`
							Pulje         string `json:"pulje"`
						}
						store := &Put{}

						if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
							logger.Error("Failed to read signals", "error", readSignalErr)
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
							return
						}
						ctx := r.Context()
						userInfo := userctx.GetUserRequestInfo(ctx)
						billettholderId, billettholderIdErr := billettholderService.GetBillettholderByUserId(db, logger, userInfo.Id)

						if billettholderIdErr != nil {
							logger.Error("Failed to get billettholder ID", "error", billettholderIdErr)
							http.Error(w, "Failed to get billettholder ID", http.StatusInternalServerError)
							return
						}

						eventID := chi.URLParam(r, "idx")
						value := r.URL.Query().Get("pulje")
						sessionID, mvc, mvcErr := mvcSession(w, r)
						if mvcErr != nil {
							http.Error(w, mvcErr.Error(), http.StatusInternalServerError)
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
	Low    string `json:"litt_interessert"`
	Medium string `json:"middels_interessert"`
	High   string `json:"veldig_interessert"`
}

func updateInterest(
	db *sql.DB,
	logger *slog.Logger,
	billettholder_id string,
	eventID string,
	interest InterestLevels,
	pulje string,
) error {
	puljeQuery := `SELECT EXISTS (SELECT * FROM event_puljer WHERE event_id = $1 AND pulje_id = $2)`
	_, puljerErr := db.Query(puljeQuery, eventID, pulje)
	if puljerErr != nil {
		logger.Info("failed to check if pulje exists", "error", puljerErr)
		return puljerErr
	}

	logger.Info(
		"updating interest",
		"eventID", eventID,
		"interest", interest,
		"pulje", pulje,
		"billettholder_id", billettholder_id,
	)
	updateQuery := `
                IF EXISTS (SELECT * FROM interests WHERE event_id = $1 AND pulje = $2)
                BEGIN
                    UPDATE interests
                    SET billettholder_id = $3, event_id = $1, interest_level = $4
                    WHERE event_id = $1 AND pulje = $2 AND billettholder_id = $3
                END
                ELSE
                BEGIN
                    INSERT INTO interests (billettholder_id, event_id, interest_level)
                    VALUES ($3, $1, $4)
                END
            `
	updateRows, updateAffectedErr := db.Exec(updateQuery, eventID, pulje, billettholder_id, interest)
	if updateAffectedErr != nil {
		logger.Info("failed to update interest", "error", updateAffectedErr)
		return updateAffectedErr
	}

	updateAffected, updateAffectedErr := updateRows.RowsAffected()
	if updateAffectedErr != nil {
		logger.Info("failed to get affected rows", "error", updateAffectedErr)
		return updateAffectedErr
	}

	if updateAffected == 0 {
		logger.Info("no rows were updated")
		return nil
	}

	return nil
}
