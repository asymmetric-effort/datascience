package numgo

import "math"

// Zeros returns an NDArray of the given shape filled with zeros.
func Zeros(shape ...int) *NDArray {
	return NewNDArray(shape, nil)
}

// Ones returns an NDArray of the given shape filled with ones.
func Ones(shape ...int) *NDArray {
	return Full(1.0, shape...)
}

// Full returns an NDArray of the given shape filled with the specified value.
func Full(value float64, shape ...int) *NDArray {
	size := product(shape)
	data := make([]float64, size)
	for i := range data {
		data[i] = value
	}
	return NewNDArray(shape, data)
}

// Eye returns a 2-D identity matrix of size n x n.
func Eye(n int) *NDArray {
	a := Zeros(n, n)
	for i := 0; i < n; i++ {
		a.Set(1.0, i, i)
	}
	return a
}

// FromSlice creates a 1-D NDArray from a float64 slice.
func FromSlice(data []float64) *NDArray {
	return NewNDArray([]int{len(data)}, data)
}

// FromSlice2D creates a 2-D NDArray from a slice of slices.
// All rows must have the same length.
func FromSlice2D(data [][]float64) *NDArray {
	if len(data) == 0 {
		return NewNDArray([]int{0, 0}, nil)
	}
	rows := len(data)
	cols := len(data[0])
	flat := make([]float64, 0, rows*cols)
	for i, row := range data {
		if len(row) != cols {
			panic("numgo: FromSlice2D requires all rows to have equal length")
		}
		_ = i
		flat = append(flat, row...)
	}
	return NewNDArray([]int{rows, cols}, flat)
}

// Empty returns a zero-initialized NDArray of the given shape.
// In Go, float64 slices are zero-initialized, so this is identical to Zeros.
func Empty(shape ...int) *NDArray {
	return Zeros(shape...)
}

// Identity returns an n x n identity matrix. Alias for Eye.
func Identity(n int) *NDArray {
	return Eye(n)
}

// Arange returns a 1D array of evenly spaced values in [start, stop) with the given step.
func Arange(start, stop, step float64) *NDArray {
	if step == 0 {
		panic("numgo.Arange: step must be non-zero")
	}
	var data []float64
	if step > 0 {
		for v := start; v < stop; v += step {
			data = append(data, v)
		}
	} else {
		for v := start; v > stop; v += step {
			data = append(data, v)
		}
	}
	if len(data) == 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	return NewNDArray([]int{len(data)}, data)
}

// Linspace returns num evenly spaced values over [start, stop].
func Linspace(start, stop float64, num int) *NDArray {
	if num <= 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	if num == 1 {
		return NewNDArray([]int{1}, []float64{start})
	}
	data := make([]float64, num)
	step := (stop - start) / float64(num-1)
	for i := 0; i < num; i++ {
		data[i] = start + float64(i)*step
	}
	return NewNDArray([]int{num}, data)
}

// Logspace returns num values spaced evenly on a log scale from 10^start to 10^stop.
func Logspace(start, stop float64, num int) *NDArray {
	lin := Linspace(start, stop, num)
	data := make([]float64, lin.Size())
	for i, v := range lin.data {
		data[i] = math.Pow(10, v)
	}
	return NewNDArray([]int{lin.Size()}, data)
}

// Geomspace returns num values spaced evenly on a log scale (geometric progression)
// from start to stop. Both start and stop must be positive.
func Geomspace(start, stop float64, num int) *NDArray {
	if start <= 0 || stop <= 0 {
		panic("numgo.Geomspace: start and stop must be positive")
	}
	if num <= 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	if num == 1 {
		return NewNDArray([]int{1}, []float64{start})
	}
	data := make([]float64, num)
	logStart := math.Log(start)
	logStop := math.Log(stop)
	step := (logStop - logStart) / float64(num-1)
	for i := 0; i < num; i++ {
		data[i] = math.Exp(logStart + float64(i)*step)
	}
	return NewNDArray([]int{num}, data)
}

// Meshgrid returns coordinate matrices from coordinate vectors.
// Given N 1D arrays, returns N NDArrays each with N dimensions.
func Meshgrid(xi ...*NDArray) []*NDArray {
	n := len(xi)
	if n == 0 {
		return nil
	}
	// All inputs must be 1D.
	sizes := make([]int, n)
	for i, x := range xi {
		if x.Ndim() != 1 {
			panic("numgo.Meshgrid: all inputs must be 1D")
		}
		sizes[i] = x.shape[0]
	}

	// Output shape for 2D case (most common): (len(x1), len(x0)) — "xy" indexing variant.
	// We use "ij" indexing: shape = (sizes[0], sizes[1], ..., sizes[n-1]).
	outShape := make([]int, n)
	copy(outShape, sizes)
	totalSize := productUnsafe(outShape)
	outStrides := computeStrides(outShape)

	result := make([]*NDArray, n)
	for g := 0; g < n; g++ {
		data := make([]float64, totalSize)
		idx := make([]int, n)
		for flat := 0; flat < totalSize; flat++ {
			// Decompose flat index.
			rem := flat
			for d := 0; d < n; d++ {
				idx[d] = rem / outStrides[d]
				rem %= outStrides[d]
			}
			data[flat] = xi[g].data[idx[g]]

			// No need to manually increment; flat loop handles it.
		}
		result[g] = NewNDArray(outShape, data)
	}
	return result
}

// Diag extracts or constructs a diagonal.
//   - If a is 1D, returns a 2D matrix with a on the k-th diagonal.
//   - If a is 2D, extracts the k-th diagonal as a 1D array.
func Diag(a *NDArray, k int) *NDArray {
	if a.Ndim() == 1 {
		n := a.shape[0] + abs(k)
		result := Zeros(n, n)
		for i := 0; i < a.shape[0]; i++ {
			if k >= 0 {
				result.Set(a.data[i], i, i+k)
			} else {
				result.Set(a.data[i], i-k, i)
			}
		}
		return result
	}
	if a.Ndim() == 2 {
		m, n := a.shape[0], a.shape[1]
		var diagLen int
		if k >= 0 {
			diagLen = min(m, n-k)
		} else {
			diagLen = min(m+k, n)
		}
		if diagLen <= 0 {
			return NewNDArray([]int{0}, []float64{})
		}
		data := make([]float64, diagLen)
		for i := 0; i < diagLen; i++ {
			if k >= 0 {
				data[i] = a.Get(i, i+k)
			} else {
				data[i] = a.Get(i-k, i)
			}
		}
		return NewNDArray([]int{diagLen}, data)
	}
	panic("numgo.Diag: input must be 1D or 2D")
}

// abs returns absolute value of an int.
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Diagflat creates a 2D array with the flattened input as the k-th diagonal.
func Diagflat(a *NDArray, k int) *NDArray {
	flat := a.Flatten()
	return Diag(flat, k)
}

// Tri returns an n x m matrix with ones at and below the k-th diagonal.
func Tri(n, m, k int) *NDArray {
	data := make([]float64, n*m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j <= i+k {
				data[i*m+j] = 1
			}
		}
	}
	return NewNDArray([]int{n, m}, data)
}

// Tril returns the lower triangle of a 2D array. Elements above the k-th diagonal are zeroed.
func Tril(a *NDArray, k int) *NDArray {
	if a.Ndim() != 2 {
		panic("numgo.Tril: input must be 2D")
	}
	m, n := a.shape[0], a.shape[1]
	data := make([]float64, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if j <= i+k {
				data[i*n+j] = a.Get(i, j)
			}
		}
	}
	return NewNDArray([]int{m, n}, data)
}

// Triu returns the upper triangle of a 2D array. Elements below the k-th diagonal are zeroed.
func Triu(a *NDArray, k int) *NDArray {
	if a.Ndim() != 2 {
		panic("numgo.Triu: input must be 2D")
	}
	m, n := a.shape[0], a.shape[1]
	data := make([]float64, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if j >= i+k {
				data[i*n+j] = a.Get(i, j)
			}
		}
	}
	return NewNDArray([]int{m, n}, data)
}

// Vander returns the Vandermonde matrix of a 1D input.
// Column j of the output is x^(n-1-j). If n <= 0, n defaults to len(x).
func Vander(x *NDArray, n int) *NDArray {
	if x.Ndim() != 1 {
		panic("numgo.Vander: input must be 1D")
	}
	m := x.shape[0]
	if n <= 0 {
		n = m
	}
	data := make([]float64, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			data[i*n+j] = math.Pow(x.data[i], float64(n-1-j))
		}
	}
	return NewNDArray([]int{m, n}, data)
}

// FromFunction constructs an NDArray by calling fn for each set of indices.
func FromFunction(shape []int, fn func(indices []int) float64) *NDArray {
	size := product(shape)
	data := make([]float64, size)
	strides := computeStrides(shape)
	idx := make([]int, len(shape))
	for flat := 0; flat < size; flat++ {
		rem := flat
		for d := 0; d < len(shape); d++ {
			idx[d] = rem / strides[d]
			rem %= strides[d]
		}
		// Copy indices so fn cannot mutate our slice.
		idxCopy := make([]int, len(shape))
		copy(idxCopy, idx)
		data[flat] = fn(idxCopy)
	}
	return NewNDArray(shape, data)
}

// FromIter constructs a 1D NDArray by reading count values from a channel.
func FromIter(ch <-chan float64, count int) *NDArray {
	data := make([]float64, 0, count)
	for i := 0; i < count; i++ {
		v, ok := <-ch
		if !ok {
			break
		}
		data = append(data, v)
	}
	return NewNDArray([]int{len(data)}, data)
}
