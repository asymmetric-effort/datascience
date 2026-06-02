package numgo

import "fmt"

// BroadcastShapes computes the result shape from two input shapes using numpy-style
// broadcasting rules:
//  1. If arrays differ in number of dimensions, pad the shorter shape with 1s on the left.
//  2. Dimensions with size 1 are stretched to match the other array's size.
//  3. If sizes differ and neither is 1, an error is returned.
func BroadcastShapes(a, b []int) ([]int, error) {
	ndim := len(a)
	if len(b) > ndim {
		ndim = len(b)
	}

	result := make([]int, ndim)
	for i := 0; i < ndim; i++ {
		// Index from the right: dimension (ndim-1-i) in the result corresponds
		// to dimension (len(x)-1-i) in each input, or 1 if out of range (left-pad).
		da := 1
		if idx := len(a) - 1 - i; idx >= 0 {
			da = a[idx]
		}
		db := 1
		if idx := len(b) - 1 - i; idx >= 0 {
			db = b[idx]
		}

		switch {
		case da == db:
			result[ndim-1-i] = da
		case da == 1:
			result[ndim-1-i] = db
		case db == 1:
			result[ndim-1-i] = da
		default:
			return nil, fmt.Errorf("numgo: cannot broadcast shapes %v and %v: dimension %d has size %d vs %d",
				a, b, ndim-1-i, da, db)
		}
	}
	return result, nil
}

// BroadcastTo broadcasts an NDArray to the given target shape, returning a new
// NDArray with data repeated as needed. The source array must be broadcast-compatible
// with the target shape (each source dimension must be 1 or equal to the target dimension).
func BroadcastTo(a *NDArray, shape []int) (*NDArray, error) {
	srcShape := a.shape

	if len(shape) < len(srcShape) {
		return nil, fmt.Errorf("numgo: cannot broadcast shape %v to %v: target has fewer dimensions", srcShape, shape)
	}

	// Left-pad source shape with 1s to match target ndim.
	padded := make([]int, len(shape))
	offset := len(shape) - len(srcShape)
	for i := 0; i < offset; i++ {
		padded[i] = 1
	}
	copy(padded[offset:], srcShape)

	// Validate compatibility.
	for i := range shape {
		if padded[i] != shape[i] && padded[i] != 1 {
			return nil, fmt.Errorf("numgo: cannot broadcast shape %v to %v: dimension %d is %d, cannot stretch to %d",
				srcShape, shape, i, padded[i], shape[i])
		}
	}

	// Compute strides for the padded source: 0 for broadcast (size-1) dims.
	// We need logical strides into the original flat data.
	paddedStrides := computeStrides(padded)
	srcStrides := make([]int, len(shape))
	for i := range srcStrides {
		if padded[i] == 1 {
			srcStrides[i] = 0 // broadcast: index stays 0 along this dim
		} else {
			srcStrides[i] = paddedStrides[i]
		}
	}

	totalSize := product(shape)
	data := make([]float64, totalSize)
	outStrides := computeStrides(shape)
	ndim := len(shape)

	for outIdx := 0; outIdx < totalSize; outIdx++ {
		// Decompose outIdx into multi-dim indices, then compute source flat index.
		srcFlat := 0
		rem := outIdx
		for d := 0; d < ndim; d++ {
			coord := rem / outStrides[d]
			rem %= outStrides[d]
			srcFlat += coord * srcStrides[d]
		}
		data[outIdx] = a.data[srcFlat]
	}

	return NewNDArray(shape, data), nil
}

// broadcastElementWise broadcasts a and b to a common shape, then applies op element-wise.
func broadcastElementWise(a, b *NDArray, op func(float64, float64) float64) *NDArray {
	resultShape, err := BroadcastShapes(a.shape, b.shape)
	if err != nil {
		panic(err.Error())
	}

	ab, err := BroadcastTo(a, resultShape)
	if err != nil {
		panic(err.Error())
	}
	bb, err := BroadcastTo(b, resultShape)
	if err != nil {
		panic(err.Error())
	}

	data := make([]float64, ab.Size())
	for i := range data {
		data[i] = op(ab.data[i], bb.data[i])
	}
	return NewNDArray(resultShape, data)
}
