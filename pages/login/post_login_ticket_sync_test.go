package login

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/checkIn"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/descope/go-sdk/descope"
)

type ticketSyncCall struct {
	userID string
	email  string
}

func TestPostLogin_RedirectsWithoutWaitingForTicketSync(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at henting av billetter fra CheckIn tar lang tid.",
		When:  "Når en bruker logger inn.",
		Then:  "Så skal brukeren sendes videre med en gang, og billettene hentes for brukeren i bakgrunnen.",
	})

	// Given
	expectedCall := ticketSyncCall{userID: "user-ticket-sync", email: "ungdom@example.com"}
	expectedLocation := "/"
	validator := &fakeSessionValidator{
		sessionOK:    true,
		sessionToken: &descope.Token{ID: expectedCall.userID, JWT: "session-jwt", Claims: map[string]any{"email": expectedCall.email}},
	}
	router := authTestRouterWithDB(t, validator, "post_login_ticket_sync")
	calls := make(chan ticketSyncCall, 1)
	releaseSync := make(chan struct{})
	t.Cleanup(func() { close(releaseSync) })
	replacePostLoginTicketSync(t, func(_ context.Context, userID string, email string, _ *sql.DB, _ *slog.Logger) (checkIn.UserTicketImportResult, error) {
		calls <- ticketSyncCall{userID: userID, email: email}
		<-releaseSync
		return checkIn.UserTicketImportResult{}, nil
	})
	request := httptest.NewRequest(http.MethodGet, "/auth/post-login", nil)
	request.AddCookie(&http.Cookie{Name: authctx.SessionCookieName, Value: "session"})
	request.AddCookie(&http.Cookie{Name: authctx.RefreshCookieName, Value: "refresh"})
	recorder := httptest.NewRecorder()

	// When
	served := make(chan struct{})
	go func() {
		router.ServeHTTP(recorder, request)
		close(served)
	}()

	// Then
	select {
	case <-served:
	case <-time.After(5 * time.Second):
		t.Fatal("post-login waited for the ticket sync to finish")
	}
	if recorder.Code != http.StatusSeeOther || recorder.Header().Get("Location") != expectedLocation {
		t.Fatalf("expected redirect to %q, got status=%d location=%q", expectedLocation, recorder.Code, recorder.Header().Get("Location"))
	}
	select {
	case actualCall := <-calls:
		if actualCall != expectedCall {
			t.Fatalf("ticket sync call mismatch\nexpected: %+v\nactual:   %+v", expectedCall, actualCall)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("expected tickets to be synced after login")
	}
}

// replacePostLoginTicketSync swaps the background "Hent billetter" run so tests
// do not reach CheckIn, and restores it afterwards.
func replacePostLoginTicketSync(t *testing.T, sync func(context.Context, string, string, *sql.DB, *slog.Logger) (checkIn.UserTicketImportResult, error)) {
	t.Helper()

	original := syncUserTicketsFromCheckIn
	syncUserTicketsFromCheckIn = sync
	t.Cleanup(func() { syncUserTicketsFromCheckIn = original })
}
