//go:build unit

package scigo

import (
	"math"
	"testing"
)

func TestChiSquareTest(t *testing.T) {
	// Fair die: observed matches expected perfectly
	observed := []float64{10, 10, 10, 10, 10, 10}
	expected := []float64{10, 10, 10, 10, 10, 10}
	stat, pval := ChiSquareTest(observed, expected)
	if stat != 0 {
		t.Errorf("Perfect fit: statistic = %v, want 0", stat)
	}
	if !approxEqual(pval, 1, 1e-8) {
		t.Errorf("Perfect fit: p-value = %v, want 1", pval)
	}

	// Known example: biased die
	observed2 := []float64{16, 18, 16, 14, 12, 12}
	expected2 := []float64{14.67, 14.67, 14.67, 14.67, 14.67, 14.67}
	stat2, pval2 := ChiSquareTest(observed2, expected2)
	if stat2 <= 0 {
		t.Error("Biased die should have positive statistic")
	}
	// With 5 df and small statistic, p-value should be > 0.05
	// statistic ≈ (1.33^2+3.33^2+1.33^2+0.67^2+2.67^2+2.67^2)/14.67 ≈ 2.046
	if !approxEqual(stat2, 2.046, 0.1) {
		t.Errorf("Biased die: statistic = %v, want ~2.046", stat2)
	}
	// p-value for chi2(5) at 2.046 should be large (fail to reject)
	if pval2 < 0.5 {
		t.Errorf("Biased die: p-value = %v, want > 0.5", pval2)
	}

	// Highly skewed: should reject null hypothesis
	observed3 := []float64{50, 0}
	expected3 := []float64{25, 25}
	stat3, pval3 := ChiSquareTest(observed3, expected3)
	if stat3 != 50 {
		t.Errorf("Skewed: statistic = %v, want 50", stat3)
	}
	if pval3 > 0.001 {
		t.Errorf("Skewed: p-value = %v, want < 0.001", pval3)
	}
}

func TestChiSquareTestPanic(t *testing.T) {
	// Different lengths
	defer func() {
		if r := recover(); r == nil {
			t.Error("Should panic on different lengths")
		}
	}()
	ChiSquareTest([]float64{1, 2}, []float64{1})
}

func TestChiSquareTestTooFew(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Should panic on fewer than 2 categories")
		}
	}()
	ChiSquareTest([]float64{1}, []float64{1})
}

func TestGTest(t *testing.T) {
	// Perfect fit: G = 0
	observed := []float64{10, 10, 10}
	expected := []float64{10, 10, 10}
	stat, pval := GTest(observed, expected)
	if stat != 0 {
		t.Errorf("Perfect fit: G statistic = %v, want 0", stat)
	}
	if !approxEqual(pval, 1, 1e-8) {
		t.Errorf("Perfect fit: p-value = %v, want 1", pval)
	}

	// Known computation: G = 2 * sum(O * ln(O/E))
	observed2 := []float64{20, 30}
	expected2 := []float64{25, 25}
	stat2, pval2 := GTest(observed2, expected2)
	wantG := 2 * (20*math.Log(20.0/25.0) + 30*math.Log(30.0/25.0))
	if !approxEqual(stat2, wantG, 1e-10) {
		t.Errorf("G statistic = %v, want %v", stat2, wantG)
	}
	// This is a mild deviation; p-value should be moderate
	if pval2 <= 0 || pval2 >= 1 {
		t.Errorf("G-test p-value = %v, expected in (0, 1)", pval2)
	}

	// Highly skewed
	observed3 := []float64{100, 0.1}
	expected3 := []float64{50, 50}
	stat3, pval3 := GTest(observed3, expected3)
	if stat3 <= 0 {
		t.Error("Skewed G-test should have positive statistic")
	}
	if pval3 > 0.001 {
		t.Errorf("Skewed G-test: p-value = %v, want < 0.001", pval3)
	}
}

func TestGTestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Should panic on different lengths")
		}
	}()
	GTest([]float64{1, 2}, []float64{1})
}

func TestGTestTooFew(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Should panic on fewer than 2 categories")
		}
	}()
	GTest([]float64{1}, []float64{1})
}
