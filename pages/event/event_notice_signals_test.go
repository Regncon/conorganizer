package event

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

func TestSelectedInterest_SwitchingBillettholderClearsAssignmentSignals(t *testing.T) {
	for _, assignment := range []struct{ role, source string }{{"GM", "solver"}, {"Player", "manual"}} {
		t.Run(assignment.role+"/"+assignment.source, func(t *testing.T) {
			bdd.Behavior(t, bdd.BDD{
				Given: "An assigned billettholder and another without an assignment.",
				When:  "The selected billettholder changes while the cookie still identifies the previous holder.",
				Then:  "Assignment signals clear and the new holder's interest is returned without an HTML patch.",
			})
			// Given
			expectedTitle := `Assigned "Event"`
			db := createEventInterestTestDB(t)
			fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
			seedNoticeBillettholder(t, db, 902, true)
			seedNoticeAssignment(t, db, fixture, assignment.role, assignment.source, expectedTitle)

			// When
			assigned := decodeNoticeSignals(t, requestNoticeSignals(t, db, fixture, fixture.billettholderID))
			unassigned := decodeNoticeSignals(t, requestNoticeSignals(t, db, fixture, 902))

			// Then
			if !assigned.ShowAssigned || assigned.CanChoose || assigned.Assigned.Title != expectedTitle || assigned.Assigned.EventID != "assigned-event" || assigned.Assigned.Role != assignment.role {
				t.Fatalf("unexpected assignment signals: %+v", assigned)
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
	actual := decodeNoticeSignals(t, requestNoticeSignals(t, db, fixture, expectedBillettholderID))

	// Then
	if !actual.ShowAssigned || actual.CanChoose || actual.Assigned.EventID != fixture.eventID || actual.Assigned.Source != "manual" {
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
	actual := decodeNoticeSignals(t, requestNoticeSignals(t, db, fixture, 902))

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
	actual := decodeNoticeSignals(t, requestNoticeSignals(t, db, fixture, fixture.billettholderID))

	// Then
	if actual.ShowAssigned || !actual.CanChoose || actual.Assigned != expectedAssigned || actual.Interest != models.InterestLevelHigh {
		t.Fatalf("solver player should keep interest choices, got %+v", actual)
	}
}

func TestSelectedInterest_SwitchingMinorAndAdultUpdatesAgeSignal(t *testing.T) {
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
	minor := decodeNoticeSignals(t, requestNoticeSignals(t, db, fixture, 902))
	adult := decodeNoticeSignals(t, requestNoticeSignals(t, db, fixture, fixture.billettholderID))

	// Then
	if !minor.ShowUnder18 || minor.CanChoose || adult.ShowUnder18 || !adult.CanChoose || adult.Interest != expectedAdultInterest {
		t.Fatalf("incorrect age or interest signals: minor=%+v adult=%+v", minor, adult)
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
	actual := decodeNoticeSignals(t, requestNoticeSignals(t, db, fixture, fixture.billettholderID))

	// Then
	if actual.ShowAssigned || actual.CanChoose != expectedCanChoose {
		t.Fatalf("another pulje's assignment leaked into current selection: %+v", actual)
	}
}

func TestSelectedInterest_UnrelatedBillettholderDoesNotExposeNoticeSignals(t *testing.T) {
	// Given
	expectedStatus := http.StatusForbidden
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	seedNoticeBillettholder(t, db, 902, false)
	mustExecEventInterestTest(t, db, `DELETE FROM relation_billettholdere_users WHERE billettholder_id = 902`)

	// When
	response := requestNoticeSignals(t, db, fixture, 902)

	// Then
	if response.Code != expectedStatus || strings.Contains(response.Body.String(), "datastar-patch-") {
		t.Fatalf("expected forbidden without signals, got %d: %s", response.Code, response.Body.String())
	}
}

func TestEventInterests_RendersBothNoticesWithServerSignalBindings(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "The dialog initially has no selected billettholder.",
		When:  "Its server HTML is rendered.",
		Then:  "Both notice layouts and the buttons exist, with server signal bindings preserved across live renders.",
	})
	// Given
	expectedBindings := []string{"$showAssignedEvent", "$showUnder18", "$canChooseInterest"}
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, models.InterestLevelHigh)
	request := httptest.NewRequest(http.MethodGet, "/event/"+fixture.eventID, nil)

	// When
	doc := templtest.Render(t, event_components.EventInterests(requestctx.UserRequestInfo{}, fixture.eventID, string(fixture.puljeID), "Event", models.AgeGroupAdultsOnly, nil, db, request, testutil.NewTestLogger()))

	// Then
	for _, binding := range expectedBindings {
		selection := doc.Find(`[data-show="` + binding + `"]`)
		if selection.Length() != 1 {
			t.Errorf("expected element bound to %s", binding)
		}
	}
	if doc.Find(".interest-already-assigned").Length() != 1 || doc.Find(".interest-under-18").Length() != 1 || doc.Find(".interest-buttons").Length() != 1 {
		t.Fatal("notices and interest buttons must exist before selection changes")
	}
	initial, ok := doc.Find("dialog").Attr("data-signals__ifmissing")
	if !ok || !strings.Contains(initial, `"showAssignedEvent"`) || !strings.Contains(initial, `"assignedEvent"`) {
		t.Fatalf("expected initial notice signals preserved across live renders: %s", initial)
	}
}

type noticeAssignment struct {
	Title    string `json:"title"`
	System   string `json:"system"`
	Role     string `json:"role"`
	Source   string `json:"source"`
	EventID  string `json:"eventId"`
	ImageURL string `json:"imageUrl"`
}

type noticeResponse struct {
	ShowAssigned bool                 `json:"showAssignedEvent"`
	ShowUnder18  bool                 `json:"showUnder18"`
	CanChoose    bool                 `json:"canChooseInterest"`
	Assigned     noticeAssignment     `json:"assignedEvent"`
	Interest     models.InterestLevel `json:"selectedInterestLevel"`
}

func requestNoticeSignals(t *testing.T, db *sql.DB, fixture eventInterestUpdateFixture, selectedID int) *httptest.ResponseRecorder {
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

func decodeNoticeSignals(t *testing.T, response *httptest.ResponseRecorder) noticeResponse {
	t.Helper()
	body := response.Body.String()
	if response.Code != http.StatusOK || strings.Contains(body, "datastar-patch-elements") {
		t.Fatalf("expected signal-only response, got %d: %s", response.Code, body)
	}
	for _, line := range strings.Split(body, "\n") {
		if payload, ok := strings.CutPrefix(line, "data: signals "); ok {
			var fields map[string]json.RawMessage
			if err := json.Unmarshal([]byte(payload), &fields); err != nil {
				t.Fatal(err)
			}
			for _, key := range []string{"showAssignedEvent", "showUnder18", "canChooseInterest", "assignedEvent"} {
				if _, exists := fields[key]; !exists {
					t.Fatalf("missing server signal %s: %s", key, payload)
				}
			}
			var actual noticeResponse
			if err := json.Unmarshal([]byte(payload), &actual); err != nil {
				t.Fatal(err)
			}
			return actual
		}
	}
	t.Fatalf("missing signal patch: %s", body)
	return noticeResponse{}
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
