//go:build unit

package tabgo

import (
	"reflect"
	"testing"
)

func TestIloc(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{
			{1, "x", 10},
			{2, "y", 20},
			{3, "z", 30},
		},
	)

	// Select rows 0, 2 and cols 0, 2.
	sub := df.Iloc([]int{0, 2}, []int{0, 2})
	if sub.Shape() != [2]int{2, 2} {
		t.Fatalf("Iloc Shape = %v, want [2 2]", sub.Shape())
	}
	if !reflect.DeepEqual(sub.Columns(), []string{"a", "c"}) {
		t.Fatalf("Iloc Columns = %v", sub.Columns())
	}
	if !reflect.DeepEqual(sub.Column("a").Values(), []any{1, 3}) {
		t.Fatalf("Iloc a = %v", sub.Column("a").Values())
	}
	if !reflect.DeepEqual(sub.Column("c").Values(), []any{10, 30}) {
		t.Fatalf("Iloc c = %v", sub.Column("c").Values())
	}
}

func TestIlocNilRows(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}, {2}, {3}})
	sub := df.Iloc(nil, []int{0})
	if sub.Len() != 3 {
		t.Fatalf("Iloc nil rows Len = %d, want 3", sub.Len())
	}
}

func TestIlocNilCols(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}})
	sub := df.Iloc([]int{0}, nil)
	if len(sub.Columns()) != 2 {
		t.Fatalf("Iloc nil cols ncols = %d, want 2", len(sub.Columns()))
	}
}

func TestIlocPanicsOnBadRow(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for out-of-range row")
		}
	}()
	df.Iloc([]int{5}, nil)
}

func TestIlocPanicsOnBadCol(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for out-of-range col")
		}
	}()
	df.Iloc(nil, []int{5})
}

func TestSample(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"v"},
		[][]any{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}},
	)
	s := df.Sample(3, 42)
	if s.Len() != 3 {
		t.Fatalf("Sample Len = %d, want 3", s.Len())
	}

	// Same seed should produce same result.
	s2 := df.Sample(3, 42)
	if !reflect.DeepEqual(s.Column("v").Values(), s2.Column("v").Values()) {
		t.Fatal("Sample with same seed should produce same result")
	}

	// Different seed should (very likely) produce different result.
	s3 := df.Sample(3, 99)
	// Not guaranteed but very likely different
	_ = s3
}

func TestSampleExceedsRows(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}, {2}})
	s := df.Sample(100, 1)
	if s.Len() != 2 {
		t.Fatalf("Sample exceeding rows Len = %d, want 2", s.Len())
	}
}

func TestSampleZero(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	s := df.Sample(0, 1)
	if s.Len() != 0 {
		t.Fatalf("Sample(0) Len = %d, want 0", s.Len())
	}
}

func TestNlargest(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"name", "score"},
		[][]any{
			{"Alice", 85},
			{"Bob", 92},
			{"Carol", 78},
			{"Dave", 95},
			{"Eve", 88},
		},
	)
	top := df.Nlargest(3, "score")
	if top.Len() != 3 {
		t.Fatalf("Nlargest Len = %d, want 3", top.Len())
	}
	scores := top.Column("score").Float64()
	// Should be descending.
	if scores[0] < scores[1] || scores[1] < scores[2] {
		t.Fatalf("Nlargest not in descending order: %v", scores)
	}
	if scores[0] != 95 {
		t.Fatalf("Nlargest top = %v, want 95", scores[0])
	}
}

func TestNlargestExceedsRows(t *testing.T) {
	df := NewDataFrameFromRows([]string{"v"}, [][]any{{1}, {2}})
	top := df.Nlargest(10, "v")
	if top.Len() != 2 {
		t.Fatalf("Nlargest exceeding rows Len = %d, want 2", top.Len())
	}
}

func TestNsmallest(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"name", "score"},
		[][]any{
			{"Alice", 85},
			{"Bob", 92},
			{"Carol", 78},
			{"Dave", 95},
			{"Eve", 88},
		},
	)
	bottom := df.Nsmallest(2, "score")
	if bottom.Len() != 2 {
		t.Fatalf("Nsmallest Len = %d, want 2", bottom.Len())
	}
	scores := bottom.Column("score").Float64()
	if scores[0] > scores[1] {
		t.Fatalf("Nsmallest not ascending: %v", scores)
	}
	if scores[0] != 78 {
		t.Fatalf("Nsmallest bottom = %v, want 78", scores[0])
	}
}

func TestWhere(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1, 10},
			{2, 20},
			{3, 30},
		},
	)
	result := df.Where(func(row map[string]any) bool {
		return toFloat64(row["a"]) > 1
	}, nil)

	if result.Len() != 3 {
		t.Fatalf("Where should preserve all rows, got %d", result.Len())
	}
	// Row 0 should have nil values (condition false).
	aVals := result.Column("a").Values()
	if aVals[0] != nil {
		t.Fatalf("Where: row 0 col a = %v, want nil", aVals[0])
	}
	// Rows 1,2 should have original values.
	if aVals[1] != 2 || aVals[2] != 3 {
		t.Fatalf("Where: rows 1,2 a = %v, %v", aVals[1], aVals[2])
	}
}

func TestWhereWithOtherValue(t *testing.T) {
	df := NewDataFrameFromRows([]string{"v"}, [][]any{{1}, {2}, {3}})
	result := df.Where(func(row map[string]any) bool {
		return toFloat64(row["v"]) >= 2
	}, -1)

	vals := result.Column("v").Values()
	if vals[0] != -1 {
		t.Fatalf("Where: row 0 = %v, want -1", vals[0])
	}
	if vals[1] != 2 || vals[2] != 3 {
		t.Fatalf("Where: rows 1,2 = %v, %v", vals[1], vals[2])
	}
}
