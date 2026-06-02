package scigo

import "math"

// ---------------------------------------------------------------------------
// F-Distribution
// ---------------------------------------------------------------------------

// FDistribution represents an F-distribution with df1 and df2 degrees of freedom.
type FDistribution struct {
	df1 float64
	df2 float64
}

// NewFDistribution creates an F-distribution with the given degrees of freedom.
// Panics if df1 <= 0 or df2 <= 0.
func NewFDistribution(df1, df2 float64) *FDistribution {
	if df1 <= 0 {
		panic("scigo: FDistribution df1 must be positive")
	}
	if df2 <= 0 {
		panic("scigo: FDistribution df2 must be positive")
	}
	return &FDistribution{df1: df1, df2: df2}
}

// PDF returns the probability density of the F-distribution at x.
func (f *FDistribution) PDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	d1, d2 := f.df1, f.df2
	logNum := (d1/2)*math.Log(d1) + (d2/2)*math.Log(d2) + (d1/2-1)*math.Log(x)
	logDen := ((d1 + d2) / 2) * math.Log(d1*x+d2)
	logBeta := Gammaln(d1/2) + Gammaln(d2/2) - Gammaln((d1+d2)/2)
	return math.Exp(logNum - logDen - logBeta)
}

// CDF returns the cumulative distribution function of the F-distribution at x.
// Uses the regularized incomplete beta function.
func (f *FDistribution) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	d1, d2 := f.df1, f.df2
	z := d1 * x / (d1*x + d2)
	return RegularizedIncompleteBeta(z, d1/2, d2/2)
}

// PPF returns the percent point function (inverse CDF) using Newton's method.
func (f *FDistribution) PPF(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return math.Inf(1)
	}
	// Initial guess from chi-squared ratio approximation
	x := f.df1 / f.df2 // start near the mean for df2>2
	if f.df2 > 2 {
		x = f.df2 / (f.df2 - 2) // mean
	}
	for i := 0; i < 100; i++ {
		cdfVal := f.CDF(x)
		pdfVal := f.PDF(x)
		if pdfVal == 0 {
			break
		}
		dx := (cdfVal - p) / pdfVal
		x -= dx
		if x <= 0 {
			x = 1e-10
		}
		if math.Abs(dx) < 1e-12*(1+x) {
			break
		}
	}
	return x
}

// LogPDF returns the natural log of the PDF at x.
func (f *FDistribution) LogPDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	d1, d2 := f.df1, f.df2
	logNum := (d1/2)*math.Log(d1) + (d2/2)*math.Log(d2) + (d1/2-1)*math.Log(x)
	logDen := ((d1 + d2) / 2) * math.Log(d1*x+d2)
	logBeta := Gammaln(d1/2) + Gammaln(d2/2) - Gammaln((d1+d2)/2)
	return logNum - logDen - logBeta
}

// Mean returns the mean of the F-distribution. Defined for df2 > 2.
func (f *FDistribution) Mean() float64 {
	if f.df2 <= 2 {
		return math.NaN()
	}
	return f.df2 / (f.df2 - 2)
}

// Var returns the variance of the F-distribution. Defined for df2 > 4.
func (f *FDistribution) Var() float64 {
	if f.df2 <= 4 {
		return math.NaN()
	}
	d1, d2 := f.df1, f.df2
	return 2 * d2 * d2 * (d1 + d2 - 2) / (d1 * (d2 - 2) * (d2 - 2) * (d2 - 4))
}

// ---------------------------------------------------------------------------
// Log-Normal Distribution
// ---------------------------------------------------------------------------

// Lognormal represents a log-normal distribution with parameters mu and sigma.
type Lognormal struct {
	mu    float64
	sigma float64
}

// NewLognormal creates a Lognormal distribution. Panics if sigma <= 0.
func NewLognormal(mu, sigma float64) *Lognormal {
	if sigma <= 0 {
		panic("scigo: Lognormal sigma must be positive")
	}
	return &Lognormal{mu: mu, sigma: sigma}
}

// PDF returns the probability density of the log-normal distribution at x.
func (ln *Lognormal) PDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := (math.Log(x) - ln.mu) / ln.sigma
	return math.Exp(-0.5*z*z) / (x * ln.sigma * math.Sqrt(2*math.Pi))
}

// CDF returns the cumulative distribution function at x.
func (ln *Lognormal) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := (math.Log(x) - ln.mu) / (ln.sigma * math.Sqrt2)
	return 0.5 * (1 + math.Erf(z))
}

// PPF returns the percent point function for probability p.
func (ln *Lognormal) PPF(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return math.Inf(1)
	}
	return math.Exp(ln.mu + ln.sigma*math.Sqrt2*Erfinv(2*p-1))
}

// LogPDF returns the log of the PDF at x.
func (ln *Lognormal) LogPDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	z := (math.Log(x) - ln.mu) / ln.sigma
	return -0.5*z*z - math.Log(x) - math.Log(ln.sigma) - 0.5*math.Log(2*math.Pi)
}

// Mean returns the mean of the log-normal distribution.
func (ln *Lognormal) Mean() float64 {
	return math.Exp(ln.mu + 0.5*ln.sigma*ln.sigma)
}

// Var returns the variance of the log-normal distribution.
func (ln *Lognormal) Var() float64 {
	s2 := ln.sigma * ln.sigma
	return (math.Exp(s2) - 1) * math.Exp(2*ln.mu+s2)
}

// ---------------------------------------------------------------------------
// Weibull Distribution (minimum)
// ---------------------------------------------------------------------------

// Weibull represents a Weibull minimum distribution with shape k and scale lambda.
type Weibull struct {
	shape float64
	scale float64
}

// NewWeibull creates a Weibull distribution. Panics if shape <= 0 or scale <= 0.
func NewWeibull(shape, scale float64) *Weibull {
	if shape <= 0 {
		panic("scigo: Weibull shape must be positive")
	}
	if scale <= 0 {
		panic("scigo: Weibull scale must be positive")
	}
	return &Weibull{shape: shape, scale: scale}
}

// PDF returns the probability density at x.
func (w *Weibull) PDF(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x == 0 {
		if w.shape < 1 {
			return math.Inf(1)
		}
		if w.shape == 1 {
			return 1.0 / w.scale
		}
		return 0
	}
	k, lam := w.shape, w.scale
	z := x / lam
	return (k / lam) * math.Pow(z, k-1) * math.Exp(-math.Pow(z, k))
}

// CDF returns the cumulative distribution function at x.
func (w *Weibull) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return 1 - math.Exp(-math.Pow(x/w.scale, w.shape))
}

// PPF returns the percent point function for probability p.
func (w *Weibull) PPF(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return math.Inf(1)
	}
	return w.scale * math.Pow(-math.Log(1-p), 1.0/w.shape)
}

// LogPDF returns the log of the PDF at x.
func (w *Weibull) LogPDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	k, lam := w.shape, w.scale
	z := x / lam
	return math.Log(k) - math.Log(lam) + (k-1)*math.Log(z) - math.Pow(z, k)
}

// Mean returns the mean of the Weibull distribution.
func (w *Weibull) Mean() float64 {
	return w.scale * math.Gamma(1+1.0/w.shape)
}

// Var returns the variance of the Weibull distribution.
func (w *Weibull) Var() float64 {
	g1 := math.Gamma(1 + 1.0/w.shape)
	g2 := math.Gamma(1 + 2.0/w.shape)
	return w.scale * w.scale * (g2 - g1*g1)
}

// ---------------------------------------------------------------------------
// Pareto Distribution
// ---------------------------------------------------------------------------

// Pareto represents a Pareto distribution with shape alpha and minimum xm.
type Pareto struct {
	alpha float64
	xm    float64
}

// NewPareto creates a Pareto distribution. Panics if alpha <= 0 or xm <= 0.
func NewPareto(alpha, xm float64) *Pareto {
	if alpha <= 0 {
		panic("scigo: Pareto alpha must be positive")
	}
	if xm <= 0 {
		panic("scigo: Pareto xm must be positive")
	}
	return &Pareto{alpha: alpha, xm: xm}
}

// PDF returns the probability density at x.
func (p *Pareto) PDF(x float64) float64 {
	if x < p.xm {
		return 0
	}
	return p.alpha * math.Pow(p.xm, p.alpha) / math.Pow(x, p.alpha+1)
}

// CDF returns the cumulative distribution function at x.
func (p *Pareto) CDF(x float64) float64 {
	if x < p.xm {
		return 0
	}
	return 1 - math.Pow(p.xm/x, p.alpha)
}

// PPF returns the percent point function for probability pr.
func (p *Pareto) PPF(pr float64) float64 {
	if pr <= 0 {
		return p.xm
	}
	if pr >= 1 {
		return math.Inf(1)
	}
	return p.xm / math.Pow(1-pr, 1.0/p.alpha)
}

// LogPDF returns the log of the PDF at x.
func (p *Pareto) LogPDF(x float64) float64 {
	if x < p.xm {
		return math.Inf(-1)
	}
	return math.Log(p.alpha) + p.alpha*math.Log(p.xm) - (p.alpha+1)*math.Log(x)
}

// Mean returns the mean. Defined for alpha > 1.
func (p *Pareto) Mean() float64 {
	if p.alpha <= 1 {
		return math.Inf(1)
	}
	return p.alpha * p.xm / (p.alpha - 1)
}

// Var returns the variance. Defined for alpha > 2.
func (p *Pareto) Var() float64 {
	if p.alpha <= 2 {
		return math.Inf(1)
	}
	return p.xm * p.xm * p.alpha / ((p.alpha - 1) * (p.alpha - 1) * (p.alpha - 2))
}

// ---------------------------------------------------------------------------
// Cauchy Distribution
// ---------------------------------------------------------------------------

// Cauchy represents a Cauchy (Lorentz) distribution with location and scale parameters.
type Cauchy struct {
	loc   float64
	scale float64
}

// NewCauchy creates a Cauchy distribution. Panics if scale <= 0.
func NewCauchy(loc, scale float64) *Cauchy {
	if scale <= 0 {
		panic("scigo: Cauchy scale must be positive")
	}
	return &Cauchy{loc: loc, scale: scale}
}

// PDF returns the probability density at x.
func (c *Cauchy) PDF(x float64) float64 {
	z := (x - c.loc) / c.scale
	return 1.0 / (math.Pi * c.scale * (1 + z*z))
}

// CDF returns the cumulative distribution function at x.
func (c *Cauchy) CDF(x float64) float64 {
	return 0.5 + math.Atan((x-c.loc)/c.scale)/math.Pi
}

// PPF returns the percent point function for probability p.
func (c *Cauchy) PPF(p float64) float64 {
	if p <= 0 {
		return math.Inf(-1)
	}
	if p >= 1 {
		return math.Inf(1)
	}
	return c.loc + c.scale*math.Tan(math.Pi*(p-0.5))
}

// LogPDF returns the log of the PDF at x.
func (c *Cauchy) LogPDF(x float64) float64 {
	z := (x - c.loc) / c.scale
	return -math.Log(math.Pi) - math.Log(c.scale) - math.Log(1+z*z)
}

// Mean returns NaN (undefined for Cauchy).
func (c *Cauchy) Mean() float64 {
	return math.NaN()
}

// Var returns NaN (undefined for Cauchy).
func (c *Cauchy) Var() float64 {
	return math.NaN()
}

// ---------------------------------------------------------------------------
// Laplace Distribution
// ---------------------------------------------------------------------------

// Laplace represents a Laplace (double exponential) distribution.
type Laplace struct {
	loc   float64
	scale float64
}

// NewLaplace creates a Laplace distribution. Panics if scale <= 0.
func NewLaplace(loc, scale float64) *Laplace {
	if scale <= 0 {
		panic("scigo: Laplace scale must be positive")
	}
	return &Laplace{loc: loc, scale: scale}
}

// PDF returns the probability density at x.
func (l *Laplace) PDF(x float64) float64 {
	return math.Exp(-math.Abs(x-l.loc)/l.scale) / (2 * l.scale)
}

// CDF returns the cumulative distribution function at x.
func (l *Laplace) CDF(x float64) float64 {
	if x < l.loc {
		return 0.5 * math.Exp((x-l.loc)/l.scale)
	}
	return 1 - 0.5*math.Exp(-(x-l.loc)/l.scale)
}

// PPF returns the percent point function for probability p.
func (l *Laplace) PPF(p float64) float64 {
	if p <= 0 {
		return math.Inf(-1)
	}
	if p >= 1 {
		return math.Inf(1)
	}
	if p < 0.5 {
		return l.loc + l.scale*math.Log(2*p)
	}
	return l.loc - l.scale*math.Log(2*(1-p))
}

// LogPDF returns the log of the PDF at x.
func (l *Laplace) LogPDF(x float64) float64 {
	return -math.Abs(x-l.loc)/l.scale - math.Log(2*l.scale)
}

// Mean returns the mean.
func (l *Laplace) Mean() float64 {
	return l.loc
}

// Var returns the variance.
func (l *Laplace) Var() float64 {
	return 2 * l.scale * l.scale
}

// ---------------------------------------------------------------------------
// Logistic Distribution
// ---------------------------------------------------------------------------

// Logistic represents a logistic distribution.
type Logistic struct {
	loc   float64
	scale float64
}

// NewLogistic creates a Logistic distribution. Panics if scale <= 0.
func NewLogistic(loc, scale float64) *Logistic {
	if scale <= 0 {
		panic("scigo: Logistic scale must be positive")
	}
	return &Logistic{loc: loc, scale: scale}
}

// PDF returns the probability density at x.
func (l *Logistic) PDF(x float64) float64 {
	z := (x - l.loc) / l.scale
	ez := math.Exp(-z)
	return ez / (l.scale * (1 + ez) * (1 + ez))
}

// CDF returns the cumulative distribution function at x.
func (l *Logistic) CDF(x float64) float64 {
	z := (x - l.loc) / l.scale
	return 1.0 / (1.0 + math.Exp(-z))
}

// PPF returns the percent point function for probability p.
func (l *Logistic) PPF(p float64) float64 {
	if p <= 0 {
		return math.Inf(-1)
	}
	if p >= 1 {
		return math.Inf(1)
	}
	return l.loc + l.scale*math.Log(p/(1-p))
}

// LogPDF returns the log of the PDF at x.
func (l *Logistic) LogPDF(x float64) float64 {
	z := (x - l.loc) / l.scale
	return -z - math.Log(l.scale) - 2*math.Log(1+math.Exp(-z))
}

// Mean returns the mean.
func (l *Logistic) Mean() float64 {
	return l.loc
}

// Var returns the variance.
func (l *Logistic) Var() float64 {
	return l.scale * l.scale * math.Pi * math.Pi / 3
}

// ---------------------------------------------------------------------------
// Gumbel Distribution (Extreme Value Type I)
// ---------------------------------------------------------------------------

// Gumbel represents a Gumbel distribution (Type I extreme value, for maxima).
type Gumbel struct {
	loc   float64
	scale float64
}

// NewGumbel creates a Gumbel distribution. Panics if scale <= 0.
func NewGumbel(loc, scale float64) *Gumbel {
	if scale <= 0 {
		panic("scigo: Gumbel scale must be positive")
	}
	return &Gumbel{loc: loc, scale: scale}
}

// PDF returns the probability density at x.
func (g *Gumbel) PDF(x float64) float64 {
	z := (x - g.loc) / g.scale
	return math.Exp(-z-math.Exp(-z)) / g.scale
}

// CDF returns the cumulative distribution function at x.
func (g *Gumbel) CDF(x float64) float64 {
	z := (x - g.loc) / g.scale
	return math.Exp(-math.Exp(-z))
}

// PPF returns the percent point function for probability p.
func (g *Gumbel) PPF(p float64) float64 {
	if p <= 0 {
		return math.Inf(-1)
	}
	if p >= 1 {
		return math.Inf(1)
	}
	return g.loc - g.scale*math.Log(-math.Log(p))
}

// LogPDF returns the log of the PDF at x.
func (g *Gumbel) LogPDF(x float64) float64 {
	z := (x - g.loc) / g.scale
	return -z - math.Exp(-z) - math.Log(g.scale)
}

// Mean returns the mean.
func (g *Gumbel) Mean() float64 {
	// Euler-Mascheroni constant
	return g.loc + g.scale*0.5772156649015329
}

// Var returns the variance.
func (g *Gumbel) Var() float64 {
	return g.scale * g.scale * math.Pi * math.Pi / 6
}

// ---------------------------------------------------------------------------
// Rayleigh Distribution
// ---------------------------------------------------------------------------

// Rayleigh represents a Rayleigh distribution with parameter sigma.
type Rayleigh struct {
	sigma float64
}

// NewRayleigh creates a Rayleigh distribution. Panics if sigma <= 0.
func NewRayleigh(sigma float64) *Rayleigh {
	if sigma <= 0 {
		panic("scigo: Rayleigh sigma must be positive")
	}
	return &Rayleigh{sigma: sigma}
}

// PDF returns the probability density at x.
func (r *Rayleigh) PDF(x float64) float64 {
	if x < 0 {
		return 0
	}
	s2 := r.sigma * r.sigma
	return (x / s2) * math.Exp(-x*x/(2*s2))
}

// CDF returns the cumulative distribution function at x.
func (r *Rayleigh) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return 1 - math.Exp(-x*x/(2*r.sigma*r.sigma))
}

// PPF returns the percent point function for probability p.
func (r *Rayleigh) PPF(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return math.Inf(1)
	}
	return r.sigma * math.Sqrt(-2*math.Log(1-p))
}

// LogPDF returns the log of the PDF at x.
func (r *Rayleigh) LogPDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	s2 := r.sigma * r.sigma
	return math.Log(x) - math.Log(s2) - x*x/(2*s2)
}

// Mean returns the mean of the Rayleigh distribution.
func (r *Rayleigh) Mean() float64 {
	return r.sigma * math.Sqrt(math.Pi/2)
}

// Var returns the variance of the Rayleigh distribution.
func (r *Rayleigh) Var() float64 {
	return (2 - math.Pi/2) * r.sigma * r.sigma
}

// ---------------------------------------------------------------------------
// Rice Distribution
// ---------------------------------------------------------------------------

// Rice represents a Rice (Rician) distribution with parameters nu and sigma.
type Rice struct {
	nu    float64
	sigma float64
}

// NewRice creates a Rice distribution. Panics if nu < 0 or sigma <= 0.
func NewRice(nu, sigma float64) *Rice {
	if nu < 0 {
		panic("scigo: Rice nu must be non-negative")
	}
	if sigma <= 0 {
		panic("scigo: Rice sigma must be positive")
	}
	return &Rice{nu: nu, sigma: sigma}
}

// besselI0 computes the modified Bessel function of the first kind, order 0.
// Uses polynomial approximation (Abramowitz and Stegun).
func besselI0(x float64) float64 {
	ax := math.Abs(x)
	if ax < 3.75 {
		t := x / 3.75
		t2 := t * t
		return 1 + t2*(3.5156229+t2*(3.0899424+t2*(1.2067492+
			t2*(0.2659732+t2*(0.0360768+t2*0.0045813)))))
	}
	t := 3.75 / ax
	return (math.Exp(ax) / math.Sqrt(ax)) * (0.39894228 + t*(0.01328592+
		t*(0.00225319+t*(-0.00157565+t*(0.00916281+t*(-0.02057706+
			t*(0.02635537+t*(-0.01647633+t*0.00392377))))))))
}

// besselI1 computes the modified Bessel function of the first kind, order 1.
func besselI1(x float64) float64 {
	ax := math.Abs(x)
	var ans float64
	if ax < 3.75 {
		t := x / 3.75
		t2 := t * t
		ans = ax * (0.5 + t2*(0.87890594+t2*(0.51498869+t2*(0.15084934+
			t2*(0.02658733+t2*(0.00301532+t2*0.00032411))))))
	} else {
		t := 3.75 / ax
		ans = (math.Exp(ax) / math.Sqrt(ax)) * (0.39894228 + t*(-0.03988024+
			t*(-0.00362018+t*(0.00163801+t*(-0.01031555+t*(0.02282967+
				t*(-0.02895312+t*(0.01787654+t*(-0.00420059)))))))))
	}
	if x < 0 {
		return -ans
	}
	return ans
}

// PDF returns the probability density at x.
func (r *Rice) PDF(x float64) float64 {
	if x < 0 {
		return 0
	}
	if x == 0 {
		if r.nu == 0 {
			return 0 // Rayleigh limit at x=0 gives 0
		}
		return 0
	}
	s2 := r.sigma * r.sigma
	return (x / s2) * math.Exp(-(x*x+r.nu*r.nu)/(2*s2)) * besselI0(x*r.nu/s2)
}

// CDF returns the cumulative distribution function at x.
// Uses numerical integration via Simpson's rule.
func (r *Rice) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	// Numerical integration using adaptive Simpson's rule
	n := 1000
	h := x / float64(n)
	sum := r.PDF(0) + r.PDF(x)
	for i := 1; i < n; i++ {
		xi := float64(i) * h
		if i%2 == 0 {
			sum += 2 * r.PDF(xi)
		} else {
			sum += 4 * r.PDF(xi)
		}
	}
	return sum * h / 3
}

// PPF returns the percent point function using Newton's method.
func (r *Rice) PPF(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return math.Inf(1)
	}
	// Initial guess from mean
	x := r.Mean()
	if x <= 0 {
		x = r.sigma
	}
	for i := 0; i < 100; i++ {
		cdfVal := r.CDF(x)
		pdfVal := r.PDF(x)
		if pdfVal == 0 {
			break
		}
		dx := (cdfVal - p) / pdfVal
		x -= dx
		if x <= 0 {
			x = 1e-10
		}
		if math.Abs(dx) < 1e-10*(1+x) {
			break
		}
	}
	return x
}

// LogPDF returns the log of the PDF at x.
func (r *Rice) LogPDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	s2 := r.sigma * r.sigma
	return math.Log(x) - math.Log(s2) - (x*x+r.nu*r.nu)/(2*s2) + math.Log(besselI0(x*r.nu/s2))
}

// Mean returns the mean of the Rice distribution.
// Uses the approximation: mean = sigma * sqrt(pi/2) * L_{1/2}(-nu^2/(2*sigma^2))
// where L is the Laguerre function. Approximated numerically.
func (r *Rice) Mean() float64 {
	// Use numerical integration for mean
	// For Rice, mean = sigma*sqrt(pi/2)*L_{1/2}(-nu^2/(2*sigma^2))
	// where L_{1/2}(x) = e^{x/2} * [(1-x)*I_0(-x/2) - x*I_1(-x/2)]
	s2 := r.sigma * r.sigma
	arg := -r.nu * r.nu / (2 * s2)
	halfArg := -arg / 2 // = nu^2/(4*sigma^2)
	l := math.Exp(arg/2) * ((1-arg)*besselI0(halfArg) - arg*besselI1(halfArg))
	return r.sigma * math.Sqrt(math.Pi/2) * l
}

// Var returns the variance of the Rice distribution.
func (r *Rice) Var() float64 {
	m := r.Mean()
	return 2*r.sigma*r.sigma + r.nu*r.nu - m*m
}

// ---------------------------------------------------------------------------
// Nakagami Distribution
// ---------------------------------------------------------------------------

// Nakagami represents a Nakagami-m distribution with shape m and spread omega.
type Nakagami struct {
	m     float64
	omega float64
}

// NewNakagami creates a Nakagami distribution. Panics if m < 0.5 or omega <= 0.
func NewNakagami(m, omega float64) *Nakagami {
	if m < 0.5 {
		panic("scigo: Nakagami m must be >= 0.5")
	}
	if omega <= 0 {
		panic("scigo: Nakagami omega must be positive")
	}
	return &Nakagami{m: m, omega: omega}
}

// PDF returns the probability density at x.
func (n *Nakagami) PDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	m, om := n.m, n.omega
	return math.Exp(math.Log(2) + m*math.Log(m) - m*math.Log(om) + (2*m-1)*math.Log(x) -
		m*x*x/om - Gammaln(m))
}

// CDF returns the cumulative distribution function at x.
func (n *Nakagami) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return RegularizedIncompleteGamma(n.m, n.m*x*x/n.omega)
}

// PPF returns the percent point function using Newton's method.
func (n *Nakagami) PPF(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return math.Inf(1)
	}
	x := n.Mean()
	for i := 0; i < 100; i++ {
		cdfVal := n.CDF(x)
		pdfVal := n.PDF(x)
		if pdfVal == 0 {
			break
		}
		dx := (cdfVal - p) / pdfVal
		x -= dx
		if x <= 0 {
			x = 1e-10
		}
		if math.Abs(dx) < 1e-12*(1+x) {
			break
		}
	}
	return x
}

// LogPDF returns the log of the PDF at x.
func (n *Nakagami) LogPDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	m, om := n.m, n.omega
	return math.Log(2) + m*math.Log(m) - m*math.Log(om) + (2*m-1)*math.Log(x) -
		m*x*x/om - Gammaln(m)
}

// Mean returns the mean.
func (n *Nakagami) Mean() float64 {
	return math.Exp(Gammaln(n.m+0.5)-Gammaln(n.m)) * math.Sqrt(n.omega/n.m)
}

// Var returns the variance.
func (n *Nakagami) Var() float64 {
	return n.omega * (1 - (1.0/n.m)*math.Exp(2*(Gammaln(n.m+0.5)-Gammaln(n.m))))
}

// ---------------------------------------------------------------------------
// Von Mises Distribution
// ---------------------------------------------------------------------------

// VonMises represents a Von Mises distribution (circular) with mean mu and concentration kappa.
type VonMises struct {
	mu    float64
	kappa float64
}

// NewVonMises creates a Von Mises distribution. Panics if kappa < 0.
func NewVonMises(mu, kappa float64) *VonMises {
	if kappa < 0 {
		panic("scigo: VonMises kappa must be non-negative")
	}
	return &VonMises{mu: mu, kappa: kappa}
}

// PDF returns the probability density at x.
func (v *VonMises) PDF(x float64) float64 {
	return math.Exp(v.kappa*math.Cos(x-v.mu)) / (2 * math.Pi * besselI0(v.kappa))
}

// CDF returns the cumulative distribution function at x.
// Numerically integrates from -pi to x.
func (v *VonMises) CDF(x float64) float64 {
	// Normalize x to [-pi, pi]
	xn := x
	for xn > math.Pi {
		xn -= 2 * math.Pi
	}
	for xn < -math.Pi {
		xn += 2 * math.Pi
	}
	// Simpson's rule from -pi to xn
	n := 1000
	if n%2 != 0 {
		n++
	}
	a := -math.Pi
	b := xn
	h := (b - a) / float64(n)
	sum := v.PDF(a) + v.PDF(b)
	for i := 1; i < n; i++ {
		xi := a + float64(i)*h
		if i%2 == 0 {
			sum += 2 * v.PDF(xi)
		} else {
			sum += 4 * v.PDF(xi)
		}
	}
	return sum * h / 3
}

// PPF returns the percent point function using Newton's method.
func (v *VonMises) PPF(p float64) float64 {
	if p <= 0 {
		return -math.Pi
	}
	if p >= 1 {
		return math.Pi
	}
	x := v.mu
	for i := 0; i < 100; i++ {
		cdfVal := v.CDF(x)
		pdfVal := v.PDF(x)
		if pdfVal == 0 {
			break
		}
		dx := (cdfVal - p) / pdfVal
		x -= dx
		if math.Abs(dx) < 1e-10 {
			break
		}
	}
	return x
}

// LogPDF returns the log of the PDF at x.
func (v *VonMises) LogPDF(x float64) float64 {
	return v.kappa*math.Cos(x-v.mu) - math.Log(2*math.Pi) - math.Log(besselI0(v.kappa))
}

// Mean returns the mean (circular mean = mu).
func (v *VonMises) Mean() float64 {
	return v.mu
}

// Var returns the circular variance = 1 - I_1(kappa)/I_0(kappa).
func (v *VonMises) Var() float64 {
	return 1 - besselI1(v.kappa)/besselI0(v.kappa)
}

// ---------------------------------------------------------------------------
// Wald (Inverse Gaussian) Distribution
// ---------------------------------------------------------------------------

// Wald represents an inverse Gaussian distribution with mean mu and shape lambda.
type Wald struct {
	mu     float64
	lambda float64
}

// NewWald creates a Wald distribution. Panics if mu <= 0 or lambda <= 0.
func NewWald(mu, lambda float64) *Wald {
	if mu <= 0 {
		panic("scigo: Wald mu must be positive")
	}
	if lambda <= 0 {
		panic("scigo: Wald lambda must be positive")
	}
	return &Wald{mu: mu, lambda: lambda}
}

// PDF returns the probability density at x.
func (w *Wald) PDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Sqrt(w.lambda/(2*math.Pi*x*x*x)) *
		math.Exp(-w.lambda*(x-w.mu)*(x-w.mu)/(2*w.mu*w.mu*x))
}

// CDF returns the cumulative distribution function at x.
func (w *Wald) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	sqrtLx := math.Sqrt(w.lambda / x)
	t1 := sqrtLx * (x/w.mu - 1)
	t2 := sqrtLx * (x/w.mu + 1)
	n := NewNormal(0, 1)
	return n.CDF(t1) + math.Exp(2*w.lambda/w.mu)*n.CDF(-t2)
}

// PPF returns the percent point function using bisection then Newton's method.
func (w *Wald) PPF(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return math.Inf(1)
	}
	// Use bisection to get a good initial bracket
	lo, hi := 1e-10, w.mu*10
	for w.CDF(hi) < p {
		hi *= 2
	}
	for i := 0; i < 60; i++ {
		mid := (lo + hi) / 2
		if w.CDF(mid) < p {
			lo = mid
		} else {
			hi = mid
		}
		if hi-lo < 1e-12*(1+lo) {
			break
		}
	}
	x := (lo + hi) / 2
	// Polish with Newton's method
	for i := 0; i < 20; i++ {
		cdfVal := w.CDF(x)
		pdfVal := w.PDF(x)
		if pdfVal == 0 {
			break
		}
		dx := (cdfVal - p) / pdfVal
		x -= dx
		if x <= 0 {
			x = 1e-10
		}
		if math.Abs(dx) < 1e-12*(1+x) {
			break
		}
	}
	return x
}

// LogPDF returns the log of the PDF at x.
func (w *Wald) LogPDF(x float64) float64 {
	if x <= 0 {
		return math.Inf(-1)
	}
	return 0.5*(math.Log(w.lambda)-math.Log(2*math.Pi)-3*math.Log(x)) -
		w.lambda*(x-w.mu)*(x-w.mu)/(2*w.mu*w.mu*x)
}

// Mean returns the mean.
func (w *Wald) Mean() float64 {
	return w.mu
}

// Var returns the variance.
func (w *Wald) Var() float64 {
	return w.mu * w.mu * w.mu / w.lambda
}

// ---------------------------------------------------------------------------
// Half-Normal Distribution
// ---------------------------------------------------------------------------

// HalfNormal represents a half-normal distribution with parameter sigma.
type HalfNormal struct {
	sigma float64
}

// NewHalfNormal creates a HalfNormal distribution. Panics if sigma <= 0.
func NewHalfNormal(sigma float64) *HalfNormal {
	if sigma <= 0 {
		panic("scigo: HalfNormal sigma must be positive")
	}
	return &HalfNormal{sigma: sigma}
}

// PDF returns the probability density at x.
func (h *HalfNormal) PDF(x float64) float64 {
	if x < 0 {
		return 0
	}
	z := x / h.sigma
	return math.Sqrt(2/(math.Pi)) * math.Exp(-0.5*z*z) / h.sigma
}

// CDF returns the cumulative distribution function at x.
func (h *HalfNormal) CDF(x float64) float64 {
	if x <= 0 {
		return 0
	}
	return math.Erf(x / (h.sigma * math.Sqrt2))
}

// PPF returns the percent point function for probability p.
func (h *HalfNormal) PPF(p float64) float64 {
	if p <= 0 {
		return 0
	}
	if p >= 1 {
		return math.Inf(1)
	}
	return h.sigma * math.Sqrt2 * Erfinv(p)
}

// LogPDF returns the log of the PDF at x.
func (h *HalfNormal) LogPDF(x float64) float64 {
	if x < 0 {
		return math.Inf(-1)
	}
	z := x / h.sigma
	return 0.5*math.Log(2/math.Pi) - math.Log(h.sigma) - 0.5*z*z
}

// Mean returns the mean.
func (h *HalfNormal) Mean() float64 {
	return h.sigma * math.Sqrt(2/math.Pi)
}

// Var returns the variance.
func (h *HalfNormal) Var() float64 {
	return h.sigma * h.sigma * (1 - 2/math.Pi)
}

// ---------------------------------------------------------------------------
// Truncated Normal Distribution
// ---------------------------------------------------------------------------

// TruncatedNormal represents a normal distribution truncated to [a, b].
type TruncatedNormal struct {
	mu    float64
	sigma float64
	a     float64
	b     float64
	norm  *Normal
	phiA  float64 // CDF(a) of the underlying normal
	phiB  float64 // CDF(b) of the underlying normal
	z     float64 // phiB - phiA (normalizing constant)
}

// NewTruncatedNormal creates a TruncatedNormal distribution.
// Panics if sigma <= 0 or a >= b.
func NewTruncatedNormal(mu, sigma, a, b float64) *TruncatedNormal {
	if sigma <= 0 {
		panic("scigo: TruncatedNormal sigma must be positive")
	}
	if a >= b {
		panic("scigo: TruncatedNormal requires a < b")
	}
	n := NewNormal(mu, sigma)
	phiA := n.CDF(a)
	phiB := n.CDF(b)
	return &TruncatedNormal{
		mu: mu, sigma: sigma, a: a, b: b,
		norm: n, phiA: phiA, phiB: phiB, z: phiB - phiA,
	}
}

// PDF returns the probability density at x.
func (tn *TruncatedNormal) PDF(x float64) float64 {
	if x < tn.a || x > tn.b {
		return 0
	}
	return tn.norm.PDF(x) / tn.z
}

// CDF returns the cumulative distribution function at x.
func (tn *TruncatedNormal) CDF(x float64) float64 {
	if x <= tn.a {
		return 0
	}
	if x >= tn.b {
		return 1
	}
	return (tn.norm.CDF(x) - tn.phiA) / tn.z
}

// PPF returns the percent point function for probability p.
func (tn *TruncatedNormal) PPF(p float64) float64 {
	if p <= 0 {
		return tn.a
	}
	if p >= 1 {
		return tn.b
	}
	return tn.norm.PPF(tn.phiA + p*tn.z)
}

// LogPDF returns the log of the PDF at x.
func (tn *TruncatedNormal) LogPDF(x float64) float64 {
	if x < tn.a || x > tn.b {
		return math.Inf(-1)
	}
	return tn.norm.LogPDF(x) - math.Log(tn.z)
}

// Mean returns the mean of the truncated normal.
func (tn *TruncatedNormal) Mean() float64 {
	alphaA := (tn.a - tn.mu) / tn.sigma
	alphaB := (tn.b - tn.mu) / tn.sigma
	stdNorm := NewNormal(0, 1)
	phiA := stdNorm.PDF(alphaA)
	phiB := stdNorm.PDF(alphaB)
	return tn.mu + tn.sigma*(phiA-phiB)/tn.z
}

// Var returns the variance of the truncated normal.
func (tn *TruncatedNormal) Var() float64 {
	alphaA := (tn.a - tn.mu) / tn.sigma
	alphaB := (tn.b - tn.mu) / tn.sigma
	stdNorm := NewNormal(0, 1)
	phiA := stdNorm.PDF(alphaA)
	phiB := stdNorm.PDF(alphaB)
	z := tn.z
	term1 := (alphaA*phiA - alphaB*phiB) / z
	term2 := ((phiA - phiB) / z) * ((phiA - phiB) / z)
	return tn.sigma * tn.sigma * (1 + term1 - term2)
}

// ---------------------------------------------------------------------------
// Skew-Normal Distribution
// ---------------------------------------------------------------------------

// SkewNormal represents a skew-normal distribution with location, scale, and shape (alpha).
type SkewNormal struct {
	loc   float64
	scale float64
	alpha float64
}

// NewSkewNormal creates a SkewNormal distribution. Panics if scale <= 0.
func NewSkewNormal(loc, scale, alpha float64) *SkewNormal {
	if scale <= 0 {
		panic("scigo: SkewNormal scale must be positive")
	}
	return &SkewNormal{loc: loc, scale: scale, alpha: alpha}
}

// PDF returns the probability density at x.
func (sn *SkewNormal) PDF(x float64) float64 {
	z := (x - sn.loc) / sn.scale
	stdNorm := NewNormal(0, 1)
	return 2 * stdNorm.PDF(z) * stdNorm.CDF(sn.alpha*z) / sn.scale
}

// CDF returns the cumulative distribution function at x.
// Uses numerical integration via Simpson's rule.
func (sn *SkewNormal) CDF(x float64) float64 {
	// Use Owen's T function approach or numerical integration
	// Numerical integration from -inf is expensive, so we use a practical lower bound
	lower := sn.loc - 10*sn.scale
	if x <= lower {
		return 0
	}
	n := 1000
	if n%2 != 0 {
		n++
	}
	h := (x - lower) / float64(n)
	sum := sn.PDF(lower) + sn.PDF(x)
	for i := 1; i < n; i++ {
		xi := lower + float64(i)*h
		if i%2 == 0 {
			sum += 2 * sn.PDF(xi)
		} else {
			sum += 4 * sn.PDF(xi)
		}
	}
	return sum * h / 3
}

// PPF returns the percent point function using Newton's method.
func (sn *SkewNormal) PPF(p float64) float64 {
	if p <= 0 {
		return math.Inf(-1)
	}
	if p >= 1 {
		return math.Inf(1)
	}
	// Initial guess from normal approximation
	n := NewNormal(sn.loc, sn.scale)
	x := n.PPF(p)
	for i := 0; i < 100; i++ {
		cdfVal := sn.CDF(x)
		pdfVal := sn.PDF(x)
		if pdfVal == 0 {
			break
		}
		dx := (cdfVal - p) / pdfVal
		x -= dx
		if math.Abs(dx) < 1e-10*(1+math.Abs(x)) {
			break
		}
	}
	return x
}

// LogPDF returns the log of the PDF at x.
func (sn *SkewNormal) LogPDF(x float64) float64 {
	z := (x - sn.loc) / sn.scale
	stdNorm := NewNormal(0, 1)
	return math.Log(2) + stdNorm.LogPDF(z) + math.Log(stdNorm.CDF(sn.alpha*z)) - math.Log(sn.scale)
}

// Mean returns the mean.
func (sn *SkewNormal) Mean() float64 {
	delta := sn.alpha / math.Sqrt(1+sn.alpha*sn.alpha)
	return sn.loc + sn.scale*delta*math.Sqrt(2/math.Pi)
}

// Var returns the variance.
func (sn *SkewNormal) Var() float64 {
	delta := sn.alpha / math.Sqrt(1+sn.alpha*sn.alpha)
	return sn.scale * sn.scale * (1 - 2*delta*delta/math.Pi)
}

// ---------------------------------------------------------------------------
// Generalized Extreme Value Distribution
// ---------------------------------------------------------------------------

// GeneralizedExtremeValue represents a GEV distribution with location mu, scale sigma, shape xi.
type GeneralizedExtremeValue struct {
	mu    float64
	sigma float64
	xi    float64
}

// NewGeneralizedExtremeValue creates a GEV distribution. Panics if sigma <= 0.
func NewGeneralizedExtremeValue(mu, sigma, xi float64) *GeneralizedExtremeValue {
	if sigma <= 0 {
		panic("scigo: GEV sigma must be positive")
	}
	return &GeneralizedExtremeValue{mu: mu, sigma: sigma, xi: xi}
}

// gevT computes t(x) for the GEV distribution.
func (gev *GeneralizedExtremeValue) gevT(x float64) float64 {
	z := (x - gev.mu) / gev.sigma
	if math.Abs(gev.xi) < 1e-10 {
		return math.Exp(-z)
	}
	arg := 1 + gev.xi*z
	if arg <= 0 {
		return math.Inf(1) // outside support
	}
	return math.Pow(arg, -1.0/gev.xi)
}

// PDF returns the probability density at x.
func (gev *GeneralizedExtremeValue) PDF(x float64) float64 {
	z := (x - gev.mu) / gev.sigma
	if math.Abs(gev.xi) < 1e-10 {
		// Gumbel case
		return math.Exp(-z-math.Exp(-z)) / gev.sigma
	}
	arg := 1 + gev.xi*z
	if arg <= 0 {
		return 0
	}
	t := math.Pow(arg, -1.0/gev.xi)
	return math.Pow(arg, -1.0/gev.xi-1) * math.Exp(-t) / gev.sigma
}

// CDF returns the cumulative distribution function at x.
func (gev *GeneralizedExtremeValue) CDF(x float64) float64 {
	t := gev.gevT(x)
	if math.IsInf(t, 1) {
		return 0
	}
	return math.Exp(-t)
}

// PPF returns the percent point function for probability p.
func (gev *GeneralizedExtremeValue) PPF(p float64) float64 {
	if p <= 0 {
		if gev.xi > 0 {
			return gev.mu - gev.sigma/gev.xi
		}
		return math.Inf(-1)
	}
	if p >= 1 {
		if gev.xi < 0 {
			return gev.mu - gev.sigma/gev.xi
		}
		return math.Inf(1)
	}
	if math.Abs(gev.xi) < 1e-10 {
		return gev.mu - gev.sigma*math.Log(-math.Log(p))
	}
	return gev.mu + gev.sigma*(math.Pow(-math.Log(p), -gev.xi)-1)/gev.xi
}

// LogPDF returns the log of the PDF at x.
func (gev *GeneralizedExtremeValue) LogPDF(x float64) float64 {
	z := (x - gev.mu) / gev.sigma
	if math.Abs(gev.xi) < 1e-10 {
		return -z - math.Exp(-z) - math.Log(gev.sigma)
	}
	arg := 1 + gev.xi*z
	if arg <= 0 {
		return math.Inf(-1)
	}
	t := math.Pow(arg, -1.0/gev.xi)
	return (-1.0/gev.xi-1)*math.Log(arg) - t - math.Log(gev.sigma)
}

// Mean returns the mean. Defined for xi < 1.
func (gev *GeneralizedExtremeValue) Mean() float64 {
	if math.Abs(gev.xi) < 1e-10 {
		return gev.mu + gev.sigma*0.5772156649015329 // Euler-Mascheroni
	}
	if gev.xi >= 1 {
		return math.Inf(1)
	}
	return gev.mu + gev.sigma*(math.Gamma(1-gev.xi)-1)/gev.xi
}

// Var returns the variance. Defined for xi < 0.5.
func (gev *GeneralizedExtremeValue) Var() float64 {
	if math.Abs(gev.xi) < 1e-10 {
		return gev.sigma * gev.sigma * math.Pi * math.Pi / 6
	}
	if gev.xi >= 0.5 {
		return math.Inf(1)
	}
	g1 := math.Gamma(1 - gev.xi)
	g2 := math.Gamma(1 - 2*gev.xi)
	return gev.sigma * gev.sigma * (g2 - g1*g1) / (gev.xi * gev.xi)
}

// ---------------------------------------------------------------------------
// Levy Distribution
// ---------------------------------------------------------------------------

// Levy represents a Levy distribution with location and scale parameters.
type Levy struct {
	loc   float64
	scale float64
}

// NewLevy creates a Levy distribution. Panics if scale <= 0.
func NewLevy(loc, scale float64) *Levy {
	if scale <= 0 {
		panic("scigo: Levy scale must be positive")
	}
	return &Levy{loc: loc, scale: scale}
}

// PDF returns the probability density at x.
func (l *Levy) PDF(x float64) float64 {
	if x <= l.loc {
		return 0
	}
	d := x - l.loc
	return math.Sqrt(l.scale/(2*math.Pi)) * math.Exp(-l.scale/(2*d)) / (d * math.Sqrt(d))
}

// CDF returns the cumulative distribution function at x.
func (l *Levy) CDF(x float64) float64 {
	if x <= l.loc {
		return 0
	}
	return math.Erfc(math.Sqrt(l.scale / (2 * (x - l.loc))))
}

// PPF returns the percent point function for probability p.
func (l *Levy) PPF(p float64) float64 {
	if p <= 0 {
		return l.loc
	}
	if p >= 1 {
		return math.Inf(1)
	}
	// CDF(x) = erfc(sqrt(c/(2(x-mu)))) = p
	// sqrt(c/(2(x-mu))) = erfinv(1-p) * sqrt(2)... using erfc inverse
	// erfc(z) = p => z = erfinv(1-p) (approximately, using erfc^{-1})
	// Actually erfc(z) = 1 - erf(z) = p => erf(z) = 1-p => z = erfinv(1-p)
	z := Erfinv(1 - p)
	return l.loc + l.scale/(2*z*z)
}

// LogPDF returns the log of the PDF at x.
func (l *Levy) LogPDF(x float64) float64 {
	if x <= l.loc {
		return math.Inf(-1)
	}
	d := x - l.loc
	return 0.5*math.Log(l.scale/(2*math.Pi)) - l.scale/(2*d) - 1.5*math.Log(d)
}

// Mean returns the mean (infinite for Levy).
func (l *Levy) Mean() float64 {
	return math.Inf(1)
}

// Var returns the variance (infinite for Levy).
func (l *Levy) Var() float64 {
	return math.Inf(1)
}

// ---------------------------------------------------------------------------
// Power Law Distribution
// ---------------------------------------------------------------------------

// PowerLaw represents a power law distribution with exponent alpha and minimum xmin.
// PDF(x) = (alpha-1)/xmin * (x/xmin)^{-alpha} for x >= xmin.
type PowerLaw struct {
	alpha float64
	xmin  float64
}

// NewPowerLaw creates a PowerLaw distribution. Panics if alpha <= 1 or xmin <= 0.
func NewPowerLaw(alpha, xmin float64) *PowerLaw {
	if alpha <= 1 {
		panic("scigo: PowerLaw alpha must be > 1")
	}
	if xmin <= 0 {
		panic("scigo: PowerLaw xmin must be positive")
	}
	return &PowerLaw{alpha: alpha, xmin: xmin}
}

// PDF returns the probability density at x.
func (pl *PowerLaw) PDF(x float64) float64 {
	if x < pl.xmin {
		return 0
	}
	return (pl.alpha - 1) / pl.xmin * math.Pow(x/pl.xmin, -pl.alpha)
}

// CDF returns the cumulative distribution function at x.
func (pl *PowerLaw) CDF(x float64) float64 {
	if x < pl.xmin {
		return 0
	}
	return 1 - math.Pow(pl.xmin/x, pl.alpha-1)
}

// PPF returns the percent point function for probability p.
func (pl *PowerLaw) PPF(p float64) float64 {
	if p <= 0 {
		return pl.xmin
	}
	if p >= 1 {
		return math.Inf(1)
	}
	return pl.xmin / math.Pow(1-p, 1.0/(pl.alpha-1))
}

// LogPDF returns the log of the PDF at x.
func (pl *PowerLaw) LogPDF(x float64) float64 {
	if x < pl.xmin {
		return math.Inf(-1)
	}
	return math.Log(pl.alpha-1) - math.Log(pl.xmin) - pl.alpha*math.Log(x/pl.xmin)
}

// Mean returns the mean. Defined for alpha > 2.
func (pl *PowerLaw) Mean() float64 {
	if pl.alpha <= 2 {
		return math.Inf(1)
	}
	return (pl.alpha - 1) * pl.xmin / (pl.alpha - 2)
}

// Var returns the variance. Defined for alpha > 3.
func (pl *PowerLaw) Var() float64 {
	if pl.alpha <= 3 {
		return math.Inf(1)
	}
	return pl.xmin * pl.xmin * (pl.alpha - 1) /
		((pl.alpha - 2) * (pl.alpha - 2) * (pl.alpha - 3))
}
