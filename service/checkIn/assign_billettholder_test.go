package checkIn

import (
	"testing"

	"github.com/Regncon/conorganizer/service"
	"github.com/google/uuid"
)

func TestAssociateTicketsWithBillettholder(t *testing.T) {
	// Arrange
	uniqueDatabaseName := "test_associate_billettholders_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/tests/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom("../../database/events.db", testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	const targetEmail = "test@regncon.com"

	// Generate some fake people
	// Generate some fake tickets
	// Slize some tickets and call convert() to make tickets into billettholders

	// generate some fake users
	// generate some fake billettholder_emails

	// Generate some fake billettholder_users

	// Act
	err = AssociateTicketsWithBillettholder(generatedTickets, targetEmail)
	if err != nil {
		t.Fatalf("failed to associate email with billettholder: %v", err)
	}

	// Assert
}
