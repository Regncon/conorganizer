package billettholderadmin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/live"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/go-chi/chi/v5"
)

func TestAddEmailToBilettholderRoute_AddsEmailToRequestedAdminCard(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at admin legger til en manuell e-postadresse på én av flere billettholdere.",
		When:  "Når handlingen lykkes.",
		Then:  "Så skal bekreftelsen gjelde riktig kort og adressen lagres på riktig billettholder.",
	})

	// Given
	expectedEmail := "new@example.com"
	expectedSuccess := "Epostadressen new@example.com er lagt til"
	db, router := setupAdminBillettholderEmailRouteTest(t)
	insertAdminRouteTestBillettholder(t, db, 42)
	insertAdminRouteTestBillettholder(t, db, 99)
	insertAdminRouteTestUser(t, db, 7, "New@Example.com")

	// When
	recorder := postAdminDatastarSignals(t, router, "/new-email/42/", map[string]string{
		"newEmail-42": expectedEmail,
		"newEmail-99": "wrong-card@example.com",
	})

	// Then
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected successful add status %d, got %d\nbody: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	assertAdminPatchedSignal(t, recorder.Body.String(), "successMessage-42", expectedSuccess)
	assertAdminPatchedSignal(t, recorder.Body.String(), "errorMessage-42", "")
	assertAdminPatchedSignal(t, recorder.Body.String(), "newEmail-42", "")
	assertAdminBillettholderEmailCount(t, db, 42, expectedEmail, 1)
	assertAdminBillettholderEmailCount(t, db, 99, "wrong-card@example.com", 0)
	assertAdminBillettholderUserAssociationCount(t, db, 42, 7, 1)
}

func TestAddEmailToBilettholderRoute_RejectsEmptyAdminEmail(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at admin forsøker å legge til en tom e-postadresse.",
		When:  "Når handlingen utføres.",
		Then:  "Så skal riktig kort få feilmelding og ingen adresse lagres.",
	})

	// Given
	expectedError := "Tomt felt for epostadresse"
	db, router := setupAdminBillettholderEmailRouteTest(t)
	insertAdminRouteTestBillettholder(t, db, 42)

	// When
	recorder := postAdminDatastarSignals(t, router, "/new-email/42/", map[string]string{
		"newEmail-42": "",
	})

	// Then
	assertAdminPatchedSignal(t, recorder.Body.String(), "errorMessage-42", expectedError)
	assertAdminPatchedSignal(t, recorder.Body.String(), "successMessage-42", "")
	assertAdminManualEmailCount(t, db, 42, 0)
}

func TestAddEmailToBilettholderRoute_RejectsDuplicateAdminEmail(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at admin forsøker å legge til en e-postadresse som allerede finnes på samme billettholder.",
		When:  "Når handlingen utføres.",
		Then:  "Så skal duplikatet avvises tydelig uten å lagre en ekstra adresse.",
	})

	// Given
	duplicateEmail := "dupe@example.com"
	expectedError := "Epostadressen dupe@example.com finnes allerede for denne bilettholderen"
	db, router := setupAdminBillettholderEmailRouteTest(t)
	insertAdminRouteTestBillettholder(t, db, 42)
	insertAdminRouteTestEmail(t, db, 42, duplicateEmail, models.BillettholderEmailKindManual)

	// When
	recorder := postAdminDatastarSignals(t, router, "/new-email/42/", map[string]string{
		"newEmail-42": duplicateEmail,
	})

	// Then
	assertAdminPatchedSignal(t, recorder.Body.String(), "errorMessage-42", expectedError)
	assertAdminPatchedSignal(t, recorder.Body.String(), "successMessage-42", "")
	assertAdminBillettholderEmailCount(t, db, 42, duplicateEmail, 1)
}

func TestDeleteEmailFromBillettholderRoute_DeletesRequestedAdminManualEmail(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at admin sletter en manuell e-postadresse på én av flere billettholdere.",
		When:  "Når handlingen lykkes.",
		Then:  "Så skal riktig adresse fjernes fra riktig kort.",
	})

	// Given
	deletedEmail := "delete-me@example.com"
	otherEmail := "keep-me@example.com"
	expectedSuccess := "Epostadressen delete-me@example.com er slettet"
	db, router := setupAdminBillettholderEmailRouteTest(t)
	insertAdminRouteTestBillettholder(t, db, 42)
	insertAdminRouteTestBillettholder(t, db, 99)
	deletedEmailID := insertAdminRouteTestEmail(t, db, 42, deletedEmail, models.BillettholderEmailKindManual)
	insertAdminRouteTestEmail(t, db, 99, otherEmail, models.BillettholderEmailKindManual)

	// When
	recorder := postAdminDatastarSignals(t, router, fmt.Sprintf("/delete-email/42/%d/", deletedEmailID), map[string]string{})

	// Then
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected successful delete status %d, got %d\nbody: %s", http.StatusOK, recorder.Code, recorder.Body.String())
	}
	assertAdminPatchedSignal(t, recorder.Body.String(), "successMessage-42", expectedSuccess)
	assertAdminPatchedSignal(t, recorder.Body.String(), "errorMessage-42", "")
	assertAdminBillettholderEmailCount(t, db, 42, deletedEmail, 0)
	assertAdminBillettholderEmailCount(t, db, 99, otherEmail, 1)
}

func TestDeleteEmailFromBillettholderRoute_RemovesAdminUserAssociationWhenNoMatchingEmailRemains(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at sletting av manuell e-postadresse gjør at bruker-tilknytningen ikke lenger har en matchende adresse.",
		When:  "Når admin sletter adressen.",
		Then:  "Så skal bruker-tilknytningen ryddes opp.",
	})

	// Given
	deletedEmail := "participant@example.com"
	db, router := setupAdminBillettholderEmailRouteTest(t)
	insertAdminRouteTestBillettholder(t, db, 42)
	insertAdminRouteTestUser(t, db, 7, deletedEmail)
	deletedEmailID := insertAdminRouteTestEmail(t, db, 42, deletedEmail, models.BillettholderEmailKindManual)
	insertAdminRouteTestBillettholderUserAssociation(t, db, 42, 7)

	// When
	postAdminDatastarSignals(t, router, fmt.Sprintf("/delete-email/42/%d/", deletedEmailID), map[string]string{})

	// Then
	assertAdminBillettholderEmailCount(t, db, 42, deletedEmail, 0)
	assertAdminBillettholderUserAssociationCount(t, db, 42, 7, 0)
}

func TestDeleteEmailFromBillettholderRoute_KeepsAdminUserAssociationWhenMatchingEmailRemains(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at en billettholder har en bruker-tilknytning og en annen gjenværende e-postadresse som fortsatt matcher brukeren.",
		When:  "Når admin sletter én av adressene.",
		Then:  "Så skal bruker-tilknytningen beholdes.",
	})

	// Given
	deletedEmail := "participant@example.com"
	remainingEmail := "PARTICIPANT@example.com"
	db, router := setupAdminBillettholderEmailRouteTest(t)
	insertAdminRouteTestBillettholder(t, db, 42)
	insertAdminRouteTestUser(t, db, 7, deletedEmail)
	deletedEmailID := insertAdminRouteTestEmail(t, db, 42, deletedEmail, models.BillettholderEmailKindManual)
	insertAdminRouteTestEmail(t, db, 42, remainingEmail, models.BillettholderEmailKindManual)
	insertAdminRouteTestBillettholderUserAssociation(t, db, 42, 7)

	// When
	postAdminDatastarSignals(t, router, fmt.Sprintf("/delete-email/42/%d/", deletedEmailID), map[string]string{})

	// Then
	assertAdminBillettholderEmailIDCount(t, db, deletedEmailID, 0)
	assertAdminBillettholderEmailCount(t, db, 42, remainingEmail, 1)
	assertAdminBillettholderUserAssociationCount(t, db, 42, 7, 1)
}

func setupAdminBillettholderEmailRouteTest(t testing.TB) (*sql.DB, chi.Router) {
	t.Helper()

	db, logger := testutil.CreateTestDBAndLogger(t, "admin_billettholder_email_routes")
	router := chi.NewRouter()
	liveManager := &live.Manager{}
	addEmailToBilettholderRoute(router, db, logger, liveManager)
	deleteEmailFromBillettholderRoute(router, db, logger, liveManager)
	return db, router
}

func postAdminDatastarSignals(t testing.TB, router http.Handler, path string, signals map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	body, err := json.Marshal(signals)
	if err != nil {
		t.Fatalf("failed to marshal Datastar signals: %v", err)
	}
	request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(string(body)))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	return recorder
}

func insertAdminRouteTestBillettholder(t testing.TB, db *sql.DB, billettholderID int) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, 'Test', 'Participant', 100, 'Festivalpass', 1, ?, ?)
	`, billettholderID, 10000+billettholderID, 20000+billettholderID)
}

func insertAdminRouteTestUser(t testing.TB, db *sql.DB, userID int, email string) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO users (id, external_id, email, is_admin)
		VALUES (?, ?, ?, 0)
	`, userID, fmt.Sprintf("admin-route-user-%d", userID), email)
}

func insertAdminRouteTestEmail(t testing.TB, db *sql.DB, billettholderID int, email string, kind models.BillettholderEmailKind) int {
	t.Helper()

	result, err := db.Exec(`
		INSERT INTO relation_billettholder_emails (billettholder_id, email, kind)
		VALUES (?, ?, ?)
	`, billettholderID, email, kind)
	if err != nil {
		t.Fatalf("failed to insert billettholder email: %v", err)
	}
	emailID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("failed to get billettholder email id: %v", err)
	}
	return int(emailID)
}

func insertAdminRouteTestBillettholderUserAssociation(t testing.TB, db *sql.DB, billettholderID int, userID int) {
	t.Helper()

	testutil.MustExec(t, db, `
		INSERT INTO relation_billettholdere_users (billettholder_id, user_id)
		VALUES (?, ?)
	`, billettholderID, userID)
}

func assertAdminPatchedSignal(t testing.TB, body string, key string, expected string) {
	t.Helper()

	signals := adminPatchedSignals(t, body)
	actual, ok := signals[key]
	if !ok {
		t.Fatalf("expected patched signal %q\navailable signals: %#v\nbody: %s", key, signals, body)
	}
	if actual != expected {
		t.Fatalf("patched signal %q mismatch\nexpected: %q\nactual:   %q\nbody: %s", key, expected, actual, body)
	}
}

func adminPatchedSignals(t testing.TB, body string) map[string]string {
	t.Helper()

	signals := map[string]string{}
	for line := range strings.SplitSeq(body, "\n") {
		payload, ok := strings.CutPrefix(line, "data: signals ")
		if !ok {
			continue
		}
		patch := map[string]string{}
		if err := json.Unmarshal([]byte(payload), &patch); err != nil {
			t.Fatalf("failed to unmarshal Datastar signal patch %q: %v", payload, err)
		}
		maps.Copy(signals, patch)
	}
	return signals
}

func assertAdminBillettholderEmailCount(t testing.TB, db *sql.DB, billettholderID int, email string, expected int) {
	t.Helper()

	actual := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_billettholder_emails
		WHERE billettholder_id = ? AND email = ? COLLATE NOCASE
	`, billettholderID, email)
	if actual != expected {
		t.Fatalf("billettholder email count mismatch for %d/%s\nexpected: %d\nactual:   %d", billettholderID, email, expected, actual)
	}
}

func assertAdminBillettholderEmailIDCount(t testing.TB, db *sql.DB, emailID int, expected int) {
	t.Helper()

	actual := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_billettholder_emails
		WHERE id = ?
	`, emailID)
	if actual != expected {
		t.Fatalf("billettholder email id count mismatch for %d\nexpected: %d\nactual:   %d", emailID, expected, actual)
	}
}

func assertAdminManualEmailCount(t testing.TB, db *sql.DB, billettholderID int, expected int) {
	t.Helper()

	actual := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_billettholder_emails
		WHERE billettholder_id = ? AND kind = ?
	`, billettholderID, models.BillettholderEmailKindManual)
	if actual != expected {
		t.Fatalf("manual billettholder email count mismatch for %d\nexpected: %d\nactual:   %d", billettholderID, expected, actual)
	}
}

func assertAdminBillettholderUserAssociationCount(t testing.TB, db *sql.DB, billettholderID int, userID int, expected int) {
	t.Helper()

	actual := testutil.QueryInt(t, db, `
		SELECT COUNT(*)
		FROM relation_billettholdere_users
		WHERE billettholder_id = ? AND user_id = ?
	`, billettholderID, userID)
	if actual != expected {
		t.Fatalf("billettholder user association count mismatch for %d/%d\nexpected: %d\nactual:   %d", billettholderID, userID, expected, actual)
	}
}
