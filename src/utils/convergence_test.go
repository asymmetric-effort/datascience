//go:build unit

package utils

import "testing"

func TestConvergenceCheckerBasic(t *testing.T) {
	cc := NewConvergenceChecker(0.01, 1000)

	// First update never converges (no previous value).
	if cc.Update(10.0) {
		t.Error("should not converge on first update")
	}
	if cc.Iterations() != 1 {
		t.Errorf("iterations = %d, want 1", cc.Iterations())
	}

	// Large change: not converged.
	if cc.Update(5.0) {
		t.Error("should not converge with large change")
	}

	// Small change: converged.
	if !cc.Update(5.005) {
		t.Error("should converge with change < tolerance")
	}
	if !cc.Converged() {
		t.Error("Converged() should be true")
	}
	if cc.Iterations() != 3 {
		t.Errorf("iterations = %d, want 3", cc.Iterations())
	}
}

func TestConvergenceCheckerMaxIter(t *testing.T) {
	cc := NewConvergenceChecker(1e-15, 3)

	cc.Update(1.0)
	cc.Update(2.0)
	// Third update hits maxIter.
	if !cc.Update(3.0) {
		t.Error("should converge at max iterations")
	}
	if !cc.Converged() {
		t.Error("Converged() should be true at max iterations")
	}
}

func TestConvergenceCheckerReset(t *testing.T) {
	cc := NewConvergenceChecker(0.01, 1000)
	cc.Update(1.0)
	cc.Update(1.0001)
	if !cc.Converged() {
		t.Fatal("should be converged")
	}

	cc.Reset()
	if cc.Converged() {
		t.Error("should not be converged after reset")
	}
	if cc.Iterations() != 0 {
		t.Errorf("iterations = %d, want 0 after reset", cc.Iterations())
	}

	// First update after reset should not converge.
	if cc.Update(100.0) {
		t.Error("should not converge on first update after reset")
	}
}

func TestConvergenceCheckerExactTolerance(t *testing.T) {
	cc := NewConvergenceChecker(0.1, 1000)
	cc.Update(1.0)
	// Change of exactly 0.1 is NOT less than tolerance, so no convergence.
	if cc.Update(1.1) {
		t.Error("change == tolerance should not converge (requires strictly less than)")
	}
	// Change of 0.09 is less than tolerance.
	if !cc.Update(1.19) {
		t.Error("change < tolerance should converge")
	}
}

func TestConvergenceCheckerStaysConverged(t *testing.T) {
	cc := NewConvergenceChecker(0.01, 1000)
	cc.Update(1.0)
	cc.Update(1.001) // converges

	// Even with a large change after convergence, Update still returns true.
	if !cc.Update(999.0) {
		t.Error("once converged, should stay converged")
	}
}
