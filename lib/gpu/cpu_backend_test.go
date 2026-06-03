//go:build unit

package gpu

import (
	"math"
	"testing"
)

const epsilon = 1e-12

func floatsEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > epsilon {
			return false
		}
	}
	return true
}

func floatsNear(a, b []float64, tol float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if math.Abs(a[i]-b[i]) > tol {
			return false
		}
	}
	return true
}

func intsEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// --- Original tests ---

func TestCPUBackendName(t *testing.T) {
	b := NewCPUBackend()
	if b.Name() != "cpu" {
		t.Errorf("expected 'cpu', got %q", b.Name())
	}
}

func TestCPUBackendIsAvailable(t *testing.T) {
	b := NewCPUBackend()
	if !b.IsAvailable() {
		t.Error("expected CPU backend to always be available")
	}
}

func TestCPUBackendClose(t *testing.T) {
	b := NewCPUBackend()
	if err := b.Close(); err != nil {
		t.Errorf("expected nil error from Close, got %v", err)
	}
}

func TestCPUBackendMatMul(t *testing.T) {
	b := NewCPUBackend()
	a := []float64{1, 2, 3, 4}
	bm := []float64{5, 6, 7, 8}
	result := b.MatMul(a, bm, 2, 2, 2)
	expected := []float64{19, 22, 43, 50}
	if !floatsEqual(result, expected) {
		t.Errorf("MatMul: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendMatMulNonSquare(t *testing.T) {
	b := NewCPUBackend()
	a := []float64{1, 2, 3, 4, 5, 6}
	bm := []float64{7, 8, 9, 10, 11, 12}
	result := b.MatMul(a, bm, 2, 3, 2)
	expected := []float64{58, 64, 139, 154}
	if !floatsEqual(result, expected) {
		t.Errorf("MatMul non-square: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendElementWiseMul(t *testing.T) {
	b := NewCPUBackend()
	a := []float64{1, 2, 3, 4}
	bv := []float64{5, 6, 7, 8}
	result := b.ElementWiseMul(a, bv)
	expected := []float64{5, 12, 21, 32}
	if !floatsEqual(result, expected) {
		t.Errorf("ElementWiseMul: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendElementWiseMulEmpty(t *testing.T) {
	b := NewCPUBackend()
	result := b.ElementWiseMul([]float64{}, []float64{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}

func TestCPUBackendSum(t *testing.T) {
	b := NewCPUBackend()
	result := b.Sum([]float64{1, 2, 3, 4})
	if math.Abs(result-10) > epsilon {
		t.Errorf("Sum: expected 10, got %v", result)
	}
}

func TestCPUBackendSumEmpty(t *testing.T) {
	b := NewCPUBackend()
	result := b.Sum([]float64{})
	if result != 0 {
		t.Errorf("Sum empty: expected 0, got %v", result)
	}
}

func TestCPUBackendNormalize(t *testing.T) {
	b := NewCPUBackend()
	result := b.Normalize([]float64{1, 2, 3, 4})
	expected := []float64{0.1, 0.2, 0.3, 0.4}
	if !floatsEqual(result, expected) {
		t.Errorf("Normalize: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendNormalizeZeroSum(t *testing.T) {
	b := NewCPUBackend()
	result := b.Normalize([]float64{0, 0, 0})
	expected := []float64{0, 0, 0}
	if !floatsEqual(result, expected) {
		t.Errorf("Normalize zero sum: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendNormalizeSumsToOne(t *testing.T) {
	b := NewCPUBackend()
	result := b.Normalize([]float64{3, 7, 5, 15})
	sum := b.Sum(result)
	if math.Abs(sum-1.0) > epsilon {
		t.Errorf("Normalize result should sum to 1, got %v", sum)
	}
}

func TestCPUBackendFactorProduct(t *testing.T) {
	b := NewCPUBackend()
	aValues := []float64{0.4, 0.6}
	aShape := []int{2}
	bValues := []float64{0.1, 0.2, 0.7}
	bShape := []int{3}
	resultShape := []int{2, 3}
	result := b.FactorProduct(aValues, aShape, bValues, bShape, resultShape)
	expected := []float64{0.04, 0.08, 0.28, 0.06, 0.12, 0.42}
	if !floatsEqual(result, expected) {
		t.Errorf("FactorProduct: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendFactorProductScalar(t *testing.T) {
	b := NewCPUBackend()
	aValues := []float64{2.0}
	aShape := []int{1}
	bValues := []float64{3.0, 4.0}
	bShape := []int{2}
	resultShape := []int{1, 2}
	result := b.FactorProduct(aValues, aShape, bValues, bShape, resultShape)
	expected := []float64{6.0, 8.0}
	if !floatsEqual(result, expected) {
		t.Errorf("FactorProduct scalar: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendMarginalize(t *testing.T) {
	b := NewCPUBackend()
	values := []float64{1, 2, 3, 4, 5, 6}
	shape := []int{2, 3}
	result, newShape := b.Marginalize(values, shape, 0)
	expectedVals := []float64{5, 7, 9}
	expectedShape := []int{3}
	if !floatsEqual(result, expectedVals) {
		t.Errorf("Marginalize axis 0 values: expected %v, got %v", expectedVals, result)
	}
	if !intsEqual(newShape, expectedShape) {
		t.Errorf("Marginalize axis 0 shape: expected %v, got %v", expectedShape, newShape)
	}
}

func TestCPUBackendMarginalizeAxis1(t *testing.T) {
	b := NewCPUBackend()
	values := []float64{1, 2, 3, 4, 5, 6}
	shape := []int{2, 3}
	result, newShape := b.Marginalize(values, shape, 1)
	expectedVals := []float64{6, 15}
	expectedShape := []int{2}
	if !floatsEqual(result, expectedVals) {
		t.Errorf("Marginalize axis 1 values: expected %v, got %v", expectedVals, result)
	}
	if !intsEqual(newShape, expectedShape) {
		t.Errorf("Marginalize axis 1 shape: expected %v, got %v", expectedShape, newShape)
	}
}

func TestCPUBackendMarginalize3D(t *testing.T) {
	b := NewCPUBackend()
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	shape := []int{2, 2, 2}
	result, newShape := b.Marginalize(values, shape, 1)
	expectedVals := []float64{4, 6, 12, 14}
	expectedShape := []int{2, 2}
	if !floatsEqual(result, expectedVals) {
		t.Errorf("Marginalize 3D axis 1 values: expected %v, got %v", expectedVals, result)
	}
	if !intsEqual(newShape, expectedShape) {
		t.Errorf("Marginalize 3D axis 1 shape: expected %v, got %v", expectedShape, newShape)
	}
}

func TestCPUBackendImplementsInterface(t *testing.T) {
	var _ Backend = (*CPUBackend)(nil)
}

// --- New tensor operation tests ---

func TestCPUBackendElementWiseAdd(t *testing.T) {
	b := NewCPUBackend()
	result := b.ElementWiseAdd([]float64{1, 2, 3}, []float64{4, 5, 6})
	expected := []float64{5, 7, 9}
	if !floatsEqual(result, expected) {
		t.Errorf("ElementWiseAdd: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendElementWiseAddEmpty(t *testing.T) {
	b := NewCPUBackend()
	result := b.ElementWiseAdd([]float64{}, []float64{})
	if len(result) != 0 {
		t.Errorf("expected empty, got %v", result)
	}
}

func TestCPUBackendElementWiseSub(t *testing.T) {
	b := NewCPUBackend()
	result := b.ElementWiseSub([]float64{10, 20, 30}, []float64{1, 2, 3})
	expected := []float64{9, 18, 27}
	if !floatsEqual(result, expected) {
		t.Errorf("ElementWiseSub: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendElementWiseDiv(t *testing.T) {
	b := NewCPUBackend()
	result := b.ElementWiseDiv([]float64{10, 20, 30}, []float64{2, 5, 10})
	expected := []float64{5, 4, 3}
	if !floatsEqual(result, expected) {
		t.Errorf("ElementWiseDiv: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendElementWiseDivByZero(t *testing.T) {
	b := NewCPUBackend()
	result := b.ElementWiseDiv([]float64{1}, []float64{0})
	if !math.IsInf(result[0], 1) {
		t.Errorf("expected +Inf, got %v", result[0])
	}
}

func TestCPUBackendScalarMul(t *testing.T) {
	b := NewCPUBackend()
	result := b.ScalarMul([]float64{1, 2, 3}, 3.0)
	expected := []float64{3, 6, 9}
	if !floatsEqual(result, expected) {
		t.Errorf("ScalarMul: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendScalarAdd(t *testing.T) {
	b := NewCPUBackend()
	result := b.ScalarAdd([]float64{1, 2, 3}, 10.0)
	expected := []float64{11, 12, 13}
	if !floatsEqual(result, expected) {
		t.Errorf("ScalarAdd: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendExp(t *testing.T) {
	b := NewCPUBackend()
	result := b.Exp([]float64{0, 1, 2})
	expected := []float64{1.0, math.E, math.E * math.E}
	if !floatsNear(result, expected, 1e-10) {
		t.Errorf("Exp: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendLog(t *testing.T) {
	b := NewCPUBackend()
	result := b.Log([]float64{1, math.E, math.E * math.E})
	expected := []float64{0, 1, 2}
	if !floatsNear(result, expected, 1e-10) {
		t.Errorf("Log: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendSqrt(t *testing.T) {
	b := NewCPUBackend()
	result := b.Sqrt([]float64{0, 1, 4, 9, 16})
	expected := []float64{0, 1, 2, 3, 4}
	if !floatsEqual(result, expected) {
		t.Errorf("Sqrt: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendAbs(t *testing.T) {
	b := NewCPUBackend()
	result := b.Abs([]float64{-3, -1, 0, 1, 3})
	expected := []float64{3, 1, 0, 1, 3}
	if !floatsEqual(result, expected) {
		t.Errorf("Abs: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendMax(t *testing.T) {
	b := NewCPUBackend()
	result := b.Max([]float64{3, 1, 4, 1, 5, 9, 2, 6})
	if result != 9 {
		t.Errorf("Max: expected 9, got %v", result)
	}
}

func TestCPUBackendMaxSingleElement(t *testing.T) {
	b := NewCPUBackend()
	result := b.Max([]float64{42})
	if result != 42 {
		t.Errorf("Max single: expected 42, got %v", result)
	}
}

func TestCPUBackendMin(t *testing.T) {
	b := NewCPUBackend()
	result := b.Min([]float64{3, 1, 4, 1, 5, 9, 2, 6})
	if result != 1 {
		t.Errorf("Min: expected 1, got %v", result)
	}
}

func TestCPUBackendMinNegative(t *testing.T) {
	b := NewCPUBackend()
	result := b.Min([]float64{-5, -1, 0, 3})
	if result != -5 {
		t.Errorf("Min negative: expected -5, got %v", result)
	}
}

func TestCPUBackendArgMax(t *testing.T) {
	b := NewCPUBackend()
	result := b.ArgMax([]float64{1, 3, 2, 5, 4})
	if result != 3 {
		t.Errorf("ArgMax: expected 3, got %v", result)
	}
}

func TestCPUBackendArgMaxFirst(t *testing.T) {
	b := NewCPUBackend()
	result := b.ArgMax([]float64{9, 1, 2})
	if result != 0 {
		t.Errorf("ArgMax first: expected 0, got %v", result)
	}
}

func TestCPUBackendArgMin(t *testing.T) {
	b := NewCPUBackend()
	result := b.ArgMin([]float64{5, 3, 1, 2, 4})
	if result != 2 {
		t.Errorf("ArgMin: expected 2, got %v", result)
	}
}

func TestCPUBackendDot(t *testing.T) {
	b := NewCPUBackend()
	result := b.Dot([]float64{1, 2, 3}, []float64{4, 5, 6})
	expected := 32.0 // 1*4 + 2*5 + 3*6
	if math.Abs(result-expected) > epsilon {
		t.Errorf("Dot: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendDotOrthogonal(t *testing.T) {
	b := NewCPUBackend()
	result := b.Dot([]float64{1, 0}, []float64{0, 1})
	if result != 0 {
		t.Errorf("Dot orthogonal: expected 0, got %v", result)
	}
}

// --- Factor operations tests ---

func TestCPUBackendFactorReduce(t *testing.T) {
	b := NewCPUBackend()
	// shape [2, 3]: [[1,2,3],[4,5,6]]
	// Reduce axis 0, index 1 => [4,5,6]
	values := []float64{1, 2, 3, 4, 5, 6}
	result, newShape := b.FactorReduce(values, []int{2, 3}, 0, 1)
	expected := []float64{4, 5, 6}
	expectedShape := []int{3}
	if !floatsEqual(result, expected) {
		t.Errorf("FactorReduce axis 0 idx 1: expected %v, got %v", expected, result)
	}
	if !intsEqual(newShape, expectedShape) {
		t.Errorf("FactorReduce shape: expected %v, got %v", expectedShape, newShape)
	}
}

func TestCPUBackendFactorReduceAxis1(t *testing.T) {
	b := NewCPUBackend()
	// shape [2, 3]: [[1,2,3],[4,5,6]]
	// Reduce axis 1, index 2 => [3, 6]
	values := []float64{1, 2, 3, 4, 5, 6}
	result, newShape := b.FactorReduce(values, []int{2, 3}, 1, 2)
	expected := []float64{3, 6}
	expectedShape := []int{2}
	if !floatsEqual(result, expected) {
		t.Errorf("FactorReduce axis 1 idx 2: expected %v, got %v", expected, result)
	}
	if !intsEqual(newShape, expectedShape) {
		t.Errorf("FactorReduce shape: expected %v, got %v", expectedShape, newShape)
	}
}

func TestCPUBackendFactorReduce3D(t *testing.T) {
	b := NewCPUBackend()
	// shape [2,2,2]: [[[1,2],[3,4]],[[5,6],[7,8]]]
	// Reduce axis 1, index 0 => [[1,2],[5,6]]
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	result, newShape := b.FactorReduce(values, []int{2, 2, 2}, 1, 0)
	expected := []float64{1, 2, 5, 6}
	expectedShape := []int{2, 2}
	if !floatsEqual(result, expected) {
		t.Errorf("FactorReduce 3D: expected %v, got %v", expected, result)
	}
	if !intsEqual(newShape, expectedShape) {
		t.Errorf("FactorReduce 3D shape: expected %v, got %v", expectedShape, newShape)
	}
}

func TestCPUBackendFactorMaximize(t *testing.T) {
	b := NewCPUBackend()
	// shape [2, 3]: [[1,2,3],[4,5,6]]
	// Maximize axis 0 => [4,5,6]
	values := []float64{1, 2, 3, 4, 5, 6}
	result, newShape := b.FactorMaximize(values, []int{2, 3}, 0)
	expected := []float64{4, 5, 6}
	expectedShape := []int{3}
	if !floatsEqual(result, expected) {
		t.Errorf("FactorMaximize axis 0: expected %v, got %v", expected, result)
	}
	if !intsEqual(newShape, expectedShape) {
		t.Errorf("FactorMaximize shape: expected %v, got %v", expectedShape, newShape)
	}
}

func TestCPUBackendFactorMaximizeAxis1(t *testing.T) {
	b := NewCPUBackend()
	// shape [2, 3]: [[1,2,3],[4,5,6]]
	// Maximize axis 1 => [3, 6]
	values := []float64{1, 2, 3, 4, 5, 6}
	result, newShape := b.FactorMaximize(values, []int{2, 3}, 1)
	expected := []float64{3, 6}
	expectedShape := []int{2}
	if !floatsEqual(result, expected) {
		t.Errorf("FactorMaximize axis 1: expected %v, got %v", expected, result)
	}
	if !intsEqual(newShape, expectedShape) {
		t.Errorf("FactorMaximize shape: expected %v, got %v", expectedShape, newShape)
	}
}

func TestCPUBackendLogSumExp(t *testing.T) {
	b := NewCPUBackend()
	// log(exp(1) + exp(2) + exp(3))
	result := b.LogSumExp([]float64{1, 2, 3})
	// = 3 + log(exp(-2) + exp(-1) + 1) = 3 + log(1 + e^-1 + e^-2)
	expected := 3 + math.Log(1+math.Exp(-1)+math.Exp(-2))
	if math.Abs(result-expected) > 1e-10 {
		t.Errorf("LogSumExp: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendLogSumExpLargeValues(t *testing.T) {
	b := NewCPUBackend()
	// Numerical stability test: large values should not overflow.
	result := b.LogSumExp([]float64{1000, 1001, 1002})
	expected := 1002 + math.Log(1+math.Exp(-1)+math.Exp(-2))
	if math.Abs(result-expected) > 1e-6 {
		t.Errorf("LogSumExp large: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendSoftmax(t *testing.T) {
	b := NewCPUBackend()
	result := b.Softmax([]float64{1, 2, 3})
	// Should sum to 1.
	sum := b.Sum(result)
	if math.Abs(sum-1.0) > 1e-10 {
		t.Errorf("Softmax should sum to 1, got %v", sum)
	}
	// Monotonically increasing input => monotonically increasing output.
	if result[0] >= result[1] || result[1] >= result[2] {
		t.Errorf("Softmax should be monotonically increasing, got %v", result)
	}
}

func TestCPUBackendSoftmaxUniform(t *testing.T) {
	b := NewCPUBackend()
	result := b.Softmax([]float64{0, 0, 0})
	expected := []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}
	if !floatsNear(result, expected, 1e-10) {
		t.Errorf("Softmax uniform: expected %v, got %v", expected, result)
	}
}

// --- Batch operations tests ---

func TestCPUBackendBatchMatMul(t *testing.T) {
	b := NewCPUBackend()
	// Batch of 2 matrices: each 2x2 @ 2x2.
	a := []float64{
		1, 2, 3, 4, // batch 0
		5, 6, 7, 8, // batch 1
	}
	bm := []float64{
		1, 0, 0, 1, // batch 0: identity
		2, 0, 0, 2, // batch 1: scale by 2
	}
	result := b.BatchMatMul(a, bm, 2, 2, 2, 2)
	expected := []float64{
		1, 2, 3, 4, // batch 0: same as input
		10, 12, 14, 16, // batch 1: scaled by 2
	}
	if !floatsEqual(result, expected) {
		t.Errorf("BatchMatMul: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendBatchMatMulSingle(t *testing.T) {
	b := NewCPUBackend()
	a := []float64{1, 2, 3, 4}
	bm := []float64{5, 6, 7, 8}
	result := b.BatchMatMul(a, bm, 1, 2, 2, 2)
	expected := b.MatMul(a, bm, 2, 2, 2)
	if !floatsEqual(result, expected) {
		t.Errorf("BatchMatMul single: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendBatchNormalize(t *testing.T) {
	b := NewCPUBackend()
	a := []float64{1, 2, 3, 4, 10, 20, 30, 40}
	result := b.BatchNormalize(a, 2, 4)
	// Batch 0: sum=10, so [0.1, 0.2, 0.3, 0.4]
	// Batch 1: sum=100, so [0.1, 0.2, 0.3, 0.4]
	expected := []float64{0.1, 0.2, 0.3, 0.4, 0.1, 0.2, 0.3, 0.4}
	if !floatsNear(result, expected, 1e-10) {
		t.Errorf("BatchNormalize: expected %v, got %v", expected, result)
	}
}

func TestCPUBackendBatchNormalizeZeroSum(t *testing.T) {
	b := NewCPUBackend()
	a := []float64{0, 0, 0, 1, 2, 3}
	result := b.BatchNormalize(a, 2, 3)
	// Batch 0: zero sum => all zeros.
	// Batch 1: sum=6.
	if result[0] != 0 || result[1] != 0 || result[2] != 0 {
		t.Errorf("BatchNormalize zero batch: expected zeros, got %v", result[:3])
	}
	if math.Abs(b.Sum(result[3:])-1.0) > 1e-10 {
		t.Errorf("BatchNormalize non-zero batch should sum to 1")
	}
}

// --- Memory management tests ---

func TestCPUBackendAlloc(t *testing.T) {
	b := NewCPUBackend()
	data := b.Alloc(100)
	if len(data) != 100 {
		t.Errorf("Alloc: expected len 100, got %d", len(data))
	}
	for i, v := range data {
		if v != 0 {
			t.Errorf("Alloc: expected zero at index %d, got %v", i, v)
			break
		}
	}
}

func TestCPUBackendFree(t *testing.T) {
	b := NewCPUBackend()
	data := b.Alloc(10)
	// Free should not panic.
	b.Free(data)
}

func TestCPUBackendCopyToDevice(t *testing.T) {
	b := NewCPUBackend()
	orig := []float64{1, 2, 3}
	copied := b.CopyToDevice(orig)
	if !floatsEqual(orig, copied) {
		t.Errorf("CopyToDevice: expected %v, got %v", orig, copied)
	}
	// Verify it is a different slice.
	copied[0] = 999
	if orig[0] == 999 {
		t.Error("CopyToDevice should return a separate slice")
	}
}

func TestCPUBackendCopyFromDevice(t *testing.T) {
	b := NewCPUBackend()
	orig := []float64{4, 5, 6}
	copied := b.CopyFromDevice(orig)
	if !floatsEqual(orig, copied) {
		t.Errorf("CopyFromDevice: expected %v, got %v", orig, copied)
	}
	copied[0] = 999
	if orig[0] == 999 {
		t.Error("CopyFromDevice should return a separate slice")
	}
}

// --- Device info tests ---

func TestCPUBackendDeviceCount(t *testing.T) {
	b := NewCPUBackend()
	if b.DeviceCount() != 1 {
		t.Errorf("DeviceCount: expected 1, got %d", b.DeviceCount())
	}
}

func TestCPUBackendDeviceName(t *testing.T) {
	b := NewCPUBackend()
	if b.DeviceName(0) != "cpu" {
		t.Errorf("DeviceName(0): expected 'cpu', got %q", b.DeviceName(0))
	}
	if b.DeviceName(1) != "" {
		t.Errorf("DeviceName(1): expected '', got %q", b.DeviceName(1))
	}
}

func TestCPUBackendMemoryUsed(t *testing.T) {
	b := NewCPUBackend()
	mem := b.MemoryUsed()
	if mem <= 0 {
		t.Errorf("MemoryUsed: expected positive value, got %d", mem)
	}
}

func TestCPUBackendMemoryTotal(t *testing.T) {
	b := NewCPUBackend()
	mem := b.MemoryTotal()
	if mem <= 0 {
		t.Errorf("MemoryTotal: expected positive value, got %d", mem)
	}
}
