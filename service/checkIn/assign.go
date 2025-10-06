package checkIn

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

// uniquePair is used for getting unique values when comparing tickets and billettholder
type uniquePair struct {
	TicketID int
	Ticket   CheckInTicket
}

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

// AssociateTicketsWithBillettholder is responsible for finding tickets registered on an email and
// inserting new unique tickets into billettholder
func AssociateTicketsWithBillettholder(tickets []CheckInTicket, email string, db *sql.DB, logger *slog.Logger) error {
	// Filtrer tickets til de som er registrert på user email
	associatedTickets, err := AssociateTicketsWithEmail(tickets, email)
	if err != nil {
		// No associated tickets found, quitting early
		fmt.Printf("Found no tickets associated with %s, quitting early\n", email)
		return nil
	}
	fmt.Printf("Found %d/%d tickets associated with %s\n", len(associatedTickets), len(tickets), email)

	// get list of ticket ids to exclude when quering billettholder

	// get existing billetterholdere registered to user email
	var billettholdereIDs []models.Billettholder
	rows, err := db.Query(`
        SELECT DISTINCT b.ticket_id
        FROM billettholder_emails e
        JOIN billettholdere b ON b.id = e.billettholder_id
        WHERE e.email = ? COLLATE NOCASE;
    `, email)
	if err != nil {
		return fmt.Errorf("unable to query billettholder: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var result models.Billettholder
		err := rows.Scan(&result.TicketID)
		if err != nil {
			return fmt.Errorf("unable to scan billettholder: %v", err)
		}
		billettholdereIDs = append(billettholdereIDs, result)
	}
	fmt.Printf("Found %d existing billettholdere with email: %s\n", len(billettholdereIDs), email)

	// Create new map of unique non-existing tickets
	var uniqueNewTicketsMap = map[uniquePair]struct{}{}
	for _, associatedTicket := range associatedTickets {
		var billettholderExists = false

		for _, billetholderID := range billettholdereIDs {
			if billetholderID.TicketID == associatedTicket.ID {
				billettholderExists = true
				break
			}
		}

		if !billettholderExists {
			uniqueNewTicketsMap[uniquePair{TicketID: associatedTicket.ID, Ticket: associatedTicket}] = struct{}{}
		}
	}

	// Convert unique map back to array
	var uniqueNewTickets []CheckInTicket
	for pair := range uniqueNewTicketsMap {
		uniqueNewTickets = append(uniqueNewTickets, pair.Ticket)
	}

	// No new unique tickets, quitting early
	if len(uniqueNewTickets) == 0 {
		fmt.Println("Found no new unique tickets to add, quitting early")
		return nil
	}

	// Enter new array into billettholdere ... newTicket needs to be unique?
	for _, ticket := range uniqueNewTickets {
		err = converTicketIdToNewBillettholder(ticket.ID, uniqueNewTickets, db, logger)
		fmt.Printf("Adding Ticket: %+v \n", ticket)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Added %d new billettholdere from %d tickets\n", len(uniqueNewTickets), len(tickets))
	return nil
}

// AssociateUserWithBillettholder uses userID string from users table to match billettholders
// and combine ids to billettholder_users for later lookup
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
