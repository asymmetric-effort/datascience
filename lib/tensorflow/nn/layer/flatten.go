package layer

import (
	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// Flatten reshapes a tensor to 2D, preserving the batch dimension.
// Input shape: (batch, d1, d2, ...). Output shape: (batch, d1*d2*...).
// Analogous to tf.keras.layers.Flatten.
type Flatten struct{}

// NewFlatten creates a new Flatten layer.
func NewFlatten() *Flatten {
	return &Flatten{}
}

// Forward flattens the input, keeping the first dimension (batch).
func (f *Flatten) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	shape := input.Shape()
	batch := shape[0]
	flatSize := input.Size() / batch
	data := input.Data()
	return numgo.NewNDArray([]int{batch, flatSize}, data), nil
}
