package login

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/Regncon/conorganizer/components/redirect"
	"github.com/Regncon/conorganizer/layouts"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/checkIn"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/service/userctx"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type sessionRequest struct {
	SessionJWT string `json:"sessionJwt"`
	RefreshJWT string `json:"refreshJwt"`
}

// nesteQueryParam is the query parameter carrying where login should return
// the user to. "Neste" is Norwegian for "next".
const nesteQueryParam = "neste"

// safeReturnPath validates a "neste" return target from a query parameter.
// Only a safe, same-site relative path is accepted; anything else (an empty
// value, a scheme-relative "//host" path, a backslash-escaped path, a path
// containing control characters or whitespace, or an absolute URL) falls
// back to "/" so login can never be used to redirect a user off-site.
func safeReturnPath(raw string) string {
	// Browsers strip tab, CR and LF from URLs and treat "\" as "/", so
	// "/\t/evil.com" or "/\\evil.com" would become "//evil.com". Reject any
	// control character or whitespace outright, then check the path with
	// backslashes normalised to slashes.
	if strings.ContainsFunc(raw, func(r rune) bool { return unicode.IsControl(r) || unicode.IsSpace(r) }) {
		return "/"
	}
	normalised := strings.ReplaceAll(raw, `\`, "/")
	if !strings.HasPrefix(normalised, "/") || strings.HasPrefix(normalised, "//") {
		return "/"
	}
	parsed, err := url.Parse(normalised)
	if err != nil || parsed.Scheme != "" || parsed.Host != "" {
		return "/"
	}
	return raw
}

// authPathWithNeste builds the "/auth" path to retry login while keeping the
// return target, so a failed post-login sync does not lose where the user
// was headed.
func authPathWithNeste(neste string) string {
	if neste == "" || neste == "/" {
		return "/auth"
	}
	return "/auth?" + nesteQueryParam + "=" + url.QueryEscape(neste)
}

func SetupAuthRoute(publicRouter, authenticatedRouter chi.Router, liveManager *live.Manager, db *sql.DB, logger *slog.Logger, sessionValidator authctx.SessionValidator) error {
	logger = logger.With("component", "auth")
	publicRouter.Post("/auth/session", func(w http.ResponseWriter, r *http.Request) {
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
				logger.Warn("failed to validate login session",
					"session_error", sessionErr,
					"refresh_error", refreshErr,
					"request_id", middleware.GetReqID(r.Context()),
				)
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

	publicRouter.Get("/auth/logout", func(w http.ResponseWriter, r *http.Request) {
		authctx.ClearAuthCookies(w, r)
		requestctx.ClearBillettholderSelectionCookie(w)

		redirectUrl := "/"
		var ctx = r.Context()
		if err := layouts.Base("Logging you out",
			requestctx.UserRequestInfo{},
			db,
			logger,
			redirect.Redirect(redirectUrl),
		).Render(ctx, w); err != nil {
			logger.Error(fmt.Errorf("failed to render logout page: %w", err).Error())
		}
	})

	authenticatedRouter.Route("/auth", func(authRouter chi.Router) {
		authRouter.Get("/", func(w http.ResponseWriter, r *http.Request) {
			var ctx = r.Context()
			userToken, _ := authctx.GetUserTokenFromContext(r.Context())
			neste := safeReturnPath(r.URL.Query().Get(nesteQueryParam))

			if userToken != nil {
				if err := layouts.Base(
					"Velkommen tilbake til Regncon 2026!",
					userctx.GetUserRequestInfo(ctx),
					db,
					logger,
					alreadyLogedIn(neste),
				).Render(ctx, w); err != nil {
					logger.Error(fmt.Errorf("failed to render already loged in page: %w", err).Error())
				}
			} else {
				if err := layouts.Base(
					"Innlogging til Regncon 2026!",
					userctx.GetUserRequestInfo(ctx),
					db,
					logger,
					loginForm(neste),
				).Render(ctx, w); err != nil {
					logger.Error(fmt.Errorf("failed to render login page: %w", err).Error())
				}

			}

		})

		authRouter.Group(func(protectedRoute chi.Router) {
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
					db,
					logger,
					testComp,
				).Render(ctx, w); err != nil {
					logger.Error(fmt.Errorf("failed to render auth test page: %w", err).Error())
				}
			})

			protectedRoute.Get("/post-login", func(w http.ResponseWriter, r *http.Request) {
				isAdmin := authctx.GetAdminFromUserToken(r.Context())
				userToken, userTokenErr := authctx.GetUserTokenFromContext(r.Context())
				neste := safeReturnPath(r.URL.Query().Get(nesteQueryParam))
				if userTokenErr != nil {
					logger.Error(fmt.Errorf("failed to get user token from context: %w", userTokenErr).Error())
					http.Redirect(w, r, authPathWithNeste(neste), http.StatusSeeOther)
					return
				}

				email, emailOk := userToken.Claims["email"].(string)
				userID, _ := authctx.GetUserIDFromToken(r.Context())

				if emailOk && email != "" && userID != "" {
					if err := syncPostLoginUser(db, userID, email, isAdmin, logger); err != nil {
						logger.Error(fmt.Errorf("failed to sync post-login user %q: %w", userID, err).Error())
						http.Redirect(w, r, authPathWithNeste(neste), http.StatusSeeOther)
						return
					}
					http.Redirect(w, r, neste, http.StatusSeeOther)
					go syncTicketsAfterLogin(context.WithoutCancel(r.Context()), syncUserTicketsFromCheckIn, liveManager.Broadcast, db, userID, email, logger)
					return
				}
				http.Redirect(w, r, neste, http.StatusSeeOther)
			})

		})
	})

	return nil
}

func syncPostLoginUser(db *sql.DB, userID string, email string, isAdmin bool, logger *slog.Logger) error {
	exists, err := userExistsByExternalID(db, userID)
	if err != nil {
		return err
	}
	if !exists {
		insertUser(db, userID, email, isAdmin, logger)
	} else {
		updateUserEmail(db, userID, email, logger)
	}
	updateUserAdmin(db, userID, isAdmin, logger)
	return nil
}

// postLoginTicketSyncTimeout bounds the background "Hent billetter" run, so a
// slow CheckIn cannot keep it alive indefinitely.
const postLoginTicketSyncTimeout = time.Minute

// syncUserTicketsFromCheckIn is the "Hent billetter" flow run after login.
// Tests replace it so they do not reach CheckIn.
var syncUserTicketsFromCheckIn = checkIn.SyncUserTicketsFromCheckIn

// syncTicketsAfterLogin runs "Hent billetter" for a user who just logged in. It
// runs after the redirect so login never waits for CheckIn, and its errors are
// only logged. Like "Hent billetter", it broadcasts the billettholder bucket when
// something was created, so open menus and admin pages pick up the change.
func syncTicketsAfterLogin(
	ctx context.Context,
	syncUserTickets func(context.Context, string, string, *sql.DB, *slog.Logger) (checkIn.UserTicketImportResult, error),
	broadcast func(context.Context, ...live.Bucket) error,
	db *sql.DB,
	userID string,
	email string,
	logger *slog.Logger,
) {
	logger = logger.With("component", "post_login_ticket_sync", "user_id", userID)
	defer func() {
		if recovered := recover(); recovered != nil {
			logger.Error("Post-login ticket sync panicked", "panic", recovered)
		}
	}()

	ctx, cancel := context.WithTimeout(ctx, postLoginTicketSyncTimeout)
	defer cancel()

	result, err := syncUserTickets(ctx, userID, email, db, logger)
	if err != nil {
		logger.Error(fmt.Errorf("failed to sync tickets after login: %w", err).Error())
	}
	if result.CreatedBillettholders == 0 && result.CreatedUserAssociations == 0 {
		return
	}
	logger.Info("Synced tickets after login",
		"created_billettholdere", result.CreatedBillettholders,
		"created_user_associations", result.CreatedUserAssociations,
	)
	if err := broadcast(ctx, live.BucketBillettholders); err != nil {
		logger.Error(fmt.Errorf("failed to broadcast tickets synced after login: %w", err).Error())
	}
}

func normalizeToken(token string) string {
	token = strings.TrimSpace(token)
	if len(token) >= len("bearer ") && strings.EqualFold(token[:len("bearer ")], "bearer ") {
		return strings.TrimSpace(token[len("bearer "):])
	}
	return token
}

func userExistsByExternalID(db *sql.DB, externalID string) (bool, error) {
	var exists int
	err := db.QueryRow("SELECT 1 FROM users WHERE external_id = ?", externalID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to query user %q: %w", externalID, err)
	}
	return true, nil
}

// updateUserEmail keeps users.email in step with the login provider, since
// "Hent billetter" matches tickets against it.
func updateUserEmail(db *sql.DB, externalID, email string, logger *slog.Logger) {
	result, err := db.Exec("UPDATE users SET email = ? WHERE external_id = ? AND email <> ?", email, externalID, email)
	if err != nil {
		logger.Error(fmt.Errorf("failed to update user email: %w", err).Error(), "external_id", externalID)
		return
	}
	if rowsAffected, err := result.RowsAffected(); err == nil && rowsAffected > 0 {
		logger.Info("Updated user email", "external_id", externalID)
	}
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
