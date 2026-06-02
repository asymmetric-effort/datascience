//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// TTestInd
// ---------------------------------------------------------------------------

func TestTTestInd_EqualSamples(t *testing.T) {
	// Two identical samples should have t=0, p=1
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{1, 2, 3, 4, 5}
	stat, pval := TTestInd(x, y)
	if !approxEqual(stat, 0, 1e-10) {
		t.Errorf("TTestInd identical: stat=%v, want 0", stat)
	}
	if !approxEqual(pval, 1, 1e-6) {
		t.Errorf("TTestInd identical: p=%v, want 1", pval)
	}
}

func TestTTestInd_ClearDifference(t *testing.T) {
	// Two well-separated samples: should reject H0
	x := []float64{10, 11, 12, 13, 14}
	y := []float64{1, 2, 3, 4, 5}
	stat, pval := TTestInd(x, y)
	if stat <= 0 {
		t.Errorf("TTestInd separated: stat=%v, want positive", stat)
	}
	if pval > 0.01 {
		t.Errorf("TTestInd separated: p=%v, want < 0.01", pval)
	}
}

func TestTTestInd_PanicSmall(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("TTestInd should panic with < 2 elements")
		}
	}()
	TTestInd([]float64{1}, []float64{2, 3})
}

func TestTTestInd_ConstantData(t *testing.T) {
	x := []float64{5, 5, 5}
	y := []float64{5, 5, 5}
	stat, pval := TTestInd(x, y)
	if stat != 0 {
		t.Errorf("TTestInd constant: stat=%v, want 0", stat)
	}
	if pval != 1 {
		t.Errorf("TTestInd constant: p=%v, want 1", pval)
	}
}

// ---------------------------------------------------------------------------
// TTest1Samp
// ---------------------------------------------------------------------------

func TestTTest1Samp_MeanEqual(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	// mean = 3
	stat, pval := TTest1Samp(x, 3.0)
	if !approxEqual(stat, 0, 1e-10) {
		t.Errorf("TTest1Samp mu=mean: stat=%v, want 0", stat)
	}
	if !approxEqual(pval, 1, 1e-6) {
		t.Errorf("TTest1Samp mu=mean: p=%v, want 1", pval)
	}
}

func TestTTest1Samp_MeanFarAway(t *testing.T) {
	x := []float64{10, 11, 12, 13, 14}
	stat, pval := TTest1Samp(x, 0)
	if stat <= 0 {
		t.Errorf("TTest1Samp far: stat=%v, want positive", stat)
	}
	if pval > 0.001 {
		t.Errorf("TTest1Samp far: p=%v, want < 0.001", pval)
	}
}

func TestTTest1Samp_PanicSmall(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("TTest1Samp should panic with < 2 elements")
		}
	}()
	TTest1Samp([]float64{1}, 0)
}

// ---------------------------------------------------------------------------
// TTestRel
// ---------------------------------------------------------------------------

func TestTTestRel_NoDifference(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{1, 2, 3, 4, 5}
	stat, pval := TTestRel(x, y)
	// All differences are 0 -> constant => se=0
	// Special case: should be 0, 1
	_ = stat
	_ = pval
}

func TestTTestRel_ClearDifference(t *testing.T) {
	x := []float64{10, 11, 12, 13, 14}
	y := []float64{1, 2, 3, 4, 5}
	stat, pval := TTestRel(x, y)
	if stat <= 0 {
		t.Errorf("TTestRel diff: stat=%v, want positive", stat)
	}
	if pval > 0.001 {
		t.Errorf("TTestRel diff: p=%v, want < 0.001", pval)
	}
}

func TestTTestRel_PanicDiffLen(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("TTestRel should panic on different lengths")
		}
	}()
	TTestRel([]float64{1, 2}, []float64{1})
}

func TestTTestRel_SmallEffect(t *testing.T) {
	// Differences with high variance relative to mean -> not significant
	x := []float64{1.1, 1.9, 3.2, 3.8, 5.1}
	y := []float64{1.0, 2.1, 3.0, 4.0, 5.0}
	// Mixed positive and negative differences, should not be significant
	_, pval := TTestRel(x, y)
	if pval < 0.05 {
		t.Errorf("TTestRel small effect: p=%v, expect > 0.05", pval)
	}
}

// ---------------------------------------------------------------------------
// MannWhitneyU
// ---------------------------------------------------------------------------

func TestMannWhitneyU_Identical(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{1, 2, 3, 4, 5}
	_, pval := MannWhitneyU(x, y)
	// Identical distributions: p-value should be high
	if pval < 0.05 {
		t.Errorf("MannWhitneyU identical: p=%v, want > 0.05", pval)
	}
}

func TestMannWhitneyU_ClearDifference(t *testing.T) {
	x := []float64{10, 11, 12, 13, 14}
	y := []float64{1, 2, 3, 4, 5}
	stat, pval := MannWhitneyU(x, y)
	// U should be 0 (no overlap)
	if stat != 0 {
		t.Errorf("MannWhitneyU separated: stat=%v, want 0", stat)
	}
	if pval > 0.05 {
		t.Errorf("MannWhitneyU separated: p=%v, want < 0.05", pval)
	}
}

func TestMannWhitneyU_PanicEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MannWhitneyU should panic on empty sample")
		}
	}()
	MannWhitneyU([]float64{}, []float64{1})
}

// ---------------------------------------------------------------------------
// WilcoxonSignedRank
// ---------------------------------------------------------------------------

func TestWilcoxonSignedRank_Symmetric(t *testing.T) {
	// Symmetric around zero: should not reject
	x := []float64{-3, -1, 1, 3}
	_, pval := WilcoxonSignedRank(x)
	if pval < 0.05 {
		t.Errorf("Wilcoxon symmetric: p=%v, want > 0.05", pval)
	}
}

func TestWilcoxonSignedRank_AllPositive(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	stat, pval := WilcoxonSignedRank(x)
	// All positive: T+ = n*(n+1)/2 = 55
	if stat != 55 {
		t.Errorf("Wilcoxon all positive: stat=%v, want 55", stat)
	}
	if pval > 0.05 {
		t.Errorf("Wilcoxon all positive: p=%v, want < 0.05", pval)
	}
}

func TestWilcoxonSignedRank_AllZeros(t *testing.T) {
	x := []float64{0, 0, 0}
	stat, pval := WilcoxonSignedRank(x)
	if stat != 0 {
		t.Errorf("Wilcoxon zeros: stat=%v, want 0", stat)
	}
	if pval != 1 {
		t.Errorf("Wilcoxon zeros: p=%v, want 1", pval)
	}
}

func TestWilcoxonSignedRank_PanicEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("WilcoxonSignedRank should panic on empty")
		}
	}()
	WilcoxonSignedRank([]float64{})
}

// ---------------------------------------------------------------------------
// KruskalWallis
// ---------------------------------------------------------------------------

func TestKruskalWallis_EqualGroups(t *testing.T) {
	g1 := []float64{1, 2, 3, 4, 5}
	g2 := []float64{1, 2, 3, 4, 5}
	_, pval := KruskalWallis(g1, g2)
	if pval < 0.05 {
		t.Errorf("KruskalWallis equal: p=%v, want > 0.05", pval)
	}
}

func TestKruskalWallis_DifferentGroups(t *testing.T) {
	g1 := []float64{1, 2, 3, 4, 5}
	g2 := []float64{10, 11, 12, 13, 14}
	g3 := []float64{20, 21, 22, 23, 24}
	stat, pval := KruskalWallis(g1, g2, g3)
	if stat <= 0 {
		t.Errorf("KruskalWallis different: stat=%v, want positive", stat)
	}
	if pval > 0.01 {
		t.Errorf("KruskalWallis different: p=%v, want < 0.01", pval)
	}
}

func TestKruskalWallis_PanicOneGroup(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("KruskalWallis should panic with < 2 groups")
		}
	}()
	KruskalWallis([]float64{1, 2, 3})
}

// ---------------------------------------------------------------------------
// FriedmanChiSquare
// ---------------------------------------------------------------------------

func TestFriedmanChiSquare_NoDifference(t *testing.T) {
	g1 := []float64{1, 2, 3, 4, 5}
	g2 := []float64{1, 2, 3, 4, 5}
	stat, pval := FriedmanChiSquare(g1, g2)
	if stat != 0 {
		t.Errorf("Friedman equal: stat=%v, want 0", stat)
	}
	if !approxEqual(pval, 1, 1e-6) {
		t.Errorf("Friedman equal: p=%v, want 1", pval)
	}
}

func TestFriedmanChiSquare_ClearDifference(t *testing.T) {
	// Group 2 always ranks higher
	g1 := []float64{1, 2, 3, 4, 5}
	g2 := []float64{10, 20, 30, 40, 50}
	stat, pval := FriedmanChiSquare(g1, g2)
	if stat <= 0 {
		t.Errorf("Friedman different: stat=%v, want positive", stat)
	}
	if pval > 0.05 {
		t.Errorf("Friedman different: p=%v, want < 0.05", pval)
	}
}

func TestFriedmanChiSquare_PanicDiffLen(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("FriedmanChiSquare should panic on different lengths")
		}
	}()
	FriedmanChiSquare([]float64{1, 2}, []float64{1, 2, 3})
}

// ---------------------------------------------------------------------------
// KSTest (one-sample)
// ---------------------------------------------------------------------------

func TestKSTest_PerfectNormal(t *testing.T) {
	// Generate quantiles from standard normal
	norm := NewNormal(0, 1)
	n := 100
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		p := (float64(i) + 0.5) / float64(n)
		x[i] = norm.PPF(p)
	}
	stat, pval := KSTest(x, norm.CDF)
	if stat > 0.1 {
		t.Errorf("KSTest normal quantiles: stat=%v, want < 0.1", stat)
	}
	if pval < 0.05 {
		t.Errorf("KSTest normal quantiles: p=%v, want > 0.05", pval)
	}
}

func TestKSTest_WrongDistribution(t *testing.T) {
	// Uniform data tested against normal CDF
	n := 50
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i) / float64(n)
	}
	norm := NewNormal(0, 1)
	stat, _ := KSTest(x, norm.CDF)
	if stat < 0.1 {
		t.Errorf("KSTest wrong dist: stat=%v, want > 0.1", stat)
	}
}

func TestKSTest_PanicEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("KSTest should panic on empty")
		}
	}()
	KSTest([]float64{}, func(x float64) float64 { return x })
}

// ---------------------------------------------------------------------------
// KS2Samp
// ---------------------------------------------------------------------------

func TestKS2Samp_SameDistribution(t *testing.T) {
	norm := NewNormal(0, 1)
	n := 100
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = norm.PPF((float64(i) + 0.5) / float64(n))
		y[i] = norm.PPF((float64(i) + 0.25) / float64(n))
	}
	_, pval := KS2Samp(x, y)
	if pval < 0.05 {
		t.Errorf("KS2Samp same dist: p=%v, want > 0.05", pval)
	}
}

func TestKS2Samp_DifferentDistribution(t *testing.T) {
	n := 50
	x := make([]float64, n)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i) / float64(n)       // Uniform [0,1]
		y[i] = float64(i)/float64(n)*10 + 10 // Uniform [10,20]
	}
	stat, _ := KS2Samp(x, y)
	if stat < 0.5 {
		t.Errorf("KS2Samp different: stat=%v, want > 0.5", stat)
	}
}

func TestKS2Samp_PanicEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("KS2Samp should panic on empty")
		}
	}()
	KS2Samp([]float64{}, []float64{1})
}

// ---------------------------------------------------------------------------
// ShapiroWilk
// ---------------------------------------------------------------------------

func TestShapiroWilk_NormalData(t *testing.T) {
	norm := NewNormal(0, 1)
	n := 50
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = norm.PPF((float64(i) + 0.5) / float64(n))
	}
	stat, pval := ShapiroWilk(x)
	if stat < 0.9 {
		t.Errorf("ShapiroWilk normal: stat=%v, want > 0.9", stat)
	}
	if pval < 0.05 {
		t.Errorf("ShapiroWilk normal: p=%v, want > 0.05", pval)
	}
}

func TestShapiroWilk_UniformData(t *testing.T) {
	// Uniform data: should be non-normal
	n := 50
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i) / float64(n-1)
	}
	stat, _ := ShapiroWilk(x)
	// W for uniform data should be notably less than 1
	if stat > 0.99 {
		t.Errorf("ShapiroWilk uniform: stat=%v, want < 0.99", stat)
	}
}

func TestShapiroWilk_PanicSmall(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("ShapiroWilk should panic with < 3 elements")
		}
	}()
	ShapiroWilk([]float64{1, 2})
}

func TestShapiroWilk_Constant(t *testing.T) {
	x := []float64{5, 5, 5, 5, 5}
	stat, pval := ShapiroWilk(x)
	if stat != 1 {
		t.Errorf("ShapiroWilk constant: stat=%v, want 1", stat)
	}
	if pval != 1 {
		t.Errorf("ShapiroWilk constant: p=%v, want 1", pval)
	}
}

// ---------------------------------------------------------------------------
// NormalTest
// ---------------------------------------------------------------------------

func TestNormalTest_NormalData(t *testing.T) {
	norm := NewNormal(0, 1)
	n := 200
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = norm.PPF((float64(i) + 0.5) / float64(n))
	}
	_, pval := NormalTest(x)
	if pval < 0.05 {
		t.Errorf("NormalTest normal data: p=%v, want > 0.05", pval)
	}
}

func TestNormalTest_SkewedData(t *testing.T) {
	// Exponential data: highly skewed
	n := 200
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		// Exponential quantiles: -ln(1 - p)
		p := (float64(i) + 0.5) / float64(n)
		x[i] = -math.Log(1 - p)
	}
	_, pval := NormalTest(x)
	if pval > 0.05 {
		t.Errorf("NormalTest skewed: p=%v, want < 0.05", pval)
	}
}

func TestNormalTest_PanicSmall(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NormalTest should panic with < 20 elements")
		}
	}()
	NormalTest(make([]float64, 19))
}

func TestNormalTest_ConstantData(t *testing.T) {
	x := make([]float64, 30)
	for i := range x {
		x[i] = 5
	}
	stat, pval := NormalTest(x)
	if stat != 0 {
		t.Errorf("NormalTest constant: stat=%v, want 0", stat)
	}
	if pval != 1 {
		t.Errorf("NormalTest constant: p=%v, want 1", pval)
	}
}

// ---------------------------------------------------------------------------
// AndersonDarling
// ---------------------------------------------------------------------------

func TestAndersonDarling_NormalData(t *testing.T) {
	norm := NewNormal(0, 1)
	n := 100
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = norm.PPF((float64(i) + 0.5) / float64(n))
	}
	stat, critical := AndersonDarling(x)
	if len(critical) != 5 {
		t.Fatalf("AndersonDarling: expected 5 critical values, got %d", len(critical))
	}
	// Normal data should pass at 5% level (stat < critical[2]=0.787)
	if stat > critical[2] {
		t.Errorf("AndersonDarling normal: stat=%v > critical[5%%]=%v", stat, critical[2])
	}
}

func TestAndersonDarling_NonNormalData(t *testing.T) {
	// Uniform data
	n := 100
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = float64(i)/float64(n-1)*6 - 3 // [-3, 3]
	}
	stat, critical := AndersonDarling(x)
	// Uniform data should fail normality test
	if stat < critical[4] {
		t.Errorf("AndersonDarling uniform: stat=%v < critical[1%%]=%v, expected to reject", stat, critical[4])
	}
}

func TestAndersonDarling_PanicSmall(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("AndersonDarling should panic with < 7 elements")
		}
	}()
	AndersonDarling([]float64{1, 2, 3, 4, 5, 6})
}

// ---------------------------------------------------------------------------
// BartlettTest
// ---------------------------------------------------------------------------

func TestBartlettTest_EqualVariances(t *testing.T) {
	g1 := []float64{1, 2, 3, 4, 5}
	g2 := []float64{6, 7, 8, 9, 10}
	g3 := []float64{11, 12, 13, 14, 15}
	stat, pval := BartlettTest(g1, g2, g3)
	if !approxEqual(stat, 0, 1e-8) {
		t.Errorf("Bartlett equal var: stat=%v, want ~0", stat)
	}
	if pval < 0.9 {
		t.Errorf("Bartlett equal var: p=%v, want > 0.9", pval)
	}
}

func TestBartlettTest_UnequalVariances(t *testing.T) {
	g1 := []float64{1, 2, 3, 4, 5}   // var=2.5
	g2 := []float64{1, 5, 9, 13, 17} // var=40
	_, pval := BartlettTest(g1, g2)
	if pval > 0.1 {
		t.Errorf("Bartlett unequal var: p=%v, want < 0.1", pval)
	}
}

func TestBartlettTest_PanicSmallGroup(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("BartlettTest should panic with group < 2 elements")
		}
	}()
	BartlettTest([]float64{1}, []float64{1, 2})
}

// ---------------------------------------------------------------------------
// LeveneTest
// ---------------------------------------------------------------------------

func TestLeveneTest_EqualVariances(t *testing.T) {
	g1 := []float64{1, 2, 3, 4, 5}
	g2 := []float64{6, 7, 8, 9, 10}
	_, pval := LeveneTest(g1, g2)
	if pval < 0.05 {
		t.Errorf("Levene equal var: p=%v, want > 0.05", pval)
	}
}

func TestLeveneTest_UnequalVariances(t *testing.T) {
	g1 := []float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5}
	g2 := []float64{-20, -10, 0, 10, 20, 30, 40, 50, 60, 70}
	_, pval := LeveneTest(g1, g2)
	if pval > 0.05 {
		t.Errorf("Levene unequal var: p=%v, want < 0.05", pval)
	}
}

func TestLeveneTest_PanicFewGroups(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("LeveneTest should panic with < 2 groups")
		}
	}()
	LeveneTest([]float64{1, 2, 3})
}

// ---------------------------------------------------------------------------
// FlignerKilleen
// ---------------------------------------------------------------------------

func TestFlignerKilleen_EqualVariances(t *testing.T) {
	g1 := []float64{1, 2, 3, 4, 5}
	g2 := []float64{6, 7, 8, 9, 10}
	_, pval := FlignerKilleen(g1, g2)
	if pval < 0.05 {
		t.Errorf("FlignerKilleen equal var: p=%v, want > 0.05", pval)
	}
}

func TestFlignerKilleen_UnequalVariances(t *testing.T) {
	g1 := []float64{1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5, 5.5}
	g2 := []float64{-20, -10, 0, 10, 20, 30, 40, 50, 60, 70}
	_, pval := FlignerKilleen(g1, g2)
	if pval > 0.05 {
		t.Errorf("FlignerKilleen unequal var: p=%v, want < 0.05", pval)
	}
}

func TestFlignerKilleen_PanicSmallGroup(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("FlignerKilleen should panic with group < 2")
		}
	}()
	FlignerKilleen([]float64{1}, []float64{1, 2})
}

// ---------------------------------------------------------------------------
// MoodTest
// ---------------------------------------------------------------------------

func TestMoodTest_EqualScale(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	y := []float64{11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	_, pval := MoodTest(x, y)
	// Same spread, just shifted
	if pval < 0.05 {
		t.Errorf("Mood equal scale: p=%v, want > 0.05", pval)
	}
}

func TestMoodTest_DifferentScale(t *testing.T) {
	x := []float64{4.9, 4.95, 5.0, 5.05, 5.1, 5.15, 5.2, 5.25, 5.3, 5.35}
	y := []float64{-100, -50, 0, 5, 10, 15, 20, 50, 100, 200}
	_, pval := MoodTest(x, y)
	if pval > 0.05 {
		t.Errorf("Mood different scale: p=%v, want < 0.05", pval)
	}
}

func TestMoodTest_PanicSmall(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MoodTest should panic with < 3 elements")
		}
	}()
	MoodTest([]float64{1, 2}, []float64{3, 4, 5})
}

// ---------------------------------------------------------------------------
// SpearmanR
// ---------------------------------------------------------------------------

func TestSpearmanR_PerfectMonotone(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2, 4, 8, 16, 32} // monotonically increasing (but not linear)
	r, pval := SpearmanR(x, y)
	if !approxEqual(r, 1.0, 1e-10) {
		t.Errorf("SpearmanR perfect monotone: r=%v, want 1", r)
	}
	if pval > 0.01 {
		t.Errorf("SpearmanR perfect monotone: p=%v, want < 0.01", pval)
	}
}

func TestSpearmanR_PerfectNegative(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{10, 8, 6, 4, 2}
	r, _ := SpearmanR(x, y)
	if !approxEqual(r, -1.0, 1e-10) {
		t.Errorf("SpearmanR perfect neg: r=%v, want -1", r)
	}
}

func TestSpearmanR_Uncorrelated(t *testing.T) {
	x := []float64{1, 2, 3, 4}
	y := []float64{2, 4, 1, 3}
	r, _ := SpearmanR(x, y)
	// Not perfectly correlated
	if math.Abs(r) > 0.9 {
		t.Errorf("SpearmanR mixed: r=%v, want closer to 0", r)
	}
}

func TestSpearmanR_PanicDiffLen(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("SpearmanR should panic on different lengths")
		}
	}()
	SpearmanR([]float64{1, 2, 3}, []float64{1, 2})
}

// ---------------------------------------------------------------------------
// KendallTau
// ---------------------------------------------------------------------------

func TestKendallTau_PerfectConcordance(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{1, 2, 3, 4, 5}
	tau, pval := KendallTau(x, y)
	if !approxEqual(tau, 1.0, 1e-10) {
		t.Errorf("KendallTau perfect: tau=%v, want 1", tau)
	}
	if pval > 0.05 {
		t.Errorf("KendallTau perfect: p=%v, want < 0.05", pval)
	}
}

func TestKendallTau_PerfectDiscordance(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{5, 4, 3, 2, 1}
	tau, _ := KendallTau(x, y)
	if !approxEqual(tau, -1.0, 1e-10) {
		t.Errorf("KendallTau perfect neg: tau=%v, want -1", tau)
	}
}

func TestKendallTau_PanicSmall(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("KendallTau should panic with < 3 elements")
		}
	}()
	KendallTau([]float64{1, 2}, []float64{1, 2})
}

func TestKendallTau_WithTies(t *testing.T) {
	x := []float64{1, 2, 2, 3, 4}
	y := []float64{1, 2, 2, 3, 4}
	tau, _ := KendallTau(x, y)
	if tau < 0.8 {
		t.Errorf("KendallTau ties: tau=%v, want > 0.8", tau)
	}
}

// ---------------------------------------------------------------------------
// Linregress
// ---------------------------------------------------------------------------

func TestLinregress_PerfectLine(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{3, 5, 7, 9, 11} // y = 2x + 1
	slope, intercept, r, pval, stderr := Linregress(x, y)
	if !approxEqual(slope, 2, 1e-10) {
		t.Errorf("Linregress perfect: slope=%v, want 2", slope)
	}
	if !approxEqual(intercept, 1, 1e-10) {
		t.Errorf("Linregress perfect: intercept=%v, want 1", intercept)
	}
	if !approxEqual(r, 1, 1e-10) {
		t.Errorf("Linregress perfect: r=%v, want 1", r)
	}
	if pval > 0.001 {
		t.Errorf("Linregress perfect: p=%v, want < 0.001", pval)
	}
	if !approxEqual(stderr, 0, 1e-10) {
		t.Errorf("Linregress perfect: stderr=%v, want 0", stderr)
	}
}

func TestLinregress_KnownValues(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{2.1, 3.9, 6.2, 7.8, 10.1}
	slope, intercept, r, pval, stderr := Linregress(x, y)
	// slope should be close to 2, intercept close to 0
	if math.Abs(slope-2) > 0.3 {
		t.Errorf("Linregress known: slope=%v, want ~2", slope)
	}
	if math.Abs(intercept) > 1 {
		t.Errorf("Linregress known: intercept=%v, want ~0", intercept)
	}
	if r < 0.99 {
		t.Errorf("Linregress known: r=%v, want > 0.99", r)
	}
	if pval > 0.01 {
		t.Errorf("Linregress known: p=%v, want < 0.01", pval)
	}
	if stderr <= 0 {
		t.Errorf("Linregress known: stderr=%v, want > 0", stderr)
	}
}

func TestLinregress_NegativeSlope(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5}
	y := []float64{10, 8, 6, 4, 2}
	slope, _, r, _, _ := Linregress(x, y)
	if !approxEqual(slope, -2, 1e-10) {
		t.Errorf("Linregress neg slope: slope=%v, want -2", slope)
	}
	if !approxEqual(r, -1, 1e-10) {
		t.Errorf("Linregress neg slope: r=%v, want -1", r)
	}
}

func TestLinregress_PanicSmall(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Linregress should panic with < 3 elements")
		}
	}()
	Linregress([]float64{1, 2}, []float64{3, 4})
}

func TestLinregress_ConstantX(t *testing.T) {
	x := []float64{5, 5, 5, 5}
	y := []float64{1, 2, 3, 4}
	slope, intercept, _, _, _ := Linregress(x, y)
	if slope != 0 {
		t.Errorf("Linregress constant x: slope=%v, want 0", slope)
	}
	if !approxEqual(intercept, 2.5, 1e-10) {
		t.Errorf("Linregress constant x: intercept=%v, want 2.5", intercept)
	}
}

// ---------------------------------------------------------------------------
// Entropy
// ---------------------------------------------------------------------------

func TestEntropy_Uniform(t *testing.T) {
	// Uniform over 4 outcomes: H = ln(4)
	pk := []float64{0.25, 0.25, 0.25, 0.25}
	h := Entropy(pk, nil)
	if !approxEqual(h, math.Log(4), 1e-10) {
		t.Errorf("Entropy uniform: H=%v, want ln(4)=%v", h, math.Log(4))
	}
}

func TestEntropy_Certain(t *testing.T) {
	// All mass on one outcome: H = 0
	pk := []float64{1, 0, 0, 0}
	h := Entropy(pk, nil)
	if !approxEqual(h, 0, 1e-10) {
		t.Errorf("Entropy certain: H=%v, want 0", h)
	}
}

func TestEntropy_KLDivergence(t *testing.T) {
	// KL(pk || qk) where pk=qk should be 0
	pk := []float64{0.25, 0.25, 0.25, 0.25}
	kl := Entropy(pk, pk)
	if !approxEqual(kl, 0, 1e-10) {
		t.Errorf("Entropy KL same: D=%v, want 0", kl)
	}
}

func TestEntropy_KLDivergencePositive(t *testing.T) {
	pk := []float64{0.5, 0.5}
	qk := []float64{0.25, 0.75}
	kl := Entropy(pk, qk)
	if kl <= 0 {
		t.Errorf("Entropy KL different: D=%v, want > 0", kl)
	}
}

func TestEntropy_KLDivergenceZeroQ(t *testing.T) {
	pk := []float64{0.5, 0.5}
	qk := []float64{1, 0}
	kl := Entropy(pk, qk)
	if !math.IsInf(kl, 1) {
		t.Errorf("Entropy KL zero q: D=%v, want +Inf", kl)
	}
}

func TestEntropy_PanicEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Entropy should panic on empty")
		}
	}()
	Entropy([]float64{}, nil)
}

func TestEntropy_PanicDiffLen(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Entropy should panic on different lengths")
		}
	}()
	Entropy([]float64{0.5, 0.5}, []float64{0.3, 0.3, 0.4})
}

// ---------------------------------------------------------------------------
// WassersteinDistance
// ---------------------------------------------------------------------------

func TestWassersteinDistance_Identical(t *testing.T) {
	u := []float64{1, 2, 3, 4, 5}
	v := []float64{1, 2, 3, 4, 5}
	d := WassersteinDistance(u, v)
	if !approxEqual(d, 0, 1e-10) {
		t.Errorf("Wasserstein identical: d=%v, want 0", d)
	}
}

func TestWassersteinDistance_Shifted(t *testing.T) {
	u := []float64{0, 1, 2, 3, 4}
	v := []float64{1, 2, 3, 4, 5}
	d := WassersteinDistance(u, v)
	// Shifted by 1, distance should be 1
	if !approxEqual(d, 1, 1e-10) {
		t.Errorf("Wasserstein shifted: d=%v, want 1", d)
	}
}

func TestWassersteinDistance_Positive(t *testing.T) {
	u := []float64{1, 2, 3}
	v := []float64{4, 5, 6}
	d := WassersteinDistance(u, v)
	if d <= 0 {
		t.Errorf("Wasserstein different: d=%v, want > 0", d)
	}
}

func TestWassersteinDistance_PanicEmpty(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("WassersteinDistance should panic on empty")
		}
	}()
	WassersteinDistance([]float64{}, []float64{1})
}

// ---------------------------------------------------------------------------
// Helper function tests
// ---------------------------------------------------------------------------

func TestMeanVar(t *testing.T) {
	x := []float64{2, 4, 4, 4, 5, 5, 7, 9}
	m, v := meanVar(x)
	if !approxEqual(m, 5.0, 1e-10) {
		t.Errorf("meanVar: mean=%v, want 5", m)
	}
	// Sample variance = 32/7 = 4.571...
	if !approxEqual(v, 32.0/7.0, 1e-10) {
		t.Errorf("meanVar: var=%v, want %v", v, 32.0/7.0)
	}
}

func TestMedian_Odd(t *testing.T) {
	x := []float64{3, 1, 2}
	m := median(x)
	if m != 2 {
		t.Errorf("median odd: %v, want 2", m)
	}
}

func TestMedian_Even(t *testing.T) {
	x := []float64{4, 1, 3, 2}
	m := median(x)
	if m != 2.5 {
		t.Errorf("median even: %v, want 2.5", m)
	}
}

func TestRankData(t *testing.T) {
	x := []float64{3, 1, 4, 1, 5}
	r := rankData(x)
	// Sorted: 1, 1, 3, 4, 5
	// Ranks: 1.5, 1.5, 3, 4, 5
	// Original positions: 3->3, 1->1.5, 4->4, 1->1.5, 5->5
	if !approxEqual(r[0], 3, 1e-10) { // 3 is rank 3
		t.Errorf("rankData[0]=%v, want 3", r[0])
	}
	if !approxEqual(r[1], 1.5, 1e-10) { // first 1 is rank 1.5
		t.Errorf("rankData[1]=%v, want 1.5", r[1])
	}
	if !approxEqual(r[2], 4, 1e-10) { // 4 is rank 4
		t.Errorf("rankData[2]=%v, want 4", r[2])
	}
	if !approxEqual(r[3], 1.5, 1e-10) { // second 1 is rank 1.5
		t.Errorf("rankData[3]=%v, want 1.5", r[3])
	}
	if !approxEqual(r[4], 5, 1e-10) { // 5 is rank 5
		t.Errorf("rankData[4]=%v, want 5", r[4])
	}
}

func TestFDistSurvival(t *testing.T) {
	// F(0) should give survival = 1
	p := fDistSurvival(0, 5, 10)
	if !approxEqual(p, 1, 1e-10) {
		t.Errorf("fDistSurvival(0, 5, 10)=%v, want 1", p)
	}

	// Very large F should give p ~ 0
	p = fDistSurvival(100, 5, 10)
	if p > 0.001 {
		t.Errorf("fDistSurvival(100, 5, 10)=%v, want < 0.001", p)
	}
}
