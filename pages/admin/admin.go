package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Regncon/conorganizer/components/formsubmission"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/pages/admin/approval"
	edit_form "github.com/Regncon/conorganizer/pages/admin/approval/editForm"
	"github.com/Regncon/conorganizer/pages/admin/rooms"
	"github.com/Regncon/conorganizer/service/live"
	roomService "github.com/Regncon/conorganizer/service/rooms"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	datastar "github.com/starfederation/datastar-go/datastar"
)

func SetupAdminRoute(router chi.Router, logger *slog.Logger, liveManager *live.Manager, db *sql.DB, eventImageDir *string) error {
	baseLogger := logger
	logger = logger.With("component", "admin")


	router.Route("/admin", func(adminRouter chi.Router) {
		adminLayoutRoute(adminRouter, db, logger)
		puljefordelingStatusRoute(adminRouter, db, liveManager, logger)
		programPublishingRoute(adminRouter, db, liveManager, logger)
		adminRouter.Get("/api/", func(w http.ResponseWriter, r *http.Request) {
            var ctx = r.Context()
            userInfo := userctx.GetUserRequestInfo(ctx)
			liveManager.Stream(w, r, live.Page{
				Buckets: []live.Bucket{live.BucketEvents},
				Render: func(ctx context.Context, r *http.Request) templ.Component {
					return adminPage(userInfo, db, logger)
				},
			})
		})

		adminRouter.Route("/approval/", func(approvalRouter chi.Router) {
			approvalRouter.Route("/api/", func(apiRouter chi.Router) {
				apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
					liveManager.Stream(w, r, live.Page{
						Buckets: []live.Bucket{live.BucketEvents, live.BucketInterests, live.BucketBillettholders},
						Render: func(ctx context.Context, r *http.Request) templ.Component {
							return approval.ApprovalPage(db, baseLogger)
						},
					})
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
						if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
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
						if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
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
						if err := liveManager.Broadcast(r.Context(), live.BucketInterests); err != nil {
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
						eventId := chi.URLParam(r, "id")
						if eventId == "" {
							http.Error(w, "Event ID is required. Got: "+eventId, http.StatusBadRequest)
							return
						}

						liveManager.Stream(w, r, live.Page{
							Buckets: []live.Bucket{live.BucketEvents, live.BucketRooms},
							Render: func(ctx context.Context, r *http.Request) templ.Component {
								return edit_form.EditEventFormPage(ctx, eventId, db, eventImageDir, baseLogger)
							},
						})
					})
				})
			})
			approval.ApprovalLayoutRoute(approvalRouter, db, baseLogger)
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
							"floor":                1,
							"max_concurrent_games": 1,
							"notes":                "",
							"is_disabled":          false,
							"errors": map[string]string{
								"room_number":          "",
								"max_concurrent_games": "",
								"error":                "",
							},
						})

						if err := sse.PatchSignals(payload); err != nil {
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
						errors := map[string]string{
							"room_number":          "",
							"max_concurrent_games": "",
							"error":                "",
						}

						if !strings.HasPrefix(store.RoomNumber, fmt.Sprintf("%d", store.Floor)) {
							errors["room_number"] = "Romnummer må starte med etasje som første tall"
						}

						if strings.TrimSpace(store.RoomNumber) == "" {
							errors["room_number"] = "Romnummer er påkrevd"
						}

						if store.MaxConcurrentGames < 1 {
							errors["max_concurrent_games"] = "Maks samtidige spill må være minst 1"
						}

						hasErrors := false
						for _, msg := range errors {
							if msg != "" {
								hasErrors = true
								break
							}
						}

						if hasErrors {
							payload, _ := json.Marshal(map[string]any{
								"errors": errors,
							})

							if err := sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							}
							return
						}

						// Create room
						_, err := roomService.CreateRoom(db, *store)
						if err != nil {
							payload, _ := json.Marshal(map[string]any{
								"errors": map[string]string{
									"room_number":          "",
									"max_concurrent_games": "",
									"error":                err.Error(),
								},
							})

							if err := sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							}
							return
						}

						if err := liveManager.Broadcast(r.Context(), live.BucketRooms); err != nil {
							logger.Error("Failed to broadcast room creation", "error", err.Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
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
							"delete_url":   fmt.Sprintf("/admin/rooms/api/delete/%d", room.ID),

							"id":                   room.ID,
							"name":                 room.Name,
							"room_number":          room.RoomNumber,
							"floor":                room.Floor,
							"max_concurrent_games": room.MaxConcurrentGames,
							"notes":                room.Notes,
							"is_disabled":          room.IsDisabled,
							"errors": map[string]string{
								"room_number":          "",
								"max_concurrent_games": "",
								"error":                "",
							},
						})

						if err := sse.PatchSignals(payload); err != nil {
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
						errors := map[string]string{
							"room_number":          "",
							"max_concurrent_games": "",
							"error":                "",
						}

						if !strings.HasPrefix(store.RoomNumber, fmt.Sprintf("%d", store.Floor)) {
							errors["room_number"] = "Romnummer må starte med etasje som første tall"
						}

						if strings.TrimSpace(store.RoomNumber) == "" {
							errors["room_number"] = "Romnummer er påkrevd"
						}

						if store.MaxConcurrentGames < 1 {
							errors["max_concurrent_games"] = "Maks samtidige spill må være minst 1"
						}

						parsedID, err := strconv.ParseInt(roomID, 10, 0)
						if err != nil {
							errors["error"] = err.Error()
						}

						hasErrors := false
						for _, msg := range errors {
							if msg != "" {
								hasErrors = true
								break
							}
						}

						if hasErrors {
							payload, _ := json.Marshal(map[string]any{
								"errors": errors,
							})

							if err := sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							}
							return
						}

						// Update room
						store.ID = int(parsedID)
						_, err = roomService.UpdateRoom(db, *store)
						if err != nil {
							payload, _ := json.Marshal(map[string]string{
								"error": err.Error(),
							})

							if err := sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)

							}
							return
						}

						if err := liveManager.Broadcast(r.Context(), live.BucketRooms); err != nil {
							logger.Error("Failed to broadcast room update", "error", err.Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}

						// Ridirect on success
						_ = sse.Redirect("/admin/rooms")
					})
				})

				roomApiRouter.Route("/delete/{id}", func(deleteRoomRoute chi.Router) {
					deleteRoomRoute.Post("/", func(w http.ResponseWriter, r *http.Request) {
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

							if err := sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							}
							return
						}

						// delete room
						err = roomService.DeleteRoom(db, int(parsedID))
						if err != nil {
							payload, _ := json.Marshal(map[string]string{
								"error": err.Error(),
							})

							if err := sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)

							}
							return
						}

						if err := liveManager.Broadcast(r.Context(), live.BucketRooms); err != nil {
							logger.Error("Failed to broadcast room deletion", "error", err.Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}

						// Ridirect on success
						_ = sse.Redirect("/admin/rooms")
					})
				})
			})

			rooms.RoomsLayoutRoute(roomRouter, db, baseLogger, eventImageDir)
		})
	})

	return nil
}
