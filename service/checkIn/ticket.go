package checkIn

import (
	"database/sql"
	"log/slog"
)

type CheckInTicket struct {
	ID        int
	OrderID   int
	TypeId    int
	Type      string
	FirstName string
	LastName  string
	Email     string
	IsOver18  bool
}

func GetTicketsFromCheckIn(logger *slog.Logger, searchTerm string) ([]CheckInTicket, error) {

	return ticketCache.Get(logger, searchTerm)
}

func ConvertTicketToBillettholder(ticketId int, db *sql.DB, logger *slog.Logger) (*CheckInTicket, error) {
	tickets, err := GetTicketsFromCheckIn(logger, "")
	if err != nil {
		return nil, err
	}

	converTicketIdToNewBillettholder(ticketId, tickets, db, logger)
	return nil, nil
}
