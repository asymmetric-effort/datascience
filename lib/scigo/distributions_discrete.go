package scigo

import "math"

// ---------------------------------------------------------------------------
// Negative Binomial Distribution
// ---------------------------------------------------------------------------

// NegativeBinomial represents a negative binomial distribution.
// Models the number of failures before r successes, with success probability p.
type NegativeBinomial struct {
	r float64
	p float64
}

// NewNegativeBinomial creates a NegativeBinomial distribution.
// Panics if r <= 0 or p is not in (0, 1].
func NewNegativeBinomial(r, p float64) *NegativeBinomial {
	if r <= 0 {
		panic("scigo: NegativeBinomial r must be positive")
	}
	if p <= 0 || p > 1 {
		panic("scigo: NegativeBinomial p must be in (0, 1]")
	}
	return &NegativeBinomial{r: r, p: p}
}

// PMF returns the probability mass function at k (number of failures).
func (nb *NegativeBinomial) PMF(k int) float64 {
	if k < 0 {
		return 0
	}
	// PMF = C(k+r-1, k) * p^r * (1-p)^k
	logC := Gammaln(float64(k)+nb.r) - Gammaln(float64(k)+1) - Gammaln(nb.r)
	return math.Exp(logC + nb.r*math.Log(nb.p) + float64(k)*math.Log(1-nb.p))
}

// CDF returns the cumulative distribution function at k.
func (nb *NegativeBinomial) CDF(k int) float64 {
	if k < 0 {
		return 0
	}
	// CDF(k) = I_p(r, k+1) using the regularized incomplete beta function
	return RegularizedIncompleteBeta(nb.p, nb.r, float64(k)+1)
}

// Mean returns the mean: r*(1-p)/p.
func (nb *NegativeBinomial) Mean() float64 {
	return nb.r * (1 - nb.p) / nb.p
}

// Var returns the variance: r*(1-p)/p^2.
func (nb *NegativeBinomial) Var() float64 {
	return nb.r * (1 - nb.p) / (nb.p * nb.p)
}

// ---------------------------------------------------------------------------
// Geometric Distribution
// ---------------------------------------------------------------------------

// Geometric represents a geometric distribution with success probability p.
// Models the number of failures before the first success (support: k = 0, 1, 2, ...).
type Geometric struct {
	p float64
}

// NewGeometric creates a Geometric distribution. Panics if p is not in (0, 1].
func NewGeometric(p float64) *Geometric {
	if p <= 0 || p > 1 {
		panic("scigo: Geometric p must be in (0, 1]")
	}
	return &Geometric{p: p}
}

// PMF returns the probability mass function at k (number of failures before first success).
func (g *Geometric) PMF(k int) float64 {
	if k < 0 {
		return 0
	}
	return g.p * math.Pow(1-g.p, float64(k))
}

// CDF returns the cumulative distribution function at k.
func (g *Geometric) CDF(k int) float64 {
	if k < 0 {
		return 0
	}
	return 1 - math.Pow(1-g.p, float64(k)+1)
}

// Mean returns the mean: (1-p)/p.
func (g *Geometric) Mean() float64 {
	return (1 - g.p) / g.p
}

// Var returns the variance: (1-p)/p^2.
func (g *Geometric) Var() float64 {
	return (1 - g.p) / (g.p * g.p)
}

// ---------------------------------------------------------------------------
// Hypergeometric Distribution
// ---------------------------------------------------------------------------

// Hypergeometric represents a hypergeometric distribution.
// N = population size, K = number of success states, n = number of draws.
type Hypergeometric struct {
	popN  int
	succK int
	draws int
}

// NewHypergeometric creates a Hypergeometric distribution.
// Panics if parameters are invalid.
func NewHypergeometric(N, K, n int) *Hypergeometric {
	if N < 0 {
		panic("scigo: Hypergeometric N must be non-negative")
	}
	if K < 0 || K > N {
		panic("scigo: Hypergeometric K must be in [0, N]")
	}
	if n < 0 || n > N {
		panic("scigo: Hypergeometric n must be in [0, N]")
	}
	return &Hypergeometric{popN: N, succK: K, draws: n}
}

// PMF returns the probability mass function at k.
func (h *Hypergeometric) PMF(k int) float64 {
	N, K, n := h.popN, h.succK, h.draws
	lo := max(0, n-(N-K))
	hi := min(K, n)
	if k < lo || k > hi {
		return 0
	}
	// C(K,k) * C(N-K, n-k) / C(N,n)
	logP := Gammaln(float64(K)+1) - Gammaln(float64(k)+1) - Gammaln(float64(K-k)+1) +
		Gammaln(float64(N-K)+1) - Gammaln(float64(n-k)+1) - Gammaln(float64(N-K-n+k)+1) -
		Gammaln(float64(N)+1) + Gammaln(float64(n)+1) + Gammaln(float64(N-n)+1)
	return math.Exp(logP)
}

// CDF returns the cumulative distribution function at k.
func (h *Hypergeometric) CDF(k int) float64 {
	N, K, n := h.popN, h.succK, h.draws
	lo := max(0, n-(N-K))
	if k < lo {
		return 0
	}
	hi := min(K, n)
	if k >= hi {
		return 1
	}
	sum := 0.0
	for i := lo; i <= k; i++ {
		sum += h.PMF(i)
	}
	return sum
}

// Mean returns the mean: n*K/N.
func (h *Hypergeometric) Mean() float64 {
	return float64(h.draws) * float64(h.succK) / float64(h.popN)
}

// Var returns the variance: n*K*(N-K)*(N-n) / (N^2*(N-1)).
func (h *Hypergeometric) Var() float64 {
	N := float64(h.popN)
	K := float64(h.succK)
	n := float64(h.draws)
	if N <= 1 {
		return 0
	}
	return n * K * (N - K) * (N - n) / (N * N * (N - 1))
}

// ---------------------------------------------------------------------------
// Bernoulli Distribution
// ---------------------------------------------------------------------------

// Bernoulli represents a Bernoulli distribution (single trial with probability p).
type Bernoulli struct {
	p float64
}

// NewBernoulli creates a Bernoulli distribution. Panics if p is not in [0, 1].
func NewBernoulli(p float64) *Bernoulli {
	if p < 0 || p > 1 {
		panic("scigo: Bernoulli p must be in [0, 1]")
	}
	return &Bernoulli{p: p}
}

// PMF returns the probability mass function at k (0 or 1).
func (b *Bernoulli) PMF(k int) float64 {
	if k == 0 {
		return 1 - b.p
	}
	if k == 1 {
		return b.p
	}
	return 0
}

// CDF returns the cumulative distribution function at k.
func (b *Bernoulli) CDF(k int) float64 {
	if k < 0 {
		return 0
	}
	if k >= 1 {
		return 1
	}
	return 1 - b.p
}

// Mean returns the mean: p.
func (b *Bernoulli) Mean() float64 {
	return b.p
}

// Var returns the variance: p*(1-p).
func (b *Bernoulli) Var() float64 {
	return b.p * (1 - b.p)
}

// ---------------------------------------------------------------------------
// Zipf Distribution
// ---------------------------------------------------------------------------

// Zipf represents a Zipf distribution with exponent s over ranks 1..n.
type Zipf struct {
	s float64
	n int
	h float64 // generalized harmonic number H_{n,s}
}

// NewZipf creates a Zipf distribution. Panics if s <= 0 or n <= 0.
func NewZipf(s float64, n int) *Zipf {
	if s <= 0 {
		panic("scigo: Zipf s must be positive")
	}
	if n <= 0 {
		panic("scigo: Zipf n must be positive")
	}
	h := 0.0
	for k := 1; k <= n; k++ {
		h += math.Pow(float64(k), -s)
	}
	return &Zipf{s: s, n: n, h: h}
}

// PMF returns the probability mass function at k (rank 1..n).
func (z *Zipf) PMF(k int) float64 {
	if k < 1 || k > z.n {
		return 0
	}
	return math.Pow(float64(k), -z.s) / z.h
}

// CDF returns the cumulative distribution function at k.
func (z *Zipf) CDF(k int) float64 {
	if k < 1 {
		return 0
	}
	if k >= z.n {
		return 1
	}
	sum := 0.0
	for i := 1; i <= k; i++ {
		sum += math.Pow(float64(i), -z.s)
	}
	return sum / z.h
}

// Mean returns the mean: H_{n,s-1} / H_{n,s}.
func (z *Zipf) Mean() float64 {
	num := 0.0
	for k := 1; k <= z.n; k++ {
		num += float64(k) * math.Pow(float64(k), -z.s)
	}
	return num / z.h
}

// Var returns the variance.
func (z *Zipf) Var() float64 {
	m := z.Mean()
	e2 := 0.0
	for k := 1; k <= z.n; k++ {
		fk := float64(k)
		e2 += fk * fk * math.Pow(fk, -z.s)
	}
	e2 /= z.h
	return e2 - m*m
}

// ---------------------------------------------------------------------------
// Discrete Uniform Distribution
// ---------------------------------------------------------------------------

// DiscreteUniform represents a discrete uniform distribution on [low, high].
type DiscreteUniform struct {
	low  int
	high int
}

// NewDiscreteUniform creates a DiscreteUniform distribution on {low, low+1, ..., high}.
// Panics if low > high.
func NewDiscreteUniform(low, high int) *DiscreteUniform {
	if low > high {
		panic("scigo: DiscreteUniform requires low <= high")
	}
	return &DiscreteUniform{low: low, high: high}
}

// PMF returns the probability mass function at k.
func (d *DiscreteUniform) PMF(k int) float64 {
	if k < d.low || k > d.high {
		return 0
	}
	return 1.0 / float64(d.high-d.low+1)
}

// CDF returns the cumulative distribution function at k.
func (d *DiscreteUniform) CDF(k int) float64 {
	if k < d.low {
		return 0
	}
	if k >= d.high {
		return 1
	}
	return float64(k-d.low+1) / float64(d.high-d.low+1)
}

// Mean returns the mean: (low+high)/2.
func (d *DiscreteUniform) Mean() float64 {
	return float64(d.low+d.high) / 2.0
}

// Var returns the variance: ((high-low+1)^2 - 1) / 12.
func (d *DiscreteUniform) Var() float64 {
	n := float64(d.high - d.low + 1)
	return (n*n - 1) / 12.0
}

// ---------------------------------------------------------------------------
// Boltzmann Distribution (Truncated Discrete Exponential)
// ---------------------------------------------------------------------------

// Boltzmann represents a truncated discrete exponential (Boltzmann) distribution
// on {0, 1, ..., n-1} with rate parameter lambda.
// PMF(k) = exp(-lambda*k) / Z, where Z = sum_{k=0}^{n-1} exp(-lambda*k).
type Boltzmann struct {
	lambda float64
	n      int
	z      float64 // normalization constant
}

// NewBoltzmann creates a Boltzmann distribution. Panics if lambda <= 0 or n <= 0.
func NewBoltzmann(lambda float64, n int) *Boltzmann {
	if lambda <= 0 {
		panic("scigo: Boltzmann lambda must be positive")
	}
	if n <= 0 {
		panic("scigo: Boltzmann n must be positive")
	}
	// Z = sum_{k=0}^{n-1} exp(-lambda*k) = (1 - exp(-lambda*n)) / (1 - exp(-lambda))
	var z float64
	el := math.Exp(-lambda)
	if math.Abs(el-1) < 1e-15 {
		z = float64(n)
	} else {
		z = (1 - math.Pow(el, float64(n))) / (1 - el)
	}
	return &Boltzmann{lambda: lambda, n: n, z: z}
}

// PMF returns the probability mass function at k.
func (b *Boltzmann) PMF(k int) float64 {
	if k < 0 || k >= b.n {
		return 0
	}
	return math.Exp(-b.lambda*float64(k)) / b.z
}

// CDF returns the cumulative distribution function at k.
func (b *Boltzmann) CDF(k int) float64 {
	if k < 0 {
		return 0
	}
	if k >= b.n-1 {
		return 1
	}
	// Sum of geometric series: sum_{i=0}^{k} exp(-lambda*i) = (1 - exp(-lambda*(k+1))) / (1 - exp(-lambda))
	el := math.Exp(-b.lambda)
	if math.Abs(el-1) < 1e-15 {
		return float64(k+1) / float64(b.n)
	}
	num := (1 - math.Pow(el, float64(k+1))) / (1 - el)
	return num / b.z
}

// Mean returns the mean.
func (b *Boltzmann) Mean() float64 {
	sum := 0.0
	for k := 0; k < b.n; k++ {
		sum += float64(k) * math.Exp(-b.lambda*float64(k))
	}
	return sum / b.z
}

// Var returns the variance.
func (b *Boltzmann) Var() float64 {
	m := b.Mean()
	sum := 0.0
	for k := 0; k < b.n; k++ {
		sum += float64(k) * float64(k) * math.Exp(-b.lambda*float64(k))
	}
	return sum/b.z - m*m
}
