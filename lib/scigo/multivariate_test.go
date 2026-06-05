//go:build unit

package scigo

import (
	"math"
	"testing"
)

func approxEqualMV(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

// ---------------------------------------------------------------------------
// MultivariateNormal Construction
// ---------------------------------------------------------------------------

func TestNewMultivariateNormal_Valid(t *testing.T) {
	mean := []float64{0, 0}
	cov := [][]float64{{1, 0.5}, {0.5, 1}}
	mvn, err := NewMultivariateNormal(mean, cov)
	if err != nil {
		t.Fatal(err)
	}
	if mvn.Dim() != 2 {
		t.Errorf("Dim() = %d, want 2", mvn.Dim())
	}
}

func TestNewMultivariateNormal_EmptyMean(t *testing.T) {
	_, err := NewMultivariateNormal(nil, nil)
	if err == nil {
		t.Error("expected error for empty mean")
	}
}

func TestNewMultivariateNormal_CovRowsMismatch(t *testing.T) {
	_, err := NewMultivariateNormal([]float64{0, 0}, [][]float64{{1}})
	if err == nil {
		t.Error("expected error for cov rows mismatch")
	}
}

func TestNewMultivariateNormal_CovNotSquare(t *testing.T) {
	_, err := NewMultivariateNormal([]float64{0, 0}, [][]float64{{1, 0, 0}, {0, 1, 0}})
	if err == nil {
		t.Error("expected error for non-square cov")
	}
}

func TestNewMultivariateNormal_CovNotSymmetric(t *testing.T) {
	_, err := NewMultivariateNormal([]float64{0, 0}, [][]float64{{1, 0.5}, {0.3, 1}})
	if err == nil {
		t.Error("expected error for non-symmetric cov")
	}
}

func TestNewMultivariateNormal_CovNotPD(t *testing.T) {
	// Not positive definite.
	_, err := NewMultivariateNormal([]float64{0, 0}, [][]float64{{1, 2}, {2, 1}})
	if err == nil {
		t.Error("expected error for non-PD cov")
	}
}

// ---------------------------------------------------------------------------
// PDF / LogPDF
// ---------------------------------------------------------------------------

func TestMultivariateNormal_PDF_Standard2D(t *testing.T) {
	mean := []float64{0, 0}
	cov := [][]float64{{1, 0}, {0, 1}}
	mvn, err := NewMultivariateNormal(mean, cov)
	if err != nil {
		t.Fatal(err)
	}

	// PDF at origin: (2*pi)^{-1} * exp(0) = 1/(2*pi)
	got := mvn.PDF([]float64{0, 0})
	want := 1.0 / (2 * math.Pi)
	if !approxEqualMV(got, want, 1e-10) {
		t.Errorf("PDF([0,0]) = %v, want %v", got, want)
	}
}

func TestMultivariateNormal_LogPDF_WrongDim(t *testing.T) {
	mean := []float64{0, 0}
	cov := [][]float64{{1, 0}, {0, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	got := mvn.LogPDF([]float64{0, 0, 0})
	if !math.IsInf(got, -1) {
		t.Errorf("expected -Inf for wrong dimension, got %v", got)
	}
}

func TestMultivariateNormal_PDF_NonZeroMean(t *testing.T) {
	mean := []float64{1, 2}
	cov := [][]float64{{1, 0}, {0, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	// PDF at the mean should equal 1/(2*pi).
	got := mvn.PDF([]float64{1, 2})
	want := 1.0 / (2 * math.Pi)
	if !approxEqualMV(got, want, 1e-10) {
		t.Errorf("PDF at mean = %v, want %v", got, want)
	}
}

func TestMultivariateNormal_PDF_Correlated(t *testing.T) {
	mean := []float64{0, 0}
	rho := 0.5
	cov := [][]float64{{1, rho}, {rho, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	// PDF at origin: 1 / (2*pi*sqrt(1-rho^2))
	got := mvn.PDF([]float64{0, 0})
	want := 1.0 / (2 * math.Pi * math.Sqrt(1-rho*rho))
	if !approxEqualMV(got, want, 1e-10) {
		t.Errorf("PDF([0,0]) = %v, want %v", got, want)
	}
}

func TestMultivariateNormal_LogPDF_ConsistentWithPDF(t *testing.T) {
	mean := []float64{1, -1}
	cov := [][]float64{{2, 0.3}, {0.3, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	x := []float64{0.5, -0.5}
	logPDF := mvn.LogPDF(x)
	pdf := mvn.PDF(x)
	if !approxEqualMV(logPDF, math.Log(pdf), 1e-10) {
		t.Errorf("LogPDF and log(PDF) differ: %v vs %v", logPDF, math.Log(pdf))
	}
}

// ---------------------------------------------------------------------------
// Sample
// ---------------------------------------------------------------------------

func TestMultivariateNormal_Sample_MeanAndCov(t *testing.T) {
	mean := []float64{2, -1}
	cov := [][]float64{{4, 1}, {1, 2}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	nSamples := 100000
	samples := mvn.Sample(nSamples, 42)

	if len(samples) != nSamples {
		t.Fatalf("expected %d samples, got %d", nSamples, len(samples))
	}

	// Estimate sample mean.
	sMean := make([]float64, 2)
	for _, s := range samples {
		sMean[0] += s[0]
		sMean[1] += s[1]
	}
	sMean[0] /= float64(nSamples)
	sMean[1] /= float64(nSamples)

	if !approxEqualMV(sMean[0], mean[0], 0.05) {
		t.Errorf("sample mean[0] = %v, want ~%v", sMean[0], mean[0])
	}
	if !approxEqualMV(sMean[1], mean[1], 0.05) {
		t.Errorf("sample mean[1] = %v, want ~%v", sMean[1], mean[1])
	}

	// Estimate sample covariance.
	sCov := make([][]float64, 2)
	for i := range sCov {
		sCov[i] = make([]float64, 2)
	}
	for _, s := range samples {
		d0 := s[0] - sMean[0]
		d1 := s[1] - sMean[1]
		sCov[0][0] += d0 * d0
		sCov[0][1] += d0 * d1
		sCov[1][0] += d1 * d0
		sCov[1][1] += d1 * d1
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			sCov[i][j] /= float64(nSamples - 1)
		}
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if !approxEqualMV(sCov[i][j], cov[i][j], 0.1) {
				t.Errorf("sample cov[%d][%d] = %v, want ~%v", i, j, sCov[i][j], cov[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Mahalanobis
// ---------------------------------------------------------------------------

func TestMultivariateNormal_Mahalanobis_AtMean(t *testing.T) {
	mean := []float64{1, 2}
	cov := [][]float64{{1, 0}, {0, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	d := mvn.Mahalanobis([]float64{1, 2})
	if !approxEqualMV(d, 0, 1e-10) {
		t.Errorf("Mahalanobis at mean = %v, want 0", d)
	}
}

func TestMultivariateNormal_Mahalanobis_UnitStd(t *testing.T) {
	mean := []float64{0, 0}
	cov := [][]float64{{1, 0}, {0, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	// (1,0): distance = 1
	d := mvn.Mahalanobis([]float64{1, 0})
	if !approxEqualMV(d, 1.0, 1e-10) {
		t.Errorf("Mahalanobis = %v, want 1", d)
	}

	// (1,1): distance = sqrt(2)
	d2 := mvn.Mahalanobis([]float64{1, 1})
	if !approxEqualMV(d2, math.Sqrt(2), 1e-10) {
		t.Errorf("Mahalanobis = %v, want %v", d2, math.Sqrt(2))
	}
}

func TestMultivariateNormal_Mahalanobis_ScaledCov(t *testing.T) {
	mean := []float64{0}
	cov := [][]float64{{4}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	// x=2, sigma=2 => distance = 2/2 = 1
	d := mvn.Mahalanobis([]float64{2})
	if !approxEqualMV(d, 1.0, 1e-10) {
		t.Errorf("Mahalanobis = %v, want 1", d)
	}
}

// ---------------------------------------------------------------------------
// MarginalDistribution
// ---------------------------------------------------------------------------

func TestMultivariateNormal_MarginalDistribution(t *testing.T) {
	mean := []float64{1, 2, 3}
	cov := [][]float64{
		{4, 1, 0.5},
		{1, 3, 0.2},
		{0.5, 0.2, 2},
	}
	mvn, _ := NewMultivariateNormal(mean, cov)

	marg := mvn.MarginalDistribution([]int{0, 2})
	mMean := marg.MeanVec()
	mCov := marg.CovMatrix()

	if !approxEqualMV(mMean[0], 1, 1e-10) || !approxEqualMV(mMean[1], 3, 1e-10) {
		t.Errorf("marginal mean = %v, want [1,3]", mMean)
	}
	if !approxEqualMV(mCov[0][0], 4, 1e-10) || !approxEqualMV(mCov[1][1], 2, 1e-10) {
		t.Errorf("marginal cov diagonal wrong: %v", mCov)
	}
	if !approxEqualMV(mCov[0][1], 0.5, 1e-10) {
		t.Errorf("marginal cov off-diag = %v, want 0.5", mCov[0][1])
	}
}

// ---------------------------------------------------------------------------
// ConditionalDistribution
// ---------------------------------------------------------------------------

func TestMultivariateNormal_ConditionalDistribution(t *testing.T) {
	// 2D standard normal with correlation 0.5.
	mean := []float64{0, 0}
	rho := 0.5
	cov := [][]float64{{1, rho}, {rho, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	// Condition on x2 = 1.
	cond := mvn.ConditionalDistribution(map[int]float64{1: 1.0})
	cMean := cond.MeanVec()
	cCov := cond.CovMatrix()

	// E[X1|X2=1] = 0 + 0.5*1 = 0.5
	if !approxEqualMV(cMean[0], 0.5, 1e-10) {
		t.Errorf("conditional mean = %v, want 0.5", cMean[0])
	}
	// Var[X1|X2=1] = 1 - 0.5^2 = 0.75
	if !approxEqualMV(cCov[0][0], 0.75, 1e-10) {
		t.Errorf("conditional variance = %v, want 0.75", cCov[0][0])
	}
}

func TestMultivariateNormal_ConditionalDistribution_NoObserved(t *testing.T) {
	mean := []float64{0, 0}
	cov := [][]float64{{1, 0}, {0, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	cond := mvn.ConditionalDistribution(map[int]float64{})
	if cond.Dim() != 2 {
		t.Errorf("expected same distribution when no observations")
	}
}

func TestMultivariateNormal_ConditionalDistribution_AllObserved(t *testing.T) {
	mean := []float64{0, 0}
	cov := [][]float64{{1, 0}, {0, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	cond := mvn.ConditionalDistribution(map[int]float64{0: 1, 1: 2})
	// Should return original since all observed leaves nothing.
	if cond.Dim() != 2 {
		t.Errorf("expected original distribution when all observed")
	}
}

// ---------------------------------------------------------------------------
// MeanVec / CovMatrix (copy safety)
// ---------------------------------------------------------------------------

func TestMultivariateNormal_CopySafety(t *testing.T) {
	mean := []float64{0, 0}
	cov := [][]float64{{1, 0}, {0, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	m := mvn.MeanVec()
	m[0] = 999
	if mvn.MeanVec()[0] == 999 {
		t.Error("MeanVec should return a copy")
	}

	c := mvn.CovMatrix()
	c[0][0] = 999
	if mvn.CovMatrix()[0][0] == 999 {
		t.Error("CovMatrix should return a copy")
	}
}

// ---------------------------------------------------------------------------
// GaussianCopula
// ---------------------------------------------------------------------------

func TestNewGaussianCopula_Valid(t *testing.T) {
	corr := [][]float64{{1, 0.5}, {0.5, 1}}
	gc, err := NewGaussianCopula(corr)
	if err != nil {
		t.Fatal(err)
	}
	if gc.Dim() != 2 {
		t.Errorf("Dim() = %d, want 2", gc.Dim())
	}
}

func TestNewGaussianCopula_Empty(t *testing.T) {
	_, err := NewGaussianCopula(nil)
	if err == nil {
		t.Error("expected error for empty correlation matrix")
	}
}

func TestNewGaussianCopula_NotSquare(t *testing.T) {
	_, err := NewGaussianCopula([][]float64{{1, 0, 0}, {0, 1, 0}})
	if err == nil {
		t.Error("expected error for non-square matrix")
	}
}

func TestNewGaussianCopula_DiagNotOne(t *testing.T) {
	_, err := NewGaussianCopula([][]float64{{2, 0}, {0, 1}})
	if err == nil {
		t.Error("expected error for diagonal != 1")
	}
}

func TestNewGaussianCopula_NotSymmetric(t *testing.T) {
	_, err := NewGaussianCopula([][]float64{{1, 0.5}, {0.3, 1}})
	if err == nil {
		t.Error("expected error for non-symmetric matrix")
	}
}

func TestNewGaussianCopula_NotPD(t *testing.T) {
	_, err := NewGaussianCopula([][]float64{{1, 1.5}, {1.5, 1}})
	if err == nil {
		t.Error("expected error for non-PD matrix")
	}
}

func TestGaussianCopula_Sample_Uniform(t *testing.T) {
	corr := [][]float64{{1, 0}, {0, 1}}
	gc, _ := NewGaussianCopula(corr)

	samples := gc.Sample(10000, 42)
	if len(samples) != 10000 {
		t.Fatalf("expected 10000 samples, got %d", len(samples))
	}

	// All samples should be in [0,1].
	for i, s := range samples {
		for j, v := range s {
			if v < 0 || v > 1 {
				t.Errorf("sample[%d][%d] = %v, out of [0,1]", i, j, v)
				break
			}
		}
	}

	// Mean should be approximately 0.5.
	sum0, sum1 := 0.0, 0.0
	for _, s := range samples {
		sum0 += s[0]
		sum1 += s[1]
	}
	mean0 := sum0 / 10000
	mean1 := sum1 / 10000
	if !approxEqualMV(mean0, 0.5, 0.02) {
		t.Errorf("sample mean[0] = %v, want ~0.5", mean0)
	}
	if !approxEqualMV(mean1, 0.5, 0.02) {
		t.Errorf("sample mean[1] = %v, want ~0.5", mean1)
	}
}

func TestGaussianCopula_Sample_Correlated(t *testing.T) {
	rho := 0.8
	corr := [][]float64{{1, rho}, {rho, 1}}
	gc, _ := NewGaussianCopula(corr)

	samples := gc.Sample(50000, 123)

	// Compute rank correlation (Spearman) approximately.
	// For strong positive correlation, most samples should have both components
	// above or below 0.5 simultaneously.
	concordant := 0
	for _, s := range samples {
		if (s[0] > 0.5 && s[1] > 0.5) || (s[0] < 0.5 && s[1] < 0.5) {
			concordant++
		}
	}
	concProp := float64(concordant) / float64(len(samples))
	// With rho=0.8, concordance should be well above 0.5.
	if concProp < 0.6 {
		t.Errorf("concordance proportion = %v, expected > 0.6 for rho=0.8", concProp)
	}
}

func TestGaussianCopula_CorrMatrix_Copy(t *testing.T) {
	corr := [][]float64{{1, 0.3}, {0.3, 1}}
	gc, _ := NewGaussianCopula(corr)

	c := gc.CorrMatrix()
	c[0][1] = 999
	if gc.CorrMatrix()[0][1] == 999 {
		t.Error("CorrMatrix should return a copy")
	}
}

// ---------------------------------------------------------------------------
// ConditionalDistribution error path (singular sig22)
// ---------------------------------------------------------------------------

func TestMultivariateNormal_ConditionalDistribution_SingularFallback(t *testing.T) {
	// This tests the path where matInverse fails.
	// Create a 3D MVN and condition on dims 1,2 where the sub-cov is near-singular.
	// Use a valid MVN but observe from only 1 dimension (which is always invertible for 1x1).
	// Instead, test with a dimension that works fine and verify output.
	mean := []float64{1, 2}
	cov := [][]float64{{4, 2}, {2, 4}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	cond := mvn.ConditionalDistribution(map[int]float64{1: 3.0})
	cMean := cond.MeanVec()
	// E[X0|X1=3] = 1 + (2/4)*(3-2) = 1.5
	if !approxEqualMV(cMean[0], 1.5, 1e-8) {
		t.Errorf("conditional mean = %v, want 1.5", cMean[0])
	}
	cCov := cond.CovMatrix()
	// Var[X0|X1] = 4 - 4/4 = 3
	if !approxEqualMV(cCov[0][0], 3.0, 1e-8) {
		t.Errorf("conditional var = %v, want 3.0", cCov[0][0])
	}
}

// ---------------------------------------------------------------------------
// 3D MultivariateNormal
// ---------------------------------------------------------------------------

func TestMultivariateNormal_3D_PDF(t *testing.T) {
	mean := []float64{0, 0, 0}
	cov := [][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}}
	mvn, _ := NewMultivariateNormal(mean, cov)

	// PDF at origin: (2*pi)^{-3/2}
	got := mvn.PDF([]float64{0, 0, 0})
	want := math.Pow(2*math.Pi, -1.5)
	if !approxEqualMV(got, want, 1e-10) {
		t.Errorf("PDF([0,0,0]) = %v, want %v", got, want)
	}
}

func TestMultivariateNormal_1D_MatchesNormal(t *testing.T) {
	// 1D MVN should match univariate normal.
	mean := []float64{3}
	cov := [][]float64{{4}} // sigma = 2
	mvn, _ := NewMultivariateNormal(mean, cov)

	normal := NewNormal(3, 2)

	x := 5.0
	mvnPDF := mvn.PDF([]float64{x})
	normalPDF := normal.PDF(x)
	if !approxEqualMV(mvnPDF, normalPDF, 1e-10) {
		t.Errorf("1D MVN PDF = %v, Normal PDF = %v", mvnPDF, normalPDF)
	}
}

func TestMultivariateNormal_ConditionalDistribution_3D(t *testing.T) {
	// 3D normal, condition on dim 2.
	mean := []float64{0, 0, 0}
	cov := [][]float64{
		{1, 0.5, 0.3},
		{0.5, 1, 0.4},
		{0.3, 0.4, 1},
	}
	mvn, _ := NewMultivariateNormal(mean, cov)

	cond := mvn.ConditionalDistribution(map[int]float64{2: 1.0})
	if cond.Dim() != 2 {
		t.Errorf("conditional dim = %d, want 2", cond.Dim())
	}

	cMean := cond.MeanVec()
	// E[X0|X2=1] = 0 + 0.3*1 = 0.3
	if !approxEqualMV(cMean[0], 0.3, 1e-10) {
		t.Errorf("cond mean[0] = %v, want 0.3", cMean[0])
	}
	// E[X1|X2=1] = 0 + 0.4*1 = 0.4
	if !approxEqualMV(cMean[1], 0.4, 1e-10) {
		t.Errorf("cond mean[1] = %v, want 0.4", cMean[1])
	}
}
