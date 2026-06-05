//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func TestRollingCorrSymmetricDiagonalOne(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{
			{1.0, 2.0, 3.0},
			{2.0, 4.0, 5.0},
			{3.0, 6.0, 7.0},
			{4.0, 8.0, 9.0},
			{5.0, 10.0, 11.0},
		},
	)
	result := RollingCorr(df, 3)

	// First two should be nil
	if result[0] != nil || result[1] != nil {
		t.Error("RollingCorr: first 2 should be nil")
	}

	for row := 2; row < 5; row++ {
		corrDf := result[row]
		if corrDf == nil {
			t.Errorf("RollingCorr[%d] should not be nil", row)
			continue
		}
		cols := corrDf.Columns()
		nCols := len(cols)

		// Check diagonal = 1
		for ci, c := range cols {
			vals := corrDf.Column(c).Values()
			diag := toFloat64(vals[ci])
			if !almostEqual(diag, 1.0, 1e-10) {
				t.Errorf("RollingCorr[%d] diagonal(%s) = %v, want 1.0", row, c, diag)
			}
		}

		// Check symmetry
		for ci := 0; ci < nCols; ci++ {
			for cj := ci + 1; cj < nCols; cj++ {
				vij := toFloat64(corrDf.Column(cols[ci]).Values()[cj])
				vji := toFloat64(corrDf.Column(cols[cj]).Values()[ci])
				if !almostEqual(vij, vji, 1e-10) {
					t.Errorf("RollingCorr[%d] asymmetric: corr(%s,%s)=%v != corr(%s,%s)=%v",
						row, cols[ci], cols[cj], vij, cols[cj], cols[ci], vji)
				}
			}
		}
	}
}

func TestRollingCorrPerfectCorrelation(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 2.0},
			{2.0, 4.0},
			{3.0, 6.0},
			{4.0, 8.0},
		},
	)
	result := RollingCorr(df, 3)
	corrDf := result[2]

	// Perfect positive correlation between a and b
	cols := corrDf.Columns()
	var aCol, bCol string
	for _, c := range cols {
		if c == "a" {
			aCol = c
		}
		if c == "b" {
			bCol = c
		}
	}
	_ = bCol

	aVals := corrDf.Column(aCol).Values()
	bIdx := -1
	for i, c := range cols {
		if c == "b" {
			bIdx = i
		}
	}
	corrAB := toFloat64(aVals[bIdx])
	if !almostEqual(corrAB, 1.0, 1e-10) {
		t.Errorf("RollingCorr perfect correlation = %v, want 1.0", corrAB)
	}
}

func TestRollingCorrPair(t *testing.T) {
	s1 := NewSeries("a", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	s2 := NewSeries("b", []any{2.0, 4.0, 6.0, 8.0, 10.0})
	result := RollingCorrPair(s1, s2, 3)
	vals := result.Values()

	if vals[0] != nil || vals[1] != nil {
		t.Error("RollingCorrPair: first 2 should be nil")
	}

	// Perfect correlation
	for i := 2; i < 5; i++ {
		v := toFloat64(vals[i])
		if !almostEqual(v, 1.0, 1e-10) {
			t.Errorf("RollingCorrPair[%d] = %v, want 1.0", i, v)
		}
	}
}

func TestRollingCorrPairNegative(t *testing.T) {
	s1 := NewSeries("a", []any{1.0, 2.0, 3.0, 4.0})
	s2 := NewSeries("b", []any{10.0, 8.0, 6.0, 4.0})
	result := RollingCorrPair(s1, s2, 3)
	vals := result.Values()

	v := toFloat64(vals[2])
	if !almostEqual(v, -1.0, 1e-10) {
		t.Errorf("RollingCorrPair negative = %v, want -1.0", v)
	}
}

func TestRollingCorrPairUnequalLength(t *testing.T) {
	s1 := NewSeries("a", []any{1.0, 2.0, 3.0})
	s2 := NewSeries("b", []any{2.0, 4.0})
	result := RollingCorrPair(s1, s2, 2)
	if result.Len() != 2 {
		t.Errorf("RollingCorrPair unequal: got len=%d, want 2", result.Len())
	}
}

func TestRollingCorrPairSmallWindow(t *testing.T) {
	s1 := NewSeries("a", []any{1.0, 2.0})
	s2 := NewSeries("b", []any{3.0, 4.0})
	result := RollingCorrPair(s1, s2, 1)
	vals := result.Values()
	// Window of 1 has < 2 elements -> nil
	if vals[0] != nil {
		t.Errorf("RollingCorrPair window=1: expected nil, got %v", vals[0])
	}
}

func TestEWMCorr(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 2.0},
			{2.0, 4.0},
			{3.0, 6.0},
			{4.0, 8.0},
		},
	)
	result := EWMCorr(df, 3)

	if len(result) != 4 {
		t.Fatalf("EWMCorr: expected 4 results, got %d", len(result))
	}

	// First result should exist (identity-like)
	if result[0] == nil {
		t.Fatal("EWMCorr[0] should not be nil")
	}

	// Check diagonal is always 1
	for row := 0; row < 4; row++ {
		corrDf := result[row]
		cols := corrDf.Columns()
		for ci, c := range cols {
			vals := corrDf.Column(c).Values()
			diag := toFloat64(vals[ci])
			if !almostEqual(diag, 1.0, 1e-10) {
				t.Errorf("EWMCorr[%d] diagonal(%s) = %v, want 1.0", row, c, diag)
			}
		}
	}

	// For perfect correlation, off-diagonal should approach 1
	lastCorr := result[3]
	cols := lastCorr.Columns()
	aIdx, bIdx := -1, -1
	for i, c := range cols {
		if c == "a" {
			aIdx = i
		}
		if c == "b" {
			bIdx = i
		}
	}
	_ = aIdx
	corr := toFloat64(lastCorr.Column(cols[0]).Values()[bIdx])
	if corr < 0.9 {
		t.Errorf("EWMCorr perfect correlation = %v, want close to 1", corr)
	}
}

func TestEWMCorrSymmetry(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x", "y"},
		[][]any{
			{1.0, 5.0},
			{3.0, 2.0},
			{5.0, 8.0},
			{2.0, 1.0},
		},
	)
	result := EWMCorr(df, 2)

	for row := 1; row < 4; row++ {
		corrDf := result[row]
		cols := corrDf.Columns()
		xIdx, yIdx := -1, -1
		for i, c := range cols {
			if c == "x" {
				xIdx = i
			}
			if c == "y" {
				yIdx = i
			}
		}
		cxy := toFloat64(corrDf.Column(cols[xIdx]).Values()[yIdx])
		cyx := toFloat64(corrDf.Column(cols[yIdx]).Values()[xIdx])
		if !almostEqual(cxy, cyx, 1e-10) {
			t.Errorf("EWMCorr[%d] asymmetric: xy=%v, yx=%v", row, cxy, cyx)
		}
	}
}

func TestEWMCorrZeroVariance(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 5.0},
			{1.0, 5.0},
		},
	)
	result := EWMCorr(df, 2)
	// With zero variance, off-diagonal should be 0
	corrDf := result[1]
	cols := corrDf.Columns()
	bIdx := -1
	for i, c := range cols {
		if c == "b" {
			bIdx = i
		}
	}
	v := toFloat64(corrDf.Column(cols[0]).Values()[bIdx])
	if !almostEqual(v, 0.0, 1e-10) {
		t.Errorf("EWMCorr zero variance off-diagonal = %v, want 0", v)
	}
}

func TestPearsonCorrSmall(t *testing.T) {
	// Less than 2 elements
	if pearsonCorr([]float64{1.0}, []float64{2.0}) != 0 {
		t.Error("pearsonCorr single element should be 0")
	}
}

func TestPearsonCorrConstant(t *testing.T) {
	if pearsonCorr([]float64{5, 5, 5}, []float64{1, 2, 3}) != 0 {
		t.Error("pearsonCorr constant x should be 0")
	}
}

func TestRollingCorrBoundInRange(t *testing.T) {
	s1 := NewSeries("a", []any{1.0, 3.0, 2.0, 5.0, 4.0})
	s2 := NewSeries("b", []any{5.0, 1.0, 3.0, 2.0, 4.0})
	result := RollingCorrPair(s1, s2, 3)
	vals := result.Values()
	for i := 2; i < 5; i++ {
		v := toFloat64(vals[i])
		if v < -1.0-1e-10 || v > 1.0+1e-10 {
			t.Errorf("RollingCorrPair[%d] = %v, out of [-1, 1]", i, v)
		}
	}
}

func TestRollingCorrSingleColumn(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}},
	)
	result := RollingCorr(df, 2)
	corrDf := result[1]
	if corrDf == nil {
		t.Fatal("RollingCorr single col should not be nil")
	}
	v := toFloat64(corrDf.Column("x").Values()[0])
	if !almostEqual(v, 1.0, 1e-10) {
		t.Errorf("RollingCorr single col diagonal = %v, want 1", v)
	}
}

func TestRollingCorrPairNameFormat(t *testing.T) {
	s1 := NewSeries("alpha", []any{1.0, 2.0, 3.0})
	s2 := NewSeries("beta", []any{3.0, 2.0, 1.0})
	result := RollingCorrPair(s1, s2, 2)
	if result.Name() != "alpha_beta_corr" {
		t.Errorf("RollingCorrPair name = %q, want %q", result.Name(), "alpha_beta_corr")
	}
}

func TestEWMCorrSingleRow(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{{1.0, 2.0}},
	)
	result := EWMCorr(df, 3)
	if len(result) != 1 {
		t.Fatalf("EWMCorr single row: got %d results", len(result))
	}
	// First row: diagonal should be 1
	corrDf := result[0]
	cols := corrDf.Columns()
	for ci, c := range cols {
		diag := toFloat64(corrDf.Column(c).Values()[ci])
		if !almostEqual(diag, 1.0, 1e-10) {
			t.Errorf("EWMCorr[0] diagonal(%s) = %v, want 1", c, diag)
		}
	}
}

func TestRollingCorrPairMatchesManual(t *testing.T) {
	// Verify rolling corr matches manual Cov/Var computation
	s1 := NewSeries("a", []any{2.0, 4.0, 6.0, 8.0, 10.0})
	s2 := NewSeries("b", []any{1.0, 3.0, 2.0, 5.0, 4.0})

	corrResult := RollingCorrPair(s1, s2, 3)
	covResult := RollingCov(s1, s2, 3)

	for i := 2; i < 5; i++ {
		corr := toFloat64(corrResult.Values()[i])

		// Manual: corr = cov/(std1 * std2)
		w1 := []float64{toFloat64(s1.Values()[i-2]), toFloat64(s1.Values()[i-1]), toFloat64(s1.Values()[i])}
		w2 := []float64{toFloat64(s2.Values()[i-2]), toFloat64(s2.Values()[i-1]), toFloat64(s2.Values()[i])}
		cov := toFloat64(covResult.Values()[i])
		std1 := stddev(w1)
		std2 := stddev(w2)
		if std1 > 0 && std2 > 0 {
			manualCorr := cov / (std1 * std2)
			if !almostEqual(corr, manualCorr, 1e-10) {
				t.Errorf("RollingCorrPair[%d] = %v, manual = %v", i, corr, manualCorr)
			}
		}
	}
}

func TestRollingCorrLarger(t *testing.T) {
	// 3 columns, verify all properties hold
	df := NewDataFrameFromRows(
		[]string{"x", "y", "z"},
		[][]any{
			{1.0, 10.0, 5.0},
			{2.0, 8.0, 6.0},
			{3.0, 6.0, 7.0},
			{4.0, 4.0, 8.0},
			{5.0, 2.0, 9.0},
		},
	)
	result := RollingCorr(df, 3)
	for row := 2; row < 5; row++ {
		corrDf := result[row]
		cols := corrDf.Columns()

		// All correlations should be in [-1, 1]
		for _, c := range cols {
			vals := corrDf.Column(c).Values()
			for _, v := range vals {
				fv := toFloat64(v)
				if fv < -1.0-1e-10 || fv > 1.0+1e-10 {
					t.Errorf("RollingCorr[%d] corr out of bounds: %v", row, fv)
				}
			}
		}

		// x and y should be negatively correlated (x increasing, y decreasing)
		xIdx, yIdx := -1, -1
		for i, c := range cols {
			if c == "x" {
				xIdx = i
			}
			if c == "y" {
				yIdx = i
			}
		}
		cxy := toFloat64(corrDf.Column(cols[xIdx]).Values()[yIdx])
		if !almostEqual(cxy, -1.0, 1e-10) {
			t.Errorf("RollingCorr x,y perfect negative = %v, want -1", cxy)
		}
	}
}

func TestEWMCorrCorrelationBounds(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 5.0},
			{3.0, 2.0},
			{5.0, 8.0},
			{2.0, 1.0},
			{4.0, 6.0},
		},
	)
	result := EWMCorr(df, 3)
	for row := 0; row < 5; row++ {
		corrDf := result[row]
		cols := corrDf.Columns()
		for _, c := range cols {
			vals := corrDf.Column(c).Values()
			for _, v := range vals {
				fv := toFloat64(v)
				if fv < -1.0-1e-10 || fv > 1.0+1e-10 {
					t.Errorf("EWMCorr[%d] correlation out of bounds: %v", row, fv)
				}
			}
		}
	}
}

func TestRollingCorrPairIdenticalSeries(t *testing.T) {
	s := NewSeries("x", []any{1.0, 3.0, 2.0, 5.0, 4.0})
	result := RollingCorrPair(s, s, 3)
	vals := result.Values()
	for i := 2; i < 5; i++ {
		v := toFloat64(vals[i])
		if !almostEqual(v, 1.0, 1e-10) {
			t.Errorf("RollingCorrPair self[%d] = %v, want 1.0", i, v)
		}
	}
}

func TestRollingCorrPairZeroStd(t *testing.T) {
	s1 := NewSeries("a", []any{5.0, 5.0, 5.0})
	s2 := NewSeries("b", []any{1.0, 2.0, 3.0})
	result := RollingCorrPair(s1, s2, 3)
	vals := result.Values()
	// s1 has 0 std -> correlation should be 0
	v := toFloat64(vals[2])
	if v != 0 {
		t.Errorf("RollingCorrPair zero std = %v, want 0", v)
	}
}

// Dummy usage to avoid import errors
var _ = math.Abs
