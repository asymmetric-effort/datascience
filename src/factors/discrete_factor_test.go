//go:build unit

package factors

import (
	"fmt"
	"math"
	"testing"
)

const epsilon = 1e-9

func floatEq(a, b float64) bool {
	return math.Abs(a-b) < epsilon
}

// ---------------------------------------------------------------------------
// NewDiscreteFactor
// ---------------------------------------------------------------------------

func TestNewDiscreteFactor_Basic(t *testing.T) {
	// Factor over A(card=2), B(card=3)
	f, err := NewDiscreteFactor(
		[]string{"A", "B"},
		[]int{2, 3},
		[]float64{1, 2, 3, 4, 5, 6},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := f.Variables(); len(got) != 2 || got[0] != "A" || got[1] != "B" {
		t.Errorf("Variables() = %v", got)
	}
	if got := f.Cardinality(); len(got) != 2 || got[0] != 2 || got[1] != 3 {
		t.Errorf("Cardinality() = %v", got)
	}
	if got := f.Values().Size(); got != 6 {
		t.Errorf("Values().Size() = %d, want 6", got)
	}
}

func TestNewDiscreteFactor_Errors(t *testing.T) {
	tests := []struct {
		name string
		vars []string
		card []int
		vals []float64
	}{
		{"mismatched lengths", []string{"A"}, []int{2, 3}, []float64{1, 2}},
		{"wrong value count", []string{"A", "B"}, []int{2, 3}, []float64{1, 2, 3}},
		{"zero cardinality", []string{"A"}, []int{0}, []float64{}},
		{"negative cardinality", []string{"A"}, []int{-1}, []float64{}},
		{"duplicate variables", []string{"A", "A"}, []int{2, 2}, []float64{1, 2, 3, 4}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewDiscreteFactor(tc.vars, tc.card, tc.vals)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetValue / SetValue
// ---------------------------------------------------------------------------

func TestGetSetValue(t *testing.T) {
	f, err := NewDiscreteFactor(
		[]string{"A", "B"},
		[]int{2, 3},
		[]float64{1, 2, 3, 4, 5, 6},
	)
	if err != nil {
		t.Fatal(err)
	}

	// Row-major: index [0,0]=1, [0,1]=2, [0,2]=3, [1,0]=4, [1,1]=5, [1,2]=6
	cases := []struct {
		a, b int
		want float64
	}{
		{0, 0, 1}, {0, 1, 2}, {0, 2, 3},
		{1, 0, 4}, {1, 1, 5}, {1, 2, 6},
	}
	for _, c := range cases {
		got := f.GetValue(map[string]int{"A": c.a, "B": c.b})
		if !floatEq(got, c.want) {
			t.Errorf("GetValue(A=%d,B=%d) = %f, want %f", c.a, c.b, got, c.want)
		}
	}

	f.SetValue(map[string]int{"A": 1, "B": 2}, 99.0)
	if got := f.GetValue(map[string]int{"A": 1, "B": 2}); !floatEq(got, 99.0) {
		t.Errorf("after SetValue, got %f, want 99", got)
	}
}

// ---------------------------------------------------------------------------
// Marginalize
// ---------------------------------------------------------------------------

func TestMarginalize_SumOutOne(t *testing.T) {
	// f(A,B) with A in {0,1}, B in {0,1,2}
	// values: A=0,B=0:1  A=0,B=1:2  A=0,B=2:3  A=1,B=0:4  A=1,B=1:5  A=1,B=2:6
	f, _ := NewDiscreteFactor(
		[]string{"A", "B"},
		[]int{2, 3},
		[]float64{1, 2, 3, 4, 5, 6},
	)

	// Marginalize out B: result over A.
	// f(A=0) = 1+2+3 = 6, f(A=1) = 4+5+6 = 15
	result, err := f.Marginalize([]string{"B"})
	if err != nil {
		t.Fatal(err)
	}
	if vars := result.Variables(); len(vars) != 1 || vars[0] != "A" {
		t.Errorf("expected variables=[A], got %v", vars)
	}
	if !floatEq(result.GetValue(map[string]int{"A": 0}), 6) {
		t.Errorf("f(A=0) = %f, want 6", result.GetValue(map[string]int{"A": 0}))
	}
	if !floatEq(result.GetValue(map[string]int{"A": 1}), 15) {
		t.Errorf("f(A=1) = %f, want 15", result.GetValue(map[string]int{"A": 1}))
	}

	// Marginalize out A: result over B.
	// f(B=0)=1+4=5, f(B=1)=2+5=7, f(B=2)=3+6=9
	result2, err := f.Marginalize([]string{"A"})
	if err != nil {
		t.Fatal(err)
	}
	if vars := result2.Variables(); len(vars) != 1 || vars[0] != "B" {
		t.Errorf("expected variables=[B], got %v", vars)
	}
	expected := []float64{5, 7, 9}
	for i, want := range expected {
		got := result2.GetValue(map[string]int{"B": i})
		if !floatEq(got, want) {
			t.Errorf("f(B=%d) = %f, want %f", i, got, want)
		}
	}
}

func TestMarginalize_ThreeVars(t *testing.T) {
	// f(A,B,C) A=2, B=2, C=2, values 0..7
	vals := make([]float64, 8)
	for i := range vals {
		vals[i] = float64(i)
	}
	f, _ := NewDiscreteFactor([]string{"A", "B", "C"}, []int{2, 2, 2}, vals)

	// Marginalize out C.
	// f(A=0,B=0) = 0+1 = 1
	// f(A=0,B=1) = 2+3 = 5
	// f(A=1,B=0) = 4+5 = 9
	// f(A=1,B=1) = 6+7 = 13
	result, err := f.Marginalize([]string{"C"})
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]float64{
		"0,0": 1, "0,1": 5, "1,0": 9, "1,1": 13,
	}
	for key, want := range expected {
		var a, b int
		fmt.Sscanf(key, "%d,%d", &a, &b)
		got := result.GetValue(map[string]int{"A": a, "B": b})
		if !floatEq(got, want) {
			t.Errorf("f(A=%d,B=%d) = %f, want %f", a, b, got, want)
		}
	}

	// Marginalize out A and C from f(A,B,C).
	// f(B=0) = 0+1+4+5 = 10
	// f(B=1) = 2+3+6+7 = 18
	result2, err := f.Marginalize([]string{"A", "C"})
	if err != nil {
		t.Fatal(err)
	}
	if !floatEq(result2.GetValue(map[string]int{"B": 0}), 10) {
		t.Errorf("f(B=0) = %f, want 10", result2.GetValue(map[string]int{"B": 0}))
	}
	if !floatEq(result2.GetValue(map[string]int{"B": 1}), 18) {
		t.Errorf("f(B=1) = %f, want 18", result2.GetValue(map[string]int{"B": 1}))
	}
}

func TestMarginalize_EmptyList(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A"}, []int{2}, []float64{3, 7})
	result, err := f.Marginalize(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !floatEq(result.GetValue(map[string]int{"A": 0}), 3) {
		t.Error("expected copy with same values")
	}
}

func TestMarginalize_Errors(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A", "B"}, []int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	_, err := f.Marginalize([]string{"C"})
	if err == nil {
		t.Error("expected error for unknown variable")
	}
	_, err = f.Marginalize([]string{"A", "B"})
	if err == nil {
		t.Error("expected error for marginalizing all variables")
	}
}

// ---------------------------------------------------------------------------
// Reduce
// ---------------------------------------------------------------------------

func TestReduce_SingleEvidence(t *testing.T) {
	// f(A,B) A=2, B=3, vals=[1,2,3,4,5,6]
	f, _ := NewDiscreteFactor([]string{"A", "B"}, []int{2, 3}, []float64{1, 2, 3, 4, 5, 6})

	// Reduce B=1: result is f(A) = [2, 5]
	result, err := f.Reduce(map[string]int{"B": 1})
	if err != nil {
		t.Fatal(err)
	}
	if vars := result.Variables(); len(vars) != 1 || vars[0] != "A" {
		t.Errorf("expected [A], got %v", vars)
	}
	if !floatEq(result.GetValue(map[string]int{"A": 0}), 2) {
		t.Errorf("f(A=0|B=1) = %f, want 2", result.GetValue(map[string]int{"A": 0}))
	}
	if !floatEq(result.GetValue(map[string]int{"A": 1}), 5) {
		t.Errorf("f(A=1|B=1) = %f, want 5", result.GetValue(map[string]int{"A": 1}))
	}
}

func TestReduce_MultipleEvidence(t *testing.T) {
	vals := make([]float64, 8)
	for i := range vals {
		vals[i] = float64(i + 1)
	}
	f, _ := NewDiscreteFactor([]string{"A", "B", "C"}, []int{2, 2, 2}, vals)

	// Reduce A=1, C=0: remain B
	// Original: A=1,B=0,C=0 -> index 4 -> val 5; A=1,B=1,C=0 -> index 6 -> val 7
	result, err := f.Reduce(map[string]int{"A": 1, "C": 0})
	if err != nil {
		t.Fatal(err)
	}
	if vars := result.Variables(); len(vars) != 1 || vars[0] != "B" {
		t.Errorf("expected [B], got %v", vars)
	}
	if !floatEq(result.GetValue(map[string]int{"B": 0}), 5) {
		t.Errorf("f(B=0|A=1,C=0) = %f, want 5", result.GetValue(map[string]int{"B": 0}))
	}
	if !floatEq(result.GetValue(map[string]int{"B": 1}), 7) {
		t.Errorf("f(B=1|A=1,C=0) = %f, want 7", result.GetValue(map[string]int{"B": 1}))
	}
}

func TestReduce_Empty(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A"}, []int{2}, []float64{3, 7})
	result, err := f.Reduce(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !floatEq(result.GetValue(map[string]int{"A": 0}), 3) {
		t.Error("expected copy")
	}
}

func TestReduce_Errors(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A", "B"}, []int{2, 3}, []float64{1, 2, 3, 4, 5, 6})
	_, err := f.Reduce(map[string]int{"C": 0})
	if err == nil {
		t.Error("expected error for unknown variable")
	}
	_, err = f.Reduce(map[string]int{"B": 5})
	if err == nil {
		t.Error("expected error for out-of-range value")
	}
}

// ---------------------------------------------------------------------------
// Normalize
// ---------------------------------------------------------------------------

func TestNormalize(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	f.Normalize()
	data := f.Values().Data()
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	if !floatEq(sum, 1.0) {
		t.Errorf("after normalize, sum = %f, want 1.0", sum)
	}
	// 1/10=0.1, 2/10=0.2, 3/10=0.3, 4/10=0.4
	expected := []float64{0.1, 0.2, 0.3, 0.4}
	for i, want := range expected {
		if !floatEq(data[i], want) {
			t.Errorf("data[%d] = %f, want %f", i, data[i], want)
		}
	}
}

func TestNormalize_AllZeros(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A"}, []int{3}, []float64{0, 0, 0})
	f.Normalize() // Should not panic or produce NaN.
	data := f.Values().Data()
	for i, v := range data {
		if v != 0 {
			t.Errorf("data[%d] = %f, want 0 (all-zero normalization)", i, v)
		}
	}
}

// ---------------------------------------------------------------------------
// Copy
// ---------------------------------------------------------------------------

func TestCopy(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	c := f.Copy()
	// Modify original.
	f.SetValue(map[string]int{"A": 0, "B": 0}, 99)
	if c.GetValue(map[string]int{"A": 0, "B": 0}) == 99 {
		t.Error("copy was affected by modifying original")
	}
}

// ---------------------------------------------------------------------------
// String
// ---------------------------------------------------------------------------

func TestString(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.3, 0.7})
	s := f.String()
	if s == "" {
		t.Error("String() returned empty")
	}
}
