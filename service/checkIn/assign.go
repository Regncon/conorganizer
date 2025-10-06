package checkIn

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

// AssociateTicketsWithEmail takes a list of tickets and matches with email supplied, returns matches
func AssociateTicketsWithEmail(tickets []CheckInTicket, email string) ([]CheckInTicket, error) {
	var result []CheckInTicket
	for _, ticket := range tickets {
		if ticket.TypeId == TicketTypeMiddag {
			continue
		}

		if strings.EqualFold(ticket.Email, email) {
			result = append(result, ticket)
		}
	}

	if len(result) < 1 {
		return nil, errors.New("found 0 tickets registered on: " + email)
	}

	return result, nil
}

// todo AssociateBillettholderWithEmail should be called and then feed inn result to billettholder table
// todo etter funksjon er called og data skal føres inn i db. Ikke opprett nye billettholdere som eksisterer fra før, og ikke koblinger (billettholder_users)

func AssociateTicketsWithBillettholder(tickets []CheckInTicket, email string) error {
	// Filtrer tickets til de som er registrert på user email
	_, err := AssociateTicketsWithEmail(tickets, email)
	if err != nil {
		return err
	}

	// get existing billetterholdere registered to user email

	// Create new array of unique non-existing tickets

	// Enter new array into billettholdere

	// call AssociateUserWithBillettholder

	return nil
}

// AssociateUserWithBillettholder uses userID string from users table to match billettholders
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
        SELECT id, billettholder_id, email, kind, inserted_time FROM billettholder_emails WHERE email = ? COLLATE NOCASE
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

	if len(billettholdere) < 1 {
		//return errors.New("no billettholdere found on user: " + userID)
		return nil
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
		fmt.Printf("UserID: %s has id: %d \n", userID, user.ID)
		for _, billet := range billettholdere {
			fmt.Printf("Billettholdere: %+v \n", billet)
		}
		fmt.Println(baseQuery)
		return fmt.Errorf("unable to insert into billettholder_users: %v", err)
	}

	return nil
}
