package prediction

import (
	"fmt"
	"math"
	"strings"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// DoubleMLRegressor implements Double Machine Learning for causal effect estimation.
// It uses cross-fitting with configurable number of folds to estimate the Average Treatment Effect (ATE).
type DoubleMLRegressor struct {
	treatment   string
	outcome     string
	confounders []string
	nSplits     int
	ate         float64
	se          float64
	yResiduals  []float64 // pooled outcome residuals from cross-fitting
	tResiduals  []float64 // pooled treatment residuals from cross-fitting
	fitted      bool
}

// NewDoubleMLRegressor creates a new DoubleMLRegressor.
// It defaults to 2-fold cross-fitting. Use SetNSplits to change.
func NewDoubleMLRegressor(treatment, outcome string, confounders []string) *DoubleMLRegressor {
	c := make([]string, len(confounders))
	copy(c, confounders)
	return &DoubleMLRegressor{
		treatment:   treatment,
		outcome:     outcome,
		confounders: c,
		nSplits:     2,
	}
}

// SetNSplits sets the number of cross-fitting folds. Must be >= 2.
func (d *DoubleMLRegressor) SetNSplits(n int) {
	if n < 2 {
		n = 2
	}
	d.nSplits = n
}

// Fit performs DML estimation with cross-fitting on the given data.
//
// Algorithm (K-fold cross-fitting):
//  1. Split data into K folds.
//  2. For each fold k: train outcome and treatment models on all other folds,
//     compute residuals on fold k.
//  3. Pool all residuals and estimate ATE as the OLS coefficient of treatment
//     residuals on outcome residuals (no intercept).
//  4. Compute standard error from residuals.
func (d *DoubleMLRegressor) Fit(data *tabgo.DataFrame) error {
	n := data.Len()
	if n < 2*d.nSplits {
		return fmt.Errorf("prediction: need at least %d observations for %d folds, got %d", 2*d.nSplits, d.nSplits, n)
	}

	K := d.nSplits
	// Compute fold boundaries.
	foldSize := n / K
	foldStarts := make([]int, K+1)
	for i := 0; i < K; i++ {
		foldStarts[i] = i * foldSize
	}
	foldStarts[K] = n

	var allYResid, allTResid []float64

	for k := 0; k < K; k++ {
		testStart := foldStarts[k]
		testEnd := foldStarts[k+1]
		testData := sliceDataFrame(data, testStart, testEnd)

		// Build training data: all rows except fold k.
		trainData := excludeSlice(data, testStart, testEnd)

		yR, tR, err := crossFitResiduals(trainData, testData, d.outcome, d.treatment, d.confounders)
		if err != nil {
			return fmt.Errorf("prediction: cross-fit fold %d: %w", k, err)
		}
		allYResid = append(allYResid, yR...)
		allTResid = append(allTResid, tR...)
	}

	// Estimate ATE: regress outcome residuals on treatment residuals (no intercept).
	// beta = sum(tResid * yResid) / sum(tResid * tResid)
	num := 0.0
	den := 0.0
	for i := range allYResid {
		num += allTResid[i] * allYResid[i]
		den += allTResid[i] * allTResid[i]
	}
	if den == 0 {
		return fmt.Errorf("prediction: treatment residuals are all zero; cannot estimate ATE")
	}
	d.ate = num / den

	// Store residuals for SE computation.
	d.yResiduals = allYResid
	d.tResiduals = allTResid

	// Compute standard error.
	nResid := len(allYResid)
	// Residuals from the final regression: e_i = yResid_i - ate * tResid_i
	sse := 0.0
	for i := 0; i < nResid; i++ {
		e := allYResid[i] - d.ate*allTResid[i]
		sse += e * e
	}
	sigma2 := sse / float64(nResid-1)
	d.se = math.Sqrt(sigma2 / den)

	d.fitted = true
	return nil
}

// excludeSlice builds a DataFrame from all rows of df except [start, end).
func excludeSlice(df *tabgo.DataFrame, start, end int) *tabgo.DataFrame {
	n := df.Len()
	names := df.Columns()
	colMap := make(map[string]*tabgo.Series, len(names))
	for _, name := range names {
		vals := df.Column(name).Values()
		excluded := make([]any, 0, n-(end-start))
		excluded = append(excluded, vals[:start]...)
		excluded = append(excluded, vals[end:]...)
		colMap[name] = tabgo.NewSeries(name, excluded)
	}
	return tabgo.NewDataFrame(colMap)
}

// ATE returns the estimated Average Treatment Effect.
func (d *DoubleMLRegressor) ATE() float64 {
	return d.ate
}

// Predict returns counterfactual outcome predictions for each row:
// predicted_outcome = ATE * treatment_value.
// This is a simplified prediction that applies the estimated treatment effect.
func (d *DoubleMLRegressor) Predict(data *tabgo.DataFrame) ([]float64, error) {
	if !d.fitted {
		return nil, fmt.Errorf("prediction: model not fitted")
	}
	t, err := extractColumnFloat64(data, d.treatment)
	if err != nil {
		return nil, err
	}
	preds := make([]float64, len(t))
	for i, tv := range t {
		preds[i] = d.ate * tv
	}
	return preds, nil
}

// SE returns the standard error of the ATE estimate.
func (d *DoubleMLRegressor) SE() float64 {
	return d.se
}

// ConfidenceInterval returns the (lower, upper) bounds of a confidence interval
// for the ATE at the given significance level alpha (e.g., 0.05 for 95% CI).
// Uses a normal approximation.
func (d *DoubleMLRegressor) ConfidenceInterval(alpha float64) (float64, float64) {
	z := normalQuantile(1 - alpha/2)
	return d.ate - z*d.se, d.ate + z*d.se
}

// PValue returns the two-sided p-value for testing H0: ATE = 0.
func (d *DoubleMLRegressor) PValue() float64 {
	if d.se == 0 {
		return 0
	}
	z := math.Abs(d.ate / d.se)
	return 2 * (1 - normalCDF(z))
}

// EstimateCate estimates conditional average treatment effects (CATE) for each
// observation. It computes observation-level treatment effects as:
// cate_i = yResid_i / tResid_i, smoothed by local weighting.
// For a linear DML model, the CATE is constant and equals the ATE;
// this method returns per-observation influence-weighted effects.
func (d *DoubleMLRegressor) EstimateCate() ([]float64, error) {
	if !d.fitted {
		return nil, fmt.Errorf("prediction: model not fitted")
	}
	n := len(d.tResiduals)
	cate := make([]float64, n)

	// Compute observation-level scores: psi_i = tResid_i * (yResid_i - ate * tResid_i)
	// CATE_i = ATE + influence_i where influence is the IF contribution.
	// For linear DML: theta_i = ATE + psi_i / E[tResid^2]
	den := 0.0
	for i := 0; i < n; i++ {
		den += d.tResiduals[i] * d.tResiduals[i]
	}
	den /= float64(n)

	for i := 0; i < n; i++ {
		psi := d.tResiduals[i] * (d.yResiduals[i] - d.ate*d.tResiduals[i])
		cate[i] = d.ate + psi/den
	}
	return cate, nil
}

// Summary returns a formatted summary string with ATE, SE, CI, and p-value.
func (d *DoubleMLRegressor) Summary() string {
	if !d.fitted {
		return "DoubleMLRegressor: not fitted"
	}
	lo, hi := d.ConfidenceInterval(0.05)
	var sb strings.Builder
	sb.WriteString("DoubleMLRegressor Summary\n")
	sb.WriteString("========================\n")
	fmt.Fprintf(&sb, "Treatment:   %s\n", d.treatment)
	fmt.Fprintf(&sb, "Outcome:     %s\n", d.outcome)
	fmt.Fprintf(&sb, "Confounders: %v\n", d.confounders)
	fmt.Fprintf(&sb, "N Splits:    %d\n", d.nSplits)
	fmt.Fprintf(&sb, "N Obs:       %d\n", len(d.yResiduals))
	sb.WriteString("------------------------\n")
	fmt.Fprintf(&sb, "ATE:         %.6f\n", d.ate)
	fmt.Fprintf(&sb, "Std. Error:  %.6f\n", d.se)
	fmt.Fprintf(&sb, "95%% CI:      [%.6f, %.6f]\n", lo, hi)
	fmt.Fprintf(&sb, "P-value:     %.6f\n", d.PValue())
	return sb.String()
}

// crossFitResiduals trains OLS models on trainData and computes residuals on testData.
func crossFitResiduals(
	trainData, testData *tabgo.DataFrame,
	outcome, treatment string,
	confounders []string,
) (yResid, tResid []float64, err error) {
	// Build training design matrix (intercept + confounders).
	Xtrain, err := buildDesignMatrix(trainData, confounders)
	if err != nil {
		return nil, nil, err
	}

	// Fit outcome ~ confounders on training data.
	yTrain, err := extractColumnFloat64(trainData, outcome)
	if err != nil {
		return nil, nil, err
	}
	betaY := olsFit(yTrain, Xtrain)

	// Fit treatment ~ confounders on training data.
	tTrain, err := extractColumnFloat64(trainData, treatment)
	if err != nil {
		return nil, nil, err
	}
	betaT := olsFit(tTrain, Xtrain)

	// Build test design matrix.
	Xtest, err := buildDesignMatrix(testData, confounders)
	if err != nil {
		return nil, nil, err
	}

	// Compute residuals on test data.
	yTest, err := extractColumnFloat64(testData, outcome)
	if err != nil {
		return nil, nil, err
	}
	tTest, err := extractColumnFloat64(testData, treatment)
	if err != nil {
		return nil, nil, err
	}

	nTest := len(yTest)
	yResid = make([]float64, nTest)
	tResid = make([]float64, nTest)
	for i := 0; i < nTest; i++ {
		yPred := dotProduct(betaY, Xtest[i])
		tPred := dotProduct(betaT, Xtest[i])
		yResid[i] = yTest[i] - yPred
		tResid[i] = tTest[i] - tPred
	}
	return yResid, tResid, nil
}

// dotProduct computes the dot product of two slices of equal length.
func dotProduct(a, b []float64) float64 {
	sum := 0.0
	for i := range a {
		sum += a[i] * b[i]
	}
	return sum
}

// sliceDataFrame returns rows [start, end) of a DataFrame.
func sliceDataFrame(df *tabgo.DataFrame, start, end int) *tabgo.DataFrame {
	names := df.Columns()
	colMap := make(map[string]*tabgo.Series, len(names))
	for _, name := range names {
		vals := df.Column(name).Values()
		sliced := vals[start:end]
		colMap[name] = tabgo.NewSeries(name, sliced)
	}
	return tabgo.NewDataFrame(colMap)
}
