package checkIn

import (
	"fmt"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/google/uuid"
)

func GetAllBillettHolderByUserEmail() ([]int, error) {
	// hent ut alle billetter som matcher userID.email fra billettholder_emails, returnerer billettholder(id) array
	return []int{}, nil

}

func InsertBilletHolderIDSFromUserEmail() error {
	return nil
}

/* func GetTicketsFromEmail() (CheckInTicket, error) {

return nil, nil
} */

func TestAssociateUserWithBillettholder(t *testing.T) {
	// Arrange

	/* Create temp db for testing */
	uniqueDatabaseName := "test_associate_billettholders_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/tests/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom("../../database/events.db", testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	/* Add billettholdere test data */
	expectedBillettholders := []models.BillettholderUsers{
		{BillettholderID: 1, UserID: "1"},
		{BillettholderID: 2, UserID: "2"},
	}

	_, err = db.Exec(`
    INSERT INTO billettholdere (
        id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, "Ola", "Nordmann", 1, "Test", 1, 1, 1), (2, "kari", "NordVPN", 2, "test", 0, 2, 2)`,
		expectedBillettholders[0].BillettholderID, expectedBillettholders[1].BillettholderID,
	)
	if err != nil {
		fmt.Println("failed to insert billettholder", "error", err)
		return
	}

	const email = "test@regncon.no"
	billettholderEmails := []models.BillettholderEmail{
		{
			BillettholderID: expectedBillettholders[0].BillettholderID,
			Email:           email,
		},
	}

	/* Attempt to insert into billettholder_emails */
	_, err = db.Exec(`
		INSERT INTO billettholder_emails (
        billettholder_id, email, kind
		) VALUES (?, ?, "Manual")`,
		billettholderEmails[0].BillettholderID, billettholderEmails[0].Email,
	)
	if err != nil {
		fmt.Println("failed to insert billettholder_emails", "error", err)
		return
	}

	user := models.User{
		ID:      1,
		UserID:  expectedBillettholders[0].UserID,
		Email:   email,
		IsAdmin: true,
	}

	/* Attempt to insert into users*/
	_, err = db.Exec(`
		INSERT INTO users (
        id, user_id, email, is_admin
		) VALUES (?,?,?,?)`,
		user.ID, user.UserID, user.Email, user.IsAdmin,
	)
	if err != nil {
		fmt.Println("failed to insert users", "error", err)
		return
	}

	// Act
	sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl)

	err = AssociateUserWithBillettholder(expectedBillettholders[0].UserID, db, slogger)
	if err != nil {
		t.Fatalf("failed to convert ticketId to billettholder: %v", err)
	}

	// Assert
	// query billettholder_users
	// valider data mot arrange
}
