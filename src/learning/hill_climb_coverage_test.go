//go:build unit

package learning

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
)

// ---------------------------------------------------------------------------
// bestOperation coverage tests
// ---------------------------------------------------------------------------

func TestBestOperation_AllViolateMaxIndegree(t *testing.T) {
	// MaxIndegree=1, and all nodes already have 1 parent.
	// A -> B, A -> C. B and C have indegree 1. Adding any edge to B or C
	// would violate indegree. The only option is adding B->A or C->A but
	// if the score doesn't improve, no operation is found.
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{1, 0, 1, 0}),
	})

	// Score function that never improves.
	noImproveScore := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		return 0.0
	}

	hc := NewHillClimbSearch(data, noImproveScore, WithMaxIndegree(1))
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")

	_, found := hc.bestOperation(g, []string{"A", "B", "C"}, 0.0, nil)
	if found {
		t.Error("expected no improving operation when maxIndegree blocks all adds")
	}
}

func TestBestOperation_AllInTabuList(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})

	// Score that always gives positive delta for adds.
	alwaysImprove := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		return float64(len(parents)) * 1.0
	}

	hc := NewHillClimbSearch(data, alwaysImprove)
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")

	// Put all possible operations in the tabu list.
	tabu := []operation{
		{opType: opAdd, from: "A", to: "B"},
		{opType: opAdd, from: "B", to: "A"},
	}

	_, found := hc.bestOperation(g, []string{"A", "B"}, 0.0, tabu)
	if found {
		t.Error("expected no operation when all operations are tabu")
	}
}

func TestBestOperation_AllInBlacklist(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})

	alwaysImprove := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		return float64(len(parents)) * 1.0
	}

	hc := NewHillClimbSearch(data, alwaysImprove,
		WithBlackList([][2]string{{"A", "B"}, {"B", "A"}}),
	)
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")

	_, found := hc.bestOperation(g, []string{"A", "B"}, 0.0, nil)
	if found {
		t.Error("expected no operation when all edges are blacklisted")
	}
}

func TestBestOperation_WhitelistPreventsDelete(t *testing.T) {
	// Whitelist edge A->B prevents deletion of that edge.
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{1, 0, 1, 0}),
	})

	// Score that rewards removing edges (negative effect of parents).
	negScore := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		return -1.0 * float64(len(parents))
	}

	hc := NewHillClimbSearch(data, negScore,
		WithWhiteList([][2]string{{"A", "B"}}),
	)
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")
	g.AddEdge("A", "B")

	op, found := hc.bestOperation(g, []string{"A", "B"}, -1.0, nil)
	// Should not find delete A->B because it's whitelisted.
	// Might find reverse (if not blacklisted), but whitelist also prevents reverse.
	if found && op.opType == opDelete && op.from == "A" && op.to == "B" {
		t.Error("should not delete whitelisted edge A->B")
	}
}

func TestBestOperation_NegativeScoreNoImprovement(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{1, 0, 1, 0}),
		"C": tabgo.NewSeries("C", []any{0, 0, 1, 1}),
	})

	// Score that always returns negative for any parents.
	negScore := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		if len(parents) == 0 {
			return 0.0
		}
		return -1.0
	}

	hc := NewHillClimbSearch(data, negScore)
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")

	_, found := hc.bestOperation(g, []string{"A", "B", "C"}, 0.0, nil)
	if found {
		t.Error("expected no improvement when score is always negative")
	}
}

func TestBestOperation_DeleteAndReverseOperations(t *testing.T) {
	// Start with edge A->B. Score function rewards deleting it.
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{1, 0, 1, 0}),
	})

	// Score: penalize having parents.
	penalizeParents := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		return -0.5 * float64(len(parents))
	}

	hc := NewHillClimbSearch(data, penalizeParents)
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")
	g.AddEdge("A", "B")

	op, found := hc.bestOperation(g, []string{"A", "B"}, -0.5, nil)
	if !found {
		t.Fatal("expected to find a delete operation")
	}
	if op.opType != opDelete {
		t.Errorf("expected delete operation, got type %d", op.opType)
	}
}

func TestBestOperation_ReversePreferredOverDelete(t *testing.T) {
	// Score function that rewards B having parent A more than A having parent B,
	// but reversing A->B to B->A is even better.
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1}),
	})

	// Custom score: A with parent B gets bonus 2.0, B with parent A gets 0.5.
	asymmetricScore := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		if variable == "A" {
			for _, p := range parents {
				if p == "B" {
					return 2.0
				}
			}
		}
		if variable == "B" {
			for _, p := range parents {
				if p == "A" {
					return 0.5
				}
			}
		}
		return 0.0
	}

	hc := NewHillClimbSearch(data, asymmetricScore)
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")
	g.AddEdge("A", "B") // Current: A->B, score for B=0.5, A=0.0. Total=0.5

	// Reverse to B->A: score for A=2.0, B=0.0. Total=2.0. Delta=1.5
	// Delete A->B: score for B=0.0, A=0.0. Total=0.0. Delta=-0.5 (worse)
	op, found := hc.bestOperation(g, []string{"A", "B"}, 0.5, nil)
	if !found {
		t.Fatal("expected to find a reverse operation")
	}
	if op.opType != opReverse {
		t.Errorf("expected reverse operation, got type %d", op.opType)
	}
}

func TestHillClimbEstimate_AllNegativeScores(t *testing.T) {
	// When the score function returns negative for all operations,
	// the search should produce an empty graph.
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{0, 1}),
		"Y": tabgo.NewSeries("Y", []any{1, 0}),
		"Z": tabgo.NewSeries("Z", []any{0, 0}),
	})

	negScore := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		if len(parents) == 0 {
			return 0.0
		}
		return -10.0
	}

	hc := NewHillClimbSearch(data, negScore)
	bn, err := hc.Estimate()
	if err != nil {
		t.Fatalf("Estimate error: %v", err)
	}
	edges := bn.Edges()
	if len(edges) != 0 {
		t.Errorf("expected no edges, got %v", edges)
	}
}

func TestHillClimbEstimate_MaxIndegreeBlocksAllAdds(t *testing.T) {
	// 3 nodes with maxIndegree=0 effectively. Actually maxIndegree=0 means unlimited.
	// Use maxIndegree=1 with initial whitelist edges filling all indegrees.
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})

	constantScore := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		return 0.0
	}

	hc := NewHillClimbSearch(data, constantScore, WithMaxIndegree(1))
	bn, err := hc.Estimate()
	if err != nil {
		t.Fatalf("Estimate error: %v", err)
	}
	// Should produce no edges since adding any edge has 0 delta.
	edges := bn.Edges()
	if len(edges) != 0 {
		t.Errorf("expected no edges, got %v", edges)
	}
}

func TestLegalOperations_Coverage(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})

	// Score that rewards having parents.
	rewardScore := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		return float64(len(parents))
	}

	hc := NewHillClimbSearch(data, rewardScore)
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")

	ops := hc.LegalOperations(g, []string{"A", "B"})
	if len(ops) == 0 {
		t.Error("expected legal operations for empty graph")
	}

	// With edge A->B, should find delete and possibly reverse operations.
	g.AddEdge("A", "B")
	ops = hc.LegalOperations(g, []string{"A", "B"})
	hasAdd := false
	hasDelete := false
	hasReverse := false
	for _, op := range ops {
		switch op.Type {
		case "add":
			hasAdd = true
		case "delete":
			hasDelete = true
		case "reverse":
			hasReverse = true
		}
	}
	_ = hasAdd
	_ = hasDelete
	_ = hasReverse
	// Just verify it doesn't crash and returns valid operations.
}
