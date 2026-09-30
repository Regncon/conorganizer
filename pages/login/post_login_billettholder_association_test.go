package login

import (
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/descope/go-sdk/descope"
	"github.com/go-chi/chi/v5"
)

// This reproduces a support ticket: a parent adds their teenager's email as a
// secondary ("Manual") email on a billettholder they bought, but the teen,
// once they log in, cannot "meld interesse" (register interest) even though
// the "Meld interesse" dialog (components/ticket_holder.GetTicketHolders,
// keyed on relation_billettholder_emails) shows them the billettholder.
//
// The access check that gates registering interest
// (pages/event.updateInterest) is keyed on a different table,
// relation_billettholdere_users, which is only populated for a user that
// already existed in `users` at the moment the secondary email was added
// (service/checkIn.AssociateUsersWithBillettholderEmail). Post-login sync
// never backfilled it for a user created afterwards, so the two tables could
// disagree about which billettholdere a user may act for.

func TestPostLogin_WhenSecondaryEmailWasAddedBeforeFirstLogin_GrantsBillettholderAccess(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "En forelder har lagt til tenåringens e-post som sekundær e-post på en billettholder, før tenåringen noen gang har logget inn.",
		When:  "Tenåringen logger inn for første gang.",
		Then:  "Så skal tenåringen kunne endre interessen til billettholderen, ikke bare se den.",
	})

	// Given
	const billettholderID = 9201
	const teenExternalID = "teen-first-login-ext"
	const teenEmail = "teen-first-login@example.com"

	db, logger := testutil.CreateTestDBAndLogger(t, "post_login_billettholder_association")
	seedOwnedBillettholder(t, db, billettholderID)
	addSecondaryBillettholderEmail(t, db, billettholderID, teenEmail)
	assertNoUserExists(t, db, teenEmail)

	validator := &fakeSessionValidator{
		sessionOK: true,
		sessionToken: &descope.Token{
			ID:     teenExternalID,
			JWT:    "session-jwt",
			Claims: map[string]any{"email": teenEmail},
		},
	}
	router := authRouterWithDB(t, db, logger, validator)
	request := httptest.NewRequest(http.MethodGet, "/auth/post-login", nil)
	request.AddCookie(&http.Cookie{Name: authctx.SessionCookieName, Value: "session"})
	request.AddCookie(&http.Cookie{Name: authctx.RefreshCookieName, Value: "refresh"})
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected first login to redirect (status %d), got %d", http.StatusSeeOther, recorder.Code)
	}
	if !canActForBillettholder(t, db, billettholderID, teenExternalID) {
		t.Errorf("expected the teen to be able to act for (register interest on) the billettholder after their first login, matching what they are already shown in the 'Meld interesse' picker")
	}
}

func TestPostLogin_WhenNoSecondaryEmailExists_DoesNotGrantBillettholderAccess(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "En bruker logger inn uten å være knyttet til noen billettholder.",
		When:  "Post-login synkroniserer brukeren.",
		Then:  "Så skal ingen billettholder-tilgang opprettes for brukeren.",
	})

	// Given
	const billettholderID = 9202
	const userExternalID = "unrelated-user-ext"
	const userEmail = "unrelated-user@example.com"

	db, logger := testutil.CreateTestDBAndLogger(t, "post_login_billettholder_no_association")
	seedOwnedBillettholder(t, db, billettholderID)

	validator := &fakeSessionValidator{
		sessionOK: true,
		sessionToken: &descope.Token{
			ID:     userExternalID,
			JWT:    "session-jwt",
			Claims: map[string]any{"email": userEmail},
		},
	}
	router := authRouterWithDB(t, db, logger, validator)
	request := httptest.NewRequest(http.MethodGet, "/auth/post-login", nil)
	request.AddCookie(&http.Cookie{Name: authctx.SessionCookieName, Value: "session"})
	request.AddCookie(&http.Cookie{Name: authctx.RefreshCookieName, Value: "refresh"})
	recorder := httptest.NewRecorder()

	// When
	router.ServeHTTP(recorder, request)

	// Then
	if recorder.Code != http.StatusSeeOther {
		t.Fatalf("expected login to redirect (status %d), got %d", http.StatusSeeOther, recorder.Code)
	}
	if canActForBillettholder(t, db, billettholderID, userExternalID) {
		t.Errorf("expected a user with no matching billettholder email to get no billettholder access")
	}
}

// seedOwnedBillettholder inserts a minimal billettholder row, as if imported
// from a ticket, for the given id.
func seedOwnedBillettholder(t *testing.T, db *sql.DB, billettholderID int) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO billettholdere (
			id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id
		) VALUES (?, 'Teen', 'Owned', 1, 'Ticket', 1, ?, ?)
	`, billettholderID, 7000+billettholderID, 8000+billettholderID)
}

// addSecondaryBillettholderEmail reproduces the raw INSERT that the "add
// secondary email" route (pages/profile/tickets.addEmailToBilettholderRoute
// / pages/admin/billettholder_admin.addEmailToBilettholderRoute) performs
// before calling checkIn.AssociateUsersWithBillettholderEmail. It is run
// without a matching `users` row, exactly like a parent adding a teen's
// email before the teen has ever logged in.
func addSecondaryBillettholderEmail(t *testing.T, db *sql.DB, billettholderID int, email string) {
	t.Helper()
	testutil.MustExec(t, db, `
		INSERT INTO relation_billettholder_emails (billettholder_id, email, kind) VALUES (?, ?, ?)
	`, billettholderID, email, models.BillettholderEmailKindManual)
}

func assertNoUserExists(t *testing.T, db *sql.DB, email string) {
	t.Helper()
	if count := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM users WHERE email = ?`, email); count != 0 {
		t.Fatalf("expected no user row to exist yet for %q, found %d", email, count)
	}
}

// canActForBillettholder mirrors, verbatim, the access check
// pages/event.updateInterest performs before letting a user register
// interest for a billettholder.
func canActForBillettholder(t *testing.T, db *sql.DB, billettholderID int, externalID string) bool {
	t.Helper()
	var hasAccess bool
	err := db.QueryRow(`
		SELECT EXISTS
			(SELECT 1
				FROM relation_billettholdere_users bu
				JOIN users u ON bu.user_id = u.id
				WHERE bu.billettholder_id = ? AND u.external_id = ?)
	`, billettholderID, externalID).Scan(&hasAccess)
	if err != nil {
		t.Fatalf("failed to check billettholder access: %v", err)
	}
	return hasAccess
}

// authRouterWithDB builds the real /auth route wiring against a caller-owned
// test database, so the test can assert on its state after the request.
func authRouterWithDB(t *testing.T, db *sql.DB, logger *slog.Logger, validator authctx.SessionValidator) chi.Router {
	t.Helper()
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	router := chi.NewRouter()
	authenticatedRouter := router.With(authctx.AuthMiddleware(validator, logger))
	if err := SetupAuthRoute(router, authenticatedRouter, db, logger, validator); err != nil {
		t.Fatalf("expected auth route setup to succeed: %v", err)
	}
	return router
}
