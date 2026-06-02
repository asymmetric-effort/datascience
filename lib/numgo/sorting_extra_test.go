//go:build unit

package numgo

import (
	"testing"
)

// ---------------------------------------------------------------------------
// ArgMin
// ---------------------------------------------------------------------------

func TestArgMin1D(t *testing.T) {
	a := FromSlice([]float64{30, 10, 20})
	got := ArgMin(a, 0)
	assertData(t, "ArgMin1D", got, []float64{1})
}

func TestArgMin2DAxis0(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{3, 1, 1, 3})
	got := ArgMin(a, 0)
	// col 0: min at row 1, col 1: min at row 0
	assertData(t, "ArgMin2DAxis0", got, []float64{1, 0})
}

func TestArgMin2DAxis1(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{3, 1, 2, 6, 4, 5})
	got := ArgMin(a, 1)
	// row 0: min at col 1, row 1: min at col 1
	assertData(t, "ArgMin2DAxis1", got, []float64{1, 1})
}

// ---------------------------------------------------------------------------
// CountNonzero
// ---------------------------------------------------------------------------

func TestCountNonzeroGlobal(t *testing.T) {
	a := FromSlice([]float64{0, 1, 0, 3, 0})
	got := CountNonzero(a)
	assertData(t, "CountNonzeroGlobal", got, []float64{2})
}

func TestCountNonzeroAxis(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{0, 1, 0, 1, 0, 1})
	got := CountNonzero(a, 0)
	assertData(t, "CountNonzeroAxis0", got, []float64{1, 1, 1})
}

func TestCountNonzeroAllZero(t *testing.T) {
	a := FromSlice([]float64{0, 0, 0})
	got := CountNonzero(a)
	assertData(t, "CountNonzeroAllZero", got, []float64{0})
}

// ---------------------------------------------------------------------------
// Extract
// ---------------------------------------------------------------------------

func TestExtract(t *testing.T) {
	cond := FromSlice([]float64{1, 0, 1, 0, 1})
	a := FromSlice([]float64{10, 20, 30, 40, 50})
	got := Extract(cond, a)
	assertData(t, "Extract", got, []float64{10, 30, 50})
}

func TestExtractNone(t *testing.T) {
	cond := FromSlice([]float64{0, 0, 0})
	a := FromSlice([]float64{1, 2, 3})
	got := Extract(cond, a)
	if got.Size() != 0 {
		t.Fatalf("ExtractNone: expected size 0, got %d", got.Size())
	}
}

// ---------------------------------------------------------------------------
// Flatnonzero
// ---------------------------------------------------------------------------

func TestFlatnonzero(t *testing.T) {
	a := FromSlice([]float64{0, 3, 0, 5, 0})
	got := Flatnonzero(a)
	assertData(t, "Flatnonzero", got, []float64{1, 3})
}

func TestFlatnonzeroAllZero(t *testing.T) {
	a := FromSlice([]float64{0, 0})
	got := Flatnonzero(a)
	if got.Size() != 0 {
		t.Fatalf("FlatnonzeroAllZero: expected size 0, got %d", got.Size())
	}
}

// ---------------------------------------------------------------------------
// Argwhere
// ---------------------------------------------------------------------------

func TestArgwhere(t *testing.T) {
	a := FromSlice([]float64{0, 1, 0, 3})
	got := Argwhere(a)
	if len(got) != 2 {
		t.Fatalf("Argwhere: expected 2 results, got %d", len(got))
	}
	assertCoord(t, "Argwhere[0]", got[0], []int{1})
	assertCoord(t, "Argwhere[1]", got[1], []int{3})
}

// ---------------------------------------------------------------------------
// Lexsort
// ---------------------------------------------------------------------------

func TestLexsort(t *testing.T) {
	// Sort by last name (primary), then first name (secondary).
	// Keys: primary (last) = keys[1], secondary (first) = keys[0].
	lastName := FromSlice([]float64{2, 1, 2, 1})  // B, A, B, A
	firstName := FromSlice([]float64{1, 2, 2, 1}) // a, b, b, a
	got := Lexsort([]*NDArray{firstName, lastName})
	// Expected sort: (A,a)=idx3, (A,b)=idx1, (B,a)=idx0, (B,b)=idx2
	assertData(t, "Lexsort", got, []float64{3, 1, 0, 2})
}

func TestLexsortSingle(t *testing.T) {
	k := FromSlice([]float64{3, 1, 2})
	got := Lexsort([]*NDArray{k})
	assertData(t, "LexsortSingle", got, []float64{1, 2, 0})
}

// ---------------------------------------------------------------------------
// Partition
// ---------------------------------------------------------------------------

func TestPartition(t *testing.T) {
	a := FromSlice([]float64{3, 1, 2, 5, 4})
	got := Partition(a, 2, 0)
	d := got.Data()
	// Element at index 2 should be 3 (the 3rd smallest).
	if d[2] != 3 {
		t.Fatalf("Partition: element at kth=%d is %g, want 3", 2, d[2])
	}
	// All elements before kth should be <= d[kth].
	for i := 0; i < 2; i++ {
		if d[i] > d[2] {
			t.Fatalf("Partition: d[%d]=%g > d[kth]=%g", i, d[i], d[2])
		}
	}
	// All elements after kth should be >= d[kth].
	for i := 3; i < 5; i++ {
		if d[i] < d[2] {
			t.Fatalf("Partition: d[%d]=%g < d[kth]=%g", i, d[i], d[2])
		}
	}
}

// ---------------------------------------------------------------------------
// Argpartition
// ---------------------------------------------------------------------------

func TestArgpartition(t *testing.T) {
	a := FromSlice([]float64{3, 1, 2, 5, 4})
	got := Argpartition(a, 2, 0)
	d := got.Data()
	// The index at position kth should point to the kth smallest element.
	kthIdx := int(d[2])
	if a.Data()[kthIdx] != 3 {
		t.Fatalf("Argpartition: element at kth index is %g, want 3", a.Data()[kthIdx])
	}
}

func TestPartitionPanicsOnBadAxis(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for bad axis")
		}
	}()
	Partition(FromSlice([]float64{1, 2}), 0, 1)
}

func TestPartitionPanicsOnBadKth(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for bad kth")
		}
	}()
	Partition(FromSlice([]float64{1, 2}), 5, 0)
}
