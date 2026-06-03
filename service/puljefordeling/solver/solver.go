package solver // solver is defined in flow.go

import (
	"math/rand/v2"
	"sort"

	"github.com/Regncon/conorganizer/service/puljefordeling/solver/model"
)

// dmBonus is the flat adjustment added to every edge involving a player who
// runs at least one event during the weekend. It rewards their contribution
// by giving them priority for the slots in which they themselves are players.
const dmBonus = 10

// minViablePlayers is the threshold below which an event is flagged for
// organiser review (not cancelled). Three players is enough for most
// tabletop games to be worth running; below that the organiser should look
// at the assignment and decide.
const minViablePlayers = 3

// State carries satisfaction data and seeding parameters forward across slots.
// A player is satisfied once they have received at least one assignment
// to an event they scored 5.
type State struct {
	satisfied map[string]bool
	year      int
	slotIDs   []string
	slotIndex int

	// dmSlots[playerID] is the set of slot IDs in which playerID is DMing
	// an event. Players in this set are excluded from being assigned during
	// those slots.
	dmSlots map[string]map[string]bool

	// isDM[playerID] is true if playerID DMs at least one event anywhere
	// in the weekend. Such players receive a fixed bonus on every edge.
	isDM map[string]bool
}

// NewState returns a fresh State for the start of a weekend.
//
// year is used to derive per-slot tie-breaking seeds (seed = year×1000 + slotIndex),
// making results deterministic within a year but different across years.
//
// weekend provides the full slot schedule so the solver can: (1) compute the
// number of remaining score-5 opportunities for each unsatisfied player,
// (2) recognise which players DM events in which slots, and (3) apply the DM
// bonus across the weekend.
func NewState(year int, weekend model.Weekend) *State {
	slotIDs := make([]string, len(weekend.Slots))
	dmSlots := make(map[string]map[string]bool)
	isDM := make(map[string]bool)
	for i, sl := range weekend.Slots {
		slotIDs[i] = sl.ID
		for _, ev := range sl.Events {
			if ev.DMID == "" {
				continue
			}
			isDM[ev.DMID] = true
			if dmSlots[ev.DMID] == nil {
				dmSlots[ev.DMID] = make(map[string]bool)
			}
			dmSlots[ev.DMID][sl.ID] = true
		}
	}
	return &State{
		satisfied: make(map[string]bool),
		year:      year,
		slotIDs:   slotIDs,
		dmSlots:   dmSlots,
		isDM:      isDM,
	}
}

// IsSatisfied reports whether playerID has received a score-5 assignment.
func (s *State) IsSatisfied(playerID string) bool {
	return s.satisfied[playerID]
}

// SatisfiedCount returns the number of satisfied players.
func (s *State) SatisfiedCount() int {
	return len(s.satisfied)
}

// IsDM reports whether playerID runs at least one event in the weekend.
func (s *State) IsDM(playerID string) bool {
	return s.isDM[playerID]
}

// SolveSlot assigns players to events for one slot, updates the satisfaction
// state, and returns the result.
//
// Players who are DMing any event in this slot are excluded from the player
// pool — they cannot also be assigned as participants in the same slot.
//
// lateBoost, when true, doubles all scores for unsatisfied non-DM players (in
// addition to the permanent score-5 doubling that always applies).
//
// Events whose MinPlayers threshold is not met are NOT cancelled — they are
// reported in UndersubscribedEvents for the organiser to review manually.
// This deliberately avoids the cascading edge cases that automated
// cancellation could produce; organisers can talk to players, allow swaps,
// merge groups, or accept a smaller table at their discretion.
func (s *State) SolveSlot(slot model.Slot, players []model.Player, lateBoost bool) model.SlotResult {
	currentIndex := s.slotIndex
	seed := int64(s.year)*1000 + int64(currentIndex)
	s.slotIndex++

	result := model.SlotResult{
		SlotID:      slot.ID,
		Assignments: make(map[string][]string),
		Seed:        seed,
	}

	// Players DMing in this slot are unavailable as players.
	dmingHere := make(map[string]bool)
	for _, ev := range slot.Events {
		if ev.DMID != "" {
			dmingHere[ev.DMID] = true
		}
	}

	// Only consider players who expressed interest and are not DMing here.
	interested := make([]model.Player, 0, len(players))
	for _, p := range players {
		if dmingHere[p.ID] {
			continue
		}
		if len(p.Prefs[slot.ID]) > 0 {
			interested = append(interested, p)
		}
	}
	if len(interested) == 0 {
		return result
	}

	// Shuffle players before building the graph so that equal-score ties are
	// broken randomly but reproducibly.
	rng := rand.New(rand.NewPCG(uint64(seed), 0)) //nolint:gosec
	rng.Shuffle(len(interested), func(i, j int) {
		interested[i], interested[j] = interested[j], interested[i]
	})

	// Compute remaining score-5 opportunities (including current slot) for
	// each unsatisfied player, excluding slots where the player is DMing.
	opportunities := s.remainingOpportunities(currentIndex, interested)

	// Single MCMF run — no cancellation cascade.
	assignments := s.runMCMF(slot.ID, slot.Events, interested, lateBoost, opportunities)
	result.Assignments = assignments

	// Flag events with fewer than minViablePlayers for human review.
	for _, ev := range slot.Events {
		if len(assignments[ev.ID]) < minViablePlayers {
			result.UndersubscribedEvents = append(result.UndersubscribedEvents, ev.ID)
		}
	}

	// Update satisfaction state and collect totals.
	assigned := make(map[string]bool, len(interested))
	for evID, playerIDs := range assignments {
		for _, pid := range playerIDs {
			score := s.lookupScore(pid, slot.ID, evID, interested)
			result.TotalScore += int(score)
			if score == model.MaxScore && !s.satisfied[pid] {
				s.satisfied[pid] = true
				result.NewlySatisfied = append(result.NewlySatisfied, pid)
			}
			assigned[pid] = true
		}
	}

	// Collect unassigned players (had interest but no seat was available).
	for _, p := range interested {
		if !assigned[p.ID] {
			result.Unassigned = append(result.Unassigned, p.ID)
		}
	}

	// Sort all output slices for deterministic results.
	sort.Strings(result.NewlySatisfied)
	sort.Strings(result.Unassigned)
	sort.Strings(result.UndersubscribedEvents)
	for evID := range result.Assignments {
		sort.Strings(result.Assignments[evID])
	}

	return result
}

// remainingOpportunities counts, for each unsatisfied player, how many of the
// remaining slots (starting at fromIndex) contain at least one event they
// scored 5 — excluding slots in which the player is DMing. Satisfied players
// get 0.
func (s *State) remainingOpportunities(fromIndex int, players []model.Player) map[string]int {
	out := make(map[string]int, len(players))
	if fromIndex >= len(s.slotIDs) {
		return out
	}
	remaining := s.slotIDs[fromIndex:]
	for _, p := range players {
		if s.satisfied[p.ID] {
			continue
		}
		n := 0
		for _, sid := range remaining {
			if s.dmSlots[p.ID][sid] {
				continue
			}
			for _, score := range p.Prefs[sid] {
				if score == model.MaxScore {
					n++
					break
				}
			}
		}
		out[p.ID] = n
	}
	return out
}

// runMCMF builds and solves the flow network for the given events and players,
// returning the raw assignment map (eventID -> []playerID, unsorted).
func (s *State) runMCMF(
	slotID string,
	events []model.Event,
	players []model.Player,
	lateBoost bool,
	opportunities map[string]int,
) map[string][]string {
	assignments := make(map[string][]string)
	if len(events) == 0 {
		return assignments
	}

	// Node layout:
	//   0           → source
	//   1 .. P      → one node per player
	//   P+1 .. P+E  → one node per event
	//   P+E+1       → sink
	P := len(players)
	E := len(events)
	source := 0
	sink := P + E + 1
	g := newFlowGraph(sink + 1)

	for i := range players {
		g.addEdge(source, i+1, 1, 0)
	}

	for j, ev := range events {
		g.addEdge(P+1+j, sink, ev.Capacity, 0)
	}

	// Iterate events (slice, deterministic) for each player rather than the
	// player's preference map so the edge addition order is identical
	// run-to-run.
	for i, p := range players {
		for j, ev := range events {
			score, ok := p.Prefs[slotID][ev.ID]
			if !ok {
				continue
			}
			adj := adjustScore(score, s.satisfied[p.ID], lateBoost, opportunities[p.ID], s.isDM[p.ID])
			g.addEdge(i+1, P+1+j, 1, -adj)
		}
	}

	g.minCostFlow(source, sink)

	for i, p := range players {
		for _, eid := range g.adj[i+1] {
			e := g.edges[eid]
			if e.flow != 1 || e.to < P+1 || e.to > P+E {
				continue
			}
			evID := events[e.to-P-1].ID
			assignments[evID] = append(assignments[evID], p.ID)
		}
	}

	return assignments
}

// lookupScore returns the raw preference score for a player in a slot/event.
func (s *State) lookupScore(playerID, slotID, eventID string, players []model.Player) model.Score {
	for _, p := range players {
		if p.ID == playerID {
			return p.Prefs[slotID][eventID]
		}
	}
	return 0
}

// adjustScore returns the adjusted edge weight for a (player, event) pair.
//
//   - Satisfied players use their raw score.
//   - Unsatisfied players on a score-5 edge get base 10 plus a scarcity bonus
//     max(0, 5 - opportunities). Fewer remaining score-5 opportunities means
//     higher priority.
//   - Unsatisfied players on score 1–4 edges get those scores doubled when
//     lateBoost is enabled, otherwise raw.
//   - DMs (any player running at least one event in the weekend) receive a
//     flat +dmBonus on every edge, rewarding their contribution.
func adjustScore(score model.Score, satisfied, lateBoost bool, opportunities int, isDM bool) int {
	base := rawAdjusted(score, satisfied, lateBoost, opportunities)
	if isDM {
		base += dmBonus
	}
	return base
}

func rawAdjusted(score model.Score, satisfied, lateBoost bool, opportunities int) int {
	if satisfied {
		return int(score)
	}
	if score == model.MaxScore {
		bonus := 5 - opportunities
		if bonus < 0 {
			bonus = 0
		}
		return int(model.MaxScore)*2 + bonus
	}
	if lateBoost {
		return int(score) * 2
	}
	return int(score)
}
