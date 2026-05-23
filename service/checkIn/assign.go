package checkIn

import (
	"database/sql"
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
		return nil, fmt.Errorf("found 0 tickets registered on: %s", email)
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
		// fmt.Printf("Found no tickets associated with %s, quitting early\n", email)
		return nil
	}
	// fmt.Printf("Found %d/%d tickets associated with %s\n", len(associatedTickets), len(tickets), email)

	// get list of ticket ids to exclude when quering billettholder

	// get existing billetterholdere registered to user email
	var billettholdereIDs []models.Billettholder
	rows, err := db.Query(`
        SELECT DISTINCT b.ticket_id
        FROM relation_billettholder_emails e
        JOIN billettholdere b ON b.id = e.billettholder_id
        WHERE e.email = ? COLLATE NOCASE;
    `, email)
	if err != nil {
		return fmt.Errorf("unable to query billettholder for email %q: %w", email, err)
	}
	defer rows.Close()
	for rows.Next() {
		var result models.Billettholder
		err := rows.Scan(&result.TicketID)
		if err != nil {
			return fmt.Errorf("unable to scan billettholder for email %q: %w", email, err)
		}
		billettholdereIDs = append(billettholdereIDs, result)
	}
	// fmt.Printf("Found %d existing billettholdere with email: %s\n", len(billettholdereIDs), email)

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
		// fmt.Println("Found no new unique tickets to add, quitting early")
		return nil
	}

	// Enter new array into billettholdere ... newTicket needs to be unique?
	for _, ticket := range uniqueNewTickets {
		err = converTicketIdToNewBillettholder(ticket.ID, uniqueNewTickets, db, logger)
		if err != nil {
			return fmt.Errorf("unable to convert ticket %d to billettholder for email %q: %w", ticket.ID, email, err)
		}
	}

	// fmt.Printf("Added %d new billettholdere from %d tickets\n", len(uniqueNewTickets), len(tickets))
	return nil
}

// AssociateUserWithBillettholder uses userID string from users table to match billettholders
// and combine ids to billettholder_users for later lookup
func AssociateUserWithBillettholder(userID string, db *sql.DB, logger *slog.Logger) error {
	logger = logger.With("component", "checkin_assign")
	logger.Debug("Associating user with billettholder", "user_id", userID)

	// Get user
	var user models.User
	err := db.QueryRow(`
        SELECT id, email FROM users WHERE external_id = ?;
    `, userID).Scan(&user.ID, &user.Email)
	if err != nil {
		return fmt.Errorf("failed to get user %q: %w", userID, err)
	}

	// Get associated billettholdere
	var billettholdere []models.BillettholderEmail
	rows, err := db.Query(`
        SELECT id, billettholder_id, email, kind, created_at, updated_at, created_by_id, updated_by_id FROM relation_billettholder_emails WHERE email = ? COLLATE NOCASE
    `, user.Email)
	if err != nil {
		return fmt.Errorf("unable to query relation_billettholder_emails for email %q: %w", user.Email, err)
	}
	defer rows.Close()

	for rows.Next() {
		var result models.BillettholderEmail
		err := rows.Scan(&result.ID, &result.BillettholderID, &result.Email, &result.Kind, &result.CreatedAt, &result.UpdatedAt, &result.CreatedByID, &result.UpdatedByID)
		if err != nil {
			return fmt.Errorf("unable to scan relation_billettholder_emails for email %q: %w", user.Email, err)
		}
		billettholdere = append(billettholdere, result)
	}

	if len(billettholdere) < 1 {
		return nil
	}

	// Insert into billettholder_users the new data
	var lines []string
	for _, billettholder := range billettholdere {
		lines = append(lines, fmt.Sprintf(`(%d, %d)`, billettholder.BillettholderID, user.ID))
	}
	var baseQuery = fmt.Sprintf(`
        INSERT OR IGNORE INTO relation_billettholdere_users (
            billettholder_id, user_id
        ) VALUES %s
    `, strings.Join(lines, ", "))

	_, err = db.Exec(baseQuery)
	if err != nil {
		fmt.Printf("UserID: %s has id: %d \n", userID, user.ID)
		for _, billet := range billettholdere {
			fmt.Printf("Billettholdere: %+v \n", billet)
		}

		return fmt.Errorf("unable to insert into relation_billettholdere_users: %v", err)
	}

	return nil
}

// AssociateUsersWithBillettholderEmail links an existing billettholder to all users
// whose email matches the supplied billettholder email.
func AssociateUsersWithBillettholderEmail(billettholderID int, email string, db *sql.DB, logger *slog.Logger) error {
	logger = logger.With("component", "checkin_assign")
	logger.Debug("Associating users with billettholder email", "billettholder_id", billettholderID)

	result, err := db.Exec(`
		INSERT OR IGNORE INTO relation_billettholdere_users (
			billettholder_id, user_id
		)
		SELECT ?, id
		FROM users
		WHERE email = ? COLLATE NOCASE
	`, billettholderID, email)
	if err != nil {
		return fmt.Errorf("unable to associate users with billettholder %d by email: %w", billettholderID, err)
	}

	if rowsAffected, err := result.RowsAffected(); err == nil {
		if rowsAffected > 0 {
			logger.Info("Created billettholder user associations",
				"billettholder_id", billettholderID,
				"association_flow", "billettholder_email",
				"created_associations", rowsAffected,
			)
		} else {
			logger.Debug("No billettholder user associations created", "billettholder_id", billettholderID)
		}
	} else {
		logger.Debug("Unable to read created association count",
			"billettholder_id", billettholderID,
			"error", err,
		)
	}

	return nil
}

// DisassociateUsersFromBillettholderEmail removes user links for a removed
// billettholder email when no remaining email still matches the same users.
func DisassociateUsersFromBillettholderEmail(billettholderID int, email string, db *sql.DB, logger *slog.Logger) error {
	logger = logger.With("component", "checkin_assign")
	logger.Debug("Disassociating users from billettholder email", "billettholder_id", billettholderID)

	result, err := db.Exec(`
		DELETE FROM relation_billettholdere_users
		WHERE billettholder_id = ?
		AND user_id IN (
			SELECT id
			FROM users
			WHERE email = ? COLLATE NOCASE
		)
		AND NOT EXISTS (
			SELECT 1
			FROM relation_billettholder_emails
			WHERE billettholder_id = ?
			AND email = ? COLLATE NOCASE
		)
	`, billettholderID, email, billettholderID, email)
	if err != nil {
		return fmt.Errorf("unable to disassociate users from billettholder %d by email: %w", billettholderID, err)
	}

	if rowsAffected, err := result.RowsAffected(); err == nil {
		if rowsAffected > 0 {
			logger.Info("Removed billettholder user associations",
				"billettholder_id", billettholderID,
				"association_flow", "billettholder_email",
				"removed_associations", rowsAffected,
			)
		} else {
			logger.Debug("No billettholder user associations removed", "billettholder_id", billettholderID)
		}
	} else {
		logger.Debug("Unable to read removed association count",
			"billettholder_id", billettholderID,
			"error", err,
		)
	}

	return nil
}
