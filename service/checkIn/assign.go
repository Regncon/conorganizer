package checkIn

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Regncon/conorganizer/models"
)

type TicketAssociationResult struct {
	BillettholderIDs      []int
	CreatedBillettholders int
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

// AssociateTicketsWithBillettholder imports every non-dinner ticket from an
// order containing a ticket registered to email.
func AssociateTicketsWithBillettholder(tickets []CheckInTicket, email string, db *sql.DB, logger *slog.Logger) (TicketAssociationResult, error) {
	var result TicketAssociationResult

	associatedTickets, err := AssociateTicketsWithEmail(tickets, email)
	if err != nil {
		return result, nil
	}

	associatedOrderIDs := make(map[int]struct{}, len(associatedTickets))
	for _, associatedTicket := range associatedTickets {
		associatedOrderIDs[associatedTicket.OrderID] = struct{}{}
	}

	for _, ticket := range tickets {
		if ticket.TypeId == TicketTypeMiddag {
			continue
		}
		if _, belongsToAssociatedOrder := associatedOrderIDs[ticket.OrderID]; !belongsToAssociatedOrder {
			continue
		}

		conversionResult, err := converTicketIdToNewBillettholder(ticket.ID, tickets, db, logger)
		if err != nil {
			return result, fmt.Errorf("unable to convert ticket %d to billettholder for email %q: %w", ticket.ID, email, err)
		}
		result.BillettholderIDs = append(result.BillettholderIDs, conversionResult.BillettholderID)
		result.CreatedBillettholders += conversionResult.CreatedBillettholders
	}

	return result, nil
}

// AssociateUserWithBillettholder uses userID string from users table to match billettholders
// and combine ids to billettholder_users for later lookup
func AssociateUserWithBillettholder(userID string, db *sql.DB, logger *slog.Logger) (int, error) {
	logger = logger.With("component", "checkin_assign")
	logger.Debug("Associating user with billettholder", "user_id", userID)

	// Get user
	var user models.User
	err := db.QueryRow(`
		SELECT id, email FROM users WHERE external_id = ?;
	`, userID).Scan(&user.ID, &user.Email)
	if err != nil {
		return 0, fmt.Errorf("failed to get user %q: %w", userID, err)
	}

	// Get associated billettholdere
	var billettholdere []models.BillettholderEmail
	rows, err := db.Query(`
		SELECT id, billettholder_id, email, kind, created_at, updated_at, created_by_id, updated_by_id
		FROM relation_billettholder_emails
		WHERE email = ? COLLATE NOCASE
	`, user.Email)
	if err != nil {
		return 0, fmt.Errorf("unable to query relation_billettholder_emails for email %q: %w", user.Email, err)
	}
	defer rows.Close()

	for rows.Next() {
		var result models.BillettholderEmail
		err := rows.Scan(&result.ID, &result.BillettholderID, &result.Email, &result.Kind, &result.CreatedAt, &result.UpdatedAt, &result.CreatedByID, &result.UpdatedByID)
		if err != nil {
			return 0, fmt.Errorf("unable to scan relation_billettholder_emails for email %q: %w", user.Email, err)
		}
		billettholdere = append(billettholdere, result)
	}
	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("unable to iterate relation_billettholder_emails for email %q: %w", user.Email, err)
	}

	if len(billettholdere) < 1 {
		return 0, nil
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

	insertResult, err := db.Exec(baseQuery)
	if err != nil {
		fmt.Printf("UserID: %s has id: %d \n", userID, user.ID)
		for _, billet := range billettholdere {
			fmt.Printf("Billettholdere: %+v \n", billet)
		}

		return 0, fmt.Errorf("unable to insert into relation_billettholdere_users: %v", err)
	}

	rowsAffected, err := insertResult.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to read created billettholder user association count: %w", err)
	}

	return int(rowsAffected), nil
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

// LinkUsersToBillettholdere links every existing user whose email is on one of
// the billettholdere, whatever the email kind. It returns the number of new links.
func LinkUsersToBillettholdere(billettholderIDs []int, db *sql.DB, logger *slog.Logger) (int, error) {
	if len(billettholderIDs) == 0 {
		return 0, nil
	}
	logger = logger.With("component", "checkin_assign")

	placeholders := make([]string, 0, len(billettholderIDs))
	args := make([]any, 0, len(billettholderIDs))
	for _, billettholderID := range billettholderIDs {
		placeholders = append(placeholders, "?")
		args = append(args, billettholderID)
	}

	result, err := db.Exec(fmt.Sprintf(`
		INSERT OR IGNORE INTO relation_billettholdere_users (
			billettholder_id, user_id
		)
		SELECT DISTINCT e.billettholder_id, u.id
		FROM relation_billettholder_emails AS e
		JOIN users AS u ON u.email = e.email COLLATE NOCASE
		WHERE e.billettholder_id IN (%s)
	`, strings.Join(placeholders, ", ")), args...)
	if err != nil {
		return 0, fmt.Errorf("unable to link users to billettholdere %v by email: %w", billettholderIDs, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("unable to read created billettholder user association count: %w", err)
	}
	if rowsAffected > 0 {
		logger.Info("Created billettholder user associations",
			"billettholder_ids", billettholderIDs,
			"association_flow", "billettholder_emails",
			"created_associations", rowsAffected,
		)
	}

	return int(rowsAffected), nil
}

type UserTicketImportResult struct {
	CreatedBillettholders   int
	CreatedUserAssociations int
}

// ImportUserTickets is the "Hent billetter" flow for one user. It imports every
// order with a ticket registered to email, links the user to every billettholder
// carrying their email, and links any other existing users whose email is on the
// imported billettholdere. CreatedUserAssociations only counts the user's own links.
func ImportUserTickets(tickets []CheckInTicket, userID string, email string, db *sql.DB, logger *slog.Logger) (UserTicketImportResult, error) {
	var result UserTicketImportResult

	ticketResult, err := AssociateTicketsWithBillettholder(tickets, email, db, logger)
	if err != nil {
		return result, err
	}
	result.CreatedBillettholders = ticketResult.CreatedBillettholders

	result.CreatedUserAssociations, err = AssociateUserWithBillettholder(userID, db, logger)
	if err != nil {
		return result, err
	}

	if _, err := LinkUsersToBillettholdere(ticketResult.BillettholderIDs, db, logger); err != nil {
		return result, err
	}

	return result, nil
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
