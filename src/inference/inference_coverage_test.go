//go:build unit

package inference

import (
	"strings"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// ---------------------------------------------------------------------------
// parseEdgeKey / normalizeEdgeKey coverage
// ---------------------------------------------------------------------------

func TestParseEdgeKey_HyphenFormat(t *testing.T) {
	a, b := parseEdgeKey("3-7")
	if a != 3 || b != 7 {
		t.Errorf("expected (3,7), got (%d,%d)", a, b)
	}
}

func TestParseEdgeKey_NULFormat(t *testing.T) {
	a, b := parseEdgeKey("1\x002")
	if a != 1 || b != 2 {
		t.Errorf("expected (1,2), got (%d,%d)", a, b)
	}
}

func TestParseEdgeKey_Invalid(t *testing.T) {
	a, b := parseEdgeKey("garbage")
	if a != -1 || b != -1 {
		t.Errorf("expected (-1,-1), got (%d,%d)", a, b)
	}
}

func TestParseEdgeKey_BadNUL(t *testing.T) {
	a, b := parseEdgeKey("abc\x00def")
	if a != -1 || b != -1 {
		t.Errorf("expected (-1,-1) for non-numeric NUL parts, got (%d,%d)", a, b)
	}
}

func TestNormalizeEdgeKey_HyphenFormat(t *testing.T) {
	k := normalizeEdgeKey("5-3")
	if k != "3-5" {
		t.Errorf("expected '3-5', got %q", k)
	}
}

func TestNormalizeEdgeKey_Fallback(t *testing.T) {
	k := normalizeEdgeKey("garbage")
	if k != "garbage" {
		t.Errorf("expected 'garbage', got %q", k)
	}
}

// ---------------------------------------------------------------------------
// BeliefPropagation coverage
// ---------------------------------------------------------------------------

func makeSimpleBP(t *testing.T) *BeliefPropagation {
	t.Helper()
	// Two cliques: {A, B} and {B, C}, separator B.
	cliques := [][]string{{"A", "B"}, {"B", "C"}}
	separators := map[string][]string{
		"0-1": {"B"},
	}

	fAB, err := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.3, 0.7, 0.4, 0.6})
	if err != nil {
		t.Fatal(err)
	}
	fBC, err := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{0.2, 0.8, 0.5, 0.5})
	if err != nil {
		t.Fatal(err)
	}

	cliqueFactors := map[int][]*factors.DiscreteFactor{
		0: {fAB},
		1: {fBC},
	}

	return NewBeliefPropagation(cliques, separators, cliqueFactors)
}

func TestBP_InitializePotentials_MultipleFactors(t *testing.T) {
	// Test multiple factors per clique (factor product path).
	fA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	fA2, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})

	cliques := [][]string{{"A"}}
	cliqueFactors := map[int][]*factors.DiscreteFactor{
		0: {fA, fA2}, // two factors in one clique
	}

	bp := NewBeliefPropagation(cliques, nil, cliqueFactors)
	err := bp.Calibrate()
	if err != nil {
		t.Fatalf("Calibrate failed: %v", err)
	}
}

func TestBP_SingleClique(t *testing.T) {
	// Single clique => trivially calibrated.
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	cliques := [][]string{{"A"}}
	cliqueFactors := map[int][]*factors.DiscreteFactor{0: {f}}

	bp := NewBeliefPropagation(cliques, nil, cliqueFactors)
	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}
	if !bp.IsCalibrated() {
		t.Error("expected calibrated")
	}
}

func TestBP_EmptyCliques(t *testing.T) {
	bp := NewBeliefPropagation(nil, nil, nil)
	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}
	if !bp.IsCalibrated() {
		t.Error("expected calibrated for empty cliques")
	}
}

func TestBP_GetCliqueBelief_OutOfRange(t *testing.T) {
	bp := makeSimpleBP(t)
	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}
	if bp.GetCliqueBelief(-1) != nil {
		t.Error("expected nil for negative index")
	}
	if bp.GetCliqueBelief(999) != nil {
		t.Error("expected nil for out-of-range index")
	}
}

func TestBP_GetCliqueBelief_NilPotential(t *testing.T) {
	bp := makeSimpleBP(t)
	// Don't calibrate => potentials are nil.
	if bp.GetCliqueBelief(0) != nil {
		t.Error("expected nil for uncalibrated clique")
	}
}

func TestBP_GetCliqueBeliefs(t *testing.T) {
	bp := makeSimpleBP(t)
	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}
	beliefs := bp.GetCliqueBeliefs()
	if len(beliefs) != 2 {
		t.Errorf("expected 2 beliefs, got %d", len(beliefs))
	}
}

func TestBP_GetCliques(t *testing.T) {
	bp := makeSimpleBP(t)
	cliques := bp.GetCliques()
	if len(cliques) != 2 {
		t.Errorf("expected 2 cliques, got %d", len(cliques))
	}
}

func TestBP_GetSepsetBeliefs_Calibrated(t *testing.T) {
	bp := makeSimpleBP(t)
	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}
	beliefs := bp.GetSepsetBeliefs()
	if len(beliefs) != 1 {
		t.Errorf("expected 1 separator belief, got %d", len(beliefs))
	}
	for _, v := range beliefs {
		if v == nil {
			t.Error("expected non-nil separator belief")
		}
	}
}

func TestBP_GetSepsetBeliefs_Uncalibrated(t *testing.T) {
	bp := makeSimpleBP(t)
	beliefs := bp.GetSepsetBeliefs()
	for _, v := range beliefs {
		if v != nil {
			t.Error("expected nil for uncalibrated separator")
		}
	}
}

func TestBP_String(t *testing.T) {
	bp := makeSimpleBP(t)
	s := bp.String()
	if !strings.Contains(s, "BeliefPropagation") {
		t.Error("expected 'BeliefPropagation' in string")
	}
	if !strings.Contains(s, "clique") {
		t.Error("expected 'clique' in string")
	}
	if !strings.Contains(s, "separator") {
		t.Error("expected 'separator' in string")
	}
}

func TestBP_Query_WithEvidence(t *testing.T) {
	bp := makeSimpleBP(t)
	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}
	result, err := bp.Query([]string{"A"}, map[string]int{"B": 0})
	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Error("expected non-nil result")
	}
}

func TestBP_MAPQuery(t *testing.T) {
	bp := makeSimpleBP(t)
	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}
	result, err := bp.MAPQuery([]string{"A"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["A"]; !ok {
		t.Error("expected 'A' in MAP result")
	}
}

func TestBP_MAPQuery_WithEvidence(t *testing.T) {
	bp := makeSimpleBP(t)
	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}
	result, err := bp.MAPQuery([]string{"A"}, map[string]int{"C": 0})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["A"]; !ok {
		t.Error("expected 'A' in MAP result")
	}
}

func TestBP_MAPQuery_EmptyVars(t *testing.T) {
	bp := makeSimpleBP(t)
	_, err := bp.MAPQuery(nil, nil)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}
}

func TestBP_MaxCalibrate_SingleClique(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	cliques := [][]string{{"A"}}
	cliqueFactors := map[int][]*factors.DiscreteFactor{0: {f}}
	bp := NewBeliefPropagation(cliques, nil, cliqueFactors)
	if err := bp.MaxCalibrate(); err != nil {
		t.Fatal(err)
	}
}

func TestBP_MaxCalibrate_Empty(t *testing.T) {
	bp := NewBeliefPropagation(nil, nil, nil)
	if err := bp.MaxCalibrate(); err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// CausalInference coverage
// ---------------------------------------------------------------------------

func makeSimpleCausalBN(t *testing.T) *models.BayesianNetwork {
	t.Helper()
	bn := models.NewBayesianNetwork()
	bn.AddNode("Z")
	bn.AddNode("X")
	bn.AddNode("Y")
	bn.AddEdge("Z", "X")
	bn.AddEdge("X", "Y")
	bn.SetStates("Z", []string{"z0", "z1"})
	bn.SetStates("X", []string{"x0", "x1"})
	bn.SetStates("Y", []string{"y0", "y1"})

	cpdZ, _ := factors.NewTabularCPD("Z", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	_ = bn.AddCPD(cpdZ)
	cpdX, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.8, 0.2}, {0.2, 0.8}}, []string{"Z"}, []int{2})
	_ = bn.AddCPD(cpdX)
	cpdY, _ := factors.NewTabularCPD("Y", 2, [][]float64{{0.9, 0.3}, {0.1, 0.7}}, []string{"X"}, []int{2})
	_ = bn.AddCPD(cpdY)
	return bn
}

func TestCausalInference_NilBN(t *testing.T) {
	_, err := NewCausalInference(nil)
	if err == nil {
		t.Error("expected error for nil BN")
	}
}

func TestCausalInference_IsValidFrontdoor(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}

	// X is a mediator between Z and Y. But Z->X->Y means X is a valid frontdoor from Z to Y.
	valid := ci.IsValidFrontdoorAdjustmentSet("Z", "Y", []string{"X"})
	// For Z->X->Y, X intercepts all paths from Z to Y and satisfies frontdoor.
	_ = valid
}

func TestCausalInference_IsValidFrontdoor_EmptySet(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}
	valid := ci.IsValidFrontdoorAdjustmentSet("Z", "Y", nil)
	if valid {
		t.Error("expected false for empty frontdoor set")
	}
}

func TestCausalInference_GetAllFrontdoorSets(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}
	sets := ci.GetAllFrontdoorAdjustmentSets("Z", "Y")
	_ = sets // just ensure no crash
}

func TestCausalInference_GetAllBackdoorSets(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}
	sets := ci.GetAllBackdoorAdjustmentSets("X", "Y")
	_ = sets
}

func TestCausalInference_GetScalingIndicators(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}
	indicators := ci.GetScalingIndicators("X", "Y")
	// Z is the only root (exogenous) variable.
	if len(indicators) != 1 || indicators[0] != "Z" {
		t.Errorf("expected [Z], got %v", indicators)
	}
}

func TestCausalInference_GetIVs(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}
	ivs := ci.GetIVs("X", "Y")
	// Z should be an IV for X->Y: Z is associated with X, affects Y only through X.
	found := false
	for _, iv := range ivs {
		if iv == "Z" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected Z as IV, got %v", ivs)
	}
}

func TestCausalInference_GetConditionalIVs(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}
	ivs := ci.GetConditionalIVs("X", "Y", nil)
	_ = ivs
}

func TestCausalInference_Query_EmptyQueryVars(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatal(err)
	}
	_, err = ci.Query(nil, map[string]int{"X": 0}, nil)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}
}

func TestInterceptsAllPaths(t *testing.T) {
	bn := makeSimpleCausalBN(t)
	g := bnToDigraph(bn)

	// {X} intercepts all paths from Z to Y.
	intercepts := interceptsAllPaths(g, "Z", "Y", map[string]bool{"X": true})
	if !intercepts {
		t.Error("expected X to intercept all Z->Y paths")
	}

	// Empty set does NOT intercept (there is a path Z->X->Y).
	intercepts = interceptsAllPaths(g, "Z", "Y", map[string]bool{})
	if intercepts {
		t.Error("expected empty set to not intercept paths")
	}
}

// ---------------------------------------------------------------------------
// ApproxInference coverage
// ---------------------------------------------------------------------------

func TestApproxInference_GetDistribution(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	ai := NewApproxInference([]*factors.DiscreteFactor{f}, 42)
	dist, err := ai.GetDistribution(1000)
	if err != nil {
		t.Fatal(err)
	}
	if dist == nil {
		t.Error("expected non-nil distribution")
	}
}

func TestApproxInference_MAPQuery(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	ai := NewApproxInference([]*factors.DiscreteFactor{f}, 42)
	result, err := ai.MAPQuery([]string{"A"}, nil, 1000)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := result["A"]; !ok {
		t.Error("expected 'A' in MAP result")
	}
}
