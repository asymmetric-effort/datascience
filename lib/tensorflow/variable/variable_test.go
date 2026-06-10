package variable

import (
	"math"
	"reflect"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

func makeArray(t *testing.T, shape []int, data []float64) *numgo.NDArray {
	t.Helper()
	return numgo.NewNDArray(shape, data)
}

func TestNew(t *testing.T) {
	initial := makeArray(t, []int{3}, []float64{1, 2, 3})
	v := New(initial, "weights", true)
	if v.Name() != "weights" {
		t.Errorf("Name() = %q, want %q", v.Name(), "weights")
	}
	if !v.Trainable() {
		t.Error("expected trainable=true")
	}
	if v.ID() <= 0 {
		t.Error("expected positive ID")
	}
	if v.NumElements() != 3 {
		t.Errorf("NumElements() = %d, want 3", v.NumElements())
	}
}

func TestNewUniqueIDs(t *testing.T) {
	initial := makeArray(t, []int{1}, []float64{0})
	a := New(initial, "a", true)
	b := New(initial, "b", true)
	if a.ID() == b.ID() {
		t.Error("variables should have unique IDs")
	}
}

func TestValue(t *testing.T) {
	initial := makeArray(t, []int{2}, []float64{10, 20})
	v := New(initial, "v", true)
	val := v.Value()
	initialData := initial.Data()
	for i, got := range val.Data() {
		if got != initialData[i] {
			t.Errorf("Value()[%d] = %f, want %f", i, got, initialData[i])
		}
	}
}

func TestValueIsolation(t *testing.T) {
	initial := makeArray(t, []int{2}, []float64{1, 2})
	v := New(initial, "v", true)
	val := v.Value()
	// Modifying returned value should not affect variable.
	val.Set(99, 0)
	if v.Value().Get(0) == 99 {
		t.Error("modifying Value() affected internal state")
	}
}

func TestAssign(t *testing.T) {
	initial := makeArray(t, []int{2}, []float64{1, 2})
	v := New(initial, "v", true)
	newVal := makeArray(t, []int{2}, []float64{10, 20})
	err := v.Assign(newVal)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	newValData := newVal.Data()
	for i, got := range v.Value().Data() {
		if got != newValData[i] {
			t.Errorf("data[%d] = %f, want %f", i, got, newValData[i])
		}
	}
}

func TestAssignShapeMismatch(t *testing.T) {
	initial := makeArray(t, []int{2}, []float64{1, 2})
	v := New(initial, "v", true)
	bad := makeArray(t, []int{3}, []float64{1, 2, 3})
	err := v.Assign(bad)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestAssignAdd(t *testing.T) {
	initial := makeArray(t, []int{3}, []float64{1, 2, 3})
	v := New(initial, "v", true)
	delta := makeArray(t, []int{3}, []float64{10, 20, 30})
	err := v.AssignAdd(delta)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{11, 22, 33}
	for i, got := range v.Value().Data() {
		if math.Abs(got-want[i]) > 1e-10 {
			t.Errorf("data[%d] = %f, want %f", i, got, want[i])
		}
	}
}

func TestAssignAddShapeMismatch(t *testing.T) {
	initial := makeArray(t, []int{2}, []float64{1, 2})
	v := New(initial, "v", true)
	bad := makeArray(t, []int{3}, []float64{1, 2, 3})
	err := v.AssignAdd(bad)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestAssignSub(t *testing.T) {
	initial := makeArray(t, []int{3}, []float64{10, 20, 30})
	v := New(initial, "v", true)
	delta := makeArray(t, []int{3}, []float64{1, 2, 3})
	err := v.AssignSub(delta)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []float64{9, 18, 27}
	for i, got := range v.Value().Data() {
		if math.Abs(got-want[i]) > 1e-10 {
			t.Errorf("data[%d] = %f, want %f", i, got, want[i])
		}
	}
}

func TestAssignSubShapeMismatch(t *testing.T) {
	initial := makeArray(t, []int{2}, []float64{1, 2})
	v := New(initial, "v", true)
	bad := makeArray(t, []int{3}, []float64{1, 2, 3})
	err := v.AssignSub(bad)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestShape(t *testing.T) {
	initial := makeArray(t, []int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	v := New(initial, "v", false)
	expectedShape := []int{2, 3}
	if !reflect.DeepEqual(v.Shape(), expectedShape) {
		t.Errorf("Shape() = %v, want %v", v.Shape(), expectedShape)
	}
	if v.Trainable() {
		t.Error("expected trainable=false")
	}
}
