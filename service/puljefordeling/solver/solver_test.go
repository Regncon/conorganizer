package solver

import (
	"fmt"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/service/puljefordeling/solver/model"
)

// --- helpers ----------------------------------------------------------------

func slot(id string, events ...model.Event) model.Slot {
	return model.Slot{ID: id, Name: id, Events: events}
}

func event(id string, capacity int) model.Event {
	return model.Event{ID: id, Name: id, Capacity: capacity}
}

func player(id string, prefs map[string]map[string]model.Score) model.Player {
	return model.Player{ID: id, Name: id, Prefs: prefs}
}

func prefs(slotID string, scores map[string]model.Score) map[string]map[string]model.Score {
	return map[string]map[string]model.Score{slotID: scores}
}

func assigned(result model.SlotResult, eventID string) []string {
	return result.Assignments[eventID]
}

// weekendOf wraps slots in a minimal Weekend for NewState.
func weekendOf(slots ...model.Slot) model.Weekend {
	return model.Weekend{Slots: slots}
}

// --- adjustScore ------------------------------------------------------------

func TestAdjustScore_SatisfiedAlwaysRaw(t *testing.T) {
	for _, score := range []model.Score{1, 2, 3, 4, 5} {
		got := adjustScore(score, true, true, 1, false)
		if got != int(score) {
			t.Errorf("satisfied score %d: want %d, got %d", score, score, got)
		}
	}
}

func TestAdjustScore_ScarcityBonusDecays(t *testing.T) {
	// Fewer remaining opportunities → larger bonus. Base 10 + max(0, 5-opps).
	cases := []struct {
		opps int
		want int
	}{
		{1, 14},
		{2, 13},
		{3, 12},
		{4, 11},
		{5, 10},
		{10, 10},
	}
	for _, c := range cases {
		got := adjustScore(5, false, false, c.opps, false)
		if got != c.want {
			t.Errorf("unsatisfied 5 with %d opps: want %d, got %d", c.opps, c.want, got)
		}
	}
}

func TestAdjustScore_LateBoostDoublesLowerScores(t *testing.T) {
	for _, score := range []model.Score{1, 2, 3, 4} {
		got := adjustScore(score, false, true, 5, false)
		want := int(score) * 2
		if got != want {
			t.Errorf("unsatisfied score %d with boost: want %d, got %d", score, want, got)
		}
	}
}

func TestAdjustScore_NoBoostLowerScoresUnchanged(t *testing.T) {
	for _, score := range []model.Score{1, 2, 3, 4} {
		got := adjustScore(score, false, false, 5, false)
		if got != int(score) {
			t.Errorf("unsatisfied score %d no boost: want %d, got %d", score, score, got)
		}
	}
}

func TestAdjustScore_FiveBeatsBoostedFour(t *testing.T) {
	// Even with many opportunities (bonus minimal), unsatisfied score-5 must
	// still beat unsatisfied score-4 with late boost (4×2=8).
	for opps := 1; opps <= 100; opps++ {
		got := adjustScore(5, false, false, opps, false)
		if got <= 8 {
			t.Errorf("unsatisfied 5 with %d opps must beat boosted 4 (8), got %d", opps, got)
		}
	}
}

// --- SolveSlot --------------------------------------------------------------

func TestSolveSlot_BasicAssignment(t *testing.T) {
	// Two players each prefer a different event — both should get their pick.
	sl := slot("s1", event("A", 1), event("B", 1))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"B": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if !slices.Contains(assigned(result, "A"), "alice") {
		t.Error("alice should be assigned to A")
	}
	if !slices.Contains(assigned(result, "B"), "bob") {
		t.Error("bob should be assigned to B")
	}
	if len(result.Unassigned) != 0 {
		t.Errorf("expected no unassigned, got %v", result.Unassigned)
	}
}

func TestSolveSlot_CapacityRespected(t *testing.T) {
	// Three players all want the same event with capacity 2.
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("p1", prefs("s1", map[string]model.Score{"A": 5})),
		player("p2", prefs("s1", map[string]model.Score{"A": 5})),
		player("p3", prefs("s1", map[string]model.Score{"A": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if len(assigned(result, "A")) != 2 {
		t.Errorf("event A capacity 2: want 2 assigned, got %d", len(assigned(result, "A")))
	}
	if len(result.Unassigned) != 1 {
		t.Errorf("want 1 unassigned, got %d", len(result.Unassigned))
	}
}

func TestSolveSlot_NoInterestSkipped(t *testing.T) {
	// One player has interest, one does not.
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 4})),
		player("bob", map[string]map[string]model.Score{}), // no interest in s1
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if !slices.Contains(assigned(result, "A"), "alice") {
		t.Error("alice should be assigned")
	}
	if slices.Contains(assigned(result, "A"), "bob") {
		t.Error("bob has no interest and should not be assigned")
	}
}

func TestSolveSlot_HighScoreWinsContention(t *testing.T) {
	// Capacity 1: player with score 5 should beat player with score 3.
	sl := slot("s1", event("A", 1))
	players := []model.Player{
		player("low", prefs("s1", map[string]model.Score{"A": 3})),
		player("high", prefs("s1", map[string]model.Score{"A": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if !slices.Contains(assigned(result, "A"), "high") {
		t.Error("high scorer should win the seat")
	}
	if !slices.Contains(result.Unassigned, "low") {
		t.Error("low scorer should be unassigned")
	}
}

func TestSolveSlot_SatisfactionTracked(t *testing.T) {
	// Player gets a score-5 event → should appear in NewlySatisfied.
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"A": 3})),
	}

	st := NewState(2026, weekendOf(sl))
	result := st.SolveSlot(sl, players, false)

	if !slices.Contains(result.NewlySatisfied, "alice") {
		t.Error("alice should be newly satisfied")
	}
	if slices.Contains(result.NewlySatisfied, "bob") {
		t.Error("bob scored 3, not 5 — should not be satisfied")
	}
	if !st.IsSatisfied("alice") {
		t.Error("state should mark alice as satisfied")
	}
}

func TestSolveSlot_SatisfiedPlayerDeprioritised(t *testing.T) {
	// Slot 1: alice gets her score-5 event → satisfied.
	// Slot 2: same event, capacity 1. alice (satisfied, score 5)
	// vs charlie (unsatisfied, score 5). Charlie should win.
	sl1 := slot("s1", event("A", 1))
	sl2 := slot("s2", event("A", 1))

	alice := model.Player{
		ID:   "alice",
		Name: "alice",
		Prefs: map[string]map[string]model.Score{
			"s1": {"A": 5},
			"s2": {"A": 5},
		},
	}
	charlie := model.Player{
		ID:   "charlie",
		Name: "charlie",
		Prefs: map[string]map[string]model.Score{
			"s2": {"A": 5},
		},
	}

	st := NewState(2026, weekendOf(sl1, sl2))
	st.SolveSlot(sl1, []model.Player{alice}, false)

	if !st.IsSatisfied("alice") {
		t.Fatal("alice should be satisfied after slot 1")
	}

	result := st.SolveSlot(sl2, []model.Player{alice, charlie}, false)

	if !slices.Contains(assigned(result, "A"), "charlie") {
		t.Error("unsatisfied charlie should win the seat over satisfied alice")
	}
	if slices.Contains(assigned(result, "A"), "alice") {
		t.Error("satisfied alice should lose the seat to unsatisfied charlie")
	}
}

func TestSolveSlot_LateBoostPrioritisesUnsatisfied(t *testing.T) {
	// Without boost: satisfied alice (score 5) beats unsatisfied bob (score 4).
	// With boost:    unsatisfied bob (score 4→8) beats satisfied alice (score 5).
	sl := slot("s1", event("A", 1))
	alice := model.Player{
		ID:   "alice",
		Name: "alice",
		Prefs: map[string]map[string]model.Score{
			"s1": {"A": 5},
		},
	}
	bob := model.Player{
		ID:   "bob",
		Name: "bob",
		Prefs: map[string]map[string]model.Score{
			"s1": {"A": 4},
		},
	}

	// Without boost: alice wins.
	stNoBoost := NewState(2026, weekendOf(sl))
	stNoBoost.satisfied["alice"] = true // pre-satisfy alice
	resultNoBoost := stNoBoost.SolveSlot(sl, []model.Player{alice, bob}, false)
	if !slices.Contains(assigned(resultNoBoost, "A"), "alice") {
		t.Error("without boost: satisfied alice (score 5) should beat unsatisfied bob (score 4)")
	}

	// With boost: bob wins (4×2=8 > 5).
	stBoost := NewState(2026, weekendOf(sl))
	stBoost.satisfied["alice"] = true // pre-satisfy alice
	resultBoost := stBoost.SolveSlot(sl, []model.Player{alice, bob}, true)
	if !slices.Contains(assigned(resultBoost, "A"), "bob") {
		t.Error("with boost: unsatisfied bob (score 4→8) should beat satisfied alice (score 5)")
	}
}

func TestSolveSlot_TotalScoreIsUnadjusted(t *testing.T) {
	// TotalScore must reflect actual preference scores, not adjusted ones.
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"A": 3})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if result.TotalScore != 8 {
		t.Errorf("TotalScore: want 8 (5+3), got %d", result.TotalScore)
	}
}

func TestSolveSlot_EmptySlot(t *testing.T) {
	sl := slot("s1", event("A", 4))
	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, []model.Player{}, false)

	if len(result.Assignments) != 0 {
		t.Error("no players means no assignments")
	}
}

func TestSolveSlot_UndersubscribedEventFlagged(t *testing.T) {
	// Only 1 player interested in event A — fewer than the hardcoded threshold
	// of 3. The lone player is still assigned, and the event is flagged.
	sl := slot("s1", event("A", 4), event("B", 4))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5, "B": 3})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if !slices.Contains(result.UndersubscribedEvents, "A") {
		t.Error("event A should be flagged (1 player < 3)")
	}
	if !slices.Contains(assigned(result, "A"), "alice") {
		t.Error("alice should still be assigned to her preferred event A")
	}
}

func TestSolveSlot_FullyEmptyEventFlagged(t *testing.T) {
	// Event B has no interested players — flagged.
	sl := slot("s1", event("A", 4), event("B", 4))
	players := []model.Player{
		player("p1", prefs("s1", map[string]model.Score{"A": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if !slices.Contains(result.UndersubscribedEvents, "B") {
		t.Error("event B should be flagged (0 players < 3)")
	}
}

func TestSolveSlot_ThreePlayersIsViable(t *testing.T) {
	// Exactly 3 players → at the threshold, not flagged.
	sl := slot("s1", event("A", 4))
	players := []model.Player{
		player("p1", prefs("s1", map[string]model.Score{"A": 5})),
		player("p2", prefs("s1", map[string]model.Score{"A": 5})),
		player("p3", prefs("s1", map[string]model.Score{"A": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if slices.Contains(result.UndersubscribedEvents, "A") {
		t.Errorf("3 players meets the threshold, should not be flagged")
	}
}

func TestSolveSlot_TopScoreAlwaysBeatsFallback(t *testing.T) {
	// 10 players all have X=5 (adjusted to 10). X has capacity for all 10.
	// 4 of them also have Y=4. Nobody should be routed to Y while X has room —
	// the MCMF must never trade a score-10 edge for a score-4 edge.
	sl := model.Slot{
		ID:   "s1",
		Name: "s1",
		Events: []model.Event{
			event("X", 10),
			event("Y", 4),
		},
	}
	players := make([]model.Player, 10)
	for i := range players {
		p := model.Player{
			ID:   fmt.Sprintf("p%d", i),
			Name: fmt.Sprintf("p%d", i),
			Prefs: map[string]map[string]model.Score{
				"s1": {"X": 5},
			},
		}
		if i < 4 {
			p.Prefs["s1"]["Y"] = 4
		}
		players[i] = p
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if len(assigned(result, "X")) != 10 {
		t.Errorf("all 10 players should be on X (capacity allows it), got %d", len(assigned(result, "X")))
	}
	if len(assigned(result, "Y")) != 0 {
		t.Errorf("no player should be routed to Y while X has capacity, got %v", assigned(result, "Y"))
	}
}

func TestSolveSlot_TopScoreLoserGetsFallback(t *testing.T) {
	// X has capacity 4, 6 players want it at score 5. The 2 who lose the
	// lottery and also have Y=4 should be reassigned to Y, not left unassigned.
	sl := model.Slot{
		ID:   "s1",
		Name: "s1",
		Events: []model.Event{
			event("X", 4),
			event("Y", 4),
		},
	}
	players := []model.Player{
		player("p1", prefs("s1", map[string]model.Score{"X": 5, "Y": 4})),
		player("p2", prefs("s1", map[string]model.Score{"X": 5, "Y": 4})),
		player("p3", prefs("s1", map[string]model.Score{"X": 5})),
		player("p4", prefs("s1", map[string]model.Score{"X": 5})),
		player("p5", prefs("s1", map[string]model.Score{"X": 5})),
		player("p6", prefs("s1", map[string]model.Score{"X": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	if len(assigned(result, "X")) != 4 {
		t.Errorf("X capacity 4: want 4 assigned, got %d", len(assigned(result, "X")))
	}
	// p1 and p2 have a fallback; whichever of them lost the X lottery should
	// end up on Y, not unassigned.
	onY := assigned(result, "Y")
	onX := assigned(result, "X")
	for _, pid := range []string{"p1", "p2"} {
		if !slices.Contains(onX, pid) && !slices.Contains(onY, pid) {
			t.Errorf("%s lost the X lottery but was not reassigned to Y", pid)
		}
	}
}

func TestSolveSlot_ScarcePlayerBeatsAbundantOne(t *testing.T) {
	// Two slots. Event A in s1 has capacity 1.
	// "scarce" only has a score-5 in s1.
	// "abundant" has score-5 in both s1 and s2.
	// scarce must win the s1 seat — they have no other chance.
	sl1 := slot("s1", event("A", 1))
	sl2 := slot("s2", event("B", 1))

	scarce := model.Player{
		ID:   "scarce",
		Name: "scarce",
		Prefs: map[string]map[string]model.Score{
			"s1": {"A": 5},
		},
	}
	abundant := model.Player{
		ID:   "abundant",
		Name: "abundant",
		Prefs: map[string]map[string]model.Score{
			"s1": {"A": 5},
			"s2": {"B": 5},
		},
	}

	st := NewState(2026, weekendOf(sl1, sl2))
	r1 := st.SolveSlot(sl1, []model.Player{scarce, abundant}, false)

	if !slices.Contains(assigned(r1, "A"), "scarce") {
		t.Errorf("scarce player (1 opportunity) should win the s1 seat over abundant (2), got %v", assigned(r1, "A"))
	}
}

func TestAdjustScore_DMBonusOnAllEdges(t *testing.T) {
	// DMs get +10 on every edge type.
	cases := []struct {
		name string
		got  int
		want int
	}{
		{"satisfied score-5", adjustScore(5, true, false, 5, true), 15},
		{"satisfied score-1", adjustScore(1, true, false, 5, true), 11},
		{"unsatisfied score-5 abundant", adjustScore(5, false, false, 5, true), 20},
		{"unsatisfied score-5 scarce", adjustScore(5, false, false, 1, true), 24},
		{"unsatisfied score-4 boost", adjustScore(4, false, true, 5, true), 18},
		{"unsatisfied score-4 no boost", adjustScore(4, false, false, 5, true), 14},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s: want %d, got %d", c.name, c.want, c.got)
		}
	}
}

func TestSolveSlot_DMExcludedFromOwnSlot(t *testing.T) {
	// "alice" DMs event A in slot s1. She has Prefs for s1 but must NOT
	// be assigned because she is busy running A.
	dmEvent := model.Event{ID: "A", Name: "A", Capacity: 4, DMID: "alice"}
	sl := slot("s1", dmEvent, event("B", 4))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5, "B": 5})),
		player("bob", prefs("s1", map[string]model.Score{"A": 5, "B": 4})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players, false)

	for _, ev := range []string{"A", "B"} {
		if slices.Contains(assigned(result, ev), "alice") {
			t.Errorf("alice DMs A in this slot and must not be assigned to %s", ev)
		}
	}
	if !slices.Contains(assigned(result, "A"), "bob") {
		t.Error("bob should be assigned to A (alice is DMing it, doesn't count as participant)")
	}
}

func TestSolveSlot_DMPriorityBeatsRegularPlayer(t *testing.T) {
	// Two players compete for one seat. dm is a DM elsewhere in the weekend,
	// regular is not. Both score 5. DM should win because +10 bonus.
	dmEventInOtherSlot := model.Event{ID: "Z", Name: "Z", Capacity: 4, DMID: "dm"}
	sl1 := slot("s1", event("A", 1))
	sl2 := slot("s2", dmEventInOtherSlot)

	dm := model.Player{
		ID:   "dm",
		Name: "dm",
		Prefs: map[string]map[string]model.Score{
			"s1": {"A": 5},
		},
	}
	regular := model.Player{
		ID:   "regular",
		Name: "regular",
		Prefs: map[string]map[string]model.Score{
			"s1": {"A": 5},
		},
	}

	st := NewState(2026, weekendOf(sl1, sl2))
	result := st.SolveSlot(sl1, []model.Player{dm, regular}, false)

	if !slices.Contains(assigned(result, "A"), "dm") {
		t.Errorf("DM should beat regular player for the score-5 seat, got %v", assigned(result, "A"))
	}
}

func TestSolveSlot_DMOpportunitiesExcludesDMingSlots(t *testing.T) {
	// dm DMs s2 and s3, leaving only s1 and s4 to play.
	// They have score-5 events in s1, s2, s3, s4 — but s2 and s3 don't count
	// as opportunities. So playable opportunities = 2 (s1, s4).
	// That gives them a bigger scarcity bonus than if all 4 slots counted.
	dmEvent2 := model.Event{ID: "Y", Name: "Y", Capacity: 4, DMID: "dm"}
	dmEvent3 := model.Event{ID: "Z", Name: "Z", Capacity: 4, DMID: "dm"}
	sl1 := slot("s1", event("A", 1))
	sl2 := slot("s2", dmEvent2)
	sl3 := slot("s3", dmEvent3)
	sl4 := slot("s4", event("D", 1))

	dm := model.Player{
		ID:   "dm",
		Name: "dm",
		Prefs: map[string]map[string]model.Score{
			"s1": {"A": 5},
			"s2": {"Y": 5}, // DMing — shouldn't count
			"s3": {"Z": 5}, // DMing — shouldn't count
			"s4": {"D": 5},
		},
	}

	st := NewState(2026, weekendOf(sl1, sl2, sl3, sl4))
	// Compute opportunities at the start (slotIndex 0).
	opps := st.remainingOpportunities(0, []model.Player{dm})

	// Should be 2: s1 and s4. NOT 4.
	if opps["dm"] != 2 {
		t.Errorf("DM playable opportunities: want 2 (s1, s4 — s2/s3 are DMing), got %d", opps["dm"])
	}
}
