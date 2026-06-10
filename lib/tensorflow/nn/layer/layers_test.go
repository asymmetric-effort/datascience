package layer

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

func makeTestTensor(t *testing.T, dims []int, data []float64) *numgo.NDArray {
	t.Helper()
	return numgo.NewNDArray(dims, data)
}

// --- Flatten Tests ---

func TestFlatten(t *testing.T) {
	input := makeTestTensor(t, []int{2, 3, 4}, make([]float64, 24))
	f := NewFlatten()
	result, err := f.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{2, 12}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestFlatten2D(t *testing.T) {
	input := makeTestTensor(t, []int{3, 5}, make([]float64, 15))
	f := NewFlatten()
	result, err := f.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{3, 5}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

// --- Dropout Tests ---

func TestDropout(t *testing.T) {
	counter := 0
	rng := func() float64 {
		counter++
		if counter%2 == 0 {
			return 0.1 // below rate=0.5 -> drop
		}
		return 0.9 // above rate=0.5 -> keep
	}
	d, err := NewDropout(0.5, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := makeTestTensor(t, []int{4}, []float64{1, 1, 1, 1})
	result, err := d.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Some elements should be 0, others scaled by 2.
	data := result.Data()
	hasZero := false
	hasScaled := false
	for _, v := range data {
		if v == 0 {
			hasZero = true
		}
		if math.Abs(v-2.0) < 1e-10 {
			hasScaled = true
		}
	}
	if !hasZero || !hasScaled {
		t.Errorf("expected mix of zeroed and scaled values, got %v", data)
	}
}

func TestDropoutDisabled(t *testing.T) {
	rng := func() float64 { return 0 }
	d, _ := NewDropout(0.5, rng)
	d.SetTraining(false)
	input := makeTestTensor(t, []int{3}, []float64{1, 2, 3})
	result, err := d.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inputData := input.Data()
	for i, v := range result.Data() {
		if v != inputData[i] {
			t.Errorf("disabled dropout changed data[%d]: %f", i, v)
		}
	}
}

func TestDropoutZeroRate(t *testing.T) {
	rng := func() float64 { return 0 }
	d, _ := NewDropout(0.0, rng)
	input := makeTestTensor(t, []int{3}, []float64{1, 2, 3})
	result, err := d.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inputData := input.Data()
	for i, v := range result.Data() {
		if v != inputData[i] {
			t.Errorf("zero-rate dropout changed data[%d]: %f", i, v)
		}
	}
}

func TestDropoutInvalidRate(t *testing.T) {
	rng := func() float64 { return 0 }
	_, err := NewDropout(1.0, rng)
	if err == nil {
		t.Error("expected error for rate=1.0")
	}
	_, err = NewDropout(-0.1, rng)
	if err == nil {
		t.Error("expected error for negative rate")
	}
}

func TestDropoutRate(t *testing.T) {
	rng := func() float64 { return 0 }
	d, _ := NewDropout(0.3, rng)
	if d.Rate() != 0.3 {
		t.Errorf("Rate() = %f, want 0.3", d.Rate())
	}
}

// --- BatchNorm Tests ---

func TestBatchNorm(t *testing.T) {
	bn, err := NewBatchNorm(2, 1e-5, 0.1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	input := makeTestTensor(t, []int{4, 2}, []float64{
		1, 10,
		2, 20,
		3, 30,
		4, 40,
	})
	result, err := bn.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// After batch normalization, each feature should have mean≈0, var≈1.
	data := result.Data()
	for f := range 2 {
		sum := 0.0
		for b := range 4 {
			sum += data[b*2+f]
		}
		mean := sum / 4
		if math.Abs(mean) > 1e-5 {
			t.Errorf("feature %d mean = %f, expected ~0", f, mean)
		}
	}
}

func TestBatchNormInference(t *testing.T) {
	bn, _ := NewBatchNorm(2, 1e-5, 0.1)
	// Train first to set running stats.
	trainInput := makeTestTensor(t, []int{4, 2}, []float64{1, 10, 2, 20, 3, 30, 4, 40})
	_, _ = bn.Forward(trainInput)

	// Switch to inference mode.
	bn.SetTraining(false)
	testInput := makeTestTensor(t, []int{2, 2}, []float64{2.5, 25, 3.5, 35})
	result, err := bn.Forward(testInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Size() != 4 {
		t.Errorf("expected 4 elements, got %d", result.Size())
	}
}

func TestBatchNormWrongShape(t *testing.T) {
	bn, _ := NewBatchNorm(3, 1e-5, 0.1)
	input := makeTestTensor(t, []int{2, 2}, []float64{1, 2, 3, 4})
	_, err := bn.Forward(input)
	if err == nil {
		t.Error("expected error for wrong feature count")
	}
}

func TestBatchNormInvalidFeatures(t *testing.T) {
	_, err := NewBatchNorm(0, 1e-5, 0.1)
	if err == nil {
		t.Error("expected error for zero features")
	}
}

func TestBatchNormNumFeatures(t *testing.T) {
	bn, _ := NewBatchNorm(5, 1e-5, 0.1)
	if bn.NumFeatures() != 5 {
		t.Errorf("NumFeatures() = %d, want 5", bn.NumFeatures())
	}
}

// --- Conv2D Tests ---

func TestConv2D(t *testing.T) {
	rng := deterministicRNG()
	conv, err := NewConv2D(1, 1, 3, 3, 1, 1, false, rng)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Input: batch=1, 5x5, 1 channel.
	inData := make([]float64, 25)
	for i := range 25 {
		inData[i] = 1.0
	}
	input := numgo.NewNDArray([]int{1, 5, 5, 1}, inData)
	result, err := conv.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Output should be 3x3 (valid padding).
	expectedShape := []int{1, 3, 3, 1}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestConv2DSamePadding(t *testing.T) {
	rng := deterministicRNG()
	conv, _ := NewConv2D(1, 2, 3, 3, 1, 1, true, rng)
	input := makeTestTensor(t, []int{1, 4, 4, 1}, make([]float64, 16))
	result, err := conv.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Same padding should preserve spatial dims.
	expectedShape := []int{1, 4, 4, 2}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestConv2DWrongInput(t *testing.T) {
	rng := deterministicRNG()
	conv, _ := NewConv2D(3, 1, 3, 3, 1, 1, false, rng)
	input := makeTestTensor(t, []int{1, 4, 4, 1}, make([]float64, 16))
	_, err := conv.Forward(input)
	if err == nil {
		t.Error("expected error for wrong channel count")
	}
}

func TestConv2DInvalidParams(t *testing.T) {
	rng := deterministicRNG()
	_, err := NewConv2D(0, 1, 3, 3, 1, 1, false, rng)
	if err == nil {
		t.Error("expected error for zero inChannels")
	}
	_, err = NewConv2D(1, 1, 3, 3, 0, 1, false, rng)
	if err == nil {
		t.Error("expected error for zero stride")
	}
}

func TestConv2DFilters(t *testing.T) {
	rng := deterministicRNG()
	conv, _ := NewConv2D(1, 2, 3, 3, 1, 1, false, rng)
	f := conv.Filters()
	expectedShape := []int{3, 3, 1, 2}
	if !shapeEq(f.Shape(), expectedShape) {
		t.Errorf("filters shape = %v, want %v", f.Shape(), expectedShape)
	}
	if conv.OutChannels() != 2 {
		t.Errorf("OutChannels() = %d, want 2", conv.OutChannels())
	}
}

// --- MaxPool2D Tests ---

func TestMaxPool2D(t *testing.T) {
	pool, err := NewMaxPool2D(2, 2, 2, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Input: batch=1, 4x4, 1 channel.
	input := makeTestTensor(t, []int{1, 4, 4, 1}, []float64{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
	})
	result, err := pool.Forward(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{1, 2, 2, 1}
	if !shapeEq(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
	want := []float64{6, 8, 14, 16}
	for i, v := range result.Data() {
		if v != want[i] {
			t.Errorf("data[%d] = %f, want %f", i, v, want[i])
		}
	}
}

func TestMaxPool2DWrongRank(t *testing.T) {
	pool, _ := NewMaxPool2D(2, 2, 2, 2)
	input := makeTestTensor(t, []int{4, 4}, make([]float64, 16))
	_, err := pool.Forward(input)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

func TestMaxPool2DInvalidParams(t *testing.T) {
	_, err := NewMaxPool2D(0, 2, 2, 2)
	if err == nil {
		t.Error("expected error for zero pool size")
	}
}

func TestMaxPool2DOutputTooSmall(t *testing.T) {
	pool, _ := NewMaxPool2D(5, 5, 1, 1)
	input := makeTestTensor(t, []int{1, 3, 3, 1}, make([]float64, 9))
	_, err := pool.Forward(input)
	if err == nil {
		t.Error("expected error for output dimensions <= 0")
	}
}
