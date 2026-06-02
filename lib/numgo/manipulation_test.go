//go:build unit

package numgo

import (
	"testing"
)

func sliceEqual(a, b []float64) bool {
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

func TestRavel(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	r := Ravel(a)
	if r.Ndim() != 1 {
		t.Fatalf("expected 1D, got %dD", r.Ndim())
	}
	if r.Size() != 6 {
		t.Fatalf("expected size 6, got %d", r.Size())
	}
	expected := []float64{1, 2, 3, 4, 5, 6}
	if !sliceEqual(r.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, r.Data())
	}
	// Ensure it's a copy.
	r.Set(99, 0)
	if a.Get(0, 0) == 99 {
		t.Fatal("Ravel should return a copy, not a view")
	}
}

func TestSwapaxes(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	s := Swapaxes(a, 0, 1)
	if !shapeEqual(s.Shape(), []int{3, 2}) {
		t.Fatalf("expected shape [3 2], got %v", s.Shape())
	}
	// Swapping axes of a 2D array is a transpose.
	if s.Get(0, 0) != 1 || s.Get(1, 0) != 2 || s.Get(0, 1) != 4 {
		t.Fatalf("unexpected values after swap: %v", s.Data())
	}
}

func TestMoveaxis(t *testing.T) {
	a := NewNDArray([]int{2, 3, 4}, nil)
	m := Moveaxis(a, 0, 2)
	if !shapeEqual(m.Shape(), []int{3, 4, 2}) {
		t.Fatalf("expected shape [3 4 2], got %v", m.Shape())
	}
}

func TestRollaxis(t *testing.T) {
	a := NewNDArray([]int{2, 3, 4}, nil)
	r := Rollaxis(a, 2, 0)
	if !shapeEqual(r.Shape(), []int{4, 2, 3}) {
		t.Fatalf("expected shape [4 2 3], got %v", r.Shape())
	}
}

func TestExpandDims(t *testing.T) {
	a := NewNDArray([]int{3, 4}, []float64{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
	})
	e := ExpandDims(a, 0)
	if !shapeEqual(e.Shape(), []int{1, 3, 4}) {
		t.Fatalf("expected shape [1 3 4], got %v", e.Shape())
	}
	e2 := ExpandDims(a, 2)
	if !shapeEqual(e2.Shape(), []int{3, 4, 1}) {
		t.Fatalf("expected shape [3 4 1], got %v", e2.Shape())
	}
}

func TestSqueeze(t *testing.T) {
	a := NewNDArray([]int{1, 3, 1, 4}, nil)
	s := Squeeze(a)
	if !shapeEqual(s.Shape(), []int{3, 4}) {
		t.Fatalf("expected shape [3 4], got %v", s.Shape())
	}
}

func TestSqueezeAllOnes(t *testing.T) {
	a := NewNDArray([]int{1, 1, 1}, []float64{42})
	s := Squeeze(a)
	if !shapeEqual(s.Shape(), []int{1}) {
		t.Fatalf("expected shape [1], got %v", s.Shape())
	}
	if s.Get(0) != 42 {
		t.Fatalf("expected 42, got %f", s.Get(0))
	}
}

func TestConcatenate1D(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	b := NewNDArray([]int{2}, []float64{4, 5})
	c, err := Concatenate([]*NDArray{a, b}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(c.Shape(), []int{5}) {
		t.Fatalf("expected shape [5], got %v", c.Shape())
	}
	expected := []float64{1, 2, 3, 4, 5}
	if !sliceEqual(c.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, c.Data())
	}
}

func TestConcatenate2D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{2, 2}, []float64{7, 8, 9, 10})
	c, err := Concatenate([]*NDArray{a, b}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(c.Shape(), []int{2, 5}) {
		t.Fatalf("expected shape [2 5], got %v", c.Shape())
	}
	expected := []float64{1, 2, 3, 7, 8, 4, 5, 6, 9, 10}
	if !sliceEqual(c.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, c.Data())
	}
}

func TestConcatenate2DAxis0(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{1, 3}, []float64{7, 8, 9})
	c, err := Concatenate([]*NDArray{a, b}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(c.Shape(), []int{3, 3}) {
		t.Fatalf("expected shape [3 3], got %v", c.Shape())
	}
	expected := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !sliceEqual(c.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, c.Data())
	}
}

func TestStack1D(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	b := NewNDArray([]int{3}, []float64{4, 5, 6})
	s, err := Stack([]*NDArray{a, b}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(s.Shape(), []int{2, 3}) {
		t.Fatalf("expected shape [2 3], got %v", s.Shape())
	}
	expected := []float64{1, 2, 3, 4, 5, 6}
	if !sliceEqual(s.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, s.Data())
	}
}

func TestStack1DAxis1(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	b := NewNDArray([]int{3}, []float64{4, 5, 6})
	s, err := Stack([]*NDArray{a, b}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(s.Shape(), []int{3, 2}) {
		t.Fatalf("expected shape [3 2], got %v", s.Shape())
	}
	// Column-stacked: [[1,4],[2,5],[3,6]]
	expected := []float64{1, 4, 2, 5, 3, 6}
	if !sliceEqual(s.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, s.Data())
	}
}

func TestStack2D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{2, 3}, []float64{7, 8, 9, 10, 11, 12})
	s, err := Stack([]*NDArray{a, b}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(s.Shape(), []int{2, 2, 3}) {
		t.Fatalf("expected shape [2 2 3], got %v", s.Shape())
	}
	expected := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	if !sliceEqual(s.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, s.Data())
	}
}

func TestVstack(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	b := NewNDArray([]int{3}, []float64{4, 5, 6})
	v, err := Vstack([]*NDArray{a, b})
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(v.Shape(), []int{2, 3}) {
		t.Fatalf("expected shape [2 3], got %v", v.Shape())
	}
}

func TestHstack(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	b := NewNDArray([]int{2}, []float64{4, 5})
	h, err := Hstack([]*NDArray{a, b})
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(h.Shape(), []int{5}) {
		t.Fatalf("expected shape [5], got %v", h.Shape())
	}
}

func TestHstack2D(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	b := NewNDArray([]int{2, 1}, []float64{5, 6})
	h, err := Hstack([]*NDArray{a, b})
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(h.Shape(), []int{2, 3}) {
		t.Fatalf("expected shape [2 3], got %v", h.Shape())
	}
}

func TestDstack(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{2, 3}, []float64{7, 8, 9, 10, 11, 12})
	d, err := Dstack([]*NDArray{a, b})
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(d.Shape(), []int{2, 3, 2}) {
		t.Fatalf("expected shape [2 3 2], got %v", d.Shape())
	}
}

func TestSplitAndConcatenateRoundTrip(t *testing.T) {
	a := NewNDArray([]int{6}, []float64{1, 2, 3, 4, 5, 6})
	parts, err := Split(a, 3, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts, got %d", len(parts))
	}
	// Each part should be [2].
	for i, p := range parts {
		if !shapeEqual(p.Shape(), []int{2}) {
			t.Fatalf("part %d: expected shape [2], got %v", i, p.Shape())
		}
	}
	// Round-trip: concatenate should give back original.
	restored, err := Concatenate(parts, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !sliceEqual(restored.Data(), a.Data()) {
		t.Fatalf("round-trip failed: expected %v, got %v", a.Data(), restored.Data())
	}
}

func TestSplit2DRoundTrip(t *testing.T) {
	a := NewNDArray([]int{4, 3}, []float64{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
		10, 11, 12,
	})
	parts, err := Split(a, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	for i, p := range parts {
		if !shapeEqual(p.Shape(), []int{2, 3}) {
			t.Fatalf("part %d: expected shape [2 3], got %v", i, p.Shape())
		}
	}
	restored, err := Concatenate(parts, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !sliceEqual(restored.Data(), a.Data()) {
		t.Fatalf("round-trip failed: expected %v, got %v", a.Data(), restored.Data())
	}
}

func TestHsplit(t *testing.T) {
	a := NewNDArray([]int{2, 4}, []float64{1, 2, 3, 4, 5, 6, 7, 8})
	parts, err := Hsplit(a, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	if !shapeEqual(parts[0].Shape(), []int{2, 2}) {
		t.Fatalf("expected shape [2 2], got %v", parts[0].Shape())
	}
}

func TestVsplit(t *testing.T) {
	a := NewNDArray([]int{4, 2}, []float64{1, 2, 3, 4, 5, 6, 7, 8})
	parts, err := Vsplit(a, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
}

func TestDsplit(t *testing.T) {
	a := NewNDArray([]int{2, 2, 4}, nil)
	parts, err := Dsplit(a, 2)
	if err != nil {
		t.Fatal(err)
	}
	if len(parts) != 2 {
		t.Fatalf("expected 2 parts, got %d", len(parts))
	}
	if !shapeEqual(parts[0].Shape(), []int{2, 2, 2}) {
		t.Fatalf("expected shape [2 2 2], got %v", parts[0].Shape())
	}
}

func TestTile(t *testing.T) {
	a := NewNDArray([]int{2}, []float64{1, 2})
	tiled := Tile(a, []int{3})
	if !shapeEqual(tiled.Shape(), []int{6}) {
		t.Fatalf("expected shape [6], got %v", tiled.Shape())
	}
	expected := []float64{1, 2, 1, 2, 1, 2}
	if !sliceEqual(tiled.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, tiled.Data())
	}
}

func TestTile2D(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	tiled := Tile(a, []int{2, 3})
	if !shapeEqual(tiled.Shape(), []int{4, 6}) {
		t.Fatalf("expected shape [4 6], got %v", tiled.Shape())
	}
}

func TestRepeat(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	r := Repeat(a, 2, 0)
	if !shapeEqual(r.Shape(), []int{6}) {
		t.Fatalf("expected shape [6], got %v", r.Shape())
	}
	expected := []float64{1, 1, 2, 2, 3, 3}
	if !sliceEqual(r.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, r.Data())
	}
}

func TestRepeat2D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	r := Repeat(a, 2, 0)
	if !shapeEqual(r.Shape(), []int{4, 3}) {
		t.Fatalf("expected shape [4 3], got %v", r.Shape())
	}
	expected := []float64{1, 2, 3, 1, 2, 3, 4, 5, 6, 4, 5, 6}
	if !sliceEqual(r.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, r.Data())
	}
}

func TestDelete(t *testing.T) {
	a := NewNDArray([]int{5}, []float64{1, 2, 3, 4, 5})
	d, err := Delete(a, []int{1, 3}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(d.Shape(), []int{3}) {
		t.Fatalf("expected shape [3], got %v", d.Shape())
	}
	expected := []float64{1, 3, 5}
	if !sliceEqual(d.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, d.Data())
	}
}

func TestInsert(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	vals := NewNDArray([]int{2}, []float64{10, 20})
	ins, err := Insert(a, 1, vals, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(ins.Shape(), []int{5}) {
		t.Fatalf("expected shape [5], got %v", ins.Shape())
	}
	expected := []float64{1, 10, 20, 2, 3}
	if !sliceEqual(ins.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, ins.Data())
	}
}

func TestAppend(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	vals := NewNDArray([]int{2}, []float64{4, 5})
	app, err := Append(a, vals, 0)
	if err != nil {
		t.Fatal(err)
	}
	expected := []float64{1, 2, 3, 4, 5}
	if !sliceEqual(app.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, app.Data())
	}
}

func TestFlip(t *testing.T) {
	a := NewNDArray([]int{4}, []float64{1, 2, 3, 4})
	f := Flip(a, 0)
	expected := []float64{4, 3, 2, 1}
	if !sliceEqual(f.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, f.Data())
	}
}

func TestFliplr(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	f := Fliplr(a)
	expected := []float64{3, 2, 1, 6, 5, 4}
	if !sliceEqual(f.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, f.Data())
	}
}

func TestFlipud(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	f := Flipud(a)
	expected := []float64{4, 5, 6, 1, 2, 3}
	if !sliceEqual(f.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, f.Data())
	}
}

func TestRot90(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	r := Rot90(a, 1)
	if !shapeEqual(r.Shape(), []int{3, 2}) {
		t.Fatalf("expected shape [3 2], got %v", r.Shape())
	}
	// 90 CCW of [[1,2,3],[4,5,6]] = [[3,6],[2,5],[1,4]]
	expected := []float64{3, 6, 2, 5, 1, 4}
	if !sliceEqual(r.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, r.Data())
	}
}

func TestRot90_360(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	r := Rot90(a, 4)
	if !sliceEqual(r.Data(), a.Data()) {
		t.Fatalf("Rot90 by 4 should be identity, got %v", r.Data())
	}
}

func TestRoll(t *testing.T) {
	a := NewNDArray([]int{5}, []float64{1, 2, 3, 4, 5})
	r := Roll(a, 2, 0)
	expected := []float64{4, 5, 1, 2, 3}
	if !sliceEqual(r.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, r.Data())
	}
}

func TestRollNegative(t *testing.T) {
	a := NewNDArray([]int{5}, []float64{1, 2, 3, 4, 5})
	r := Roll(a, -2, 0)
	expected := []float64{3, 4, 5, 1, 2}
	if !sliceEqual(r.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, r.Data())
	}
}

func TestRoll2D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	r := Roll(a, 1, 1)
	// Each row shifted right by 1: [[3,1,2],[6,4,5]]
	expected := []float64{3, 1, 2, 6, 4, 5}
	if !sliceEqual(r.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, r.Data())
	}
}

func TestRollFlat(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	r := Roll(a, 2, -1)
	if !shapeEqual(r.Shape(), []int{2, 3}) {
		t.Fatalf("expected shape [2 3], got %v", r.Shape())
	}
	expected := []float64{5, 6, 1, 2, 3, 4}
	if !sliceEqual(r.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, r.Data())
	}
}

func TestConcatenateError(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	b := NewNDArray([]int{3, 3}, nil)
	_, err := Concatenate([]*NDArray{a, b}, 1)
	if err == nil {
		t.Fatal("expected error for shape mismatch on non-concat axis")
	}
}

func TestStackError(t *testing.T) {
	a := NewNDArray([]int{3}, []float64{1, 2, 3})
	b := NewNDArray([]int{4}, []float64{1, 2, 3, 4})
	_, err := Stack([]*NDArray{a, b}, 0)
	if err == nil {
		t.Fatal("expected error for shape mismatch")
	}
}

func TestSplitError(t *testing.T) {
	a := NewNDArray([]int{5}, []float64{1, 2, 3, 4, 5})
	_, err := Split(a, 3, 0)
	if err == nil {
		t.Fatal("expected error for uneven split")
	}
}

func TestDelete2D(t *testing.T) {
	a := NewNDArray([]int{3, 3}, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9})
	d, err := Delete(a, []int{1}, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(d.Shape(), []int{2, 3}) {
		t.Fatalf("expected shape [2 3], got %v", d.Shape())
	}
	expected := []float64{1, 2, 3, 7, 8, 9}
	if !sliceEqual(d.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, d.Data())
	}
}

func TestInsert2D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	vals := NewNDArray([]int{1, 3}, []float64{10, 20, 30})
	ins, err := Insert(a, 1, vals, 0)
	if err != nil {
		t.Fatal(err)
	}
	if !shapeEqual(ins.Shape(), []int{3, 3}) {
		t.Fatalf("expected shape [3 3], got %v", ins.Shape())
	}
	expected := []float64{1, 2, 3, 10, 20, 30, 4, 5, 6}
	if !sliceEqual(ins.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, ins.Data())
	}
}

func TestFlip2D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	f := Flip(a, 0)
	expected := []float64{4, 5, 6, 1, 2, 3}
	if !sliceEqual(f.Data(), expected) {
		t.Fatalf("expected %v, got %v", expected, f.Data())
	}
}
