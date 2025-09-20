package checkIn

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/Regncon/conorganizer/models"
)

func converTicketIdToNewBillettholder(ticketId int, tickets []CheckInTicket, db *sql.DB, logger *slog.Logger) error {
	logger.Info("Converting ticket to billettholder", "ticketID", ticketId)

	var ticket *CheckInTicket
	for _, t := range tickets {
		if t.ID == ticketId {
			ticket = &t
			break
		}
	}

	if ticket == nil {
		logger.Error("ticket not found", "ticketId", ticketId)
		return errors.New("ticket not found")
	}
	const TicketTypeMiddag = 193284
	if ticket.TypeId == TicketTypeMiddag {
		logger.Error("cannot convert 'Middag' ticket to billettholder", "ticketId", ticketId)
		return errors.New("cannot convert 'Middag' ticket to billettholder")
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
		logger.Info("failed to insert billettholder", "error", err)
		return err
	}

	billettholderID, lastIdErr := result.LastInsertId()
	if lastIdErr != nil {
		logger.Info("failed to fetch last insert ID", "error", lastIdErr)
		return lastIdErr
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
		_, err := db.Exec(`
			INSERT INTO billettholder_emails (
				billettholder_id, email, kind
			) VALUES (?, ?, ?)
		`, email.BillettholderID, email.Email, email.Kind)
		if err != nil {
			logger.Info("failed to insert billettholder email", "error", err)
			return err
		}
	}

	logger.Info("successfully inserted billettholder", "ticketId", ticketId, "billettholderId", billettholderID)
	return nil
}
