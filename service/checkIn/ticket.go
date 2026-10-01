package checkIn

import (
	"context"
	"database/sql"
	"errors"
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

	conversionResult, err := converTicketIdToNewBillettholder(ticketId, tickets, db, logger)
	if err != nil {
		return fmt.Errorf("failed to convert ticket %d to billettholder: %w", ticketId, err)
	}
	if _, err := LinkUsersToBillettholdere([]int{conversionResult.BillettholderID}, db, logger); err != nil {
		return fmt.Errorf("failed to link users to converted ticket %d: %w", ticketId, err)
	}
	return nil
}

// SyncUserTicketsFromCheckIn runs "Hent billetter" for a user without a request
// to report back to, such as right after login. When CheckIn is unavailable the
// user is still linked to billettholdere that already carry their email.
func SyncUserTicketsFromCheckIn(ctx context.Context, userID string, email string, db *sql.DB, logger *slog.Logger) (UserTicketImportResult, error) {
	fetchResult, fetchErr := GetTicketsFromCheckIn(ctx, logger, "")
	if fetchErr != nil && !fetchResult.UsedStaleCache {
		fetchErr = fmt.Errorf("failed to fetch tickets from check-in: %w", fetchErr)
	} else {
		fetchErr = nil
	}

	result, importErr := ImportUserTickets(fetchResult.Tickets, userID, email, db, logger)
	return result, errors.Join(fetchErr, importErr)
}
