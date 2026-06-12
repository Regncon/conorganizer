package login

import (
	"database/sql"
	"testing"

	"github.com/Regncon/conorganizer/testutil"
)

type postLoginUser struct {
	externalID string
	email      string
	isAdmin    int
}

func TestSyncPostLoginUser_WhenUserDoesNotExist_CreatesLocalUser(t *testing.T) {
	// Gitt at en Descope-bruker logger inn for første gang,
	// når post-login synkroniserer brukeren,
	// så skal brukeren lagres lokalt med riktig e-post og adminstatus.

	// Given
	expectedUser := postLoginUser{
		externalID: "descope-user-1",
		email:      "new-user@example.com",
		isAdmin:    1,
	}

	db, logger := testutil.CreateTestDBAndLogger(t, "post_login_create_user")

	// When
	err := syncPostLoginUser(db, expectedUser.externalID, expectedUser.email, true, logger)

	// Then
	if err != nil {
		t.Fatalf("expected post-login user sync to succeed: %v", err)
	}
	assertPostLoginUser(t, db, expectedUser)
}

func TestSyncPostLoginUser_WhenUserExists_UpdatesAdminStatusWithoutDuplicatingUser(t *testing.T) {
	// Gitt at en lokal bruker finnes fra før,
	// når post-login synkroniserer en endret adminstatus,
	// så skal eksisterende bruker oppdateres uten duplikat.

	// Given
	expectedUser := postLoginUser{
		externalID: "descope-user-2",
		email:      "existing-user@example.com",
		isAdmin:    1,
	}
	expectedUserCount := 1

	db, logger := testutil.CreateTestDBAndLogger(t, "post_login_update_user")
	testutil.MustExec(t, db, `
		INSERT INTO users(external_id, email, is_admin)
		VALUES(?, ?, 0)
	`, expectedUser.externalID, expectedUser.email)

	// When
	err := syncPostLoginUser(db, expectedUser.externalID, expectedUser.email, true, logger)
	actualUserCount := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM users WHERE email = ?`, expectedUser.email)

	// Then
	if err != nil {
		t.Fatalf("expected post-login user sync to succeed: %v", err)
	}
	assertPostLoginUser(t, db, expectedUser)
	if actualUserCount != expectedUserCount {
		t.Fatalf("user count mismatch\nexpected: %d\nactual:   %d", expectedUserCount, actualUserCount)
	}
}

func assertPostLoginUser(t *testing.T, db *sql.DB, expectedUser postLoginUser) {
	t.Helper()

	var actualUser postLoginUser
	if err := db.QueryRow(`
		SELECT external_id, email, is_admin
		FROM users
		WHERE external_id = ?
	`, expectedUser.externalID).Scan(
		&actualUser.externalID,
		&actualUser.email,
		&actualUser.isAdmin,
	); err != nil {
		t.Fatalf("failed to query post-login user: %v", err)
	}

	if actualUser != expectedUser {
		t.Fatalf("post-login user mismatch\nexpected: %+v\nactual:   %+v", expectedUser, actualUser)
	}
}
