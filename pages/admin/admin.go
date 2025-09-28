package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/pages/admin/approval"
	"github.com/Regncon/conorganizer/pages/admin/approval/editForm"
	"github.com/Regncon/conorganizer/pages/index"
	"github.com/delaneyj/toolbelt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupAdminRoute(router chi.Router, store sessions.Store, logger *slog.Logger, kv jetstream.KeyValueConfig, db *sql.DB, eventImageDir *string) error {

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

	// eventLayoutRoute(router, db, err)
	// newEvent.NewEventLayoutRoute(router, db, err)

	router.Route("/admin", func(adminRouter chi.Router) {
		adminLayoutRoute(adminRouter, db, logger, err)
		adminRouter.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
			sse := datastar.NewSSE(w, r)

			sessionID, mvc, err := mvcSession(w, r)
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
					c := adminPage(db)
					if err := sse.PatchElementTempl(c); err != nil {
						sse.ConsoleError(err)
						return
					}
				}
			}
		})

		adminRouter.Route("/approval/", func(approvalRouter chi.Router) {
			approvalRouter.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
				sse := datastar.NewSSE(w, r)

				sessionID, mvc, err := mvcSession(w, r)
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
						c := approval.ApprovalPage(db, logger)
						if err := sse.PatchElementTempl(c); err != nil {
							sse.ConsoleError(err)
							return
						}
					}
				}
			})

			approvalRouter.Route("/edit", func(editEventRouter chi.Router) {
				editEventRouter.Route("/{id}", func(newIdRoute chi.Router) {
					edit_form.EditFormLayoutRoute(newIdRoute, db, eventImageDir, logger)
				})
				editEventRouter.Route("/api/{id}", func(newApiIdRouter chi.Router) {
					newApiIdRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
						sessionID, mvc, err := mvcSession(w, r)
						if err != nil {
							http.Error(w, fmt.Sprintf("failed to get session id: %v", err), http.StatusInternalServerError)
							return
						}

						eventId := chi.URLParam(r, "id")
						if eventId == "" {
							http.Error(w, "Event ID is required. Got: "+eventId, http.StatusBadRequest)
							return
						}

						sse := datastar.NewSSE(w, r)

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

								c := edit_form.EditEventFormPage(ctx, eventId, db, eventImageDir, logger)
								if err := sse.PatchElementTempl(c); err != nil {
									sse.ConsoleError(err)
									return
								}
							}
						}
					})
				})
			})
			approval.ApprovalLayoutRoute(approvalRouter, db, logger, err)
		})
	})

	return nil
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
