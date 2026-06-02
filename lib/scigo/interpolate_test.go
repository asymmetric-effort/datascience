//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// Interp1D Tests
// ---------------------------------------------------------------------------

func TestInterp1DLinear(t *testing.T) {
	x := []float64{0, 1, 2, 3}
	y := []float64{0, 1, 4, 9}
	f := Interp1D(x, y, "linear")

	// At data points.
	for i := range x {
		if !approxEqual(f(x[i]), y[i], 1e-14) {
			t.Errorf("Interp1D(%v) = %v, want %v", x[i], f(x[i]), y[i])
		}
	}

	// Between points.
	mid := f(0.5)
	if !approxEqual(mid, 0.5, 1e-14) {
		t.Errorf("Interp1D(0.5) = %v, want 0.5", mid)
	}

	mid = f(1.5)
	expected := 2.5 // linear interp between (1,1) and (2,4)
	if !approxEqual(mid, expected, 1e-14) {
		t.Errorf("Interp1D(1.5) = %v, want %v", mid, expected)
	}
}

func TestInterp1DNearest(t *testing.T) {
	x := []float64{0, 1, 2}
	y := []float64{10, 20, 30}
	f := Interp1D(x, y, "nearest")

	if !approxEqual(f(0.3), 10, 1e-14) {
		t.Errorf("nearest(0.3) = %v, want 10", f(0.3))
	}
	if !approxEqual(f(0.6), 20, 1e-14) {
		t.Errorf("nearest(0.6) = %v, want 20", f(0.6))
	}
	if !approxEqual(f(1.8), 30, 1e-14) {
		t.Errorf("nearest(1.8) = %v, want 30", f(1.8))
	}
}

func TestInterp1DExtrapolation(t *testing.T) {
	x := []float64{0, 1, 2}
	y := []float64{10, 20, 30}
	f := Interp1D(x, y, "linear")

	// Below range: clamp to first value.
	if !approxEqual(f(-1), 10, 1e-14) {
		t.Errorf("Interp1D(-1) = %v, want 10", f(-1))
	}
	// Above range: clamp to last value.
	if !approxEqual(f(5), 30, 1e-14) {
		t.Errorf("Interp1D(5) = %v, want 30", f(5))
	}
}

func TestInterp1DSinglePoint(t *testing.T) {
	f := Interp1D([]float64{1}, []float64{42}, "linear")
	if !approxEqual(f(100), 42, 1e-14) {
		t.Errorf("single point: got %v, want 42", f(100))
	}
}

// ---------------------------------------------------------------------------
// CubicSpline Tests
// ---------------------------------------------------------------------------

func TestCubicSpline(t *testing.T) {
	// Interpolate a quadratic: y = x^2.
	x := []float64{0, 1, 2, 3, 4}
	y := []float64{0, 1, 4, 9, 16}
	f := CubicSpline(x, y)

	// At data points.
	for i := range x {
		if !approxEqual(f(x[i]), y[i], 1e-8) {
			t.Errorf("CubicSpline(%v) = %v, want %v", x[i], f(x[i]), y[i])
		}
	}

	// Between points: cubic spline of x^2 should be accurate.
	if !approxEqual(f(0.5), 0.25, 0.1) {
		t.Errorf("CubicSpline(0.5) = %v, want ~0.25", f(0.5))
	}
	if !approxEqual(f(2.5), 6.25, 0.1) {
		t.Errorf("CubicSpline(2.5) = %v, want ~6.25", f(2.5))
	}
}

func TestCubicSplineLinearData(t *testing.T) {
	x := []float64{0, 1, 2, 3}
	y := []float64{0, 2, 4, 6}
	f := CubicSpline(x, y)

	// For linear data, cubic spline should match exactly.
	for _, xq := range []float64{0.5, 1.0, 1.5, 2.5} {
		expected := 2 * xq
		if !approxEqual(f(xq), expected, 1e-10) {
			t.Errorf("CubicSpline(%v) = %v, want %v", xq, f(xq), expected)
		}
	}
}

func TestCubicSplineTwoPoints(t *testing.T) {
	f := CubicSpline([]float64{0, 1}, []float64{0, 1})
	if !approxEqual(f(0.5), 0.5, 1e-14) {
		t.Errorf("CubicSpline two points(0.5) = %v, want 0.5", f(0.5))
	}
}

// ---------------------------------------------------------------------------
// BSpline Tests
// ---------------------------------------------------------------------------

func TestBSplineDegree1(t *testing.T) {
	x := []float64{0, 1, 2}
	y := []float64{0, 1, 4}
	f := BSpline(x, y, 1)

	if !approxEqual(f(0.5), 0.5, 1e-14) {
		t.Errorf("BSpline degree 1 (0.5) = %v, want 0.5", f(0.5))
	}
}

func TestBSplineDegree3(t *testing.T) {
	x := []float64{0, 1, 2, 3, 4}
	y := []float64{0, 1, 4, 9, 16}
	f := BSpline(x, y, 3)

	// Should match CubicSpline behavior.
	for _, xq := range x {
		if !approxEqual(f(xq), xq*xq, 0.2) {
			t.Errorf("BSpline degree 3 (%v) = %v, want ~%v", xq, f(xq), xq*xq)
		}
	}
}

func TestBSplineDegree2(t *testing.T) {
	x := []float64{0, 1, 2, 3}
	y := []float64{0, 1, 4, 9}
	f := BSpline(x, y, 2)

	// At data points, Lagrange interpolation should be exact.
	for i := range x {
		if !approxEqual(f(x[i]), y[i], 1e-8) {
			t.Errorf("BSpline degree 2 (%v) = %v, want %v", x[i], f(x[i]), y[i])
		}
	}
}

// ---------------------------------------------------------------------------
// RBFInterpolator Tests
// ---------------------------------------------------------------------------

func TestRBFInterpolator(t *testing.T) {
	// 1D points.
	points := [][]float64{{0}, {1}, {2}, {3}}
	values := []float64{0, 1, 4, 9}

	for _, kernel := range []string{"linear", "cubic", "gaussian", "multiquadric"} {
		f := RBFInterpolator(points, values, kernel)

		// At data points, should be exact (within tolerance).
		for i := range points {
			result := f(points[i])
			if !approxEqual(result, values[i], 1e-6) {
				t.Errorf("RBF(%s) at point %v = %v, want %v", kernel, points[i], result, values[i])
			}
		}
	}
}

func TestRBFInterpolator2D(t *testing.T) {
	points := [][]float64{
		{0, 0}, {1, 0}, {0, 1}, {1, 1},
	}
	values := []float64{0, 1, 1, 2}

	f := RBFInterpolator(points, values, "multiquadric")

	// At data points.
	for i := range points {
		result := f(points[i])
		if !approxEqual(result, values[i], 1e-6) {
			t.Errorf("RBF 2D at point %v = %v, want %v", points[i], result, values[i])
		}
	}
}

func TestRBFInterpolatorEmpty(t *testing.T) {
	f := RBFInterpolator(nil, nil, "linear")
	if !math.IsNaN(f([]float64{0})) {
		t.Error("expected NaN for empty RBF")
	}
}
