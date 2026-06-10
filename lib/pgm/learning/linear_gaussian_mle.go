package learning

import (
	"fmt"
	"math"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// LinearGaussianMLE estimates LinearGaussianCPD parameters for a
// LinearGaussianBayesianNetwork from observed continuous data using
// ordinary least squares (OLS) regression.
type LinearGaussianMLE struct {
	bn   *models.LinearGaussianBayesianNetwork
	data *tabgo.DataFrame
}

// NewLinearGaussianMLE creates a new LinearGaussianMLE estimator.
// The network must already have nodes and edges defined. The DataFrame
// should have float64-convertible columns matching the node names.
func NewLinearGaussianMLE(bn *models.LinearGaussianBayesianNetwork, data *tabgo.DataFrame) *LinearGaussianMLE {
	return &LinearGaussianMLE{
		bn:   bn,
		data: data,
	}
}

// Estimate fits LinearGaussianCPDs for all nodes in the network from the data
// using OLS regression. For each node it:
//  1. Gets the node's parents.
//  2. Extracts the node column as Y and parent columns as X (with intercept).
//  3. Computes beta = (X'X)^{-1} X'Y via Gaussian elimination.
//  4. Computes residual variance = sum((Y - X*beta)^2) / (n - k).
//  5. Creates and stores the resulting LinearGaussianCPD.
func (mle *LinearGaussianMLE) Estimate() error {
	if mle.bn == nil {
		return fmt.Errorf("learning: LinearGaussianBayesianNetwork is nil")
	}
	if mle.data == nil {
		return fmt.Errorf("learning: data is nil")
	}

	nodes := mle.bn.Nodes()
	if len(nodes) == 0 {
		return fmt.Errorf("learning: LinearGaussianBayesianNetwork has no nodes")
	}

	// Validate that all required columns exist.
	dataColumns := make(map[string]bool)
	for _, c := range mle.data.Columns() {
		dataColumns[c] = true
	}
	for _, node := range nodes {
		if !dataColumns[node] {
			return fmt.Errorf("learning: data is missing column for node %q", node)
		}
	}

	for _, node := range nodes {
		cpd, err := mle.estimateNode(node)
		if err != nil {
			return fmt.Errorf("learning: failed to estimate LG CPD for %q: %w", node, err)
		}
		if err := mle.bn.AddLinearGaussianCPD(cpd); err != nil {
			return fmt.Errorf("learning: failed to add LG CPD for %q: %w", node, err)
		}
	}
	return nil
}

// GetParameters estimates and returns the LinearGaussianCPD for a single node.
func (mle *LinearGaussianMLE) GetParameters(node string) (*factors.LinearGaussianCPD, error) {
	if mle.bn == nil {
		return nil, fmt.Errorf("learning: LinearGaussianBayesianNetwork is nil")
	}
	if mle.data == nil {
		return nil, fmt.Errorf("learning: data is nil")
	}

	found := false
	for _, n := range mle.bn.Nodes() {
		if n == node {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("learning: node %q not found in LinearGaussianBayesianNetwork", node)
	}

	dataColumns := make(map[string]bool)
	for _, c := range mle.data.Columns() {
		dataColumns[c] = true
	}
	if !dataColumns[node] {
		return nil, fmt.Errorf("learning: data is missing column for node %q", node)
	}
	parents := mle.bn.Parents(node)
	for _, p := range parents {
		if !dataColumns[p] {
			return nil, fmt.Errorf("learning: data is missing column for parent %q", p)
		}
	}

	return mle.estimateNode(node)
}

// estimateNode computes the OLS-based LinearGaussianCPD for a single node.
func (mle *LinearGaussianMLE) estimateNode(node string) (*factors.LinearGaussianCPD, error) {
	parents := mle.bn.Parents(node) // sorted
	n := mle.data.Len()
	k := len(parents) + 1 // number of parameters (intercept + betas)

	if n < k {
		return nil, fmt.Errorf("learning: insufficient data rows (%d) for %d parameters", n, k)
	}

	// Extract Y vector (node column).
	y := mle.data.Column(node).Float64()

	// Build X matrix: n x k, first column is all 1s (intercept).
	// X is stored row-major: X[i*k + j]
	X := make([]float64, n*k)
	for i := 0; i < n; i++ {
		X[i*k] = 1.0 // intercept column
	}
	for j, p := range parents {
		col := mle.data.Column(p).Float64()
		for i := 0; i < n; i++ {
			X[i*k+(j+1)] = col[i]
		}
	}

	// Compute X'X (k x k) and X'Y (k x 1).
	xtx := make([]float64, k*k)
	xty := make([]float64, k)
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			sum := 0.0
			for r := 0; r < n; r++ {
				sum += X[r*k+i] * X[r*k+j]
			}
			xtx[i*k+j] = sum
		}
		sum := 0.0
		for r := 0; r < n; r++ {
			sum += X[r*k+i] * y[r]
		}
		xty[i] = sum
	}

	// Solve (X'X) beta = X'Y using Gaussian elimination with partial pivoting.
	beta, err := solveLinearSystem(xtx, xty, k)
	if err != nil {
		return nil, fmt.Errorf("learning: OLS solve failed for %q: %w", node, err)
	}

	// Compute residual variance: var = sum((Y - X*beta)^2) / (n - k).
	var residualSS float64
	for i := 0; i < n; i++ {
		predicted := 0.0
		for j := 0; j < k; j++ {
			predicted += X[i*k+j] * beta[j]
		}
		residual := y[i] - predicted
		residualSS += residual * residual
	}

	denominator := n - k
	if denominator <= 0 {
		denominator = 1 // avoid division by zero; biased estimate when n == k
	}
	variance := residualSS / float64(denominator)
	if variance <= 0 {
		// Guard against perfect fit giving zero variance.
		variance = 1e-10
	}

	// beta[0] is the intercept (mean), beta[1:] are the parent betas.
	mean := beta[0]
	betas := make([]float64, len(parents))
	copy(betas, beta[1:])

	return factors.NewLinearGaussianCPD(node, mean, betas, variance, parents)
}

// solveLinearSystem solves A*x = b where A is k x k (row-major) using
// Gaussian elimination with partial pivoting. Returns the solution vector x.
func solveLinearSystem(A, b []float64, k int) ([]float64, error) {
	// Build augmented matrix [A | b], k x (k+1).
	aug := make([]float64, k*(k+1))
	for i := 0; i < k; i++ {
		for j := 0; j < k; j++ {
			aug[i*(k+1)+j] = A[i*k+j]
		}
		aug[i*(k+1)+k] = b[i]
	}

	// Forward elimination with partial pivoting.
	for col := 0; col < k; col++ {
		// Find pivot row.
		maxVal := math.Abs(aug[col*(k+1)+col])
		maxRow := col
		for row := col + 1; row < k; row++ {
			v := math.Abs(aug[row*(k+1)+col])
			if v > maxVal {
				maxVal = v
				maxRow = row
			}
		}
		if maxVal < 1e-15 {
			return nil, fmt.Errorf("singular or near-singular matrix")
		}

		// Swap rows.
		if maxRow != col {
			for j := 0; j <= k; j++ {
				aug[col*(k+1)+j], aug[maxRow*(k+1)+j] = aug[maxRow*(k+1)+j], aug[col*(k+1)+j]
			}
		}

		// Eliminate below.
		pivot := aug[col*(k+1)+col]
		for row := col + 1; row < k; row++ {
			factor := aug[row*(k+1)+col] / pivot
			for j := col; j <= k; j++ {
				aug[row*(k+1)+j] -= factor * aug[col*(k+1)+j]
			}
		}
	}

	// Back substitution.
	x := make([]float64, k)
	for i := k - 1; i >= 0; i-- {
		x[i] = aug[i*(k+1)+k]
		for j := i + 1; j < k; j++ {
			x[i] -= aug[i*(k+1)+j] * x[j]
		}
		x[i] /= aug[i*(k+1)+i]
	}

	return x, nil
}
