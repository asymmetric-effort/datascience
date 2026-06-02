//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func TestCorr(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 2.0},
			{2.0, 4.0},
			{3.0, 6.0},
			{4.0, 8.0},
		},
	)
	result, err := Corr(df)
	if err != nil {
		t.Fatalf("Corr error: %v", err)
	}
	// a-a correlation should be 1
	aVals := result.Column("a").Values()
	if !almostEqual(toFloat64(aVals[0]), 1.0, 0.001) {
		t.Errorf("Corr a-a = %v, want 1.0", aVals[0])
	}
	// a-b correlation should be 1 (perfectly correlated)
	bVals := result.Column("b").Values()
	if !almostEqual(toFloat64(bVals[0]), 1.0, 0.001) {
		t.Errorf("Corr a-b = %v, want 1.0", bVals[0])
	}
}

func TestCorrNoNumericColumns(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"name"},
		[][]any{{"a"}, {"b"}},
	)
	_, err := Corr(df)
	if err == nil {
		t.Error("Corr should fail with no numeric columns")
	}
}

func TestCov(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1.0, 2.0},
			{2.0, 4.0},
			{3.0, 6.0},
		},
	)
	result, err := Cov(df)
	if err != nil {
		t.Fatalf("Cov error: %v", err)
	}
	// a variance = 1.0
	aVals := result.Column("a").Values()
	if !almostEqual(toFloat64(aVals[0]), 1.0, 0.001) {
		t.Errorf("Cov a-a = %v, want 1.0", aVals[0])
	}
	// cov(a,b) = 2.0
	bVals := result.Column("b").Values()
	if !almostEqual(toFloat64(bVals[0]), 2.0, 0.001) {
		t.Errorf("Cov a-b = %v, want 2.0", bVals[0])
	}
}

func TestCumsum(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}},
	)
	result := Cumsum(df)
	vals := result.Column("x").Float64()
	expected := []float64{1.0, 3.0, 6.0}
	for i, v := range vals {
		if v != expected[i] {
			t.Errorf("Cumsum[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestCumprod(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {3.0}},
	)
	result := Cumprod(df)
	vals := result.Column("x").Float64()
	expected := []float64{1.0, 2.0, 6.0}
	for i, v := range vals {
		if v != expected[i] {
			t.Errorf("Cumprod[%d] = %v, want %v", i, v, expected[i])
		}
	}
}

func TestDiff(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {3.0}, {6.0}, {10.0}},
	)
	result := Diff(df, 1)
	vals := result.Column("x").Values()
	if vals[0] != nil {
		t.Errorf("Diff[0] = %v, want nil", vals[0])
	}
	if toFloat64(vals[1]) != 2.0 {
		t.Errorf("Diff[1] = %v, want 2.0", vals[1])
	}
	if toFloat64(vals[2]) != 3.0 {
		t.Errorf("Diff[2] = %v, want 3.0", vals[2])
	}
	if toFloat64(vals[3]) != 4.0 {
		t.Errorf("Diff[3] = %v, want 4.0", vals[3])
	}
}

func TestDiffPeriods2(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {5.0}},
	)
	result := Diff(df, 2)
	vals := result.Column("x").Values()
	if vals[0] != nil || vals[1] != nil {
		t.Error("Diff periods=2, first 2 should be nil")
	}
	if toFloat64(vals[2]) != 4.0 {
		t.Errorf("Diff[2] = %v, want 4.0", vals[2])
	}
}

func TestPctChange(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{100.0}, {110.0}, {121.0}},
	)
	result := PctChange(df, 1)
	vals := result.Column("x").Values()
	if vals[0] != nil {
		t.Errorf("PctChange[0] = %v, want nil", vals[0])
	}
	if !almostEqual(toFloat64(vals[1]), 0.1, 0.001) {
		t.Errorf("PctChange[1] = %v, want 0.1", vals[1])
	}
	if !almostEqual(toFloat64(vals[2]), 0.1, 0.001) {
		t.Errorf("PctChange[2] = %v, want 0.1", vals[2])
	}
}

func TestPctChangeZeroDenom(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{0.0}, {5.0}},
	)
	result := PctChange(df, 1)
	vals := result.Column("x").Values()
	if vals[1] != nil {
		t.Errorf("PctChange with zero denom = %v, want nil", vals[1])
	}
}

func TestRankDF(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{3.0}, {1.0}, {2.0}},
	)
	result := Rank(df)
	vals := result.Column("x").Values()
	// 1.0->rank 1, 2.0->rank 2, 3.0->rank 3
	if toFloat64(vals[0]) != 3.0 {
		t.Errorf("Rank[0] = %v, want 3.0", vals[0])
	}
	if toFloat64(vals[1]) != 1.0 {
		t.Errorf("Rank[1] = %v, want 1.0", vals[1])
	}
	if toFloat64(vals[2]) != 2.0 {
		t.Errorf("Rank[2] = %v, want 2.0", vals[2])
	}
}

func TestRankDFTies(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{1.0}, {2.0}, {2.0}, {4.0}},
	)
	result := Rank(df)
	vals := result.Column("x").Values()
	// 2.0 ties at positions 2,3 -> avg rank 2.5
	if toFloat64(vals[1]) != 2.5 {
		t.Errorf("Rank tied[1] = %v, want 2.5", vals[1])
	}
	if toFloat64(vals[2]) != 2.5 {
		t.Errorf("Rank tied[2] = %v, want 2.5", vals[2])
	}
}

func TestRankDFWithNil(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x"},
		[][]any{{3.0}, {nil}, {1.0}},
	)
	result := Rank(df)
	vals := result.Column("x").Values()
	if vals[1] != nil {
		t.Errorf("Rank nil = %v, want nil", vals[1])
	}
}

// Ensure unused import is used
var _ = math.Abs
