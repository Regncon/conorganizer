package checkIn

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func TestAssociateTicketsWithBillettholder(t *testing.T) {
	// Arrange
	db, slogger, err := testutil.CreateTemporaryDBAndLogger("test_associate_tickets", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}
	defer db.Close()

	// Test variables
	const targetEmail = "test@regncon.com"
	const fakePeopleAmount = 100
	const billettholderConversionRatio = 0.5

	// Happy user will never be included in conversion!
	var happyPerson = testutil.GenerateFakePerson()
	happyPerson.Email = targetEmail
	var happyPersonTicket = CheckInTicket{
		ID:        1,
		OrderID:   1,
		TypeId:    8999,
		Type:      "Manuell billett",
		FirstName: happyPerson.FirstName,
		LastName:  happyPerson.LastName,
		Email:     happyPerson.Email,
		IsOver18:  true,
	}

	// Generate test data
	var generatededTickets []CheckInTicket
	generatededPeople := testutil.GeneratePeople(fakePeopleAmount)
	for i, generatedPerson := range generatededPeople {
		// Tie 10% of tickets with our target email
		var emailValue = targetEmail
		if rand.Intn(10) > 1 {
			emailValue = generatedPerson.Email
		}

		// Start at ID 2+ to allow happy person
		generatededTickets = append(generatededTickets, CheckInTicket{
			ID:        i + 2,
			OrderID:   i + 2,
			TypeId:    9000,
			FirstName: generatedPerson.FirstName,
			LastName:  generatedPerson.LastName,
			Type:      "Test billet",
			Email:     emailValue,
			IsOver18:  rand.Intn(10) > 2,
		})
	}

	// Add happy ticket to the end of our tickets array
	generatededTickets = append(generatededTickets, happyPersonTicket)

	// How many tickets have the targetEmail as their email?
	var expectedTargetEmailCount int
	for _, targetEmailCount := range generatededTickets {
		if targetEmailCount.Email == targetEmail {
			expectedTargetEmailCount++
		}
	}

	// Slize generated tickets from begining according to conversion
	// ammount and write them to billettholders table
	var conversionAmmount = fakePeopleAmount * billettholderConversionRatio
	billettholderConversion := generatededTickets[:int(conversionAmmount)]
	for _, ticket := range billettholderConversion {
		// fmt.Printf("Preparing billettholder: %+v\n", ticket)
		err = converTicketIdToNewBillettholder(ticket.ID, billettholderConversion, db, slogger)
		if err != nil {
			fmt.Println(err)
		}
	}

	// How many tickets with targetedEmail was converted to existing billettholdere
	var expectedConvertedTargetedEmailCount int
	for _, billettholderConverted := range billettholderConversion {
		if billettholderConverted.Email == targetEmail {
			expectedConvertedTargetedEmailCount++
		}
	}

	// Remaining tickets after conversion
	// var remainingTickets = generatededTickets[int(conversionAmmount):]
	// remainingTickets = append(remainingTickets, happyPersonTicket)

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
	/* for _, expectedUser := range expectedUsers {
		err = AssociateUserWithBillettholder(expectedUser.UserID, db, slogger)
		if err != nil {
			fmt.Println(err)
		}
	} */

	// Act
	err = AssociateTicketsWithBillettholder(generatededTickets, targetEmail, db, slogger)
	if err != nil {
		t.Fatalf("failed to associate ticket with billettholder: %v", err)
	}

	// Assert

	// Check that billettholder contains expected ammount
	var resultBillettholderCount int
	resultBillettholderRow := db.QueryRow("SELECT COUNT(id) FROM billettholdere")
	if err := resultBillettholderRow.Scan(&resultBillettholderCount); err != nil {
		t.Fatal(resultBillettholderRow.Err())
	}
	if resultBillettholderCount != (expectedTargetEmailCount-expectedConvertedTargetedEmailCount)+int(conversionAmmount) {
		t.Fatalf("expected %d billettholders, got %d", (expectedTargetEmailCount-expectedConvertedTargetedEmailCount)+int(conversionAmmount), resultBillettholderCount)
	}

	// Check billettholder contains a certain amount of email: targetEmail
	var resultTargetEmailCount int
	resultTargetEmailRow := db.QueryRow("SELECT COUNT(email) FROM billettholder_emails WHERE email = ?", targetEmail)
	if err = resultTargetEmailRow.Scan(&resultTargetEmailCount); err != nil {
		t.Fatal(resultTargetEmailRow.Err())
	}
	if resultTargetEmailCount != expectedTargetEmailCount {
		t.Fatalf("expected %d billettholders_email with %s, got %d", expectedTargetEmailCount, targetEmail, resultTargetEmailCount)
	}

	// todo add more checks
}
