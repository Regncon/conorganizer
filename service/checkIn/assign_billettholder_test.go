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

func TestAssociateTicketsWithBillettholder(t *testing.T) {
	// Arrange
	sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl)

	uniqueDatabaseName := "test_associate_tickets_" + t.Name() + "_" + uuid.New().String() + ".db"
	testDBPath := "../../database/tests/" + uniqueDatabaseName

	db, err := service.InitTestDBFrom("../../database/events.db", testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// Test variables
	const targetEmail = "test@regncon.com"
	const fakePeopleAmount = 10
	const billettholderConversionRatio = 0.2

	// Generate test data
	generatededPeople := testutil.GeneratePeople(fakePeopleAmount)
	var generatededTickets []CheckInTicket
	for i, generatedPerson := range generatededPeople {
		// Tie 10% of tickets with our target email
		var emailValue = targetEmail
		if rand.Intn(10) > 1 {
			emailValue = generatedPerson.Email
		}

		generatededTickets = append(generatededTickets, CheckInTicket{
			ID:        i + 1,
			OrderID:   i + 1,
			TypeId:    i + 1,
			FirstName: generatedPerson.FirstName,
			LastName:  generatedPerson.LastName,
			Type:      "Test billet",
			Email:     emailValue,
			IsOver18:  rand.Intn(10) > 2,
		})
	}

	// Splice some generated tickets and write them to billettholders table
	var conversionAmmount = fakePeopleAmount * billettholderConversionRatio
	billettholderConversion := generatededTickets[:int(conversionAmmount)]
	for _, ticket := range billettholderConversion {
		err = converTicketIdToNewBillettholder(ticket.ID, billettholderConversion, db, slogger)
		if err != nil {
			fmt.Println(err)
		}
	}

	// generate some fake users
	var expectedUsers []models.User
	for i, holder := range generatededTickets {
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
	queryBase := fmt.Sprintf(`INSERT INTO users (id, user_id, email, is_admin) VALUES %s`, strings.Join(queryUsers, ", "))

	_, err = db.Exec(queryBase)
	if err != nil {
		fmt.Println("failed to insert users", "error", err)
		return
	}

	// Generate some fake billettholder_users
	for _, expectedUser := range expectedUsers {
		err = AssociateUserWithBillettholder(expectedUser.UserID, db, slogger)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Act
	err = AssociateTicketsWithBillettholder(generatededTickets, targetEmail)
	if err != nil {
		t.Fatalf("failed to associate ticket with billettholder: %v", err)
	}

	// Assert
}
