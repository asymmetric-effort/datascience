package scigo

import "math"

// SignalConvolve computes the full discrete linear convolution of a and b.
// The output length is len(a) + len(b) - 1.
func SignalConvolve(a, b []float64) []float64 {
	na, nb := len(a), len(b)
	if na == 0 || nb == 0 {
		return nil
	}
	n := na + nb - 1
	out := make([]float64, n)
	for i := 0; i < na; i++ {
		for j := 0; j < nb; j++ {
			out[i+j] += a[i] * b[j]
		}
	}
	return out
}

// SignalCorrelate computes the full discrete linear cross-correlation of a and b.
// Equivalent to convolving a with the reversed b.
func SignalCorrelate(a, b []float64) []float64 {
	nb := len(b)
	if nb == 0 {
		return nil
	}
	rev := make([]float64, nb)
	for i := range rev {
		rev[i] = b[nb-1-i]
	}
	return SignalConvolve(a, rev)
}

// FFTConvolve computes the convolution of a and b. This simplified implementation
// uses direct convolution (same result as FFT-based, suitable for moderate sizes).
func FFTConvolve(a, b []float64) []float64 {
	return SignalConvolve(a, b)
}

// Welch estimates the power spectral density using the periodogram method.
// x is the input signal, fs is the sampling frequency, nperseg is the segment length.
// Returns frequency bins and corresponding PSD estimates.
func Welch(x []float64, fs float64, nperseg int) (freqs, psd []float64) {
	n := len(x)
	if n == 0 || nperseg <= 0 {
		return nil, nil
	}
	if nperseg > n {
		nperseg = n
	}

	// Number of frequency bins (one-sided).
	nfreqs := nperseg/2 + 1
	freqs = make([]float64, nfreqs)
	psd = make([]float64, nfreqs)

	for i := 0; i < nfreqs; i++ {
		freqs[i] = float64(i) * fs / float64(nperseg)
	}

	// Overlap: 50%.
	step := nperseg / 2
	if step == 0 {
		step = 1
	}

	nSegments := 0
	for start := 0; start+nperseg <= n; start += step {
		seg := x[start : start+nperseg]

		// Apply Hann window.
		windowed := make([]float64, nperseg)
		winSum := 0.0
		for i := 0; i < nperseg; i++ {
			w := 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(nperseg-1)))
			windowed[i] = seg[i] * w
			winSum += w * w
		}

		// DFT for each frequency bin.
		for k := 0; k < nfreqs; k++ {
			re, im := 0.0, 0.0
			for t := 0; t < nperseg; t++ {
				angle := 2 * math.Pi * float64(k) * float64(t) / float64(nperseg)
				re += windowed[t] * math.Cos(angle)
				im -= windowed[t] * math.Sin(angle)
			}
			power := (re*re + im*im) / (fs * winSum)
			// Double for one-sided, except DC and Nyquist.
			if k > 0 && k < nfreqs-1 {
				power *= 2
			}
			psd[k] += power
		}
		nSegments++
	}

	// Average over segments.
	if nSegments > 0 {
		for i := range psd {
			psd[i] /= float64(nSegments)
		}
	}

	return freqs, psd
}

// Spectrogram computes a simplified spectrogram (STFT magnitude squared).
// Returns time bins, frequency bins, and the spectrogram matrix.
func Spectrogram(x []float64, fs float64, nperseg int) (times, freqs []float64, sxx [][]float64) {
	n := len(x)
	if n == 0 || nperseg <= 0 {
		return nil, nil, nil
	}
	if nperseg > n {
		nperseg = n
	}

	nfreqs := nperseg/2 + 1
	freqs = make([]float64, nfreqs)
	for i := 0; i < nfreqs; i++ {
		freqs[i] = float64(i) * fs / float64(nperseg)
	}

	step := nperseg / 2
	if step == 0 {
		step = 1
	}

	for start := 0; start+nperseg <= n; start += step {
		seg := x[start : start+nperseg]
		times = append(times, (float64(start)+float64(nperseg)/2)/fs)

		// Apply Hann window.
		windowed := make([]float64, nperseg)
		winSum := 0.0
		for i := 0; i < nperseg; i++ {
			w := 0.5 * (1 - math.Cos(2*math.Pi*float64(i)/float64(nperseg-1)))
			windowed[i] = seg[i] * w
			winSum += w * w
		}

		spectrum := make([]float64, nfreqs)
		for k := 0; k < nfreqs; k++ {
			re, im := 0.0, 0.0
			for t := 0; t < nperseg; t++ {
				angle := 2 * math.Pi * float64(k) * float64(t) / float64(nperseg)
				re += windowed[t] * math.Cos(angle)
				im -= windowed[t] * math.Sin(angle)
			}
			power := (re*re + im*im) / (fs * winSum)
			if k > 0 && k < nfreqs-1 {
				power *= 2
			}
			spectrum[k] = power
		}
		sxx = append(sxx, spectrum)
	}

	return times, freqs, sxx
}

// Butter computes the Butterworth lowpass filter coefficients for the given order,
// cutoff frequency, and sampling frequency. Returns the numerator (b) and
// denominator (a) of the transfer function in the z-domain.
// Uses the bilinear transform of the analog prototype.
func Butter(order int, cutoff float64, fs float64) (b, a []float64) {
	if order < 1 || cutoff <= 0 || fs <= 0 || cutoff >= fs/2 {
		return nil, nil
	}

	// Pre-warp the cutoff frequency.
	wc := 2 * fs * math.Tan(math.Pi*cutoff/fs)

	// Analog prototype poles (Butterworth): s_k = wc * exp(j*pi*(2k+order+1)/(2*order))
	// for k = 0..order-1. We keep the left half-plane poles.
	type complex struct{ re, im float64 }
	poles := make([]complex, order)
	for k := 0; k < order; k++ {
		angle := math.Pi * float64(2*k+order+1) / float64(2*order)
		poles[k] = complex{wc * math.Cos(angle), wc * math.Sin(angle)}
	}

	// Bilinear transform: z = (1 + s/(2*fs)) / (1 - s/(2*fs))
	// Map analog poles to digital poles.
	type complexVal struct{ re, im float64 }
	zPoles := make([]complexVal, order)
	for k := 0; k < order; k++ {
		// s -> (2*fs * (z-1)/(z+1)), so z = (1 + s/2fs) / (1 - s/2fs)
		sre := poles[k].re / (2 * fs)
		sim := poles[k].im / (2 * fs)
		// (1 + s) / (1 - s) where s = sre + j*sim
		numRe := 1 + sre
		numIm := sim
		denRe := 1 - sre
		denIm := -sim
		denom := denRe*denRe + denIm*denIm
		zPoles[k] = complexVal{
			re: (numRe*denRe + numIm*denIm) / denom,
			im: (numIm*denRe - numRe*denIm) / denom,
		}
	}

	// All zeros at z = -1 for lowpass.
	// Expand polynomial from poles.
	// a[z] = prod(z - zp_k)
	// Start with [1].
	aCoeffs := []complexVal{{1, 0}}
	for k := 0; k < order; k++ {
		newA := make([]complexVal, len(aCoeffs)+1)
		for i := range aCoeffs {
			// Multiply by (z - pole)
			newA[i] = complexVal{
				re: newA[i].re + aCoeffs[i].re,
				im: newA[i].im + aCoeffs[i].im,
			}
			newA[i+1] = complexVal{
				re: newA[i+1].re - aCoeffs[i].re*zPoles[k].re + aCoeffs[i].im*zPoles[k].im,
				im: newA[i+1].im - aCoeffs[i].re*zPoles[k].im - aCoeffs[i].im*zPoles[k].re,
			}
		}
		aCoeffs = newA
	}

	// Extract real parts (should be real for conjugate pairs).
	a = make([]float64, len(aCoeffs))
	for i := range aCoeffs {
		a[i] = aCoeffs[i].re
	}

	// b coefficients: all zeros at z = -1 => b[z] = (1+z^{-1})^order = sum C(order,k) z^{-k}
	b = make([]float64, order+1)
	for k := 0; k <= order; k++ {
		b[k] = binomial(order, k)
	}

	// Normalize so that sum(b)/sum(a) matches DC gain = 1.
	sumB := 0.0
	sumA := 0.0
	for i := range b {
		sumB += b[i]
	}
	for i := range a {
		sumA += a[i]
	}
	if sumA != 0 {
		gain := sumA / sumB
		for i := range b {
			b[i] *= gain
		}
	}

	return b, a
}

// LFilter applies a linear filter to x using the difference equation
// (direct form II transposed). b and a are the numerator and denominator
// coefficients of the transfer function.
func LFilter(b, a, x []float64) []float64 {
	if len(a) == 0 || len(b) == 0 || len(x) == 0 {
		return nil
	}

	// Normalize by a[0].
	a0 := a[0]
	if a0 == 0 {
		return nil
	}

	nb := len(b)
	na := len(a)
	n := len(x)

	// Normalized coefficients.
	bn := make([]float64, nb)
	an := make([]float64, na)
	for i := range bn {
		bn[i] = b[i] / a0
	}
	for i := range an {
		an[i] = a[i] / a0
	}

	y := make([]float64, n)
	// State for direct form II transposed.
	nstate := nb
	if na > nstate {
		nstate = na
	}
	d := make([]float64, nstate)

	for i := 0; i < n; i++ {
		y[i] = bn[0]*x[i] + d[0]
		// Shift state.
		for j := 0; j < nstate-1; j++ {
			d[j] = d[j+1]
			if j+1 < nb {
				d[j] += bn[j+1] * x[i]
			}
			if j+1 < na {
				d[j] -= an[j+1] * y[i]
			}
		}
		d[nstate-1] = 0
		if nstate < nb {
			d[nstate-1] += bn[nstate] * x[i]
		}
		if nstate < na {
			d[nstate-1] -= an[nstate] * y[i]
		}
	}

	return y
}
