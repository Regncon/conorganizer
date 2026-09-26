package admin

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Regncon/conorganizer/components/errorfeedback"
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
		puljefordelingRoute(adminRouter, db, liveManager, baseLogger, eventImageDir)
		puljeoppsettRoute(adminRouter, db, liveManager, baseLogger, eventImageDir)
		programPublishingRoute(adminRouter, db, liveManager, logger)
		feedbackAdminRoute(adminRouter, db, baseLogger)
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
							logger.Error("User tried to fetch a room without ID", "url", r.URL)
							http.Error(w, fmt.Sprintf("Room ID is required, got: %s", idQuery), http.StatusBadRequest)
							return
						}
						roomID, err := strconv.ParseInt(idQuery, 10, 0)
						if err != nil {
							logger.Error("User tried to fetch a room with invalid ID", "url", r.URL, "error", err.Error())
							http.Error(w, fmt.Sprintf("Unable to parse given room ID, error: %s", err.Error()), http.StatusBadRequest)
							return
						}

						// Initiate new signals
						sse := datastar.NewSSE(w, r)

						// Clear error messages
						roomErrors := models.RoomFormErrors{}
						roomErrors.ResetErrors()
						feedback := errorfeedback.New(roomErrors.GetKeys()...)

						// Fill form bindings
						store := models.RoomFormSignals{}
						if roomID == 0 {
							store.FormTitle = "Legg til et nytt rom"
							store.ButtonLabel = "Legg til"
							store.Mode = "create"
						} else {
							store.ButtonLabel = "Oppdater"
							store.Mode = "edit"

							room, err := roomService.GetRoomByID(db, int(roomID))
							if err != nil {
								store.FormTitle = "Finner ikke rom"
								feedback.Set("error", err.Error())
							} else {
								store.FormTitle = "Oppdaterer rom " + room.RoomNumber

								// Get updated values from database
								store.ID = room.ID
								store.Name = room.Name
								store.RoomNumber = room.RoomNumber
								store.Floor = room.Floor
								store.Notes = room.Notes
							}
						}

						// Patch signals
						if err := feedback.Patch(sse); err != nil {
							logger.Error("Failed to marshal and patch feedback signals", "error", err.Error())
							http.Error(w, "Failed to marshal and patch feedback signals", http.StatusInternalServerError)
							return
						}
						if err := sse.MarshalAndPatchSignals(store); err != nil {
							logger.Error("Failed to marshal and patch signals", "error", err.Error())
							http.Error(w, "Failed to marshal and patch signals", http.StatusInternalServerError)
							return
						}
					})

					roomApiRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
						// Validate id param from URL
						idQuery := chi.URLParam(r, "id")
						if idQuery == "" {
							logger.Error("User tried to update room without ID", "url", r.URL)
							http.Error(w, fmt.Sprintf("Room ID is required, got: %s", idQuery), http.StatusBadRequest)
							return
						}
						roomID, err := strconv.ParseInt(idQuery, 10, 0)
						if err != nil {
							logger.Error("User tried to update room with invalid ID", "url", r.URL, "error", err.Error())
							http.Error(w, fmt.Sprintf("Unable to parse given room ID, error: %s", err.Error()), http.StatusBadRequest)
							return
						}

						// Read data post submission
						store := &models.RoomFormSignals{}
						if readSignalErr := datastar.ReadSignals(r, store); readSignalErr != nil {
							http.Error(w, readSignalErr.Error(), http.StatusBadRequest)
							return
						}
						room := models.Room{
							ID:         int(roomID),
							Name:       store.Name,
							RoomNumber: store.RoomNumber,
							Floor:      store.Floor,
							Notes:      store.Notes,
						}

						// Decide between create and update based on room ID
						var roomErrors models.RoomFormErrors
						if room.ID == 0 {
							_, roomErrors = roomService.CreateRoom(db, room)
						} else {
							_, roomErrors = roomService.UpdateRoom(db, room)
						}

						// Set up sse signals for response
						sse := datastar.NewSSE(w, r)

						feedbackKeys := roomErrors.GetKeys()
						feedback := errorfeedback.New(feedbackKeys...)

						if roomErrors.HasErrors() {
							for errorKey, errorMsg := range roomErrors {
								feedback.Set(string(errorKey), errorMsg)
							}

							if err := feedback.Patch(sse); err != nil {
								logger.Error("Failed to marshal and patch feedback signals", "error", err.Error())
								http.Error(w, "Failed to marshal and patch feedback signals", http.StatusInternalServerError)

							}
							return
						}

						// Close modal on success
						if err := sse.ExecuteScript(`document.getElementById('room-dialog').close()`); err != nil {
							logger.Error("Failed to execute sse script", "error", err.Error())
							http.Error(w, "Failed to execute sse script", http.StatusInternalServerError)
						}

						// Broadcast that data has been changed, triggering all clients to update
						if err := liveManager.Broadcast(r.Context(), live.BucketRooms); err != nil {
							logger.Error(fmt.Errorf("failed to broadcast update: %w", err).Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}

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
						if err := sse.ExecuteScript(`document.getElementById('room-dialog').close()`); err != nil {
							logger.Error("Failed to execute sse script", "error", err.Error())
							http.Error(w, "Failed to execute sse script", http.StatusInternalServerError)
						}
					})
				})

				roomsApiRouter.Route("/assignment/{pulje}", func(roomsAssignmentRouter chi.Router) {
					roomsAssignmentRouter.Delete("/{event}/{room}", func(w http.ResponseWriter, r *http.Request) {
						puljeID, valid := models.ParsePulje(chi.URLParam(r, "pulje"))
						roomID, err := strconv.ParseInt(chi.URLParam(r, "room"), 10, 64)
						if !valid || err != nil || roomID <= 0 {
							http.Error(w, "Ugyldig pulje eller rom.", http.StatusBadRequest)
							return
						}
						eventID := chi.URLParam(r, "event")
						logger := logger.With("component", "room_assignment", "event_id", eventID, "pulje_id", puljeID, "room_id", roomID)
						// Match the old room too: a stale card must not undo another admin's move.
						result, err := db.ExecContext(r.Context(), `
							UPDATE relation_event_puljer SET room_id = NULL
							WHERE event_id = ? AND pulje_id = ? AND room_id = ? AND is_in_pulje = 1
						`, eventID, puljeID, roomID)
						if err != nil {
							logger.Error(fmt.Errorf("failed to remove room assignment: %w", err).Error())
							http.Error(w, "Klarte ikke å fjerne romtildelingen.", http.StatusInternalServerError)
							return
						}
						changed, err := result.RowsAffected()
						if err != nil {
							logger.Error(fmt.Errorf("failed to verify room removal: %w", err).Error())
							http.Error(w, "Klarte ikke å bekrefte romendringen.", http.StatusInternalServerError)
							return
						}
						if changed == 0 {
							http.Error(w, "Romtildelingen er endret. Last siden på nytt.", http.StatusConflict)
							return
						}
						if err := liveManager.Broadcast(r.Context(), live.BucketRooms, live.BucketEvents); err != nil {
							logger.Error(fmt.Errorf("failed to broadcast room removal: %w", err).Error())
							http.Error(w, "Romtildelingen ble fjernet, men visningen kunne ikke oppdateres. Last siden på nytt.", http.StatusInternalServerError)
							return
						}
						sse := datastar.NewSSE(w, r)
						_ = sse.ExecuteScript(`window.dispatchEvent(new Event('room-saved'))`)
					})

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
								live.BucketEvents,
							},
							Render: func(ctx context.Context, r *http.Request) templ.Component {
								return rooms.RoomsAssignmentPageContent(db, logger, puljeID, eventImageDir)
							},
						})
					})

					roomsAssignmentRouter.Post("/{event}/{room}", func(w http.ResponseWriter, r *http.Request) {
						puljeQuery := chi.URLParam(r, "pulje")
						puljeID, isPujeIDValid := models.ParsePulje(puljeQuery)
						if !isPujeIDValid {
							http.Error(w, fmt.Sprintf("Expected a valid pulje ID: %v", puljeQuery), http.StatusBadRequest)
							return
						}

						eventQuery := chi.URLParam(r, "event")
						if eventQuery == "" {
							http.Error(w, "Event ID is required", http.StatusBadRequest)
							return
						}

						roomQuery := chi.URLParam(r, "room")
						if roomQuery == "" {
							http.Error(w, "Room ID is required", http.StatusBadRequest)
							return
						}
						roomID, err := strconv.ParseInt(roomQuery, 10, 0)
						if err != nil {
							http.Error(w, fmt.Sprintf("Unable to parse roomID: %v", err.Error()), http.StatusBadRequest)
							return
						}

						// Check that room exists
						_, err = roomService.GetRoomByID(db, int(roomID))
						if err != nil {
							http.Error(w, fmt.Sprintf("Room with id %d not found - error: %v", roomID, err.Error()), http.StatusBadRequest)
							return
						}

						// Only approved events belong in the room assignment flow.
						var eventStatus models.EventStatus
						err = db.QueryRowContext(r.Context(), `SELECT status FROM events WHERE id = ?`, eventQuery).Scan(&eventStatus)
						if err == sql.ErrNoRows {
							http.Error(w, "Arrangementet ble ikke funnet.", http.StatusConflict)
							return
						}
						if err != nil {
							http.Error(w, fmt.Sprintf("Unable to check event: %v", err), http.StatusInternalServerError)
							return
						}
						if eventStatus != models.EventStatusApproved && eventStatus != models.EventStatusAnnounced {
							http.Error(w, "Arrangementet er ikke godkjent.", http.StatusConflict)
							return
						}

						// Assign the event to this pulje even when it has no active
						// relation here yet.
						result, err := db.ExecContext(r.Context(), `
							INSERT INTO relation_event_puljer (event_id, pulje_id, is_in_pulje, room_id)
							VALUES (?, ?, 1, ?)
							ON CONFLICT(event_id, pulje_id) DO UPDATE SET
								is_in_pulje = 1,
								room_id = excluded.room_id
						`, eventQuery, puljeID, roomID)
						if err != nil {
							http.Error(w, fmt.Sprintf("Unable to assign room: %v", err.Error()), http.StatusBadRequest)
							return
						}
						rowsAffected, err := result.RowsAffected()
						if err != nil {
							http.Error(w, "Unable to verify room assignment", http.StatusInternalServerError)
							return
						}
						if rowsAffected == 0 {
							http.Error(w, "Arrangementet er ikke lenger i denne puljen. Last siden på nytt.", http.StatusConflict)
							return
						}

						// Stream update
						if err := liveManager.Broadcast(r.Context(), live.BucketRooms, live.BucketEvents); err != nil {
							logger.Error(fmt.Errorf("failed to broadcast update: %w", err).Error())
							http.Error(w, "Failed to broadcast update", http.StatusInternalServerError)
							return
						}

						// Close modal on success
						sse := datastar.NewSSE(w, r)
						_ = sse.ExecuteScript(`document.getElementById('assignment-dialog')?.close(); window.dispatchEvent(new Event('room-saved'))`)
					})
				})
			})

			rooms.RoomsLayoutRoute(roomsRouter, db, logger, eventImageDir)
		})
	})

	return nil
}
