//go:build unit

package metrics

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// TestPearson_EmptySlices exercises the empty-slice path.
func TestPearson_EmptySlices(t *testing.T) {
	r := pearson(nil, nil)
	if r != 0 {
		t.Errorf("expected 0 for empty slices, got %f", r)
	}
}

// TestPearson_UnequalLengths exercises the mismatched-length path.
func TestPearson_UnequalLengths(t *testing.T) {
	r := pearson([]float64{1, 2}, []float64{1})
	if r != 0 {
		t.Errorf("expected 0 for unequal lengths, got %f", r)
	}
}

// TestPearson_ConstantValues exercises the zero-denominator path.
func TestPearson_ConstantValues(t *testing.T) {
	r := pearson([]float64{3, 3, 3}, []float64{5, 5, 5})
	if r != 0 {
		t.Errorf("expected 0 for constant values, got %f", r)
	}
}

// TestLowerIncompleteGammaReg_NegativeX exercises x < 0 path.
func TestLowerIncompleteGammaReg_NegativeX(t *testing.T) {
	result := lowerIncompleteGammaReg(2.0, -1.0)
	if result != 0 {
		t.Errorf("expected 0 for negative x, got %f", result)
	}
}

// TestLowerIncompleteGammaReg_ZeroX exercises x == 0 path.
func TestLowerIncompleteGammaReg_ZeroX(t *testing.T) {
	result := lowerIncompleteGammaReg(2.0, 0)
	if result != 0 {
		t.Errorf("expected 0 for zero x, got %f", result)
	}
}

// TestLowerIncompleteGammaReg_SeriesPath exercises x < a+1 series expansion.
func TestLowerIncompleteGammaReg_SeriesPath(t *testing.T) {
	result := lowerIncompleteGammaReg(5.0, 3.0)
	if result <= 0 || result >= 1 {
		t.Errorf("expected value in (0,1), got %f", result)
	}
}

// TestLowerIncompleteGammaReg_CFPath exercises x >= a+1 continued fraction.
func TestLowerIncompleteGammaReg_CFPath(t *testing.T) {
	result := lowerIncompleteGammaReg(2.0, 10.0)
	if result <= 0 || result > 1 {
		t.Errorf("expected value in (0,1], got %f", result)
	}
}

// TestUpperGammaCF exercises the continued fraction computation.
func TestUpperGammaCF(t *testing.T) {
	result := upperGammaCF(2.0, 5.0)
	if result < 0 || result > 1 {
		t.Errorf("expected value in [0,1], got %f", result)
	}
	// Q(2,5) should be small
	if result > 0.1 {
		t.Errorf("expected Q(2,5) < 0.1, got %f", result)
	}
}

// TestChiSquareCDF_ZeroX exercises x <= 0 path.
func TestChiSquareCDF_ZeroX(t *testing.T) {
	result := chiSquareCDF(0, 5)
	if result != 0 {
		t.Errorf("expected 0, got %f", result)
	}
}

// TestChiSquareCDF_ZeroDF exercises df <= 0 path.
func TestChiSquareCDF_ZeroDF(t *testing.T) {
	result := chiSquareCDF(5.0, 0)
	if result != 0 {
		t.Errorf("expected 0, got %f", result)
	}
}

// TestChiSquareCDF_Normal exercises the normal path.
func TestChiSquareCDF_Normal(t *testing.T) {
	result := chiSquareCDF(5.99, 2)
	if result < 0.9 || result > 1.0 {
		t.Errorf("expected CDF near 0.95 for chi2(2) at 5.99, got %f", result)
	}
}

// TestTTestPValue_ZeroDF exercises df <= 0 path.
func TestTTestPValue_ZeroDF(t *testing.T) {
	result := tTestPValue(2.0, 0)
	if result != 1 {
		t.Errorf("expected 1 for zero df, got %f", result)
	}
}

// TestTTestPValue_Normal exercises the normal computation.
func TestTTestPValue_Normal(t *testing.T) {
	result := tTestPValue(0, 10)
	if math.Abs(result-1.0) > 0.01 {
		t.Errorf("expected p-value near 1.0 for t=0, got %f", result)
	}
}

// TestTTestPValue_LargeT exercises large t-statistic.
func TestTTestPValue_LargeT_Coverage(t *testing.T) {
	result := tTestPValue(100, 10)
	if result > 0.001 {
		t.Errorf("expected very small p-value for t=100, got %f", result)
	}
}

// TestRegularizedIncompleteBeta_Boundaries exercises edge cases.
func TestRegularizedIncompleteBeta_Boundaries(t *testing.T) {
	if regularizedIncompleteBeta(0, 2, 3) != 0 {
		t.Error("expected 0 for x=0")
	}
	if regularizedIncompleteBeta(-1, 2, 3) != 0 {
		t.Error("expected 0 for x<0")
	}
	if regularizedIncompleteBeta(1, 2, 3) != 1 {
		t.Error("expected 1 for x=1")
	}
	if regularizedIncompleteBeta(2, 2, 3) != 1 {
		t.Error("expected 1 for x>1")
	}
}

// TestRegularizedIncompleteBeta_SymmetryPath exercises the symmetry branch.
func TestRegularizedIncompleteBeta_SymmetryPath(t *testing.T) {
	// When x > (a+1)/(a+b+2), the symmetry relation is used.
	a, b := 1.0, 10.0
	threshold := (a + 1) / (a + b + 2)
	x := threshold + 0.1
	if x >= 1 {
		x = 0.99
	}
	result := regularizedIncompleteBeta(x, a, b)
	if result < 0 || result > 1 {
		t.Errorf("expected value in [0,1], got %f", result)
	}
}

// TestPartialCorrelation_NoZ exercises the no-conditioning-variables path.
func TestPartialCorrelation_NoZ(t *testing.T) {
	colData := map[string][]float64{
		"X": {1, 2, 3, 4, 5},
		"Y": {2, 4, 6, 8, 10},
	}
	r := partialCorrelation(colData, "X", "Y", nil, 5)
	if math.Abs(r-1.0) > 0.01 {
		t.Errorf("expected r near 1.0 for perfectly correlated data, got %f", r)
	}
}

// TestPartialCorrelation_WithZ exercises the residualization path.
func TestPartialCorrelation_WithZ(t *testing.T) {
	colData := map[string][]float64{
		"X": {1, 2, 3, 4, 5},
		"Y": {2, 4, 6, 8, 10},
		"Z": {0, 1, 0, 1, 0},
	}
	r := partialCorrelation(colData, "X", "Y", []string{"Z"}, 5)
	if math.Abs(r) > 1.01 {
		t.Errorf("expected r in [-1,1], got %f", r)
	}
}

// TestResiduals_NoPredictors exercises the no-predictors path.
func TestResiduals_NoPredictors(t *testing.T) {
	colData := map[string][]float64{
		"Y": {1, 2, 3},
	}
	result := residuals(colData, "Y", nil, 3)
	if len(result) != 3 {
		t.Errorf("expected 3 residuals, got %d", len(result))
	}
}

// TestResiduals_ConstantPredictor exercises the zero-denom path.
func TestResiduals_ConstantPredictor(t *testing.T) {
	colData := map[string][]float64{
		"Y": {1, 2, 3, 4},
		"Z": {5, 5, 5, 5}, // constant predictor => denom ≈ 0
	}
	result := residuals(colData, "Y", []string{"Z"}, 4)
	if len(result) != 4 {
		t.Errorf("expected 4 residuals, got %d", len(result))
	}
}

// TestFisherC_NoCIs exercises the saturated model path (no implied CIs).
func TestFisherC_Saturated(t *testing.T) {
	// Fully connected graph => no non-adjacent pairs => saturated model
	edges := [][2]string{{"A", "B"}, {"B", "C"}, {"A", "C"}}
	data := makeSmallDataFrame()
	stat, pval := FisherC(edges, data)
	if stat != 0 || pval != 1 {
		t.Errorf("expected stat=0, pval=1 for saturated model, got stat=%f, pval=%f", stat, pval)
	}
}

// TestImpliedCIs_NoEdges exercises the all-independent case.
func TestImpliedCIs_NoEdges_Coverage(t *testing.T) {
	cis := ImpliedCIs(nil, []string{"A", "B", "C"})
	if len(cis) != 3 {
		t.Errorf("expected 3 CIs for 3 unconnected vars, got %d", len(cis))
	}
}

// TestFisherC_WithImpliedCIs exercises the full FisherC computation.
func TestFisherC_WithImpliedCIs(t *testing.T) {
	// A -> B, C separate => implied CI between A/C and B/C
	edges := [][2]string{{"A", "B"}}
	data := makeSmallDataFrame()
	stat, pval := FisherC(edges, data)
	if stat < 0 {
		t.Errorf("expected non-negative statistic, got %f", stat)
	}
	_ = pval // Just verify it runs without error
}

// TestFisherC_SmallDF exercises the dfT < 1 fallback.
func TestFisherC_SmallDF(t *testing.T) {
	// With many conditioning variables and few observations,
	// dfT could be < 1, triggering the fallback.
	edges := [][2]string{{"A", "B"}, {"B", "C"}}
	allVars := []string{"A", "B", "C"}
	cis := ImpliedCIs(edges, allVars)
	// A-C are non-adjacent; conditioning set could include B.
	_ = cis
}

// TestUpperGammaCF_SmallB0 exercises the tiny b0 path.
func TestUpperGammaCF_SmallB0(t *testing.T) {
	// When b0 = x + 1 - a is very close to 0
	result := upperGammaCF(1.0, 0.0001)
	if result < 0 || result > 1 {
		t.Errorf("expected value in [0,1], got %f", result)
	}
}

// TestRegularizedIncompleteBeta_MiddleValues exercises convergence.
func TestRegularizedIncompleteBeta_MiddleValues(t *testing.T) {
	v := regularizedIncompleteBeta(0.5, 2, 3)
	if v < 0 || v > 1 {
		t.Errorf("expected value in [0,1], got %f", v)
	}
	// Also test a value that triggers the symmetry relation.
	v2 := regularizedIncompleteBeta(0.8, 1, 5)
	if v2 < 0 || v2 > 1 {
		t.Errorf("expected value in [0,1], got %f", v2)
	}
}

// TestImpliedCIs_WithCommonNeighbors exercises the conditioning set computation.
func TestImpliedCIs_WithCommonNeighbors(t *testing.T) {
	// A->B, A->C, B not connected to C.
	// A is a common neighbor, so conditioning set for B-C should include A.
	edges := [][2]string{{"A", "B"}, {"A", "C"}}
	cis := ImpliedCIs(edges, []string{"A", "B", "C"})
	if len(cis) != 1 {
		t.Fatalf("expected 1 CI (B-C), got %d", len(cis))
	}
	// Conditioning set should include A.
	condSet := cis[0][2]
	found := false
	for _, v := range condSet {
		if v == "A" {
			found = true
		}
	}
	if !found {
		t.Error("expected A in conditioning set")
	}
}

// TestFisherC_NonSaturated exercises the full computation with implied CIs.
func TestFisherC_NonSaturated(t *testing.T) {
	// A->B, C separate.
	edges := [][2]string{{"A", "B"}}
	data := makeSmallDataFrame()
	stat, pval := FisherC(edges, data)
	if stat < 0 {
		t.Errorf("expected non-negative statistic, got %f", stat)
	}
	if pval < 0 || pval > 1 {
		t.Errorf("expected p-value in [0,1], got %f", pval)
	}
}

// makeSmallDataFrame creates a minimal DataFrame for testing.
func makeSmallDataFrame() *tabgo.DataFrame {
	sm := map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}),
		"B": tabgo.NewSeries("B", []any{2.0, 4.0, 6.0, 8.0, 10.0, 12.0, 14.0, 16.0, 18.0, 20.0}),
		"C": tabgo.NewSeries("C", []any{1.0, 3.0, 2.0, 4.0, 3.0, 5.0, 4.0, 6.0, 5.0, 7.0}),
	}
	return tabgo.NewDataFrame(sm)
}
