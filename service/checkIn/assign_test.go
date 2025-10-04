package checkIn

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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
	var expectedGeneratedUsers []testutil.Person
	var expectedBillettholders []models.Billettholder

	var ammountToGenerate = 99
	for i := range ammountToGenerate {
		// Create fake users which billettholders can access (this is done to retain email)
		var generatedPerson = testutil.GenerateFakePerson()
		expectedGeneratedUsers = append(expectedGeneratedUsers, generatedPerson)

		// Create billettholders
		expectedBillettholders = append(expectedBillettholders, models.Billettholder{
			ID:           i + 1,
			FirstName:    generatedPerson.FirstName,
			LastName:     generatedPerson.LastName,
			TicketTypeId: i + 190001,
			TicketType:   "Test",
			IsOver18:     rand.Intn(100) > 20,
			OrderID:      i + 13200001,
			TicketID:     i + 4800001,
		})
	}

	// construct query for inserting billettholdere
	var queryBillettholder []string
	for _, billettholder := range expectedBillettholders {
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

	/* Attempt to insert into billettholder_emails */
	var expectedBillettholderEmails []models.BillettholderEmail
	for _, person := range expectedGeneratedUsers {
		billettholderEmail := models.BillettholderEmail{
			BillettholderID: rand.Intn(len(expectedGeneratedUsers)-1) + 1,
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

	/* Attempt to insert into expectedUsers */
	var expectedUsers []models.User
	for i, holder := range expectedGeneratedUsers {
		expectedUsers = append(expectedUsers, models.User{
			ID:      i + 1,
			UserID:  holder.FirstName + strconv.Itoa(i+1),
			Email:   holder.Email,
			IsAdmin: rand.Intn(100) > 10,
		})
	}

	var queryUsers []string
	for _, user := range expectedUsers {
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

	for _, user := range expectedUsers {
		// fmt.Printf("Calling AssociateUserWithBillettholder() on: %s (%s)\n", user.UserID, user.Email)
		err = AssociateUserWithBillettholder(user.UserID, db, slogger)
		if err != nil {
			t.Fatalf("failed to convert ticketId to billettholder: %v", err)
		}
	}

	// Assert
	var expectedBillettholderUsers []models.BillettholderUsers
	for _, expectedUser := range expectedUsers {
		for _, generatedBilletholder := range expectedBillettholderEmails {
			if expectedUser.Email == generatedBilletholder.Email {
				expectedBillettholderUsers = append(expectedBillettholderUsers, models.BillettholderUsers{
					BillettholderID: generatedBilletholder.BillettholderID,
					UserID:          expectedUser.UserID,
				})
			}
		}
	}

	/* Check that billettholder_users got populated */
	var billettholderUsers []models.BillettholderUsers

	rows, err := db.Query(`
        SELECT billettholder_id, user_id FROM billettholdere_users
    `)
	if err != nil {
		t.Fatalf("Failed to get rows from billettholder_users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var result models.BillettholderUsers
		err = rows.Scan(&result.BillettholderID, &result.UserID)
		if err != nil {
			t.Fatalf("Failed scan rows in billettholder_users: %v", err)
		}
		billettholderUsers = append(billettholderUsers, result)
	}

	// compare expected billettholderUsers with billettholderUsers
	if len(billettholderUsers) != len(expectedBillettholderUsers) {
		t.Fatalf("expected %d billettholder users, got %d", len(expectedBillettholderUsers), len(billettholderUsers))
	}

	// Kan det eksistere billettholder_emails med duplicate e-post?
}
