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
	sl := &testutil.StubLogger{}

	// expectedBillettholder := models.Billettholder{
	// 	TicketID: 42,
	// }

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
		{ID: 42, OrderID: 1, Type: "Adult", Name: "John Doe", Email: "test@test.test", IsAdult: true},
		{ID: 43, OrderID: 1, Type: "Child", Name: "Jane Doe", Email: "test2@test.test", IsAdult: false},
	}

	uniqueDatabaseName := "test_convert_ticket_" + t.Name() + "_" + uuid.New().String() + ".db"
	db, err := service.InitDB("../../database/"+uniqueDatabaseName, "../../initialize.sql")
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// ❷ Act
	slogger := testutil.NewSlogAdapter(sl)

	err = converTicketIdToNewBillettholder(42, tickets, db, slogger)
	if err != nil {
		t.Fatalf("failed to convert ticketId to billettholder: %v", err)
	}

	// ❸ Assert
	var billettholder models.Billettholder
	err = db.QueryRow("SELECT id, ticket_id FROM billettholdere WHERE ticket_id = ?", 42).Scan(
		// err = db.QueryRow("SELECT id, ticket_id FROM billettholdere").Scan(
		&billettholder.ID,
		&billettholder.TicketID,
	)

	if err != nil {
		t.Fatalf("failed to find billettholder with ticketId 42: %v", err)
	}

	if billettholder.TicketID != 42 {
		t.Errorf("expected billettholder with ticketId 42, got %d", billettholder.TicketID)
	}
}
