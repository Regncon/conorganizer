package checkIn

import (
	"database/sql"
	"log/slog"
)

func converTicketIdToNewBilettholder(ticketId int, tickets []CheckInTicket, db *sql.DB, logger *slog.Logger) {
	logger.Info("Converting ticket to bilettholder", "ticketID", ticketId)
}
