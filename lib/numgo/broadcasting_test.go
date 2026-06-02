//go:build unit

package numgo

import (
	"math"
	"testing"
)

func sliceEq(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func floatSliceEq(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > 1e-12 {
			return false
		}
	}
	return true
}

// --- BroadcastShapes tests ---

func TestBroadcastShapes_SameShape(t *testing.T) {
	s, err := BroadcastShapes([]int{3, 4}, []int{3, 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sliceEq(s, []int{3, 4}) {
		t.Fatalf("expected [3 4], got %v", s)
	}
}

func TestBroadcastShapes_3x1_1x4(t *testing.T) {
	s, err := BroadcastShapes([]int{3, 1}, []int{1, 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sliceEq(s, []int{3, 4}) {
		t.Fatalf("expected [3 4], got %v", s)
	}
}

func TestBroadcastShapes_ScalarAnd1D(t *testing.T) {
	// scalar (shape [1]) + (5,) -> (5,)
	s, err := BroadcastShapes([]int{1}, []int{5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sliceEq(s, []int{5}) {
		t.Fatalf("expected [5], got %v", s)
	}
}

func TestBroadcastShapes_2x3_3(t *testing.T) {
	// (2,3) + (3,) -> (2,3)
	s, err := BroadcastShapes([]int{2, 3}, []int{3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sliceEq(s, []int{2, 3}) {
		t.Fatalf("expected [2 3], got %v", s)
	}
}

func TestBroadcastShapes_2x3_2x1(t *testing.T) {
	// (2,3) + (2,1) -> (2,3)
	s, err := BroadcastShapes([]int{2, 3}, []int{2, 1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sliceEq(s, []int{2, 3}) {
		t.Fatalf("expected [2 3], got %v", s)
	}
}

func TestBroadcastShapes_DifferentNdim(t *testing.T) {
	// (3,) + (2,3) -> (2,3)
	s, err := BroadcastShapes([]int{3}, []int{2, 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sliceEq(s, []int{2, 3}) {
		t.Fatalf("expected [2 3], got %v", s)
	}
}

func TestBroadcastShapes_3D(t *testing.T) {
	// (1,3,1) + (2,1,4) -> (2,3,4)
	s, err := BroadcastShapes([]int{1, 3, 1}, []int{2, 1, 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sliceEq(s, []int{2, 3, 4}) {
		t.Fatalf("expected [2 3 4], got %v", s)
	}
}

func TestBroadcastShapes_Incompatible(t *testing.T) {
	_, err := BroadcastShapes([]int{3}, []int{4})
	if err == nil {
		t.Fatal("expected error for incompatible shapes (3,) and (4,)")
	}
}

func TestBroadcastShapes_Incompatible2D(t *testing.T) {
	_, err := BroadcastShapes([]int{2, 3}, []int{3, 2})
	if err == nil {
		t.Fatal("expected error for incompatible shapes (2,3) and (3,2)")
	}
}

// --- BroadcastTo tests ---

func TestBroadcastTo_1Dto2D(t *testing.T) {
	// (3,) -> (2,3)
	a := FromSlice([]float64{1, 2, 3})
	b, err := BroadcastTo(a, []int{2, 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sliceEq(b.Shape(), []int{2, 3}) {
		t.Fatalf("expected shape [2 3], got %v", b.Shape())
	}
	expect := []float64{1, 2, 3, 1, 2, 3}
	if !floatSliceEq(b.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, b.Data())
	}
}

func TestBroadcastTo_Column(t *testing.T) {
	// (3,1) -> (3,4)
	a := NewNDArray([]int{3, 1}, []float64{10, 20, 30})
	b, err := BroadcastTo(a, []int{3, 4})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := []float64{10, 10, 10, 10, 20, 20, 20, 20, 30, 30, 30, 30}
	if !floatSliceEq(b.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, b.Data())
	}
}

func TestBroadcastTo_ScalarTo1D(t *testing.T) {
	a := NewNDArray([]int{1}, []float64{7})
	b, err := BroadcastTo(a, []int{5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expect := []float64{7, 7, 7, 7, 7}
	if !floatSliceEq(b.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, b.Data())
	}
}

func TestBroadcastTo_Incompatible(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	_, err := BroadcastTo(a, []int{2, 4})
	if err == nil {
		t.Fatal("expected error broadcasting (3,) to (2,4)")
	}
}

// --- Arithmetic with broadcasting tests ---

func TestAdd_3x1_plus_1x4(t *testing.T) {
	// (3,1) + (1,4) -> (3,4)
	a := NewNDArray([]int{3, 1}, []float64{1, 2, 3})
	b := NewNDArray([]int{1, 4}, []float64{10, 20, 30, 40})
	c := Add(a, b)
	if !sliceEq(c.Shape(), []int{3, 4}) {
		t.Fatalf("expected shape [3 4], got %v", c.Shape())
	}
	expect := []float64{
		11, 21, 31, 41,
		12, 22, 32, 42,
		13, 23, 33, 43,
	}
	if !floatSliceEq(c.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, c.Data())
	}
}

func TestAdd_ScalarBroadcast(t *testing.T) {
	// (5,) + scalar (shape [1]) -> (5,)
	a := FromSlice([]float64{1, 2, 3, 4, 5})
	scalar := NewNDArray([]int{1}, []float64{10})
	c := Add(a, scalar)
	if !sliceEq(c.Shape(), []int{5}) {
		t.Fatalf("expected shape [5], got %v", c.Shape())
	}
	expect := []float64{11, 12, 13, 14, 15}
	if !floatSliceEq(c.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, c.Data())
	}
}

func TestAdd_2x3_plus_3(t *testing.T) {
	// (2,3) + (3,) -> (2,3)
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := FromSlice([]float64{10, 20, 30})
	c := Add(a, b)
	if !sliceEq(c.Shape(), []int{2, 3}) {
		t.Fatalf("expected shape [2 3], got %v", c.Shape())
	}
	expect := []float64{11, 22, 33, 14, 25, 36}
	if !floatSliceEq(c.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, c.Data())
	}
}

func TestAdd_2x3_plus_2x1(t *testing.T) {
	// (2,3) + (2,1) -> (2,3)
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{2, 1}, []float64{10, 20})
	c := Add(a, b)
	if !sliceEq(c.Shape(), []int{2, 3}) {
		t.Fatalf("expected shape [2 3], got %v", c.Shape())
	}
	expect := []float64{11, 12, 13, 24, 25, 26}
	if !floatSliceEq(c.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, c.Data())
	}
}

func TestSub_Broadcast(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{10, 20, 30, 40, 50, 60})
	b := FromSlice([]float64{1, 2, 3})
	c := Sub(a, b)
	expect := []float64{9, 18, 27, 39, 48, 57}
	if !floatSliceEq(c.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, c.Data())
	}
}

func TestMul_Broadcast(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{2, 1}, []float64{10, 100})
	c := Mul(a, b)
	expect := []float64{10, 20, 30, 400, 500, 600}
	if !floatSliceEq(c.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, c.Data())
	}
}

func TestDiv_Broadcast(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{10, 20, 30, 40, 50, 60})
	b := NewNDArray([]int{1, 3}, []float64{10, 10, 10})
	c := Div(a, b)
	expect := []float64{1, 2, 3, 4, 5, 6}
	if !floatSliceEq(c.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, c.Data())
	}
}

func TestAdd_SameShape_StillWorks(t *testing.T) {
	// Ensure same-shape operations still work after the refactor.
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{4, 5, 6})
	c := Add(a, b)
	expect := []float64{5, 7, 9}
	if !floatSliceEq(c.Data(), expect) {
		t.Fatalf("expected %v, got %v", expect, c.Data())
	}
}

func TestAdd_IncompatiblePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for incompatible shapes")
		}
	}()
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{1, 2, 3, 4})
	Add(a, b)
}
