package layer

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// MultiHeadAttention implements multi-head scaled dot-product attention.
// Input: query, key, value tensors of shape (batch, seqLen, dModel).
// Output: (batch, seqLen, dModel).
// Analogous to tf.keras.layers.MultiHeadAttention.
type MultiHeadAttention struct {
	numHeads int
	dModel   int
	dK       int       // dModel / numHeads
	wQ       []float64 // (dModel, dModel)
	wK       []float64
	wV       []float64
	wO       []float64
}

// NewMultiHeadAttention creates a new multi-head attention layer.
func NewMultiHeadAttention(dModel, numHeads int, rng func() float64) (*MultiHeadAttention, error) {
	if dModel <= 0 || numHeads <= 0 {
		return nil, fmt.Errorf("attention: dModel and numHeads must be positive")
	}
	if dModel%numHeads != 0 {
		return nil, fmt.Errorf("attention: dModel (%d) must be divisible by numHeads (%d)", dModel, numHeads)
	}

	dK := dModel / numHeads
	scale := math.Sqrt(2.0 / float64(dModel))
	size := dModel * dModel

	initWeights := func() []float64 {
		w := make([]float64, size)
		for i := range w {
			w[i] = (2*rng() - 1) * scale
		}
		return w
	}

	return &MultiHeadAttention{
		numHeads: numHeads,
		dModel:   dModel,
		dK:       dK,
		wQ:       initWeights(),
		wK:       initWeights(),
		wV:       initWeights(),
		wO:       initWeights(),
	}, nil
}

// Forward computes multi-head attention.
// query, key, value all have shape (batch, seqLen, dModel).
// For self-attention, pass the same tensor for all three.
func (mha *MultiHeadAttention) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	return mha.ForwardQKV(input, input, input)
}

// ForwardQKV computes multi-head attention with separate Q, K, V inputs.
func (mha *MultiHeadAttention) ForwardQKV(query, key, value *numgo.NDArray) (*numgo.NDArray, error) {
	qShape := query.Shape()
	kShape := key.Shape()
	vShape := value.Shape()

	if len(qShape) != 3 || len(kShape) != 3 || len(vShape) != 3 {
		return nil, fmt.Errorf("attention: expected rank-3 inputs, got %d, %d, %d", len(qShape), len(kShape), len(vShape))
	}
	if qShape[2] != mha.dModel || kShape[2] != mha.dModel || vShape[2] != mha.dModel {
		return nil, fmt.Errorf("attention: feature dim must be %d", mha.dModel)
	}
	if qShape[0] != kShape[0] || kShape[0] != vShape[0] {
		return nil, fmt.Errorf("attention: batch size mismatch")
	}
	if kShape[1] != vShape[1] {
		return nil, fmt.Errorf("attention: key and value sequence lengths must match")
	}

	batch := qShape[0]
	qLen := qShape[1]
	kvLen := kShape[1]
	qData := query.Data()
	kData := key.Data()
	vData := value.Data()

	// Project Q, K, V: (batch, seqLen, dModel) @ (dModel, dModel) = (batch, seqLen, dModel)
	projQ := mha.project(qData, mha.wQ, batch, qLen)
	projK := mha.project(kData, mha.wK, batch, kvLen)
	projV := mha.project(vData, mha.wV, batch, kvLen)

	// Compute attention per head.
	scaleFactor := 1.0 / math.Sqrt(float64(mha.dK))
	attnOut := make([]float64, batch*qLen*mha.dModel)

	for b := range batch {
		for head := range mha.numHeads {
			headOff := head * mha.dK

			// For each query position, compute attention over all key positions.
			for qi := range qLen {
				// Compute attention scores.
				scores := make([]float64, kvLen)
				for ki := range kvLen {
					dot := 0.0
					for d := range mha.dK {
						qIdx := b*qLen*mha.dModel + qi*mha.dModel + headOff + d
						kIdx := b*kvLen*mha.dModel + ki*mha.dModel + headOff + d
						dot += projQ[qIdx] * projK[kIdx]
					}
					scores[ki] = dot * scaleFactor
				}

				// Softmax over scores.
				maxScore := scores[0]
				for _, s := range scores[1:] {
					if s > maxScore {
						maxScore = s
					}
				}
				expSum := 0.0
				for i := range scores {
					scores[i] = math.Exp(scores[i] - maxScore)
					expSum += scores[i]
				}
				for i := range scores {
					scores[i] /= expSum
				}

				// Weighted sum of values.
				for d := range mha.dK {
					sum := 0.0
					for ki := range kvLen {
						vIdx := b*kvLen*mha.dModel + ki*mha.dModel + headOff + d
						sum += scores[ki] * projV[vIdx]
					}
					outIdx := b*qLen*mha.dModel + qi*mha.dModel + headOff + d
					attnOut[outIdx] = sum
				}
			}
		}
	}

	// Output projection: attnOut @ wO.
	result := mha.project(attnOut, mha.wO, batch, qLen)

	return numgo.NewNDArray([]int{batch, qLen, mha.dModel}, result), nil
}

// NumHeads returns the number of attention heads.
func (mha *MultiHeadAttention) NumHeads() int {
	return mha.numHeads
}

// DModel returns the model dimension.
func (mha *MultiHeadAttention) DModel() int {
	return mha.dModel
}

func (mha *MultiHeadAttention) project(data, weights []float64, batch, seqLen int) []float64 {
	out := make([]float64, batch*seqLen*mha.dModel)
	for b := range batch {
		for s := range seqLen {
			inOff := b*seqLen*mha.dModel + s*mha.dModel
			outOff := inOff
			for o := range mha.dModel {
				sum := 0.0
				for k := range mha.dModel {
					sum += data[inOff+k] * weights[k*mha.dModel+o]
				}
				out[outOff+o] = sum
			}
		}
	}
	return out
}
