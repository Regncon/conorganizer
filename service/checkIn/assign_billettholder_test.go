package checkIn

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestAssociateTicketsWithBillettholder_WhenSomeMatchingTicketsAreNew_ConvertsOnlyNewTickets(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given tickets for a target email where one ticket is already converted.",
		When:  "When the target email is associated with billettholdere.",
		Then:  "Then only new non-dinner target tickets are inserted.",
	})

	// Given
	expectedBillettholderCount := 3
	expectedCreatedBillettholders := 2
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
	actualResult, err := AssociateTicketsWithBillettholder(tickets, targetEmail, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected ticket association to succeed: %v", err)
	}
	if actualResult.CreatedBillettholders != expectedCreatedBillettholders {
		t.Fatalf("created billettholder count mismatch\nexpected: %d\nactual:   %d", expectedCreatedBillettholders, actualResult.CreatedBillettholders)
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

func TestAssociateTicketsWithBillettholder_WhenNoTicketsMatch_ReturnsNoCreatedBillettholders(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given CheckIn tickets registered to other emails.",
		When:  "When the target email is associated with billettholdere.",
		Then:  "Then no billettholdere are inserted and the operation is neutral.",
	})

	// Given
	expectedCreatedBillettholders := 0
	targetEmail := "test@regncon.com"
	tickets := []CheckInTicket{
		{ID: 101, OrderID: 1, TypeId: 9000, Type: "Festivalpass", FirstName: "Other", LastName: "Person", Email: "other@regncon.com", IsOver18: true},
	}
	db, logger := createCheckInTestDB(t)

	// When
	actualResult, err := AssociateTicketsWithBillettholder(tickets, targetEmail, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected neutral ticket association to succeed: %v", err)
	}
	if actualResult.CreatedBillettholders != expectedCreatedBillettholders {
		t.Fatalf("created billettholder count mismatch\nexpected: %d\nactual:   %d", expectedCreatedBillettholders, actualResult.CreatedBillettholders)
	}
	actualBillettholderCount := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM billettholdere`)
	if actualBillettholderCount != 0 {
		t.Fatalf("expected no billettholdere to be inserted, got %d", actualBillettholderCount)
	}
}

func TestAssociateTicketsWithBillettholder_WhenMatchingTicketSharesOrderWithAnotherEmail_ImportsAndAssociatesEntireOrder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a user whose email matches one non-dinner ticket in a multi-email order.",
		When:  "When that user's tickets are imported and associated.",
		Then:  "Then every non-dinner ticket in that order is stored and linked to the user.",
	})

	// Given
	expectedCreatedBillettholders := 2
	expectedCreatedAssociations := 2
	expectedTicketIDs := []int{101, 102}
	userEmail := "user-a@example.com"
	tickets := []CheckInTicket{
		{ID: 101, OrderID: 10, TypeId: 9000, Type: "Festivalpass", FirstName: "User", LastName: "A", Email: userEmail, IsOver18: true},
		{ID: 102, OrderID: 10, TypeId: 9000, Type: "Festivalpass", FirstName: "User", LastName: "B", Email: "user-b@example.com", IsOver18: true},
		{ID: 103, OrderID: 10, TypeId: TicketTypeMiddag, Type: "Middag", FirstName: "Dinner", LastName: "Guest", Email: "user-b@example.com", IsOver18: true},
		{ID: 104, OrderID: 11, TypeId: 9000, Type: "Festivalpass", FirstName: "Unrelated", LastName: "Person", Email: "other@example.com", IsOver18: true},
	}
	db, logger := createCheckInTestDB(t)
	insertUser(t, db, 1, "user-a", userEmail)

	// When
	associationResult, err := AssociateTicketsWithBillettholder(tickets, userEmail, db, logger)
	createdAssociations, associationErr := AssociateUserWithBillettholder("user-a", db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected ticket import to succeed: %v", err)
	}
	if associationErr != nil {
		t.Fatalf("expected user association to succeed: %v", associationErr)
	}
	if associationResult.CreatedBillettholders != expectedCreatedBillettholders {
		t.Fatalf("created billettholder count mismatch\nexpected: %d\nactual:   %d", expectedCreatedBillettholders, associationResult.CreatedBillettholders)
	}
	if createdAssociations != expectedCreatedAssociations {
		t.Fatalf("created user association count mismatch\nexpected: %d\nactual:   %d", expectedCreatedAssociations, createdAssociations)
	}
	actualTicketIDs := queryBillettholderTicketIDs(t, db)
	if !slices.Equal(expectedTicketIDs, actualTicketIDs) {
		t.Fatalf("billettholder ticket IDs mismatch\nexpected: %v\nactual:   %v", expectedTicketIDs, actualTicketIDs)
	}
	actualAssociatedTicketIDs := queryUserBillettholderTicketIDs(t, db, 1)
	if !slices.Equal(expectedTicketIDs, actualAssociatedTicketIDs) {
		t.Fatalf("associated billettholder ticket IDs mismatch\nexpected: %v\nactual:   %v", expectedTicketIDs, actualAssociatedTicketIDs)
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

func queryUserBillettholderTicketIDs(t testing.TB, db *sql.DB, userID int) []int {
	t.Helper()

	rows, err := db.Query(`
		SELECT b.ticket_id
		FROM billettholdere AS b
		JOIN relation_billettholdere_users AS bu ON bu.billettholder_id = b.id
		WHERE bu.user_id = ?
		ORDER BY b.ticket_id
	`, userID)
	if err != nil {
		t.Fatalf("failed to query user's billettholder ticket IDs: %v", err)
	}
	defer rows.Close()

	var ticketIDs []int
	for rows.Next() {
		var ticketID int
		if err := rows.Scan(&ticketID); err != nil {
			t.Fatalf("failed to scan user's billettholder ticket ID: %v", err)
		}
		ticketIDs = append(ticketIDs, ticketID)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("failed to iterate user's billettholder ticket IDs: %v", err)
	}

	return ticketIDs
}
