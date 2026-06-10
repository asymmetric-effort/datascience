package keras

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

func deterministicRNG() func() float64 {
	state := 0.5
	return func() float64 {
		state = math.Mod(state*1.1+0.3, 1.0)
		return state
	}
}

func makeArray(t *testing.T, dims []int, data []float64) *numgo.NDArray {
	t.Helper()
	return numgo.NewNDArray(dims, data)
}

// testDense is a local dense layer that implements the keras.Layer interface
// using *numgo.NDArray, since nn/layer still uses tensor.Tensor.
type testDense struct {
	weights *numgo.NDArray
	bias    *numgo.NDArray
	inSize  int
	outSize int
}

func newTestDense(inSize, outSize int, rng func() float64) (*testDense, error) {
	if inSize <= 0 || outSize <= 0 {
		return nil, fmt.Errorf("dense layer sizes must be positive: in=%d, out=%d", inSize, outSize)
	}
	scale := math.Sqrt(6.0 / float64(inSize+outSize))
	wData := make([]float64, inSize*outSize)
	for i := range wData {
		wData[i] = (2*rng() - 1) * scale
	}
	weights := numgo.NewNDArray([]int{inSize, outSize}, wData)
	bias := numgo.Zeros(outSize)
	return &testDense{weights: weights, bias: bias, inSize: inSize, outSize: outSize}, nil
}

func (d *testDense) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	if len(shape) != 2 || shape[1] != d.inSize {
		return nil, fmt.Errorf("dense forward: expected input shape (batch, %d), got %v", d.inSize, shape)
	}
	batch := shape[0]
	outData := make([]float64, batch*d.outSize)
	inData := input.Data()
	wData := d.weights.Data()
	bData := d.bias.Data()
	for b := range batch {
		for o := range d.outSize {
			sum := bData[o]
			for k := range d.inSize {
				sum += inData[b*d.inSize+k] * wData[k*d.outSize+o]
			}
			outData[b*d.outSize+o] = sum
		}
	}
	return numgo.NewNDArray([]int{batch, d.outSize}, outData), nil
}

func (d *testDense) Weights() *numgo.NDArray {
	return numgo.NewNDArray([]int{d.inSize, d.outSize}, d.weights.Data())
}

func (d *testDense) Bias() *numgo.NDArray {
	return numgo.NewNDArray([]int{d.outSize}, d.bias.Data())
}

func (d *testDense) SetWeights(w *numgo.NDArray) error {
	expected := []int{d.inSize, d.outSize}
	if !shapeEqual(w.Shape(), expected) {
		return fmt.Errorf("expected weight shape %v, got %v", expected, w.Shape())
	}
	d.weights = numgo.NewNDArray(expected, w.Data())
	return nil
}

func (d *testDense) SetBias(b *numgo.NDArray) error {
	expected := []int{d.outSize}
	if !shapeEqual(b.Shape(), expected) {
		return fmt.Errorf("expected bias shape %v, got %v", expected, b.Shape())
	}
	d.bias = numgo.NewNDArray(expected, b.Data())
	return nil
}

func TestNewSequential(t *testing.T) {
	model := NewSequential()
	if model == nil {
		t.Fatal("NewSequential returned nil")
	}
	if model.NumLayers() != 0 {
		t.Errorf("NumLayers() = %d, want 0", model.NumLayers())
	}
	if model.Trained() {
		t.Error("new model should not be trained")
	}
}

func TestSequentialAdd(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(3, 2, deterministicRNG())
	model.Add(dense)
	if model.NumLayers() != 1 {
		t.Errorf("NumLayers() = %d, want 1", model.NumLayers())
	}
}

func TestSequentialPredict(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(3, 2, deterministicRNG())
	model.Add(dense)

	input := makeArray(t, []int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	result, err := model.Predict(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedShape := []int{2, 2}
	if !shapeEqual(result.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", result.Shape(), expectedShape)
	}
}

func TestSequentialEvaluate(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(2, 1, deterministicRNG())
	model.Add(dense)
	model.Compile(MeanSquaredError, 0.01)

	x := makeArray(t, []int{4, 2}, []float64{1, 0, 0, 1, 1, 1, 0, 0})
	y := makeArray(t, []int{4, 1}, []float64{1, 0, 1, 0})
	loss, err := model.Evaluate(x, y)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.IsNaN(loss) || math.IsInf(loss, 0) {
		t.Errorf("loss is %f, expected finite", loss)
	}
}

func TestSequentialEvaluateNotCompiled(t *testing.T) {
	model := NewSequential()
	x := makeArray(t, []int{1, 2}, []float64{1, 2})
	y := makeArray(t, []int{1, 1}, []float64{1})
	_, err := model.Evaluate(x, y)
	if err == nil {
		t.Error("expected error for uncompiled model")
	}
}

func TestSequentialFit(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(2, 1, deterministicRNG())
	model.Add(dense)
	model.Compile(MeanSquaredError, 0.1)

	// Simple linear data: y = x1 + x2.
	x := makeArray(t, []int{4, 2}, []float64{
		0, 0,
		1, 0,
		0, 1,
		1, 1,
	})
	y := makeArray(t, []int{4, 1}, []float64{0, 1, 1, 2})

	history, err := model.Fit(x, y, 3, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 3 {
		t.Errorf("history length = %d, want 3", len(history))
	}
	if !model.Trained() {
		t.Error("model should be marked as trained")
	}
	// Loss should generally decrease (not guaranteed in 3 epochs but should be finite).
	for _, loss := range history {
		if math.IsNaN(loss) || math.IsInf(loss, 0) {
			t.Errorf("loss is %f, expected finite", loss)
		}
	}
}

func TestSequentialFitNotCompiled(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(2, 1, deterministicRNG())
	model.Add(dense)
	x := makeArray(t, []int{1, 2}, []float64{1, 2})
	y := makeArray(t, []int{1, 1}, []float64{1})
	_, err := model.Fit(x, y, 1, 1)
	if err == nil {
		t.Error("expected error for uncompiled model")
	}
}

func TestSequentialFitNoLayers(t *testing.T) {
	model := NewSequential()
	model.Compile(MeanSquaredError, 0.01)
	x := makeArray(t, []int{1, 2}, []float64{1, 2})
	y := makeArray(t, []int{1, 1}, []float64{1})
	_, err := model.Fit(x, y, 1, 1)
	if err == nil {
		t.Error("expected error for model with no layers")
	}
}

func TestSequentialFitInvalidBatchSize(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(2, 1, deterministicRNG())
	model.Add(dense)
	model.Compile(MeanSquaredError, 0.01)
	x := makeArray(t, []int{1, 2}, []float64{1, 2})
	y := makeArray(t, []int{1, 1}, []float64{1})
	_, err := model.Fit(x, y, 1, 0)
	if err == nil {
		t.Error("expected error for zero batch size")
	}
}

func TestSequentialFitInvalidEpochs(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(2, 1, deterministicRNG())
	model.Add(dense)
	model.Compile(MeanSquaredError, 0.01)
	x := makeArray(t, []int{1, 2}, []float64{1, 2})
	y := makeArray(t, []int{1, 1}, []float64{1})
	_, err := model.Fit(x, y, 0, 1)
	if err == nil {
		t.Error("expected error for zero epochs")
	}
}

func TestSequentialFitNoSamples(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(2, 1, deterministicRNG())
	model.Add(dense)
	model.Compile(MeanSquaredError, 0.01)
	x := makeArray(t, []int{0, 2}, []float64{})
	y := makeArray(t, []int{0, 1}, []float64{})
	_, err := model.Fit(x, y, 1, 1)
	if err == nil {
		t.Error("expected error for no samples")
	}
}

func TestSequentialSummary(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(3, 2, deterministicRNG())
	model.Add(dense)
	model.Compile(MeanSquaredError, 0.01)
	summary := model.Summary()
	if !strings.Contains(summary, "Sequential") {
		t.Error("summary should contain 'Sequential'")
	}
	if !strings.Contains(summary, "Layers: 1") {
		t.Error("summary should show 1 layer")
	}
	if !strings.Contains(summary, "0.01") {
		t.Error("summary should show learning rate")
	}
}

func TestMeanSquaredErrorLoss(t *testing.T) {
	a := makeArray(t, []int{2}, []float64{1, 3})
	b := makeArray(t, []int{2}, []float64{2, 1})
	loss, err := MeanSquaredError(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if math.Abs(loss-2.5) > 1e-10 {
		t.Errorf("MSE = %f, want 2.5", loss)
	}
}

func TestMeanSquaredErrorShapeMismatch(t *testing.T) {
	a := makeArray(t, []int{2}, []float64{1, 2})
	b := makeArray(t, []int{3}, []float64{1, 2, 3})
	_, err := MeanSquaredError(a, b)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestBinaryCrossEntropyLoss(t *testing.T) {
	a := makeArray(t, []int{2}, []float64{0.9, 0.1})
	b := makeArray(t, []int{2}, []float64{1, 0})
	loss, err := BinaryCrossEntropy(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loss > 0.2 {
		t.Errorf("BCE = %f, expected low loss for good predictions", loss)
	}
}

func TestBinaryCrossEntropyShapeMismatch(t *testing.T) {
	a := makeArray(t, []int{2}, []float64{0.5, 0.5})
	b := makeArray(t, []int{3}, []float64{1, 0, 1})
	_, err := BinaryCrossEntropy(a, b)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

// failLayer always returns an error from Forward.
type failLayer struct{}

func (f failLayer) Forward(_ *numgo.NDArray) (*numgo.NDArray, error) {
	return nil, fmt.Errorf("intentional failure")
}

func TestSequentialForwardError(t *testing.T) {
	model := NewSequential()
	model.Add(failLayer{})
	input := makeArray(t, []int{1, 2}, []float64{1, 2})
	_, err := model.Predict(input)
	if err == nil {
		t.Error("expected error from failing layer")
	}
}

func TestSequentialEvaluateForwardError(t *testing.T) {
	model := NewSequential()
	model.Add(failLayer{})
	model.Compile(MeanSquaredError, 0.01)
	x := makeArray(t, []int{1, 2}, []float64{1, 2})
	y := makeArray(t, []int{1, 1}, []float64{1})
	_, err := model.Evaluate(x, y)
	if err == nil {
		t.Error("expected error from failing layer in evaluate")
	}
}

func TestSequentialFitForwardError(t *testing.T) {
	model := NewSequential()
	model.Add(failLayer{})
	model.Compile(MeanSquaredError, 0.01)
	x := makeArray(t, []int{2, 2}, []float64{1, 2, 3, 4})
	y := makeArray(t, []int{2, 1}, []float64{1, 2})
	_, err := model.Fit(x, y, 1, 2)
	if err == nil {
		t.Error("expected error from failing layer in fit")
	}
}

func TestSequentialFitLossError(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(2, 3, deterministicRNG()) // output 3
	model.Add(dense)
	badLoss := func(a, b *numgo.NDArray) (float64, error) {
		return 0, fmt.Errorf("loss shape error")
	}
	model.Compile(badLoss, 0.01)
	x := makeArray(t, []int{2, 2}, []float64{1, 2, 3, 4})
	y := makeArray(t, []int{2, 1}, []float64{1, 2})
	_, err := model.Fit(x, y, 1, 2)
	if err == nil {
		t.Error("expected error from loss function in fit")
	}
}

func TestSequentialFitBatchSizeLargerThanSamples(t *testing.T) {
	model := NewSequential()
	dense, _ := newTestDense(2, 1, deterministicRNG())
	model.Add(dense)
	model.Compile(MeanSquaredError, 0.01)
	x := makeArray(t, []int{2, 2}, []float64{1, 2, 3, 4})
	y := makeArray(t, []int{2, 1}, []float64{1, 2})
	history, err := model.Fit(x, y, 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(history) != 1 {
		t.Errorf("history length = %d, want 1", len(history))
	}
}

func TestExtractBatch(t *testing.T) {
	arr := makeArray(t, []int{4, 2}, []float64{1, 2, 3, 4, 5, 6, 7, 8})
	batch := extractBatch(arr, 1, 3)
	expectedShape := []int{2, 2}
	if !shapeEqual(batch.Shape(), expectedShape) {
		t.Errorf("shape = %v, want %v", batch.Shape(), expectedShape)
	}
	want := []float64{3, 4, 5, 6}
	for i, v := range batch.Data() {
		if v != want[i] {
			t.Errorf("data[%d] = %f, want %f", i, v, want[i])
		}
	}
}
