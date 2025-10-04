package checkIn

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

// AssociateUserWithBillettholder uses userID string from
func AssociateUserWithBillettholder(userID string, db *sql.DB, logger *slog.Logger) error {
	logger.Info("Associating userID with billettholder", "userID", userID)

	// Get user
	var user models.User
	err := db.QueryRow(`
        SELECT user_id, email FROM users WHERE user_id = ?;
    `, userID).Scan(&user.UserID, &user.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Get associated billettholdere
	var billettholderIDs []models.BillettholderEmail
	rows, err := db.Query(`
        SELECT billettholder_id FROM billettholder_emails WHERE email = ?
    `, user.Email)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var billettholderID models.BillettholderEmail
		err := rows.Scan(&billettholderID.BillettholderID)
		if err != nil {
			return err
		}
		billettholderIDs = append(billettholderIDs, billettholderID)
	}

	// Insert into billettholder_users the new data
	var lines []string

	for _, billettholderID := range billettholderIDs {
		lines = append(lines, fmt.Sprintf(`(%d, '%s')`, billettholderID.BillettholderID, user.Email))
	}

	lines = append(lines, "(1, 'lars@regncon')")

	var baseQuery = fmt.Sprintf(`
        INSERT INTO billettholdere_users (
            billettholder_id, user_id
        ) VALUES %s
    `, strings.Join(lines, ", "))

	_, err = db.Exec(baseQuery)
	if err != nil {
		return err
	}

	return nil
}
