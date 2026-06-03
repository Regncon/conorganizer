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

// minCostFlow runs successive shortest-path augmentation until no
// augmenting path remains. It returns (totalFlow, totalCost).
func (g *flowGraph) minCostFlow(s, t int) (int, int) {
	totalFlow, totalCost := 0, 0
	for {
		f, c, ok := g.spfaAugment(s, t)
		if !ok {
			break
		}
		totalFlow += f
		totalCost += c
	}
	return totalFlow, totalCost
}

// spfaAugment finds the minimum-cost augmenting path from s to t using SPFA,
// pushes as much flow as the path bottleneck allows, and returns
// (flow, cost, true). Returns (0, 0, false) when no path exists.
func (g *flowGraph) spfaAugment(s, t int) (int, int, bool) {
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

	if dist[t] == math.MaxInt {
		return 0, 0, false
	}

	// Bottleneck capacity along the shortest path.
	pushFlow := math.MaxInt
	for v := t; v != s; v = prevv[v] {
		eid := preve[v]
		if rem := g.edges[eid].cap - g.edges[eid].flow; rem < pushFlow {
			pushFlow = rem
		}
	}

	// Augment.
	for v := t; v != s; v = prevv[v] {
		eid := preve[v]
		g.edges[eid].flow += pushFlow
		g.edges[eid^1].flow -= pushFlow
	}

	return pushFlow, pushFlow * dist[t], true
}
