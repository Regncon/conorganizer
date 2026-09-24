package solver

import (
	"testing"

	"github.com/Regncon/conorganizer/service/puljefordeling/solver/model"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

func TestSolveSlot_RecordsScoreBreakdownForSeatedPlayers(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "two players want the same single top-choice seat in the first slot, and both want a roomy game in the second",
		When:  "both slots are solved",
		Then:  "the second slot records the loser's miss and never-seated bumps, and the winner's satisfied band",
	})

	// Given
	expectedLoser := model.ScoreBreakdown{Score: 5, Band: bandUnsatVeldig, Misses: 1, MissBonus: missStep, NeverSeatedBump: neverSeatedBump, Total: bandUnsatVeldig + missStep + neverSeatedBump}
	expectedWinner := model.ScoreBreakdown{Score: 5, Satisfied: true, Band: bandSatVeldig, Total: bandSatVeldig}
	s1 := slot("s1", event("e1", 1))
	s2 := slot("s2", event("e2", 2))
	both := func(id string) model.Player {
		return model.Player{ID: id, Name: id, Prefs: map[string]map[string]model.Score{
			"s1": {"e1": 5},
			"s2": {"e2": 5},
		}}
	}
	players := []model.Player{both("p1"), both("p2")}
	state := NewState(2026, weekendOf(s1, s2))

	// When
	first := state.SolveSlot(s1, players)
	second := state.SolveSlot(s2, players)

	// Then
	winner := first.Assignments["e1"][0]
	loser := first.Unassigned[0]
	assertScoreBreakdown(t, second, loser, expectedLoser)
	assertScoreBreakdown(t, second, winner, expectedWinner)
}

func TestSolveSlotFixed_PinnedPlayersHaveNoScoreBreakdown(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "a player pinned into a game by an admin",
		When:  "the slot is solved",
		Then:  "no score breakdown is recorded, because the solver never valued that seat",
	})

	// Given
	s1 := slot("s1", event("e1", 2))
	players := []model.Player{player("p1", prefs("s1", map[string]model.Score{"e1": 5}))}
	state := NewState(2026, weekendOf(s1))

	// When
	result := state.SolveSlotFixed(s1, players, map[string][]string{"p1": {"e1"}})

	// Then
	if _, ok := result.Scores["p1"]; ok {
		t.Fatalf("expected no score breakdown for a pinned player, got %+v", result.Scores["p1"])
	}
}

func TestScoreBreakdown_TotalMatchesAdjustScore(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "an unsatisfied, never-seated GM with two prior misses",
		When:  "their top-choice seat is scored",
		Then:  "the breakdown lists every bump and its total equals the solver's edge weight",
	})

	// Given
	expected := model.ScoreBreakdown{Score: 5, Band: bandUnsatVeldig, Misses: 2, MissBonus: 2 * missStep, NeverSeatedBump: neverSeatedBump, DMBump: dmBump, Total: adjustScore(5, false, true, true, 2)}

	// When
	actual := scoreBreakdown(5, false, true, true, 2)

	// Then
	if actual != expected {
		t.Fatalf("expected %+v, got %+v", expected, actual)
	}
}

func assertScoreBreakdown(t *testing.T, result model.SlotResult, playerID string, expected model.ScoreBreakdown) {
	t.Helper()
	actual, ok := result.Scores[playerID]
	if !ok {
		t.Fatalf("expected a score breakdown for %s", playerID)
	}
	if actual != expected {
		t.Fatalf("%s: expected %+v, got %+v", playerID, expected, actual)
	}
}
