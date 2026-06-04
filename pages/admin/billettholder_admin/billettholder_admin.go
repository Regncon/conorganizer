package billettholderadmin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	addbillettholder "github.com/Regncon/conorganizer/pages/admin/billettholder_admin/add"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupBillettholderAdminRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, logger *slog.Logger, db *sql.DB) error {
	baseLogger := logger
	logger = logger.With("component", "billettholder_admin")
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
		TTL:         24 * time.Hour, // match the "connections" cookie lifetime so open pages keep state
		MaxBytes:    16 * 1024 * 1024,
	})
	if err != nil {
		return fmt.Errorf("error creating key value: %w", err)
	}

	notifyUpdate := func(sessionID string) {
		subj := fmt.Sprintf("billettholder.%s.updated", sessionID)
		logger.Debug("Publishing billettholder update")
		if err := nc.Publish(subj, nil); err != nil {
			logger.Error(fmt.Errorf("failed to publish billettholder page update for session %s: %w", sessionID, err).Error())
		}
	}

	resetAndSaveMVC := func(ctx context.Context, mvc *root.TodoMVC, sessionID string) error {
		*mvc = root.TodoMVC{}
		if err := saveMVC(ctx, mvc, sessionID, kv, notifyUpdate); err != nil {
			return fmt.Errorf("failed to save mvc: %w", err)
		}
		return nil
	}

	applyMVCEntry := func(ctx context.Context, mvc *root.TodoMVC, sessionID string, entry jetstream.KeyValueEntry) error {
		if entry.Operation() != jetstream.KeyValuePut {
			logger.Debug("resetting billettholder admin live update state after KV operation", "operation", entry.Operation().String())
			return resetAndSaveMVC(ctx, mvc, sessionID)
		}
		if err := json.Unmarshal(entry.Value(), mvc); err != nil {
			logger.Debug("resetting billettholder admin live update state after invalid KV value", "error", err.Error())
			return resetAndSaveMVC(ctx, mvc, sessionID)
		}
		return nil
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
			if err := resetAndSaveMVC(ctx, mvc, sessionID); err != nil {
				return "", nil, err
			}
		} else {
			if err := applyMVCEntry(ctx, mvc, sessionID, entry); err != nil {
				return "", nil, err
			}
		}
		return sessionID, mvc, nil
	}

	indexRoute(router, db, logger, err)

	router.Route("/admin/billettholder/api/", func(billettholderAdminRouter chi.Router) {
		billettholderAdminRouter.With(authctx.RequireAdmin(baseLogger)).Get("/", func(w http.ResponseWriter, r *http.Request) {
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
			defer func() {
				if err := sub.Unsubscribe(); err != nil {
					if errors.Is(err, nats.ErrBadSubscription) || ctx.Err() != nil {
						return
					}
					logger.Error(fmt.Errorf("failed to unsubscribe billettholder admin stream: %w", err).Error())
				}
			}()
			sse := datastar.NewSSE(w, r)

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
		addBillettholderRouter.With(authctx.RequireAdmin(baseLogger)).Get("/", func(w http.ResponseWriter, r *http.Request) {
			sessionID, _, err := session(w, r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			subj := fmt.Sprintf("billettholder.%s.updated", sessionID)
			logger.Debug("Subscribing add billettholder page")
			sub, err := nc.SubscribeSync(subj)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer func() {
				if err := sub.Unsubscribe(); err != nil {
					if errors.Is(err, nats.ErrBadSubscription) || ctx.Err() != nil {
						return
					}
					logger.Error(fmt.Errorf("failed to unsubscribe add-billettholder stream: %w", err).Error())
				}
			}()
			sse := datastar.NewSSE(w, r)

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
		id = uuid.NewString()
		sess.Values["id"] = id
		if err := sess.Save(r, w); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}
	}
	return id, nil
}
