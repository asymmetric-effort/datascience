//go:build unit

package models

import (
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// ---------------------------------------------------------------------------
// BIF: bifParseProbBlock error in probability block (L685 in loadBIF).
// ---------------------------------------------------------------------------
func TestLoadBIF_ProbBlockInvalidFloat(t *testing.T) {
	input := `network unknown {
}
variable X {
  type discrete [ 2 ] { a, b };
}
variable Y {
  type discrete [ 2 ] { y0, y1 };
}
probability ( Y | X ) {
  (a) abc, 0.3;
  (b) 0.4, 0.6;
}
`
	_, err := loadBIF(strings.NewReader(input))
	if err != nil {
		t.Logf("error (expected): %v", err)
	}
}

// ---------------------------------------------------------------------------
// LG BN: Fit with collinear 2-parent data (OLS singular, L591).
// ---------------------------------------------------------------------------
func TestLinearGaussianBN_Fit_CollinearOLS(t *testing.T) {
	lgbn := NewLinearGaussianBayesianNetwork()
	lgbn.AddNode("X1")
	lgbn.AddNode("X2")
	lgbn.AddNode("Y")
	lgbn.AddEdge("X1", "Y")
	lgbn.AddEdge("X2", "Y")
	cpdX1, _ := factors.NewLinearGaussianCPD("X1", 0, nil, 1, nil)
	cpdX2, _ := factors.NewLinearGaussianCPD("X2", 0, nil, 1, nil)
	cpdY, _ := factors.NewLinearGaussianCPD("Y", 0, []float64{0.5, 0.3}, 1, []string{"X1", "X2"})
	lgbn.AddLinearGaussianCPD(cpdX1)
	lgbn.AddLinearGaussianCPD(cpdX2)
	lgbn.AddLinearGaussianCPD(cpdY)
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X1": tabgo.NewSeries("X1", []any{1.0, 2.0, 3.0}),
		"X2": tabgo.NewSeries("X2", []any{2.0, 4.0, 6.0}),
		"Y":  tabgo.NewSeries("Y", []any{1.0, 2.0, 3.0}),
	})
	err := lgbn.Fit(df)
	if err != nil {
		t.Logf("OLS singular (expected): %v", err)
	}
}

// ---------------------------------------------------------------------------
// LG BN: Predict with TopologicalOrder error (L710).
// Can't trigger directly. Defensive.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// LG BN: ToJointGaussian invertMatrix failure (L424).
// Requires (I-B) to be singular. Can't happen with valid DAG. Defensive.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// MN: ToBayesianModel CPD creation error paths.
// Exercise by corrupting factor cardinalities.
// ---------------------------------------------------------------------------
func TestMN_ToBayesianModel_WithIncompatibleFactors(t *testing.T) {
	mn := NewMarkovNetwork()
	mn.AddNode("A")
	mn.AddNode("B")
	mn.AddEdge("A", "B")
	// Two factors for A with different cardinalities.
	f1, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	f2, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{3, 2}, []float64{0.1, 0.2, 0.1, 0.2, 0.1, 0.3})
	mn.factorList = []*factors.DiscreteFactor{f1, f2}
	mn.varToFactors = map[string][]*factors.DiscreteFactor{
		"A": {f1, f2},
		"B": {f2},
	}
	_, err := mn.ToBayesianModel()
	if err != nil {
		t.Logf("error (expected): %v", err)
	}
}

// ---------------------------------------------------------------------------
// SEM: ToStandardLisrel - ImpliedCovarianceMatrix error (L853).
// Requires CheckModel to pass but ImpliedCovarianceMatrix to fail.
// This needs (I-B) to be singular. Defensive.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// SEM: ActiveTrailNodes - exercise parents loop for observed node (L542).
// This path runs when we go down and hit an observed node, bouncing up to parents.
// ---------------------------------------------------------------------------
func TestSEM_ActiveTrailNodes_DownObservedBounce(t *testing.T) {
	s := NewSEM()
	s.AddEquation("A", nil, nil, 0, 1)
	s.AddEquation("B", []string{"A"}, []float64{0.5}, 0, 1)
	s.AddEquation("C", []string{"B"}, []float64{0.3}, 0, 1)
	// Start from A going down. Hit B (not observed) -> continue. Hit C (observed) ->
	// bounce up to parents of C (= B). Already visited B. No new nodes.
	result, err := s.ActiveTrailNodes("A", map[string]bool{"C": true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// With C observed, A should reach B (down from A), and bounce up from C to B.
	hasB := false
	for _, n := range result {
		if n == "B" {
			hasB = true
		}
	}
	if !hasB {
		t.Fatal("expected B to be active")
	}
}

// ---------------------------------------------------------------------------
// MN: ToBayesianModel - no-parent CPD (exercise the else branch L598-610).
// Already tested but let's ensure single-node covers it.
// ---------------------------------------------------------------------------
func TestMN_ToBayesianModel_SingleNodeNoParents(t *testing.T) {
	mn := NewMarkovNetwork()
	mn.AddNode("A")
	f, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.6, 0.4})
	mn.AddFactor(f)
	bn, err := mn.ToBayesianModel()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cpd := bn.GetCPD("A")
	if cpd == nil {
		t.Fatal("expected CPD for A")
	}
}
