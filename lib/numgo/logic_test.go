//go:build unit

package numgo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// All / Any
// ---------------------------------------------------------------------------

func TestAllTrue(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	got := All(a)
	assertData(t, "AllTrue", got, []float64{1})
}

func TestAllFalse(t *testing.T) {
	a := FromSlice([]float64{1, 0, 3})
	got := All(a)
	assertData(t, "AllFalse", got, []float64{0})
}

func TestAllAxis(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 0, 1, 1})
	got := All(a, 0) // col 0: all nonzero, col 1: not
	assertData(t, "AllAxis0", got, []float64{1, 0})
}

func TestAnyTrue(t *testing.T) {
	a := FromSlice([]float64{0, 0, 1})
	got := Any(a)
	assertData(t, "AnyTrue", got, []float64{1})
}

func TestAnyFalse(t *testing.T) {
	a := FromSlice([]float64{0, 0, 0})
	got := Any(a)
	assertData(t, "AnyFalse", got, []float64{0})
}

func TestAnyAxis(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{0, 0, 1, 0})
	got := Any(a, 0) // col 0: has 1, col 1: all 0
	assertData(t, "AnyAxis0", got, []float64{1, 0})
}

// ---------------------------------------------------------------------------
// Isnan, Isinf, Isfinite, Isneginf, Isposinf
// ---------------------------------------------------------------------------

func TestIsnan(t *testing.T) {
	a := FromSlice([]float64{1, math.NaN(), 3})
	got := Isnan(a)
	assertData(t, "Isnan", got, []float64{0, 1, 0})
}

func TestIsinf(t *testing.T) {
	a := FromSlice([]float64{math.Inf(1), 1, math.Inf(-1)})
	got := Isinf(a)
	assertData(t, "Isinf", got, []float64{1, 0, 1})
}

func TestIsfinite(t *testing.T) {
	a := FromSlice([]float64{1, math.Inf(1), math.NaN()})
	got := Isfinite(a)
	assertData(t, "Isfinite", got, []float64{1, 0, 0})
}

func TestIsneginf(t *testing.T) {
	a := FromSlice([]float64{math.Inf(-1), math.Inf(1), 0})
	got := Isneginf(a)
	assertData(t, "Isneginf", got, []float64{1, 0, 0})
}

func TestIsposinf(t *testing.T) {
	a := FromSlice([]float64{math.Inf(-1), math.Inf(1), 0})
	got := Isposinf(a)
	assertData(t, "Isposinf", got, []float64{0, 1, 0})
}

// ---------------------------------------------------------------------------
// Logical operations
// ---------------------------------------------------------------------------

func TestLogicalAnd(t *testing.T) {
	a := FromSlice([]float64{1, 0, 1, 0})
	b := FromSlice([]float64{1, 1, 0, 0})
	got := LogicalAnd(a, b)
	assertData(t, "LogicalAnd", got, []float64{1, 0, 0, 0})
}

func TestLogicalOr(t *testing.T) {
	a := FromSlice([]float64{1, 0, 1, 0})
	b := FromSlice([]float64{1, 1, 0, 0})
	got := LogicalOr(a, b)
	assertData(t, "LogicalOr", got, []float64{1, 1, 1, 0})
}

func TestLogicalNot(t *testing.T) {
	a := FromSlice([]float64{1, 0, 3, 0})
	got := LogicalNot(a)
	assertData(t, "LogicalNot", got, []float64{0, 1, 0, 1})
}

func TestLogicalXor(t *testing.T) {
	a := FromSlice([]float64{1, 0, 1, 0})
	b := FromSlice([]float64{1, 1, 0, 0})
	got := LogicalXor(a, b)
	assertData(t, "LogicalXor", got, []float64{0, 1, 1, 0})
}

// ---------------------------------------------------------------------------
// ArrayEqual / ArrayEquiv
// ---------------------------------------------------------------------------

func TestArrayEqualTrue(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{1, 2, 3})
	if !ArrayEqual(a, b) {
		t.Fatal("ArrayEqual: expected true")
	}
}

func TestArrayEqualFalse(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{1, 2, 4})
	if ArrayEqual(a, b) {
		t.Fatal("ArrayEqual: expected false")
	}
}

func TestArrayEqualDiffShape(t *testing.T) {
	a := FromSlice([]float64{1, 2})
	b := FromSlice([]float64{1, 2, 3})
	if ArrayEqual(a, b) {
		t.Fatal("ArrayEqual: expected false for different shapes")
	}
}

func TestArrayEquivBroadcast(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{5, 5, 5, 5})
	b := FromSlice([]float64{5})
	if !ArrayEquiv(a, b) {
		t.Fatal("ArrayEquiv: expected true after broadcasting")
	}
}

func TestArrayEquivFalse(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{5, 5, 5, 6})
	b := FromSlice([]float64{5})
	if ArrayEquiv(a, b) {
		t.Fatal("ArrayEquiv: expected false")
	}
}

// ---------------------------------------------------------------------------
// Comparison functions
// ---------------------------------------------------------------------------

func TestGreater(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{2, 2, 1})
	got := Greater(a, b)
	assertData(t, "Greater", got, []float64{0, 0, 1})
}

func TestLess(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{2, 2, 1})
	got := Less(a, b)
	assertData(t, "Less", got, []float64{1, 0, 0})
}

func TestEqual(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{1, 3, 3})
	got := Equal(a, b)
	assertData(t, "Equal", got, []float64{1, 0, 1})
}

func TestNotEqual(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{1, 3, 3})
	got := NotEqual(a, b)
	assertData(t, "NotEqual", got, []float64{0, 1, 0})
}

func TestGreaterEqual(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{2, 2, 1})
	got := GreaterEqual(a, b)
	assertData(t, "GreaterEqual", got, []float64{0, 1, 1})
}

func TestLessEqual(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{2, 2, 1})
	got := LessEqual(a, b)
	assertData(t, "LessEqual", got, []float64{1, 1, 0})
}

// ---------------------------------------------------------------------------
// Isclose
// ---------------------------------------------------------------------------

func TestIsclose(t *testing.T) {
	a := FromSlice([]float64{1.0, 2.0, 3.0})
	b := FromSlice([]float64{1.0001, 2.0, 3.1})
	got := Isclose(a, b, 0.001, 0)
	assertData(t, "Isclose", got, []float64{1, 1, 0})
}

func TestIscloseRtol(t *testing.T) {
	a := FromSlice([]float64{100})
	b := FromSlice([]float64{101})
	// rtol=0.02 => threshold = 0 + 0.02*101 = 2.02, |100-101|=1 < 2.02
	got := Isclose(a, b, 0, 0.02)
	assertData(t, "IscloseRtol", got, []float64{1})
}

// ---------------------------------------------------------------------------
// Broadcasting comparisons
// ---------------------------------------------------------------------------

func TestGreaterBroadcast(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 5, 3, 7})
	b := FromSlice([]float64{4})
	got := Greater(a, b)
	assertData(t, "GreaterBroadcast", got, []float64{0, 1, 0, 1})
}
