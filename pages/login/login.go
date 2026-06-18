package login

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Regncon/conorganizer/components/redirect"
	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

var newSessionValidator = authctx.NewSessionValidatorFromEnv

type sessionRequest struct {
	SessionJWT string `json:"sessionJwt"`
	RefreshJWT string `json:"refreshJwt"`
}

func SetupAuthRoute(router chi.Router, db *sql.DB, logger *slog.Logger) error {
	baseLogger := logger
	logger = logger.With("component", "auth")
	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			userToken, _ := authctx.GetUserTokenFromContext(r.Context())

			if userToken != nil {
				if err := layouts.Base(
					"Velkomen tilbake til Regncon 2026!",
					userctx.GetUserRequestInfo(ctx),
					alreadyLogedIn(),
				).Render(ctx, w); err != nil {
					logger.Error(fmt.Errorf("failed to render already loged in page: %w", err).Error())
				}
			} else {
				if err := layouts.Base(
					"Innlogging til Regncon 2026!",
					userctx.GetUserRequestInfo(ctx),
					loginForm(),
				).Render(ctx, w); err != nil {
					logger.Error(fmt.Errorf("failed to render login page: %w", err).Error())
				}

			}

		})

		authRouter.Post("/session", func(w http.ResponseWriter, r *http.Request) {
			request := sessionRequest{}
			decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
			decoder.DisallowUnknownFields()
			if err := decoder.Decode(&request); err != nil {
				http.Error(w, "invalid session request", http.StatusBadRequest)
				return
			}
			if err := decoder.Decode(&struct{}{}); err != io.EOF {
				http.Error(w, "invalid session request", http.StatusBadRequest)
				return
			}
			request.SessionJWT = normalizeToken(request.SessionJWT)
			request.RefreshJWT = normalizeToken(request.RefreshJWT)
			if request.SessionJWT == "" || request.RefreshJWT == "" {
				http.Error(w, "missing session tokens", http.StatusBadRequest)
				return
			}

			sessionValidator, err := newSessionValidator()
			if err != nil {
				logger.Error(fmt.Errorf("failed to create auth session validator: %w", err).Error())
				http.Error(w, "authentication unavailable", http.StatusInternalServerError)
				return
			}

			userOK, userToken, sessionErr := sessionValidator.ValidateSessionWithToken(
				r.Context(),
				request.SessionJWT,
			)
			if sessionErr == nil && (!userOK || userToken == nil) {
				sessionErr = fmt.Errorf("session token rejected")
			}
			if sessionErr != nil || !userOK || userToken == nil {
				refreshedOK, refreshedToken, refreshErr := sessionValidator.RefreshSessionWithToken(
					r.Context(),
					request.RefreshJWT,
				)
				if refreshErr == nil && (!refreshedOK || refreshedToken == nil) {
					refreshErr = fmt.Errorf("refresh token rejected")
				}
				if refreshErr != nil || !refreshedOK || refreshedToken == nil {
					if sessionErr != nil || refreshErr != nil {
						logger.Warn("failed to validate login session", "session_error", sessionErr, "refresh_error", refreshErr)
					} else {
						logger.Warn("login session was rejected")
					}
					http.Error(w, "invalid session", http.StatusUnauthorized)
					return
				}

				userToken = refreshedToken
			}

			sessionJWT := request.SessionJWT
			if userToken.JWT != "" {
				sessionJWT = userToken.JWT
			}

			authctx.SetAuthCookies(w, r, sessionJWT, request.RefreshJWT)
			w.WriteHeader(http.StatusNoContent)
		})

		authRouter.Group(func(protectedRoute chi.Router) {
			protectedRoute.Use(authctx.AuthMiddleware(baseLogger))

			protectedRoute.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				userToken, userTokenErr := authctx.GetUserTokenFromContext(r.Context())
				if userTokenErr != nil {
					http.Error(w, userTokenErr.Error(), http.StatusUnauthorized)
					return
				}

				isAdmin := authctx.GetAdminFromUserToken(r.Context())
				testComp := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
					_, err := io.WriteString(w, fmt.Sprintf("Test successful! Authenticated as: %v, and is admin: %v", userToken.Claims["email"], isAdmin))
					if err != nil {
						return fmt.Errorf("write auth test response: %w", err)
					}
					return nil
				})

				var ctx = r.Context()
				if err := layouts.Base(
					"Is logged in test",
					userctx.GetUserRequestInfo(ctx),
					testComp,
				).Render(ctx, w); err != nil {
					logger.Error(fmt.Errorf("failed to render auth test page: %w", err).Error())
				}
			})

			protectedRoute.Get("/post-login", func(w http.ResponseWriter, r *http.Request) {
				isAdmin := authctx.GetAdminFromUserToken(r.Context())
				userToken, userTokenErr := authctx.GetUserTokenFromContext(r.Context())
				if userTokenErr != nil {
					logger.Error(fmt.Errorf("failed to get user token from context: %w", userTokenErr).Error())
					http.Redirect(w, r, "/auth", http.StatusSeeOther)
					return
				}

				email, emailOk := userToken.Claims["email"].(string)
				userID, _ := authctx.GetUserIDFromToken(r.Context())

				if emailOk && email != "" && userID != "" {
					if err := syncPostLoginUser(db, userID, email, isAdmin, logger); err != nil {
						logger.Error(fmt.Errorf("failed to sync post-login user %q: %w", userID, err).Error())
						http.Redirect(w, r, "/auth", http.StatusSeeOther)
						return
					}
				}
				http.Redirect(w, r, "/", http.StatusSeeOther)
			})

		})

		authRouter.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			authctx.ClearAuthCookies(w, r)

			redirectUrl := "/"
			var ctx = r.Context()
			if err := layouts.Base("Logging you out",
				userctx.GetUserRequestInfo(ctx),
				redirect.Redirect(redirectUrl),
			).Render(ctx, w); err != nil {
				logger.Error(fmt.Errorf("failed to render logout page: %w", err).Error())
			}
		})
	})

	return nil
}

func syncPostLoginUser(db *sql.DB, userID string, email string, isAdmin bool, logger *slog.Logger) error {
	exists, err := userExistsByEmail(db, email)
	if err != nil {
		return err
	}
	if !exists {
		insertUser(db, userID, email, isAdmin, logger)
	}
	updateUserAdmin(db, userID, isAdmin, logger)
	return nil
}

func normalizeToken(token string) string {
	token = strings.TrimSpace(token)
	if len(token) >= len("bearer ") && strings.EqualFold(token[:len("bearer ")], "bearer ") {
		return strings.TrimSpace(token[len("bearer "):])
	}
	return token
}

func userExistsByEmail(db *sql.DB, email string) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT 1 FROM users WHERE email = ?", email).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to query user by email %q: %w", email, err)
	}
	return true, nil
}

func insertUser(db *sql.DB, externalID, email string, isAdmin bool, logger *slog.Logger) {
	_, err := db.Exec("INSERT INTO users (external_id, email, is_admin) VALUES (?, ?, ?)", externalID, email, isAdmin)
	if err != nil {
		logger.Error(fmt.Errorf("failed to insert new user %q: %w", externalID, err).Error())
		return
	}

	logger.Info("Inserted new user", "email", email, "external_id", externalID, "is_admin", isAdmin)
}

func updateUserAdmin(db *sql.DB, externalID string, isAdmin bool, logger *slog.Logger) {
	var currentIsAdmin bool
	err := db.QueryRow("SELECT is_admin FROM users WHERE external_id = ?", externalID).Scan(&currentIsAdmin)
	if err == sql.ErrNoRows {
		logger.Error("user not found for admin update", "user_id", externalID)
		return
	}
	if err != nil {
		logger.Error(fmt.Errorf("failed to fetch current is_admin: %w", err).Error(), "external_id", externalID)
		return
	}
	if currentIsAdmin == isAdmin {
		logger.Info("No change to user admin status", "external_id", externalID, "is_admin", isAdmin)
		return
	}
	_, updateErr := db.Exec("UPDATE users SET is_admin = ? WHERE external_id = ?", isAdmin, externalID)
	if updateErr != nil {
		logger.Error(fmt.Errorf("failed to update user: %w", updateErr).Error(), "external_id", externalID)
		return
	}
	logger.Info("Updated user admin status", "external_id", externalID, "is_admin", isAdmin)
}
