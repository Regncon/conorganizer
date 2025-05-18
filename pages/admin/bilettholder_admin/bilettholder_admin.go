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

// SetupBilettholderAdminRoute wires the admin UI (SSE) + the API routes.
//
// The page is refreshed whenever we receive a core‑NATS message on
//
//	"bilettholder.<sessionID>.updated".
//
// These messages are published by:
//   - saveMVC – whenever Todo‑state is changed
//   - the ticket‑check‑in search route (via a callback we pass in)
func SetupBilettholderAdminRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, logger *slog.Logger, db *sql.DB) error {
	// --------------------------------------------------------------------------------
	// NATS set‑up (core client + JetStream for KV snapshots)
	// --------------------------------------------------------------------------------
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

	// Helper that publishes the session‑scoped update poke
	notifyUpdate := func(sessionID string) {
		subj := fmt.Sprintf("bilettholder.%s.updated", sessionID)
		fmt.Println("update bilettholeder subj", subj)
		if err := nc.Publish(subj, nil); err != nil {
			logger.Error("failed to publish page update", "err", err, "session", sessionID)
		}
	}

	// -----------------------------------------------------------------------------
	// Session helper – identical to the previous version except it no longer needs
	// the KV watcher to always stay perfectly in‑sync with DB. We persist a JSON
	// snapshot only because the SSE handlers still expect one.
	// -----------------------------------------------------------------------------
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

	// initial index page (non‑SSE)
	indexRoute(router, db, err)

	// -----------------------------------------------------------------------------
	// /admin/bilettholder – main admin list UI (SSE)
	// -----------------------------------------------------------------------------
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

	// -----------------------------------------------------------------------------
	// /admin/bilettholder/add – add/check‑in UI (SSE)
	// -----------------------------------------------------------------------------
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

		// Register the search/check‑in routes and tell them how to broadcast updates
		addbilettholder.CheckInTicketsSearchRoute(addBilettholderRouter, db, logger, store, notifyUpdate)
	})

	return nil
}

// saveMVC writes the JSON snapshot + pokes subscribers so that any tab listening
// to "bilettholder.<sessionID>.updated" re‑renders.
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
