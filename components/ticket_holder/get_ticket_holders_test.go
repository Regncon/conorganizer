package ticketholder

import (
	"database/sql"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

const getTicketHoldersTestUserID = "teen-user"
const getTicketHoldersTestUserEmail = "ungdom@example.com"

func TestGetTicketHolders_WhenBillettholderIsLinkedToTheUser_ListsItWithTheTicketEmail(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en ungdom er knyttet til en billettholder som står på forelderens e-post.",
		When:  "Når billettholderne til ungdommen hentes.",
		Then:  "Så skal billettholderen vises med billettens e-post.",
	})

	// Given
	expected := []BillettHolder{{Id: 801, Email: "forelder@example.com", Name: "Ungdom Familie", TicketType: "Ungdom"}}
	db := createGetTicketHoldersTestDB(t)
	insertGetTicketHoldersBillettholder(t, db, expected[0].Id, "forelder@example.com")
	linkGetTicketHoldersBillettholder(t, db, expected[0].Id)

	// When
	actual, err := GetTicketHolders(getTicketHoldersTestUserInfo(), db)

	// Then
	if err != nil {
		t.Fatalf("expected ticket holders to load: %v", err)
	}
	assertTicketHolders(t, expected, actual)
}

func TestGetTicketHolders_WhenOnlyTheEmailMatches_DoesNotListTheBillettholder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder har brukerens e-post, men brukeren ikke er knyttet til den.",
		When:  "Når billettholderne til brukeren hentes.",
		Then:  "Så skal billettholderen ikke vises, siden brukeren ikke kan melde interesse for den.",
	})

	// Given
	var expected []BillettHolder
	db := createGetTicketHoldersTestDB(t)
	insertGetTicketHoldersBillettholder(t, db, 801, getTicketHoldersTestUserEmail)

	// When
	actual, err := GetTicketHolders(getTicketHoldersTestUserInfo(), db)

	// Then
	if err != nil {
		t.Fatalf("expected ticket holders to load: %v", err)
	}
	assertTicketHolders(t, expected, actual)
}

func createGetTicketHoldersTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, _ := testutil.CreateTestDBAndLogger(t, "get_ticket_holders")
	testutil.MustExec(t, db, `
		INSERT INTO users (external_id, email) VALUES (?, ?)
	`, getTicketHoldersTestUserID, getTicketHoldersTestUserEmail)
	return db
}

func getTicketHoldersTestUserInfo() requestctx.UserRequestInfo {
	return requestctx.UserRequestInfo{IsLoggedIn: true, Id: getTicketHoldersTestUserID, Email: getTicketHoldersTestUserEmail}
}

func insertGetTicketHoldersBillettholder(t *testing.T, db *sql.DB, billettholderID int, ticketEmail string) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO billettholdere (id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id)
		VALUES (?, 'Ungdom', 'Familie', 9001, 'Ungdom', 0, 60, ?)
	`, billettholderID, billettholderID+1000)
	testutil.MustExec(t, db, `
		INSERT INTO relation_billettholder_emails (billettholder_id, email, kind) VALUES (?, ?, ?)
	`, billettholderID, ticketEmail, models.BillettholderEmailKindTicket)
}

func linkGetTicketHoldersBillettholder(t *testing.T, db *sql.DB, billettholderID int) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO relation_billettholdere_users (billettholder_id, user_id)
		SELECT ?, id FROM users WHERE external_id = ?
	`, billettholderID, getTicketHoldersTestUserID)
}

func assertTicketHolders(t *testing.T, expected []BillettHolder, actual []BillettHolder) {
	t.Helper()

	type ticketHolderSummary struct {
		Id         int
		Email      string
		Name       string
		TicketType string
	}
	summarize := func(holders []BillettHolder) []ticketHolderSummary {
		var summaries []ticketHolderSummary
		for _, holder := range holders {
			summaries = append(summaries, ticketHolderSummary{holder.Id, holder.Email, holder.Name, holder.TicketType})
		}
		return summaries
	}
	if !slices.Equal(summarize(expected), summarize(actual)) {
		t.Fatalf("ticket holders mismatch\nexpected: %+v\nactual:   %+v", summarize(expected), summarize(actual))
	}
}
