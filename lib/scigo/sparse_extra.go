package scigo

import "fmt"

// EyeSparse creates a sparse identity matrix of size n x n in CSR format.
func EyeSparse(n int) *CSR {
	if n <= 0 {
		return &CSR{indptr: []int{0}, indices: nil, data: nil, shape: [2]int{1, 1}}
	}
	indptr := make([]int, n+1)
	indices := make([]int, n)
	data := make([]float64, n)
	for i := 0; i < n; i++ {
		indptr[i+1] = i + 1
		indices[i] = i
		data[i] = 1.0
	}
	return &CSR{indptr: indptr, indices: indices, data: data, shape: [2]int{n, n}}
}

// Diags creates a sparse diagonal matrix from the given diagonals and offsets.
// diagonals[k] is placed at offset offsets[k].
// offset 0 is the main diagonal, positive offsets are above, negative below.
// n is the size of the resulting square matrix.
// Returns a CSR sparse matrix.
func Diags(diagonals [][]float64, offsets []int, n int) (*CSR, error) {
	if len(diagonals) != len(offsets) {
		return nil, fmt.Errorf("scigo: Diags: diagonals and offsets must have the same length")
	}
	if n <= 0 {
		return nil, fmt.Errorf("scigo: Diags: n must be positive")
	}

	// Accumulate entries into a map
	type key struct{ r, c int }
	entries := make(map[key]float64)
	for k, offset := range offsets {
		diag := diagonals[k]
		for i, v := range diag {
			var r, c int
			if offset >= 0 {
				r = i
				c = i + offset
			} else {
				r = i - offset
				c = i
			}
			if r >= 0 && r < n && c >= 0 && c < n {
				entries[key{r, c}] += v
			}
		}
	}

	// Build CSR
	indptr := make([]int, n+1)
	var indices []int
	var data []float64

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if v, ok := entries[key{i, j}]; ok && v != 0 {
				indices = append(indices, j)
				data = append(data, v)
			}
		}
		indptr[i+1] = len(data)
	}

	return &CSR{indptr: indptr, indices: indices, data: data, shape: [2]int{n, n}}, nil
}

// KronSparse computes the Kronecker product of two CSR sparse matrices.
// The result is a CSR sparse matrix of shape (a.rows*b.rows, a.cols*b.cols).
func KronSparse(a, b *CSR) *CSR {
	aShape := a.Shape()
	bShape := b.Shape()
	rRows := aShape[0] * bShape[0]
	rCols := aShape[1] * bShape[1]

	indptr := make([]int, rRows+1)
	var indices []int
	var data []float64

	for ia := 0; ia < aShape[0]; ia++ {
		aStart, aEnd := a.indptr[ia], a.indptr[ia+1]
		for ib := 0; ib < bShape[0]; ib++ {
			bStart, bEnd := b.indptr[ib], b.indptr[ib+1]
			row := ia*bShape[0] + ib
			for aj := aStart; aj < aEnd; aj++ {
				ac := a.indices[aj]
				av := a.data[aj]
				for bj := bStart; bj < bEnd; bj++ {
					bc := b.indices[bj]
					bv := b.data[bj]
					col := ac*bShape[1] + bc
					indices = append(indices, col)
					data = append(data, av*bv)
				}
			}
			indptr[row+1] = len(data)
		}
	}

	return &CSR{indptr: indptr, indices: indices, data: data, shape: [2]int{rRows, rCols}}
}

// HStackSparse horizontally stacks two CSR sparse matrices (side by side).
// Both matrices must have the same number of rows.
// Panics if the row counts do not match.
func HStackSparse(a, b *CSR) *CSR {
	aShape := a.Shape()
	bShape := b.Shape()
	if aShape[0] != bShape[0] {
		panic(fmt.Sprintf("scigo: HStackSparse: row count mismatch %d vs %d", aShape[0], bShape[0]))
	}

	nrows := aShape[0]
	ncols := aShape[1] + bShape[1]
	indptr := make([]int, nrows+1)
	var indices []int
	var data []float64

	for i := 0; i < nrows; i++ {
		// Add entries from a
		as, ae := a.indptr[i], a.indptr[i+1]
		for j := as; j < ae; j++ {
			indices = append(indices, a.indices[j])
			data = append(data, a.data[j])
		}
		// Add entries from b (shifted by a's column count)
		bs, be := b.indptr[i], b.indptr[i+1]
		for j := bs; j < be; j++ {
			indices = append(indices, b.indices[j]+aShape[1])
			data = append(data, b.data[j])
		}
		indptr[i+1] = len(data)
	}

	return &CSR{indptr: indptr, indices: indices, data: data, shape: [2]int{nrows, ncols}}
}

// VStackSparse vertically stacks two CSR sparse matrices (one on top of the other).
// Both matrices must have the same number of columns.
// Panics if the column counts do not match.
func VStackSparse(a, b *CSR) *CSR {
	aShape := a.Shape()
	bShape := b.Shape()
	if aShape[1] != bShape[1] {
		panic(fmt.Sprintf("scigo: VStackSparse: column count mismatch %d vs %d", aShape[1], bShape[1]))
	}

	nrows := aShape[0] + bShape[0]
	ncols := aShape[1]
	indptr := make([]int, nrows+1)

	// Copy a's data
	indices := make([]int, len(a.indices)+len(b.indices))
	data := make([]float64, len(a.data)+len(b.data))
	copy(indices, a.indices)
	copy(data, a.data)
	copy(indices[len(a.indices):], b.indices)
	copy(data[len(a.data):], b.data)

	// Build indptr
	for i := 0; i <= aShape[0]; i++ {
		indptr[i] = a.indptr[i]
	}
	offset := a.indptr[aShape[0]]
	for i := 0; i <= bShape[0]; i++ {
		indptr[aShape[0]+i] = b.indptr[i] + offset
	}

	return &CSR{indptr: indptr, indices: indices, data: data, shape: [2]int{nrows, ncols}}
}

// SparseEntry represents a nonzero entry in a sparse matrix.
type SparseEntry struct {
	Row int
	Col int
	Val float64
}

// FindSparse returns the nonzero entries of a CSR sparse matrix as a list of
// (row, col, value) triples, analogous to scipy.sparse.find.
func FindSparse(m *CSR) []SparseEntry {
	var entries []SparseEntry
	shape := m.Shape()
	for i := 0; i < shape[0]; i++ {
		start, end := m.indptr[i], m.indptr[i+1]
		for j := start; j < end; j++ {
			if m.data[j] != 0 {
				entries = append(entries, SparseEntry{Row: i, Col: m.indices[j], Val: m.data[j]})
			}
		}
	}
	return entries
}
