package bilettholderadmin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	addbilettholder "github.com/Regncon/conorganizer/pages/admin/bilettholder_admin/add"
	"github.com/Regncon/conorganizer/pages/index"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar/sdk/go"
)

func SetupBilettholderAdminRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, logger *slog.Logger, db *sql.DB) error {
	nc, err := ns.Client()
	if err != nil {
		return fmt.Errorf("error creating nats client: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("error creating jetstream client: %w", err)
	}

	kv, err := js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "bilettholder",
		Description: "Bilettholder data",
		Compression: true,
		TTL:         time.Hour,
		MaxBytes:    16 * 1024 * 1024,
	})
	if err != nil {
		return fmt.Errorf("error creating key value: %w", err)
	}

	notifyUpdate := func(sessionID string) {
		subj := fmt.Sprintf("bilettholder.%s.updated", sessionID)
		fmt.Println("update bilettholeder subj", subj)
		if err := nc.Publish(subj, nil); err != nil {
			logger.Error("failed to publish page update", "err", err, "session", sessionID)
		}
	}

	session := func(w http.ResponseWriter, r *http.Request) (string, *index.TodoMVC, error) {
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
			// first visit ⇒ create an empty snapshot so the SSE loop can unmarshal
			if err := saveMVC(ctx, mvc, sessionID, kv, notifyUpdate); err != nil {
				return "", nil, fmt.Errorf("failed to save mvc: %w", err)
			}
		} else {
			if err := json.Unmarshal(entry.Value(), mvc); err != nil {
				return "", nil, fmt.Errorf("failed to unmarshal mvc: %w", err)
			}
		}
		return sessionID, mvc, nil
	}

	indexRoute(router, db, err)

	router.Route("/admin/bilettholder/api/", func(bilettholderAdminRouter chi.Router) {
		bilettholderAdminRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			sse := datastar.NewSSE(w, r)
			sessionID, _, err := session(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			subj := fmt.Sprintf("bilettholder.%s.updated", sessionID)
			sub, err := nc.SubscribeSync(subj)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer sub.Unsubscribe()

			// send the first render immediately
			if err := sse.MergeFragmentTempl(BilettholderAdminPage(db)); err != nil {
				_ = sse.ConsoleError(err)
				return
			}

			for {
				if _, err := sub.NextMsgWithContext(ctx); err != nil {
					return // context cancelled or sub closed
				}
				if err := sse.MergeFragmentTempl(BilettholderAdminPage(db)); err != nil {
					_ = sse.ConsoleError(err)
					return
				}
			}
		})
	})

	addbilettholder.AddBilettholderRoute(router, db, err)

	router.Route("/admin/bilettholder/add/api/", func(addBilettholderRouter chi.Router) {
		addBilettholderRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			sse := datastar.NewSSE(w, r)
			sessionID, _, err := session(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			subj := fmt.Sprintf("bilettholder.%s.updated", sessionID)
			fmt.Println("add bilettholder page subj", subj)
			sub, err := nc.SubscribeSync(subj)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer sub.Unsubscribe()

			// initial render
			if err := sse.MergeFragmentTempl(addbilettholder.AddBilettholderAdminPage(db, logger)); err != nil {
				_ = sse.ConsoleError(err)
				return
			}

			for {
				if _, err := sub.NextMsgWithContext(ctx); err != nil {
					return
				}
				if err := sse.MergeFragmentTempl(addbilettholder.AddBilettholderAdminPage(db, logger)); err != nil {
					_ = sse.ConsoleError(err)
					return
				}
			}
		})

		addbilettholder.CheckInTicketsSearchRoute(addBilettholderRouter, db, logger, store, notifyUpdate)
	})

	return nil
}

func saveMVC(ctx context.Context, mvc *index.TodoMVC, sessionID string, kv jetstream.KeyValue, poke func(string)) error {
	b, err := json.Marshal(mvc)
	if err != nil {
		return fmt.Errorf("failed to marshal mvc: %w", err)
	}
	if _, err := kv.Put(ctx, sessionID, b); err != nil {
		return fmt.Errorf("failed to put key value: %w", err)
	}
	if poke != nil {
		poke(sessionID)
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
