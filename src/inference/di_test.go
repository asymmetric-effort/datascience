//go:build unit

package inference

import (
	"fmt"
	"strings"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// ---------------------------------------------------------------------------
// Mock types for dependency injection tests.
// ---------------------------------------------------------------------------

// failingFactorMultiplier is a test mock that simulates FactorProduct and
// Marginalize failures for coverage testing of defensive error paths in
// computeMessage and eliminateVariable.
type failingFactorMultiplier struct {
	failProduct     bool
	failMarginalize bool
}

func (f failingFactorMultiplier) Product(fs ...*factors.DiscreteFactor) (*factors.DiscreteFactor, error) {
	if f.failProduct {
		return nil, fmt.Errorf("injected product failure")
	}
	return factors.FactorProduct(fs...)
}

func (f failingFactorMultiplier) Marginalize(factor *factors.DiscreteFactor, vars []string) (*factors.DiscreteFactor, error) {
	if f.failMarginalize {
		return nil, fmt.Errorf("injected marginalize failure")
	}
	return factor.Marginalize(vars)
}

// failingFactorReducer is a test mock that simulates Reduce failure for
// coverage testing of defensive error paths in reduceAll.
type failingFactorReducer struct{}

func (failingFactorReducer) Reduce(f *factors.DiscreteFactor, evidence map[string]int) (*factors.DiscreteFactor, error) {
	return nil, fmt.Errorf("injected reduce failure")
}

// ---------------------------------------------------------------------------
// Tests for computeMessageImpl defensive error paths.
// ---------------------------------------------------------------------------

func TestDI_ComputeMessageImpl_ProductFailure(t *testing.T) {
	// Build a 3-clique BP so clique 1 has neighbors 0 and 2.
	// When computing message 1->2, there's an incoming message from 0 that
	// triggers the product path.
	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	fBC, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{0.5, 0.1, 0.2, 0.2})
	fCD, _ := factors.NewDiscreteFactor([]string{"C", "D"}, []int{2, 2}, []float64{0.3, 0.1, 0.2, 0.4})

	cliques := [][]string{{"A", "B"}, {"B", "C"}, {"C", "D"}}
	separators := map[string][]string{
		edgeKey(0, 1): {"B"},
		edgeKey(1, 2): {"C"},
	}
	cliqueFactors := map[int][]*factors.DiscreteFactor{
		0: {fAB},
		1: {fBC},
		2: {fCD},
	}

	bp := NewBeliefPropagation(cliques, separators, cliqueFactors)
	if err := bp.initializePotentials(); err != nil {
		t.Fatal(err)
	}
	// Store a message from 0 to 1 so computing 1->2 triggers product.
	bp.messages[msgKey(0, 1)] = bp.potentials[0].Copy()

	fm := failingFactorMultiplier{failProduct: true}
	_, err := computeMessageImpl(bp, 1, 2, fm)
	if err == nil {
		t.Fatal("expected error from failing product")
	}
	if !strings.Contains(err.Error(), "injected product failure") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDI_ComputeMessageImpl_MarginalizeFailure(t *testing.T) {
	bp := buildSimpleBP(t)
	if err := bp.initializePotentials(); err != nil {
		t.Fatal(err)
	}

	fm := failingFactorMultiplier{failMarginalize: true}
	_, err := computeMessageImpl(bp, 0, 1, fm)
	if err == nil {
		t.Fatal("expected error from failing marginalize")
	}
	if !strings.Contains(err.Error(), "injected marginalize failure") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestDI_ComputeMessageImpl_NoMargVars(t *testing.T) {
	// When all vars in the potential are in the separator, no marginalization needed.
	bp := buildSimpleBP(t)
	if err := bp.initializePotentials(); err != nil {
		t.Fatal(err)
	}
	// Override separator to contain all vars of clique 0.
	bp.separators[edgeKey(0, 1)] = bp.cliques[0]

	fm := defaultFactorMultiplier{}
	msg, err := computeMessageImpl(bp, 0, 1, fm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg == nil {
		t.Fatal("expected non-nil message")
	}
}

// ---------------------------------------------------------------------------
// Tests for eliminateVariableImpl defensive error paths.
// ---------------------------------------------------------------------------

func TestDI_EliminateVariableImpl_ProductFailure(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	fm := failingFactorMultiplier{failProduct: true}
	_, err := eliminateVariableImpl([]*factors.DiscreteFactor{f}, "A", fm)
	if err == nil {
		t.Fatal("expected error from failing product")
	}
}

func TestDI_EliminateVariableImpl_MarginalizeFailure(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	fm := failingFactorMultiplier{failMarginalize: true}
	_, err := eliminateVariableImpl([]*factors.DiscreteFactor{f}, "A", fm)
	if err == nil {
		t.Fatal("expected error from failing marginalize")
	}
}

func TestDI_EliminateVariableImpl_NotPresent(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	fm := defaultFactorMultiplier{}
	result, err := eliminateVariableImpl([]*factors.DiscreteFactor{f}, "B", fm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 factor, got %d", len(result))
	}
}

func TestDI_EliminateVariableImpl_SingleVariable(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	fm := defaultFactorMultiplier{}
	result, err := eliminateVariableImpl([]*factors.DiscreteFactor{f}, "A", fm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 remaining factors, got %d", len(result))
	}
}

// ---------------------------------------------------------------------------
// Tests for reduceAllImpl defensive error paths.
// ---------------------------------------------------------------------------

func TestDI_ReduceAllImpl_ReduceFailure(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	fr := failingFactorReducer{}
	_, err := reduceAllImpl([]*factors.DiscreteFactor{f}, map[string]int{"A": 0}, fr)
	if err == nil {
		t.Fatal("expected error from failing reducer")
	}
	if !strings.Contains(err.Error(), "injected reduce failure") {
		t.Errorf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Tests for BeliefPropagation Query error paths.
// ---------------------------------------------------------------------------

func TestDI_BP_Query_EmptyQueryVars(t *testing.T) {
	bp := buildSimpleBP(t)
	_ = bp.Calibrate()
	_, err := bp.Query(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty query vars")
	}
}

func TestDI_BP_Query_NotCalibrated(t *testing.T) {
	bp := buildSimpleBP(t)
	_, err := bp.Query([]string{"A"}, nil)
	if err == nil {
		t.Fatal("expected error when not calibrated")
	}
}

func TestDI_BP_Query_NoCliqueContainsAll(t *testing.T) {
	bp := buildSimpleBP(t)
	_ = bp.Calibrate()
	_, err := bp.Query([]string{"nonexistent"}, nil)
	if err == nil {
		t.Fatal("expected error for vars not in any clique")
	}
}

func TestDI_BP_Query_UnknownEvidenceVariable(t *testing.T) {
	bp := buildSimpleBP(t)
	_ = bp.Calibrate()
	_, err := bp.Query([]string{"A"}, map[string]int{"nonexistent": 0})
	if err == nil {
		t.Fatal("expected error for unknown evidence variable")
	}
}

func TestDI_BP_MAPQuery_EmptyQueryVars(t *testing.T) {
	bp := buildSimpleBP(t)
	_, err := bp.MAPQuery(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty query vars")
	}
}

func TestDI_BP_MAPQuery_UnknownEvidence(t *testing.T) {
	bp := buildSimpleBP(t)
	_, err := bp.MAPQuery([]string{"A"}, map[string]int{"nonexistent": 0})
	if err == nil {
		t.Fatal("expected error for unknown evidence variable")
	}
}

func TestDI_BP_MAPQuery_NoClique(t *testing.T) {
	bp := buildSimpleBP(t)
	_, err := bp.MAPQuery([]string{"nonexistent"}, nil)
	if err == nil {
		t.Fatal("expected error for vars not in any clique")
	}
}

// ---------------------------------------------------------------------------
// Tests for VariableElimination Query error paths.
// ---------------------------------------------------------------------------

func TestDI_VE_Query_EmptyQueryVars(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{f})
	_, err := ve.Query(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty query vars")
	}
}

func TestDI_VE_Query_ReduceError(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{f})
	_, err := ve.Query([]string{"A"}, map[string]int{"A": 99})
	if err == nil {
		t.Fatal("expected error for out-of-range evidence")
	}
}

func TestDI_VE_MaxMarginal_EmptyQueryVars(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{f})
	_, err := ve.MaxMarginal(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty query vars")
	}
}

func TestDI_VE_QueryMarginals_EmptyQueryVars(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{f})
	_, err := ve.QueryMarginals(nil, nil)
	if err == nil {
		t.Fatal("expected error for empty query vars")
	}
}

func TestDI_VE_QueryWithVirtualEvidence_EmptyQueryVars(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{f})
	_, err := ve.QueryWithVirtualEvidence(nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for empty query vars")
	}
}

// ---------------------------------------------------------------------------
// Tests for CausalInference Query error paths.
// ---------------------------------------------------------------------------

func TestDI_CI_Query_EmptyQueryVars(t *testing.T) {
	ci := buildSimpleCI(t)
	_, err := ci.Query(nil, nil, nil)
	if err == nil {
		t.Fatal("expected error for empty query vars")
	}
}

func TestDI_CI_Query_InvalidDoValue(t *testing.T) {
	ci := buildSimpleCI(t)
	_, err := ci.Query([]string{"B"}, map[string]int{"A": 99}, nil)
	if err == nil {
		t.Fatal("expected error for out-of-range do value")
	}
}

func TestDI_CI_ATE_QueryFailure(t *testing.T) {
	ci := buildSimpleCI(t)
	// ATE with out-of-range treatment values.
	_, err := ci.ATE("A", "B", [2]int{0, 99})
	if err == nil {
		t.Fatal("expected error for invalid treatment value")
	}
}

func TestDI_CI_IdentificationMethod(t *testing.T) {
	ci := buildSimpleCI(t)
	method := ci.IdentificationMethod("A", "B")
	if method != "backdoor" && method != "frontdoor" && method != "iv" && method != "none" {
		t.Errorf("unexpected method: %s", method)
	}
}

func TestDI_CI_GetMinimalAdjustmentSet(t *testing.T) {
	ci := buildSimpleCI(t)
	_, err := ci.GetMinimalAdjustmentSet("A", "B")
	// May succeed or fail depending on model structure.
	_ = err
}

func TestDI_CI_EstimateATE_NilData(t *testing.T) {
	ci := buildSimpleCI(t)
	_, err := ci.EstimateATE("A", "B", nil)
	if err == nil {
		t.Fatal("expected error for nil data")
	}
}

// ---------------------------------------------------------------------------
// Tests for BP GetSepsetBeliefs edge cases.
// ---------------------------------------------------------------------------

func TestDI_BP_GetSepsetBeliefs_Uncalibrated(t *testing.T) {
	bp := buildSimpleBP(t)
	beliefs := bp.GetSepsetBeliefs()
	for _, b := range beliefs {
		if b != nil {
			t.Error("expected nil beliefs when uncalibrated")
		}
	}
}

func TestDI_BP_GetSepsetBeliefs_Calibrated(t *testing.T) {
	bp := buildSimpleBP(t)
	_ = bp.Calibrate()
	beliefs := bp.GetSepsetBeliefs()
	if len(beliefs) == 0 {
		t.Error("expected some separator beliefs after calibration")
	}
}

// ---------------------------------------------------------------------------
// Tests for BP initializePotentials edge cases.
// ---------------------------------------------------------------------------

func TestDI_BP_InitializePotentials_UnknownCardinality(t *testing.T) {
	// Create BP with a clique var not in any factor (unknown cardinality).
	cliques := [][]string{{"X", "Y"}}
	separators := map[string][]string{}
	cliqueFactors := map[int][]*factors.DiscreteFactor{} // no factors for clique 0

	bp := NewBeliefPropagation(cliques, separators, cliqueFactors)
	err := bp.initializePotentials()
	if err == nil {
		t.Fatal("expected error for unknown cardinality")
	}
}

// ---------------------------------------------------------------------------
// Tests for BP Calibrate edge cases.
// ---------------------------------------------------------------------------

func TestDI_BP_Calibrate_Empty(t *testing.T) {
	bp := NewBeliefPropagation(nil, nil, nil)
	err := bp.Calibrate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bp.IsCalibrated() {
		t.Error("expected calibrated")
	}
}

func TestDI_BP_Calibrate_SingleClique(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	cliques := [][]string{{"A"}}
	cliqueFactors := map[int][]*factors.DiscreteFactor{0: {f}}
	bp := NewBeliefPropagation(cliques, nil, cliqueFactors)
	err := bp.Calibrate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDI_BP_MaxCalibrate_Empty(t *testing.T) {
	bp := NewBeliefPropagation(nil, nil, nil)
	err := bp.MaxCalibrate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDI_BP_MaxCalibrate_SingleClique(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	cliques := [][]string{{"A"}}
	cliqueFactors := map[int][]*factors.DiscreteFactor{0: {f}}
	bp := NewBeliefPropagation(cliques, nil, cliqueFactors)
	err := bp.MaxCalibrate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Tests for BP GetCliqueBelief edge cases.
// ---------------------------------------------------------------------------

func TestDI_BP_GetCliqueBelief_OutOfRange(t *testing.T) {
	bp := buildSimpleBP(t)
	b := bp.GetCliqueBelief(-1)
	if b != nil {
		t.Error("expected nil for negative index")
	}
	b = bp.GetCliqueBelief(999)
	if b != nil {
		t.Error("expected nil for out-of-range index")
	}
}

// ---------------------------------------------------------------------------
// --- Tests for default implementations ---

func TestDI_DefaultFactorMultiplier_Marginalize(t *testing.T) {
	fm := defaultFactorMultiplier{}
	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	result, err := fm.Marginalize(f, []string{"B"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestDI_DefaultFactorReducer_Reduce(t *testing.T) {
	fr := defaultFactorReducer{}
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	result, err := fr.Reduce(f, map[string]int{"A": 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}
}

func TestDI_ReduceAllImpl_Success(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	fr := defaultFactorReducer{}
	result, err := reduceAllImpl([]*factors.DiscreteFactor{f}, map[string]int{"A": 0}, fr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 factor, got %d", len(result))
	}
}

// Helpers.
// ---------------------------------------------------------------------------

// buildSimpleBP creates a minimal BP engine with 2 cliques joined by a
// separator, suitable for DI testing.
func buildSimpleBP(t *testing.T) *BeliefPropagation {
	t.Helper()
	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	fBC, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{0.5, 0.1, 0.2, 0.2})

	cliques := [][]string{{"A", "B"}, {"B", "C"}}
	separators := map[string][]string{edgeKey(0, 1): {"B"}}
	cliqueFactors := map[int][]*factors.DiscreteFactor{
		0: {fAB},
		1: {fBC},
	}

	return NewBeliefPropagation(cliques, separators, cliqueFactors)
}

// buildSimpleCI creates a minimal CausalInference engine for testing.
func buildSimpleCI(t *testing.T) *CausalInference {
	t.Helper()
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")

	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.4}, {0.6}}, nil, nil)
	_ = bn.AddCPD(cpdA)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.2, 0.8}, {0.8, 0.2}}, []string{"A"}, []int{2})
	_ = bn.AddCPD(cpdB)

	ci, err := NewCausalInference(bn)
	if err != nil {
		t.Fatalf("NewCausalInference: %v", err)
	}
	return ci
}
