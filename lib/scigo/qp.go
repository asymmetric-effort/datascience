package scigo

import (
	"errors"
	"math"
)

// QPResult holds the result of a quadratic programming solver.
type QPResult struct {
	X          []float64 // Solution vector
	Fun        float64   // Objective function value at the solution
	Success    bool      // Whether the solver converged
	Iterations int       // Number of iterations performed
}

// QPSolve solves a quadratic programming problem of the form:
//
//	minimize    0.5 * x^T * Q * x + c^T * x
//	subject to  A * x <= b        (inequality constraints)
//	            Aeq * x = beq     (equality constraints)
//	            lb <= x <= ub     (bound constraints)
//
// Q is stored as a flat array of length Qn*Qn, where Qn is the dimension.
// A is a slice of inequality constraint rows (each of length Qn); b is the RHS.
// Aeq/beq are equality constraints. lb/ub are per-variable bounds (nil = no bounds).
//
// Uses an active-set method. When no inequality constraints are present,
// the KKT system is solved directly.
func QPSolve(Q, c []float64, Qn int, A [][]float64, b []float64,
	Aeq [][]float64, beq []float64, lb, ub []float64) (*QPResult, error) {

	if Qn <= 0 {
		return nil, errors.New("scigo.QPSolve: Qn must be positive")
	}
	if len(Q) != Qn*Qn {
		return nil, errors.New("scigo.QPSolve: Q length must equal Qn*Qn")
	}
	if len(c) != Qn {
		return nil, errors.New("scigo.QPSolve: c length must equal Qn")
	}
	if len(A) != len(b) {
		return nil, errors.New("scigo.QPSolve: A and b must have same length")
	}
	for i, row := range A {
		if len(row) != Qn {
			return nil, errors.New("scigo.QPSolve: A row " + itoa(i) + " has wrong length")
		}
	}
	if len(Aeq) != len(beq) {
		return nil, errors.New("scigo.QPSolve: Aeq and beq must have same length")
	}
	for i, row := range Aeq {
		if len(row) != Qn {
			return nil, errors.New("scigo.QPSolve: Aeq row " + itoa(i) + " has wrong length")
		}
	}
	if lb != nil && len(lb) != Qn {
		return nil, errors.New("scigo.QPSolve: lb length must equal Qn")
	}
	if ub != nil && len(ub) != Qn {
		return nil, errors.New("scigo.QPSolve: ub length must equal Qn")
	}

	// Convert bound constraints to inequality constraints: x_i <= ub_i  =>  e_i * x <= ub_i
	// and -x_i <= -lb_i  =>  -e_i * x <= -lb_i
	allA := make([][]float64, 0, len(A))
	allB := make([]float64, 0, len(b))
	allA = append(allA, A...)
	allB = append(allB, b...)

	if ub != nil {
		for i := 0; i < Qn; i++ {
			if !math.IsInf(ub[i], 1) {
				row := make([]float64, Qn)
				row[i] = 1.0
				allA = append(allA, row)
				allB = append(allB, ub[i])
			}
		}
	}
	if lb != nil {
		for i := 0; i < Qn; i++ {
			if !math.IsInf(lb[i], -1) {
				row := make([]float64, Qn)
				row[i] = -1.0
				allA = append(allA, row)
				allB = append(allB, -lb[i])
			}
		}
	}

	nIneq := len(allA)
	nEq := len(Aeq)

	// If no inequality constraints, solve equality-constrained QP via KKT.
	if nIneq == 0 {
		return qpSolveEqualityOnly(Q, c, Qn, Aeq, beq)
	}

	// Active-set method.
	return qpActiveSet(Q, c, Qn, allA, allB, Aeq, beq, nIneq, nEq)
}

// qpSolveEqualityOnly solves min 0.5*x'Qx+c'x s.t. Aeq*x=beq via the KKT system.
func qpSolveEqualityOnly(Q, c []float64, n int, Aeq [][]float64, beq []float64) (*QPResult, error) {
	nEq := len(Aeq)
	dim := n + nEq

	// Build KKT matrix:
	// [ Q    Aeq^T ] [ x ]   = [ -c  ]
	// [ Aeq  0     ] [ lam]     [ beq ]
	kkt := make([][]float64, dim)
	for i := 0; i < dim; i++ {
		kkt[i] = make([]float64, dim)
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			kkt[i][j] = Q[i*n+j]
		}
	}
	for k := 0; k < nEq; k++ {
		for j := 0; j < n; j++ {
			kkt[j][n+k] = Aeq[k][j]
			kkt[n+k][j] = Aeq[k][j]
		}
	}

	rhs := make([]float64, dim)
	for i := 0; i < n; i++ {
		rhs[i] = -c[i]
	}
	for k := 0; k < nEq; k++ {
		rhs[n+k] = beq[k]
	}

	sol, err := solveLinearSystem(kkt, rhs)
	if err != nil {
		return nil, errors.New("scigo.QPSolve: KKT system is singular")
	}

	x := sol[:n]
	fun := qpObjective(Q, c, x, n)

	return &QPResult{X: x, Fun: fun, Success: true, Iterations: 1}, nil
}

// qpActiveSet implements the active-set method for inequality-constrained QP.
func qpActiveSet(Q, c []float64, n int, A [][]float64, b []float64,
	Aeq [][]float64, beq []float64, nIneq, nEq int) (*QPResult, error) {

	const maxIter = 10000

	// Start with a feasible point. Try x=0 first; if infeasible, find one.
	x := make([]float64, n)

	// Try to find a feasible starting point.
	x = qpFindFeasible(x, n, A, b, Aeq, beq)

	// Active set: indices of inequality constraints currently treated as equalities.
	active := make(map[int]bool)

	// Initialize active set: constraints that are (nearly) tight at x.
	for i := 0; i < nIneq; i++ {
		val := dotVec(A[i], x)
		if val >= b[i]-1e-10 {
			active[i] = true
		}
	}

	var iter int
	for iter = 0; iter < maxIter; iter++ {
		// Build working equality constraints: Aeq plus active inequalities.
		nActive := len(active)
		nWork := nEq + nActive
		workA := make([][]float64, nWork)
		workB := make([]float64, nWork)

		idx := 0
		for i := 0; i < nEq; i++ {
			workA[idx] = Aeq[i]
			workB[idx] = beq[i]
			idx++
		}
		activeList := make([]int, 0, nActive)
		for k := range active {
			activeList = append(activeList, k)
		}
		// Sort for determinism.
		for i := 0; i < len(activeList); i++ {
			for j := i + 1; j < len(activeList); j++ {
				if activeList[j] < activeList[i] {
					activeList[i], activeList[j] = activeList[j], activeList[i]
				}
			}
		}
		for _, k := range activeList {
			workA[idx] = A[k]
			workB[idx] = b[k]
			idx++
		}

		// Solve the equality-constrained subproblem for the step p:
		// min 0.5*(x+p)'Q(x+p) + c'(x+p) s.t. workA*(x+p) = workB
		// Equivalently: min 0.5*p'Qp + g'p s.t. workA*p = 0
		// where g = Qx + c
		g := make([]float64, n)
		for i := 0; i < n; i++ {
			g[i] = c[i]
			for j := 0; j < n; j++ {
				g[i] += Q[i*n+j] * x[j]
			}
		}

		p, lambdas, err := qpSolveEqSubproblem(Q, g, n, workA, nWork)
		if err != nil {
			// If we can't solve the subproblem, return current x.
			break
		}

		pNorm := 0.0
		for _, v := range p {
			pNorm += v * v
		}
		pNorm = math.Sqrt(pNorm)

		if pNorm < 1e-10 {
			// p is essentially zero; check multipliers for active inequality constraints.
			// If all multipliers for active inequality constraints are >= 0, we're optimal.
			allNonNeg := true
			worstLam := 0.0
			worstIdx := -1
			for i, k := range activeList {
				lam := lambdas[nEq+i]
				if lam < worstLam {
					worstLam = lam
					worstIdx = k
					allNonNeg = false
				}
			}
			if allNonNeg {
				// Optimal.
				break
			}
			// Remove the most negative multiplier constraint from active set.
			delete(active, worstIdx)
		} else {
			// Compute step length alpha: largest alpha in [0,1] s.t. all inactive constraints remain feasible.
			alpha := 1.0
			blockingConstraint := -1
			for i := 0; i < nIneq; i++ {
				if active[i] {
					continue
				}
				ai_p := dotVec(A[i], p)
				if ai_p > 1e-14 {
					slack := b[i] - dotVec(A[i], x)
					maxAlpha := slack / ai_p
					if maxAlpha < alpha {
						alpha = maxAlpha
						blockingConstraint = i
					}
				}
			}

			if alpha < 0 {
				alpha = 0
			}

			// Take the step.
			for i := 0; i < n; i++ {
				x[i] += alpha * p[i]
			}

			if blockingConstraint >= 0 && alpha < 1.0-1e-12 {
				active[blockingConstraint] = true
			}
		}
	}

	fun := qpObjective(Q, c, x, n)
	return &QPResult{X: x, Fun: fun, Success: iter < maxIter, Iterations: iter}, nil
}

// qpSolveEqSubproblem solves: min 0.5*p'Qp + g'p s.t. workA*p = 0
// Returns p and the Lagrange multipliers.
func qpSolveEqSubproblem(Q, g []float64, n int, workA [][]float64, nWork int) ([]float64, []float64, error) {
	dim := n + nWork

	kkt := make([][]float64, dim)
	for i := 0; i < dim; i++ {
		kkt[i] = make([]float64, dim)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			kkt[i][j] = Q[i*n+j]
		}
	}
	for k := 0; k < nWork; k++ {
		for j := 0; j < n; j++ {
			kkt[j][n+k] = workA[k][j]
			kkt[n+k][j] = workA[k][j]
		}
	}

	rhs := make([]float64, dim)
	for i := 0; i < n; i++ {
		rhs[i] = -g[i]
	}
	// RHS for constraint rows is 0 (we want workA*p = 0).

	sol, err := solveLinearSystem(kkt, rhs)
	if err != nil {
		return nil, nil, err
	}

	p := sol[:n]
	lambdas := sol[n:]
	return p, lambdas, nil
}

// qpFindFeasible attempts to find a feasible point given constraints.
func qpFindFeasible(x0 []float64, n int, A [][]float64, b []float64,
	Aeq [][]float64, beq []float64) []float64 {
	x := make([]float64, n)
	copy(x, x0)

	// First satisfy equality constraints if any.
	nEq := len(Aeq)
	if nEq > 0 {
		// Solve Aeq * x = beq via least-norm solution: x = Aeq^T * (Aeq * Aeq^T)^{-1} * beq
		aat := make([][]float64, nEq)
		for i := 0; i < nEq; i++ {
			aat[i] = make([]float64, nEq)
			for j := 0; j < nEq; j++ {
				for k := 0; k < n; k++ {
					aat[i][j] += Aeq[i][k] * Aeq[j][k]
				}
			}
		}
		sol, err := solveLinearSystem(aat, beq)
		if err == nil {
			for i := 0; i < n; i++ {
				x[i] = 0
				for k := 0; k < nEq; k++ {
					x[i] += Aeq[k][i] * sol[k]
				}
			}
		}
	}

	// Project to satisfy inequality constraints iteratively.
	for iter := 0; iter < 100; iter++ {
		feasible := true
		for i, row := range A {
			val := dotVec(row, x)
			if val > b[i]+1e-12 {
				feasible = false
				// Project: x = x - ((a'x - b) / (a'a)) * a
				normSq := dotVec(row, row)
				if normSq > 1e-14 {
					scale := (val - b[i]) / normSq
					for j := 0; j < n; j++ {
						x[j] -= scale * row[j]
					}
				}
			}
		}
		if feasible {
			break
		}
	}

	return x
}

func qpObjective(Q, c, x []float64, n int) float64 {
	val := 0.0
	for i := 0; i < n; i++ {
		val += c[i] * x[i]
		for j := 0; j < n; j++ {
			val += 0.5 * Q[i*n+j] * x[i] * x[j]
		}
	}
	return val
}

func dotVec(a, b []float64) float64 {
	s := 0.0
	for i := range a {
		s += a[i] * b[i]
	}
	return s
}

// solveLinearSystem solves Ax = b using LU decomposition.
func solveLinearSystem(A [][]float64, b []float64) ([]float64, error) {
	lu, piv, err := LUFactor(A)
	if err != nil {
		return nil, err
	}
	return LUSolve(lu, piv, b)
}

// itoa converts an int to a string without importing strconv.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	digits := make([]byte, 0, 10)
	for n > 0 {
		digits = append(digits, byte('0'+n%10))
		n /= 10
	}
	if neg {
		digits = append(digits, '-')
	}
	// Reverse.
	for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
		digits[i], digits[j] = digits[j], digits[i]
	}
	return string(digits)
}
