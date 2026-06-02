package scigo

import (
	"math"
	"sort"
)

// ---------------------------------------------------------------------------
// T-Tests
// ---------------------------------------------------------------------------

// TTestInd performs an independent two-sample t-test (Welch's t-test).
// It tests whether the means of two independent samples differ.
// Returns the t-statistic and two-tailed p-value.
// Panics if either sample has fewer than 2 elements.
func TTestInd(x, y []float64) (statistic, pvalue float64) {
	nx, ny := len(x), len(y)
	if nx < 2 || ny < 2 {
		panic("scigo: TTestInd: each sample must have at least 2 elements")
	}

	mx, vx := meanVar(x)
	my, vy := meanVar(y)

	fnx, fny := float64(nx), float64(ny)
	se := math.Sqrt(vx/fnx + vy/fny)
	if se == 0 {
		return 0, 1
	}
	statistic = (mx - my) / se

	// Welch-Satterthwaite degrees of freedom
	num := (vx/fnx + vy/fny) * (vx/fnx + vy/fny)
	denom := (vx*vx)/(fnx*fnx*(fnx-1)) + (vy*vy)/(fny*fny*(fny-1))
	if denom == 0 {
		return 0, 1
	}
	df := num / denom

	td := NewTDistribution(df)
	pvalue = 2 * td.SurvivalFunction(math.Abs(statistic))
	return
}

// TTest1Samp performs a one-sample t-test.
// It tests whether the mean of x equals mu.
// Returns the t-statistic and two-tailed p-value.
// Panics if x has fewer than 2 elements.
func TTest1Samp(x []float64, mu float64) (statistic, pvalue float64) {
	n := len(x)
	if n < 2 {
		panic("scigo: TTest1Samp: sample must have at least 2 elements")
	}

	m, v := meanVar(x)
	fn := float64(n)
	se := math.Sqrt(v / fn)
	if se == 0 {
		if m == mu {
			return 0, 1
		}
		return math.Inf(1), 0
	}
	statistic = (m - mu) / se
	df := fn - 1

	td := NewTDistribution(df)
	pvalue = 2 * td.SurvivalFunction(math.Abs(statistic))
	return
}

// TTestRel performs a paired (related) samples t-test.
// It tests whether the mean of the differences x[i]-y[i] is zero.
// Returns the t-statistic and two-tailed p-value.
// Panics if x and y have different lengths or fewer than 2 elements.
func TTestRel(x, y []float64) (statistic, pvalue float64) {
	n := len(x)
	if n != len(y) {
		panic("scigo: TTestRel: x and y must have the same length")
	}
	if n < 2 {
		panic("scigo: TTestRel: samples must have at least 2 elements")
	}

	d := make([]float64, n)
	for i := range d {
		d[i] = x[i] - y[i]
	}
	return TTest1Samp(d, 0)
}

// ---------------------------------------------------------------------------
// Non-parametric Tests
// ---------------------------------------------------------------------------

// MannWhitneyU performs the Mann-Whitney U test (Wilcoxon rank-sum test).
// It tests whether two independent samples come from the same distribution.
// Returns the U statistic and the p-value (normal approximation).
// Panics if either sample is empty.
func MannWhitneyU(x, y []float64) (statistic, pvalue float64) {
	nx, ny := len(x), len(y)
	if nx == 0 || ny == 0 {
		panic("scigo: MannWhitneyU: samples must not be empty")
	}

	// Pool and rank
	type valGroup struct {
		val   float64
		group int // 0 = x, 1 = y
	}
	all := make([]valGroup, 0, nx+ny)
	for _, v := range x {
		all = append(all, valGroup{v, 0})
	}
	for _, v := range y {
		all = append(all, valGroup{v, 1})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].val < all[j].val })

	// Assign average ranks
	n := nx + ny
	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i + 1
		for j < n && all[j].val == all[i].val {
			j++
		}
		avgRank := float64(i+j+1) / 2.0
		for k := i; k < j; k++ {
			ranks[k] = avgRank
		}
		i = j
	}

	// Sum ranks for group x
	r1 := 0.0
	for i, vg := range all {
		if vg.group == 0 {
			r1 += ranks[i]
		}
	}

	fnx, fny := float64(nx), float64(ny)
	u1 := r1 - fnx*(fnx+1)/2
	u2 := fnx*fny - u1
	statistic = math.Min(u1, u2)

	// Normal approximation
	mu := fnx * fny / 2
	sigma := math.Sqrt(fnx * fny * (fnx + fny + 1) / 12)
	if sigma == 0 {
		return statistic, 1
	}
	z := (statistic - mu) / sigma
	norm := NewNormal(0, 1)
	pvalue = 2 * norm.CDF(z) // z is negative since statistic <= mu
	return
}

// WilcoxonSignedRank performs the Wilcoxon signed-rank test for a single sample.
// It tests whether the median of x is zero.
// Returns the test statistic T+ and p-value (normal approximation).
// Zeros are excluded. Panics if x is empty.
func WilcoxonSignedRank(x []float64) (statistic, pvalue float64) {
	if len(x) == 0 {
		panic("scigo: WilcoxonSignedRank: sample must not be empty")
	}

	// Remove zeros
	type absSigned struct {
		absVal float64
		sign   float64
	}
	var data []absSigned
	for _, v := range x {
		if v != 0 {
			s := 1.0
			if v < 0 {
				s = -1.0
			}
			data = append(data, absSigned{math.Abs(v), s})
		}
	}

	n := len(data)
	if n == 0 {
		return 0, 1
	}

	// Sort by absolute value and assign ranks
	sort.Slice(data, func(i, j int) bool { return data[i].absVal < data[j].absVal })
	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i + 1
		for j < n && data[j].absVal == data[i].absVal {
			j++
		}
		avgRank := float64(i+j+1) / 2.0
		for k := i; k < j; k++ {
			ranks[k] = avgRank
		}
		i = j
	}

	// T+ = sum of ranks for positive values
	tPlus := 0.0
	for i, d := range data {
		if d.sign > 0 {
			tPlus += ranks[i]
		}
	}
	statistic = tPlus

	// Normal approximation
	fn := float64(n)
	mu := fn * (fn + 1) / 4
	sigma := math.Sqrt(fn * (fn + 1) * (2*fn + 1) / 24)
	if sigma == 0 {
		return statistic, 1
	}
	z := (statistic - mu) / sigma
	norm := NewNormal(0, 1)
	pvalue = 2 * math.Min(norm.CDF(z), norm.CDF(-z))
	return
}

// KruskalWallis performs the Kruskal-Wallis H test.
// It is a non-parametric alternative to one-way ANOVA for comparing
// the medians of two or more independent groups.
// Returns the H statistic and the p-value (chi-squared approximation).
// Panics if fewer than 2 groups or any group is empty.
func KruskalWallis(groups ...[]float64) (statistic, pvalue float64) {
	k := len(groups)
	if k < 2 {
		panic("scigo: KruskalWallis: need at least 2 groups")
	}
	for i, g := range groups {
		if len(g) == 0 {
			panic("scigo: KruskalWallis: group " + string(rune('0'+i)) + " is empty")
		}
	}

	// Pool all values
	type valGroup struct {
		val   float64
		group int
	}
	var all []valGroup
	ns := make([]int, k)
	for gi, g := range groups {
		ns[gi] = len(g)
		for _, v := range g {
			all = append(all, valGroup{v, gi})
		}
	}
	N := len(all)

	sort.Slice(all, func(i, j int) bool { return all[i].val < all[j].val })

	// Average ranks
	ranks := make([]float64, N)
	i := 0
	for i < N {
		j := i + 1
		for j < N && all[j].val == all[i].val {
			j++
		}
		avg := float64(i+j+1) / 2.0
		for m := i; m < j; m++ {
			ranks[m] = avg
		}
		i = j
	}

	// Sum of ranks per group
	rankSums := make([]float64, k)
	for i, vg := range all {
		rankSums[vg.group] += ranks[i]
	}

	// H statistic
	fN := float64(N)
	H := 0.0
	for gi := 0; gi < k; gi++ {
		ni := float64(ns[gi])
		H += (rankSums[gi] * rankSums[gi]) / ni
	}
	H = 12.0/(fN*(fN+1))*H - 3*(fN+1)

	statistic = H
	df := float64(k - 1)
	chi2 := NewChiSquared(df)
	pvalue = chi2.SurvivalFunction(statistic)
	return
}

// FriedmanChiSquare performs the Friedman chi-squared test.
// It is a non-parametric test for differences across repeated measures
// (randomized block design). Each group must have the same length (number of blocks).
// Returns the chi-squared statistic and p-value.
// Panics if fewer than 2 groups or groups have different lengths.
func FriedmanChiSquare(groups ...[]float64) (statistic, pvalue float64) {
	k := len(groups)
	if k < 2 {
		panic("scigo: FriedmanChiSquare: need at least 2 groups")
	}
	n := len(groups[0])
	if n == 0 {
		panic("scigo: FriedmanChiSquare: groups must not be empty")
	}
	for i := 1; i < k; i++ {
		if len(groups[i]) != n {
			panic("scigo: FriedmanChiSquare: all groups must have the same length")
		}
	}

	fk := float64(k)
	fn := float64(n)

	// Rank within each block (row)
	rankSums := make([]float64, k)
	for b := 0; b < n; b++ {
		// Get values for this block
		type idxVal struct {
			idx int
			val float64
		}
		row := make([]idxVal, k)
		for j := 0; j < k; j++ {
			row[j] = idxVal{j, groups[j][b]}
		}
		sort.Slice(row, func(i, j int) bool { return row[i].val < row[j].val })

		// Assign average ranks
		ranks := make([]float64, k)
		i := 0
		for i < k {
			j := i + 1
			for j < k && row[j].val == row[i].val {
				j++
			}
			avg := float64(i+j+1) / 2.0
			for m := i; m < j; m++ {
				ranks[row[m].idx] = avg
			}
			i = j
		}
		for j := 0; j < k; j++ {
			rankSums[j] += ranks[j]
		}
	}

	// Friedman statistic
	ssRanks := 0.0
	for j := 0; j < k; j++ {
		diff := rankSums[j] - fn*(fk+1)/2
		ssRanks += diff * diff
	}
	statistic = 12.0 / (fn * fk * (fk + 1)) * ssRanks

	df := fk - 1
	chi2 := NewChiSquared(df)
	pvalue = chi2.SurvivalFunction(statistic)
	return
}

// ---------------------------------------------------------------------------
// Kolmogorov-Smirnov Tests
// ---------------------------------------------------------------------------

// KSTest performs a one-sample Kolmogorov-Smirnov test.
// It tests whether x is drawn from the distribution described by cdf.
// cdf should be a continuous CDF function.
// Returns the KS statistic D and the p-value.
// Panics if x is empty.
func KSTest(x []float64, cdf func(float64) float64) (statistic, pvalue float64) {
	n := len(x)
	if n == 0 {
		panic("scigo: KSTest: sample must not be empty")
	}

	sorted := make([]float64, n)
	copy(sorted, x)
	sort.Float64s(sorted)

	fn := float64(n)
	d := 0.0
	for i, v := range sorted {
		fi := float64(i)
		cdfVal := cdf(v)
		// D+ = max(i/n - F(x_i))
		dp := (fi+1)/fn - cdfVal
		// D- = max(F(x_i) - (i-1)/n)
		dm := cdfVal - fi/fn
		if dp > d {
			d = dp
		}
		if dm > d {
			d = dm
		}
	}
	statistic = d
	pvalue = ksPvalue(d, n)
	return
}

// KS2Samp performs a two-sample Kolmogorov-Smirnov test.
// It tests whether x and y come from the same continuous distribution.
// Returns the KS statistic D and the p-value.
// Panics if either sample is empty.
func KS2Samp(x, y []float64) (statistic, pvalue float64) {
	nx, ny := len(x), len(y)
	if nx == 0 || ny == 0 {
		panic("scigo: KS2Samp: samples must not be empty")
	}

	sx := make([]float64, nx)
	copy(sx, x)
	sort.Float64s(sx)

	sy := make([]float64, ny)
	copy(sy, y)
	sort.Float64s(sy)

	fnx, fny := float64(nx), float64(ny)
	d := 0.0
	i, j := 0, 0
	for i < nx && j < ny {
		var val float64
		if sx[i] <= sy[j] {
			val = sx[i]
			i++
		} else {
			val = sy[j]
			j++
		}
		// Advance past ties
		for i < nx && sx[i] == val {
			i++
		}
		for j < ny && sy[j] == val {
			j++
		}
		diff := math.Abs(float64(i)/fnx - float64(j)/fny)
		if diff > d {
			d = diff
		}
	}
	statistic = d

	// Effective sample size for p-value
	en := math.Sqrt(fnx * fny / (fnx + fny))
	pvalue = ksProb((en + 0.12 + 0.11/en) * statistic)
	return
}

// ksPvalue computes the p-value for the one-sample KS test using the
// asymptotic formula with finite-sample correction.
func ksPvalue(d float64, n int) float64 {
	fn := float64(n)
	en := math.Sqrt(fn)
	return ksProb((en + 0.12 + 0.11/en) * d)
}

// ksProb computes the Kolmogorov-Smirnov survival function P(D > d)
// using the series approximation: Q(lambda) = 2 * sum_{j=1}^{inf} (-1)^{j-1} exp(-2*j^2*lambda^2)
func ksProb(lam float64) float64 {
	if lam <= 0 {
		return 1
	}
	if lam > 8 {
		return 0
	}
	sum := 0.0
	for j := 1; j <= 100; j++ {
		fj := float64(j)
		term := math.Exp(-2 * fj * fj * lam * lam)
		if j%2 == 0 {
			sum -= term
		} else {
			sum += term
		}
		if term < 1e-15 {
			break
		}
	}
	return 2 * sum
}

// ---------------------------------------------------------------------------
// Normality Tests
// ---------------------------------------------------------------------------

// ShapiroWilk performs a simplified Shapiro-Wilk test for normality.
// This implementation uses the approximation method suitable for moderate sample sizes.
// Returns the W statistic and an approximate p-value.
// Panics if x has fewer than 3 elements or more than 5000.
func ShapiroWilk(x []float64) (statistic, pvalue float64) {
	n := len(x)
	if n < 3 {
		panic("scigo: ShapiroWilk: need at least 3 data points")
	}
	if n > 5000 {
		panic("scigo: ShapiroWilk: sample size must not exceed 5000")
	}

	sorted := make([]float64, n)
	copy(sorted, x)
	sort.Float64s(sorted)

	// Compute mean and SS
	mean := 0.0
	for _, v := range sorted {
		mean += v
	}
	mean /= float64(n)

	ss := 0.0
	for _, v := range sorted {
		d := v - mean
		ss += d * d
	}
	if ss == 0 {
		return 1, 1 // constant data
	}

	// Compute expected normal order statistics (Blom's approximation)
	fn := float64(n)
	m := make([]float64, n)
	norm := NewNormal(0, 1)
	for i := 0; i < n; i++ {
		p := (float64(i) + 0.375) / (fn + 0.25)
		m[i] = norm.PPF(p)
	}

	// Compute sum of m^2
	mss := 0.0
	for _, v := range m {
		mss += v * v
	}

	// Compute a_i coefficients (simplified)
	a := make([]float64, n)
	for i := 0; i < n; i++ {
		a[i] = m[i] / math.Sqrt(mss)
	}

	// W = (sum(a_i * x_(i)))^2 / SS
	num := 0.0
	for i := 0; i < n; i++ {
		num += a[i] * sorted[i]
	}
	statistic = (num * num) / ss

	// Approximate p-value using transformation to normality
	// Based on Royston (1992) approximation
	if statistic > 1 {
		statistic = 1
	}
	if statistic > 0.999 {
		pvalue = 1
		return
	}

	lnN := math.Log(fn)
	if fn <= 11 {
		// Small sample: use -log(1-W) transformation
		gamma := 0.459*fn - 2.273
		y := -math.Log(1 - statistic)
		mu := -1.2725 + 1.0521*(gamma-2.0)/gamma
		sigma := 1.0308 - 0.26758*(gamma-2.0)/gamma
		z := (math.Log(y) - mu) / sigma
		normDist := NewNormal(0, 1)
		pvalue = 1 - normDist.CDF(z)
	} else {
		// Larger sample: use log(1-W) transformation
		y := math.Log(1 - statistic)
		mu := 0.0038915*lnN*lnN*lnN - 0.083751*lnN*lnN - 0.31082*lnN - 1.5861
		sigma := math.Exp(0.0030302*lnN*lnN - 0.082676*lnN - 0.4803)
		z := (y - mu) / sigma
		normDist := NewNormal(0, 1)
		pvalue = 1 - normDist.CDF(z)
	}

	if pvalue > 1 {
		pvalue = 1
	}
	if pvalue < 0 {
		pvalue = 0
	}
	return
}

// NormalTest performs D'Agostino-Pearson omnibus normality test.
// It combines the skewness and kurtosis tests into a single statistic.
// Returns the K^2 statistic and p-value (chi-squared with 2 df).
// Panics if x has fewer than 20 elements (needed for reliable estimates).
func NormalTest(x []float64) (statistic, pvalue float64) {
	n := len(x)
	if n < 20 {
		panic("scigo: NormalTest: need at least 20 data points")
	}

	fn := float64(n)

	// Compute mean
	mean := 0.0
	for _, v := range x {
		mean += v
	}
	mean /= fn

	// Compute central moments
	m2, m3, m4 := 0.0, 0.0, 0.0
	for _, v := range x {
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

	// Sample skewness (G1) and kurtosis (G2)
	g1 := m3 / math.Pow(m2, 1.5)
	g2 := m4/(m2*m2) - 3

	// D'Agostino's skewness test: transform G1 to Z1
	y := g1 * math.Sqrt((fn+1)*(fn+3)/(6*(fn-2)))
	beta2 := 3 * (fn*fn + 27*fn - 70) * (fn + 1) * (fn + 3) / ((fn - 2) * (fn + 5) * (fn + 7) * (fn + 9))
	w2 := -1 + math.Sqrt(2*(beta2-1))
	delta := 1.0 / math.Sqrt(math.Log(math.Sqrt(w2)))
	alpha := math.Sqrt(2.0 / (w2 - 1))
	yAlpha := y / alpha
	z1 := delta * math.Log(yAlpha+math.Sqrt(yAlpha*yAlpha+1))

	// Anscombe-Glynn kurtosis test: transform G2 to Z2
	meanG2 := -6.0 / (fn + 1)
	varG2 := 24 * fn * (fn - 2) * (fn - 3) / ((fn + 1) * (fn + 1) * (fn + 3) * (fn + 5))
	x0 := (g2 - meanG2) / math.Sqrt(varG2)
	sqrtBeta1 := 6 * (fn*fn - 5*fn + 2) / ((fn + 7) * (fn + 9)) * math.Sqrt(6*(fn+3)*(fn+5)/(fn*(fn-2)*(fn-3)))
	A := 6 + 8/sqrtBeta1*(2/sqrtBeta1+math.Sqrt(1+4/(sqrtBeta1*sqrtBeta1)))
	denom := 1 - 2.0/A
	if denom <= 0 {
		denom = 1e-10
	}
	num := 1 - 2.0/(9*A)
	cube := num - math.Pow(denom/(1+x0*math.Sqrt(2.0/(A-4))), 1.0/3.0)
	z2 := cube / math.Sqrt(2.0/(9*A))

	statistic = z1*z1 + z2*z2
	chi2 := NewChiSquared(2)
	pvalue = chi2.SurvivalFunction(statistic)
	return
}

// AndersonDarling performs the Anderson-Darling test for normality.
// It returns the A^2 statistic and critical values at significance levels
// [15%, 10%, 5%, 2.5%, 1%]. If A^2 exceeds a critical value, the null
// hypothesis of normality is rejected at that significance level.
// Panics if x has fewer than 7 elements.
func AndersonDarling(x []float64) (statistic float64, critical []float64) {
	n := len(x)
	if n < 7 {
		panic("scigo: AndersonDarling: need at least 7 data points")
	}

	sorted := make([]float64, n)
	copy(sorted, x)
	sort.Float64s(sorted)

	// Standardize: compute mean and std
	mean, std := 0.0, 0.0
	for _, v := range sorted {
		mean += v
	}
	mean /= float64(n)
	for _, v := range sorted {
		d := v - mean
		std += d * d
	}
	std = math.Sqrt(std / float64(n-1))
	if std == 0 {
		return math.Inf(1), []float64{0.576, 0.656, 0.787, 0.918, 1.092}
	}

	z := make([]float64, n)
	norm := NewNormal(0, 1)
	for i, v := range sorted {
		z[i] = norm.CDF((v - mean) / std)
		if z[i] <= 0 {
			z[i] = 1e-15
		}
		if z[i] >= 1 {
			z[i] = 1 - 1e-15
		}
	}

	fn := float64(n)
	s := 0.0
	for i := 0; i < n; i++ {
		fi := float64(i)
		s += (2*fi + 1) * (math.Log(z[i]) + math.Log(1-z[n-1-i]))
	}
	statistic = -fn - s/fn

	// Apply correction factor for finite sample size
	statistic *= (1 + 0.75/fn + 2.25/(fn*fn))

	// Critical values for normal distribution test at 15%, 10%, 5%, 2.5%, 1%
	critical = []float64{0.576, 0.656, 0.787, 0.918, 1.092}
	return
}

// ---------------------------------------------------------------------------
// Variance Tests
// ---------------------------------------------------------------------------

// BartlettTest performs Bartlett's test for equal variances.
// It tests whether all groups have the same variance.
// Assumes the data is normally distributed.
// Returns the test statistic and p-value (chi-squared distribution).
// Panics if fewer than 2 groups or any group has fewer than 2 elements.
func BartlettTest(groups ...[]float64) (statistic, pvalue float64) {
	k := len(groups)
	if k < 2 {
		panic("scigo: BartlettTest: need at least 2 groups")
	}

	ns := make([]int, k)
	vars := make([]float64, k)
	N := 0
	for i, g := range groups {
		if len(g) < 2 {
			panic("scigo: BartlettTest: each group must have at least 2 elements")
		}
		ns[i] = len(g)
		N += ns[i]
		_, vars[i] = meanVar(g)
	}

	// Pooled variance
	sp2 := 0.0
	for i := 0; i < k; i++ {
		sp2 += float64(ns[i]-1) * vars[i]
	}
	fNk := float64(N - k)
	sp2 /= fNk

	if sp2 == 0 {
		return 0, 1
	}

	// Bartlett's statistic
	num := 0.0
	sumInv := 0.0
	for i := 0; i < k; i++ {
		ni := float64(ns[i] - 1)
		if vars[i] > 0 {
			num += ni * math.Log(vars[i])
		} else {
			num += ni * math.Log(1e-300)
		}
		sumInv += 1.0 / ni
	}
	num = fNk*math.Log(sp2) - num
	denom := 1 + (sumInv-1.0/fNk)/(3*float64(k-1))
	statistic = num / denom

	df := float64(k - 1)
	chi2 := NewChiSquared(df)
	pvalue = chi2.SurvivalFunction(statistic)
	return
}

// LeveneTest performs Levene's test for equal variances.
// Uses the mean-based version (Brown-Forsythe uses medians).
// Returns the test statistic and p-value (F distribution approximation).
// Panics if fewer than 2 groups or any group has fewer than 2 elements.
func LeveneTest(groups ...[]float64) (statistic, pvalue float64) {
	k := len(groups)
	if k < 2 {
		panic("scigo: LeveneTest: need at least 2 groups")
	}

	N := 0
	means := make([]float64, k)
	ns := make([]int, k)
	for i, g := range groups {
		if len(g) < 2 {
			panic("scigo: LeveneTest: each group must have at least 2 elements")
		}
		ns[i] = len(g)
		N += ns[i]
		s := 0.0
		for _, v := range g {
			s += v
		}
		means[i] = s / float64(ns[i])
	}

	// Compute |x_ij - mean_i| for each observation
	zGroups := make([][]float64, k)
	for i, g := range groups {
		zGroups[i] = make([]float64, len(g))
		for j, v := range g {
			zGroups[i][j] = math.Abs(v - means[i])
		}
	}

	// Mean of z values per group and overall
	zMeans := make([]float64, k)
	zOverall := 0.0
	for i, zg := range zGroups {
		s := 0.0
		for _, v := range zg {
			s += v
		}
		zMeans[i] = s / float64(ns[i])
		zOverall += s
	}
	zOverall /= float64(N)

	// F statistic
	fN := float64(N)
	fk := float64(k)
	numSS := 0.0
	for i := 0; i < k; i++ {
		d := zMeans[i] - zOverall
		numSS += float64(ns[i]) * d * d
	}
	denomSS := 0.0
	for i, zg := range zGroups {
		for _, v := range zg {
			d := v - zMeans[i]
			denomSS += d * d
		}
	}

	if denomSS == 0 {
		return 0, 1
	}

	statistic = ((fN - fk) / (fk - 1)) * (numSS / denomSS)
	pvalue = fDistSurvival(statistic, fk-1, fN-fk)
	return
}

// FlignerKilleen performs Fligner-Killeen's test for equal variances.
// This is a non-parametric test that uses ranks of absolute deviations from
// group medians. Returns the test statistic and p-value (chi-squared approximation).
// Panics if fewer than 2 groups or any group has fewer than 2 elements.
func FlignerKilleen(groups ...[]float64) (statistic, pvalue float64) {
	k := len(groups)
	if k < 2 {
		panic("scigo: FlignerKilleen: need at least 2 groups")
	}

	N := 0
	ns := make([]int, k)
	medians := make([]float64, k)
	for i, g := range groups {
		if len(g) < 2 {
			panic("scigo: FlignerKilleen: each group must have at least 2 elements")
		}
		ns[i] = len(g)
		N += ns[i]
		medians[i] = median(g)
	}

	// Compute |x_ij - median_i|
	type absDevGroup struct {
		absDev float64
		group  int
	}
	allDevs := make([]absDevGroup, 0, N)
	for i, g := range groups {
		for _, v := range g {
			allDevs = append(allDevs, absDevGroup{math.Abs(v - medians[i]), i})
		}
	}

	// Rank the absolute deviations
	sort.Slice(allDevs, func(i, j int) bool { return allDevs[i].absDev < allDevs[j].absDev })
	ranks := make([]float64, N)
	i := 0
	for i < N {
		j := i + 1
		for j < N && allDevs[j].absDev == allDevs[i].absDev {
			j++
		}
		avg := float64(i+j+1) / 2.0
		for m := i; m < j; m++ {
			ranks[m] = avg
		}
		i = j
	}

	// Convert ranks to normal scores a_i = Phi^{-1}((rank_i/(N+1) + 1)/2)
	norm := NewNormal(0, 1)
	fN := float64(N)
	scores := make([]float64, N)
	for i := 0; i < N; i++ {
		p := (ranks[i]/(fN+1) + 1) / 2
		scores[i] = norm.PPF(p)
	}

	// Mean of all scores
	aMean := 0.0
	for _, s := range scores {
		aMean += s
	}
	aMean /= fN

	// Sum of scores per group
	groupScoreSums := make([]float64, k)
	idx := 0
	// We need to track which score goes to which group
	for i := 0; i < N; i++ {
		groupScoreSums[allDevs[i].group] += scores[i]
	}
	_ = idx

	// Group means of scores
	groupScoreMeans := make([]float64, k)
	for i := 0; i < k; i++ {
		groupScoreMeans[i] = groupScoreSums[i] / float64(ns[i])
	}

	// Variance of all scores
	sVar := 0.0
	for _, s := range scores {
		d := s - aMean
		sVar += d * d
	}
	sVar /= (fN - 1)

	if sVar == 0 {
		return 0, 1
	}

	// Test statistic
	statistic = 0
	for i := 0; i < k; i++ {
		d := groupScoreMeans[i] - aMean
		statistic += float64(ns[i]) * d * d
	}
	statistic /= sVar

	df := float64(k - 1)
	chi2 := NewChiSquared(df)
	pvalue = chi2.SurvivalFunction(statistic)
	return
}

// ---------------------------------------------------------------------------
// Mood's Test
// ---------------------------------------------------------------------------

// MoodTest performs Mood's test for equal scale parameters.
// It tests whether two samples have the same scale (dispersion).
// Returns the Z statistic and p-value (normal approximation).
// Panics if either sample has fewer than 3 elements.
func MoodTest(x, y []float64) (statistic, pvalue float64) {
	nx, ny := len(x), len(y)
	if nx < 3 || ny < 3 {
		panic("scigo: MoodTest: each sample must have at least 3 elements")
	}

	// Pool and rank
	N := nx + ny
	type valGroup struct {
		val   float64
		group int // 0=x
	}
	all := make([]valGroup, 0, N)
	for _, v := range x {
		all = append(all, valGroup{v, 0})
	}
	for _, v := range y {
		all = append(all, valGroup{v, 1})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].val < all[j].val })

	ranks := make([]float64, N)
	i := 0
	for i < N {
		j := i + 1
		for j < N && all[j].val == all[i].val {
			j++
		}
		avg := float64(i+j+1) / 2.0
		for m := i; m < j; m++ {
			ranks[m] = avg
		}
		i = j
	}

	// M = sum of (rank_i - (N+1)/2)^2 for sample x
	fN := float64(N)
	center := (fN + 1) / 2
	M := 0.0
	for i, vg := range all {
		if vg.group == 0 {
			d := ranks[i] - center
			M += d * d
		}
	}

	// Expected value and variance under H0
	fnx := float64(nx)
	eM := fnx * (fN*fN - 1) / 12
	varM := fnx * float64(ny) * (fN + 1) * (fN*fN - 4) / 180

	if varM == 0 {
		return 0, 1
	}

	z := (M - eM) / math.Sqrt(varM)
	statistic = z
	norm := NewNormal(0, 1)
	pvalue = 2 * math.Min(norm.CDF(z), norm.CDF(-z))
	return
}

// ---------------------------------------------------------------------------
// Correlation Measures
// ---------------------------------------------------------------------------

// SpearmanR computes the Spearman rank correlation coefficient and its p-value.
// It is a non-parametric measure of monotonic association.
// Returns the correlation coefficient r and two-tailed p-value.
// Panics if x and y have different lengths or fewer than 3 elements.
func SpearmanR(x, y []float64) (r, pvalue float64) {
	n := len(x)
	if n != len(y) {
		panic("scigo: SpearmanR: x and y must have the same length")
	}
	if n < 3 {
		panic("scigo: SpearmanR: need at least 3 data points")
	}

	rx := rankData(x)
	ry := rankData(y)
	return PearsonCorrelation(rx, ry)
}

// KendallTau computes the Kendall rank correlation coefficient (tau-b) and its p-value.
// Returns the tau statistic and two-tailed p-value (normal approximation).
// Panics if x and y have different lengths or fewer than 3 elements.
func KendallTau(x, y []float64) (tau, pvalue float64) {
	n := len(x)
	if n != len(y) {
		panic("scigo: KendallTau: x and y must have the same length")
	}
	if n < 3 {
		panic("scigo: KendallTau: need at least 3 data points")
	}

	concordant, discordant := 0, 0
	tiedX, tiedY := 0, 0
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			dx := x[i] - x[j]
			dy := y[i] - y[j]
			if dx == 0 && dy == 0 {
				tiedX++
				tiedY++
			} else if dx == 0 {
				tiedX++
			} else if dy == 0 {
				tiedY++
			} else if (dx > 0 && dy > 0) || (dx < 0 && dy < 0) {
				concordant++
			} else {
				discordant++
			}
		}
	}

	n0 := n * (n - 1) / 2
	denomX := math.Sqrt(float64(n0 - tiedX))
	denomY := math.Sqrt(float64(n0 - tiedY))
	if denomX == 0 || denomY == 0 {
		return 0, 1
	}
	tau = float64(concordant-discordant) / (denomX * denomY)

	// Normal approximation for p-value
	fn := float64(n)
	var0 := (2 * (2*fn + 5)) / (9 * fn * (fn - 1))
	z := tau / math.Sqrt(var0)
	norm := NewNormal(0, 1)
	pvalue = 2 * math.Min(norm.CDF(z), norm.CDF(-z))
	return
}

// ---------------------------------------------------------------------------
// Linear Regression
// ---------------------------------------------------------------------------

// Linregress performs simple linear regression of y on x.
// Returns slope, intercept, Pearson r, two-tailed p-value for the slope,
// and standard error of the slope estimate.
// Panics if x and y have different lengths or fewer than 3 elements.
func Linregress(x, y []float64) (slope, intercept, r, pvalue, stderr float64) {
	n := len(x)
	if n != len(y) {
		panic("scigo: Linregress: x and y must have the same length")
	}
	if n < 3 {
		panic("scigo: Linregress: need at least 3 data points")
	}

	fn := float64(n)
	mx, my := 0.0, 0.0
	for i := range x {
		mx += x[i]
		my += y[i]
	}
	mx /= fn
	my /= fn

	var sxx, sxy, syy float64
	for i := range x {
		dx := x[i] - mx
		dy := y[i] - my
		sxx += dx * dx
		sxy += dx * dy
		syy += dy * dy
	}

	if sxx == 0 {
		return 0, my, 0, 1, math.Inf(1)
	}

	slope = sxy / sxx
	intercept = my - slope*mx

	if syy == 0 {
		r = 0
	} else {
		r = sxy / math.Sqrt(sxx*syy)
	}
	if r > 1 {
		r = 1
	} else if r < -1 {
		r = -1
	}

	// Residual sum of squares
	rss := 0.0
	for i := range x {
		resid := y[i] - (slope*x[i] + intercept)
		rss += resid * resid
	}

	df := fn - 2
	if df <= 0 {
		return slope, intercept, r, math.NaN(), math.NaN()
	}
	mse := rss / df
	stderr = math.Sqrt(mse / sxx)

	if stderr == 0 {
		pvalue = 0
	} else {
		tstat := slope / stderr
		td := NewTDistribution(df)
		pvalue = 2 * td.SurvivalFunction(math.Abs(tstat))
	}
	return
}

// ---------------------------------------------------------------------------
// Information Theory
// ---------------------------------------------------------------------------

// Entropy computes the Shannon entropy of a probability distribution pk.
// If qk is provided (non-nil and same length), computes the Kullback-Leibler
// divergence D_KL(pk || qk) instead.
// Probabilities should be non-negative. Zero probabilities are handled as 0*log(0) = 0.
// Panics if pk is empty or if qk is provided with a different length.
func Entropy(pk, qk []float64) float64 {
	if len(pk) == 0 {
		panic("scigo: Entropy: pk must not be empty")
	}
	if qk != nil && len(qk) != len(pk) {
		panic("scigo: Entropy: pk and qk must have the same length")
	}

	if qk == nil {
		// Shannon entropy: H = -sum(p * log(p))
		h := 0.0
		for _, p := range pk {
			if p > 0 {
				h -= p * math.Log(p)
			}
		}
		return h
	}

	// KL divergence: D_KL = sum(p * log(p/q))
	kl := 0.0
	for i, p := range pk {
		if p > 0 {
			if qk[i] <= 0 {
				return math.Inf(1)
			}
			kl += p * math.Log(p/qk[i])
		}
	}
	return kl
}

// ---------------------------------------------------------------------------
// Wasserstein Distance
// ---------------------------------------------------------------------------

// WassersteinDistance computes the first Wasserstein distance (earth mover's distance)
// between two 1D distributions given as samples u and v.
// This is equivalent to the integral of the absolute difference of the CDFs.
// Panics if either sample is empty.
func WassersteinDistance(u, v []float64) float64 {
	nu, nv := len(u), len(v)
	if nu == 0 || nv == 0 {
		panic("scigo: WassersteinDistance: samples must not be empty")
	}

	su := make([]float64, nu)
	copy(su, u)
	sort.Float64s(su)

	sv := make([]float64, nv)
	copy(sv, v)
	sort.Float64s(sv)

	// Merge all unique values
	all := make([]float64, 0, nu+nv)
	all = append(all, su...)
	all = append(all, sv...)
	sort.Float64s(all)

	// Compute the integral of |CDF_u - CDF_v|
	dist := 0.0
	iu, iv := 0, 0
	fnu, fnv := float64(nu), float64(nv)
	prev := all[0]
	for _, val := range all {
		if val > prev {
			cdfU := float64(iu) / fnu
			cdfV := float64(iv) / fnv
			dist += math.Abs(cdfU-cdfV) * (val - prev)
			prev = val
		}
		for iu < nu && su[iu] <= val {
			iu++
		}
		for iv < nv && sv[iv] <= val {
			iv++
		}
	}
	return dist
}

// ---------------------------------------------------------------------------
// Helper Functions
// ---------------------------------------------------------------------------

// meanVar computes the sample mean and sample variance (Bessel-corrected) of x.
func meanVar(x []float64) (mean, variance float64) {
	n := float64(len(x))
	for _, v := range x {
		mean += v
	}
	mean /= n
	for _, v := range x {
		d := v - mean
		variance += d * d
	}
	if n > 1 {
		variance /= (n - 1)
	}
	return
}

// median returns the median of x (does not modify x).
func median(x []float64) float64 {
	sorted := make([]float64, len(x))
	copy(sorted, x)
	sort.Float64s(sorted)
	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

// rankData assigns average ranks to the data.
func rankData(x []float64) []float64 {
	n := len(x)
	type idxVal struct {
		idx int
		val float64
	}
	data := make([]idxVal, n)
	for i, v := range x {
		data[i] = idxVal{i, v}
	}
	sort.Slice(data, func(i, j int) bool { return data[i].val < data[j].val })

	ranks := make([]float64, n)
	i := 0
	for i < n {
		j := i + 1
		for j < n && data[j].val == data[i].val {
			j++
		}
		avg := float64(i+j+1) / 2.0
		for k := i; k < j; k++ {
			ranks[data[k].idx] = avg
		}
		i = j
	}
	return ranks
}

// fDistSurvival computes 1 - CDF of the F-distribution with df1, df2 degrees of freedom at x.
// Uses the relationship: F CDF = I_{df1*x/(df1*x+df2)}(df1/2, df2/2)
func fDistSurvival(x, df1, df2 float64) float64 {
	if x <= 0 {
		return 1
	}
	bt := df1 * x / (df1*x + df2)
	return 1 - RegularizedIncompleteBeta(bt, df1/2, df2/2)
}
