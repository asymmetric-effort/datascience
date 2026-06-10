//go:build unit

package metrics

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// betaCFImpl — DI tests with large tiny threshold to trigger underflow guards
// ---------------------------------------------------------------------------

// TestBetaCFImpl_InitDTiny exercises the initial d underflow guard (line 82-84)
// by injecting a large tiny threshold that forces |d| < tiny on the initial
// computation.
func TestBetaCFImpl_InitDTiny(t *testing.T) {
	// With tiny = 1e10, any d value with |d| < 1e10 triggers the guard.
	// d_init = 1 - (a+b)*x/(a+1). For a=1, b=1, x=0.5:
	// d_init = 1 - 2*0.5/2 = 1 - 0.5 = 0.5, which is < 1e10 => triggers guard.
	result := betaCFImpl(0.5, 1.0, 1.0, 1e10)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		t.Errorf("expected finite result, got %f", result)
	}
}

// TestBetaCFImpl_LoopDTiny exercises the even-step and odd-step d underflow
// guards by injecting a large tiny threshold.
func TestBetaCFImpl_LoopDTiny(t *testing.T) {
	// With tiny = 1e10, all d and c values in the CF loop will trigger the
	// underflow guard since they are typically small numbers.
	result := betaCFImpl(0.3, 2.0, 3.0, 1e10)
	if math.IsNaN(result) {
		t.Errorf("expected non-NaN result, got NaN")
	}
}

// TestBetaCFImpl_LoopCTiny exercises the even-step and odd-step c underflow
// guards by injecting a large tiny threshold.
func TestBetaCFImpl_LoopCTiny(t *testing.T) {
	// Multiple parameter sets to ensure all CF loop guards are hit.
	testCases := []struct {
		x, a, b float64
	}{
		{0.1, 0.5, 0.5},
		{0.3, 1.0, 2.0},
		{0.01, 5.0, 10.0},
		{0.2, 0.1, 0.1},
		{0.4, 3.0, 7.0},
	}
	for _, tc := range testCases {
		result := betaCFImpl(tc.x, tc.a, tc.b, 1e10)
		if math.IsNaN(result) {
			t.Errorf("NaN for x=%v, a=%v, b=%v with large tiny", tc.x, tc.a, tc.b)
		}
	}
}

// TestBetaCFImpl_SymmetryWithLargeTiny exercises the symmetry relation path
// when the recursive call also uses the large tiny threshold.
func TestBetaCFImpl_SymmetryWithLargeTiny(t *testing.T) {
	// x > (a+1)/(a+b+2) triggers symmetry. With a=1, b=10: threshold = 2/13 ~ 0.154.
	// x=0.2 > 0.154 => symmetry path taken, recursive call also uses large tiny.
	result := betaCFImpl(0.2, 1.0, 10.0, 1e10)
	if math.IsNaN(result) {
		t.Errorf("expected non-NaN result, got NaN")
	}
}

// TestBetaCFImpl_AllGuardsTriggered uses a moderately large tiny value to
// ensure every guard in the CF loop fires for multiple iterations.
func TestBetaCFImpl_AllGuardsTriggered(t *testing.T) {
	// Use tiny = 1e5 to force guards while keeping numerical stability.
	result := betaCFImpl(0.25, 3.0, 5.0, 1e5)
	if math.IsNaN(result) {
		t.Errorf("expected non-NaN result, got NaN")
	}
}

// TestBetaCFImpl_BoundaryValues exercises the boundary conditions with large tiny.
func TestBetaCFImpl_BoundaryValues(t *testing.T) {
	// x <= 0 returns 0.
	if betaCFImpl(0, 1, 1, 1e10) != 0 {
		t.Error("expected 0 for x=0")
	}
	if betaCFImpl(-1, 1, 1, 1e10) != 0 {
		t.Error("expected 0 for x<0")
	}
	// x >= 1 returns 1.
	if betaCFImpl(1, 1, 1, 1e10) != 1 {
		t.Error("expected 1 for x=1")
	}
}

// ---------------------------------------------------------------------------
// upperGammaCFImpl — DI tests with large tiny threshold to trigger guards
// ---------------------------------------------------------------------------

// TestUpperGammaCFImpl_DTiny exercises the inner loop d underflow guard
// by injecting a large tiny threshold.
func TestUpperGammaCFImpl_DTiny(t *testing.T) {
	// With tiny = 1e10, |d| will almost always be < tiny, triggering the guard.
	result := upperGammaCFImpl(2.0, 5.0, 1e10)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		t.Errorf("expected finite result, got %f", result)
	}
}

// TestUpperGammaCFImpl_CTiny exercises the inner loop c underflow guard
// by injecting a large tiny threshold.
func TestUpperGammaCFImpl_CTiny(t *testing.T) {
	testCases := []struct {
		a, x float64
	}{
		{2.0, 5.0},
		{0.5, 1.0},
		{10.0, 8.0},
		{1.0, 0.5},
		{5.0, 3.0},
	}
	for _, tc := range testCases {
		result := upperGammaCFImpl(tc.a, tc.x, 1e10)
		if math.IsNaN(result) {
			t.Errorf("NaN for a=%v, x=%v with large tiny", tc.a, tc.x)
		}
	}
}

// TestUpperGammaCFImpl_B0Tiny exercises the b0 underflow guard where
// b0 = x + 1 - a is very small relative to the tiny threshold.
func TestUpperGammaCFImpl_B0Tiny(t *testing.T) {
	// b0 = x + 1 - a. For a=2, x=1: b0 = 0 => |b0| < any positive tiny.
	result := upperGammaCFImpl(2.0, 1.0, 1e10)
	if math.IsNaN(result) || math.IsInf(result, 0) {
		t.Errorf("expected finite result, got %f", result)
	}
}

// TestUpperGammaCFImpl_AllGuardsTriggered exercises all guards with a
// moderately large tiny value.
func TestUpperGammaCFImpl_AllGuardsTriggered(t *testing.T) {
	result := upperGammaCFImpl(3.0, 2.0, 1e5)
	if math.IsNaN(result) {
		t.Errorf("expected non-NaN result, got NaN")
	}
}

// ---------------------------------------------------------------------------
// Verify correctness: results with tiny=1e-30 should match original behavior
// ---------------------------------------------------------------------------

// TestBetaCFImpl_MatchesOriginal verifies that betaCFImpl with the default
// tiny threshold produces the same results as the original implementation.
func TestBetaCFImpl_MatchesOriginal(t *testing.T) {
	testCases := []struct {
		x, a, b float64
	}{
		{0.3, 2, 3},
		{0.5, 1, 1},
		{0.1, 5, 10},
		{0.9, 0.5, 0.5},
	}
	for _, tc := range testCases {
		got := betaCFImpl(tc.x, tc.a, tc.b, 1e-30)
		want := regularizedIncompleteBeta(tc.x, tc.a, tc.b)
		if math.Abs(got-want) > 1e-12 {
			t.Errorf("betaCFImpl(%v,%v,%v,1e-30) = %v, regularizedIncompleteBeta = %v",
				tc.x, tc.a, tc.b, got, want)
		}
	}
}

// TestUpperGammaCFImpl_MatchesOriginal verifies that upperGammaCFImpl with
// the default tiny threshold produces the same results as the original.
func TestUpperGammaCFImpl_MatchesOriginal(t *testing.T) {
	testCases := []struct {
		a, x float64
	}{
		{2.0, 5.0},
		{5.0, 3.0},
		{1.0, 1.0},
		{10.0, 8.0},
	}
	for _, tc := range testCases {
		got := upperGammaCFImpl(tc.a, tc.x, 1e-30)
		want := upperGammaCF(tc.a, tc.x)
		if math.Abs(got-want) > 1e-12 {
			t.Errorf("upperGammaCFImpl(%v,%v,1e-30) = %v, upperGammaCF = %v",
				tc.a, tc.x, got, want)
		}
	}
}
