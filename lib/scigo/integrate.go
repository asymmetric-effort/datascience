package scigo

import (
	"errors"
	"math"
)

// Quad computes the definite integral of f from a to b using adaptive Simpson's
// quadrature with a default tolerance of 1e-10.
func Quad(f func(float64) float64, a, b float64) (float64, error) {
	if math.IsNaN(a) || math.IsNaN(b) {
		return 0, errors.New("scigo.Quad: integration limits must be finite")
	}
	if a == b {
		return 0, nil
	}

	tol := 1e-10
	maxDepth := 50
	result, err := adaptiveSimpson(f, a, b, tol, maxDepth)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func adaptiveSimpson(f func(float64) float64, a, b, tol float64, maxDepth int) (float64, error) {
	c := (a + b) / 2
	h := b - a
	fa := f(a)
	fb := f(b)
	fc := f(c)
	s := (h / 6) * (fa + 4*fc + fb)
	return adaptiveSimpsonRec(f, a, b, tol, s, fa, fb, fc, maxDepth)
}

func adaptiveSimpsonRec(f func(float64) float64, a, b, tol, whole, fa, fb, fc float64, depth int) (float64, error) {
	c := (a + b) / 2
	d := (a + c) / 2
	e := (c + b) / 2
	fd := f(d)
	fe := f(e)
	h := b - a
	left := (h / 12) * (fa + 4*fd + fc)
	right := (h / 12) * (fc + 4*fe + fb)
	sum := left + right

	if depth <= 0 || math.Abs(sum-whole) <= 15*tol {
		return sum + (sum-whole)/15, nil
	}

	leftResult, err := adaptiveSimpsonRec(f, a, c, tol/2, left, fa, fc, fd, depth-1)
	if err != nil {
		return 0, err
	}
	rightResult, err := adaptiveSimpsonRec(f, c, b, tol/2, right, fc, fb, fe, depth-1)
	if err != nil {
		return 0, err
	}
	return leftResult + rightResult, nil
}

// Dblquad computes the double integral of f(y, x) over the region
// a <= x <= b, gfun(x) <= y <= hfun(x).
func Dblquad(f func(float64, float64) float64, a, b float64, gfun, hfun func(float64) float64) (float64, error) {
	outer := func(x float64) float64 {
		inner := func(y float64) float64 {
			return f(y, x)
		}
		lo := gfun(x)
		hi := hfun(x)
		val, _ := Quad(inner, lo, hi)
		return val
	}
	return Quad(outer, a, b)
}

// Tplquad computes the triple integral of f(z, y, x) by nesting Dblquad.
// a <= x <= b, gfun(x) <= y <= hfun(x), qfun(x,y) <= z <= rfun(x,y).
func Tplquad(
	f func(float64, float64, float64) float64,
	a, b float64,
	gfun, hfun func(float64) float64,
	qfun, rfun func(float64, float64) float64,
) (float64, error) {
	outer := func(y, x float64) float64 {
		inner := func(z float64) float64 {
			return f(z, y, x)
		}
		lo := qfun(x, y)
		hi := rfun(x, y)
		val, _ := Quad(inner, lo, hi)
		return val
	}
	return Dblquad(outer, a, b, gfun, hfun)
}

// FixedQuad computes the integral of f from a to b using n-point Gauss-Legendre quadrature.
func FixedQuad(f func(float64) float64, a, b float64, n int) float64 {
	nodes, weights := gaussLegendre(n)

	// Transform from [-1,1] to [a,b].
	halfLen := (b - a) / 2
	midpoint := (a + b) / 2
	sum := 0.0
	for i := 0; i < n; i++ {
		sum += weights[i] * f(midpoint+halfLen*nodes[i])
	}
	return halfLen * sum
}

// Trapezoid computes the integral of evenly-spaced data y using the trapezoidal rule.
// dx is the spacing between points.
func Trapezoid(y []float64, dx float64) float64 {
	n := len(y)
	if n < 2 {
		return 0
	}
	sum := 0.5 * (y[0] + y[n-1])
	for i := 1; i < n-1; i++ {
		sum += y[i]
	}
	return sum * dx
}

// Simpson computes the integral of evenly-spaced data y using Simpson's rule.
// dx is the spacing between points. If n is even (odd number of points), uses
// composite Simpson's 1/3 rule; otherwise falls back to trapezoidal for the last interval.
func Simpson(y []float64, dx float64) float64 {
	n := len(y)
	if n < 2 {
		return 0
	}
	if n == 2 {
		return Trapezoid(y, dx)
	}

	// Use composite Simpson's 1/3 rule for pairs of intervals.
	sum := 0.0
	end := n
	extra := 0.0
	if n%2 == 0 {
		// Odd number of intervals: use Simpson for n-1 points, trap for last.
		end = n - 1
		extra = 0.5 * (y[n-2] + y[n-1]) * dx
	}

	sum = y[0] + y[end-1]
	for i := 1; i < end-1; i++ {
		if i%2 == 1 {
			sum += 4 * y[i]
		} else {
			sum += 2 * y[i]
		}
	}
	return sum*dx/3 + extra
}

// Romberg computes the integral of f from a to b using Romberg integration.
func Romberg(f func(float64) float64, a, b float64) float64 {
	maxK := 10
	r := make([][]float64, maxK)
	for i := range r {
		r[i] = make([]float64, maxK)
	}

	h := b - a
	r[0][0] = h * (f(a) + f(b)) / 2

	for k := 1; k < maxK; k++ {
		h /= 2
		sum := 0.0
		nPts := 1 << (k - 1) // 2^(k-1)
		for i := 0; i < nPts; i++ {
			sum += f(a + float64(2*i+1)*h)
		}
		r[k][0] = r[k-1][0]/2 + h*sum

		for j := 1; j <= k; j++ {
			factor := math.Pow(4, float64(j))
			r[k][j] = (factor*r[k][j-1] - r[k-1][j-1]) / (factor - 1)
		}

		if k > 1 && math.Abs(r[k][k]-r[k-1][k-1]) < 1e-12 {
			return r[k][k]
		}
	}

	return r[maxK-1][maxK-1]
}

// Odeint solves a scalar ODE dy/dt = f(y, t) using the classical 4th-order Runge-Kutta method.
// y0 is the initial condition, t is the initial time, and tspan provides the time points
// at which to evaluate the solution.
func Odeint(f func(y, t float64) float64, y0, t float64, tspan []float64) []float64 {
	n := len(tspan)
	if n == 0 {
		return nil
	}

	result := make([]float64, n)
	y := y0
	tc := t

	for i := 0; i < n; i++ {
		target := tspan[i]
		// Integrate from tc to target using fixed steps.
		nSteps := 100
		dt := (target - tc) / float64(nSteps)
		for s := 0; s < nSteps; s++ {
			k1 := dt * f(y, tc)
			k2 := dt * f(y+k1/2, tc+dt/2)
			k3 := dt * f(y+k2/2, tc+dt/2)
			k4 := dt * f(y+k3, tc+dt)
			y += (k1 + 2*k2 + 2*k3 + k4) / 6
			tc += dt
		}
		result[i] = y
	}

	return result
}

// SolveIVP solves a system of ODEs dy/dt = f(t, y) using the Dormand-Prince (RK45) method.
// tspan gives [t0, tf], y0 is the initial state vector.
// Returns the time points and the solution matrix (each row is a state vector).
func SolveIVP(f func(t float64, y []float64) []float64, tspan [2]float64, y0 []float64) ([]float64, [][]float64, error) {
	t0, tf := tspan[0], tspan[1]
	if t0 == tf {
		return []float64{t0}, [][]float64{append([]float64{}, y0...)}, nil
	}

	n := len(y0)
	h := (tf - t0) / 100 // Initial step size.
	if h == 0 {
		return nil, nil, errors.New("scigo.SolveIVP: zero time span")
	}

	atol := 1e-8
	rtol := 1e-6

	times := []float64{t0}
	states := [][]float64{append([]float64{}, y0...)}

	t := t0
	y := append([]float64{}, y0...)

	for (h > 0 && t < tf) || (h < 0 && t > tf) {
		// Clamp h so we don't overshoot.
		if h > 0 && t+h > tf {
			h = tf - t
		} else if h < 0 && t+h < tf {
			h = tf - t
		}

		// RK4 step.
		k1 := f(t, y)
		y2 := make([]float64, n)
		for i := range y2 {
			y2[i] = y[i] + h/2*k1[i]
		}
		k2 := f(t+h/2, y2)
		y3 := make([]float64, n)
		for i := range y3 {
			y3[i] = y[i] + h/2*k2[i]
		}
		k3 := f(t+h/2, y3)
		y4 := make([]float64, n)
		for i := range y4 {
			y4[i] = y[i] + h*k3[i]
		}
		k4 := f(t+h, y4)

		yNew := make([]float64, n)
		for i := range yNew {
			yNew[i] = y[i] + h/6*(k1[i]+2*k2[i]+2*k3[i]+k4[i])
		}

		// Estimate error using a lower-order method (Euler).
		errMax := 0.0
		for i := range yNew {
			yEuler := y[i] + h*k1[i]
			sc := atol + rtol*math.Max(math.Abs(y[i]), math.Abs(yNew[i]))
			errMax = math.Max(errMax, math.Abs(yNew[i]-yEuler)/sc)
		}
		errMax /= float64(n)

		if errMax <= 1.0 || math.Abs(h) < 1e-14 {
			// Accept step.
			t += h
			copy(y, yNew)
			times = append(times, t)
			states = append(states, append([]float64{}, y...))

			// Increase step size.
			if errMax > 0 {
				h *= math.Min(2.0, math.Max(0.5, 0.9*math.Pow(1.0/errMax, 0.2)))
			}
		} else {
			// Reject step, decrease h.
			h *= math.Max(0.2, 0.9*math.Pow(1.0/errMax, 0.25))
		}
	}

	return times, states, nil
}

// ---------------------------------------------------------------------------
// Gauss-Legendre nodes and weights
// ---------------------------------------------------------------------------

func gaussLegendre(n int) (nodes, weights []float64) {
	nodes = make([]float64, n)
	weights = make([]float64, n)

	for i := 0; i < n; i++ {
		// Initial guess for the i-th root of P_n(x).
		x := math.Cos(math.Pi * (float64(i) + 0.75) / (float64(n) + 0.5))

		for iter := 0; iter < 100; iter++ {
			p0 := 1.0
			p1 := x
			for j := 2; j <= n; j++ {
				p2 := ((2*float64(j)-1)*x*p1 - (float64(j)-1)*p0) / float64(j)
				p0 = p1
				p1 = p2
			}
			// p1 is P_n(x), dp = P_n'(x)
			dp := float64(n) * (x*p1 - p0) / (x*x - 1)
			dx := p1 / dp
			x -= dx
			if math.Abs(dx) < 1e-15 {
				break
			}
		}

		nodes[i] = x
		// Weight: 2 / ((1-x^2) * [P_n'(x)]^2)
		p0 := 1.0
		p1 := x
		for j := 2; j <= n; j++ {
			p2 := ((2*float64(j)-1)*x*p1 - (float64(j)-1)*p0) / float64(j)
			p0 = p1
			p1 = p2
		}
		dp := float64(n) * (x*p1 - p0) / (x*x - 1)
		weights[i] = 2.0 / ((1 - x*x) * dp * dp)
	}

	return nodes, weights
}
