// Package solver implements the min-cost max-flow based event assignment algorithm.
package solver

import "math"

// edge is a directed edge in the flow network.
// Forward and reverse edges are always added in consecutive pairs so that
// the reverse of edge at index i is at index i^1.
type edge struct {
	to   int
	cap  int
	cost int
	flow int
}

// flowGraph is a min-cost max-flow network solved with SPFA
// (Bellman-Ford based shortest-path augmentation).
// Negative edge costs are supported, which allows negated preference
// scores to be used directly.
type flowGraph struct {
	n     int
	edges []edge
	adj   [][]int
}

func newFlowGraph(n int) *flowGraph {
	return &flowGraph{
		n:   n,
		adj: make([][]int, n),
	}
}

// addEdge adds a directed edge from→to with the given capacity and cost,
// together with its zero-capacity reverse edge.
func (g *flowGraph) addEdge(from, to, capacity, cost int) {
	g.adj[from] = append(g.adj[from], len(g.edges))
	g.edges = append(g.edges, edge{to, capacity, cost, 0})
	g.adj[to] = append(g.adj[to], len(g.edges))
	g.edges = append(g.edges, edge{from, 0, -cost, 0})
}

// minCostFlow runs successive shortest-path augmentation, pushing flow only
// while it improves total value, and returns (totalFlow, totalCost, reduced).
// Because the participation bonus is baked into the assignment-edge costs, this
// stops once the cheapest remaining augmenting path is non-negative — i.e. it
// fills chairs only when a new seat is worth more than it costs, rather than
// always pushing to maximum flow.
//
// reduced lists the forward-edge ids whose flow was pushed back by a residual
// (reverse) edge during augmentation — i.e. flow that was tentatively routed one
// way and later rerouted to improve the global objective. The caller interprets
// these in domain terms (which players were bumped off which events).
func (g *flowGraph) minCostFlow(s, t int) (int, int, []int) {
	totalFlow, totalCost := 0, 0
	var reduced []int
	for {
		f, c, red, ok := g.spfaAugment(s, t)
		if !ok {
			break
		}
		totalFlow += f
		totalCost += c
		reduced = append(reduced, red...)
	}
	return totalFlow, totalCost, reduced
}

// spfaAugment finds the minimum-cost augmenting path from s to t using SPFA,
// pushes as much flow as the path bottleneck allows, and returns
// (flow, cost, reduced, true). reduced holds the forward-edge ids whose flow was
// pushed back along this path's reverse edges. Returns (0, 0, nil, false) when
// no path exists.
func (g *flowGraph) spfaAugment(s, t int) (int, int, []int, bool) {
	dist := make([]int, g.n)
	for i := range dist {
		dist[i] = math.MaxInt
	}
	dist[s] = 0

	prevv := make([]int, g.n)
	preve := make([]int, g.n)
	for i := range prevv {
		prevv[i] = -1
	}

	inq := make([]bool, g.n)
	inq[s] = true
	queue := []int{s}

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		inq[v] = false

		for _, eid := range g.adj[v] {
			e := g.edges[eid]
			if e.cap <= e.flow || dist[v] == math.MaxInt {
				continue
			}
			if next := dist[v] + e.cost; next < dist[e.to] {
				dist[e.to] = next
				prevv[e.to] = v
				preve[e.to] = eid
				if !inq[e.to] {
					queue = append(queue, e.to)
					inq[e.to] = true
				}
			}
		}
	}

	// Stop when no welfare-improving augmentation remains. A path with
	// non-negative total cost would not increase total value (the
	// participation bonus is already folded into the edge costs), so we leave
	// the remaining seats empty rather than make a losing trade. SSP visits
	// augmenting paths in non-decreasing cost order, so once the cheapest is
	// non-negative we are at the optimal flow value.
	if dist[t] == math.MaxInt || dist[t] >= 0 {
		return 0, 0, nil, false
	}

	// Bottleneck capacity along the shortest path.
	pushFlow := math.MaxInt
	for v := t; v != s; v = prevv[v] {
		eid := preve[v]
		if rem := g.edges[eid].cap - g.edges[eid].flow; rem < pushFlow {
			pushFlow = rem
		}
	}

	// Augment. A reverse edge (odd id) on the path pushes flow back along its
	// forward edge (eid^1) — record that forward edge as reduced so the caller
	// can see which tentative assignment was undone.
	var reduced []int
	for v := t; v != s; v = prevv[v] {
		eid := preve[v]
		g.edges[eid].flow += pushFlow
		g.edges[eid^1].flow -= pushFlow
		if eid&1 == 1 {
			reduced = append(reduced, eid^1)
		}
	}

	return pushFlow, pushFlow * dist[t], reduced, true
}
