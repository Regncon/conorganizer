package checkIn

import (
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestConvertTicketIdToNewBillettholder_CreatesBillettholderAndEmails(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given one ticket in an order with a second ticket using another email.",
		When:  "When the ticket is converted to a billettholder.",
		Then:  "Then the billettholder stores the ticket data and both related emails.",
	})

	// Given
	const ticketID = 42
	expectedBillettholder := models.Billettholder{
		FirstName:    "John",
		LastName:     "Doe",
		TicketTypeId: 1,
		TicketType:   "Adult",
		OrderID:      1,
		TicketID:     ticketID,
		IsOver18:     true,
	}
	expectedEmails := []models.BillettholderEmail{
		{Email: "ticket_email@test.test", Kind: models.BillettholderEmailKindTicket},
		{Email: "associated_email@test.test", Kind: models.BillettholderEmailKindAssociated},
	}
	tickets := []CheckInTicket{
		{
			ID:        ticketID,
			OrderID:   1,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "ticket_email@test.test",
			IsOver18:  true,
		},
		{
			ID:        43,
			OrderID:   1,
			TypeId:    2,
			Type:      "Child",
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "associated_email@test.test",
			IsOver18:  false,
		},
		{
			ID:        44,
			OrderID:   2,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "Not",
			LastName:  "Associated",
			Email:     "not_associated_email@test.test",
			IsOver18:  false,
		},
	}
	db, logger := createCheckInTestDB(t)

	// When
	err := converTicketIdToNewBillettholder(ticketID, tickets, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected ticket conversion to succeed: %v", err)
	}
	actualBillettholder := queryBillettholderByTicketID(t, db, ticketID)
	assertBillettholderMatches(t, expectedBillettholder, actualBillettholder)
	assertBillettholderEmails(t, db, actualBillettholder.ID, expectedEmails)
}

func TestConvertTicketIdToNewBillettholder_WhenTicketIsDinner_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a dinner ticket.",
		When:  "When it is converted to a billettholder.",
		Then:  "Then conversion is rejected before database writes are attempted.",
	})

	// Given
	const ticketID = 42
	expectedError := "cannot convert 'Middag' ticket to billettholder"
	tickets := []CheckInTicket{
		{
			ID:        ticketID,
			OrderID:   1,
			TypeId:    TicketTypeMiddag,
			Type:      "Middag",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "ticket_email@test.test",
			IsOver18:  true,
		},
	}

	// When
	err := converTicketIdToNewBillettholder(ticketID, tickets, nil, testutil.NewTestLogger())

	// Then
	if err == nil {
		t.Fatalf("expected dinner ticket conversion to fail")
	}
	if err.Error() != expectedError {
		t.Fatalf("error mismatch\nexpected: %q\nactual:   %q", expectedError, err.Error())
	}
}

func TestConvertTicketIdToNewBillettholder_WhenAssociatedEmailsRepeat_InsertsEachEmailOnce(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given multiple tickets in the same order with the same associated email.",
		When:  "When the main ticket is converted to a billettholder.",
		Then:  "Then each billettholder email is stored once.",
	})

	// Given
	const ticketID = 42
	expectedEmails := []models.BillettholderEmail{
		{Email: "ticket_email@test.test", Kind: models.BillettholderEmailKindTicket},
		{Email: "associated_email@test.test", Kind: models.BillettholderEmailKindAssociated},
	}
	tickets := []CheckInTicket{
		{
			ID:        ticketID,
			OrderID:   1,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "ticket_email@test.test",
			IsOver18:  true,
		},
		{
			ID:        43,
			OrderID:   1,
			TypeId:    2,
			Type:      "Child",
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "associated_email@test.test",
			IsOver18:  false,
		},
		{
			ID:        44,
			OrderID:   1,
			TypeId:    2,
			Type:      "Child",
			FirstName: "Same",
			LastName:  "Associated",
			Email:     "associated_email@test.test",
			IsOver18:  false,
		},
		{
			ID:        45,
			OrderID:   2,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "Not",
			LastName:  "Associated",
			Email:     "not_associated_email@test.test",
			IsOver18:  false,
		},
	}
	db, logger := createCheckInTestDB(t)

	// When
	err := converTicketIdToNewBillettholder(ticketID, tickets, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected ticket conversion to succeed: %v", err)
	}
	actualBillettholder := queryBillettholderByTicketID(t, db, ticketID)
	assertBillettholderEmails(t, db, actualBillettholder.ID, expectedEmails)
}

func assertBillettholderMatches(t testing.TB, expected models.Billettholder, actual models.Billettholder) {
	t.Helper()

	if actual.FirstName != expected.FirstName ||
		actual.LastName != expected.LastName ||
		actual.TicketTypeId != expected.TicketTypeId ||
		actual.TicketType != expected.TicketType ||
		actual.IsOver18 != expected.IsOver18 ||
		actual.OrderID != expected.OrderID ||
		actual.TicketID != expected.TicketID {
		t.Fatalf("billettholder mismatch\nexpected: %+v\nactual:   %+v", expected, actual)
	}
}
