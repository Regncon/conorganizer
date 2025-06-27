package myevents

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/pages/index"
	"github.com/Regncon/conorganizer/pages/myprofile/myevents/formsubmission"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/delaneyj/toolbelt"
	"github.com/delaneyj/toolbelt/embeddednats"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nats-io/nats.go/jetstream"
	datastar "github.com/starfederation/datastar/sdk/go"
)

func SetupMyEventsRoute(router chi.Router, store sessions.Store, ns *embeddednats.Server, db *sql.DB, logger *slog.Logger) error {
	kv, kvErr := SetupNats(ns)
	if kvErr != nil {
		return fmt.Errorf("error setting up nats: %w", kvErr)
	}

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

	router.Route("/my-events", func(myeventsRouter chi.Router) {
		myeventsLayoutRoute(myeventsRouter)
		myeventsRouter.Route("/api", func(apiRouter chi.Router) {
			apiRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
				sessionID, mvc, err := mvcSession(w, r)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to get session id: %v", err), http.StatusInternalServerError)
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

						c := myEventsPage(userctx.GetUserRequestInfo(r.Context()).Id, db, logger)
						if err := sse.MergeFragmentTempl(c); err != nil {
							sse.ConsoleError(err)
							return
						}
					}
				}
			})

			apiRouter.Route("/new", func(newRouter chi.Router) {
				newRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
					sessionID, mvc, err := mvcSession(w, r)
					if err != nil {
						http.Error(w, fmt.Sprintf("failed to get session id: %v", err), http.StatusInternalServerError)
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

							c := formsubmission.NewEventFormPage()
							if err := sse.MergeFragmentTempl(c); err != nil {
								sse.ConsoleError(err)
								return
							}
						}
					}
				})

				formsubmission.SetupExampleInlineValidation(db, newRouter, logger)
			})

			apiRouter.Post("/create", func(w http.ResponseWriter, r *http.Request) {
				createNewEventFormSubmission(db, logger, w, r)
			})

		})

		myeventsRouter.Route("/new", func(newRouter chi.Router) {
			formsubmission.NewEventLayoutRoute(newRouter, db)
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

func createNewEventFormSubmission(db *sql.DB, logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	userInfo := userctx.GetUserRequestInfo(r.Context())
	if userInfo.Id == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Creating new event form submission for user: %s\n, and email: %s\n", userInfo.Id, userInfo.Email)
	logger.Info("Creating new event form submission")
	sqlStatement := `
	INSERT INTO events (
		host, email, status, title, description, host_name, phone_number, max_players,
		child_friendly, adults_only, beginner_friendly, experienced_only,
		can_be_run_in_english, long_running, short_running
	) VALUES (
		$1, $2, $3, '', '', '', 0, 0, false, false, false, false, false, false, false
	) RETURNING id`
	var eventID int64
	err := db.QueryRow(sqlStatement, userInfo.Id, userInfo.Email, EventStatusDraft).Scan(&eventID)
	if err != nil {
		logger.Error("Failed to create new event form submission", "error", err)
		return
	}
	logger.Info("New event form submission created", "eventID", eventID)
	http.Redirect(w, r, fmt.Sprintf("/my-events/%d", eventID), http.StatusSeeOther)
}
