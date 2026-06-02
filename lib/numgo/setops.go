package numgo

import "sort"

// Intersect1D returns a sorted 1D NDArray of values common to both a and b.
// Both inputs are treated as flattened 1D arrays.
func Intersect1D(a, b *NDArray) *NDArray {
	setB := make(map[float64]bool, b.Size())
	for _, v := range b.data {
		setB[v] = true
	}

	seen := make(map[float64]bool)
	var result []float64
	for _, v := range a.data {
		if setB[v] && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	sort.Float64s(result)
	if len(result) == 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	return NewNDArray([]int{len(result)}, result)
}

// Union1D returns a sorted 1D NDArray of the unique values from both a and b.
// Both inputs are treated as flattened 1D arrays.
func Union1D(a, b *NDArray) *NDArray {
	seen := make(map[float64]bool)
	var result []float64

	for _, v := range a.data {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	for _, v := range b.data {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	sort.Float64s(result)
	if len(result) == 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	return NewNDArray([]int{len(result)}, result)
}

// In1d returns a 1-D NDArray with 1.0 where the corresponding element
// of a is found in b, and 0.0 otherwise. Both inputs are treated as flat.
func In1d(a, b *NDArray) *NDArray {
	setB := make(map[float64]bool, b.Size())
	for _, v := range b.data {
		setB[v] = true
	}
	data := make([]float64, a.Size())
	for i, v := range a.data {
		if setB[v] {
			data[i] = 1
		}
	}
	return FromSlice(data)
}

// Setxor1d returns a sorted 1D NDArray of values that are in exactly one
// of a or b (symmetric difference). Both inputs are treated as flat.
func Setxor1d(a, b *NDArray) *NDArray {
	countA := make(map[float64]bool)
	for _, v := range a.data {
		countA[v] = true
	}
	countB := make(map[float64]bool)
	for _, v := range b.data {
		countB[v] = true
	}

	seen := make(map[float64]bool)
	var result []float64
	for v := range countA {
		if !countB[v] && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	for v := range countB {
		if !countA[v] && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	sort.Float64s(result)
	if len(result) == 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	return NewNDArray([]int{len(result)}, result)
}

// SetDiff1D returns a sorted 1D NDArray of values in a that are not in b.
// Both inputs are treated as flattened 1D arrays.
func SetDiff1D(a, b *NDArray) *NDArray {
	setB := make(map[float64]bool, b.Size())
	for _, v := range b.data {
		setB[v] = true
	}

	seen := make(map[float64]bool)
	var result []float64
	for _, v := range a.data {
		if !setB[v] && !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}

	sort.Float64s(result)
	if len(result) == 0 {
		return NewNDArray([]int{0}, []float64{})
	}
	return NewNDArray([]int{len(result)}, result)
}
