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
		FirstName:   ticket.FirstName,
		LastName:    ticket.LastName,
		OrderID:     ticket.OrderID,
		TicketID:    ticket.ID,
		IsOver18:    ticket.IsOver18,
		TicketEmail: ticket.Email,
	}

	_, err := db.Exec(`
		INSERT INTO billettholdere (
        first_name, last_name, ticket_type,
        ticket_id, is_over_18, order_id,
        ticket_email, order_email, ticket_category_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		billettholder.FirstName, billettholder.LastName, billettholder.TicketType,
		billettholder.TicketID, billettholder.IsOver18, billettholder.OrderID,
		billettholder.TicketEmail, billettholder.OrderEmail, billettholder.TicketCategoryID,
	)

	if err != nil {
		logger.Info("failed to insert billettholder", "error", err)
		return err
	}

	logger.Info("successfully inserted billettholder", "ticketId", ticketId)
	return nil
}
