package gpu

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
		// Decompose flatIdx into multi-dimensional indices.
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
