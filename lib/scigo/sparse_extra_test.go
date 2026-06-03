//go:build unit

package scigo

import (
	"math"
	"testing"
)

func approxEqualSP(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

// ---------------------------------------------------------------------------
// EyeSparse
// ---------------------------------------------------------------------------

func TestEyeSparse(t *testing.T) {
	eye := EyeSparse(3)
	if eye.Shape() != [2]int{3, 3} {
		t.Errorf("EyeSparse shape=%v, want [3,3]", eye.Shape())
	}
	if eye.NNZ() != 3 {
		t.Errorf("EyeSparse NNZ=%v, want 3", eye.NNZ())
	}
	dense := eye.ToDense()
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if dense[i][j] != expected {
				t.Errorf("EyeSparse[%d][%d]=%v, want %v", i, j, dense[i][j], expected)
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Diags
// ---------------------------------------------------------------------------

func TestDiags_MainDiagonal(t *testing.T) {
	d, err := Diags([][]float64{{1, 2, 3}}, []int{0}, 3)
	if err != nil {
		t.Fatal(err)
	}
	dense := d.ToDense()
	for i := 0; i < 3; i++ {
		if dense[i][i] != float64(i+1) {
			t.Errorf("Diags main[%d][%d]=%v, want %v", i, i, dense[i][i], float64(i+1))
		}
	}
}

func TestDiags_MultiDiagonals(t *testing.T) {
	// Tridiagonal matrix
	main := []float64{2, 2, 2, 2}
	upper := []float64{-1, -1, -1}
	lower := []float64{-1, -1, -1}
	d, err := Diags([][]float64{lower, main, upper}, []int{-1, 0, 1}, 4)
	if err != nil {
		t.Fatal(err)
	}
	dense := d.ToDense()
	// Check main diagonal
	for i := 0; i < 4; i++ {
		if dense[i][i] != 2 {
			t.Errorf("Diags tridiag[%d][%d]=%v, want 2", i, i, dense[i][i])
		}
	}
	// Check super-diagonal
	for i := 0; i < 3; i++ {
		if dense[i][i+1] != -1 {
			t.Errorf("Diags tridiag[%d][%d]=%v, want -1", i, i+1, dense[i][i+1])
		}
	}
	// Check sub-diagonal
	for i := 1; i < 4; i++ {
		if dense[i][i-1] != -1 {
			t.Errorf("Diags tridiag[%d][%d]=%v, want -1", i, i-1, dense[i][i-1])
		}
	}
}

// ---------------------------------------------------------------------------
// KronSparse
// ---------------------------------------------------------------------------

func TestKronSparse_Identity(t *testing.T) {
	i2 := EyeSparse(2)
	i3 := EyeSparse(3)
	kron := KronSparse(i2, i3)
	shape := kron.Shape()
	if shape != [2]int{6, 6} {
		t.Errorf("KronSparse(I2,I3) shape=%v, want [6,6]", shape)
	}
	if kron.NNZ() != 6 {
		t.Errorf("KronSparse(I2,I3) NNZ=%v, want 6", kron.NNZ())
	}
	// Should be I6
	dense := kron.ToDense()
	for i := 0; i < 6; i++ {
		for j := 0; j < 6; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if dense[i][j] != expected {
				t.Errorf("KronSparse[%d][%d]=%v, want %v", i, j, dense[i][j], expected)
			}
		}
	}
}

func TestKronSparse_Values(t *testing.T) {
	// 2x2 matrix [[1,2],[3,4]] kron identity 2x2
	a, _ := NewCSR([]int{0, 2, 4}, []int{0, 1, 0, 1}, []float64{1, 2, 3, 4}, [2]int{2, 2})
	i2 := EyeSparse(2)
	kron := KronSparse(a, i2)
	shape := kron.Shape()
	if shape != [2]int{4, 4} {
		t.Errorf("KronSparse shape=%v, want [4,4]", shape)
	}
	dense := kron.ToDense()
	// A kron I2 = [[1,0,2,0],[0,1,0,2],[3,0,4,0],[0,3,0,4]]
	expected := [][]float64{
		{1, 0, 2, 0},
		{0, 1, 0, 2},
		{3, 0, 4, 0},
		{0, 3, 0, 4},
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if !approxEqualSP(dense[i][j], expected[i][j], 1e-10) {
				t.Errorf("KronSparse[%d][%d]=%v, want %v", i, j, dense[i][j], expected[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// HStackSparse and VStackSparse
// ---------------------------------------------------------------------------

func TestHStackSparse(t *testing.T) {
	a := EyeSparse(2)
	b := EyeSparse(2)
	h := HStackSparse(a, b)
	shape := h.Shape()
	if shape != [2]int{2, 4} {
		t.Errorf("HStackSparse shape=%v, want [2,4]", shape)
	}
	dense := h.ToDense()
	expected := [][]float64{
		{1, 0, 1, 0},
		{0, 1, 0, 1},
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 4; j++ {
			if dense[i][j] != expected[i][j] {
				t.Errorf("HStackSparse[%d][%d]=%v, want %v", i, j, dense[i][j], expected[i][j])
			}
		}
	}
}

func TestVStackSparse(t *testing.T) {
	a := EyeSparse(2)
	b := EyeSparse(2)
	v := VStackSparse(a, b)
	shape := v.Shape()
	if shape != [2]int{4, 2} {
		t.Errorf("VStackSparse shape=%v, want [4,2]", shape)
	}
	dense := v.ToDense()
	expected := [][]float64{
		{1, 0},
		{0, 1},
		{1, 0},
		{0, 1},
	}
	for i := 0; i < 4; i++ {
		for j := 0; j < 2; j++ {
			if dense[i][j] != expected[i][j] {
				t.Errorf("VStackSparse[%d][%d]=%v, want %v", i, j, dense[i][j], expected[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// FindSparse
// ---------------------------------------------------------------------------

func TestFindSparse(t *testing.T) {
	m, _ := NewCSR(
		[]int{0, 2, 3, 4},
		[]int{0, 2, 1, 0},
		[]float64{1, 3, 5, 7},
		[2]int{3, 3},
	)
	entries := FindSparse(m)
	if len(entries) != 4 {
		t.Errorf("FindSparse len=%v, want 4", len(entries))
	}
	// Check that we found the expected entries
	found := make(map[[2]int]float64)
	for _, e := range entries {
		found[[2]int{e.Row, e.Col}] = e.Val
	}
	if found[[2]int{0, 0}] != 1 {
		t.Errorf("FindSparse missing (0,0)=1")
	}
	if found[[2]int{0, 2}] != 3 {
		t.Errorf("FindSparse missing (0,2)=3")
	}
	if found[[2]int{1, 1}] != 5 {
		t.Errorf("FindSparse missing (1,1)=5")
	}
	if found[[2]int{2, 0}] != 7 {
		t.Errorf("FindSparse missing (2,0)=7")
	}
}

func TestFindSparse_Empty(t *testing.T) {
	m := EyeSparse(3)
	entries := FindSparse(m)
	if len(entries) != 3 {
		t.Errorf("FindSparse(eye3) len=%v, want 3", len(entries))
	}
}
