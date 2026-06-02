//go:build unit

package numgo

import (
	"math"
	"testing"
)

func TestEmpty(t *testing.T) {
	a := Empty(2, 3)
	if a.Size() != 6 {
		t.Fatalf("expected size 6, got %d", a.Size())
	}
	for _, v := range a.Data() {
		if v != 0 {
			t.Fatal("expected all zeros")
		}
	}
}

func TestIdentity(t *testing.T) {
	a := Identity(3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if a.Get(i, j) != expected {
				t.Fatalf("unexpected value at (%d,%d): %f", i, j, a.Get(i, j))
			}
		}
	}
}

func TestArange(t *testing.T) {
	a := Arange(0, 5, 1)
	if a.Size() != 5 {
		t.Fatalf("expected 5 elements, got %d", a.Size())
	}
	for i := 0; i < 5; i++ {
		if a.data[i] != float64(i) {
			t.Fatalf("expected %d, got %f", i, a.data[i])
		}
	}
}

func TestArangeStep(t *testing.T) {
	a := Arange(0, 10, 2.5)
	// 0, 2.5, 5, 7.5
	if a.Size() != 4 {
		t.Fatalf("expected 4 elements, got %d", a.Size())
	}
}

func TestArangeNegativeStep(t *testing.T) {
	a := Arange(5, 0, -1)
	// 5, 4, 3, 2, 1
	if a.Size() != 5 {
		t.Fatalf("expected 5 elements, got %d", a.Size())
	}
	if a.data[0] != 5 || a.data[4] != 1 {
		t.Fatalf("unexpected values: %v", a.Data())
	}
}

func TestLinspace(t *testing.T) {
	a := Linspace(0, 1, 5)
	if a.Size() != 5 {
		t.Fatalf("expected 5 elements, got %d", a.Size())
	}
	expected := []float64{0, 0.25, 0.5, 0.75, 1.0}
	for i, e := range expected {
		if math.Abs(a.data[i]-e) > 1e-10 {
			t.Fatalf("at %d: expected %f, got %f", i, e, a.data[i])
		}
	}
}

func TestLinspaceSingle(t *testing.T) {
	a := Linspace(5, 10, 1)
	if a.Size() != 1 || a.data[0] != 5 {
		t.Fatalf("unexpected: %v", a.Data())
	}
}

func TestLogspace(t *testing.T) {
	a := Logspace(0, 2, 3)
	// 10^0=1, 10^1=10, 10^2=100
	expected := []float64{1, 10, 100}
	for i, e := range expected {
		if math.Abs(a.data[i]-e) > 1e-8 {
			t.Fatalf("at %d: expected %f, got %f", i, e, a.data[i])
		}
	}
}

func TestGeomspace(t *testing.T) {
	a := Geomspace(1, 1000, 4)
	// 1, 10, 100, 1000
	expected := []float64{1, 10, 100, 1000}
	for i, e := range expected {
		if math.Abs(a.data[i]-e)/e > 1e-8 {
			t.Fatalf("at %d: expected %f, got %f", i, e, a.data[i])
		}
	}
}

func TestMeshgrid(t *testing.T) {
	x := FromSlice([]float64{1, 2, 3})
	y := FromSlice([]float64{4, 5})
	grids := Meshgrid(x, y)
	if len(grids) != 2 {
		t.Fatalf("expected 2 grids, got %d", len(grids))
	}
	// Shape should be (3, 2) with "ij" indexing.
	s := grids[0].Shape()
	if s[0] != 3 || s[1] != 2 {
		t.Fatalf("unexpected shape: %v", s)
	}
	// First grid: x values broadcast along columns.
	if grids[0].Get(0, 0) != 1 || grids[0].Get(1, 0) != 2 || grids[0].Get(2, 0) != 3 {
		t.Fatalf("unexpected x grid: %v", grids[0].Data())
	}
	// Second grid: y values broadcast along rows.
	if grids[1].Get(0, 0) != 4 || grids[1].Get(0, 1) != 5 {
		t.Fatalf("unexpected y grid: %v", grids[1].Data())
	}
}

func TestDiag1D(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	d := Diag(a, 0)
	if d.Ndim() != 2 || d.Shape()[0] != 3 || d.Shape()[1] != 3 {
		t.Fatalf("unexpected shape: %v", d.Shape())
	}
	for i := 0; i < 3; i++ {
		if d.Get(i, i) != float64(i+1) {
			t.Fatalf("unexpected diagonal value at %d", i)
		}
	}
}

func TestDiag2D(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	d := Diag(a, 0)
	expected := []float64{1, 5, 9}
	for i, e := range expected {
		if d.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, d.data[i])
		}
	}
}

func TestDiag2DOffset(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	d := Diag(a, 1)
	// Superdiagonal: 2, 6
	if d.Size() != 2 || d.data[0] != 2 || d.data[1] != 6 {
		t.Fatalf("unexpected: %v", d.Data())
	}
}

func TestDiagflat(t *testing.T) {
	a := FromSlice([]float64{1, 2})
	d := Diagflat(a, 1)
	// 3x3 matrix with [1,2] on superdiagonal.
	if d.Shape()[0] != 3 || d.Shape()[1] != 3 {
		t.Fatalf("unexpected shape: %v", d.Shape())
	}
	if d.Get(0, 1) != 1 || d.Get(1, 2) != 2 {
		t.Fatalf("unexpected values: %v", d.Data())
	}
}

func TestTri(t *testing.T) {
	a := Tri(3, 3, 0)
	// Lower triangular with ones.
	expected := []float64{1, 0, 0, 1, 1, 0, 1, 1, 1}
	for i, e := range expected {
		if a.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, a.data[i])
		}
	}
}

func TestTril(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	l := Tril(a, 0)
	expected := []float64{1, 0, 0, 4, 5, 0, 7, 8, 9}
	for i, e := range expected {
		if l.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, l.data[i])
		}
	}
}

func TestTriu(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	u := Triu(a, 0)
	expected := []float64{1, 2, 3, 0, 5, 6, 0, 0, 9}
	for i, e := range expected {
		if u.data[i] != e {
			t.Fatalf("at %d: expected %f, got %f", i, e, u.data[i])
		}
	}
}

func TestVander(t *testing.T) {
	x := FromSlice([]float64{1, 2, 3})
	v := Vander(x, 3)
	// [[1,1,1],[4,2,1],[9,3,1]]
	expected := []float64{1, 1, 1, 4, 2, 1, 9, 3, 1}
	for i, e := range expected {
		if math.Abs(v.data[i]-e) > 1e-10 {
			t.Fatalf("at %d: expected %f, got %f", i, e, v.data[i])
		}
	}
}

func TestFromFunction(t *testing.T) {
	a := FromFunction([]int{3, 3}, func(idx []int) float64 {
		if idx[0] == idx[1] {
			return 1
		}
		return 0
	})
	if !matApproxEqual(a, Eye(3), 1e-10) {
		t.Fatal("FromFunction identity test failed")
	}
}

func TestFromIter(t *testing.T) {
	ch := make(chan float64, 5)
	for i := 0; i < 5; i++ {
		ch <- float64(i)
	}
	close(ch)
	a := FromIter(ch, 5)
	if a.Size() != 5 {
		t.Fatalf("expected 5 elements, got %d", a.Size())
	}
	for i := 0; i < 5; i++ {
		if a.data[i] != float64(i) {
			t.Fatalf("at %d: expected %f, got %f", i, float64(i), a.data[i])
		}
	}
}
