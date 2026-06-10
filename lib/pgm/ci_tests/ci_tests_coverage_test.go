//go:build unit

package ci_tests

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// TestZKey_SingleColumn exercises the single-column path of zKey.
func TestZKey_SingleColumn(t *testing.T) {
	cols := [][]string{{"a", "b", "c"}}
	key := zKey(cols, 1)
	if key != "b" {
		t.Errorf("expected 'b', got %q", key)
	}
}

// TestZKey_MultipleColumns exercises the multi-column path with separator.
func TestZKey_MultipleColumns(t *testing.T) {
	cols := [][]string{{"a", "b"}, {"x", "y"}, {"1", "2"}}
	key := zKey(cols, 0)
	if key != "a|x|1" {
		t.Errorf("expected 'a|x|1', got %q", key)
	}
	key = zKey(cols, 1)
	if key != "b|y|2" {
		t.Errorf("expected 'b|y|2', got %q", key)
	}
}

// TestComputeExpected_ZeroTotal exercises the zero total path.
func TestComputeExpected_ZeroTotal(t *testing.T) {
	table := []float64{0, 0, 0, 0}
	exp := computeExpected(table, 2, 2)
	for _, v := range exp {
		if v != 0 {
			t.Errorf("expected 0 for zero total, got %f", v)
		}
	}
}

// TestNewDecisionTree_DefaultParams exercises the default parameter paths.
func TestNewDecisionTree_DefaultParams(t *testing.T) {
	dt := newDecisionTree(0, 0)
	if dt.maxDepth != 10 {
		t.Errorf("expected default maxDepth=10, got %d", dt.maxDepth)
	}
	if dt.minSamples != 2 {
		t.Errorf("expected default minSamples=2, got %d", dt.minSamples)
	}
}

// TestNewDecisionTree_CustomParams exercises the normal parameter path.
func TestNewDecisionTree_CustomParams(t *testing.T) {
	dt := newDecisionTree(5, 10)
	if dt.maxDepth != 5 {
		t.Errorf("expected maxDepth=5, got %d", dt.maxDepth)
	}
	if dt.minSamples != 10 {
		t.Errorf("expected minSamples=10, got %d", dt.minSamples)
	}
}

// TestNewDecisionTree_MinSamplesOne exercises minSamples < 2 boundary.
func TestNewDecisionTree_MinSamplesOne(t *testing.T) {
	dt := newDecisionTree(5, 1)
	if dt.minSamples != 2 {
		t.Errorf("expected minSamples clamped to 2, got %d", dt.minSamples)
	}
}

// TestDecisionTree_FitEmpty exercises the empty data path.
func TestDecisionTree_FitEmpty(t *testing.T) {
	dt := newDecisionTree(5, 2)
	dt.fit(nil, nil)
	if dt.root == nil {
		t.Fatal("expected non-nil root")
	}
	if !dt.root.isLeaf {
		t.Error("expected leaf node for empty data")
	}
	if dt.root.value != 0 {
		t.Errorf("expected value=0, got %f", dt.root.value)
	}
}

// TestSolveLinearSystem_Singular exercises the singular matrix path.
func TestSolveLinearSystem_Singular(t *testing.T) {
	// 2x2 singular matrix: rows are identical.
	A := []float64{1, 1, 1, 1}
	b := []float64{2, 2}
	x := solveLinearSystem(A, b, 2)
	if len(x) != 2 {
		t.Fatalf("expected 2 results, got %d", len(x))
	}
	// Should not crash; result is undefined for singular systems.
}

// TestSolveLinearSystem_RowSwap exercises the row-swap path.
func TestSolveLinearSystem_RowSwap(t *testing.T) {
	// System where pivot row needs swapping: [0, 1; 1, 0] * x = [3, 2]
	A := []float64{0, 1, 1, 0}
	b := []float64{3, 2}
	x := solveLinearSystem(A, b, 2)
	if math.Abs(x[0]-2) > 1e-10 || math.Abs(x[1]-3) > 1e-10 {
		t.Errorf("expected [2, 3], got [%f, %f]", x[0], x[1])
	}
}

// TestMultivariateCIBase_InsufficientDF exercises the df < 1 path.
func TestMultivariateCIBase_InsufficientDF(t *testing.T) {
	// With very few data points and many conditioning variables,
	// df should be < 1.
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0}),
		"Y": tabgo.NewSeries("Y", []any{3.0, 4.0}),
		"Z": tabgo.NewSeries("Z", []any{5.0, 6.0}),
	}
	data := tabgo.NewDataFrame(sm)

	// n=2, k=1, df = 2-2-1 = -1 < 1
	_, _, _, ok := multivariateCIBase("X", "Y", []string{"Z"}, data)
	if ok {
		t.Error("expected ok=false for insufficient df")
	}
}

// TestPowerDivergence_EmptyData exercises PowerDivergence with empty data.
func TestPowerDivergence_EmptyData(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{}),
		"Y": tabgo.NewSeries("Y", []any{}),
	}
	data := tabgo.NewDataFrame(sm)

	test := PowerDivergence(0.5)
	stat, pval, indep := test("X", "Y", nil, data, 0.05)
	if stat != 0 || pval != 1 || !indep {
		t.Errorf("expected (0, 1, true) for empty data, got (%f, %f, %v)", stat, pval, indep)
	}
}

// TestPearsonrEquivalence_AbsRGEEpsilon exercises the |r| >= epsilon path.
func TestPearsonrEquivalence_AbsRGEEpsilon(t *testing.T) {
	// Create strongly correlated data where |r| will be large.
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0, 4.0, 6.0, 8.0, 10.0, 12.0, 14.0, 16.0, 18.0, 20.0}),
	}
	data := tabgo.NewDataFrame(sm)

	test := PearsonrEquivalence(0.01) // very small epsilon
	stat, pval, indep := test("X", "Y", nil, data, 0.05)
	// |r| ≈ 1.0 >= 0.01 => cannot claim equivalence
	if indep {
		t.Error("expected indep=false for strongly correlated data with tiny epsilon")
	}
	if pval != 1 {
		t.Errorf("expected pval=1, got %f", pval)
	}
	_ = stat
}

// TestPearsonrEquivalence_SEZero exercises the se=0 fallback path.
func TestPearsonrEquivalence_SEZero(t *testing.T) {
	// With perfectly correlated data, 1-r^2 could be 0.
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{0.0, 1.0, 2.0, 3.0, 4.0}),
		"Y": tabgo.NewSeries("Y", []any{0.0, 0.0, 0.0, 0.0, 0.0}), // constant => r=0, se = sqrt(1/3)
	}
	data := tabgo.NewDataFrame(sm)

	test := PearsonrEquivalence(0.5)
	stat, pval, indep := test("X", "Y", nil, data, 0.05)
	_ = stat
	_ = pval
	_ = indep
}

// TestPearsonr_InsufficientDF exercises the df < 1 path.
func TestPearsonr_InsufficientDF(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0}),
		"Y": tabgo.NewSeries("Y", []any{3.0, 4.0}),
	}
	data := tabgo.NewDataFrame(sm)

	stat, pval, indep := Pearsonr("X", "Y", nil, data, 0.05)
	// n=2, df=2-2-0=0 < 1
	if stat != 0 || pval != 1 || !indep {
		t.Errorf("expected (0, 1, true) for insufficient df, got (%f, %f, %v)", stat, pval, indep)
	}
}

// TestPartialCorFromData_InsufficientData exercises the n <= k+2 path.
func TestPartialCorFromData_InsufficientData(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0}),
	}
	data := tabgo.NewDataFrame(sm)

	r, ok := partialCorFromData("X", "Y", nil, data)
	if ok {
		t.Errorf("expected ok=false, got r=%f", r)
	}
}

// TestDecisionTree_CandidateThresholds_AllSame exercises the degenerate case.
func TestDecisionTree_CandidateThresholds_AllSame(t *testing.T) {
	dt := newDecisionTree(5, 2)
	X := [][]float64{{5.0}, {5.0}, {5.0}}
	indices := []int{0, 1, 2}
	thresholds := dt.candidateThresholds(X, indices, 0)
	if len(thresholds) != 0 {
		t.Errorf("expected 0 thresholds for constant values, got %d", len(thresholds))
	}
}

// TestBuildContingencyTables_WithZ exercises the z-conditioning path.
func TestBuildContingencyTables_WithZ(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{"a", "b", "a", "b", "a", "b"}),
		"Y": tabgo.NewSeries("Y", []any{"x", "y", "x", "y", "y", "x"}),
		"Z": tabgo.NewSeries("Z", []any{"0", "0", "1", "1", "0", "1"}),
	}
	data := tabgo.NewDataFrame(sm)

	tables, xLevels, yLevels, dfs := buildContingencyTables("X", "Y", []string{"Z"}, data)
	if len(tables) != 2 {
		t.Errorf("expected 2 strata, got %d", len(tables))
	}
	if len(xLevels) != 2 {
		t.Errorf("expected 2 x levels, got %d", len(xLevels))
	}
	if len(yLevels) != 2 {
		t.Errorf("expected 2 y levels, got %d", len(yLevels))
	}
	if len(dfs) != 2 {
		t.Errorf("expected 2 df values, got %d", len(dfs))
	}
}

// TestMultivariateCIBase_ValidCase exercises the normal computation path.
func TestMultivariateCIBase_ValidCase(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0, 4.0, 6.0, 8.0, 10.0, 12.0, 14.0, 16.0, 18.0, 20.0}),
	}
	data := tabgo.NewDataFrame(sm)

	rSq, fStat, df2, ok := multivariateCIBase("X", "Y", nil, data)
	if !ok {
		t.Error("expected ok=true")
	}
	if rSq < 0 || rSq > 1 {
		t.Errorf("expected rSq in [0,1], got %f", rSq)
	}
	_ = fStat
	_ = df2
}

// TestMultivariateCIBase_WithZ exercises conditioning.
func TestMultivariateCIBase_WithZ(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0, 4.0, 6.0, 8.0, 10.0, 12.0, 14.0, 16.0, 18.0, 20.0}),
		"Z": tabgo.NewSeries("Z", []any{0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0, 0.0, 1.0}),
	}
	data := tabgo.NewDataFrame(sm)

	rSq, _, _, ok := multivariateCIBase("X", "Y", []string{"Z"}, data)
	if !ok {
		t.Error("expected ok=true")
	}
	_ = rSq
}

// TestPearsonrEquivalence_NormalCase exercises normal computation.
func TestPearsonrEquivalence_NormalCase(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0, 10.0}),
		"Y": tabgo.NewSeries("Y", []any{5.0, 3.0, 7.0, 1.0, 8.0, 2.0, 6.0, 4.0, 9.0, 0.0}),
	}
	data := tabgo.NewDataFrame(sm)

	test := PearsonrEquivalence(0.99)
	stat, pval, indep := test("X", "Y", nil, data, 0.05)
	_ = stat
	_ = pval
	_ = indep
}

// TestPearsonrEquivalence_InsufficientDF exercises insufficient df.
func TestPearsonrEquivalence_InsufficientDF(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0}),
		"Y": tabgo.NewSeries("Y", []any{3.0, 4.0}),
	}
	data := tabgo.NewDataFrame(sm)

	test := PearsonrEquivalence(0.5)
	stat, pval, indep := test("X", "Y", nil, data, 0.05)
	if stat != 0 || pval != 1 || !indep {
		t.Errorf("expected (0, 1, true), got (%f, %f, %v)", stat, pval, indep)
	}
}

// TestResiduals_WithPredictors exercises the OLS residuals path.
func TestResiduals_WithPredictors(t *testing.T) {
	target := []float64{1, 2, 3, 4, 5}
	predictors := [][]float64{{1, 2, 3, 4, 5}}
	res := residuals(target, predictors)
	if len(res) != 5 {
		t.Fatalf("expected 5 residuals, got %d", len(res))
	}
}

// TestPowerDivergence_WithConditioning exercises conditioned path.
func TestPowerDivergence_WithConditioning(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{"a", "b", "a", "b", "a", "b", "a", "b"}),
		"Y": tabgo.NewSeries("Y", []any{"x", "y", "x", "y", "y", "x", "y", "x"}),
		"Z": tabgo.NewSeries("Z", []any{"0", "0", "0", "0", "1", "1", "1", "1"}),
	}
	data := tabgo.NewDataFrame(sm)

	test := PowerDivergence(0.5)
	stat, pval, indep := test("X", "Y", []string{"Z"}, data, 0.05)
	_ = stat
	_ = pval
	_ = indep
}

// TestBuildContingencyTables_EmptyData exercises empty data path.
func TestBuildContingencyTables_EmptyData(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{}),
		"Y": tabgo.NewSeries("Y", []any{}),
	}
	data := tabgo.NewDataFrame(sm)

	tables, xLevels, yLevels, dfs := buildContingencyTables("X", "Y", nil, data)
	if tables != nil || xLevels != nil || yLevels != nil || dfs != nil {
		t.Error("expected all nil for empty data")
	}
}

// TestChiSquare_WithConditioning exercises the z-conditioned chi-square test.
func TestChiSquare_WithConditioning_Coverage(t *testing.T) {
	sm := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{"a", "b", "a", "b", "a", "b", "a", "b"}),
		"Y": tabgo.NewSeries("Y", []any{"x", "y", "x", "y", "y", "x", "y", "x"}),
		"Z": tabgo.NewSeries("Z", []any{"0", "0", "0", "0", "1", "1", "1", "1"}),
	}
	data := tabgo.NewDataFrame(sm)

	stat, pval, indep := ChiSquare("X", "Y", []string{"Z"}, data, 0.05)
	_ = stat
	_ = pval
	_ = indep
}
