package scigo

import "math"

// FFT computes the one-dimensional discrete Fourier transform using the
// Cooley-Tukey radix-2 algorithm. Input length must be a power of 2; if not,
// the input is zero-padded to the next power of 2.
func FFT(x []complex128) []complex128 {
	n := len(x)
	if n == 0 {
		return nil
	}
	// Pad to next power of 2
	m := nextPow2(n)
	data := make([]complex128, m)
	copy(data, x)
	fftInPlace(data, false)
	return data
}

// IFFT computes the one-dimensional inverse discrete Fourier transform.
func IFFT(x []complex128) []complex128 {
	n := len(x)
	if n == 0 {
		return nil
	}
	m := nextPow2(n)
	data := make([]complex128, m)
	copy(data, x)
	fftInPlace(data, true)
	// Scale by 1/N
	fn := complex(float64(m), 0)
	for i := range data {
		data[i] /= fn
	}
	return data
}

// FFTReal computes the FFT of real-valued input data (convenience wrapper).
// Returns the full complex spectrum.
func FFTReal(x []float64) []complex128 {
	cx := make([]complex128, len(x))
	for i, v := range x {
		cx[i] = complex(v, 0)
	}
	return FFT(cx)
}

// RFFT computes the one-dimensional discrete Fourier transform of real input,
// returning only the non-negative frequency terms. For input of length N,
// returns N/2+1 complex values.
func RFFT(x []float64) []complex128 {
	full := FFTReal(x)
	if full == nil {
		return nil
	}
	n := len(full)
	return full[:n/2+1]
}

// IRFFT computes the inverse of RFFT. Takes the non-negative frequency terms
// and returns real-valued output of length n.
// If n is 0, it is inferred as 2*(len(x)-1).
func IRFFT(x []complex128, n int) []float64 {
	if len(x) == 0 {
		return nil
	}
	if n == 0 {
		n = 2 * (len(x) - 1)
	}

	// Reconstruct full spectrum using Hermitian symmetry
	m := nextPow2(n)
	full := make([]complex128, m)
	for i := 0; i < len(x) && i < m; i++ {
		full[i] = x[i]
	}
	// Fill conjugate symmetric part
	for i := 1; i < m/2; i++ {
		if i < len(x) {
			full[m-i] = cmplxConj(x[i])
		}
	}

	result := IFFT(full)
	out := make([]float64, n)
	for i := 0; i < n && i < len(result); i++ {
		out[i] = real(result[i])
	}
	return out
}

// FFT2 computes the two-dimensional discrete Fourier transform.
// Input is a 2D slice of complex numbers (row-major).
func FFT2(x [][]complex128) [][]complex128 {
	rows := len(x)
	if rows == 0 {
		return nil
	}
	cols := len(x[0])
	if cols == 0 {
		return nil
	}

	// FFT along rows
	result := make([][]complex128, rows)
	for i := 0; i < rows; i++ {
		result[i] = FFT(x[i])
	}

	// FFT along columns
	newCols := len(result[0])
	for j := 0; j < newCols; j++ {
		col := make([]complex128, rows)
		for i := 0; i < rows; i++ {
			col[i] = result[i][j]
		}
		colFFT := FFT(col)
		for i := 0; i < len(colFFT); i++ {
			if i < rows {
				result[i][j] = colFFT[i]
			}
		}
	}

	return result
}

// IFFT2 computes the two-dimensional inverse discrete Fourier transform.
func IFFT2(x [][]complex128) [][]complex128 {
	rows := len(x)
	if rows == 0 {
		return nil
	}
	cols := len(x[0])
	if cols == 0 {
		return nil
	}

	// IFFT along rows
	result := make([][]complex128, rows)
	for i := 0; i < rows; i++ {
		result[i] = IFFT(x[i])
	}

	// IFFT along columns
	newCols := len(result[0])
	for j := 0; j < newCols; j++ {
		col := make([]complex128, rows)
		for i := 0; i < rows; i++ {
			col[i] = result[i][j]
		}
		colIFFT := IFFT(col)
		for i := 0; i < len(colIFFT); i++ {
			if i < rows {
				result[i][j] = colIFFT[i]
			}
		}
	}

	return result
}

// FFTFreq returns the DFT sample frequencies for a signal of length n with
// sample spacing d. Returns frequencies in cycles per unit of the sample spacing.
// The returned array has length n and follows the convention:
// [0, 1, ..., n/2-1, -n/2, ..., -1] / (d*n) for even n
// [0, 1, ..., (n-1)/2, -(n-1)/2, ..., -1] / (d*n) for odd n
func FFTFreq(n int, d float64) []float64 {
	if n <= 0 {
		return nil
	}
	if d == 0 {
		d = 1
	}
	freq := make([]float64, n)
	val := 1.0 / (float64(n) * d)
	mid := (n + 1) / 2
	for i := 0; i < mid; i++ {
		freq[i] = float64(i) * val
	}
	for i := mid; i < n; i++ {
		freq[i] = float64(i-n) * val
	}
	return freq
}

// RFFTFreq returns the DFT sample frequencies for use with RFFT.
// Returns n/2+1 frequencies in the range [0, fs/2] where fs = 1/d.
func RFFTFreq(n int, d float64) []float64 {
	if n <= 0 {
		return nil
	}
	if d == 0 {
		d = 1
	}
	nf := n/2 + 1
	freq := make([]float64, nf)
	val := 1.0 / (float64(n) * d)
	for i := 0; i < nf; i++ {
		freq[i] = float64(i) * val
	}
	return freq
}

// ---------------------------------------------------------------------------
// Internal FFT helpers
// ---------------------------------------------------------------------------

// fftInPlace performs an in-place FFT using the Cooley-Tukey radix-2 algorithm.
// inverse=true computes the inverse (without the 1/N scaling).
func fftInPlace(data []complex128, inverse bool) {
	n := len(data)
	if n <= 1 {
		return
	}

	// Bit-reversal permutation
	j := 0
	for i := 0; i < n; i++ {
		if i < j {
			data[i], data[j] = data[j], data[i]
		}
		m := n >> 1
		for m >= 1 && j >= m {
			j -= m
			m >>= 1
		}
		j += m
	}

	// Cooley-Tukey butterfly
	sign := -1.0
	if inverse {
		sign = 1.0
	}

	for size := 2; size <= n; size <<= 1 {
		halfSize := size >> 1
		angle := sign * 2 * math.Pi / float64(size)
		wn := complex(math.Cos(angle), math.Sin(angle))
		for start := 0; start < n; start += size {
			w := complex(1, 0)
			for k := 0; k < halfSize; k++ {
				u := data[start+k]
				t := w * data[start+k+halfSize]
				data[start+k] = u + t
				data[start+k+halfSize] = u - t
				w *= wn
			}
		}
	}
}

// nextPow2 returns the smallest power of 2 >= n.
func nextPow2(n int) int {
	if n <= 1 {
		return 1
	}
	p := 1
	for p < n {
		p <<= 1
	}
	return p
}

// cmplxConj returns the complex conjugate.
func cmplxConj(z complex128) complex128 {
	return complex(real(z), -imag(z))
}
