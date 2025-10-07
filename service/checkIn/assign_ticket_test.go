package checkIn

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
)

func TestAssociateTicketsWithEmail(t *testing.T) {
	// Arrange
	var ammountOfFakeTickets = 99
	const targetEmail = "test@regncon.no"
	const targetEmailUppercase = "Test@regncon.no"

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

	var expectedMatches []CheckInTicket
	for _, generatedTicket := range generatedTickets {
		if generatedTicket.TypeId == TicketTypeMiddag {
			continue
		}

		if strings.EqualFold(generatedTicket.Email, targetEmail) {
			expectedMatches = append(expectedMatches, generatedTicket)
		}
	}

	// Act
	/* sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl) */

	result, err := AssociateTicketsWithEmail(generatedTickets, targetEmail)
	if err != nil {
		t.Fatalf("failed to associate email with billettholder: %v", err)
	}

	// Case sensitivity
	resultUppercase, err := AssociateTicketsWithEmail(generatedTickets, targetEmailUppercase)
	if err != nil {
		t.Fatalf("failed to associate email with billettholder: %v", err)
	}

	// Assert
	if len(result) != len(expectedMatches) {
		t.Fatalf("expected %d tickets, got %d", len(expectedMatches), len(result))
	}

	// Case sensitivity
	if len(resultUppercase) != len(expectedMatches) {
		t.Fatalf("expected %d uppercase tickets, got %d", len(expectedMatches), len(resultUppercase))
	}
}
