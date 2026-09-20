package admin

import (
	"database/sql"
	"encoding/json"
	"html"
	"net/http"
	"net/http/httptest"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestTildeling_DragMovesOnlySelectedPlayerPin(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A billettholder has manual Player pins in A and B in the same pulje.",
		When:  "An admin drags the Player tile from A to C.",
		Then:  "Only the pin in A moves, leaving manual pins in B and C.",
	})

	// Given
	expectedPins := []string{"evB:manual", "evC:manual"}
	db, router := tildelingDragFixture(t)

	// When
	response := postTildelingDrag(t, router, "evA", "evC", "")

	// Then
	if response.Code != http.StatusNoContent {
		t.Fatalf("drag selected pin: %d %s", response.Code, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedPins)
}

func TestTildeling_DragOntoExistingPlayerPinIsNoOp(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A billettholder has manual Player pins in A and B.",
		When:  "An admin drops the Player tile from A onto B.",
		Then:  "Both original pins remain without adding another assignment.",
	})

	// Given
	expectedPins := []string{"evA:manual", "evB:manual"}
	db, router := tildelingDragFixture(t)

	// When
	response := postTildelingDrag(t, router, "evA", "evB", "")

	// Then
	if response.Code != http.StatusNoContent {
		t.Fatalf("drop on existing pin: %d %s", response.Code, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedPins)
}

func TestTildeling_DragRejectsAmbiguousOrStaleSourceWithoutChangingPins(t *testing.T) {
	for _, tc := range []struct {
		name   string
		source string
	}{
		{name: "missing source", source: ""},
		{name: "unknown event", source: "unknown-event"},
		{name: "event without current pin", source: "evD"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			bdd.Behavior(t, bdd.BDD{
				Given: "A billettholder has two manual Player pins and the drag has no valid current source pin.",
				When:  "The drag request is submitted.",
				Then:  "The route rejects the request without changing either pin.",
			})

			// Given
			expectedPins := []string{"evA:manual", "evB:manual"}
			const expectedStatus = http.StatusBadRequest
			db, router := tildelingDragFixture(t)
			testutil.MustExec(t, db, `INSERT INTO events(id,title,intro,description,host_name,email,phone_number,max_players,is_in_puljefordeling) VALUES ('evD','Former source','','','','','',4,1)`)
			testutil.MustExec(t, db, `INSERT INTO relation_event_puljer(event_id,pulje_id,is_in_pulje) VALUES ('evD','FredagKveld',1)`)

			// When
			response := postTildelingDrag(t, router, tc.source, "evC", "")

			// Then
			if response.Code != expectedStatus {
				t.Errorf("invalid drag source: got %d, want %d: %s", response.Code, expectedStatus, response.Body.String())
			}
			assertTildelingDragPins(t, db, expectedPins)
		})
	}
}

func TestTildeling_DragWithoutSourceMovesSolePinForOlderClient(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A billettholder has exactly one manual Player pin in A.",
		When:  "An older client drags it to C without sending the source event.",
		Then:  "The sole pin moves to C.",
	})

	// Given
	expectedPins := []string{"evC:manual"}
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA','FredagKveld',1,'Player','manual')`)

	// When
	response := postTildelingDrag(t, router, "", "evC", "")

	// Then
	if response.Code != http.StatusNoContent {
		t.Fatalf("older-client drag: %d %s", response.Code, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedPins)
}

func TestTildeling_DragDeletedLastManualPinRejectsStaleTile(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An admin has a manual Player tile from A open after its only persisted pin has been removed.",
		When:  "The stale tile is dragged to C with its manual-source flag.",
		Then:  "The route rejects the stale source and does not create a destination pin.",
	})

	// Given
	const expectedStatus = http.StatusBadRequest
	var expectedPins []string
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA','FredagKveld',1,'Player','manual')`)
	testutil.MustExec(t, db, `DELETE FROM relation_events_players WHERE event_id='evA' AND pulje_id='FredagKveld' AND billettholder_id=1 AND role='Player'`)

	// When
	response := postTildelingDrag(t, router, "evA", "evC", "", true)

	// Then
	if response.Code != expectedStatus {
		t.Errorf("stale manual tile: got %d, want %d: %s", response.Code, expectedStatus, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedPins)
}

func TestTildeling_DragUnsavedPreviewCreatesDestinationPin(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A billettholder is automatically previewed in A and has no persisted Player seats.",
		When:  "An admin drags the preview tile from A to C.",
		Then:  "The route creates a manual pin in C even though the source seat has not been saved.",
	})

	// Given
	expectedPins := []string{"evC:manual"}
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO interests(billettholder_id,event_id,pulje_id,interest_level) VALUES (1,'evA','FredagKveld',?)`, models.InterestLevelHigh)
	emulation, err := puljefordeling.EmulateSeatings(db)
	if err != nil {
		t.Fatalf("preview initial source seat: %v", err)
	}
	if len(emulation.Puljer) != 1 || len(emulation.Puljer[0].Events) == 0 {
		t.Fatal("expected preview containing the source event")
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1`); got != 0 {
		t.Fatal("source preview unexpectedly wrote a persisted assignment")
	}

	// When
	response := postTildelingDrag(t, router, "evA", "evC", "")

	// Then
	if response.Code != http.StatusNoContent {
		t.Fatalf("drag unsaved preview: %d %s", response.Code, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedPins)
}

func TestTildeling_DragStillBlocksGMOverlapWithExplicitSource(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A billettholder is GM in C and has manual Player pins in A and B.",
		When:  "An admin drags their Player tile from A to C with an explicit source event.",
		Then:  "The route directs the admin to Legg til and preserves all Player and GM assignments.",
	})

	// Given
	expectedPins := []string{"evA:manual", "evB:manual"}
	db, router := tildelingDragFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evC','FredagKveld',1,'GM','manual')`)

	// When
	response := postTildelingDrag(t, router, "evA", "evC", "")

	// Then
	if !strings.Contains(response.Body.String(), "Legg til") || strings.Contains(response.Body.String(), "data-confirmation=") {
		t.Fatalf("GM overlap must direct the admin to Legg til: %d %s", response.Code, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedPins)
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND event_id='evC' AND role='GM'`); got != 1 {
		t.Fatal("blocked drag changed the GM assignment")
	}
}

func TestTildeling_DragCapacityConfirmationRetainsSourceAndOtherPins(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A billettholder has manual Player pins in A and B, and C is full with another manual Player.",
		When:  "An admin drags A to C and confirms the returned capacity warning.",
		Then:  "The warning keeps both original pins and carries the dragged source into the retry, which moves only A to C.",
	})

	// Given
	expectedBefore := []string{"evA:manual", "evB:manual"}
	expectedAfter := []string{"evB:manual", "evC:manual"}
	db, router := tildelingDragFixture(t)
	seedTildelingAdult(t, db, 2)
	testutil.MustExec(t, db, `UPDATE events SET max_players=1 WHERE id='evC'`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evC','FredagKveld',2,'Player','manual')`)

	// When
	warning := postTildelingDrag(t, router, "evA", "evC", "", true)
	assertTildelingDragPins(t, db, expectedBefore)
	if !strings.Contains(warning.Body.String(), "kapasitet") {
		t.Fatalf("drag into full event must show capacity warning: %d %s", warning.Code, warning.Body.String())
	}
	confirmation := confirmationFromResponse(t, warning)
	source := tildelingDragSourceFromResponse(t, warning)
	if source != "evA" {
		t.Fatalf("retry source: got %q, want evA", source)
	}
	manualSource := tildelingDragManualSourceFromResponse(t, warning)
	if !manualSource {
		t.Fatal("capacity retry lost the manual-source flag")
	}
	response := postTildelingDrag(t, router, source, "evC", confirmation, manualSource)

	// Then
	if response.Code != http.StatusNoContent {
		t.Fatalf("confirmed drag: %d %s", response.Code, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedAfter)
	assertTildelingManualPins(t, db, "evC", []int{1, 2})
}

func TestTildeling_DragCapacityConfirmationCannotSwitchSourcePin(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A capacity warning confirms moving the manual Player pin in A to a full event C while another pin remains in B.",
		When:  "The confirmation is submitted for moving B instead of A.",
		Then:  "The route preserves both pins and requires a fresh confirmation, which moves B while retaining A.",
	})

	// Given
	expectedBefore := []string{"evA:manual", "evB:manual"}
	expectedAfter := []string{"evA:manual", "evC:manual"}
	db, router := tildelingDragFixture(t)
	seedTildelingAdult(t, db, 2)
	testutil.MustExec(t, db, `UPDATE events SET max_players=1 WHERE id='evC'`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evC','FredagKveld',2,'Player','manual')`)
	warning := postTildelingDrag(t, router, "evA", "evC", "")
	originalConfirmation := confirmationFromResponse(t, warning)

	// When
	refreshed := postTildelingDrag(t, router, "evB", "evC", originalConfirmation)

	// Then
	assertTildelingDragPins(t, db, expectedBefore)
	if refreshed.Code != http.StatusOK {
		t.Fatalf("changed drag source must refresh the warning: %d %s", refreshed.Code, refreshed.Body.String())
	}
	freshConfirmation := confirmationFromResponse(t, refreshed)
	if freshConfirmation == originalConfirmation {
		t.Fatal("capacity confirmation was reusable for a different source pin")
	}
	source := tildelingDragSourceFromResponse(t, refreshed)
	if source != "evB" {
		t.Fatalf("fresh retry source: got %q, want evB", source)
	}
	response := postTildelingDrag(t, router, source, "evC", freshConfirmation)
	if response.Code != http.StatusNoContent {
		t.Fatalf("confirm changed source: %d %s", response.Code, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedAfter)
	assertTildelingManualPins(t, db, "evC", []int{1, 2})
}

func TestTildeling_DragUnsavedPreviewAgeWarningRetainsSource(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A minor has an unsaved automatic Player preview in A and no persisted pins.",
		When:  "An admin drags the preview to adults-only event C and confirms the age warning.",
		Then:  "The warning carries the source event into the retry and the confirmed request pins the minor in C.",
	})

	// Given
	expectedPins := []string{"evC:manual"}
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `UPDATE billettholdere SET is_over_18=0 WHERE id=1`)
	testutil.MustExec(t, db, `UPDATE events SET age_group=? WHERE id='evC'`, models.AgeGroupAdultsOnly)
	testutil.MustExec(t, db, `INSERT INTO interests(billettholder_id,event_id,pulje_id,interest_level) VALUES (1,'evA','FredagKveld',?)`, models.InterestLevelHigh)

	// When
	warning := postTildelingDrag(t, router, "evA", "evC", "")

	// Then
	assertTildelingDragPins(t, db, nil)
	var signals struct {
		FromEventID string `json:"ageWarningFromEventId"`
		Text        string `json:"ageWarningText"`
	}
	for _, line := range strings.Split(warning.Body.String(), "\n") {
		if encoded, ok := strings.CutPrefix(line, "data: signals "); ok {
			if err := json.Unmarshal([]byte(encoded), &signals); err != nil {
				t.Fatalf("decode age warning signals: %v", err)
			}
		}
	}
	if signals.FromEventID != "evA" || !strings.Contains(signals.Text, "under 18") {
		t.Fatalf("age warning must retain source evA and explain the age limit: %+v; response %s", signals, warning.Body.String())
	}
	response := postApprovalSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 1, "evC", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":false,"assignmentFromEventId":"`+signals.FromEventID+`","assignmentAgeConfirmed":true`)
	if response.Code != http.StatusNoContent {
		t.Fatalf("confirmed age warning: %d %s", response.Code, response.Body.String())
	}
	assertTildelingDragPins(t, db, expectedPins)
}

func tildelingDragFixture(t *testing.T) (*sql.DB, http.Handler) {
	t.Helper()
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES
		('evA','FredagKveld',1,'Player','manual'),
		('evB','FredagKveld',1,'Player','manual')`)
	return db, router
}

func postTildelingDrag(t *testing.T, router http.Handler, source, destination, confirmation string, fromManualSeat ...bool) *httptest.ResponseRecorder {
	t.Helper()
	extra := `,"assignmentRole":"Player","assignmentFromAddMenu":false`
	if len(fromManualSeat) > 0 {
		extra += `,"assignmentFromManualSeat":` + strconv.FormatBool(fromManualSeat[0])
	}
	if source != "" {
		extra += `,"assignmentFromEventId":"` + source + `"`
	}
	if confirmation != "" {
		extra += `,"assignmentConfirmation":"` + confirmation + `"`
	}
	return postApprovalSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 1, destination, "FredagKveld", extra)
}

func assertTildelingDragPins(t *testing.T, db *sql.DB, expected []string) {
	t.Helper()
	rows, err := db.Query(`SELECT event_id,source FROM relation_events_players WHERE billettholder_id=1 AND pulje_id='FredagKveld' AND role='Player' ORDER BY event_id`)
	if err != nil {
		t.Fatalf("read dragged billettholder's pins: %v", err)
	}
	defer rows.Close()
	var actual []string
	for rows.Next() {
		var eventID, source string
		if err := rows.Scan(&eventID, &source); err != nil {
			t.Fatalf("read pin: %v", err)
		}
		actual = append(actual, eventID+":"+source)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("read pins: %v", err)
	}
	if !slices.Equal(actual, expected) {
		t.Fatalf("Player pins after drag: got %v, want %v", actual, expected)
	}
}

func tildelingDragSourceFromResponse(t *testing.T, response *httptest.ResponseRecorder) string {
	t.Helper()
	clickAttributes := regexp.MustCompile(`data-on:click="([^"]+)"`).FindAllStringSubmatch(response.Body.String(), -1)
	sourceAssignment := regexp.MustCompile(`\$assignmentFromEventId\s*=\s*"([^"]*)"`)
	for _, attribute := range clickAttributes {
		if match := sourceAssignment.FindStringSubmatch(html.UnescapeString(attribute[1])); len(match) == 2 {
			return match[1]
		}
	}
	t.Fatalf("rendered retry action omitted the dragged source: %s", response.Body.String())
	return ""
}

func tildelingDragManualSourceFromResponse(t *testing.T, response *httptest.ResponseRecorder) bool {
	t.Helper()
	clickAttributes := regexp.MustCompile(`data-on:click="([^"]+)"`).FindAllStringSubmatch(response.Body.String(), -1)
	manualSourceAssignment := regexp.MustCompile(`\$assignmentFromManualSeat\s*=\s*(true|false)`)
	for _, attribute := range clickAttributes {
		if match := manualSourceAssignment.FindStringSubmatch(html.UnescapeString(attribute[1])); len(match) == 2 {
			return match[1] == "true"
		}
	}
	t.Fatalf("rendered retry action omitted the manual-source flag: %s", response.Body.String())
	return false
}
