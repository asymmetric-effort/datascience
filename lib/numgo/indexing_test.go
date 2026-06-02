//go:build unit

package numgo

import (
	"math"
	"testing"
)

func TestTakeFlat(t *testing.T) {
	a := FromSlice([]float64{10, 20, 30, 40, 50})
	r, err := Take(a, []int{0, 2, 4}, -1)
	if err != nil {
		t.Fatal(err)
	}
	expected := []float64{10, 30, 50}
	for i, e := range expected {
		if r.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, r.data[i])
		}
	}
}

func TestTakeAxis0(t *testing.T) {
	a := NewNDArray([]int{3, 2}, []float64{1, 2, 3, 4, 5, 6})
	r, err := Take(a, []int{0, 2}, 0)
	if err != nil {
		t.Fatal(err)
	}
	// Rows 0 and 2: [[1,2],[5,6]]
	if r.Shape()[0] != 2 || r.Shape()[1] != 2 {
		t.Fatalf("unexpected shape: %v", r.Shape())
	}
	expected := []float64{1, 2, 5, 6}
	for i, e := range expected {
		if r.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, r.data[i])
		}
	}
}

func TestTakeAxis1(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	r, err := Take(a, []int{0, 2}, 1)
	if err != nil {
		t.Fatal(err)
	}
	// Columns 0 and 2: [[1,3],[4,6]]
	expected := []float64{1, 3, 4, 6}
	for i, e := range expected {
		if r.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, r.data[i])
		}
	}
}

func TestTakeOutOfRange(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	_, err := Take(a, []int{5}, -1)
	if err == nil {
		t.Fatal("expected error for out-of-range index")
	}
}

func TestTakeAlongAxis(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{10, 20, 30, 40, 50, 60})
	indices := NewNDArray([]int{2, 1}, []float64{1, 2})
	r, err := TakeAlongAxis(a, indices, 1)
	if err != nil {
		t.Fatal(err)
	}
	// row0: col1=20, row1: col2=60
	if r.data[0] != 20 || r.data[1] != 60 {
		t.Fatalf("unexpected: %v", r.Data())
	}
}

func TestChoose(t *testing.T) {
	indices := FromSlice([]float64{0, 1, 2, 0})
	choices := []*NDArray{
		FromSlice([]float64{10, 10, 10, 10}),
		FromSlice([]float64{20, 20, 20, 20}),
		FromSlice([]float64{30, 30, 30, 30}),
	}
	r, err := Choose(indices, choices)
	if err != nil {
		t.Fatal(err)
	}
	expected := []float64{10, 20, 30, 10}
	for i, e := range expected {
		if r.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, r.data[i])
		}
	}
}

func TestCompressFlat(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4, 5})
	r, err := Compress([]bool{true, false, true, false, true}, a, -1)
	if err != nil {
		t.Fatal(err)
	}
	expected := []float64{1, 3, 5}
	if r.Size() != 3 {
		t.Fatalf("expected 3 elements, got %d", r.Size())
	}
	for i, e := range expected {
		if r.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, r.data[i])
		}
	}
}

func TestCompressAxis0(t *testing.T) {
	a := NewNDArray([]int{3, 2}, []float64{1, 2, 3, 4, 5, 6})
	r, err := Compress([]bool{true, false, true}, a, 0)
	if err != nil {
		t.Fatal(err)
	}
	// Rows 0 and 2: [[1,2],[5,6]]
	expected := []float64{1, 2, 5, 6}
	for i, e := range expected {
		if r.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, r.data[i])
		}
	}
}

func TestDiagonal(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	d, err := Diagonal(a, 0, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	expected := []float64{1, 5, 9}
	for i, e := range expected {
		if d.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, d.data[i])
		}
	}
}

func TestDiagonalOffset(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	d, err := Diagonal(a, 1, 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	expected := []float64{2, 6}
	if d.Size() != 2 {
		t.Fatalf("expected 2 elements, got %d", d.Size())
	}
	for i, e := range expected {
		if d.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, d.data[i])
		}
	}
}

func TestSelect(t *testing.T) {
	cond1 := FromSlice([]float64{1, 0, 0, 0})
	cond2 := FromSlice([]float64{0, 1, 0, 0})
	choice1 := FromSlice([]float64{10, 10, 10, 10})
	choice2 := FromSlice([]float64{20, 20, 20, 20})

	r, err := Select([]*NDArray{cond1, cond2}, []*NDArray{choice1, choice2}, -1)
	if err != nil {
		t.Fatal(err)
	}
	expected := []float64{10, 20, -1, -1}
	for i, e := range expected {
		if r.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, r.data[i])
		}
	}
}

func TestAsStrided(t *testing.T) {
	a := FromSlice([]float64{0, 1, 2, 3, 4, 5})
	// Create a 2x3 view with default row-major strides.
	r := AsStrided(a, []int{2, 3}, []int{3, 1})
	if r.Shape()[0] != 2 || r.Shape()[1] != 3 {
		t.Fatalf("unexpected shape: %v", r.Shape())
	}
	if r.Get(0, 0) != 0 || r.Get(0, 2) != 2 || r.Get(1, 0) != 3 || r.Get(1, 2) != 5 {
		t.Fatalf("unexpected values: %v", r.Data())
	}
}

func TestAsStridedRepeat(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	// Repeat row: strides[0]=0 means every row sees the same data.
	r := AsStrided(a, []int{3, 3}, []int{0, 1})
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if math.Abs(r.Get(i, j)-float64(j+1)) > 1e-10 {
				t.Fatalf("at (%d,%d): expected %f, got %f", i, j, float64(j+1), r.Get(i, j))
			}
		}
	}
}
