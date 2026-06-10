package loss

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

func makeArr(t *testing.T, dims []int, data []float64) *numgo.NDArray {
	t.Helper()
	return numgo.NewNDArray(dims, data)
}

func TestMeanSquaredError(t *testing.T) {
	preds := makeArr(t, []int{4}, []float64{1, 2, 3, 4})
	targets := makeArr(t, []int{4}, []float64{1, 2, 3, 4})
	loss, err := MeanSquaredError(preds, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loss != 0 {
		t.Errorf("MSE of identical tensors = %f, want 0", loss)
	}
}

func TestMeanSquaredErrorNonZero(t *testing.T) {
	preds := makeArr(t, []int{2}, []float64{1, 3})
	targets := makeArr(t, []int{2}, []float64{2, 1})
	loss, err := MeanSquaredError(preds, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// (1-2)^2 + (3-1)^2 = 1 + 4 = 5, mean = 2.5
	if math.Abs(loss-2.5) > 1e-10 {
		t.Errorf("MSE = %f, want 2.5", loss)
	}
}

func TestMeanSquaredErrorShapeMismatch(t *testing.T) {
	preds := makeArr(t, []int{3}, []float64{1, 2, 3})
	targets := makeArr(t, []int{2}, []float64{1, 2})
	_, err := MeanSquaredError(preds, targets)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestMeanAbsoluteError(t *testing.T) {
	preds := makeArr(t, []int{3}, []float64{1, 5, 3})
	targets := makeArr(t, []int{3}, []float64{2, 3, 3})
	loss, err := MeanAbsoluteError(preds, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// |1-2| + |5-3| + |3-3| = 1 + 2 + 0 = 3, mean = 1.0
	if math.Abs(loss-1.0) > 1e-10 {
		t.Errorf("MAE = %f, want 1.0", loss)
	}
}

func TestMeanAbsoluteErrorShapeMismatch(t *testing.T) {
	preds := makeArr(t, []int{3}, []float64{1, 2, 3})
	targets := makeArr(t, []int{2}, []float64{1, 2})
	_, err := MeanAbsoluteError(preds, targets)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestBinaryCrossEntropy(t *testing.T) {
	// Perfect predictions should give near-zero loss.
	preds := makeArr(t, []int{2}, []float64{0.999, 0.001})
	targets := makeArr(t, []int{2}, []float64{1, 0})
	loss, err := BinaryCrossEntropy(preds, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loss > 0.01 {
		t.Errorf("BCE for near-perfect predictions = %f, expected near 0", loss)
	}
}

func TestBinaryCrossEntropyBadPredictions(t *testing.T) {
	// Opposite predictions should give high loss.
	preds := makeArr(t, []int{2}, []float64{0.001, 0.999})
	targets := makeArr(t, []int{2}, []float64{1, 0})
	loss, err := BinaryCrossEntropy(preds, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loss < 1.0 {
		t.Errorf("BCE for wrong predictions = %f, expected > 1.0", loss)
	}
}

func TestBinaryCrossEntropyShapeMismatch(t *testing.T) {
	preds := makeArr(t, []int{3}, []float64{0.5, 0.5, 0.5})
	targets := makeArr(t, []int{2}, []float64{1, 0})
	_, err := BinaryCrossEntropy(preds, targets)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestCategoricalCrossEntropy(t *testing.T) {
	// Correct one-hot with high probability.
	preds := makeArr(t, []int{2, 3}, []float64{
		0.9, 0.05, 0.05,
		0.1, 0.1, 0.8,
	})
	targets := makeArr(t, []int{2, 3}, []float64{
		1, 0, 0,
		0, 0, 1,
	})
	loss, err := CategoricalCrossEntropy(preds, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loss > 0.3 {
		t.Errorf("CCE for good predictions = %f, expected < 0.3", loss)
	}
}

func TestCategoricalCrossEntropyShapeMismatch(t *testing.T) {
	preds := makeArr(t, []int{2, 3}, []float64{0.5, 0.3, 0.2, 0.1, 0.8, 0.1})
	targets := makeArr(t, []int{2, 2}, []float64{1, 0, 0, 1})
	_, err := CategoricalCrossEntropy(preds, targets)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestCategoricalCrossEntropyWrongRank(t *testing.T) {
	preds := makeArr(t, []int{6}, []float64{0.1, 0.2, 0.3, 0.1, 0.2, 0.1})
	targets := makeArr(t, []int{6}, []float64{1, 0, 0, 0, 0, 0})
	_, err := CategoricalCrossEntropy(preds, targets)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

func TestHuberLoss(t *testing.T) {
	preds := makeArr(t, []int{2}, []float64{1, 10})
	targets := makeArr(t, []int{2}, []float64{1.5, 1})
	loss, err := HuberLoss(preds, targets, 1.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Element 0: |0.5| <= 1.0, so 0.5*0.25 = 0.125
	// Element 1: |9| > 1.0, so 1.0*(9 - 0.5) = 8.5
	// Mean = (0.125 + 8.5) / 2 = 4.3125
	if math.Abs(loss-4.3125) > 1e-10 {
		t.Errorf("Huber = %f, want 4.3125", loss)
	}
}

func TestHuberLossShapeMismatch(t *testing.T) {
	preds := makeArr(t, []int{3}, []float64{1, 2, 3})
	targets := makeArr(t, []int{2}, []float64{1, 2})
	_, err := HuberLoss(preds, targets, 1.0)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}
