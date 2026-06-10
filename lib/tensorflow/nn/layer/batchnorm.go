package layer

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// BatchNorm implements batch normalization.
// Input shape: (batch, features). Normalizes across the batch dimension.
// Analogous to tf.keras.layers.BatchNormalization.
type BatchNorm struct {
	numFeatures int
	epsilon     float64
	momentum    float64
	gamma       []float64 // scale parameter
	beta        []float64 // shift parameter
	runningMean []float64
	runningVar  []float64
	training    bool
}

// NewBatchNorm creates a new BatchNorm layer for the given number of features.
func NewBatchNorm(numFeatures int, epsilon, momentum float64) (*BatchNorm, error) {
	if numFeatures <= 0 {
		return nil, fmt.Errorf("batchnorm: numFeatures must be positive, got %d", numFeatures)
	}
	gamma := make([]float64, numFeatures)
	beta := make([]float64, numFeatures)
	runningMean := make([]float64, numFeatures)
	runningVar := make([]float64, numFeatures)
	for i := range numFeatures {
		gamma[i] = 1.0
		runningVar[i] = 1.0
	}
	return &BatchNorm{
		numFeatures: numFeatures,
		epsilon:     epsilon,
		momentum:    momentum,
		gamma:       gamma,
		beta:        beta,
		runningMean: runningMean,
		runningVar:  runningVar,
		training:    true,
	}, nil
}

// Forward applies batch normalization.
// Input shape: (batch, features).
func (bn *BatchNorm) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	if len(shape) != 2 || shape[1] != bn.numFeatures {
		return nil, fmt.Errorf("batchnorm: expected shape (batch, %d), got %v", bn.numFeatures, shape)
	}
	batch := shape[0]
	data := input.Data()
	outData := make([]float64, len(data))

	if bn.training {
		// Compute batch mean and variance per feature.
		mean := make([]float64, bn.numFeatures)
		variance := make([]float64, bn.numFeatures)

		for f := range bn.numFeatures {
			sum := 0.0
			for b := range batch {
				sum += data[b*bn.numFeatures+f]
			}
			mean[f] = sum / float64(batch)
		}
		for f := range bn.numFeatures {
			sum := 0.0
			for b := range batch {
				diff := data[b*bn.numFeatures+f] - mean[f]
				sum += diff * diff
			}
			variance[f] = sum / float64(batch)
		}

		// Normalize and apply gamma/beta.
		for b := range batch {
			for f := range bn.numFeatures {
				idx := b*bn.numFeatures + f
				normalized := (data[idx] - mean[f]) / math.Sqrt(variance[f]+bn.epsilon)
				outData[idx] = bn.gamma[f]*normalized + bn.beta[f]
			}
		}

		// Update running statistics.
		for f := range bn.numFeatures {
			bn.runningMean[f] = bn.momentum*bn.runningMean[f] + (1-bn.momentum)*mean[f]
			bn.runningVar[f] = bn.momentum*bn.runningVar[f] + (1-bn.momentum)*variance[f]
		}
	} else {
		// Use running statistics for inference.
		for b := range batch {
			for f := range bn.numFeatures {
				idx := b*bn.numFeatures + f
				normalized := (data[idx] - bn.runningMean[f]) / math.Sqrt(bn.runningVar[f]+bn.epsilon)
				outData[idx] = bn.gamma[f]*normalized + bn.beta[f]
			}
		}
	}

	return numgo.NewNDArray(shape, outData), nil
}

// SetTraining enables or disables training mode.
func (bn *BatchNorm) SetTraining(training bool) {
	bn.training = training
}

// NumFeatures returns the number of features.
func (bn *BatchNorm) NumFeatures() int {
	return bn.numFeatures
}
