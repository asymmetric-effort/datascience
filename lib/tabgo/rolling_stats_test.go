//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func TestRollingSkew(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0, 100.0})
	result := RollingSkew(s, 3)
	vals := result.Values()

	// First two should be nil
	if vals[0] != nil || vals[1] != nil {
		t.Errorf("RollingSkew: first 2 should be nil")
	}

	// Symmetric window [1,2,3] should have ~0 skew
	v2 := toFloat64(vals[2])
	if !almostEqual(v2, 0, 0.01) {
		t.Errorf("RollingSkew [1,2,3] = %v, want ~0", v2)
	}

	// Window with 100 should have positive skew
	v5 := toFloat64(vals[5])
	if v5 <= 0 {
		t.Errorf("RollingSkew [4,5,100] = %v, want > 0", v5)
	}
}

func TestRollingSkewSmallWindow(t *testing.T) {
	// Window of 2 means < 3 non-nil values, should get nil
	s := NewSeries("x", []any{1.0, 2.0, 3.0})
	result := RollingSkew(s, 2)
	vals := result.Values()
	// index 1: window [1,2] only 2 elements, need 3 for skew
	if vals[1] != nil {
		t.Errorf("RollingSkew window=2: expected nil for 2-element window, got %v", vals[1])
	}
}

func TestRollingKurtosis(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0})
	result := RollingKurtosis(s, 5)
	vals := result.Values()

	// First 4 should be nil
	for i := 0; i < 4; i++ {
		if vals[i] != nil {
			t.Errorf("RollingKurtosis[%d] should be nil", i)
		}
	}

	// Uniform-ish distribution should have negative excess kurtosis
	v4 := toFloat64(vals[4])
	if v4 > 0.5 {
		t.Errorf("RollingKurtosis [1..5] = %v, expect near -1.2 (platykurtic)", v4)
	}
}

func TestRollingKurtosisSmallWindow(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0})
	result := RollingKurtosis(s, 3)
	vals := result.Values()
	// Window of 3 has only 3 elements, need 4 for kurtosis
	if vals[2] != nil {
		t.Errorf("RollingKurtosis window=3: expected nil for 3-element window, got %v", vals[2])
	}
}

func TestRollingQuantile(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	result := RollingQuantile(s, 3, 0.5)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Errorf("RollingQuantile: first 2 should be nil")
	}

	// Median of [1,2,3] = 2
	if v := toFloat64(vals[2]); v != 2.0 {
		t.Errorf("RollingQuantile(0.5) [1,2,3] = %v, want 2.0", v)
	}

	// Median of [3,4,5] = 4
	if v := toFloat64(vals[4]); v != 4.0 {
		t.Errorf("RollingQuantile(0.5) [3,4,5] = %v, want 4.0", v)
	}
}

func TestRollingQuantileExtremes(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0})

	// q=0 should give min
	r0 := RollingQuantile(s, 3, 0.0)
	if v := toFloat64(r0.Values()[2]); v != 1.0 {
		t.Errorf("RollingQuantile(0) = %v, want 1.0", v)
	}

	// q=1 should give max
	r1 := RollingQuantile(s, 3, 1.0)
	if v := toFloat64(r1.Values()[2]); v != 3.0 {
		t.Errorf("RollingQuantile(1) = %v, want 3.0", v)
	}
}

func TestRollingZscore(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	result := RollingZscore(s, 3)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Errorf("RollingZscore: first 2 should be nil")
	}

	// For window [1,2,3], last value is 3, mean=2, std=1
	// z = (3-2)/1 = 1
	if v := toFloat64(vals[2]); !almostEqual(v, 1.0, 0.01) {
		t.Errorf("RollingZscore[2] = %v, want 1.0", v)
	}
}

func TestRollingZscoreConstant(t *testing.T) {
	s := NewSeries("x", []any{5.0, 5.0, 5.0})
	result := RollingZscore(s, 3)
	vals := result.Values()
	// Std = 0, should return 0
	if v := toFloat64(vals[2]); v != 0.0 {
		t.Errorf("RollingZscore constant = %v, want 0.0", v)
	}
}

func TestRollingCov(t *testing.T) {
	s1 := NewSeries("a", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	s2 := NewSeries("b", []any{2.0, 4.0, 6.0, 8.0, 10.0})
	result := RollingCov(s1, s2, 3)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Errorf("RollingCov: first 2 should be nil")
	}

	// Cov([1,2,3], [2,4,6]) = 2 * Var([1,2,3]) = 2 * 1 = 2
	if v := toFloat64(vals[2]); !almostEqual(v, 2.0, 0.01) {
		t.Errorf("RollingCov[2] = %v, want 2.0", v)
	}
}

func TestRollingCovUnequalLength(t *testing.T) {
	s1 := NewSeries("a", []any{1.0, 2.0, 3.0})
	s2 := NewSeries("b", []any{2.0, 4.0})
	result := RollingCov(s1, s2, 2)
	if result.Len() != 2 {
		t.Errorf("RollingCov unequal length: got %d, want 2", result.Len())
	}
}

func TestRollingSharpe(t *testing.T) {
	// Returns: all 0.1
	s := NewSeries("r", []any{0.1, 0.1, 0.1, 0.1, 0.1})
	result := RollingSharpe(s, 3, 0.0)
	vals := result.Values()

	// Std = 0, constant returns -> Sharpe = 0
	if v := toFloat64(vals[2]); v != 0.0 {
		t.Errorf("RollingSharpe constant = %v, want 0.0", v)
	}

	// Varying returns
	s2 := NewSeries("r", []any{0.05, 0.10, 0.15, 0.20, 0.25})
	result2 := RollingSharpe(s2, 3, 0.0)
	vals2 := result2.Values()
	v := toFloat64(vals2[2])
	if v <= 0 {
		t.Errorf("RollingSharpe positive returns = %v, want > 0", v)
	}
}

func TestRollingSortino(t *testing.T) {
	// Mix of positive and negative returns
	s := NewSeries("r", []any{0.05, -0.10, 0.15, -0.05, 0.20})
	result := RollingSortino(s, 3, 0.0)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Errorf("RollingSortino: first 2 should be nil")
	}

	// Window [0.05, -0.10, 0.15]: mean=0.0333, has downside
	v := toFloat64(vals[2])
	// Should be a finite number
	if math.IsNaN(v) || math.IsInf(v, 0) {
		t.Errorf("RollingSortino[2] = %v, want finite", v)
	}
}

func TestRollingSortinoNoDownside(t *testing.T) {
	s := NewSeries("r", []any{0.1, 0.2, 0.3})
	result := RollingSortino(s, 3, 0.0)
	vals := result.Values()
	// All positive returns, no downside -> dd=0 -> Sortino=0
	if v := toFloat64(vals[2]); v != 0.0 {
		t.Errorf("RollingSortino no downside = %v, want 0.0", v)
	}
}

func TestRollingMaxDrawdown(t *testing.T) {
	s := NewSeries("p", []any{100.0, 110.0, 90.0, 95.0, 105.0})
	result := RollingMaxDrawdown(s, 3)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Errorf("RollingMaxDrawdown: first 2 should be nil")
	}

	// Window [100, 110, 90]: peak=110, trough=90, dd = 20/110 = 0.1818
	v2 := toFloat64(vals[2])
	if !almostEqual(v2, 20.0/110.0, 0.001) {
		t.Errorf("RollingMaxDrawdown[2] = %v, want %v", v2, 20.0/110.0)
	}
}

func TestRollingMaxDrawdownFlat(t *testing.T) {
	s := NewSeries("p", []any{100.0, 100.0, 100.0})
	result := RollingMaxDrawdown(s, 3)
	vals := result.Values()
	if v := toFloat64(vals[2]); v != 0.0 {
		t.Errorf("RollingMaxDrawdown flat = %v, want 0.0", v)
	}
}

func TestRollingVaR(t *testing.T) {
	// Simple data: sorted is [-0.05, -0.02, 0.01, 0.03, 0.05]
	s := NewSeries("r", []any{0.01, -0.02, 0.05, 0.03, -0.05})
	result := RollingVaR(s, 5, 0.95)
	vals := result.Values()

	for i := 0; i < 4; i++ {
		if vals[i] != nil {
			t.Errorf("RollingVaR[%d] should be nil", i)
		}
	}

	// 5% quantile of [-0.05, -0.02, 0.01, 0.03, 0.05] at q=0.05
	// VaR = -quantile(data, 0.05)
	v := toFloat64(vals[4])
	if v < 0 {
		t.Errorf("RollingVaR = %v, want >= 0", v)
	}
}

func TestRollingCVaR(t *testing.T) {
	s := NewSeries("r", []any{-0.10, -0.05, 0.01, 0.03, 0.05})
	result := RollingCVaR(s, 5, 0.95)
	vals := result.Values()

	for i := 0; i < 4; i++ {
		if vals[i] != nil {
			t.Errorf("RollingCVaR[%d] should be nil", i)
		}
	}

	v := toFloat64(vals[4])
	// CVaR should be >= VaR
	varResult := RollingVaR(s, 5, 0.95)
	varV := toFloat64(varResult.Values()[4])
	if v < varV-0.01 {
		t.Errorf("RollingCVaR = %v < VaR = %v", v, varV)
	}
}

func TestSkewnessZeroVariance(t *testing.T) {
	vals := []float64{5, 5, 5}
	if skewness(vals) != 0 {
		t.Errorf("skewness of constant = %v, want 0", skewness(vals))
	}
}

func TestKurtosisZeroVariance(t *testing.T) {
	vals := []float64{5, 5, 5, 5}
	if kurtosis(vals) != 0 {
		t.Errorf("kurtosis of constant = %v, want 0", kurtosis(vals))
	}
}

func TestQuantileSingleElement(t *testing.T) {
	v := quantile([]float64{42.0}, 0.5)
	if v != 42.0 {
		t.Errorf("quantile single = %v, want 42.0", v)
	}
}

func TestDownsideDeviationNoDownside(t *testing.T) {
	dd := downsideDeviation([]float64{0.1, 0.2, 0.3}, 0.0)
	if dd != 0 {
		t.Errorf("downsideDeviation no downside = %v, want 0", dd)
	}
}

func TestDownsideDeviationSingleDownside(t *testing.T) {
	dd := downsideDeviation([]float64{-0.1, 0.2, 0.3}, 0.0)
	// Only one value below threshold, n<2, returns 0
	if dd != 0 {
		t.Errorf("downsideDeviation single = %v, want 0", dd)
	}
}

func TestMaxDrawdownEmpty(t *testing.T) {
	if maxDrawdown([]float64{}) != 0 {
		t.Error("maxDrawdown empty should be 0")
	}
}

func TestMaxDrawdownMonotonic(t *testing.T) {
	// Strictly increasing -> no drawdown
	if maxDrawdown([]float64{1, 2, 3, 4, 5}) != 0 {
		t.Error("maxDrawdown increasing should be 0")
	}
}

func TestExtractWindowNils(t *testing.T) {
	vals := []any{nil, 1.0, nil, 2.0}
	w := extractWindow(vals, 3, 4)
	if len(w) != 2 {
		t.Errorf("extractWindow with nils: got %d elements, want 2", len(w))
	}
}

func TestRollingCovNilValues(t *testing.T) {
	s1 := NewSeries("a", []any{nil, 2.0, 3.0})
	s2 := NewSeries("b", []any{nil, 4.0, 6.0})
	result := RollingCov(s1, s2, 2)
	// Only 1 non-nil value in window at index 1, should still work
	vals := result.Values()
	if vals[0] != nil {
		t.Errorf("RollingCov[0] should be nil")
	}
}

func TestRollingQuantileWithNils(t *testing.T) {
	s := NewSeries("x", []any{nil, nil, nil})
	result := RollingQuantile(s, 2, 0.5)
	vals := result.Values()
	// No non-nil values in window -> nil
	if vals[1] != nil {
		t.Errorf("RollingQuantile all nils should be nil, got %v", vals[1])
	}
}

func TestRollingSharpeSmallWindow(t *testing.T) {
	s := NewSeries("r", []any{0.1, 0.2})
	result := RollingSharpe(s, 1, 0.0)
	vals := result.Values()
	// Window=1 means only 1 element, < 2 -> nil
	if vals[0] != nil {
		t.Errorf("RollingSharpe window=1: expected nil, got %v", vals[0])
	}
}

func TestRollingMaxDrawdownSingleWindow(t *testing.T) {
	s := NewSeries("p", []any{100.0})
	result := RollingMaxDrawdown(s, 1)
	vals := result.Values()
	if v := toFloat64(vals[0]); v != 0.0 {
		t.Errorf("RollingMaxDrawdown single = %v, want 0", v)
	}
}

func TestRollingVaRSmallWindow(t *testing.T) {
	s := NewSeries("r", []any{-0.05, 0.05, 0.10})
	result := RollingVaR(s, 2, 0.95)
	vals := result.Values()
	if vals[0] != nil {
		t.Errorf("RollingVaR[0] should be nil")
	}
	// Index 1: window [-0.05, 0.05], VaR should be positive
	v := toFloat64(vals[1])
	if v < 0 {
		t.Errorf("RollingVaR[1] = %v, want >= 0", v)
	}
}

func TestRollingCVaRSmallWindow(t *testing.T) {
	s := NewSeries("r", []any{-0.10, -0.05, 0.01})
	result := RollingCVaR(s, 2, 0.95)
	vals := result.Values()
	if vals[0] != nil {
		t.Errorf("RollingCVaR[0] should be nil")
	}
}

func TestRollingSortinoMultipleDownside(t *testing.T) {
	s := NewSeries("r", []any{-0.05, -0.10, -0.15, 0.20, 0.10})
	result := RollingSortino(s, 4, 0.0)
	vals := result.Values()
	v := toFloat64(vals[3])
	// Mean is negative-ish, dd > 0, should give a negative Sortino
	if math.IsNaN(v) || math.IsInf(v, 0) {
		t.Errorf("RollingSortino = %v, want finite", v)
	}
}

func TestRollingZscoreWithNils(t *testing.T) {
	// Window with nils producing < 2 non-nil values
	s := NewSeries("x", []any{nil, nil, 1.0})
	result := RollingZscore(s, 2)
	vals := result.Values()
	// Window at index 1: [nil, nil] -> len(w)=0 < 2 -> nil
	if vals[1] != nil {
		t.Errorf("RollingZscore nils: expected nil, got %v", vals[1])
	}
}

func TestRollingSortinoWithNils(t *testing.T) {
	s := NewSeries("r", []any{nil, nil, 0.1})
	result := RollingSortino(s, 2, 0.0)
	vals := result.Values()
	if vals[1] != nil {
		t.Errorf("RollingSortino nils: expected nil, got %v", vals[1])
	}
}

func TestRollingMaxDrawdownWithNils(t *testing.T) {
	s := NewSeries("p", []any{nil, nil, 100.0})
	result := RollingMaxDrawdown(s, 2)
	vals := result.Values()
	// Window at index 1: [nil, nil] -> len=0 -> nil
	if vals[1] != nil {
		t.Errorf("RollingMaxDrawdown nils: expected nil, got %v", vals[1])
	}
}

func TestRollingVaRWithNils(t *testing.T) {
	s := NewSeries("r", []any{nil, nil, 0.01})
	result := RollingVaR(s, 2, 0.95)
	vals := result.Values()
	if vals[1] != nil {
		t.Errorf("RollingVaR nils: expected nil, got %v", vals[1])
	}
}

func TestRollingCVaRWithNils(t *testing.T) {
	s := NewSeries("r", []any{nil, nil, 0.01})
	result := RollingCVaR(s, 2, 0.95)
	vals := result.Values()
	if vals[1] != nil {
		t.Errorf("RollingCVaR nils: expected nil, got %v", vals[1])
	}
}

func TestRollingCVaRNoValuesBelow(t *testing.T) {
	// All values above the threshold quantile
	s := NewSeries("r", []any{0.10, 0.10, 0.10})
	result := RollingCVaR(s, 3, 0.5)
	vals := result.Values()
	// All values identical = same as threshold, so count should be > 0
	v := toFloat64(vals[2])
	if math.IsNaN(v) || math.IsInf(v, 0) {
		t.Errorf("RollingCVaR constant = %v, want finite", v)
	}
}

func TestRollingCovWithNilsInWindow(t *testing.T) {
	s1 := NewSeries("a", []any{nil, nil, 3.0})
	s2 := NewSeries("b", []any{nil, nil, 6.0})
	result := RollingCov(s1, s2, 2)
	vals := result.Values()
	// Window at index 1: both nil -> len < 2 -> nil
	if vals[1] != nil {
		t.Errorf("RollingCov nils: expected nil, got %v", vals[1])
	}
}

func TestRollingCovMismatchedNils(t *testing.T) {
	s1 := NewSeries("a", []any{nil, 2.0, 3.0})
	s2 := NewSeries("b", []any{1.0, nil, 6.0})
	result := RollingCov(s1, s2, 2)
	vals := result.Values()
	// At index 1: s1 window=[nil,2.0] -> [2.0], s2 window=[1.0,nil] -> [1.0]
	// Lengths match (1 each) but < 2, so nil
	if vals[1] != nil {
		t.Errorf("RollingCov mismatched nils: expected nil, got %v", vals[1])
	}
}
