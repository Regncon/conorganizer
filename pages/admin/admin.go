package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Regncon/conorganizer/components/formsubmission"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/pages/admin/approval"
	edit_form "github.com/Regncon/conorganizer/pages/admin/approval/editForm"
	"github.com/Regncon/conorganizer/pages/admin/rooms"
	"github.com/Regncon/conorganizer/service/live"
	roomService "github.com/Regncon/conorganizer/service/rooms"
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
			liveManager.Stream(w, r, live.Page{
				Buckets: []live.Bucket{live.BucketEvents},
				Render: func(ctx context.Context, r *http.Request) templ.Component {
					return adminPage(db)
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

		adminRouter.Route("/rooms", func(roomsRouter chi.Router) {
			roomsRouter.Route("/api", func(roomsApiRouter chi.Router) {
				roomsApiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
					liveManager.Stream(w, r, live.Page{
						Buckets: []live.Bucket{
							live.BucketRooms,
						},
						Render: func(ctx context.Context, r *http.Request) templ.Component {
							return rooms.RoomsPageContent(db, logger)
						},
					})
				})

				roomsApiRouter.Route("/{id}", func(roomApiRouter chi.Router) {
					// This route is used for getting form data when creating or updating rooms
					roomApiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
						// Validate id param from URL
						idQuery := chi.URLParam(r, "id")
						if idQuery == "" {
							http.Error(w, "Room ID is required. Got: "+idQuery, http.StatusBadRequest)
							return
						}
						roomID, err := strconv.ParseInt(idQuery, 10, 0)
						if err != nil {
							http.Error(w, "Unable to parse roomID, error: "+err.Error(), http.StatusBadRequest)
							return
						}

						// Handle signals depending on create or update room
						store := models.RoomFormSignals{}
						if roomID == 0 {
							store.FormTitle = "Legg til et nytt rom"
							store.ButtonLabel = "Legg til"
							store.Mode = "create"
						} else {
							room, err := roomService.GetRoomByID(db, int(roomID))
							if err != nil {
								store.Errors.Error = "Unable to get room"
							}

							store.FormTitle = "Oppdaterer rom " + room.RoomNumber
							store.ButtonLabel = "Oppdater"
							store.Mode = "edit"

							// Get updated values from database
							store.ID = room.ID
							store.Name = room.Name
							store.RoomNumber = room.RoomNumber
							store.Floor = room.Floor
							store.MaxConcurrentGames = room.MaxConcurrentGames
							store.IsDisabled = room.IsDisabled
							store.Notes = room.Notes
						}

						// Patch signals
						sse := datastar.NewSSE(w, r)
						payload, err := json.Marshal(store)
						if err != nil {
							logger.Error("Failed to marshal signals", "error", err.Error())
							http.Error(w, "Failed to marshal signals", http.StatusInternalServerError)
						}

						if err := sse.PatchSignals(payload); err != nil {
							logger.Error("Failed to patch signals", "error", err.Error())
							http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							return
						}
					})

					roomApiRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
						// Validate id param from URL
						idQuery := chi.URLParam(r, "id")
						if idQuery == "" {
							http.Error(w, "Room ID is required. Got: "+idQuery, http.StatusBadRequest)
							return
						}
						roomID, err := strconv.ParseInt(idQuery, 10, 0)
						if err != nil {
							http.Error(w, "Unable to parse roomID, error: "+err.Error(), http.StatusBadRequest)
							return
						}

						// Read data-star post submission
						store := &models.RoomFormSignals{}
						if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
							return
						}
						room := models.Room{
							ID:                 int(roomID),
							Name:               store.Name,
							RoomNumber:         store.RoomNumber,
							Floor:              store.Floor,
							MaxConcurrentGames: store.MaxConcurrentGames,
							IsDisabled:         store.IsDisabled,
							Notes:              store.Notes,
						}

						// Set up sse signals for form errors
						sse := datastar.NewSSE(w, r)

						// Validate form input
						validationErrors := roomService.ValidateRooms(room)
						store.Errors = validationErrors

						if err := validationErrors.HasErrors(); err {
							payload, err := json.Marshal(store)
							if err != nil {
								logger.Error("Failed to marshal signals", "error", err.Error())
								http.Error(w, "Failed to marshal signals", http.StatusInternalServerError)
							}

							if err := sse.PatchSignals(payload); err != nil {
								logger.Error("Failed to patch signals", "error", err.Error())
								http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
							}

							// Stop signals and let user update form
							return
						}

						// Update database
						if room.ID == 0 {
							_, err := roomService.CreateRoom(db, room)
							if err != nil {
								store.Errors.Error = err.Error()
								payload, err := json.Marshal(store)
								if err != nil {
									logger.Error("Failed to marshal signals", "error", err.Error())
									http.Error(w, "Failed to marshal signals", http.StatusInternalServerError)
								}

								if err := sse.PatchSignals(payload); err != nil {
									logger.Error("Failed to patch signals", "error", err.Error())
									http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
								}

								// Stop signals and let user update form
								return
							}
						} else {
							_, err := roomService.UpdateRoom(db, room)
							if err != nil {
								store.Errors.Error = err.Error()
								payload, err := json.Marshal(store)
								if err != nil {
									logger.Error("Failed to marshal signals", "error", err.Error())
									http.Error(w, "Failed to marshal signals", http.StatusInternalServerError)
								}

								if err := sse.PatchSignals(payload); err != nil {
									logger.Error("Failed to patch signals", "error", err.Error())
									http.Error(w, "Failed to patch signals", http.StatusInternalServerError)
								}

								// Stop signals and let user update form
								return
							}
						}

						// Broadcast that data has been changed, triggering all clients to update
						if err := liveManager.Broadcast(r.Context(), live.BucketRooms); err != nil {
							logger.Error(fmt.Errorf("failed to broadcast update: %w", err).Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}

						// Close modal on success
						_ = sse.ExecuteScript(`document.getElementById('room-dialog').close()`)
					})

					roomApiRouter.Delete("/", func(w http.ResponseWriter, r *http.Request) {
						// Validate id param from URL
						idQuery := chi.URLParam(r, "id")
						if idQuery == "" {
							http.Error(w, "Room ID is required. Got: "+idQuery, http.StatusBadRequest)
							return
						}
						roomID, err := strconv.ParseInt(idQuery, 10, 0)
						if err != nil {
							http.Error(w, "Unable to parse roomID, error: "+err.Error(), http.StatusBadRequest)
							return
						}

						// Delete room
						err = roomService.DeleteRoom(db, int(roomID))
						if err != nil {
							http.Error(w, "Unable to deleto room with ID, error: "+err.Error(), http.StatusBadRequest)
							return
						}

						if err := liveManager.Broadcast(r.Context(), live.BucketRooms); err != nil {
							logger.Error(fmt.Errorf("failed to broadcast update: %w", err).Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}

						// Close modal on success
						sse := datastar.NewSSE(w, r)
						_ = sse.ExecuteScript(`document.getElementById('room-dialog').close()`)
					})
				})

				roomsApiRouter.Route("/assignment/{pulje}", func(roomsAssignmentRouter chi.Router) {
					roomsAssignmentRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
						puljeQuery := chi.URLParam(r, "pulje")
						puljeID, isPujeIDValid := models.ParsePulje(puljeQuery)
						if !isPujeIDValid {
							http.Error(w, "Expected a valid pulje ID, got: "+puljeQuery, http.StatusBadRequest)
							return
						}

						liveManager.Stream(w, r, live.Page{
							Buckets: []live.Bucket{
								live.BucketRooms,
							},
							Render: func(ctx context.Context, r *http.Request) templ.Component {
								return rooms.RoomsAssignmentPageContent(db, logger, puljeID, eventImageDir)
							},
						})
					})

					roomsAssignmentRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
						puljeQuery := chi.URLParam(r, "pulje")
						puljeID, isPujeIDValid := models.ParsePulje(puljeQuery)
						if !isPujeIDValid {
							http.Error(w, "Expected a valid pulje ID, got: "+puljeQuery, http.StatusBadRequest)
							return
						}

						// todo: handle in page event assignment by clicking add to room
						fmt.Println(puljeID)
					})
				})
			})

			rooms.RoomsLayoutRoute(roomsRouter, db, logger, eventImageDir)
		})
	})

	return nil
}
