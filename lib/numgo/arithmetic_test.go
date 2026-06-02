//go:build unit

package numgo

import (
	"math"
	"testing"
)

func TestAdd(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{4, 5, 6})
	c := Add(a, b)
	expect := []float64{5, 7, 9}
	for i, v := range c.Data() {
		if v != expect[i] {
			t.Fatalf("Add: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestAddScalar(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	c := AddScalar(a, 10)
	expect := []float64{11, 12, 13}
	for i, v := range c.Data() {
		if v != expect[i] {
			t.Fatalf("AddScalar: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestSub(t *testing.T) {
	a := FromSlice([]float64{10, 20, 30})
	b := FromSlice([]float64{1, 2, 3})
	c := Sub(a, b)
	expect := []float64{9, 18, 27}
	for i, v := range c.Data() {
		if v != expect[i] {
			t.Fatalf("Sub: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestSubScalar(t *testing.T) {
	a := FromSlice([]float64{10, 20})
	c := SubScalar(a, 5)
	if c.Data()[0] != 5 || c.Data()[1] != 15 {
		t.Fatalf("SubScalar: unexpected result %v", c.Data())
	}
}

func TestMul(t *testing.T) {
	a := FromSlice([]float64{2, 3})
	b := FromSlice([]float64{4, 5})
	c := Mul(a, b)
	if c.Data()[0] != 8 || c.Data()[1] != 15 {
		t.Fatalf("Mul: unexpected result %v", c.Data())
	}
}

func TestMulScalar(t *testing.T) {
	a := FromSlice([]float64{2, 3})
	c := MulScalar(a, 10)
	if c.Data()[0] != 20 || c.Data()[1] != 30 {
		t.Fatal("MulScalar: wrong")
	}
}

func TestDiv(t *testing.T) {
	a := FromSlice([]float64{10, 20})
	b := FromSlice([]float64{2, 5})
	c := Div(a, b)
	if c.Data()[0] != 5 || c.Data()[1] != 4 {
		t.Fatalf("Div: unexpected result %v", c.Data())
	}
}

func TestDivScalar(t *testing.T) {
	a := FromSlice([]float64{10, 20})
	c := DivScalar(a, 2)
	if c.Data()[0] != 5 || c.Data()[1] != 10 {
		t.Fatal("DivScalar: wrong")
	}
}

func TestShapeMismatchPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on shape mismatch")
		}
	}()
	a := FromSlice([]float64{1, 2})
	b := FromSlice([]float64{1, 2, 3})
	Add(a, b)
}

func TestSumAll(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	s := Sum(a)
	if s.Data()[0] != 21 {
		t.Fatalf("Sum all: expected 21, got %f", s.Data()[0])
	}
}

func TestSumAxis0(t *testing.T) {
	// [[1,2,3],[4,5,6]] -> sum axis 0 -> [5,7,9]
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	s := Sum(a, 0)
	expect := []float64{5, 7, 9}
	if s.Ndim() != 1 || s.Size() != 3 {
		t.Fatalf("Sum axis 0: unexpected shape %v", s.Shape())
	}
	for i, v := range s.Data() {
		if v != expect[i] {
			t.Fatalf("Sum axis 0: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestSumAxis1(t *testing.T) {
	// [[1,2,3],[4,5,6]] -> sum axis 1 -> [6,15]
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	s := Sum(a, 1)
	expect := []float64{6, 15}
	if s.Size() != 2 {
		t.Fatalf("Sum axis 1: unexpected shape %v", s.Shape())
	}
	for i, v := range s.Data() {
		if v != expect[i] {
			t.Fatalf("Sum axis 1: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestProdAll(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	p := Prod(a)
	if p.Data()[0] != 24 {
		t.Fatalf("Prod all: expected 24, got %f", p.Data()[0])
	}
}

func TestProdAxis(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	p := Prod(a, 0) // [4,10,18]
	expect := []float64{4, 10, 18}
	for i, v := range p.Data() {
		if v != expect[i] {
			t.Fatalf("Prod axis 0: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestMaxAll(t *testing.T) {
	a := FromSlice([]float64{3, 1, 4, 1, 5, 9})
	m := Max(a)
	if m.Data()[0] != 9 {
		t.Fatalf("Max all: expected 9, got %f", m.Data()[0])
	}
}

func TestMaxAxis(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 5, 3, 4, 2, 6})
	m := Max(a, 1) // [5,6]
	expect := []float64{5, 6}
	for i, v := range m.Data() {
		if v != expect[i] {
			t.Fatalf("Max axis 1: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestArgMax(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 5, 3, 4, 2, 6})
	am := ArgMax(a, 1) // [1,2]
	expect := []float64{1, 2}
	for i, v := range am.Data() {
		if v != expect[i] {
			t.Fatalf("ArgMax axis 1: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestArgMaxAxis0(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 5, 3, 4, 2, 6})
	am := ArgMax(a, 0) // [1,0,1]
	expect := []float64{1, 0, 1}
	for i, v := range am.Data() {
		if v != expect[i] {
			t.Fatalf("ArgMax axis 0: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestSumMultipleAxes(t *testing.T) {
	// 2x3x4 array, sum along axes 0 and 2
	data := make([]float64, 24)
	for i := range data {
		data[i] = float64(i)
	}
	a := NewNDArray([]int{2, 3, 4}, data)
	s := Sum(a, 0, 2)
	// After reducing axis 0 on 2x3x4 -> 3x4
	// After reducing axis 1 (originally 2) on 3x4 -> 3
	// Expected: sum over first and last axes for each middle index
	// s[j] = sum over i,k of a[i][j][k]
	// j=0: 0+1+2+3+12+13+14+15 = 60
	// j=1: 4+5+6+7+16+17+18+19 = 92
	// j=2: 8+9+10+11+20+21+22+23 = 124
	if s.Size() != 3 {
		t.Fatalf("expected size 3, got %d (shape %v)", s.Size(), s.Shape())
	}
	expect := []float64{60, 92, 124}
	for i, v := range s.Data() {
		if math.Abs(v-expect[i]) > 1e-10 {
			t.Fatalf("Sum multi-axes: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}

func TestAxisOutOfRangePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on axis out of range")
		}
	}()
	a := FromSlice([]float64{1, 2, 3})
	Sum(a, 1)
}

func TestAdd2D(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	b := NewNDArray([]int{2, 2}, []float64{10, 20, 30, 40})
	c := Add(a, b)
	expect := []float64{11, 22, 33, 44}
	for i, v := range c.Data() {
		if v != expect[i] {
			t.Fatalf("Add2D: at %d expected %f, got %f", i, expect[i], v)
		}
	}
}
