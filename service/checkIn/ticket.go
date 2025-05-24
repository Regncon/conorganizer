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

type crm struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Born      string `json:"born"`
}

type eventTicket struct {
	ID         int    `json:"id"`
	Category   string `json:"category"`
	CategoryID int    `json:"category_id"`
	Crm        crm    `json:"crm"`
	OrderID    int    `json:"order_id"`
}

type queryResult struct {
	Data struct {
		EventTickets []eventTicket `json:"eventTickets"`
	} `json:"data"`
}

func GetTicketsFromCheckIn(logger *slog.Logger, searchTerm string) ([]CheckInTicket, error) {

	return ticketCache.Get(logger, searchTerm)
}

func ConvertTicketToBilettholder(ticketId int, db *sql.DB, logger *slog.Logger) (*CheckInTicket, error) {
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
