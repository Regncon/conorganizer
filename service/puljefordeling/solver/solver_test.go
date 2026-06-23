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

func weekendOf(slots ...model.Slot) model.Weekend {
	return model.Weekend{Slots: slots}
}

// --- adjustScore: bands -----------------------------------------------------

func TestAdjustScore_Bands(t *testing.T) {
	cases := []struct {
		name      string
		score     model.Score
		satisfied bool
		want      int
	}{
		{"unsatisfied Veldig", 5, false, bandUnsatVeldig},
		{"satisfied Veldig", 5, true, bandSatVeldig},
		{"Middels unsatisfied", 3, false, bandMiddels},
		{"Middels satisfied", 3, true, bandMiddels}, // flattened: same band
		{"Litt unsatisfied", 1, false, bandLitt},
		{"Litt satisfied", 1, true, bandLitt}, // flattened: same band
	}
	for _, c := range cases {
		got := adjustScore(c.score, c.satisfied, false, false, 0)
		if got != c.want {
			t.Errorf("%s: want %d, got %d", c.name, c.want, got)
		}
	}
}

func TestAdjustScore_MiddelsAndLittIgnoreSatisfaction(t *testing.T) {
	// The unsatisfied advantage exists only on the top choice.
	if adjustScore(3, false, false, false, 0) != adjustScore(3, true, false, false, 0) {
		t.Error("Middels weight must not depend on satisfaction")
	}
	if adjustScore(1, false, false, false, 0) != adjustScore(1, true, false, false, 0) {
		t.Error("Litt weight must not depend on satisfaction")
	}
}

func TestAdjustScore_UnmetVeldigBeatsAnyDMLowerInterest(t *testing.T) {
	regularUnmetVeldig := adjustScore(5, false, false, false, 0) // 800, no bumps
	// The strongest a DM's non-top edge can ever be:
	dmMiddelsMax := adjustScore(3, false, true, true, 0) // never-seated + DM
	dmLittMax := adjustScore(1, false, true, true, 0)
	if regularUnmetVeldig <= dmMiddelsMax {
		t.Errorf("unmet Veldig (%d) must beat any DM Middels (%d)", regularUnmetVeldig, dmMiddelsMax)
	}
	if regularUnmetVeldig <= dmLittMax {
		t.Errorf("unmet Veldig (%d) must beat any DM Litt (%d)", regularUnmetVeldig, dmLittMax)
	}
}

func TestAdjustScore_DMBump(t *testing.T) {
	if got := adjustScore(5, false, false, true, 0); got != bandUnsatVeldig+dmBump {
		t.Errorf("DM unsat Veldig: want %d, got %d", bandUnsatVeldig+dmBump, got)
	}
	if got := adjustScore(3, false, false, true, 0); got != bandMiddels+dmBump {
		t.Errorf("DM Middels: want %d, got %d", bandMiddels+dmBump, got)
	}
}

func TestAdjustScore_NeverSeatedBump(t *testing.T) {
	if got := adjustScore(3, false, true, false, 0); got != bandMiddels+neverSeatedBump {
		t.Errorf("never-seated Middels: want %d, got %d", bandMiddels+neverSeatedBump, got)
	}
	// Satisfied players never get the never-seated bump.
	if got := adjustScore(5, true, true, false, 0); got != bandSatVeldig {
		t.Errorf("satisfied gets no never-seated bump: want %d, got %d", bandSatVeldig, got)
	}
}

func TestAdjustScore_MissBonusGrowsAndCaps(t *testing.T) {
	cases := []struct {
		misses int
		want   int
	}{
		{0, bandUnsatVeldig},
		{1, bandUnsatVeldig + 20},
		{2, bandUnsatVeldig + 40},
		{3, bandUnsatVeldig + 60}, // capped
		{4, bandUnsatVeldig + 60},
		{10, bandUnsatVeldig + 60},
	}
	for _, c := range cases {
		if got := adjustScore(5, false, false, false, c.misses); got != c.want {
			t.Errorf("unsat Veldig %d misses: want %d, got %d", c.misses, c.want, got)
		}
	}
	// Miss bonus only applies to unsatisfied top-choice edges.
	if got := adjustScore(3, false, false, false, 10); got != bandMiddels {
		t.Errorf("misses must not boost Middels: want %d, got %d", bandMiddels, got)
	}
	if got := adjustScore(5, true, false, false, 10); got != bandSatVeldig {
		t.Errorf("misses must not boost a satisfied player: want %d, got %d", bandSatVeldig, got)
	}
}

func TestAdjustScore_MaxBumpStaysInBand(t *testing.T) {
	// The largest a Middels edge can get must stay below satisfied-Veldig, which
	// must stay below unmet-Veldig — bumps never cross a band.
	maxMiddels := adjustScore(3, false, true, true, 0) // never-seated + DM
	if !(maxMiddels < bandSatVeldig && bandSatVeldig < bandUnsatVeldig) {
		t.Errorf("band separation broken: maxMiddels=%d satVeldig=%d unsatVeldig=%d",
			maxMiddels, bandSatVeldig, bandUnsatVeldig)
	}
}

// --- SolveSlot --------------------------------------------------------------

func TestSolveSlot_BasicAssignment(t *testing.T) {
	sl := slot("s1", event("A", 1), event("B", 1))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"B": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

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
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("p1", prefs("s1", map[string]model.Score{"A": 5})),
		player("p2", prefs("s1", map[string]model.Score{"A": 5})),
		player("p3", prefs("s1", map[string]model.Score{"A": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if len(assigned(result, "A")) != 2 {
		t.Errorf("event A capacity 2: want 2 assigned, got %d", len(assigned(result, "A")))
	}
	if len(result.Unassigned) != 1 {
		t.Errorf("want 1 unassigned, got %d", len(result.Unassigned))
	}
}

func TestSolveSlot_NoInterestSkipped(t *testing.T) {
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 3})),
		player("bob", map[string]map[string]model.Score{}),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if !slices.Contains(assigned(result, "A"), "alice") {
		t.Error("alice should be assigned")
	}
	if slices.Contains(assigned(result, "A"), "bob") {
		t.Error("bob has no interest and should not be assigned")
	}
}

func TestSolveSlot_HighScoreWinsContention(t *testing.T) {
	sl := slot("s1", event("A", 1))
	players := []model.Player{
		player("low", prefs("s1", map[string]model.Score{"A": 3})),
		player("high", prefs("s1", map[string]model.Score{"A": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if !slices.Contains(assigned(result, "A"), "high") {
		t.Error("high scorer should win the seat")
	}
	if !slices.Contains(result.Unassigned, "low") {
		t.Error("low scorer should be unassigned")
	}
}

func TestSolveSlot_SatisfactionTracked(t *testing.T) {
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"A": 3})),
	}

	st := NewState(2026, weekendOf(sl))
	result := st.SolveSlot(sl, players)

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
	// Slot 1: alice gets her top choice → satisfied. Slot 2: same event, cap 1.
	// alice (satisfied) vs charlie (unsatisfied). Charlie should win.
	sl1 := slot("s1", event("A", 1))
	sl2 := slot("s2", event("A", 1))

	alice := model.Player{ID: "alice", Name: "alice", Prefs: map[string]map[string]model.Score{
		"s1": {"A": 5}, "s2": {"A": 5},
	}}
	charlie := model.Player{ID: "charlie", Name: "charlie", Prefs: map[string]map[string]model.Score{
		"s2": {"A": 5},
	}}

	st := NewState(2026, weekendOf(sl1, sl2))
	st.SolveSlot(sl1, []model.Player{alice})
	if !st.IsSatisfied("alice") {
		t.Fatal("alice should be satisfied after slot 1")
	}

	result := st.SolveSlot(sl2, []model.Player{alice, charlie})
	if !slices.Contains(assigned(result, "A"), "charlie") {
		t.Error("unsatisfied charlie should win the seat over satisfied alice")
	}
	if slices.Contains(assigned(result, "A"), "alice") {
		t.Error("satisfied alice should lose the seat to unsatisfied charlie")
	}
}

func TestSolveSlot_MissBoostsNextSlot(t *testing.T) {
	// p misses a top choice in s1 (loses to a DM), so in s2 their miss bonus
	// must let them beat a fresh, never-missed player for a top-choice seat.
	sl1 := slot("s1", event("A", 1))
	sl2 := slot("s2", event("B", 1))
	sl3 := slot("s3", model.Event{ID: "Z", Name: "Z", Capacity: 4, DMID: "dm"})

	dm := model.Player{ID: "dm", Name: "dm", Prefs: map[string]map[string]model.Score{
		"s1": {"A": 5}, // DM bump → wins A in s1
	}}
	p := model.Player{ID: "p", Name: "p", Prefs: map[string]map[string]model.Score{
		"s1": {"A": 5}, // loses A to the DM → records a miss
		"s2": {"B": 5},
	}}
	fresh := model.Player{ID: "fresh", Name: "fresh", Prefs: map[string]map[string]model.Score{
		"s2": {"B": 5}, // no miss history
	}}

	st := NewState(2026, weekendOf(sl1, sl2, sl3))

	r1 := st.SolveSlot(sl1, []model.Player{dm, p})
	if !slices.Contains(assigned(r1, "A"), "dm") {
		t.Fatalf("DM should win A in s1, got %v", assigned(r1, "A"))
	}

	r2 := st.SolveSlot(sl2, []model.Player{p, fresh})
	if !slices.Contains(assigned(r2, "B"), "p") {
		t.Errorf("p missed in s1 and should win B in s2 over a fresh player, got %v", assigned(r2, "B"))
	}
}

func TestSolveSlot_TotalScoreIsUnadjusted(t *testing.T) {
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"A": 3})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if result.TotalScore != 8 {
		t.Errorf("TotalScore: want 8 (5+3), got %d", result.TotalScore)
	}
}

func TestSolveSlot_EmptySlot(t *testing.T) {
	sl := slot("s1", event("A", 4))
	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, []model.Player{})

	if len(result.Assignments) != 0 {
		t.Error("no players means no assignments")
	}
}

func TestSolveSlot_UndersubscribedEventFlagged(t *testing.T) {
	sl := slot("s1", event("A", 4), event("B", 4))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5, "B": 3})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if !slices.Contains(result.UndersubscribedEvents, "A") {
		t.Error("event A should be flagged (1 player < 3)")
	}
	if !slices.Contains(assigned(result, "A"), "alice") {
		t.Error("alice should be assigned to her preferred event A")
	}
}

func TestSolveSlot_FullyEmptyEventFlagged(t *testing.T) {
	sl := slot("s1", event("A", 4), event("B", 4))
	players := []model.Player{
		player("p1", prefs("s1", map[string]model.Score{"A": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if !slices.Contains(result.UndersubscribedEvents, "B") {
		t.Error("event B should be flagged (0 players < 3)")
	}
}

func TestSolveSlot_ThreePlayersIsViable(t *testing.T) {
	sl := slot("s1", event("A", 4))
	players := []model.Player{
		player("p1", prefs("s1", map[string]model.Score{"A": 5})),
		player("p2", prefs("s1", map[string]model.Score{"A": 5})),
		player("p3", prefs("s1", map[string]model.Score{"A": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if slices.Contains(result.UndersubscribedEvents, "A") {
		t.Error("3 players meets the threshold, should not be flagged")
	}
}

func TestSolveSlot_TopScoreAlwaysBeatsFallback(t *testing.T) {
	// 10 players all want X (cap 10). 4 also have a Middels fallback Y. Nobody
	// should be routed to Y while X has room for them.
	sl := model.Slot{ID: "s1", Name: "s1", Events: []model.Event{event("X", 10), event("Y", 4)}}
	players := make([]model.Player, 10)
	for i := range players {
		p := model.Player{
			ID: fmt.Sprintf("p%d", i), Name: fmt.Sprintf("p%d", i),
			Prefs: map[string]map[string]model.Score{"s1": {"X": 5}},
		}
		if i < 4 {
			p.Prefs["s1"]["Y"] = 3
		}
		players[i] = p
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if len(assigned(result, "X")) != 10 {
		t.Errorf("all 10 players should be on X, got %d", len(assigned(result, "X")))
	}
	if len(assigned(result, "Y")) != 0 {
		t.Errorf("no player should be routed to Y while X has capacity, got %v", assigned(result, "Y"))
	}
}

func TestSolveSlot_TopScoreLoserGetsFallback(t *testing.T) {
	// X cap 4, 6 want X. The 2 who also have a fallback Y should end up seated
	// (in X or Y), not unassigned.
	sl := model.Slot{ID: "s1", Name: "s1", Events: []model.Event{event("X", 4), event("Y", 4)}}
	players := []model.Player{
		player("p1", prefs("s1", map[string]model.Score{"X": 5, "Y": 3})),
		player("p2", prefs("s1", map[string]model.Score{"X": 5, "Y": 3})),
		player("p3", prefs("s1", map[string]model.Score{"X": 5})),
		player("p4", prefs("s1", map[string]model.Score{"X": 5})),
		player("p5", prefs("s1", map[string]model.Score{"X": 5})),
		player("p6", prefs("s1", map[string]model.Score{"X": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if len(assigned(result, "X")) != 4 {
		t.Errorf("X capacity 4: want 4 assigned, got %d", len(assigned(result, "X")))
	}
	onX, onY := assigned(result, "X"), assigned(result, "Y")
	for _, pid := range []string{"p1", "p2"} {
		if !slices.Contains(onX, pid) && !slices.Contains(onY, pid) {
			t.Errorf("%s has a fallback but was not seated", pid)
		}
	}
}

func TestSolveSlot_DMExcludedFromOwnSlot(t *testing.T) {
	dmEvent := model.Event{ID: "A", Name: "A", Capacity: 4, DMID: "alice"}
	sl := slot("s1", dmEvent, event("B", 4))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5, "B": 5})),
		player("bob", prefs("s1", map[string]model.Score{"A": 5, "B": 3})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	for _, ev := range []string{"A", "B"} {
		if slices.Contains(assigned(result, ev), "alice") {
			t.Errorf("alice DMs A this slot and must not be assigned to %s", ev)
		}
	}
	if !slices.Contains(assigned(result, "A"), "bob") {
		t.Error("bob should be assigned to A")
	}
}

func TestSolveSlot_ReverseEdgeBumpMarksMoved(t *testing.T) {
	// "a" runs a game in another slot (DM bump), so it outranks "b" for the
	// single X seat and is tentatively seated there first. "b" wants only X, so
	// the solver bumps "a" off X — via a reverse/residual edge — down to its Y
	// fallback to seat "b". "a" is therefore "moved"; "b" got its top wish
	// directly and is not.
	sl1 := slot("s1", event("X", 1), event("Y", 1))
	sl2 := slot("s2", model.Event{ID: "Z", Name: "Z", Capacity: 4, DMID: "a"})

	a := model.Player{ID: "a", Name: "a", Prefs: map[string]map[string]model.Score{"s1": {"X": 5, "Y": 3}}}
	b := model.Player{ID: "b", Name: "b", Prefs: map[string]map[string]model.Score{"s1": {"X": 5}}}

	st := NewState(2026, weekendOf(sl1, sl2))
	result := st.SolveSlot(sl1, []model.Player{a, b})

	if !slices.Contains(assigned(result, "X"), "b") {
		t.Errorf("b should win the X seat, got %v", assigned(result, "X"))
	}
	if !slices.Contains(assigned(result, "Y"), "a") {
		t.Errorf("a should be bumped down to its Y fallback, got %v", assigned(result, "Y"))
	}
	if !slices.Contains(result.MovedPlayers, "a") {
		t.Errorf("a was bumped off X by a reverse edge and must be marked moved, got %v", result.MovedPlayers)
	}
	if slices.Contains(result.MovedPlayers, "b") {
		t.Errorf("b got its top wish directly and must not be marked moved, got %v", result.MovedPlayers)
	}
}

func TestSolveSlot_NoBumpLeavesMovedEmpty(t *testing.T) {
	// Two players, two disjoint top choices: each is seated directly, no reverse
	// edge is ever used, so nobody is marked moved.
	sl := slot("s1", event("A", 1), event("B", 1))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"B": 5})),
	}

	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)

	if len(result.MovedPlayers) != 0 {
		t.Errorf("no contention means no bumps, want empty MovedPlayers, got %v", result.MovedPlayers)
	}
}

func TestSolveSlot_DMPriorityBeatsRegularPlayer(t *testing.T) {
	// dm runs a game elsewhere; both want the single A seat at the top level.
	// The DM bump should win it.
	dmElsewhere := model.Event{ID: "Z", Name: "Z", Capacity: 4, DMID: "dm"}
	sl1 := slot("s1", event("A", 1))
	sl2 := slot("s2", dmElsewhere)

	dm := model.Player{ID: "dm", Name: "dm", Prefs: map[string]map[string]model.Score{"s1": {"A": 5}}}
	regular := model.Player{ID: "regular", Name: "regular", Prefs: map[string]map[string]model.Score{"s1": {"A": 5}}}

	st := NewState(2026, weekendOf(sl1, sl2))
	result := st.SolveSlot(sl1, []model.Player{dm, regular})

	if !slices.Contains(assigned(result, "A"), "dm") {
		t.Errorf("DM should beat regular player for the seat, got %v", assigned(result, "A"))
	}
}

func TestSolveSlotFixed_PinnedSeatHonoredWithoutPreference(t *testing.T) {
	// "kid" has NO preference for A, but is manually pinned there. The pin must be
	// honored, and the seat must reduce A's effective capacity.
	sl := slot("s1", event("A", 2))
	players := []model.Player{
		player("alice", prefs("s1", map[string]model.Score{"A": 5})),
		player("bob", prefs("s1", map[string]model.Score{"A": 5})),
		// "kid" intentionally absent from players: a manual placement with no interest.
	}

	result := NewState(2026, weekendOf(sl)).SolveSlotFixed(sl, players, map[string]string{"kid": "A"})

	if !slices.Contains(assigned(result, "A"), "kid") {
		t.Errorf("pinned kid must be seated in A, got %v", assigned(result, "A"))
	}
	if len(assigned(result, "A")) != 2 {
		t.Errorf("A cap 2: pin takes one seat, solver fills one more, want 2 total, got %d", len(assigned(result, "A")))
	}
	// Exactly one of alice/bob got the remaining seat; the other is unassigned.
	if len(result.Unassigned) != 1 {
		t.Errorf("want 1 unassigned (only 1 free seat after the pin), got %v", result.Unassigned)
	}
}

func TestSolveSlotFixed_PinnedSeatCountsAsSeated(t *testing.T) {
	// A pinned player who got their top choice is satisfied; one with no/low
	// interest is seated but not satisfied.
	sl := slot("s1", event("A", 4))
	players := []model.Player{
		player("top", prefs("s1", map[string]model.Score{"A": 5})),
	}

	st := NewState(2026, weekendOf(sl))
	result := st.SolveSlotFixed(sl, players, map[string]string{"top": "A", "nopref": "A"})

	if !st.IsSatisfied("top") {
		t.Error("pinned player whose pinned event is their top choice should be satisfied")
	}
	if st.IsSatisfied("nopref") {
		t.Error("pinned player with no top-choice preference must not be satisfied")
	}
	if !slices.Contains(result.NewlySatisfied, "top") {
		t.Errorf("top should be newly satisfied, got %v", result.NewlySatisfied)
	}
}

func TestSolveSlot_NilFixedRegression(t *testing.T) {
	// The 2-arg SolveSlot must behave exactly as before (delegates with nil).
	sl := slot("s1", event("A", 1))
	players := []model.Player{
		player("low", prefs("s1", map[string]model.Score{"A": 3})),
		player("high", prefs("s1", map[string]model.Score{"A": 5})),
	}
	result := NewState(2026, weekendOf(sl)).SolveSlot(sl, players)
	if !slices.Contains(assigned(result, "A"), "high") {
		t.Error("high scorer should still win with nil fixed")
	}
}
