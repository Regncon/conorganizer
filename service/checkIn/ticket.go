package checkIn

import (
	"database/sql"
	"fmt"
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

const TicketTypeMiddag = 193284

func GetTicketsFromCheckIn(logger *slog.Logger, searchTerm string) ([]CheckInTicket, error) {

	return ticketCache.Get(logger, searchTerm)
}

func ConvertTicketToBillettholder(ticketId int, db *sql.DB, logger *slog.Logger) error {
	tickets, err := GetTicketsFromCheckIn(logger, "")
	if err != nil {
		return fmt.Errorf("failed to fetch tickets from check-in: %w", err)
	}

	if err := converTicketIdToNewBillettholder(ticketId, tickets, db, logger); err != nil {
		return fmt.Errorf("failed to convert ticket %d to billettholder: %w", ticketId, err)
	}
	return nil
}
