// Package metrics provides evaluation metrics for model performance,
// analogous to tf.keras.metrics.
package metrics

import (
	"fmt"
	"math"
	"sort"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// Metric is the interface all metrics implement.
type Metric interface {
	// Update adds a batch of predictions and targets.
	Update(predictions, targets *numgo.NDArray) error
	// Result returns the current metric value.
	Result() float64
	// Reset clears accumulated state.
	Reset()
	// Name returns the metric name.
	Name() string
}

// shapeEqual returns true if two shapes are identical.
func shapeEqual(a, b []int) bool {
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

// MeanMetric accumulates values and returns their mean.
type MeanMetric struct {
	name  string
	total float64
	count int
}

// NewMeanMetric creates a new MeanMetric.
func NewMeanMetric(name string) *MeanMetric {
	return &MeanMetric{name: name}
}

// Update adds values to the metric.
func (m *MeanMetric) Update(predictions, targets *numgo.NDArray) error {
	data := predictions.Data()
	for _, v := range data {
		m.total += v
		m.count++
	}
	return nil
}

// Result returns the mean of all accumulated values.
func (m *MeanMetric) Result() float64 {
	if m.count == 0 {
		return 0
	}
	return m.total / float64(m.count)
}

// Reset clears accumulated state.
func (m *MeanMetric) Reset() {
	m.total = 0
	m.count = 0
}

// Name returns the metric name.
func (m *MeanMetric) Name() string {
	return m.name
}

// BinaryAccuracy computes the fraction of correct binary predictions.
// Predictions > threshold are treated as 1, else 0.
type BinaryAccuracy struct {
	threshold float64
	correct   int
	total     int
}

// NewBinaryAccuracy creates a new BinaryAccuracy metric.
func NewBinaryAccuracy(threshold float64) *BinaryAccuracy {
	return &BinaryAccuracy{threshold: threshold}
}

// Update adds a batch of predictions and targets.
func (a *BinaryAccuracy) Update(predictions, targets *numgo.NDArray) error {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return fmt.Errorf("binary_accuracy: shape mismatch %v vs %v", predictions.Shape(), targets.Shape())
	}
	pData := predictions.Data()
	tData := targets.Data()
	for i := range pData {
		predicted := 0.0
		if pData[i] > a.threshold {
			predicted = 1.0
		}
		if predicted == tData[i] {
			a.correct++
		}
		a.total++
	}
	return nil
}

// Result returns the accuracy.
func (a *BinaryAccuracy) Result() float64 {
	if a.total == 0 {
		return 0
	}
	return float64(a.correct) / float64(a.total)
}

// Reset clears accumulated state.
func (a *BinaryAccuracy) Reset() {
	a.correct = 0
	a.total = 0
}

// Name returns "binary_accuracy".
func (a *BinaryAccuracy) Name() string {
	return "binary_accuracy"
}

// CategoricalAccuracy computes accuracy for one-hot encoded targets.
// Both inputs have shape (batch, numClasses).
type CategoricalAccuracy struct {
	correct int
	total   int
}

// NewCategoricalAccuracy creates a new CategoricalAccuracy metric.
func NewCategoricalAccuracy() *CategoricalAccuracy {
	return &CategoricalAccuracy{}
}

// Update adds a batch of predictions and targets.
func (a *CategoricalAccuracy) Update(predictions, targets *numgo.NDArray) error {
	pShape := predictions.Shape()
	tShape := targets.Shape()
	if !shapeEqual(pShape, tShape) || len(pShape) != 2 {
		return fmt.Errorf("categorical_accuracy: expected matching rank-2 arrays, got %v and %v", pShape, tShape)
	}
	batch := pShape[0]
	numClasses := pShape[1]
	pData := predictions.Data()
	tData := targets.Data()

	for b := range batch {
		pArgMax := 0
		tArgMax := 0
		for c := 1; c < numClasses; c++ {
			if pData[b*numClasses+c] > pData[b*numClasses+pArgMax] {
				pArgMax = c
			}
			if tData[b*numClasses+c] > tData[b*numClasses+tArgMax] {
				tArgMax = c
			}
		}
		if pArgMax == tArgMax {
			a.correct++
		}
		a.total++
	}
	return nil
}

// Result returns the accuracy.
func (a *CategoricalAccuracy) Result() float64 {
	if a.total == 0 {
		return 0
	}
	return float64(a.correct) / float64(a.total)
}

// Reset clears accumulated state.
func (a *CategoricalAccuracy) Reset() {
	a.correct = 0
	a.total = 0
}

// Name returns "categorical_accuracy".
func (a *CategoricalAccuracy) Name() string {
	return "categorical_accuracy"
}

// Precision computes the precision: tp / (tp + fp).
type Precision struct {
	threshold float64
	truePos   int
	falsePos  int
}

// NewPrecision creates a new Precision metric.
func NewPrecision(threshold float64) *Precision {
	return &Precision{threshold: threshold}
}

// Update adds a batch of binary predictions and targets.
func (p *Precision) Update(predictions, targets *numgo.NDArray) error {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return fmt.Errorf("precision: shape mismatch")
	}
	pData := predictions.Data()
	tData := targets.Data()
	for i := range pData {
		predicted := pData[i] > p.threshold
		actual := tData[i] == 1.0
		if predicted && actual {
			p.truePos++
		} else if predicted && !actual {
			p.falsePos++
		}
	}
	return nil
}

// Result returns the precision value.
func (p *Precision) Result() float64 {
	denom := p.truePos + p.falsePos
	if denom == 0 {
		return 0
	}
	return float64(p.truePos) / float64(denom)
}

// Reset clears accumulated state.
func (p *Precision) Reset() {
	p.truePos = 0
	p.falsePos = 0
}

// Name returns "precision".
func (p *Precision) Name() string {
	return "precision"
}

// Recall computes the recall: tp / (tp + fn).
type Recall struct {
	threshold float64
	truePos   int
	falseNeg  int
}

// NewRecall creates a new Recall metric.
func NewRecall(threshold float64) *Recall {
	return &Recall{threshold: threshold}
}

// Update adds a batch of binary predictions and targets.
func (r *Recall) Update(predictions, targets *numgo.NDArray) error {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return fmt.Errorf("recall: shape mismatch")
	}
	pData := predictions.Data()
	tData := targets.Data()
	for i := range pData {
		predicted := pData[i] > r.threshold
		actual := tData[i] == 1.0
		if predicted && actual {
			r.truePos++
		} else if !predicted && actual {
			r.falseNeg++
		}
	}
	return nil
}

// Result returns the recall value.
func (r *Recall) Result() float64 {
	denom := r.truePos + r.falseNeg
	if denom == 0 {
		return 0
	}
	return float64(r.truePos) / float64(denom)
}

// Reset clears accumulated state.
func (r *Recall) Reset() {
	r.truePos = 0
	r.falseNeg = 0
}

// Name returns "recall".
func (r *Recall) Name() string {
	return "recall"
}

// F1Score computes the F1 score: 2 * (precision * recall) / (precision + recall).
type F1Score struct {
	precision *Precision
	recall    *Recall
}

// NewF1Score creates a new F1Score metric.
func NewF1Score(threshold float64) *F1Score {
	return &F1Score{
		precision: NewPrecision(threshold),
		recall:    NewRecall(threshold),
	}
}

// Update adds a batch of predictions and targets.
func (f *F1Score) Update(predictions, targets *numgo.NDArray) error {
	if err := f.precision.Update(predictions, targets); err != nil {
		return err
	}
	return f.recall.Update(predictions, targets)
}

// Result returns the F1 score.
func (f *F1Score) Result() float64 {
	p := f.precision.Result()
	r := f.recall.Result()
	if p+r == 0 {
		return 0
	}
	return 2 * p * r / (p + r)
}

// Reset clears accumulated state.
func (f *F1Score) Reset() {
	f.precision.Reset()
	f.recall.Reset()
}

// Name returns "f1_score".
func (f *F1Score) Name() string {
	return "f1_score"
}

// AUC computes the Area Under the ROC Curve using the trapezoidal rule.
type AUC struct {
	predictions []float64
	targets     []float64
}

// NewAUC creates a new AUC metric.
func NewAUC() *AUC {
	return &AUC{}
}

// Update adds a batch of predictions and targets.
func (a *AUC) Update(predictions, targets *numgo.NDArray) error {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return fmt.Errorf("auc: shape mismatch")
	}
	a.predictions = append(a.predictions, predictions.Data()...)
	a.targets = append(a.targets, targets.Data()...)
	return nil
}

// Result computes the AUC-ROC.
func (a *AUC) Result() float64 {
	n := len(a.predictions)
	if n == 0 {
		return 0
	}

	// Create index sorted by prediction descending.
	indices := make([]int, n)
	for i := range indices {
		indices[i] = i
	}
	preds := a.predictions
	sort.Slice(indices, func(i, j int) bool {
		return preds[indices[i]] > preds[indices[j]]
	})

	totalPos := 0
	totalNeg := 0
	for _, t := range a.targets {
		if t == 1 {
			totalPos++
		} else {
			totalNeg++
		}
	}
	if totalPos == 0 || totalNeg == 0 {
		return 0
	}

	// Walk thresholds and compute trapezoidal AUC.
	auc := 0.0
	tp := 0
	fp := 0
	prevTPR := 0.0
	prevFPR := 0.0

	for i, idx := range indices {
		if a.targets[idx] == 1 {
			tp++
		} else {
			fp++
		}

		// At each threshold change, compute a trapezoid.
		if i == n-1 || preds[indices[i]] != preds[indices[i+1]] {
			tpr := float64(tp) / float64(totalPos)
			fpr := float64(fp) / float64(totalNeg)
			auc += (fpr - prevFPR) * (tpr + prevTPR) / 2
			prevTPR = tpr
			prevFPR = fpr
		}
	}

	return math.Min(math.Max(auc, 0), 1)
}

// Reset clears accumulated state.
func (a *AUC) Reset() {
	a.predictions = a.predictions[:0]
	a.targets = a.targets[:0]
}

// Name returns "auc".
func (a *AUC) Name() string {
	return "auc"
}
