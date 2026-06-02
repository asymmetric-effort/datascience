//go:build unit

package numgo

import "testing"

func TestZeros(t *testing.T) {
	a := Zeros(3, 4)
	if a.Size() != 12 {
		t.Fatalf("expected size 12, got %d", a.Size())
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 4; j++ {
			if a.Get(i, j) != 0 {
				t.Fatalf("Zeros: non-zero at (%d,%d)", i, j)
			}
		}
	}
}

func TestOnes(t *testing.T) {
	a := Ones(2, 2)
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if a.Get(i, j) != 1 {
				t.Fatalf("Ones: expected 1 at (%d,%d), got %f", i, j, a.Get(i, j))
			}
		}
	}
}

func TestFull(t *testing.T) {
	a := Full(3.14, 2, 3)
	if a.Shape()[0] != 2 || a.Shape()[1] != 3 {
		t.Fatalf("Full: unexpected shape %v", a.Shape())
	}
	for i := 0; i < a.Size(); i++ {
		if a.Data()[i] != 3.14 {
			t.Fatalf("Full: expected 3.14, got %f", a.Data()[i])
		}
	}
}

func TestEye(t *testing.T) {
	a := Eye(3)
	if a.Shape()[0] != 3 || a.Shape()[1] != 3 {
		t.Fatalf("Eye: unexpected shape %v", a.Shape())
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			expected := 0.0
			if i == j {
				expected = 1.0
			}
			if a.Get(i, j) != expected {
				t.Fatalf("Eye: at (%d,%d) expected %f, got %f", i, j, expected, a.Get(i, j))
			}
		}
	}
}

func TestFromSlice(t *testing.T) {
	a := FromSlice([]float64{10, 20, 30})
	if a.Ndim() != 1 || a.Size() != 3 {
		t.Fatalf("FromSlice: unexpected shape %v", a.Shape())
	}
	if a.Get(1) != 20 {
		t.Fatal("FromSlice: wrong data")
	}
}

func TestFromSlice2D(t *testing.T) {
	a := FromSlice2D([][]float64{
		{1, 2, 3},
		{4, 5, 6},
	})
	if a.Ndim() != 2 {
		t.Fatalf("FromSlice2D: expected 2D, got %dD", a.Ndim())
	}
	if a.Shape()[0] != 2 || a.Shape()[1] != 3 {
		t.Fatalf("FromSlice2D: unexpected shape %v", a.Shape())
	}
	if a.Get(1, 2) != 6 {
		t.Fatalf("FromSlice2D: expected 6, got %f", a.Get(1, 2))
	}
}

func TestFromSlice2DEmpty(t *testing.T) {
	a := FromSlice2D([][]float64{})
	if a.Size() != 0 {
		t.Fatal("FromSlice2D empty should have size 0")
	}
}

func TestFromSlice2DPanicsJagged(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on jagged input")
		}
	}()
	FromSlice2D([][]float64{
		{1, 2},
		{3},
	})
}
