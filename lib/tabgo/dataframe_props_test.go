//go:build unit

package tabgo

import (
	"reflect"
	"strings"
	"testing"
)

func TestShape(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}, {3, 4}, {5, 6}})
	got := df.Shape()
	if got != [2]int{3, 2} {
		t.Fatalf("Shape = %v, want [3 2]", got)
	}
}

func TestShapeEmpty(t *testing.T) {
	df := NewDataFrame(map[string]*Series{})
	got := df.Shape()
	if got != [2]int{0, 0} {
		t.Fatalf("Shape = %v, want [0 0]", got)
	}
}

func TestValues(t *testing.T) {
	df := NewDataFrameFromRows([]string{"x", "y"}, [][]any{{1, "a"}, {2, "b"}})
	vals := df.Values()
	expected := [][]any{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(vals, expected) {
		t.Fatalf("Values = %v, want %v", vals, expected)
	}
}

func TestValuesEmpty(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, nil)
	vals := df.Values()
	if len(vals) != 0 {
		t.Fatalf("Values of empty df should be empty, got %v", vals)
	}
}

func TestSize(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b", "c"}, [][]any{{1, 2, 3}, {4, 5, 6}})
	if df.Size() != 6 {
		t.Fatalf("Size = %d, want 6", df.Size())
	}
}

func TestSizeEmpty(t *testing.T) {
	df := NewDataFrame(map[string]*Series{})
	if df.Size() != 0 {
		t.Fatalf("Size = %d, want 0", df.Size())
	}
}

func TestNdim(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{1}})
	if df.Ndim() != 2 {
		t.Fatalf("Ndim = %d, want 2", df.Ndim())
	}
}

func TestEmptyProp(t *testing.T) {
	tests := []struct {
		name string
		df   *DataFrame
		want bool
	}{
		{"no columns", NewDataFrame(map[string]*Series{}), true},
		{"no rows", NewDataFrameFromRows([]string{"a"}, nil), true},
		{"has data", NewDataFrameFromRows([]string{"a"}, [][]any{{1}}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.df.Empty(); got != tt.want {
				t.Fatalf("Empty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTranspose(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, [][]any{{1, 2}, {3, 4}})
	tr := df.T()

	// Transposed should have 2 rows (original columns) and 2 columns named "0", "1".
	if tr.Shape() != [2]int{2, 2} {
		t.Fatalf("T() Shape = %v, want [2 2]", tr.Shape())
	}
	cols := tr.Columns()
	if !reflect.DeepEqual(cols, []string{"0", "1"}) {
		t.Fatalf("T() Columns = %v, want [0, 1]", cols)
	}

	// Col "0" of transposed = first element of each original column: a[0]=1, b[0]=2
	col0 := tr.Column("0").Values()
	if !reflect.DeepEqual(col0, []any{1, 2}) {
		t.Fatalf("T() col 0 = %v, want [1, 2]", col0)
	}
	// Col "1" of transposed = second element of each original column: a[1]=3, b[1]=4
	col1 := tr.Column("1").Values()
	if !reflect.DeepEqual(col1, []any{3, 4}) {
		t.Fatalf("T() col 1 = %v, want [3, 4]", col1)
	}
}

func TestTransposeEmpty(t *testing.T) {
	df := NewDataFrame(map[string]*Series{})
	tr := df.T()
	if tr.Len() != 0 {
		t.Fatalf("T() of empty should have 0 rows, got %d", tr.Len())
	}
}

func TestToString(t *testing.T) {
	df := NewDataFrameFromRows([]string{"name", "age"}, [][]any{{"Alice", 30}, {"Bob", 25}})
	s := df.ToString()
	if !strings.Contains(s, "name") || !strings.Contains(s, "age") {
		t.Fatalf("ToString missing headers: %s", s)
	}
	if !strings.Contains(s, "Alice") || !strings.Contains(s, "Bob") {
		t.Fatalf("ToString missing values: %s", s)
	}
	if !strings.Contains(s, "30") || !strings.Contains(s, "25") {
		t.Fatalf("ToString missing numeric values: %s", s)
	}
}

func TestToStringEmpty(t *testing.T) {
	df := NewDataFrame(map[string]*Series{})
	s := df.ToString()
	if s != "Empty DataFrame" {
		t.Fatalf("ToString of empty = %q, want 'Empty DataFrame'", s)
	}
}
