package gpu

import (
	"math"
	"runtime"
)

// CPUBackend is a pure-Go compute backend that runs all operations on the CPU.
type CPUBackend struct{}

// NewCPUBackend creates a new CPU-based compute backend.
func NewCPUBackend() *CPUBackend {
	return &CPUBackend{}
}

// Name returns "cpu".
func (c *CPUBackend) Name() string {
	return "cpu"
}

// IsAvailable always returns true for the CPU backend.
func (c *CPUBackend) IsAvailable() bool {
	return true
}

// MatMul performs matrix multiplication of a (m x k) and b (k x n).
// Both a and b are in row-major order. Returns a flat slice of length m*n.
func (c *CPUBackend) MatMul(a, b []float64, m, k, n int) []float64 {
	result := make([]float64, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			var sum float64
			for p := 0; p < k; p++ {
				sum += a[i*k+p] * b[p*n+j]
			}
			result[i*n+j] = sum
		}
	}
	return result
}

// ElementWiseMul returns the element-wise product of a and b.
func (c *CPUBackend) ElementWiseMul(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] * b[i]
	}
	return result
}

// Sum returns the sum of all elements.
func (c *CPUBackend) Sum(a []float64) float64 {
	var s float64
	for _, v := range a {
		s += v
	}
	return s
}

// Normalize divides each element by the total sum, producing a distribution
// that sums to 1. Returns a new slice.
func (c *CPUBackend) Normalize(a []float64) []float64 {
	s := c.Sum(a)
	result := make([]float64, len(a))
	if s == 0 {
		return result
	}
	for i, v := range a {
		result[i] = v / s
	}
	return result
}

// FactorProduct computes the product of two discrete factors.
//
// Each factor is represented as a flat value array with an associated shape.
// The result shape must be provided and defines the dimensionality of the output.
//
// The algorithm treats the first factor's dimensions as the leading axes and
// the second factor's dimensions as the trailing axes of the result tensor,
// computing the outer product and storing it in row-major order.
func (c *CPUBackend) FactorProduct(aValues []float64, aShape []int, bValues []float64, bShape []int, resultShape []int) []float64 {
	resultSize := 1
	for _, d := range resultShape {
		resultSize *= d
	}
	result := make([]float64, resultSize)

	bSize := 1
	for _, d := range bShape {
		bSize *= d
	}

	for i, av := range aValues {
		for j, bv := range bValues {
			result[i*bSize+j] = av * bv
		}
	}
	return result
}

// Marginalize sums out the given axis from a multi-dimensional array.
// Returns the reduced values and the new shape with the axis removed.
//
// The values slice is interpreted as a row-major tensor with the given shape.
// The axis parameter specifies which dimension to sum over (0-indexed).
func (c *CPUBackend) Marginalize(values []float64, shape []int, axis int) ([]float64, []int) {
	ndim := len(shape)

	// Compute strides for the input tensor.
	strides := make([]int, ndim)
	strides[ndim-1] = 1
	for i := ndim - 2; i >= 0; i-- {
		strides[i] = strides[i+1] * shape[i+1]
	}

	// Build the new shape with the axis removed.
	newShape := make([]int, 0, ndim-1)
	for i, d := range shape {
		if i != axis {
			newShape = append(newShape, d)
		}
	}

	newSize := 1
	for _, d := range newShape {
		newSize *= d
	}
	result := make([]float64, newSize)

	// Compute new strides.
	newStrides := make([]int, len(newShape))
	if len(newStrides) > 0 {
		newStrides[len(newStrides)-1] = 1
		for i := len(newStrides) - 2; i >= 0; i-- {
			newStrides[i] = newStrides[i+1] * newShape[i+1]
		}
	}

	// Iterate over every element of the original tensor and accumulate
	// into the appropriate position in the result.
	totalSize := len(values)
	for flatIdx := 0; flatIdx < totalSize; flatIdx++ {
		remaining := flatIdx
		newFlatIdx := 0
		newDim := 0
		for d := 0; d < ndim; d++ {
			coord := remaining / strides[d]
			remaining %= strides[d]
			if d != axis {
				newFlatIdx += coord * newStrides[newDim]
				newDim++
			}
		}
		result[newFlatIdx] += values[flatIdx]
	}

	return result, newShape
}

// Close is a no-op for the CPU backend.
func (c *CPUBackend) Close() error {
	return nil
}

// --- Tensor operations ---

// ElementWiseAdd returns the element-wise sum of a and b.
func (c *CPUBackend) ElementWiseAdd(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] + b[i]
	}
	return result
}

// ElementWiseSub returns the element-wise difference a - b.
func (c *CPUBackend) ElementWiseSub(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] - b[i]
	}
	return result
}

// ElementWiseDiv returns the element-wise quotient a / b.
func (c *CPUBackend) ElementWiseDiv(a, b []float64) []float64 {
	result := make([]float64, len(a))
	for i := range a {
		result[i] = a[i] / b[i]
	}
	return result
}

// ScalarMul multiplies every element of a by the scalar s.
func (c *CPUBackend) ScalarMul(a []float64, s float64) []float64 {
	result := make([]float64, len(a))
	for i, v := range a {
		result[i] = v * s
	}
	return result
}

// ScalarAdd adds the scalar s to every element of a.
func (c *CPUBackend) ScalarAdd(a []float64, s float64) []float64 {
	result := make([]float64, len(a))
	for i, v := range a {
		result[i] = v + s
	}
	return result
}

// Exp returns the element-wise exponential of a.
func (c *CPUBackend) Exp(a []float64) []float64 {
	result := make([]float64, len(a))
	for i, v := range a {
		result[i] = math.Exp(v)
	}
	return result
}

// Log returns the element-wise natural logarithm of a.
func (c *CPUBackend) Log(a []float64) []float64 {
	result := make([]float64, len(a))
	for i, v := range a {
		result[i] = math.Log(v)
	}
	return result
}

// Sqrt returns the element-wise square root of a.
func (c *CPUBackend) Sqrt(a []float64) []float64 {
	result := make([]float64, len(a))
	for i, v := range a {
		result[i] = math.Sqrt(v)
	}
	return result
}

// Abs returns the element-wise absolute value of a.
func (c *CPUBackend) Abs(a []float64) []float64 {
	result := make([]float64, len(a))
	for i, v := range a {
		result[i] = math.Abs(v)
	}
	return result
}

// Max returns the maximum element in a.
// Panics if a is empty.
func (c *CPUBackend) Max(a []float64) float64 {
	m := a[0]
	for _, v := range a[1:] {
		if v > m {
			m = v
		}
	}
	return m
}

// Min returns the minimum element in a.
// Panics if a is empty.
func (c *CPUBackend) Min(a []float64) float64 {
	m := a[0]
	for _, v := range a[1:] {
		if v < m {
			m = v
		}
	}
	return m
}

// ArgMax returns the index of the maximum element in a.
// Panics if a is empty.
func (c *CPUBackend) ArgMax(a []float64) int {
	idx := 0
	for i, v := range a {
		if v > a[idx] {
			idx = i
		}
	}
	return idx
}

// ArgMin returns the index of the minimum element in a.
// Panics if a is empty.
func (c *CPUBackend) ArgMin(a []float64) int {
	idx := 0
	for i, v := range a {
		if v < a[idx] {
			idx = i
		}
	}
	return idx
}

// Dot returns the dot product of a and b.
func (c *CPUBackend) Dot(a, b []float64) float64 {
	var s float64
	for i := range a {
		s += a[i] * b[i]
	}
	return s
}

// --- Factor operations ---

// FactorReduce fixes the variable at the given axis to the specified index,
// returning the reduced values and the new shape with that axis removed.
func (c *CPUBackend) FactorReduce(values []float64, shape []int, axis int, index int) ([]float64, []int) {
	ndim := len(shape)

	strides := make([]int, ndim)
	strides[ndim-1] = 1
	for i := ndim - 2; i >= 0; i-- {
		strides[i] = strides[i+1] * shape[i+1]
	}

	newShape := make([]int, 0, ndim-1)
	for i, d := range shape {
		if i != axis {
			newShape = append(newShape, d)
		}
	}

	newSize := 1
	for _, d := range newShape {
		newSize *= d
	}
	result := make([]float64, newSize)

	newStrides := make([]int, len(newShape))
	if len(newStrides) > 0 {
		newStrides[len(newStrides)-1] = 1
		for i := len(newStrides) - 2; i >= 0; i-- {
			newStrides[i] = newStrides[i+1] * newShape[i+1]
		}
	}

	totalSize := len(values)
	for flatIdx := 0; flatIdx < totalSize; flatIdx++ {
		remaining := flatIdx
		axisCoord := -1
		newFlatIdx := 0
		newDim := 0
		for d := 0; d < ndim; d++ {
			coord := remaining / strides[d]
			remaining %= strides[d]
			if d == axis {
				axisCoord = coord
			} else {
				newFlatIdx += coord * newStrides[newDim]
				newDim++
			}
		}
		if axisCoord == index {
			result[newFlatIdx] = values[flatIdx]
		}
	}

	return result, newShape
}

// FactorMaximize takes the maximum over the given axis,
// returning the reduced values and the new shape with that axis removed.
func (c *CPUBackend) FactorMaximize(values []float64, shape []int, axis int) ([]float64, []int) {
	ndim := len(shape)

	strides := make([]int, ndim)
	strides[ndim-1] = 1
	for i := ndim - 2; i >= 0; i-- {
		strides[i] = strides[i+1] * shape[i+1]
	}

	newShape := make([]int, 0, ndim-1)
	for i, d := range shape {
		if i != axis {
			newShape = append(newShape, d)
		}
	}

	newSize := 1
	for _, d := range newShape {
		newSize *= d
	}
	result := make([]float64, newSize)
	initialized := make([]bool, newSize)

	newStrides := make([]int, len(newShape))
	if len(newStrides) > 0 {
		newStrides[len(newStrides)-1] = 1
		for i := len(newStrides) - 2; i >= 0; i-- {
			newStrides[i] = newStrides[i+1] * newShape[i+1]
		}
	}

	totalSize := len(values)
	for flatIdx := 0; flatIdx < totalSize; flatIdx++ {
		remaining := flatIdx
		newFlatIdx := 0
		newDim := 0
		for d := 0; d < ndim; d++ {
			coord := remaining / strides[d]
			remaining %= strides[d]
			if d != axis {
				newFlatIdx += coord * newStrides[newDim]
				newDim++
			}
		}
		if !initialized[newFlatIdx] || values[flatIdx] > result[newFlatIdx] {
			result[newFlatIdx] = values[flatIdx]
			initialized[newFlatIdx] = true
		}
	}

	return result, newShape
}

// LogSumExp computes log(sum(exp(a))) in a numerically stable way.
func (c *CPUBackend) LogSumExp(a []float64) float64 {
	m := c.Max(a)
	var s float64
	for _, v := range a {
		s += math.Exp(v - m)
	}
	return m + math.Log(s)
}

// Softmax returns the softmax of a: exp(a_i - max(a)) / sum(exp(a - max(a))).
func (c *CPUBackend) Softmax(a []float64) []float64 {
	m := c.Max(a)
	result := make([]float64, len(a))
	var s float64
	for i, v := range a {
		result[i] = math.Exp(v - m)
		s += result[i]
	}
	for i := range result {
		result[i] /= s
	}
	return result
}

// --- Batch operations ---

// BatchMatMul performs batched matrix multiplication. a and b contain
// batchSize matrices packed contiguously.
func (c *CPUBackend) BatchMatMul(a, b []float64, batchSize, m, k, n int) []float64 {
	aStride := m * k
	bStride := k * n
	rStride := m * n
	result := make([]float64, batchSize*rStride)
	for batch := 0; batch < batchSize; batch++ {
		aOff := batch * aStride
		bOff := batch * bStride
		rOff := batch * rStride
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				var sum float64
				for p := 0; p < k; p++ {
					sum += a[aOff+i*k+p] * b[bOff+p*n+j]
				}
				result[rOff+i*n+j] = sum
			}
		}
	}
	return result
}

// BatchNormalize normalizes batchSize vectors of length n packed contiguously.
func (c *CPUBackend) BatchNormalize(a []float64, batchSize, n int) []float64 {
	result := make([]float64, len(a))
	for batch := 0; batch < batchSize; batch++ {
		off := batch * n
		var s float64
		for i := 0; i < n; i++ {
			s += a[off+i]
		}
		if s == 0 {
			continue
		}
		for i := 0; i < n; i++ {
			result[off+i] = a[off+i] / s
		}
	}
	return result
}

// --- Memory management ---

// Alloc allocates a zeroed slice of the given size. On CPU this is make.
func (c *CPUBackend) Alloc(size int) []float64 {
	return make([]float64, size)
}

// Free releases device memory. On CPU this is a no-op; the GC handles it.
func (c *CPUBackend) Free(data []float64) {
	// no-op for CPU
}

// CopyToDevice copies data to device memory. On CPU this copies the slice.
func (c *CPUBackend) CopyToDevice(data []float64) []float64 {
	dst := make([]float64, len(data))
	copy(dst, data)
	return dst
}

// CopyFromDevice copies data from device memory. On CPU this copies the slice.
func (c *CPUBackend) CopyFromDevice(data []float64) []float64 {
	dst := make([]float64, len(data))
	copy(dst, data)
	return dst
}

// --- Device info ---

// DeviceCount returns 1 for the CPU backend.
func (c *CPUBackend) DeviceCount() int {
	return 1
}

// DeviceName returns "cpu" for index 0.
func (c *CPUBackend) DeviceName(index int) string {
	if index == 0 {
		return "cpu"
	}
	return ""
}

// MemoryUsed returns the approximate heap memory in use (bytes).
func (c *CPUBackend) MemoryUsed() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.Alloc)
}

// MemoryTotal returns 0 for the CPU backend (system RAM is not tracked).
func (c *CPUBackend) MemoryTotal() int64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int64(m.Sys)
}
