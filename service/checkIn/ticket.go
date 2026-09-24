package checkIn

import (
	"context"
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

// TicketFetchResult makes stale-cache use explicit to callers. Tickets may be
// usable even when err reports that the latest refresh failed.
type TicketFetchResult struct {
	Tickets        []CheckInTicket
	UsedStaleCache bool
}

const TicketTypeMiddag = 251934

func GetTicketsFromCheckIn(ctx context.Context, logger *slog.Logger, searchTerm string) (TicketFetchResult, error) {
	return ticketCache.Get(ctx, logger, searchTerm)
}

func ConvertTicketToBillettholder(ctx context.Context, ticketId int, db *sql.DB, logger *slog.Logger) error {
	result, err := GetTicketsFromCheckIn(ctx, logger, "")
	tickets := result.Tickets
	if err != nil && !result.UsedStaleCache {
		return fmt.Errorf("failed to fetch tickets from check-in: %w", err)
	}

	if _, err := converTicketIdToNewBillettholder(ticketId, tickets, db, logger); err != nil {
		return fmt.Errorf("failed to convert ticket %d to billettholder: %w", ticketId, err)
	}
	return nil
}
