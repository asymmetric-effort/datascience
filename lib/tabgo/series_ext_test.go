//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func almostEqual(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

func TestSeriesSum(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0})
	if got := s.Sum(); got != 10.0 {
		t.Errorf("Sum() = %v, want 10.0", got)
	}
}

func TestSeriesSumWithNil(t *testing.T) {
	s := NewSeries("x", []any{1.0, nil, 3.0})
	if got := s.Sum(); got != 4.0 {
		t.Errorf("Sum() = %v, want 4.0", got)
	}
}

func TestSeriesMean(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0})
	if got := s.Mean(); got != 2.5 {
		t.Errorf("Mean() = %v, want 2.5", got)
	}
}

func TestSeriesMeanEmpty(t *testing.T) {
	s := NewSeries("x", []any{})
	if got := s.Mean(); got != 0 {
		t.Errorf("Mean() = %v, want 0", got)
	}
}

func TestSeriesStd(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	got := s.Std()
	// sample std dev of [1,2,3,4,5] = sqrt(2.5) ~ 1.5811
	if !almostEqual(got, 1.5811, 0.001) {
		t.Errorf("Std() = %v, want ~1.5811", got)
	}
}

func TestSeriesVar(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	got := s.Var()
	// sample variance of [1,2,3,4,5] = 10/4 = 2.5
	if !almostEqual(got, 2.5, 0.001) {
		t.Errorf("Var() = %v, want 2.5", got)
	}
}

func TestSeriesMin(t *testing.T) {
	s := NewSeries("x", []any{3.0, 1.0, 4.0, 1.5})
	if got := s.Min(); got != 1.0 {
		t.Errorf("Min() = %v, want 1.0", got)
	}
}

func TestSeriesMinEmpty(t *testing.T) {
	s := NewSeries("x", []any{})
	if got := s.Min(); got != 0 {
		t.Errorf("Min() = %v, want 0", got)
	}
}

func TestSeriesMax(t *testing.T) {
	s := NewSeries("x", []any{3.0, 1.0, 4.0, 1.5})
	if got := s.Max(); got != 4.0 {
		t.Errorf("Max() = %v, want 4.0", got)
	}
}

func TestSeriesMedian(t *testing.T) {
	s := NewSeries("x", []any{1.0, 3.0, 2.0})
	if got := s.Median(); got != 2.0 {
		t.Errorf("Median() = %v, want 2.0", got)
	}
}

func TestSeriesMedianEven(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0})
	if got := s.Median(); got != 2.5 {
		t.Errorf("Median() = %v, want 2.5", got)
	}
}

func TestSeriesCount(t *testing.T) {
	s := NewSeries("x", []any{1, nil, 3, nil, 5})
	if got := s.Count(); got != 3 {
		t.Errorf("Count() = %v, want 3", got)
	}
}

func TestSeriesDescribe(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	desc := s.Describe()
	if desc["count"] != 5 {
		t.Errorf("Describe count = %v, want 5", desc["count"])
	}
	if desc["mean"] != 3.0 {
		t.Errorf("Describe mean = %v, want 3", desc["mean"])
	}
	if desc["min"] != 1.0 {
		t.Errorf("Describe min = %v, want 1", desc["min"])
	}
	if desc["max"] != 5.0 {
		t.Errorf("Describe max = %v, want 5", desc["max"])
	}
	if desc["50%"] != 3.0 {
		t.Errorf("Describe 50%% = %v, want 3", desc["50%"])
	}
}

func TestSeriesApply(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0})
	result := s.Apply(func(v any) any {
		return toFloat64(v) * 2
	})
	vals := result.Float64()
	expected := []float64{2.0, 4.0, 6.0}
	for i, v := range vals {
		if v != expected[i] {
			t.Errorf("Apply[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestSeriesMap(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0})
	result := s.Map(func(v any) any {
		return toFloat64(v) + 10
	})
	vals := result.Float64()
	if vals[0] != 11.0 || vals[1] != 12.0 {
		t.Errorf("Map() = %v, want [11, 12]", vals)
	}
}

func TestSeriesReplace(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 1, 3})
	result := s.Replace(1, 99)
	vals := result.Values()
	if vals[0] != 99 || vals[1] != 2 || vals[2] != 99 || vals[3] != 3 {
		t.Errorf("Replace() = %v, want [99 2 99 3]", vals)
	}
}

func TestSeriesClip(t *testing.T) {
	s := NewSeries("x", []any{-1.0, 2.0, 5.0, 10.0})
	result := s.Clip(0, 6)
	vals := result.Float64()
	expected := []float64{0, 2, 5, 6}
	for i, v := range vals {
		if v != expected[i] {
			t.Errorf("Clip[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestSeriesClipNil(t *testing.T) {
	s := NewSeries("x", []any{nil, 5.0})
	result := s.Clip(0, 10)
	vals := result.Values()
	if vals[0] != nil {
		t.Errorf("Clip nil should stay nil, got %v", vals[0])
	}
}

func TestSeriesAbs(t *testing.T) {
	s := NewSeries("x", []any{-3.0, 2.0, -1.0})
	result := s.Abs()
	vals := result.Float64()
	expected := []float64{3.0, 2.0, 1.0}
	for i, v := range vals {
		if v != expected[i] {
			t.Errorf("Abs[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestSeriesRound(t *testing.T) {
	s := NewSeries("x", []any{1.234, 2.567, 3.891})
	result := s.Round(1)
	vals := result.Float64()
	expected := []float64{1.2, 2.6, 3.9}
	for i, v := range vals {
		if !almostEqual(v, expected[i], 0.001) {
			t.Errorf("Round[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestSeriesSort(t *testing.T) {
	s := NewSeries("x", []any{3.0, 1.0, 2.0})
	asc := s.Sort(true)
	vals := asc.Float64()
	if vals[0] != 1.0 || vals[1] != 2.0 || vals[2] != 3.0 {
		t.Errorf("Sort(asc) = %v, want [1 2 3]", vals)
	}

	desc := s.Sort(false)
	vals = desc.Float64()
	if vals[0] != 3.0 || vals[1] != 2.0 || vals[2] != 1.0 {
		t.Errorf("Sort(desc) = %v, want [3 2 1]", vals)
	}
}

func TestSeriesRank(t *testing.T) {
	s := NewSeries("x", []any{3.0, 1.0, 2.0, 1.0})
	result := s.Rank()
	vals := result.Values()
	// 1.0 appears at indices 1,3 -> ranks 1,2 -> avg 1.5
	// 2.0 at index 2 -> rank 3
	// 3.0 at index 0 -> rank 4
	if toFloat64(vals[0]) != 4.0 {
		t.Errorf("Rank[0] = %v, want 4.0", vals[0])
	}
	if toFloat64(vals[1]) != 1.5 {
		t.Errorf("Rank[1] = %v, want 1.5", vals[1])
	}
	if toFloat64(vals[2]) != 3.0 {
		t.Errorf("Rank[2] = %v, want 3.0", vals[2])
	}
	if toFloat64(vals[3]) != 1.5 {
		t.Errorf("Rank[3] = %v, want 1.5", vals[3])
	}
}

func TestSeriesRankWithNil(t *testing.T) {
	s := NewSeries("x", []any{3.0, nil, 1.0})
	result := s.Rank()
	vals := result.Values()
	if vals[1] != nil {
		t.Errorf("Rank nil should be nil, got %v", vals[1])
	}
}

func TestSeriesIsin(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3, 4})
	result := s.Isin([]any{2, 4})
	expected := []bool{false, true, false, true}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Isin[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestSeriesBetween(t *testing.T) {
	s := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	result := s.Between(2.0, 4.0)
	expected := []bool{false, true, true, true, false}
	for i, v := range result {
		if v != expected[i] {
			t.Errorf("Between[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestSeriesBetweenNil(t *testing.T) {
	s := NewSeries("x", []any{nil, 3.0})
	result := s.Between(1.0, 5.0)
	if result[0] != false {
		t.Error("Between nil should be false")
	}
	if result[1] != true {
		t.Error("Between 3.0 in [1,5] should be true")
	}
}

func TestSeriesCorr(t *testing.T) {
	x := NewSeries("x", []any{1.0, 2.0, 3.0, 4.0, 5.0})
	y := NewSeries("y", []any{2.0, 4.0, 6.0, 8.0, 10.0})
	corr := x.Corr(y)
	if !almostEqual(corr, 1.0, 0.001) {
		t.Errorf("Corr() = %v, want 1.0 (perfect positive)", corr)
	}
}

func TestSeriesCorrNegative(t *testing.T) {
	x := NewSeries("x", []any{1.0, 2.0, 3.0})
	y := NewSeries("y", []any{3.0, 2.0, 1.0})
	corr := x.Corr(y)
	if !almostEqual(corr, -1.0, 0.001) {
		t.Errorf("Corr() = %v, want -1.0", corr)
	}
}

func TestSeriesCorrMismatchLen(t *testing.T) {
	x := NewSeries("x", []any{1.0, 2.0})
	y := NewSeries("y", []any{1.0})
	if x.Corr(y) != 0 {
		t.Error("Corr mismatched lengths should return 0")
	}
}

func TestSeriesDropNA(t *testing.T) {
	s := NewSeries("x", []any{1, nil, 3, nil})
	result := s.DropNA()
	if result.Len() != 2 {
		t.Errorf("DropNA Len = %v, want 2", result.Len())
	}
	vals := result.Values()
	if vals[0] != 1 || vals[1] != 3 {
		t.Errorf("DropNA values = %v, want [1, 3]", vals)
	}
}

func TestSeriesFillNA(t *testing.T) {
	s := NewSeries("x", []any{1, nil, 3, nil})
	result := s.FillNA(0)
	vals := result.Values()
	if vals[1] != 0 || vals[3] != 0 {
		t.Errorf("FillNA values = %v, want [1, 0, 3, 0]", vals)
	}
}

func TestSeriesSumInts(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3})
	if got := s.Sum(); got != 6.0 {
		t.Errorf("Sum() = %v, want 6.0", got)
	}
}

func TestSeriesVarSingleElement(t *testing.T) {
	s := NewSeries("x", []any{5.0})
	if got := s.Var(); got != 0 {
		t.Errorf("Var() single element = %v, want 0", got)
	}
}

func TestSeriesAbsNil(t *testing.T) {
	s := NewSeries("x", []any{nil, -5.0})
	result := s.Abs()
	vals := result.Values()
	if vals[0] != nil {
		t.Error("Abs nil should stay nil")
	}
	if toFloat64(vals[1]) != 5.0 {
		t.Errorf("Abs(-5) = %v, want 5", vals[1])
	}
}

func TestSeriesRoundNil(t *testing.T) {
	s := NewSeries("x", []any{nil, 1.555})
	result := s.Round(2)
	vals := result.Values()
	if vals[0] != nil {
		t.Error("Round nil should stay nil")
	}
	if !almostEqual(toFloat64(vals[1]), 1.56, 0.001) {
		t.Errorf("Round(1.555, 2) = %v, want 1.56", vals[1])
	}
}
