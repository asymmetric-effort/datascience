//go:build unit

package prediction

import (
	"math"
	"math/rand"
	"testing"
)

// TestNormalQuantile_EdgeCases exercises the boundary conditions.
func TestNormalQuantile_EdgeCases(t *testing.T) {
	if !math.IsInf(normalQuantile(0), -1) {
		t.Error("expected -Inf for p=0")
	}
	if !math.IsInf(normalQuantile(1), 1) {
		t.Error("expected +Inf for p=1")
	}
	// Lower tail (p < 0.02425)
	v := normalQuantile(0.001)
	if v >= 0 {
		t.Errorf("expected negative quantile for p=0.001, got %f", v)
	}
	// Upper tail (p > 1 - 0.02425)
	v = normalQuantile(0.999)
	if v <= 0 {
		t.Errorf("expected positive quantile for p=0.999, got %f", v)
	}
	// Central region
	v = normalQuantile(0.5)
	if math.Abs(v) > 0.01 {
		t.Errorf("expected near 0 for p=0.5, got %f", v)
	}
}

// TestBuildDesignMatrixNoIntercept exercises the no-intercept design matrix builder.
func TestBuildDesignMatrixNoIntercept(t *testing.T) {
	df := makeDF(map[string][]float64{
		"X1": {1, 2, 3},
		"X2": {4, 5, 6},
	})
	X, err := buildDesignMatrixNoIntercept(df, []string{"X1", "X2"})
	if err != nil {
		t.Fatal(err)
	}
	if len(X) != 3 {
		t.Errorf("expected 3 rows, got %d", len(X))
	}
	if len(X[0]) != 2 {
		t.Errorf("expected 2 cols, got %d", len(X[0]))
	}
	if X[0][0] != 1 || X[0][1] != 4 {
		t.Errorf("unexpected values in first row: %v", X[0])
	}
}

// TestBuildDesignMatrixNoIntercept_SingleCol exercises single column path.
func TestBuildDesignMatrixNoIntercept_SingleCol(t *testing.T) {
	df := makeDF(map[string][]float64{
		"X1": {1, 2, 3},
	})
	X, err := buildDesignMatrixNoIntercept(df, []string{"X1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(X) != 3 || len(X[0]) != 1 {
		t.Errorf("expected 3x1 matrix, got %dx%d", len(X), len(X[0]))
	}
}

// TestExtractColumnFloat64_ValidColumn exercises the valid column path.
func TestExtractColumnFloat64_ValidColumn(t *testing.T) {
	df := makeDF(map[string][]float64{
		"X": {1, 2, 3},
	})
	vals, err := extractColumnFloat64(df, "X")
	if err != nil {
		t.Fatal(err)
	}
	if len(vals) != 3 {
		t.Errorf("expected 3 values, got %d", len(vals))
	}
}

// TestSetNSplits_Clamped exercises the n < 2 clamping.
func TestSetNSplits_Clamped(t *testing.T) {
	d := NewDoubleMLRegressor("T", "Y", []string{"C"})
	d.SetNSplits(1)
	if d.nSplits != 2 {
		t.Errorf("expected nSplits clamped to 2, got %d", d.nSplits)
	}
	d.SetNSplits(5)
	if d.nSplits != 5 {
		t.Errorf("expected nSplits=5, got %d", d.nSplits)
	}
}

// TestDoubleML_PValue_ZeroSE exercises the se=0 path.
func TestDoubleML_PValue_ZeroSE(t *testing.T) {
	d := &DoubleMLRegressor{
		ate: 1.0,
		se:  0,
	}
	p := d.PValue()
	if p != 0 {
		t.Errorf("expected 0 for se=0, got %f", p)
	}
}

// TestDoubleML_Predict_NotFitted exercises the not-fitted path.
func TestDoubleML_Predict_NotFitted(t *testing.T) {
	d := NewDoubleMLRegressor("T", "Y", []string{"C"})
	_, err := d.Predict(nil)
	if err == nil {
		t.Error("expected error for not-fitted model")
	}
}

// TestNaiveAdjustment_PValue_ZeroSE exercises the se=0 path.
func TestNaiveAdjustment_PValue_ZeroSE(t *testing.T) {
	r := &NaiveAdjustmentRegressor{
		se:           0,
		fitted:       true,
		coefficients: []float64{0, 1.0},
	}
	p := r.PValue()
	if p != 0 {
		t.Errorf("expected 0 for se=0, got %f", p)
	}
}

// TestNaiveAdjustment_Predict_NotFitted exercises the not-fitted error path.
func TestNaiveAdjustment_Predict_NotFitted(t *testing.T) {
	r := NewNaiveAdjustmentRegressor("T", "Y", []string{"C"})
	_, err := r.Predict(nil)
	if err == nil {
		t.Error("expected error for not-fitted model")
	}
}

// TestNaiveIV_PValue_ZeroSE exercises the se=0 path.
func TestNaiveIV_PValue_ZeroSE(t *testing.T) {
	r := &NaiveIVRegressor{
		ate:    1.0,
		se:     0,
		fitted: true,
	}
	p := r.PValue()
	if p != 0 {
		t.Errorf("expected 0 for se=0, got %f", p)
	}
}

// TestNaiveIV_Predict_NotFitted exercises the not-fitted error path.
func TestNaiveIV_Predict_NotFitted(t *testing.T) {
	r := NewNaiveIVRegressor("T", "Y", []string{"Z"})
	_, err := r.Predict(nil)
	if err == nil {
		t.Error("expected error for not-fitted model")
	}
}

// TestNaiveIV_FirstStageFStat_NotFitted exercises the not-fitted path.
func TestNaiveIV_FirstStageFStat_NotFitted(t *testing.T) {
	r := NewNaiveIVRegressor("T", "Y", []string{"Z"})
	f := r.FirstStageFStat()
	if f != 0 {
		t.Errorf("expected 0 for not-fitted, got %f", f)
	}
}

// TestNaiveIV_Fit_Empty exercises the empty DataFrame error path.
func TestNaiveIV_Fit_Empty(t *testing.T) {
	r := NewNaiveIVRegressor("T", "Y", []string{"Z"})
	df := makeDF(map[string][]float64{
		"T": {},
		"Y": {},
		"Z": {},
	})
	err := r.Fit(df)
	if err == nil {
		t.Error("expected error for empty DataFrame")
	}
}

// TestNaiveAdjustment_Fit_Empty exercises the empty DataFrame error path.
func TestNaiveAdjustment_Fit_Empty(t *testing.T) {
	r := NewNaiveAdjustmentRegressor("T", "Y", []string{"C"})
	df := makeDF(map[string][]float64{
		"T": {},
		"Y": {},
		"C": {},
	})
	err := r.Fit(df)
	if err == nil {
		t.Error("expected error for empty DataFrame")
	}
}

// TestDoubleML_Fit_InsufficientData exercises the too-few-rows error path.
func TestDoubleML_Fit_InsufficientData(t *testing.T) {
	d := NewDoubleMLRegressor("T", "Y", []string{"C"})
	d.SetNSplits(3)
	df := makeDF(map[string][]float64{
		"T": {1, 2, 3},
		"Y": {1, 2, 3},
		"C": {1, 2, 3},
	})
	err := d.Fit(df)
	if err == nil {
		t.Error("expected error for insufficient data")
	}
}

// TestNaiveAdjustment_FullWorkflow exercises the full Fit -> PValue -> Predict workflow.
func TestNaiveAdjustment_FullWorkflow(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	n := 100
	treatment := make([]float64, n)
	outcome := make([]float64, n)
	confounder := make([]float64, n)
	for i := 0; i < n; i++ {
		confounder[i] = rng.Float64()
		treatment[i] = confounder[i] + rng.Float64()*0.1
		outcome[i] = 2.0*treatment[i] + confounder[i] + rng.Float64()*0.1
	}
	df := makeDF(map[string][]float64{
		"T": treatment,
		"Y": outcome,
		"C": confounder,
	})
	r := NewNaiveAdjustmentRegressor("T", "Y", []string{"C"})
	if err := r.Fit(df); err != nil {
		t.Fatal(err)
	}
	pval := r.PValue()
	if pval >= 1 || pval < 0 {
		t.Errorf("expected p-value in [0,1), got %f", pval)
	}
	preds, err := r.Predict(df)
	if err != nil {
		t.Fatal(err)
	}
	if len(preds) != n {
		t.Errorf("expected %d predictions, got %d", n, len(preds))
	}
}

// TestNaiveIV_FullWorkflow exercises the full Fit -> PValue -> Predict -> FStat workflow.
func TestNaiveIV_FullWorkflow(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	n := 100
	instrument := make([]float64, n)
	treatment := make([]float64, n)
	outcome := make([]float64, n)
	for i := 0; i < n; i++ {
		instrument[i] = rng.Float64()
		treatment[i] = instrument[i]*3 + rng.Float64()*0.1
		outcome[i] = 2.0*treatment[i] + rng.Float64()*0.1
	}
	df := makeDF(map[string][]float64{
		"T": treatment,
		"Y": outcome,
		"Z": instrument,
	})
	r := NewNaiveIVRegressor("T", "Y", []string{"Z"})
	if err := r.Fit(df); err != nil {
		t.Fatal(err)
	}
	pval := r.PValue()
	if pval >= 1 || pval < 0 {
		t.Errorf("expected p-value in [0,1), got %f", pval)
	}
	preds, err := r.Predict(df)
	if err != nil {
		t.Fatal(err)
	}
	if len(preds) != n {
		t.Errorf("expected %d predictions, got %d", n, len(preds))
	}
	fstat := r.FirstStageFStat()
	if fstat <= 0 {
		t.Errorf("expected positive F-stat, got %f", fstat)
	}
}

// TestInvertMatrix_SmallCase exercises the invertMatrix function.
func TestInvertMatrix_SmallCase(t *testing.T) {
	// 2x2 identity matrix
	A := [][]float64{{1, 0}, {0, 1}}
	inv := invertMatrix(A)
	if math.Abs(inv[0][0]-1) > 1e-10 || math.Abs(inv[1][1]-1) > 1e-10 {
		t.Error("expected identity inverse")
	}
	if math.Abs(inv[0][1]) > 1e-10 || math.Abs(inv[1][0]) > 1e-10 {
		t.Error("expected zero off-diagonals")
	}
}

// TestComputeCoefficientSE exercises the SE computation.
func TestComputeCoefficientSE(t *testing.T) {
	X := [][]float64{{1, 1}, {1, 2}, {1, 3}}
	// sigma2 = variance of residuals
	sigma2 := 0.01
	se0 := computeCoefficientSE(X, sigma2, 0)
	se1 := computeCoefficientSE(X, sigma2, 1)
	if se0 < 0 {
		t.Error("SE should be non-negative")
	}
	if se1 < 0 {
		t.Error("SE should be non-negative")
	}
}

// TestNormalCDF_Values exercises the normalCDF function.
func TestNormalCDF_Values(t *testing.T) {
	if math.Abs(normalCDF(0)-0.5) > 0.01 {
		t.Errorf("expected CDF(0) near 0.5, got %f", normalCDF(0))
	}
	if normalCDF(-10) > 0.001 {
		t.Errorf("expected CDF(-10) near 0, got %f", normalCDF(-10))
	}
	if normalCDF(10) < 0.999 {
		t.Errorf("expected CDF(10) near 1, got %f", normalCDF(10))
	}
}
