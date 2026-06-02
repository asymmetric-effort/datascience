//go:build unit

package tabgo

import (
	"testing"
)

func TestMelt(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"id", "a", "b"},
		[][]any{
			{1, 10, 20},
			{2, 30, 40},
		},
	)
	result, err := Melt(df, []string{"id"}, []string{"a", "b"})
	if err != nil {
		t.Fatalf("Melt error: %v", err)
	}
	if result.Len() != 4 {
		t.Errorf("Melt rows = %d, want 4", result.Len())
	}
	cols := result.Columns()
	if len(cols) != 3 {
		t.Errorf("Melt cols = %d, want 3 (id, variable, value)", len(cols))
	}
	// Check variable column contains "a" and "b"
	varVals := result.Column("variable").Values()
	if varVals[0] != "a" || varVals[1] != "a" || varVals[2] != "b" || varVals[3] != "b" {
		t.Errorf("Melt variable values = %v", varVals)
	}
}

func TestMeltInvalidColumn(t *testing.T) {
	df := NewDataFrameFromRows([]string{"id"}, [][]any{{1}})
	_, err := Melt(df, []string{"missing"}, []string{"id"})
	if err == nil {
		t.Error("Melt should fail with missing column")
	}
}

func TestPivotTableSum(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"row", "col", "val"},
		[][]any{
			{"r1", "c1", 1.0},
			{"r1", "c2", 2.0},
			{"r2", "c1", 3.0},
			{"r2", "c2", 4.0},
			{"r1", "c1", 5.0},
		},
	)
	result, err := PivotTable(df, "row", "col", "val", "sum")
	if err != nil {
		t.Fatalf("PivotTable error: %v", err)
	}
	if result.Len() != 2 {
		t.Errorf("PivotTable rows = %d, want 2", result.Len())
	}
	// r1, c1 should be 1+5=6
	c1Vals := result.Column("c1").Values()
	if toFloat64(c1Vals[0]) != 6.0 {
		t.Errorf("PivotTable r1,c1 = %v, want 6.0", c1Vals[0])
	}
}

func TestPivotTableMean(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"row", "col", "val"},
		[][]any{
			{"r1", "c1", 2.0},
			{"r1", "c1", 4.0},
		},
	)
	result, err := PivotTable(df, "row", "col", "val", "mean")
	if err != nil {
		t.Fatalf("PivotTable error: %v", err)
	}
	c1Vals := result.Column("c1").Values()
	if toFloat64(c1Vals[0]) != 3.0 {
		t.Errorf("PivotTable mean = %v, want 3.0", c1Vals[0])
	}
}

func TestPivotTableCount(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"row", "col", "val"},
		[][]any{
			{"r1", "c1", 1.0},
			{"r1", "c1", 2.0},
			{"r1", "c2", 3.0},
		},
	)
	result, err := PivotTable(df, "row", "col", "val", "count")
	if err != nil {
		t.Fatalf("PivotTable error: %v", err)
	}
	c1Vals := result.Column("c1").Values()
	if toFloat64(c1Vals[0]) != 2.0 {
		t.Errorf("PivotTable count = %v, want 2.0", c1Vals[0])
	}
}

func TestPivotTableInvalidColumn(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{"x"}})
	_, err := PivotTable(df, "a", "b", "c", "sum")
	if err == nil {
		t.Error("PivotTable should fail with missing column")
	}
}

func TestCrosstab(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"gender", "hand"},
		[][]any{
			{"M", "R"},
			{"M", "L"},
			{"F", "R"},
			{"F", "R"},
			{"M", "R"},
		},
	)
	result, err := Crosstab(df, "gender", "hand")
	if err != nil {
		t.Fatalf("Crosstab error: %v", err)
	}
	if result.Len() != 2 {
		t.Errorf("Crosstab rows = %d, want 2", result.Len())
	}
	// F: L=0, R=2
	// M: L=1, R=2
	rVals := result.Column("R").Values()
	lVals := result.Column("L").Values()
	// Sorted: F first, M second
	if toInt(rVals[0]) != 2 {
		t.Errorf("Crosstab F,R = %v, want 2", rVals[0])
	}
	if toInt(lVals[0]) != 0 {
		t.Errorf("Crosstab F,L = %v, want 0", lVals[0])
	}
	if toInt(rVals[1]) != 2 {
		t.Errorf("Crosstab M,R = %v, want 2", rVals[1])
	}
	if toInt(lVals[1]) != 1 {
		t.Errorf("Crosstab M,L = %v, want 1", lVals[1])
	}
}

func TestCrosstabInvalidColumn(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a"}, [][]any{{"x"}})
	_, err := Crosstab(df, "a", "missing")
	if err == nil {
		t.Error("Crosstab should fail with missing column")
	}
}
