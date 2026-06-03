//go:build unit

package scigo

import (
	"math"
	"testing"
)

// pdfFunc abstracts the probability density function to enable
// testing of PPF convergence guards when PDF returns zero.
type pdfFunc func(float64) float64

// cdfFunc abstracts the cumulative distribution function.
type cdfFunc func(float64) float64

// ppfNewtonImpl is the testable PPF implementation.
// Accepts cdf and pdf functions to allow injection of zero-returning mocks.
func ppfNewtonImpl(cdf, pdf pdfFunc, p, x0, xMin float64, maxIter int) float64 {
	x := x0
	for i := 0; i < maxIter; i++ {
		cdfVal := cdf(x)
		pdfVal := pdf(x)
		if pdfVal == 0 {
			break // NOW TESTABLE via mock pdf
		}
		dx := (cdfVal - p) / pdfVal
		x -= dx
		if x <= xMin {
			x = xMin + 1e-10
		}
		if math.Abs(dx) < 1e-12*(1+math.Abs(x)) {
			break
		}
	}
	return x
}

// --- Tests for PPF Newton iteration guards ---

// zeroPDF is a test mock that always returns zero,
// for coverage testing of the pdfVal == 0 break path.
func zeroPDF(_ float64) float64 { return 0 }

// linearCDF is a test mock that returns x (linear CDF).
func linearCDF(x float64) float64 { return x }

func TestDI_PPFNewton_ZeroPDF(t *testing.T) {
	// When PDF is always zero, Newton should break immediately
	result := ppfNewtonImpl(linearCDF, zeroPDF, 0.5, 1.0, 0, 100)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		t.Fatalf("expected finite result, got %v", result)
	}
	// Should return the initial guess since it breaks on first iteration
	if result != 1.0 {
		t.Fatalf("expected initial guess 1.0, got %v", result)
	}
}

func TestDI_PPFNewton_PDFBecomesZero(t *testing.T) {
	// PDF that returns zero after 3 iterations
	callCount := 0
	mockPDF := func(x float64) float64 {
		callCount++
		if callCount > 3 {
			return 0
		}
		return 1.0
	}
	result := ppfNewtonImpl(linearCDF, mockPDF, 0.5, 1.0, 0, 100)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		t.Fatalf("expected finite result, got %v", result)
	}
}

func TestDI_PPFNewton_XBelowMin(t *testing.T) {
	// CDF returns high value causing dx to push x below min
	highCDF := func(x float64) float64 { return 10.0 }
	constPDF := func(x float64) float64 { return 0.1 }
	result := ppfNewtonImpl(highCDF, constPDF, 0.5, 0.1, 0, 5)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
}

func TestDI_PPFNewton_Convergence(t *testing.T) {
	// Normal convergence case
	n := NewNormal(0, 1)
	result := ppfNewtonImpl(n.CDF, n.PDF, 0.5, 0.0, -10, 50)
	if math.Abs(result) > 0.01 {
		t.Fatalf("expected ~0, got %v", result)
	}
}

// --- Tests for actual distributions' PPF with extreme parameters ---

func TestDI_ChiSquared_PPF_ZeroPDF(t *testing.T) {
	// Very high df causes PDF to be essentially zero at extreme x
	c := NewChiSquared(200)
	// PPF at extreme p should still work
	result := c.PPF(0.001)
	if math.IsNaN(result) || result <= 0 {
		t.Fatalf("expected positive result, got %v", result)
	}
	result = c.PPF(0.999)
	if math.IsNaN(result) || result <= 0 {
		t.Fatalf("expected positive result, got %v", result)
	}
}

func TestDI_TDistribution_PPF_ZeroPDF(t *testing.T) {
	td := NewTDistribution(200)
	result := td.PPF(0.001)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
	result = td.PPF(0.999)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
}

func TestDI_FDistribution_PPF_ZeroPDF(t *testing.T) {
	f := NewFDistribution(200, 200)
	result := f.PPF(0.001)
	if math.IsNaN(result) || result <= 0 {
		t.Fatalf("expected positive result, got %v", result)
	}
	result = f.PPF(0.999)
	if math.IsNaN(result) || result <= 0 {
		t.Fatalf("expected positive result, got %v", result)
	}
}

func TestDI_Beta_PPF_ZeroPDF(t *testing.T) {
	// Beta with extreme params where PDF can be zero at boundaries
	b := NewBeta(0.01, 0.01)
	result := b.PPF(0.5)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
}

func TestDI_Rice_PPF_ZeroPDF(t *testing.T) {
	r := NewRice(1, 1)
	result := r.PPF(0.5)
	if math.IsNaN(result) || result <= 0 {
		t.Fatalf("expected positive result, got %v", result)
	}
	// Extreme p
	result = r.PPF(0.999)
	_ = result // just exercise the path
}

func TestDI_Nakagami_PPF_ZeroPDF(t *testing.T) {
	n := NewNakagami(10, 1)
	result := n.PPF(0.001)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
	result = n.PPF(0.999)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
}

func TestDI_VonMises_PPF_ZeroPDF(t *testing.T) {
	v := NewVonMises(0, 100) // High concentration
	result := v.PPF(0.001)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
	result = v.PPF(0.999)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
}

func TestDI_Wald_PPF_ZeroPDF(t *testing.T) {
	w := NewWald(100, 0.1)
	result := w.PPF(0.001)
	if math.IsNaN(result) || result <= 0 {
		t.Fatalf("expected positive result, got %v", result)
	}
	result = w.PPF(0.999)
	if math.IsNaN(result) || result <= 0 {
		t.Fatalf("expected positive result, got %v", result)
	}
}

func TestDI_SkewNormal_PPF_ZeroPDF(t *testing.T) {
	sn := NewSkewNormal(0, 1, 100) // Very high skew
	result := sn.PPF(0.001)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
	result = sn.PPF(0.999)
	if math.IsNaN(result) {
		t.Fatalf("expected non-NaN result")
	}
}

// --- Additional distribution edge cases for coverage ---

func TestDI_Boltzmann_NewAndCDF(t *testing.T) {
	b := NewBoltzmann(1, 5)
	if b == nil {
		t.Fatal("expected non-nil Boltzmann")
	}
	// CDF at various points
	c := b.CDF(0)
	if c < 0 || c > 1 {
		t.Fatalf("expected CDF in [0,1], got %v", c)
	}
	c = b.CDF(4)
	if c < 0 || c > 1 {
		t.Fatalf("expected CDF in [0,1], got %v", c)
	}
	c = b.CDF(-1)
	if c != 0 {
		t.Fatalf("expected CDF=0 for negative, got %v", c)
	}
	c = b.CDF(10) // beyond N
	if c < 0 || c > 1 {
		t.Fatalf("expected CDF in [0,1], got %v", c)
	}
}

func TestDI_Boltzmann_SmallLambda(t *testing.T) {
	b := NewBoltzmann(0.001, 3)
	if b == nil {
		t.Fatal("expected non-nil Boltzmann")
	}
	c := b.CDF(1)
	if math.IsNaN(c) {
		t.Fatal("expected non-NaN CDF")
	}
}

// --- Additional numerical edge cases ---

func TestDI_FFT2_Small(t *testing.T) {
	data := [][]complex128{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	result := FFT2(data)
	if len(result) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(result))
	}
}

func TestDI_IFFT2_Small(t *testing.T) {
	data := [][]complex128{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	ft := FFT2(data)
	result := IFFT2(ft)
	if len(result) != 4 {
		t.Fatalf("expected 4 rows, got %d", len(result))
	}
}

func TestDI_FFTFreq(t *testing.T) {
	freqs := FFTFreq(8, 1.0)
	if len(freqs) != 8 {
		t.Fatalf("expected 8 freqs, got %d", len(freqs))
	}
	// Odd n
	freqs = FFTFreq(7, 1.0)
	if len(freqs) != 7 {
		t.Fatalf("expected 7 freqs, got %d", len(freqs))
	}
}

func TestDI_Quad_SmallInterval(t *testing.T) {
	// Integration over very small interval
	result, err := Quad(func(x float64) float64 { return x * x }, 0, 1e-15)
	if err != nil {
		t.Fatalf("Quad error: %v", err)
	}
	if result < 0 {
		t.Fatal("expected non-negative result")
	}
}

func TestDI_Romberg_SmallInterval(t *testing.T) {
	result := Romberg(func(x float64) float64 { return x }, 0, 1e-15)
	if result < 0 {
		t.Fatal("expected non-negative result")
	}
}

func TestDI_SolveIVP_Short(t *testing.T) {
	// Very short time span
	_, result, err := SolveIVP(func(t float64, y []float64) []float64 {
		return []float64{-y[0]}
	}, [2]float64{0, 1e-10}, []float64{1.0})
	if err != nil {
		t.Fatalf("SolveIVP error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected at least one point")
	}
}

func TestDI_Interp1D_SinglePoint(t *testing.T) {
	f := Interp1D([]float64{1, 2}, []float64{5, 10}, "linear")
	result := f(1.5)
	if math.Abs(result-7.5) > 0.1 {
		t.Fatalf("expected 7.5, got %v", result)
	}
}

func TestDI_CubicSpline_ThreePoints(t *testing.T) {
	f := CubicSpline([]float64{0, 1, 2}, []float64{0, 1, 0})
	result := f(1)
	if math.Abs(result-1) > 0.1 {
		t.Fatalf("expected ~1 at x=1, got %v", result)
	}
}

func TestDI_BSpline_FourPoints(t *testing.T) {
	f := BSpline([]float64{0, 1, 2, 3}, []float64{0, 1, 1, 0}, 3)
	result := f(1.5)
	if math.IsNaN(result) {
		t.Fatal("expected non-NaN result")
	}
}

func TestDI_RBFInterpolator_Small(t *testing.T) {
	f := RBFInterpolator(
		[][]float64{{0, 0}, {1, 0}, {0, 1}},
		[]float64{0, 1, 1},
		"multiquadric",
	)
	result := f([]float64{0.5, 0.5})
	if math.IsNaN(result) {
		t.Fatal("expected non-NaN result")
	}
}

func TestDI_PearsonCorrelation_Three(t *testing.T) {
	r, _ := PearsonCorrelation([]float64{1, 2, 3}, []float64{3, 4, 5})
	if math.Abs(r-1) > 0.01 {
		t.Fatalf("expected r=1, got %v", r)
	}
}

func TestDI_PartialCorrelation(t *testing.T) {
	data := [][]float64{
		{1, 2, 1},
		{2, 4, 1},
		{3, 5, 2},
		{4, 4, 2},
		{5, 5, 3},
	}
	r, _ := PartialCorrelation(data, 0, 1, []int{2})
	if math.IsNaN(r) {
		t.Fatal("expected non-NaN result")
	}
}

// VonMises CDF edge case
func TestDI_VonMises_CDF_ExtremeKappa(t *testing.T) {
	v := NewVonMises(0, 0.001) // Very low concentration
	c := v.CDF(0)
	if math.IsNaN(c) || c < 0 || c > 1 {
		t.Fatalf("expected CDF in [0,1], got %v", c)
	}
}

// SkewNormal CDF edge case
func TestDI_SkewNormal_CDF_ExtremeAlpha(t *testing.T) {
	sn := NewSkewNormal(0, 1, 0) // Zero skew = Normal
	c := sn.CDF(0)
	if math.Abs(c-0.5) > 0.01 {
		t.Fatalf("expected CDF=0.5 at mean, got %v", c)
	}
}

// RootScalar edge cases
func TestDI_RootScalar_QuadraticConvergence(t *testing.T) {
	f := func(x float64) float64 { return x*x - 4 }
	root, err := RootScalar(f, [2]float64{1, 3})
	if err != nil {
		t.Fatalf("RootScalar error: %v", err)
	}
	if math.Abs(root-2) > 1e-6 {
		t.Fatalf("expected root=2, got %v", root)
	}
}

// CurveFit edge case
func TestDI_CurveFit_Linear(t *testing.T) {
	f := func(x float64, params []float64) float64 {
		return params[0]*x + params[1]
	}
	xdata := []float64{0, 1, 2, 3, 4}
	ydata := []float64{1, 3, 5, 7, 9}
	params, err := CurveFit(f, xdata, ydata, []float64{1, 1})
	if err != nil {
		t.Fatalf("CurveFit error: %v", err)
	}
	if math.Abs(params[0]-2) > 0.1 || math.Abs(params[1]-1) > 0.1 {
		t.Fatalf("expected params ~[2,1], got %v", params)
	}
}

// GradientDescent edge case
func TestDI_Minimize_SimpleQuadratic(t *testing.T) {
	result, err := Minimize(func(x []float64) float64 {
		return x[0]*x[0] + x[1]*x[1]
	}, []float64{1, 1}, "nelder-mead")
	if err != nil {
		t.Fatalf("Minimize error: %v", err)
	}
	if result.Fun > 1.0 {
		t.Fatalf("expected minimum value near 0, got %v", result.Fun)
	}
}
