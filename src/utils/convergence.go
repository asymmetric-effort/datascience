package utils

import "math"

// ConvergenceChecker tracks a scalar value across iterations and determines
// whether it has converged (i.e., the absolute change falls below a tolerance).
type ConvergenceChecker struct {
	tolerance float64
	maxIter   int
	iteration int
	previous  float64
	converged bool
	started   bool
}

// NewConvergenceChecker creates a ConvergenceChecker with the given tolerance
// and maximum iteration count.
func NewConvergenceChecker(tolerance float64, maxIter int) *ConvergenceChecker {
	return &ConvergenceChecker{
		tolerance: tolerance,
		maxIter:   maxIter,
	}
}

// Update records a new value and returns true if the checker has converged
// (the absolute change from the previous value is less than tolerance) or
// the maximum number of iterations has been reached.
func (c *ConvergenceChecker) Update(value float64) bool {
	c.iteration++
	if c.started {
		if math.Abs(value-c.previous) < c.tolerance {
			c.converged = true
		}
	}
	c.started = true
	c.previous = value
	if c.iteration >= c.maxIter {
		c.converged = true
	}
	return c.converged
}

// Iterations returns the number of Update calls so far.
func (c *ConvergenceChecker) Iterations() int {
	return c.iteration
}

// Converged returns whether convergence has been detected.
func (c *ConvergenceChecker) Converged() bool {
	return c.converged
}

// Reset clears the checker state so it can be reused.
func (c *ConvergenceChecker) Reset() {
	c.iteration = 0
	c.previous = 0
	c.converged = false
	c.started = false
}
