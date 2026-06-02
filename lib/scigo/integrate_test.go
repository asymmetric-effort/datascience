//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// Quad Tests
// ---------------------------------------------------------------------------

func TestQuadPolynomial(t *testing.T) {
	// Integrate x^2 from 0 to 1. Expected: 1/3.
	result, err := Quad(func(x float64) float64 { return x * x }, 0, 1)
	if err != nil {
		t.Fatalf("Quad: %v", err)
	}
	if !approxEqual(result, 1.0/3.0, 1e-8) {
		t.Errorf("Quad(x^2, 0, 1) = %v, want %v", result, 1.0/3.0)
	}
}

func TestQuadSin(t *testing.T) {
	// Integrate sin(x) from 0 to pi. Expected: 2.
	result, err := Quad(math.Sin, 0, math.Pi)
	if err != nil {
		t.Fatalf("Quad: %v", err)
	}
	if !approxEqual(result, 2.0, 1e-8) {
		t.Errorf("Quad(sin, 0, pi) = %v, want 2.0", result)
	}
}

func TestQuadExp(t *testing.T) {
	// Integrate e^x from 0 to 1. Expected: e - 1.
	result, err := Quad(math.Exp, 0, 1)
	if err != nil {
		t.Fatalf("Quad: %v", err)
	}
	if !approxEqual(result, math.E-1, 1e-8) {
		t.Errorf("Quad(exp, 0, 1) = %v, want %v", result, math.E-1)
	}
}

func TestQuadZeroInterval(t *testing.T) {
	result, err := Quad(func(x float64) float64 { return x }, 5, 5)
	if err != nil {
		t.Fatalf("Quad: %v", err)
	}
	if result != 0 {
		t.Errorf("Quad(x, 5, 5) = %v, want 0", result)
	}
}

// ---------------------------------------------------------------------------
// Dblquad Tests
// ---------------------------------------------------------------------------

func TestDblquad(t *testing.T) {
	// Integrate 1 over the unit square. Expected: 1.
	result, err := Dblquad(
		func(y, x float64) float64 { return 1 },
		0, 1,
		func(x float64) float64 { return 0 },
		func(x float64) float64 { return 1 },
	)
	if err != nil {
		t.Fatalf("Dblquad: %v", err)
	}
	if !approxEqual(result, 1.0, 1e-6) {
		t.Errorf("Dblquad(1, unit square) = %v, want 1.0", result)
	}
}

func TestDblquadTriangle(t *testing.T) {
	// Integrate 1 over the triangle 0<=x<=1, 0<=y<=x. Expected: 0.5.
	result, err := Dblquad(
		func(y, x float64) float64 { return 1 },
		0, 1,
		func(x float64) float64 { return 0 },
		func(x float64) float64 { return x },
	)
	if err != nil {
		t.Fatalf("Dblquad: %v", err)
	}
	if !approxEqual(result, 0.5, 1e-6) {
		t.Errorf("Dblquad(1, triangle) = %v, want 0.5", result)
	}
}

// ---------------------------------------------------------------------------
// Tplquad Tests
// ---------------------------------------------------------------------------

func TestTplquad(t *testing.T) {
	// Integrate 1 over the unit cube. Expected: 1.
	result, err := Tplquad(
		func(z, y, x float64) float64 { return 1 },
		0, 1,
		func(x float64) float64 { return 0 },
		func(x float64) float64 { return 1 },
		func(x, y float64) float64 { return 0 },
		func(x, y float64) float64 { return 1 },
	)
	if err != nil {
		t.Fatalf("Tplquad: %v", err)
	}
	if !approxEqual(result, 1.0, 1e-4) {
		t.Errorf("Tplquad(1, unit cube) = %v, want 1.0", result)
	}
}

// ---------------------------------------------------------------------------
// FixedQuad Tests
// ---------------------------------------------------------------------------

func TestFixedQuad(t *testing.T) {
	// Gauss-Legendre should be exact for polynomials of degree <= 2n-1.
	// n=3 => exact for degree <= 5.
	result := FixedQuad(func(x float64) float64 { return x * x * x * x }, 0, 1, 3)
	expected := 0.2 // 1/5
	if !approxEqual(result, expected, 1e-10) {
		t.Errorf("FixedQuad(x^4, 0, 1, 3) = %v, want %v", result, expected)
	}
}

func TestFixedQuadSin(t *testing.T) {
	result := FixedQuad(math.Sin, 0, math.Pi, 10)
	if !approxEqual(result, 2.0, 1e-8) {
		t.Errorf("FixedQuad(sin, 0, pi, 10) = %v, want 2.0", result)
	}
}

// ---------------------------------------------------------------------------
// Trapezoid Tests
// ---------------------------------------------------------------------------

func TestTrapezoid(t *testing.T) {
	// Integrate y = [1, 2, 3] with dx = 1. Expected: 4.
	result := Trapezoid([]float64{1, 2, 3}, 1)
	if !approxEqual(result, 4.0, 1e-14) {
		t.Errorf("Trapezoid([1,2,3], 1) = %v, want 4.0", result)
	}
}

func TestTrapezoidLinear(t *testing.T) {
	// y = x from 0 to 1 with 101 points. Expected: 0.5.
	n := 101
	dx := 1.0 / float64(n-1)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = float64(i) * dx
	}
	result := Trapezoid(y, dx)
	if !approxEqual(result, 0.5, 1e-10) {
		t.Errorf("Trapezoid(x, 0, 1) = %v, want 0.5", result)
	}
}

// ---------------------------------------------------------------------------
// Simpson Tests
// ---------------------------------------------------------------------------

func TestSimpson(t *testing.T) {
	// Integrate y = x^2 with uniform spacing.
	n := 101 // odd number of points = even number of intervals
	dx := 1.0 / float64(n-1)
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		x := float64(i) * dx
		y[i] = x * x
	}
	result := Simpson(y, dx)
	if !approxEqual(result, 1.0/3.0, 1e-6) {
		t.Errorf("Simpson(x^2, 0, 1) = %v, want %v", result, 1.0/3.0)
	}
}

func TestSimpsonConstant(t *testing.T) {
	y := []float64{5, 5, 5, 5, 5}
	result := Simpson(y, 1)
	if !approxEqual(result, 20.0, 1e-10) {
		t.Errorf("Simpson(5, dx=1, n=5) = %v, want 20.0", result)
	}
}

// ---------------------------------------------------------------------------
// Romberg Tests
// ---------------------------------------------------------------------------

func TestRomberg(t *testing.T) {
	result := Romberg(func(x float64) float64 { return x * x }, 0, 1)
	if !approxEqual(result, 1.0/3.0, 1e-10) {
		t.Errorf("Romberg(x^2, 0, 1) = %v, want %v", result, 1.0/3.0)
	}
}

func TestRombergExp(t *testing.T) {
	result := Romberg(math.Exp, 0, 1)
	if !approxEqual(result, math.E-1, 1e-10) {
		t.Errorf("Romberg(exp, 0, 1) = %v, want %v", result, math.E-1)
	}
}

// ---------------------------------------------------------------------------
// Odeint Tests
// ---------------------------------------------------------------------------

func TestOdeint(t *testing.T) {
	// dy/dt = -y, y(0) = 1 => y(t) = e^{-t}.
	f := func(y, tval float64) float64 { return -y }
	tspan := []float64{0.5, 1.0, 2.0}
	result := Odeint(f, 1.0, 0.0, tspan)

	for i, tval := range tspan {
		expected := math.Exp(-tval)
		if !approxEqual(result[i], expected, 1e-4) {
			t.Errorf("Odeint at t=%v: %v, want %v", tval, result[i], expected)
		}
	}
}

func TestOdeintLinear(t *testing.T) {
	// dy/dt = 1, y(0) = 0 => y(t) = t.
	f := func(y, tval float64) float64 { return 1 }
	tspan := []float64{1, 2, 3}
	result := Odeint(f, 0, 0, tspan)
	for i, tval := range tspan {
		if !approxEqual(result[i], tval, 1e-6) {
			t.Errorf("Odeint linear at t=%v: %v, want %v", tval, result[i], tval)
		}
	}
}

// ---------------------------------------------------------------------------
// SolveIVP Tests
// ---------------------------------------------------------------------------

func TestSolveIVP(t *testing.T) {
	// dy/dt = -y, y(0) = 1 => y(t) = e^{-t}.
	f := func(tval float64, y []float64) []float64 {
		return []float64{-y[0]}
	}
	times, states, err := SolveIVP(f, [2]float64{0, 2}, []float64{1.0})
	if err != nil {
		t.Fatalf("SolveIVP: %v", err)
	}
	if len(times) < 2 {
		t.Fatal("SolveIVP: too few time points")
	}

	// Check the final value.
	finalT := times[len(times)-1]
	finalY := states[len(states)-1][0]
	expected := math.Exp(-finalT)
	if !approxEqual(finalY, expected, 1e-3) {
		t.Errorf("SolveIVP at t=%v: %v, want %v", finalT, finalY, expected)
	}
}

func TestSolveIVPSystem(t *testing.T) {
	// Harmonic oscillator: y'' = -y => [y, v]' = [v, -y]
	// y(0) = 1, v(0) = 0 => y(t) = cos(t)
	f := func(tval float64, y []float64) []float64 {
		return []float64{y[1], -y[0]}
	}
	times, states, err := SolveIVP(f, [2]float64{0, math.Pi}, []float64{1, 0})
	if err != nil {
		t.Fatalf("SolveIVP: %v", err)
	}

	// At t = pi, y should be approximately -1.
	finalY := states[len(states)-1][0]
	if !approxEqual(finalY, -1, 0.05) {
		t.Errorf("SolveIVP harmonic at t=pi: y = %v, want ~-1", finalY)
	}

	_ = times
}

func TestSolveIVPZeroSpan(t *testing.T) {
	f := func(tval float64, y []float64) []float64 { return []float64{0} }
	times, states, err := SolveIVP(f, [2]float64{0, 0}, []float64{5})
	if err != nil {
		t.Fatalf("SolveIVP: %v", err)
	}
	if len(times) != 1 || states[0][0] != 5 {
		t.Error("SolveIVP zero span: unexpected result")
	}
}
