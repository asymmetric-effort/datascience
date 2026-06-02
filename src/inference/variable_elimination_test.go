//go:build unit

package inference

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

const epsilon = 1e-9

func floatEq(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

// floatNear allows a looser tolerance for accumulated floating-point error.
func floatNear(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

// ---------------------------------------------------------------------------
// Student network helpers
// ---------------------------------------------------------------------------
//
// Classic student Bayesian network (Koller & Friedman):
//
//   Difficulty(D) -> Grade(G) <- Intelligence(I)
//                    Grade(G) -> Letter(L)
//                    Intelligence(I) -> SAT(S)
//
// Variable cardinalities:
//   D: 2 (d0=easy, d1=hard)
//   I: 2 (i0=low, i1=high)
//   G: 3 (g1, g2, g3)
//   L: 2 (l0=weak, l1=strong)
//   S: 2 (s0=low, s1=high)

func studentFactors() []*factors.DiscreteFactor {
	// P(D)
	pD, _ := factors.NewDiscreteFactor(
		[]string{"D"}, []int{2},
		[]float64{0.6, 0.4},
	)

	// P(I)
	pI, _ := factors.NewDiscreteFactor(
		[]string{"I"}, []int{2},
		[]float64{0.7, 0.3},
	)

	// P(G | D, I) — factor over (G, D, I), shape 3x2x2
	// Rows = G states, columns = (D,I) configs in row-major: (d0,i0),(d0,i1),(d1,i0),(d1,i1)
	// P(g1|d0,i0)=0.3, P(g1|d0,i1)=0.9, P(g1|d1,i0)=0.05, P(g1|d1,i1)=0.5
	// P(g2|d0,i0)=0.4, P(g2|d0,i1)=0.08,P(g2|d1,i0)=0.25, P(g2|d1,i1)=0.3
	// P(g3|d0,i0)=0.3, P(g3|d0,i1)=0.02,P(g3|d1,i0)=0.7,  P(g3|d1,i1)=0.2
	// Flat row-major over (G,D,I): G varies slowest.
	pGDI, _ := factors.NewDiscreteFactor(
		[]string{"G", "D", "I"}, []int{3, 2, 2},
		[]float64{
			// G=0 (g1): D=0,I=0; D=0,I=1; D=1,I=0; D=1,I=1
			0.3, 0.9, 0.05, 0.5,
			// G=1 (g2)
			0.4, 0.08, 0.25, 0.3,
			// G=2 (g3)
			0.3, 0.02, 0.7, 0.2,
		},
	)

	// P(L | G) — factor over (L, G), shape 2x3
	// P(l0|g1)=0.1, P(l0|g2)=0.4, P(l0|g3)=0.99
	// P(l1|g1)=0.9, P(l1|g2)=0.6, P(l1|g3)=0.01
	pLG, _ := factors.NewDiscreteFactor(
		[]string{"L", "G"}, []int{2, 3},
		[]float64{
			// L=0: G=0,G=1,G=2
			0.1, 0.4, 0.99,
			// L=1
			0.9, 0.6, 0.01,
		},
	)

	// P(S | I) — factor over (S, I), shape 2x2
	// P(s0|i0)=0.95, P(s0|i1)=0.2
	// P(s1|i0)=0.05, P(s1|i1)=0.8
	pSI, _ := factors.NewDiscreteFactor(
		[]string{"S", "I"}, []int{2, 2},
		[]float64{
			0.95, 0.2,
			0.05, 0.8,
		},
	)

	return []*factors.DiscreteFactor{pD, pI, pGDI, pLG, pSI}
}

// ---------------------------------------------------------------------------
// NewVariableElimination
// ---------------------------------------------------------------------------

func TestNewVariableElimination(t *testing.T) {
	fl := studentFactors()
	ve := NewVariableElimination(fl)
	if ve == nil {
		t.Fatal("NewVariableElimination returned nil")
	}
	if len(ve.factors) != 5 {
		t.Errorf("expected 5 factors, got %d", len(ve.factors))
	}
	// Verify deep copy: modifying original should not affect VE.
	fl[0].Normalize()
	origVal := ve.factors[0].GetValue(map[string]int{"D": 0})
	if !floatEq(origVal, 0.6) {
		t.Errorf("expected deep copy to preserve 0.6, got %f", origVal)
	}
}

// ---------------------------------------------------------------------------
// Query — marginal without evidence
// ---------------------------------------------------------------------------

func TestQuery_MarginalD(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	result, err := ve.Query([]string{"D"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	// P(D) should be [0.6, 0.4] since no evidence changes the prior.
	assertSumsToOne(t, result)
	assertNear(t, result.GetValue(map[string]int{"D": 0}), 0.6, 1e-6, "P(D=0)")
	assertNear(t, result.GetValue(map[string]int{"D": 1}), 0.4, 1e-6, "P(D=1)")
}

func TestQuery_MarginalI(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	result, err := ve.Query([]string{"I"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)
	assertNear(t, result.GetValue(map[string]int{"I": 0}), 0.7, 1e-6, "P(I=0)")
	assertNear(t, result.GetValue(map[string]int{"I": 1}), 0.3, 1e-6, "P(I=1)")
}

func TestQuery_MarginalG(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	result, err := ve.Query([]string{"G"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// P(G) = sum over D,I of P(G|D,I)*P(D)*P(I)
	// P(G=0) = 0.3*0.6*0.7 + 0.9*0.6*0.3 + 0.05*0.4*0.7 + 0.5*0.4*0.3
	//        = 0.126 + 0.162 + 0.014 + 0.06 = 0.362
	// P(G=1) = 0.4*0.6*0.7 + 0.08*0.6*0.3 + 0.25*0.4*0.7 + 0.3*0.4*0.3
	//        = 0.168 + 0.0144 + 0.07 + 0.036 = 0.2884
	// P(G=2) = 0.3*0.6*0.7 + 0.02*0.6*0.3 + 0.7*0.4*0.7 + 0.2*0.4*0.3
	//        = 0.126 + 0.0036 + 0.196 + 0.024 = 0.3496
	assertNear(t, result.GetValue(map[string]int{"G": 0}), 0.362, 1e-6, "P(G=0)")
	assertNear(t, result.GetValue(map[string]int{"G": 1}), 0.2884, 1e-6, "P(G=1)")
	assertNear(t, result.GetValue(map[string]int{"G": 2}), 0.3496, 1e-6, "P(G=2)")
}

func TestQuery_MarginalL(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	result, err := ve.Query([]string{"L"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// P(L=0) = sum_G P(L=0|G)*P(G)
	//        = 0.1*0.362 + 0.4*0.2884 + 0.99*0.3496
	//        = 0.0362 + 0.11536 + 0.346104 = 0.497664
	assertNear(t, result.GetValue(map[string]int{"L": 0}), 0.497664, 1e-5, "P(L=0)")
	assertNear(t, result.GetValue(map[string]int{"L": 1}), 1.0-0.497664, 1e-5, "P(L=1)")
}

func TestQuery_MarginalS(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	result, err := ve.Query([]string{"S"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// P(S=0) = 0.95*0.7 + 0.2*0.3 = 0.665 + 0.06 = 0.725
	assertNear(t, result.GetValue(map[string]int{"S": 0}), 0.725, 1e-6, "P(S=0)")
	assertNear(t, result.GetValue(map[string]int{"S": 1}), 0.275, 1e-6, "P(S=1)")
}

// ---------------------------------------------------------------------------
// Query — with evidence
// ---------------------------------------------------------------------------

func TestQuery_GradeGivenDifficulty(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	// P(G | D=1) — hard difficulty
	result, err := ve.Query([]string{"G"}, map[string]int{"D": 1})
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// P(G|D=1) = sum_I P(G|D=1,I)*P(I)
	// P(G=0|D=1) = 0.05*0.7 + 0.5*0.3 = 0.035 + 0.15 = 0.185
	// P(G=1|D=1) = 0.25*0.7 + 0.3*0.3 = 0.175 + 0.09 = 0.265
	// P(G=2|D=1) = 0.7*0.7  + 0.2*0.3 = 0.49  + 0.06 = 0.55
	assertNear(t, result.GetValue(map[string]int{"G": 0}), 0.185, 1e-6, "P(G=0|D=1)")
	assertNear(t, result.GetValue(map[string]int{"G": 1}), 0.265, 1e-6, "P(G=1|D=1)")
	assertNear(t, result.GetValue(map[string]int{"G": 2}), 0.55, 1e-6, "P(G=2|D=1)")
}

func TestQuery_IntelligenceGivenGrade(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	// P(I | G=0) — good grade, should shift toward high intelligence
	result, err := ve.Query([]string{"I"}, map[string]int{"G": 0})
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// P(I=0|G=0) = P(G=0|I=0)*P(I=0) / P(G=0)
	// P(G=0,I=0) = sum_D P(G=0|D,I=0)*P(D) = 0.3*0.6 + 0.05*0.4 = 0.18 + 0.02 = 0.2
	// P(G=0,I=1) = sum_D P(G=0|D,I=1)*P(D) = 0.9*0.6 + 0.5*0.4  = 0.54 + 0.2  = 0.74 (wait, wrong)
	// Actually P(G=0,I=1) is not normalized over I. Let me recompute:
	// P(G=0,I=0) = 0.2 * P(I=0) ... no.
	// Joint: P(G=0,I=0) = sum_D P(G=0|D,I=0)*P(D)*P(I=0)
	//                    = (0.3*0.6 + 0.05*0.4)*0.7 = 0.2*0.7 = 0.14
	// P(G=0,I=1) = (0.9*0.6 + 0.5*0.4)*0.3 = 0.74*0.3 = 0.222
	// P(G=0) = 0.14 + 0.222 = 0.362
	// P(I=0|G=0) = 0.14/0.362 ≈ 0.38674
	// P(I=1|G=0) = 0.222/0.362 ≈ 0.61326
	assertNear(t, result.GetValue(map[string]int{"I": 0}), 0.14/0.362, 1e-5, "P(I=0|G=0)")
	assertNear(t, result.GetValue(map[string]int{"I": 1}), 0.222/0.362, 1e-5, "P(I=1|G=0)")
}

func TestQuery_LetterGivenIntelligence(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	// P(L | I=1) — high intelligence
	result, err := ve.Query([]string{"L"}, map[string]int{"I": 1})
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)

	// P(G|I=1) = sum_D P(G|D,I=1)*P(D)
	// P(G=0|I=1) = 0.9*0.6 + 0.5*0.4 = 0.74
	// P(G=1|I=1) = 0.08*0.6 + 0.3*0.4 = 0.168
	// P(G=2|I=1) = 0.02*0.6 + 0.2*0.4 = 0.092
	// P(L=0|I=1) = 0.1*0.74 + 0.4*0.168 + 0.99*0.092
	//            = 0.074 + 0.0672 + 0.09108 = 0.23228
	assertNear(t, result.GetValue(map[string]int{"L": 0}), 0.23228, 1e-5, "P(L=0|I=1)")
	assertNear(t, result.GetValue(map[string]int{"L": 1}), 1.0-0.23228, 1e-5, "P(L=1|I=1)")
}

func TestQuery_MultipleEvidence(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	// P(G | D=0, I=1) — should just be the CPD column
	result, err := ve.Query([]string{"G"}, map[string]int{"D": 0, "I": 1})
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)
	assertNear(t, result.GetValue(map[string]int{"G": 0}), 0.9, 1e-6, "P(G=0|D=0,I=1)")
	assertNear(t, result.GetValue(map[string]int{"G": 1}), 0.08, 1e-6, "P(G=1|D=0,I=1)")
	assertNear(t, result.GetValue(map[string]int{"G": 2}), 0.02, 1e-6, "P(G=2|D=0,I=1)")
}

// ---------------------------------------------------------------------------
// Query — joint distributions
// ---------------------------------------------------------------------------

func TestQuery_JointDI(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	// P(D, I) — should be product of priors since they're independent.
	result, err := ve.Query([]string{"D", "I"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)
	assertNear(t, result.GetValue(map[string]int{"D": 0, "I": 0}), 0.42, 1e-6, "P(D=0,I=0)")
	assertNear(t, result.GetValue(map[string]int{"D": 0, "I": 1}), 0.18, 1e-6, "P(D=0,I=1)")
	assertNear(t, result.GetValue(map[string]int{"D": 1, "I": 0}), 0.28, 1e-6, "P(D=1,I=0)")
	assertNear(t, result.GetValue(map[string]int{"D": 1, "I": 1}), 0.12, 1e-6, "P(D=1,I=1)")
}

// ---------------------------------------------------------------------------
// Query — error cases
// ---------------------------------------------------------------------------

func TestQuery_EmptyQueryVars(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	_, err := ve.Query(nil, nil)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}
}

// ---------------------------------------------------------------------------
// MAP
// ---------------------------------------------------------------------------

func TestMAP_NoEvidence(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	assignment, err := ve.MAP([]string{"D"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	// P(D=0)=0.6 > P(D=1)=0.4, so MAP should be D=0.
	if assignment["D"] != 0 {
		t.Errorf("expected MAP D=0, got D=%d", assignment["D"])
	}
}

func TestMAP_GradeGivenHardAndLowIntelligence(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	// P(G | D=1, I=0) = [0.05, 0.25, 0.7] -> MAP is G=2
	assignment, err := ve.MAP([]string{"G"}, map[string]int{"D": 1, "I": 0})
	if err != nil {
		t.Fatal(err)
	}
	if assignment["G"] != 2 {
		t.Errorf("expected MAP G=2, got G=%d", assignment["G"])
	}
}

func TestMAP_GradeGivenEasyAndHighIntelligence(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	// P(G | D=0, I=1) = [0.9, 0.08, 0.02] -> MAP is G=0
	assignment, err := ve.MAP([]string{"G"}, map[string]int{"D": 0, "I": 1})
	if err != nil {
		t.Fatal(err)
	}
	if assignment["G"] != 0 {
		t.Errorf("expected MAP G=0, got G=%d", assignment["G"])
	}
}

func TestMAP_MultipleVars(t *testing.T) {
	ve := NewVariableElimination(studentFactors())
	// MAP(D, I) with no evidence: should be D=0, I=0 (both priors favor 0)
	assignment, err := ve.MAP([]string{"D", "I"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if assignment["D"] != 0 || assignment["I"] != 0 {
		t.Errorf("expected MAP D=0,I=0, got D=%d,I=%d", assignment["D"], assignment["I"])
	}
}

// ---------------------------------------------------------------------------
// Simple two-variable network
// ---------------------------------------------------------------------------

func TestQuery_SimpleNetwork(t *testing.T) {
	// A -> B
	// P(A) = [0.4, 0.6]
	// P(B|A): P(B=0|A=0)=0.2, P(B=0|A=1)=0.3, P(B=1|A=0)=0.8, P(B=1|A=1)=0.7
	pA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.4, 0.6})
	pBA, _ := factors.NewDiscreteFactor([]string{"B", "A"}, []int{2, 2}, []float64{0.2, 0.3, 0.8, 0.7})

	ve := NewVariableElimination([]*factors.DiscreteFactor{pA, pBA})

	// P(B) = sum_A P(B|A)*P(A)
	// P(B=0) = 0.2*0.4 + 0.3*0.6 = 0.08 + 0.18 = 0.26
	// P(B=1) = 0.8*0.4 + 0.7*0.6 = 0.32 + 0.42 = 0.74
	result, err := ve.Query([]string{"B"}, nil)
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result)
	assertNear(t, result.GetValue(map[string]int{"B": 0}), 0.26, 1e-6, "P(B=0)")
	assertNear(t, result.GetValue(map[string]int{"B": 1}), 0.74, 1e-6, "P(B=1)")

	// P(A | B=0)
	// P(A=0|B=0) = P(B=0|A=0)*P(A=0)/P(B=0) = 0.08/0.26
	result2, err := ve.Query([]string{"A"}, map[string]int{"B": 0})
	if err != nil {
		t.Fatal(err)
	}
	assertSumsToOne(t, result2)
	assertNear(t, result2.GetValue(map[string]int{"A": 0}), 0.08/0.26, 1e-6, "P(A=0|B=0)")
	assertNear(t, result2.GetValue(map[string]int{"A": 1}), 0.18/0.26, 1e-6, "P(A=1|B=0)")
}

// ---------------------------------------------------------------------------
// Elimination order tests
// ---------------------------------------------------------------------------

func TestMinNeighborsOrder_Empty(t *testing.T) {
	order := MinNeighborsOrder(nil, nil)
	if len(order) != 0 {
		t.Errorf("expected empty order, got %v", order)
	}
}

func TestMinNeighborsOrder_SingleVar(t *testing.T) {
	f, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	order := MinNeighborsOrder([]*factors.DiscreteFactor{f}, []string{"A"})
	if len(order) != 1 || order[0] != "A" {
		t.Errorf("expected [A], got %v", order)
	}
}

func TestMinNeighborsOrder_PrefersFewer(t *testing.T) {
	// f1(A,B), f2(A,C), f3(B,D) — A appears in 2, B appears in 2, D appears in 1
	// Eliminating {D, A}: D should come first (appears in fewer factors).
	f1, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	f2, _ := factors.NewDiscreteFactor([]string{"A", "C"}, []int{2, 2}, []float64{1, 2, 3, 4})
	f3, _ := factors.NewDiscreteFactor([]string{"B", "D"}, []int{2, 2}, []float64{1, 2, 3, 4})

	order := MinNeighborsOrder([]*factors.DiscreteFactor{f1, f2, f3}, []string{"D", "A"})
	if len(order) != 2 {
		t.Fatalf("expected 2 vars, got %v", order)
	}
	if order[0] != "D" {
		t.Errorf("expected D first (fewest factors), got %v", order)
	}
}

func TestMinNeighborsOrder_AllStudentVars(t *testing.T) {
	fl := studentFactors()
	// Eliminate all non-query vars for a query on G.
	order := MinNeighborsOrder(fl, []string{"D", "I", "L", "S"})
	if len(order) != 4 {
		t.Fatalf("expected 4 vars in order, got %d", len(order))
	}
	// Verify all vars are present.
	seen := make(map[string]bool)
	for _, v := range order {
		seen[v] = true
	}
	for _, v := range []string{"D", "I", "L", "S"} {
		if !seen[v] {
			t.Errorf("missing %q in elimination order", v)
		}
	}
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

func assertSumsToOne(t *testing.T, f *factors.DiscreteFactor) {
	t.Helper()
	data := f.Values().Data()
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	if !floatNear(sum, 1.0, 1e-6) {
		t.Errorf("factor values sum to %f, want 1.0", sum)
	}
}

func assertNear(t *testing.T, got, want, tol float64, label string) {
	t.Helper()
	if !floatNear(got, want, tol) {
		t.Errorf("%s = %f, want %f (tol=%e)", label, got, want, tol)
	}
}
