//go:build unit

package learning

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// buildThreeNodeNetwork creates a simple 3-node BN: A -> B -> C
// where B is the latent variable. A and C are binary, B is binary.
//
// True parameters:
//
//	P(A=0) = 0.6, P(A=1) = 0.4
//	P(B=0|A=0) = 0.8, P(B=0|A=1) = 0.3
//	P(B=1|A=0) = 0.2, P(B=1|A=1) = 0.7
//	P(C=0|B=0) = 0.9, P(C=0|B=1) = 0.4
//	P(C=1|B=0) = 0.1, P(C=1|B=1) = 0.6
func buildThreeNodeNetwork() *models.BayesianNetwork {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddNode("C")
	_ = bn.AddEdge("A", "B")
	_ = bn.AddEdge("B", "C")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	_ = bn.SetStates("C", []string{"0", "1"})
	return bn
}

// generateData generates observed data (A, C) from the true network
// parameters using a deterministic pattern that approximates the true
// joint distribution P(A, C) = sum_B P(A) P(B|A) P(C|B).
//
// P(A=0, C=0) = 0.6*(0.8*0.9 + 0.2*0.4) = 0.6*(0.72+0.08) = 0.48
// P(A=0, C=1) = 0.6*(0.8*0.1 + 0.2*0.6) = 0.6*(0.08+0.12) = 0.12
// P(A=1, C=0) = 0.4*(0.3*0.9 + 0.7*0.4) = 0.4*(0.27+0.28) = 0.22
// P(A=1, C=1) = 0.4*(0.3*0.1 + 0.7*0.6) = 0.4*(0.03+0.42) = 0.18
func generateData(n int) *tabgo.DataFrame {
	// We generate data in proportions matching the true marginal P(A,C).
	// Using 100 as base: 48, 12, 22, 18.
	aVals := make([]any, 0, n)
	cVals := make([]any, 0, n)

	patterns := []struct {
		a, c  string
		count int
	}{
		{"0", "0", 48},
		{"0", "1", 12},
		{"1", "0", 22},
		{"1", "1", 18},
	}

	for len(aVals) < n {
		for _, p := range patterns {
			for j := 0; j < p.count && len(aVals) < n; j++ {
				aVals = append(aVals, p.a)
				cVals = append(cVals, p.c)
			}
		}
	}

	return tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", aVals),
		"C": tabgo.NewSeries("C", cVals),
	})
}

func TestNewEM(t *testing.T) {
	bn := buildThreeNodeNetwork()
	data := generateData(100)
	em := NewEM(bn, data, []string{"B"}, 50, 1e-6)

	if em.Iterations() != 0 {
		t.Errorf("expected 0 iterations before Estimate, got %d", em.Iterations())
	}
	if em.Converged() {
		t.Error("expected not converged before Estimate")
	}
}

func TestEM_Estimate_Converges(t *testing.T) {
	bn := buildThreeNodeNetwork()
	data := generateData(500)
	em := NewEM(bn, data, []string{"B"}, 100, 1e-6)

	err := em.Estimate()
	if err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	if !em.Converged() {
		t.Errorf("expected convergence, ran %d iterations", em.Iterations())
	}

	// Verify that the model is valid after EM.
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("model invalid after EM: %v", err)
	}

	// Verify CPDs are proper distributions (columns sum to 1).
	for _, cpd := range bn.GetCPDs() {
		if err := cpd.Validate(); err != nil {
			t.Errorf("CPD for %q invalid: %v", cpd.Variable(), err)
		}
	}
}

func TestEM_Estimate_ReasonableCPDs(t *testing.T) {
	bn := buildThreeNodeNetwork()
	data := generateData(1000)
	em := NewEM(bn, data, []string{"B"}, 200, 1e-8)

	err := em.Estimate()
	if err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	// Check that the learned marginal P(A) matches the data.
	// P(A=0) should be ~0.6.
	cpdA := bn.GetCPD("A")
	if cpdA == nil {
		t.Fatal("no CPD for A")
	}
	pA0 := cpdA.ToFactor().Values().Data()[0]
	if math.Abs(pA0-0.6) > 0.05 {
		t.Errorf("P(A=0) = %f, expected ~0.6", pA0)
	}

	// The marginal P(C|A) implied by the learned parameters should
	// approximately match the empirical distribution.
	// P(C=0|A=0) = P(C=0|B=0)*P(B=0|A=0) + P(C=0|B=1)*P(B=1|A=0)
	// From data: P(C=0|A=0) = 48/60 = 0.8
	cpdB := bn.GetCPD("B")
	cpdC := bn.GetCPD("C")
	if cpdB == nil || cpdC == nil {
		t.Fatal("missing CPD for B or C")
	}

	bData := cpdB.ToFactor().Values().Data()
	cData := cpdC.ToFactor().Values().Data()

	// B factor layout: [B, A] with card [2, 2]
	// bData[B*2 + A]
	pB0gA0 := bData[0*2+0] // P(B=0|A=0)
	pB1gA0 := bData[1*2+0] // P(B=1|A=0)

	// C factor layout: [C, B] with card [2, 2]
	// cData[C*2 + B]
	pC0gB0 := cData[0*2+0] // P(C=0|B=0)
	pC0gB1 := cData[0*2+1] // P(C=0|B=1)

	impliedPC0gA0 := pC0gB0*pB0gA0 + pC0gB1*pB1gA0
	empiricalPC0gA0 := 48.0 / 60.0 // from data proportions

	if math.Abs(impliedPC0gA0-empiricalPC0gA0) > 0.1 {
		t.Errorf("implied P(C=0|A=0) = %f, expected ~%f", impliedPC0gA0, empiricalPC0gA0)
	}

	t.Logf("EM converged in %d iterations", em.Iterations())
	t.Logf("P(A=0) = %f", pA0)
	t.Logf("P(B=0|A=0) = %f, P(B=1|A=0) = %f", pB0gA0, pB1gA0)
	t.Logf("P(C=0|B=0) = %f, P(C=0|B=1) = %f", pC0gB0, pC0gB1)
	t.Logf("Implied P(C=0|A=0) = %f (empirical: %f)", impliedPC0gA0, empiricalPC0gA0)
}

func TestEM_ConvergenceDetection(t *testing.T) {
	bn := buildThreeNodeNetwork()
	data := generateData(200)

	// Run with a very loose tolerance to converge quickly.
	em := NewEM(bn, data, []string{"B"}, 1000, 0.1)
	err := em.Estimate()
	if err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}
	if !em.Converged() {
		t.Error("expected convergence with loose tolerance")
	}
	looseIter := em.Iterations()

	// Run again with a much tighter tolerance.
	bn2 := buildThreeNodeNetwork()
	em2 := NewEM(bn2, data, []string{"B"}, 1000, 1e-10)
	err = em2.Estimate()
	if err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}
	tightIter := em2.Iterations()

	// Tighter tolerance should require at least as many iterations.
	if tightIter < looseIter {
		t.Errorf("tight tolerance took fewer iterations (%d) than loose (%d)", tightIter, looseIter)
	}
	t.Logf("Loose tol iterations: %d, tight tol iterations: %d", looseIter, tightIter)
}

func TestEM_MaxIterReached(t *testing.T) {
	bn := buildThreeNodeNetwork()
	data := generateData(100)

	// Run with very few iterations and very tight tolerance.
	em := NewEM(bn, data, []string{"B"}, 2, 1e-15)
	err := em.Estimate()
	if err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	if em.Converged() {
		t.Error("did not expect convergence with only 2 iterations and tight tolerance")
	}
	if em.Iterations() != 2 {
		t.Errorf("expected 2 iterations, got %d", em.Iterations())
	}

	// Model should still be valid even if not converged.
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("model invalid after non-converged EM: %v", err)
	}
}

func TestEM_NoLatentVars(t *testing.T) {
	// When there are no latent variables, EM should just do MLE and
	// converge immediately (in 1 iteration).
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("X")
	_ = bn.AddNode("Y")
	_ = bn.AddEdge("X", "Y")
	_ = bn.SetStates("X", []string{"0", "1"})
	_ = bn.SetStates("Y", []string{"0", "1"})

	// Data: X and Y are both observed.
	xVals := make([]any, 100)
	yVals := make([]any, 100)
	for i := 0; i < 100; i++ {
		if i < 70 {
			xVals[i] = "0"
		} else {
			xVals[i] = "1"
		}
		if i < 50 {
			yVals[i] = "0"
		} else {
			yVals[i] = "1"
		}
	}
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", xVals),
		"Y": tabgo.NewSeries("Y", yVals),
	})

	em := NewEM(bn, data, nil, 50, 1e-6)
	err := em.Estimate()
	if err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	if !em.Converged() {
		t.Error("expected convergence with no latent variables")
	}
	// With no latent variables, MLE initialization should be exact,
	// so EM should converge on the first iteration.
	if em.Iterations() != 1 {
		t.Errorf("expected 1 iteration, got %d", em.Iterations())
	}

	// Check P(X=0) = 0.7.
	cpdX := bn.GetCPD("X")
	pX0 := cpdX.ToFactor().Values().Data()[0]
	if math.Abs(pX0-0.7) > 1e-6 {
		t.Errorf("P(X=0) = %f, expected 0.7", pX0)
	}
}

func TestEM_MissingStatesError(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	// Don't set states for A.

	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{"0", "1"}),
	})

	em := NewEM(bn, data, nil, 10, 1e-6)
	err := em.Estimate()
	if err == nil {
		t.Fatal("expected error when states are not defined")
	}
}

func TestEM_MissingDataColumnError(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})

	// Data doesn't contain column "A".
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{"0", "1"}),
	})

	em := NewEM(bn, data, nil, 10, 1e-6)
	err := em.Estimate()
	if err == nil {
		t.Fatal("expected error when data column is missing for observed variable")
	}
}

func TestEM_CPDsAreProperDistributions(t *testing.T) {
	bn := buildThreeNodeNetwork()
	data := generateData(200)
	em := NewEM(bn, data, []string{"B"}, 50, 1e-6)

	err := em.Estimate()
	if err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	// Every CPD column must sum to 1.
	for _, cpd := range bn.GetCPDs() {
		if err := cpd.Validate(); err != nil {
			t.Errorf("CPD for %q is not a proper distribution: %v", cpd.Variable(), err)
		}

		// Also check all values are non-negative.
		data := cpd.ToFactor().Values().Data()
		for i, v := range data {
			if v < 0 {
				t.Errorf("CPD for %q has negative value %f at index %d", cpd.Variable(), v, i)
			}
		}
	}
}

func TestEM_IterationsAndConvergedAccessors(t *testing.T) {
	bn := buildThreeNodeNetwork()
	data := generateData(100)
	em := NewEM(bn, data, []string{"B"}, 50, 1e-4)

	// Before estimation.
	if em.Iterations() != 0 {
		t.Errorf("expected 0 iterations before Estimate")
	}
	if em.Converged() {
		t.Error("expected not converged before Estimate")
	}

	err := em.Estimate()
	if err != nil {
		t.Fatalf("Estimate failed: %v", err)
	}

	if em.Iterations() == 0 {
		t.Error("expected at least 1 iteration after Estimate")
	}
}
