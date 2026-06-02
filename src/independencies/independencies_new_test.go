//go:build unit

package independencies

import (
	"strings"
	"testing"
)

func TestGetAllVariables(t *testing.T) {
	ind := NewIndependencies()
	ind.Add(NewIndependenceAssertion([]string{"A"}, []string{"B"}, []string{"C"}))
	ind.Add(NewIndependenceAssertion([]string{"D"}, []string{"E"}, nil))

	vars := ind.GetAllVariables()
	expected := map[string]bool{"A": true, "B": true, "C": true, "D": true, "E": true}

	if len(vars) != len(expected) {
		t.Errorf("expected %d variables, got %d: %v", len(expected), len(vars), vars)
	}
	for _, v := range vars {
		if !expected[v] {
			t.Errorf("unexpected variable %q", v)
		}
	}
}

func TestGetAllVariablesEmpty(t *testing.T) {
	ind := NewIndependencies()
	vars := ind.GetAllVariables()
	if len(vars) != 0 {
		t.Errorf("expected 0 variables, got %d", len(vars))
	}
}

func TestClosure(t *testing.T) {
	ind := NewIndependencies()
	// A _|_ {B, C}
	ind.Add(NewIndependenceAssertion([]string{"A"}, []string{"B", "C"}, nil))

	closure := ind.Closure()

	// Should contain the original.
	orig := NewIndependenceAssertion([]string{"A"}, []string{"B", "C"}, nil)
	if !closure.Contains(orig) {
		t.Error("closure should contain original assertion")
	}

	// Symmetry: {B, C} _|_ A
	sym := NewIndependenceAssertion([]string{"B", "C"}, []string{"A"}, nil)
	if !closure.Contains(sym) {
		t.Error("closure should contain symmetric assertion")
	}

	// Decomposition: A _|_ B and A _|_ C
	decompB := NewIndependenceAssertion([]string{"A"}, []string{"B"}, nil)
	if !closure.Contains(decompB) {
		t.Error("closure should contain decomposition A _|_ B")
	}

	decompC := NewIndependenceAssertion([]string{"A"}, []string{"C"}, nil)
	if !closure.Contains(decompC) {
		t.Error("closure should contain decomposition A _|_ C")
	}
}

func TestClosureEmpty(t *testing.T) {
	ind := NewIndependencies()
	closure := ind.Closure()
	if closure.Len() != 0 {
		t.Errorf("expected empty closure, got %d assertions", closure.Len())
	}
}

func TestEntails(t *testing.T) {
	ind := NewIndependencies()
	// A _|_ {B, C}
	ind.Add(NewIndependenceAssertion([]string{"A"}, []string{"B", "C"}, nil))

	// Should entail A _|_ B (subset).
	sub := NewIndependenceAssertion([]string{"A"}, []string{"B"}, nil)
	if !ind.Entails(sub) {
		t.Error("expected A _|_ {B,C} to entail A _|_ B")
	}

	// Should not entail D _|_ E.
	other := NewIndependenceAssertion([]string{"D"}, []string{"E"}, nil)
	if ind.Entails(other) {
		t.Error("should not entail unrelated assertion")
	}
}

func TestEntailsNil(t *testing.T) {
	ind := NewIndependencies()
	if !ind.Entails(nil) {
		t.Error("nil assertion should always be entailed")
	}
}

func TestReduce(t *testing.T) {
	ind := NewIndependencies()
	// A _|_ {B, C} entails A _|_ B (via containment), so the latter is redundant.
	ind.Add(NewIndependenceAssertion([]string{"A"}, []string{"B", "C"}, nil))
	ind.Add(NewIndependenceAssertion([]string{"A"}, []string{"B"}, nil))

	reduced := ind.Reduce()
	// Only the broader assertion should remain.
	if reduced.Len() != 1 {
		t.Errorf("expected 1 assertion after reduction, got %d", reduced.Len())
	}
}

func TestReduceNothingToRemove(t *testing.T) {
	ind := NewIndependencies()
	ind.Add(NewIndependenceAssertion([]string{"A"}, []string{"B"}, nil))
	ind.Add(NewIndependenceAssertion([]string{"C"}, []string{"D"}, nil))

	reduced := ind.Reduce()
	if reduced.Len() != 2 {
		t.Errorf("expected 2 assertions, got %d", reduced.Len())
	}
}

func TestLatexString(t *testing.T) {
	ind := NewIndependencies()
	ind.Add(NewIndependenceAssertion([]string{"A"}, []string{"B"}, []string{"C"}))

	latex := ind.LatexString()
	if !strings.Contains(latex, "\\perp") {
		t.Error("expected \\perp in LaTeX output")
	}
	if !strings.Contains(latex, "\\mid") {
		t.Error("expected \\mid in LaTeX output")
	}
}

func TestLatexStringEmpty(t *testing.T) {
	ind := NewIndependencies()
	latex := ind.LatexString()
	if latex != "\\emptyset" {
		t.Errorf("expected \\emptyset, got %q", latex)
	}
}

func TestGetFactorizedProduct(t *testing.T) {
	ind := NewIndependencies()
	ind.Add(NewIndependenceAssertion([]string{"A"}, []string{"B"}, []string{"C"}))

	product := ind.GetFactorizedProduct()
	if !strings.Contains(product, "P(A | C)") {
		t.Errorf("expected P(A | C) in product, got %q", product)
	}
}

func TestGetFactorizedProductEmpty(t *testing.T) {
	ind := NewIndependencies()
	product := ind.GetFactorizedProduct()
	if product != "P(V)" {
		t.Errorf("expected P(V), got %q", product)
	}
}
