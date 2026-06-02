//go:build unit

package scigo

import (
	"math"
	"testing"
)

// ---------------------------------------------------------------------------
// Convolution and Correlation Tests
// ---------------------------------------------------------------------------

func TestSignalConvolve(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{0, 1, 0.5}
	result := SignalConvolve(a, b)
	// Expected: [0, 1, 2.5, 4, 1.5]
	expected := []float64{0, 1, 2.5, 4, 1.5}
	if len(result) != len(expected) {
		t.Fatalf("SignalConvolve: length = %d, want %d", len(result), len(expected))
	}
	for i := range expected {
		if !approxEqual(result[i], expected[i], 1e-12) {
			t.Errorf("SignalConvolve[%d] = %v, want %v", i, result[i], expected[i])
		}
	}
}

func TestSignalConvolveIdentity(t *testing.T) {
	a := []float64{1, 2, 3, 4}
	delta := []float64{1}
	result := SignalConvolve(a, delta)
	if len(result) != len(a) {
		t.Fatalf("length mismatch")
	}
	for i := range a {
		if result[i] != a[i] {
			t.Errorf("result[%d] = %v, want %v", i, result[i], a[i])
		}
	}
}

func TestSignalConvolveEmpty(t *testing.T) {
	if SignalConvolve(nil, []float64{1}) != nil {
		t.Error("expected nil for empty a")
	}
	if SignalConvolve([]float64{1}, nil) != nil {
		t.Error("expected nil for empty b")
	}
}

func TestSignalCorrelate(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{0, 1, 0.5}
	result := SignalCorrelate(a, b)
	// Correlation = convolve(a, reverse(b)) = convolve([1,2,3], [0.5,1,0])
	expected := SignalConvolve(a, []float64{0.5, 1, 0})
	if len(result) != len(expected) {
		t.Fatalf("SignalCorrelate: length mismatch")
	}
	for i := range expected {
		if !approxEqual(result[i], expected[i], 1e-12) {
			t.Errorf("SignalCorrelate[%d] = %v, want %v", i, result[i], expected[i])
		}
	}
}

func TestFFTConvolve(t *testing.T) {
	a := []float64{1, 2, 3}
	b := []float64{4, 5}
	direct := SignalConvolve(a, b)
	fft := FFTConvolve(a, b)
	if len(direct) != len(fft) {
		t.Fatalf("length mismatch")
	}
	for i := range direct {
		if !approxEqual(direct[i], fft[i], 1e-12) {
			t.Errorf("FFTConvolve[%d] = %v, want %v", i, fft[i], direct[i])
		}
	}
}

// ---------------------------------------------------------------------------
// Welch Tests
// ---------------------------------------------------------------------------

func TestWelch(t *testing.T) {
	// Generate a simple sine wave.
	fs := 100.0
	n := 256
	x := make([]float64, n)
	freq := 10.0 // 10 Hz
	for i := 0; i < n; i++ {
		x[i] = math.Sin(2 * math.Pi * freq * float64(i) / fs)
	}

	freqs, psd := Welch(x, fs, 64)
	if len(freqs) == 0 || len(psd) == 0 {
		t.Fatal("Welch: empty result")
	}
	if len(freqs) != len(psd) {
		t.Fatal("Welch: freqs and psd length mismatch")
	}

	// The peak should be near 10 Hz.
	maxIdx := 0
	maxVal := psd[0]
	for i, v := range psd {
		if v > maxVal {
			maxVal = v
			maxIdx = i
		}
	}
	peakFreq := freqs[maxIdx]
	if math.Abs(peakFreq-freq) > 5 {
		t.Errorf("Welch: peak frequency = %v, expected near %v", peakFreq, freq)
	}
}

func TestWelchEmpty(t *testing.T) {
	f, p := Welch(nil, 1, 1)
	if f != nil || p != nil {
		t.Error("Welch: expected nil for empty input")
	}
}

// ---------------------------------------------------------------------------
// Spectrogram Tests
// ---------------------------------------------------------------------------

func TestSpectrogram(t *testing.T) {
	fs := 100.0
	n := 256
	x := make([]float64, n)
	for i := 0; i < n; i++ {
		x[i] = math.Sin(2 * math.Pi * 10 * float64(i) / fs)
	}

	times, freqs, sxx := Spectrogram(x, fs, 64)
	if len(times) == 0 || len(freqs) == 0 || len(sxx) == 0 {
		t.Fatal("Spectrogram: empty result")
	}
	if len(sxx) != len(times) {
		t.Errorf("Spectrogram: sxx rows (%d) != times length (%d)", len(sxx), len(times))
	}
	if len(sxx[0]) != len(freqs) {
		t.Errorf("Spectrogram: sxx cols (%d) != freqs length (%d)", len(sxx[0]), len(freqs))
	}
}

// ---------------------------------------------------------------------------
// Butter Tests
// ---------------------------------------------------------------------------

func TestButter(t *testing.T) {
	b, a := Butter(2, 10, 100)
	if b == nil || a == nil {
		t.Fatal("Butter: nil result")
	}
	if len(b) != 3 {
		t.Errorf("Butter: len(b) = %d, want 3", len(b))
	}
	if len(a) != 3 {
		t.Errorf("Butter: len(a) = %d, want 3", len(a))
	}

	// DC gain should be 1.
	sumB := 0.0
	sumA := 0.0
	for _, v := range b {
		sumB += v
	}
	for _, v := range a {
		sumA += v
	}
	dcGain := sumB / sumA
	if !approxEqual(dcGain, 1.0, 1e-10) {
		t.Errorf("Butter: DC gain = %v, want 1.0", dcGain)
	}
}

func TestButterInvalid(t *testing.T) {
	b, a := Butter(0, 10, 100)
	if b != nil || a != nil {
		t.Error("Butter: expected nil for order 0")
	}
	b, a = Butter(2, 50, 100) // cutoff = Nyquist
	if b != nil || a != nil {
		t.Error("Butter: expected nil for cutoff >= Nyquist")
	}
}

// ---------------------------------------------------------------------------
// LFilter Tests
// ---------------------------------------------------------------------------

func TestLFilter(t *testing.T) {
	// Simple FIR filter: moving average of 3.
	b := []float64{1.0 / 3, 1.0 / 3, 1.0 / 3}
	a := []float64{1}
	x := []float64{1, 2, 3, 4, 5}
	y := LFilter(b, a, x)
	if len(y) != len(x) {
		t.Fatalf("LFilter: length = %d, want %d", len(y), len(x))
	}
	// y[0] = x[0]/3 = 1/3
	// y[1] = (x[0]+x[1])/3 = 1
	// y[2] = (x[0]+x[1]+x[2])/3 = 2
	expected := []float64{1.0 / 3, 1.0, 2.0, 3.0, 4.0}
	for i := range expected {
		if !approxEqual(y[i], expected[i], 1e-10) {
			t.Errorf("LFilter[%d] = %v, want %v", i, y[i], expected[i])
		}
	}
}

func TestLFilterPassthrough(t *testing.T) {
	// Identity filter.
	b := []float64{1}
	a := []float64{1}
	x := []float64{1, 2, 3, 4}
	y := LFilter(b, a, x)
	for i := range x {
		if !approxEqual(y[i], x[i], 1e-14) {
			t.Errorf("LFilter passthrough[%d] = %v, want %v", i, y[i], x[i])
		}
	}
}

func TestLFilterEmpty(t *testing.T) {
	if LFilter(nil, []float64{1}, []float64{1}) != nil {
		t.Error("expected nil for empty b")
	}
}
