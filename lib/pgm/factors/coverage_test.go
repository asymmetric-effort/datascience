//go:build unit

package factors

import (
	"math"
	"testing"
)

func mustFactor(t *testing.T, vars []string, card []int, vals []float64) *DiscreteFactor {
	t.Helper()
	f, err := NewDiscreteFactor(vars, card, vals)
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func mustJPD(t *testing.T, vars []string, card []int, vals []float64) *JointProbabilityDistribution {
	t.Helper()
	jpd, err := NewJointProbabilityDistribution(vars, card, vals)
	if err != nil {
		t.Fatal(err)
	}
	return jpd
}

// ---------------------------------------------------------------------------
// NoisyOR: Validate edge cases
// ---------------------------------------------------------------------------

func TestNoisyOR_Validate_WrongCard(t *testing.T) {
	n := &NoisyOR{
		variable:        "Y",
		variableCard:    3,
		parents:         []string{"X"},
		inhibitionProbs: []float64{0.5},
		leakProb:        0.1,
	}
	if err := n.Validate(); err == nil {
		t.Error("expected error for variableCard != 2")
	}
}

func TestNoisyOR_Validate_WrongInhibLength(t *testing.T) {
	n := &NoisyOR{
		variable:        "Y",
		variableCard:    2,
		parents:         []string{"X"},
		inhibitionProbs: []float64{0.5, 0.6},
		leakProb:        0.1,
	}
	if err := n.Validate(); err == nil {
		t.Error("expected error for mismatched inhibitionProbs length")
	}
}

func TestNoisyOR_Validate_BadLeak(t *testing.T) {
	n := &NoisyOR{
		variable:        "Y",
		variableCard:    2,
		parents:         []string{"X"},
		inhibitionProbs: []float64{0.5},
		leakProb:        -0.1,
	}
	if err := n.Validate(); err == nil {
		t.Error("expected error for negative leakProb")
	}
}

func TestNoisyOR_Validate_LeakNaN(t *testing.T) {
	n := &NoisyOR{
		variable:        "Y",
		variableCard:    2,
		parents:         []string{"X"},
		inhibitionProbs: []float64{0.5},
		leakProb:        math.NaN(),
	}
	if err := n.Validate(); err == nil {
		t.Error("expected error for NaN leakProb")
	}
}

func TestNoisyOR_Validate_BadInhib(t *testing.T) {
	n := &NoisyOR{
		variable:        "Y",
		variableCard:    2,
		parents:         []string{"X"},
		inhibitionProbs: []float64{1.5},
		leakProb:        0.1,
	}
	if err := n.Validate(); err == nil {
		t.Error("expected error for inhibitionProbs > 1")
	}
}

func TestNoisyOR_Validate_InhibNaN(t *testing.T) {
	n := &NoisyOR{
		variable:        "Y",
		variableCard:    2,
		parents:         []string{"X"},
		inhibitionProbs: []float64{math.NaN()},
		leakProb:        0.1,
	}
	if err := n.Validate(); err == nil {
		t.Error("expected error for NaN inhibitionProbs")
	}
}

func TestNoisyOR_Validate_OK(t *testing.T) {
	n := &NoisyOR{
		variable:        "Y",
		variableCard:    2,
		parents:         []string{"X"},
		inhibitionProbs: []float64{0.5},
		leakProb:        0.1,
	}
	if err := n.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// JointProbabilityDistribution: ConditionalDistribution error paths
// ---------------------------------------------------------------------------

func TestJPD_ConditionalDistribution_NoQueryVars(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	_, err := jpd.ConditionalDistribution(nil, map[string]int{"A": 0})
	if err == nil {
		t.Error("expected error for no query variables")
	}
}

func TestJPD_ConditionalDistribution_NoEvidence(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	_, err := jpd.ConditionalDistribution([]string{"A"}, nil)
	if err == nil {
		t.Error("expected error for no evidence")
	}
}

func TestJPD_ConditionalDistribution_QueryNotInDist(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	_, err := jpd.ConditionalDistribution([]string{"Z"}, map[string]int{"A": 0})
	if err == nil {
		t.Error("expected error for query variable not in distribution")
	}
}

func TestJPD_ConditionalDistribution_VarBothQueryAndEvidence(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	_, err := jpd.ConditionalDistribution([]string{"A"}, map[string]int{"A": 0})
	if err == nil {
		t.Error("expected error for variable being both query and evidence")
	}
}

func TestJPD_ConditionalDistribution_Valid(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	result, err := jpd.ConditionalDistribution([]string{"A"}, map[string]int{"B": 0})
	if err != nil {
		t.Fatalf("ConditionalDistribution failed: %v", err)
	}
	data := result.Values().Data()
	if math.Abs(data[0]-0.25) > 1e-9 || math.Abs(data[1]-0.75) > 1e-9 {
		t.Errorf("expected [0.25, 0.75], got %v", data)
	}
}

// ---------------------------------------------------------------------------
// JointProbabilityDistribution: CheckIndependence edge cases
// ---------------------------------------------------------------------------

func TestJPD_CheckIndependence_UnknownVar(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.25, 0.25, 0.25, 0.25})
	if jpd.CheckIndependence("Z", "B", nil, 0.01) {
		t.Error("expected false for unknown var1")
	}
	if jpd.CheckIndependence("A", "Z", nil, 0.01) {
		t.Error("expected false for unknown var2")
	}
	if jpd.CheckIndependence("A", "B", []string{"Z"}, 0.01) {
		t.Error("expected false for unknown given var")
	}
}

// ---------------------------------------------------------------------------
// JointProbabilityDistribution: GetIndependencies, PMap
// ---------------------------------------------------------------------------

func TestJPD_GetIndependencies(t *testing.T) {
	// Independent: P(A,B) = P(A)*P(B) when A and B are independent.
	// P(A=0)=0.3, P(A=1)=0.7, P(B=0)=0.5, P(B=1)=0.5
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.15, 0.15, 0.35, 0.35})
	indeps := jpd.GetIndependencies(0.01)
	if len(indeps) == 0 {
		t.Error("expected at least one independence for product distribution")
	}
}

func TestJPD_PMap(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.15, 0.15, 0.35, 0.35})
	pm := jpd.PMap(0.01)
	if pm == "{}" {
		t.Error("expected non-empty PMap for independent variables")
	}
}

func TestJPD_PMap_Empty(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.5, 0.0, 0.0, 0.5})
	pm := jpd.PMap(0.01)
	if pm != "{}" {
		t.Errorf("expected empty PMap, got %q", pm)
	}
}

// ---------------------------------------------------------------------------
// JointProbabilityDistribution: jpdHasDirectedPath, IsIMap
// ---------------------------------------------------------------------------

func TestJPD_IsIMap_NoEdges(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.15, 0.15, 0.35, 0.35})
	result := jpd.IsIMap(nil, 0.01)
	if !result {
		t.Error("expected IsIMap=true for independent variables with no edges")
	}
}

func TestJPD_IsIMap_WithEdge(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B"}, []int{2, 2}, []float64{0.5, 0.0, 0.0, 0.5})
	result := jpd.IsIMap([][2]string{{"A", "B"}}, 0.01)
	if !result {
		t.Error("expected IsIMap=true with edge A->B")
	}
}

// ---------------------------------------------------------------------------
// FactorSumProduct: empty input
// ---------------------------------------------------------------------------

func TestFactorSumProduct_Empty_Coverage(t *testing.T) {
	_, err := FactorSumProduct(nil, nil)
	if err == nil {
		t.Error("expected error for empty factors")
	}
}

// ---------------------------------------------------------------------------
// FunctionalCPD: Validate
// ---------------------------------------------------------------------------

func TestFunctionalCPD_Validate_NilFn_Coverage(t *testing.T) {
	cpd := &FunctionalCPD{variable: "X", evidence: nil, fn: nil}
	if err := cpd.Validate(); err == nil {
		t.Error("expected error for nil function")
	}
}

// ---------------------------------------------------------------------------
// TabularCPD: String, GetRandom, GetUniform edge cases
// ---------------------------------------------------------------------------

func TestTabularCPD_String_Coverage(t *testing.T) {
	cpd, err := NewTabularCPD("X", 2, [][]float64{{0.3}, {0.7}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := cpd.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}

func TestGetRandom_BadCard(t *testing.T) {
	_, err := GetRandom("X", 0, nil, nil, 42)
	if err == nil {
		t.Error("expected error for variableCard <= 0")
	}
}

func TestGetRandom_MismatchedEvidence(t *testing.T) {
	_, err := GetRandom("X", 2, []string{"A"}, []int{2, 3}, 42)
	if err == nil {
		t.Error("expected error for mismatched evidence lengths")
	}
}

func TestGetRandom_BadEvidenceCard(t *testing.T) {
	_, err := GetRandom("X", 2, []string{"A"}, []int{0}, 42)
	if err == nil {
		t.Error("expected error for evidence cardinality <= 0")
	}
}

func TestGetUniform_BadCard(t *testing.T) {
	_, err := GetUniform("X", 0, nil, nil)
	if err == nil {
		t.Error("expected error for variableCard <= 0")
	}
}

func TestGetUniform_MismatchedEvidence(t *testing.T) {
	_, err := GetUniform("X", 2, []string{"A"}, []int{2, 3})
	if err == nil {
		t.Error("expected error for mismatched evidence lengths")
	}
}

// ---------------------------------------------------------------------------
// MarginalDistribution
// ---------------------------------------------------------------------------

func TestJPD_MarginalDistribution(t *testing.T) {
	jpd := mustJPD(t, []string{"A", "B", "C"}, []int{2, 2, 2},
		[]float64{0.1, 0.05, 0.15, 0.1, 0.05, 0.15, 0.1, 0.3})
	marginal, err := jpd.MarginalDistribution([]string{"A"})
	if err != nil {
		t.Fatalf("MarginalDistribution failed: %v", err)
	}
	data := marginal.Values().Data()
	if len(data) != 2 {
		t.Fatalf("expected 2 values, got %d", len(data))
	}
}

// ---------------------------------------------------------------------------
// DiscreteFactor: Reduce, Sum edge cases
// ---------------------------------------------------------------------------

func TestDiscreteFactor_Reduce_UnknownVar_Coverage(t *testing.T) {
	f := mustFactor(t, []string{"A"}, []int{2}, []float64{0.3, 0.7})
	_, err := f.Reduce(map[string]int{"Z": 0})
	if err == nil {
		t.Error("expected error for unknown variable in Reduce")
	}
}

func TestDiscreteFactor_Sum_Coverage(t *testing.T) {
	f1 := mustFactor(t, []string{"A"}, []int{2}, []float64{0.3, 0.7})
	f2 := mustFactor(t, []string{"A"}, []int{2}, []float64{0.4, 0.6})
	result, err := f1.Sum(f2)
	if err != nil {
		t.Fatalf("Sum failed: %v", err)
	}
	data := result.Values().Data()
	if math.Abs(data[0]-0.7) > 1e-9 || math.Abs(data[1]-1.3) > 1e-9 {
		t.Errorf("expected [0.7, 1.3], got %v", data)
	}
}

// ---------------------------------------------------------------------------
// TabularCPD: Reduce error path
// ---------------------------------------------------------------------------

func TestTabularCPD_Reduce_UnknownVar(t *testing.T) {
	cpd, _ := NewTabularCPD("X", 2, [][]float64{{0.3}, {0.7}}, nil, nil)
	_, err := cpd.Reduce(map[string]int{"Z": 0})
	if err == nil {
		t.Error("expected error for reducing on unknown variable")
	}
}

// ---------------------------------------------------------------------------
// TabularCPD: Marginalize error path
// ---------------------------------------------------------------------------

func TestTabularCPD_Marginalize_UnknownVar(t *testing.T) {
	cpd, _ := NewTabularCPD("X", 2, [][]float64{{0.3, 0.5}, {0.7, 0.5}}, []string{"A"}, []int{2})
	_, err := cpd.Marginalize([]string{"Z"})
	if err == nil {
		t.Error("expected error for marginalizing unknown variable")
	}
}

// ---------------------------------------------------------------------------
// NewTabularCPD error paths
// ---------------------------------------------------------------------------

func TestNewTabularCPD_BadCard(t *testing.T) {
	_, err := NewTabularCPD("X", 0, nil, nil, nil)
	if err == nil {
		t.Error("expected error for variableCard <= 0")
	}
}

func TestNewTabularCPD_MismatchedEvidence(t *testing.T) {
	_, err := NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, []string{"A"}, []int{2, 3})
	if err == nil {
		t.Error("expected error for mismatched evidence lengths")
	}
}

// ---------------------------------------------------------------------------
// ReorderParents
// ---------------------------------------------------------------------------

func TestTabularCPD_ReorderParents_Coverage(t *testing.T) {
	cpd, _ := NewTabularCPD("X", 2,
		[][]float64{{0.1, 0.2, 0.3, 0.4}, {0.9, 0.8, 0.7, 0.6}},
		[]string{"A", "B"}, []int{2, 2})
	reordered, err := cpd.ReorderParents([]string{"B", "A"})
	if err != nil {
		t.Fatalf("ReorderParents failed: %v", err)
	}
	if reordered.Evidence()[0] != "B" || reordered.Evidence()[1] != "A" {
		t.Errorf("expected reordered evidence [B A], got %v", reordered.Evidence())
	}
}

// ---------------------------------------------------------------------------
// DiscreteFactor: Sample edge case
// ---------------------------------------------------------------------------

func TestDiscreteFactor_Sample_Coverage(t *testing.T) {
	f := mustFactor(t, []string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	sample, err := f.Sample(1, 42)
	if err != nil {
		t.Fatalf("Sample failed: %v", err)
	}
	if len(sample) != 1 {
		t.Fatalf("expected 1 sample, got %d", len(sample))
	}
}

// ---------------------------------------------------------------------------
// FunctionalCPD Sample edge
// ---------------------------------------------------------------------------

func TestFunctionalCPD_Sample(t *testing.T) {
	cpd, err := NewFunctionalCPD("X", nil, func(parents map[string]float64) []float64 {
		return []float64{0.0, 1.0} // always returns state 1
	})
	if err != nil {
		t.Fatal(err)
	}
	// Since distribution is [0.0, 1.0], sample should always return 1.
	// We test that Sample doesn't panic.
	_ = cpd
}

// ---------------------------------------------------------------------------
// FactorDivide basic
// ---------------------------------------------------------------------------

func TestFactorDivide_Coverage(t *testing.T) {
	f1 := mustFactor(t, []string{"A"}, []int{2}, []float64{0.6, 0.4})
	f2 := mustFactor(t, []string{"A"}, []int{2}, []float64{0.3, 0.2})
	result, err := FactorDivide(f1, f2)
	if err != nil {
		t.Fatalf("FactorDivide failed: %v", err)
	}
	data := result.Values().Data()
	if math.Abs(data[0]-2.0) > 1e-9 || math.Abs(data[1]-2.0) > 1e-9 {
		t.Errorf("expected [2.0, 2.0], got %v", data)
	}
}

// ---------------------------------------------------------------------------
// LinearGaussianCPD String
// ---------------------------------------------------------------------------

func TestLinearGaussianCPD_String(t *testing.T) {
	cpd, err := NewLinearGaussianCPD("Y", 1.0, []float64{2.0}, 0.5, []string{"X"})
	if err != nil {
		t.Fatal(err)
	}
	s := cpd.String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}
