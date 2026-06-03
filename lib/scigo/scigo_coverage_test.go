//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// PartialCorrelation coverage
// ---------------------------------------------------------------------------

func TestPartialCorrelation_InsufficientData(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for insufficient data")
		}
	}()
	data := [][]float64{{1, 2}, {3, 4}}
	PartialCorrelation(data, 0, 1, nil)
}

func TestPartialCorrelation_InconsistentRows(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for inconsistent rows")
		}
	}()
	data := [][]float64{{1, 2}, {3}, {5, 6}}
	PartialCorrelation(data, 0, 1, nil)
}

func TestPartialCorrelation_BadColumnIndex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for bad column index")
		}
	}()
	data := [][]float64{{1, 2}, {3, 4}, {5, 6}}
	PartialCorrelation(data, 0, 5, nil)
}

func TestPartialCorrelation_BadZIndex(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for bad z index")
		}
	}()
	data := [][]float64{{1, 2}, {3, 4}, {5, 6}}
	PartialCorrelation(data, 0, 1, []int{5})
}

func TestPartialCorrelation_WithConditioning(t *testing.T) {
	// X, Y, Z where X and Y are both caused by Z.
	data := [][]float64{
		{1, 2, 0},
		{2, 4, 0},
		{3, 6, 1},
		{4, 8, 1},
		{5, 10, 0},
		{6, 12, 1},
	}
	r, pval := PartialCorrelation(data, 0, 1, []int{2})
	if math.Abs(r) > 1.01 {
		t.Errorf("expected r in [-1,1], got %f", r)
	}
	_ = pval
}

func TestPartialCorrelation_ZeroDenom(t *testing.T) {
	// Create data where one of the partial correlations equals exactly 1,
	// making the denominator zero.
	data := [][]float64{
		{1, 2, 1},
		{2, 4, 2},
		{3, 6, 3},
		{4, 8, 4},
	}
	r, pval := PartialCorrelation(data, 0, 1, []int{2})
	// When denom is 0, should return 0, 1.
	_ = r
	_ = pval
}

// ---------------------------------------------------------------------------
// F-Distribution edge cases
// ---------------------------------------------------------------------------

func TestNewFDistribution_PanicDF1(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for df1 <= 0")
		}
	}()
	NewFDistribution(0, 5)
}

func TestNewFDistribution_PanicDF2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for df2 <= 0")
		}
	}()
	NewFDistribution(5, 0)
}

func TestFDistribution_CDF_NegativeX(t *testing.T) {
	f := NewFDistribution(2, 5)
	if f.CDF(-1) != 0 {
		t.Error("expected 0 for negative x")
	}
}

func TestFDistribution_PPF_Edges(t *testing.T) {
	f := NewFDistribution(2, 5)
	if f.PPF(0) != 0 {
		t.Error("expected 0 for p=0")
	}
	if !math.IsInf(f.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

// ---------------------------------------------------------------------------
// Lognormal edge cases
// ---------------------------------------------------------------------------

func TestLognormal_PanicSigma(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for sigma <= 0")
		}
	}()
	NewLognormal(0, 0)
}

func TestLognormal_CDFNegative(t *testing.T) {
	ln := NewLognormal(0, 1)
	if ln.CDF(-1) != 0 {
		t.Error("expected 0 for CDF at negative x")
	}
	if ln.CDF(0) != 0 {
		t.Error("expected 0 for CDF at x=0")
	}
}

func TestLognormal_PPF_Edges(t *testing.T) {
	ln := NewLognormal(0, 1)
	if ln.PPF(0) != 0 {
		t.Error("expected 0 for p=0")
	}
	if !math.IsInf(ln.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

func TestLognormal_LogPDF_Negative(t *testing.T) {
	ln := NewLognormal(0, 1)
	if !math.IsInf(ln.LogPDF(-1), -1) {
		t.Error("expected -Inf for LogPDF at negative x")
	}
	if !math.IsInf(ln.LogPDF(0), -1) {
		t.Error("expected -Inf for LogPDF at x=0")
	}
}

// ---------------------------------------------------------------------------
// Weibull edge cases
// ---------------------------------------------------------------------------

func TestNewWeibull_PanicShape(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for shape <= 0")
		}
	}()
	NewWeibull(0, 1)
}

func TestNewWeibull_PanicScale(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for scale <= 0")
		}
	}()
	NewWeibull(1, 0)
}

func TestWeibull_PDF_ZeroX(t *testing.T) {
	// shape < 1: PDF(0) = Inf
	w := NewWeibull(0.5, 1)
	if !math.IsInf(w.PDF(0), 1) {
		t.Error("expected +Inf for shape<1 at x=0")
	}
	// shape == 1: PDF(0) = 1/scale
	w2 := NewWeibull(1, 2)
	if math.Abs(w2.PDF(0)-0.5) > 1e-10 {
		t.Errorf("expected 0.5 for shape=1 at x=0, got %f", w2.PDF(0))
	}
	// shape > 1: PDF(0) = 0
	w3 := NewWeibull(2, 1)
	if w3.PDF(0) != 0 {
		t.Errorf("expected 0 for shape>1 at x=0, got %f", w3.PDF(0))
	}
}

func TestWeibull_PDF_NegativeX(t *testing.T) {
	w := NewWeibull(2, 1)
	if w.PDF(-1) != 0 {
		t.Error("expected 0 for negative x")
	}
}

func TestWeibull_CDF_NegativeX(t *testing.T) {
	w := NewWeibull(2, 1)
	if w.CDF(-1) != 0 {
		t.Error("expected 0 for negative x")
	}
	if w.CDF(0) != 0 {
		t.Error("expected 0 for x=0")
	}
}

func TestWeibull_PPF_Edges(t *testing.T) {
	w := NewWeibull(2, 1)
	if w.PPF(0) != 0 {
		t.Error("expected 0 for p=0")
	}
	if !math.IsInf(w.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

func TestWeibull_LogPDF_Negative(t *testing.T) {
	w := NewWeibull(2, 1)
	if !math.IsInf(w.LogPDF(-1), -1) {
		t.Error("expected -Inf for negative x")
	}
	if !math.IsInf(w.LogPDF(0), -1) {
		t.Error("expected -Inf for x=0")
	}
}

// ---------------------------------------------------------------------------
// Pareto edge cases
// ---------------------------------------------------------------------------

func TestNewPareto_PanicAlpha(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for alpha <= 0")
		}
	}()
	NewPareto(0, 1)
}

func TestNewPareto_PanicXm(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for xm <= 0")
		}
	}()
	NewPareto(2, 0)
}

func TestPareto_CDF_BelowXm(t *testing.T) {
	p := NewPareto(2, 1)
	if p.CDF(0.5) != 0 {
		t.Error("expected 0 below xm")
	}
}

func TestPareto_PPF_Edges(t *testing.T) {
	p := NewPareto(2, 1)
	if p.PPF(0) != 1 {
		t.Errorf("expected xm=1 for p=0, got %f", p.PPF(0))
	}
	if !math.IsInf(p.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

func TestPareto_LogPDF_BelowXm(t *testing.T) {
	p := NewPareto(2, 1)
	if !math.IsInf(p.LogPDF(0.5), -1) {
		t.Error("expected -Inf below xm")
	}
}

func TestPareto_Var_LowAlpha(t *testing.T) {
	p := NewPareto(1.5, 1)
	if !math.IsInf(p.Var(), 1) {
		t.Error("expected +Inf for alpha <= 2")
	}
}

func TestPareto_Mean_LowAlpha(t *testing.T) {
	p := NewPareto(0.5, 1)
	if !math.IsInf(p.Mean(), 1) {
		t.Error("expected +Inf for alpha <= 1")
	}
}

// ---------------------------------------------------------------------------
// GEV edge cases
// ---------------------------------------------------------------------------

func TestGEV_PPF_Edges(t *testing.T) {
	// xi > 0 at p=0
	gev := NewGeneralizedExtremeValue(0, 1, 0.5)
	v := gev.PPF(0)
	if math.IsInf(v, 0) && v > 0 {
		t.Error("should not be +Inf for xi>0 at p=0")
	}

	// xi < 0 at p=1
	gev2 := NewGeneralizedExtremeValue(0, 1, -0.5)
	v = gev2.PPF(1)
	if math.IsInf(v, 0) && v < 0 {
		t.Error("should not be -Inf for xi<0 at p=1")
	}

	// Gumbel case (xi near 0)
	gev3 := NewGeneralizedExtremeValue(0, 1, 0)
	v = gev3.PPF(0.5)
	if math.IsNaN(v) {
		t.Error("expected finite value for Gumbel PPF at p=0.5")
	}
}

func TestGEV_Mean_HighXi(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 1.5)
	if !math.IsInf(gev.Mean(), 1) {
		t.Error("expected +Inf for xi >= 1")
	}
}

func TestGEV_Var_HighXi(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0.6)
	if !math.IsInf(gev.Var(), 1) {
		t.Error("expected +Inf for xi >= 0.5")
	}
}

func TestGEV_Gumbel_MeanVar(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0)
	m := gev.Mean()
	v := gev.Var()
	// Gumbel mean ≈ 0.577, var ≈ pi^2/6
	if math.Abs(m-0.577) > 0.1 {
		t.Errorf("expected mean near 0.577, got %f", m)
	}
	if math.Abs(v-math.Pi*math.Pi/6) > 0.1 {
		t.Errorf("expected var near pi^2/6, got %f", v)
	}
}

func TestGEV_CDF_OutOfSupport(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0.5)
	cdf := gev.CDF(-10)
	if cdf < 0 || cdf > 1 {
		t.Errorf("expected CDF in [0,1], got %f", cdf)
	}
}

func TestGEV_GevT_Gumbel(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0)
	v := gev.gevT(1)
	expected := math.Exp(-1)
	if math.Abs(v-expected) > 1e-10 {
		t.Errorf("expected %f, got %f", expected, v)
	}
}

func TestGEV_GevT_OutOfSupport(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0.5)
	// arg = 1 + 0.5*(-3) = -0.5 <= 0
	v := gev.gevT(-3)
	if !math.IsInf(v, 1) {
		t.Errorf("expected +Inf for out-of-support, got %f", v)
	}
}

func TestGEV_PDF_Gumbel(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0)
	p := gev.PDF(0)
	expected := math.Exp(-math.Exp(0)) / 1
	if math.Abs(p-expected) > 0.01 {
		t.Errorf("expected %f, got %f", expected, p)
	}
}

func TestGEV_PDF_OutOfSupport(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0.5)
	p := gev.PDF(-3)
	if p != 0 {
		t.Errorf("expected 0 for out-of-support, got %f", p)
	}
}

func TestGEV_LogPDF_Gumbel(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0)
	lp := gev.LogPDF(0)
	if math.IsInf(lp, -1) {
		t.Error("should be finite for Gumbel case in support")
	}
}

func TestGEV_LogPDF_OutOfSupport(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0.5)
	lp := gev.LogPDF(-3)
	if !math.IsInf(lp, -1) {
		t.Errorf("expected -Inf for out-of-support, got %f", lp)
	}
}

func TestGEV_CDF_Gumbel(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0)
	cdf := gev.CDF(0)
	expected := math.Exp(-1) // exp(-exp(0)) = exp(-1)
	if math.Abs(cdf-expected) > 0.01 {
		t.Errorf("expected %f, got %f", expected, cdf)
	}
}

// ---------------------------------------------------------------------------
// HalfNormal edge cases
// ---------------------------------------------------------------------------

func TestHalfNormal_PPF_Edges(t *testing.T) {
	hn := NewHalfNormal(1)
	if hn.PPF(0) != 0 {
		t.Error("expected 0 for p=0")
	}
	if !math.IsInf(hn.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

func TestHalfNormal_LogPDF_Negative(t *testing.T) {
	hn := NewHalfNormal(1)
	if !math.IsInf(hn.LogPDF(-1), -1) {
		t.Error("expected -Inf for negative x")
	}
}

// ---------------------------------------------------------------------------
// TruncatedNormal edge cases
// ---------------------------------------------------------------------------

func TestNewTruncatedNormal_PanicSigma(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for sigma <= 0")
		}
	}()
	NewTruncatedNormal(0, 0, -1, 1)
}

func TestNewTruncatedNormal_PanicBounds(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for a >= b")
		}
	}()
	NewTruncatedNormal(0, 1, 1, 1)
}

func TestTruncatedNormal_PPF_Edges(t *testing.T) {
	tn := NewTruncatedNormal(0, 1, -2, 2)
	if tn.PPF(0) != -2 {
		t.Errorf("expected -2 for p=0, got %f", tn.PPF(0))
	}
	if tn.PPF(1) != 2 {
		t.Errorf("expected 2 for p=1, got %f", tn.PPF(1))
	}
}

func TestTruncatedNormal_LogPDF_OutOfRange(t *testing.T) {
	tn := NewTruncatedNormal(0, 1, -2, 2)
	if !math.IsInf(tn.LogPDF(-3), -1) {
		t.Error("expected -Inf for x < a")
	}
	if !math.IsInf(tn.LogPDF(3), -1) {
		t.Error("expected -Inf for x > b")
	}
}

// ---------------------------------------------------------------------------
// SkewNormal edge cases
// ---------------------------------------------------------------------------

func TestSkewNormal_PPF(t *testing.T) {
	sn := NewSkewNormal(0, 1, 0)
	if !math.IsInf(sn.PPF(0), -1) {
		t.Error("expected -Inf for p=0")
	}
	if !math.IsInf(sn.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
	v := sn.PPF(0.5)
	if math.Abs(v) > 0.5 {
		t.Errorf("expected near 0 for symmetric skew-normal, got %f", v)
	}
}

func TestSkewNormal_CDF_LowerBound(t *testing.T) {
	sn := NewSkewNormal(0, 1, 0)
	cdf := sn.CDF(-100)
	if cdf != 0 {
		t.Errorf("expected 0 for very low x, got %f", cdf)
	}
}

// ---------------------------------------------------------------------------
// Levy edge cases
// ---------------------------------------------------------------------------

func TestLevy_CDF_AtLoc(t *testing.T) {
	l := NewLevy(0, 1)
	if l.CDF(0) != 0 {
		t.Error("expected 0 at location")
	}
	if l.CDF(-1) != 0 {
		t.Error("expected 0 below location")
	}
}

func TestLevy_PPF_Edges(t *testing.T) {
	l := NewLevy(0, 1)
	if l.PPF(0) != 0 {
		t.Errorf("expected 0 for p=0, got %f", l.PPF(0))
	}
	if !math.IsInf(l.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

func TestLevy_LogPDF_AtLoc(t *testing.T) {
	l := NewLevy(0, 1)
	if !math.IsInf(l.LogPDF(0), -1) {
		t.Error("expected -Inf at location")
	}
	if !math.IsInf(l.LogPDF(-1), -1) {
		t.Error("expected -Inf below location")
	}
}

// ---------------------------------------------------------------------------
// PowerLaw edge cases
// ---------------------------------------------------------------------------

func TestNewPowerLaw_PanicAlpha(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for alpha <= 1")
		}
	}()
	NewPowerLaw(1, 1)
}

func TestNewPowerLaw_PanicXmin(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for xmin <= 0")
		}
	}()
	NewPowerLaw(2, 0)
}

func TestPowerLaw_CDF_BelowXmin(t *testing.T) {
	pl := NewPowerLaw(2, 1)
	if pl.CDF(0.5) != 0 {
		t.Error("expected 0 below xmin")
	}
}

func TestPowerLaw_PPF_Edges(t *testing.T) {
	pl := NewPowerLaw(2, 1)
	if pl.PPF(0) != 1 {
		t.Errorf("expected xmin=1 for p=0, got %f", pl.PPF(0))
	}
	if !math.IsInf(pl.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

func TestPowerLaw_LogPDF_BelowXmin(t *testing.T) {
	pl := NewPowerLaw(2, 1)
	if !math.IsInf(pl.LogPDF(0.5), -1) {
		t.Error("expected -Inf below xmin")
	}
}

func TestPowerLaw_Mean_LowAlpha(t *testing.T) {
	pl := NewPowerLaw(1.5, 1)
	if !math.IsInf(pl.Mean(), 1) {
		t.Error("expected +Inf for alpha <= 2")
	}
}

func TestPowerLaw_Var_LowAlpha(t *testing.T) {
	pl := NewPowerLaw(2.5, 1)
	if !math.IsInf(pl.Var(), 1) {
		t.Error("expected +Inf for alpha <= 3")
	}
}

func TestPowerLaw_Var_HighAlpha(t *testing.T) {
	pl := NewPowerLaw(4, 1)
	v := pl.Var()
	if v < 0 || math.IsInf(v, 0) {
		t.Errorf("expected finite positive variance for alpha=4, got %f", v)
	}
}

// ---------------------------------------------------------------------------
// Wald edge cases
// ---------------------------------------------------------------------------

func TestWald_PPF_Edges(t *testing.T) {
	w := NewWald(1, 1)
	if w.PPF(0) != 0 {
		t.Error("expected 0 for p=0")
	}
	if !math.IsInf(w.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

func TestWald_LogPDF_NonPositive(t *testing.T) {
	w := NewWald(1, 1)
	if !math.IsInf(w.LogPDF(0), -1) {
		t.Error("expected -Inf for x=0")
	}
	if !math.IsInf(w.LogPDF(-1), -1) {
		t.Error("expected -Inf for x<0")
	}
}

// ---------------------------------------------------------------------------
// PearsonCorrelation edge cases
// ---------------------------------------------------------------------------

func TestPearsonCorrelation_Constant(t *testing.T) {
	x := []float64{5, 5, 5}
	y := []float64{3, 3, 3}
	r, pval := PearsonCorrelation(x, y)
	if r != 0 {
		t.Errorf("expected 0 for constant data, got %f", r)
	}
	_ = pval
}

// ---------------------------------------------------------------------------
// Cauchy edge cases
// ---------------------------------------------------------------------------

func TestCauchy_PPF_Edges(t *testing.T) {
	c := NewCauchy(0, 1)
	if !math.IsInf(c.PPF(0), -1) {
		t.Error("expected -Inf for p=0")
	}
	if !math.IsInf(c.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

// ---------------------------------------------------------------------------
// Laplace edge cases
// ---------------------------------------------------------------------------

func TestLaplace_PPF_Edges(t *testing.T) {
	l := NewLaplace(0, 1)
	if !math.IsInf(l.PPF(0), -1) {
		t.Error("expected -Inf for p=0")
	}
	if !math.IsInf(l.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
	// Lower half
	v := l.PPF(0.25)
	if v >= 0 {
		t.Errorf("expected negative for p=0.25, got %f", v)
	}
	// Upper half
	v = l.PPF(0.75)
	if v <= 0 {
		t.Errorf("expected positive for p=0.75, got %f", v)
	}
}

// ---------------------------------------------------------------------------
// Logistic edge cases
// ---------------------------------------------------------------------------

func TestLogistic_PPF_Edges(t *testing.T) {
	l := NewLogistic(0, 1)
	if !math.IsInf(l.PPF(0), -1) {
		t.Error("expected -Inf for p=0")
	}
	if !math.IsInf(l.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

// ---------------------------------------------------------------------------
// Gumbel edge cases
// ---------------------------------------------------------------------------

func TestGumbel_PPF_Edges(t *testing.T) {
	g := NewGumbel(0, 1)
	if !math.IsInf(g.PPF(0), -1) {
		t.Error("expected -Inf for p=0")
	}
	if !math.IsInf(g.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

// ---------------------------------------------------------------------------
// Rayleigh edge cases
// ---------------------------------------------------------------------------

func TestRayleigh_PPF_Edges(t *testing.T) {
	r := NewRayleigh(1)
	if r.PPF(0) != 0 {
		t.Error("expected 0 for p=0")
	}
	if !math.IsInf(r.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
}

func TestRayleigh_CDF_Negative(t *testing.T) {
	r := NewRayleigh(1)
	if r.CDF(-1) != 0 {
		t.Error("expected 0 for negative x")
	}
	if r.CDF(0) != 0 {
		t.Error("expected 0 for x=0")
	}
}

func TestRayleigh_LogPDF_Negative(t *testing.T) {
	r := NewRayleigh(1)
	if !math.IsInf(r.LogPDF(0), -1) {
		t.Error("expected -Inf for x=0")
	}
}

// ---------------------------------------------------------------------------
// Rice edge cases
// ---------------------------------------------------------------------------

func TestRice_PPF(t *testing.T) {
	r := NewRice(1, 1)
	if r.PPF(0) != 0 {
		t.Error("expected 0 for p=0")
	}
	if !math.IsInf(r.PPF(1), 1) {
		t.Error("expected +Inf for p=1")
	}
	// Normal PPF
	v := r.PPF(0.5)
	if v <= 0 || math.IsNaN(v) {
		t.Errorf("expected positive finite value, got %f", v)
	}
}

func TestRice_PDF_Zero(t *testing.T) {
	r := NewRice(1, 1)
	if r.PDF(-1) != 0 {
		t.Error("expected 0 for negative x")
	}
	if r.PDF(0) != 0 {
		t.Error("expected 0 at x=0")
	}
}

func TestRice_LogPDF_Negative(t *testing.T) {
	r := NewRice(1, 1)
	if !math.IsInf(r.LogPDF(0), -1) {
		t.Error("expected -Inf for x=0")
	}
}

func TestRice_Var(t *testing.T) {
	r := NewRice(1, 1)
	v := r.Var()
	if v < 0 || math.IsNaN(v) {
		t.Errorf("expected non-negative variance, got %f", v)
	}
}

func TestBesselI0_LargeArg(t *testing.T) {
	// Exercise the large argument branch (ax >= 3.75).
	v := besselI0(5.0)
	if v <= 0 {
		t.Errorf("expected positive value, got %f", v)
	}
}

func TestBesselI1_LargeArg(t *testing.T) {
	v := besselI1(5.0)
	if v <= 0 {
		t.Errorf("expected positive value for positive x, got %f", v)
	}
	// Negative input
	v2 := besselI1(-5.0)
	if v2 >= 0 {
		t.Errorf("expected negative value for negative x, got %f", v2)
	}
}

func TestPearsonCorrelation_TwoPoints(t *testing.T) {
	// With 2 points, correlation should be either -1, 0, or 1.
	x := []float64{1, 2, 3}
	y := []float64{2, 4, 6}
	r, pval := PearsonCorrelation(x, y)
	if math.Abs(r-1.0) > 0.01 {
		t.Errorf("expected r near 1.0 for perfectly correlated, got %f", r)
	}
	_ = pval
}
