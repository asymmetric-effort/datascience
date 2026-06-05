//go:build unit

package scigo

import (
	"math"
	"testing"
)

func approxEqualQP(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

// ---------------------------------------------------------------------------
// QPSolve Tests
// ---------------------------------------------------------------------------

func TestQPSolve_SimpleUnconstrained(t *testing.T) {
	// min 0.5*(2x1^2 + 2x2^2) = min x1^2 + x2^2
	// Q = [[2,0],[0,2]], c = [0,0]
	// Solution: x = [0,0], fun = 0
	Q := []float64{2, 0, 0, 2}
	c := []float64{0, 0}
	result, err := QPSolve(Q, c, 2, nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Success {
		t.Error("expected Success=true")
	}
	if !approxEqualQP(result.X[0], 0, 1e-6) || !approxEqualQP(result.X[1], 0, 1e-6) {
		t.Errorf("expected [0,0], got %v", result.X)
	}
	if !approxEqualQP(result.Fun, 0, 1e-8) {
		t.Errorf("expected fun=0, got %v", result.Fun)
	}
}

func TestQPSolve_WithLinearTerm(t *testing.T) {
	// min 0.5*(2x^2) + (-2)x = x^2 - 2x
	// Derivative: 2x - 2 = 0 => x = 1, fun = -1
	Q := []float64{2}
	c := []float64{-2}
	result, err := QPSolve(Q, c, 1, nil, nil, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(result.X[0], 1, 1e-6) {
		t.Errorf("expected x=1, got %v", result.X[0])
	}
	if !approxEqualQP(result.Fun, -1, 1e-6) {
		t.Errorf("expected fun=-1, got %v", result.Fun)
	}
}

func TestQPSolve_EqualityConstraint(t *testing.T) {
	// min 0.5*(x1^2 + x2^2) s.t. x1 + x2 = 1
	// Solution: x1 = x2 = 0.5, fun = 0.25
	Q := []float64{1, 0, 0, 1}
	c := []float64{0, 0}
	Aeq := [][]float64{{1, 1}}
	beq := []float64{1}

	result, err := QPSolve(Q, c, 2, nil, nil, Aeq, beq, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !result.Success {
		t.Error("expected Success=true")
	}
	if !approxEqualQP(result.X[0], 0.5, 1e-6) || !approxEqualQP(result.X[1], 0.5, 1e-6) {
		t.Errorf("expected [0.5,0.5], got %v", result.X)
	}
	if !approxEqualQP(result.Fun, 0.25, 1e-6) {
		t.Errorf("expected fun=0.25, got %v", result.Fun)
	}
}

func TestQPSolve_InequalityConstraint(t *testing.T) {
	// min 0.5*(2x^2) - 3x  =>  x^2 - 3x  => unconstrained min at x=1.5
	// s.t. x <= 1
	// Solution: x = 1, fun = 1 - 3 = -2
	Q := []float64{2}
	c := []float64{-3}
	A := [][]float64{{1}}
	b := []float64{1}

	result, err := QPSolve(Q, c, 1, A, b, nil, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(result.X[0], 1.0, 1e-4) {
		t.Errorf("expected x=1.0, got %v", result.X[0])
	}
	if !approxEqualQP(result.Fun, -2.0, 1e-4) {
		t.Errorf("expected fun=-2.0, got %v", result.Fun)
	}
}

func TestQPSolve_BoundConstraints(t *testing.T) {
	// min 0.5*(2x^2) - 5x => x^2 - 5x => unconstrained min at x=2.5
	// lb=0, ub=2 => solution x=2, fun = 4-10 = -6
	Q := []float64{2}
	c := []float64{-5}
	lb := []float64{0}
	ub := []float64{2}

	result, err := QPSolve(Q, c, 1, nil, nil, nil, nil, lb, ub)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(result.X[0], 2.0, 1e-4) {
		t.Errorf("expected x=2.0, got %v", result.X[0])
	}
	if !approxEqualQP(result.Fun, -6.0, 1e-4) {
		t.Errorf("expected fun=-6.0, got %v", result.Fun)
	}
}

func TestQPSolve_2D_WithBounds(t *testing.T) {
	// min 0.5*(x1^2 + x2^2) - x1 - x2  s.t. 0 <= x1,x2 <= 0.5
	// Unconstrained: x1=x2=1, but bounded to 0.5.
	Q := []float64{1, 0, 0, 1}
	c := []float64{-1, -1}
	lb := []float64{0, 0}
	ub := []float64{0.5, 0.5}

	result, err := QPSolve(Q, c, 2, nil, nil, nil, nil, lb, ub)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(result.X[0], 0.5, 1e-4) || !approxEqualQP(result.X[1], 0.5, 1e-4) {
		t.Errorf("expected [0.5,0.5], got %v", result.X)
	}
}

func TestQPSolve_SimplexConstraint(t *testing.T) {
	// min 0.5 * w' * [[2, 1],[1, 2]] * w  s.t. w1+w2=1, w>=0
	// Analytical: w = [0.5, 0.5], fun = 0.5*(0.5+0.5+0.5) = 0.75
	Q := []float64{2, 1, 1, 2}
	c := []float64{0, 0}
	Aeq := [][]float64{{1, 1}}
	beq := []float64{1}
	lb := []float64{0, 0}
	ub := []float64{1, 1}

	result, err := QPSolve(Q, c, 2, nil, nil, Aeq, beq, lb, ub)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(result.X[0], 0.5, 1e-3) || !approxEqualQP(result.X[1], 0.5, 1e-3) {
		t.Errorf("expected [0.5,0.5], got %v", result.X)
	}
}

// ---------------------------------------------------------------------------
// QPSolve Error Cases
// ---------------------------------------------------------------------------

func TestQPSolve_InvalidDimension(t *testing.T) {
	_, err := QPSolve(nil, nil, 0, nil, nil, nil, nil, nil, nil)
	if err == nil {
		t.Error("expected error for Qn=0")
	}
}

func TestQPSolve_QWrongLength(t *testing.T) {
	_, err := QPSolve([]float64{1, 2}, []float64{0, 0}, 2, nil, nil, nil, nil, nil, nil)
	if err == nil {
		t.Error("expected error for wrong Q length")
	}
}

func TestQPSolve_CWrongLength(t *testing.T) {
	_, err := QPSolve([]float64{1, 0, 0, 1}, []float64{0}, 2, nil, nil, nil, nil, nil, nil)
	if err == nil {
		t.Error("expected error for wrong c length")
	}
}

func TestQPSolve_ABMismatch(t *testing.T) {
	_, err := QPSolve([]float64{1}, []float64{0}, 1,
		[][]float64{{1}}, []float64{1, 2}, nil, nil, nil, nil)
	if err == nil {
		t.Error("expected error for A/b mismatch")
	}
}

func TestQPSolve_AeqBeqMismatch(t *testing.T) {
	_, err := QPSolve([]float64{1}, []float64{0}, 1,
		nil, nil, [][]float64{{1}}, []float64{1, 2}, nil, nil)
	if err == nil {
		t.Error("expected error for Aeq/beq mismatch")
	}
}

func TestQPSolve_ARowWrongLength(t *testing.T) {
	_, err := QPSolve([]float64{1, 0, 0, 1}, []float64{0, 0}, 2,
		[][]float64{{1}}, []float64{1}, nil, nil, nil, nil)
	if err == nil {
		t.Error("expected error for A row wrong length")
	}
}

func TestQPSolve_AeqRowWrongLength(t *testing.T) {
	_, err := QPSolve([]float64{1, 0, 0, 1}, []float64{0, 0}, 2,
		nil, nil, [][]float64{{1}}, []float64{1}, nil, nil)
	if err == nil {
		t.Error("expected error for Aeq row wrong length")
	}
}

func TestQPSolve_LbWrongLength(t *testing.T) {
	_, err := QPSolve([]float64{1, 0, 0, 1}, []float64{0, 0}, 2,
		nil, nil, nil, nil, []float64{0}, nil)
	if err == nil {
		t.Error("expected error for lb wrong length")
	}
}

func TestQPSolve_UbWrongLength(t *testing.T) {
	_, err := QPSolve([]float64{1, 0, 0, 1}, []float64{0, 0}, 2,
		nil, nil, nil, nil, nil, []float64{1})
	if err == nil {
		t.Error("expected error for ub wrong length")
	}
}

func TestQPSolve_InfBounds(t *testing.T) {
	// Test with infinite bounds (should be ignored).
	Q := []float64{2}
	c := []float64{-2}
	lb := []float64{math.Inf(-1)}
	ub := []float64{math.Inf(1)}

	result, err := QPSolve(Q, c, 1, nil, nil, nil, nil, lb, ub)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(result.X[0], 1.0, 1e-6) {
		t.Errorf("expected x=1, got %v", result.X[0])
	}
}

func TestQPSolve_MultipleInequalityConstraints(t *testing.T) {
	// min 0.5*(x1^2 + x2^2)  s.t. x1 + x2 >= 1 (i.e., -x1 - x2 <= -1)
	// and x1 >= 0, x2 >= 0
	// Solution: x1 = x2 = 0.5
	Q := []float64{1, 0, 0, 1}
	c := []float64{0, 0}
	A := [][]float64{{-1, -1}}
	b := []float64{-1}
	lb := []float64{0, 0}

	result, err := QPSolve(Q, c, 2, A, b, nil, nil, lb, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(result.X[0], 0.5, 1e-3) || !approxEqualQP(result.X[1], 0.5, 1e-3) {
		t.Errorf("expected [0.5,0.5], got %v", result.X)
	}
}

// ---------------------------------------------------------------------------
// itoa tests
// ---------------------------------------------------------------------------

func TestItoa(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "0"},
		{1, "1"},
		{-5, "-5"},
		{123, "123"},
	}
	for _, tt := range tests {
		got := itoa(tt.n)
		if got != tt.want {
			t.Errorf("itoa(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestDotVec(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 5, 6}
	got := dotVec(a, b)
	want := 32.0
	if !approxEqualQP(got, want, 1e-10) {
		t.Errorf("dotVec = %v, want %v", got, want)
	}
}

func TestQPObjective(t *testing.T) {
	Q := []float64{2, 0, 0, 2}
	c := []float64{-1, -1}
	x := []float64{1, 1}
	// 0.5 * (2+2) + (-1-1) = 2 - 2 = 0
	got := qpObjective(Q, c, x, 2)
	if !approxEqualQP(got, 0, 1e-10) {
		t.Errorf("qpObjective = %v, want 0", got)
	}
}

func TestQPSolve_3DWithEqAndIneq(t *testing.T) {
	// min 0.5*||x||^2 s.t. x1+x2+x3=1, x1 <= 0.4, all >= 0
	// Without x1<=0.4: x=[1/3,1/3,1/3]. With it: x1=0.4 would not bind since 1/3 < 0.4.
	// So solution is still [1/3,1/3,1/3].
	Q := []float64{1, 0, 0, 0, 1, 0, 0, 0, 1}
	c := []float64{0, 0, 0}
	A := [][]float64{{1, 0, 0}} // x1 <= 0.4
	b := []float64{0.4}
	Aeq := [][]float64{{1, 1, 1}}
	beq := []float64{1}
	lb := []float64{0, 0, 0}

	result, err := QPSolve(Q, c, 3, A, b, Aeq, beq, lb, nil)
	if err != nil {
		t.Fatal(err)
	}
	third := 1.0 / 3.0
	for i := 0; i < 3; i++ {
		if !approxEqualQP(result.X[i], third, 1e-2) {
			t.Errorf("X[%d] = %v, want ~%v", i, result.X[i], third)
		}
	}
}

func TestSolveLinearSystem(t *testing.T) {
	A := [][]float64{{2, 1}, {1, 3}}
	b := []float64{5, 7}
	// 2x+y=5, x+3y=7 => x=1.6, y=1.8
	x, err := solveLinearSystem(A, b)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(x[0], 1.6, 1e-10) || !approxEqualQP(x[1], 1.8, 1e-10) {
		t.Errorf("got %v, want [1.6, 1.8]", x)
	}
}

func TestQPFindFeasible(t *testing.T) {
	// Simple: x >= 1 as -x <= -1
	x := qpFindFeasible([]float64{0}, 1,
		[][]float64{{-1}}, []float64{-1},
		nil, nil)
	if x[0] < 0.99 {
		t.Errorf("expected feasible x >= 1, got %v", x[0])
	}
}

func TestQPSolveEqualityOnly_Singular(t *testing.T) {
	// Singular KKT system: Q is zero, Aeq is zero row.
	Q := []float64{0, 0, 0, 0}
	c := []float64{0, 0}
	Aeq := [][]float64{{0, 0}}
	beq := []float64{1}

	_, err := QPSolve(Q, c, 2, nil, nil, Aeq, beq, nil, nil)
	if err == nil {
		t.Error("expected error for singular KKT")
	}
}

func TestQPSolve_EqualityOnlyDirect(t *testing.T) {
	// min 0.5 * (x1^2 + x2^2) s.t. x1 - x2 = 0
	// Solution: x = [0, 0]
	Q := []float64{1, 0, 0, 1}
	c := []float64{0, 0}
	Aeq := [][]float64{{1, -1}}
	beq := []float64{0}

	result, err := QPSolve(Q, c, 2, nil, nil, Aeq, beq, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !approxEqualQP(result.X[0], 0, 1e-8) || !approxEqualQP(result.X[1], 0, 1e-8) {
		t.Errorf("expected [0,0], got %v", result.X)
	}
	if result.Iterations != 1 {
		t.Errorf("equality-only should take 1 iteration, got %d", result.Iterations)
	}
}
