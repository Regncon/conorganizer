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
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     authctx.RefreshCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
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

func insertUser(db *sql.DB, userID, email string, isAdmin bool, logger *slog.Logger) {
	_, err := db.Exec("INSERT INTO users (user_id, email, is_admin) VALUES (?, ?, ?)", userID, email, isAdmin)
	if err != nil {
		logger.Error(fmt.Errorf("failed to insert new user %q: %w", userID, err).Error())
		return
	}
	logger.Info("Inserted new user", "user_id", userID, "is_admin", isAdmin)
}

func updateUserAdmin(db *sql.DB, userID string, isAdmin bool, logger *slog.Logger) {
	var currentIsAdmin bool
	err := db.QueryRow("SELECT is_admin FROM users WHERE user_id = ?", userID).Scan(&currentIsAdmin)
	if err == sql.ErrNoRows {
		logger.Error("User not found for admin update", "user_id", userID)
		return
	}
	if err != nil {
		logger.Error(fmt.Errorf("failed to fetch current is_admin for user %q: %w", userID, err).Error())
		return
	}
	if currentIsAdmin == isAdmin {
		logger.Info("No change to user admin status", "user_id", userID, "is_admin", isAdmin)
		return
	}
	_, updateErr := db.Exec("UPDATE users SET is_admin = ? WHERE user_id = ?", isAdmin, userID)
	if updateErr != nil {
		logger.Error(fmt.Errorf("failed to update user %q: %w", userID, updateErr).Error())
		return
	}
	logger.Info("Updated user admin status", "user_id", userID, "is_admin", isAdmin)
}
