package learning

import (
	"fmt"
	"sort"

	"github.com/asymmetric-effort/datascience/lib/graphgo"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// GES implements the Greedy Equivalence Search algorithm for structure learning.
// GES searches the space of Markov equivalence classes (represented as PDAGs)
// using a two-phase approach: a forward phase that greedily adds edges, and a
// backward phase that greedily removes edges.
type GES struct {
	data    *tabgo.DataFrame
	scoreFn ScoreFunc
}

// NewGES creates a new GES instance with the given data and scoring function.
func NewGES(data *tabgo.DataFrame, scoreFn ScoreFunc) *GES {
	return &GES{
		data:    data,
		scoreFn: scoreFn,
	}
}

// Estimate runs the GES algorithm and returns a PDAG representing the learned
// equivalence class.
//
// The algorithm proceeds in two phases:
//  1. Forward phase: start with an empty graph, greedily add edges that improve
//     the score until no improvement is possible.
//  2. Backward phase: greedily remove edges that improve the score until no
//     improvement is possible.
func (g *GES) Estimate() (*graphgo.PDAG, error) {
	columns := g.data.Columns()
	if len(columns) < 2 {
		return nil, fmt.Errorf("learning: GES requires at least 2 variables, got %d", len(columns))
	}

	// We work on a DiGraph internally and convert to PDAG at the end.
	dag := graphgo.NewDiGraph()
	for _, col := range columns {
		dag.AddNode(col)
	}

	// Forward phase: greedily add edges.
	for {
		bestFrom, bestTo := "", ""
		bestDelta := 0.0

		nodes := make([]string, len(columns))
		copy(nodes, columns)
		sort.Strings(nodes)

		for _, u := range nodes {
			for _, v := range nodes {
				if u == v {
					continue
				}
				if dag.HasEdge(u, v) || dag.HasEdge(v, u) {
					continue
				}

				// Try adding u -> v.
				oldParents := sortedParents(dag, v)
				oldScore := g.scoreFn(v, oldParents, g.data)

				dag.AddEdge(u, v)
				if !graphgo.IsDAG(dag) {
					_ = dag.RemoveEdge(u, v)
					continue
				}

				newParents := sortedParents(dag, v)
				newScore := g.scoreFn(v, newParents, g.data)
				_ = dag.RemoveEdge(u, v)

				delta := newScore - oldScore
				if delta > bestDelta {
					bestDelta = delta
					bestFrom = u
					bestTo = v
				}
			}
		}

		if bestDelta <= 0 {
			break
		}
		dag.AddEdge(bestFrom, bestTo)
	}

	// Backward phase: greedily remove edges.
	for {
		bestFrom, bestTo := "", ""
		bestDelta := 0.0

		edges := dag.Edges()
		sort.Slice(edges, func(i, j int) bool {
			if edges[i].Src != edges[j].Src {
				return edges[i].Src < edges[j].Src
			}
			return edges[i].Dst < edges[j].Dst
		})

		for _, e := range edges {
			u, v := e.Src, e.Dst

			oldParents := sortedParents(dag, v)
			oldScore := g.scoreFn(v, oldParents, g.data)

			_ = dag.RemoveEdge(u, v)
			newParents := sortedParents(dag, v)
			newScore := g.scoreFn(v, newParents, g.data)
			dag.AddEdge(u, v)

			delta := newScore - oldScore
			if delta > bestDelta {
				bestDelta = delta
				bestFrom = u
				bestTo = v
			}
		}

		if bestDelta <= 0 {
			break
		}
		_ = dag.RemoveEdge(bestFrom, bestTo)
	}

	// Convert the DAG to a PDAG (equivalence class representation).
	pdag := g.dagToPDAG(dag, columns)
	return pdag, nil
}

// Insert adds a directed edge from u to v in the given DAG, checking that the
// result remains a DAG. Returns the updated score delta, or an error if the
// insertion would create a cycle or the edge already exists.
func (g *GES) Insert(dag *graphgo.DiGraph, u, v string) (float64, error) {
	if dag.HasEdge(u, v) || dag.HasEdge(v, u) {
		return 0, fmt.Errorf("learning: edge between %q and %q already exists", u, v)
	}

	oldParents := sortedParents(dag, v)
	oldScore := g.scoreFn(v, oldParents, g.data)

	dag.AddEdge(u, v)
	if !graphgo.IsDAG(dag) {
		_ = dag.RemoveEdge(u, v)
		return 0, fmt.Errorf("learning: inserting edge %q -> %q would create a cycle", u, v)
	}

	newParents := sortedParents(dag, v)
	newScore := g.scoreFn(v, newParents, g.data)

	return newScore - oldScore, nil
}

// Delete removes a directed edge from u to v in the given DAG. Returns the
// score delta, or an error if the edge does not exist.
func (g *GES) Delete(dag *graphgo.DiGraph, u, v string) (float64, error) {
	if !dag.HasEdge(u, v) {
		return 0, fmt.Errorf("learning: edge %q -> %q does not exist", u, v)
	}

	oldParents := sortedParents(dag, v)
	oldScore := g.scoreFn(v, oldParents, g.data)

	_ = dag.RemoveEdge(u, v)

	newParents := sortedParents(dag, v)
	newScore := g.scoreFn(v, newParents, g.data)

	return newScore - oldScore, nil
}

// Turn reverses the direction of an existing edge from u to v (making it v to u)
// in the given DAG. Returns the score delta, or an error if the reversal is not
// possible (edge doesn't exist or reversal would create a cycle).
func (g *GES) Turn(dag *graphgo.DiGraph, u, v string) (float64, error) {
	if !dag.HasEdge(u, v) {
		return 0, fmt.Errorf("learning: edge %q -> %q does not exist", u, v)
	}

	// Compute old scores for both endpoints.
	oldParentsV := sortedParents(dag, v)
	oldScoreV := g.scoreFn(v, oldParentsV, g.data)
	oldParentsU := sortedParents(dag, u)
	oldScoreU := g.scoreFn(u, oldParentsU, g.data)

	_ = dag.RemoveEdge(u, v)
	dag.AddEdge(v, u)

	if !graphgo.IsDAG(dag) {
		_ = dag.RemoveEdge(v, u)
		dag.AddEdge(u, v)
		return 0, fmt.Errorf("learning: turning edge %q -> %q would create a cycle", u, v)
	}

	newParentsV := sortedParents(dag, v)
	newScoreV := g.scoreFn(v, newParentsV, g.data)
	newParentsU := sortedParents(dag, u)
	newScoreU := g.scoreFn(u, newParentsU, g.data)

	return (newScoreV + newScoreU) - (oldScoreV + oldScoreU), nil
}

// dagToPDAG converts a DAG to its PDAG (CPDAG) representation.
// Compelled edges (those present in every DAG of the equivalence class) remain
// directed; reversible edges become undirected.
func (g *GES) dagToPDAG(dag *graphgo.DiGraph, columns []string) *graphgo.PDAG {
	pdag := graphgo.NewPDAG()
	for _, col := range columns {
		pdag.AddNode(col)
	}

	edges := dag.Edges()
	sort.Slice(edges, func(i, j int) bool {
		if edges[i].Src != edges[j].Src {
			return edges[i].Src < edges[j].Src
		}
		return edges[i].Dst < edges[j].Dst
	})

	// Identify v-structures: for each node, if two parents are not adjacent,
	// both edges into the node are compelled.
	compelled := make(map[[2]string]bool)

	for _, v := range columns {
		parents := sortedParents(dag, v)
		for i := 0; i < len(parents); i++ {
			for j := i + 1; j < len(parents); j++ {
				pi, pj := parents[i], parents[j]
				if !dag.HasEdge(pi, pj) && !dag.HasEdge(pj, pi) {
					// v-structure: pi -> v <- pj
					compelled[[2]string{pi, v}] = true
					compelled[[2]string{pj, v}] = true
				}
			}
		}
	}

	// Add edges: compelled edges are directed, others are undirected.
	for _, e := range edges {
		if compelled[[2]string{e.Src, e.Dst}] {
			pdag.AddDirectedEdge(e.Src, e.Dst)
		} else {
			pdag.AddUndirectedEdge(e.Src, e.Dst)
		}
	}

	// Apply Meek rules to orient additional edges that must be directed.
	graphgo.ApplyMeekRules(pdag)

	return pdag
}
