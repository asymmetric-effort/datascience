//go:build unit

package numgo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// Seed
// ---------------------------------------------------------------------------

func TestSeed(t *testing.T) {
	Seed(42)
	// Should not panic.
}

// ---------------------------------------------------------------------------
// Random (alias for Rand)
// ---------------------------------------------------------------------------

func TestRandom(t *testing.T) {
	rng := NewRNG(42)
	a := rng.Random(10)
	if a.Size() != 10 {
		t.Fatalf("Random: expected size 10, got %d", a.Size())
	}
	for _, v := range a.Data() {
		if v < 0 || v >= 1 {
			t.Fatalf("Random: value %g out of [0, 1)", v)
		}
	}
}

// ---------------------------------------------------------------------------
// Permutation
// ---------------------------------------------------------------------------

func TestPermutation(t *testing.T) {
	rng := NewRNG(42)
	p := rng.Permutation(10)
	if p.Size() != 10 {
		t.Fatalf("Permutation: expected size 10, got %d", p.Size())
	}
	// Check all values 0..9 are present.
	seen := make(map[float64]bool)
	for _, v := range p.Data() {
		seen[v] = true
	}
	for i := 0; i < 10; i++ {
		if !seen[float64(i)] {
			t.Fatalf("Permutation: missing value %d", i)
		}
	}
}

// ---------------------------------------------------------------------------
// Exponential
// ---------------------------------------------------------------------------

func TestExponential(t *testing.T) {
	rng := NewRNG(42)
	a := rng.Exponential(1.0, 1000)
	if a.Size() != 1000 {
		t.Fatalf("Exponential: expected size 1000, got %d", a.Size())
	}
	// All values should be positive.
	for _, v := range a.Data() {
		if v <= 0 {
			t.Fatalf("Exponential: got non-positive value %g", v)
		}
	}
	// Mean should be close to scale=1.0.
	m := Mean(a).Data()[0]
	if math.Abs(m-1.0) > 0.15 {
		t.Fatalf("Exponential: mean %g not close to 1.0", m)
	}
}

// ---------------------------------------------------------------------------
// Poisson
// ---------------------------------------------------------------------------

func TestPoisson(t *testing.T) {
	rng := NewRNG(42)
	a := rng.Poisson(5.0, 1000)
	if a.Size() != 1000 {
		t.Fatalf("Poisson: expected size 1000, got %d", a.Size())
	}
	// All values should be non-negative integers.
	for _, v := range a.Data() {
		if v < 0 || v != math.Floor(v) {
			t.Fatalf("Poisson: got non-integer or negative value %g", v)
		}
	}
	// Mean should be close to lam=5.0.
	m := Mean(a).Data()[0]
	if math.Abs(m-5.0) > 1.0 {
		t.Fatalf("Poisson: mean %g not close to 5.0", m)
	}
}

// ---------------------------------------------------------------------------
// BinomialSample
// ---------------------------------------------------------------------------

func TestBinomialSample(t *testing.T) {
	rng := NewRNG(42)
	a := rng.BinomialSample(10, 0.5, 1000)
	if a.Size() != 1000 {
		t.Fatalf("BinomialSample: expected size 1000, got %d", a.Size())
	}
	for _, v := range a.Data() {
		if v < 0 || v > 10 || v != math.Floor(v) {
			t.Fatalf("BinomialSample: invalid value %g", v)
		}
	}
	m := Mean(a).Data()[0]
	if math.Abs(m-5.0) > 1.0 {
		t.Fatalf("BinomialSample: mean %g not close to 5.0", m)
	}
}

// ---------------------------------------------------------------------------
// Beta
// ---------------------------------------------------------------------------

func TestBeta(t *testing.T) {
	rng := NewRNG(42)
	a := rng.Beta(2, 5, 1000)
	if a.Size() != 1000 {
		t.Fatalf("Beta: expected size 1000, got %d", a.Size())
	}
	for _, v := range a.Data() {
		if v < 0 || v > 1 {
			t.Fatalf("Beta: value %g out of [0, 1]", v)
		}
	}
	// Mean of Beta(2,5) = 2/7 ~ 0.2857
	m := Mean(a).Data()[0]
	if math.Abs(m-2.0/7.0) > 0.05 {
		t.Fatalf("Beta: mean %g not close to %g", m, 2.0/7.0)
	}
}

// ---------------------------------------------------------------------------
// Gamma
// ---------------------------------------------------------------------------

func TestGamma(t *testing.T) {
	rng := NewRNG(42)
	a := rng.Gamma(2.0, 1.0, 1000)
	if a.Size() != 1000 {
		t.Fatalf("Gamma: expected size 1000, got %d", a.Size())
	}
	for _, v := range a.Data() {
		if v < 0 {
			t.Fatalf("Gamma: got negative value %g", v)
		}
	}
	// Mean of Gamma(2, 1) = 2.
	m := Mean(a).Data()[0]
	if math.Abs(m-2.0) > 0.3 {
		t.Fatalf("Gamma: mean %g not close to 2.0", m)
	}
}

// ---------------------------------------------------------------------------
// Chisquare
// ---------------------------------------------------------------------------

func TestChisquare(t *testing.T) {
	rng := NewRNG(42)
	a := rng.Chisquare(5.0, 1000)
	if a.Size() != 1000 {
		t.Fatalf("Chisquare: expected size 1000, got %d", a.Size())
	}
	// Mean of chi-squared(df=5) = 5.
	m := Mean(a).Data()[0]
	if math.Abs(m-5.0) > 1.0 {
		t.Fatalf("Chisquare: mean %g not close to 5.0", m)
	}
}

// ---------------------------------------------------------------------------
// StandardNormal
// ---------------------------------------------------------------------------

func TestStandardNormal(t *testing.T) {
	rng := NewRNG(42)
	a := rng.StandardNormal(1000)
	m := Mean(a).Data()[0]
	if math.Abs(m) > 0.15 {
		t.Fatalf("StandardNormal: mean %g not close to 0", m)
	}
}

// ---------------------------------------------------------------------------
// StandardT
// ---------------------------------------------------------------------------

func TestStandardT(t *testing.T) {
	rng := NewRNG(42)
	a := rng.StandardT(30, 1000)
	if a.Size() != 1000 {
		t.Fatalf("StandardT: expected size 1000, got %d", a.Size())
	}
	// With high df, mean should be near 0.
	m := Mean(a).Data()[0]
	if math.Abs(m) > 0.2 {
		t.Fatalf("StandardT: mean %g not close to 0", m)
	}
}
