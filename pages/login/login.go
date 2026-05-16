package login

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/components/redirect"
	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
)

func SetupAuthRoute(router chi.Router, db *sql.DB, logger *slog.Logger) error {
	baseLogger := logger
	logger = logger.With("component", "auth")
	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			if err := layouts.Base(
				"Innlogging til Regncon 2025!",
				userctx.GetUserRequestInfo(ctx),
				loginForm(),
			).Render(ctx, w); err != nil {
				logger.Error(fmt.Errorf("failed to render login page: %w", err).Error())
			}
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
					exists, err := userExistsByEmail(db, email)
					if err != nil {
						logger.Error(fmt.Errorf("failed to check if user %q exists: %w", userID, err).Error())
						http.Redirect(w, r, "/auth", http.StatusSeeOther)
						return
					}
					if !exists {
						insertUser(db, userID, email, isAdmin, logger)
					}
					updateUserAdmin(db, userID, isAdmin, logger)
				}
				http.Redirect(w, r, "/", http.StatusSeeOther)
			})

		})

		authRouter.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, &http.Cookie{
				Name:     authctx.SessionCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     authctx.RefreshCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteLaxMode,
			})

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
