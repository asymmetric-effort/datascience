// Package keras provides high-level model building and training APIs,
// analogous to tf.keras.
package keras

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// maxBatchSize limits the maximum batch size to prevent unbounded allocation.
const maxBatchSize = 1 << 16 // 65536

// Layer is the interface that all neural network layers must implement.
type Layer interface {
	Forward(input *numgo.NDArray) (*numgo.NDArray, error)
}

// LossFunc computes the loss between predictions and targets.
type LossFunc func(predictions, targets *numgo.NDArray) (float64, error)

// Sequential is a linear stack of layers.
// Analogous to tf.keras.Sequential.
type Sequential struct {
	layers  []Layer
	lossFn  LossFunc
	lr      float64
	trained bool
}

// NewSequential creates a new empty Sequential model.
func NewSequential() *Sequential {
	return &Sequential{}
}

// Add appends a layer to the model.
func (s *Sequential) Add(layer Layer) {
	s.layers = append(s.layers, layer)
}

// Compile configures the model for training with a loss function and learning rate.
func (s *Sequential) Compile(lossFn LossFunc, lr float64) {
	s.lossFn = lossFn
	s.lr = lr
}

// forward runs the input through all layers.
func (s *Sequential) forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	current := input
	for i, layer := range s.layers {
		var err error
		current, err = layer.Forward(current)
		if err != nil {
			return nil, fmt.Errorf("layer %d: %w", i, err)
		}
	}
	return current, nil
}

// Predict runs forward inference on the input data.
func (s *Sequential) Predict(input *numgo.NDArray) (*numgo.NDArray, error) {
	return s.forward(input)
}

// Evaluate computes the loss on the given data without training.
func (s *Sequential) Evaluate(x, y *numgo.NDArray) (float64, error) {
	if s.lossFn == nil {
		return 0, fmt.Errorf("model not compiled: call Compile() first")
	}
	predictions, err := s.forward(x)
	if err != nil {
		return 0, fmt.Errorf("evaluate forward: %w", err)
	}
	return s.lossFn(predictions, y)
}

// Fit trains the model using numerical gradient descent.
// Returns the loss history for each epoch.
// x shape: (samples, features...), y shape: (samples, outputs...).
func (s *Sequential) Fit(x, y *numgo.NDArray, epochs, batchSize int) ([]float64, error) {
	if s.lossFn == nil {
		return nil, fmt.Errorf("model not compiled: call Compile() first")
	}
	if len(s.layers) == 0 {
		return nil, fmt.Errorf("model has no layers")
	}
	if batchSize <= 0 || batchSize > maxBatchSize {
		return nil, fmt.Errorf("batch size must be in [1, %d], got %d", maxBatchSize, batchSize)
	}
	if epochs <= 0 {
		return nil, fmt.Errorf("epochs must be positive, got %d", epochs)
	}

	xShape := x.Shape()
	numSamples := xShape[0]
	if numSamples == 0 {
		return nil, fmt.Errorf("no training samples")
	}

	lossHistory := make([]float64, 0, epochs)

	for epoch := range epochs {
		epochLoss := 0.0
		numBatches := 0

		for start := 0; start < numSamples; start += batchSize {
			end := start + batchSize
			if end > numSamples {
				end = numSamples
			}

			batchX := extractBatch(x, start, end)
			batchY := extractBatch(y, start, end)

			predictions, err := s.forward(batchX)
			if err != nil {
				return nil, fmt.Errorf("epoch %d: forward: %w", epoch, err)
			}

			loss, err := s.lossFn(predictions, batchY)
			if err != nil {
				return nil, fmt.Errorf("epoch %d: loss: %w", epoch, err)
			}
			epochLoss += loss
			numBatches++

			// Compute numerical gradients and update weights using finite differences.
			err = s.numericalUpdate(batchX, batchY)
			if err != nil {
				return nil, fmt.Errorf("epoch %d: update: %w", epoch, err)
			}
		}

		lossHistory = append(lossHistory, epochLoss/float64(numBatches))
	}

	s.trained = true
	return lossHistory, nil
}

// Trained returns whether the model has been trained.
func (s *Sequential) Trained() bool {
	return s.trained
}

// NumLayers returns the number of layers in the model.
func (s *Sequential) NumLayers() int {
	return len(s.layers)
}

// numericalUpdate performs a weight update using finite-difference gradient estimation.
// This is a simple approach suitable for demonstrating the training API.
func (s *Sequential) numericalUpdate(x, y *numgo.NDArray) error {
	epsilon := 1e-5

	for _, layer := range s.layers {
		type paramAccessor interface {
			Weights() *numgo.NDArray
			Bias() *numgo.NDArray
			SetWeights(w *numgo.NDArray) error
			SetBias(b *numgo.NDArray) error
		}

		pa, ok := layer.(paramAccessor)
		if !ok {
			continue
		}

		// Update weights.
		weights := pa.Weights()
		wData := weights.Data()
		wGrad := make([]float64, len(wData))

		for i := range wData {
			original := wData[i]

			wData[i] = original + epsilon
			wPlus := numgo.NewNDArray(weights.Shape(), wData)
			_ = pa.SetWeights(wPlus)
			predPlus, _ := s.forward(x)
			lossPlus, _ := s.lossFn(predPlus, y)

			wData[i] = original - epsilon
			wMinus := numgo.NewNDArray(weights.Shape(), wData)
			_ = pa.SetWeights(wMinus)
			predMinus, _ := s.forward(x)
			lossMinus, _ := s.lossFn(predMinus, y)

			wGrad[i] = (lossPlus - lossMinus) / (2 * epsilon)
			wData[i] = original
		}

		// Apply gradient to weights.
		for i := range wData {
			wData[i] -= s.lr * wGrad[i]
		}
		updated := numgo.NewNDArray(weights.Shape(), wData)
		_ = pa.SetWeights(updated)

		// Update bias.
		bias := pa.Bias()
		bData := bias.Data()
		bGrad := make([]float64, len(bData))

		for i := range bData {
			original := bData[i]

			bData[i] = original + epsilon
			bPlus := numgo.NewNDArray(bias.Shape(), bData)
			_ = pa.SetBias(bPlus)
			predPlus, _ := s.forward(x)
			lossPlus, _ := s.lossFn(predPlus, y)

			bData[i] = original - epsilon
			bMinus := numgo.NewNDArray(bias.Shape(), bData)
			_ = pa.SetBias(bMinus)
			predMinus, _ := s.forward(x)
			lossMinus, _ := s.lossFn(predMinus, y)

			bGrad[i] = (lossPlus - lossMinus) / (2 * epsilon)
			bData[i] = original
		}

		for i := range bData {
			bData[i] -= s.lr * bGrad[i]
		}
		updatedBias := numgo.NewNDArray(bias.Shape(), bData)
		_ = pa.SetBias(updatedBias)
	}
	return nil
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

// extractBatch slices a batch from the NDArray along axis 0.
func extractBatch(a *numgo.NDArray, start, end int) *numgo.NDArray {
	shape := a.Shape()
	rank := len(shape)

	batchDims := make([]int, rank)
	batchDims[0] = end - start
	for i := 1; i < rank; i++ {
		batchDims[i] = shape[i]
	}

	elementsPerSample := 1
	for i := 1; i < rank; i++ {
		elementsPerSample *= shape[i]
	}

	data := a.Data()
	batchData := make([]float64, (end-start)*elementsPerSample)
	copy(batchData, data[start*elementsPerSample:end*elementsPerSample])

	return numgo.NewNDArray(batchDims, batchData)
}

// Summary returns a string describing the model architecture.
func (s *Sequential) Summary() string {
	result := "Sequential Model\n"
	result += fmt.Sprintf("Layers: %d\n", len(s.layers))
	for i, layer := range s.layers {
		result += fmt.Sprintf("  [%d] %T\n", i, layer)
	}
	if s.lossFn != nil {
		result += fmt.Sprintf("Learning Rate: %g\n", s.lr)
	}
	result += fmt.Sprintf("Trained: %v\n", s.trained)
	return result
}

// MeanSquaredError is a convenience loss function for use with Compile.
func MeanSquaredError(predictions, targets *numgo.NDArray) (float64, error) {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return 0, fmt.Errorf("mse: shape mismatch %v vs %v", predictions.Shape(), targets.Shape())
	}
	pData := predictions.Data()
	tData := targets.Data()
	sum := 0.0
	for i := range pData {
		diff := pData[i] - tData[i]
		sum += diff * diff
	}
	return sum / float64(len(pData)), nil
}

// BinaryCrossEntropy is a convenience loss function for classification.
func BinaryCrossEntropy(predictions, targets *numgo.NDArray) (float64, error) {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return 0, fmt.Errorf("bce: shape mismatch %v vs %v", predictions.Shape(), targets.Shape())
	}
	pData := predictions.Data()
	tData := targets.Data()
	sum := 0.0
	eps := 1e-7
	for i := range pData {
		p := math.Max(eps, math.Min(1-eps, pData[i]))
		sum += -(tData[i]*math.Log(p) + (1-tData[i])*math.Log(1-p))
	}
	return sum / float64(len(pData)), nil
}
