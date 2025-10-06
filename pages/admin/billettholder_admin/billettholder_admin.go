package billettholderadmin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	addbillettholder "github.com/Regncon/conorganizer/pages/admin/billettholder_admin/add"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupBillettholderAdminRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, logger *slog.Logger, db *sql.DB) error {
	nc, err := ns.Client()
	if err != nil {
		return fmt.Errorf("error creating nats client: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return fmt.Errorf("error creating jetstream client: %w", err)
	}

	kv, err := js.CreateOrUpdateKeyValue(context.Background(), jetstream.KeyValueConfig{
		Bucket:      "billettholder",
		Description: "Billettholder data",
		Compression: true,
		TTL:         time.Hour,
		MaxBytes:    16 * 1024 * 1024,
	})
	if err != nil {
		return fmt.Errorf("error creating key value: %w", err)
	}

	notifyUpdate := func(sessionID string) {
		subj := fmt.Sprintf("billettholder.%s.updated", sessionID)
		fmt.Println("update billettholder subj", subj)
		if err := nc.Publish(subj, nil); err != nil {
			logger.Error("failed to publish page update", "err", err, "session", sessionID)
		}
	}

	session := func(w http.ResponseWriter, r *http.Request) (string, *root.TodoMVC, error) {
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

	router.Route("/admin/billettholder/api/", func(billettholderAdminRouter chi.Router) {
		billettholderAdminRouter.With(authctx.RequireAdmin(logger)).Get("/", func(w http.ResponseWriter, r *http.Request) {
			sse := datastar.NewSSE(w, r)
			sessionID, _, err := session(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			subj := fmt.Sprintf("billettholder.%s.updated", sessionID)
			sub, err := nc.SubscribeSync(subj)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer sub.Unsubscribe()

			// send the first render immediately
			if err := sse.PatchElementTempl(BillettholderAdminPage(db, logger)); err != nil {
				_ = sse.ConsoleError(err)
				return
			}

			for {
				if _, err := sub.NextMsgWithContext(ctx); err != nil {
					return // context cancelled or sub closed
				}
				if err := sse.PatchElementTempl(BillettholderAdminPage(db, logger)); err != nil {
					_ = sse.ConsoleError(err)
					return
				}
			}
		})
		billettholdereSearchRoute(billettholderAdminRouter, store, notifyUpdate)
		addEmailToBilettholderRoute(billettholderAdminRouter, db, logger, store, notifyUpdate)
		deleteEmailFromBillettholderRoute(billettholderAdminRouter, db, logger, store, notifyUpdate)
	})

	addbillettholder.AddBillettholderRoute(router, db, logger, err)

	router.Route("/admin/billettholder/add/api/", func(addBillettholderRouter chi.Router) {
		addBillettholderRouter.With(authctx.RequireAdmin(logger)).Get("/", func(w http.ResponseWriter, r *http.Request) {
			sse := datastar.NewSSE(w, r)
			sessionID, _, err := session(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			subj := fmt.Sprintf("billettholder.%s.updated", sessionID)
			fmt.Println("add billettholder page subj", subj)
			sub, err := nc.SubscribeSync(subj)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer sub.Unsubscribe()

			// initial render
			if err := sse.PatchElementTempl(addbillettholder.AddBillettholderAdminPage(db, logger)); err != nil {
				_ = sse.ConsoleError(err)
				return
			}

			for {
				if _, err := sub.NextMsgWithContext(ctx); err != nil {
					return
				}
				if err := sse.PatchElementTempl(addbillettholder.AddBillettholderAdminPage(db, logger)); err != nil {
					_ = sse.ConsoleError(err)
					return
				}
			}
		})

		addbillettholder.CheckInTicketsSearchRoute(addBillettholderRouter, db, logger, store, notifyUpdate)
		addbillettholder.ConvertTicketToBillettholderRoute(addBillettholderRouter, db, store, notifyUpdate, logger)
	})

	return nil
}

func saveMVC(ctx context.Context, mvc *root.TodoMVC, sessionID string, kv jetstream.KeyValue, poke func(string)) error {
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
