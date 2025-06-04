package auth

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/pages/auth/redirect"
	"github.com/Regncon/conorganizer/service"
	"github.com/go-chi/chi/v5"
)

func SetupAuthRoute(router chi.Router, db *sql.DB, logger *slog.Logger) error {
	router.Route("/auth", func(authRouter chi.Router) {
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			loginForm().Render(r.Context(), w)
		})

		authRouter.Group(func(protectedRoute chi.Router) {
			protectedRoute.Use(service.AuthMiddleware(logger))

			protectedRoute.Get("/test", func(w http.ResponseWriter, r *http.Request) {
				userToken, err := service.GetUserTokenFromContext(r.Context())
				if err != nil {
					http.Error(w, err.Error(), http.StatusUnauthorized)
					return
				}

				isAdmin := service.GetAdminFromUserToken(r.Context())
				layouts.Base("Is logged in test").Render(r.Context(), w)
				w.Write(fmt.Appendf(nil, "Test successful! Authenticated as: %v, and is admin: %v", userToken.Claims["email"], isAdmin))
			})

			protectedRoute.Get("/post-login", func(w http.ResponseWriter, r *http.Request) {
				isAdmin := service.GetAdminFromUserToken(r.Context())
				userToken, userTokenErr := service.GetUserTokenFromContext(r.Context())
				if userTokenErr != nil {
					logger.Error("Failed to get user token from context", "error", userTokenErr)
					http.Redirect(w, r, "/auth", http.StatusSeeOther)
					return
				}

				email, emailOk := userToken.Claims["email"].(string)
				userID, _ := service.GetUserIDFromToken(r.Context())

				if emailOk && email != "" && userID != "" {
					exists, err := userExistsByEmail(db, email)
					if err != nil {
						logger.Error("Failed to check if user exists", "error", err, "email", email)
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
				Name:     service.SessionCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			http.SetCookie(w, &http.Cookie{
				Name:     service.RefreshCookieName,
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})

			redirectUrl := "/"
			redirect.Redirect(redirectUrl, "Logging you out").Render(r.Context(), w)
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
		return false, err
	}
	return true, nil
}

func insertUser(db *sql.DB, userID, email string, isAdmin bool, logger *slog.Logger) {
	_, err := db.Exec("INSERT INTO users (user_id, email, is_admin) VALUES (?, ?, ?)", userID, email, isAdmin)
	if err != nil {
		logger.Error("Failed to insert new user", "error", err, "email", email)
		return
	}
	logger.Info("Inserted new user", "email", email, "user_id", userID, "is_admin", isAdmin)
}

func updateUserAdmin(db *sql.DB, userID string, isAdmin bool, logger *slog.Logger) {
	var currentIsAdmin bool
	err := db.QueryRow("SELECT is_admin FROM users WHERE user_id = ?", userID).Scan(&currentIsAdmin)
	if err == sql.ErrNoRows {
		logger.Error("User not found for admin update", "user_id", userID)
		return
	}
	if err != nil {
		logger.Error("Failed to fetch current is_admin", "error", err, "user_id", userID)
		return
	}
	if currentIsAdmin == isAdmin {
		logger.Info("No change to user admin status", "user_id", userID, "is_admin", isAdmin)
		return
	}
	_, updateErr := db.Exec("UPDATE users SET is_admin = ? WHERE user_id = ?", isAdmin, userID)
	if updateErr != nil {
		logger.Error("Failed to update user", "error", updateErr, "user_id", userID)
		return
	}
	logger.Info("Updated user admin status", "user_id", userID, "is_admin", isAdmin)
}
