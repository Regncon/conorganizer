package checkIn

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"
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

// ConvertOrderToBillettholdere is the admin conversion. It converts every
// non-dinner ticket in the order that ticketID belongs to, like "Hent billetter"
// does, and links every existing user whose email is on those billettholdere.
func ConvertOrderToBillettholdere(ctx context.Context, ticketID int, db *sql.DB, logger *slog.Logger) (TicketAssociationResult, error) {
	fetchResult, err := GetTicketsFromCheckIn(ctx, logger, "")
	tickets := fetchResult.Tickets
	if err != nil && !fetchResult.UsedStaleCache {
		return TicketAssociationResult{}, fmt.Errorf("failed to fetch tickets from check-in: %w", err)
	}

	ticketIndex := slices.IndexFunc(tickets, func(ticket CheckInTicket) bool { return ticket.ID == ticketID })
	if ticketIndex == -1 {
		return TicketAssociationResult{}, fmt.Errorf("ticket %d not found", ticketID)
	}
	orderID := tickets[ticketIndex].OrderID

	result, err := convertOrders(map[int]struct{}{orderID: {}}, tickets, db, logger)
	if err != nil {
		return result, fmt.Errorf("failed to convert order %d to billettholdere: %w", orderID, err)
	}
	if len(result.BillettholderIDs) == 0 {
		return result, fmt.Errorf("order %d has no tickets that can become billettholdere", orderID)
	}
	if _, err := LinkUsersToBillettholdere(result.BillettholderIDs, db, logger); err != nil {
		return result, fmt.Errorf("failed to link users to order %d: %w", orderID, err)
	}
	return result, nil
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
