//go:build unit

package numgo

import (
	"math"
	"testing"
)

const tol = 1e-12

func assertClose(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Fatalf("%s: expected %v, got %v", name, want, got)
	}
}

func assertSliceClose(t *testing.T, name string, got []float64, want []float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: length mismatch %d vs %d", name, len(got), len(want))
	}
	for i := range got {
		if math.Abs(got[i]-want[i]) > tol {
			t.Fatalf("%s: at index %d expected %v, got %v", name, i, want[i], got[i])
		}
	}
}

func TestSin(t *testing.T) {
	a := FromSlice([]float64{0, math.Pi / 2, math.Pi})
	r := Sin(a).Data()
	assertSliceClose(t, "Sin", r, []float64{0, 1, 0})
}

func TestCos(t *testing.T) {
	a := FromSlice([]float64{0, math.Pi / 2, math.Pi})
	r := Cos(a).Data()
	assertSliceClose(t, "Cos", r, []float64{1, 0, -1})
}

func TestTan(t *testing.T) {
	a := FromSlice([]float64{0, math.Pi / 4})
	r := Tan(a).Data()
	assertSliceClose(t, "Tan", r, []float64{0, 1})
}

func TestArcsin(t *testing.T) {
	a := FromSlice([]float64{0, 1})
	r := Arcsin(a).Data()
	assertSliceClose(t, "Arcsin", r, []float64{0, math.Pi / 2})
}

func TestArccos(t *testing.T) {
	a := FromSlice([]float64{1, 0})
	r := Arccos(a).Data()
	assertSliceClose(t, "Arccos", r, []float64{0, math.Pi / 2})
}

func TestArctan(t *testing.T) {
	a := FromSlice([]float64{0, 1})
	r := Arctan(a).Data()
	assertSliceClose(t, "Arctan", r, []float64{0, math.Pi / 4})
}

func TestArctan2(t *testing.T) {
	y := FromSlice([]float64{1, 0, -1})
	x := FromSlice([]float64{0, 1, 0})
	r := Arctan2(y, x).Data()
	assertSliceClose(t, "Arctan2", r, []float64{math.Pi / 2, 0, -math.Pi / 2})
}

func TestArctan2Broadcast(t *testing.T) {
	y := NewNDArray([]int{3, 1}, []float64{1, 0, -1})
	x := NewNDArray([]int{1, 2}, []float64{1, -1})
	r := Arctan2(y, x)
	if r.Shape()[0] != 3 || r.Shape()[1] != 2 {
		t.Fatalf("Arctan2 broadcast shape: expected [3 2], got %v", r.Shape())
	}
	assertClose(t, "Arctan2 broadcast [0,0]", r.Get(0, 0), math.Atan2(1, 1))
	assertClose(t, "Arctan2 broadcast [0,1]", r.Get(0, 1), math.Atan2(1, -1))
}

func TestSinh(t *testing.T) {
	a := FromSlice([]float64{0, 1})
	r := Sinh(a).Data()
	assertSliceClose(t, "Sinh", r, []float64{0, math.Sinh(1)})
}

func TestCosh(t *testing.T) {
	a := FromSlice([]float64{0, 1})
	r := Cosh(a).Data()
	assertSliceClose(t, "Cosh", r, []float64{1, math.Cosh(1)})
}

func TestTanh(t *testing.T) {
	a := FromSlice([]float64{0, 1})
	r := Tanh(a).Data()
	assertSliceClose(t, "Tanh", r, []float64{0, math.Tanh(1)})
}

func TestExp(t *testing.T) {
	a := FromSlice([]float64{0, 1, 2})
	r := Exp(a).Data()
	assertSliceClose(t, "Exp", r, []float64{1, math.E, math.E * math.E})
}

func TestExp2(t *testing.T) {
	a := FromSlice([]float64{0, 1, 3})
	r := Exp2(a).Data()
	assertSliceClose(t, "Exp2", r, []float64{1, 2, 8})
}

func TestExpm1(t *testing.T) {
	a := FromSlice([]float64{0, 1})
	r := Expm1(a).Data()
	assertSliceClose(t, "Expm1", r, []float64{0, math.E - 1})
}

func TestLog(t *testing.T) {
	a := FromSlice([]float64{1, math.E})
	r := Log(a).Data()
	assertSliceClose(t, "Log", r, []float64{0, 1})
}

func TestLog2(t *testing.T) {
	a := FromSlice([]float64{1, 2, 8})
	r := Log2(a).Data()
	assertSliceClose(t, "Log2", r, []float64{0, 1, 3})
}

func TestLog10(t *testing.T) {
	a := FromSlice([]float64{1, 10, 100})
	r := Log10(a).Data()
	assertSliceClose(t, "Log10", r, []float64{0, 1, 2})
}

func TestLog1p(t *testing.T) {
	a := FromSlice([]float64{0, math.E - 1})
	r := Log1p(a).Data()
	assertSliceClose(t, "Log1p", r, []float64{0, 1})
}

func TestPower(t *testing.T) {
	base := FromSlice([]float64{2, 3, 4})
	exp := FromSlice([]float64{3, 2, 0.5})
	r := Power(base, exp).Data()
	assertSliceClose(t, "Power", r, []float64{8, 9, 2})
}

func TestPowerBroadcast(t *testing.T) {
	base := NewNDArray([]int{2, 1}, []float64{2, 3})
	exp := NewNDArray([]int{1, 3}, []float64{1, 2, 3})
	r := Power(base, exp)
	if r.Shape()[0] != 2 || r.Shape()[1] != 3 {
		t.Fatalf("Power broadcast shape: expected [2 3], got %v", r.Shape())
	}
	assertClose(t, "Power broadcast [0,0]", r.Get(0, 0), 2)
	assertClose(t, "Power broadcast [0,2]", r.Get(0, 2), 8)
	assertClose(t, "Power broadcast [1,1]", r.Get(1, 1), 9)
}

func TestSqrt(t *testing.T) {
	a := FromSlice([]float64{0, 1, 4, 9})
	r := Sqrt(a).Data()
	assertSliceClose(t, "Sqrt", r, []float64{0, 1, 2, 3})
}

func TestCbrt(t *testing.T) {
	a := FromSlice([]float64{0, 1, 8, 27})
	r := Cbrt(a).Data()
	assertSliceClose(t, "Cbrt", r, []float64{0, 1, 2, 3})
}

func TestSquare(t *testing.T) {
	a := FromSlice([]float64{-3, 0, 4})
	r := Square(a).Data()
	assertSliceClose(t, "Square", r, []float64{9, 0, 16})
}

func TestAbsolute(t *testing.T) {
	a := FromSlice([]float64{-3, 0, 4})
	r := Absolute(a).Data()
	assertSliceClose(t, "Absolute", r, []float64{3, 0, 4})
}

func TestFabs(t *testing.T) {
	a := FromSlice([]float64{-2.5, 0, 3.5})
	r := Fabs(a).Data()
	assertSliceClose(t, "Fabs", r, []float64{2.5, 0, 3.5})
}

func TestSign(t *testing.T) {
	a := FromSlice([]float64{-5, 0, 7})
	r := Sign(a).Data()
	assertSliceClose(t, "Sign", r, []float64{-1, 0, 1})
}

func TestHeaviside(t *testing.T) {
	x := FromSlice([]float64{-1, 0, 1})
	h0 := FromSlice([]float64{0.5, 0.5, 0.5})
	r := Heaviside(x, h0).Data()
	assertSliceClose(t, "Heaviside", r, []float64{0, 0.5, 1})
}

func TestHeavisideBroadcast(t *testing.T) {
	x := NewNDArray([]int{3, 1}, []float64{-1, 0, 1})
	h0 := NewNDArray([]int{1, 1}, []float64{0.5})
	r := Heaviside(x, h0)
	if r.Shape()[0] != 3 || r.Shape()[1] != 1 {
		t.Fatalf("Heaviside broadcast shape: expected [3 1], got %v", r.Shape())
	}
	assertClose(t, "Heaviside broadcast [0,0]", r.Get(0, 0), 0)
	assertClose(t, "Heaviside broadcast [1,0]", r.Get(1, 0), 0.5)
	assertClose(t, "Heaviside broadcast [2,0]", r.Get(2, 0), 1)
}

func TestFmod(t *testing.T) {
	x := FromSlice([]float64{5, -5, 7.5})
	y := FromSlice([]float64{3, 3, 2.5})
	r := Fmod(x, y).Data()
	assertSliceClose(t, "Fmod", r, []float64{2, -2, 0})
}

func TestFmodBroadcast(t *testing.T) {
	x := NewNDArray([]int{2, 1}, []float64{7, 10})
	y := NewNDArray([]int{1, 2}, []float64{3, 4})
	r := Fmod(x, y)
	if r.Shape()[0] != 2 || r.Shape()[1] != 2 {
		t.Fatalf("Fmod broadcast shape: expected [2 2], got %v", r.Shape())
	}
	assertClose(t, "Fmod broadcast [0,0]", r.Get(0, 0), math.Mod(7, 3))
	assertClose(t, "Fmod broadcast [1,1]", r.Get(1, 1), math.Mod(10, 4))
}

func TestModf(t *testing.T) {
	a := FromSlice([]float64{1.5, -2.7, 3.0})
	frac, integer := Modf(a)
	fracData := frac.Data()
	intData := integer.Data()
	for i, v := range a.Data() {
		wantInt, wantFrac := math.Modf(v)
		assertClose(t, "Modf frac", fracData[i], wantFrac)
		assertClose(t, "Modf integer", intData[i], wantInt)
	}
}

func TestRemainder(t *testing.T) {
	x := FromSlice([]float64{5, -5, 7})
	y := FromSlice([]float64{3, 3, 4})
	r := Remainder(x, y).Data()
	assertSliceClose(t, "Remainder", r, []float64{
		math.Remainder(5, 3),
		math.Remainder(-5, 3),
		math.Remainder(7, 4),
	})
}

func TestRemainderBroadcast(t *testing.T) {
	x := NewNDArray([]int{2, 1}, []float64{5, 10})
	y := NewNDArray([]int{1, 2}, []float64{3, 4})
	r := Remainder(x, y)
	if r.Shape()[0] != 2 || r.Shape()[1] != 2 {
		t.Fatalf("Remainder broadcast shape: expected [2 2], got %v", r.Shape())
	}
	assertClose(t, "Remainder broadcast [0,0]", r.Get(0, 0), math.Remainder(5, 3))
}

func TestClip(t *testing.T) {
	a := FromSlice([]float64{-3, 0, 5, 10})
	r := Clip(a, 0, 7).Data()
	assertSliceClose(t, "Clip", r, []float64{0, 0, 5, 7})
}

func TestAround(t *testing.T) {
	a := FromSlice([]float64{1.234, 2.555, -3.678})
	r := Around(a, 2).Data()
	assertSliceClose(t, "Around", r, []float64{1.23, 2.56, -3.68})
}

func TestAroundZeroDecimals(t *testing.T) {
	a := FromSlice([]float64{1.5, 2.5, 3.5, 4.5})
	r := Around(a, 0).Data()
	// Banker's rounding: .5 rounds to even
	assertSliceClose(t, "Around(0)", r, []float64{2, 2, 4, 4})
}

func TestRint(t *testing.T) {
	a := FromSlice([]float64{1.5, 2.5, 3.1, -1.7})
	r := Rint(a).Data()
	assertSliceClose(t, "Rint", r, []float64{2, 2, 3, -2})
}

func TestFloor(t *testing.T) {
	a := FromSlice([]float64{1.7, -1.7, 0})
	r := Floor(a).Data()
	assertSliceClose(t, "Floor", r, []float64{1, -2, 0})
}

func TestCeil(t *testing.T) {
	a := FromSlice([]float64{1.1, -1.7, 0})
	r := Ceil(a).Data()
	assertSliceClose(t, "Ceil", r, []float64{2, -1, 0})
}
