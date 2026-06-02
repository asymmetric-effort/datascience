package numgo

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
