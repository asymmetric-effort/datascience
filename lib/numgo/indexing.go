package numgo

import (
	"fmt"
)

// Take returns elements from a along the given axis at the specified indices.
// If axis < 0, operates on the flattened array.
func Take(a *NDArray, indices []int, axis int) (*NDArray, error) {
	if axis < 0 {
		// Operate on flattened array.
		data := make([]float64, len(indices))
		for i, idx := range indices {
			if idx < 0 || idx >= a.Size() {
				return nil, fmt.Errorf("numgo.Take: index %d out of range for size %d", idx, a.Size())
			}
			data[i] = a.data[idx]
		}
		return NewNDArray([]int{len(indices)}, data), nil
	}

	if axis >= a.Ndim() {
		return nil, fmt.Errorf("numgo.Take: axis %d out of range for %dD array", axis, a.Ndim())
	}

	// Build output shape: replace shape[axis] with len(indices).
	outShape := make([]int, a.Ndim())
	copy(outShape, a.shape)
	outShape[axis] = len(indices)
	outSize := product(outShape)

	data := make([]float64, outSize)
	outStrides := computeStrides(outShape)

	// Iterate over all output positions.
	outIdx := make([]int, a.Ndim())
	for flat := 0; flat < outSize; flat++ {
		// Decompose flat index.
		rem := flat
		for d := 0; d < a.Ndim(); d++ {
			outIdx[d] = rem / outStrides[d]
			rem %= outStrides[d]
		}
		// Map axis index through indices array.
		srcIdx := make([]int, a.Ndim())
		copy(srcIdx, outIdx)
		idx := indices[outIdx[axis]]
		if idx < 0 || idx >= a.shape[axis] {
			return nil, fmt.Errorf("numgo.Take: index %d out of range for axis %d with size %d", idx, axis, a.shape[axis])
		}
		srcIdx[axis] = idx
		data[flat] = a.Get(srcIdx...)
	}
	return NewNDArray(outShape, data), nil
}

// TakeAlongAxis gathers elements from a along the given axis using indices array.
// indices must have the same number of dimensions as a.
func TakeAlongAxis(a, indices *NDArray, axis int) (*NDArray, error) {
	if a.Ndim() != indices.Ndim() {
		return nil, fmt.Errorf("numgo.TakeAlongAxis: a and indices must have same ndim, got %d and %d", a.Ndim(), indices.Ndim())
	}
	if axis < 0 || axis >= a.Ndim() {
		return nil, fmt.Errorf("numgo.TakeAlongAxis: axis %d out of range for %dD array", axis, a.Ndim())
	}

	outShape := indices.Shape()
	outSize := indices.Size()
	data := make([]float64, outSize)
	outStrides := computeStrides(outShape)

	idx := make([]int, a.Ndim())
	for flat := 0; flat < outSize; flat++ {
		rem := flat
		for d := 0; d < a.Ndim(); d++ {
			idx[d] = rem / outStrides[d]
			rem %= outStrides[d]
		}
		// Replace axis dimension with value from indices.
		srcIdx := make([]int, a.Ndim())
		copy(srcIdx, idx)
		srcIdx[axis] = int(indices.data[flat])
		data[flat] = a.Get(srcIdx...)
	}
	return NewNDArray(outShape, data), nil
}

// Choose selects elements from choices based on indices.
// indices is a 1D array of ints selecting which choice array to pick from.
// All choices must have the same shape as indices.
func Choose(indices *NDArray, choices []*NDArray) (*NDArray, error) {
	if indices.Ndim() != 1 {
		return nil, fmt.Errorf("numgo.Choose: indices must be 1D")
	}
	n := indices.Size()
	data := make([]float64, n)
	for i := 0; i < n; i++ {
		ci := int(indices.data[i])
		if ci < 0 || ci >= len(choices) {
			return nil, fmt.Errorf("numgo.Choose: index %d out of range for %d choices", ci, len(choices))
		}
		if choices[ci].Size() != n {
			return nil, fmt.Errorf("numgo.Choose: choice %d size %d != indices size %d", ci, choices[ci].Size(), n)
		}
		data[i] = choices[ci].data[i]
	}
	return NewNDArray([]int{n}, data), nil
}

// Compress selects elements from a along the given axis where condition is true.
// If axis < 0, operates on the flattened array.
func Compress(condition []bool, a *NDArray, axis int) (*NDArray, error) {
	if axis < 0 {
		// Flatten.
		flat := a.Flatten()
		if len(condition) > flat.Size() {
			return nil, fmt.Errorf("numgo.Compress: condition length %d exceeds array size %d", len(condition), flat.Size())
		}
		var data []float64
		for i, c := range condition {
			if c {
				data = append(data, flat.data[i])
			}
		}
		if data == nil {
			data = []float64{}
		}
		return NewNDArray([]int{len(data)}, data), nil
	}

	if axis >= a.Ndim() {
		return nil, fmt.Errorf("numgo.Compress: axis %d out of range for %dD array", axis, a.Ndim())
	}
	if len(condition) > a.shape[axis] {
		return nil, fmt.Errorf("numgo.Compress: condition length %d exceeds axis %d size %d", len(condition), axis, a.shape[axis])
	}

	// Count true values.
	trueCount := 0
	for _, c := range condition {
		if c {
			trueCount++
		}
	}

	outShape := make([]int, a.Ndim())
	copy(outShape, a.shape)
	outShape[axis] = trueCount
	outSize := product(outShape)
	data := make([]float64, outSize)

	outStrides := computeStrides(outShape)
	outIdx := make([]int, a.Ndim())

	for flat := 0; flat < outSize; flat++ {
		rem := flat
		for d := 0; d < a.Ndim(); d++ {
			outIdx[d] = rem / outStrides[d]
			rem %= outStrides[d]
		}
		// Map the axis index to the original index.
		trueI := outIdx[axis]
		count := 0
		origIdx := -1
		for ci, c := range condition {
			if c {
				if count == trueI {
					origIdx = ci
					break
				}
				count++
			}
		}
		srcIdx := make([]int, a.Ndim())
		copy(srcIdx, outIdx)
		srcIdx[axis] = origIdx
		data[flat] = a.Get(srcIdx...)
	}
	return NewNDArray(outShape, data), nil
}

// Diagonal extracts the diagonal from a 2D array.
// offset > 0 selects superdiagonals, offset < 0 selects subdiagonals.
// axis1 and axis2 specify the 2D sub-array to extract from (for higher dims).
func Diagonal(a *NDArray, offset, axis1, axis2 int) (*NDArray, error) {
	if a.Ndim() < 2 {
		return nil, fmt.Errorf("numgo.Diagonal: input must be at least 2D")
	}
	if axis1 < 0 || axis1 >= a.Ndim() || axis2 < 0 || axis2 >= a.Ndim() || axis1 == axis2 {
		return nil, fmt.Errorf("numgo.Diagonal: invalid axes")
	}

	dim1 := a.shape[axis1]
	dim2 := a.shape[axis2]

	var diagLen int
	if offset >= 0 {
		diagLen = min(dim1, dim2-offset)
	} else {
		diagLen = min(dim1+offset, dim2)
	}
	if diagLen <= 0 {
		return NewNDArray([]int{0}, []float64{}), nil
	}

	// For 2D case (most common).
	if a.Ndim() == 2 {
		data := make([]float64, diagLen)
		for i := 0; i < diagLen; i++ {
			var r, c int
			if offset >= 0 {
				r, c = i, i+offset
			} else {
				r, c = i-offset, i
			}
			idx := make([]int, 2)
			if axis1 < axis2 {
				idx[axis1] = r
				idx[axis2] = c
			} else {
				idx[axis1] = r
				idx[axis2] = c
			}
			data[i] = a.Get(idx...)
		}
		return NewNDArray([]int{diagLen}, data), nil
	}

	return nil, fmt.Errorf("numgo.Diagonal: only 2D arrays are currently supported")
}

// Select returns elements from choices based on conditions.
// The first true condition selects the corresponding choice.
// If no condition is true, defaultVal is used.
func Select(conditions []*NDArray, choices []*NDArray, defaultVal float64) (*NDArray, error) {
	if len(conditions) != len(choices) {
		return nil, fmt.Errorf("numgo.Select: conditions and choices must have same length")
	}
	if len(conditions) == 0 {
		return nil, fmt.Errorf("numgo.Select: empty conditions")
	}
	n := conditions[0].Size()
	for i, c := range conditions {
		if c.Size() != n {
			return nil, fmt.Errorf("numgo.Select: condition %d size mismatch", i)
		}
	}
	for i, c := range choices {
		if c.Size() != n {
			return nil, fmt.Errorf("numgo.Select: choice %d size mismatch", i)
		}
	}

	data := make([]float64, n)
	for i := 0; i < n; i++ {
		data[i] = defaultVal
		for j, cond := range conditions {
			if cond.data[i] != 0 {
				data[i] = choices[j].data[i]
				break
			}
		}
	}
	shape := conditions[0].Shape()
	return NewNDArray(shape, data), nil
}

// AsStrided creates a view of the array with the given shape and strides.
// This is an unsafe operation: the returned array shares the same underlying data.
// Out-of-bounds strides can cause reads beyond the original data.
func AsStrided(a *NDArray, shape, strides []int) *NDArray {
	s := make([]int, len(shape))
	copy(s, shape)
	st := make([]int, len(strides))
	copy(st, strides)

	// Compute the required data size based on shape and strides.
	size := product(s)
	data := make([]float64, size)

	// Fill by iterating over all positions in the output.
	idx := make([]int, len(s))
	for flat := 0; flat < size; flat++ {
		// Compute source offset using custom strides.
		srcOffset := 0
		for d := 0; d < len(s); d++ {
			srcOffset += idx[d] * st[d]
		}
		if srcOffset >= 0 && srcOffset < len(a.data) {
			data[flat] = a.data[srcOffset]
		}

		// Increment indices (last dimension fastest).
		for d := len(s) - 1; d >= 0; d-- {
			idx[d]++
			if idx[d] < s[d] {
				break
			}
			idx[d] = 0
		}
	}

	return NewNDArray(s, data)
}
