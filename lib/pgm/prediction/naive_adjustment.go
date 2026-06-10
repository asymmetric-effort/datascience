package prediction

import (
	"fmt"
	"math"
	"strings"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// NaiveAdjustmentRegressor implements naive back-door adjustment for causal
// effect estimation via OLS regression: outcome ~ treatment + adjustmentSet.
type NaiveAdjustmentRegressor struct {
	treatment     string
	outcome       string
	adjustmentSet []string
	coefficients  []float64 // [intercept, treatment, adj1, adj2, ...]
	se            float64   // standard error of the treatment coefficient
	nObs          int       // number of observations used in fit
	residuals     []float64 // OLS residuals
	designMatrix  [][]float64
	fitted        bool
}

// NewNaiveAdjustmentRegressor creates a new NaiveAdjustmentRegressor.
func NewNaiveAdjustmentRegressor(treatment, outcome string, adjustmentSet []string) *NaiveAdjustmentRegressor {
	a := make([]string, len(adjustmentSet))
	copy(a, adjustmentSet)
	return &NaiveAdjustmentRegressor{
		treatment:     treatment,
		outcome:       outcome,
		adjustmentSet: a,
	}
}

// Fit performs OLS regression: outcome ~ intercept + treatment + adjustmentSet.
func (r *NaiveAdjustmentRegressor) Fit(data *tabgo.DataFrame) error {
	n := data.Len()
	if n == 0 {
		return fmt.Errorf("prediction: empty DataFrame")
	}

	// Build column list: treatment first, then adjustment set.
	regressors := make([]string, 0, 1+len(r.adjustmentSet))
	regressors = append(regressors, r.treatment)
	regressors = append(regressors, r.adjustmentSet...)

	// Build design matrix with intercept.
	X, err := buildDesignMatrix(data, regressors)
	if err != nil {
		return err
	}

	y, err := extractColumnFloat64(data, r.outcome)
	if err != nil {
		return err
	}

	r.coefficients = olsFit(y, X)
	r.nObs = n

	// Compute residuals and standard error of the treatment coefficient.
	// Rebuild design matrix (olsFit consumed it via solveLinearSystem).
	regressors2 := make([]string, 0, 1+len(r.adjustmentSet))
	regressors2 = append(regressors2, r.treatment)
	regressors2 = append(regressors2, r.adjustmentSet...)
	X2, _ := buildDesignMatrix(data, regressors2)
	r.designMatrix = X2

	y2, _ := extractColumnFloat64(data, r.outcome)
	p := len(r.coefficients)
	r.residuals = make([]float64, n)
	sse := 0.0
	for i := 0; i < n; i++ {
		pred := dotProduct(r.coefficients, X2[i])
		r.residuals[i] = y2[i] - pred
		sse += r.residuals[i] * r.residuals[i]
	}
	sigma2 := sse / float64(n-p)

	// SE of treatment coefficient = sqrt(sigma2 * (X'X)^{-1}[1,1])
	// Recompute X'X and invert to get diagonal element.
	r.se = computeCoefficientSE(X2, sigma2, 1) // index 1 = treatment

	r.fitted = true
	return nil
}

// ATE returns the estimated Average Treatment Effect, which is the OLS
// coefficient on the treatment variable.
func (r *NaiveAdjustmentRegressor) ATE() float64 {
	if !r.fitted {
		return 0
	}
	// coefficients[0] = intercept, coefficients[1] = treatment
	return r.coefficients[1]
}

// Predict returns predicted outcome values for each row.
func (r *NaiveAdjustmentRegressor) Predict(data *tabgo.DataFrame) ([]float64, error) {
	if !r.fitted {
		return nil, fmt.Errorf("prediction: model not fitted")
	}
	regressors := make([]string, 0, 1+len(r.adjustmentSet))
	regressors = append(regressors, r.treatment)
	regressors = append(regressors, r.adjustmentSet...)
	X, err := buildDesignMatrix(data, regressors)
	if err != nil {
		return nil, err
	}
	preds := make([]float64, len(X))
	for i, row := range X {
		preds[i] = dotProduct(r.coefficients, row)
	}
	return preds, nil
}

// SE returns the standard error of the ATE estimate.
func (r *NaiveAdjustmentRegressor) SE() float64 {
	return r.se
}

// ConfidenceInterval returns (lower, upper) bounds for the ATE at significance level alpha.
func (r *NaiveAdjustmentRegressor) ConfidenceInterval(alpha float64) (float64, float64) {
	z := normalQuantile(1 - alpha/2)
	ate := r.ATE()
	return ate - z*r.se, ate + z*r.se
}

// PValue returns the two-sided p-value for testing H0: ATE = 0.
func (r *NaiveAdjustmentRegressor) PValue() float64 {
	if r.se == 0 {
		return 0
	}
	z := math.Abs(r.ATE() / r.se)
	return 2 * (1 - normalCDF(z))
}

// Summary returns a formatted summary string.
func (r *NaiveAdjustmentRegressor) Summary() string {
	if !r.fitted {
		return "NaiveAdjustmentRegressor: not fitted"
	}
	lo, hi := r.ConfidenceInterval(0.05)
	var sb strings.Builder
	sb.WriteString("NaiveAdjustmentRegressor Summary\n")
	sb.WriteString("================================\n")
	fmt.Fprintf(&sb, "Treatment:      %s\n", r.treatment)
	fmt.Fprintf(&sb, "Outcome:        %s\n", r.outcome)
	fmt.Fprintf(&sb, "Adjustment Set: %v\n", r.adjustmentSet)
	fmt.Fprintf(&sb, "N Obs:          %d\n", r.nObs)
	sb.WriteString("--------------------------------\n")
	fmt.Fprintf(&sb, "ATE:            %.6f\n", r.ATE())
	fmt.Fprintf(&sb, "Std. Error:     %.6f\n", r.se)
	fmt.Fprintf(&sb, "95%% CI:         [%.6f, %.6f]\n", lo, hi)
	fmt.Fprintf(&sb, "P-value:        %.6f\n", r.PValue())
	return sb.String()
}
