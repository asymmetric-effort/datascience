//go:build unit

package numgo

import "testing"

func TestAllCloseIdentical(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	b := FromSlice([]float64{1, 2, 3})
	if !AllClose(a, b, 0, 0) {
		t.Fatal("identical arrays should be AllClose")
	}
}

func TestAllCloseWithTolerance(t *testing.T) {
	a := FromSlice([]float64{1.0, 2.0, 3.0})
	b := FromSlice([]float64{1.01, 2.01, 3.01})
	if !AllClose(a, b, 0.02, 0) {
		t.Fatal("should be close with atol=0.02")
	}
	if AllClose(a, b, 0.001, 0) {
		t.Fatal("should NOT be close with atol=0.001")
	}
}

func TestAllCloseRelativeTolerance(t *testing.T) {
	a := FromSlice([]float64{100})
	b := FromSlice([]float64{101})
	// |100-101| = 1 <= 0 + 0.02*101 = 2.02 -> true
	if !AllClose(a, b, 0, 0.02) {
		t.Fatal("should be close with rtol=0.02")
	}
	// |100-101| = 1 <= 0 + 0.005*101 = 0.505 -> false
	if AllClose(a, b, 0, 0.005) {
		t.Fatal("should NOT be close with rtol=0.005")
	}
}

func TestAllCloseShapeMismatch(t *testing.T) {
	a := FromSlice([]float64{1, 2})
	b := FromSlice([]float64{1, 2, 3})
	if AllClose(a, b, 1e-8, 1e-5) {
		t.Fatal("different shapes should not be AllClose")
	}
}

func TestAllClose2D(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	b := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4.0001})
	if !AllClose(a, b, 1e-3, 0) {
		t.Fatal("should be close with atol=1e-3")
	}
	if AllClose(a, b, 1e-5, 0) {
		t.Fatal("should NOT be close with atol=1e-5")
	}
}

func TestAllCloseZeros(t *testing.T) {
	a := Zeros(5)
	b := Zeros(5)
	if !AllClose(a, b, 0, 0) {
		t.Fatal("two zero arrays should be AllClose")
	}
}

func TestAllCloseDifferentDimensions(t *testing.T) {
	a := NewNDArray([]int{4}, []float64{1, 2, 3, 4})
	b := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	if AllClose(a, b, 0, 0) {
		t.Fatal("different ndim should not be AllClose")
	}
}
