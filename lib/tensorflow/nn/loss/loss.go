// Package loss provides loss functions for training neural networks,
// analogous to tf.keras.losses.
package loss

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// MeanSquaredError computes MSE: mean((predictions - targets)^2).
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

// MeanAbsoluteError computes MAE: mean(|predictions - targets|).
func MeanAbsoluteError(predictions, targets *numgo.NDArray) (float64, error) {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return 0, fmt.Errorf("mae: shape mismatch %v vs %v", predictions.Shape(), targets.Shape())
	}
	pData := predictions.Data()
	tData := targets.Data()
	sum := 0.0
	for i := range pData {
		sum += math.Abs(pData[i] - tData[i])
	}
	return sum / float64(len(pData)), nil
}

// BinaryCrossEntropy computes binary cross-entropy loss.
func BinaryCrossEntropy(predictions, targets *numgo.NDArray) (float64, error) {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return 0, fmt.Errorf("binary_crossentropy: shape mismatch %v vs %v", predictions.Shape(), targets.Shape())
	}
	pData := predictions.Data()
	tData := targets.Data()
	sum := 0.0
	epsilon := 1e-7
	for i := range pData {
		p := math.Max(epsilon, math.Min(1-epsilon, pData[i]))
		sum += -(tData[i]*math.Log(p) + (1-tData[i])*math.Log(1-p))
	}
	return sum / float64(len(pData)), nil
}

// CategoricalCrossEntropy computes categorical cross-entropy loss.
func CategoricalCrossEntropy(predictions, targets *numgo.NDArray) (float64, error) {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return 0, fmt.Errorf("categorical_crossentropy: shape mismatch %v vs %v", predictions.Shape(), targets.Shape())
	}
	pShape := predictions.Shape()
	if len(pShape) != 2 {
		return 0, fmt.Errorf("categorical_crossentropy: expected rank 2, got %d", len(pShape))
	}
	pData := predictions.Data()
	tData := targets.Data()
	batch := pShape[0]
	numClasses := pShape[1]
	epsilon := 1e-7
	totalLoss := 0.0
	for b := range batch {
		sampleLoss := 0.0
		for c := range numClasses {
			idx := b*numClasses + c
			p := math.Max(epsilon, math.Min(1-epsilon, pData[idx]))
			sampleLoss += -tData[idx] * math.Log(p)
		}
		totalLoss += sampleLoss
	}
	return totalLoss / float64(batch), nil
}

// HuberLoss computes the Huber loss.
func HuberLoss(predictions, targets *numgo.NDArray, delta float64) (float64, error) {
	if !shapeEqual(predictions.Shape(), targets.Shape()) {
		return 0, fmt.Errorf("huber: shape mismatch %v vs %v", predictions.Shape(), targets.Shape())
	}
	pData := predictions.Data()
	tData := targets.Data()
	sum := 0.0
	for i := range pData {
		diff := math.Abs(pData[i] - tData[i])
		if diff <= delta {
			sum += 0.5 * diff * diff
		} else {
			sum += delta * (diff - 0.5*delta)
		}
	}
	return sum / float64(len(pData)), nil
}

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
