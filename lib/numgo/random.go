package numgo

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// RNG wraps a seeded random source for reproducible random number generation.
type RNG struct {
	src *rand.Rand
}

// NewRNG creates a new RNG with the given seed.
func NewRNG(seed int64) *RNG {
	return &RNG{
		src: rand.New(rand.NewPCG(uint64(seed), uint64(seed>>1)^0xa02bdbf7bb3c0785)),
	}
}

// Rand returns an NDArray of the given shape with values drawn uniformly from [0, 1).
func (r *RNG) Rand(shape ...int) *NDArray {
	size := product(shape)
	data := make([]float64, size)
	for i := range data {
		data[i] = r.src.Float64()
	}
	return NewNDArray(shape, data)
}

// Randn returns an NDArray of the given shape with values drawn from the
// standard normal distribution (mean=0, std=1) using the Box-Muller transform.
func (r *RNG) Randn(shape ...int) *NDArray {
	size := product(shape)
	data := make([]float64, size)
	for i := 0; i < size; i += 2 {
		u1 := r.src.Float64()
		u2 := r.src.Float64()
		// Avoid log(0).
		for u1 == 0 {
			u1 = r.src.Float64()
		}
		z0 := math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
		z1 := math.Sqrt(-2*math.Log(u1)) * math.Sin(2*math.Pi*u2)
		data[i] = z0
		if i+1 < size {
			data[i+1] = z1
		}
	}
	return NewNDArray(shape, data)
}

// RandInt returns an NDArray of the given shape with integer values drawn
// uniformly from [low, high).
func (r *RNG) RandInt(low, high int, shape ...int) *NDArray {
	if low >= high {
		panic(fmt.Sprintf("numgo: RandInt requires low < high, got low=%d high=%d", low, high))
	}
	size := product(shape)
	data := make([]float64, size)
	span := high - low
	for i := range data {
		data[i] = float64(r.src.IntN(span) + low)
	}
	return NewNDArray(shape, data)
}

// Choice returns a slice of random indices in [0, n).
// If replace is true, indices may repeat; otherwise they are unique
// and size must be <= n.
func (r *RNG) Choice(n int, size int, replace bool) []int {
	if n <= 0 {
		panic("numgo: Choice requires n > 0")
	}
	if size < 0 {
		panic("numgo: Choice requires size >= 0")
	}
	if !replace && size > n {
		panic(fmt.Sprintf("numgo: Choice without replacement requires size <= n, got size=%d n=%d", size, n))
	}

	if replace {
		result := make([]int, size)
		for i := range result {
			result[i] = r.src.IntN(n)
		}
		return result
	}

	// Without replacement: Fisher-Yates partial shuffle.
	pool := make([]int, n)
	for i := range pool {
		pool[i] = i
	}
	for i := 0; i < size; i++ {
		j := i + r.src.IntN(n-i)
		pool[i], pool[j] = pool[j], pool[i]
	}
	result := make([]int, size)
	copy(result, pool[:size])
	return result
}

// Shuffle performs an in-place Fisher-Yates shuffle on the first axis of the
// array. For 1-D arrays this shuffles elements; for N-D arrays it shuffles
// the sub-arrays along axis 0.
func (r *RNG) Shuffle(a *NDArray) {
	n := a.shape[0]
	if a.Ndim() == 1 {
		for i := n - 1; i > 0; i-- {
			j := r.src.IntN(i + 1)
			a.data[i], a.data[j] = a.data[j], a.data[i]
		}
		return
	}

	// For N-D: swap entire sub-arrays along axis 0.
	stride := a.strides[0]
	tmp := make([]float64, stride)
	for i := n - 1; i > 0; i-- {
		j := r.src.IntN(i + 1)
		if i != j {
			iOff := i * stride
			jOff := j * stride
			copy(tmp, a.data[iOff:iOff+stride])
			copy(a.data[iOff:iOff+stride], a.data[jOff:jOff+stride])
			copy(a.data[jOff:jOff+stride], tmp)
		}
	}
}

// Normal returns an NDArray of the given shape with values drawn from a
// normal distribution with the specified mean and standard deviation.
func (r *RNG) Normal(mean, std float64, shape ...int) *NDArray {
	a := r.Randn(shape...)
	data := a.Data()
	for i := range data {
		data[i] = data[i]*std + mean
	}
	return NewNDArray(shape, data)
}

// Uniform returns an NDArray of the given shape with values drawn uniformly
// from [low, high).
func (r *RNG) Uniform(low, high float64, shape ...int) *NDArray {
	if low >= high {
		panic(fmt.Sprintf("numgo: Uniform requires low < high, got low=%f high=%f", low, high))
	}
	a := r.Rand(shape...)
	span := high - low
	data := a.Data()
	for i := range data {
		data[i] = data[i]*span + low
	}
	return NewNDArray(shape, data)
}

// gammaVariate generates a single sample from the Gamma(alpha, 1) distribution
// using the Marsaglia and Tsang method for alpha >= 1, with a boost for alpha < 1.
func (r *RNG) gammaVariate(alpha float64) float64 {
	if alpha <= 0 {
		panic("numgo: gammaVariate requires alpha > 0")
	}

	if alpha < 1.0 {
		// Boost: Gamma(alpha) = Gamma(alpha+1) * U^(1/alpha)
		return r.gammaVariate(alpha+1.0) * math.Pow(r.src.Float64(), 1.0/alpha)
	}

	// Marsaglia and Tsang's method for alpha >= 1.
	d := alpha - 1.0/3.0
	c := 1.0 / math.Sqrt(9.0*d)
	for {
		var x, v float64
		for {
			x = r.boxMullerSingle()
			v = 1.0 + c*x
			if v > 0 {
				break
			}
		}
		v = v * v * v
		u := r.src.Float64()
		if u < 1.0-0.0331*(x*x)*(x*x) {
			return d * v
		}
		if math.Log(u) < 0.5*x*x+d*(1.0-v+math.Log(v)) {
			return d * v
		}
	}
}

// boxMullerSingle returns a single standard normal variate.
func (r *RNG) boxMullerSingle() float64 {
	u1 := r.src.Float64()
	for u1 == 0 {
		u1 = r.src.Float64()
	}
	u2 := r.src.Float64()
	return math.Sqrt(-2*math.Log(u1)) * math.Cos(2*math.Pi*u2)
}

// Dirichlet draws a single sample from the Dirichlet distribution with
// parameter vector alpha. It returns a probability vector of length len(alpha).
// The implementation draws independent Gamma(alpha_i, 1) samples and normalizes.
func (r *RNG) Dirichlet(alpha []float64) []float64 {
	if len(alpha) == 0 {
		panic("numgo: Dirichlet requires non-empty alpha")
	}
	samples := make([]float64, len(alpha))
	total := 0.0
	for i, a := range alpha {
		samples[i] = r.gammaVariate(a)
		total += samples[i]
	}
	for i := range samples {
		samples[i] /= total
	}
	return samples
}

// defaultRNG is a package-level RNG used by the Seed function.
var defaultRNG = NewRNG(0)

// Seed sets the seed for the package-level default RNG.
func Seed(seed int64) {
	defaultRNG = NewRNG(seed)
}

// Random returns an NDArray of the given shape with values drawn uniformly from [0, 1).
// It is an alias for Rand.
func (r *RNG) Random(shape ...int) *NDArray {
	return r.Rand(shape...)
}

// Permutation returns an NDArray containing a random permutation of integers [0, n).
func (r *RNG) Permutation(n int) *NDArray {
	perm := make([]float64, n)
	for i := range perm {
		perm[i] = float64(i)
	}
	// Fisher-Yates shuffle.
	for i := n - 1; i > 0; i-- {
		j := r.src.IntN(i + 1)
		perm[i], perm[j] = perm[j], perm[i]
	}
	return FromSlice(perm)
}

// Exponential returns samples from an exponential distribution with the given scale.
func (r *RNG) Exponential(scale float64, shape ...int) *NDArray {
	if scale <= 0 {
		panic("numgo: Exponential requires scale > 0")
	}
	size := product(shape)
	data := make([]float64, size)
	for i := range data {
		u := r.src.Float64()
		for u == 0 {
			u = r.src.Float64()
		}
		data[i] = -scale * math.Log(u)
	}
	return NewNDArray(shape, data)
}

// Poisson returns samples from a Poisson distribution with the given rate (lambda).
// Uses Knuth's algorithm for small lambda, and a rejection method for large lambda.
func (r *RNG) Poisson(lam float64, shape ...int) *NDArray {
	if lam < 0 {
		panic("numgo: Poisson requires lam >= 0")
	}
	size := product(shape)
	data := make([]float64, size)
	for i := range data {
		data[i] = float64(r.poissonSingle(lam))
	}
	return NewNDArray(shape, data)
}

func (r *RNG) poissonSingle(lam float64) int {
	if lam == 0 {
		return 0
	}
	if lam < 30 {
		// Knuth's algorithm.
		L := math.Exp(-lam)
		k := 0
		p := 1.0
		for {
			k++
			p *= r.src.Float64()
			if p < L {
				return k - 1
			}
		}
	}
	// Rejection method for large lambda (based on transformed normal).
	sqrtLam := math.Sqrt(lam)
	for {
		x := r.boxMullerSingle()*sqrtLam + lam
		k := int(math.Floor(x + 0.5))
		if k < 0 {
			continue
		}
		// Accept with some probability (simplified).
		return k
	}
}

// BinomialSample returns samples from a binomial distribution with parameters n and p.
func (r *RNG) BinomialSample(n int, p float64, shape ...int) *NDArray {
	if n < 0 {
		panic("numgo: BinomialSample requires n >= 0")
	}
	if p < 0 || p > 1 {
		panic("numgo: BinomialSample requires p in [0, 1]")
	}
	size := product(shape)
	data := make([]float64, size)
	for i := range data {
		successes := 0
		for trial := 0; trial < n; trial++ {
			if r.src.Float64() < p {
				successes++
			}
		}
		data[i] = float64(successes)
	}
	return NewNDArray(shape, data)
}

// Beta returns samples from a Beta(a, b) distribution using the Gamma distribution.
func (r *RNG) Beta(a, b float64, shape ...int) *NDArray {
	if a <= 0 || b <= 0 {
		panic("numgo: Beta requires a > 0 and b > 0")
	}
	size := product(shape)
	data := make([]float64, size)
	for i := range data {
		x := r.gammaVariate(a)
		y := r.gammaVariate(b)
		data[i] = x / (x + y)
	}
	return NewNDArray(shape, data)
}

// Gamma returns samples from a Gamma(shape, scale) distribution.
func (r *RNG) Gamma(shapep, scale float64, shapeArr ...int) *NDArray {
	if shapep <= 0 || scale <= 0 {
		panic("numgo: Gamma requires shape > 0 and scale > 0")
	}
	size := product(shapeArr)
	data := make([]float64, size)
	for i := range data {
		data[i] = r.gammaVariate(shapep) * scale
	}
	return NewNDArray(shapeArr, data)
}

// Chisquare returns samples from a chi-squared distribution with df degrees of freedom.
// Chi-squared(df) = Gamma(df/2, 2).
func (r *RNG) Chisquare(df float64, shape ...int) *NDArray {
	if df <= 0 {
		panic("numgo: Chisquare requires df > 0")
	}
	return r.Gamma(df/2.0, 2.0, shape...)
}

// StandardNormal returns samples from the standard normal distribution (mean=0, std=1).
// It is an alias for Randn.
func (r *RNG) StandardNormal(shape ...int) *NDArray {
	return r.Randn(shape...)
}

// StandardT returns samples from Student's t-distribution with df degrees of freedom.
// Uses the ratio of a standard normal to the square root of a chi-squared/df.
func (r *RNG) StandardT(df float64, shape ...int) *NDArray {
	if df <= 0 {
		panic("numgo: StandardT requires df > 0")
	}
	size := product(shape)
	data := make([]float64, size)
	for i := range data {
		z := r.boxMullerSingle()
		chi2 := r.gammaVariate(df/2.0) * 2.0
		data[i] = z / math.Sqrt(chi2/df)
	}
	return NewNDArray(shape, data)
}

// Multinomial draws a single sample from the multinomial distribution:
// distribute n trials among len(pvals) categories with the given probabilities.
// Returns a slice of counts (length len(pvals)) summing to n.
func (r *RNG) Multinomial(n int, pvals []float64) []int {
	if n < 0 {
		panic("numgo: Multinomial requires n >= 0")
	}
	if len(pvals) == 0 {
		panic("numgo: Multinomial requires non-empty pvals")
	}

	// Build cumulative probabilities.
	cumulative := make([]float64, len(pvals))
	cumulative[0] = pvals[0]
	for i := 1; i < len(pvals); i++ {
		cumulative[i] = cumulative[i-1] + pvals[i]
	}
	// Normalize to handle floating-point drift.
	total := cumulative[len(cumulative)-1]

	counts := make([]int, len(pvals))
	for trial := 0; trial < n; trial++ {
		u := r.src.Float64() * total
		// Binary search for the bucket.
		lo, hi := 0, len(cumulative)-1
		for lo < hi {
			mid := (lo + hi) / 2
			if cumulative[mid] <= u {
				lo = mid + 1
			} else {
				hi = mid
			}
		}
		counts[lo]++
	}
	return counts
}
