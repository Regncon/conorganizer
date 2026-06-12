package checkIn

import (
	"slices"
	"testing"
)

func TestAssociateTicketsWithEmail_WhenEmailMatchesCaseInsensitively_ReturnsNonDinnerTickets(t *testing.T) {
	// Given tickets with matching, differently-cased, dinner, and unrelated emails,
	// when tickets are associated by email,
	// then only non-dinner tickets matching the email are returned.

	// Given
	expectedTicketIDs := []int{101, 102}
	tickets := []CheckInTicket{
		{ID: 101, TypeId: 1, Email: "test@regncon.no"},
		{ID: 102, TypeId: 2, Email: "Test@Regncon.no"},
		{ID: 103, TypeId: TicketTypeMiddag, Email: "test@regncon.no"},
		{ID: 104, TypeId: 1, Email: "other@regncon.no"},
	}

	// When
	actualTickets, err := AssociateTicketsWithEmail(tickets, "TEST@REGNCON.NO")

	// Then
	if err != nil {
		t.Fatalf("expected ticket association to succeed: %v", err)
	}
	actualTicketIDs := ticketIDs(actualTickets)
	if !slices.Equal(expectedTicketIDs, actualTicketIDs) {
		t.Fatalf("ticket IDs mismatch\nexpected: %v\nactual:   %v", expectedTicketIDs, actualTicketIDs)
	}
}

func TestAssociateTicketsWithEmail_WhenNoTicketMatches_ReturnsError(t *testing.T) {
	// Given tickets registered to other emails,
	// when tickets are associated by an unknown email,
	// then the caller receives an error.

	// Given
	expectedError := true
	tickets := []CheckInTicket{
		{ID: 101, TypeId: 1, Email: "other@regncon.no"},
	}

	// When
	_, err := AssociateTicketsWithEmail(tickets, "missing@regncon.no")
	actualError := err != nil

	// Then
	if actualError != expectedError {
		t.Fatalf("error presence mismatch\nexpected: %v\nactual:   %v", expectedError, actualError)
	}
}

func ticketIDs(tickets []CheckInTicket) []int {
	ids := make([]int, 0, len(tickets))
	for _, ticket := range tickets {
		ids = append(ids, ticket.ID)
	}
	return ids
}
