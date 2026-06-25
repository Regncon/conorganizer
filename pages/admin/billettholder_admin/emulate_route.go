package billettholderadmin

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Regncon/conorganizer/service/emulatectx"
)

// EmulatePlayerRoute mounts POST /admin/emulate/{billettholderID}/, which starts
// read-only "view as player" emulation for the user account linked to the given
// billettholder. The route is admin-only (the router is expected to carry
// RequireAdmin). On success it sets the emulate_target cookie to the linked
// user's id and redirects to the player home page.
func EmulatePlayerRoute(router chi.Router, db *sql.DB, logger *slog.Logger) {
	logger = logger.With("component", "billettholder_admin")
	router.Post("/admin/emulate/{billettholderID}/", func(w http.ResponseWriter, r *http.Request) {
		billettholderID := chi.URLParam(r, "billettholderID")
		if _, err := strconv.Atoi(billettholderID); err != nil {
			http.Error(w, "Ugyldig billettholder-id", http.StatusBadRequest)
			return
		}

		userID, err := linkedNonAdminUserID(db, billettholderID)
		if err == sql.ErrNoRows {
			http.Error(w, "Denne billettholderen har ingen brukerkonto å se som.", http.StatusNotFound)
			return
		}
		if err != nil {
			logger.Error("failed to resolve linked user for emulation", "billettholder_id", billettholderID, "error", err)
			http.Error(w, "Noe gikk galt, prøv igjen", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     emulatectx.CookieName,
			Value:    strconv.Itoa(userID),
			Path:     "/",
			HttpOnly: true,
			Secure:   r.TLS != nil,
			SameSite: http.SameSiteLaxMode,
		})
		logger.Info("started player emulation", "billettholder_id", billettholderID, "user_id", userID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}

// linkedNonAdminUserID returns the id of the non-admin user account linked to the
// billettholder. Emulation targets a non-admin player; if several are linked the
// lowest user id wins. Returns sql.ErrNoRows when no non-admin account is linked.
func linkedNonAdminUserID(db *sql.DB, billettholderID string) (int, error) {
	const query = `
		SELECT u.id
		FROM relation_billettholdere_users r
		JOIN users u ON u.id = r.user_id
		WHERE r.billettholder_id = ? AND u.is_admin = 0
		ORDER BY u.id ASC
		LIMIT 1
	`
	var userID int
	if err := db.QueryRow(query, billettholderID).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil
}
