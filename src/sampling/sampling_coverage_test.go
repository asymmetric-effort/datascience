//go:build unit

package sampling

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/numgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

func makeTwoNodeBN(t *testing.T) *models.BayesianNetwork {
	t.Helper()
	bn := models.NewBayesianNetwork()
	bn.AddNode("A")
	bn.AddNode("B")
	bn.AddEdge("A", "B")
	bn.SetStates("A", []string{"a0", "a1"})
	bn.SetStates("B", []string{"b0", "b1"})

	cpdA, err := factors.NewTabularCPD("A", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = bn.AddCPD(cpdA)

	cpdB, err := factors.NewTabularCPD("B", 2, [][]float64{{0.8, 0.3}, {0.2, 0.7}}, []string{"A"}, []int{2})
	if err != nil {
		t.Fatal(err)
	}
	_ = bn.AddCPD(cpdB)
	return bn
}

func makeSingleNodeBN(t *testing.T) *models.BayesianNetwork {
	t.Helper()
	bn := models.NewBayesianNetwork()
	bn.AddNode("X")
	bn.SetStates("X", []string{"x0", "x1"})
	cpd, err := factors.NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = bn.AddCPD(cpd)
	return bn
}

// TestGibbsSampling_WithEvidence exercises the Gibbs sampler with evidence.
func TestGibbsSampling_WithEvidence_Coverage(t *testing.T) {
	bn := makeTwoNodeBN(t)

	gs, err := NewGibbsSampling(bn, 42)
	if err != nil {
		t.Fatalf("NewGibbsSampling failed: %v", err)
	}

	df, err := gs.Sample(5, 10, 1, map[string]int{"A": 0})
	if err != nil {
		t.Fatalf("Sample failed: %v", err)
	}
	if df.Len() != 5 {
		t.Errorf("expected 5 samples, got %d", df.Len())
	}
}

// TestGibbsSampling_HighThinning exercises thinning > 1.
func TestGibbsSampling_HighThinning(t *testing.T) {
	bn := makeTwoNodeBN(t)

	gs, err := NewGibbsSampling(bn, 42)
	if err != nil {
		t.Fatalf("NewGibbsSampling failed: %v", err)
	}

	df, err := gs.Sample(3, 5, 3, nil)
	if err != nil {
		t.Fatalf("Sample failed: %v", err)
	}
	if df.Len() != 3 {
		t.Errorf("expected 3 samples, got %d", df.Len())
	}
}

// TestSampleCategorical_LastIndex tests the fallback to last index.
func TestSampleCategorical_LastIndex(t *testing.T) {
	rng := numgo.NewRNG(12345)
	probs := []float64{0.0, 0.0, 0.0}
	idx := sampleCategorical(rng, probs)
	if idx != 2 {
		t.Errorf("expected last index 2, got %d", idx)
	}
}

// TestSampleCategorical_SingleElement tests single element.
func TestSampleCategorical_SingleElement(t *testing.T) {
	rng := numgo.NewRNG(42)
	probs := []float64{1.0}
	idx := sampleCategorical(rng, probs)
	if idx != 0 {
		t.Errorf("expected index 0, got %d", idx)
	}
}

// TestForwardSample_NegativeN tests error path.
func TestForwardSample_NegativeN(t *testing.T) {
	bn := makeSingleNodeBN(t)
	bms, err := NewBayesianModelSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}
	_, err = bms.ForwardSample(-1)
	if err == nil {
		t.Error("expected error for negative n")
	}
}

// TestRejectionSample_NegativeN tests error path.
func TestRejectionSample_NegativeN(t *testing.T) {
	bn := makeSingleNodeBN(t)
	bms, err := NewBayesianModelSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}
	_, err = bms.RejectionSample(-1, nil)
	if err == nil {
		t.Error("expected error for negative n")
	}
}

// TestRejectionSample_BadEvidence tests error for unknown evidence variable.
func TestRejectionSample_BadEvidence(t *testing.T) {
	bn := makeSingleNodeBN(t)
	bms, err := NewBayesianModelSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}
	_, err = bms.RejectionSample(5, map[string]int{"NonExistent": 0})
	if err == nil {
		t.Error("expected error for unknown evidence variable")
	}
}

// TestLikelihoodWeightedSample_NegativeN tests error path.
func TestLikelihoodWeightedSample_NegativeN(t *testing.T) {
	bn := makeSingleNodeBN(t)
	bms, err := NewBayesianModelSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = bms.LikelihoodWeightedSample(-1, nil)
	if err == nil {
		t.Error("expected error for negative n")
	}
}

// TestLikelihoodWeightedSample_BadEvidence tests error for unknown evidence.
func TestLikelihoodWeightedSample_BadEvidence(t *testing.T) {
	bn := makeSingleNodeBN(t)
	bms, err := NewBayesianModelSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}
	_, _, err = bms.LikelihoodWeightedSample(5, map[string]int{"NonExistent": 0})
	if err == nil {
		t.Error("expected error for unknown evidence variable")
	}
}

// TestGibbsSampling_ComputeFullConditional_ErrorPath exercises reduce error.
func TestGibbsSampling_ComputeFullConditional_ErrorPath(t *testing.T) {
	// Build a 3-node network to ensure diverse factor reduction paths.
	bn := models.NewBayesianNetwork()
	bn.AddNode("A")
	bn.AddNode("B")
	bn.AddNode("C")
	bn.AddEdge("A", "B")
	bn.AddEdge("B", "C")
	bn.SetStates("A", []string{"a0", "a1"})
	bn.SetStates("B", []string{"b0", "b1"})
	bn.SetStates("C", []string{"c0", "c1"})

	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.3}, {0.7}}, nil, nil)
	_ = bn.AddCPD(cpdA)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.9, 0.1}, {0.1, 0.9}}, []string{"A"}, []int{2})
	_ = bn.AddCPD(cpdB)
	cpdC, _ := factors.NewTabularCPD("C", 2, [][]float64{{0.8, 0.2}, {0.2, 0.8}}, []string{"B"}, []int{2})
	_ = bn.AddCPD(cpdC)

	gs, err := NewGibbsSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}
	df, err := gs.Sample(10, 20, 2, map[string]int{"A": 1})
	if err != nil {
		t.Fatal(err)
	}
	if df.Len() != 10 {
		t.Errorf("expected 10 samples, got %d", df.Len())
	}
}

// TestTopologicalOrder_SingleNode tests topological order with single node.
func TestTopologicalOrder_SingleNode(t *testing.T) {
	bn := makeSingleNodeBN(t)
	order, err := topologicalOrder(bn)
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 1 || order[0] != "X" {
		t.Errorf("expected [X], got %v", order)
	}
}
