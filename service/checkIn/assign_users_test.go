package checkIn

import (
	"database/sql"
	"fmt"
	"log/slog"
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

	db, err := service.InitTestDBFrom(testDBPath)
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// expected billettholder_user
	var expectedBillettholderUser = models.BillettholderUsers{
		BillettholderID: 9999,
		UserID:          1,
	}

	var expectedMissMatchBillettholderUser = models.BillettholderUsers{
		BillettholderID: 8888,
		UserID:          1,
	}

	// Happy path person
	const happyPathEmail = "test@regncon.no"
	var happyPathPerson = testutil.GenerateFakePerson()
	happyPathPerson.Email = happyPathEmail

	var uassociatedPerson = testutil.GenerateFakePerson()
	var missMatchPerson = testutil.GenerateFakePerson()

	// Happy path user
	var happyPathUser = models.User{
		ID:         expectedBillettholderUser.UserID,
		ExternalID: "testuser",
		Email:      happyPathEmail,
		IsAdmin:    false,
	}

	var unassociatedPathUser = models.User{
		ID:         2,
		ExternalID: "gudrunanita",
		Email:      happyPathEmail,
		IsAdmin:    false,
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

	var missMatchBillettholder = models.Billettholder{
		ID:           expectedMissMatchBillettholderUser.BillettholderID,
		FirstName:    missMatchPerson.FirstName,
		LastName:     missMatchPerson.LastName,
		TicketTypeId: 199999,
		TicketType:   "Test",
		IsOver18:     true,
		OrderID:      19999999,
		TicketID:     4999997,
	}

	var unassociatedBillettholder = models.Billettholder{
		ID:           unassociatedPathUser.ID,
		FirstName:    uassociatedPerson.FirstName,
		LastName:     uassociatedPerson.LastName,
		TicketTypeId: 199999,
		TicketType:   "Test",
		IsOver18:     true,
		OrderID:      19999999,
		TicketID:     4999998,
	}

	// Happy path billettholder_email
	var happyPathBillettholderEmail = models.BillettholderEmail{
		BillettholderID: happyPathBillettholder.ID,
		Email:           happyPathPerson.Email,
		Kind:            models.BillettholderEmailKindManual,
	}

	var missMatchedBillettholderEmail = models.BillettholderEmail{
		BillettholderID: missMatchBillettholder.ID,
		Email:           strings.ToUpper(happyPathPerson.Email),
		Kind:            models.BillettholderEmailKindManual,
	}

	var uassociatedBillettholderEmail = models.BillettholderEmail{
		BillettholderID: unassociatedBillettholder.ID,
		Email:           uassociatedPerson.Email,
		Kind:            models.BillettholderEmailKindManual,
	}

	var testBillettholders []models.Billettholder
	testBillettholders = append(testBillettholders, happyPathBillettholder)
	testBillettholders = append(testBillettholders, missMatchBillettholder)
	testBillettholders = append(testBillettholders, unassociatedBillettholder)

	var testBillettholderEmails []models.BillettholderEmail
	testBillettholderEmails = append(testBillettholderEmails, happyPathBillettholderEmail)
	testBillettholderEmails = append(testBillettholderEmails, missMatchedBillettholderEmail)
	testBillettholderEmails = append(testBillettholderEmails, uassociatedBillettholderEmail)

	var testUsers []models.User
	testUsers = append(testUsers, happyPathUser)
	testUsers = append(testUsers, unassociatedPathUser)

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

	// Attempt to insert into relation_billettholder_emails
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
		queryBillettholderEmail = append(queryBillettholderEmail, fmt.Sprintf(`(%d, "%s", "%s")`, billettholderEmail.BillettholderID, billettholderEmail.Email, models.BillettholderEmailKindManual))
	}
	queryBase = fmt.Sprintf(`
            INSERT INTO relation_billettholder_emails (
                billettholder_id, email, kind
                ) VALUES %s`, strings.Join(queryBillettholderEmail, ", "))

	_, err = db.Exec(queryBase)
	if err != nil {
		fmt.Println("failed to insert relation_billettholder_emails", "error", err)
		return
	}

	var queryUsers []string
	for _, user := range testUsers {
		queryUsers = append(queryUsers, fmt.Sprintf(`(%d, "%s", "%s", %v)`, user.ID, user.ExternalID, user.Email, user.IsAdmin))
	}

	queryBase = fmt.Sprintf(`
                INSERT INTO users (
                    id, external_id, email, is_admin
                    ) VALUES %s`, strings.Join(queryUsers, ", "))

	_, err = db.Exec(queryBase)
	if err != nil {
		fmt.Println("failed to insert users", "error", err)
		return
	}

	// Act
	sl := &testutil.StubLogger{}
	slogger := testutil.NewSlogAdapter(sl)

	/* for _, user := range testUsers {
		// fmt.Printf("Calling AssociateUserWithBillettholder() on: %s (%s)\n", user.ExternalID, user.Email)
		err = AssociateUserWithBillettholder(user.ExternalID, db, slogger)
		if err != nil {
			t.Fatalf("failed to convert ticketId to billettholder: %v", err)
		}
	} */

	err = AssociateUserWithBillettholder(happyPathUser.ExternalID, db, slogger)
	if err != nil {
		t.Fatalf("failed to convert ticketId to billettholder: %v", err)
	}

	// Assert
	if happyPathBillettholderEmail.Email == missMatchedBillettholderEmail.Email {
		t.Fatalf("You did something wrong in var missMatchedBillettholderEmail when copying values")
	}

	var resultBillettholderUsers []models.BillettholderUsers
	rows, err := db.Query("SELECT billettholder_id, user_id FROM relation_billettholdere_users")
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

	// Both happyPathBillettholder and missMatchBillettholder should be found
	if len(resultBillettholderUsers) != 2 {
		t.Fatalf("expected 2 billettholder_user, got %d", len(resultBillettholderUsers))
	}

	// Both tickets should map to expected billettholder user
	for _, billettholderUser := range resultBillettholderUsers {
		if billettholderUser != expectedBillettholderUser && billettholderUser != expectedMissMatchBillettholderUser {
			t.Fatalf("expected billettholder_user %+v, got %+v", expectedBillettholderUser, billettholderUser)
		}
	}
}

func TestAssociateUsersWithBillettholderEmail_CreatesAssociationForMatchingUserEmail(t *testing.T) {
	// Gitt at ein billettholder har fått lagt til ei manuell e-postadresse,
	// og ein eksisterande brukar har same e-postadresse med annan casing,
	// når e-postadressa blir forsona mot brukarar,
	// så skal billettholderen få ei varig brukar-tilknyting.

	// Given
	expectedAssociation := models.BillettholderUsers{
		BillettholderID: 12345,
		UserID:          67890,
	}
	manualEmail := "participant@example.com"
	userEmail := "Participant@Example.com"

	db, slogger := createAssociationTestDB(t)
	defer db.Close()

	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertUser(t, db, expectedAssociation.UserID, "test-user", userEmail)
	insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, manualEmail)

	// When
	err := AssociateUsersWithBillettholderEmail(expectedAssociation.BillettholderID, manualEmail, db, slogger)

	// Then
	if err != nil {
		t.Fatalf("expected association to succeed: %v", err)
	}
	assertOnlyBillettholderUserAssociation(t, db, expectedAssociation)
}

func TestAssociateUsersWithBillettholderEmail_DoesNotDuplicateExistingAssociation(t *testing.T) {
	// Gitt at ein billettholder allereie er knytt til ein brukar via ei manuell e-postadresse,
	// når same e-postforsoning køyrer på nytt,
	// så skal det framleis berre finnast éi brukar-tilknyting.

	// Given
	expectedAssociation := models.BillettholderUsers{
		BillettholderID: 12345,
		UserID:          67890,
	}
	expectedAssociationCount := 1
	manualEmail := "participant@example.com"

	db, slogger := createAssociationTestDB(t)
	defer db.Close()

	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertUser(t, db, expectedAssociation.UserID, "test-user", manualEmail)
	insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, manualEmail)
	insertBillettholderUserAssociation(t, db, expectedAssociation)

	// When
	err := AssociateUsersWithBillettholderEmail(expectedAssociation.BillettholderID, manualEmail, db, slogger)

	// Then
	if err != nil {
		t.Fatalf("expected repeated association to succeed: %v", err)
	}
	assertBillettholderUserAssociationCount(t, db, expectedAssociation, expectedAssociationCount)
}

func TestDisassociateUsersFromBillettholderEmail_RemovesAssociationWhenNoRemainingEmailMatchesUser(t *testing.T) {
	// Gitt at ei manuell e-postadresse er fjerna frå ein billettholder,
	// og ingen attverande e-postadresser på billettholderen samsvarer med brukaren,
	// når e-postadressa blir forsona mot brukar-tilknytingar,
	// så skal den varige brukar-tilknytinga fjernast.

	// Given
	expectedAssociation := models.BillettholderUsers{
		BillettholderID: 12345,
		UserID:          67890,
	}
	expectedAssociationCount := 0
	manualEmail := "participant@example.com"
	userEmail := "Participant@Example.com"

	db, slogger := createAssociationTestDB(t)
	defer db.Close()

	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertUser(t, db, expectedAssociation.UserID, "test-user", userEmail)
	removedEmailID := insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, manualEmail)
	insertBillettholderUserAssociation(t, db, expectedAssociation)
	deleteBillettholderEmailByID(t, db, removedEmailID)

	// When
	err := DisassociateUsersFromBillettholderEmail(expectedAssociation.BillettholderID, manualEmail, db, slogger)

	// Then
	if err != nil {
		t.Fatalf("expected disassociation to succeed: %v", err)
	}
	assertBillettholderUserAssociationCount(t, db, expectedAssociation, expectedAssociationCount)
}

func TestDisassociateUsersFromBillettholderEmail_KeepsAssociationWhenRemainingEmailStillMatchesUser(t *testing.T) {
	// Gitt at ei manuell e-postadresse er fjerna frå ein billettholder,
	// men ei anna attverande e-postadresse på same billettholder framleis samsvarer med brukaren,
	// når e-postadressa blir forsona mot brukar-tilknytingar,
	// så skal den varige brukar-tilknytinga behaldast.

	// Given
	expectedAssociation := models.BillettholderUsers{
		BillettholderID: 12345,
		UserID:          67890,
	}
	expectedAssociationCount := 1
	removedEmail := "participant@example.com"
	remainingEmail := "PARTICIPANT@example.com"
	userEmail := "Participant@Example.com"

	db, slogger := createAssociationTestDB(t)
	defer db.Close()

	insertBillettholder(t, db, expectedAssociation.BillettholderID)
	insertUser(t, db, expectedAssociation.UserID, "test-user", userEmail)
	removedEmailID := insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, removedEmail)
	insertManualBillettholderEmail(t, db, expectedAssociation.BillettholderID, remainingEmail)
	insertBillettholderUserAssociation(t, db, expectedAssociation)
	deleteBillettholderEmailByID(t, db, removedEmailID)

	// When
	err := DisassociateUsersFromBillettholderEmail(expectedAssociation.BillettholderID, removedEmail, db, slogger)

	// Then
	if err != nil {
		t.Fatalf("expected disassociation cleanup to succeed: %v", err)
	}
	assertBillettholderUserAssociationCount(t, db, expectedAssociation, expectedAssociationCount)
}

func createAssociationTestDB(t *testing.T) (*sql.DB, *slog.Logger) {
	t.Helper()

	db, slogger, err := testutil.CreateTemporaryDBAndLogger("test_associate_users_with_billettholder_email", t)
	if err != nil {
		t.Fatalf("failed to create test database and logger: %v", err)
	}

	return db, slogger
}

func insertBillettholder(t *testing.T, db *sql.DB, billettholderID int) {
	t.Helper()

	_, err := db.Exec(`
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, billettholderID, "Test", "Participant", 199999, "Test", true, 19999999, 4999999)
	if err != nil {
		t.Fatalf("failed to insert billettholder: %v", err)
	}
}

func insertUser(t *testing.T, db *sql.DB, userID int, descopeUserID string, email string) {
	t.Helper()

	_, err := db.Exec(`
		INSERT INTO users (id, user_id, email, is_admin)
		VALUES (?, ?, ?, ?)
	`, userID, descopeUserID, email, false)
	if err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}
}

func insertManualBillettholderEmail(t *testing.T, db *sql.DB, billettholderID int, email string) int {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO billettholder_emails (billettholder_id, email, kind)
		VALUES (?, ?, 'Manual')
	`, billettholderID, email)
	if err != nil {
		t.Fatalf("failed to insert billettholder email: %v", err)
	}

	emailID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get inserted billettholder email ID: %v", err)
	}
	return int(emailID)
}

func deleteBillettholderEmailByID(t *testing.T, db *sql.DB, emailID int) {
	t.Helper()

	_, err := db.Exec(`DELETE FROM billettholder_emails WHERE id = ?`, emailID)
	if err != nil {
		t.Fatalf("failed to delete billettholder email: %v", err)
	}
}

func insertBillettholderUserAssociation(t *testing.T, db *sql.DB, association models.BillettholderUsers) {
	t.Helper()

	_, err := db.Exec(`
		INSERT INTO billettholdere_users (billettholder_id, user_id)
		VALUES (?, ?)
	`, association.BillettholderID, association.UserID)
	if err != nil {
		t.Fatalf("failed to insert billettholder user association: %v", err)
	}
}

func assertOnlyBillettholderUserAssociation(t *testing.T, db *sql.DB, expected models.BillettholderUsers) {
	t.Helper()

	assertBillettholderUserAssociationCount(t, db, expected, 1)

	var totalAssociations int
	err := db.QueryRow(`SELECT COUNT(*) FROM billettholdere_users`).Scan(&totalAssociations)
	if err != nil {
		t.Fatalf("failed to count all billettholder user associations: %v", err)
	}
	if totalAssociations != 1 {
		t.Fatalf("expected exactly 1 total billettholder user association, got %d", totalAssociations)
	}
}

func assertBillettholderUserAssociationCount(t *testing.T, db *sql.DB, expected models.BillettholderUsers, expectedCount int) {
	t.Helper()

	var associationCount int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM billettholdere_users
		WHERE billettholder_id = ? AND user_id = ?
	`, expected.BillettholderID, expected.UserID).Scan(&associationCount)
	if err != nil {
		t.Fatalf("failed to count billettholder user associations: %v", err)
	}
	if associationCount != expectedCount {
		t.Fatalf("expected %d billettholder user associations, got %d", expectedCount, associationCount)
	}
}
