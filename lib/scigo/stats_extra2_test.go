//go:build unit

package scigo

import (
	"math"
	"testing"
)

func approxEqualSE2(a, b, tol float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return math.Abs(a-b) < tol
}

// ---------------------------------------------------------------------------
// Describe
// ---------------------------------------------------------------------------

func TestDescribe(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	d := Describe(data)
	if d.Nobs != 10 {
		t.Errorf("Describe.Nobs=%v, want 10", d.Nobs)
	}
	if !approxEqualSE2(d.Min, 1, 1e-10) {
		t.Errorf("Describe.Min=%v, want 1", d.Min)
	}
	if !approxEqualSE2(d.Max, 10, 1e-10) {
		t.Errorf("Describe.Max=%v, want 10", d.Max)
	}
	if !approxEqualSE2(d.Mean, 5.5, 1e-10) {
		t.Errorf("Describe.Mean=%v, want 5.5", d.Mean)
	}
	// Variance should be 55/9 = 9.166...
	if !approxEqualSE2(d.Variance, 55.0/6, 1e-10) {
		t.Errorf("Describe.Variance=%v, want %v", d.Variance, 55.0/6)
	}
	// Symmetric data: skewness ~ 0
	if !approxEqualSE2(d.Skewness, 0, 1e-10) {
		t.Errorf("Describe.Skewness=%v, want 0", d.Skewness)
	}
}

// ---------------------------------------------------------------------------
// IQR
// ---------------------------------------------------------------------------

func TestIQR(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	got := IQR(data)
	// Q1 = 3.25, Q3 = 7.75, IQR = 4.5
	if !approxEqualSE2(got, 4.5, 1e-10) {
		t.Errorf("IQR=%v, want 4.5", got)
	}
}

func TestIQR_Single(t *testing.T) {
	got := IQR([]float64{5})
	if !approxEqualSE2(got, 0, 1e-10) {
		t.Errorf("IQR([5])=%v, want 0", got)
	}
}

// ---------------------------------------------------------------------------
// Zscore
// ---------------------------------------------------------------------------

func TestZscore(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	z := Zscore(data)
	if len(z) != 5 {
		t.Fatalf("Zscore length=%v, want 5", len(z))
	}
	// Mean z should be ~0
	sum := 0.0
	for _, v := range z {
		sum += v
	}
	if !approxEqualSE2(sum, 0, 1e-10) {
		t.Errorf("mean of z-scores=%v, want 0", sum/5)
	}
	// First element should be the most negative
	if z[0] >= z[1] {
		t.Errorf("z[0]=%v should be < z[1]=%v", z[0], z[1])
	}
}

func TestZscore_Constant(t *testing.T) {
	data := []float64{5, 5, 5}
	z := Zscore(data)
	for i, v := range z {
		if v != 0 {
			t.Errorf("Zscore constant: z[%d]=%v, want 0", i, v)
		}
	}
}

// ---------------------------------------------------------------------------
// Zmap
// ---------------------------------------------------------------------------

func TestZmap(t *testing.T) {
	scores := []float64{6, 7, 8}
	compare := []float64{1, 2, 3, 4, 5}
	z := Zmap(scores, compare)
	// compare mean=3, std=sqrt(2.5)
	mean := 3.0
	std := math.Sqrt(2.5)
	for i, s := range scores {
		expected := (s - mean) / std
		if !approxEqualSE2(z[i], expected, 1e-10) {
			t.Errorf("Zmap[%d]=%v, want %v", i, z[i], expected)
		}
	}
}

// ---------------------------------------------------------------------------
// TrimMean
// ---------------------------------------------------------------------------

func TestTrimMean(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 100}
	// 10% trim: remove bottom 1 and top 1
	got := TrimMean(data, 0.1)
	// trimmed = [2,3,4,5,6,7,8,9], mean = 44/8 = 5.5
	if !approxEqualSE2(got, 5.5, 1e-10) {
		t.Errorf("TrimMean(0.1)=%v, want 5.5", got)
	}
}

func TestTrimMean_NoTrim(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	got := TrimMean(data, 0)
	if !approxEqualSE2(got, 3, 1e-10) {
		t.Errorf("TrimMean(0)=%v, want 3", got)
	}
}

// ---------------------------------------------------------------------------
// SEM
// ---------------------------------------------------------------------------

func TestSEM(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	got := SEM(data)
	// std = sqrt(2.5), sem = sqrt(2.5)/sqrt(5) = sqrt(0.5)
	expected := math.Sqrt(2.5 / 5)
	if !approxEqualSE2(got, expected, 1e-10) {
		t.Errorf("SEM=%v, want %v", got, expected)
	}
}

// ---------------------------------------------------------------------------
// Skew
// ---------------------------------------------------------------------------

func TestSkew_Symmetric(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5}
	got := Skew(data)
	if !approxEqualSE2(got, 0, 1e-10) {
		t.Errorf("Skew(symmetric)=%v, want 0", got)
	}
}

func TestSkew_RightSkewed(t *testing.T) {
	data := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 100}
	got := Skew(data)
	if got <= 0 {
		t.Errorf("Skew(right-skewed)=%v, should be positive", got)
	}
}

// ---------------------------------------------------------------------------
// Kurtosis
// ---------------------------------------------------------------------------

func TestKurtosis_Normal(t *testing.T) {
	// For a uniform distribution {1,...,n}, excess kurtosis = -1.2*(n^2+1)/((n^2-1))
	// For simplicity, just check that kurtosis of normally-looking data is near 0
	// Uniform on [1..5]: kurtosis = -6*(5^2+1)/((5^2-1)*5) = -6*26/120 = -1.3
	data := []float64{1, 2, 3, 4, 5}
	got := Kurtosis(data)
	// Kurtosis of uniform-like data should be negative (platykurtic)
	if got >= 0 {
		t.Errorf("Kurtosis(uniform)=%v, should be negative", got)
	}
}

func TestKurtosis_HighPeak(t *testing.T) {
	// Data with a heavy tail should have positive kurtosis
	data := []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 100}
	got := Kurtosis(data)
	if got <= 0 {
		t.Errorf("Kurtosis(heavy-tail)=%v, should be positive", got)
	}
}

// ---------------------------------------------------------------------------
// JarqueBera
// ---------------------------------------------------------------------------

func TestJarqueBera_Normal(t *testing.T) {
	// Generate quasi-normal data (symmetric, moderate kurtosis)
	data := make([]float64, 100)
	for i := range data {
		// Use a simple pattern: standard normal approximation
		data[i] = float64(i-50) / 20.0
	}
	stat, pval := JarqueBera(data)
	// For roughly normal data, statistic should be small and p-value large
	if stat < 0 {
		t.Errorf("JarqueBera stat=%v, should be non-negative", stat)
	}
	if pval < 0 || pval > 1 {
		t.Errorf("JarqueBera pvalue=%v, should be in [0,1]", pval)
	}
}

func TestJarqueBera_Skewed(t *testing.T) {
	// Highly skewed data should reject normality
	data := make([]float64, 50)
	for i := 0; i < 45; i++ {
		data[i] = 1
	}
	data[45] = 100
	data[46] = 200
	data[47] = 500
	data[48] = 1000
	data[49] = 5000
	stat, pval := JarqueBera(data)
	if stat <= 0 {
		t.Errorf("JarqueBera stat=%v for skewed data, should be positive", stat)
	}
	if pval > 0.05 {
		t.Errorf("JarqueBera pvalue=%v for skewed data, should be small", pval)
	}
}

// ---------------------------------------------------------------------------
// RankData
// ---------------------------------------------------------------------------

func TestRankData_Exported(t *testing.T) {
	data := []float64{3, 1, 4, 1, 5, 9, 2, 6}
	ranks := RankData(data)
	if len(ranks) != len(data) {
		t.Fatalf("RankData length=%v, want %v", len(ranks), len(data))
	}
	// Value 1 appears twice, should have average rank (1+2)/2 = 1.5
	// Indices 1 and 3 have value 1
	if !approxEqualSE2(ranks[1], 1.5, 1e-10) {
		t.Errorf("RankData[1]=%v, want 1.5", ranks[1])
	}
	if !approxEqualSE2(ranks[3], 1.5, 1e-10) {
		t.Errorf("RankData[3]=%v, want 1.5", ranks[3])
	}
	// Value 9 is the largest, rank 8
	if !approxEqualSE2(ranks[5], 8, 1e-10) {
		t.Errorf("RankData[5]=%v, want 8", ranks[5])
	}
}

func TestRankData_Exported_NoTies(t *testing.T) {
	data := []float64{10, 20, 30}
	ranks := RankData(data)
	expected := []float64{1, 2, 3}
	for i := range ranks {
		if !approxEqualSE2(ranks[i], expected[i], 1e-10) {
			t.Errorf("RankData[%d]=%v, want %v", i, ranks[i], expected[i])
		}
	}
}
