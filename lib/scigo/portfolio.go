package scigo

import (
	"errors"
	"math"
)

// PortfolioReturn computes the expected return of a portfolio: w' * returns.
func PortfolioReturn(weights, returns []float64) float64 {
	r := 0.0
	for i := range weights {
		r += weights[i] * returns[i]
	}
	return r
}

// PortfolioRisk computes the portfolio standard deviation: sqrt(w' * cov * w).
func PortfolioRisk(weights []float64, cov [][]float64) float64 {
	n := len(weights)
	variance := 0.0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			variance += weights[i] * cov[i][j] * weights[j]
		}
	}
	if variance < 0 {
		variance = 0
	}
	return math.Sqrt(variance)
}

// SharpeRatio computes the Sharpe ratio of a portfolio:
// (portfolioReturn - riskFreeRate) / portfolioRisk.
func SharpeRatio(weights, returns []float64, cov [][]float64, riskFreeRate float64) float64 {
	pr := PortfolioReturn(weights, returns)
	risk := PortfolioRisk(weights, cov)
	if risk < 1e-15 {
		return 0
	}
	return (pr - riskFreeRate) / risk
}

// ValueAtRisk computes the parametric (Gaussian) Value at Risk for a portfolio.
// confidence is typically 0.95 or 0.99.
// VaR = -(mu + sigma * Phi^{-1}(1 - confidence))
func ValueAtRisk(weights, returns []float64, cov [][]float64, confidence float64) float64 {
	mu := PortfolioReturn(weights, returns)
	sigma := PortfolioRisk(weights, cov)
	stdNorm := NewNormal(0, 1)
	zAlpha := stdNorm.PPF(1 - confidence)
	return -(mu + sigma*zAlpha)
}

// ConditionalVaR computes the Conditional Value at Risk (Expected Shortfall)
// for a portfolio under the Gaussian assumption.
// CVaR = -mu + sigma * phi(Phi^{-1}(1-confidence)) / (1-confidence)
func ConditionalVaR(weights, returns []float64, cov [][]float64, confidence float64) float64 {
	mu := PortfolioReturn(weights, returns)
	sigma := PortfolioRisk(weights, cov)
	stdNorm := NewNormal(0, 1)
	zAlpha := stdNorm.PPF(1 - confidence)
	phiZ := stdNorm.PDF(zAlpha)
	return -mu + sigma*phiZ/(1-confidence)
}

// MinVariancePortfolio computes the minimum-variance portfolio weights.
// Constraints: weights sum to 1, weights >= 0.
// Uses QP: min 0.5 * w' * cov * w s.t. 1'w = 1, w >= 0.
func MinVariancePortfolio(returns []float64, cov [][]float64) ([]float64, error) {
	n := len(returns)
	if n == 0 {
		return nil, errors.New("scigo.MinVariancePortfolio: empty returns")
	}
	if len(cov) != n {
		return nil, errors.New("scigo.MinVariancePortfolio: cov size mismatch")
	}

	// Flatten cov for QP.
	Q := flattenCov(cov, n)
	c := make([]float64, n) // zero linear term

	// Equality constraint: sum(w) = 1.
	aeqRow := make([]float64, n)
	for i := 0; i < n; i++ {
		aeqRow[i] = 1.0
	}
	Aeq := [][]float64{aeqRow}
	beq := []float64{1.0}

	// Bounds: 0 <= w <= 1.
	lb := make([]float64, n)
	ub := make([]float64, n)
	for i := 0; i < n; i++ {
		ub[i] = 1.0
	}

	result, err := QPSolve(Q, c, n, nil, nil, Aeq, beq, lb, ub)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New("scigo.MinVariancePortfolio: QP solver did not converge")
	}

	return normalizeWeights(result.X), nil
}

// MaxSharpePortfolio computes the portfolio that maximizes the Sharpe ratio.
// Uses the tangency portfolio approach: solve for w proportional to cov^{-1} * (returns - rf).
func MaxSharpePortfolio(returns []float64, cov [][]float64, riskFreeRate float64) ([]float64, error) {
	n := len(returns)
	if n == 0 {
		return nil, errors.New("scigo.MaxSharpePortfolio: empty returns")
	}
	if len(cov) != n {
		return nil, errors.New("scigo.MaxSharpePortfolio: cov size mismatch")
	}

	// Excess returns.
	excess := make([]float64, n)
	for i := 0; i < n; i++ {
		excess[i] = returns[i] - riskFreeRate
	}

	// Solve cov * z = excess.
	covCopy := make([][]float64, n)
	for i := 0; i < n; i++ {
		covCopy[i] = make([]float64, n)
		copy(covCopy[i], cov[i])
	}

	lu, piv, err := LUFactor(covCopy)
	if err != nil {
		return nil, errors.New("scigo.MaxSharpePortfolio: singular covariance matrix")
	}
	z, err := LUSolve(lu, piv, excess)
	if err != nil {
		return nil, err
	}

	// Normalize so sum = 1.
	sumZ := 0.0
	for _, v := range z {
		sumZ += v
	}
	if math.Abs(sumZ) < 1e-14 {
		return nil, errors.New("scigo.MaxSharpePortfolio: degenerate problem")
	}

	weights := make([]float64, n)
	for i := range z {
		weights[i] = z[i] / sumZ
	}

	// If any weight is negative, fall back to constrained optimization via QP.
	hasNeg := false
	for _, w := range weights {
		if w < -1e-10 {
			hasNeg = true
			break
		}
	}
	if hasNeg {
		return maxSharpeConstrained(returns, cov, riskFreeRate, n)
	}

	return normalizeWeights(weights), nil
}

// maxSharpeConstrained solves for max Sharpe with non-negativity via QP on the
// equivalent formulation.
func maxSharpeConstrained(returns []float64, cov [][]float64, riskFreeRate float64, n int) ([]float64, error) {
	// We solve: min 0.5 * y' * cov * y  s.t. (returns - rf)' * y = 1, y >= 0
	// then w = y / sum(y).
	Q := flattenCov(cov, n)
	c := make([]float64, n)

	excess := make([]float64, n)
	for i := 0; i < n; i++ {
		excess[i] = returns[i] - riskFreeRate
	}
	Aeq := [][]float64{excess}
	beq := []float64{1.0}

	lb := make([]float64, n)
	ub := make([]float64, n)
	for i := 0; i < n; i++ {
		ub[i] = math.Inf(1)
	}

	result, err := QPSolve(Q, c, n, nil, nil, Aeq, beq, lb, ub)
	if err != nil {
		return nil, err
	}

	sumW := 0.0
	for _, v := range result.X {
		sumW += v
	}
	if math.Abs(sumW) < 1e-14 {
		return nil, errors.New("scigo.MaxSharpePortfolio: degenerate constrained problem")
	}
	weights := make([]float64, n)
	for i := range result.X {
		weights[i] = result.X[i] / sumW
	}
	return normalizeWeights(weights), nil
}

// TargetReturnPortfolio finds the minimum-variance portfolio with a given target return.
// min 0.5*w'*cov*w s.t. w'*returns = targetReturn, sum(w) = 1, w >= 0.
func TargetReturnPortfolio(returns []float64, cov [][]float64, targetReturn float64) ([]float64, error) {
	n := len(returns)
	if n == 0 {
		return nil, errors.New("scigo.TargetReturnPortfolio: empty returns")
	}
	if len(cov) != n {
		return nil, errors.New("scigo.TargetReturnPortfolio: cov size mismatch")
	}

	Q := flattenCov(cov, n)
	c := make([]float64, n)

	// Equality: sum(w)=1, w'*returns=targetReturn.
	onesRow := make([]float64, n)
	for i := 0; i < n; i++ {
		onesRow[i] = 1.0
	}
	retRow := make([]float64, n)
	copy(retRow, returns)

	Aeq := [][]float64{onesRow, retRow}
	beq := []float64{1.0, targetReturn}

	lb := make([]float64, n)
	ub := make([]float64, n)
	for i := 0; i < n; i++ {
		ub[i] = 1.0
	}

	result, err := QPSolve(Q, c, n, nil, nil, Aeq, beq, lb, ub)
	if err != nil {
		return nil, err
	}
	if !result.Success {
		return nil, errors.New("scigo.TargetReturnPortfolio: QP solver did not converge")
	}

	return normalizeWeights(result.X), nil
}

// EfficientFrontier computes the efficient frontier by solving for minimum-variance
// portfolios at nPoints equally spaced target returns between the min-variance
// return and the max individual return.
// Returns risks, rets, and weights slices.
func EfficientFrontier(returns []float64, cov [][]float64, nPoints int) ([]float64, []float64, [][]float64, error) {
	n := len(returns)
	if n == 0 {
		return nil, nil, nil, errors.New("scigo.EfficientFrontier: empty returns")
	}
	if nPoints < 2 {
		return nil, nil, nil, errors.New("scigo.EfficientFrontier: nPoints must be >= 2")
	}

	// Find the min-variance portfolio for the lower bound of target return.
	mvWeights, err := MinVariancePortfolio(returns, cov)
	if err != nil {
		return nil, nil, nil, err
	}
	minRet := PortfolioReturn(mvWeights, returns)

	// Max return among assets.
	maxRet := returns[0]
	for _, r := range returns[1:] {
		if r > maxRet {
			maxRet = r
		}
	}

	risks := make([]float64, nPoints)
	rets := make([]float64, nPoints)
	weights := make([][]float64, nPoints)

	for i := 0; i < nPoints; i++ {
		t := float64(i) / float64(nPoints-1)
		target := minRet + t*(maxRet-minRet)

		w, err := TargetReturnPortfolio(returns, cov, target)
		if err != nil {
			// Fall back to min-variance.
			w = mvWeights
		}
		weights[i] = w
		rets[i] = PortfolioReturn(w, returns)
		risks[i] = PortfolioRisk(w, cov)
	}

	return risks, rets, weights, nil
}

// flattenCov converts a 2D covariance matrix to a flat array.
func flattenCov(cov [][]float64, n int) []float64 {
	Q := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			Q[i*n+j] = cov[i][j]
		}
	}
	return Q
}

// normalizeWeights ensures weights sum to 1 and clips near-zero negatives.
func normalizeWeights(w []float64) []float64 {
	n := len(w)
	out := make([]float64, n)
	sum := 0.0
	for i := 0; i < n; i++ {
		if w[i] < 0 {
			out[i] = 0
		} else {
			out[i] = w[i]
		}
		sum += out[i]
	}
	if sum > 1e-14 {
		for i := range out {
			out[i] /= sum
		}
	}
	return out
}
