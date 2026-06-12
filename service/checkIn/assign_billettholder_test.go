package checkIn

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestAssociateTicketsWithBillettholder_WhenSomeMatchingTicketsAreNew_ConvertsOnlyNewTickets(t *testing.T) {
	// Given tickets for a target email where one ticket is already converted,
	// when the target email is associated with billettholdere,
	// then only new non-dinner target tickets are inserted.

	// Given
	expectedBillettholderCount := 3
	expectedTargetEmailCount := 3
	expectedTicketIDs := []int{101, 102, 103}
	targetEmail := "test@regncon.com"
	tickets := []CheckInTicket{
		{ID: 101, OrderID: 1, TypeId: 9000, Type: "Festivalpass", FirstName: "New", LastName: "Target", Email: targetEmail, IsOver18: true},
		{ID: 102, OrderID: 2, TypeId: 9000, Type: "Festivalpass", FirstName: "Existing", LastName: "Target", Email: targetEmail, IsOver18: true},
		{ID: 103, OrderID: 3, TypeId: 9000, Type: "Festivalpass", FirstName: "Case", LastName: "Target", Email: "TEST@REGNCON.COM", IsOver18: true},
		{ID: 104, OrderID: 4, TypeId: 9000, Type: "Festivalpass", FirstName: "Other", LastName: "Person", Email: "other@regncon.com", IsOver18: true},
		{ID: 105, OrderID: 5, TypeId: TicketTypeMiddag, Type: "Middag", FirstName: "Dinner", LastName: "Guest", Email: targetEmail, IsOver18: true},
	}
	db, logger := createCheckInTestDB(t)
	insertCheckInBillettholder(t, db, models.Billettholder{
		ID:           5000,
		FirstName:    "Existing",
		LastName:     "Target",
		TicketTypeId: 9000,
		TicketType:   "Festivalpass",
		IsOver18:     true,
		OrderID:      2,
		TicketID:     102,
	})
	insertManualBillettholderEmail(t, db, 5000, targetEmail)

	// When
	err := AssociateTicketsWithBillettholder(tickets, targetEmail, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected ticket association to succeed: %v", err)
	}
	actualBillettholderCount := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM billettholdere`)
	if actualBillettholderCount != expectedBillettholderCount {
		t.Fatalf("billettholder count mismatch\nexpected: %d\nactual:   %d", expectedBillettholderCount, actualBillettholderCount)
	}

	actualTargetEmailCount := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_billettholder_emails
		WHERE email = ? COLLATE NOCASE
	`, targetEmail)
	if actualTargetEmailCount != expectedTargetEmailCount {
		t.Fatalf("target email count mismatch\nexpected: %d\nactual:   %d", expectedTargetEmailCount, actualTargetEmailCount)
	}

	actualTicketIDs := queryBillettholderTicketIDs(t, db)
	if !slices.Equal(expectedTicketIDs, actualTicketIDs) {
		t.Fatalf("billettholder ticket IDs mismatch\nexpected: %v\nactual:   %v", expectedTicketIDs, actualTicketIDs)
	}
}

func queryBillettholderTicketIDs(t testing.TB, db *sql.DB) []int {
	t.Helper()

	rows, err := db.Query(`
		SELECT ticket_id
		FROM billettholdere
		ORDER BY ticket_id
	`)
	if err != nil {
		t.Fatalf("failed to query billettholder ticket IDs: %v", err)
	}
	defer rows.Close()

	var ticketIDs []int
	for rows.Next() {
		var ticketID int
		if err := rows.Scan(&ticketID); err != nil {
			t.Fatalf("failed to scan billettholder ticket ID: %v", err)
		}
		ticketIDs = append(ticketIDs, ticketID)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("failed to iterate billettholder ticket IDs: %v", err)
	}

	return ticketIDs
}
