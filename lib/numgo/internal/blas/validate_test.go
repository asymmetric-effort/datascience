//go:build unit

package blas

import (
	"strings"
	"testing"
)

// --- Safe functions panic on short slices ---

func TestDdotPanicsOnShortX(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if !strings.Contains(r.(string), "Ddot x") {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	Ddot(10, make([]float64, 5), 1, make([]float64, 10), 1)
}

func TestDdotPanicsOnShortY(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if !strings.Contains(r.(string), "Ddot y") {
			t.Errorf("unexpected panic: %v", r)
		}
	}()
	Ddot(10, make([]float64, 10), 1, make([]float64, 5), 1)
}

func TestDaxpyPanicsOnShortSlice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	Daxpy(10, 1.0, make([]float64, 5), 1, make([]float64, 10), 1)
}

func TestDscalPanicsOnShortSlice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	Dscal(10, 2.0, make([]float64, 5), 1)
}

func TestDnrm2PanicsOnShortSlice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	Dnrm2(10, make([]float64, 5), 1)
}

func TestDasumPanicsOnShortSlice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	Dasum(10, make([]float64, 5), 1)
}

func TestIdamaxPanicsOnShortSlice(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	Idamax(10, make([]float64, 5), 1)
}

func TestDgemvPanicsOnShortMatrix(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	// A should be 3x4=12 elements, but we pass 5
	Dgemv(false, 3, 4, 1.0, make([]float64, 5), 4,
		make([]float64, 4), 1, 0.0, make([]float64, 3), 1)
}

func TestDtrsvPanicsOnShortVector(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	// n=4 but x has 2 elements
	a := make([]float64, 16)
	for i := 0; i < 4; i++ {
		a[i*4+i] = 1 // identity
	}
	Dtrsv('U', 'N', 'N', 4, a, 4, make([]float64, 2), 1)
}

func TestDgemmPanicsOnShortC(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	// m=3, n=3, k=3, but C has only 5 elements (needs 9)
	Dgemm(false, false, 3, 3, 3, 1.0,
		make([]float64, 9), 3, make([]float64, 9), 3,
		0.0, make([]float64, 5), 3)
}

// --- Unsafe functions do NOT panic (just verify no validation) ---

func TestDdotUnsafeNoPanic(t *testing.T) {
	// Correct-length slices, should work fine.
	x := []float64{1, 2, 3}
	y := []float64{4, 5, 6}
	got := DdotUnsafe(3, x, 1, y, 1)
	want := 1*4.0 + 2*5.0 + 3*6.0
	if got != want {
		t.Errorf("DdotUnsafe = %f, want %f", got, want)
	}
}

func TestDgemmUnsafeNoPanic(t *testing.T) {
	// 2x2 identity * identity = identity
	a := []float64{1, 0, 0, 1}
	b := []float64{1, 0, 0, 1}
	c := make([]float64, 4)
	DgemmUnsafe(false, false, 2, 2, 2, 1.0, a, 2, b, 2, 0.0, c, 2)
	if c[0] != 1 || c[1] != 0 || c[2] != 0 || c[3] != 1 {
		t.Errorf("DgemmUnsafe I*I = %v, want [1 0 0 1]", c)
	}
}

// --- Stride validation ---

func TestDdotPanicsOnShortStridedVector(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for strided access beyond slice")
		}
	}()
	// n=3, incx=3 needs indices 0,3,6 → length >= 7
	Ddot(3, make([]float64, 5), 3, make([]float64, 10), 1)
}
