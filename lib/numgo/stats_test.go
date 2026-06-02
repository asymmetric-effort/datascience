//go:build unit

package numgo

import (
	"math"
	"testing"
)

func assertCloseTol(t *testing.T, label string, got, want float64, tolerance float64) {
	t.Helper()
	if math.Abs(got-want) > tolerance {
		t.Fatalf("%s: got %g, want %g (tol %g)", label, got, want, tolerance)
	}
}

// ---------------------------------------------------------------------------
// Min
// ---------------------------------------------------------------------------

func TestMinGlobal(t *testing.T) {
	a := FromSlice([]float64{3, 1, 4, 1, 5})
	got := Min(a)
	assertData(t, "MinGlobal", got, []float64{1})
}

func TestMinAxis(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{3, 1, 4, 1, 5, 9})
	got := Min(a, 0) // min along rows => [1, 1, 4]
	assertData(t, "MinAxis0", got, []float64{1, 1, 4})
}

func TestMinAxis1(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{3, 1, 4, 1, 5, 9})
	got := Min(a, 1) // min along cols => [1, 1]
	assertData(t, "MinAxis1", got, []float64{1, 1})
}

// ---------------------------------------------------------------------------
// Mean
// ---------------------------------------------------------------------------

func TestMeanGlobal(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	got := Mean(a)
	assertCloseTol(t, "MeanGlobal", got.Data()[0], 2.5, 1e-10)
}

func TestMeanAxis(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	got := Mean(a, 0) // mean along rows => [2, 3]
	assertCloseTol(t, "MeanAxis0[0]", got.Data()[0], 2.0, 1e-10)
	assertCloseTol(t, "MeanAxis0[1]", got.Data()[1], 3.0, 1e-10)
}

// ---------------------------------------------------------------------------
// Var & Std
// ---------------------------------------------------------------------------

func TestVarGlobal(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	got := Var(a)
	// population variance of [1,2,3,4] = 1.25
	assertCloseTol(t, "VarGlobal", got.Data()[0], 1.25, 1e-10)
}

func TestStdGlobal(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	got := Std(a)
	assertCloseTol(t, "StdGlobal", got.Data()[0], math.Sqrt(1.25), 1e-10)
}

func TestVarAxis(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 3, 5, 7})
	got := Var(a, 0)
	// col 0: var([1,5]) = 4.0, col 1: var([3,7]) = 4.0
	assertCloseTol(t, "VarAxis0[0]", got.Data()[0], 4.0, 1e-10)
	assertCloseTol(t, "VarAxis0[1]", got.Data()[1], 4.0, 1e-10)
}

// ---------------------------------------------------------------------------
// Cumsum & Cumprod
// ---------------------------------------------------------------------------

func TestCumsum1D(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	got := Cumsum(a, 0)
	assertData(t, "Cumsum1D", got, []float64{1, 3, 6, 10})
}

func TestCumsum2D(t *testing.T) {
	a := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	got := Cumsum(a, 1)
	assertData(t, "Cumsum2DAxis1", got, []float64{1, 3, 6, 4, 9, 15})
}

func TestCumprod1D(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	got := Cumprod(a, 0)
	assertData(t, "Cumprod1D", got, []float64{1, 2, 6, 24})
}

// ---------------------------------------------------------------------------
// Percentile, Quantile, Median
// ---------------------------------------------------------------------------

func TestPercentile(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4, 5})
	got := Percentile(a, 50)
	assertCloseTol(t, "Percentile50", got.Data()[0], 3.0, 1e-10)
}

func TestPercentile0And100(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4, 5})
	got0 := Percentile(a, 0)
	got100 := Percentile(a, 100)
	assertCloseTol(t, "Percentile0", got0.Data()[0], 1.0, 1e-10)
	assertCloseTol(t, "Percentile100", got100.Data()[0], 5.0, 1e-10)
}

func TestQuantile(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4, 5})
	got := Quantile(a, 0.25)
	assertCloseTol(t, "Quantile25", got.Data()[0], 2.0, 1e-10)
}

func TestMedian(t *testing.T) {
	a := FromSlice([]float64{1, 3, 2})
	got := Median(a)
	assertCloseTol(t, "Median", got.Data()[0], 2.0, 1e-10)
}

func TestMedianEven(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	got := Median(a)
	assertCloseTol(t, "MedianEven", got.Data()[0], 2.5, 1e-10)
}

// ---------------------------------------------------------------------------
// Average
// ---------------------------------------------------------------------------

func TestAverageUnweighted(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	got := Average(a, nil)
	assertCloseTol(t, "AverageNoWeights", got.Data()[0], 2.5, 1e-10)
}

func TestAverageWeighted(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4})
	w := FromSlice([]float64{4, 3, 2, 1})
	got := Average(a, w)
	// (1*4 + 2*3 + 3*2 + 4*1) / (4+3+2+1) = 20/10 = 2.0
	assertCloseTol(t, "AverageWeighted", got.Data()[0], 2.0, 1e-10)
}

func TestAverageWeightedAxis(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	w := FromSlice([]float64{1, 3})
	got := Average(a, w, 1)
	// row0: (1*1+2*3)/(1+3) = 7/4 = 1.75
	// row1: (3*1+4*3)/(1+3) = 15/4 = 3.75
	assertCloseTol(t, "AvgWeightedAxis[0]", got.Data()[0], 1.75, 1e-10)
	assertCloseTol(t, "AvgWeightedAxis[1]", got.Data()[1], 3.75, 1e-10)
}

// ---------------------------------------------------------------------------
// Nan* functions
// ---------------------------------------------------------------------------

func TestNanmean(t *testing.T) {
	a := FromSlice([]float64{1, math.NaN(), 3})
	got := Nanmean(a)
	assertCloseTol(t, "Nanmean", got.Data()[0], 2.0, 1e-10)
}

func TestNanstd(t *testing.T) {
	a := FromSlice([]float64{1, math.NaN(), 3})
	got := Nanstd(a)
	// std of [1,3] = 1.0
	assertCloseTol(t, "Nanstd", got.Data()[0], 1.0, 1e-10)
}

func TestNanvar(t *testing.T) {
	a := FromSlice([]float64{1, math.NaN(), 3})
	got := Nanvar(a)
	assertCloseTol(t, "Nanvar", got.Data()[0], 1.0, 1e-10)
}

func TestNanmin(t *testing.T) {
	a := FromSlice([]float64{3, math.NaN(), 1})
	got := Nanmin(a)
	assertCloseTol(t, "Nanmin", got.Data()[0], 1.0, 1e-10)
}

func TestNanmax(t *testing.T) {
	a := FromSlice([]float64{1, math.NaN(), 3})
	got := Nanmax(a)
	assertCloseTol(t, "Nanmax", got.Data()[0], 3.0, 1e-10)
}

func TestNansum(t *testing.T) {
	a := FromSlice([]float64{1, math.NaN(), 3})
	got := Nansum(a)
	assertCloseTol(t, "Nansum", got.Data()[0], 4.0, 1e-10)
}

func TestNanprod(t *testing.T) {
	a := FromSlice([]float64{2, math.NaN(), 3})
	got := Nanprod(a)
	assertCloseTol(t, "Nanprod", got.Data()[0], 6.0, 1e-10)
}

func TestNanminAllNaN(t *testing.T) {
	a := FromSlice([]float64{math.NaN(), math.NaN()})
	got := Nanmin(a)
	if !math.IsNaN(got.Data()[0]) {
		t.Fatalf("NanminAllNaN: expected NaN, got %g", got.Data()[0])
	}
}

// ---------------------------------------------------------------------------
// Histogram
// ---------------------------------------------------------------------------

func TestHistogram(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3, 4, 5})
	counts, edges := Histogram(a, 5)
	if counts.Size() != 5 {
		t.Fatalf("Histogram counts size: got %d, want 5", counts.Size())
	}
	if edges.Size() != 6 {
		t.Fatalf("Histogram edges size: got %d, want 6", edges.Size())
	}
	// Total counts should equal number of elements.
	total := 0.0
	for _, v := range counts.Data() {
		total += v
	}
	assertCloseTol(t, "HistogramTotal", total, 5.0, 1e-10)
}

// ---------------------------------------------------------------------------
// Bincount
// ---------------------------------------------------------------------------

func TestBincount(t *testing.T) {
	a := FromSlice([]float64{0, 1, 1, 3, 2, 1})
	got := Bincount(a)
	assertData(t, "Bincount", got, []float64{1, 3, 1, 1})
}

// ---------------------------------------------------------------------------
// Corrcoef
// ---------------------------------------------------------------------------

func TestCorrcoef(t *testing.T) {
	x := FromSlice([]float64{1, 2, 3, 4, 5})
	y := FromSlice([]float64{2, 4, 6, 8, 10})
	got, err := Corrcoef(x, y)
	if err != nil {
		t.Fatal(err)
	}
	assertShape(t, "Corrcoef", got, []int{2, 2})
	assertCloseTol(t, "Corrcoef[0,0]", got.Get(0, 0), 1.0, 1e-10)
	assertCloseTol(t, "Corrcoef[0,1]", got.Get(0, 1), 1.0, 1e-10)
}

// ---------------------------------------------------------------------------
// Cov
// ---------------------------------------------------------------------------

func TestCov1D(t *testing.T) {
	x := FromSlice([]float64{1, 2, 3, 4, 5})
	got, err := Cov(x)
	if err != nil {
		t.Fatal(err)
	}
	assertShape(t, "Cov1D", got, []int{1, 1})
	// sample variance of [1..5] = 2.5
	assertCloseTol(t, "Cov1DVal", got.Get(0, 0), 2.5, 1e-10)
}

func TestCov2D(t *testing.T) {
	x := NewNDArray([]int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	got, err := Cov(x)
	if err != nil {
		t.Fatal(err)
	}
	assertShape(t, "Cov2D", got, []int{2, 2})
	// Each row has sample variance 1.0
	assertCloseTol(t, "Cov2D[0,0]", got.Get(0, 0), 1.0, 1e-10)
	assertCloseTol(t, "Cov2D[1,1]", got.Get(1, 1), 1.0, 1e-10)
}

// ---------------------------------------------------------------------------
// Correlate
// ---------------------------------------------------------------------------

func TestCorrelate(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	v := FromSlice([]float64{0, 1, 0.5})
	got, err := Correlate(a, v)
	if err != nil {
		t.Fatal(err)
	}
	// Full correlation: length 3+3-1=5
	if got.Size() != 5 {
		t.Fatalf("Correlate: expected size 5, got %d", got.Size())
	}
}

// ---------------------------------------------------------------------------
// Convolve
// ---------------------------------------------------------------------------

func TestConvolve(t *testing.T) {
	a := FromSlice([]float64{1, 2, 3})
	v := FromSlice([]float64{0, 1, 0.5})
	got, err := Convolve(a, v)
	if err != nil {
		t.Fatal(err)
	}
	// [1*0, 1*1+2*0, 1*0.5+2*1+3*0, 2*0.5+3*1, 3*0.5]
	// = [0, 1, 2.5, 4, 1.5]
	want := []float64{0, 1, 2.5, 4, 1.5}
	d := got.Data()
	for i, w := range want {
		assertCloseTol(t, "Convolve", d[i], w, 1e-10)
	}
}

func TestConvolveNon1D(t *testing.T) {
	a := NewNDArray([]int{2, 2}, []float64{1, 2, 3, 4})
	v := FromSlice([]float64{1})
	_, err := Convolve(a, v)
	if err == nil {
		t.Fatal("expected error for non-1D input")
	}
}
