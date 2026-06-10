// Package initializer provides weight initialization strategies,
// analogous to tf.initializers / tf.keras.initializers.
package initializer

import (
	"math"

	"github.com/asymmetric-effort/datascience/lib/numgo"
)

// Initializer is the interface for weight initialization strategies.
type Initializer interface {
	// Initialize creates an NDArray of the given shape with initialized values.
	Initialize(shape []int) (*numgo.NDArray, error)
}

// ZerosInit initializes all weights to zero.
type ZerosInit struct{}

// Initialize returns a zero-filled NDArray.
func (z ZerosInit) Initialize(shape []int) (*numgo.NDArray, error) {
	return numgo.Zeros(shape...), nil
}

// OnesInit initializes all weights to one.
type OnesInit struct{}

// Initialize returns a ones-filled NDArray.
func (o OnesInit) Initialize(shape []int) (*numgo.NDArray, error) {
	return numgo.Ones(shape...), nil
}

// ConstantInit initializes all weights to a constant value.
type ConstantInit struct {
	Value float64
}

// Initialize returns a constant-filled NDArray.
func (c ConstantInit) Initialize(shape []int) (*numgo.NDArray, error) {
	return numgo.Full(c.Value, shape...), nil
}

// RandomNormalInit initializes weights from a normal distribution.
type RandomNormalInit struct {
	Mean   float64
	StdDev float64
	Seed   int64
}

// Initialize returns an NDArray with normally distributed values.
func (r RandomNormalInit) Initialize(shape []int) (*numgo.NDArray, error) {
	rng := numgo.NewRNG(r.Seed)
	return rng.Normal(r.Mean, r.StdDev, shape...), nil
}

// RandomUniformInit initializes weights from a uniform distribution.
type RandomUniformInit struct {
	MinVal float64
	MaxVal float64
	Seed   int64
}

// Initialize returns an NDArray with uniformly distributed values.
func (r RandomUniformInit) Initialize(shape []int) (*numgo.NDArray, error) {
	rng := numgo.NewRNG(r.Seed)
	return rng.Uniform(r.MinVal, r.MaxVal, shape...), nil
}

// GlorotUniformInit implements Xavier/Glorot uniform initialization.
// Draws from Uniform(-limit, limit) where limit = sqrt(6 / (fan_in + fan_out)).
type GlorotUniformInit struct {
	Seed int64
}

// Initialize returns an NDArray with Glorot uniform initialized values.
func (g GlorotUniformInit) Initialize(shape []int) (*numgo.NDArray, error) {
	fanIn, fanOut := computeFans(shape)
	limit := math.Sqrt(6.0 / float64(fanIn+fanOut))
	rng := numgo.NewRNG(g.Seed)
	return rng.Uniform(-limit, limit, shape...), nil
}

// GlorotNormalInit implements Xavier/Glorot normal initialization.
// Draws from Normal(0, sqrt(2 / (fan_in + fan_out))).
type GlorotNormalInit struct {
	Seed int64
}

// Initialize returns an NDArray with Glorot normal initialized values.
func (g GlorotNormalInit) Initialize(shape []int) (*numgo.NDArray, error) {
	fanIn, fanOut := computeFans(shape)
	stddev := math.Sqrt(2.0 / float64(fanIn+fanOut))
	rng := numgo.NewRNG(g.Seed)
	return rng.Normal(0, stddev, shape...), nil
}

// HeNormalInit implements He/Kaiming normal initialization.
// Draws from Normal(0, sqrt(2 / fan_in)).
type HeNormalInit struct {
	Seed int64
}

// Initialize returns an NDArray with He normal initialized values.
func (h HeNormalInit) Initialize(shape []int) (*numgo.NDArray, error) {
	fanIn, _ := computeFans(shape)
	stddev := math.Sqrt(2.0 / float64(fanIn))
	rng := numgo.NewRNG(h.Seed)
	return rng.Normal(0, stddev, shape...), nil
}

// HeUniformInit implements He/Kaiming uniform initialization.
// Draws from Uniform(-limit, limit) where limit = sqrt(6 / fan_in).
type HeUniformInit struct {
	Seed int64
}

// Initialize returns an NDArray with He uniform initialized values.
func (h HeUniformInit) Initialize(shape []int) (*numgo.NDArray, error) {
	fanIn, _ := computeFans(shape)
	limit := math.Sqrt(6.0 / float64(fanIn))
	rng := numgo.NewRNG(h.Seed)
	return rng.Uniform(-limit, limit, shape...), nil
}

// LecunNormalInit implements LeCun normal initialization.
// Draws from Normal(0, sqrt(1 / fan_in)).
type LecunNormalInit struct {
	Seed int64
}

// Initialize returns an NDArray with LeCun normal initialized values.
func (l LecunNormalInit) Initialize(shape []int) (*numgo.NDArray, error) {
	fanIn, _ := computeFans(shape)
	stddev := math.Sqrt(1.0 / float64(fanIn))
	rng := numgo.NewRNG(l.Seed)
	return rng.Normal(0, stddev, shape...), nil
}

// computeFans calculates fan_in and fan_out for a weight tensor.
// For 2D (dense): shape is (in, out).
// For 4D (conv): shape is (kH, kW, inChannels, outChannels).
func computeFans(shape []int) (int, int) {
	rank := len(shape)
	switch {
	case rank == 0:
		return 1, 1
	case rank == 1:
		return shape[0], shape[0]
	case rank == 2:
		return shape[0], shape[1]
	default:
		receptiveField := 1
		for i := 0; i < rank-2; i++ {
			receptiveField *= shape[i]
		}
		fanIn := shape[rank-2] * receptiveField
		fanOut := shape[rank-1] * receptiveField
		return fanIn, fanOut
	}
}
