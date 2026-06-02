package scigo

import "math"

// Interp1D returns a 1D interpolation function for the given (x, y) data points.
// Supported kinds: "linear" (default), "nearest".
// The returned function extrapolates using the boundary values for out-of-range inputs.
func Interp1D(x, y []float64, kind string) func(float64) float64 {
	n := len(x)
	if n == 0 || len(y) != n {
		return func(float64) float64 { return math.NaN() }
	}
	if n == 1 {
		v := y[0]
		return func(float64) float64 { return v }
	}

	// Sort by x (copy to avoid mutating input).
	xs := make([]float64, n)
	ys := make([]float64, n)
	copy(xs, x)
	copy(ys, y)
	sortPaired(xs, ys)

	switch kind {
	case "nearest":
		return func(xq float64) float64 {
			if xq <= xs[0] {
				return ys[0]
			}
			if xq >= xs[n-1] {
				return ys[n-1]
			}
			idx := searchSorted(xs, xq)
			if idx == 0 {
				return ys[0]
			}
			if math.Abs(xq-xs[idx-1]) <= math.Abs(xq-xs[idx]) {
				return ys[idx-1]
			}
			return ys[idx]
		}
	default: // "linear"
		return func(xq float64) float64 {
			if xq <= xs[0] {
				return ys[0]
			}
			if xq >= xs[n-1] {
				return ys[n-1]
			}
			idx := searchSorted(xs, xq)
			if idx == 0 {
				idx = 1
			}
			t := (xq - xs[idx-1]) / (xs[idx] - xs[idx-1])
			return ys[idx-1] + t*(ys[idx]-ys[idx-1])
		}
	}
}

// CubicSpline returns a natural cubic spline interpolation function for the given data.
// Natural boundary conditions: S”(x_0) = S”(x_n) = 0.
func CubicSpline(x, y []float64) func(float64) float64 {
	n := len(x)
	if n < 2 || len(y) != n {
		return func(float64) float64 { return math.NaN() }
	}

	// Sort by x.
	xs := make([]float64, n)
	ys := make([]float64, n)
	copy(xs, x)
	copy(ys, y)
	sortPaired(xs, ys)

	// Compute intervals.
	h := make([]float64, n-1)
	for i := 0; i < n-1; i++ {
		h[i] = xs[i+1] - xs[i]
	}

	// Set up tridiagonal system for second derivatives.
	// Natural spline: c[0] = c[n-1] = 0.
	if n == 2 {
		// Linear interpolation.
		return Interp1D(xs, ys, "linear")
	}

	// System of equations for c[1..n-2].
	m := n - 2
	diag := make([]float64, m)
	upper := make([]float64, m)
	lower := make([]float64, m)
	rhs := make([]float64, m)

	for i := 0; i < m; i++ {
		diag[i] = 2 * (h[i] + h[i+1])
		rhs[i] = 3 * ((ys[i+2]-ys[i+1])/h[i+1] - (ys[i+1]-ys[i])/h[i])
		if i < m-1 {
			upper[i] = h[i+1]
		}
		if i > 0 {
			lower[i] = h[i]
		}
	}

	// Solve tridiagonal system (Thomas algorithm).
	c := make([]float64, n)
	cInner := thomasSolve(lower, diag, upper, rhs)
	for i := 0; i < m; i++ {
		c[i+1] = cInner[i]
	}

	// Compute b and d coefficients.
	b := make([]float64, n-1)
	d := make([]float64, n-1)
	for i := 0; i < n-1; i++ {
		b[i] = (ys[i+1]-ys[i])/h[i] - h[i]*(2*c[i]+c[i+1])/3
		d[i] = (c[i+1] - c[i]) / (3 * h[i])
	}

	return func(xq float64) float64 {
		if xq <= xs[0] {
			return ys[0]
		}
		if xq >= xs[n-1] {
			return ys[n-1]
		}
		idx := searchSorted(xs, xq) - 1
		if idx < 0 {
			idx = 0
		}
		if idx >= n-1 {
			idx = n - 2
		}
		dx := xq - xs[idx]
		return ys[idx] + b[idx]*dx + c[idx]*dx*dx + d[idx]*dx*dx*dx
	}
}

// BSpline returns a simplified B-spline interpolation function.
// Uses local polynomial fitting of the given degree (1=linear, 3=cubic).
// For degree 3, delegates to CubicSpline. For other degrees, uses local polynomial.
func BSpline(x, y []float64, degree int) func(float64) float64 {
	if degree == 3 {
		return CubicSpline(x, y)
	}
	if degree == 1 {
		return Interp1D(x, y, "linear")
	}

	n := len(x)
	if n < degree+1 || len(y) != n {
		return func(float64) float64 { return math.NaN() }
	}

	xs := make([]float64, n)
	ys := make([]float64, n)
	copy(xs, x)
	copy(ys, y)
	sortPaired(xs, ys)

	// For general degree, use local polynomial interpolation.
	return func(xq float64) float64 {
		if xq <= xs[0] {
			return ys[0]
		}
		if xq >= xs[n-1] {
			return ys[n-1]
		}

		// Find the nearest points.
		idx := searchSorted(xs, xq)
		// Select degree+1 points centered on idx.
		half := degree / 2
		start := idx - half
		if start < 0 {
			start = 0
		}
		end := start + degree + 1
		if end > n {
			end = n
			start = end - degree - 1
			if start < 0 {
				start = 0
			}
		}

		// Lagrange interpolation over the selected points.
		result := 0.0
		pts := end - start
		for i := 0; i < pts; i++ {
			li := 1.0
			for j := 0; j < pts; j++ {
				if i != j {
					li *= (xq - xs[start+j]) / (xs[start+i] - xs[start+j])
				}
			}
			result += ys[start+i] * li
		}
		return result
	}
}

// RBFInterpolator returns a radial basis function interpolator.
// points is an Nx D matrix of N points in D dimensions, values has length N.
// Supported kernels: "linear", "cubic", "gaussian", "multiquadric".
func RBFInterpolator(points [][]float64, values []float64, kernel string) func([]float64) float64 {
	n := len(points)
	if n == 0 || len(values) != n {
		return func([]float64) float64 { return math.NaN() }
	}

	rbf := getRBFKernel(kernel)

	// Build the interpolation matrix and solve for weights.
	// A[i][j] = rbf(||points[i] - points[j]||)
	a := make([][]float64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			a[i][j] = rbf(pointDist(points[i], points[j]))
		}
	}

	// Solve A * w = values using LU decomposition.
	lu, piv, err := LUFactor(a)
	if err != nil {
		return func([]float64) float64 { return math.NaN() }
	}
	weights, err := LUSolve(lu, piv, values)
	if err != nil {
		return func([]float64) float64 { return math.NaN() }
	}

	ptsCopy := make([][]float64, n)
	for i := range points {
		ptsCopy[i] = make([]float64, len(points[i]))
		copy(ptsCopy[i], points[i])
	}

	return func(xq []float64) float64 {
		result := 0.0
		for i := 0; i < n; i++ {
			result += weights[i] * rbf(pointDist(xq, ptsCopy[i]))
		}
		return result
	}
}

// ---------------------------------------------------------------------------
// Interpolation helpers
// ---------------------------------------------------------------------------

func searchSorted(xs []float64, xq float64) int {
	lo, hi := 0, len(xs)
	for lo < hi {
		mid := (lo + hi) / 2
		if xs[mid] < xq {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	return lo
}

func sortPaired(x, y []float64) {
	n := len(x)
	// Simple insertion sort (data is typically small or already sorted).
	for i := 1; i < n; i++ {
		kx, ky := x[i], y[i]
		j := i - 1
		for j >= 0 && x[j] > kx {
			x[j+1] = x[j]
			y[j+1] = y[j]
			j--
		}
		x[j+1] = kx
		y[j+1] = ky
	}
}

func thomasSolve(lower, diag, upper, rhs []float64) []float64 {
	n := len(diag)
	if n == 0 {
		return nil
	}

	c := make([]float64, n)
	d := make([]float64, n)
	copy(c, upper)
	copy(d, rhs)

	// Forward sweep.
	for i := 1; i < n; i++ {
		m := lower[i] / diag[i-1]
		diag[i] -= m * c[i-1]
		d[i] -= m * d[i-1]
	}

	// Back substitution.
	x := make([]float64, n)
	x[n-1] = d[n-1] / diag[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = (d[i] - c[i]*x[i+1]) / diag[i]
	}

	return x
}

func pointDist(a, b []float64) float64 {
	s := 0.0
	for i := range a {
		if i < len(b) {
			d := a[i] - b[i]
			s += d * d
		}
	}
	return math.Sqrt(s)
}

func getRBFKernel(name string) func(float64) float64 {
	switch name {
	case "cubic":
		return func(r float64) float64 { return r * r * r }
	case "gaussian":
		return func(r float64) float64 { return math.Exp(-r * r) }
	case "multiquadric":
		return func(r float64) float64 { return math.Sqrt(1 + r*r) }
	default: // "linear"
		return func(r float64) float64 { return r }
	}
}
