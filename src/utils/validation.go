package utils

import (
	"fmt"
	"math"
)

// ValidatePositiveInt returns an error if value is not strictly positive.
func ValidatePositiveInt(name string, value int) error {
	if value <= 0 {
		return fmt.Errorf("%s must be positive, got %d", name, value)
	}
	return nil
}

// ValidateNonNegativeFloat returns an error if value is negative.
func ValidateNonNegativeFloat(name string, value float64) error {
	if value < 0 {
		return fmt.Errorf("%s must be non-negative, got %g", name, value)
	}
	return nil
}

// ValidateProbability returns an error if value is not in [0, 1].
func ValidateProbability(name string, value float64) error {
	if value < 0 || value > 1 {
		return fmt.Errorf("%s must be in [0, 1], got %g", name, value)
	}
	return nil
}

// ValidateProbabilityDistribution returns an error if the values do not sum
// to 1 within the given tolerance.
func ValidateProbabilityDistribution(name string, values []float64, tolerance float64) error {
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	if math.Abs(sum-1.0) > tolerance {
		return fmt.Errorf("%s must sum to 1 (got %g, difference %g exceeds tolerance %g)",
			name, sum, math.Abs(sum-1.0), tolerance)
	}
	return nil
}

// ValidateNoNaN returns an error if any value is NaN.
func ValidateNoNaN(name string, values []float64) error {
	for i, v := range values {
		if math.IsNaN(v) {
			return fmt.Errorf("%s[%d] is NaN", name, i)
		}
	}
	return nil
}

// ValidateNoInf returns an error if any value is infinite.
func ValidateNoInf(name string, values []float64) error {
	for i, v := range values {
		if math.IsInf(v, 0) {
			return fmt.Errorf("%s[%d] is infinite", name, i)
		}
	}
	return nil
}
