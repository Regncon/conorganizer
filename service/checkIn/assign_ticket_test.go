package checkIn

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
)

func TestAssociateBillettholderWithEmail(t *testing.T) {
	// Arrange
	var ammountOfFakeTickets = 99
	const targetEmail = "test@regncon.no"
	const targetEmail2 = "Test@regncon.no"

	var generatedTickets []CheckInTicket
	for i := range ammountOfFakeTickets {
		var generatedPerson = testutil.GenerateFakePerson()

		// Tie 10% of tickets with our target email
		var emailValue = targetEmail
		if rand.Intn(10) > 1 {
			emailValue = generatedPerson.Email
		}

		generatedTickets = append(generatedTickets, CheckInTicket{
			OrderID:   i + 1,
			TypeId:    i + 1,
			FirstName: generatedPerson.FirstName,
			LastName:  generatedPerson.LastName,
			Type:      "Test billet",
			Email:     emailValue,
			IsOver18:  rand.Intn(10) > 2,
		})
	}

	var matches []CheckInTicket
	for _, generatedTicket := range generatedTickets {
		if generatedTicket.Email == targetEmail {
			matches = append(matches, generatedTicket)
		}
	}

	// Act
	/* sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl) */

	result, err := AssociateBillettholderWithEmail(generatedTickets, targetEmail)
	if err != nil {
		t.Fatalf("failed to associate email with billettholder: %v", err)
	}

	result2, err := AssociateBillettholderWithEmail(generatedTickets, targetEmail2)
	if err != nil {
		t.Fatalf("failed to associate email with billettholder: %v", err)
	}

	// Assert
	if len(result) != len(matches) {
		t.Fatalf("expected %d tickets, got %d", len(matches), len(result))
	} else {
		fmt.Printf("AssociateBillettholderWithEmail returned %d/%d matches, total tickets: %d\n", len(result), len(matches), len(generatedTickets))
	}

	if len(result2) != len(matches) {
		t.Fatalf("expected %d tickets, got %d", len(matches), len(result2))
	} else {
		fmt.Printf("AssociateBillettholderWithEmail returned %d/%d matches, total tickets: %d\n", len(result2), len(matches), len(generatedTickets))
	}
}
