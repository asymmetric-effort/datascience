//go:build unit

package tabgo

import (
	"testing"
)

func TestMergeLeft(t *testing.T) {
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
			{1, "x"},
			{3, "z"},
		},
	)
	result, err := Merge(left, right, []string{"id"}, "left")
	if err != nil {
		t.Fatalf("Merge left error: %v", err)
	}
	if result.Len() != 3 {
		t.Errorf("Merge left rows = %d, want 3", result.Len())
	}
	rvals := result.Column("rval").Values()
	// id=2 has no match, rval should be nil
	if rvals[1] != nil {
		t.Errorf("Merge left unmatched rval = %v, want nil", rvals[1])
	}
	if rvals[0] != "x" {
		t.Errorf("Merge left matched rval[0] = %v, want x", rvals[0])
	}
}

func TestMergeRight(t *testing.T) {
	left := NewDataFrameFromRows(
		[]string{"id", "lval"},
		[][]any{
			{1, "a"},
			{3, "c"},
		},
	)
	right := NewDataFrameFromRows(
		[]string{"id", "rval"},
		[][]any{
			{1, "x"},
			{2, "y"},
			{3, "z"},
		},
	)
	result, err := Merge(left, right, []string{"id"}, "right")
	if err != nil {
		t.Fatalf("Merge right error: %v", err)
	}
	if result.Len() != 3 {
		t.Errorf("Merge right rows = %d, want 3", result.Len())
	}
	// id=2 from right has no left match -> lval should be nil
	// Find the row with id=2
	ids := result.Column("id").Values()
	lvals := result.Column("lval").Values()
	for i, id := range ids {
		if toInt(id) == 2 {
			if lvals[i] != nil {
				t.Errorf("Merge right unmatched lval = %v, want nil", lvals[i])
			}
		}
	}
}

func TestMergeOuter(t *testing.T) {
	left := NewDataFrameFromRows(
		[]string{"id", "lval"},
		[][]any{
			{1, "a"},
			{2, "b"},
		},
	)
	right := NewDataFrameFromRows(
		[]string{"id", "rval"},
		[][]any{
			{2, "y"},
			{3, "z"},
		},
	)
	result, err := Merge(left, right, []string{"id"}, "outer")
	if err != nil {
		t.Fatalf("Merge outer error: %v", err)
	}
	if result.Len() != 3 {
		t.Errorf("Merge outer rows = %d, want 3", result.Len())
	}
}

func TestMergeInnerStillWorks(t *testing.T) {
	left := NewDataFrameFromRows(
		[]string{"id", "val"},
		[][]any{{1, "a"}, {2, "b"}},
	)
	right := NewDataFrameFromRows(
		[]string{"id", "val2"},
		[][]any{{2, "x"}, {3, "y"}},
	)
	result, err := Merge(left, right, []string{"id"}, "inner")
	if err != nil {
		t.Fatalf("Merge inner error: %v", err)
	}
	if result.Len() != 1 {
		t.Errorf("Merge inner rows = %d, want 1", result.Len())
	}
}

func TestMergeLeftAllMatch(t *testing.T) {
	left := NewDataFrameFromRows(
		[]string{"id", "lval"},
		[][]any{{1, "a"}, {2, "b"}},
	)
	right := NewDataFrameFromRows(
		[]string{"id", "rval"},
		[][]any{{1, "x"}, {2, "y"}},
	)
	result, err := Merge(left, right, []string{"id"}, "left")
	if err != nil {
		t.Fatalf("Merge left all match error: %v", err)
	}
	if result.Len() != 2 {
		t.Errorf("Merge left all match rows = %d, want 2", result.Len())
	}
	// No nil values expected
	rvals := result.Column("rval").Values()
	for i, v := range rvals {
		if v == nil {
			t.Errorf("Merge left all match rval[%d] = nil, want non-nil", i)
		}
	}
}
