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

	// [1 2]   [5 6]   [1*5+2*7  1*6+2*8]   [19 22]
	// [3 4] x [7 8] = [3*5+4*7  3*6+4*8] = [43 50]
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

	// [1 2 3]   [7  8 ]   [1*7+2*9+3*11   1*8+2*10+3*12]   [ 58  64]
	// [4 5 6] x [9  10] = [4*7+5*9+6*11   4*8+5*10+6*12] = [139 154]
	//           [11 12]
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

	// Factor A: shape [2], values [0.4, 0.6]
	// Factor B: shape [3], values [0.1, 0.2, 0.7]
	// Result shape [2, 3] = outer product
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

	// Multiplying by a single-element factor should scale values.
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

	// Tensor shape [2, 3]:
	// [[1, 2, 3],
	//  [4, 5, 6]]
	values := []float64{1, 2, 3, 4, 5, 6}
	shape := []int{2, 3}

	// Marginalize axis 0 (sum rows): [5, 7, 9]
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

	// Tensor shape [2, 3]:
	// [[1, 2, 3],
	//  [4, 5, 6]]
	values := []float64{1, 2, 3, 4, 5, 6}
	shape := []int{2, 3}

	// Marginalize axis 1 (sum columns): [6, 15]
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

	// Tensor shape [2, 2, 2]:
	// [[[1,2],[3,4]], [[5,6],[7,8]]]
	values := []float64{1, 2, 3, 4, 5, 6, 7, 8}
	shape := []int{2, 2, 2}

	// Marginalize axis 1: sum over the middle axis
	// Result shape [2, 2]
	// [[1+3, 2+4], [5+7, 6+8]] = [[4, 6], [12, 14]]
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
	// Compile-time check that CPUBackend satisfies Backend.
	var _ Backend = (*CPUBackend)(nil)
}
