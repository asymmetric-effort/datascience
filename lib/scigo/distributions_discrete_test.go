//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// Negative Binomial Tests
// ---------------------------------------------------------------------------

func TestNegativeBinomialPanic(t *testing.T) {
	tests := []struct {
		name string
		r, p float64
	}{
		{"r=0", 0, 0.5},
		{"r<0", -1, 0.5},
		{"p=0", 1, 0},
		{"p>1", 1, 1.1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic")
				}
			}()
			NewNegativeBinomial(tc.r, tc.p)
		})
	}
}

func TestNegativeBinomialMeanVar(t *testing.T) {
	nb := NewNegativeBinomial(5, 0.4)
	// Mean = r*(1-p)/p = 5*0.6/0.4 = 7.5
	if !approxEqual(nb.Mean(), 7.5, 1e-12) {
		t.Errorf("NegBin(5,0.4).Mean() = %v, want 7.5", nb.Mean())
	}
	// Var = r*(1-p)/p^2 = 5*0.6/0.16 = 18.75
	if !approxEqual(nb.Var(), 18.75, 1e-12) {
		t.Errorf("NegBin(5,0.4).Var() = %v, want 18.75", nb.Var())
	}
}

func TestNegativeBinomialPMF(t *testing.T) {
	nb := NewNegativeBinomial(3, 0.5)
	// PMF(0) = C(2,0) * 0.5^3 * 0.5^0 = 0.125
	got := nb.PMF(0)
	if !approxEqual(got, 0.125, 1e-12) {
		t.Errorf("NegBin(3,0.5).PMF(0) = %v, want 0.125", got)
	}

	// PMF(1) = C(3,1) * 0.5^3 * 0.5^1 = 3*0.0625 = 0.1875
	got = nb.PMF(1)
	if !approxEqual(got, 0.1875, 1e-12) {
		t.Errorf("NegBin(3,0.5).PMF(1) = %v, want 0.1875", got)
	}

	// PMF for negative k should be 0
	if nb.PMF(-1) != 0 {
		t.Error("NegBin PMF for k<0 should be 0")
	}

	// Sum should be close to 1
	sum := 0.0
	for k := 0; k < 100; k++ {
		sum += nb.PMF(k)
	}
	if !approxEqual(sum, 1.0, 1e-8) {
		t.Errorf("Sum of NegBin(3,0.5) PMFs = %v, want ~1.0", sum)
	}
}

func TestNegativeBinomialCDF(t *testing.T) {
	nb := NewNegativeBinomial(3, 0.5)

	// CDF should match sum of PMFs
	for _, k := range []int{0, 1, 3, 5, 10} {
		sum := 0.0
		for i := 0; i <= k; i++ {
			sum += nb.PMF(i)
		}
		got := nb.CDF(k)
		if !approxEqual(got, sum, 1e-6) {
			t.Errorf("NegBin(3,0.5).CDF(%d) = %v, sum of PMFs = %v", k, got, sum)
		}
	}

	if nb.CDF(-1) != 0 {
		t.Error("NegBin CDF for k<0 should be 0")
	}
}

// ---------------------------------------------------------------------------
// Geometric Tests
// ---------------------------------------------------------------------------

func TestGeometricPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for p=0")
		}
	}()
	NewGeometric(0)
}

func TestGeometricMeanVar(t *testing.T) {
	g := NewGeometric(0.3)
	// Mean = (1-p)/p = 0.7/0.3
	wantMean := 0.7 / 0.3
	if !approxEqual(g.Mean(), wantMean, 1e-12) {
		t.Errorf("Geometric(0.3).Mean() = %v, want %v", g.Mean(), wantMean)
	}
	// Var = (1-p)/p^2 = 0.7/0.09
	wantVar := 0.7 / 0.09
	if !approxEqual(g.Var(), wantVar, 1e-10) {
		t.Errorf("Geometric(0.3).Var() = %v, want %v", g.Var(), wantVar)
	}
}

func TestGeometricPMF(t *testing.T) {
	g := NewGeometric(0.5)
	// PMF(0) = 0.5
	if !approxEqual(g.PMF(0), 0.5, 1e-12) {
		t.Errorf("Geometric(0.5).PMF(0) = %v, want 0.5", g.PMF(0))
	}
	// PMF(1) = 0.25
	if !approxEqual(g.PMF(1), 0.25, 1e-12) {
		t.Errorf("Geometric(0.5).PMF(1) = %v, want 0.25", g.PMF(1))
	}
	// PMF(k) = p*(1-p)^k
	for k := 0; k < 10; k++ {
		got := g.PMF(k)
		want := 0.5 * math.Pow(0.5, float64(k))
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Geometric(0.5).PMF(%d) = %v, want %v", k, got, want)
		}
	}

	if g.PMF(-1) != 0 {
		t.Error("Geometric PMF for k<0 should be 0")
	}

	// Sum should be 1
	sum := 0.0
	for k := 0; k < 50; k++ {
		sum += g.PMF(k)
	}
	if !approxEqual(sum, 1.0, 1e-10) {
		t.Errorf("Sum of Geometric(0.5) PMFs = %v, want ~1.0", sum)
	}
}

func TestGeometricCDF(t *testing.T) {
	g := NewGeometric(0.5)
	// CDF(k) = 1 - (1-p)^(k+1)
	for _, k := range []int{0, 1, 2, 5, 10} {
		got := g.CDF(k)
		want := 1 - math.Pow(0.5, float64(k)+1)
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("Geometric(0.5).CDF(%d) = %v, want %v", k, got, want)
		}
	}

	if g.CDF(-1) != 0 {
		t.Error("Geometric CDF for k<0 should be 0")
	}
}

func TestGeometricIsNegBinR1(t *testing.T) {
	// Geometric(p) is NegativeBinomial(1, p)
	p := 0.3
	g := NewGeometric(p)
	nb := NewNegativeBinomial(1, p)
	for k := 0; k < 10; k++ {
		if !approxEqual(g.PMF(k), nb.PMF(k), 1e-12) {
			t.Errorf("Geometric(%v).PMF(%d) = %v, NegBin(1,%v).PMF(%d) = %v",
				p, k, g.PMF(k), p, k, nb.PMF(k))
		}
	}
}

// ---------------------------------------------------------------------------
// Hypergeometric Tests
// ---------------------------------------------------------------------------

func TestHypergeometricPanic(t *testing.T) {
	tests := []struct {
		name    string
		N, K, n int
	}{
		{"N<0", -1, 0, 0},
		{"K>N", 10, 11, 5},
		{"K<0", 10, -1, 5},
		{"n>N", 10, 5, 11},
		{"n<0", 10, 5, -1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic")
				}
			}()
			NewHypergeometric(tc.N, tc.K, tc.n)
		})
	}
}

func TestHypergeometricMeanVar(t *testing.T) {
	h := NewHypergeometric(50, 20, 10)
	// Mean = n*K/N = 10*20/50 = 4
	if !approxEqual(h.Mean(), 4, 1e-12) {
		t.Errorf("Hypergeo(50,20,10).Mean() = %v, want 4", h.Mean())
	}
	// Var = n*K*(N-K)*(N-n) / (N^2*(N-1))
	// = 10*20*30*40 / (2500*49) = 240000/122500
	wantVar := 240000.0 / 122500.0
	if !approxEqual(h.Var(), wantVar, 1e-10) {
		t.Errorf("Hypergeo(50,20,10).Var() = %v, want %v", h.Var(), wantVar)
	}
}

func TestHypergeometricPMF(t *testing.T) {
	// Classic example: deck of 52 cards, 13 hearts, draw 5
	h := NewHypergeometric(52, 13, 5)

	// PMF(0) = C(13,0)*C(39,5)/C(52,5)
	got := h.PMF(0)
	// C(39,5)/C(52,5) = 575757/2598960
	want := 575757.0 / 2598960.0
	if !approxEqual(got, want, 1e-8) {
		t.Errorf("Hypergeo(52,13,5).PMF(0) = %v, want %v", got, want)
	}

	// Sum of PMFs should be 1
	sum := 0.0
	for k := 0; k <= 5; k++ {
		sum += h.PMF(k)
	}
	if !approxEqual(sum, 1.0, 1e-8) {
		t.Errorf("Sum of Hypergeo(52,13,5) PMFs = %v, want 1.0", sum)
	}

	// Out of range
	if h.PMF(-1) != 0 {
		t.Error("Hypergeometric PMF for k<0 should be 0")
	}
	if h.PMF(6) != 0 {
		t.Error("Hypergeometric PMF for k>n should be 0")
	}
}

func TestHypergeometricCDF(t *testing.T) {
	h := NewHypergeometric(20, 7, 5)
	// CDF should match sum of PMFs
	for _, k := range []int{0, 1, 2, 3, 4, 5} {
		sum := 0.0
		lo := max(0, 5-(20-7))
		for i := lo; i <= k; i++ {
			sum += h.PMF(i)
		}
		got := h.CDF(k)
		if !approxEqual(got, sum, 1e-8) {
			t.Errorf("Hypergeo(20,7,5).CDF(%d) = %v, sum = %v", k, got, sum)
		}
	}
}

// ---------------------------------------------------------------------------
// Bernoulli Tests
// ---------------------------------------------------------------------------

func TestBernoulliPanic(t *testing.T) {
	tests := []struct {
		name string
		p    float64
	}{
		{"p<0", -0.1},
		{"p>1", 1.1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic")
				}
			}()
			NewBernoulli(tc.p)
		})
	}
}

func TestBernoulliMeanVar(t *testing.T) {
	b := NewBernoulli(0.3)
	if !approxEqual(b.Mean(), 0.3, 1e-12) {
		t.Errorf("Bernoulli(0.3).Mean() = %v, want 0.3", b.Mean())
	}
	// Var = 0.3*0.7 = 0.21
	if !approxEqual(b.Var(), 0.21, 1e-12) {
		t.Errorf("Bernoulli(0.3).Var() = %v, want 0.21", b.Var())
	}
}

func TestBernoulliPMF(t *testing.T) {
	b := NewBernoulli(0.7)
	if !approxEqual(b.PMF(0), 0.3, 1e-12) {
		t.Errorf("Bernoulli(0.7).PMF(0) = %v, want 0.3", b.PMF(0))
	}
	if !approxEqual(b.PMF(1), 0.7, 1e-12) {
		t.Errorf("Bernoulli(0.7).PMF(1) = %v, want 0.7", b.PMF(1))
	}
	if b.PMF(-1) != 0 {
		t.Error("Bernoulli PMF for k!=0,1 should be 0")
	}
	if b.PMF(2) != 0 {
		t.Error("Bernoulli PMF for k!=0,1 should be 0")
	}
}

func TestBernoulliCDF(t *testing.T) {
	b := NewBernoulli(0.7)
	if b.CDF(-1) != 0 {
		t.Error("Bernoulli CDF(-1) should be 0")
	}
	if !approxEqual(b.CDF(0), 0.3, 1e-12) {
		t.Errorf("Bernoulli(0.7).CDF(0) = %v, want 0.3", b.CDF(0))
	}
	if b.CDF(1) != 1 {
		t.Errorf("Bernoulli(0.7).CDF(1) = %v, want 1", b.CDF(1))
	}
}

func TestBernoulliIsBinomialN1(t *testing.T) {
	// Bernoulli(p) should match Binomial(1, p)
	p := 0.4
	be := NewBernoulli(p)
	bi := NewBinomial(1, p)
	for k := 0; k <= 1; k++ {
		if !approxEqual(be.PMF(k), bi.PMF(k), 1e-12) {
			t.Errorf("Bernoulli(%v).PMF(%d) = %v, Binomial(1,%v).PMF(%d) = %v",
				p, k, be.PMF(k), p, k, bi.PMF(k))
		}
	}
	if !approxEqual(be.Mean(), bi.Mean(), 1e-12) {
		t.Error("Bernoulli and Binomial(1,p) mean should match")
	}
	if !approxEqual(be.Var(), bi.Var(), 1e-12) {
		t.Error("Bernoulli and Binomial(1,p) variance should match")
	}
}

// ---------------------------------------------------------------------------
// Zipf Tests
// ---------------------------------------------------------------------------

func TestZipfPanic(t *testing.T) {
	tests := []struct {
		name string
		s    float64
		n    int
	}{
		{"s=0", 0, 10},
		{"s<0", -1, 10},
		{"n=0", 1, 0},
		{"n<0", 1, -1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic")
				}
			}()
			NewZipf(tc.s, tc.n)
		})
	}
}

func TestZipfPMF(t *testing.T) {
	z := NewZipf(1, 5)
	// H_5,1 = 1 + 1/2 + 1/3 + 1/4 + 1/5 = 137/60
	h := 1.0 + 1.0/2 + 1.0/3 + 1.0/4 + 1.0/5

	// PMF(1) = 1/H
	got := z.PMF(1)
	want := 1.0 / h
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("Zipf(1,5).PMF(1) = %v, want %v", got, want)
	}

	// PMF(2) = (1/2)/H
	got = z.PMF(2)
	want = 0.5 / h
	if !approxEqual(got, want, 1e-12) {
		t.Errorf("Zipf(1,5).PMF(2) = %v, want %v", got, want)
	}

	// Sum of PMFs = 1
	sum := 0.0
	for k := 1; k <= 5; k++ {
		sum += z.PMF(k)
	}
	if !approxEqual(sum, 1.0, 1e-10) {
		t.Errorf("Sum of Zipf(1,5) PMFs = %v, want 1.0", sum)
	}

	// Out of range
	if z.PMF(0) != 0 {
		t.Error("Zipf PMF(0) should be 0")
	}
	if z.PMF(6) != 0 {
		t.Error("Zipf PMF(6) should be 0")
	}
}

func TestZipfCDF(t *testing.T) {
	z := NewZipf(1, 5)
	// CDF(5) = 1
	if !approxEqual(z.CDF(5), 1.0, 1e-10) {
		t.Errorf("Zipf(1,5).CDF(5) = %v, want 1", z.CDF(5))
	}

	// CDF monotonic
	prev := 0.0
	for k := 1; k <= 5; k++ {
		c := z.CDF(k)
		if c < prev-1e-10 {
			t.Errorf("Zipf CDF not monotonic at k=%d", k)
		}
		prev = c
	}
}

func TestZipfMeanVar(t *testing.T) {
	z := NewZipf(1, 3)
	// H_3,1 = 11/6
	h := 11.0 / 6
	// Mean = (1*1 + 2*(1/2) + 3*(1/3)) / H = (1+1+1)/H = 3/H
	wantMean := 3.0 / h
	if !approxEqual(z.Mean(), wantMean, 1e-10) {
		t.Errorf("Zipf(1,3).Mean() = %v, want %v", z.Mean(), wantMean)
	}

	// Verify variance is positive
	if z.Var() < 0 {
		t.Error("Zipf variance should be non-negative")
	}
}

// ---------------------------------------------------------------------------
// DiscreteUniform Tests
// ---------------------------------------------------------------------------

func TestDiscreteUniformPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for low > high")
		}
	}()
	NewDiscreteUniform(5, 3)
}

func TestDiscreteUniformMeanVar(t *testing.T) {
	d := NewDiscreteUniform(1, 6)
	// Mean = (1+6)/2 = 3.5
	if !approxEqual(d.Mean(), 3.5, 1e-12) {
		t.Errorf("DiscreteUniform(1,6).Mean() = %v, want 3.5", d.Mean())
	}
	// Var = (6^2-1)/12 = 35/12
	if !approxEqual(d.Var(), 35.0/12.0, 1e-12) {
		t.Errorf("DiscreteUniform(1,6).Var() = %v, want %v", d.Var(), 35.0/12.0)
	}
}

func TestDiscreteUniformPMF(t *testing.T) {
	d := NewDiscreteUniform(1, 6)
	// Each value has probability 1/6
	for k := 1; k <= 6; k++ {
		got := d.PMF(k)
		if !approxEqual(got, 1.0/6, 1e-12) {
			t.Errorf("DiscreteUniform(1,6).PMF(%d) = %v, want %v", k, got, 1.0/6)
		}
	}
	// Out of range
	if d.PMF(0) != 0 {
		t.Error("DiscreteUniform(1,6).PMF(0) should be 0")
	}
	if d.PMF(7) != 0 {
		t.Error("DiscreteUniform(1,6).PMF(7) should be 0")
	}
}

func TestDiscreteUniformCDF(t *testing.T) {
	d := NewDiscreteUniform(1, 6)
	// CDF(k) = (k-low+1)/(high-low+1)
	for k := 1; k <= 6; k++ {
		got := d.CDF(k)
		want := float64(k) / 6.0
		if !approxEqual(got, want, 1e-12) {
			t.Errorf("DiscreteUniform(1,6).CDF(%d) = %v, want %v", k, got, want)
		}
	}

	if d.CDF(0) != 0 {
		t.Error("DiscreteUniform CDF below range should be 0")
	}
	if d.CDF(7) != 1 {
		t.Error("DiscreteUniform CDF above range should be 1")
	}
}

func TestDiscreteUniformSingleValue(t *testing.T) {
	d := NewDiscreteUniform(5, 5)
	if d.PMF(5) != 1 {
		t.Error("DiscreteUniform(5,5).PMF(5) should be 1")
	}
	if d.Mean() != 5 {
		t.Error("DiscreteUniform(5,5).Mean() should be 5")
	}
	if d.Var() != 0 {
		t.Error("DiscreteUniform(5,5).Var() should be 0")
	}
}

// ---------------------------------------------------------------------------
// Boltzmann Tests
// ---------------------------------------------------------------------------

func TestBoltzmannPanic(t *testing.T) {
	tests := []struct {
		name   string
		lambda float64
		n      int
	}{
		{"lambda=0", 0, 10},
		{"lambda<0", -1, 10},
		{"n=0", 1, 0},
		{"n<0", 1, -1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Error("expected panic")
				}
			}()
			NewBoltzmann(tc.lambda, tc.n)
		})
	}
}

func TestBoltzmannPMF(t *testing.T) {
	b := NewBoltzmann(1, 5)
	// Sum of PMFs should be 1
	sum := 0.0
	for k := 0; k < 5; k++ {
		sum += b.PMF(k)
	}
	if !approxEqual(sum, 1.0, 1e-10) {
		t.Errorf("Sum of Boltzmann(1,5) PMFs = %v, want 1.0", sum)
	}

	// PMF(0) should be the largest
	if b.PMF(0) < b.PMF(1) {
		t.Error("Boltzmann PMF(0) should be >= PMF(1)")
	}

	// Out of range
	if b.PMF(-1) != 0 {
		t.Error("Boltzmann PMF(-1) should be 0")
	}
	if b.PMF(5) != 0 {
		t.Error("Boltzmann PMF(5) should be 0 (support is 0..n-1)")
	}
}

func TestBoltzmannCDF(t *testing.T) {
	b := NewBoltzmann(1, 5)

	// CDF(n-1) = 1
	if !approxEqual(b.CDF(4), 1.0, 1e-10) {
		t.Errorf("Boltzmann(1,5).CDF(4) = %v, want 1", b.CDF(4))
	}

	// CDF(-1) = 0
	if b.CDF(-1) != 0 {
		t.Error("Boltzmann CDF(-1) should be 0")
	}

	// CDF should match sum of PMFs
	for k := 0; k < 5; k++ {
		sum := 0.0
		for i := 0; i <= k; i++ {
			sum += b.PMF(i)
		}
		got := b.CDF(k)
		if !approxEqual(got, sum, 1e-10) {
			t.Errorf("Boltzmann(1,5).CDF(%d) = %v, sum = %v", k, got, sum)
		}
	}

	// Monotonic
	prev := 0.0
	for k := 0; k < 5; k++ {
		c := b.CDF(k)
		if c < prev-1e-10 {
			t.Errorf("Boltzmann CDF not monotonic at k=%d", k)
		}
		prev = c
	}
}

func TestBoltzmannMeanVar(t *testing.T) {
	b := NewBoltzmann(1, 5)

	// Mean should be between 0 and n-1
	m := b.Mean()
	if m < 0 || m > 4 {
		t.Errorf("Boltzmann(1,5).Mean() = %v, expected in [0, 4]", m)
	}

	// Variance should be positive
	v := b.Var()
	if v < 0 {
		t.Errorf("Boltzmann(1,5).Var() = %v, expected >= 0", v)
	}

	// Verify mean and variance by direct computation
	sumK := 0.0
	sumK2 := 0.0
	for k := 0; k < 5; k++ {
		p := b.PMF(k)
		sumK += float64(k) * p
		sumK2 += float64(k) * float64(k) * p
	}
	if !approxEqual(m, sumK, 1e-10) {
		t.Errorf("Boltzmann Mean %v != sum k*PMF %v", m, sumK)
	}
	wantVar := sumK2 - sumK*sumK
	if !approxEqual(v, wantVar, 1e-10) {
		t.Errorf("Boltzmann Var %v != computed %v", v, wantVar)
	}
}

func TestBoltzmannKnownValues(t *testing.T) {
	// For lambda=1, n=3: Z = 1 + e^{-1} + e^{-2}
	b := NewBoltzmann(1, 3)
	z := 1 + math.Exp(-1) + math.Exp(-2)
	// PMF(0) = 1/Z
	if !approxEqual(b.PMF(0), 1/z, 1e-12) {
		t.Errorf("Boltzmann(1,3).PMF(0) = %v, want %v", b.PMF(0), 1/z)
	}
	// PMF(1) = e^{-1}/Z
	if !approxEqual(b.PMF(1), math.Exp(-1)/z, 1e-12) {
		t.Errorf("Boltzmann(1,3).PMF(1) = %v, want %v", b.PMF(1), math.Exp(-1)/z)
	}
	// PMF(2) = e^{-2}/Z
	if !approxEqual(b.PMF(2), math.Exp(-2)/z, 1e-12) {
		t.Errorf("Boltzmann(1,3).PMF(2) = %v, want %v", b.PMF(2), math.Exp(-2)/z)
	}
}
