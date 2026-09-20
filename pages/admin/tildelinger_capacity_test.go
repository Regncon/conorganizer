package admin

import (
	"database/sql"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/service/puljefordeling"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestTildeling_AdditionalPlayerPinsSurviveRepeatedDistribution(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An adult with no GM role has a manual Player pin and interests in four events in one pulje.",
		When:  "An admin confirms two additional Player pins and repeatedly previews and commits distribution.",
		Then:  "All three pins reserve their event's only seat and the adult receives no automatic fourth seat.",
	})

	// Given
	expectedEvents := []string{"evA", "evB", "evC"}
	db, router := tildelingsFixture(t)
	testutil.MustExec(t, db, `UPDATE events SET max_players=1`)
	testutil.MustExec(t, db, `INSERT INTO events(id,title,intro,description,host_name,email,phone_number,max_players,is_in_puljefordeling) VALUES ('evD','Arrangement W','','','','','',1,1)`)
	testutil.MustExec(t, db, `INSERT INTO relation_event_puljer(event_id,pulje_id,is_in_pulje) VALUES ('evD','FredagKveld',1)`)
	seedTildelingAdult(t, db, 2)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA','FredagKveld',1,'Player','manual')`)
	for _, eventID := range []string{"evA", "evB", "evC", "evD"} {
		for _, id := range []int{1, 2} {
			testutil.MustExec(t, db, `INSERT INTO interests(billettholder_id,event_id,pulje_id,interest_level) VALUES (?,?,'FredagKveld',?)`, id, eventID, models.InterestLevelHigh)
		}
	}

	// When
	for _, eventID := range []string{"evB", "evC"} {
		warning := postTildeling(t, router, eventID, "Player", true, "")
		response := postTildeling(t, router, eventID, "Player", true, confirmationFromResponse(t, warning))
		if response.Code != http.StatusNoContent {
			t.Fatalf("confirm additional pin in %s: %d %s", eventID, response.Code, response.Body.String())
		}
	}
	for range 2 {
		emulation, err := puljefordeling.EmulateSeatings(db)
		if err != nil {
			t.Fatalf("preview additional pins: %v", err)
		}
		for _, eventID := range expectedEvents {
			assertTildelingPreviewPins(t, emulation, eventID, []int{1})
		}
		if err := puljefordeling.CommitDistribution(db, models.PuljeFredagKveld); err != nil {
			t.Fatalf("commit additional pins: %v", err)
		}
	}

	// Then
	for _, eventID := range expectedEvents {
		assertTildelingManualPins(t, db, eventID, []int{1})
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND pulje_id='FredagKveld' AND role='Player'`); got != 3 {
		t.Fatalf("expected three Player pins, got %d assignments", got)
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE event_id='evD' AND billettholder_id=2 AND role='Player' AND source='solver'`); got != 1 {
		t.Fatal("the remaining seat must go to the free player")
	}
}

func TestTildeling_AddPlayerPreservesSavedSolverSeatAsManualPin(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An adult has a saved solver Player seat in A, with unrelated saved solver seats for another person and pulje.",
		When:  "An admin confirms adding the adult as Player in B and previews and commits twice.",
		Then:  "A and B remain manual pins, and the confirmation only promotes this adult's seats in the selected pulje.",
	})

	// Given
	expectedEvents := []string{"evA", "evB"}
	db, router := tildelingSolverSeatFixture(t)
	warning := postTildeling(t, router, "evB", "Player", true, "")
	confirmation := confirmationFromResponse(t, warning)
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND pulje_id='FredagKveld' AND source='solver'`); got != 1 {
		t.Fatal("unconfirmed addition changed the existing solver seat")
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND event_id='evB' AND pulje_id='FredagKveld'`); got != 0 {
		t.Fatal("unconfirmed addition wrote the new Player seat")
	}

	// When
	response := postTildeling(t, router, "evB", "Player", true, confirmation)
	if response.Code != http.StatusNoContent {
		t.Fatalf("confirm Player addition: %d %s", response.Code, response.Body.String())
	}
	assertTildelingUnrelatedSolverSeats(t, db)
	for range 2 {
		emulation, err := puljefordeling.EmulateSeatings(db)
		if err != nil {
			t.Fatalf("preview promoted Player seats: %v", err)
		}
		for _, eventID := range expectedEvents {
			assertTildelingPreviewPins(t, emulation, eventID, []int{1})
		}
		if err := puljefordeling.CommitDistribution(db, models.PuljeFredagKveld); err != nil {
			t.Fatalf("commit promoted Player seats: %v", err)
		}
	}

	// Then
	for _, eventID := range expectedEvents {
		assertTildelingManualPins(t, db, eventID, []int{1})
	}
	if !strings.Contains(warning.Body.String(), "blir manuelle") {
		t.Fatalf("confirmation did not explain that saved solver seats become manual: %s", warning.Body.String())
	}
}

func TestTildeling_AddGMPreservesSavedSolverSeatAsManualPin(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An adult has a saved solver Player seat in A.",
		When:  "An admin confirms adding the adult as GM in B and previews and commits twice.",
		Then:  "The Player seat in A becomes a durable manual pin alongside the GM assignment in B.",
	})

	// Given
	expectedPlayerIDs := []int{1}
	const expectedGMs = 1
	db, router := tildelingSolverSeatFixture(t)
	warning := postTildeling(t, router, "evB", "GM", true, "")
	confirmation := confirmationFromResponse(t, warning)
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND pulje_id='FredagKveld' AND source='solver'`); got != 1 {
		t.Fatal("unconfirmed GM addition changed the existing solver seat")
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND role='GM'`); got != 0 {
		t.Fatal("unconfirmed addition wrote the GM assignment")
	}

	// When
	response := postTildeling(t, router, "evB", "GM", true, confirmation)
	if response.Code != http.StatusNoContent {
		t.Fatalf("confirm GM addition: %d %s", response.Code, response.Body.String())
	}
	assertTildelingUnrelatedSolverSeats(t, db)
	for range 2 {
		emulation, err := puljefordeling.EmulateSeatings(db)
		if err != nil {
			t.Fatalf("preview promoted Player seat beside GM: %v", err)
		}
		assertTildelingPreviewPins(t, emulation, "evA", expectedPlayerIDs)
		if err := puljefordeling.CommitDistribution(db, models.PuljeFredagKveld); err != nil {
			t.Fatalf("commit promoted Player seat beside GM: %v", err)
		}
	}

	// Then
	assertTildelingManualPins(t, db, "evA", expectedPlayerIDs)
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND event_id='evB' AND pulje_id='FredagKveld' AND role='GM' AND source='manual'`); got != expectedGMs {
		t.Fatalf("expected %d GM assignment beside the Player pin, got %d", expectedGMs, got)
	}
	if !strings.Contains(warning.Body.String(), "blir manuelle") {
		t.Fatalf("confirmation did not explain that saved solver seats become manual: %s", warning.Body.String())
	}
}

func TestTildeling_FifthManualPlayerWaitsForCapacityConfirmation(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A four-seat event has four manually pinned adults.",
		When:  "An admin adds a fifth adult through Legg til.",
		Then:  "The capacity warning is shown and no assignment changes before confirmation.",
	})

	// Given
	expectedPins := []int{1, 2, 3, 4}
	db, router := tildelingCapacityFixture(t)

	// When
	response := postApprovalSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 5, "evA", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":true`)

	// Then
	assertTildelingManualPins(t, db, "evA", expectedPins)
	if response.Code != http.StatusOK {
		t.Fatalf("expected capacity confirmation, got %d: %s", response.Code, response.Body.String())
	}
	confirmationFromResponse(t, response)
	if !strings.Contains(strings.ToLower(response.Body.String()), "kapasitet") {
		t.Fatalf("confirmation omitted capacity warning: %s", response.Body.String())
	}
}

func TestTildeling_ConfirmedFifthPlayerSurvivesRepeatedDistribution(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "A four-seat event has four manual Player pins and a sixth interested adult.",
		When:  "An admin confirms a fifth pin and previews and commits twice.",
		Then:  "All five manual pins remain without an automatic sixth Player.",
	})

	// Given
	expectedPins := []int{1, 2, 3, 4, 5}
	db, router := tildelingCapacityFixture(t)
	testutil.MustExec(t, db, `INSERT INTO interests(billettholder_id,event_id,pulje_id,interest_level) VALUES (6,'evA','FredagKveld',?)`, models.InterestLevelHigh)
	warning := postApprovalSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 5, "evA", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":true`)
	confirmation := confirmationFromResponse(t, warning)

	// When
	response := postApprovalSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 5, "evA", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":true,"assignmentConfirmation":"`+confirmation+`"`)
	if response.Code != http.StatusNoContent {
		t.Fatalf("confirm fifth pin: %d %s", response.Code, response.Body.String())
	}
	for range 2 {
		emulation, err := puljefordeling.EmulateSeatings(db)
		if err != nil {
			t.Fatalf("preview five pins: %v", err)
		}
		assertTildelingPreviewPins(t, emulation, "evA", expectedPins)
		if err := puljefordeling.CommitDistribution(db, models.PuljeFredagKveld); err != nil {
			t.Fatalf("commit five pins: %v", err)
		}
	}

	// Then
	assertTildelingManualPins(t, db, "evA", expectedPins)
}

func TestTildeling_AgeConfirmationCannotBypassCapacityWarning(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "An adults-only event is full with four manual Player pins.",
		When:  "An admin adds a minor and retries with only the age confirmation flag.",
		Then:  "The combined dialog warns about age and capacity and requires its confirmation before adding the fifth pin.",
	})

	// Given
	expectedPinsBeforeConfirmation := []int{1, 2, 3, 4}
	expectedPinsAfterConfirmation := []int{1, 2, 3, 4, 5}
	db, router := tildelingCapacityFixture(t)
	testutil.MustExec(t, db, `UPDATE events SET age_group=? WHERE id='evA'`, models.AgeGroupAdultsOnly)
	testutil.MustExec(t, db, `UPDATE billettholdere SET is_over_18=0 WHERE id=5`)

	// When
	warning := postApprovalSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 5, "evA", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":true`)
	ageOnly := postApprovalSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 5, "evA", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":true,"assignmentAgeConfirmed":true`)

	// Then
	assertTildelingManualPins(t, db, "evA", expectedPinsBeforeConfirmation)
	for _, expected := range []string{"tildeling-dialog-innhold", "under 18", "kapasitet"} {
		if !strings.Contains(strings.ToLower(warning.Body.String()), expected) {
			t.Errorf("combined warning omitted %q: %s", expected, warning.Body.String())
		}
	}
	if strings.Contains(warning.Body.String(), "ageWarningText") {
		t.Fatal("age and capacity must use one confirmation dialog")
	}
	confirmation := confirmationFromResponse(t, ageOnly)
	response := postApprovalSignals(t, router, http.MethodPost, "/api/puljefordeling/assign", 5, "evA", "FredagKveld", `,"assignmentRole":"Player","assignmentFromAddMenu":true,"assignmentConfirmation":"`+confirmation+`"`)
	if response.Code != http.StatusNoContent {
		t.Fatalf("confirm combined warning: %d %s", response.Code, response.Body.String())
	}
	assertTildelingManualPins(t, db, "evA", expectedPinsAfterConfirmation)
}

func tildelingCapacityFixture(t *testing.T) (*sql.DB, http.Handler) {
	t.Helper()
	db, router := tildelingsFixture(t)
	for id := 2; id <= 6; id++ {
		seedTildelingAdult(t, db, id)
	}
	for id := 1; id <= 4; id++ {
		testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES ('evA','FredagKveld',?,'Player','manual')`, id)
	}
	return db, router
}

func tildelingSolverSeatFixture(t *testing.T) (*sql.DB, http.Handler) {
	t.Helper()
	db, router := tildelingsFixture(t)
	seedTildelingAdult(t, db, 2)
	testutil.MustExec(t, db, `INSERT INTO puljer(id,name,status,start_at,end_at) VALUES ('LordagKveld','Lørdag kveld','Open','2026-09-16T18:00:00Z','2026-09-16T22:00:00Z')`)
	testutil.MustExec(t, db, `INSERT INTO relation_events_players(event_id,pulje_id,billettholder_id,role,source) VALUES
		('evA','FredagKveld',1,'Player','solver'),
		('evC','FredagKveld',2,'Player','solver'),
		('evA','LordagKveld',1,'Player','solver')`)
	return db, router
}

func assertTildelingUnrelatedSolverSeats(t *testing.T, db *sql.DB) {
	t.Helper()
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=2 AND event_id='evC' AND pulje_id='FredagKveld' AND role='Player' AND source='solver'`); got != 1 {
		t.Fatal("addition changed another person's saved solver seat")
	}
	if got := testutil.QueryInt(t, db, `SELECT COUNT(*) FROM relation_events_players WHERE billettholder_id=1 AND event_id='evA' AND pulje_id='LordagKveld' AND role='Player' AND source='solver'`); got != 1 {
		t.Fatal("addition changed the saved solver seat in another pulje")
	}
}

func seedTildelingAdult(t *testing.T, db *sql.DB, id int) {
	t.Helper()
	testutil.MustExec(t, db, `INSERT INTO billettholdere(id,first_name,last_name,ticket_type_id,ticket_type,order_id,ticket_id,is_over_18) VALUES (?,?,'Adult',0,'',0,?,1)`, id, fmt.Sprintf("Player%d", id), id)
}

func assertTildelingManualPins(t *testing.T, db *sql.DB, eventID string, expectedIDs []int) {
	t.Helper()
	rows, err := db.Query(`SELECT billettholder_id, source FROM relation_events_players WHERE event_id=? AND pulje_id='FredagKveld' AND role='Player' ORDER BY billettholder_id`, eventID)
	if err != nil {
		t.Fatalf("read Player pins: %v", err)
	}
	defer rows.Close()
	var actualIDs []int
	for rows.Next() {
		var id int
		var source string
		if err := rows.Scan(&id, &source); err != nil {
			t.Fatalf("scan Player pin: %v", err)
		}
		if source != puljefordeling.SourceManual {
			t.Errorf("event %s has automatic Player %d where only manual pins should remain", eventID, id)
		}
		actualIDs = append(actualIDs, id)
	}
	if err := rows.Err(); err != nil {
		t.Fatalf("read Player pin rows: %v", err)
	}
	if !slices.Equal(actualIDs, expectedIDs) {
		t.Fatalf("manual pins on %s: got %v, want %v", eventID, actualIDs, expectedIDs)
	}
}

func assertTildelingPreviewPins(t *testing.T, emulation puljefordeling.Emulation, eventID string, expectedIDs []int) {
	t.Helper()
	for _, pulje := range emulation.Puljer {
		if pulje.PuljeID != models.PuljeFredagKveld {
			continue
		}
		for _, event := range pulje.Events {
			if event.EventID != eventID {
				continue
			}
			var actualIDs []int
			for _, player := range event.AssignedPlayers {
				if !player.Manual {
					t.Errorf("preview for %s allocated automatic Player %d into pinned capacity", eventID, player.BillettholderID)
				}
				actualIDs = append(actualIDs, player.BillettholderID)
			}
			slices.Sort(actualIDs)
			if !slices.Equal(actualIDs, expectedIDs) {
				t.Fatalf("preview pins on %s: got %v, want %v", eventID, actualIDs, expectedIDs)
			}
			return
		}
	}
	t.Fatalf("event %s missing from preview", eventID)
}
