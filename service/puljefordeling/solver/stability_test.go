package solver

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"testing"

	"github.com/Regncon/conorganizer/service/puljefordeling/solver/model"
	"github.com/Regncon/conorganizer/testutil/bdd"
)

// randomSlot builds a slot where many players rate several games equally, so the
// solver has plenty of equally good seatings to choose between.
func randomSlot(rng *rand.Rand) (model.Slot, []model.Player) {
	s := model.Slot{ID: "s1", Name: "s1"}
	for e := range 5 {
		s.Events = append(s.Events, model.Event{ID: fmt.Sprintf("e%d", e), Name: fmt.Sprintf("e%d", e), Capacity: 2 + rng.IntN(4)})
	}
	var players []model.Player
	for p := range 20 {
		prefs := map[string]model.Score{}
		for _, ev := range s.Events {
			if rng.IntN(2) == 0 {
				prefs[ev.ID] = []model.Score{1, 3, 5}[rng.IntN(3)]
			}
		}
		if len(prefs) == 0 {
			continue
		}
		players = append(players, model.Player{ID: fmt.Sprintf("p%02d", p), Name: fmt.Sprintf("p%02d", p), Prefs: map[string]map[string]model.Score{"s1": prefs}, IsOver18: true})
	}
	return s, players
}

func seatOf(result model.SlotResult) map[string]string {
	seats := map[string]string{}
	for evID, pids := range result.Assignments {
		for _, pid := range pids {
			seats[pid] = evID
		}
	}
	return seats
}

func TestSolveSlotFixed_PinningSomeoneIntoTheirOwnSeatMovesNobodyElse(t *testing.T) {
	bdd.Behavior(t, bdd.BDD{
		Given: "a solved pulje where many players rate several games equally",
		When:  "an admin pins one player into the very seat the solver already gave them",
		Then:  "every other player keeps exactly the seat they had",
	})

	for scenario := range 200 {
		// Given
		rng := rand.New(rand.NewPCG(uint64(scenario), 7))
		s, players := randomSlot(rng)
		before := seatOf(NewState(2026, weekendOf(s)).SolveSlot(s, players))
		seated := make([]string, 0, len(before))
		for pid := range before {
			seated = append(seated, pid)
		}
		slices.Sort(seated)

		for _, pinned := range seated {
			// When
			after := seatOf(NewState(2026, weekendOf(s)).SolveSlotFixed(s, players, map[string][]string{pinned: {before[pinned]}}))

			// Then
			for pid, seat := range before {
				if after[pid] != seat {
					t.Fatalf("scenario %d: pinning %s into %s moved %s from %s to %q", scenario, pinned, before[pinned], pid, seat, after[pid])
				}
			}
		}
	}
}
