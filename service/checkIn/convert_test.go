package checkIn

import (
	"testing"

	"github.com/google/uuid"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service"
	"github.com/Regncon/conorganizer/testutil"
	_ "modernc.org/sqlite"
)

/*
	func TestConvertTicketIdToNewBillettholderError(t *testing.T) {
		// ❶ Arrange
		expectedError := "billettholder already exists"

		sl := &testutil.StubLogger{}

		tickets := []CheckInTicket{
			{ID: 42, OrderID: 1, Type: "Adult", Name: "John Doe", Email: "test@test.test", IsAdult: true},
		}
		uniqueDatabaseName := "test_convert_ticket_error_" + t.Name() + "_" + uuid.New().String() + ".db"
		db, err := service.InitDB("../../database/"+uniqueDatabaseName, "../../initialize.sql")
		if err != nil {
			t.Fatalf("failed to create test database: %v", err)
		}
		defer db.Close()

		// Insert a billettholder with ticketId 42 to simulate the error condition
		billettholder := models.Billettholder{
			TicketID: 42,
		}

		_, err = db.Exec(`
			INSERT INTO billettholdere (
	        first_name, last_name, ticket_type,
	        ticket_id, is_over_18, order_id,
	        ticket_email, order_email, ticket_category_id
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			billettholder.FirstName, billettholder.LastName, billettholder.TicketType,
			billettholder.TicketID, billettholder.IsOver18, billettholder.OrderID,
			billettholder.TicketEmail, billettholder.OrderEmail, billettholder.TicketCategoryID,
		)
		if err != nil {
			t.Fatalf("failed to insert billettholder for test: %v", err)
		}

		// ❷ Act
		slogger := testutil.NewSlogAdapter(sl)
		err = converTicketIdToNewBillettholder(42, tickets, db, slogger)
		if err == nil {
			t.Fatalf("expected error but got nil")
		}
		if err.Error() != expectedError {
			t.Errorf("expected error %q, got %q", expectedError, err.Error())
		}

}
*/
func TestConvertTicketIdToNewBillettholder(t *testing.T) {
	// ❶ Arrange
	ticketId := 42

	expectedBillettholder := models.Billettholder{
		FirstName: "John",
		LastName:  "Doe",
		OrderID:   1,
		TicketID:  ticketId,
		IsOver18:  true,
	}
	expectedBillettholderEmails := []models.BillettholderEmail{
		{BillettholderID: expectedBillettholder.ID, Email: "ticket_email@test.test", Kind: "Ticket"},
		{BillettholderID: expectedBillettholder.ID, Email: "associated_email@test.test", Kind: "Associated"},
	}

	/*
	   if (!ticket) throw new Error('ticket not found');

	   	const isOver18 = new Date().getFullYear() - new Date(ticket.crm.born).getFullYear() > 18;

	   	const orderEmails = tickets.filter((t) => t.order_id === ticket.order_id).map((t) => t.crm.email);

	   	let participant: Participant = {
	   	    firstName: ticket.crm.first_name,
	   	    lastName: ticket.crm.last_name,
	   	    over18: isOver18,
	   	    ticketEmail: ticket.crm.email,
	   	    orderEmails: orderEmails,
	   	    ticketId: ticket.id,
	   	    orderId: ticket.order_id,
	   	    ticketCategory: ticket.category,
	   	    ticketCategoryId: ticket.category_id,
	   	    createdAt: new Date().toISOString(),
	   	    createdBy: userId,
	   	    updateAt: new Date().toISOString(),
	   	    updatedBy: userId,
	   	    connectedEmails: [],
	   	};
	*/
	tickets := []CheckInTicket{
		{ID: ticketId, OrderID: 1, Type: "Adult", FirstName: "John", LastName: "Doe", Email: "ticket_email@test.test", IsOver18: true},
		{ID: 43, OrderID: 1, Type: "Child", FirstName: "Jane", LastName: "Doe", Email: "associated_email@test.test", IsOver18: false},
		{ID: 44, OrderID: 2, Type: "Adult", FirstName: "Not", LastName: "associated", Email: "not_associated_email@test.test", IsOver18: false},
	}

	uniqueDatabaseName := "test_convert_ticket_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom("../../database/events.db", testDBPath)
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
		&billettholder.TicketType,
		&billettholder.IsOver18,
		&billettholder.OrderID,
		&billettholder.TicketID,
		&billettholder.TicketCategoryID,
		&billettholder.InsertedTime,
	)

	if err != nil {
		t.Fatalf("failed to find billettholder with ticketId %d: %v", ticketId, err)
	}

	if billettholder.FirstName != expectedBillettholder.FirstName ||
		billettholder.LastName != expectedBillettholder.LastName ||
		billettholder.TicketType != expectedBillettholder.TicketType ||
		billettholder.IsOver18 != expectedBillettholder.IsOver18 ||
		billettholder.OrderID != expectedBillettholder.OrderID ||
		billettholder.TicketID != expectedBillettholder.TicketID ||
		billettholder.TicketCategoryID != expectedBillettholder.TicketCategoryID {
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
