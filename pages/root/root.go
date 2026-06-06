package root

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupRootRoute(router chi.Router, store sessions.Store, logger *slog.Logger, ns *embeddednats.Server, db *sql.DB, eventImageDir *string) error {
	logger = logger.With("component", "root")
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

	saveMVC := func(ctx context.Context, sessionID string, mvc *TodoMVC) error {
		b, err := json.Marshal(mvc)
		if err != nil {
			return fmt.Errorf("failed to marshal mvc: %w", err)
		}
		if _, err := kv.Put(ctx, sessionID, b); err != nil {
			return fmt.Errorf("failed to put key value: %w", err)
		}
		return nil
	}

	resetMVC := func(mvc *TodoMVC) {
		mvc.Mode = TodoViewModeAll
		mvc.Todos = []*Todo{
			{Text: "Learn a backend language", Completed: true},
			{Text: "Learn Datastar", Completed: false},
			{Text: "Create Hypermedia", Completed: false},
			{Text: "???", Completed: false},
			{Text: "Profit", Completed: false},
		}
		mvc.EditingIdx = -1
	}

	resetAndSaveMVC := func(ctx context.Context, sessionID string, mvc *TodoMVC) error {
		resetMVC(mvc)
		if err := saveMVC(ctx, sessionID, mvc); err != nil {
			logger.Warn("failed to save root live update state; continuing without KV session", "error", err.Error())
		}
		return nil
	}

	applyMVCEntry := func(ctx context.Context, sessionID string, mvc *TodoMVC, entry jetstream.KeyValueEntry) error {
		if entry.Operation() != jetstream.KeyValuePut {
			logger.Debug("resetting root live update state after KV operation", "operation", entry.Operation().String())
			return resetAndSaveMVC(ctx, sessionID, mvc)
		}
		if err := json.Unmarshal(entry.Value(), mvc); err != nil {
			logger.Debug("resetting root live update state after invalid KV value", "error", err.Error())
			return resetAndSaveMVC(ctx, sessionID, mvc)
		}
		return nil
	}

	loadMVCSession := func(ctx context.Context, sessionID string, mvc *TodoMVC) error {
		if entry, err := kv.Get(ctx, sessionID); err != nil {
			if err != jetstream.ErrKeyNotFound {
				logger.Warn("failed to load root live update state; resetting", "error", err.Error())
				resetMVC(mvc)
				return nil
			}
			if err := resetAndSaveMVC(ctx, sessionID, mvc); err != nil {
				return err
			}
		} else {
			if err := applyMVCEntry(ctx, sessionID, mvc, entry); err != nil {
				return err
			}
		}
		return nil
	}

	rootLayoutRoute(router, db, logger, eventImageDir, err)

	router.Route("/root", func(rootRouter chi.Router) {
		rootRouter.Route("/api", func(rootApiRouter chi.Router) {
			rootApiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				sessionID, err := upsertSessionID(store, r, w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				ctx := r.Context()
				mvc := &TodoMVC{}
				sse := datastar.NewSSE(w, r)
				renderRootPage := func() error {
					isAdmin := authctx.GetAdminFromUserToken(ctx)
					return sse.PatchElementTempl(rootPage(db, isAdmin, eventImageDir))
				}
				if err := renderRootPage(); err != nil {
					_ = sse.ConsoleError(err)
					return
				}
				if err := loadMVCSession(ctx, sessionID, mvc); err != nil {
					logger.Error(fmt.Errorf("failed to load root live update session: %w", err).Error())
					_ = sse.ConsoleError(err)
					return
				}
				watcher, err := kv.Watch(ctx, sessionID)
				if err != nil {
					logger.Error(fmt.Errorf("failed to create root watcher: %w", err).Error())
					_ = sse.ConsoleError(err)
					return
				}
				defer func() {
					if err := watcher.Stop(); err != nil {
						if errors.Is(err, nats.ErrBadSubscription) || ctx.Err() != nil {
							return
						}
						logger.Error(fmt.Errorf("failed to stop root watcher: %w", err).Error())
					}
				}()

				for {
					select {
					case <-ctx.Done():
						return
					case entry, ok := <-watcher.Updates():
						if !ok {
							return
						}
						if entry == nil {
							continue
						}
						if err := applyMVCEntry(ctx, sessionID, mvc, entry); err != nil {
							logger.Error(fmt.Errorf("failed to apply root live update state: %w", err).Error())
							_ = sse.ConsoleError(err)
							return
						}
						if err := renderRootPage(); err != nil {
							_ = sse.ConsoleError(err)
							return
						}
					}
				}
			})
		})
	})

	return nil
}

func MustJSONMarshal(v any) string {
	b, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func upsertSessionID(store sessions.Store, r *http.Request, w http.ResponseWriter) (string, error) {

	sess, err := store.Get(r, "connections")
	if err != nil {
		return "", fmt.Errorf("failed to get session: %w", err)
	}
	id, ok := sess.Values["id"].(string)
	if !ok {
		id = uuid.NewString()
		sess.Values["id"] = id
		if err := sess.Save(r, w); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}
	}
	return id, nil
}
