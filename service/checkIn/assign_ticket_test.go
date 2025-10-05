package checkIn

import (
	"testing"
)

func TestAssociateBillettholderWithEmail(t *testing.T) {
	// Arrange
	const targetEmail = "ticket_email@test.test"

	tickets := []CheckInTicket{
		{ID: 42,
			OrderID:   1,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "ticket_email@test.test",
			IsOver18:  true},
		/* {ID: 43,
			OrderID:   1,
			TypeId:    2,
			Type:      "Child",
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "associated_email@test.test",
			IsOver18:  false},
		{ID: 44,
			OrderID:   2,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "Not",
			LastName:  "associated",
			Email:     "not_associated_email@test.test",
			IsOver18:  false}, */
	}

	var matches []CheckInTicket
	matches = append(matches, tickets[0])

	// Act
	/* sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl) */

	result, err := AssociateBillettholderWithEmail(tickets, targetEmail)
	if err != nil {
		t.Fatalf("failed to associate email with billettholder: %v", err)
	}

	// Assert
	if len(result) != len(matches) {
		t.Fatalf("expected %d tickets, got %d", len(matches), len(result))
	}
}
