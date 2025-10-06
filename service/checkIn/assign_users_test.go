package checkIn

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/google/uuid"
)

func TestAssociateUserWithBillettholder(t *testing.T) {
	// Arrange
	uniqueDatabaseName := "test_associate_billettholders_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/tests/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom("../../database/events.db", testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// expected billettholder_user
	var expectedBillettholderUser = models.BillettholderUsers{
		BillettholderID: 9999,
		UserID:          1,
	}

	// Happy path person
	const happyPathEmail = "test@regncon.no"
	var happyPathPerson = testutil.GenerateFakePerson()
	happyPathPerson.Email = happyPathEmail

	// Happy path user
	var happyPathUser = models.User{
		ID:      expectedBillettholderUser.UserID,
		UserID:  "testuser",
		Email:   happyPathEmail,
		IsAdmin: false,
	}

	// Happy path billettholder
	var happyPathBillettholder = models.Billettholder{
		ID:           expectedBillettholderUser.BillettholderID,
		FirstName:    happyPathPerson.FirstName,
		LastName:     happyPathPerson.LastName,
		TicketTypeId: 199999,
		TicketType:   "Test",
		IsOver18:     true,
		OrderID:      19999999,
		TicketID:     4999999,
	}
	// Happy path billettholder_email
	var happyPathBillettholderEmail = models.BillettholderEmail{
		BillettholderID: happyPathBillettholder.ID,
		Email:           happyPathPerson.Email,
		Kind:            "Manual",
	}

	var testBillettholders []models.Billettholder
	testBillettholders = append(testBillettholders, happyPathBillettholder)

	var testBillettholderEmails []models.BillettholderEmail
	testBillettholderEmails = append(testBillettholderEmails, happyPathBillettholderEmail)

	var testUsers []models.User
	testUsers = append(testUsers, happyPathUser)

	// construct query for inserting billettholdere
	var queryBillettholder []string
	for _, billettholder := range testBillettholders {
		queryBillettholder = append(queryBillettholder, fmt.Sprintf(`(%d, "%s", "%s", %d, "%s", %v, %d, %d)`, billettholder.ID, billettholder.FirstName, billettholder.LastName, billettholder.TicketTypeId, billettholder.TicketType, billettholder.IsOver18, billettholder.OrderID, billettholder.TicketID))
	}

	var queryBase = fmt.Sprintf(`INSERT INTO billettholdere (
            id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
            ) VALUES %s`, strings.Join(queryBillettholder, ", "))

	_, err = db.Exec(queryBase)
	if err != nil {
		fmt.Println("failed to insert billettholder", "error", err)
		return
	}

	// Attempt to insert into billettholder_emails
	var expectedBillettholderEmails []models.BillettholderEmail
	for _, person := range testBillettholderEmails {
		billettholderEmail := models.BillettholderEmail{
			BillettholderID: person.BillettholderID,
			Email:           person.Email,
		}
		expectedBillettholderEmails = append(expectedBillettholderEmails, billettholderEmail)
	}

	var queryBillettholderEmail []string
	for _, billettholderEmail := range expectedBillettholderEmails {
		queryBillettholderEmail = append(queryBillettholderEmail, fmt.Sprintf(`(%d, "%s", "%s")`, billettholderEmail.BillettholderID, billettholderEmail.Email, "Manual"))
	}
	queryBase = fmt.Sprintf(`
            INSERT INTO billettholder_emails (
                billettholder_id, email, kind
                ) VALUES %s`, strings.Join(queryBillettholderEmail, ", "))

	_, err = db.Exec(queryBase)
	if err != nil {
		fmt.Println("failed to insert billettholder_emails", "error", err)
		return
	}

	var queryUsers []string
	for _, user := range testUsers {
		queryUsers = append(queryUsers, fmt.Sprintf(`(%d, "%s", "%s", %v)`, user.ID, user.UserID, user.Email, user.IsAdmin))
	}

	queryBase = fmt.Sprintf(`
                INSERT INTO users (
                    id, user_id, email, is_admin
                    ) VALUES %s`, strings.Join(queryUsers, ", "))

	_, err = db.Exec(queryBase)
	if err != nil {
		fmt.Println("failed to insert users", "error", err)
		return
	}

	// Act
	sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl)

	for _, user := range testUsers {
		// fmt.Printf("Calling AssociateUserWithBillettholder() on: %s (%s)\n", user.UserID, user.Email)
		err = AssociateUserWithBillettholder(user.UserID, db, slogger)
		if err != nil {
			t.Fatalf("failed to convert ticketId to billettholder: %v", err)
		}
	}

	// Assert

	var resultBillettholderUsers []models.BillettholderUsers
	rows, err := db.Query("SELECT billettholder_id, user_id FROM billettholdere_users")
	if err != nil {
		t.Fatalf("failed to query billettholder_users: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var user models.BillettholderUsers
		if err := rows.Scan(&user.BillettholderID, &user.UserID); err != nil {
			t.Fatalf("failed to scan billettholder_user: %v", err)
		}
		resultBillettholderUsers = append(resultBillettholderUsers, user)
	}
	if len(resultBillettholderUsers) != 1 {
		t.Fatalf("expected 1 billettholder_user, got %d", len(resultBillettholderUsers))
	}
	if resultBillettholderUsers[0] != expectedBillettholderUser {
		t.Fatalf("expected billettholder_user %+v, got %+v", expectedBillettholderUser, resultBillettholderUsers[0])
	}
}
