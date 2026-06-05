//go:build unit

package scigo

import (
	"math"
	"testing"
)

func approxEqualPF(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

// ---------------------------------------------------------------------------
// PortfolioReturn
// ---------------------------------------------------------------------------

func TestPortfolioReturn_Basic(t *testing.T) {
	w := []float64{0.5, 0.5}
	r := []float64{0.10, 0.20}
	got := PortfolioReturn(w, r)
	want := 0.15
	if !approxEqualPF(got, want, 1e-10) {
		t.Errorf("PortfolioReturn = %v, want %v", got, want)
	}
}

func TestPortfolioReturn_SingleAsset(t *testing.T) {
	w := []float64{1.0}
	r := []float64{0.08}
	got := PortfolioReturn(w, r)
	if !approxEqualPF(got, 0.08, 1e-10) {
		t.Errorf("PortfolioReturn = %v, want 0.08", got)
	}
}

// ---------------------------------------------------------------------------
// PortfolioRisk
// ---------------------------------------------------------------------------

func TestPortfolioRisk_Uncorrelated(t *testing.T) {
	w := []float64{0.5, 0.5}
	cov := [][]float64{{0.04, 0}, {0, 0.09}}
	got := PortfolioRisk(w, cov)
	// variance = 0.25*0.04 + 0.25*0.09 = 0.01 + 0.0225 = 0.0325
	want := math.Sqrt(0.0325)
	if !approxEqualPF(got, want, 1e-10) {
		t.Errorf("PortfolioRisk = %v, want %v", got, want)
	}
}

func TestPortfolioRisk_SingleAsset(t *testing.T) {
	w := []float64{1.0}
	cov := [][]float64{{0.04}}
	got := PortfolioRisk(w, cov)
	if !approxEqualPF(got, 0.2, 1e-10) {
		t.Errorf("PortfolioRisk = %v, want 0.2", got)
	}
}

func TestPortfolioRisk_Correlated(t *testing.T) {
	w := []float64{0.6, 0.4}
	sigma1, sigma2 := 0.2, 0.3
	rho := 0.5
	cov := [][]float64{
		{sigma1 * sigma1, rho * sigma1 * sigma2},
		{rho * sigma1 * sigma2, sigma2 * sigma2},
	}
	got := PortfolioRisk(w, cov)
	// var = 0.36*0.04 + 0.16*0.09 + 2*0.6*0.4*0.5*0.06 = 0.0144 + 0.0144 + 0.0144 = 0.0432
	want := math.Sqrt(0.36*0.04 + 0.16*0.09 + 2*0.6*0.4*rho*sigma1*sigma2)
	if !approxEqualPF(got, want, 1e-10) {
		t.Errorf("PortfolioRisk = %v, want %v", got, want)
	}
}

// ---------------------------------------------------------------------------
// SharpeRatio
// ---------------------------------------------------------------------------

func TestSharpeRatio_Basic(t *testing.T) {
	w := []float64{1.0}
	r := []float64{0.10}
	cov := [][]float64{{0.04}}
	rf := 0.02
	got := SharpeRatio(w, r, cov, rf)
	// (0.10 - 0.02) / 0.2 = 0.4
	if !approxEqualPF(got, 0.4, 1e-10) {
		t.Errorf("SharpeRatio = %v, want 0.4", got)
	}
}

func TestSharpeRatio_ZeroRisk(t *testing.T) {
	w := []float64{1.0}
	r := []float64{0.05}
	cov := [][]float64{{0}}
	got := SharpeRatio(w, r, cov, 0.02)
	if got != 0 {
		t.Errorf("SharpeRatio with zero risk = %v, want 0", got)
	}
}

// ---------------------------------------------------------------------------
// ValueAtRisk
// ---------------------------------------------------------------------------

func TestValueAtRisk_95(t *testing.T) {
	w := []float64{1.0}
	r := []float64{0.10}
	cov := [][]float64{{0.04}}
	got := ValueAtRisk(w, r, cov, 0.95)
	// VaR = -(mu + sigma * z_{0.05})
	// z_{0.05} = -1.6449...
	// VaR = -(0.10 + 0.2 * (-1.6449)) = -(0.10 - 0.3290) = 0.2290
	stdNorm := NewNormal(0, 1)
	z := stdNorm.PPF(0.05)
	want := -(0.10 + 0.2*z)
	if !approxEqualPF(got, want, 1e-4) {
		t.Errorf("VaR = %v, want %v", got, want)
	}
}

func TestValueAtRisk_99(t *testing.T) {
	w := []float64{1.0}
	r := []float64{0.05}
	cov := [][]float64{{0.09}}
	got := ValueAtRisk(w, r, cov, 0.99)
	stdNorm := NewNormal(0, 1)
	z := stdNorm.PPF(0.01)
	want := -(0.05 + 0.3*z)
	if !approxEqualPF(got, want, 1e-4) {
		t.Errorf("VaR = %v, want %v", got, want)
	}
}

// ---------------------------------------------------------------------------
// ConditionalVaR
// ---------------------------------------------------------------------------

func TestConditionalVaR_95(t *testing.T) {
	w := []float64{1.0}
	r := []float64{0.10}
	cov := [][]float64{{0.04}}
	got := ConditionalVaR(w, r, cov, 0.95)
	// CVaR = -mu + sigma * phi(z) / (1-conf)
	stdNorm := NewNormal(0, 1)
	z := stdNorm.PPF(0.05)
	phiZ := stdNorm.PDF(z)
	want := -0.10 + 0.2*phiZ/0.05
	if !approxEqualPF(got, want, 1e-4) {
		t.Errorf("CVaR = %v, want %v", got, want)
	}
}

func TestConditionalVaR_GreaterThanVaR(t *testing.T) {
	w := []float64{1.0}
	r := []float64{0.10}
	cov := [][]float64{{0.04}}
	var95 := ValueAtRisk(w, r, cov, 0.95)
	cvar95 := ConditionalVaR(w, r, cov, 0.95)
	// CVaR should be >= VaR.
	if cvar95 < var95-1e-10 {
		t.Errorf("CVaR (%v) < VaR (%v), expected CVaR >= VaR", cvar95, var95)
	}
}

// ---------------------------------------------------------------------------
// MinVariancePortfolio
// ---------------------------------------------------------------------------

func TestMinVariancePortfolio_2Asset_Analytical(t *testing.T) {
	// Two uncorrelated assets with variances 0.04 and 0.09.
	// Min variance: w1 = sigma2^2 / (sigma1^2 + sigma2^2) = 0.09/0.13
	returns := []float64{0.10, 0.20}
	cov := [][]float64{{0.04, 0}, {0, 0.09}}

	w, err := MinVariancePortfolio(returns, cov)
	if err != nil {
		t.Fatal(err)
	}

	w1Want := 0.09 / 0.13
	w2Want := 0.04 / 0.13
	if !approxEqualPF(w[0], w1Want, 0.02) {
		t.Errorf("w[0] = %v, want ~%v", w[0], w1Want)
	}
	if !approxEqualPF(w[1], w2Want, 0.02) {
		t.Errorf("w[1] = %v, want ~%v", w[1], w2Want)
	}

	// Weights sum to 1.
	sumW := w[0] + w[1]
	if !approxEqualPF(sumW, 1.0, 1e-6) {
		t.Errorf("weights sum = %v, want 1", sumW)
	}
}

func TestMinVariancePortfolio_Empty(t *testing.T) {
	_, err := MinVariancePortfolio(nil, nil)
	if err == nil {
		t.Error("expected error for empty returns")
	}
}

func TestMinVariancePortfolio_CovMismatch(t *testing.T) {
	_, err := MinVariancePortfolio([]float64{0.1, 0.2}, [][]float64{{0.04}})
	if err == nil {
		t.Error("expected error for cov size mismatch")
	}
}

func TestMinVariancePortfolio_SingleAsset(t *testing.T) {
	w, err := MinVariancePortfolio([]float64{0.1}, [][]float64{{0.04}})
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualPF(w[0], 1.0, 1e-6) {
		t.Errorf("single asset weight = %v, want 1", w[0])
	}
}

// ---------------------------------------------------------------------------
// MaxSharpePortfolio
// ---------------------------------------------------------------------------

func TestMaxSharpePortfolio_2Asset(t *testing.T) {
	returns := []float64{0.10, 0.20}
	cov := [][]float64{{0.04, 0}, {0, 0.09}}
	rf := 0.02

	w, err := MaxSharpePortfolio(returns, cov, rf)
	if err != nil {
		t.Fatal(err)
	}

	// Analytical: z = cov^{-1} * (r - rf)
	// z1 = (0.10-0.02)/0.04 = 2, z2 = (0.20-0.02)/0.09 = 2
	// w = z / sum(z) = [0.5, 0.5]
	if !approxEqualPF(w[0], 0.5, 0.02) {
		t.Errorf("w[0] = %v, want ~0.5", w[0])
	}
	if !approxEqualPF(w[1], 0.5, 0.02) {
		t.Errorf("w[1] = %v, want ~0.5", w[1])
	}
}

func TestMaxSharpePortfolio_Empty(t *testing.T) {
	_, err := MaxSharpePortfolio(nil, nil, 0.02)
	if err == nil {
		t.Error("expected error for empty returns")
	}
}

func TestMaxSharpePortfolio_CovMismatch(t *testing.T) {
	_, err := MaxSharpePortfolio([]float64{0.1, 0.2}, [][]float64{{0.04}}, 0.02)
	if err == nil {
		t.Error("expected error for cov size mismatch")
	}
}

// ---------------------------------------------------------------------------
// TargetReturnPortfolio
// ---------------------------------------------------------------------------

func TestTargetReturnPortfolio_Basic(t *testing.T) {
	returns := []float64{0.10, 0.20}
	cov := [][]float64{{0.04, 0.01}, {0.01, 0.09}}
	target := 0.15

	w, err := TargetReturnPortfolio(returns, cov, target)
	if err != nil {
		t.Fatal(err)
	}

	// Weights should sum to 1.
	sumW := w[0] + w[1]
	if !approxEqualPF(sumW, 1.0, 1e-6) {
		t.Errorf("weights sum = %v, want 1", sumW)
	}

	// Portfolio return should match target.
	gotRet := PortfolioReturn(w, returns)
	if !approxEqualPF(gotRet, target, 0.01) {
		t.Errorf("portfolio return = %v, want %v", gotRet, target)
	}
}

func TestTargetReturnPortfolio_Empty(t *testing.T) {
	_, err := TargetReturnPortfolio(nil, nil, 0.1)
	if err == nil {
		t.Error("expected error for empty returns")
	}
}

func TestTargetReturnPortfolio_CovMismatch(t *testing.T) {
	_, err := TargetReturnPortfolio([]float64{0.1, 0.2}, [][]float64{{0.04}}, 0.1)
	if err == nil {
		t.Error("expected error for cov size mismatch")
	}
}

// ---------------------------------------------------------------------------
// EfficientFrontier
// ---------------------------------------------------------------------------

func TestEfficientFrontier_Basic(t *testing.T) {
	returns := []float64{0.10, 0.20}
	cov := [][]float64{{0.04, 0.01}, {0.01, 0.09}}

	risks, rets, weights, err := EfficientFrontier(returns, cov, 5)
	if err != nil {
		t.Fatal(err)
	}

	if len(risks) != 5 || len(rets) != 5 || len(weights) != 5 {
		t.Fatalf("expected 5 points, got %d/%d/%d", len(risks), len(rets), len(weights))
	}

	// Returns should be non-decreasing.
	for i := 1; i < len(rets); i++ {
		if rets[i] < rets[i-1]-1e-6 {
			t.Errorf("returns not non-decreasing: rets[%d]=%v < rets[%d]=%v", i, rets[i], i-1, rets[i-1])
		}
	}

	// All risks should be positive.
	for i, r := range risks {
		if r < 0 {
			t.Errorf("risk[%d] = %v, want >= 0", i, r)
		}
	}

	// All weights should sum to 1.
	for i, w := range weights {
		sum := 0.0
		for _, v := range w {
			sum += v
		}
		if !approxEqualPF(sum, 1.0, 1e-4) {
			t.Errorf("weights[%d] sum = %v, want 1", i, sum)
		}
	}
}

func TestEfficientFrontier_Empty(t *testing.T) {
	_, _, _, err := EfficientFrontier(nil, nil, 5)
	if err == nil {
		t.Error("expected error for empty returns")
	}
}

func TestEfficientFrontier_TooFewPoints(t *testing.T) {
	_, _, _, err := EfficientFrontier([]float64{0.1, 0.2}, [][]float64{{0.04, 0}, {0, 0.09}}, 1)
	if err == nil {
		t.Error("expected error for nPoints < 2")
	}
}

// ---------------------------------------------------------------------------
// Portfolio metric consistency checks
// ---------------------------------------------------------------------------

func TestPortfolioMetrics_Consistency(t *testing.T) {
	// For a given portfolio, check that metrics are internally consistent.
	w := []float64{0.6, 0.4}
	r := []float64{0.12, 0.18}
	cov := [][]float64{{0.04, 0.01}, {0.01, 0.09}}
	rf := 0.03

	pr := PortfolioReturn(w, r)
	risk := PortfolioRisk(w, cov)
	sharpe := SharpeRatio(w, r, cov, rf)

	// Sharpe = (return - rf) / risk
	if !approxEqualPF(sharpe, (pr-rf)/risk, 1e-10) {
		t.Errorf("Sharpe ratio inconsistent: %v vs %v", sharpe, (pr-rf)/risk)
	}

	// VaR at 50% confidence should be approximately -mu (since z_{0.5} = 0).
	var50 := ValueAtRisk(w, r, cov, 0.5)
	if !approxEqualPF(var50, -pr, 1e-4) {
		t.Errorf("VaR(50%%) = %v, want ~%v", var50, -pr)
	}
}

func TestPortfolioRisk_NegativeVarianceGuard(t *testing.T) {
	// Use a non-PSD matrix that produces negative w'Cw.
	w := []float64{1, 1}
	cov := [][]float64{{1, -3}, {-3, 1}} // not PSD
	got := PortfolioRisk(w, cov)
	// w'Cw = 1+(-3)+(-3)+1 = -4, clamped to 0 => risk = 0
	if got != 0 {
		t.Errorf("PortfolioRisk with negative variance = %v, want 0", got)
	}
}

func TestEfficientFrontier_FallbackPath(t *testing.T) {
	// Two assets with very close returns; some target returns may fail,
	// exercising the fallback to mvWeights.
	returns := []float64{0.10, 0.100001}
	cov := [][]float64{{0.04, 0.039}, {0.039, 0.04}}

	// This should still succeed even if internal TargetReturn fails for some points.
	risks, rets, weights, err := EfficientFrontier(returns, cov, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(risks) != 3 || len(rets) != 3 || len(weights) != 3 {
		t.Fatalf("expected 3 points")
	}
}

func TestFlattenCov(t *testing.T) {
	cov := [][]float64{{1, 2}, {3, 4}}
	flat := flattenCov(cov, 2)
	want := []float64{1, 2, 3, 4}
	for i, v := range flat {
		if v != want[i] {
			t.Errorf("flattenCov[%d] = %v, want %v", i, v, want[i])
		}
	}
}

func TestNormalizeWeights(t *testing.T) {
	w := normalizeWeights([]float64{-0.01, 0.5, 0.51})
	sum := 0.0
	for _, v := range w {
		if v < 0 {
			t.Errorf("negative weight after normalization: %v", v)
		}
		sum += v
	}
	if !approxEqualPF(sum, 1.0, 1e-10) {
		t.Errorf("weights sum = %v, want 1", sum)
	}
}

func TestNormalizeWeights_AllZero(t *testing.T) {
	w := normalizeWeights([]float64{-1, -2})
	// All negative => all clipped to zero, sum=0. Should not panic.
	for _, v := range w {
		if v != 0 {
			t.Errorf("expected all zeros, got %v", w)
			break
		}
	}
}

// ---------------------------------------------------------------------------
// "Put-call parity equivalent": VaR/CVaR relationship
// ---------------------------------------------------------------------------

func TestVaR_CVaR_Relationship(t *testing.T) {
	// CVaR >= VaR for any confidence level.
	w := []float64{0.3, 0.4, 0.3}
	r := []float64{0.08, 0.12, 0.15}
	cov := [][]float64{
		{0.04, 0.01, 0.005},
		{0.01, 0.06, 0.02},
		{0.005, 0.02, 0.09},
	}

	for _, conf := range []float64{0.90, 0.95, 0.99} {
		var_ := ValueAtRisk(w, r, cov, conf)
		cvar := ConditionalVaR(w, r, cov, conf)
		if cvar < var_-1e-10 {
			t.Errorf("CVaR (%v) < VaR (%v) at confidence %v", cvar, var_, conf)
		}
	}
}

func TestMaxSharpePortfolio_Constrained(t *testing.T) {
	// Returns where unconstrained tangency gives negative weights,
	// forcing the constrained path.
	returns := []float64{0.20, 0.05, 0.03}
	cov := [][]float64{
		{0.04, 0.01, 0.005},
		{0.01, 0.02, 0.01},
		{0.005, 0.01, 0.01},
	}
	rf := 0.10 // high rf makes some excess returns negative

	w, err := MaxSharpePortfolio(returns, cov, rf)
	if err != nil {
		t.Fatal(err)
	}
	sum := 0.0
	for _, v := range w {
		if v < -1e-6 {
			t.Errorf("negative weight: %v", v)
		}
		sum += v
	}
	if !approxEqualPF(sum, 1.0, 1e-4) {
		t.Errorf("weights sum = %v, want 1", sum)
	}
}

func TestMaxSharpePortfolio_Singular(t *testing.T) {
	// Singular covariance should error.
	returns := []float64{0.1, 0.2}
	cov := [][]float64{{1, 1}, {1, 1}} // singular
	_, err := MaxSharpePortfolio(returns, cov, 0.02)
	if err == nil {
		t.Error("expected error for singular covariance")
	}
}

func TestMaxSharpePortfolio_Degenerate(t *testing.T) {
	// All excess returns zero => sum(z) = 0 => degenerate.
	returns := []float64{0.05, 0.05}
	cov := [][]float64{{0.04, 0}, {0, 0.09}}
	_, err := MaxSharpePortfolio(returns, cov, 0.05)
	if err == nil {
		t.Error("expected error for degenerate (all excess returns = 0)")
	}
}

func TestMinVariancePortfolio_3Asset(t *testing.T) {
	returns := []float64{0.08, 0.12, 0.15}
	cov := [][]float64{
		{0.04, 0.006, 0.002},
		{0.006, 0.06, 0.01},
		{0.002, 0.01, 0.09},
	}

	w, err := MinVariancePortfolio(returns, cov)
	if err != nil {
		t.Fatal(err)
	}

	// All weights should be non-negative and sum to 1.
	sum := 0.0
	for i, v := range w {
		if v < -1e-6 {
			t.Errorf("w[%d] = %v, want >= 0", i, v)
		}
		sum += v
	}
	if !approxEqualPF(sum, 1.0, 1e-4) {
		t.Errorf("weights sum = %v, want 1", sum)
	}

	// The min-variance portfolio should have lower or equal risk than equal-weight.
	eqW := []float64{1.0 / 3.0, 1.0 / 3.0, 1.0 / 3.0}
	mvRisk := PortfolioRisk(w, cov)
	eqRisk := PortfolioRisk(eqW, cov)
	if mvRisk > eqRisk+1e-6 {
		t.Errorf("min variance risk (%v) > equal weight risk (%v)", mvRisk, eqRisk)
	}
}
