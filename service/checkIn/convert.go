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

	billettholder := models.Billettholder{
		FirstName: ticket.FirstName,
		LastName:  ticket.LastName,
		OrderID:   ticket.OrderID,
		TicketID:  ticket.ID,
		IsOver18:  ticket.IsOver18,
	}

	result, err := db.Exec(`
		INSERT INTO billettholdere (
        first_name, last_name, ticket_type,
        ticket_id, is_over_18, order_id,
        ticket_category_id
		) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		billettholder.FirstName, billettholder.LastName, billettholder.TicketType,
		billettholder.TicketID, billettholder.IsOver18, billettholder.OrderID,
		billettholder.TicketCategoryID,
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
