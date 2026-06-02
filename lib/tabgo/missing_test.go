//go:build unit

package tabgo

import (
	"reflect"
	"testing"
)

func TestSeriesIsNA(t *testing.T) {
	s := NewSeries("x", []any{1, nil, "hello", nil, 3.14})
	got := s.IsNA()
	want := []bool{false, true, false, true, false}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("IsNA = %v, want %v", got, want)
	}
}

func TestSeriesIsNAEmpty(t *testing.T) {
	s := NewSeries("x", nil)
	got := s.IsNA()
	if len(got) != 0 {
		t.Fatalf("IsNA on empty should return empty slice, got %v", got)
	}
}

func TestSeriesIsNAAllNil(t *testing.T) {
	s := NewSeries("x", []any{nil, nil, nil})
	got := s.IsNA()
	want := []bool{true, true, true}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("IsNA = %v, want %v", got, want)
	}
}

func TestSeriesIsNANoNil(t *testing.T) {
	s := NewSeries("x", []any{1, 2, 3})
	got := s.IsNA()
	want := []bool{false, false, false}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("IsNA = %v, want %v", got, want)
	}
}

func TestDropNA(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1, "x"},
			{nil, "y"},
			{3, nil},
			{4, "w"},
		},
	)
	result := df.DropNA()
	if result.Len() != 2 {
		t.Fatalf("DropNA: expected 2 rows, got %d", result.Len())
	}
	aVals := result.Column("a").Values()
	bVals := result.Column("b").Values()
	if !reflect.DeepEqual(aVals, []any{1, 4}) {
		t.Fatalf("a values = %v, want [1 4]", aVals)
	}
	if !reflect.DeepEqual(bVals, []any{"x", "w"}) {
		t.Fatalf("b values = %v, want [x w]", bVals)
	}
}

func TestDropNANoNils(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a"},
		[][]any{{1}, {2}, {3}},
	)
	result := df.DropNA()
	if result.Len() != 3 {
		t.Fatalf("DropNA with no nils: expected 3 rows, got %d", result.Len())
	}
}

func TestDropNAAllNils(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{nil, nil},
			{nil, 1},
			{1, nil},
		},
	)
	result := df.DropNA()
	if result.Len() != 0 {
		t.Fatalf("DropNA all-nil: expected 0 rows, got %d", result.Len())
	}
}

func TestFillNA(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{1, nil},
			{nil, "y"},
			{nil, nil},
		},
	)
	result := df.FillNA(0)
	if result.Len() != 3 {
		t.Fatalf("FillNA: expected 3 rows, got %d", result.Len())
	}
	aVals := result.Column("a").Values()
	bVals := result.Column("b").Values()
	if !reflect.DeepEqual(aVals, []any{1, 0, 0}) {
		t.Fatalf("a values = %v, want [1 0 0]", aVals)
	}
	if !reflect.DeepEqual(bVals, []any{0, "y", 0}) {
		t.Fatalf("b values = %v, want [0 y 0]", bVals)
	}
}

func TestFillNADoesNotMutateOriginal(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a"},
		[][]any{{nil}, {1}},
	)
	_ = df.FillNA(99)
	// Original should still have nil
	if df.Column("a").Values()[0] != nil {
		t.Fatal("FillNA mutated original DataFrame")
	}
}

func TestFillNAColumn(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a", "b"},
		[][]any{
			{nil, nil},
			{1, nil},
			{nil, "z"},
		},
	)
	result := df.FillNAColumn("a", -1)
	aVals := result.Column("a").Values()
	bVals := result.Column("b").Values()
	if !reflect.DeepEqual(aVals, []any{-1, 1, -1}) {
		t.Fatalf("a values = %v, want [-1 1 -1]", aVals)
	}
	// b should still have nils
	if bVals[0] != nil || bVals[1] != nil {
		t.Fatalf("b values should still have nils, got %v", bVals)
	}
	if bVals[2] != "z" {
		t.Fatalf("b[2] = %v, want z", bVals[2])
	}
}

func TestFillNAColumnDoesNotMutateOriginal(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"a"},
		[][]any{{nil}},
	)
	_ = df.FillNAColumn("a", 42)
	if df.Column("a").Values()[0] != nil {
		t.Fatal("FillNAColumn mutated original DataFrame")
	}
}
