package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/Regncon/conorganizer/components/formsubmission"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/pages/admin/approval"
	edit_form "github.com/Regncon/conorganizer/pages/admin/approval/editForm"
	"github.com/Regncon/conorganizer/pages/admin/rooms"
	"github.com/Regncon/conorganizer/pages/root"
	"github.com/Regncon/conorganizer/service/keyvalue"
	roomService "github.com/Regncon/conorganizer/service/rooms"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupAdminRoute(router chi.Router, store sessions.Store, logger *slog.Logger, ns *embeddednats.Server, db *sql.DB, eventImageDir *string) error {
	baseLogger := logger
	logger = logger.With("component", "admin")
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

	resetAndSaveMVC := func(ctx context.Context, mvc *root.TodoMVC, sessionID string) error {
		resetMVC(mvc)
		if err := saveMVC(ctx, mvc, sessionID, kv); err != nil {
			logger.Warn("failed to save admin live update state; continuing without KV session", "error", err.Error())
		}
		return nil
	}

	applyMVCEntry := func(ctx context.Context, mvc *root.TodoMVC, sessionID string, entry jetstream.KeyValueEntry) error {
		if entry.Operation() != jetstream.KeyValuePut {
			logger.Debug("resetting admin live update state after KV operation", "operation", entry.Operation().String())
			return resetAndSaveMVC(ctx, mvc, sessionID)
		}
		if err := json.Unmarshal(entry.Value(), mvc); err != nil {
			logger.Debug("resetting admin live update state after invalid KV value", "error", err.Error())
			return resetAndSaveMVC(ctx, mvc, sessionID)
		}
		return nil
	}

	loadMVCSession := func(ctx context.Context, mvc *root.TodoMVC, sessionID string) error {
		if entry, err := kv.Get(ctx, sessionID); err != nil {
			if err != jetstream.ErrKeyNotFound {
				logger.Warn("failed to load admin live update state; resetting", "error", err.Error())
				resetMVC(mvc)
				return nil
			}
			if err := resetAndSaveMVC(ctx, mvc, sessionID); err != nil {
				return err
			}
		} else {
			if err := applyMVCEntry(ctx, mvc, sessionID, entry); err != nil {
				return err
			}
		}
		return nil
	}

	// eventLayoutRoute(router, db, err)
	// newEvent.NewEventLayoutRoute(router, db, err)

	router.Route("/admin", func(adminRouter chi.Router) {
		adminLayoutRoute(adminRouter, db, logger, err)
		puljefordelingStatusRoute(adminRouter, db, kv, logger)
		programPublishingRoute(adminRouter, db, kv, logger)
		adminRouter.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
			requestID := middleware.GetReqID(r.Context())
			_, connectionsCookieErr := r.Cookie("connections")
			connectionsCookiePresent := connectionsCookieErr == nil
			logger.Debug(
				"admin live update stream starting",
				"request_id", requestID,
				"connections_cookie_present", connectionsCookiePresent,
			)

			sessionID, err := upsertSessionID(store, r, w)
			if err != nil {
				logger.Error(
					fmt.Errorf("failed to get session id: %w", err).Error(),
					"request_id", requestID,
					"connections_cookie_present", connectionsCookiePresent,
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			ctx := r.Context()
			mvc := &root.TodoMVC{}
			sse := datastar.NewSSE(w, r)
			if err := sse.PatchElementTempl(adminPage(db)); err != nil {
				logger.Error(
					fmt.Errorf("failed to patch initial admin live update: %w", err).Error(),
					"request_id", requestID,
					"connections_cookie_present", connectionsCookiePresent,
				)
				_ = sse.ConsoleError(err)
				return
			}
			if err := loadMVCSession(ctx, mvc, sessionID); err != nil {
				logger.Error(
					fmt.Errorf("failed to load admin live update session: %w", err).Error(),
					"request_id", requestID,
					"connections_cookie_present", connectionsCookiePresent,
				)
				_ = sse.ConsoleError(err)
				return
			}
			watcher, err := kv.Watch(ctx, sessionID)
			if err != nil {
				logger.Error(
					fmt.Errorf("failed to create admin watcher: %w", err).Error(),
					"request_id", requestID,
					"connections_cookie_present", connectionsCookiePresent,
				)
				_ = sse.ConsoleError(err)
				return
			}
			defer func() {
				if err := watcher.Stop(); err != nil {
					if errors.Is(err, nats.ErrBadSubscription) || ctx.Err() != nil {
						logger.Debug(
							"admin watcher already stopped",
							"request_id", requestID,
							"connections_cookie_present", connectionsCookiePresent,
							"error", err.Error(),
						)
						return
					}
					logger.Error(
						fmt.Errorf("failed to stop admin watcher: %w", err).Error(),
						"request_id", requestID,
						"connections_cookie_present", connectionsCookiePresent,
					)
				}
			}()
			logger.Debug(
				"admin watcher started",
				"request_id", requestID,
				"connections_cookie_present", connectionsCookiePresent,
			)

			for {
				select {
				case <-ctx.Done():
					logger.Debug(
						"admin live update stream closed",
						"request_id", requestID,
						"connections_cookie_present", connectionsCookiePresent,
						"reason", ctx.Err().Error(),
					)
					return
				case entry, ok := <-watcher.Updates():
					if !ok {
						if ctx.Err() != nil {
							logger.Debug(
								"admin watcher updates channel closed",
								"request_id", requestID,
								"connections_cookie_present", connectionsCookiePresent,
								"reason", ctx.Err().Error(),
							)
							return
						}
						logger.Warn(
							"admin watcher updates channel closed",
							"request_id", requestID,
							"connections_cookie_present", connectionsCookiePresent,
						)
						return
					}
					if entry == nil {
						logger.Debug(
							"admin watcher initial values delivered",
							"request_id", requestID,
							"connections_cookie_present", connectionsCookiePresent,
						)
						continue
					}
					if err := applyMVCEntry(ctx, mvc, sessionID, entry); err != nil {
						logger.Error(
							fmt.Errorf("failed to apply admin live update state: %w", err).Error(),
							"request_id", requestID,
							"connections_cookie_present", connectionsCookiePresent,
							"operation", entry.Operation().String(),
							"revision", entry.Revision(),
							"value_bytes", len(entry.Value()),
						)
						_ = sse.ConsoleError(err)
						return
					}
					if err := sse.PatchElementTempl(adminPage(db)); err != nil {
						logger.Error(
							fmt.Errorf("failed to patch admin live update: %w", err).Error(),
							"request_id", requestID,
							"connections_cookie_present", connectionsCookiePresent,
							"operation", entry.Operation().String(),
							"revision", entry.Revision(),
						)
						_ = sse.ConsoleError(err)
						return
					}
					logger.Debug(
						"admin live update patch sent",
						"request_id", requestID,
						"connections_cookie_present", connectionsCookiePresent,
						"operation", entry.Operation().String(),
						"revision", entry.Revision(),
					)
				}
			}
		})

		adminRouter.Route("/approval/", func(approvalRouter chi.Router) {
			approvalRouter.Route("/api/", func(apiRouter chi.Router) {
				apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
					sessionID, err := upsertSessionID(store, r, w)
					if err != nil {
						logger.Error(fmt.Errorf("failed to get approval session id: %w", err).Error())
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					ctx := r.Context()
					mvc := &root.TodoMVC{}
					sse := datastar.NewSSE(w, r)
					if err := sse.PatchElementTempl(approval.ApprovalPage(db, baseLogger)); err != nil {
						_ = sse.ConsoleError(err)
						return
					}
					if err := loadMVCSession(ctx, mvc, sessionID); err != nil {
						logger.Error(fmt.Errorf("failed to load approval live update session: %w", err).Error())
						_ = sse.ConsoleError(err)
						return
					}
					watcher, err := kv.Watch(ctx, sessionID)
					if err != nil {
						logger.Error(fmt.Errorf("failed to create approval watcher: %w", err).Error())
						_ = sse.ConsoleError(err)
						return
					}
					defer func() {
						if err := watcher.Stop(); err != nil {
							if errors.Is(err, nats.ErrBadSubscription) || ctx.Err() != nil {
								return
							}
							logger.Error(fmt.Errorf("failed to stop approval watcher: %w", err).Error())
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
							if err := applyMVCEntry(ctx, mvc, sessionID, entry); err != nil {
								logger.Error(fmt.Errorf("failed to apply approval live update state: %w", err).Error())
								_ = sse.ConsoleError(err)
								return
							}
							if err := sse.PatchElementTempl(approval.ApprovalPage(db, baseLogger)); err != nil {
								_ = sse.ConsoleError(err)
								return
							}
						}
					}
				})

				apiRouter.Route("/event-players", func(eventPlayersRouter chi.Router) {
					eventPlayersRouter.Post("/post/add_first_choice", func(w http.ResponseWriter, r *http.Request) {
						type Store struct {
							BillettholderId int    `json:"assignmentBillettholderId"`
							EventId         string `json:"assignmentEventId"`
							PuljeId         string `json:"assignmentPuljeId"`
						}

						store := &Store{}

						if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
							return
						}
						if store.BillettholderId <= 0 {
							logger.Error(fmt.Errorf("invalid billettholder id for add first choice (event_id=%s, pulje_id=%s): invalid assignmentBillettholderId %d: must be greater than 0", store.EventId, store.PuljeId, store.BillettholderId).Error())
							http.Error(w, fmt.Errorf("invalid assignmentBillettholderId %d: must be greater than 0", store.BillettholderId).Error(), http.StatusNotFound)
							return
						}

						var addFirstChoiceErr = formsubmission.AddPlayersFirstChoice(
							store.BillettholderId,
							store.EventId,
							store.PuljeId,
							db,
							baseLogger,
						)
						if addFirstChoiceErr != nil {
							logger.Error(fmt.Errorf("failed to add player as first choice: %w", addFirstChoiceErr).Error())
							http.Error(w, addFirstChoiceErr.Error(), http.StatusInternalServerError)
							return
						}
						logger.Info(
							"Successfully added player as first choice",
							"event_id", store.EventId,
							"pulje_id", store.PuljeId,
							"billettholder_id", store.BillettholderId,
						)
						if err := keyvalue.BroadcastUpdate(kv, r); err != nil {
							logger.Error(fmt.Errorf("failed to broadcast add first choice update: %w", err).Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}

					})
					eventPlayersRouter.Post("/post/add_gm", func(w http.ResponseWriter, r *http.Request) {

						type Store struct {
							BillettholderId int    `json:"assignmentBillettholderId"`
							EventId         string `json:"assignmentEventId"`
							PuljeId         string `json:"assignmentPuljeId"`
						}
						store := &Store{}

						if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
							return
						}
						if store.BillettholderId <= 0 {
							logger.Error(fmt.Errorf("invalid billettholder id for add GM (event_id=%s, pulje_id=%s): invalid assignmentBillettholderId %d: must be greater than 0", store.EventId, store.PuljeId, store.BillettholderId).Error())
							http.Error(w, fmt.Errorf("invalid assignmentBillettholderId %d: must be greater than 0", store.BillettholderId).Error(), http.StatusNotFound)
							return
						}

						var updatePlayerStatusErr = formsubmission.UpdatePlayerStatus(
							store.EventId,
							store.PuljeId,
							store.BillettholderId,
							false,
							true,
							db,
							baseLogger,
						)
						if updatePlayerStatusErr != nil {
							logger.Error(fmt.Errorf("failed to add player as GM: %w", updatePlayerStatusErr).Error())
							http.Error(w, updatePlayerStatusErr.Error(), http.StatusInternalServerError)
							return
						}
						logger.Info(
							"Successfully Added player as GM",
							"event_id", store.EventId,
							"pulje_id", store.PuljeId,
							"billettholder_id", store.BillettholderId,
							"role", models.EventPlayerRoleGM,
						)
						if err := keyvalue.BroadcastUpdate(kv, r); err != nil {
							logger.Error(fmt.Errorf("failed to broadcast add GM update: %w", err).Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}
					})
					eventPlayersRouter.Put("/update_status", func(w http.ResponseWriter, r *http.Request) {
						type Store struct {
							BillettholderId int    `json:"assignmentBillettholderId"`
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

						var updatePlayerStatusErr = formsubmission.UpdatePlayerStatus(
							store.EventId,
							store.PuljeId,
							store.BillettholderId,
							store.IsPlayer,
							store.IsGm,
							db,
							baseLogger,
						)
						if updatePlayerStatusErr != nil {
							http.Error(w, updatePlayerStatusErr.Error(), http.StatusInternalServerError)
							return
						}
						logger.Info(
							"Successfully updated player status",
							"event_id", store.EventId,
							"pulje_id", store.PuljeId,
							"billettholder_id", store.BillettholderId,
							"assignment_is_player", store.IsPlayer,
							"assignment_is_gm", store.IsGm,
						)
						if err := keyvalue.BroadcastUpdate(kv, r); err != nil {
							logger.Error(fmt.Errorf("failed to broadcast player status update: %w", err).Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}
					})
				})
			})

			approvalRouter.Route("/edit", func(editEventRouter chi.Router) {
				editEventRouter.Route("/{id}", func(newIdRoute chi.Router) {
					edit_form.EditFormLayoutRoute(newIdRoute, db, eventImageDir, baseLogger)
				})
				editEventRouter.Route("/api/{id}", func(newApiIdRouter chi.Router) {
					newApiIdRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
						sessionID, err := upsertSessionID(store, r, w)
						if err != nil {
							http.Error(w, fmt.Sprintf("failed to get session id: %v", err), http.StatusInternalServerError)
							return
						}

						eventId := chi.URLParam(r, "id")
						if eventId == "" {
							http.Error(w, "Event ID is required. Got: "+eventId, http.StatusBadRequest)
							return
						}

						ctx := r.Context()
						mvc := &root.TodoMVC{}
						sse := datastar.NewSSE(w, r)
						if err := sse.PatchElementTempl(edit_form.EditEventFormPage(ctx, eventId, db, eventImageDir, baseLogger)); err != nil {
							_ = sse.ConsoleError(err)
							return
						}
						if err := loadMVCSession(ctx, mvc, sessionID); err != nil {
							logger.Error(fmt.Errorf("failed to load edit-form live update session: %w", err).Error(), "event_id", eventId)
							_ = sse.ConsoleError(err)
							return
						}
						watcher, err := kv.Watch(ctx, sessionID)
						if err != nil {
							logger.Error(fmt.Errorf("failed to create edit-form watcher: %w", err).Error(), "event_id", eventId)
							_ = sse.ConsoleError(err)
							return
						}
						defer func() {
							if err := watcher.Stop(); err != nil {
								if errors.Is(err, nats.ErrBadSubscription) || ctx.Err() != nil {
									return
								}
								logger.Error(fmt.Errorf("failed to stop edit-form watcher: %w", err).Error())
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
								if err := applyMVCEntry(ctx, mvc, sessionID, entry); err != nil {
									logger.Error(fmt.Errorf("failed to apply edit-form live update state: %w", err).Error(), "event_id", eventId)
									_ = sse.ConsoleError(err)
									return
								}

								c := edit_form.EditEventFormPage(ctx, eventId, db, eventImageDir, baseLogger)
								if err := sse.PatchElementTempl(c); err != nil {
									_ = sse.ConsoleError(err)
									return
								}
							}
						}
					})
				})
			})
			approval.ApprovalLayoutRoute(approvalRouter, db, baseLogger, err)
		})

		adminRouter.Route("/rooms", func(roomRouter chi.Router) {
			roomRouter.Route("/api", func(roomApiRouter chi.Router) {
				roomApiRouter.Route("/create", func(createRoomRoute chi.Router) {
					createRoomRoute.Get("/", func(w http.ResponseWriter, r *http.Request) {
						// Reset form bindings before creating new room
						sse := datastar.NewSSE(w, r)
						payload, _ := json.Marshal(map[string]any{
							"mode":         "create",
							"form_title":   "Legg til rom",
							"button_label": "Opprett rom",
							"submit_url":   "/admin/rooms/api/create",

							"id":                   0,
							"name":                 "",
							"room_number":          "",
							"floor":                0,
							"max_concurrent_games": 0,
							"notes":                "",
							"is_disabled":          false,
							"error":                "",
						})

						if err = sse.PatchSignals(payload); err != nil {
							logger.Error("Failed to patch signals", "error", err.Error())
							http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							return
						}
					})

					createRoomRoute.Post("/", func(w http.ResponseWriter, r *http.Request) {
						// Read data-star post submission
						store := &models.Room{}
						if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
							fmt.Println(readSignalErr.Error())
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
						}

						// Get ready to broadcast responses to client
						sse := datastar.NewSSE(w, r)

						// Validate input

						// Create room
						_, err := roomService.CreateRoom(db, *store)
						if err != nil {
							payload, _ := json.Marshal(map[string]string{
								"error": err.Error(),
							})

							if err = sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							}
							return
						}

						// Ridirect on success
						_ = sse.Redirect("/admin/rooms")
					})
				})

				roomApiRouter.Route("/edit/{id}", func(updateRoomRoute chi.Router) {
					updateRoomRoute.Get("/", func(w http.ResponseWriter, r *http.Request) {
						// Read url for room id
						roomID := chi.URLParam(r, "id")
						if roomID == "" {
							http.Error(w, "Room ID is required. Got: "+roomID, http.StatusBadRequest)
							return
						}
						id, err := strconv.ParseInt(roomID, 10, 0)
						if err != nil {
							http.Error(w, "", http.StatusBadRequest)
							return
						}

						// Get room from id in url param
						room, err := roomService.GetRoomByID(db, int(id))
						if err != nil {
							http.Error(w, "", http.StatusBadRequest)
							return
						}

						// Update formbindings with room data
						sse := datastar.NewSSE(w, r)
						payload, _ := json.Marshal(map[string]any{
							"mode":         "edit",
							"form_title":   "Rediger rom",
							"button_label": "Lagre endringer",
							"submit_url":   fmt.Sprintf("/admin/rooms/api/edit/%d", room.ID),

							"id":                   room.ID,
							"name":                 room.Name,
							"room_number":          room.RoomNumber,
							"floor":                room.Floor,
							"max_concurrent_games": room.MaxConcurrentGames,
							"notes":                room.Notes,
							"is_disabled":          room.IsDisabled,
							"error":                "",
						})
						if err = sse.PatchSignals(payload); err != nil {
							logger.Error("Failed to patch signals", "error", err.Error())
							http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							return
						}
					})

					updateRoomRoute.Post("/", func(w http.ResponseWriter, r *http.Request) {
						// Read url for room id
						roomID := chi.URLParam(r, "id")
						if roomID == "" {
							http.Error(w, "Room ID is required. Got: "+roomID, http.StatusBadRequest)
							return
						}

						// Read data-star post submission
						store := &models.Room{}
						if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
							fmt.Println(readSignalErr.Error())
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
						}

						// Get ready to broadcast responses to client
						sse := datastar.NewSSE(w, r)

						// Validate input
						parsedID, err := strconv.ParseInt(roomID, 10, 0)
						if err != nil {
							payload, _ := json.Marshal(map[string]string{
								"error": err.Error(),
							})

							if err = sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							}
							return
						}
						store.ID = int(parsedID)

						// Update room
						_, err = roomService.UpdateRoom(db, *store)
						if err != nil {
							payload, _ := json.Marshal(map[string]string{
								"error": err.Error(),
							})

							if err = sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)

							}
							return
						}

						// Ridirect on success
						_ = sse.Redirect("/admin/rooms")
					})
				})
			})
			rooms.RoomsLayoutRoute(roomRouter, db, baseLogger, err)
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
		id = uuid.NewString()
		sess.Values["id"] = id
		if err := sess.Save(r, w); err != nil {
			return "", fmt.Errorf("failed to save session: %w", err)
		}
	}
	return id, nil
}
