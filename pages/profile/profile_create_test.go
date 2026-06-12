package profilepage

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
)

func TestCreateNewEventFormSubmission_WhenInsertFails_ReturnsFriendlyError(t *testing.T) {
	// Gitt at brukeren sender inn nytt arrangement og databasen ikke kan opprette det,
	// når opprettelsen feiler,
	// så skal svaret være en tydelig feil og ikke en tom vellykket respons.

	// Given
	expectedStatusCode := http.StatusInternalServerError
	expectedBodyPart := createEventFailureMessage

	db, logger := testutil.CreateTestDBAndLogger(t, "profile_create")
	testutil.MustExec(t, db, `
		INSERT INTO users(id, external_id, email)
		VALUES(?, ?, ?)
	`, 501, "profile-create-user", "profile-create-user@example.com")
	testutil.MustExec(t, db, `DROP TABLE events`)
	request := httptest.NewRequest(http.MethodPost, "/profile/api/create", nil)
	recorder := httptest.NewRecorder()
	user := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		Id:         "profile-create-user",
		Email:      "profile-create-user@example.com",
	}

	// When
	createNewEventFormSubmissionForUser(db, nil, logger, recorder, request, user)

	// Then
	if recorder.Code != expectedStatusCode {
		t.Fatalf("HTTP status mismatch\nexpected: %d\nactual:   %d", expectedStatusCode, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), expectedBodyPart) {
		t.Fatalf("HTTP body mismatch\nexpected to contain: %q\nactual:              %q", expectedBodyPart, recorder.Body.String())
	}
}
