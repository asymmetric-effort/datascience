package prediction

import (
	"fmt"
	"math"
	"strings"

	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
)

// NaiveIVRegressor implements instrumental variable regression using
// two-stage least squares (2SLS) for causal effect estimation.
type NaiveIVRegressor struct {
	treatment       string
	outcome         string
	instruments     []string
	ate             float64
	se              float64
	intercept       float64 // stage 2 intercept
	betaStage1      []float64
	stage1Residuals []float64 // for F-stat
	tVals           []float64 // original treatment values
	nObs            int
	fitted          bool
}

// NewNaiveIVRegressor creates a new NaiveIVRegressor.
func NewNaiveIVRegressor(treatment, outcome string, instruments []string) *NaiveIVRegressor {
	inst := make([]string, len(instruments))
	copy(inst, instruments)
	return &NaiveIVRegressor{
		treatment:   treatment,
		outcome:     outcome,
		instruments: inst,
	}
}

// Fit performs two-stage least squares:
//
//	Stage 1: treatment ~ intercept + instruments (OLS) -> predicted treatment
//	Stage 2: outcome ~ intercept + predicted_treatment (OLS) -> ATE
func (r *NaiveIVRegressor) Fit(data *tabgo.DataFrame) error {
	n := data.Len()
	if n == 0 {
		return fmt.Errorf("prediction: empty DataFrame")
	}

	// Stage 1: regress treatment on instruments.
	Xstage1, err := buildDesignMatrix(data, r.instruments)
	if err != nil {
		return fmt.Errorf("prediction: stage 1: %w", err)
	}

	tVals, err := extractColumnFloat64(data, r.treatment)
	if err != nil {
		return fmt.Errorf("prediction: stage 1: %w", err)
	}

	betaStage1 := olsFit(tVals, Xstage1)

	// Compute predicted treatment values and stage 1 residuals.
	// Rebuild design matrix since olsFit modifies the augmented matrix.
	Xstage1b, _ := buildDesignMatrix(data, r.instruments)
	tHat := make([]float64, n)
	stage1Resid := make([]float64, n)
	for i := 0; i < n; i++ {
		tHat[i] = dotProduct(betaStage1, Xstage1b[i])
		stage1Resid[i] = tVals[i] - tHat[i]
	}

	// Stage 2: regress outcome on predicted treatment.
	Xstage2 := make([][]float64, n)
	for i := 0; i < n; i++ {
		Xstage2[i] = []float64{1.0, tHat[i]}
	}

	yVals, err := extractColumnFloat64(data, r.outcome)
	if err != nil {
		return fmt.Errorf("prediction: stage 2: %w", err)
	}

	betaStage2 := olsFit(yVals, Xstage2)

	r.ate = betaStage2[1]
	r.intercept = betaStage2[0]
	r.betaStage1 = betaStage1
	r.stage1Residuals = stage1Resid
	r.tVals = tVals
	r.nObs = n

	// Compute SE for the 2SLS estimate.
	// Use original residuals (not from fitted values):
	// e_i = y_i - intercept - ate * tHat_i
	// Rebuild Xstage2 for SE computation.
	Xstage2b := make([][]float64, n)
	for i := 0; i < n; i++ {
		Xstage2b[i] = []float64{1.0, tHat[i]}
	}
	sse := 0.0
	for i := 0; i < n; i++ {
		ei := yVals[i] - r.intercept - r.ate*tHat[i]
		sse += ei * ei
	}
	sigma2 := sse / float64(n-2)
	r.se = computeCoefficientSE(Xstage2b, sigma2, 1) // index 1 = treatment

	r.fitted = true
	return nil
}

// ATE returns the estimated Average Treatment Effect from 2SLS.
func (r *NaiveIVRegressor) ATE() float64 {
	return r.ate
}

// Predict returns predicted outcome values: intercept + ATE * treatment.
func (r *NaiveIVRegressor) Predict(data *tabgo.DataFrame) ([]float64, error) {
	if !r.fitted {
		return nil, fmt.Errorf("prediction: model not fitted")
	}
	// First stage: predict treatment from instruments.
	Xstage1, err := buildDesignMatrix(data, r.instruments)
	if err != nil {
		return nil, err
	}
	n := data.Len()
	preds := make([]float64, n)
	for i := 0; i < n; i++ {
		tHat := dotProduct(r.betaStage1, Xstage1[i])
		preds[i] = r.intercept + r.ate*tHat
	}
	return preds, nil
}

// SE returns the standard error of the ATE estimate.
func (r *NaiveIVRegressor) SE() float64 {
	return r.se
}

// ConfidenceInterval returns (lower, upper) bounds for the ATE at significance level alpha.
func (r *NaiveIVRegressor) ConfidenceInterval(alpha float64) (float64, float64) {
	z := normalQuantile(1 - alpha/2)
	return r.ate - z*r.se, r.ate + z*r.se
}

// PValue returns the two-sided p-value for testing H0: ATE = 0.
func (r *NaiveIVRegressor) PValue() float64 {
	if r.se == 0 {
		return 0
	}
	z := math.Abs(r.ate / r.se)
	return 2 * (1 - normalCDF(z))
}

// FirstStageFStat returns the F-statistic from the first stage regression,
// testing whether the instruments are jointly significant predictors of treatment.
// A value > 10 is typically considered evidence of strong instruments.
func (r *NaiveIVRegressor) FirstStageFStat() float64 {
	if !r.fitted {
		return 0
	}
	n := r.nObs
	k := len(r.instruments) // number of instruments (excluded regressors)
	if k == 0 || n <= k+1 {
		return 0
	}

	// Compute TSS (total sum of squares of treatment).
	meanT := 0.0
	for _, v := range r.tVals {
		meanT += v
	}
	meanT /= float64(n)

	tss := 0.0
	for _, v := range r.tVals {
		d := v - meanT
		tss += d * d
	}

	// RSS (residual sum of squares from stage 1).
	rss := 0.0
	for _, e := range r.stage1Residuals {
		rss += e * e
	}

	// F = ((TSS - RSS) / k) / (RSS / (n - k - 1))
	p := k + 1 // total params in stage 1 (instruments + intercept)
	if n <= p {
		return 0
	}
	fStat := ((tss - rss) / float64(k)) / (rss / float64(n-p))
	return fStat
}

// Summary returns a formatted summary string.
func (r *NaiveIVRegressor) Summary() string {
	if !r.fitted {
		return "NaiveIVRegressor: not fitted"
	}
	lo, hi := r.ConfidenceInterval(0.05)
	var sb strings.Builder
	sb.WriteString("NaiveIVRegressor Summary\n")
	sb.WriteString("========================\n")
	fmt.Fprintf(&sb, "Treatment:       %s\n", r.treatment)
	fmt.Fprintf(&sb, "Outcome:         %s\n", r.outcome)
	fmt.Fprintf(&sb, "Instruments:     %v\n", r.instruments)
	fmt.Fprintf(&sb, "N Obs:           %d\n", r.nObs)
	sb.WriteString("------------------------\n")
	fmt.Fprintf(&sb, "ATE:             %.6f\n", r.ate)
	fmt.Fprintf(&sb, "Std. Error:      %.6f\n", r.se)
	fmt.Fprintf(&sb, "95%% CI:          [%.6f, %.6f]\n", lo, hi)
	fmt.Fprintf(&sb, "P-value:         %.6f\n", r.PValue())
	fmt.Fprintf(&sb, "1st Stage F:     %.4f\n", r.FirstStageFStat())
	return sb.String()
}
