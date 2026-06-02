//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func makeGroupByTestDF() *DataFrame {
	return NewDataFrameFromRows(
		[]string{"group", "value"},
		[][]any{
			{"a", 1.0},
			{"a", 2.0},
			{"a", 3.0},
			{"b", 4.0},
			{"b", 6.0},
		},
	)
}

func TestGroupByStd(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Std("value")
	if result.Len() != 2 {
		t.Fatalf("Std rows = %d, want 2", result.Len())
	}
	vals := result.Column("value").Values()
	// group a: [1,2,3], std = 1.0
	if !almostEqual(toFloat64(vals[0]), 1.0, 0.01) {
		t.Errorf("Std group a = %v, want 1.0", vals[0])
	}
}

func TestGroupByVar(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Var("value")
	vals := result.Column("value").Values()
	// group a: var = 1.0
	if !almostEqual(toFloat64(vals[0]), 1.0, 0.01) {
		t.Errorf("Var group a = %v, want 1.0", vals[0])
	}
	// group b: [4,6], var = 2.0
	if !almostEqual(toFloat64(vals[1]), 2.0, 0.01) {
		t.Errorf("Var group b = %v, want 2.0", vals[1])
	}
}

func TestGroupByMin(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Min("value")
	vals := result.Column("value").Values()
	if toFloat64(vals[0]) != 1.0 {
		t.Errorf("Min group a = %v, want 1.0", vals[0])
	}
	if toFloat64(vals[1]) != 4.0 {
		t.Errorf("Min group b = %v, want 4.0", vals[1])
	}
}

func TestGroupByMax(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Max("value")
	vals := result.Column("value").Values()
	if toFloat64(vals[0]) != 3.0 {
		t.Errorf("Max group a = %v, want 3.0", vals[0])
	}
	if toFloat64(vals[1]) != 6.0 {
		t.Errorf("Max group b = %v, want 6.0", vals[1])
	}
}

func TestGroupByMedian(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Median("value")
	vals := result.Column("value").Values()
	// group a: [1,2,3], median = 2
	if toFloat64(vals[0]) != 2.0 {
		t.Errorf("Median group a = %v, want 2.0", vals[0])
	}
	// group b: [4,6], median = 5
	if toFloat64(vals[1]) != 5.0 {
		t.Errorf("Median group b = %v, want 5.0", vals[1])
	}
}

func TestGroupByFirst(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").First()
	if result.Len() != 2 {
		t.Fatalf("First rows = %d, want 2", result.Len())
	}
	vals := result.Column("value").Values()
	if toFloat64(vals[0]) != 1.0 {
		t.Errorf("First group a value = %v, want 1.0", vals[0])
	}
	if toFloat64(vals[1]) != 4.0 {
		t.Errorf("First group b value = %v, want 4.0", vals[1])
	}
}

func TestGroupByLast(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Last()
	vals := result.Column("value").Values()
	if toFloat64(vals[0]) != 3.0 {
		t.Errorf("Last group a value = %v, want 3.0", vals[0])
	}
	if toFloat64(vals[1]) != 6.0 {
		t.Errorf("Last group b value = %v, want 6.0", vals[1])
	}
}

func TestGroupBySize(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Size()
	vals := result.Column("size").Values()
	if toInt(vals[0]) != 3 {
		t.Errorf("Size group a = %v, want 3", vals[0])
	}
	if toInt(vals[1]) != 2 {
		t.Errorf("Size group b = %v, want 2", vals[1])
	}
}

func TestGroupByNgroups(t *testing.T) {
	df := makeGroupByTestDF()
	if n := df.GroupBy("group").Ngroups(); n != 2 {
		t.Errorf("Ngroups = %d, want 2", n)
	}
}

func TestGroupByFilter(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Filter(func(sub *DataFrame) bool {
		return sub.Len() >= 3
	})
	if result.Len() != 3 {
		t.Errorf("Filter rows = %d, want 3 (only group a)", result.Len())
	}
}

func TestGroupByFilterNone(t *testing.T) {
	df := makeGroupByTestDF()
	result := df.GroupBy("group").Filter(func(sub *DataFrame) bool {
		return false
	})
	if result.Len() != 0 {
		t.Errorf("Filter none rows = %d, want 0", result.Len())
	}
}

func TestGroupByStdSingleElement(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"g", "v"},
		[][]any{{"a", 5.0}},
	)
	result := df.GroupBy("g").Std("v")
	vals := result.Column("v").Values()
	if toFloat64(vals[0]) != 0.0 {
		t.Errorf("Std single = %v, want 0.0", vals[0])
	}
}

// Ensure unused import is used
var _ = math.Sqrt
