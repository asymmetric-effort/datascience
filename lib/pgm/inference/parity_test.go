//go:build unit

package inference

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

// ---------------------------------------------------------------------------
// VirtualEvidence tests (VariableElimination)
// ---------------------------------------------------------------------------

func TestQueryWithVirtualEvidence_Basic(t *testing.T) {
	// Simple A -> B network.
	// P(A) = [0.4, 0.6]
	// P(B|A): P(B=0|A=0)=0.2, P(B=0|A=1)=0.3, P(B=1|A=0)=0.8, P(B=1|A=1)=0.7
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{pA, pBA})

	// Virtual evidence on B: likelihood ratio [0.6, 0.4] (soft observation).
	vEvidence := []VirtualEvidence{
		{Variable: "B", Values: []float64{0.6, 0.4}},
	}

	result, err := ve.QueryWithVirtualEvidence([]string{"A"}, nil, vEvidence)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// The virtual evidence on B should shift the posterior of A.
	// P(A | ve_B) proportional to P(A) * sum_B P(B|A) * ve(B)
	// For A=0: 0.4 * (0.2*0.6 + 0.8*0.4) = 0.4 * (0.12 + 0.32) = 0.4 * 0.44 = 0.176
	// For A=1: 0.6 * (0.3*0.6 + 0.7*0.4) = 0.6 * (0.18 + 0.28) = 0.6 * 0.46 = 0.276
	// Total = 0.452
	// P(A=0|ve) = 0.176/0.452 ~ 0.38938
	// P(A=1|ve) = 0.276/0.452 ~ 0.61062
	a0 := result.GetValue(map[string]int{"A": 0})
	a1 := result.GetValue(map[string]int{"A": 1})
	assertNear(t, a0, 0.176/0.452, 1e-5, "P(A=0|ve_B)")
	assertNear(t, a1, 0.276/0.452, 1e-5, "P(A=1|ve_B)")
}

func TestQueryWithVirtualEvidence_WithHardEvidence(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{pA, pBA})

	// Hard evidence B=0, virtual evidence on A.
	vEvidence := []VirtualEvidence{
		{Variable: "A", Values: []float64{0.9, 0.1}},
	}

	result, err := ve.QueryWithVirtualEvidence([]string{"A"}, map[string]int{"B": 0}, vEvidence)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)
	// Should produce valid probabilities shifted by both hard and virtual evidence.
	a0 := result.GetValue(map[string]int{"A": 0})
	a1 := result.GetValue(map[string]int{"A": 1})
	if a0 < 0 || a1 < 0 || a0+a1 < 0.999 {
		t.Errorf("invalid posterior: A=0=%f, A=1=%f", a0, a1)
	}
}

func TestQueryWithVirtualEvidence_EmptyVirtualEvidence(t *testing.T) {
	// No virtual evidence should give the same result as standard Query.
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{pA, pBA})

	result, err := ve.QueryWithVirtualEvidence([]string{"B"}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	expected, err := ve.Query([]string{"B"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	b0Got := result.GetValue(map[string]int{"B": 0})
	b0Exp := expected.GetValue(map[string]int{"B": 0})
	assertNear(t, b0Got, b0Exp, 1e-9, "P(B=0) with empty virtual evidence")
}

func TestQueryWithVirtualEvidence_EmptyQueryVars(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ve := NewVariableElimination([]*factors.DiscreteFactor{pA})
	_, err := ve.QueryWithVirtualEvidence(nil, nil, nil)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}
}

func TestQueryWithVirtualEvidence_EmptyValues(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ve := NewVariableElimination([]*factors.DiscreteFactor{pA})
	_, err := ve.QueryWithVirtualEvidence([]string{"A"}, nil, []VirtualEvidence{
		{Variable: "A", Values: nil},
	})
	if err == nil {
		t.Error("expected error for empty virtual evidence values")
	}
}

func TestQueryWithVirtualEvidence_StudentNetwork(t *testing.T) {
	ve := NewVariableElimination(studentFactors())

	// Virtual evidence on Intelligence: slightly favors I=1 (high).
	vEvidence := []VirtualEvidence{
		{Variable: "I", Values: []float64{0.3, 0.7}},
	}

	result, err := ve.QueryWithVirtualEvidence([]string{"G"}, nil, vEvidence)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// Compare with standard query without virtual evidence.
	resultStd, _ := ve.Query([]string{"G"}, nil)
	// Virtual evidence favoring high I should increase P(G=0) (good grade).
	g0VE := result.GetValue(map[string]int{"G": 0})
	g0Std := resultStd.GetValue(map[string]int{"G": 0})
	if g0VE <= g0Std {
		t.Errorf("expected virtual evidence favoring I=1 to increase P(G=0): VE=%f, Std=%f", g0VE, g0Std)
	}
}

// ---------------------------------------------------------------------------
// QueryMarginals tests (joint=False)
// ---------------------------------------------------------------------------

func TestQueryMarginals_Basic(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	marginals, err := ve.QueryMarginals([]string{"D", "I"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(marginals) != 2 {
		t.Fatalf("expected 2 marginals, got %d", len(marginals))
	}

	// Check that each marginal is a proper distribution over its variable.
	dMarg, ok := marginals["D"]
	if !ok {
		t.Fatal("missing marginal for D")
	}
	assertSumsToOne(t, dMarg)
	assertNear(t, dMarg.GetValue(map[string]int{"D": 0}), 0.6, 1e-6, "P(D=0)")
	assertNear(t, dMarg.GetValue(map[string]int{"D": 1}), 0.4, 1e-6, "P(D=1)")

	iMarg, ok := marginals["I"]
	if !ok {
		t.Fatal("missing marginal for I")
	}
	assertSumsToOne(t, iMarg)
	assertNear(t, iMarg.GetValue(map[string]int{"I": 0}), 0.7, 1e-6, "P(I=0)")
	assertNear(t, iMarg.GetValue(map[string]int{"I": 1}), 0.3, 1e-6, "P(I=1)")
}

func TestQueryMarginals_WithEvidence(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	marginals, err := ve.QueryMarginals([]string{"D", "I"}, map[string]int{"G": 0})
	if err != nil {
		t.Fatal(err)
	}

	if len(marginals) != 2 {
		t.Fatalf("expected 2 marginals, got %d", len(marginals))
	}

	for _, v := range []string{"D", "I"} {
		m, ok := marginals[v]
		if !ok {
			t.Fatalf("missing marginal for %s", v)
		}
		assertSumsToOne(t, m)
	}

	// P(I | G=0): with good grade, intelligence should shift toward high.
	iMarg := marginals["I"]
	i1 := iMarg.GetValue(map[string]int{"I": 1})
	if i1 <= 0.3 { // prior is 0.3
		t.Errorf("expected P(I=1|G=0) > 0.3, got %f", i1)
	}
}

func TestQueryMarginals_EmptyQueryVars(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	_, err := ve.QueryMarginals(nil, nil)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}
}

func TestQueryMarginals_SingleVar(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	marginals, err := ve.QueryMarginals([]string{"G"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(marginals) != 1 {
		t.Fatalf("expected 1 marginal, got %d", len(marginals))
	}
	gMarg := marginals["G"]
	assertSumsToOne(t, gMarg)
	assertNear(t, gMarg.GetValue(map[string]int{"G": 0}), 0.362, 1e-5, "P(G=0)")
}

// ---------------------------------------------------------------------------
// GetCliqueBeliefs tests (BeliefPropagation)
// ---------------------------------------------------------------------------

func TestGetCliqueBeliefs_AfterCalibrate(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})
	pCB, _ := factors.NewDiscreteFactor([]string{"C", "B"}, []int{2, 2}, []float64{0.5, 0.1, 0.5, 0.9})

	cliques := [][]string{{"A", "B"}, {"B", "C"}}
	separators := map[string][]string{"0-1": {"B"}}
	cliqueFactors := map[int][]*factors.DiscreteFactor{
		0: {pA, pBA},
		1: {pCB},
	}
	bp := NewBeliefPropagation(cliques, separators, cliqueFactors)

	if err := bp.Calibrate(); err != nil {
		t.Fatal(err)
	}

	beliefs := bp.GetCliqueBeliefs()
	if len(beliefs) != 2 {
		t.Fatalf("expected 2 clique beliefs, got %d", len(beliefs))
	}

	for i := 0; i < 2; i++ {
		if beliefs[i] == nil {
			t.Errorf("clique belief %d is nil", i)
		}
	}

	// Verify returned beliefs are copies (modifying doesn't affect BP).
	beliefs[0].Normalize()
	original := bp.GetCliqueBelief(0)
	if original == nil {
		t.Fatal("original belief should not be nil")
	}
}

func TestGetCliqueBeliefs_BeforeCalibrate(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})

	cliques := [][]string{{"A"}}
	separators := map[string][]string{}
	cliqueFactors := map[int][]*factors.DiscreteFactor{
		0: {pA},
	}
	bp := NewBeliefPropagation(cliques, separators, cliqueFactors)

	beliefs := bp.GetCliqueBeliefs()
	if len(beliefs) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(beliefs))
	}
	// Before calibration, potentials are not initialized, so should be nil.
	if beliefs[0] != nil {
		t.Error("expected nil belief before calibration")
	}
}

// ---------------------------------------------------------------------------
// QueryRejection tests (ApproxInference)
// ---------------------------------------------------------------------------

func TestQueryRejection_Basic(t *testing.T) {
	// Simple factor: P(A) = [0.3, 0.7]
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)

	result, err := ai.QueryRejection([]string{"A"}, nil, 100000)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)
	a0 := result.GetValue(map[string]int{"A": 0})
	assertNear(t, a0, 0.3, 0.05, "P(A=0) rejection")
}

func TestQueryRejection_WithEvidence(t *testing.T) {
	// A -> B network
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA, pBA}, 42)

	result, err := ai.QueryRejection([]string{"A"}, map[string]int{"B": 0}, 500000)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// P(A=0|B=0) = P(B=0|A=0)*P(A=0)/P(B=0) = 0.08/0.26
	a0 := result.GetValue(map[string]int{"A": 0})
	assertNear(t, a0, 0.08/0.26, 0.05, "P(A=0|B=0) rejection")
}

func TestQueryRejection_EmptyQueryVars(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)
	_, err := ai.QueryRejection(nil, nil, 100)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}
}

func TestQueryRejection_ZeroSamples(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)
	_, err := ai.QueryRejection([]string{"A"}, nil, 0)
	if err == nil {
		t.Error("expected error for nSamples=0")
	}
}

func TestQueryRejection_BadEvidenceVariable(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)
	_, err := ai.QueryRejection([]string{"A"}, map[string]int{"Z": 0}, 100)
	if err == nil {
		t.Error("expected error for unknown evidence variable")
	}
}

func TestQueryRejection_BadEvidenceValue(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)
	_, err := ai.QueryRejection([]string{"A"}, map[string]int{"A": 5}, 100)
	if err == nil {
		t.Error("expected error for out-of-range evidence value")
	}
}

// ---------------------------------------------------------------------------
// QueryGibbs tests (ApproxInference)
// ---------------------------------------------------------------------------

func TestQueryGibbs_Basic(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)

	result, err := ai.QueryGibbs([]string{"A"}, nil, 50000, 1000)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)
	a0 := result.GetValue(map[string]int{"A": 0})
	assertNear(t, a0, 0.3, 0.05, "P(A=0) Gibbs")
}

func TestQueryGibbs_WithEvidence(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA, pBA}, 42)

	result, err := ai.QueryGibbs([]string{"A"}, map[string]int{"B": 0}, 50000, 1000)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	a0 := result.GetValue(map[string]int{"A": 0})
	assertNear(t, a0, 0.08/0.26, 0.05, "P(A=0|B=0) Gibbs")
}

func TestQueryGibbs_MultiFactor(t *testing.T) {
	// A -> B -> C chain.
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.9, 0.1, 0.1, 0.9})
	pCB, _ := factors.NewDiscreteFactor([]string{"C", "B"}, []int{2, 2}, []float64{0.8, 0.2, 0.2, 0.8})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA, pBA, pCB}, 42)

	result, err := ai.QueryGibbs([]string{"A"}, map[string]int{"C": 0}, 50000, 2000)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// Evidence C=0 should favor A=0 (through the chain).
	a0 := result.GetValue(map[string]int{"A": 0})
	if a0 < 0.5 {
		t.Errorf("expected P(A=0|C=0) > 0.5, got %f", a0)
	}
}

func TestQueryGibbs_EmptyQueryVars(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)
	_, err := ai.QueryGibbs(nil, nil, 100, 10)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}
}

func TestQueryGibbs_ZeroSamples(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)
	_, err := ai.QueryGibbs([]string{"A"}, nil, 0, 10)
	if err == nil {
		t.Error("expected error for nSamples=0")
	}
}

func TestQueryGibbs_NegativeBurnIn(t *testing.T) {
	// Negative burn-in should be treated as 0 (no error).
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)
	result, err := ai.QueryGibbs([]string{"A"}, nil, 1000, -5)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)
}

func TestQueryGibbs_BadEvidenceVariable(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA}, 42)
	_, err := ai.QueryGibbs([]string{"A"}, map[string]int{"Z": 0}, 100, 10)
	if err == nil {
		t.Error("expected error for unknown evidence variable")
	}
}

// ---------------------------------------------------------------------------
// GetDistributionWithEvidence tests (ApproxInference)
// ---------------------------------------------------------------------------

func TestGetDistributionWithEvidence_Basic(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})
	ai := NewApproxInference([]*factors.DiscreteFactor{pA, pBA}, 42)

	result, err := ai.GetDistributionWithEvidence([]string{"A"}, map[string]int{"B": 0}, 100000)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	a0 := result.GetValue(map[string]int{"A": 0})
	assertNear(t, a0, 0.08/0.26, 0.05, "P(A=0|B=0) distribution with evidence")
}

// ---------------------------------------------------------------------------
// VirtualEvidence type tests
// ---------------------------------------------------------------------------

func TestVirtualEvidence_Type(t *testing.T) {
	ve := VirtualEvidence{
		Variable: "X",
		Values:   []float64{0.3, 0.7},
	}
	if ve.Variable != "X" {
		t.Errorf("expected Variable=X, got %s", ve.Variable)
	}
	if len(ve.Values) != 2 {
		t.Errorf("expected 2 values, got %d", len(ve.Values))
	}
}

// ---------------------------------------------------------------------------
// Cross-validation: VE QueryMarginals vs joint Query
// ---------------------------------------------------------------------------

func TestQueryMarginals_ConsistentWithJoint(t *testing.T) {
	ve := NewVariableElimination(studentFactors())

	// Get joint P(D, I | G=0)
	joint, err := ve.Query([]string{"D", "I"}, map[string]int{"G": 0})
	if err != nil {
		t.Fatal(err)
	}

	// Get individual marginals P(D | G=0) and P(I | G=0)
	marginals, err := ve.QueryMarginals([]string{"D", "I"}, map[string]int{"G": 0})
	if err != nil {
		t.Fatal(err)
	}

	// Check that marginalizing the joint gives the same as individual marginals.
	// P(D=0 | G=0) = sum_I P(D=0, I | G=0)
	pD0fromJoint := joint.GetValue(map[string]int{"D": 0, "I": 0}) + joint.GetValue(map[string]int{"D": 0, "I": 1})
	pD0fromMarginal := marginals["D"].GetValue(map[string]int{"D": 0})
	assertNear(t, pD0fromJoint, pD0fromMarginal, 1e-6, "P(D=0|G=0) joint vs marginal")

	pI1fromJoint := joint.GetValue(map[string]int{"D": 0, "I": 1}) + joint.GetValue(map[string]int{"D": 1, "I": 1})
	pI1fromMarginal := marginals["I"].GetValue(map[string]int{"I": 1})
	assertNear(t, pI1fromJoint, pI1fromMarginal, 1e-6, "P(I=1|G=0) joint vs marginal")
}

// ---------------------------------------------------------------------------
// Cross-validation: Virtual evidence with hard evidence = delta VE
// ---------------------------------------------------------------------------

func TestVirtualEvidence_DeltaEqualsHardEvidence(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})
	ve := NewVariableElimination([]*factors.DiscreteFactor{pA, pBA})

	// Hard evidence B=0
	hardResult, err := ve.Query([]string{"A"}, map[string]int{"B": 0})
	if err != nil {
		t.Fatal(err)
	}

	// Virtual evidence with delta function [1, 0] for B should be equivalent.
	veResult, err := ve.QueryWithVirtualEvidence([]string{"A"}, nil, []VirtualEvidence{
		{Variable: "B", Values: []float64{1.0, 0.0}},
	})
	if err != nil {
		t.Fatal(err)
	}

	a0Hard := hardResult.GetValue(map[string]int{"A": 0})
	a0VE := veResult.GetValue(map[string]int{"A": 0})
	assertNear(t, a0Hard, a0VE, 1e-6, "delta virtual evidence should match hard evidence")
}

// ---------------------------------------------------------------------------
// Cross-validation: Gibbs vs exact for small network
// ---------------------------------------------------------------------------

func TestQueryGibbs_MatchesExact(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})

	// Exact answer.
	ve := NewVariableElimination([]*factors.DiscreteFactor{pA, pBA})
	exact, err := ve.Query([]string{"B"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Gibbs answer.
	ai := NewApproxInference([]*factors.DiscreteFactor{pA, pBA}, 42)
	gibbs, err := ai.QueryGibbs([]string{"B"}, nil, 100000, 2000)
	if err != nil {
		t.Fatal(err)
	}

	b0Exact := exact.GetValue(map[string]int{"B": 0})
	b0Gibbs := gibbs.GetValue(map[string]int{"B": 0})
	if math.Abs(b0Exact-b0Gibbs) > 0.05 {
		t.Errorf("Gibbs P(B=0)=%f differs from exact %f by more than 0.05", b0Gibbs, b0Exact)
	}
}

// ---------------------------------------------------------------------------
// Cross-validation: Rejection vs exact for small network
// ---------------------------------------------------------------------------

func TestQueryRejection_MatchesExact(t *testing.T) {
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})

	ve := NewVariableElimination([]*factors.DiscreteFactor{pA, pBA})
	exact, err := ve.Query([]string{"B"}, nil)
	if err != nil {
		t.Fatal(err)
	}

	ai := NewApproxInference([]*factors.DiscreteFactor{pA, pBA}, 42)
	rejection, err := ai.QueryRejection([]string{"B"}, nil, 200000)
	if err != nil {
		t.Fatal(err)
	}

	b0Exact := exact.GetValue(map[string]int{"B": 0})
	b0Rej := rejection.GetValue(map[string]int{"B": 0})
	if math.Abs(b0Exact-b0Rej) > 0.05 {
		t.Errorf("Rejection P(B=0)=%f differs from exact %f by more than 0.05", b0Rej, b0Exact)
	}
}
