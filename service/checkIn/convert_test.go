package checkIn

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"testing"
)

// ————————————————————————————
// 1. Define the abstraction used by production code
// ————————————————————————————
type Logger interface {
	Info(msg string, keysAndValues ...interface{})
}

// ————————————————————————————
// 2. A lightweight stub that records calls
// ————————————————————————————
type stubLogger struct {
	calls []struct {
		msg           string
		keysAndValues []interface{}
	}
}

func (s *stubLogger) Info(msg string, keysAndValues ...interface{}) {
	s.calls = append(s.calls, struct {
		msg           string
		keysAndValues []interface{}
	}{msg, keysAndValues})
}

// ————————————————————————————
// 3. (Example) production function that needs a logger
// ————————————————————————————
func convertTicketIdToNewBilettholder(ticketID int, db *sql.DB, log Logger) {
	log.Info("Converting ticket to bilettholder", "ticketID", ticketID)
	// …real work goes here…
}

// ————————————————————————————
// 4. Unit test with only the standard library
// ————————————————————————————
func TestConvertTicketIdToNewBilettholder(t *testing.T) {
	// ❶ Arrange
	sl := &stubLogger{}

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
	db, err := sql.Open("sqlite", ":memory:") // requires the sqlite driver at build-time only
	if err != nil {
		t.Fatalf("opening in-memory DB: %v", err)
	}
	defer db.Close()

	// ❷ Act
	convertTicketIdToNewBilettholder(42, db, sl)

	// ❸ Assert
	if got := len(sl.calls); got != 1 {
		t.Fatalf("expected 1 Info call, got %d", got)
	}

	call := sl.calls[0]
	if call.msg != "Converting ticket to bilettholder" {
		t.Errorf("unexpected log message: %q", call.msg)
	}
	if len(call.keysAndValues) != 2 || call.keysAndValues[0] != "ticketID" || call.keysAndValues[1] != 42 {
		t.Errorf("unexpected key/values: %#v", call.keysAndValues)
	}
}
