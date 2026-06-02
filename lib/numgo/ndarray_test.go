//go:build unit

package numgo

import (
	"strings"
	"testing"
)

func TestNewNDArray(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	if a.Ndim() != 2 {
		t.Fatalf("expected Ndim=2, got %d", a.Ndim())
	}
	if a.Size() != 6 {
		t.Fatalf("expected Size=6, got %d", a.Size())
	}
	shape := a.Shape()
	if shape[0] != 2 || shape[1] != 3 {
		t.Fatalf("unexpected shape %v", shape)
	}
}

func TestNewNDArrayZeroInit(t *testing.T) {
	a := NewNDArray([]int{3, 2}, nil)
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			if a.Get(i, j) != 0 {
				t.Fatalf("expected 0 at (%d,%d), got %f", i, j, a.Get(i, j))
			}
		}
	}
}

func TestNewNDArrayPanicsOnSizeMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on size mismatch")
		}
	}()
	NewNDArray([]int{2, 2}, []float64{1, 2, 3})
}

func TestGetSet(t *testing.T) {
	a := NewNDArray([]int{2, 3}, nil)
	a.Set(42, 1, 2)
	if v := a.Get(1, 2); v != 42 {
		t.Fatalf("expected 42, got %f", v)
	}
}

func TestGetPanicsOutOfRange(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on out-of-range index")
		}
	}()
	a := NewNDArray([]int{2, 3}, nil)
	a.Get(2, 0)
}

func TestGetPanicsWrongDimCount(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on wrong index count")
		}
	}()
	a := NewNDArray([]int{2, 3}, nil)
	a.Get(0)
}

func TestReshape(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := a.Reshape(3, 2)
	if b.Shape()[0] != 3 || b.Shape()[1] != 2 {
		t.Fatalf("unexpected shape %v", b.Shape())
	}
	if b.Get(0, 0) != 1 || b.Get(2, 1) != 6 {
		t.Fatal("reshape data incorrect")
	}
}

func TestReshapePanicsOnSizeMismatch(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on reshape size mismatch")
		}
	}()
	a := NewNDArray([]int{2, 3}, nil)
	a.Reshape(4, 2)
}

func TestFlatten(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	f := a.Flatten()
	if f.Ndim() != 1 || f.Size() != 6 {
		t.Fatalf("flatten produced wrong shape: %v", f.Shape())
	}
	for i := 0; i < 6; i++ {
		if f.Get(i) != float64(i+1) {
			t.Fatalf("flatten data wrong at %d", i)
		}
	}
}

func TestCopy(t *testing.T) {
	a := NewNDArray([]int{2}, []float64{1, 2})
	b := a.Copy()
	b.Set(99, 0)
	if a.Get(0) != 1 {
		t.Fatal("copy is not independent")
	}
}

func TestTranspose2D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	at := a.T()
	if at.Shape()[0] != 3 || at.Shape()[1] != 2 {
		t.Fatalf("unexpected transposed shape %v", at.Shape())
	}
	// a[0][1]=2 should become at[1][0]=2
	if at.Get(1, 0) != 2 {
		t.Fatalf("expected at[1,0]=2, got %f", at.Get(1, 0))
	}
	// a[1][2]=6 should become at[2][1]=6
	if at.Get(2, 1) != 6 {
		t.Fatalf("expected at[2,1]=6, got %f", at.Get(2, 1))
	}
}

func TestTranspose1D(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	at := a.T()
	if at.Ndim() != 1 || at.Size() != 3 {
		t.Fatal("1-D transpose should return same shape")
	}
}

func TestTranspose3D(t *testing.T) {
	// 2x3x4 -> 4x3x2
	data := make([]float64, 24)
	for i := range data {
		data[i] = float64(i)
	}
	a := NewNDArray([]int{2, 3, 4}, data)
	at := a.T()
	if at.Shape()[0] != 4 || at.Shape()[1] != 3 || at.Shape()[2] != 2 {
		t.Fatalf("unexpected 3D transposed shape %v", at.Shape())
	}
	// a[1][2][3] = 1*12+2*4+3 = 23, should be at[3][2][1]
	if at.Get(3, 2, 1) != 23 {
		t.Fatalf("3D transpose wrong: expected 23, got %f", at.Get(3, 2, 1))
	}
}

func TestString(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	s := a.String()
	if !strings.Contains(s, "NDArray") {
		t.Fatal("String() should contain NDArray")
	}
	if !strings.Contains(s, "shape=") {
		t.Fatal("String() should contain shape info")
	}
}

func TestDataIndependence(t *testing.T) {
	original := []float64{1, 2, 3}
	a := NewNDArray([]int{3}, original)
	original[0] = 999
	if a.Get(0) != 1 {
		t.Fatal("NDArray should copy input data")
	}
	d := a.Data()
	d[0] = 888
	if a.Get(0) != 1 {
		t.Fatal("Data() should return a copy")
	}
}

func TestShapeIndependence(t *testing.T) {
	a := NewNDArray([]int{2, 3}, nil)
	s := a.Shape()
	s[0] = 99
	if a.Shape()[0] != 2 {
		t.Fatal("Shape() should return a copy")
	}
}
