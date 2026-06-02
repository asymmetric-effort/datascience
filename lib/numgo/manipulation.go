package numgo

import (
	"fmt"
)

// Ravel returns a contiguous flattened copy of the array.
func Ravel(a *NDArray) *NDArray {
	return NewNDArray([]int{a.Size()}, a.data)
}

// Swapaxes returns a new array with two axes swapped.
func Swapaxes(a *NDArray, axis1, axis2 int) *NDArray {
	ndim := a.Ndim()
	if axis1 < 0 || axis1 >= ndim {
		panic(fmt.Sprintf("numgo: axis1 %d out of range for %d dimensions", axis1, ndim))
	}
	if axis2 < 0 || axis2 >= ndim {
		panic(fmt.Sprintf("numgo: axis2 %d out of range for %d dimensions", axis2, ndim))
	}
	// Build a permutation that swaps axis1 and axis2.
	perm := make([]int, ndim)
	for i := range perm {
		perm[i] = i
	}
	perm[axis1], perm[axis2] = perm[axis2], perm[axis1]
	return transpose(a, perm)
}

// Moveaxis moves an axis from source to destination position.
func Moveaxis(a *NDArray, source, destination int) *NDArray {
	ndim := a.Ndim()
	if source < 0 || source >= ndim {
		panic(fmt.Sprintf("numgo: source axis %d out of range for %d dimensions", source, ndim))
	}
	if destination < 0 || destination >= ndim {
		panic(fmt.Sprintf("numgo: destination axis %d out of range for %d dimensions", destination, ndim))
	}
	// Build permutation: remove source, insert at destination.
	order := make([]int, 0, ndim)
	for i := 0; i < ndim; i++ {
		if i != source {
			order = append(order, i)
		}
	}
	// Insert source at destination position.
	perm := make([]int, ndim)
	copy(perm[:destination], order[:destination])
	perm[destination] = source
	copy(perm[destination+1:], order[destination:])
	return transpose(a, perm)
}

// Rollaxis rolls the specified axis backwards until it lies at start.
func Rollaxis(a *NDArray, axis, start int) *NDArray {
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		panic(fmt.Sprintf("numgo: axis %d out of range for %d dimensions", axis, ndim))
	}
	if start < 0 || start > ndim {
		panic(fmt.Sprintf("numgo: start %d out of range for %d dimensions", start, ndim))
	}
	if start > axis {
		start--
	}
	if axis == start {
		return a.Copy()
	}
	return Moveaxis(a, axis, start)
}

// transpose performs a generalized transpose given a permutation of axes.
func transpose(a *NDArray, perm []int) *NDArray {
	ndim := a.Ndim()
	newShape := make([]int, ndim)
	for i, p := range perm {
		newShape[i] = a.shape[p]
	}
	result := NewNDArray(newShape, nil)
	srcIndices := make([]int, ndim)
	dstIndices := make([]int, ndim)
	for flat := 0; flat < a.Size(); flat++ {
		// Decompose flat index into source indices.
		rem := flat
		for d := 0; d < ndim; d++ {
			srcIndices[d] = rem / a.strides[d]
			rem %= a.strides[d]
		}
		// Apply permutation.
		for d := 0; d < ndim; d++ {
			dstIndices[d] = srcIndices[perm[d]]
		}
		result.Set(a.data[flat], dstIndices...)
	}
	return result
}

// ExpandDims inserts a new axis of size 1 at the given position.
func ExpandDims(a *NDArray, axis int) *NDArray {
	ndim := a.Ndim()
	if axis < 0 || axis > ndim {
		panic(fmt.Sprintf("numgo: axis %d out of range for %d+1 dimensions", axis, ndim))
	}
	newShape := make([]int, ndim+1)
	copy(newShape[:axis], a.shape[:axis])
	newShape[axis] = 1
	copy(newShape[axis+1:], a.shape[axis:])
	return NewNDArray(newShape, a.data)
}

// Squeeze removes all size-1 dimensions from the array.
func Squeeze(a *NDArray) *NDArray {
	var newShape []int
	for _, s := range a.shape {
		if s != 1 {
			newShape = append(newShape, s)
		}
	}
	if len(newShape) == 0 {
		newShape = []int{1}
	}
	return NewNDArray(newShape, a.data)
}

// Concatenate joins a sequence of arrays along an existing axis.
func Concatenate(arrays []*NDArray, axis int) (*NDArray, error) {
	if len(arrays) == 0 {
		return nil, fmt.Errorf("numgo: need at least one array to concatenate")
	}
	ndim := arrays[0].Ndim()
	if axis < 0 || axis >= ndim {
		return nil, fmt.Errorf("numgo: axis %d out of range for %d dimensions", axis, ndim)
	}
	// Validate that all arrays have same shape except along axis.
	for i := 1; i < len(arrays); i++ {
		if arrays[i].Ndim() != ndim {
			return nil, fmt.Errorf("numgo: all arrays must have same number of dimensions")
		}
		for d := 0; d < ndim; d++ {
			if d != axis && arrays[i].shape[d] != arrays[0].shape[d] {
				return nil, fmt.Errorf("numgo: shape mismatch on axis %d: %d vs %d", d, arrays[0].shape[d], arrays[i].shape[d])
			}
		}
	}
	// Compute result shape.
	newShape := make([]int, ndim)
	copy(newShape, arrays[0].shape)
	for i := 1; i < len(arrays); i++ {
		newShape[axis] += arrays[i].shape[axis]
	}

	result := NewNDArray(newShape, nil)
	dstIndices := make([]int, ndim)
	for _, arr := range arrays {
		srcIndices := make([]int, ndim)
		for flat := 0; flat < arr.Size(); flat++ {
			// Decompose flat into source indices.
			rem := flat
			for d := 0; d < ndim; d++ {
				srcIndices[d] = rem / arr.strides[d]
				rem %= arr.strides[d]
			}
			// Map to destination: offset along the concat axis.
			copy(dstIndices, srcIndices)
			dstIndices[axis] = srcIndices[axis] + (newShape[axis] - newShape[axis]) // placeholder
			result.Set(arr.data[flat], dstIndices...)
		}
		// After processing each array, shift the axis offset for next array.
		// We need a different approach: track cumulative offset.
	}
	// Rewrite with offset tracking.
	result = NewNDArray(newShape, nil)
	axisOffset := 0
	for _, arr := range arrays {
		srcIndices := make([]int, ndim)
		for flat := 0; flat < arr.Size(); flat++ {
			rem := flat
			for d := 0; d < ndim; d++ {
				srcIndices[d] = rem / arr.strides[d]
				rem %= arr.strides[d]
			}
			copy(dstIndices, srcIndices)
			dstIndices[axis] += axisOffset
			result.Set(arr.data[flat], dstIndices...)
		}
		axisOffset += arr.shape[axis]
	}
	return result, nil
}

// Stack joins a sequence of arrays along a new axis.
func Stack(arrays []*NDArray, axis int) (*NDArray, error) {
	if len(arrays) == 0 {
		return nil, fmt.Errorf("numgo: need at least one array to stack")
	}
	ndim := arrays[0].Ndim()
	if axis < 0 || axis > ndim {
		return nil, fmt.Errorf("numgo: axis %d out of range for %d+1 dimensions", axis, ndim)
	}
	// Validate all arrays have the same shape.
	for i := 1; i < len(arrays); i++ {
		if arrays[i].Ndim() != ndim {
			return nil, fmt.Errorf("numgo: all arrays must have same number of dimensions")
		}
		for d := 0; d < ndim; d++ {
			if arrays[i].shape[d] != arrays[0].shape[d] {
				return nil, fmt.Errorf("numgo: all arrays must have same shape")
			}
		}
	}
	// Expand each array, then concatenate.
	expanded := make([]*NDArray, len(arrays))
	for i, arr := range arrays {
		expanded[i] = ExpandDims(arr, axis)
	}
	return Concatenate(expanded, axis)
}

// Vstack stacks arrays vertically (along axis 0).
func Vstack(arrays []*NDArray) (*NDArray, error) {
	if len(arrays) == 0 {
		return nil, fmt.Errorf("numgo: need at least one array")
	}
	// For 1-D arrays, reshape to (1, N) then concatenate along axis 0.
	if arrays[0].Ndim() == 1 {
		reshaped := make([]*NDArray, len(arrays))
		for i, arr := range arrays {
			reshaped[i] = arr.Reshape(1, arr.Size())
		}
		return Concatenate(reshaped, 0)
	}
	return Concatenate(arrays, 0)
}

// Hstack stacks arrays horizontally.
// For 1-D arrays, concatenate. For 2D+, concatenate along axis 1.
func Hstack(arrays []*NDArray) (*NDArray, error) {
	if len(arrays) == 0 {
		return nil, fmt.Errorf("numgo: need at least one array")
	}
	if arrays[0].Ndim() == 1 {
		return Concatenate(arrays, 0)
	}
	return Concatenate(arrays, 1)
}

// Dstack stacks arrays along the third axis (axis 2).
func Dstack(arrays []*NDArray) (*NDArray, error) {
	if len(arrays) == 0 {
		return nil, fmt.Errorf("numgo: need at least one array")
	}
	ndim := arrays[0].Ndim()
	if ndim < 3 {
		// Expand to at least 3D.
		expanded := make([]*NDArray, len(arrays))
		for i, arr := range arrays {
			a := arr
			for a.Ndim() < 3 {
				a = ExpandDims(a, a.Ndim())
			}
			expanded[i] = a
		}
		return Concatenate(expanded, 2)
	}
	return Concatenate(arrays, 2)
}

// Split divides an array into equal sections along the given axis.
func Split(a *NDArray, sections int, axis int) ([]*NDArray, error) {
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		return nil, fmt.Errorf("numgo: axis %d out of range for %d dimensions", axis, ndim)
	}
	if sections <= 0 {
		return nil, fmt.Errorf("numgo: sections must be positive")
	}
	axisLen := a.shape[axis]
	if axisLen%sections != 0 {
		return nil, fmt.Errorf("numgo: array of size %d on axis %d cannot be split into %d equal sections", axisLen, axis, sections)
	}
	chunkSize := axisLen / sections
	results := make([]*NDArray, sections)
	for sec := 0; sec < sections; sec++ {
		// Build slice shape.
		sliceShape := make([]int, ndim)
		copy(sliceShape, a.shape)
		sliceShape[axis] = chunkSize
		result := NewNDArray(sliceShape, nil)
		srcIndices := make([]int, ndim)
		dstIndices := make([]int, ndim)
		for flat := 0; flat < result.Size(); flat++ {
			rem := flat
			for d := 0; d < ndim; d++ {
				dstIndices[d] = rem / result.strides[d]
				rem %= result.strides[d]
			}
			copy(srcIndices, dstIndices)
			srcIndices[axis] += sec * chunkSize
			result.data[flat] = a.Get(srcIndices...)
		}
		results[sec] = result
	}
	return results, nil
}

// Hsplit splits an array horizontally. For 1-D, splits along axis 0.
// For 2D+, splits along axis 1.
func Hsplit(a *NDArray, sections int) ([]*NDArray, error) {
	if a.Ndim() == 1 {
		return Split(a, sections, 0)
	}
	return Split(a, sections, 1)
}

// Vsplit splits an array vertically (along axis 0).
func Vsplit(a *NDArray, sections int) ([]*NDArray, error) {
	if a.Ndim() < 2 {
		return nil, fmt.Errorf("numgo: Vsplit requires at least 2 dimensions")
	}
	return Split(a, sections, 0)
}

// Dsplit splits an array along axis 2.
func Dsplit(a *NDArray, sections int) ([]*NDArray, error) {
	if a.Ndim() < 3 {
		return nil, fmt.Errorf("numgo: Dsplit requires at least 3 dimensions")
	}
	return Split(a, sections, 2)
}

// Tile constructs an array by repeating a the number of times given by reps.
func Tile(a *NDArray, reps []int) *NDArray {
	// Pad shape or reps to match in length.
	ndim := a.Ndim()
	nreps := len(reps)
	maxDim := ndim
	if nreps > maxDim {
		maxDim = nreps
	}

	// Pad shape with leading 1s.
	paddedShape := make([]int, maxDim)
	offset := maxDim - ndim
	for i := 0; i < offset; i++ {
		paddedShape[i] = 1
	}
	copy(paddedShape[offset:], a.shape)

	// Pad reps with leading 1s.
	paddedReps := make([]int, maxDim)
	rOffset := maxDim - nreps
	for i := 0; i < rOffset; i++ {
		paddedReps[i] = 1
	}
	copy(paddedReps[rOffset:], reps)

	// Reshape a to padded shape.
	src := NewNDArray(paddedShape, a.data)

	// Compute result shape.
	resultShape := make([]int, maxDim)
	for i := 0; i < maxDim; i++ {
		resultShape[i] = paddedShape[i] * paddedReps[i]
	}

	result := NewNDArray(resultShape, nil)
	dstIndices := make([]int, maxDim)
	srcIndices := make([]int, maxDim)
	for flat := 0; flat < result.Size(); flat++ {
		rem := flat
		for d := 0; d < maxDim; d++ {
			dstIndices[d] = rem / result.strides[d]
			rem %= result.strides[d]
		}
		for d := 0; d < maxDim; d++ {
			srcIndices[d] = dstIndices[d] % paddedShape[d]
		}
		result.data[flat] = src.Get(srcIndices...)
	}
	return result
}

// Repeat repeats elements of an array along the given axis.
func Repeat(a *NDArray, repeats int, axis int) *NDArray {
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		panic(fmt.Sprintf("numgo: axis %d out of range for %d dimensions", axis, ndim))
	}
	newShape := make([]int, ndim)
	copy(newShape, a.shape)
	newShape[axis] = a.shape[axis] * repeats

	result := NewNDArray(newShape, nil)
	srcIndices := make([]int, ndim)
	dstIndices := make([]int, ndim)
	for flat := 0; flat < a.Size(); flat++ {
		rem := flat
		for d := 0; d < ndim; d++ {
			srcIndices[d] = rem / a.strides[d]
			rem %= a.strides[d]
		}
		for r := 0; r < repeats; r++ {
			copy(dstIndices, srcIndices)
			dstIndices[axis] = srcIndices[axis]*repeats + r
			result.Set(a.data[flat], dstIndices...)
		}
	}
	return result
}

// Delete removes elements at the given indices along the specified axis.
func Delete(a *NDArray, indices []int, axis int) (*NDArray, error) {
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		return nil, fmt.Errorf("numgo: axis %d out of range for %d dimensions", axis, ndim)
	}
	// Build a set of indices to delete.
	delSet := make(map[int]bool)
	for _, idx := range indices {
		if idx < 0 || idx >= a.shape[axis] {
			return nil, fmt.Errorf("numgo: index %d out of range for axis %d with size %d", idx, axis, a.shape[axis])
		}
		delSet[idx] = true
	}
	newAxisLen := a.shape[axis] - len(delSet)
	if newAxisLen < 0 {
		newAxisLen = 0
	}
	newShape := make([]int, ndim)
	copy(newShape, a.shape)
	newShape[axis] = newAxisLen

	result := NewNDArray(newShape, nil)
	srcIndices := make([]int, ndim)
	dstIndices := make([]int, ndim)
	// Map from source axis index to dest axis index.
	axisMap := make([]int, a.shape[axis])
	di := 0
	for i := 0; i < a.shape[axis]; i++ {
		if delSet[i] {
			axisMap[i] = -1
		} else {
			axisMap[i] = di
			di++
		}
	}

	for flat := 0; flat < a.Size(); flat++ {
		rem := flat
		for d := 0; d < ndim; d++ {
			srcIndices[d] = rem / a.strides[d]
			rem %= a.strides[d]
		}
		if axisMap[srcIndices[axis]] < 0 {
			continue
		}
		copy(dstIndices, srcIndices)
		dstIndices[axis] = axisMap[srcIndices[axis]]
		result.Set(a.data[flat], dstIndices...)
	}
	return result, nil
}

// Insert inserts values before the given index along the specified axis.
func Insert(a *NDArray, index int, values *NDArray, axis int) (*NDArray, error) {
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		return nil, fmt.Errorf("numgo: axis %d out of range for %d dimensions", axis, ndim)
	}
	if index < 0 || index > a.shape[axis] {
		return nil, fmt.Errorf("numgo: index %d out of range for axis %d with size %d", index, axis, a.shape[axis])
	}
	// Values must have same shape except possibly along the insertion axis.
	if values.Ndim() != ndim {
		return nil, fmt.Errorf("numgo: values ndim %d does not match array ndim %d", values.Ndim(), ndim)
	}
	for d := 0; d < ndim; d++ {
		if d != axis && values.shape[d] != a.shape[d] {
			return nil, fmt.Errorf("numgo: shape mismatch on axis %d", d)
		}
	}

	newShape := make([]int, ndim)
	copy(newShape, a.shape)
	newShape[axis] += values.shape[axis]

	result := NewNDArray(newShape, nil)
	srcIndices := make([]int, ndim)
	dstIndices := make([]int, ndim)

	// Copy elements before index.
	for flat := 0; flat < a.Size(); flat++ {
		rem := flat
		for d := 0; d < ndim; d++ {
			srcIndices[d] = rem / a.strides[d]
			rem %= a.strides[d]
		}
		copy(dstIndices, srcIndices)
		if srcIndices[axis] >= index {
			dstIndices[axis] += values.shape[axis]
		}
		result.Set(a.data[flat], dstIndices...)
	}
	// Copy inserted values.
	valIndices := make([]int, ndim)
	for flat := 0; flat < values.Size(); flat++ {
		rem := flat
		for d := 0; d < ndim; d++ {
			valIndices[d] = rem / values.strides[d]
			rem %= values.strides[d]
		}
		copy(dstIndices, valIndices)
		dstIndices[axis] += index
		result.Set(values.data[flat], dstIndices...)
	}
	return result, nil
}

// Append appends values to the end of the array along the given axis.
func Append(a, values *NDArray, axis int) (*NDArray, error) {
	return Insert(a, a.shape[axis], values, axis)
}

// Flip reverses the order of elements along the given axis.
func Flip(a *NDArray, axis int) *NDArray {
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		panic(fmt.Sprintf("numgo: axis %d out of range for %d dimensions", axis, ndim))
	}
	result := NewNDArray(a.shape, nil)
	srcIndices := make([]int, ndim)
	dstIndices := make([]int, ndim)
	for flat := 0; flat < a.Size(); flat++ {
		rem := flat
		for d := 0; d < ndim; d++ {
			srcIndices[d] = rem / a.strides[d]
			rem %= a.strides[d]
		}
		copy(dstIndices, srcIndices)
		dstIndices[axis] = a.shape[axis] - 1 - srcIndices[axis]
		result.Set(a.data[flat], dstIndices...)
	}
	return result
}

// Fliplr flips the array left-right (reverses axis 1).
func Fliplr(a *NDArray) *NDArray {
	if a.Ndim() < 2 {
		panic("numgo: Fliplr requires at least 2 dimensions")
	}
	return Flip(a, 1)
}

// Flipud flips the array up-down (reverses axis 0).
func Flipud(a *NDArray) *NDArray {
	return Flip(a, 0)
}

// Rot90 rotates the array 90 degrees counter-clockwise k times in the plane
// defined by axes 0 and 1.
func Rot90(a *NDArray, k int) *NDArray {
	if a.Ndim() < 2 {
		panic("numgo: Rot90 requires at least 2 dimensions")
	}
	k = ((k % 4) + 4) % 4 // normalize to [0,3]
	if k == 0 {
		return a.Copy()
	}
	result := a
	for i := 0; i < k; i++ {
		// One 90-degree CCW rotation: transpose then flip axis 0.
		result = Flipud(Swapaxes(result, 0, 1))
	}
	return result
}

// Roll performs a circular shift of elements along the given axis.
// If axis is -1, the array is flattened, rolled, then reshaped.
func Roll(a *NDArray, shift int, axis int) *NDArray {
	if axis == -1 {
		flat := Ravel(a)
		rolled := Roll(flat, shift, 0)
		return NewNDArray(a.shape, rolled.data)
	}
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		panic(fmt.Sprintf("numgo: axis %d out of range for %d dimensions", axis, ndim))
	}
	axisLen := a.shape[axis]
	if axisLen == 0 {
		return a.Copy()
	}
	// Normalize shift.
	shift = ((shift % axisLen) + axisLen) % axisLen
	if shift == 0 {
		return a.Copy()
	}

	result := NewNDArray(a.shape, nil)
	srcIndices := make([]int, ndim)
	dstIndices := make([]int, ndim)
	for flat := 0; flat < a.Size(); flat++ {
		rem := flat
		for d := 0; d < ndim; d++ {
			srcIndices[d] = rem / a.strides[d]
			rem %= a.strides[d]
		}
		copy(dstIndices, srcIndices)
		dstIndices[axis] = (srcIndices[axis] + shift) % axisLen
		result.Set(a.data[flat], dstIndices...)
	}
	return result
}
