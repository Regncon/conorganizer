package profilepage

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Regncon/conorganizer/components/formsubmission"
	eventimgupload "github.com/Regncon/conorganizer/components/formsubmission/event_img_upload"
	profilecomponent "github.com/Regncon/conorganizer/components/profile"
	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/pages/profile/newevent"
	profileticketspage "github.com/Regncon/conorganizer/pages/profile/tickets"
	billettholderService "github.com/Regncon/conorganizer/service/billettholder"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

const createEventFailureMessage = "Vi klarte ikke å opprette arrangementet akkurat nå. Prøv igjen om litt."

func SetupProfileRoute(router chi.Router, liveManager *live.Manager, db *sql.DB, eventImageDir *string, logger *slog.Logger) error {
	var profileTicketsErr error
	router.Route("/profile", func(profileRouter chi.Router) {
		profileRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			requestLogger := logger.With("component", "profile")
			ctx := r.Context()
			user := userctx.GetUserRequestInfo(ctx)
			events := GetEventsByExternalID(user.Id, db, requestLogger)
			billettholdere, err := billettholderService.GetBillettholdere(user.Id, db)
			if err != nil {
				requestLogger.Error(err.Error(), "user_id", user.Id)
			}

			selectedBillettholderID := selectedBillettholderIDFromRequest(r, user, billettholdere, requestLogger)
			validBillettholderIDs := billettholderIDs(billettholdere)

			tickets := make([]profilecomponent.TicketHolder, 0, len(billettholdere))
			for _, billettholder := range billettholdere {
				email := ""
				if len(billettholder.Emails) > 0 {
					email = billettholder.Emails[0].Email
				}
				tickets = append(tickets, profilecomponent.TicketHolder{
					Name:   strings.TrimSpace(billettholder.FirstName + " " + billettholder.LastName),
					Ticket: billettholder.TicketType,
					Email:  email,
				})
			}

			if err := layouts.Base(
				"Min profil side",
				user,
				ProfilePage(user, events, tickets, selectedBillettholderID, validBillettholderIDs, db, requestLogger, eventImageDir),
			).Render(ctx, w); err != nil {
				requestLogger.Error(fmt.Errorf("failed to render profile page: %w", err).Error(), "user_id", user.Id)
			}
		})
		profileRouter.Get("/descope-profile", func(w http.ResponseWriter, r *http.Request) {
			requestLogger := logger.With("component", "profile")
			ctx := r.Context()
			user := userctx.GetUserRequestInfo(ctx)

			if err := layouts.Base(
				"Min profil side",
				user,
				DescopeProfilePage(),
			).Render(ctx, w); err != nil {
				requestLogger.Error(fmt.Errorf("failed to render profile page: %w", err).Error(), "user_id", user.Id)
			}
		})

		profileRouter.Route("/api", func(profileApiRouter chi.Router) {
			profileApiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				requestLogger := logger.With("component", "profile")
				liveManager.Stream(w, r, live.Page{
					Buckets: []live.Bucket{live.BucketEvents, live.BucketInterests, live.BucketBillettholders},
					Render: func(ctx context.Context, r *http.Request) templ.Component {
						user := userctx.GetUserRequestInfo(ctx)
						events := GetEventsByExternalID(user.Id, db, requestLogger)
						billettholdere, err := billettholderService.GetBillettholdere(user.Id, db)
						if err != nil {
							requestLogger.Error(err.Error(), "user_id", user.Id)
						}
						selectedBillettholderID := selectedBillettholderIDFromRequest(r, user, billettholdere, requestLogger)
						return ProfileMainColumn(user, events, selectedBillettholderID, db, requestLogger, eventImageDir)
					},
				})
			})

			profileApiRouter.Post("/create", func(w http.ResponseWriter, r *http.Request) {
				createNewEventFormSubmission(db, liveManager, logger, w, r)
			})

			profileApiRouter.Route("/new", func(newApiRouter chi.Router) {
				newApiRouter.Route("/{id}", func(newApiIdRouter chi.Router) {
					newApiIdRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
						eventId := chi.URLParam(r, "id")
						if eventId == "" {
							http.Error(w, "Event ID is required. Got: "+eventId, http.StatusBadRequest)
							return
						}

						liveManager.Stream(w, r, live.Page{
							Buckets: []live.Bucket{live.BucketEvents, live.BucketRooms},
							Render: func(ctx context.Context, r *http.Request) templ.Component {
								userId := userctx.GetUserRequestInfo(r.Context()).Id
								return newevent.NewEventFormPage(eventId, userId, ctx, db, eventImageDir, logger)
							},
						})
					})
					if err := formsubmission.SetupExampleInlineValidation(db, newApiIdRouter, logger); err != nil {
						logger.Error(fmt.Errorf("failed to set up inline validation: %w", err).Error())
					}

					newApiIdRouter.Route("/event-in-pulje", func(putEventInPuljeRouter chi.Router) {
						formsubmission.UpdateEventInPulje(putEventInPuljeRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/is-published", func(putIsPublishedRouter chi.Router) {
						formsubmission.UpdateIsPublished(putIsPublishedRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/assign-room", func(putAssignedRoomRouter chi.Router) {
						formsubmission.UpdateRoomInPulje(putAssignedRoomRouter, db, liveManager, logger)
					})

					newApiIdRouter.Route("/status", func(putStatusRouter chi.Router) {
						formsubmission.UpdateStatus(putStatusRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/name", func(putNameRouter chi.Router) {
						formsubmission.UpdateName(putNameRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/email", func(putEmailRouter chi.Router) {
						formsubmission.UpdateEmail(putEmailRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/phone", func(putPhoneRouter chi.Router) {
						formsubmission.UpdatePhone(putPhoneRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/title", func(putTitleRouter chi.Router) {
						formsubmission.UpdateTitle(putTitleRouter, db, liveManager, logger)
					})

					newApiIdRouter.Route("/intro", func(putIntroRouter chi.Router) {
						formsubmission.UpdateIntro(putIntroRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/system", func(putSystemRouter chi.Router) {
						formsubmission.UpdateSystem(putSystemRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/type", func(putTypeRouter chi.Router) {
						formsubmission.UpdateType(putTypeRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/description", func(putDescriptionRouter chi.Router) {
						formsubmission.UpdateDescription(putDescriptionRouter, db, liveManager, logger)
					})

					// should be age-group
					newApiIdRouter.Route("/ageGroup", func(putAgeGroupRouter chi.Router) {
						formsubmission.UpdateAgeGroup(putAgeGroupRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/runtime", func(putRuntimeRouter chi.Router) {
						formsubmission.UpdateRuntime(putRuntimeRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/beginner-friendly", func(putBeginnerFriendlyRouter chi.Router) {
						formsubmission.UpdateBeginnerFriendly(putBeginnerFriendlyRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/can-be-run-in-english", func(putCanBeRunInEnglishRouter chi.Router) {
						formsubmission.UpdateCanBeRunInEnglish(putCanBeRunInEnglishRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/max-players", func(putMaxPlayersRouter chi.Router) {
						formsubmission.UpdateMaxPlayers(putMaxPlayersRouter, db, liveManager, logger)
					})
					newApiIdRouter.Route("/notes", func(putNotesRouter chi.Router) {
						formsubmission.UpdateNotes(putNotesRouter, db, liveManager, logger)
					})

					newApiIdRouter.Route("/upload", func(uploadRouter chi.Router) {
						eventimgupload.EventImageFormSubmission(uploadRouter, db, eventImageDir, logger)
					})
					newApiIdRouter.Route("/upload-cropped", func(uploadCroppedRouter chi.Router) {
						eventimgupload.EventImageCroppedSubmission(uploadCroppedRouter, db, eventImageDir, logger)
					})
					newApiIdRouter.Route("/submit", func(newApiIdRouter chi.Router) {
						formsubmission.SubmitFormRoute(newApiIdRouter, db, liveManager, logger)
					})
				})
			})
		})

		profileTicketsErr = profileticketspage.ProfileTicketsRoute(profileRouter, liveManager, db, logger)

		profileRouter.Route("/new", func(newRouter chi.Router) {
			newRouter.Route("/{id}", func(newIdRoute chi.Router) {
				newevent.NewEventLayoutRoute(newIdRoute, db, eventImageDir, logger)

				newIdRoute.Route("/image", func(imageRouter chi.Router) {
					eventimgupload.EventImageRoute(imageRouter, db, eventImageDir, logger)
				})
			})
		})

	})
	if profileTicketsErr != nil {
		return fmt.Errorf("error setting up profile tickets route: %w", profileTicketsErr)
	}

	return nil
}

func GetEventsByExternalID(externalID string, db *sql.DB, logger *slog.Logger) []models.EventCardModel {
	logger = logger.With("component", "profile")
	var events []models.EventCardModel

	// Get events where event created id is the same as user
	userID, err := userctx.GetUserIDFromExternalID(externalID, db, logger)
	if err != nil {
		logger.Error(fmt.Errorf("failed to resolve user_id for external_id %q: %w", externalID, err).Error())
		return events
	}

	// Query for events created by user
	eventsQuery := "SELECT id, title, intro, status, system, host_name, beginner_friendly, event_type, age_group, event_runtime, can_be_run_in_english FROM events WHERE user_id = ?"
	rows, eventsQueryErr := db.Query(eventsQuery, userID)
	if eventsQueryErr != nil {
		logger.Error(fmt.Errorf("failed to query events for external_id %q: %w", externalID, eventsQueryErr).Error())
		return events
	}
	defer rows.Close()

	// Validate database query return
	for rows.Next() {
		var event models.EventCardModel
		if scanErr := rows.Scan(&event.Id, &event.Title, &event.Intro, &event.Status, &event.System, &event.HostName, &event.BeginnerFriendly, &event.EventType, &event.AgeGroup, &event.Runtime, &event.CanBeRunInEnglish); scanErr != nil {
			logger.Error(fmt.Errorf("failed to scan event row for external_id %q: %w", externalID, scanErr).Error())
			return events
		}
		events = append(events, event)
	}

	return events
}

func selectedBillettholderIDFromRequest(r *http.Request, user requestctx.UserRequestInfo, billettholdere []models.Billettholder, logger *slog.Logger) int {
	rawID := strings.TrimSpace(r.URL.Query().Get("b_id"))
	if rawID == "" {
		return defaultSelectedBillettholderID(user, billettholdere, logger)
	}

	billettholderID, err := strconv.Atoi(rawID)
	if err != nil {
		logger.Debug("Ignoring invalid profile billettholder id", "user_id", user.Id, "b_id", rawID)
		return defaultSelectedBillettholderID(user, billettholdere, logger)
	}

	if hasBillettholderID(billettholdere, billettholderID) {
		return billettholderID
	}

	logger.Debug("Ignoring profile billettholder id without user relation", "user_id", user.Id, "b_id", billettholderID)
	return defaultSelectedBillettholderID(user, billettholdere, logger)
}

func hasBillettholderID(billettholdere []models.Billettholder, billettholderID int) bool {
	for _, billettholder := range billettholdere {
		if billettholder.ID == billettholderID {
			return true
		}
	}
	return false
}

func billettholderIDs(billettholdere []models.Billettholder) []int {
	ids := make([]int, 0, len(billettholdere))
	for _, billettholder := range billettholdere {
		ids = append(ids, billettholder.ID)
	}
	return ids
}

func defaultSelectedBillettholderID(user requestctx.UserRequestInfo, billettholdere []models.Billettholder, logger *slog.Logger) int {
	for _, billettholder := range billettholdere {
		for _, email := range billettholder.Emails {
			if strings.EqualFold(email.Email, user.Email) {
				return billettholder.ID
			}
		}
	}

	if len(billettholdere) == 0 {
		logger.Debug("No billettholder available for profile selection", "user_id", user.Id)
		return 0
	}
	return billettholdere[0].ID
}

func createNewEventFormSubmission(db *sql.DB, liveManager *live.Manager, logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	logger = logger.With("component", "my_events")
	logger.Info("Creating new event form submission")
	userInfo := userctx.GetUserRequestInfo(r.Context())
	createNewEventFormSubmissionForUser(db, liveManager, logger, w, r, userInfo)
}

func createNewEventFormSubmissionForUser(db *sql.DB, liveManager *live.Manager, logger *slog.Logger, w http.ResponseWriter, r *http.Request, userInfo requestctx.UserRequestInfo) {
	if userInfo.Id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, insertError := userctx.GetUserIDFromExternalID(userInfo.Id, db, logger)
	if insertError != nil {
		logger.Error(fmt.Errorf("failed to resolve user_id for external_id %q: %w", userInfo.Id, insertError).Error())
		http.Error(w, createEventFailureMessage, http.StatusInternalServerError)
		return
	}

	logger.Info("Found user info for event creation", "external_id", userInfo.Id, "user_id", userID)
	logger.Info("Inserting new event form submission")

	// TODO: Use database relations to get foreign keys, event_type etc.
	query := `
	INSERT INTO events (
		user_id, created_by_id, updated_by_id, email, status, title, intro, description, host_name, phone_number, max_players,
		event_type, beginner_friendly,
		can_be_run_in_english
	) VALUES (
		$1, $1, $1, $2, $3, '', '', '', '', '', 6, $4, false, false
	) RETURNING id`

	var eventId string
	insertError = db.QueryRow(query, userID, userInfo.Email, models.EventStatusDraft, models.EventTypeRoleplay).Scan(&eventId)
	if insertError != nil {
		logger.Error(fmt.Errorf("failed to create new event form submission for user %q: %w", userInfo.Id, insertError).Error())
		http.Error(w, createEventFailureMessage, http.StatusInternalServerError)
		return
	}

	logger.Info("New event form submission created", "event_id", eventId, "user_id", userInfo.Id)
	if err := liveManager.Broadcast(r.Context(), live.BucketEvents); err != nil {
		logger.Error(fmt.Errorf("failed to broadcast new event creation: %w", err).Error(), "event_id", eventId, "user_id", userInfo.Id)
	}
	http.Redirect(w, r, fmt.Sprintf("/profile/new/%s", eventId), http.StatusSeeOther)
}
