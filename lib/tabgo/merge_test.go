//go:build unit

package tabgo

import (
	"reflect"
	"testing"
)

func TestMergeInnerSingleColumn(t *testing.T) {
	left := NewDataFrameFromRows(
		[]string{"id", "lval"},
		[][]any{
			{1, "a"},
			{2, "b"},
			{3, "c"},
		},
	)
	right := NewDataFrameFromRows(
		[]string{"id", "rval"},
		[][]any{
			{2, "x"},
			{3, "y"},
			{4, "z"},
		},
	)
	result, err := Merge(left, right, []string{"id"}, "inner")
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	if result.Len() != 2 {
		t.Fatalf("expected 2 rows, got %d", result.Len())
	}

	idVals := result.Column("id").Values()
	lVals := result.Column("lval").Values()
	rVals := result.Column("rval").Values()

	// Row order follows left table order
	if !reflect.DeepEqual(idVals, []any{2, 3}) {
		t.Fatalf("id values = %v, want [2 3]", idVals)
	}
	if !reflect.DeepEqual(lVals, []any{"b", "c"}) {
		t.Fatalf("lval values = %v, want [b c]", lVals)
	}
	if !reflect.DeepEqual(rVals, []any{"x", "y"}) {
		t.Fatalf("rval values = %v, want [x y]", rVals)
	}
}

func TestMergeInnerMultiColumn(t *testing.T) {
	left := NewDataFrameFromRows(
		[]string{"k1", "k2", "v1"},
		[][]any{
			{"a", 1, "L1"},
			{"a", 2, "L2"},
			{"b", 1, "L3"},
		},
	)
	right := NewDataFrameFromRows(
		[]string{"k1", "k2", "v2"},
		[][]any{
			{"a", 1, "R1"},
			{"b", 2, "R2"},
			{"a", 2, "R3"},
		},
	)
	result, err := Merge(left, right, []string{"k1", "k2"}, "inner")
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	if result.Len() != 2 {
		t.Fatalf("expected 2 rows, got %d", result.Len())
	}
	v1Vals := result.Column("v1").Values()
	v2Vals := result.Column("v2").Values()
	if !reflect.DeepEqual(v1Vals, []any{"L1", "L2"}) {
		t.Fatalf("v1 values = %v, want [L1 L2]", v1Vals)
	}
	if !reflect.DeepEqual(v2Vals, []any{"R1", "R3"}) {
		t.Fatalf("v2 values = %v, want [R1 R3]", v2Vals)
	}
}

func TestMergeInnerDuplicateKeys(t *testing.T) {
	left := NewDataFrameFromRows(
		[]string{"id", "lv"},
		[][]any{
			{1, "a"},
			{1, "b"},
		},
	)
	right := NewDataFrameFromRows(
		[]string{"id", "rv"},
		[][]any{
			{1, "x"},
			{1, "y"},
		},
	)
	result, err := Merge(left, right, []string{"id"}, "inner")
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	// Cartesian product: 2 * 2 = 4
	if result.Len() != 4 {
		t.Fatalf("expected 4 rows, got %d", result.Len())
	}
}

func TestMergeNoOverlap(t *testing.T) {
	left := NewDataFrameFromRows(
		[]string{"id", "v"},
		[][]any{{1, "a"}},
	)
	right := NewDataFrameFromRows(
		[]string{"id", "v2"},
		[][]any{{2, "b"}},
	)
	result, err := Merge(left, right, []string{"id"}, "inner")
	if err != nil {
		t.Fatalf("Merge error: %v", err)
	}
	if result.Len() != 0 {
		t.Fatalf("expected 0 rows, got %d", result.Len())
	}
}

func TestMergeUnsupportedHow(t *testing.T) {
	df := NewDataFrameFromRows([]string{"id"}, [][]any{{1}})
	_, err := Merge(df, df, []string{"id"}, "outer")
	if err == nil {
		t.Fatal("expected error for unsupported merge type")
	}
}

func TestMergeMissingColumn(t *testing.T) {
	left := NewDataFrameFromRows([]string{"id"}, [][]any{{1}})
	right := NewDataFrameFromRows([]string{"key"}, [][]any{{1}})
	_, err := Merge(left, right, []string{"id"}, "inner")
	if err == nil {
		t.Fatal("expected error for missing column in right")
	}
}

func TestMergeNoOnColumns(t *testing.T) {
	df := NewDataFrameFromRows([]string{"id"}, [][]any{{1}})
	_, err := Merge(df, df, nil, "inner")
	if err == nil {
		t.Fatal("expected error for empty on columns")
	}
}

func TestConcatBasic(t *testing.T) {
	df1 := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1, "x"},
			{2, "y"},
		},
	)
	df2 := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{3, "z"},
		},
	)
	result, err := Concat([]*DataFrame{df1, df2})
	if err != nil {
		t.Fatalf("Concat error: %v", err)
	}
	if result.Len() != 3 {
		t.Fatalf("expected 3 rows, got %d", result.Len())
	}
	aVals := result.Column("a").Values()
	if !reflect.DeepEqual(aVals, []any{1, 2, 3}) {
		t.Fatalf("a values = %v, want [1 2 3]", aVals)
	}
	bVals := result.Column("b").Values()
	if !reflect.DeepEqual(bVals, []any{"x", "y", "z"}) {
		t.Fatalf("b values = %v, want [x y z]", bVals)
	}
}

func TestConcatMultiple(t *testing.T) {
	frames := make([]*DataFrame, 3)
	for i := range frames {
		frames[i] = NewDataFrameFromRows(
			[]string{"v"},
			[][]any{{i}},
		)
	}
	result, err := Concat(frames)
	if err != nil {
		t.Fatalf("Concat error: %v", err)
	}
	if result.Len() != 3 {
		t.Fatalf("expected 3 rows, got %d", result.Len())
	}
}

func TestConcatEmpty(t *testing.T) {
	result, err := Concat(nil)
	if err != nil {
		t.Fatalf("Concat nil error: %v", err)
	}
	if result.Len() != 0 {
		t.Fatalf("expected 0 rows, got %d", result.Len())
	}
}

func TestConcatMismatchedColumns(t *testing.T) {
	df1 := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}})
	df2 := NewDataFrameFromRows([]string{"a", "c"}, [][]any{{3, 4}})
	_, err := Concat([]*DataFrame{df1, df2})
	if err == nil {
		t.Fatal("expected error for mismatched columns")
	}
}

func TestConcatDifferentColumnCount(t *testing.T) {
	df1 := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	df2 := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}})
	_, err := Concat([]*DataFrame{df1, df2})
	if err == nil {
		t.Fatal("expected error for different column count")
	}
}

func TestConcatPreservesColumnOrder(t *testing.T) {
	df1 := NewDataFrameFromRows(
		[]string{"b", "a"},
		[][]any{
			{1, 2},
		},
	)
	df2 := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{3, 4},
		},
	)
	result, err := Concat([]*DataFrame{df1, df2})
	if err != nil {
		t.Fatalf("Concat error: %v", err)
	}
	// Column order from first frame
	cols := result.Columns()
	if !reflect.DeepEqual(cols, []string{"b", "a"}) {
		t.Fatalf("columns = %v, want [b a]", cols)
	}
	// Values should be aligned by column name, not position
	bVals := result.Column("b").Values()
	aVals := result.Column("a").Values()
	if !reflect.DeepEqual(bVals, []any{1, 4}) {
		t.Fatalf("b values = %v, want [1 4]", bVals)
	}
	if !reflect.DeepEqual(aVals, []any{2, 3}) {
		t.Fatalf("a values = %v, want [2 3]", aVals)
	}
}
