package event

import (
	"slices"
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

func TestGetSelectedInterestState_ReturnsConfirmedAssignmentsAndIgnoresSolverResult(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt manuelle tildelingar, påmeldingar, eit skjult solverresultat og ei spilledertildeling i same pulje.",
		When:  "Når tilstanden til interesseveljaren blir henta.",
		Then:  "Så skal berre stadfesta tildelingar og spillederrolla blokkere vanlege interesser.",
	})

	// Given
	expectedInterest := models.InterestLevelHigh
	expectedPlayerEventIDs := []string{"event-interest-event", "registered-event"}
	expectedGamemasterEventIDs := []string{"gm-event"}
	expectedAssignedToEvent := true
	db := createEventInterestTestDB(t)
	fixture := seedEventInterestUpdateFixture(t, db, models.PuljeStatusOpen, expectedInterest)
	seedEventInterestBlockingEvent(t, db, "registered-event", "Registered Event", fixture.puljeID)
	seedEventInterestBlockingEvent(t, db, "solver-event", "Solver Event", fixture.puljeID)
	seedEventInterestBlockingEvent(t, db, "gm-event", "GM Event", fixture.puljeID)
	for _, assignment := range []struct {
		eventID string
		role    models.EventPlayerRole
		source  models.EventPlayerSource
	}{
		{eventID: fixture.eventID, role: models.EventPlayerRolePlayer, source: models.EventPlayerSourceManual},
		{eventID: "registered-event", role: models.EventPlayerRolePlayer, source: models.EventPlayerSourceRegistration},
		{eventID: "solver-event", role: models.EventPlayerRolePlayer, source: models.EventPlayerSourceSolver},
		{eventID: "gm-event", role: models.EventPlayerRoleGM, source: models.EventPlayerSourceManual},
	} {
		mustExecEventInterestTest(t, db, `
			INSERT INTO relation_events_players (event_id, pulje_id, billettholder_id, role, source)
			VALUES (?, ?, ?, ?, ?)
		`, assignment.eventID, fixture.puljeID, fixture.billettholderID, assignment.role, assignment.source)
	}

	// When
	state, err := getSelectedInterestState(fixture.eventID, fixture.billettholderID, string(fixture.puljeID), db)
	actualPlayerEventIDs := assignmentEventIDs(state.PlayerAssignments)
	actualGamemasterEventIDs := assignmentEventIDs(state.GamemasterAssignments)

	// Then
	if err != nil {
		t.Fatalf("expected selected interest state to load: %v", err)
	}
	if state.InterestLevel != expectedInterest {
		t.Fatalf("interest mismatch\nexpected: %s\nactual:   %s", expectedInterest, state.InterestLevel)
	}
	if state.IsAssignedToEvent != expectedAssignedToEvent {
		t.Fatalf("exact assignment mismatch\nexpected: %t\nactual:   %t", expectedAssignedToEvent, state.IsAssignedToEvent)
	}
	if !slices.Equal(actualPlayerEventIDs, expectedPlayerEventIDs) {
		t.Fatalf("player assignments mismatch\nexpected: %v\nactual:   %v", expectedPlayerEventIDs, actualPlayerEventIDs)
	}
	if !slices.Equal(actualGamemasterEventIDs, expectedGamemasterEventIDs) {
		t.Fatalf("gamemaster assignments mismatch\nexpected: %v\nactual:   %v", expectedGamemasterEventIDs, actualGamemasterEventIDs)
	}
	if !state.ordinaryInterestsDisabled() {
		t.Fatal("expected confirmed assignment state to disable ordinary interests")
	}
}

func TestSelectedInterestWarning_WhenPlayerHasMultipleAssignments_ListsLinkedEvents(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt ein billettheldar som er påmeldt fleire arrangement i same pulje.",
		When:  "Når åtvaringa over interesseveljaren blir vist.",
		Then:  "Så skal alle arrangementa visast som lenker slik at påmeldinga er lett å finne.",
	})

	// Given
	expectedHrefs := []string{
		"/event/blood-clocktower?pulje=FredagKveld",
		"/event/cosplay?pulje=FredagKveld",
	}
	expectedText := []string{"Blood in the Clocktower", "Cosplay Competition", "Vanlege interesser kan ikkje endrast"}
	state := selectedInterestState{
		interestParticipationState: interestParticipationState{
			PlayerAssignments: []interestAssignmentLink{
				{EventID: "blood-clocktower", Title: "Blood in the Clocktower"},
				{EventID: "cosplay", Title: "Cosplay Competition"},
			},
		},
	}

	// When
	doc := templtest.Render(t, selectedInterestWarning(state, string(models.PuljeFredagKveld)))
	actualHrefs := templtest.CollectUniqueHrefs(doc)
	actualText := strings.Join(templtest.CollectTexts(doc, "#interest-participation-warning"), " ")

	// Then
	if !slices.Equal(actualHrefs, expectedHrefs) {
		t.Fatalf("warning links mismatch\nexpected: %v\nactual:   %v", expectedHrefs, actualHrefs)
	}
	for _, text := range expectedText {
		if !strings.Contains(actualText, text) {
			t.Fatalf("expected warning text to contain %q\nactual: %s", text, actualText)
		}
	}
}

func assignmentEventIDs(assignments []interestAssignmentLink) []string {
	eventIDs := make([]string, 0, len(assignments))
	for _, assignment := range assignments {
		eventIDs = append(eventIDs, assignment.EventID)
	}
	return eventIDs
}
