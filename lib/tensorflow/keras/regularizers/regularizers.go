// Package regularizers provides weight regularization functions,
// analogous to tf.keras.regularizers.
package regularizers

import "math"

// Regularizer computes a penalty term from weight values.
type Regularizer interface {
	// Penalty computes the regularization penalty for the given weights.
	Penalty(weights []float64) float64
	// Name returns the regularizer name.
	Name() string
}

// L1 computes L1 regularization: l1 * sum(|w|).
type L1 struct {
	L1Factor float64
}

// NewL1 creates a new L1 regularizer.
func NewL1(l1 float64) *L1 {
	return &L1{L1Factor: l1}
}

// Penalty returns the L1 penalty.
func (r *L1) Penalty(weights []float64) float64 {
	sum := 0.0
	for _, w := range weights {
		sum += math.Abs(w)
	}
	return r.L1Factor * sum
}

// Name returns "l1".
func (r *L1) Name() string {
	return "l1"
}

// L2 computes L2 regularization: l2 * sum(w^2).
type L2 struct {
	L2Factor float64
}

// NewL2 creates a new L2 regularizer.
func NewL2(l2 float64) *L2 {
	return &L2{L2Factor: l2}
}

// Penalty returns the L2 penalty.
func (r *L2) Penalty(weights []float64) float64 {
	sum := 0.0
	for _, w := range weights {
		sum += w * w
	}
	return r.L2Factor * sum
}

// Name returns "l2".
func (r *L2) Name() string {
	return "l2"
}

// L1L2 computes combined L1 and L2 regularization.
type L1L2 struct {
	L1Factor float64
	L2Factor float64
}

// NewL1L2 creates a new L1L2 regularizer.
func NewL1L2(l1, l2 float64) *L1L2 {
	return &L1L2{L1Factor: l1, L2Factor: l2}
}

// Penalty returns the combined L1+L2 penalty.
func (r *L1L2) Penalty(weights []float64) float64 {
	l1Sum := 0.0
	l2Sum := 0.0
	for _, w := range weights {
		l1Sum += math.Abs(w)
		l2Sum += w * w
	}
	return r.L1Factor*l1Sum + r.L2Factor*l2Sum
}

// Name returns "l1_l2".
func (r *L1L2) Name() string {
	return "l1_l2"
}

// OrthogonalRegularizer encourages weight matrices to be orthogonal.
// Penalty = factor * ||W^T W - I||^2_F.
type OrthogonalRegularizer struct {
	Factor float64
	N      int // matrix is NxN
}

// NewOrthogonalRegularizer creates a new orthogonal regularizer for NxN weight matrices.
func NewOrthogonalRegularizer(factor float64, n int) *OrthogonalRegularizer {
	return &OrthogonalRegularizer{Factor: factor, N: n}
}

// Penalty computes ||W^T W - I||^2_F.
func (r *OrthogonalRegularizer) Penalty(weights []float64) float64 {
	n := r.N
	if len(weights) != n*n {
		return 0
	}
	// Compute W^T W.
	wtw := make([]float64, n*n)
	for i := range n {
		for j := range n {
			sum := 0.0
			for k := range n {
				sum += weights[k*n+i] * weights[k*n+j]
			}
			wtw[i*n+j] = sum
		}
	}
	// Compute ||W^T W - I||^2_F.
	penalty := 0.0
	for i := range n {
		for j := range n {
			target := 0.0
			if i == j {
				target = 1.0
			}
			diff := wtw[i*n+j] - target
			penalty += diff * diff
		}
	}
	return r.Factor * penalty
}

// Name returns "orthogonal".
func (r *OrthogonalRegularizer) Name() string {
	return "orthogonal"
}
