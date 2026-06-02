//go:build unit

package utils

import (
	"math"
	"testing"
)

func TestValidatePositiveInt(t *testing.T) {
	if err := ValidatePositiveInt("x", 1); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := ValidatePositiveInt("x", 100); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := ValidatePositiveInt("x", 0); err == nil {
		t.Error("expected error for 0")
	}
	if err := ValidatePositiveInt("x", -5); err == nil {
		t.Error("expected error for -5")
	}
}

func TestValidateNonNegativeFloat(t *testing.T) {
	if err := ValidateNonNegativeFloat("x", 0.0); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := ValidateNonNegativeFloat("x", 3.14); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := ValidateNonNegativeFloat("x", -0.001); err == nil {
		t.Error("expected error for negative value")
	}
}

func TestValidateProbability(t *testing.T) {
	valid := []float64{0.0, 0.5, 1.0, 0.001, 0.999}
	for _, v := range valid {
		if err := ValidateProbability("p", v); err != nil {
			t.Errorf("unexpected error for %g: %v", v, err)
		}
	}
	invalid := []float64{-0.1, 1.1, -1.0, 2.0}
	for _, v := range invalid {
		if err := ValidateProbability("p", v); err == nil {
			t.Errorf("expected error for %g", v)
		}
	}
}

func TestValidateProbabilityDistribution(t *testing.T) {
	// valid distribution
	if err := ValidateProbabilityDistribution("d", []float64{0.3, 0.7}, 1e-9); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// within tolerance
	if err := ValidateProbabilityDistribution("d", []float64{0.3, 0.7 + 1e-12}, 1e-9); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// outside tolerance
	if err := ValidateProbabilityDistribution("d", []float64{0.3, 0.6}, 1e-9); err == nil {
		t.Error("expected error for sum != 1")
	}

	// empty sums to 0
	if err := ValidateProbabilityDistribution("d", []float64{}, 1e-9); err == nil {
		t.Error("expected error for empty distribution")
	}

	// single element
	if err := ValidateProbabilityDistribution("d", []float64{1.0}, 1e-9); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateNoNaN(t *testing.T) {
	if err := ValidateNoNaN("x", []float64{1.0, 2.0, 3.0}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := ValidateNoNaN("x", []float64{}); err != nil {
		t.Errorf("unexpected error for empty: %v", err)
	}
	if err := ValidateNoNaN("x", []float64{1.0, math.NaN(), 3.0}); err == nil {
		t.Error("expected error for NaN")
	}
	// check that error message contains the index
	err := ValidateNoNaN("vals", []float64{math.NaN()})
	if err == nil {
		t.Fatal("expected error")
	}
	if got := err.Error(); got != "vals[0] is NaN" {
		t.Errorf("unexpected error message: %s", got)
	}
}

func TestValidateNoInf(t *testing.T) {
	if err := ValidateNoInf("x", []float64{1.0, 2.0, 3.0}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := ValidateNoInf("x", []float64{}); err != nil {
		t.Errorf("unexpected error for empty: %v", err)
	}
	if err := ValidateNoInf("x", []float64{1.0, math.Inf(1)}); err == nil {
		t.Error("expected error for +Inf")
	}
	if err := ValidateNoInf("x", []float64{math.Inf(-1), 2.0}); err == nil {
		t.Error("expected error for -Inf")
	}
}
