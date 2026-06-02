//go:build unit

package numgo

import (
	"testing"
)

// ---------------------------------------------------------------------------
// In1d
// ---------------------------------------------------------------------------

func TestIn1d(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4, 5})
	b := FromSlice([]float64{2, 4, 6})
	got := In1d(a, b)
	assertData(t, "In1d", got, []float64{0, 1, 0, 1, 0})
}

func TestIn1dEmpty(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := NewNDArray([]int{0}, []float64{})
	got := In1d(a, b)
	assertData(t, "In1dEmpty", got, []float64{0, 0, 0})
}

func TestIn1dAllMatch(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{1, 2, 3, 4})
	got := In1d(a, b)
	assertData(t, "In1dAllMatch", got, []float64{1, 1, 1})
}

// ---------------------------------------------------------------------------
// Setxor1d
// ---------------------------------------------------------------------------

func TestSetxor1d(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	b := FromSlice([]float64{3, 4, 5, 6})
	got := Setxor1d(a, b)
	assertData(t, "Setxor1d", got, []float64{1, 2, 5, 6})
}

func TestSetxor1dNoOverlap(t *testing.T) {
	a := FromSlice([]float64{1, 2})
	b := FromSlice([]float64{3, 4})
	got := Setxor1d(a, b)
	assertData(t, "Setxor1dNoOverlap", got, []float64{1, 2, 3, 4})
}

func TestSetxor1dFullOverlap(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{1, 2, 3})
	got := Setxor1d(a, b)
	if got.Size() != 0 {
		t.Fatalf("Setxor1dFullOverlap: expected empty, got size %d", got.Size())
	}
}

func TestSetxor1dWithDuplicates(t *testing.T) {
	a := FromSlice([]float64{1, 1, 2, 2, 3})
	b := FromSlice([]float64{2, 3, 3, 4})
	got := Setxor1d(a, b)
	assertData(t, "Setxor1dDups", got, []float64{1, 4})
}
