package checkIn

import (
	"database/sql"
	"log/slog"
)

type CheckInTicket struct {
	ID      int
	OrderID int
	Type    string
	Name    string
	Email   string
	IsAdult bool
}

func GetTicketsFromCheckIn(logger *slog.Logger, searchTerm string) ([]CheckInTicket, error) {

	return ticketCache.Get(logger, searchTerm)
}

func ConvertTicketToBilettholder(ticketId int, db *sql.DB, logger *slog.Logger) (*CheckInTicket, error) {
	logger.Info("Converting ticket to bilettholder", "ticketID", ticketId)
	tickets, err := GetTicketsFromCheckIn(logger, "")
	if err != nil {
		return nil, err
	}

	for _, ticket := range tickets {
		if ticket.ID == ticketId {
			return &ticket, nil
		}
	}

	logger.Error("ticket not found", "ticketId", ticketId)
	return nil, nil
}
