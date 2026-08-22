package formsubmission

import (
	"strings"
	"testing"

	"github.com/Regncon/conorganizer/models"
	"github.com/Regncon/conorganizer/testutil"
	"github.com/Regncon/conorganizer/testutil/bdd"
	"github.com/Regncon/conorganizer/testutil/templtest"
)

const (
	eventE1 = "E1"
	eventE2 = "E2"
	eventE3 = "E3"
	eventE4 = "E4"

	puljeP1 = "P1"
	puljeP2 = "P2"
	puljeP3 = "P3"
	puljeP4 = "P4"

	idPlayerAssigned             = 1
	idGMAssigned                 = 2
	idNotVeryInterested          = 3
	idUnassigned                 = 4
	idSameEventAssignee          = 5
	idGMPlayer                   = 6
	idGMAndPlayerDifferentEvents = 7
	idGMOnlyVeryInterestedOther  = 8
)

// FirstChoice rules:
// - The event you were assigned to does not mark you as FirstChoice there because you are
//   already placed; the flag only becomes meaningful when you appear in other event interest lists.
// - GM status alone never sets FirstChoice; only player assignments do, and GM-only involvement
//   should never show you as FirstChoice.

func TestGetInterestsForEvent_FirstChoiceRules(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Given a mix of player assignments, GM assignments, and interest levels across events.",
		When:  "When interests and assignees are fetched for each event.",
		Then:  "Then first-choice flags only come from high-interest player assignments in other events.",
	})

	// Given
	expectedAbsentFromE1 := idPlayerAssigned
	expectedE1AssigneeFirstChoice := []firstChoiceCase{
		{id: idPlayerAssigned, want: false, name: "player assigned in current event should not mark first choice"},
	}
	expectedAbsentFromE2 := idSameEventAssignee
	expectedPresentInE2 := []int{
		idPlayerAssigned,
		idGMAssigned,
		idNotVeryInterested,
		idUnassigned,
		idGMPlayer,
		idGMAndPlayerDifferentEvents,
		idGMOnlyVeryInterestedOther,
	}
	expectedE2FirstChoice := []firstChoiceCase{
		{id: idPlayerAssigned, want: true, name: "player assigned to other event"},
		{id: idGMAssigned, want: false, name: "gm assigned to other event"},
		{id: idNotVeryInterested, want: false, name: "not very interested"},
		{id: idUnassigned, want: false, name: "no assignment"},
		{id: idGMPlayer, want: true, name: "gm+player with very interested"},
		{id: idGMOnlyVeryInterestedOther, want: false, name: "gm-only with very interested in other event"},
	}
	expectedAbsentFromE3 := idGMPlayer
	expectedE3FirstChoice := []firstChoiceCase{
		{id: idPlayerAssigned, want: true, name: "player assigned to other event"},
		{id: idGMAssigned, want: false, name: "gm assigned to other event"},
	}
	expectedE4FirstChoice := []firstChoiceCase{
		{id: idPlayerAssigned, want: true, name: "player assigned to other event"},
		{id: idGMAssigned, want: false, name: "gm assigned to other event"},
		{id: idUnassigned, want: false, name: "no assignment"},
	}

	db, logger := testutil.CreateTestDBAndLogger(t, "first-choice")

	seedBaseTables(t, db)
	seedBillettholdere(t, db, append(
		playerFixtures(),
		gmFixtures()...,
	))
	seedInterests(t, db, append(
		interestsForE1(),
		append(interestsForE2(), append(interestsForE3(), interestsForE4()...)...)...,
	))
	assignmentRows := append(assignmentsE1(), assignmentsE2()...)
	assignmentRows = append(assignmentRows, assignmentsE3()...)
	seedAssignments(t, db, assignmentRows)

	// When
	actualE1Interests := interestIndexForEvent(t, eventE1, db, logger)
	actualE1Assignees := assigneeIndexForEvent(t, eventE1, db, logger)
	actualE2Interests := interestIndexForEvent(t, eventE2, db, logger)
	actualE3Interests := interestIndexForEvent(t, eventE3, db, logger)
	actualE4Interests := interestIndexForEvent(t, eventE4, db, logger)

	// Then
	// E1 inclusion check confirms same-event assignees are excluded from interests.
	t.Run("E1 includes/excludes correct billettholders", func(t *testing.T) {
		expectAbsent(t, actualE1Interests, expectedAbsentFromE1, "expected assigned-to-same-event billettholder to be excluded for E1")
	})

	// E1 first-choice check confirms current-event assignees are not marked as first choice.
	t.Run("E1 assignees should not show first-choice for current event", func(t *testing.T) {
		for _, tc := range expectedE1AssigneeFirstChoice {
			expectFirstChoice(t, actualE1Assignees, tc)
		}
	})

	// E2 inclusion checks confirm same-event assignees are excluded from interests.
	t.Run("E2 includes/excludes correct billettholders", func(t *testing.T) {
		expectAbsent(t, actualE2Interests, expectedAbsentFromE2, "expected assigned-to-same-event billettholder to be excluded")
		for _, expectedID := range expectedPresentInE2 {
			expectPresent(t, actualE2Interests, expectedID, "expected billettholder to be returned")
		}
	})

	// E2 first-choice checks focus on the CASE logic in queryFirstChoice:
	// - Highest interest + assigned as player in a different event => FirstChoice should be true.
	// - GM-only in a different event should NOT count as FirstChoice.
	// - Any lower interest should NOT be FirstChoice, even if assigned elsewhere.
	// - No assignment at all should NOT be FirstChoice.
	t.Run("E2 first-choice rules", func(t *testing.T) {
		for _, tc := range expectedE2FirstChoice {
			expectFirstChoice(t, actualE2Interests, tc)
		}
	})

	// E3 inclusion check confirms same-event assignees are excluded from interests.
	t.Run("E3 includes/excludes correct billettholders", func(t *testing.T) {
		expectAbsent(t, actualE3Interests, expectedAbsentFromE3, "expected assigned-to-same-event billettholder to be excluded for E3")
	})

	// E3 first-choice checks re-run the same CASE rules against a different event to confirm
	// the logic is not accidentally tied to E2-only data setup.
	t.Run("E3 first-choice rules", func(t *testing.T) {
		for _, tc := range expectedE3FirstChoice {
			expectFirstChoice(t, actualE3Interests, tc)
		}
	})

	// E4 first-choice checks cover an interest mix with an explicit "no assignment" case to ensure
	// the FirstChoice flag remains false when the participant has no cross-event player assignment.
	t.Run("E4 first-choice rules", func(t *testing.T) {
		for _, tc := range expectedE4FirstChoice {
			expectFirstChoice(t, actualE4Interests, tc)
		}
	})
}

func TestPuljeSection_FlagsMissingDm(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der arrangementet ikke har en tildelt spilleder.",
		When:  "Når puljeseksjonen på godkjenningssiden rendres.",
		Then:  "Så skal den vise 'Mangler spilleder'.",
	})

	pulje := models.PuljeRow{ID: models.PuljeFredagKveld, Name: "Fredag Kveld"}
	assignees := []InterestWithHolder{
		{BillettholderId: 1, EventId: "evA", PuljeId: string(models.PuljeFredagKveld), IsGamemaster: false, FirstName: "Kari", LastName: "Nordmann"},
	}

	doc := templtest.Render(t, puljeSection("evA", pulje, nil, assignees, true, nil))
	html, err := doc.Html()
	if err != nil {
		t.Fatalf("render html: %v", err)
	}
	if !strings.Contains(html, "Mangler spilleder") {
		t.Fatalf("pulje without a DM should show 'Mangler spilleder'")
	}
}

func TestPuljeSection_NoMissingDmFlagWhenGmAssigned(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "Gitt en pulje der arrangementet har en tildelt spilleder.",
		When:  "Når puljeseksjonen på godkjenningssiden rendres.",
		Then:  "Så skal den ikke vise 'Mangler spilleder'.",
	})

	pulje := models.PuljeRow{ID: models.PuljeFredagKveld, Name: "Fredag Kveld"}
	assignees := []InterestWithHolder{
		{BillettholderId: 9, EventId: "evA", PuljeId: string(models.PuljeFredagKveld), IsGamemaster: true, FirstName: "Game", LastName: "Master"},
	}

	doc := templtest.Render(t, puljeSection("evA", pulje, nil, assignees, true, nil))
	html, err := doc.Html()
	if err != nil {
		t.Fatalf("render html: %v", err)
	}
	if strings.Contains(html, "Mangler spilleder") {
		t.Fatalf("pulje with a GM should not show 'Mangler spilleder'")
	}
}
