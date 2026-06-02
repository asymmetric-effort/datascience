//go:build unit

package tabgo

import (
	"testing"
)

func TestIterRows(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1, 2},
			{3, 4},
		},
	)
	rows := IterRows(df)
	if len(rows) != 2 {
		t.Fatalf("IterRows len = %d, want 2", len(rows))
	}
	if rows[0]["a"] != 1 || rows[0]["b"] != 2 {
		t.Errorf("IterRows[0] = %v, want {a:1, b:2}", rows[0])
	}
	if rows[1]["a"] != 3 || rows[1]["b"] != 4 {
		t.Errorf("IterRows[1] = %v, want {a:3, b:4}", rows[1])
	}
}

func TestIterRowsEmpty(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, nil)
	rows := IterRows(df)
	if len(rows) != 0 {
		t.Errorf("IterRows empty len = %d, want 0", len(rows))
	}
}

func TestIterCols(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"x", "y"},
		[][]any{
			{1, 2},
			{3, 4},
		},
	)
	cols := IterCols(df)
	if len(cols) != 2 {
		t.Fatalf("IterCols len = %d, want 2", len(cols))
	}
	xSeries, ok := cols["x"]
	if !ok {
		t.Fatal("IterCols missing 'x'")
	}
	if xSeries.Len() != 2 {
		t.Errorf("IterCols x len = %d, want 2", xSeries.Len())
	}
	vals := xSeries.Values()
	if vals[0] != 1 || vals[1] != 3 {
		t.Errorf("IterCols x values = %v, want [1, 3]", vals)
	}
}
