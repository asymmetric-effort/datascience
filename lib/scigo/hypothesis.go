package scigo

import "math"

// ChiSquareTest performs a Pearson chi-squared goodness-of-fit test.
// It compares observed frequencies against expected frequencies and returns
// the test statistic and the p-value.
// The degrees of freedom is len(observed) - 1.
// Panics if the slices have different lengths or fewer than 2 elements.
func ChiSquareTest(observed, expected []float64) (statistic, pvalue float64) {
	if len(observed) != len(expected) {
		panic("scigo: ChiSquareTest: observed and expected must have the same length")
	}
	if len(observed) < 2 {
		panic("scigo: ChiSquareTest: need at least 2 categories")
	}

	statistic = 0
	for i := range observed {
		diff := observed[i] - expected[i]
		statistic += diff * diff / expected[i]
	}

	df := float64(len(observed) - 1)
	chi2 := NewChiSquared(df)
	pvalue = chi2.SurvivalFunction(statistic)
	return statistic, pvalue
}

// GTest performs a G-test (log-likelihood ratio test) for goodness of fit.
// It compares observed frequencies against expected frequencies and returns
// the test statistic G = 2 * sum(observed * ln(observed/expected)) and the p-value.
// The degrees of freedom is len(observed) - 1.
// Panics if the slices have different lengths or fewer than 2 elements.
func GTest(observed, expected []float64) (statistic, pvalue float64) {
	if len(observed) != len(expected) {
		panic("scigo: GTest: observed and expected must have the same length")
	}
	if len(observed) < 2 {
		panic("scigo: GTest: need at least 2 categories")
	}

	statistic = 0
	for i := range observed {
		if observed[i] > 0 {
			statistic += observed[i] * math.Log(observed[i]/expected[i])
		}
	}
	statistic *= 2

	df := float64(len(observed) - 1)
	chi2 := NewChiSquared(df)
	pvalue = chi2.SurvivalFunction(statistic)
	return statistic, pvalue
}
