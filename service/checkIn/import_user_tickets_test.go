package checkIn

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

const (
	parentEmail = "forelder@example.com"
	teenEmail   = "ungdom@example.com"
)

// familyOrderTickets is one CheckIn order where a parent bought tickets for
// themselves and a teen, with the teen's own e-post on the teen's ticket.
var familyOrderTickets = []CheckInTicket{
	{ID: 601, OrderID: 60, TypeId: 9000, Type: "Festivalpass", FirstName: "Forelder", LastName: "Familie", Email: parentEmail, IsOver18: true},
	{ID: 602, OrderID: 60, TypeId: 9001, Type: "Ungdom", FirstName: "Ungdom", LastName: "Familie", Email: teenEmail, IsOver18: false},
}

const familyOrderCheckinResponse = `{"data":{"eventTickets":[` +
	`{"id":601,"order_id":60,"category":"Festivalpass","category_id":9000,"crm":{"first_name":"Forelder","last_name":"Familie","email":"forelder@example.com","born":"1980-01-01"}},` +
	`{"id":602,"order_id":60,"category":"Ungdom","category_id":9001,"crm":{"first_name":"Ungdom","last_name":"Familie","email":"ungdom@example.com","born":"2010-01-01"}}` +
	`]}}`

func TestImportUserTickets_WhenAnotherUserOnTheOrderAlreadyExists_LinksThemToo(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en ungdom allerede har logget inn, og forelderen har kjøpt billetter til begge på samme bestilling.",
		When:  "Når forelderen trykker «Hent billetter».",
		Then:  "Så skal både forelderen og ungdommen knyttes til billettholderne på bestillingen, og bare forelderens nye koblinger telles.",
	})

	// Given
	parentUserID, teenUserID := 1, 2
	expectedCreatedUserAssociations := 2
	db, logger := createCheckInTestDB(t)
	insertUser(t, db, parentUserID, "parent-user", parentEmail)
	insertUser(t, db, teenUserID, "teen-user", teenEmail)

	// When
	result, err := ImportUserTickets(familyOrderTickets, "parent-user", parentEmail, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected ticket import to succeed: %v", err)
	}
	if result.CreatedUserAssociations != expectedCreatedUserAssociations {
		t.Fatalf("created user association count mismatch\nexpected: %d\nactual:   %d", expectedCreatedUserAssociations, result.CreatedUserAssociations)
	}
	parentTicket := queryBillettholderByTicketID(t, db, 601)
	teenTicket := queryBillettholderByTicketID(t, db, 602)
	assertBillettholderUserAssociations(t, db, []models.BillettholderUsers{
		{BillettholderID: parentTicket.ID, UserID: parentUserID},
		{BillettholderID: parentTicket.ID, UserID: teenUserID},
		{BillettholderID: teenTicket.ID, UserID: parentUserID},
		{BillettholderID: teenTicket.ID, UserID: teenUserID},
	})
}

func TestImportUserTickets_WhenManualEmailWasAddedBeforeFirstLogin_LinksTheUser(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at forelderen la til ungdommens e-post på billetten før ungdommen hadde logget inn.",
		When:  "Når ungdommens billetter hentes etter første innlogging, uten egne billetter i CheckIn.",
		Then:  "Så skal ungdommen knyttes til billettholderen med e-posten sin.",
	})

	// Given
	expectedAssociation := models.BillettholderUsers{BillettholderID: 700, UserID: 2}
	db, logger := createCheckInTestDB(t)
	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, teenEmail)
	insertUser(t, db, expectedAssociation.UserID, "teen-user", teenEmail)

	// When
	_, err := ImportUserTickets(nil, "teen-user", teenEmail, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected ticket import to succeed: %v", err)
	}
	assertOnlyBillettholderUserAssociation(t, db, expectedAssociation)
}

func TestConvertOrderToBillettholdere_WhenAdminConvertsOneTicket_ConvertsTheWholeOrderAndLinksUsers(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at både forelderen og ungdommen har logget inn, og billettene deres står på samme bestilling.",
		When:  "Når admin konverterer bestillingen fra ungdommens billett.",
		Then:  "Så skal alle billettene på bestillingen bli billettholdere, og alle brukere med en e-post på dem skal knyttes til dem.",
	})

	// Given
	parentUserID, teenUserID := 1, 2
	expectedCreatedBillettholders := 2
	db, logger := createCheckInTestDB(t)
	insertUser(t, db, parentUserID, "parent-user", parentEmail)
	insertUser(t, db, teenUserID, "teen-user", teenEmail)
	useCheckInResponse(t, familyOrderCheckinResponse)

	// When
	result, err := ConvertOrderToBillettholdere(context.Background(), 602, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected order conversion to succeed: %v", err)
	}
	if result.CreatedBillettholders != expectedCreatedBillettholders {
		t.Fatalf("created billettholder count mismatch\nexpected: %d\nactual:   %d", expectedCreatedBillettholders, result.CreatedBillettholders)
	}
	parentTicket := queryBillettholderByTicketID(t, db, 601)
	teenTicket := queryBillettholderByTicketID(t, db, 602)
	assertBillettholderUserAssociations(t, db, []models.BillettholderUsers{
		{BillettholderID: parentTicket.ID, UserID: parentUserID},
		{BillettholderID: parentTicket.ID, UserID: teenUserID},
		{BillettholderID: teenTicket.ID, UserID: parentUserID},
		{BillettholderID: teenTicket.ID, UserID: teenUserID},
	})
}

func TestConvertOrderToBillettholdere_WhenTicketIsUnknown_ReturnsError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at billetten admin konverterer ikke finnes i CheckIn.",
		When:  "Når admin konverterer bestillingen.",
		Then:  "Så skal konverteringen feile uten å opprette billettholdere.",
	})

	// Given
	expectedBillettholderCount := 0
	db, logger := createCheckInTestDB(t)
	useCheckInResponse(t, familyOrderCheckinResponse)

	// When
	_, err := ConvertOrderToBillettholdere(context.Background(), 999, db, logger)

	// Then
	if err == nil {
		t.Fatal("expected converting an unknown ticket to fail")
	}
	actualBillettholderCount := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM billettholdere`)
	if actualBillettholderCount != expectedBillettholderCount {
		t.Fatalf("billettholder count mismatch\nexpected: %d\nactual:   %d", expectedBillettholderCount, actualBillettholderCount)
	}
}

func TestLinkUsersToBillettholdere_WhenRunTwice_CreatesNoDuplicates(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en billettholder med e-posten til en bruker, skrevet med andre store og små bokstaver.",
		When:  "Når brukerne knyttes til billettholderen to ganger.",
		Then:  "Så skal brukeren være knyttet én gang, og andre kjøring skal ikke lage nye koblinger.",
	})

	// Given
	expectedAssociation := models.BillettholderUsers{BillettholderID: 700, UserID: 2}
	expectedSecondRunLinks := 0
	db, logger := createCheckInTestDB(t)
	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, "UNGDOM@example.com")
	insertUser(t, db, expectedAssociation.UserID, "teen-user", teenEmail)
	if _, err := LinkUsersToBillettholdere([]int{expectedAssociation.BillettholderID}, db, logger); err != nil {
		t.Fatalf("expected first link run to succeed: %v", err)
	}

	// When
	actualSecondRunLinks, err := LinkUsersToBillettholdere([]int{expectedAssociation.BillettholderID}, db, logger)

	// Then
	if err != nil {
		t.Fatalf("expected second link run to succeed: %v", err)
	}
	if actualSecondRunLinks != expectedSecondRunLinks {
		t.Fatalf("second run link count mismatch\nexpected: %d\nactual:   %d", expectedSecondRunLinks, actualSecondRunLinks)
	}
	assertOnlyBillettholderUserAssociation(t, db, expectedAssociation)
}

// useCheckInResponse points the shared CheckIn ticket cache at a fake CheckIn
// that always answers with response, and restores the real cache afterwards.
func useCheckInResponse(t *testing.T, response string) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)

	originalCache := ticketCache
	ticketCache = newTestCache(server.URL, server.Client())
	t.Cleanup(func() { ticketCache = originalCache })
}
