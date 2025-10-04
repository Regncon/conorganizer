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
        SELECT id, email FROM users WHERE user_id = ?;
    `, userID).Scan(&user.ID, &user.Email)
	if err != nil {
		return errors.New("user not found")
	}

	// Get associated billettholdere
	var billettholdere []models.BillettholderEmail
	rows, err := db.Query(`
        SELECT id, billettholder_id, email, kind, inserted_time FROM billettholder_emails WHERE email = ?
    `, user.Email)
	if err != nil {
		return fmt.Errorf("unable to query billettholder_emails: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result models.BillettholderEmail
		err := rows.Scan(&result.ID, &result.BillettholderID, &result.Email, &result.Kind, &result.InsertedTime)
		if err != nil {
			return fmt.Errorf("unable to scan billettholder_emails: %v", err)
		}
		billettholdere = append(billettholdere, result)
	}

	// Insert into billettholder_users the new data
	var lines []string

	for _, billettholder := range billettholdere {
		lines = append(lines, fmt.Sprintf(`(%d, %d)`, billettholder.BillettholderID, user.ID))
	}

	var baseQuery = fmt.Sprintf(`
        INSERT INTO billettholdere_users (
            billettholder_id, user_id
        ) VALUES %s
    `, strings.Join(lines, ", "))

	_, err = db.Exec(baseQuery)
	if err != nil {
		fmt.Printf("UserID: %s has id: %d \n\n\n", userID, user.ID)
		/* var outputTest []int
		for _, billetID := range billettholdere {
			outputTest = append(outputTest, billetID.BillettholderID)
		} */
		for _, billet := range billettholdere {
			fmt.Printf("Billettholdere: %+v \n", billet)
		}
		fmt.Print("\n\n\n")
		fmt.Println(baseQuery)
		return fmt.Errorf("unable to insert into billettholder_users: %v", err)
	}

	return nil
}
