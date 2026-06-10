//go:build unit

package prediction

import (
	"math"
	"math/rand"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// tolerance for floating-point comparisons in tests.
const tol = 0.3

// Helper: create a DataFrame from column name -> []float64 mapping.
func makeDF(cols map[string][]float64) *tabgo.DataFrame {
	m := make(map[string]*tabgo.Series, len(cols))
	for name, vals := range cols {
		anyVals := make([]any, len(vals))
		for i, v := range vals {
			anyVals[i] = v
		}
		m[name] = tabgo.NewSeries(name, anyVals)
	}
	return tabgo.NewDataFrame(m)
}

// --- OLS helper tests ---

func TestOlsFitSimple(t *testing.T) {
	// y = 2 + 3*x, no noise
	n := 100
	X := make([][]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x := float64(i) / float64(n)
		X[i] = []float64{1.0, x} // intercept, x
		y[i] = 2.0 + 3.0*x
	}
	beta := olsFit(y, X)
	if len(beta) != 2 {
		t.Fatalf("expected 2 coefficients, got %d", len(beta))
	}
	if math.Abs(beta[0]-2.0) > 1e-10 {
		t.Errorf("intercept: got %f, want 2.0", beta[0])
	}
	if math.Abs(beta[1]-3.0) > 1e-10 {
		t.Errorf("slope: got %f, want 3.0", beta[1])
	}
}

func TestOlsFitMultiple(t *testing.T) {
	// y = 1 + 2*x1 + 3*x2
	n := 200
	rng := rand.New(rand.NewSource(42))
	X := make([][]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x1 := rng.Float64()
		x2 := rng.Float64()
		X[i] = []float64{1.0, x1, x2}
		y[i] = 1.0 + 2.0*x1 + 3.0*x2
	}
	beta := olsFit(y, X)
	if math.Abs(beta[0]-1.0) > 1e-10 {
		t.Errorf("intercept: got %f, want 1.0", beta[0])
	}
	if math.Abs(beta[1]-2.0) > 1e-10 {
		t.Errorf("beta1: got %f, want 2.0", beta[1])
	}
	if math.Abs(beta[2]-3.0) > 1e-10 {
		t.Errorf("beta2: got %f, want 3.0", beta[2])
	}
}

// --- DoubleML tests ---

func TestDoubleMLSyntheticData(t *testing.T) {
	// Generate data with known treatment effect = 2.0.
	// DGP:
	//   confounder ~ Uniform(0,1)
	//   treatment = 0.5 * confounder + noise
	//   outcome = 2.0 * treatment + 1.0 * confounder + noise
	// True ATE = 2.0
	n := 1000
	rng := rand.New(rand.NewSource(123))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 4
		tr := 0.5*c + rng.NormFloat64()*0.1
		y := 2.0*tr + 1.0*c + rng.NormFloat64()*0.1
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	dml := NewDoubleMLRegressor("treatment", "outcome", []string{"confounder"})
	err := dml.Fit(df)
	if err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	ate := dml.ATE()
	if math.Abs(ate-2.0) > tol {
		t.Errorf("DML ATE: got %f, want ~2.0 (tolerance %f)", ate, tol)
	}
}

func TestDoubleMLPredict(t *testing.T) {
	n := 200
	rng := rand.New(rand.NewSource(456))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 2
		tr := 0.5*c + rng.NormFloat64()*0.1
		y := 3.0*tr + c + rng.NormFloat64()*0.1
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	dml := NewDoubleMLRegressor("treatment", "outcome", []string{"confounder"})
	if err := dml.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	preds, err := dml.Predict(df)
	if err != nil {
		t.Fatalf("Predict error: %v", err)
	}
	if len(preds) != n {
		t.Fatalf("expected %d predictions, got %d", n, len(preds))
	}

	// Each prediction should be ATE * treatment
	ate := dml.ATE()
	for i, p := range preds {
		expected := ate * treatments[i]
		if math.Abs(p-expected) > 1e-10 {
			t.Errorf("prediction[%d]: got %f, want %f", i, p, expected)
			break
		}
	}
}

func TestDoubleMLPredictNotFitted(t *testing.T) {
	dml := NewDoubleMLRegressor("t", "y", []string{"c"})
	_, err := dml.Predict(makeDF(map[string][]float64{
		"t": {1, 2}, "y": {1, 2}, "c": {1, 2},
	}))
	if err == nil {
		t.Error("expected error for unfitted model")
	}
}

func TestDoubleMLTooFewObservations(t *testing.T) {
	df := makeDF(map[string][]float64{
		"t": {1, 2}, "y": {3, 4}, "c": {5, 6},
	})
	dml := NewDoubleMLRegressor("t", "y", []string{"c"})
	err := dml.Fit(df)
	if err == nil {
		t.Error("expected error for too few observations")
	}
}

// --- NaiveAdjustment tests ---

func TestNaiveAdjustmentSimple(t *testing.T) {
	// y = 3*treatment + 2*confounder + 1 (no noise)
	n := 200
	rng := rand.New(rand.NewSource(789))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 3
		tr := rng.Float64() * 2
		y := 3.0*tr + 2.0*c + 1.0
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	adj := NewNaiveAdjustmentRegressor("treatment", "outcome", []string{"confounder"})
	err := adj.Fit(df)
	if err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	ate := adj.ATE()
	if math.Abs(ate-3.0) > 0.01 {
		t.Errorf("NaiveAdjustment ATE: got %f, want 3.0", ate)
	}
}

func TestNaiveAdjustmentWithNoise(t *testing.T) {
	// y = 5*treatment + 1*confounder + noise
	n := 500
	rng := rand.New(rand.NewSource(101))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 4
		tr := 0.3*c + rng.NormFloat64()*0.5
		y := 5.0*tr + 1.0*c + rng.NormFloat64()*0.2
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	adj := NewNaiveAdjustmentRegressor("treatment", "outcome", []string{"confounder"})
	if err := adj.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	ate := adj.ATE()
	if math.Abs(ate-5.0) > tol {
		t.Errorf("NaiveAdjustment ATE: got %f, want ~5.0 (tolerance %f)", ate, tol)
	}
}

func TestNaiveAdjustmentNotFitted(t *testing.T) {
	adj := NewNaiveAdjustmentRegressor("t", "y", []string{"c"})
	ate := adj.ATE()
	if ate != 0 {
		t.Errorf("expected 0 for unfitted model, got %f", ate)
	}
}

func TestNaiveAdjustmentMultipleConfounders(t *testing.T) {
	// y = 4*treatment + 2*c1 + 3*c2 + 10
	n := 300
	rng := rand.New(rand.NewSource(202))

	c1 := make([]float64, n)
	c2 := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c1[i] = rng.Float64() * 2
		c2[i] = rng.Float64() * 3
		treatments[i] = rng.Float64() * 2
		outcomes[i] = 4.0*treatments[i] + 2.0*c1[i] + 3.0*c2[i] + 10.0
	}

	df := makeDF(map[string][]float64{
		"c1":        c1,
		"c2":        c2,
		"treatment": treatments,
		"outcome":   outcomes,
	})

	adj := NewNaiveAdjustmentRegressor("treatment", "outcome", []string{"c1", "c2"})
	if err := adj.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}
	if math.Abs(adj.ATE()-4.0) > 0.01 {
		t.Errorf("ATE: got %f, want 4.0", adj.ATE())
	}
}

// --- NaiveIV tests ---

func TestNaiveIVSimple(t *testing.T) {
	// DGP with instrument:
	//   z ~ Uniform(0, 5)           (instrument)
	//   u ~ Normal(0, 0.5)          (unobserved confounder)
	//   treatment = 2*z + u
	//   outcome = 3*treatment + u   (true ATE = 3.0)
	//
	// OLS would be biased because u affects both treatment and outcome.
	// 2SLS with z as instrument should recover ATE = 3.0.
	n := 1000
	rng := rand.New(rand.NewSource(303))

	instruments := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		z := rng.Float64() * 5
		u := rng.NormFloat64() * 0.5
		tr := 2.0*z + u
		y := 3.0*tr + u
		instruments[i] = z
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"instrument": instruments,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	iv := NewNaiveIVRegressor("treatment", "outcome", []string{"instrument"})
	err := iv.Fit(df)
	if err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	ate := iv.ATE()
	if math.Abs(ate-3.0) > tol {
		t.Errorf("NaiveIV ATE: got %f, want ~3.0 (tolerance %f)", ate, tol)
	}
}

func TestNaiveIVMultipleInstruments(t *testing.T) {
	// Two instruments, true ATE = 2.0.
	n := 800
	rng := rand.New(rand.NewSource(404))

	z1 := make([]float64, n)
	z2 := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		z1[i] = rng.Float64() * 3
		z2[i] = rng.Float64() * 2
		u := rng.NormFloat64() * 0.3
		treatments[i] = 1.0*z1[i] + 0.5*z2[i] + u
		outcomes[i] = 2.0*treatments[i] + u
	}

	df := makeDF(map[string][]float64{
		"z1":        z1,
		"z2":        z2,
		"treatment": treatments,
		"outcome":   outcomes,
	})

	iv := NewNaiveIVRegressor("treatment", "outcome", []string{"z1", "z2"})
	if err := iv.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	ate := iv.ATE()
	if math.Abs(ate-2.0) > tol {
		t.Errorf("NaiveIV ATE: got %f, want ~2.0 (tolerance %f)", ate, tol)
	}
}

func TestNaiveIVNotFitted(t *testing.T) {
	iv := NewNaiveIVRegressor("t", "y", []string{"z"})
	if iv.ATE() != 0 {
		t.Errorf("expected 0 for unfitted model, got %f", iv.ATE())
	}
}

func TestNaiveIVvsOLSBias(t *testing.T) {
	// Demonstrate that IV corrects for confounding bias.
	// DGP: u confounds both treatment and outcome.
	// OLS (naive adjustment without u) will be biased.
	// IV with a valid instrument should be closer to the true effect.
	n := 1000
	rng := rand.New(rand.NewSource(505))
	trueATE := 4.0

	instruments := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		z := rng.Float64() * 5
		u := rng.NormFloat64() * 1.0
		tr := 1.5*z + 2.0*u
		y := trueATE*tr + 3.0*u + rng.NormFloat64()*0.1
		instruments[i] = z
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"instrument": instruments,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	iv := NewNaiveIVRegressor("treatment", "outcome", []string{"instrument"})
	if err := iv.Fit(df); err != nil {
		t.Fatalf("IV Fit error: %v", err)
	}

	ivATE := iv.ATE()
	ivErr := math.Abs(ivATE - trueATE)

	// OLS without controlling for u should be more biased.
	ols := NewNaiveAdjustmentRegressor("treatment", "outcome", nil)
	if err := ols.Fit(df); err != nil {
		t.Fatalf("OLS Fit error: %v", err)
	}
	olsATE := ols.ATE()
	olsErr := math.Abs(olsATE - trueATE)

	t.Logf("True ATE=%.2f, IV ATE=%.4f (err=%.4f), OLS ATE=%.4f (err=%.4f)", trueATE, ivATE, ivErr, olsATE, olsErr)

	if ivErr > tol {
		t.Errorf("IV estimate too far from true ATE: got %f, want ~%f", ivATE, trueATE)
	}
	if ivErr >= olsErr {
		t.Logf("WARNING: IV not more accurate than OLS in this sample (IV err=%f, OLS err=%f)", ivErr, olsErr)
	}
}

// --- Gaussian elimination / linear system tests ---

func TestSolveLinearSystem(t *testing.T) {
	// 2x + y = 5
	// x + 3y = 7
	// Solution: x=1.6, y=1.8
	A := [][]float64{
		{2, 1},
		{1, 3},
	}
	b := []float64{5, 7}
	x := solveLinearSystem(A, b)
	if math.Abs(x[0]-1.6) > 1e-10 {
		t.Errorf("x[0]: got %f, want 1.6", x[0])
	}
	if math.Abs(x[1]-1.8) > 1e-10 {
		t.Errorf("x[1]: got %f, want 1.8", x[1])
	}
}

// --- Normal distribution helper tests ---

func TestNormalCDF(t *testing.T) {
	// CDF(0) = 0.5
	if math.Abs(normalCDF(0)-0.5) > 1e-10 {
		t.Errorf("normalCDF(0) = %f, want 0.5", normalCDF(0))
	}
	// CDF(-inf) -> 0
	if normalCDF(-10) > 1e-10 {
		t.Errorf("normalCDF(-10) should be near 0, got %f", normalCDF(-10))
	}
	// CDF(1.96) ~ 0.975
	if math.Abs(normalCDF(1.96)-0.975) > 0.001 {
		t.Errorf("normalCDF(1.96) = %f, want ~0.975", normalCDF(1.96))
	}
}

func TestNormalQuantile(t *testing.T) {
	// quantile(0.5) = 0
	if math.Abs(normalQuantile(0.5)) > 1e-10 {
		t.Errorf("normalQuantile(0.5) = %f, want 0", normalQuantile(0.5))
	}
	// quantile(0.975) ~ 1.96
	if math.Abs(normalQuantile(0.975)-1.96) > 0.01 {
		t.Errorf("normalQuantile(0.975) = %f, want ~1.96", normalQuantile(0.975))
	}
	// quantile(0.025) ~ -1.96
	if math.Abs(normalQuantile(0.025)+1.96) > 0.01 {
		t.Errorf("normalQuantile(0.025) = %f, want ~-1.96", normalQuantile(0.025))
	}
}

func TestInvertMatrix(t *testing.T) {
	// 2x2 identity
	A := [][]float64{{1, 0}, {0, 1}}
	inv := invertMatrix(A)
	if inv == nil {
		t.Fatal("expected non-nil inverse")
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if math.Abs(inv[i][j]-expected) > 1e-10 {
				t.Errorf("inv[%d][%d] = %f, want %f", i, j, inv[i][j], expected)
			}
		}
	}

	// 2x2 non-trivial: [[2,1],[1,3]], inverse = [[3/5, -1/5],[-1/5, 2/5]]
	B := [][]float64{{2, 1}, {1, 3}}
	inv2 := invertMatrix(B)
	if inv2 == nil {
		t.Fatal("expected non-nil inverse")
	}
	if math.Abs(inv2[0][0]-0.6) > 1e-10 {
		t.Errorf("inv[0][0] = %f, want 0.6", inv2[0][0])
	}
	if math.Abs(inv2[0][1]+0.2) > 1e-10 {
		t.Errorf("inv[0][1] = %f, want -0.2", inv2[0][1])
	}
}

// --- DoubleML SE/CI/PValue/Summary/CATE/NSplits tests ---

func TestDoubleMLSEAndCI(t *testing.T) {
	n := 1000
	rng := rand.New(rand.NewSource(777))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 4
		tr := 0.5*c + rng.NormFloat64()*0.1
		y := 2.0*tr + 1.0*c + rng.NormFloat64()*0.1
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	dml := NewDoubleMLRegressor("treatment", "outcome", []string{"confounder"})
	if err := dml.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	// SE should be positive and small relative to ATE.
	se := dml.SE()
	if se <= 0 {
		t.Errorf("SE should be positive, got %f", se)
	}
	if se > 1.0 {
		t.Errorf("SE unexpectedly large: %f", se)
	}

	// 95% CI should contain true ATE of 2.0.
	lo, hi := dml.ConfidenceInterval(0.05)
	if lo >= hi {
		t.Errorf("CI lower (%f) >= upper (%f)", lo, hi)
	}
	if lo > 2.0 || hi < 2.0 {
		t.Logf("WARNING: 95%% CI [%f, %f] does not contain true ATE 2.0", lo, hi)
	}

	// P-value should be small (ATE is significantly different from 0).
	pval := dml.PValue()
	if pval < 0 || pval > 1 {
		t.Errorf("P-value out of range: %f", pval)
	}
	if pval > 0.05 {
		t.Errorf("P-value should be <0.05 for large effect, got %f", pval)
	}
}

func TestDoubleMLNSplits(t *testing.T) {
	n := 600
	rng := rand.New(rand.NewSource(888))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 4
		tr := 0.5*c + rng.NormFloat64()*0.1
		y := 2.0*tr + 1.0*c + rng.NormFloat64()*0.1
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	// Test with 5-fold cross-fitting.
	dml := NewDoubleMLRegressor("treatment", "outcome", []string{"confounder"})
	dml.SetNSplits(5)
	if err := dml.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	ate := dml.ATE()
	if math.Abs(ate-2.0) > tol {
		t.Errorf("DML ATE with 5 folds: got %f, want ~2.0", ate)
	}
}

func TestDoubleMLEstimateCate(t *testing.T) {
	n := 500
	rng := rand.New(rand.NewSource(999))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 4
		tr := 0.5*c + rng.NormFloat64()*0.1
		y := 2.0*tr + 1.0*c + rng.NormFloat64()*0.1
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	dml := NewDoubleMLRegressor("treatment", "outcome", []string{"confounder"})
	if err := dml.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	cate, err := dml.EstimateCate()
	if err != nil {
		t.Fatalf("EstimateCate error: %v", err)
	}
	if len(cate) != n {
		t.Fatalf("expected %d CATE values, got %d", n, len(cate))
	}

	// For a linear model, CATE values should be centered around the ATE.
	mean := 0.0
	for _, v := range cate {
		mean += v
	}
	mean /= float64(len(cate))
	if math.Abs(mean-dml.ATE()) > 0.5 {
		t.Errorf("mean CATE (%f) should be close to ATE (%f)", mean, dml.ATE())
	}
}

func TestDoubleMLEstimateCateNotFitted(t *testing.T) {
	dml := NewDoubleMLRegressor("t", "y", []string{"c"})
	_, err := dml.EstimateCate()
	if err == nil {
		t.Error("expected error for unfitted model")
	}
}

func TestDoubleMLSummary(t *testing.T) {
	n := 200
	rng := rand.New(rand.NewSource(111))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 4
		tr := 0.5*c + rng.NormFloat64()*0.1
		y := 2.0*tr + 1.0*c + rng.NormFloat64()*0.1
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	dml := NewDoubleMLRegressor("treatment", "outcome", []string{"confounder"})
	// Not fitted.
	s := dml.Summary()
	if s != "DoubleMLRegressor: not fitted" {
		t.Errorf("unexpected summary for unfitted model: %s", s)
	}

	if err := dml.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	s = dml.Summary()
	if len(s) == 0 {
		t.Error("summary should not be empty")
	}
	// Check key fields are present.
	for _, substr := range []string{"ATE:", "Std. Error:", "95%", "P-value:", "treatment", "outcome"} {
		if !containsSubstring(s, substr) {
			t.Errorf("summary missing %q", substr)
		}
	}
}

// --- NaiveAdjustment Predict/SE/CI/PValue/Summary tests ---

func TestNaiveAdjustmentPredict(t *testing.T) {
	// y = 3*treatment + 2*confounder + 1 (no noise)
	n := 200
	rng := rand.New(rand.NewSource(789))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 3
		tr := rng.Float64() * 2
		y := 3.0*tr + 2.0*c + 1.0
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	adj := NewNaiveAdjustmentRegressor("treatment", "outcome", []string{"confounder"})
	if err := adj.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	preds, err := adj.Predict(df)
	if err != nil {
		t.Fatalf("Predict error: %v", err)
	}
	if len(preds) != n {
		t.Fatalf("expected %d predictions, got %d", n, len(preds))
	}

	// Predictions should be close to actual outcomes (no noise case).
	for i := 0; i < n; i++ {
		if math.Abs(preds[i]-outcomes[i]) > 0.01 {
			t.Errorf("prediction[%d]: got %f, want %f", i, preds[i], outcomes[i])
			break
		}
	}
}

func TestNaiveAdjustmentPredictNotFitted(t *testing.T) {
	adj := NewNaiveAdjustmentRegressor("t", "y", []string{"c"})
	_, err := adj.Predict(makeDF(map[string][]float64{
		"t": {1, 2}, "y": {1, 2}, "c": {1, 2},
	}))
	if err == nil {
		t.Error("expected error for unfitted model")
	}
}

func TestNaiveAdjustmentSEAndCI(t *testing.T) {
	n := 500
	rng := rand.New(rand.NewSource(222))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 4
		tr := rng.Float64() * 2
		y := 5.0*tr + 1.0*c + rng.NormFloat64()*0.2
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	adj := NewNaiveAdjustmentRegressor("treatment", "outcome", []string{"confounder"})
	if err := adj.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	se := adj.SE()
	if se <= 0 {
		t.Errorf("SE should be positive, got %f", se)
	}

	lo, hi := adj.ConfidenceInterval(0.05)
	if lo >= hi {
		t.Errorf("CI lower (%f) >= upper (%f)", lo, hi)
	}
	ate := adj.ATE()
	if lo > ate || hi < ate {
		t.Errorf("CI [%f, %f] should contain ATE %f", lo, hi, ate)
	}

	pval := adj.PValue()
	if pval < 0 || pval > 1 {
		t.Errorf("P-value out of range: %f", pval)
	}
	// With true ATE=5 and low noise, p-value should be very small.
	if pval > 0.05 {
		t.Errorf("P-value should be <0.05, got %f", pval)
	}
}

func TestNaiveAdjustmentSummary(t *testing.T) {
	n := 200
	rng := rand.New(rand.NewSource(333))

	confounders := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		c := rng.Float64() * 3
		tr := rng.Float64() * 2
		y := 3.0*tr + 2.0*c + 1.0 + rng.NormFloat64()*0.1
		confounders[i] = c
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"confounder": confounders,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	adj := NewNaiveAdjustmentRegressor("treatment", "outcome", []string{"confounder"})
	// Not fitted.
	s := adj.Summary()
	if s != "NaiveAdjustmentRegressor: not fitted" {
		t.Errorf("unexpected summary for unfitted model: %s", s)
	}

	if err := adj.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	s = adj.Summary()
	if len(s) == 0 {
		t.Error("summary should not be empty")
	}
	for _, substr := range []string{"ATE:", "Std. Error:", "95%", "P-value:", "treatment", "outcome"} {
		if !containsSubstring(s, substr) {
			t.Errorf("summary missing %q", substr)
		}
	}
}

// --- NaiveIV Predict/SE/CI/PValue/FirstStageFStat/Summary tests ---

func TestNaiveIVPredict(t *testing.T) {
	n := 500
	rng := rand.New(rand.NewSource(444))

	instruments := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		z := rng.Float64() * 5
		u := rng.NormFloat64() * 0.3
		tr := 2.0*z + u
		y := 3.0*tr + u
		instruments[i] = z
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"instrument": instruments,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	iv := NewNaiveIVRegressor("treatment", "outcome", []string{"instrument"})
	if err := iv.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	preds, err := iv.Predict(df)
	if err != nil {
		t.Fatalf("Predict error: %v", err)
	}
	if len(preds) != n {
		t.Fatalf("expected %d predictions, got %d", n, len(preds))
	}

	// Predictions should be correlated with actual outcomes.
	// Check they are not all zero.
	allZero := true
	for _, p := range preds {
		if p != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Error("all predictions are zero")
	}
}

func TestNaiveIVPredictNotFitted(t *testing.T) {
	iv := NewNaiveIVRegressor("t", "y", []string{"z"})
	_, err := iv.Predict(makeDF(map[string][]float64{
		"t": {1, 2}, "y": {1, 2}, "z": {1, 2},
	}))
	if err == nil {
		t.Error("expected error for unfitted model")
	}
}

func TestNaiveIVSEAndCI(t *testing.T) {
	n := 1000
	rng := rand.New(rand.NewSource(555))

	instruments := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		z := rng.Float64() * 5
		u := rng.NormFloat64() * 0.5
		tr := 2.0*z + u
		y := 3.0*tr + u
		instruments[i] = z
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"instrument": instruments,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	iv := NewNaiveIVRegressor("treatment", "outcome", []string{"instrument"})
	if err := iv.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	se := iv.SE()
	if se <= 0 {
		t.Errorf("SE should be positive, got %f", se)
	}

	lo, hi := iv.ConfidenceInterval(0.05)
	if lo >= hi {
		t.Errorf("CI lower (%f) >= upper (%f)", lo, hi)
	}

	pval := iv.PValue()
	if pval < 0 || pval > 1 {
		t.Errorf("P-value out of range: %f", pval)
	}
	if pval > 0.05 {
		t.Errorf("P-value should be <0.05 for large effect, got %f", pval)
	}
}

func TestNaiveIVFirstStageFStat(t *testing.T) {
	n := 1000
	rng := rand.New(rand.NewSource(666))

	instruments := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		z := rng.Float64() * 5
		u := rng.NormFloat64() * 0.5
		tr := 2.0*z + u // strong instrument
		y := 3.0*tr + u
		instruments[i] = z
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"instrument": instruments,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	iv := NewNaiveIVRegressor("treatment", "outcome", []string{"instrument"})
	if err := iv.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	fstat := iv.FirstStageFStat()
	// With a strong instrument (coeff=2.0, low noise), F-stat should be >> 10.
	if fstat < 10 {
		t.Errorf("F-stat should be >> 10 for strong instrument, got %f", fstat)
	}
	t.Logf("First stage F-stat: %.2f", fstat)
}

func TestNaiveIVFirstStageFStatNotFitted(t *testing.T) {
	iv := NewNaiveIVRegressor("t", "y", []string{"z"})
	if iv.FirstStageFStat() != 0 {
		t.Error("expected 0 for unfitted model")
	}
}

func TestNaiveIVSummary(t *testing.T) {
	n := 500
	rng := rand.New(rand.NewSource(777))

	instruments := make([]float64, n)
	treatments := make([]float64, n)
	outcomes := make([]float64, n)

	for i := 0; i < n; i++ {
		z := rng.Float64() * 5
		u := rng.NormFloat64() * 0.5
		tr := 2.0*z + u
		y := 3.0*tr + u
		instruments[i] = z
		treatments[i] = tr
		outcomes[i] = y
	}

	df := makeDF(map[string][]float64{
		"instrument": instruments,
		"treatment":  treatments,
		"outcome":    outcomes,
	})

	iv := NewNaiveIVRegressor("treatment", "outcome", []string{"instrument"})
	// Not fitted.
	s := iv.Summary()
	if s != "NaiveIVRegressor: not fitted" {
		t.Errorf("unexpected summary for unfitted model: %s", s)
	}

	if err := iv.Fit(df); err != nil {
		t.Fatalf("Fit error: %v", err)
	}

	s = iv.Summary()
	if len(s) == 0 {
		t.Error("summary should not be empty")
	}
	for _, substr := range []string{"ATE:", "Std. Error:", "95%", "P-value:", "1st Stage F:", "treatment", "outcome"} {
		if !containsSubstring(s, substr) {
			t.Errorf("summary missing %q", substr)
		}
	}
}

// containsSubstring checks if s contains substr.
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
