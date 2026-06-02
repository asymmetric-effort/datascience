//go:build unit

package tabgo

import (
	"math"
	"testing"
)

func makeGroupByDF() *DataFrame {
	return NewDataFrameFromRows(
		[]string{"dept", "role", "salary"},
		[][]any{
			{"eng", "dev", 100.0},
			{"eng", "dev", 120.0},
			{"eng", "mgr", 150.0},
			{"sales", "rep", 80.0},
			{"sales", "rep", 90.0},
			{"sales", "mgr", 130.0},
		},
	)
}

func TestGroupBySingleColumn(t *testing.T) {
	df := makeGroupByDF()
	g := df.GroupBy("dept")
	groups := g.Groups()
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	eng := groups["eng"]
	if eng == nil || eng.Len() != 3 {
		t.Fatalf("eng group: expected 3 rows, got %v", eng)
	}
	sales := groups["sales"]
	if sales == nil || sales.Len() != 3 {
		t.Fatalf("sales group: expected 3 rows, got %v", sales)
	}
}

func TestGroupByMultiColumn(t *testing.T) {
	df := makeGroupByDF()
	g := df.GroupBy("dept", "role")
	groups := g.Groups()
	if len(groups) != 4 {
		t.Fatalf("expected 4 groups, got %d", len(groups))
	}
	engDev := groups["eng|dev"]
	if engDev == nil || engDev.Len() != 2 {
		t.Fatalf("eng|dev group: expected 2 rows, got %v", engDev)
	}
	salesMgr := groups["sales|mgr"]
	if salesMgr == nil || salesMgr.Len() != 1 {
		t.Fatalf("sales|mgr group: expected 1 row, got %v", salesMgr)
	}
}

func TestGroupByCount(t *testing.T) {
	df := makeGroupByDF()
	result := df.GroupBy("dept").Count()

	if result.Len() != 2 {
		t.Fatalf("expected 2 rows, got %d", result.Len())
	}
	// Rows are sorted by key, so eng < sales
	deptVals := result.Column("dept").Values()
	countVals := result.Column("count").Values()

	foundEng := false
	foundSales := false
	for i, d := range deptVals {
		switch d {
		case "eng":
			foundEng = true
			if countVals[i] != 3 {
				t.Errorf("eng count = %v, want 3", countVals[i])
			}
		case "sales":
			foundSales = true
			if countVals[i] != 3 {
				t.Errorf("sales count = %v, want 3", countVals[i])
			}
		}
	}
	if !foundEng || !foundSales {
		t.Fatal("missing expected groups in Count result")
	}
}

func TestGroupByMultiColumnCount(t *testing.T) {
	df := makeGroupByDF()
	result := df.GroupBy("dept", "role").Count()
	if result.Len() != 4 {
		t.Fatalf("expected 4 rows, got %d", result.Len())
	}
	// Verify columns exist
	cols := result.Columns()
	expected := map[string]bool{"dept": true, "role": true, "count": true}
	for _, c := range cols {
		if !expected[c] {
			t.Errorf("unexpected column %q", c)
		}
	}
}

func TestGroupBySum(t *testing.T) {
	df := makeGroupByDF()
	result := df.GroupBy("dept").Sum("salary")

	if result.Len() != 2 {
		t.Fatalf("expected 2 rows, got %d", result.Len())
	}

	deptVals := result.Column("dept").Values()
	salaryVals := result.Column("salary").Values()
	for i, d := range deptVals {
		s := toFloat64(salaryVals[i])
		switch d {
		case "eng":
			if math.Abs(s-370.0) > 0.01 {
				t.Errorf("eng salary sum = %f, want 370", s)
			}
		case "sales":
			if math.Abs(s-300.0) > 0.01 {
				t.Errorf("sales salary sum = %f, want 300", s)
			}
		}
	}
}

func TestGroupByMean(t *testing.T) {
	df := makeGroupByDF()
	result := df.GroupBy("dept").Mean("salary")

	if result.Len() != 2 {
		t.Fatalf("expected 2 rows, got %d", result.Len())
	}

	deptVals := result.Column("dept").Values()
	salaryVals := result.Column("salary").Values()
	for i, d := range deptVals {
		m := toFloat64(salaryVals[i])
		switch d {
		case "eng":
			want := 370.0 / 3.0
			if math.Abs(m-want) > 0.01 {
				t.Errorf("eng salary mean = %f, want %f", m, want)
			}
		case "sales":
			want := 300.0 / 3.0
			if math.Abs(m-want) > 0.01 {
				t.Errorf("sales salary mean = %f, want %f", m, want)
			}
		}
	}
}

func TestGroupByMeanMultiColumn(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"g", "x", "y"},
		[][]any{
			{"a", 10.0, 100.0},
			{"a", 20.0, 200.0},
			{"b", 30.0, 300.0},
		},
	)
	result := df.GroupBy("g").Mean("x", "y")
	if result.Len() != 2 {
		t.Fatalf("expected 2 rows, got %d", result.Len())
	}
	gVals := result.Column("g").Values()
	xVals := result.Column("x").Values()
	yVals := result.Column("y").Values()
	for i, g := range gVals {
		xm := toFloat64(xVals[i])
		ym := toFloat64(yVals[i])
		switch g {
		case "a":
			if math.Abs(xm-15.0) > 0.01 || math.Abs(ym-150.0) > 0.01 {
				t.Errorf("group a: x=%f y=%f, want 15, 150", xm, ym)
			}
		case "b":
			if math.Abs(xm-30.0) > 0.01 || math.Abs(ym-300.0) > 0.01 {
				t.Errorf("group b: x=%f y=%f, want 30, 300", xm, ym)
			}
		}
	}
}

func TestGroupByApply(t *testing.T) {
	df := makeGroupByDF()
	// Apply function that returns Head(1) for each group
	result := df.GroupBy("dept").Apply(func(sub *DataFrame) *DataFrame {
		return sub.Head(1)
	})
	if result.Len() != 2 {
		t.Fatalf("expected 2 rows, got %d", result.Len())
	}
}

func TestGroupByApplyTransform(t *testing.T) {
	df := NewDataFrameFromRows(
		[]string{"g", "v"},
		[][]any{
			{"a", 1},
			{"a", 2},
			{"b", 10},
		},
	)
	// Identity apply should return all rows
	result := df.GroupBy("g").Apply(func(sub *DataFrame) *DataFrame {
		return sub
	})
	if result.Len() != 3 {
		t.Fatalf("expected 3 rows, got %d", result.Len())
	}
}

func TestGroupByEmpty(t *testing.T) {
	df := NewDataFrameFromRows([]string{"a", "b"}, nil)
	g := df.GroupBy("a")
	groups := g.Groups()
	if len(groups) != 0 {
		t.Fatalf("expected 0 groups, got %d", len(groups))
	}
	count := g.Count()
	if count.Len() != 0 {
		t.Fatalf("expected 0 rows in count, got %d", count.Len())
	}
}
