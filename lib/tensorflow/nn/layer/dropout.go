package layer

import (
	"fmt"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// Dropout randomly sets elements to zero during training.
// Analogous to tf.keras.layers.Dropout.
type Dropout struct {
	rate    float64
	rng     func() float64
	enabled bool
}

// NewDropout creates a new Dropout layer with the given rate.
// rate must be in [0, 1). rng should return uniform random values in [0, 1).
func NewDropout(rate float64, rng func() float64) (*Dropout, error) {
	if rate < 0 || rate >= 1 {
		return nil, fmt.Errorf("dropout rate must be in [0, 1), got %f", rate)
	}
	return &Dropout{rate: rate, rng: rng, enabled: true}, nil
}

// Forward applies dropout during training. Elements are zeroed with
// probability `rate` and surviving elements are scaled by 1/(1-rate).
func (d *Dropout) Forward(input *numgo.NDArray) (*numgo.NDArray, error) {
	data := input.Data()
	if !d.enabled || d.rate == 0 {
		return numgo.NewNDArray(input.Shape(), data), nil
	}
	scale := 1.0 / (1.0 - d.rate)
	for i := range data {
		if d.rng() < d.rate {
			data[i] = 0
		} else {
			data[i] *= scale
		}
	}
	return numgo.NewNDArray(input.Shape(), data), nil
}

// SetTraining enables or disables dropout.
func (d *Dropout) SetTraining(training bool) {
	d.enabled = training
}

// Rate returns the dropout rate.
func (d *Dropout) Rate() float64 {
	return d.rate
}
