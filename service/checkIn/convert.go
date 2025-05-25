package checkIn

import (
	"database/sql"
	"log/slog"
)

func converTicketIdToNewBillettholder(ticketId int, tickets []CheckInTicket, db *sql.DB, logger *slog.Logger) {
	logger.Info("Converting ticket to billettholder", "ticketID", ticketId)
}
