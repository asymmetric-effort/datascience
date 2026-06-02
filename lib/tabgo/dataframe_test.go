//go:build unit

package tabgo

import (
	"reflect"
	"testing"
)

func TestNewDataFrame(t *testing.T) {
	df := NewDataFrame(map[string]*Series{
		"a": NewSeries("a", []any{1, 2}),
		"b": NewSeries("b", []any{"x", "y"}),
	})
	if df.Len() != 2 {
		t.Fatalf("Len = %d, want 2", df.Len())
	}
	cols := df.Columns()
	// alphabetical order
	if !reflect.DeepEqual(cols, []string{"a", "b"}) {
		t.Fatalf("Columns = %v", cols)
	}
}

func TestNewDataFrameFromRows(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x", "y"},
		[][]any{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		},
	)
	if df.Len() != 3 {
		t.Fatalf("Len = %d, want 3", df.Len())
	}
	xVals := df.Column("x").Values()
	if !reflect.DeepEqual(xVals, []any{1, 2, 3}) {
		t.Fatalf("x values = %v", xVals)
	}
}

func TestNewDataFrameFromRowsShortRow(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{
			{1, 2}, // missing third column
		},
	)
	if df.Column("c").Values()[0] != nil {
		t.Fatal("short row should fill with nil")
	}
}

func TestColumnPanics(t *testing.T) {
	df := NewDataFrame(map[string]*Series{
		"a": NewSeries("a", []any{1}),
	})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for missing column")
		}
	}()
	df.Column("nonexistent")
}

func TestSelect(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b", "c"},
		[][]any{{1, 2, 3}},
	)
	sub := df.Select("c", "a")
	if !reflect.DeepEqual(sub.Columns(), []string{"c", "a"}) {
		t.Fatalf("Select columns = %v", sub.Columns())
	}
	if sub.Len() != 1 {
		t.Fatalf("Select Len = %d", sub.Len())
	}
}

func TestFilter(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"name", "age"},
		[][]any{
			{"Alice", 30},
			{"Bob", 20},
			{"Carol", 25},
		},
	)
	filtered := df.Filter(func(row map[string]any) bool {
		return row["age"].(int) >= 25
	})
	if filtered.Len() != 2 {
		t.Fatalf("Filter Len = %d, want 2", filtered.Len())
	}
	names := filtered.Column("name").Values()
	if !reflect.DeepEqual(names, []any{"Alice", "Carol"}) {
		t.Fatalf("filtered names = %v", names)
	}
}

func TestFilterEmpty(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}, {2}})
	empty := df.Filter(func(row map[string]any) bool { return false })
	if empty.Len() != 0 {
		t.Fatalf("expected 0 rows, got %d", empty.Len())
	}
}

func TestHead(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"v"},
		[][]any{{1}, {2}, {3}, {4}, {5}},
	)
	h := df.Head(3)
	if h.Len() != 3 {
		t.Fatalf("Head(3) Len = %d", h.Len())
	}
	got := h.Column("v").Values()
	if !reflect.DeepEqual(got, []any{1, 2, 3}) {
		t.Fatalf("Head values = %v", got)
	}

	// n > len
	h2 := df.Head(100)
	if h2.Len() != 5 {
		t.Fatalf("Head(100) Len = %d", h2.Len())
	}
}

func TestTail(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"v"},
		[][]any{{1}, {2}, {3}, {4}, {5}},
	)
	tl := df.Tail(2)
	if tl.Len() != 2 {
		t.Fatalf("Tail(2) Len = %d", tl.Len())
	}
	got := tl.Column("v").Values()
	if !reflect.DeepEqual(got, []any{4, 5}) {
		t.Fatalf("Tail values = %v", got)
	}
}

func TestCopy(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a"},
		[][]any{{1}, {2}},
	)
	cp := df.Copy()
	if cp.Len() != df.Len() {
		t.Fatal("Copy Len mismatch")
	}
	// verify independence
	if &df.columns[0] == &cp.columns[0] {
		t.Fatal("Copy shares column slice")
	}
}

func TestEmptyDataFrame(t *testing.T) {
	df := NewDataFrame(map[string]*Series{})
	if df.Len() != 0 {
		t.Fatalf("empty df Len = %d", df.Len())
	}
	if len(df.Columns()) != 0 {
		t.Fatal("empty df should have no columns")
	}
}
