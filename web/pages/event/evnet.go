package event

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/web/pages/index"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar/sdk/go"
)

func SetupEventRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, db *sql.DB) error {
	nc, err := ns.Client()
	if err != nil {
		return fmt.Errorf("error creating nats client: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("error creating jetstream client: %w", err)
	}

	kv, err := js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "todos",
		Description: "Datastar Todos",
		Compression: true,
		TTL:         time.Hour,
		MaxBytes:    16 * 1024 * 1024,
	})

	if err != nil {
		return fmt.Errorf("error creating key value: %w", err)
	}

	resetMVC := func(mvc *index.TodoMVC) {
		mvc.Mode = index.TodoViewModeAll
		mvc.Todos = []*index.Todo{
			{Text: "Learn a backend language", Completed: true},
			{Text: "Learn Datastar", Completed: false},
			{Text: "Create Hypermedia", Completed: false},
			{Text: "???", Completed: false},
			{Text: "Profit", Completed: false},
		}
		mvc.EditingIdx = -1
	}

	mvcSession := func(w http.ResponseWriter, r *http.Request) (string, *index.TodoMVC, error) {
		ctx := r.Context()
		sessionID, err := upsertSessionID(store, r, w)
		if err != nil {
			return "", nil, fmt.Errorf("failed to get session id: %w", err)
		}

		mvc := &index.TodoMVC{}
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
	eventLayoutRoute(router, db, err)

	router.Route("/event/api/{idx}/", func(eventRouter chi.Router) {
		eventRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
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
					c := event_page(eventID, db)
					if err := sse.MergeFragmentTempl(c); err != nil {
						sse.ConsoleError(err)
						return
					}
				}
			}

		})
		editRout(eventRouter, db, kv)
	})

	return nil
}

func getEventByID(id string, db *sql.DB) (*models.Event, error) {
	query := "SELECT id, title, description FROM events WHERE id = ?"
	row := db.QueryRow(query, id)

	var event models.Event
	if err := row.Scan(&event.ID, &event.Title, &event.Description); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No event found
		}
		return nil, err
	}

	return &event, nil
}

func saveMVC(ctx context.Context, mvc *index.TodoMVC, sessionID string, kv jetstream.KeyValue) error {
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
