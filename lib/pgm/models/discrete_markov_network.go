package models

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

// DiscreteMarkovNetwork is a Markov network where all variables are discrete.
// It embeds *MarkovNetwork and adds discrete-specific validation.
type DiscreteMarkovNetwork struct {
	*MarkovNetwork
}

// NewDiscreteMarkovNetwork creates a new empty DiscreteMarkovNetwork.
func NewDiscreteMarkovNetwork() *DiscreteMarkovNetwork {
	return &DiscreteMarkovNetwork{
		MarkovNetwork: NewMarkovNetwork(),
	}
}

// AddFactor adds a factor to the network with additional discrete-specific
// checks: all cardinalities must be positive and all values must be finite.
func (dmn *DiscreteMarkovNetwork) AddFactor(f *factors.DiscreteFactor) error {
	if f == nil {
		return fmt.Errorf("models: factor must not be nil")
	}

	// Validate cardinalities are positive.
	for i, c := range f.Cardinality() {
		if c <= 0 {
			return fmt.Errorf("models: cardinality at index %d must be positive, got %d", i, c)
		}
	}

	// Validate values contain no NaN or Inf.
	for _, v := range f.Values().Data() {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("models: factor contains NaN or Inf values")
		}
	}

	return dmn.MarkovNetwork.AddFactor(f)
}

// CheckModel validates the discrete Markov network. It calls
// MarkovNetwork.CheckModel() and then performs additional discrete-specific
// checks:
//   - All factor cardinalities must be positive integers.
//   - All factor values must be finite (no NaN or Inf).
//   - All factor values must be non-negative (valid potentials).
func (dmn *DiscreteMarkovNetwork) CheckModel() error {
	if err := dmn.MarkovNetwork.CheckModel(); err != nil {
		return err
	}

	for i, f := range dmn.MarkovNetwork.GetFactors() {
		// Verify cardinalities are positive.
		for j, c := range f.Cardinality() {
			if c <= 0 {
				return fmt.Errorf("models: factor %d has non-positive cardinality %d at index %d", i, c, j)
			}
		}

		// Verify values are finite and non-negative.
		for _, v := range f.Values().Data() {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return fmt.Errorf("models: factor %d contains NaN or Inf values", i)
			}
			if v < 0 {
				return fmt.Errorf("models: factor %d contains negative value %f", i, v)
			}
		}
	}

	return nil
}

// Copy returns a deep copy of the DiscreteMarkovNetwork.
func (dmn *DiscreteMarkovNetwork) Copy() *DiscreteMarkovNetwork {
	return &DiscreteMarkovNetwork{
		MarkovNetwork: dmn.MarkovNetwork.Copy(),
	}
}
