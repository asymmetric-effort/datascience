package layer

import (
	"fmt"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// Embedding maps integer indices to dense vectors.
// Input shape: (batch, seqLen) with float64 values representing integer indices.
// Output shape: (batch, seqLen, embeddingDim).
// Analogous to tf.keras.layers.Embedding.
type Embedding struct {
	weights      *numgo.NDArray // shape: (vocabSize, embeddingDim)
	vocabSize    int
	embeddingDim int
}

// NewEmbedding creates a new Embedding layer.
func NewEmbedding(vocabSize, embeddingDim int, rng func() float64) (*Embedding, error) {
	if vocabSize <= 0 || embeddingDim <= 0 {
		return nil, fmt.Errorf("embedding: vocabSize and embeddingDim must be positive, got %d, %d", vocabSize, embeddingDim)
	}
	data := make([]float64, vocabSize*embeddingDim)
	for i := range data {
		data[i] = (2*rng() - 1) * 0.05
	}
	weights := numgo.NewNDArray([]int{vocabSize, embeddingDim}, data)
	return &Embedding{
		weights:      weights,
		vocabSize:    vocabSize,
		embeddingDim: embeddingDim,
	}, nil
}

// Forward looks up embeddings for the input indices.
// Input values are treated as integer indices into the embedding table.
func (e *Embedding) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	if len(shape) != 2 {
		return nil, fmt.Errorf("embedding: expected rank-2 input (batch, seqLen), got rank %d", len(shape))
	}
	batch, seqLen := shape[0], shape[1]
	outSize := batch * seqLen * e.embeddingDim
	outData := make([]float64, outSize)
	inData := input.Data()
	wData := e.weights.Data()

	for b := range batch {
		for s := range seqLen {
			idx := int(inData[b*seqLen+s])
			if idx < 0 || idx >= e.vocabSize {
				return nil, fmt.Errorf("embedding: index %d out of range [0, %d)", idx, e.vocabSize)
			}
			outOff := (b*seqLen + s) * e.embeddingDim
			wOff := idx * e.embeddingDim
			copy(outData[outOff:outOff+e.embeddingDim], wData[wOff:wOff+e.embeddingDim])
		}
	}

	return numgo.NewNDArray([]int{batch, seqLen, e.embeddingDim}, outData), nil
}

// VocabSize returns the vocabulary size.
func (e *Embedding) VocabSize() int {
	return e.vocabSize
}

// EmbeddingDim returns the embedding dimension.
func (e *Embedding) EmbeddingDim() int {
	return e.embeddingDim
}

// Weights returns a copy of the embedding weight matrix.
func (e *Embedding) Weights() *numgo.NDArray {
	return e.weights.Copy()
}
