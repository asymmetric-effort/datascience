//go:build unit

package tabgo

import (
	"reflect"
	"testing"
)

func TestAssignNew(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}, {2}})
	df2 := df.Assign("b", []any{10, 20})
	if len(df2.Columns()) != 2 {
		t.Fatalf("Assign new column: ncols = %d, want 2", len(df2.Columns()))
	}
	bVals := df2.Column("b").Values()
	if !reflect.DeepEqual(bVals, []any{10, 20}) {
		t.Fatalf("Assign b = %v", bVals)
	}
	// Original unchanged.
	if len(df.Columns()) != 1 {
		t.Fatal("Assign mutated original")
	}
}

func TestAssignReplace(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}, {3, 4}})
	df2 := df.Assign("a", []any{10, 30})
	aVals := df2.Column("a").Values()
	if !reflect.DeepEqual(aVals, []any{10, 30}) {
		t.Fatalf("Assign replace a = %v", aVals)
	}
	// b unchanged.
	bVals := df2.Column("b").Values()
	if !reflect.DeepEqual(bVals, []any{2, 4}) {
		t.Fatalf("Assign replace b = %v", bVals)
	}
}

func TestAssignPanicsOnLengthMismatch(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}, {2}})
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for length mismatch")
		}
	}()
	df.Assign("b", []any{1})
}

func TestInsert(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "c"}, [][]any{{1, 3}, {4, 6}})
	df2, err := df.Insert(1, "b", []any{2, 5})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(df2.Columns(), []string{"a", "b", "c"}) {
		t.Fatalf("Insert columns = %v", df2.Columns())
	}
	bVals := df2.Column("b").Values()
	if !reflect.DeepEqual(bVals, []any{2, 5}) {
		t.Fatalf("Insert b = %v", bVals)
	}
}

func TestInsertAtStart(t *testing.T) {
	df := NewDataFrameFromRows([]string{"b"}, [][]any{{2}})
	df2, err := df.Insert(0, "a", []any{1})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(df2.Columns(), []string{"a", "b"}) {
		t.Fatalf("Insert at 0 columns = %v", df2.Columns())
	}
}

func TestInsertAtEnd(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	df2, err := df.Insert(1, "b", []any{2})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(df2.Columns(), []string{"a", "b"}) {
		t.Fatalf("Insert at end columns = %v", df2.Columns())
	}
}

func TestInsertErrors(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	// Bad loc.
	if _, err := df.Insert(-1, "b", []any{2}); err == nil {
		t.Fatal("expected error for negative loc")
	}
	if _, err := df.Insert(5, "b", []any{2}); err == nil {
		t.Fatal("expected error for loc > ncols")
	}
	// Duplicate column.
	if _, err := df.Insert(0, "a", []any{2}); err == nil {
		t.Fatal("expected error for duplicate column")
	}
	// Length mismatch.
	if _, err := df.Insert(0, "b", []any{1, 2}); err == nil {
		t.Fatal("expected error for length mismatch")
	}
}

func TestDrop(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b", "c"}, [][]any{{1, 2, 3}})
	df2 := df.Drop("b")
	if !reflect.DeepEqual(df2.Columns(), []string{"a", "c"}) {
		t.Fatalf("Drop columns = %v", df2.Columns())
	}
}

func TestDropMultiple(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b", "c"}, [][]any{{1, 2, 3}})
	df2 := df.Drop("a", "c")
	if !reflect.DeepEqual(df2.Columns(), []string{"b"}) {
		t.Fatalf("Drop multiple columns = %v", df2.Columns())
	}
}

func TestDropAll(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	df2 := df.Drop("a")
	if len(df2.Columns()) != 0 {
		t.Fatalf("Drop all columns = %v", df2.Columns())
	}
}

func TestRename(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}})
	df2 := df.Rename(map[string]string{"a": "x", "b": "y"})
	if !reflect.DeepEqual(df2.Columns(), []string{"x", "y"}) {
		t.Fatalf("Rename columns = %v", df2.Columns())
	}
	// Values preserved.
	if !reflect.DeepEqual(df2.Column("x").Values(), []any{1}) {
		t.Fatalf("Rename x values = %v", df2.Column("x").Values())
	}
}

func TestRenamePartial(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}})
	df2 := df.Rename(map[string]string{"a": "x"})
	if !reflect.DeepEqual(df2.Columns(), []string{"x", "b"}) {
		t.Fatalf("Rename partial columns = %v", df2.Columns())
	}
}

func TestSortValues(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"name", "score"},
		[][]any{
			{"Carol", 78},
			{"Alice", 92},
			{"Bob", 85},
		},
	)
	asc := df.SortValues("score", true)
	scores := asc.Column("score").Float64()
	if scores[0] != 78 || scores[1] != 85 || scores[2] != 92 {
		t.Fatalf("SortValues asc = %v", scores)
	}
	names := asc.Column("name").Values()
	if names[0] != "Carol" || names[1] != "Bob" || names[2] != "Alice" {
		t.Fatalf("SortValues asc names = %v", names)
	}

	desc := df.SortValues("score", false)
	scoresDesc := desc.Column("score").Float64()
	if scoresDesc[0] != 92 || scoresDesc[1] != 85 || scoresDesc[2] != 78 {
		t.Fatalf("SortValues desc = %v", scoresDesc)
	}
}

func TestAppend(t *testing.T) {
	df1 := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}})
	df2 := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{3, 4}})
	result, err := df1.Append(df2)
	if err != nil {
		t.Fatal(err)
	}
	if result.Len() != 2 {
		t.Fatalf("Append Len = %d, want 2", result.Len())
	}
	aVals := result.Column("a").Values()
	if !reflect.DeepEqual(aVals, []any{1, 3}) {
		t.Fatalf("Append a = %v", aVals)
	}
}

func TestAppendMismatchedColumns(t *testing.T) {
	df1 := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	df2 := NewDataFrameFromRows([]string{"b"}, [][]any{{2}})
	_, err := df1.Append(df2)
	if err == nil {
		t.Fatal("expected error for mismatched columns")
	}
}
