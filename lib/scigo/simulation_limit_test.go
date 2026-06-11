//go:build unit

package scigo

import (
	"strings"
	"testing"
)

// --- checkSimulationSize tests ---

func TestCheckSimulationSizeAcceptsWithinLimit(t *testing.T) {
	// Should not panic.
	checkSimulationSize(100, 100, 10001)
	checkSimulationSize(1, 1, 1)
	checkSimulationSize(10000, 10000, MaxSimulationElements)
}

func TestCheckSimulationSizePanicsOverLimit(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for oversized simulation")
		}
		msg := r.(string)
		if !strings.Contains(msg, "exceeds limit") {
			t.Errorf("unexpected panic message: %s", msg)
		}
	}()
	checkSimulationSize(100_000, 10_000, 100_000) // 1 billion > 100K
}

func TestCheckSimulationSizePanicsOnOverflow(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for integer overflow")
		}
		msg := r.(string)
		if !strings.Contains(msg, "overflows") {
			t.Errorf("unexpected panic message: %s", msg)
		}
	}()
	checkSimulationSize(1<<40, 1<<40, MaxSimulationElements)
}

// --- Default limit enforcement on public functions ---

func TestEulerMaruyamaPanicsOnOversizedSimulation(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("EulerMaruyama should panic on oversized simulation")
		}
	}()
	drift := func(x, t float64) float64 { return 0 }
	diff := func(x, t float64) float64 { return 1 }
	// 1M paths * 1000 steps = 1B elements > MaxSimulationElements
	EulerMaruyama(drift, diff, 0, [2]float64{0, 1}, 0.001, 1_000_000, 42)
}

func TestBrownianMotionPanicsOnOversizedSimulation(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("BrownianMotion should panic on oversized simulation")
		}
	}()
	BrownianMotion(1.0, 1_000_000, 1_000, 42)
}

// --- WithLimit variants allow larger simulations ---

func TestEulerMaruyamaWithLimitAllowsLarger(t *testing.T) {
	drift := func(x, t float64) float64 { return 0 }
	diff := func(x, t float64) float64 { return 0 }
	// 10 paths * 11 steps = 110 elements, default would allow this too,
	// but test the WithLimit plumbing works.
	result := EulerMaruyamaWithLimit(drift, diff, 0, [2]float64{0, 1}, 0.1, 10, 42, 200)
	if len(result.X) != 10 {
		t.Errorf("expected 10 paths, got %d", len(result.X))
	}
}

func TestEulerMaruyamaWithLimitRejectsOverCustom(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic when exceeding custom limit")
		}
	}()
	drift := func(x, t float64) float64 { return 0 }
	diff := func(x, t float64) float64 { return 0 }
	// 10 paths * 11 steps = 110 elements > limit of 50
	EulerMaruyamaWithLimit(drift, diff, 0, [2]float64{0, 1}, 0.1, 10, 42, 50)
}

func TestBrownianMotionWithLimitAllowsLarger(t *testing.T) {
	result := BrownianMotionWithLimit(1.0, 100, 5, 42, 1000)
	if len(result.X) != 5 {
		t.Errorf("expected 5 paths, got %d", len(result.X))
	}
}

func TestMilsteinWithLimitWorks(t *testing.T) {
	drift := func(x, t float64) float64 { return 0 }
	diff := func(x, t float64) float64 { return 1 }
	diffD := func(x, t float64) float64 { return 0 }
	result := MilsteinWithLimit(drift, diff, diffD, 0, [2]float64{0, 1}, 0.5, 2, 42, 100)
	if len(result.X) != 2 {
		t.Errorf("expected 2 paths, got %d", len(result.X))
	}
}

func TestGBMWithLimitWorks(t *testing.T) {
	result := GeometricBrownianMotionWithLimit(100, 0.05, 0.2, 1.0, 10, 3, 42, 100)
	if len(result.X) != 3 {
		t.Errorf("expected 3 paths, got %d", len(result.X))
	}
}

func TestOUWithLimitWorks(t *testing.T) {
	result := OrnsteinUhlenbeckWithLimit(0, 1, 0, 0.5, 1.0, 10, 3, 42, 100)
	if len(result.X) != 3 {
		t.Errorf("expected 3 paths, got %d", len(result.X))
	}
}

func TestBrownianBridgeWithLimitWorks(t *testing.T) {
	bridge := BrownianBridgeWithLimit(1.0, 10, 0, 1, 42, 100)
	if len(bridge) != 11 {
		t.Errorf("expected 11 points, got %d", len(bridge))
	}
}

// --- MonteCarloPricing has no limit (O(1) memory) ---

func TestMonteCarloPricingNoLimitNeeded(t *testing.T) {
	// Even with huge nPaths, it's O(1) memory — just a running sum.
	// Use a modest number to keep the test fast.
	payoff := func(sT float64) float64 {
		if sT > 100 {
			return sT - 100
		}
		return 0
	}
	price := MonteCarloPricing(100, 100, 1.0, 0.05, 0.2, 10000, 42, payoff)
	if price <= 0 {
		t.Errorf("expected positive price, got %f", price)
	}
}
