//go:build unit

package gpu

import (
	"math"
	"testing"
)

func TestNewTensorNilData(t *testing.T) {
	b := NewCPUBackend()
	tensor := NewTensor([]int{2, 3}, nil, b)
	if len(tensor.Data()) != 6 {
		t.Errorf("expected 6 elements, got %d", len(tensor.Data()))
	}
	for _, v := range tensor.Data() {
		if v != 0 {
			t.Errorf("expected zero-initialized, got %v", v)
			break
		}
	}
}

func TestNewTensorWithData(t *testing.T) {
	b := NewCPUBackend()
	data := []float64{1, 2, 3, 4, 5, 6}
	tensor := NewTensor([]int{2, 3}, data, b)
	result := tensor.Data()
	if !floatsEqual(result, data) {
		t.Errorf("expected %v, got %v", data, result)
	}
}

func TestTensorDataIsCopy(t *testing.T) {
	b := NewCPUBackend()
	data := []float64{1, 2, 3}
	tensor := NewTensor([]int{3}, data, b)
	d := tensor.Data()
	d[0] = 999
	if tensor.Data()[0] == 999 {
		t.Error("Data() should return a copy")
	}
}

func TestTensorShape(t *testing.T) {
	b := NewCPUBackend()
	tensor := NewTensor([]int{2, 3, 4}, nil, b)
	shape := tensor.Shape()
	if !intsEqual(shape, []int{2, 3, 4}) {
		t.Errorf("expected [2,3,4], got %v", shape)
	}
	// Mutating returned shape should not affect tensor.
	shape[0] = 99
	if tensor.Shape()[0] == 99 {
		t.Error("Shape() should return a copy")
	}
}

func TestTensorSize(t *testing.T) {
	b := NewCPUBackend()
	tensor := NewTensor([]int{2, 3, 4}, nil, b)
	if tensor.Size() != 24 {
		t.Errorf("expected size 24, got %d", tensor.Size())
	}
}

func TestTensorAdd(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{3}, []float64{1, 2, 3}, b)
	other := NewTensor([]int{3}, []float64{4, 5, 6}, b)
	result := a.Add(other)
	expected := []float64{5, 7, 9}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.Add: expected %v, got %v", expected, result.Data())
	}
}

func TestTensorSub(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{3}, []float64{10, 20, 30}, b)
	other := NewTensor([]int{3}, []float64{1, 2, 3}, b)
	result := a.Sub(other)
	expected := []float64{9, 18, 27}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.Sub: expected %v, got %v", expected, result.Data())
	}
}

func TestTensorMul(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{3}, []float64{2, 3, 4}, b)
	other := NewTensor([]int{3}, []float64{5, 6, 7}, b)
	result := a.Mul(other)
	expected := []float64{10, 18, 28}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.Mul: expected %v, got %v", expected, result.Data())
	}
}

func TestTensorMatMul(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{2, 2}, []float64{1, 2, 3, 4}, b)
	other := NewTensor([]int{2, 2}, []float64{5, 6, 7, 8}, b)
	result := a.MatMul(other)
	expected := []float64{19, 22, 43, 50}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.MatMul: expected %v, got %v", expected, result.Data())
	}
	if !intsEqual(result.Shape(), []int{2, 2}) {
		t.Errorf("Tensor.MatMul shape: expected [2,2], got %v", result.Shape())
	}
}

func TestTensorMatMulNonSquare(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6}, b)
	other := NewTensor([]int{3, 1}, []float64{1, 2, 3}, b)
	result := a.MatMul(other)
	expected := []float64{14, 32}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.MatMul non-square: expected %v, got %v", expected, result.Data())
	}
	if !intsEqual(result.Shape(), []int{2, 1}) {
		t.Errorf("shape: expected [2,1], got %v", result.Shape())
	}
}

func TestTensorScalarMul(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{3}, []float64{1, 2, 3}, b)
	result := a.ScalarMul(2.5)
	expected := []float64{2.5, 5.0, 7.5}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.ScalarMul: expected %v, got %v", expected, result.Data())
	}
}

func TestTensorSum(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6}, b)
	result := a.Sum(0)
	expected := []float64{5, 7, 9}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.Sum axis 0: expected %v, got %v", expected, result.Data())
	}
	if !intsEqual(result.Shape(), []int{3}) {
		t.Errorf("Tensor.Sum shape: expected [3], got %v", result.Shape())
	}
}

func TestTensorMax(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{2, 3}, []float64{1, 5, 3, 4, 2, 6}, b)
	result := a.Max(0)
	expected := []float64{4, 5, 6}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.Max axis 0: expected %v, got %v", expected, result.Data())
	}
}

func TestTensorReshape(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6}, b)
	result := a.Reshape([]int{3, 2})
	if !intsEqual(result.Shape(), []int{3, 2}) {
		t.Errorf("Reshape shape: expected [3,2], got %v", result.Shape())
	}
	if !floatsEqual(result.Data(), []float64{1, 2, 3, 4, 5, 6}) {
		t.Errorf("Reshape data should be unchanged")
	}
}

func TestTensorClone(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{3}, []float64{1, 2, 3}, b)
	clone := a.Clone()
	if !floatsEqual(a.Data(), clone.Data()) {
		t.Errorf("Clone data mismatch")
	}
	if !intsEqual(a.Shape(), clone.Shape()) {
		t.Errorf("Clone shape mismatch")
	}
}

func TestTensorNormalize(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{4}, []float64{1, 2, 3, 4}, b)
	result := a.Normalize()
	expected := []float64{0.1, 0.2, 0.3, 0.4}
	if !floatsEqual(result.Data(), expected) {
		t.Errorf("Tensor.Normalize: expected %v, got %v", expected, result.Data())
	}
}

func TestTensorExp(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{2}, []float64{0, 1}, b)
	result := a.Exp()
	expected := []float64{1.0, math.E}
	if !floatsNear(result.Data(), expected, 1e-10) {
		t.Errorf("Tensor.Exp: expected %v, got %v", expected, result.Data())
	}
}

func TestTensorLog(t *testing.T) {
	b := NewCPUBackend()
	a := NewTensor([]int{2}, []float64{1, math.E}, b)
	result := a.Log()
	expected := []float64{0, 1}
	if !floatsNear(result.Data(), expected, 1e-10) {
		t.Errorf("Tensor.Log: expected %v, got %v", expected, result.Data())
	}
}

func TestTensorToDevice(t *testing.T) {
	b1 := NewCPUBackend()
	b2 := NewCPUBackend()
	a := NewTensor([]int{3}, []float64{1, 2, 3}, b1)
	moved := a.ToDevice(b2)
	if !floatsEqual(a.Data(), moved.Data()) {
		t.Errorf("ToDevice data mismatch")
	}
}
