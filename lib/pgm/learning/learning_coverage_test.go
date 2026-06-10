//go:build unit

package learning

import (
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/graphgo"
	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// ---------------------------------------------------------------------------
// buildVariableContext coverage
// ---------------------------------------------------------------------------

func TestBuildVariableContext(t *testing.T) {
	result := buildVariableContext([]string{"A", "B", "C"})
	if !strings.Contains(result, "A") || !strings.Contains(result, "B") || !strings.Contains(result, "C") {
		t.Errorf("expected all variables in context, got %q", result)
	}
}

// ---------------------------------------------------------------------------
// HillClimbSearch — applyOperation, isTabu coverage
// ---------------------------------------------------------------------------

func TestApplyOperation_Add(t *testing.T) {
	data := makeSmallDiscreteDF()
	scoreFn := func(v string, parents []string, d *tabgo.DataFrame) float64 { return 0 }
	hc := NewHillClimbSearch(data, scoreFn)
	_ = hc

	// Create a graph and apply operations.
	g := makeSmallDigraph()

	// Test isTabu with empty list.
	op := operation{opType: opAdd, from: "A", to: "B", delta: 1.0}
	if hc.isTabu(nil, op) {
		t.Error("expected false for empty tabu list")
	}

	// Test isTabu with matching operation.
	tabu := []operation{op}
	if !hc.isTabu(tabu, op) {
		t.Error("expected true for matching tabu operation")
	}

	// Test apply add.
	hc.applyOperation(g, operation{opType: opAdd, from: "B", to: "C"})
	if !g.HasEdge("B", "C") {
		t.Error("expected edge B->C after add")
	}

	// Test apply delete.
	hc.applyOperation(g, operation{opType: opDelete, from: "B", to: "C"})
	if g.HasEdge("B", "C") {
		t.Error("expected no edge B->C after delete")
	}

	// Test apply reverse.
	g.AddEdge("A", "C")
	hc.applyOperation(g, operation{opType: opReverse, from: "A", to: "C"})
	if g.HasEdge("A", "C") {
		t.Error("expected no edge A->C after reverse")
	}
	if !g.HasEdge("C", "A") {
		t.Error("expected edge C->A after reverse")
	}
}

func TestLegalOperations(t *testing.T) {
	data := makeSmallDiscreteDF()
	scoreFn := func(v string, parents []string, d *tabgo.DataFrame) float64 {
		// Simple score: 0 for no parents, -1 for any parents.
		if len(parents) > 0 {
			return -1.0
		}
		return 0
	}
	hc := NewHillClimbSearch(data, scoreFn)
	g := makeSmallDigraph()
	g.AddEdge("A", "B") // This has a negative delta so delete should have positive delta.

	ops := hc.LegalOperations(g, []string{"A", "B"})
	_ = ops // Just verify it runs.
}

// ---------------------------------------------------------------------------
// ExhaustiveSearch — AllScores, AllDAGs coverage
// ---------------------------------------------------------------------------

func TestExhaustiveSearch_AllScores_Coverage(t *testing.T) {
	data := makeTwoVarDF()
	scoreFn := func(v string, parents []string, d *tabgo.DataFrame) float64 { return 0 }
	es := NewExhaustiveSearch(data, scoreFn)
	scores, err := es.AllScores()
	if err != nil {
		t.Fatal(err)
	}
	if len(scores) == 0 {
		t.Error("expected non-empty scores map")
	}
}

func TestExhaustiveSearch_AllDAGs(t *testing.T) {
	data := makeTwoVarDF()
	scoreFn := func(v string, parents []string, d *tabgo.DataFrame) float64 { return 0 }
	es := NewExhaustiveSearch(data, scoreFn)
	dags, err := es.AllDAGs()
	if err != nil {
		t.Fatal(err)
	}
	if len(dags) == 0 {
		t.Error("expected non-empty DAG list")
	}
}

func TestExhaustiveSearch_AllScores_EmptyData(t *testing.T) {
	data := makeEmptyDF()
	scoreFn := func(v string, parents []string, d *tabgo.DataFrame) float64 { return 0 }
	es := NewExhaustiveSearch(data, scoreFn)
	_, err := es.AllScores()
	if err == nil {
		t.Error("expected error for empty data")
	}
}

func TestExhaustiveSearch_AllDAGs_EmptyData(t *testing.T) {
	data := makeEmptyDF()
	scoreFn := func(v string, parents []string, d *tabgo.DataFrame) float64 { return 0 }
	es := NewExhaustiveSearch(data, scoreFn)
	_, err := es.AllDAGs()
	if err == nil {
		t.Error("expected error for empty data")
	}
}

// ---------------------------------------------------------------------------
// BayesianEstimator — pseudoCount coverage
// ---------------------------------------------------------------------------

func TestBayesianEstimator_BDeuPseudoCount(t *testing.T) {
	bn := makeSimpleBN()
	data := makeSmallDiscreteDF()
	be := NewBayesianEstimator(bn, data, BDeu, 5.0)
	err := be.Estimate()
	if err != nil {
		t.Fatalf("BDeu estimate failed: %v", err)
	}
}

func TestBayesianEstimator_K2PseudoCount(t *testing.T) {
	bn := makeSimpleBN()
	data := makeSmallDiscreteDF()
	be := NewBayesianEstimator(bn, data, K2, 1.0)
	err := be.Estimate()
	if err != nil {
		t.Fatalf("K2 estimate failed: %v", err)
	}
}

// ---------------------------------------------------------------------------
// EM — computeLatentPosterior edge case
// ---------------------------------------------------------------------------

func TestEM_NoLatentVars_Coverage(t *testing.T) {
	bn := makeSimpleBN()
	data := makeSmallDiscreteDF()
	em := NewEM(bn, data, nil, 10, 1e-4)
	err := em.Estimate()
	// With no latent vars, EM should converge immediately.
	if err != nil {
		t.Fatalf("EM with no latent vars failed: %v", err)
	}
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func makeSmallDiscreteDF() *tabgo.DataFrame {
	sm := map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1, 0, 1, 1, 0, 0, 1}),
	}
	return tabgo.NewDataFrame(sm)
}

func makeTwoVarDF() *tabgo.DataFrame {
	sm := map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1, 0, 1, 1, 0, 0, 1}),
	}
	return tabgo.NewDataFrame(sm)
}

func makeEmptyDF() *tabgo.DataFrame {
	return tabgo.NewDataFrame(map[string]*tabgo.Series{})
}

func makeSimpleBN() *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()
	bn.AddNode("A")
	bn.AddNode("B")
	bn.AddEdge("A", "B")
	bn.SetStates("A", []string{"0", "1"})
	bn.SetStates("B", []string{"0", "1"})

	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	_ = bn.AddCPD(cpdA)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.8, 0.3}, {0.2, 0.7}}, []string{"A"}, []int{2})
	_ = bn.AddCPD(cpdB)
	return bn
}

func makeSmallDigraph() *graphgo.DiGraph {
	g := graphgo.NewDiGraph()
	g.AddNode("A")
	g.AddNode("B")
	g.AddNode("C")
	g.AddEdge("A", "B")
	return g
}
