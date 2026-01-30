package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Regncon/conorganizer/components/formsubmission"
	"github.com/Regncon/conorganizer/pages/admin/approval"
	edit_form "github.com/Regncon/conorganizer/pages/admin/approval/editForm"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service/keyvalue"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupAdminRoute(router chi.Router, store sessions.Store, logger *slog.Logger, ns *embeddednats.Server, db *sql.DB, eventImageDir *string) error {
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
			approvalRouter.Route("/api/", func(apiRouter chi.Router) {
				apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
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

				apiRouter.Route("/events_players", func(eventsPlayersRouter chi.Router) {
					eventsPlayersRouter.Post("/post/add_gm", func(w http.ResponseWriter, r *http.Request) {
						// sse := datastar.NewSSE(w, r)
						// sessionID, mvc, err := mvcSession(w, r)
						// ctx := r.Context()
						type Store struct {
							BillettholderId int    `json:"assignmentGmSelect"`
							EventId         string `json:"assignmentEventId"`
							PuljeId         string `json:"assignmentPuljeId"`
							IsPlayer        bool   `json:"assignmentIsPlayer"`
							IsGm            bool   `json:"assignmentIsGm"`
						}
						store := &Store{}

						if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
							return
						}
						fmt.Printf("Store: %+v\n", store)
						// get from form data
						if err := r.ParseForm(); err != nil {
							http.Error(w, "Failed to parse form data", http.StatusBadRequest)
							return
						}

						fmt.Printf("Form data: %+v\n", r.Form.Get("name"))
						fmt.Printf("Form data: %+v\n", r.Form.Get("billettholder-6"))

						http.Error(w, "Not implemented", http.StatusNotImplemented)
						return
					})
					eventsPlayersRouter.Put("/put/{eventId}/{puljeId}/{billettholderId}/{isPlayer}/{isGm}", func(w http.ResponseWriter, r *http.Request) {
						// sse := datastar.NewSSE(w, r)
						// sessionID, mvc, err := mvcSession(w, r)
						// ctx := r.Context()

						eventId := chi.URLParam(r, "eventId")
						puljeId := chi.URLParam(r, "puljeId")
						billettholderId, billettholderIdToIntErr := strconv.Atoi(chi.URLParam(r, "billettholderId"))
						if billettholderIdToIntErr != nil {
							logger.Error("Invalid billettholder ID", slog.String("billettholderId", chi.URLParam(r, "billettholderId")), "billettholderIdToIntErr", billettholderIdToIntErr)
							http.Error(w, "Invalid billettholder ID", http.StatusBadRequest)
							return
						}
						isGm, isGmToBoolErr := strconv.ParseBool(chi.URLParam(r, "isGm"))
						if isGmToBoolErr != nil {
							logger.Error("Invalid isGm value", slog.String("isGm", chi.URLParam(r, "isGm")), "isGmToBoolErr", isGmToBoolErr)
							http.Error(w, "Invalid isGm value", http.StatusBadRequest)
							return
						}
						isPlayer, isPlayerToBoolErr := strconv.ParseBool(chi.URLParam(r, "isPlayer"))
						if isPlayerToBoolErr != nil {
							logger.Error("Invalid isPlayer value", slog.String("isPlayer", chi.URLParam(r, "isPlayer")), "isPlayerToBoolErr", isPlayerToBoolErr)
							http.Error(w, "Invalid isPlayer value", http.StatusBadRequest)
							return
						}
						// type Store struct {
						// 	Input string `json:"title"`
						// }
						// store := &Store{}

						// if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
						// 	http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
						// 	return
						// }

						var updatePlayerStatusErr = formsubmission.UpdatePlayerStatus(
							eventId,
							puljeId,
							billettholderId,
							isPlayer,
							isGm,
							db,
							logger,
						)
						if updatePlayerStatusErr != nil {
							http.Error(w, updatePlayerStatusErr.Error(), http.StatusInternalServerError)
							return
						}
						logger.Info(
							"Successfully updated player status",
							"eventId", eventId,
							"puljeId", puljeId,
							"billettholderId", billettholderId,
							"isPlayer", isPlayer,
							"isGm", isGm,
						)
						keyvalue.BroadcastUpdate(kv, r)
					})
				})
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
