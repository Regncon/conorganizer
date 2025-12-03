package checkIn

import (
	"testing"

	"github.com/google/uuid"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service"
	"github.com/Regncon/conorganizer/testutil"
	_ "modernc.org/sqlite"
)

func TestConvertTicketIdToNewBillettholder(t *testing.T) {
	// ❶ Arrange
	const ticketId = 42

	expectedBillettholder := models.Billettholder{
		FirstName:    "John",
		LastName:     "Doe",
		TicketTypeId: 1,
		TicketType:   "Adult",
		OrderID:      1,
		TicketID:     ticketId,
		IsOver18:     true,
	}
	expectedBillettholderEmails := []models.BillettholderEmail{
		{BillettholderID: expectedBillettholder.ID, Email: "ticket_email@test.test", Kind: "Ticket"},
		{BillettholderID: expectedBillettholder.ID, Email: "associated_email@test.test", Kind: "Associated"},
	}

	tickets := []CheckInTicket{
		{ID: ticketId,
			OrderID:   1,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "ticket_email@test.test",
			IsOver18:  true},
		{ID: 43,
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
			IsOver18:  false},
	}

	uniqueDatabaseName := "test_convert_ticket_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/tests/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom(testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// ❷ Act
	sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl)

	err = converTicketIdToNewBillettholder(ticketId, tickets, db, slogger)
	if err != nil {
		t.Fatalf("failed to convert ticketId to billettholder: %v", err)
	}

	// ❸ Assert
	var billettholder models.Billettholder
	err = db.QueryRow("SELECT * FROM billettholdere WHERE ticket_id = ?", ticketId).Scan(
		&billettholder.ID,
		&billettholder.FirstName,
		&billettholder.LastName,
		&billettholder.TicketTypeId,
		&billettholder.TicketType,
		&billettholder.IsOver18,
		&billettholder.OrderID,
		&billettholder.TicketID,
		&billettholder.InsertedTime,
	)

	if err != nil {
		t.Fatalf("failed to find billettholder with ticketId %d: %v", ticketId, err)
	}

	if billettholder.FirstName != expectedBillettholder.FirstName ||
		billettholder.LastName != expectedBillettholder.LastName ||
		billettholder.TicketTypeId != expectedBillettholder.TicketTypeId ||
		billettholder.TicketType != expectedBillettholder.TicketType ||
		billettholder.IsOver18 != expectedBillettholder.IsOver18 ||
		billettholder.OrderID != expectedBillettholder.OrderID ||
		billettholder.TicketID != expectedBillettholder.TicketID {
		t.Errorf("expected billettholder %+v, got %+v", expectedBillettholder, billettholder)
	}

	var billettholderEmails []models.BillettholderEmail
	rows, err := db.Query("SELECT id, billettholder_id, email, kind, inserted_time FROM billettholder_emails WHERE billettholder_id = ?", billettholder.ID)
	if err != nil {
		t.Fatalf("failed to query billettholder emails: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var email models.BillettholderEmail
		if err := rows.Scan(&email.ID, &email.BillettholderID, &email.Email, &email.Kind, &email.InsertedTime); err != nil {
			t.Fatalf("failed to scan billettholder email: %v", err)
		}
		billettholderEmails = append(billettholderEmails, email)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("error iterating over billettholder emails: %v", err)
	}
	if len(billettholderEmails) != len(expectedBillettholderEmails) {
		t.Fatalf("expected %d billettholder emails, got %d", len(expectedBillettholderEmails), len(billettholderEmails))
	}
	for i, email := range billettholderEmails {
		expectedEmail := expectedBillettholderEmails[i]
		if email.Email != expectedEmail.Email || email.Kind != expectedEmail.Kind {
			t.Errorf("expected billettholder email %+v kind %+v, got %+v kind %+v", expectedEmail.Email, expectedEmail.Kind, email.Email, email.Kind)
		}
	}
}

func TestDoNotConvertTicketsOfTypeMiddag(t *testing.T) {
	// ❶ Arrange
	expectedError := "cannot convert 'Middag' ticket to billettholder"

	ticketId := 42
	tickets := []CheckInTicket{
		{ID: ticketId,
			OrderID:   1,
			TypeId:    TicketTypeMiddag,
			Type:      "Middag",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "ticket_email@test.test",
			IsOver18:  true},
	}

	// ❷ Act
	sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl)

	err := converTicketIdToNewBillettholder(ticketId, tickets, nil, slogger)

	// ❸ Assert
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
	if err.Error() != expectedError {
		t.Errorf("expected error %q, got %q", expectedError, err.Error())
	}
}

func TestDontAddDuplicateAssociatedEmails(t *testing.T) {
	// ❶ Arrange
	const ticketId = 42

	expectedBillettholderEmails := []models.BillettholderEmail{
		{BillettholderID: 0, Email: "ticket_email@test.test", Kind: "Ticket"},
		{BillettholderID: 0, Email: "associated_email@test.test", Kind: "Associated"},
	}

	tickets := []CheckInTicket{
		{ID: ticketId,
			OrderID:   1,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "ticket_email@test.test",
			IsOver18:  true},
		{ID: 43,
			OrderID:   1,
			TypeId:    2,
			Type:      "Child",
			FirstName: "Jane",
			LastName:  "Doe",
			Email:     "associated_email@test.test",
			IsOver18:  false},
		{ID: 44,
			OrderID:   1,
			TypeId:    2,
			Type:      "Child",
			FirstName: "Same as previous",
			LastName:  "associated email",
			Email:     "associated_email@test.test",
			IsOver18:  false},
		{ID: 45,
			OrderID:   2,
			TypeId:    1,
			Type:      "Adult",
			FirstName: "Not",
			LastName:  "associated",
			Email:     "not_associated_email@test.test",
			IsOver18:  false},
	}

	uniqueDatabaseName := "test_convert_ticket_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/tests/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom(testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// ❷ Act
	sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl)

	err = converTicketIdToNewBillettholder(ticketId, tickets, db, slogger)
	if err != nil {
		t.Fatalf("failed to convert ticketId to billettholder: %v", err)
	}

	// ❸ Assert
	var billettholderEmails []models.BillettholderEmail
	rows, err := db.Query("SELECT id, billettholder_id, email, kind, inserted_time FROM billettholder_emails")
	if err != nil {
		t.Fatalf("failed to query billettholder emails: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var email models.BillettholderEmail
		if err := rows.Scan(&email.ID, &email.BillettholderID, &email.Email, &email.Kind, &email.InsertedTime); err != nil {
			t.Fatalf("failed to scan billettholder email: %v", err)
		}
		billettholderEmails = append(billettholderEmails, email)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("error iterating over billettholder emails: %v", err)
	}
	if len(billettholderEmails) != len(expectedBillettholderEmails) {
		t.Fatalf("expected %d billettholder emails, got %d", len(expectedBillettholderEmails), len(billettholderEmails))
	}
	for i, email := range billettholderEmails {
		expectedEmail := expectedBillettholderEmails[i]
		if email.Email != expectedEmail.Email || email.Kind != expectedEmail.Kind {
			t.Errorf("expected billettholder email %+v kind %+v, got %+v kind %+v", expectedEmail.Email, expectedEmail.Kind, email.Email, email.Kind)
		}
	}
}
