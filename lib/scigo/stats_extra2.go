package scigo

import (
	"math"
	"sort"
)

// DescribeResult holds descriptive statistics for a dataset.
type DescribeResult struct {
	Nobs     int     // Number of observations
	Min      float64 // Minimum value
	Max      float64 // Maximum value
	Mean     float64 // Arithmetic mean
	Variance float64 // Sample variance (Bessel-corrected)
	Skewness float64 // Sample skewness
	Kurtosis float64 // Sample excess kurtosis
}

// Describe computes summary statistics for a dataset, analogous to scipy.stats.describe.
// Panics if data has fewer than 2 elements.
func Describe(data []float64) DescribeResult {
	n := len(data)
	if n < 2 {
		panic("scigo: Describe: need at least 2 data points")
	}

	fn := float64(n)
	mn := data[0]
	mx := data[0]
	sum := 0.0
	for _, v := range data {
		sum += v
		if v < mn {
			mn = v
		}
		if v > mx {
			mx = v
		}
	}
	mean := sum / fn

	m2, m3, m4 := 0.0, 0.0, 0.0
	for _, v := range data {
		d := v - mean
		d2 := d * d
		m2 += d2
		m3 += d2 * d
		m4 += d2 * d2
	}

	variance := m2 / (fn - 1) // Bessel-corrected

	// Sample skewness (adjusted Fisher-Pearson)
	skewness := 0.0
	if m2 > 0 {
		skewness = (m3 / fn) / math.Pow(m2/fn, 1.5)
	}

	// Sample excess kurtosis
	kurtosis := 0.0
	if m2 > 0 {
		kurtosis = (m4/fn)/(m2/fn*m2/fn) - 3
	}

	return DescribeResult{
		Nobs:     n,
		Min:      mn,
		Max:      mx,
		Mean:     mean,
		Variance: variance,
		Skewness: skewness,
		Kurtosis: kurtosis,
	}
}

// IQR computes the interquartile range of the data.
// IQR = Q3 - Q1 where Q1 and Q3 are the 25th and 75th percentiles.
// Panics if data is empty.
func IQR(data []float64) float64 {
	n := len(data)
	if n == 0 {
		panic("scigo: IQR: data must not be empty")
	}
	sorted := make([]float64, n)
	copy(sorted, data)
	sort.Float64s(sorted)

	q1 := percentile(sorted, 25)
	q3 := percentile(sorted, 75)
	return q3 - q1
}

// percentile computes the p-th percentile from sorted data using linear interpolation.
func percentile(sorted []float64, p float64) float64 {
	n := len(sorted)
	if n == 1 {
		return sorted[0]
	}
	// Use linear interpolation method (same as numpy default)
	rank := p / 100 * float64(n-1)
	lo := int(math.Floor(rank))
	hi := lo + 1
	if hi >= n {
		return sorted[n-1]
	}
	frac := rank - float64(lo)
	return sorted[lo]*(1-frac) + sorted[hi]*frac
}

// Zscore computes the z-scores of the data: z_i = (x_i - mean) / std.
// Panics if data has fewer than 2 elements.
func Zscore(data []float64) []float64 {
	n := len(data)
	if n < 2 {
		panic("scigo: Zscore: need at least 2 data points")
	}
	mean, variance := meanVar(data)
	std := math.Sqrt(variance)
	if std == 0 {
		result := make([]float64, n)
		return result
	}
	result := make([]float64, n)
	for i, v := range data {
		result[i] = (v - mean) / std
	}
	return result
}

// Zmap maps scores to z-scores using the mean and standard deviation of a
// comparison array. z_i = (scores_i - mean(compare)) / std(compare).
// Panics if compare has fewer than 2 elements.
func Zmap(scores, compare []float64) []float64 {
	if len(compare) < 2 {
		panic("scigo: Zmap: compare must have at least 2 data points")
	}
	mean, variance := meanVar(compare)
	std := math.Sqrt(variance)
	result := make([]float64, len(scores))
	if std == 0 {
		return result
	}
	for i, v := range scores {
		result[i] = (v - mean) / std
	}
	return result
}

// TrimMean computes the trimmed mean, removing the given proportional fraction
// from each end of the sorted data. proportiontocut should be in [0, 0.5).
// For example, 0.1 removes the bottom and top 10%.
// Panics if data is empty or proportiontocut is invalid.
func TrimMean(data []float64, proportiontocut float64) float64 {
	n := len(data)
	if n == 0 {
		panic("scigo: TrimMean: data must not be empty")
	}
	if proportiontocut < 0 || proportiontocut >= 0.5 {
		panic("scigo: TrimMean: proportiontocut must be in [0, 0.5)")
	}

	sorted := make([]float64, n)
	copy(sorted, data)
	sort.Float64s(sorted)

	ncut := int(math.Floor(proportiontocut * float64(n)))
	trimmed := sorted[ncut : n-ncut]
	if len(trimmed) == 0 {
		trimmed = sorted
	}

	sum := 0.0
	for _, v := range trimmed {
		sum += v
	}
	return sum / float64(len(trimmed))
}

// SEM computes the standard error of the mean: SEM = std(data) / sqrt(n).
// Uses the sample standard deviation (Bessel-corrected).
// Panics if data has fewer than 2 elements.
func SEM(data []float64) float64 {
	n := len(data)
	if n < 2 {
		panic("scigo: SEM: need at least 2 data points")
	}
	_, variance := meanVar(data)
	return math.Sqrt(variance / float64(n))
}

// Skew computes the sample skewness (Fisher-Pearson coefficient of skewness).
// Panics if data has fewer than 3 elements.
func Skew(data []float64) float64 {
	n := len(data)
	if n < 3 {
		panic("scigo: Skew: need at least 3 data points")
	}
	fn := float64(n)

	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / fn

	m2, m3 := 0.0, 0.0
	for _, v := range data {
		d := v - mean
		m2 += d * d
		m3 += d * d * d
	}
	m2 /= fn
	m3 /= fn

	if m2 == 0 {
		return 0
	}
	return m3 / math.Pow(m2, 1.5)
}

// Kurtosis computes the sample excess kurtosis (Fisher definition).
// kurtosis = m4/m2^2 - 3 where m2, m4 are the 2nd and 4th central moments.
// Panics if data has fewer than 4 elements.
func Kurtosis(data []float64) float64 {
	n := len(data)
	if n < 4 {
		panic("scigo: Kurtosis: need at least 4 data points")
	}
	fn := float64(n)

	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / fn

	m2, m4 := 0.0, 0.0
	for _, v := range data {
		d := v - mean
		d2 := d * d
		m2 += d2
		m4 += d2 * d2
	}
	m2 /= fn
	m4 /= fn

	if m2 == 0 {
		return 0
	}
	return m4/(m2*m2) - 3
}

// JarqueBera performs the Jarque-Bera test for normality.
// The test statistic is JB = (n/6) * (S^2 + (K^2)/4) where S is skewness
// and K is excess kurtosis.
// Returns the test statistic and p-value (chi-squared with 2 df).
// Panics if data has fewer than 8 elements.
func JarqueBera(data []float64) (statistic, pvalue float64) {
	n := len(data)
	if n < 8 {
		panic("scigo: JarqueBera: need at least 8 data points")
	}
	fn := float64(n)

	sum := 0.0
	for _, v := range data {
		sum += v
	}
	mean := sum / fn

	m2, m3, m4 := 0.0, 0.0, 0.0
	for _, v := range data {
		d := v - mean
		d2 := d * d
		m2 += d2
		m3 += d2 * d
		m4 += d2 * d2
	}
	m2 /= fn
	m3 /= fn
	m4 /= fn

	if m2 == 0 {
		return 0, 1
	}

	s := m3 / math.Pow(m2, 1.5) // skewness
	k := m4/(m2*m2) - 3         // excess kurtosis
	statistic = fn / 6 * (s*s + k*k/4)

	chi2 := NewChiSquared(2)
	pvalue = chi2.SurvivalFunction(statistic)
	return
}

// RankData assigns average ranks to the data (exported version).
// Ties are handled by averaging the ranks of tied values.
// This is equivalent to scipy.stats.rankdata with method='average'.
func RankData(data []float64) []float64 {
	return rankData(data)
}
