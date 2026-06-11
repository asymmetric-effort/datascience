package numgo

import (
	"fmt"
	"math"
	"strings"
)

// NDArray is a multidimensional array backed by a flat []float64 slice.
type NDArray struct {
	data    []float64
	shape   []int
	strides []int
}

// computeStrides returns row-major (C-order) strides for the given shape.
func computeStrides(shape []int) []int {
	n := len(shape)
	strides := make([]int, n)
	if n == 0 {
		return strides
	}
	strides[n-1] = 1
	for i := n - 2; i >= 0; i-- {
		strides[i] = strides[i+1] * shape[i+1]
	}
	return strides
}

// product returns the product of ints with overflow and negative dimension checks.
// It panics if any dimension is negative or if the product overflows int.
// Zero dimensions are allowed (they produce a zero-size array).
func product(vals []int) int {
	p := 1
	for _, v := range vals {
		if v < 0 {
			panic(fmt.Sprintf("numgo: shape dimensions must be non-negative, got %d", v))
		}
		if v == 0 {
			return 0
		}
		if p > math.MaxInt/v {
			panic(fmt.Sprintf("numgo: shape product overflows int: %v", vals))
		}
		p *= v
	}
	return p
}

// productUnsafe returns the product of ints without overflow or validation checks.
// Use only when dimensions have already been validated (e.g., internal reshapes).
func productUnsafe(vals []int) int {
	p := 1
	for _, v := range vals {
		p *= v
	}
	return p
}

// NewNDArray creates an NDArray with the given shape and optional data.
// If data is nil, the array is zero-initialized.
// If data is provided, its length must equal the product of shape dimensions.
func NewNDArray(shape []int, data []float64) *NDArray {
	size := product(shape)
	if data == nil {
		data = make([]float64, size)
	}
	if len(data) != size {
		panic(fmt.Sprintf("numgo: data length %d does not match shape %v (size %d)", len(data), shape, size))
	}
	s := make([]int, len(shape))
	copy(s, shape)
	d := make([]float64, len(data))
	copy(d, data)
	return &NDArray{
		data:    d,
		shape:   s,
		strides: computeStrides(s),
	}
}

// Shape returns a copy of the array's shape.
func (a *NDArray) Shape() []int {
	s := make([]int, len(a.shape))
	copy(s, a.shape)
	return s
}

// Ndim returns the number of dimensions.
func (a *NDArray) Ndim() int {
	return len(a.shape)
}

// Size returns the total number of elements.
func (a *NDArray) Size() int {
	return len(a.data)
}

// Data returns a copy of the underlying data slice.
func (a *NDArray) Data() []float64 {
	d := make([]float64, len(a.data))
	copy(d, a.data)
	return d
}

// flatIndex converts multidimensional indices to a flat index.
func (a *NDArray) flatIndex(indices []int) int {
	if len(indices) != len(a.shape) {
		panic(fmt.Sprintf("numgo: expected %d indices, got %d", len(a.shape), len(indices)))
	}
	idx := 0
	for i, v := range indices {
		if v < 0 || v >= a.shape[i] {
			panic(fmt.Sprintf("numgo: index %d out of range [0, %d) for axis %d", v, a.shape[i], i))
		}
		idx += v * a.strides[i]
	}
	return idx
}

// Get returns the element at the given indices.
func (a *NDArray) Get(indices ...int) float64 {
	return a.data[a.flatIndex(indices)]
}

// Set sets the element at the given indices.
func (a *NDArray) Set(value float64, indices ...int) {
	a.data[a.flatIndex(indices)] = value
}

// Reshape returns a new NDArray with the same data but a different shape.
// The total number of elements must remain the same.
func (a *NDArray) Reshape(shape ...int) *NDArray {
	newSize := product(shape)
	if newSize != a.Size() {
		panic(fmt.Sprintf("numgo: cannot reshape size %d into shape %v (size %d)", a.Size(), shape, newSize))
	}
	return NewNDArray(shape, a.data)
}

// Flatten returns a 1-D copy of the array.
func (a *NDArray) Flatten() *NDArray {
	return NewNDArray([]int{a.Size()}, a.data)
}

// Copy returns a deep copy of the array.
func (a *NDArray) Copy() *NDArray {
	return NewNDArray(a.shape, a.data)
}

// T returns the transpose of the array.
// For 1-D arrays it returns a copy. For N-D arrays it reverses the axes.
func (a *NDArray) T() *NDArray {
	ndim := a.Ndim()
	if ndim <= 1 {
		return a.Copy()
	}

	newShape := make([]int, ndim)
	for i := range newShape {
		newShape[i] = a.shape[ndim-1-i]
	}

	result := NewNDArray(newShape, nil)
	indices := make([]int, ndim)
	transposed := make([]int, ndim)

	for flat := 0; flat < a.Size(); flat++ {
		// Decompose flat index into original indices.
		rem := flat
		for d := 0; d < ndim; d++ {
			indices[d] = rem / a.strides[d]
			rem %= a.strides[d]
		}
		// Reverse indices for transposed position.
		for d := 0; d < ndim; d++ {
			transposed[d] = indices[ndim-1-d]
		}
		result.Set(a.data[flat], transposed...)
	}
	return result
}

// String returns a human-readable representation of the array.
func (a *NDArray) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("NDArray(shape=%v, data=", a.shape))
	if a.Ndim() <= 2 && a.Size() <= 100 {
		a.writeFormatted(&b, 0, 0)
	} else {
		b.WriteString(fmt.Sprintf("%v", a.data))
	}
	b.WriteString(")")
	return b.String()
}

// writeFormatted writes a nested bracket representation.
func (a *NDArray) writeFormatted(b *strings.Builder, axis int, offset int) {
	if axis == a.Ndim()-1 || a.Ndim() == 0 {
		b.WriteString("[")
		for i := 0; i < a.shape[axis]; i++ {
			if i > 0 {
				b.WriteString(" ")
			}
			b.WriteString(fmt.Sprintf("%g", a.data[offset+i]))
		}
		b.WriteString("]")
		return
	}
	b.WriteString("[")
	for i := 0; i < a.shape[axis]; i++ {
		if i > 0 {
			b.WriteString(" ")
		}
		a.writeFormatted(b, axis+1, offset+i*a.strides[axis])
	}
	b.WriteString("]")
}
