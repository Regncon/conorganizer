package event

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/Regncon/conorganizer/components/event_components"
	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/authctx"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/service/requestctx"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
	"github.com/go-chi/chi/v5"
)

func TestSelectedInterest_SwitchingBillettholderReplacesNoticeWithChoices(t *testing.T) {
	for _, assignment := range []struct{ role, source string }{{"GM", "solver"}, {"Player", "manual"}} {
		t.Run(assignment.role+"/"+assignment.source, func(t *testing.T) {
			bdd.Behavior(t, bdd.BDD{
				Given: "An assigned billettholder and another without an assignment.",
				When:  "The selected billettholder changes while the cookie still identifies the previous holder.",
				Then:  "The HTML changes from the assigned event notice to the new holder's interest choices.",
			})
			// Given
			expectedTitle := `Assigned "Event"`
			db := createEventInterestTestDB(t)
			fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
			seedNoticeBillettholder(t, db, 902, true)
			seedNoticeAssignment(t, db, fixture, assignment.role, assignment.source, expectedTitle)

			// When
			assigned := decodeInterestContent(t, requestInterestContent(t, db, fixture, fixture.billettholderID))
			unassigned := decodeInterestContent(t, requestInterestContent(t, db, fixture, 902))

			// Then
			if !assigned.ShowAssigned || assigned.CanChoose || assigned.Assigned.Title != expectedTitle || assigned.Assigned.EventID != "assigned-event" || assigned.Assigned.Role != assignment.role {
				t.Fatalf("unexpected assigned event content: %+v", assigned)
			}
			if unassigned.ShowAssigned || !unassigned.CanChoose || unassigned.Assigned != (noticeAssignment{}) || unassigned.Interest != models.InterestLevelNone {
				t.Fatalf("previous selection was not cleared: %+v", unassigned)
			}
		})
	}
}

func TestSelectedInterest_PuljefordelingAssignmentAppearsForCurrentBillettholder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "The current billettholder differs from the one in the request cookie.",
		When:  "Puljefordeling assigns the current billettholder and the browser refreshes notices.",
		Then:  "The response shows their manual assignment without changing the selected billettholder.",
	})
	// Given
	expectedBillettholderID := 902
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	seedNoticeBillettholder(t, db, expectedBillettholderID, true)

	// When
	if err := puljefordeling.AddManualSeat(db, fixture.puljeID, fixture.eventID, expectedBillettholderID); err != nil {
		t.Fatal(err)
	}
	actual := decodeInterestContent(t, requestInterestContent(t, db, fixture, expectedBillettholderID))

	// Then
	if !actual.ShowAssigned || actual.CanChoose || actual.Assigned.EventID != fixture.eventID || actual.Assigned.Role != "Player" {
		t.Fatalf("expected the current billettholder's manual assignment, got %+v", actual)
	}
}

func TestSelectedInterest_PuljefordelingRemovalClearsCurrentAssignment(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "The current billettholder has a manual assignment and the request cookie identifies someone else.",
		When:  "Puljefordeling removes that assignment and the browser refreshes notices.",
		Then:  "The assignment notice and its event data clear for the current billettholder.",
	})
	// Given
	expectedAssigned := noticeAssignment{}
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	seedNoticeBillettholder(t, db, 902, true)
	mustExecEventInterestTest(t, db, `INSERT INTO relation_events_players(event_id, pulje_id, billettholder_id, role, source) VALUES (?, ?, 902, 'Player', 'manual')`, fixture.eventID, fixture.puljeID)

	// When
	if err := puljefordeling.RemoveManualSeat(db, fixture.puljeID, fixture.eventID, 902); err != nil {
		t.Fatal(err)
	}
	actual := decodeInterestContent(t, requestInterestContent(t, db, fixture, 902))

	// Then
	if actual.ShowAssigned || !actual.CanChoose || actual.Assigned != expectedAssigned {
		t.Fatalf("expected removed assignment to clear, got %+v", actual)
	}
}

func TestSelectedInterest_SolverPlayerDoesNotShowAssignmentNotice(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A billettholder has a persisted solver placement as a player.",
		When:  "Their notices are refreshed.",
		Then:  "They keep interest choices without an assignment notice.",
	})
	// Given
	expectedAssigned := noticeAssignment{}
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	seedNoticeAssignment(t, db, fixture, "Player", "solver", "Solver event")

	// When
	actual := decodeInterestContent(t, requestInterestContent(t, db, fixture, fixture.billettholderID))

	// Then
	if actual.ShowAssigned || !actual.CanChoose || actual.Assigned != expectedAssigned || actual.Interest != models.InterestLevelHigh {
		t.Fatalf("solver player should keep interest choices, got %+v", actual)
	}
}

func TestSelectedInterest_SwitchingMinorAndAdultReplacesAgeNoticeWithChoices(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An adults-only event and associated adult and minor billettholdere.",
		When:  "The minor is selected and then the adult is selected again.",
		Then:  "The age notice appears and clears while interest selection remains correct.",
	})
	// Given
	expectedAdultInterest := models.InterestLevelHigh
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, expectedAdultInterest)
	seedNoticeBillettholder(t, db, 902, false)
	mustExecEventInterestTest(t, db, `INSERT OR IGNORE INTO age_groups(age_group) VALUES (?)`, models.AgeGroupAdultsOnly)
	mustExecEventInterestTest(t, db, `UPDATE events SET age_group = ? WHERE id = ?`, models.AgeGroupAdultsOnly, fixture.eventID)

	// When
	minor := decodeInterestContent(t, requestInterestContent(t, db, fixture, 902))
	adult := decodeInterestContent(t, requestInterestContent(t, db, fixture, fixture.billettholderID))

	// Then
	if !minor.ShowUnder18 || minor.CanChoose || adult.ShowUnder18 || !adult.CanChoose || adult.Interest != expectedAdultInterest {
		t.Fatalf("incorrect age or interest content: minor=%+v adult=%+v", minor, adult)
	}
}

func TestSelectedInterest_AssignmentOnlyBlocksItsOwnPulje(t *testing.T) {
	// Given
	expectedCanChoose := true
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	seedNoticeAssignment(t, db, fixture, "GM", "solver", "GM event")
	otherPulje := models.PuljeLordagMorgen
	mustExecEventInterestTest(t, db, `INSERT INTO puljer(id, name, status, start_at, end_at) VALUES (?, 'Lørdag morgen', ?, '2026-10-10T10:00:00+02:00', '2026-10-10T14:00:00+02:00')`, otherPulje, models.PuljeStatusOpen)
	mustExecEventInterestTest(t, db, `INSERT INTO relation_event_puljer(event_id, pulje_id, is_in_pulje, is_published) VALUES (?, ?, 1, 1)`, fixture.eventID, otherPulje)
	fixture.puljeID = otherPulje

	// When
	actual := decodeInterestContent(t, requestInterestContent(t, db, fixture, fixture.billettholderID))

	// Then
	if actual.ShowAssigned || actual.CanChoose != expectedCanChoose {
		t.Fatalf("another pulje's assignment leaked into current selection: %+v", actual)
	}
}

func TestSelectedInterest_UnrelatedBillettholderDoesNotExposeNoticeContent(t *testing.T) {
	// Given
	expectedStatus := http.StatusForbidden
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	seedNoticeBillettholder(t, db, 902, false)
	mustExecEventInterestTest(t, db, `DELETE FROM relation_billettholdere_users WHERE billettholder_id = 902`)

	// When
	response := requestInterestContent(t, db, fixture, 902)

	// Then
	if response.Code != expectedStatus || strings.Contains(response.Body.String(), "datastar-patch-") {
		t.Fatalf("expected forbidden without content patches, got %d: %s", response.Code, response.Body.String())
	}
}

func TestEventInterests_RendersPermanentWrapperBeforeBillettholderSelection(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "The dialog initially has no selected billettholder.",
		When:  "Its server HTML is rendered.",
		Then:  "A protected content wrapper exists, ready to receive the selected billettholder's HTML.",
	})
	// Given
	expectedWrapperCount := 1
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	request := httptest.NewRequest(http.MethodGet, "/event/"+fixture.eventID, nil)

	// When
	doc := templtest.Render(t, event_components.EventInterests(requestctx.UserRequestInfo{}, fixture.eventID, string(fixture.puljeID), "Event", models.AgeGroupAdultsOnly, nil, db, request, testutil.NewTestLogger()))

	// Then
	if doc.Find("#interest-content[data-ignore-morph]").Length() != expectedWrapperCount {
		t.Fatal("expected a protected wrapper before selection is available")
	}
	if doc.Find("#interest-content .interest-already-assigned, #interest-content .interest-under-18, #interest-content .interest-buttons").Length() != 0 {
		t.Fatal("no selected billettholder should render no notice or choices")
	}
}

func TestEventInterests_InitialRenderShowsAssignmentForSelectedBillettholder(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "The selected billettholder cookie identifies a manually assigned player.",
		When:  "The interest dialog is first rendered on the server.",
		Then:  "Their assigned event notice is present and the interest choices are absent before any browser request.",
	})
	// Given
	expectedTitle := "Already assigned event"
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	seedNoticeAssignment(t, db, fixture, "Player", "manual", expectedTitle)
	request := httptest.NewRequest(http.MethodGet, "/event/"+fixture.eventID, nil)
	request.AddCookie(&http.Cookie{Name: requestctx.SelectedBillettholderCookieName, Value: fmt.Sprint(fixture.billettholderID)})
	var doc *goquery.Document
	handler := requestctx.BillettholderSelectionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		doc = templtest.Render(t, event_components.EventInterests(requestctx.UserRequestInfo{}, fixture.eventID, string(fixture.puljeID), "Viewed event", models.AgeGroupAdultsOnly, nil, db, r, testutil.NewTestLogger()))
	}))

	// When
	handler.ServeHTTP(httptest.NewRecorder(), request)

	// Then
	if actualTitle := strings.TrimSpace(doc.Find("#interest-content .clickable-event-card strong").Text()); actualTitle != expectedTitle {
		t.Fatalf("expected assigned event %q in initial HTML, got %q", expectedTitle, actualTitle)
	}
	if doc.Find("#interest-content .interest-buttons").Length() != 0 {
		t.Fatal("assigned billettholder must not receive interest choices in initial HTML")
	}
}

type noticeAssignment struct {
	Title   string
	Role    string
	EventID string
}

type noticeResponse struct {
	ShowAssigned bool
	ShowUnder18  bool
	CanChoose    bool
	Assigned     noticeAssignment
	Interest     models.InterestLevel
}

func requestInterestContent(t *testing.T, db *sql.DB, fixture eventInterestUpdateFixture, selectedID int) *httptest.ResponseRecorder {
	t.Helper()
	router := chi.NewRouter()
	router.Use(requestctx.BillettholderSelectionMiddleware)
	router.Put("/event/api/{idx}/interest/selected-interest", selectedInterestHandler(db, testutil.NewTestLogger(), nil))
	body := fmt.Sprintf(`{"billettHolderId":%d,"puljeId":%q}`, selectedID, fixture.puljeID)
	request := httptest.NewRequest(http.MethodPut, "/event/api/"+fixture.eventID+"/interest/selected-interest", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: requestctx.SelectedBillettholderCookieName, Value: fmt.Sprint(fixture.billettholderID)})
	request = request.WithContext(authctx.WithUserToken(request.Context(), fixture.userExternalID, "event-interest-user@example.com"))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func decodeInterestContent(t *testing.T, response *httptest.ResponseRecorder) noticeResponse {
	t.Helper()
	body := response.Body.String()
	if response.Code != http.StatusOK || !strings.Contains(body, "event: datastar-patch-elements") || !strings.Contains(body, "data: mode replace") || strings.Contains(body, "datastar-patch-signals") {
		t.Fatalf("expected replacement HTML without a signal patch, got %d: %s", response.Code, body)
	}
	var elements []string
	for _, line := range strings.Split(body, "\n") {
		if payload, ok := strings.CutPrefix(line, "data: elements "); ok {
			elements = append(elements, payload)
		}
	}
	html := strings.Join(elements, "\n")
	for _, obsoleteSignal := range []string{"$assignedEvent", "$showAssignedEvent", "$showUnder18", "$canChooseInterest"} {
		if strings.Contains(html, obsoleteSignal) {
			t.Fatalf("notice HTML must not depend on %s", obsoleteSignal)
		}
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatal(err)
	}
	wrapper := doc.Find("#interest-content[data-ignore-morph]")
	if wrapper.Length() != 1 || doc.Find("dialog, .ticket-holder-container, .ticket-holder-pulje-picker").Length() != 0 {
		t.Fatalf("expected only the protected interest content: %s", html)
	}
	actual := noticeResponse{
		ShowAssigned: wrapper.Find(".interest-already-assigned").Length() == 1,
		ShowUnder18:  wrapper.Find(".interest-under-18").Length() == 1,
		CanChoose:    wrapper.Find(".interest-buttons").Length() == 1,
	}
	if actual.ShowAssigned {
		card := wrapper.Find(".clickable-event-card")
		actual.Assigned.Title = strings.TrimSpace(card.Find("strong").Text())
		actual.Assigned.EventID = strings.TrimPrefix(card.Find("a").AttrOr("href", ""), "/event/")
		actual.Assigned.Role = "Player"
		if strings.Contains(wrapper.Find(".interest-already-assigned h3").Text(), "arrangør") {
			actual.Assigned.Role = "GM"
		}
	}
	if err := json.Unmarshal([]byte(wrapper.AttrOr("data-signals:selected-interest-level", "")), &actual.Interest); err != nil {
		t.Fatalf("missing saved interest in content HTML: %v", err)
	}
	return actual
}

func seedNoticeBillettholder(t *testing.T, db *sql.DB, id int, over18 bool) {
	t.Helper()
	mustExecEventInterestTest(t, db, `INSERT INTO billettholdere(id, first_name, last_name, ticket_type_id, ticket_type, is_over_18, order_id, ticket_id) VALUES (?, 'Selected', 'Holder', 1, 'Ticket', ?, 7002, ?)`, id, over18, id+8000)
	mustExecEventInterestTest(t, db, `INSERT INTO relation_billettholdere_users(billettholder_id, user_id) VALUES (?, 501)`, id)
}

func seedNoticeAssignment(t *testing.T, db *sql.DB, fixture eventInterestUpdateFixture, role, source, title string) {
	t.Helper()
	mustExecEventInterestTest(t, db, `
		INSERT INTO events(id, title, intro, description, system, event_type, age_group, event_runtime, host_name, email, phone_number, max_players, beginner_friendly, can_be_run_in_english, status)
		SELECT 'assigned-event', ?, intro, description, 'Assigned system', event_type, age_group, event_runtime, host_name, email, phone_number, max_players, beginner_friendly, can_be_run_in_english, status FROM events WHERE id = ?
	`, title, fixture.eventID)
	mustExecEventInterestTest(t, db, `INSERT INTO relation_events_players(event_id, pulje_id, billettholder_id, role, source) VALUES ('assigned-event', ?, ?, ?, ?)`, fixture.puljeID, fixture.billettholderID, role, source)
}
