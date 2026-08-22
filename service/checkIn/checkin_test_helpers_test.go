package checkIn

import (
	"database/sql"
	"log/slog"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
)

func createCheckInTestDB(t testing.TB) (*sql.DB, *slog.Logger) {
	t.Helper()

	return testutil.CreateTestDBAndLogger(t, "checkin")
}

func insertBillettholder(t testing.TB, db *sql.DB, billettholderID int) {
	t.Helper()

	insertCheckInBillettholder(t, db, models.Billettholder{
		ID:           billettholderID,
		FirstName:    "Test",
		LastName:     "Participant",
		TicketTypeId: 199999,
		TicketType:   "Test",
		IsOver18:     true,
		OrderID:      19999999 + billettholderID,
		TicketID:     4999999 + billettholderID,
	})
}

func insertCheckInBillettholder(t testing.TB, db *sql.DB, billettholder models.Billettholder) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, billettholder.ID, billettholder.FirstName, billettholder.LastName, billettholder.TicketTypeId, billettholder.TicketType, billettholder.IsOver18, billettholder.OrderID, billettholder.TicketID)
}

func insertUser(t testing.TB, db *sql.DB, userID int, externalID string, email string) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO users (id, external_id, email, is_admin)
		VALUES (?, ?, ?, ?)
	`, userID, externalID, email, false)
}

func insertManualBillettholderEmail(t testing.TB, db *sql.DB, billettholderID int, email string) int {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO relation_billettholder_emails (billettholder_id, email, kind)
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

func deleteBillettholderEmailByID(t testing.TB, db *sql.DB, emailID int) {
	t.Helper()

	testutil.MustExec(t, db, `DELETE FROM relation_billettholder_emails WHERE id = ?`, emailID)
}

func insertBillettholderUserAssociation(t testing.TB, db *sql.DB, association models.BillettholderUsers) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO relation_billettholdere_users (billettholder_id, user_id)
		VALUES (?, ?)
	`, association.BillettholderID, association.UserID)
}

func queryBillettholderByTicketID(t testing.TB, db *sql.DB, ticketID int) models.Billettholder {
	t.Helper()

	var billettholder models.Billettholder
	err := db.QueryRow(`
		SELECT id, first_name, last_name, ticket_type_id, ticket_type,
			is_over_18, order_id, ticket_id, created_at, updated_at, created_by_id, updated_by_id
		FROM billettholdere
		WHERE ticket_id = ?
	`, ticketID).Scan(
		&billettholder.ID,
		&billettholder.FirstName,
		&billettholder.LastName,
		&billettholder.TicketTypeId,
		&billettholder.TicketType,
		&billettholder.IsOver18,
		&billettholder.OrderID,
		&billettholder.TicketID,
		&billettholder.CreatedAt,
		&billettholder.UpdatedAt,
		&billettholder.CreatedByID,
		&billettholder.UpdatedByID,
	)
	if err != nil {
		t.Fatalf("failed to query billettholder by ticket ID %d: %v", ticketID, err)
	}

	return billettholder
}

func queryBillettholderEmails(t testing.TB, db *sql.DB, billettholderID int) []models.BillettholderEmail {
	t.Helper()

	rows, err := db.Query(`
		SELECT id, billettholder_id, email, kind, created_at, updated_at, created_by_id, updated_by_id
		FROM relation_billettholder_emails
		WHERE billettholder_id = ?
		ORDER BY id
	`, billettholderID)
	if err != nil {
		t.Fatalf("failed to query billettholder emails: %v", err)
	}
	defer rows.Close()

	var emails []models.BillettholderEmail
	for rows.Next() {
		var email models.BillettholderEmail
		if err := rows.Scan(&email.ID, &email.BillettholderID, &email.Email, &email.Kind, &email.CreatedAt, &email.UpdatedAt, &email.CreatedByID, &email.UpdatedByID); err != nil {
			t.Fatalf("failed to scan billettholder email: %v", err)
		}
		emails = append(emails, email)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("failed to iterate billettholder emails: %v", err)
	}

	return emails
}

func assertBillettholderEmails(t testing.TB, db *sql.DB, billettholderID int, expected []models.BillettholderEmail) {
	t.Helper()

	actual := queryBillettholderEmails(t, db, billettholderID)
	if len(actual) != len(expected) {
		t.Fatalf("billettholder email count mismatch\nexpected: %d\nactual:   %d", len(expected), len(actual))
	}

	for i, expectedEmail := range expected {
		actualEmail := actual[i]
		if actualEmail.Email != expectedEmail.Email || actualEmail.Kind != expectedEmail.Kind {
			t.Fatalf("billettholder email mismatch at index %d\nexpected: %s/%s\nactual:   %s/%s", i, expectedEmail.Email, expectedEmail.Kind, actualEmail.Email, actualEmail.Kind)
		}
	}
}

func queryBillettholderUserAssociations(t testing.TB, db *sql.DB) []models.BillettholderUsers {
	t.Helper()

	rows, err := db.Query(`
		SELECT billettholder_id, user_id
		FROM relation_billettholdere_users
		ORDER BY billettholder_id, user_id
	`)
	if err != nil {
		t.Fatalf("failed to query billettholder user associations: %v", err)
	}
	defer rows.Close()

	var associations []models.BillettholderUsers
	for rows.Next() {
		var association models.BillettholderUsers
		if err := rows.Scan(&association.BillettholderID, &association.UserID); err != nil {
			t.Fatalf("failed to scan billettholder user association: %v", err)
		}
		associations = append(associations, association)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("failed to iterate billettholder user associations: %v", err)
	}

	return associations
}

func assertBillettholderUserAssociations(t testing.TB, db *sql.DB, expected []models.BillettholderUsers) {
	t.Helper()

	actual := queryBillettholderUserAssociations(t, db)
	if !slices.Equal(expected, actual) {
		t.Fatalf("billettholder user associations mismatch\nexpected: %+v\nactual:   %+v", expected, actual)
	}
}

func assertOnlyBillettholderUserAssociation(t testing.TB, db *sql.DB, expected models.BillettholderUsers) {
	t.Helper()

	assertBillettholderUserAssociationCount(t, db, expected, 1)

	totalAssociations := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_billettholdere_users`)
	if totalAssociations != 1 {
		t.Fatalf("expected exactly 1 total billettholder user association, got %d", totalAssociations)
	}
}

func assertBillettholderUserAssociationCount(t testing.TB, db *sql.DB, expected models.BillettholderUsers, expectedCount int) {
	t.Helper()

	associationCount := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_billettholdere_users
		WHERE billettholder_id = ? AND user_id = ?
	`, expected.BillettholderID, expected.UserID)
	if associationCount != expectedCount {
		t.Fatalf("expected %d billettholder user associations, got %d", expectedCount, associationCount)
	}
}
