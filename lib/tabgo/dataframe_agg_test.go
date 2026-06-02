//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func makeNumericDF() *DataFrame {
	return NewDataFrameFromRows(
		[]string{"a", "b", "name"},
		[][]any{
			{1, 10.0, "x"},
			{2, 20.0, "y"},
			{3, 30.0, "z"},
			{4, 40.0, "w"},
		},
	)
}

func approxEqual(a, b, tol float64) bool {
	return math.Abs(a-b) < tol
}

func TestSum(t *testing.T) {
	df := makeNumericDF()
	s := df.Sum()
	if s["a"] != 10 {
		t.Fatalf("Sum a = %v, want 10", s["a"])
	}
	if s["b"] != 100 {
		t.Fatalf("Sum b = %v, want 100", s["b"])
	}
	if _, ok := s["name"]; ok {
		t.Fatal("Sum should not include non-numeric column 'name'")
	}
}

func TestMeanAll(t *testing.T) {
	df := makeNumericDF()
	m := df.MeanAll()
	if m["a"] != 2.5 {
		t.Fatalf("MeanAll a = %v, want 2.5", m["a"])
	}
	if m["b"] != 25.0 {
		t.Fatalf("MeanAll b = %v, want 25.0", m["b"])
	}
}

func TestMedianAll(t *testing.T) {
	df := makeNumericDF()
	m := df.MedianAll()
	// Median of [1,2,3,4] = (2+3)/2 = 2.5
	if m["a"] != 2.5 {
		t.Fatalf("MedianAll a = %v, want 2.5", m["a"])
	}
}

func TestMedianAllOdd(t *testing.T) {
	df := NewDataFrameFromRows([]string{"v"}, [][]any{{1}, {3}, {5}})
	m := df.MedianAll()
	if m["v"] != 3 {
		t.Fatalf("MedianAll v = %v, want 3", m["v"])
	}
}

func TestStdAll(t *testing.T) {
	df := makeNumericDF()
	s := df.StdAll()
	// Sample std of [1,2,3,4]: sqrt(((1-2.5)^2+(2-2.5)^2+(3-2.5)^2+(4-2.5)^2)/3) = sqrt(5/3)
	expected := math.Sqrt(5.0 / 3.0)
	if !approxEqual(s["a"], expected, 1e-10) {
		t.Fatalf("StdAll a = %v, want %v", s["a"], expected)
	}
}

func TestVarAll(t *testing.T) {
	df := makeNumericDF()
	v := df.VarAll()
	// Sample variance of [1,2,3,4] = 5/3
	expected := 5.0 / 3.0
	if !approxEqual(v["a"], expected, 1e-10) {
		t.Fatalf("VarAll a = %v, want %v", v["a"], expected)
	}
}

func TestMinAll(t *testing.T) {
	df := makeNumericDF()
	m := df.MinAll()
	if m["a"] != 1 {
		t.Fatalf("MinAll a = %v, want 1", m["a"])
	}
	if m["b"] != 10 {
		t.Fatalf("MinAll b = %v, want 10", m["b"])
	}
}

func TestMaxAll(t *testing.T) {
	df := makeNumericDF()
	m := df.MaxAll()
	if m["a"] != 4 {
		t.Fatalf("MaxAll a = %v, want 4", m["a"])
	}
	if m["b"] != 40 {
		t.Fatalf("MaxAll b = %v, want 40", m["b"])
	}
}

func TestCount(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1, nil},
			{nil, 2},
			{3, 4},
		},
	)
	c := df.Count()
	if c["a"] != 2 {
		t.Fatalf("Count a = %d, want 2", c["a"])
	}
	if c["b"] != 2 {
		t.Fatalf("Count b = %d, want 2", c["b"])
	}
}

func TestCountAllPresent(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}, {2}, {3}})
	c := df.Count()
	if c["a"] != 3 {
		t.Fatalf("Count a = %d, want 3", c["a"])
	}
}

func TestDescribe(t *testing.T) {
	df := makeNumericDF()
	desc := df.Describe()

	// Should have 8 rows (count, mean, std, min, 25%, 50%, 75%, max).
	if desc.Len() != 8 {
		t.Fatalf("Describe Len = %d, want 8", desc.Len())
	}

	// Check column names: "stat", then numeric columns.
	cols := desc.Columns()
	if cols[0] != "stat" {
		t.Fatalf("Describe first column = %q, want 'stat'", cols[0])
	}

	// Verify some specific values.
	statVals := desc.Column("stat").Values()
	if statVals[0] != "count" || statVals[1] != "mean" {
		t.Fatalf("Describe stat = %v", statVals)
	}

	// Check count for column "a" = 4.
	aVals := desc.Column("a").Values()
	if aVals[0] != 4.0 {
		t.Fatalf("Describe count a = %v, want 4.0", aVals[0])
	}
	// Check mean = 2.5.
	if aVals[1] != 2.5 {
		t.Fatalf("Describe mean a = %v, want 2.5", aVals[1])
	}
}

func TestDescribeNoNumericColumns(t *testing.T) {
	df := NewDataFrameFromRows([]string{"name"}, [][]any{{"a"}, {"b"}})
	desc := df.Describe()
	// Should have only the "stat" column.
	if len(desc.Columns()) != 1 {
		t.Fatalf("Describe with no numeric cols = %v", desc.Columns())
	}
}

func TestApply(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 10}, {2, 20}})
	result := df.Apply(func(s *Series) *Series {
		vals := s.Values()
		newVals := make([]any, len(vals))
		for i, v := range vals {
			newVals[i] = toFloat64(v) * 2
		}
		return NewSeries(s.Name(), newVals)
	})
	aVals := result.Column("a").Float64()
	if aVals[0] != 2 || aVals[1] != 4 {
		t.Fatalf("Apply a = %v, want [2 4]", aVals)
	}
	bVals := result.Column("b").Float64()
	if bVals[0] != 20 || bVals[1] != 40 {
		t.Fatalf("Apply b = %v, want [20 40]", bVals)
	}
}

func TestApplymap(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, "x"}, {2, "y"}})
	result := df.Applymap(func(v any) any {
		if n, ok := v.(int); ok {
			return n * 10
		}
		if s, ok := v.(string); ok {
			return s + "!"
		}
		return v
	})
	aVals := result.Column("a").Values()
	if aVals[0] != 10 || aVals[1] != 20 {
		t.Fatalf("Applymap a = %v", aVals)
	}
	bVals := result.Column("b").Values()
	if bVals[0] != "x!" || bVals[1] != "y!" {
		t.Fatalf("Applymap b = %v", bVals)
	}
}

func TestSumEmpty(t *testing.T) {
	df := NewDataFrame(map[string]*Series{})
	s := df.Sum()
	if len(s) != 0 {
		t.Fatalf("Sum of empty = %v", s)
	}
}

func TestAggWithNils(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}, {nil}, {3}})
	s := df.Sum()
	if s["a"] != 4 {
		t.Fatalf("Sum with nils a = %v, want 4", s["a"])
	}
	m := df.MeanAll()
	if m["a"] != 2 {
		t.Fatalf("MeanAll with nils a = %v, want 2", m["a"])
	}
}

func TestPercentile(t *testing.T) {
	// Test with known values: [1,2,3,4,5]
	vals := []float64{1, 2, 3, 4, 5}
	p25 := percentile(vals, 25)
	p50 := percentile(vals, 50)
	p75 := percentile(vals, 75)

	if p50 != 3.0 {
		t.Fatalf("percentile 50 = %v, want 3.0", p50)
	}
	if p25 != 2.0 {
		t.Fatalf("percentile 25 = %v, want 2.0", p25)
	}
	if p75 != 4.0 {
		t.Fatalf("percentile 75 = %v, want 4.0", p75)
	}
}

func TestPercentileEmpty(t *testing.T) {
	if percentile(nil, 50) != 0 {
		t.Fatal("percentile of empty should be 0")
	}
}

func TestVarianceSingleElement(t *testing.T) {
	if variance([]float64{5}) != 0 {
		t.Fatal("variance of single element should be 0")
	}
}
