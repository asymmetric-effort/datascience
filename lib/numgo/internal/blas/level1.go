package blas

import "math"

// Ddot computes the dot product of two vectors with bounds validation.
func Ddot(n int, x []float64, incx int, y []float64, incy int) float64 {
	if n <= 0 {
		return 0
	}
	validateVector("Ddot x", x, n, incx)
	validateVector("Ddot y", y, n, incy)
	return DdotUnsafe(n, x, incx, y, incy)
}

// DdotUnsafe computes the dot product without bounds checks.
// Use only when slice lengths have already been validated.
func DdotUnsafe(n int, x []float64, incx int, y []float64, incy int) float64 {
	if n <= 0 {
		return 0
	}
	if incx == 1 && incy == 1 {
		return ddotUnit(n, x, y)
	}
	var sum float64
	ix, iy := 0, 0
	for i := 0; i < n; i++ {
		sum += x[ix] * y[iy]
		ix += incx
		iy += incy
	}
	return sum
}

// ddotUnit is the unit-stride fast path with 4x unrolling.
func ddotUnit(n int, x, y []float64) float64 {
	var s0, s1, s2, s3 float64
	m := n - n%4
	for i := 0; i < m; i += 4 {
		s0 += x[i] * y[i]
		s1 += x[i+1] * y[i+1]
		s2 += x[i+2] * y[i+2]
		s3 += x[i+3] * y[i+3]
	}
	sum := s0 + s1 + s2 + s3
	for i := m; i < n; i++ {
		sum += x[i] * y[i]
	}
	return sum
}

// Daxpy computes y = alpha*x + y with bounds validation.
func Daxpy(n int, alpha float64, x []float64, incx int, y []float64, incy int) {
	if n <= 0 || alpha == 0 {
		return
	}
	validateVector("Daxpy x", x, n, incx)
	validateVector("Daxpy y", y, n, incy)
	DaxpyUnsafe(n, alpha, x, incx, y, incy)
}

// DaxpyUnsafe computes y = alpha*x + y without bounds checks.
func DaxpyUnsafe(n int, alpha float64, x []float64, incx int, y []float64, incy int) {
	if n <= 0 || alpha == 0 {
		return
	}
	if incx == 1 && incy == 1 {
		m := n - n%4
		for i := 0; i < m; i += 4 {
			y[i] += alpha * x[i]
			y[i+1] += alpha * x[i+1]
			y[i+2] += alpha * x[i+2]
			y[i+3] += alpha * x[i+3]
		}
		for i := m; i < n; i++ {
			y[i] += alpha * x[i]
		}
		return
	}
	ix, iy := 0, 0
	for i := 0; i < n; i++ {
		y[iy] += alpha * x[ix]
		ix += incx
		iy += incy
	}
}

// Dscal scales a vector by a constant with bounds validation.
func Dscal(n int, alpha float64, x []float64, incx int) {
	if n <= 0 {
		return
	}
	validateVector("Dscal x", x, n, incx)
	DscalUnsafe(n, alpha, x, incx)
}

// DscalUnsafe scales a vector by a constant without bounds checks.
func DscalUnsafe(n int, alpha float64, x []float64, incx int) {
	if n <= 0 {
		return
	}
	if incx == 1 {
		m := n - n%4
		for i := 0; i < m; i += 4 {
			x[i] *= alpha
			x[i+1] *= alpha
			x[i+2] *= alpha
			x[i+3] *= alpha
		}
		for i := m; i < n; i++ {
			x[i] *= alpha
		}
		return
	}
	ix := 0
	for i := 0; i < n; i++ {
		x[ix] *= alpha
		ix += incx
	}
}

// Dnrm2 computes the Euclidean norm with bounds validation.
func Dnrm2(n int, x []float64, incx int) float64 {
	if n <= 0 {
		return 0
	}
	validateVector("Dnrm2 x", x, n, incx)
	return Dnrm2Unsafe(n, x, incx)
}

// Dnrm2Unsafe computes the Euclidean norm without bounds checks.
func Dnrm2Unsafe(n int, x []float64, incx int) float64 {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return math.Abs(x[0])
	}
	var sum, comp float64
	ix := 0
	for i := 0; i < n; i++ {
		v := x[ix]
		prod := v * v
		y := prod - comp
		t := sum + y
		comp = (t - sum) - y
		sum = t
		ix += incx
	}
	return math.Sqrt(sum)
}

// Dasum computes the sum of absolute values with bounds validation.
func Dasum(n int, x []float64, incx int) float64 {
	if n <= 0 {
		return 0
	}
	validateVector("Dasum x", x, n, incx)
	return DasumUnsafe(n, x, incx)
}

// DasumUnsafe computes the sum of absolute values without bounds checks.
func DasumUnsafe(n int, x []float64, incx int) float64 {
	if n <= 0 {
		return 0
	}
	if incx == 1 {
		var s0, s1, s2, s3 float64
		m := n - n%4
		for i := 0; i < m; i += 4 {
			s0 += math.Abs(x[i])
			s1 += math.Abs(x[i+1])
			s2 += math.Abs(x[i+2])
			s3 += math.Abs(x[i+3])
		}
		sum := s0 + s1 + s2 + s3
		for i := m; i < n; i++ {
			sum += math.Abs(x[i])
		}
		return sum
	}
	var sum float64
	ix := 0
	for i := 0; i < n; i++ {
		sum += math.Abs(x[ix])
		ix += incx
	}
	return sum
}

// Idamax returns the index of the element with the largest absolute value
// with bounds validation. Returns -1 if n <= 0.
func Idamax(n int, x []float64, incx int) int {
	if n <= 0 {
		return -1
	}
	validateVector("Idamax x", x, n, incx)
	return IdamaxUnsafe(n, x, incx)
}

// IdamaxUnsafe returns the index of the element with the largest absolute value
// without bounds checks.
func IdamaxUnsafe(n int, x []float64, incx int) int {
	if n <= 0 {
		return -1
	}
	if n == 1 {
		return 0
	}
	maxIdx := 0
	maxVal := math.Abs(x[0])
	ix := incx
	for i := 1; i < n; i++ {
		v := math.Abs(x[ix])
		if v > maxVal {
			maxVal = v
			maxIdx = i
		}
		ix += incx
	}
	return maxIdx
}
