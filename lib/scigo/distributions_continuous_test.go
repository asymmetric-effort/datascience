//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// F-Distribution Tests
// ---------------------------------------------------------------------------

func TestFDistributionInterface(t *testing.T) {
	var _ Distribution = &FDistribution{}
}

func TestFDistributionPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for df1=0")
		}
	}()
	NewFDistribution(0, 5)
}

func TestFDistributionMeanVar(t *testing.T) {
	f := NewFDistribution(5, 10)
	// Mean = d2/(d2-2) = 10/8 = 1.25
	if !approxEqual(f.Mean(), 1.25, 1e-12) {
		t.Errorf("F(5,10).Mean() = %v, want 1.25", f.Mean())
	}
	// Var = 2*d2^2*(d1+d2-2) / (d1*(d2-2)^2*(d2-4))
	// = 2*100*13 / (5*64*6) = 2600/1920
	wantVar := 2600.0 / 1920.0
	if !approxEqual(f.Var(), wantVar, 1e-10) {
		t.Errorf("F(5,10).Var() = %v, want %v", f.Var(), wantVar)
	}

	// Mean undefined for df2 <= 2
	f2 := NewFDistribution(5, 2)
	if !math.IsNaN(f2.Mean()) {
		t.Error("F(5,2).Mean() should be NaN")
	}

	// Var undefined for df2 <= 4
	f3 := NewFDistribution(5, 4)
	if !math.IsNaN(f3.Var()) {
		t.Error("F(5,4).Var() should be NaN")
	}
}

func TestFDistributionPDF(t *testing.T) {
	f := NewFDistribution(5, 10)
	// PDF at x=1 should be a known value
	got := f.PDF(1)
	// Manual calculation: use the formula
	d1, d2 := 5.0, 10.0
	logNum := (d1/2)*math.Log(d1) + (d2/2)*math.Log(d2) + (d1/2-1)*math.Log(1)
	logDen := ((d1 + d2) / 2) * math.Log(d1+d2)
	logBeta := Gammaln(d1/2) + Gammaln(d2/2) - Gammaln((d1+d2)/2)
	want := math.Exp(logNum - logDen - logBeta)
	if !approxEqual(got, want, 1e-10) {
		t.Errorf("F(5,10).PDF(1) = %v, want %v", got, want)
	}

	// PDF for x <= 0 should be 0
	if f.PDF(0) != 0 {
		t.Error("F PDF at 0 should be 0")
	}
	if f.PDF(-1) != 0 {
		t.Error("F PDF at -1 should be 0")
	}
}

func TestFDistributionCDF(t *testing.T) {
	f := NewFDistribution(5, 10)
	// CDF(0) = 0
	if f.CDF(0) != 0 {
		t.Error("F CDF(0) should be 0")
	}
	// CDF should be monotonically increasing
	prev := 0.0
	for _, x := range []float64{0.1, 0.5, 1, 2, 5, 10} {
		c := f.CDF(x)
		if c < prev-1e-10 {
			t.Errorf("F CDF not monotonic at x=%v", x)
		}
		prev = c
	}
	// CDF at large x should be close to 1
	if !approxEqual(f.CDF(100), 1.0, 1e-4) {
		t.Errorf("F(5,10).CDF(100) = %v, want ~1", f.CDF(100))
	}
}

func TestFDistributionPPF(t *testing.T) {
	f := NewFDistribution(5, 10)
	// Round-trip
	for _, x := range []float64{0.5, 1, 2, 5} {
		p := f.CDF(x)
		got := f.PPF(p)
		if !approxEqual(got, x, 1e-4) {
			t.Errorf("F(5,10).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestFDistributionLogPDF(t *testing.T) {
	f := NewFDistribution(5, 10)
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := f.LogPDF(x)
		want := math.Log(f.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("F(5,10).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Lognormal Tests
// ---------------------------------------------------------------------------

func TestLognormalInterface(t *testing.T) {
	var _ Distribution = &Lognormal{}
}

func TestLognormalPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for sigma=0")
		}
	}()
	NewLognormal(0, 0)
}

func TestLognormalMeanVar(t *testing.T) {
	ln := NewLognormal(0, 1)
	// Mean = exp(mu + sigma^2/2) = exp(0.5)
	wantMean := math.Exp(0.5)
	if !approxEqual(ln.Mean(), wantMean, 1e-12) {
		t.Errorf("Lognormal(0,1).Mean() = %v, want %v", ln.Mean(), wantMean)
	}
	// Var = (exp(sigma^2)-1)*exp(2*mu+sigma^2) = (e-1)*e
	wantVar := (math.E - 1) * math.E
	if !approxEqual(ln.Var(), wantVar, 1e-10) {
		t.Errorf("Lognormal(0,1).Var() = %v, want %v", ln.Var(), wantVar)
	}
}

func TestLognormalPDF(t *testing.T) {
	ln := NewLognormal(0, 1)
	// At x=1: PDF = exp(-0)/sqrt(2*pi) = 1/sqrt(2*pi)
	got := ln.PDF(1)
	want := 1.0 / math.Sqrt(2*math.Pi)
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("Lognormal(0,1).PDF(1) = %v, want %v", got, want)
	}

	// PDF for x <= 0 should be 0
	if ln.PDF(0) != 0 {
		t.Error("Lognormal PDF at 0 should be 0")
	}
	if ln.PDF(-1) != 0 {
		t.Error("Lognormal PDF at -1 should be 0")
	}
}

func TestLognormalCDF(t *testing.T) {
	ln := NewLognormal(0, 1)
	// CDF(1) = Phi(0) = 0.5
	got := ln.CDF(1)
	if !approxEqual(got, 0.5, 1e-10) {
		t.Errorf("Lognormal(0,1).CDF(1) = %v, want 0.5", got)
	}

	// CDF should be monotonic
	prev := 0.0
	for _, x := range []float64{0.1, 0.5, 1, 2, 5, 10} {
		c := ln.CDF(x)
		if c < prev-1e-10 {
			t.Errorf("Lognormal CDF not monotonic at x=%v", x)
		}
		prev = c
	}
}

func TestLognormalPPF(t *testing.T) {
	ln := NewLognormal(0, 1)
	// PPF(0.5) = exp(mu) = 1
	got := ln.PPF(0.5)
	if !approxEqual(got, 1.0, 1e-10) {
		t.Errorf("Lognormal(0,1).PPF(0.5) = %v, want 1", got)
	}
	// Round-trip
	for _, x := range []float64{0.5, 1, 2, 5} {
		p := ln.CDF(x)
		got := ln.PPF(p)
		if !approxEqual(got, x, 1e-6) {
			t.Errorf("Lognormal(0,1).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestLognormalLogPDF(t *testing.T) {
	ln := NewLognormal(0, 1)
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := ln.LogPDF(x)
		want := math.Log(ln.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Lognormal(0,1).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Weibull Tests
// ---------------------------------------------------------------------------

func TestWeibullInterface(t *testing.T) {
	var _ Distribution = &Weibull{}
}

func TestWeibullPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for shape=0")
		}
	}()
	NewWeibull(0, 1)
}

func TestWeibullMeanVar(t *testing.T) {
	// Weibull(1, lambda) = Exponential(1/lambda)
	w := NewWeibull(1, 2)
	// Mean = scale * Gamma(1 + 1/shape) = 2 * Gamma(2) = 2
	if !approxEqual(w.Mean(), 2.0, 1e-12) {
		t.Errorf("Weibull(1,2).Mean() = %v, want 2", w.Mean())
	}
	// Var = scale^2 * (Gamma(1+2/k) - Gamma(1+1/k)^2) = 4*(Gamma(3)-Gamma(2)^2) = 4*(2-1) = 4
	if !approxEqual(w.Var(), 4.0, 1e-10) {
		t.Errorf("Weibull(1,2).Var() = %v, want 4", w.Var())
	}
}

func TestWeibullPDF(t *testing.T) {
	// Weibull(1, 1) = Exponential(1): PDF(x) = exp(-x)
	w := NewWeibull(1, 1)
	for _, x := range []float64{0.5, 1, 2} {
		got := w.PDF(x)
		want := math.Exp(-x)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Weibull(1,1).PDF(%v) = %v, want %v", x, got, want)
		}
	}
	if w.PDF(-1) != 0 {
		t.Error("Weibull PDF for negative x should be 0")
	}
}

func TestWeibullCDF(t *testing.T) {
	w := NewWeibull(1, 1)
	// CDF(x) = 1 - exp(-x) for shape=1
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := w.CDF(x)
		want := 1 - math.Exp(-x)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Weibull(1,1).CDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestWeibullPPF(t *testing.T) {
	w := NewWeibull(2, 3)
	for _, x := range []float64{0.5, 1, 2, 5} {
		p := w.CDF(x)
		got := w.PPF(p)
		if !approxEqual(got, x, 1e-10) {
			t.Errorf("Weibull(2,3).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestWeibullLogPDF(t *testing.T) {
	w := NewWeibull(2, 3)
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := w.LogPDF(x)
		want := math.Log(w.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Weibull(2,3).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Pareto Tests
// ---------------------------------------------------------------------------

func TestParetoInterface(t *testing.T) {
	var _ Distribution = &Pareto{}
}

func TestParetoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for alpha=0")
		}
	}()
	NewPareto(0, 1)
}

func TestParetoMeanVar(t *testing.T) {
	p := NewPareto(3, 1)
	// Mean = alpha*xm/(alpha-1) = 3/2 = 1.5
	if !approxEqual(p.Mean(), 1.5, 1e-12) {
		t.Errorf("Pareto(3,1).Mean() = %v, want 1.5", p.Mean())
	}
	// Var = xm^2 * alpha / ((alpha-1)^2*(alpha-2)) = 3/(4*1) = 0.75
	if !approxEqual(p.Var(), 0.75, 1e-12) {
		t.Errorf("Pareto(3,1).Var() = %v, want 0.75", p.Var())
	}

	// Mean infinite for alpha <= 1
	p2 := NewPareto(1, 1)
	if !math.IsInf(p2.Mean(), 1) {
		t.Error("Pareto(1,1).Mean() should be +Inf")
	}
}

func TestParetoPDF(t *testing.T) {
	p := NewPareto(2, 1)
	// PDF(x) = 2/x^3 for x >= 1
	for _, x := range []float64{1, 2, 3, 5} {
		got := p.PDF(x)
		want := 2.0 / (x * x * x)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Pareto(2,1).PDF(%v) = %v, want %v", x, got, want)
		}
	}
	if p.PDF(0.5) != 0 {
		t.Error("Pareto PDF below xm should be 0")
	}
}

func TestParetoCDF(t *testing.T) {
	p := NewPareto(2, 1)
	// CDF(x) = 1 - (1/x)^2
	for _, x := range []float64{1, 2, 3, 5} {
		got := p.CDF(x)
		want := 1 - 1.0/(x*x)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Pareto(2,1).CDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestParetoPPF(t *testing.T) {
	p := NewPareto(3, 2)
	for _, x := range []float64{2, 3, 5, 10} {
		pr := p.CDF(x)
		got := p.PPF(pr)
		if !approxEqual(got, x, 1e-10) {
			t.Errorf("Pareto(3,2).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestParetoLogPDF(t *testing.T) {
	p := NewPareto(3, 2)
	for _, x := range []float64{2, 3, 5, 10} {
		got := p.LogPDF(x)
		want := math.Log(p.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Pareto(3,2).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Cauchy Tests
// ---------------------------------------------------------------------------

func TestCauchyInterface(t *testing.T) {
	var _ Distribution = &Cauchy{}
}

func TestCauchyPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for scale=0")
		}
	}()
	NewCauchy(0, 0)
}

func TestCauchyMeanVar(t *testing.T) {
	c := NewCauchy(0, 1)
	if !math.IsNaN(c.Mean()) {
		t.Error("Cauchy Mean should be NaN")
	}
	if !math.IsNaN(c.Var()) {
		t.Error("Cauchy Var should be NaN")
	}
}

func TestCauchyPDF(t *testing.T) {
	c := NewCauchy(0, 1)
	// PDF(0) = 1/pi
	got := c.PDF(0)
	want := 1.0 / math.Pi
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("Cauchy(0,1).PDF(0) = %v, want %v", got, want)
	}
	// Symmetric
	if !approxEqual(c.PDF(1), c.PDF(-1), 1e-12) {
		t.Error("Cauchy PDF should be symmetric")
	}
	// PDF(1) = 1/(pi*2)
	if !approxEqual(c.PDF(1), 1.0/(2*math.Pi), 1e-12) {
		t.Errorf("Cauchy(0,1).PDF(1) = %v, want %v", c.PDF(1), 1.0/(2*math.Pi))
	}
}

func TestCauchyCDF(t *testing.T) {
	c := NewCauchy(0, 1)
	// CDF(0) = 0.5
	if !approxEqual(c.CDF(0), 0.5, 1e-12) {
		t.Errorf("Cauchy(0,1).CDF(0) = %v, want 0.5", c.CDF(0))
	}
	// CDF(1) = 0.5 + atan(1)/pi = 0.5 + 0.25 = 0.75
	if !approxEqual(c.CDF(1), 0.75, 1e-12) {
		t.Errorf("Cauchy(0,1).CDF(1) = %v, want 0.75", c.CDF(1))
	}
}

func TestCauchyPPF(t *testing.T) {
	c := NewCauchy(0, 1)
	// PPF(0.5) = 0
	if !approxEqual(c.PPF(0.5), 0, 1e-12) {
		t.Errorf("Cauchy(0,1).PPF(0.5) = %v, want 0", c.PPF(0.5))
	}
	// Round-trip
	for _, x := range []float64{-5, -1, 0, 1, 5} {
		p := c.CDF(x)
		got := c.PPF(p)
		if !approxEqual(got, x, 1e-10) {
			t.Errorf("Cauchy(0,1).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestCauchyLogPDF(t *testing.T) {
	c := NewCauchy(0, 1)
	for _, x := range []float64{-2, 0, 1, 5} {
		got := c.LogPDF(x)
		want := math.Log(c.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Cauchy(0,1).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Laplace Tests
// ---------------------------------------------------------------------------

func TestLaplaceInterface(t *testing.T) {
	var _ Distribution = &Laplace{}
}

func TestLaplacePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for scale=0")
		}
	}()
	NewLaplace(0, 0)
}

func TestLaplaceMeanVar(t *testing.T) {
	l := NewLaplace(3, 2)
	if l.Mean() != 3 {
		t.Errorf("Laplace(3,2).Mean() = %v, want 3", l.Mean())
	}
	// Var = 2*b^2 = 8
	if !approxEqual(l.Var(), 8, 1e-12) {
		t.Errorf("Laplace(3,2).Var() = %v, want 8", l.Var())
	}
}

func TestLaplacePDF(t *testing.T) {
	l := NewLaplace(0, 1)
	// PDF(0) = 1/(2*1) = 0.5
	if !approxEqual(l.PDF(0), 0.5, 1e-12) {
		t.Errorf("Laplace(0,1).PDF(0) = %v, want 0.5", l.PDF(0))
	}
	// PDF(1) = exp(-1)/2
	if !approxEqual(l.PDF(1), math.Exp(-1)/2, 1e-12) {
		t.Errorf("Laplace(0,1).PDF(1) = %v, want %v", l.PDF(1), math.Exp(-1)/2)
	}
	// Symmetric
	if !approxEqual(l.PDF(2), l.PDF(-2), 1e-12) {
		t.Error("Laplace PDF should be symmetric")
	}
}

func TestLaplaceCDF(t *testing.T) {
	l := NewLaplace(0, 1)
	// CDF(0) = 0.5
	if !approxEqual(l.CDF(0), 0.5, 1e-12) {
		t.Errorf("Laplace(0,1).CDF(0) = %v, want 0.5", l.CDF(0))
	}
	// CDF(-inf) -> 0, CDF(inf) -> 1
	if !approxEqual(l.CDF(-100), 0, 1e-10) {
		t.Error("Laplace CDF at large negative should be ~0")
	}
}

func TestLaplacePPF(t *testing.T) {
	l := NewLaplace(0, 1)
	// PPF(0.5) = 0
	if !approxEqual(l.PPF(0.5), 0, 1e-12) {
		t.Errorf("Laplace(0,1).PPF(0.5) = %v, want 0", l.PPF(0.5))
	}
	// Round-trip
	for _, x := range []float64{-5, -1, 0, 1, 5} {
		p := l.CDF(x)
		got := l.PPF(p)
		if !approxEqual(got, x, 1e-10) {
			t.Errorf("Laplace(0,1).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestLaplaceLogPDF(t *testing.T) {
	l := NewLaplace(0, 1)
	for _, x := range []float64{-2, 0, 1, 5} {
		got := l.LogPDF(x)
		want := math.Log(l.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Laplace(0,1).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Logistic Tests
// ---------------------------------------------------------------------------

func TestLogisticInterface(t *testing.T) {
	var _ Distribution = &Logistic{}
}

func TestLogisticPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for scale=0")
		}
	}()
	NewLogistic(0, 0)
}

func TestLogisticMeanVar(t *testing.T) {
	l := NewLogistic(2, 3)
	if l.Mean() != 2 {
		t.Errorf("Logistic(2,3).Mean() = %v, want 2", l.Mean())
	}
	// Var = s^2 * pi^2 / 3 = 9*pi^2/3 = 3*pi^2
	wantVar := 3 * math.Pi * math.Pi
	if !approxEqual(l.Var(), wantVar, 1e-10) {
		t.Errorf("Logistic(2,3).Var() = %v, want %v", l.Var(), wantVar)
	}
}

func TestLogisticPDF(t *testing.T) {
	l := NewLogistic(0, 1)
	// PDF(0) = exp(0)/(1+exp(0))^2 = 1/4 = 0.25
	if !approxEqual(l.PDF(0), 0.25, 1e-12) {
		t.Errorf("Logistic(0,1).PDF(0) = %v, want 0.25", l.PDF(0))
	}
}

func TestLogisticCDF(t *testing.T) {
	l := NewLogistic(0, 1)
	// CDF(0) = 1/(1+1) = 0.5
	if !approxEqual(l.CDF(0), 0.5, 1e-12) {
		t.Errorf("Logistic(0,1).CDF(0) = %v, want 0.5", l.CDF(0))
	}
}

func TestLogisticPPF(t *testing.T) {
	l := NewLogistic(0, 1)
	// PPF(0.5) = 0
	if !approxEqual(l.PPF(0.5), 0, 1e-12) {
		t.Errorf("Logistic(0,1).PPF(0.5) = %v, want 0", l.PPF(0.5))
	}
	// Round-trip
	for _, x := range []float64{-5, -1, 0, 1, 5} {
		p := l.CDF(x)
		got := l.PPF(p)
		if !approxEqual(got, x, 1e-10) {
			t.Errorf("Logistic(0,1).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestLogisticLogPDF(t *testing.T) {
	l := NewLogistic(0, 1)
	for _, x := range []float64{-2, 0, 1, 5} {
		got := l.LogPDF(x)
		want := math.Log(l.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Logistic(0,1).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Gumbel Tests
// ---------------------------------------------------------------------------

func TestGumbelInterface(t *testing.T) {
	var _ Distribution = &Gumbel{}
}

func TestGumbelPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for scale=0")
		}
	}()
	NewGumbel(0, 0)
}

func TestGumbelMeanVar(t *testing.T) {
	g := NewGumbel(0, 1)
	euler := 0.5772156649015329
	if !approxEqual(g.Mean(), euler, 1e-12) {
		t.Errorf("Gumbel(0,1).Mean() = %v, want %v", g.Mean(), euler)
	}
	// Var = pi^2/6
	wantVar := math.Pi * math.Pi / 6
	if !approxEqual(g.Var(), wantVar, 1e-10) {
		t.Errorf("Gumbel(0,1).Var() = %v, want %v", g.Var(), wantVar)
	}
}

func TestGumbelPDF(t *testing.T) {
	g := NewGumbel(0, 1)
	// PDF(0) = exp(-exp(0)) * exp(0) = exp(-1)
	got := g.PDF(0)
	want := math.Exp(-1)
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("Gumbel(0,1).PDF(0) = %v, want %v", got, want)
	}
}

func TestGumbelCDF(t *testing.T) {
	g := NewGumbel(0, 1)
	// CDF(0) = exp(-exp(0)) = exp(-1)
	got := g.CDF(0)
	want := math.Exp(-1)
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("Gumbel(0,1).CDF(0) = %v, want %v", got, want)
	}
}

func TestGumbelPPF(t *testing.T) {
	g := NewGumbel(0, 1)
	for _, x := range []float64{-2, 0, 1, 5} {
		p := g.CDF(x)
		got := g.PPF(p)
		if !approxEqual(got, x, 1e-10) {
			t.Errorf("Gumbel(0,1).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestGumbelLogPDF(t *testing.T) {
	g := NewGumbel(0, 1)
	for _, x := range []float64{-2, 0, 1, 5} {
		got := g.LogPDF(x)
		want := math.Log(g.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Gumbel(0,1).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Rayleigh Tests
// ---------------------------------------------------------------------------

func TestRayleighInterface(t *testing.T) {
	var _ Distribution = &Rayleigh{}
}

func TestRayleighPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for sigma=0")
		}
	}()
	NewRayleigh(0)
}

func TestRayleighMeanVar(t *testing.T) {
	r := NewRayleigh(1)
	wantMean := math.Sqrt(math.Pi / 2)
	if !approxEqual(r.Mean(), wantMean, 1e-12) {
		t.Errorf("Rayleigh(1).Mean() = %v, want %v", r.Mean(), wantMean)
	}
	wantVar := (2 - math.Pi/2)
	if !approxEqual(r.Var(), wantVar, 1e-12) {
		t.Errorf("Rayleigh(1).Var() = %v, want %v", r.Var(), wantVar)
	}
}

func TestRayleighPDF(t *testing.T) {
	r := NewRayleigh(1)
	// PDF(x) = x * exp(-x^2/2) for sigma=1
	for _, x := range []float64{0.5, 1, 2} {
		got := r.PDF(x)
		want := x * math.Exp(-x*x/2)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Rayleigh(1).PDF(%v) = %v, want %v", x, got, want)
		}
	}
	if r.PDF(-1) != 0 {
		t.Error("Rayleigh PDF for negative x should be 0")
	}
}

func TestRayleighCDF(t *testing.T) {
	r := NewRayleigh(1)
	// CDF(x) = 1 - exp(-x^2/2)
	for _, x := range []float64{0.5, 1, 2, 3} {
		got := r.CDF(x)
		want := 1 - math.Exp(-x*x/2)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Rayleigh(1).CDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestRayleighPPF(t *testing.T) {
	r := NewRayleigh(2)
	for _, x := range []float64{0.5, 1, 2, 5} {
		p := r.CDF(x)
		got := r.PPF(p)
		if !approxEqual(got, x, 1e-10) {
			t.Errorf("Rayleigh(2).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestRayleighLogPDF(t *testing.T) {
	r := NewRayleigh(2)
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := r.LogPDF(x)
		want := math.Log(r.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Rayleigh(2).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Rice Tests
// ---------------------------------------------------------------------------

func TestRiceInterface(t *testing.T) {
	var _ Distribution = &Rice{}
}

func TestRicePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for sigma=0")
		}
	}()
	NewRice(1, 0)
}

func TestRicePDF(t *testing.T) {
	// Rice(0, sigma) reduces to Rayleigh(sigma)
	r := NewRice(0, 1)
	ray := NewRayleigh(1)
	for _, x := range []float64{0.5, 1, 2} {
		got := r.PDF(x)
		want := ray.PDF(x)
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Rice(0,1).PDF(%v) = %v, want Rayleigh %v", x, got, want)
		}
	}
}

func TestRiceCDF(t *testing.T) {
	r := NewRice(0, 1)
	// Should match Rayleigh CDF
	ray := NewRayleigh(1)
	for _, x := range []float64{0.5, 1, 2} {
		got := r.CDF(x)
		want := ray.CDF(x)
		if !approxEqual(got, want, 1e-3) {
			t.Errorf("Rice(0,1).CDF(%v) = %v, want Rayleigh %v", x, got, want)
		}
	}
}

func TestRiceLogPDF(t *testing.T) {
	r := NewRice(1, 2)
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := r.LogPDF(x)
		want := math.Log(r.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Rice(1,2).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestRiceMeanVar(t *testing.T) {
	// For nu=0, Rice reduces to Rayleigh
	r := NewRice(0, 1)
	ray := NewRayleigh(1)
	if !approxEqual(r.Mean(), ray.Mean(), 1e-4) {
		t.Errorf("Rice(0,1).Mean() = %v, want Rayleigh mean %v", r.Mean(), ray.Mean())
	}
}

// ---------------------------------------------------------------------------
// Nakagami Tests
// ---------------------------------------------------------------------------

func TestNakagamiInterface(t *testing.T) {
	var _ Distribution = &Nakagami{}
}

func TestNakagamiPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for m=0")
		}
	}()
	NewNakagami(0, 1)
}

func TestNakagamiPDF(t *testing.T) {
	// Nakagami(1, omega) = Rayleigh(sqrt(omega/2)) in terms of PDF shape
	// For m=1, omega=2: Nakagami PDF(x) = 2*x*exp(-x^2/2) which is Rayleigh(1)
	n := NewNakagami(1, 2)
	ray := NewRayleigh(1)
	for _, x := range []float64{0.5, 1, 2} {
		got := n.PDF(x)
		want := ray.PDF(x)
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Nakagami(1,2).PDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestNakagamiCDF(t *testing.T) {
	n := NewNakagami(1, 2)
	ray := NewRayleigh(1)
	for _, x := range []float64{0.5, 1, 2} {
		got := n.CDF(x)
		want := ray.CDF(x)
		if !approxEqual(got, want, 1e-8) {
			t.Errorf("Nakagami(1,2).CDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestNakagamiMeanVar(t *testing.T) {
	n := NewNakagami(1, 2)
	ray := NewRayleigh(1)
	if !approxEqual(n.Mean(), ray.Mean(), 1e-10) {
		t.Errorf("Nakagami(1,2).Mean() = %v, want %v", n.Mean(), ray.Mean())
	}
	if !approxEqual(n.Var(), ray.Var(), 1e-10) {
		t.Errorf("Nakagami(1,2).Var() = %v, want %v", n.Var(), ray.Var())
	}
}

func TestNakagamiLogPDF(t *testing.T) {
	n := NewNakagami(2, 3)
	for _, x := range []float64{0.5, 1, 2} {
		got := n.LogPDF(x)
		want := math.Log(n.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Nakagami(2,3).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// VonMises Tests
// ---------------------------------------------------------------------------

func TestVonMisesInterface(t *testing.T) {
	var _ Distribution = &VonMises{}
}

func TestVonMisesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for kappa<0")
		}
	}()
	NewVonMises(0, -1)
}

func TestVonMisesMeanVar(t *testing.T) {
	v := NewVonMises(1.5, 2)
	if v.Mean() != 1.5 {
		t.Errorf("VonMises(1.5,2).Mean() = %v, want 1.5", v.Mean())
	}
	// Var = 1 - I1(kappa)/I0(kappa)
	wantVar := 1 - besselI1(2)/besselI0(2)
	if !approxEqual(v.Var(), wantVar, 1e-10) {
		t.Errorf("VonMises(1.5,2).Var() = %v, want %v", v.Var(), wantVar)
	}
}

func TestVonMisesPDF(t *testing.T) {
	// For kappa=0, VonMises is uniform on [-pi, pi]
	v := NewVonMises(0, 0)
	// PDF should be 1/(2*pi) everywhere
	for _, x := range []float64{-math.Pi, 0, math.Pi} {
		got := v.PDF(x)
		want := 1.0 / (2 * math.Pi)
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("VonMises(0,0).PDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestVonMisesCDF(t *testing.T) {
	v := NewVonMises(0, 0)
	// For uniform: CDF at pi should be 1
	got := v.CDF(math.Pi)
	if !approxEqual(got, 1.0, 1e-3) {
		t.Errorf("VonMises(0,0).CDF(pi) = %v, want 1", got)
	}
	// CDF at 0 should be 0.5
	got = v.CDF(0)
	if !approxEqual(got, 0.5, 1e-3) {
		t.Errorf("VonMises(0,0).CDF(0) = %v, want 0.5", got)
	}
}

func TestVonMisesLogPDF(t *testing.T) {
	v := NewVonMises(0, 2)
	for _, x := range []float64{-1, 0, 1} {
		got := v.LogPDF(x)
		want := math.Log(v.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("VonMises(0,2).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Wald (Inverse Gaussian) Tests
// ---------------------------------------------------------------------------

func TestWaldInterface(t *testing.T) {
	var _ Distribution = &Wald{}
}

func TestWaldPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for mu=0")
		}
	}()
	NewWald(0, 1)
}

func TestWaldMeanVar(t *testing.T) {
	w := NewWald(2, 3)
	if w.Mean() != 2 {
		t.Errorf("Wald(2,3).Mean() = %v, want 2", w.Mean())
	}
	// Var = mu^3/lambda = 8/3
	wantVar := 8.0 / 3.0
	if !approxEqual(w.Var(), wantVar, 1e-12) {
		t.Errorf("Wald(2,3).Var() = %v, want %v", w.Var(), wantVar)
	}
}

func TestWaldPDF(t *testing.T) {
	w := NewWald(1, 1)
	// PDF(x) = sqrt(1/(2*pi*x^3)) * exp(-(x-1)^2/(2*x))
	for _, x := range []float64{0.5, 1, 2, 3} {
		got := w.PDF(x)
		want := math.Sqrt(1/(2*math.Pi*x*x*x)) * math.Exp(-(x-1)*(x-1)/(2*x))
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Wald(1,1).PDF(%v) = %v, want %v", x, got, want)
		}
	}
	if w.PDF(0) != 0 {
		t.Error("Wald PDF at 0 should be 0")
	}
	if w.PDF(-1) != 0 {
		t.Error("Wald PDF for negative x should be 0")
	}
}

func TestWaldCDF(t *testing.T) {
	w := NewWald(1, 1)
	// CDF should be monotonically increasing
	prev := 0.0
	for _, x := range []float64{0.1, 0.5, 1, 2, 5, 10} {
		c := w.CDF(x)
		if c < prev-1e-10 {
			t.Errorf("Wald CDF not monotonic at x=%v", x)
		}
		prev = c
	}
	// CDF should approach 1 for large x
	if !approxEqual(w.CDF(100), 1.0, 1e-4) {
		t.Errorf("Wald(1,1).CDF(100) = %v, want ~1", w.CDF(100))
	}
}

func TestWaldPPF(t *testing.T) {
	w := NewWald(2, 3)
	for _, x := range []float64{0.5, 1, 2, 5} {
		p := w.CDF(x)
		got := w.PPF(p)
		if !approxEqual(got, x, 1e-4) {
			t.Errorf("Wald(2,3).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestWaldLogPDF(t *testing.T) {
	w := NewWald(2, 3)
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := w.LogPDF(x)
		want := math.Log(w.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Wald(2,3).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// HalfNormal Tests
// ---------------------------------------------------------------------------

func TestHalfNormalInterface(t *testing.T) {
	var _ Distribution = &HalfNormal{}
}

func TestHalfNormalPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for sigma=0")
		}
	}()
	NewHalfNormal(0)
}

func TestHalfNormalMeanVar(t *testing.T) {
	h := NewHalfNormal(1)
	wantMean := math.Sqrt(2 / math.Pi)
	if !approxEqual(h.Mean(), wantMean, 1e-12) {
		t.Errorf("HalfNormal(1).Mean() = %v, want %v", h.Mean(), wantMean)
	}
	wantVar := 1 - 2/math.Pi
	if !approxEqual(h.Var(), wantVar, 1e-12) {
		t.Errorf("HalfNormal(1).Var() = %v, want %v", h.Var(), wantVar)
	}
}

func TestHalfNormalPDF(t *testing.T) {
	h := NewHalfNormal(1)
	// PDF(0) = sqrt(2/pi)
	got := h.PDF(0)
	want := math.Sqrt(2 / math.Pi)
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("HalfNormal(1).PDF(0) = %v, want %v", got, want)
	}
	// PDF = 2 * standard normal PDF for x >= 0
	n := NewNormal(0, 1)
	for _, x := range []float64{0, 0.5, 1, 2} {
		got := h.PDF(x)
		want := 2 * n.PDF(x)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("HalfNormal(1).PDF(%v) = %v, want %v", x, got, want)
		}
	}
	if h.PDF(-1) != 0 {
		t.Error("HalfNormal PDF for negative x should be 0")
	}
}

func TestHalfNormalCDF(t *testing.T) {
	h := NewHalfNormal(1)
	// CDF(0) = 0
	if h.CDF(0) != 0 {
		t.Error("HalfNormal CDF(0) should be 0")
	}
	// CDF = erf(x/sqrt(2))
	for _, x := range []float64{0.5, 1, 2} {
		got := h.CDF(x)
		want := math.Erf(x / math.Sqrt2)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("HalfNormal(1).CDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestHalfNormalPPF(t *testing.T) {
	h := NewHalfNormal(2)
	for _, x := range []float64{0.5, 1, 2, 5} {
		p := h.CDF(x)
		got := h.PPF(p)
		if !approxEqual(got, x, 1e-8) {
			t.Errorf("HalfNormal(2).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestHalfNormalLogPDF(t *testing.T) {
	h := NewHalfNormal(2)
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := h.LogPDF(x)
		want := math.Log(h.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("HalfNormal(2).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// TruncatedNormal Tests
// ---------------------------------------------------------------------------

func TestTruncatedNormalInterface(t *testing.T) {
	var _ Distribution = &TruncatedNormal{}
}

func TestTruncatedNormalPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for a >= b")
		}
	}()
	NewTruncatedNormal(0, 1, 5, 3)
}

func TestTruncatedNormalPDF(t *testing.T) {
	tn := NewTruncatedNormal(0, 1, -1, 1)
	// PDF should be 0 outside [a, b]
	if tn.PDF(-2) != 0 {
		t.Error("TruncatedNormal PDF outside [a,b] should be 0")
	}
	if tn.PDF(2) != 0 {
		t.Error("TruncatedNormal PDF outside [a,b] should be 0")
	}
	// PDF at 0 should be Normal(0,1).PDF(0) / Z
	n := NewNormal(0, 1)
	z := n.CDF(1) - n.CDF(-1)
	got := tn.PDF(0)
	want := n.PDF(0) / z
	if !approxEqual(got, want, 1e-10) {
		t.Errorf("TruncatedNormal(0,1,-1,1).PDF(0) = %v, want %v", got, want)
	}
}

func TestTruncatedNormalCDF(t *testing.T) {
	tn := NewTruncatedNormal(0, 1, -1, 1)
	// CDF(a) = 0, CDF(b) = 1
	if tn.CDF(-1) != 0 {
		t.Error("TruncatedNormal CDF(a) should be 0")
	}
	if tn.CDF(1) != 1 {
		t.Error("TruncatedNormal CDF(b) should be 1")
	}
	// CDF(0) = (Phi(0)-Phi(-1))/(Phi(1)-Phi(-1)) = (0.5-Phi(-1))/(Phi(1)-Phi(-1))
	n := NewNormal(0, 1)
	got := tn.CDF(0)
	want := (n.CDF(0) - n.CDF(-1)) / (n.CDF(1) - n.CDF(-1))
	if !approxEqual(got, want, 1e-10) {
		t.Errorf("TruncatedNormal(0,1,-1,1).CDF(0) = %v, want %v", got, want)
	}
}

func TestTruncatedNormalPPF(t *testing.T) {
	tn := NewTruncatedNormal(0, 1, -2, 2)
	for _, x := range []float64{-1.5, -0.5, 0, 0.5, 1.5} {
		p := tn.CDF(x)
		got := tn.PPF(p)
		if !approxEqual(got, x, 1e-6) {
			t.Errorf("TruncatedNormal PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestTruncatedNormalMeanVar(t *testing.T) {
	// Symmetric truncation around mean: mean should equal mu
	tn := NewTruncatedNormal(0, 1, -2, 2)
	if !approxEqual(tn.Mean(), 0, 1e-10) {
		t.Errorf("TruncatedNormal(0,1,-2,2).Mean() = %v, want 0", tn.Mean())
	}
	// Variance should be less than 1 (truncation reduces variance)
	if tn.Var() >= 1 || tn.Var() <= 0 {
		t.Errorf("TruncatedNormal(0,1,-2,2).Var() = %v, expected in (0, 1)", tn.Var())
	}
}

func TestTruncatedNormalLogPDF(t *testing.T) {
	tn := NewTruncatedNormal(0, 1, -2, 2)
	for _, x := range []float64{-1, 0, 1} {
		got := tn.LogPDF(x)
		want := math.Log(tn.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("TruncatedNormal LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// SkewNormal Tests
// ---------------------------------------------------------------------------

func TestSkewNormalInterface(t *testing.T) {
	var _ Distribution = &SkewNormal{}
}

func TestSkewNormalPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for scale=0")
		}
	}()
	NewSkewNormal(0, 0, 1)
}

func TestSkewNormalAlphaZero(t *testing.T) {
	// With alpha=0, SkewNormal reduces to Normal
	sn := NewSkewNormal(0, 1, 0)
	n := NewNormal(0, 1)
	for _, x := range []float64{-2, -1, 0, 1, 2} {
		got := sn.PDF(x)
		want := n.PDF(x)
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("SkewNormal(0,1,0).PDF(%v) = %v, want Normal %v", x, got, want)
		}
	}
}

func TestSkewNormalMeanVar(t *testing.T) {
	sn := NewSkewNormal(0, 1, 0)
	// With alpha=0, mean=0, var=1
	if !approxEqual(sn.Mean(), 0, 1e-12) {
		t.Errorf("SkewNormal(0,1,0).Mean() = %v, want 0", sn.Mean())
	}
	if !approxEqual(sn.Var(), 1, 1e-12) {
		t.Errorf("SkewNormal(0,1,0).Var() = %v, want 1", sn.Var())
	}
}

func TestSkewNormalLogPDF(t *testing.T) {
	sn := NewSkewNormal(0, 1, 2)
	for _, x := range []float64{-1, 0, 1, 2} {
		got := sn.LogPDF(x)
		want := math.Log(sn.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("SkewNormal(0,1,2).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestSkewNormalCDF(t *testing.T) {
	sn := NewSkewNormal(0, 1, 0)
	n := NewNormal(0, 1)
	// With alpha=0, CDF should match normal CDF
	for _, x := range []float64{-2, 0, 2} {
		got := sn.CDF(x)
		want := n.CDF(x)
		if !approxEqual(got, want, 1e-3) {
			t.Errorf("SkewNormal(0,1,0).CDF(%v) = %v, want Normal %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// GEV Tests
// ---------------------------------------------------------------------------

func TestGEVInterface(t *testing.T) {
	var _ Distribution = &GeneralizedExtremeValue{}
}

func TestGEVPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for sigma=0")
		}
	}()
	NewGeneralizedExtremeValue(0, 0, 0)
}

func TestGEVGumbelCase(t *testing.T) {
	// GEV with xi=0 is Gumbel
	gev := NewGeneralizedExtremeValue(0, 1, 0)
	g := NewGumbel(0, 1)
	for _, x := range []float64{-2, 0, 1, 5} {
		if !approxEqual(gev.PDF(x), g.PDF(x), 1e-10) {
			t.Errorf("GEV(0,1,0).PDF(%v) = %v, want Gumbel %v", x, gev.PDF(x), g.PDF(x))
		}
		if !approxEqual(gev.CDF(x), g.CDF(x), 1e-10) {
			t.Errorf("GEV(0,1,0).CDF(%v) = %v, want Gumbel %v", x, gev.CDF(x), g.CDF(x))
		}
	}
}

func TestGEVMeanVar(t *testing.T) {
	// Gumbel case
	gev := NewGeneralizedExtremeValue(0, 1, 0)
	euler := 0.5772156649015329
	if !approxEqual(gev.Mean(), euler, 1e-10) {
		t.Errorf("GEV(0,1,0).Mean() = %v, want %v", gev.Mean(), euler)
	}
	wantVar := math.Pi * math.Pi / 6
	if !approxEqual(gev.Var(), wantVar, 1e-10) {
		t.Errorf("GEV(0,1,0).Var() = %v, want %v", gev.Var(), wantVar)
	}
}

func TestGEVPPF(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0.2)
	for _, x := range []float64{0, 1, 3, 5} {
		p := gev.CDF(x)
		got := gev.PPF(p)
		if !approxEqual(got, x, 1e-8) {
			t.Errorf("GEV PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestGEVLogPDF(t *testing.T) {
	gev := NewGeneralizedExtremeValue(0, 1, 0.2)
	for _, x := range []float64{0, 1, 3} {
		got := gev.LogPDF(x)
		want := math.Log(gev.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("GEV LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Levy Tests
// ---------------------------------------------------------------------------

func TestLevyInterface(t *testing.T) {
	var _ Distribution = &Levy{}
}

func TestLevyPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for scale=0")
		}
	}()
	NewLevy(0, 0)
}

func TestLevyMeanVar(t *testing.T) {
	l := NewLevy(0, 1)
	if !math.IsInf(l.Mean(), 1) {
		t.Error("Levy Mean should be +Inf")
	}
	if !math.IsInf(l.Var(), 1) {
		t.Error("Levy Var should be +Inf")
	}
}

func TestLevyPDF(t *testing.T) {
	l := NewLevy(0, 1)
	// PDF(x) = sqrt(1/(2*pi)) * exp(-1/(2*x)) / x^{3/2} for x > 0
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := l.PDF(x)
		want := math.Sqrt(1/(2*math.Pi)) * math.Exp(-1/(2*x)) / (x * math.Sqrt(x))
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Levy(0,1).PDF(%v) = %v, want %v", x, got, want)
		}
	}
	if l.PDF(0) != 0 {
		t.Error("Levy PDF at 0 should be 0")
	}
	if l.PDF(-1) != 0 {
		t.Error("Levy PDF for negative x should be 0")
	}
}

func TestLevyCDF(t *testing.T) {
	l := NewLevy(0, 1)
	// CDF(x) = erfc(sqrt(1/(2x)))
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := l.CDF(x)
		want := math.Erfc(math.Sqrt(1.0 / (2 * x)))
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Levy(0,1).CDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestLevyPPF(t *testing.T) {
	l := NewLevy(0, 1)
	for _, x := range []float64{1, 2, 5, 10} {
		p := l.CDF(x)
		got := l.PPF(p)
		if !approxEqual(got, x, 1e-6) {
			t.Errorf("Levy(0,1).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestLevyLogPDF(t *testing.T) {
	l := NewLevy(0, 1)
	for _, x := range []float64{0.5, 1, 2, 5} {
		got := l.LogPDF(x)
		want := math.Log(l.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("Levy(0,1).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// PowerLaw Tests
// ---------------------------------------------------------------------------

func TestPowerLawInterface(t *testing.T) {
	var _ Distribution = &PowerLaw{}
}

func TestPowerLawPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for alpha=1")
		}
	}()
	NewPowerLaw(1, 1)
}

func TestPowerLawMeanVar(t *testing.T) {
	pl := NewPowerLaw(3, 1)
	// Mean = (alpha-1)*xmin/(alpha-2) = 2/1 = 2
	if !approxEqual(pl.Mean(), 2, 1e-12) {
		t.Errorf("PowerLaw(3,1).Mean() = %v, want 2", pl.Mean())
	}
	// For alpha<=2, Mean is infinite
	pl2 := NewPowerLaw(2, 1)
	if !math.IsInf(pl2.Mean(), 1) {
		t.Error("PowerLaw(2,1).Mean() should be +Inf")
	}

	// Var for alpha=4, xmin=1: 1*(4-1)/((4-2)^2*(4-3)) = 3/4
	pl3 := NewPowerLaw(4, 1)
	if !approxEqual(pl3.Var(), 0.75, 1e-12) {
		t.Errorf("PowerLaw(4,1).Var() = %v, want 0.75", pl3.Var())
	}
}

func TestPowerLawPDF(t *testing.T) {
	pl := NewPowerLaw(2.5, 1)
	// PDF(x) = 1.5 * x^{-2.5} for x >= 1
	for _, x := range []float64{1, 2, 5, 10} {
		got := pl.PDF(x)
		want := 1.5 * math.Pow(x, -2.5)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("PowerLaw(2.5,1).PDF(%v) = %v, want %v", x, got, want)
		}
	}
	if pl.PDF(0.5) != 0 {
		t.Error("PowerLaw PDF below xmin should be 0")
	}
}

func TestPowerLawCDF(t *testing.T) {
	pl := NewPowerLaw(2.5, 1)
	// CDF(x) = 1 - (1/x)^{1.5}
	for _, x := range []float64{1, 2, 5, 10} {
		got := pl.CDF(x)
		want := 1 - math.Pow(1.0/x, 1.5)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("PowerLaw(2.5,1).CDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestPowerLawPPF(t *testing.T) {
	pl := NewPowerLaw(3, 2)
	for _, x := range []float64{2, 3, 5, 10} {
		p := pl.CDF(x)
		got := pl.PPF(p)
		if !approxEqual(got, x, 1e-10) {
			t.Errorf("PowerLaw(3,2).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}
}

func TestPowerLawLogPDF(t *testing.T) {
	pl := NewPowerLaw(3, 2)
	for _, x := range []float64{2, 3, 5, 10} {
		got := pl.LogPDF(x)
		want := math.Log(pl.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("PowerLaw(3,2).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

// ---------------------------------------------------------------------------
// Integration: verify all continuous distributions integrate to ~1
// ---------------------------------------------------------------------------

func TestContinuousPDFIntegratesToOne(t *testing.T) {
	// Verify that numerical integration of PDF over a wide range is close to 1
	integrate := func(d Distribution, lo, hi float64, n int) float64 {
		h := (hi - lo) / float64(n)
		sum := d.PDF(lo) + d.PDF(hi)
		for i := 1; i < n; i++ {
			x := lo + float64(i)*h
			if i%2 == 0 {
				sum += 2 * d.PDF(x)
			} else {
				sum += 4 * d.PDF(x)
			}
		}
		return sum * h / 3
	}

	tests := []struct {
		name string
		d    Distribution
		lo   float64
		hi   float64
	}{
		{"Lognormal(0,1)", NewLognormal(0, 1), 0.001, 30},
		{"Weibull(2,1)", NewWeibull(2, 1), 0.001, 10},
		{"Pareto(3,1)", NewPareto(3, 1), 1, 100},
		{"Cauchy(0,1)", NewCauchy(0, 1), -100, 100},
		{"Laplace(0,1)", NewLaplace(0, 1), -20, 20},
		{"Logistic(0,1)", NewLogistic(0, 1), -20, 20},
		{"Gumbel(0,1)", NewGumbel(0, 1), -10, 30},
		{"Rayleigh(1)", NewRayleigh(1), 0.001, 10},
		{"HalfNormal(1)", NewHalfNormal(1), 0, 10},
		{"Wald(1,1)", NewWald(1, 1), 0.001, 30},
		// Levy omitted: heavy tail requires adaptive integration
		{"PowerLaw(2.5,1)", NewPowerLaw(2.5, 1), 1, 1000},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := integrate(tc.d, tc.lo, tc.hi, 10000)
			if !approxEqual(got, 1.0, 0.02) {
				t.Errorf("%s PDF integral = %v, want ~1.0", tc.name, got)
			}
		})
	}
}
