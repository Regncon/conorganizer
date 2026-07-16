package checkIn

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Regncon/conorganizer/models"
)

type TicketConversionResult struct {
	CreatedBillettholders int
}

func converTicketIdToNewBillettholder(ticketId int, tickets []CheckInTicket, db *sql.DB, logger *slog.Logger) (TicketConversionResult, error) {
	logger = logger.With("component", "checkin_convert")
	logger.Debug("Converting ticket to billettholder", "ticket_id", ticketId)
	var result TicketConversionResult

	var ticket *CheckInTicket
	for _, t := range tickets {
		if t.ID == ticketId {
			ticket = &t
			break
		}
	}

	if ticket == nil {
		return result, fmt.Errorf("ticket %d not found", ticketId)
	}
	if ticket.TypeId == TicketTypeMiddag {
		return result, fmt.Errorf("cannot convert 'Middag' ticket to billettholder")
	}

	billettholder := models.Billettholder{
		FirstName:    ticket.FirstName,
		LastName:     ticket.LastName,
		TicketTypeId: ticket.TypeId,
		TicketType:   ticket.Type,
		OrderID:      ticket.OrderID,
		TicketID:     ticket.ID,
		IsOver18:     ticket.IsOver18,
	}

	var billettholderID int64
	selectErr := db.QueryRow(`
		SELECT id FROM billettholdere
		WHERE ticket_id = ?
	`, billettholder.TicketID).Scan(&billettholderID)

	if selectErr == sql.ErrNoRows {
		insertResult, err := db.Exec(`
			INSERT INTO billettholdere (
				first_name, last_name, ticket_type_id, ticket_type,
				ticket_id, is_over_18, order_id
			) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			billettholder.FirstName, billettholder.LastName, billettholder.TicketTypeId,
			billettholder.TicketType, billettholder.TicketID, billettholder.IsOver18,
			billettholder.OrderID,
		)
		if err != nil {
			return result, fmt.Errorf("failed to insert billettholder for ticket %d: %w", ticketId, err)
		}
		billettholderID, err = insertResult.LastInsertId()
		if err != nil {
			return result, fmt.Errorf("failed to fetch last insert ID for ticket %d: %w", ticketId, err)
		}
		result.CreatedBillettholders = 1
	} else if selectErr != nil {
		return result, fmt.Errorf("failed to select billettholder for ticket %d: %w", ticketId, selectErr)
	}

	emails := []models.BillettholderEmail{
		{
			BillettholderID: int(billettholderID),
			Email:           ticket.Email,
			Kind:            models.BillettholderEmailKindTicket,
		},
	}
	// find associated emails if any. An associated email is any email that is in a ticket with the same order ID but is not the ticket email
	for _, t := range tickets {
		if t.OrderID == ticket.OrderID && t.Email != ticket.Email {
			associatedEmail := models.BillettholderEmail{
				BillettholderID: int(billettholderID),
				Email:           t.Email,
				Kind:            models.BillettholderEmailKindAssociated,
			}
			emails = append(emails, associatedEmail)
		}
	}

	for _, email := range emails {
		exists := false
		checkErr := db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM relation_billettholder_emails
				WHERE billettholder_id = ? AND email = ?
			)
		`, email.BillettholderID, email.Email).Scan(&exists)
		if checkErr != nil {
			return result, fmt.Errorf("failed to check existing email for billettholder %d: %w", email.BillettholderID, checkErr)
		}
		if exists {
			logger.Debug("Email already exists, skipping", "billettholder_id", email.BillettholderID)
			continue
		}

		_, err := db.Exec(`
			INSERT INTO relation_billettholder_emails (
				billettholder_id, email, kind
			) VALUES (?, ?, ?)
		`, email.BillettholderID, email.Email, email.Kind)
		if err != nil {
			return result, fmt.Errorf("failed to insert billettholder email for billettholder %d: %w", email.BillettholderID, err)
		}
	}

	logger.Debug("Successfully converted ticket to billettholder", "ticket_id", ticketId, "billettholder_id", billettholderID)
	return result, nil
}
