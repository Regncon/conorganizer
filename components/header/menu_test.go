package header

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestMenuBillettholderLive_LogsTicketHolderQueryError(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a logged-in user and an unavailable database.",
		When:  "When the live menu loads the user's associated billettholdere.",
		Then:  "Then the structured log includes the underlying database error.",
	})

	// Given
	expectedError := "sql: database is closed"
	db := testutil.CreateTestDB(t, "menu_ticket_holder_query_error")
	if err := db.Close(); err != nil {
		t.Fatalf("close test database: %v", err)
	}
	var logOutput bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logOutput, nil))
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		Id:         "user-id",
		Email:      "user@example.com",
	}

	// When
	templtest.Render(t, MenuBillettholderLive(userInfo, db, logger))

	// Then
	var logEntry map[string]any
	if err := json.Unmarshal(logOutput.Bytes(), &logEntry); err != nil {
		t.Fatalf("decode structured log %q: %v", logOutput.String(), err)
	}
	actualError, ok := logEntry["error"].(string)
	if !ok {
		t.Fatalf("expected structured log to contain a string error field: %s", logOutput.String())
	}
	if !strings.Contains(actualError, expectedError) {
		t.Fatalf("expected error to contain %q\nactual: %q", expectedError, actualError)
	}
}

func TestMenu_AnonymousUserOnlyReceivesPublicNavigation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren ikke er innlogget.",
		When:  "Når hovednavigasjonen vises.",
		Then:  "Så skal brukeren bare få navigasjonslenker til forsiden og innlogging.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "test_room_services")
	expectedHrefs := []string{"/", "/auth"}
	userInfo := requestctx.UserRequestInfo{}

	// When
	doc := templtest.Render(t, Menu(userInfo, db, logger))
	actualHrefs := templtest.CollectUniqueHrefs(doc)

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
}

func TestMenu_LoggedInUserOnlyReceivesUserNavigation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren er innlogget uten adminrettigheter.",
		When:  "Når hovednavigasjonen vises.",
		Then:  "Så skal brukeren bare få navigasjonslenker til forsiden, egen profil, utlogging og vanlege spørsmål.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "test_room_services")
	expectedHrefs := []string{
		"/",
		"/profile",
		"/auth/logout",
		"https://www.regncon.no/vanlege-sporsmal/",
	}
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		IsAdmin:    false,
	}

	// When
	doc := templtest.Render(t, Menu(userInfo, db, logger))
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualExternalLinkIconVisible := doc.Find(`a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon`).Length() > 0

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	if !actualExternalLinkIconVisible {
		t.Fatalf("expected external FAQ link to include external link icon")
	}
}

func TestMenu_AdminUserReceivesUserAndAdminNavigation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt at brukeren er admin.",
		When:  "Når hovednavigasjonen vises.",
		Then:  "Så skal brukeren få navigasjonslenker til forsiden, egen profil, utlogging, adminområdene og vanlege spørsmål.",
	})

	// Given
	db, logger := testutil.CreateTestDBAndLogger(t, "test_room_services")
	expectedHrefs := []string{
		"/",
		"/profile",
		"/auth/logout",
		"/admin",
		"/admin/billettholder/",
		"/admin/approval/",
		"https://www.regncon.no/vanlege-sporsmal/",
	}
	userInfo := requestctx.UserRequestInfo{
		IsLoggedIn: true,
		IsAdmin:    true,
	}

	// When
	doc := templtest.Render(t, Menu(userInfo, db, logger))
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualExternalLinkIconVisible := doc.Find(`a[href="https://www.regncon.no/vanlege-sporsmal/"] .inline-icon`).Length() > 0

	// Then
	templtest.AssertSameHrefs(t, expectedHrefs, actualHrefs)
	if !actualExternalLinkIconVisible {
		t.Fatalf("expected external FAQ link to include external link icon")
	}
}
