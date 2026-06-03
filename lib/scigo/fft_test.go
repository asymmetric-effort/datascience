//go:build unit

package scigo

import (
	"math"
	"math/cmplx"
	"testing"
)

func approxEqualCmplx(a, b complex128, tol float64) bool {
	return cmplx.Abs(a-b) < tol
}

func approxEqualFFT(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

// ---------------------------------------------------------------------------
// FFT / IFFT
// ---------------------------------------------------------------------------

func TestFFT_DC(t *testing.T) {
	// All ones -> DC component = N, all others = 0
	x := []complex128{1, 1, 1, 1}
	y := FFT(x)
	if len(y) != 4 {
		t.Fatalf("FFT length=%v, want 4", len(y))
	}
	if !approxEqualCmplx(y[0], 4, 1e-10) {
		t.Errorf("FFT DC=%v, want 4", y[0])
	}
	for i := 1; i < 4; i++ {
		if !approxEqualCmplx(y[i], 0, 1e-10) {
			t.Errorf("FFT[%d]=%v, want 0", i, y[i])
		}
	}
}

func TestFFT_Impulse(t *testing.T) {
	// Delta function: [1,0,0,0] -> all ones
	x := []complex128{1, 0, 0, 0}
	y := FFT(x)
	for i := 0; i < 4; i++ {
		if !approxEqualCmplx(y[i], 1, 1e-10) {
			t.Errorf("FFT[%d]=%v, want 1", i, y[i])
		}
	}
}

func TestFFT_IFFT_Roundtrip(t *testing.T) {
	x := []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}
	y := FFT(x)
	z := IFFT(y)
	for i := 0; i < 4; i++ {
		if !approxEqualCmplx(z[i], x[i], 1e-10) {
			t.Errorf("IFFT(FFT(x))[%d]=%v, want %v", i, z[i], x[i])
		}
	}
}

func TestFFT_Sinusoid(t *testing.T) {
	// 8-point FFT of a sine wave at frequency 1
	n := 8
	x := make([]complex128, n)
	for i := 0; i < n; i++ {
		x[i] = complex(math.Sin(2*math.Pi*float64(i)/float64(n)), 0)
	}
	y := FFT(x)
	// Sine wave at freq 1 should have peaks at bins 1 and N-1 (bin 7)
	// with magnitude N/2 = 4
	if cmplx.Abs(y[1]) < 3 {
		t.Errorf("FFT sine: |Y[1]|=%v, should be ~4", cmplx.Abs(y[1]))
	}
	if cmplx.Abs(y[7]) < 3 {
		t.Errorf("FFT sine: |Y[7]|=%v, should be ~4", cmplx.Abs(y[7]))
	}
	// DC should be ~0
	if cmplx.Abs(y[0]) > 1e-10 {
		t.Errorf("FFT sine: |Y[0]|=%v, should be ~0", cmplx.Abs(y[0]))
	}
}

// ---------------------------------------------------------------------------
// FFTReal / RFFT / IRFFT
// ---------------------------------------------------------------------------

func TestFFTReal(t *testing.T) {
	x := []float64{1, 2, 3, 4}
	y := FFTReal(x)
	if len(y) != 4 {
		t.Fatalf("FFTReal length=%v, want 4", len(y))
	}
	// DC = sum = 10
	if !approxEqualCmplx(y[0], 10, 1e-10) {
		t.Errorf("FFTReal DC=%v, want 10", y[0])
	}
}

func TestRFFT(t *testing.T) {
	x := []float64{1, 2, 3, 4}
	y := RFFT(x)
	// For N=4, RFFT returns N/2+1 = 3 values
	// But after padding to power-of-2, N is still 4
	if len(y) != 3 {
		t.Fatalf("RFFT length=%v, want 3", len(y))
	}
	// DC = sum = 10
	if !approxEqualCmplx(y[0], 10, 1e-10) {
		t.Errorf("RFFT DC=%v, want 10", y[0])
	}
}

func TestIRFFT_Roundtrip(t *testing.T) {
	x := []float64{1, 2, 3, 4}
	y := RFFT(x)
	z := IRFFT(y, 4)
	for i := 0; i < 4; i++ {
		if !approxEqualFFT(z[i], x[i], 1e-10) {
			t.Errorf("IRFFT(RFFT(x))[%d]=%v, want %v", i, z[i], x[i])
		}
	}
}

// ---------------------------------------------------------------------------
// FFT2 / IFFT2
// ---------------------------------------------------------------------------

func TestFFT2_DC(t *testing.T) {
	// 2x2 matrix of ones -> only DC is nonzero = 4
	x := [][]complex128{
		{1, 1},
		{1, 1},
	}
	y := FFT2(x)
	if y == nil {
		t.Fatal("FFT2 returned nil")
	}
	// After padding, each row becomes length 2 (already power of 2)
	if !approxEqualCmplx(y[0][0], 4, 1e-10) {
		t.Errorf("FFT2 DC=%v, want 4", y[0][0])
	}
}

func TestFFT2_IFFT2_Roundtrip(t *testing.T) {
	x := [][]complex128{
		{1 + 0i, 2 + 0i},
		{3 + 0i, 4 + 0i},
	}
	y := FFT2(x)
	z := IFFT2(y)
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if !approxEqualCmplx(z[i][j], x[i][j], 1e-10) {
				t.Errorf("IFFT2(FFT2(x))[%d][%d]=%v, want %v", i, j, z[i][j], x[i][j])
			}
		}
	}
}

// ---------------------------------------------------------------------------
// FFTFreq / RFFTFreq
// ---------------------------------------------------------------------------

func TestFFTFreq(t *testing.T) {
	// For n=4, d=1: [0, 0.25, -0.5, -0.25]
	f := FFTFreq(4, 1.0)
	if len(f) != 4 {
		t.Fatalf("FFTFreq length=%v, want 4", len(f))
	}
	expected := []float64{0, 0.25, -0.5, -0.25}
	for i := range f {
		if !approxEqualFFT(f[i], expected[i], 1e-10) {
			t.Errorf("FFTFreq[%d]=%v, want %v", i, f[i], expected[i])
		}
	}
}

func TestFFTFreq_WithSpacing(t *testing.T) {
	// For n=4, d=0.5: frequencies are divided by 0.5*4=2
	f := FFTFreq(4, 0.5)
	expected := []float64{0, 0.5, -1.0, -0.5}
	for i := range f {
		if !approxEqualFFT(f[i], expected[i], 1e-10) {
			t.Errorf("FFTFreq(d=0.5)[%d]=%v, want %v", i, f[i], expected[i])
		}
	}
}

func TestRFFTFreq(t *testing.T) {
	// For n=8, d=1: [0, 0.125, 0.25, 0.375, 0.5]
	f := RFFTFreq(8, 1.0)
	if len(f) != 5 {
		t.Fatalf("RFFTFreq length=%v, want 5", len(f))
	}
	for i := range f {
		expected := float64(i) / 8.0
		if !approxEqualFFT(f[i], expected, 1e-10) {
			t.Errorf("RFFTFreq[%d]=%v, want %v", i, f[i], expected)
		}
	}
}

func TestFFTFreq_Odd(t *testing.T) {
	// For n=5, d=1: [0, 0.2, 0.4, -0.4, -0.2]
	f := FFTFreq(5, 1.0)
	if len(f) != 5 {
		t.Fatalf("FFTFreq(5) length=%v, want 5", len(f))
	}
	expected := []float64{0, 0.2, 0.4, -0.4, -0.2}
	for i := range f {
		if !approxEqualFFT(f[i], expected[i], 1e-10) {
			t.Errorf("FFTFreq(5)[%d]=%v, want %v", i, f[i], expected[i])
		}
	}
}

// ---------------------------------------------------------------------------
// Parseval's theorem: energy in time domain = energy in freq domain / N
// ---------------------------------------------------------------------------

func TestFFT_Parseval(t *testing.T) {
	x := []complex128{1, 2, 3, 4, 5, 6, 7, 8}
	y := FFT(x)
	n := float64(len(y))

	timeEnergy := 0.0
	for _, v := range x {
		timeEnergy += real(v)*real(v) + imag(v)*imag(v)
	}
	freqEnergy := 0.0
	for _, v := range y {
		freqEnergy += real(v)*real(v) + imag(v)*imag(v)
	}
	freqEnergy /= n

	if !approxEqualFFT(timeEnergy, freqEnergy, 1e-8) {
		t.Errorf("Parseval: time energy=%v, freq energy/N=%v", timeEnergy, freqEnergy)
	}
}

// ---------------------------------------------------------------------------
// Edge cases
// ---------------------------------------------------------------------------

func TestFFT_Empty(t *testing.T) {
	if FFT(nil) != nil {
		t.Error("FFT(nil) should return nil")
	}
}

func TestIFFT_Empty(t *testing.T) {
	if IFFT(nil) != nil {
		t.Error("IFFT(nil) should return nil")
	}
}

func TestFFTFreq_Zero(t *testing.T) {
	if FFTFreq(0, 1) != nil {
		t.Error("FFTFreq(0) should return nil")
	}
}
