//go:build unit

package factors

import (
	"fmt"
	"testing"
)

// ---------------------------------------------------------------------------
// FactorProduct
// ---------------------------------------------------------------------------

func TestFactorProduct_SharedVariable(t *testing.T) {
	// f1(A,B): A=2, B=2, values=[1,2,3,4]
	// f2(B,C): B=2, C=2, values=[5,6,7,8]
	// Product f(A,B,C): A=2, B=2, C=2
	f1, _ := NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{1, 2, 3, 4})
	f2, _ := NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{5, 6, 7, 8})

	result, err := FactorProduct(f1, f2)
	if err != nil {
		t.Fatal(err)
	}

	vars := result.Variables()
	if len(vars) != 3 || vars[0] != "A" || vars[1] != "B" || vars[2] != "C" {
		t.Fatalf("expected [A B C], got %v", vars)
	}
	card := result.Cardinality()
	if card[0] != 2 || card[1] != 2 || card[2] != 2 {
		t.Fatalf("expected [2 2 2], got %v", card)
	}

	// Verify all 8 values.
	// f(A,B,C) = f1(A,B) * f2(B,C)
	expected := map[string]float64{
		"0,0,0": 1 * 5, "0,0,1": 1 * 6, "0,1,0": 2 * 7, "0,1,1": 2 * 8,
		"1,0,0": 3 * 5, "1,0,1": 3 * 6, "1,1,0": 4 * 7, "1,1,1": 4 * 8,
	}
	for key, want := range expected {
		var a, b, c int
		fmt.Sscanf(key, "%d,%d,%d", &a, &b, &c)
		got := result.GetValue(map[string]int{"A": a, "B": b, "C": c})
		if !floatEq(got, want) {
			t.Errorf("f(A=%d,B=%d,C=%d) = %f, want %f", a, b, c, got, want)
		}
	}
}

func TestFactorProduct_NoSharedVariables(t *testing.T) {
	// f1(A): [2, 3], f2(B): [5, 7]
	// Product: f(A,B) = outer product
	f1, _ := NewDiscreteFactor([]string{"A"}, []int{2}, []float64{2, 3})
	f2, _ := NewDiscreteFactor([]string{"B"}, []int{2}, []float64{5, 7})

	result, err := FactorProduct(f1, f2)
	if err != nil {
		t.Fatal(err)
	}

	expected := map[string]float64{
		"0,0": 10, "0,1": 14, "1,0": 15, "1,1": 21,
	}
	for key, want := range expected {
		var a, b int
		fmt.Sscanf(key, "%d,%d", &a, &b)
		got := result.GetValue(map[string]int{"A": a, "B": b})
		if !floatEq(got, want) {
			t.Errorf("f(A=%d,B=%d) = %f, want %f", a, b, got, want)
		}
	}
}

func TestFactorProduct_SameVariables(t *testing.T) {
	// Both factors over the same variable A.
	f1, _ := NewDiscreteFactor([]string{"A"}, []int{3}, []float64{1, 2, 3})
	f2, _ := NewDiscreteFactor([]string{"A"}, []int{3}, []float64{4, 5, 6})

	result, err := FactorProduct(f1, f2)
	if err != nil {
		t.Fatal(err)
	}

	if vars := result.Variables(); len(vars) != 1 || vars[0] != "A" {
		t.Fatalf("expected [A], got %v", vars)
	}
	// Element-wise product.
	expected := []float64{4, 10, 18}
	for i, want := range expected {
		got := result.GetValue(map[string]int{"A": i})
		if !floatEq(got, want) {
			t.Errorf("f(A=%d) = %f, want %f", i, got, want)
		}
	}
}

func TestFactorProduct_ThreeFactors(t *testing.T) {
	f1, _ := NewDiscreteFactor([]string{"A"}, []int{2}, []float64{1, 2})
	f2, _ := NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{3, 4, 5, 6})
	f3, _ := NewDiscreteFactor([]string{"B"}, []int{2}, []float64{7, 8})

	result, err := FactorProduct(f1, f2, f3)
	if err != nil {
		t.Fatal(err)
	}

	// result(A,B) = f1(A) * f2(A,B) * f3(B)
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			want := float64((a+1)*(3+a*2+b)) * float64(7+b)
			got := result.GetValue(map[string]int{"A": a, "B": b})
			if !floatEq(got, want) {
				t.Errorf("f(A=%d,B=%d) = %f, want %f", a, b, got, want)
			}
		}
	}
}

func TestFactorProduct_SingleFactor(t *testing.T) {
	f, _ := NewDiscreteFactor([]string{"A"}, []int{2}, []float64{3, 7})
	result, err := FactorProduct(f)
	if err != nil {
		t.Fatal(err)
	}
	if !floatEq(result.GetValue(map[string]int{"A": 0}), 3) {
		t.Error("single factor product should return copy")
	}
}

func TestFactorProduct_Empty(t *testing.T) {
	_, err := FactorProduct()
	if err == nil {
		t.Error("expected error for empty FactorProduct")
	}
}

func TestFactorProduct_MismatchedCardinality(t *testing.T) {
	f1, _ := NewDiscreteFactor([]string{"A"}, []int{2}, []float64{1, 2})
	f2, _ := NewDiscreteFactor([]string{"A"}, []int{3}, []float64{1, 2, 3})
	_, err := FactorProduct(f1, f2)
	if err == nil {
		t.Error("expected error for mismatched cardinality")
	}
}

// ---------------------------------------------------------------------------
// Integration: Product then Marginalize
// ---------------------------------------------------------------------------

func TestProductThenMarginalize(t *testing.T) {
	// Classic example: P(A,B) = P(A|B) * P(B), then marginalize B to get P(A).
	// P(B=0) = 0.4, P(B=1) = 0.6
	// P(A=0|B=0) = 0.2, P(A=1|B=0) = 0.8
	// P(A=0|B=1) = 0.5, P(A=1|B=1) = 0.5
	pB, _ := NewDiscreteFactor([]string{"B"}, []int{2}, []float64{0.4, 0.6})
	// P(A|B) as factor over (A,B): [P(A=0,B=0), P(A=0,B=1), P(A=1,B=0), P(A=1,B=1)]
	pAgivenB, _ := NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.2, 0.5, 0.8, 0.5})

	joint, err := FactorProduct(pAgivenB, pB)
	if err != nil {
		t.Fatal(err)
	}

	pA, err := joint.Marginalize([]string{"B"})
	if err != nil {
		t.Fatal(err)
	}

	// P(A=0) = 0.2*0.4 + 0.5*0.6 = 0.08 + 0.30 = 0.38
	// P(A=1) = 0.8*0.4 + 0.5*0.6 = 0.32 + 0.30 = 0.62
	if !floatEq(pA.GetValue(map[string]int{"A": 0}), 0.38) {
		t.Errorf("P(A=0) = %f, want 0.38", pA.GetValue(map[string]int{"A": 0}))
	}
	if !floatEq(pA.GetValue(map[string]int{"A": 1}), 0.62) {
		t.Errorf("P(A=1) = %f, want 0.62", pA.GetValue(map[string]int{"A": 1}))
	}
}
