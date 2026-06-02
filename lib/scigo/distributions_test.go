//go:build unit

package scigo

import (
	"math"
	"math/rand"
	"testing"
)

// ---------------------------------------------------------------------------
// Normal Distribution Tests
// ---------------------------------------------------------------------------

func TestNormalInterface(t *testing.T) {
	// Verify Normal implements Distribution
	var _ Distribution = &Normal{}
}

func TestNewNormalPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewNormal(0, 0) should panic")
		}
	}()
	NewNormal(0, 0)
}

func TestNormalPDF(t *testing.T) {
	n := NewNormal(0, 1)
	// PDF at 0 for standard normal = 1/sqrt(2*pi)
	want := 1.0 / math.Sqrt(2*math.Pi)
	got := n.PDF(0)
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("Normal(0,1).PDF(0) = %v, want %v", got, want)
	}

	// PDF is symmetric
	if !approxEqual(n.PDF(1), n.PDF(-1), 1e-12) {
		t.Error("Standard normal PDF should be symmetric")
	}

	// Non-standard normal
	n2 := NewNormal(5, 2)
	got = n2.PDF(5)
	want = 1.0 / (2 * math.Sqrt(2*math.Pi))
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("Normal(5,2).PDF(5) = %v, want %v", got, want)
	}
}

func TestNormalCDF(t *testing.T) {
	n := NewNormal(0, 1)

	tests := []struct {
		x, want float64
	}{
		{0, 0.5},
		{math.Inf(1), 1.0},
		{math.Inf(-1), 0.0},
	}
	for _, tc := range tests {
		got := n.CDF(tc.x)
		if !approxEqual(got, tc.want, 1e-12) {
			t.Errorf("Normal(0,1).CDF(%v) = %v, want %v", tc.x, got, tc.want)
		}
	}

	// CDF(1.96) ≈ 0.975
	got := n.CDF(1.96)
	if !approxEqual(got, 0.975, 1e-3) {
		t.Errorf("Normal(0,1).CDF(1.96) = %v, want ~0.975", got)
	}
}

func TestNormalPPF(t *testing.T) {
	n := NewNormal(0, 1)

	// PPF(0.5) = 0 for standard normal
	got := n.PPF(0.5)
	if !approxEqual(got, 0, 1e-10) {
		t.Errorf("Normal(0,1).PPF(0.5) = %v, want 0", got)
	}

	// PPF(CDF(x)) = x (round-trip)
	for _, x := range []float64{-3, -1, 0, 0.5, 1, 2.5} {
		p := n.CDF(x)
		got := n.PPF(p)
		if !approxEqual(got, x, 1e-8) {
			t.Errorf("Normal(0,1).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}

	// Boundary values
	if !math.IsInf(n.PPF(0), -1) {
		t.Error("PPF(0) should be -Inf")
	}
	if !math.IsInf(n.PPF(1), 1) {
		t.Error("PPF(1) should be +Inf")
	}
}

func TestNormalLogPDF(t *testing.T) {
	n := NewNormal(0, 1)
	// LogPDF should equal log(PDF)
	for _, x := range []float64{-2, -1, 0, 1, 2} {
		got := n.LogPDF(x)
		want := math.Log(n.PDF(x))
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Normal(0,1).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestNormalMeanVar(t *testing.T) {
	n := NewNormal(3, 2)
	if n.Mean() != 3 {
		t.Errorf("Normal(3,2).Mean() = %v, want 3", n.Mean())
	}
	if !approxEqual(n.Var(), 4, 1e-12) {
		t.Errorf("Normal(3,2).Var() = %v, want 4", n.Var())
	}
}

func TestNormalSample(t *testing.T) {
	n := NewNormal(5, 2)
	rng := rand.New(rand.NewSource(42))
	samples := n.Sample(rng, 10000)

	if len(samples) != 10000 {
		t.Fatalf("Expected 10000 samples, got %d", len(samples))
	}

	// Check empirical mean and variance are close
	sum := 0.0
	for _, s := range samples {
		sum += s
	}
	mean := sum / float64(len(samples))
	if !approxEqual(mean, 5, 0.1) {
		t.Errorf("Sample mean = %v, want ~5", mean)
	}

	sumSq := 0.0
	for _, s := range samples {
		d := s - mean
		sumSq += d * d
	}
	variance := sumSq / float64(len(samples)-1)
	if !approxEqual(variance, 4, 0.3) {
		t.Errorf("Sample variance = %v, want ~4", variance)
	}
}

// ---------------------------------------------------------------------------
// Chi-Squared Distribution Tests
// ---------------------------------------------------------------------------

func TestChiSquaredInterface(t *testing.T) {
	var _ Distribution = &ChiSquared{}
}

func TestNewChiSquaredPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewChiSquared(0) should panic")
		}
	}()
	NewChiSquared(0)
}

func TestChiSquaredMeanVar(t *testing.T) {
	tests := []struct {
		df, wantMean, wantVar float64
	}{
		{1, 1, 2},
		{2, 2, 4},
		{5, 5, 10},
		{10, 10, 20},
	}
	for _, tc := range tests {
		c := NewChiSquared(tc.df)
		if c.Mean() != tc.wantMean {
			t.Errorf("ChiSquared(%v).Mean() = %v, want %v", tc.df, c.Mean(), tc.wantMean)
		}
		if c.Var() != tc.wantVar {
			t.Errorf("ChiSquared(%v).Var() = %v, want %v", tc.df, c.Var(), tc.wantVar)
		}
	}
}

func TestChiSquaredPDF(t *testing.T) {
	// Chi-squared(2) PDF = 0.5 * exp(-x/2), so PDF(0) = 0.5
	c := NewChiSquared(2)
	got := c.PDF(0)
	if !approxEqual(got, 0.5, 1e-10) {
		t.Errorf("ChiSquared(2).PDF(0) = %v, want 0.5", got)
	}

	// PDF at x=2 for df=2: 0.5 * exp(-1)
	got = c.PDF(2)
	want := 0.5 * math.Exp(-1)
	if !approxEqual(got, want, 1e-10) {
		t.Errorf("ChiSquared(2).PDF(2) = %v, want %v", got, want)
	}

	// PDF for negative x should be 0
	if c.PDF(-1) != 0 {
		t.Error("ChiSquared PDF for negative x should be 0")
	}
}

func TestChiSquaredCDF(t *testing.T) {
	// CDF(0) = 0
	c := NewChiSquared(2)
	if c.CDF(0) != 0 {
		t.Errorf("ChiSquared(2).CDF(0) = %v, want 0", c.CDF(0))
	}

	// For df=2, CDF(x) = 1 - exp(-x/2)
	for _, x := range []float64{0.5, 1, 2, 5, 10} {
		got := c.CDF(x)
		want := 1 - math.Exp(-x/2)
		if !approxEqual(got, want, 1e-8) {
			t.Errorf("ChiSquared(2).CDF(%v) = %v, want %v", x, got, want)
		}
	}

	// CDF for negative should be 0
	if c.CDF(-1) != 0 {
		t.Error("ChiSquared CDF for negative x should be 0")
	}
}

func TestChiSquaredPPF(t *testing.T) {
	c := NewChiSquared(2)

	// Round-trip: PPF(CDF(x)) = x
	for _, x := range []float64{0.5, 1, 2, 5, 10} {
		p := c.CDF(x)
		got := c.PPF(p)
		if !approxEqual(got, x, 1e-6) {
			t.Errorf("ChiSquared(2).PPF(CDF(%v)) = %v, want %v", x, got, x)
		}
	}

	// Boundary
	if c.PPF(0) != 0 {
		t.Errorf("ChiSquared PPF(0) = %v, want 0", c.PPF(0))
	}
	if !math.IsInf(c.PPF(1), 1) {
		t.Error("ChiSquared PPF(1) should be +Inf")
	}
}

func TestChiSquaredLogPDF(t *testing.T) {
	c := NewChiSquared(5)
	for _, x := range []float64{0.5, 1, 2, 5, 10} {
		got := c.LogPDF(x)
		want := math.Log(c.PDF(x))
		if !approxEqual(got, want, 1e-10) {
			t.Errorf("ChiSquared(5).LogPDF(%v) = %v, want %v", x, got, want)
		}
	}
}

func TestChiSquaredSurvivalFunction(t *testing.T) {
	c := NewChiSquared(2)
	// SF = 1 - CDF
	for _, x := range []float64{0.5, 1, 2, 5, 10} {
		got := c.SurvivalFunction(x)
		want := 1 - c.CDF(x)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("ChiSquared(2).SurvivalFunction(%v) = %v, want %v", x, got, want)
		}
	}

	// For df=2, SF(x) = exp(-x/2)
	for _, x := range []float64{1, 2, 5} {
		got := c.SurvivalFunction(x)
		want := math.Exp(-x / 2)
		if !approxEqual(got, want, 1e-8) {
			t.Errorf("ChiSquared(2).SurvivalFunction(%v) = %v, want %v (exp form)", x, got, want)
		}
	}
}
