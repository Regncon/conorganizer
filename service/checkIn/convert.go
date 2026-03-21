package checkIn

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Regncon/conorganizer/models"
)

func converTicketIdToNewBillettholder(ticketId int, tickets []CheckInTicket, db *sql.DB, logger *slog.Logger) error {
	logger = logger.With("component", "checkin_convert")
	logger.Info("Converting ticket to billettholder", "ticket_id", ticketId)

	var ticket *CheckInTicket
	for _, t := range tickets {
		if t.ID == ticketId {
			ticket = &t
			break
		}
	}

	if ticket == nil {
		notFoundErr := fmt.Errorf("ticket %d not found", ticketId)
		logger.Error(notFoundErr.Error())
		return notFoundErr
	}
	if ticket.TypeId == TicketTypeMiddag {
		unsupportedTicketErr := fmt.Errorf("cannot convert 'Middag' ticket to billettholder")
		logger.Error(unsupportedTicketErr.Error())
		return unsupportedTicketErr
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

	var exists bool
	billettholderExistsErr := db.QueryRow(`
        SELECT EXISTS(
            SELECT 1 FROM billettholdere
            WHERE first_name = ? AND last_name = ? AND ticket_id = ?
        )
    `, billettholder.FirstName, billettholder.LastName, billettholder.TicketID).Scan(&exists)

	if billettholderExistsErr != nil {
		checkExistsErr := fmt.Errorf("failed to check if billettholder exists for ticket %d: %w", ticketId, billettholderExistsErr)
		logger.Error(checkExistsErr.Error())
		return checkExistsErr
	}

	var billettholderID int64
	selectErr := db.QueryRow(`
		SELECT id FROM billettholdere
		WHERE first_name = ? AND last_name = ? AND ticket_id = ?
	`, billettholder.FirstName, billettholder.LastName, billettholder.TicketID).Scan(&billettholderID)

	if selectErr == sql.ErrNoRows {
		result, err := db.Exec(`
			INSERT INTO billettholdere (
				first_name, last_name, ticket_type_id, ticket_type,
				ticket_id, is_over_18, order_id
			) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			billettholder.FirstName, billettholder.LastName, billettholder.TicketTypeId,
			billettholder.TicketType, billettholder.TicketID, billettholder.IsOver18,
			billettholder.OrderID,
		)
		if err != nil {
			insertErr := fmt.Errorf("failed to insert billettholder for ticket %d: %w", ticketId, err)
			logger.Error(insertErr.Error())
			return insertErr
		}
		billettholderID, err = result.LastInsertId()
		if err != nil {
			insertIDErr := fmt.Errorf("failed to fetch last insert ID for ticket %d: %w", ticketId, err)
			logger.Error(insertIDErr.Error())
			return insertIDErr
		}
	} else if selectErr != nil {
		selectBillettholderErr := fmt.Errorf("failed to select billettholder for ticket %d: %w", ticketId, selectErr)
		logger.Error(selectBillettholderErr.Error())
		return selectBillettholderErr
	}

	emails := []models.BillettholderEmail{
		{
			BillettholderID: int(billettholderID),
			Email:           ticket.Email,
			Kind:            "Ticket",
		},
	}
	// find associated emails if any. An associated email is any email that is in a ticket with the same order ID but is not the ticket email
	for _, t := range tickets {
		if t.OrderID == ticket.OrderID && t.Email != ticket.Email {
			associatedEmail := models.BillettholderEmail{
				BillettholderID: int(billettholderID),
				Email:           t.Email,
				Kind:            "Associated",
			}
			emails = append(emails, associatedEmail)
		}
	}

	for _, email := range emails {
		exists := false
		checkErr := db.QueryRow(`
			SELECT EXISTS(
				SELECT 1 FROM billettholder_emails
				WHERE billettholder_id = ? AND email = ?
			)
		`, email.BillettholderID, email.Email).Scan(&exists)
		if checkErr != nil {
			checkEmailErr := fmt.Errorf("failed to check existing email for billettholder %d: %w", email.BillettholderID, checkErr)
			logger.Error(checkEmailErr.Error())
			return checkEmailErr
		}
		if exists {
			logger.Debug("Email already exists, skipping", "billettholder_id", email.BillettholderID)
			continue
		}

		_, err := db.Exec(`
			INSERT INTO billettholder_emails (
				billettholder_id, email, kind
			) VALUES (?, ?, ?)
		`, email.BillettholderID, email.Email, email.Kind)
		if err != nil {
			insertEmailErr := fmt.Errorf("failed to insert billettholder email for billettholder %d: %w", email.BillettholderID, err)
			logger.Error(insertEmailErr.Error())
			return insertEmailErr
		}
	}

	logger.Info("Successfully inserted billettholder", "ticket_id", ticketId, "billettholder_id", billettholderID)
	return nil
}
