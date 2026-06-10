package metrics

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

func makeArray(t *testing.T, dims []int, data []float64) *numgo.NDArray {
	t.Helper()
	return numgo.NewNDArray(dims, data)
}

// --- MeanMetric ---

func TestMeanMetric(t *testing.T) {
	m := NewMeanMetric("test_mean")
	if m.Name() != "test_mean" {
		t.Errorf("Name() = %q", m.Name())
	}
	preds := makeArray(t, []int{4}, []float64{2, 4, 6, 8})
	_ = m.Update(preds, nil)
	if math.Abs(m.Result()-5.0) > 1e-10 {
		t.Errorf("Result() = %f, want 5.0", m.Result())
	}
	m.Reset()
	if m.Result() != 0 {
		t.Errorf("after reset Result() = %f, want 0", m.Result())
	}
}

func TestMeanMetricEmpty(t *testing.T) {
	m := NewMeanMetric("empty")
	if m.Result() != 0 {
		t.Errorf("empty Result() = %f, want 0", m.Result())
	}
}

// --- BinaryAccuracy ---

func TestBinaryAccuracy(t *testing.T) {
	a := NewBinaryAccuracy(0.5)
	if a.Name() != "binary_accuracy" {
		t.Errorf("Name() = %q", a.Name())
	}
	preds := makeArray(t, []int{4}, []float64{0.9, 0.1, 0.8, 0.2})
	targets := makeArray(t, []int{4}, []float64{1, 0, 1, 0})
	err := a.Update(preds, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Result() != 1.0 {
		t.Errorf("Result() = %f, want 1.0", a.Result())
	}
}

func TestBinaryAccuracyPartial(t *testing.T) {
	a := NewBinaryAccuracy(0.5)
	preds := makeArray(t, []int{4}, []float64{0.9, 0.9, 0.1, 0.1})
	targets := makeArray(t, []int{4}, []float64{1, 0, 1, 0})
	_ = a.Update(preds, targets)
	if a.Result() != 0.5 {
		t.Errorf("Result() = %f, want 0.5", a.Result())
	}
}

func TestBinaryAccuracyEmpty(t *testing.T) {
	a := NewBinaryAccuracy(0.5)
	if a.Result() != 0 {
		t.Errorf("empty Result() = %f, want 0", a.Result())
	}
}

func TestBinaryAccuracyShapeMismatch(t *testing.T) {
	a := NewBinaryAccuracy(0.5)
	p := makeArray(t, []int{3}, []float64{1, 2, 3})
	tgt := makeArray(t, []int{2}, []float64{1, 0})
	err := a.Update(p, tgt)
	if err == nil {
		t.Error("expected error for shape mismatch")
	}
}

func TestBinaryAccuracyReset(t *testing.T) {
	a := NewBinaryAccuracy(0.5)
	preds := makeArray(t, []int{2}, []float64{0.9, 0.1})
	targets := makeArray(t, []int{2}, []float64{1, 0})
	_ = a.Update(preds, targets)
	a.Reset()
	if a.Result() != 0 {
		t.Errorf("after reset Result() = %f, want 0", a.Result())
	}
}

// --- CategoricalAccuracy ---

func TestCategoricalAccuracy(t *testing.T) {
	a := NewCategoricalAccuracy()
	if a.Name() != "categorical_accuracy" {
		t.Errorf("Name() = %q", a.Name())
	}
	preds := makeArray(t, []int{3, 3}, []float64{
		0.1, 0.8, 0.1,
		0.7, 0.2, 0.1,
		0.1, 0.1, 0.8,
	})
	targets := makeArray(t, []int{3, 3}, []float64{
		0, 1, 0,
		1, 0, 0,
		0, 0, 1,
	})
	err := a.Update(preds, targets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Result() != 1.0 {
		t.Errorf("Result() = %f, want 1.0", a.Result())
	}
}

func TestCategoricalAccuracyEmpty(t *testing.T) {
	a := NewCategoricalAccuracy()
	if a.Result() != 0 {
		t.Errorf("empty Result() = %f, want 0", a.Result())
	}
}

func TestCategoricalAccuracyWrongRank(t *testing.T) {
	a := NewCategoricalAccuracy()
	p := makeArray(t, []int{6}, []float64{1, 2, 3, 4, 5, 6})
	tgt := makeArray(t, []int{6}, []float64{1, 0, 0, 0, 1, 0})
	err := a.Update(p, tgt)
	if err == nil {
		t.Error("expected error for wrong rank")
	}
}

func TestCategoricalAccuracyReset(t *testing.T) {
	a := NewCategoricalAccuracy()
	preds := makeArray(t, []int{1, 2}, []float64{0.8, 0.2})
	targets := makeArray(t, []int{1, 2}, []float64{1, 0})
	_ = a.Update(preds, targets)
	a.Reset()
	if a.Result() != 0 {
		t.Errorf("after reset Result() = %f, want 0", a.Result())
	}
}

// --- Precision ---

func TestPrecision(t *testing.T) {
	p := NewPrecision(0.5)
	if p.Name() != "precision" {
		t.Errorf("Name() = %q", p.Name())
	}
	preds := makeArray(t, []int{6}, []float64{0.9, 0.8, 0.7, 0.1, 0.2, 0.3})
	targets := makeArray(t, []int{6}, []float64{1, 1, 0, 0, 0, 1})
	_ = p.Update(preds, targets)
	// TP=2, FP=1 => precision=2/3
	if math.Abs(p.Result()-2.0/3.0) > 1e-10 {
		t.Errorf("Result() = %f, want 0.667", p.Result())
	}
}

func TestPrecisionEmpty(t *testing.T) {
	p := NewPrecision(0.5)
	if p.Result() != 0 {
		t.Errorf("empty Result() = %f, want 0", p.Result())
	}
}

func TestPrecisionShapeMismatch(t *testing.T) {
	p := NewPrecision(0.5)
	pr := makeArray(t, []int{3}, []float64{1, 2, 3})
	tgt := makeArray(t, []int{2}, []float64{1, 0})
	err := p.Update(pr, tgt)
	if err == nil {
		t.Error("expected error")
	}
}

func TestPrecisionReset(t *testing.T) {
	p := NewPrecision(0.5)
	preds := makeArray(t, []int{2}, []float64{0.9, 0.9})
	targets := makeArray(t, []int{2}, []float64{1, 1})
	_ = p.Update(preds, targets)
	p.Reset()
	if p.Result() != 0 {
		t.Errorf("after reset Result() = %f, want 0", p.Result())
	}
}

// --- Recall ---

func TestRecall(t *testing.T) {
	r := NewRecall(0.5)
	if r.Name() != "recall" {
		t.Errorf("Name() = %q", r.Name())
	}
	preds := makeArray(t, []int{6}, []float64{0.9, 0.1, 0.8, 0.1, 0.9, 0.1})
	targets := makeArray(t, []int{6}, []float64{1, 1, 1, 0, 0, 0})
	_ = r.Update(preds, targets)
	// TP=2, FN=1 => recall=2/3
	if math.Abs(r.Result()-2.0/3.0) > 1e-10 {
		t.Errorf("Result() = %f, want 0.667", r.Result())
	}
}

func TestRecallEmpty(t *testing.T) {
	r := NewRecall(0.5)
	if r.Result() != 0 {
		t.Errorf("empty Result() = %f, want 0", r.Result())
	}
}

func TestRecallShapeMismatch(t *testing.T) {
	r := NewRecall(0.5)
	pr := makeArray(t, []int{3}, []float64{1, 2, 3})
	tgt := makeArray(t, []int{2}, []float64{1, 0})
	err := r.Update(pr, tgt)
	if err == nil {
		t.Error("expected error")
	}
}

func TestRecallReset(t *testing.T) {
	r := NewRecall(0.5)
	preds := makeArray(t, []int{2}, []float64{0.9, 0.9})
	targets := makeArray(t, []int{2}, []float64{1, 1})
	_ = r.Update(preds, targets)
	r.Reset()
	if r.Result() != 0 {
		t.Errorf("after reset Result() = %f, want 0", r.Result())
	}
}

// --- F1Score ---

func TestF1Score(t *testing.T) {
	f := NewF1Score(0.5)
	if f.Name() != "f1_score" {
		t.Errorf("Name() = %q", f.Name())
	}
	// Perfect predictions.
	preds := makeArray(t, []int{4}, []float64{0.9, 0.1, 0.9, 0.1})
	targets := makeArray(t, []int{4}, []float64{1, 0, 1, 0})
	_ = f.Update(preds, targets)
	if f.Result() != 1.0 {
		t.Errorf("Result() = %f, want 1.0", f.Result())
	}
}

func TestF1ScoreEmpty(t *testing.T) {
	f := NewF1Score(0.5)
	if f.Result() != 0 {
		t.Errorf("empty Result() = %f, want 0", f.Result())
	}
}

func TestF1ScoreReset(t *testing.T) {
	f := NewF1Score(0.5)
	preds := makeArray(t, []int{2}, []float64{0.9, 0.9})
	targets := makeArray(t, []int{2}, []float64{1, 1})
	_ = f.Update(preds, targets)
	f.Reset()
	if f.Result() != 0 {
		t.Errorf("after reset Result() = %f, want 0", f.Result())
	}
}

// --- AUC ---

func TestAUCPerfect(t *testing.T) {
	a := NewAUC()
	if a.Name() != "auc" {
		t.Errorf("Name() = %q", a.Name())
	}
	preds := makeArray(t, []int{4}, []float64{0.9, 0.8, 0.2, 0.1})
	targets := makeArray(t, []int{4}, []float64{1, 1, 0, 0})
	_ = a.Update(preds, targets)
	if a.Result() != 1.0 {
		t.Errorf("Result() = %f, want 1.0", a.Result())
	}
}

func TestAUCRandom(t *testing.T) {
	a := NewAUC()
	// All same predictions => AUC ~0.5 or 0 (degenerate).
	preds := makeArray(t, []int{4}, []float64{0.5, 0.5, 0.5, 0.5})
	targets := makeArray(t, []int{4}, []float64{1, 0, 1, 0})
	_ = a.Update(preds, targets)
	// With tied scores, AUC should still be computed.
	result := a.Result()
	if result < 0 || result > 1 {
		t.Errorf("Result() = %f, out of [0, 1]", result)
	}
}

func TestAUCEmpty(t *testing.T) {
	a := NewAUC()
	if a.Result() != 0 {
		t.Errorf("empty Result() = %f, want 0", a.Result())
	}
}

func TestAUCAllPositive(t *testing.T) {
	a := NewAUC()
	preds := makeArray(t, []int{3}, []float64{0.9, 0.8, 0.7})
	targets := makeArray(t, []int{3}, []float64{1, 1, 1})
	_ = a.Update(preds, targets)
	// No negatives => AUC undefined, returns 0.
	if a.Result() != 0 {
		t.Errorf("Result() = %f, want 0 for all positive", a.Result())
	}
}

func TestAUCShapeMismatch(t *testing.T) {
	a := NewAUC()
	p := makeArray(t, []int{3}, []float64{1, 2, 3})
	tgt := makeArray(t, []int{2}, []float64{1, 0})
	err := a.Update(p, tgt)
	if err == nil {
		t.Error("expected error")
	}
}

func TestAUCReset(t *testing.T) {
	a := NewAUC()
	preds := makeArray(t, []int{4}, []float64{0.9, 0.8, 0.2, 0.1})
	targets := makeArray(t, []int{4}, []float64{1, 1, 0, 0})
	_ = a.Update(preds, targets)
	a.Reset()
	if a.Result() != 0 {
		t.Errorf("after reset Result() = %f, want 0", a.Result())
	}
}

// --- Interface Compliance ---

func TestMetricInterface(t *testing.T) {
	var _ Metric = &MeanMetric{}
	var _ Metric = &BinaryAccuracy{}
	var _ Metric = &CategoricalAccuracy{}
	var _ Metric = &Precision{}
	var _ Metric = &Recall{}
	var _ Metric = &F1Score{}
	var _ Metric = &AUC{}
}
