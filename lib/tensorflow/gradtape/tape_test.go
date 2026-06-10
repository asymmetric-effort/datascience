package gradtape

import (
	"math"
	"testing"
)

func TestVariableValue(t *testing.T) {
	tape := New()
	v := tape.Variable([]float64{1, 2, 3})
	got := v.Value()
	want := []float64{1, 2, 3}
	for i, val := range got {
		if val != want[i] {
			t.Errorf("Value()[%d] = %f, want %f", i, val, want[i])
		}
	}
}

func TestVariableIndex(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{1})
	b := tape.Variable([]float64{2})
	if a.Index() != 0 {
		t.Errorf("a.Index() = %d, want 0", a.Index())
	}
	if b.Index() != 1 {
		t.Errorf("b.Index() = %d, want 1", b.Index())
	}
}

func TestVariableDataIsolation(t *testing.T) {
	tape := New()
	data := []float64{1, 2, 3}
	v := tape.Variable(data)
	data[0] = 99
	if v.Value()[0] == 99 {
		t.Error("modifying input data affected variable")
	}
}

func TestAdd(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{1, 2, 3})
	b := tape.Variable([]float64{4, 5, 6})
	c, err := tape.Add(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{5, 7, 9}
	got := c.Value()
	for i, v := range got {
		if v != want[i] {
			t.Errorf("result[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestAddLengthMismatch(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{1, 2})
	b := tape.Variable([]float64{1, 2, 3})
	_, err := tape.Add(a, b)
	if err == nil {
		t.Error("expected error for length mismatch")
	}
}

func TestSub(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{10, 20})
	b := tape.Variable([]float64{3, 7})
	c, err := tape.Sub(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{7, 13}
	got := c.Value()
	for i, v := range got {
		if v != want[i] {
			t.Errorf("result[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestSubLengthMismatch(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{1})
	b := tape.Variable([]float64{1, 2})
	_, err := tape.Sub(a, b)
	if err == nil {
		t.Error("expected error for length mismatch")
	}
}

func TestMul(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{2, 3})
	b := tape.Variable([]float64{4, 5})
	c, err := tape.Mul(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{8, 15}
	got := c.Value()
	for i, v := range got {
		if v != want[i] {
			t.Errorf("result[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestMulLengthMismatch(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{1})
	b := tape.Variable([]float64{1, 2})
	_, err := tape.Mul(a, b)
	if err == nil {
		t.Error("expected error for length mismatch")
	}
}

func TestNeg(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{1, -2, 3})
	b, err := tape.Neg(a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{-1, 2, -3}
	got := b.Value()
	for i, v := range got {
		if v != want[i] {
			t.Errorf("result[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestSum(t *testing.T) {
	tape := New()
	a := tape.Variable([]float64{1, 2, 3, 4})
	s, err := tape.Sum(a)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := s.Value()
	if len(got) != 1 || got[0] != 10 {
		t.Errorf("Sum() = %v, want [10]", got)
	}
}

func TestGradientAdd(t *testing.T) {
	// f = a + b, df/da = 1, df/db = 1
	tape := New()
	a := tape.Variable([]float64{3, 4})
	b := tape.Variable([]float64{5, 6})
	c, _ := tape.Add(a, b)
	grads := tape.Gradient(c)

	for i, v := range grads[a.Index()] {
		if v != 1.0 {
			t.Errorf("grad_a[%d] = %f, want 1.0", i, v)
		}
	}
	for i, v := range grads[b.Index()] {
		if v != 1.0 {
			t.Errorf("grad_b[%d] = %f, want 1.0", i, v)
		}
	}
}

func TestGradientSub(t *testing.T) {
	// f = a - b, df/da = 1, df/db = -1
	tape := New()
	a := tape.Variable([]float64{10})
	b := tape.Variable([]float64{3})
	c, _ := tape.Sub(a, b)
	grads := tape.Gradient(c)

	if grads[a.Index()][0] != 1.0 {
		t.Errorf("grad_a = %f, want 1.0", grads[a.Index()][0])
	}
	if grads[b.Index()][0] != -1.0 {
		t.Errorf("grad_b = %f, want -1.0", grads[b.Index()][0])
	}
}

func TestGradientMul(t *testing.T) {
	// f = a * b, df/da = b, df/db = a
	tape := New()
	a := tape.Variable([]float64{3})
	b := tape.Variable([]float64{5})
	c, _ := tape.Mul(a, b)
	grads := tape.Gradient(c)

	if grads[a.Index()][0] != 5.0 {
		t.Errorf("grad_a = %f, want 5.0", grads[a.Index()][0])
	}
	if grads[b.Index()][0] != 3.0 {
		t.Errorf("grad_b = %f, want 3.0", grads[b.Index()][0])
	}
}

func TestGradientNeg(t *testing.T) {
	// f = -a, df/da = -1
	tape := New()
	a := tape.Variable([]float64{7})
	b, _ := tape.Neg(a)
	grads := tape.Gradient(b)

	if grads[a.Index()][0] != -1.0 {
		t.Errorf("grad_a = %f, want -1.0", grads[a.Index()][0])
	}
}

func TestGradientSum(t *testing.T) {
	// f = sum(a), df/da_i = 1 for all i
	tape := New()
	a := tape.Variable([]float64{1, 2, 3})
	s, _ := tape.Sum(a)
	grads := tape.Gradient(s)

	for i, v := range grads[a.Index()] {
		if v != 1.0 {
			t.Errorf("grad_a[%d] = %f, want 1.0", i, v)
		}
	}
}

func TestGradientChain(t *testing.T) {
	// f = sum((a * b) + a), df/da = b + 1, df/db = a
	tape := New()
	a := tape.Variable([]float64{2, 3})
	b := tape.Variable([]float64{4, 5})
	ab, _ := tape.Mul(a, b)
	abPlusA, _ := tape.Add(ab, a)
	loss, _ := tape.Sum(abPlusA)
	grads := tape.Gradient(loss)

	wantGradA := []float64{5, 6} // b + 1
	wantGradB := []float64{2, 3} // a
	for i, v := range grads[a.Index()] {
		if math.Abs(v-wantGradA[i]) > 1e-10 {
			t.Errorf("grad_a[%d] = %f, want %f", i, v, wantGradA[i])
		}
	}
	for i, v := range grads[b.Index()] {
		if math.Abs(v-wantGradB[i]) > 1e-10 {
			t.Errorf("grad_b[%d] = %f, want %f", i, v, wantGradB[i])
		}
	}
}

func TestGradientNoOp(t *testing.T) {
	// Gradient of a variable with no ops should just be the seed.
	tape := New()
	a := tape.Variable([]float64{5})
	grads := tape.Gradient(a)
	if grads[a.Index()][0] != 1.0 {
		t.Errorf("grad_a = %f, want 1.0", grads[a.Index()][0])
	}
}

func TestTapeCapacityExceeded(t *testing.T) {
	tape := New()
	// Fill up the tape to capacity.
	tape.ops = make([]operation, maxTapeOperations)
	a := tape.Variable([]float64{1})
	b := tape.Variable([]float64{2})

	_, err := tape.Add(a, b)
	if err == nil {
		t.Error("expected capacity error from Add")
	}
	_, err = tape.Sub(a, b)
	if err == nil {
		t.Error("expected capacity error from Sub")
	}
	_, err = tape.Mul(a, b)
	if err == nil {
		t.Error("expected capacity error from Mul")
	}
	_, err = tape.Neg(a)
	if err == nil {
		t.Error("expected capacity error from Neg")
	}
	_, err = tape.Sum(a)
	if err == nil {
		t.Error("expected capacity error from Sum")
	}
}

func TestGradientDisconnected(t *testing.T) {
	// Variable not connected to the output should have no gradient.
	tape := New()
	a := tape.Variable([]float64{1, 2})
	b := tape.Variable([]float64{3, 4})
	c := tape.Variable([]float64{5, 6})
	// Only a+b is connected to output; c is disconnected.
	out, _ := tape.Add(a, b)
	grads := tape.Gradient(out)
	if _, ok := grads[c.Index()]; ok {
		t.Error("disconnected variable should not have gradient")
	}
}
