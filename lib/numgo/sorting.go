package numgo

import (
	"fmt"
	"sort"
)

// Sort returns a new NDArray with elements sorted along the given axis.
// For a 1D array, axis must be 0. For an ND array, each 1D slice along
// the specified axis is sorted independently.
func Sort(a *NDArray, axis int) *NDArray {
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		panic(fmt.Sprintf("numgo.Sort: axis %d out of range for %d dimensions", axis, ndim))
	}

	result := a.Copy()
	axisLen := result.shape[axis]

	// Iterate over every 1D slice along the given axis.
	iterateAlongAxis(result, axis, func(indices []int) {
		slice := make([]float64, axisLen)
		for i := 0; i < axisLen; i++ {
			indices[axis] = i
			slice[i] = result.data[result.flatIndex(indices)]
		}
		sort.Float64s(slice)
		for i := 0; i < axisLen; i++ {
			indices[axis] = i
			result.data[result.flatIndex(indices)] = slice[i]
		}
	})

	return result
}

// ArgSort returns a new NDArray containing the indices that would sort the
// input array along the given axis. The result has the same shape as the input,
// with float64 index values.
func ArgSort(a *NDArray, axis int) *NDArray {
	ndim := a.Ndim()
	if axis < 0 || axis >= ndim {
		panic(fmt.Sprintf("numgo.ArgSort: axis %d out of range for %d dimensions", axis, ndim))
	}

	result := NewNDArray(a.shape, nil)
	axisLen := a.shape[axis]

	iterateAlongAxis(a, axis, func(indices []int) {
		// Build index-value pairs for this slice.
		type iv struct {
			idx int
			val float64
		}
		pairs := make([]iv, axisLen)
		for i := 0; i < axisLen; i++ {
			indices[axis] = i
			pairs[i] = iv{idx: i, val: a.data[a.flatIndex(indices)]}
		}
		sort.SliceStable(pairs, func(i, j int) bool {
			return pairs[i].val < pairs[j].val
		})
		for i := 0; i < axisLen; i++ {
			indices[axis] = i
			result.data[result.flatIndex(indices)] = float64(pairs[i].idx)
		}
	})

	return result
}

// Unique returns a 1D NDArray containing the sorted unique values from a.
func Unique(a *NDArray) *NDArray {
	if a.Size() == 0 {
		return NewNDArray([]int{0}, []float64{})
	}

	sorted := make([]float64, a.Size())
	copy(sorted, a.data)
	sort.Float64s(sorted)

	unique := []float64{sorted[0]}
	for i := 1; i < len(sorted); i++ {
		if sorted[i] != sorted[i-1] {
			unique = append(unique, sorted[i])
		}
	}
	return NewNDArray([]int{len(unique)}, unique)
}

// Where performs element-wise selection: for each element, if condition is
// nonzero (true), pick from x; otherwise pick from y.
// condition, x, and y must be broadcast-compatible.
func Where(condition, x, y *NDArray) *NDArray {
	// Broadcast all three to a common shape.
	shapeAB, err := BroadcastShapes(condition.shape, x.shape)
	if err != nil {
		panic(fmt.Sprintf("numgo.Where: %v", err))
	}
	resultShape, err := BroadcastShapes(shapeAB, y.shape)
	if err != nil {
		panic(fmt.Sprintf("numgo.Where: %v", err))
	}

	cb, err := BroadcastTo(condition, resultShape)
	if err != nil {
		panic(fmt.Sprintf("numgo.Where: %v", err))
	}
	xb, err := BroadcastTo(x, resultShape)
	if err != nil {
		panic(fmt.Sprintf("numgo.Where: %v", err))
	}
	yb, err := BroadcastTo(y, resultShape)
	if err != nil {
		panic(fmt.Sprintf("numgo.Where: %v", err))
	}

	data := make([]float64, cb.Size())
	for i := range data {
		if cb.data[i] != 0 {
			data[i] = xb.data[i]
		} else {
			data[i] = yb.data[i]
		}
	}
	return NewNDArray(resultShape, data)
}

// Nonzero returns the indices of all nonzero elements in a.
// The result is a slice of coordinate tuples, where each tuple has length
// equal to a.Ndim().
func Nonzero(a *NDArray) [][]int {
	ndim := a.Ndim()
	var result [][]int

	for flat := 0; flat < a.Size(); flat++ {
		if a.data[flat] == 0 {
			continue
		}
		coords := make([]int, ndim)
		rem := flat
		for d := 0; d < ndim; d++ {
			coords[d] = rem / a.strides[d]
			rem %= a.strides[d]
		}
		result = append(result, coords)
	}
	return result
}

// SearchSorted performs a binary search on a sorted 1D array, returning the
// indices at which the given values should be inserted to maintain sort order.
// The sorted array must be 1D. The result has the same shape as values.
func SearchSorted(sorted, values *NDArray) *NDArray {
	if sorted.Ndim() != 1 {
		panic("numgo.SearchSorted: sorted array must be 1D")
	}

	data := make([]float64, values.Size())
	for i := 0; i < values.Size(); i++ {
		v := values.data[i]
		// Binary search: find leftmost insertion point.
		idx := sort.Search(sorted.Size(), func(j int) bool {
			return sorted.data[j] >= v
		})
		data[i] = float64(idx)
	}
	return NewNDArray(values.Shape(), data)
}

// ArgMin returns the index of the minimum value along the given axis.
// Only a single axis is supported. Returns an NDArray of float64 indices.
func ArgMin(a *NDArray, axis int) *NDArray {
	return reduceAxis(a, []int{axis}, func(vals []float64) float64 {
		best := 0
		for i, v := range vals {
			if v < vals[best] {
				best = i
			}
		}
		return float64(best)
	})
}

// CountNonzero counts the number of nonzero elements along the given axes.
// If no axes are given, counts all nonzero elements.
func CountNonzero(a *NDArray, axes ...int) *NDArray {
	if len(axes) == 0 {
		count := 0.0
		for _, v := range a.data {
			if v != 0 {
				count++
			}
		}
		return FromSlice([]float64{count})
	}
	return reduceAxis(a, axes, func(vals []float64) float64 {
		count := 0.0
		for _, v := range vals {
			if v != 0 {
				count++
			}
		}
		return count
	})
}

// Extract returns a 1-D array of elements from a where the corresponding
// element in condition is nonzero. Both arrays are treated as flat.
func Extract(condition, a *NDArray) *NDArray {
	var result []float64
	n := condition.Size()
	if a.Size() < n {
		n = a.Size()
	}
	for i := 0; i < n; i++ {
		if condition.data[i] != 0 {
			result = append(result, a.data[i])
		}
	}
	if len(result) == 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	return FromSlice(result)
}

// Flatnonzero returns the flat indices of nonzero elements.
func Flatnonzero(a *NDArray) *NDArray {
	var indices []float64
	for i, v := range a.data {
		if v != 0 {
			indices = append(indices, float64(i))
		}
	}
	if len(indices) == 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	return FromSlice(indices)
}

// Argwhere returns the coordinates of nonzero elements.
// This is the same as Nonzero.
func Argwhere(a *NDArray) [][]int {
	return Nonzero(a)
}

// Lexsort performs an indirect stable sort using a sequence of keys.
// The last key is the primary sort key, the second-to-last is secondary, etc.
// All keys must be 1-D arrays of the same length.
// Returns an NDArray of indices that sorts the data.
func Lexsort(keys []*NDArray) *NDArray {
	if len(keys) == 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	n := keys[0].Size()
	for _, k := range keys {
		if k.Ndim() != 1 {
			panic("numgo.Lexsort: all keys must be 1-D")
		}
		if k.Size() != n {
			panic("numgo.Lexsort: all keys must have the same length")
		}
	}

	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}

	sort.SliceStable(indices, func(i, j int) bool {
		// Last key is primary.
		for k := len(keys) - 1; k >= 0; k-- {
			vi := keys[k].data[indices[i]]
			vj := keys[k].data[indices[j]]
			if vi < vj {
				return true
			}
			if vi > vj {
				return false
			}
		}
		return false
	})

	data := make([]float64, n)
	for i, v := range indices {
		data[i] = float64(v)
	}
	return FromSlice(data)
}

// Partition rearranges elements along the given axis such that the element at
// position kth is in its sorted position, all smaller elements are before it,
// and all larger elements are after it (using introselect/quickselect).
func Partition(a *NDArray, kth int, axis int) *NDArray {
	if axis < 0 || axis >= a.Ndim() {
		panic(fmt.Sprintf("numgo.Partition: axis %d out of range for %d dimensions", axis, a.Ndim()))
	}
	axisLen := a.shape[axis]
	if kth < 0 || kth >= axisLen {
		panic(fmt.Sprintf("numgo.Partition: kth %d out of range [0, %d)", kth, axisLen))
	}

	result := a.Copy()
	iterateAlongAxis(result, axis, func(indices []int) {
		slice := make([]float64, axisLen)
		for i := 0; i < axisLen; i++ {
			indices[axis] = i
			slice[i] = result.data[result.flatIndex(indices)]
		}
		quickselect(slice, kth)
		for i := 0; i < axisLen; i++ {
			indices[axis] = i
			result.data[result.flatIndex(indices)] = slice[i]
		}
	})
	return result
}

// Argpartition returns the indices that would partition the array.
func Argpartition(a *NDArray, kth int, axis int) *NDArray {
	if axis < 0 || axis >= a.Ndim() {
		panic(fmt.Sprintf("numgo.Argpartition: axis %d out of range for %d dimensions", axis, a.Ndim()))
	}
	axisLen := a.shape[axis]
	if kth < 0 || kth >= axisLen {
		panic(fmt.Sprintf("numgo.Argpartition: kth %d out of range [0, %d)", kth, axisLen))
	}

	result := NewNDArray(a.shape, nil)
	iterateAlongAxis(a, axis, func(indices []int) {
		pairs := make([]ivPair, axisLen)
		for i := 0; i < axisLen; i++ {
			indices[axis] = i
			pairs[i] = ivPair{idx: i, val: a.data[a.flatIndex(indices)]}
		}
		quickselectPairs(pairs, kth)
		for i := 0; i < axisLen; i++ {
			indices[axis] = i
			result.data[result.flatIndex(indices)] = float64(pairs[i].idx)
		}
	})
	return result
}

// quickselect partially sorts slice so that slice[kth] is the kth smallest element,
// elements before kth are <= slice[kth], and elements after are >= slice[kth].
func quickselect(slice []float64, kth int) {
	lo, hi := 0, len(slice)-1
	for lo < hi {
		pivot := slice[hi]
		i := lo
		for j := lo; j < hi; j++ {
			if slice[j] < pivot {
				slice[i], slice[j] = slice[j], slice[i]
				i++
			}
		}
		slice[i], slice[hi] = slice[hi], slice[i]
		if i == kth {
			return
		} else if i < kth {
			lo = i + 1
		} else {
			hi = i - 1
		}
	}
}

type ivPair struct {
	idx int
	val float64
}

func quickselectPairs(pairs []ivPair, kth int) {
	lo, hi := 0, len(pairs)-1
	for lo < hi {
		pivot := pairs[hi].val
		i := lo
		for j := lo; j < hi; j++ {
			if pairs[j].val < pivot {
				pairs[i], pairs[j] = pairs[j], pairs[i]
				i++
			}
		}
		pairs[i], pairs[hi] = pairs[hi], pairs[i]
		if i == kth {
			return
		} else if i < kth {
			lo = i + 1
		} else {
			hi = i - 1
		}
	}
}

// iterateAlongAxis calls fn once for each 1D slice along the given axis.
// fn receives a mutable indices slice; fn is responsible for setting indices[axis]
// to iterate within the slice. On entry, indices[axis] is 0.
func iterateAlongAxis(a *NDArray, axis int, fn func(indices []int)) {
	ndim := a.Ndim()
	if ndim == 0 {
		return
	}

	// Total number of slices = product of all dims except axis.
	totalSlices := 1
	for d := 0; d < ndim; d++ {
		if d != axis {
			totalSlices *= a.shape[d]
		}
	}

	// Build a shape for the "other" dimensions to iterate over.
	otherDims := make([]int, 0, ndim-1)
	for d := 0; d < ndim; d++ {
		if d != axis {
			otherDims = append(otherDims, d)
		}
	}

	indices := make([]int, ndim)
	for s := 0; s < totalSlices; s++ {
		// Decompose s into coordinates for otherDims.
		rem := s
		for i := len(otherDims) - 1; i >= 0; i-- {
			d := otherDims[i]
			indices[d] = rem % a.shape[d]
			rem /= a.shape[d]
		}
		indices[axis] = 0
		fn(indices)
	}
}
