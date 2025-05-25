package checkIn

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service"
	_ "modernc.org/sqlite"
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

// 3. Adapter to make stubLogger compatible with *slog.Logger
// ————————————————————————————
type stubLoggerHandler struct {
	stub *stubLogger
}

func (h *stubLoggerHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (h *stubLoggerHandler) Handle(ctx context.Context, r slog.Record) error {
	// Extract message
	msg := r.Message

	// Extract key-value pairs
	var keyValues []interface{}
	r.Attrs(func(attr slog.Attr) bool {
		keyValues = append(keyValues, attr.Key, attr.Value.Any())
		return true
	})

	// Forward to stub logger
	h.stub.Info(msg, keyValues...)
	return nil
}

func (h *stubLoggerHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *stubLoggerHandler) WithGroup(name string) slog.Handler {
	return h
}

func newSlogAdapter(stub *stubLogger) *slog.Logger {
	handler := &stubLoggerHandler{stub: stub}
	return slog.New(handler)
}

// ————————————————————————————
// 4. Unit test with only the standard library
// ————————————————————————————
func TestConvertTicketIdToNewBillettholder(t *testing.T) {
	// ❶ Arrange
	sl := &stubLogger{}

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
	slogger := newSlogAdapter(sl)

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
