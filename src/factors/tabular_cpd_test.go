//go:build unit

package factors

import (
	"testing"
)

// ---------------------------------------------------------------------------
// NewTabularCPD
// ---------------------------------------------------------------------------

func TestNewTabularCPD_NoEvidence(t *testing.T) {
	// P(A) with A having 3 states.
	cpd, err := NewTabularCPD("A", 3,
		[][]float64{{0.2}, {0.3}, {0.5}},
		nil, nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	if cpd.Variable() != "A" {
		t.Errorf("Variable() = %q, want A", cpd.Variable())
	}
	if len(cpd.Evidence()) != 0 {
		t.Errorf("Evidence() should be empty, got %v", cpd.Evidence())
	}
	if err := cpd.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestNewTabularCPD_WithEvidence(t *testing.T) {
	// P(Grade | Difficulty, Intelligence)
	// Grade: 3 states, Difficulty: 2, Intelligence: 2
	// 4 parent configs (D=0,I=0), (D=0,I=1), (D=1,I=0), (D=1,I=1)
	cpd, err := NewTabularCPD("Grade", 3,
		[][]float64{
			{0.3, 0.05, 0.9, 0.5},  // P(Grade=0 | ...)
			{0.4, 0.25, 0.08, 0.3}, // P(Grade=1 | ...)
			{0.3, 0.7, 0.02, 0.2},  // P(Grade=2 | ...)
		},
		[]string{"Difficulty", "Intelligence"},
		[]int{2, 2},
	)
	if err != nil {
		t.Fatal(err)
	}
	if cpd.Variable() != "Grade" {
		t.Errorf("Variable() = %q", cpd.Variable())
	}
	ev := cpd.Evidence()
	if len(ev) != 2 || ev[0] != "Difficulty" || ev[1] != "Intelligence" {
		t.Errorf("Evidence() = %v", ev)
	}
	ec := cpd.EvidenceCard()
	if len(ec) != 2 || ec[0] != 2 || ec[1] != 2 {
		t.Errorf("EvidenceCard() = %v", ec)
	}
	if err := cpd.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
}

func TestNewTabularCPD_Errors(t *testing.T) {
	tests := []struct {
		name     string
		variable string
		varCard  int
		values   [][]float64
		evidence []string
		evCard   []int
	}{
		{"zero cardinality", "A", 0, nil, nil, nil},
		{"evidence mismatch", "A", 2, [][]float64{{0.5}, {0.5}}, []string{"B"}, []int{2, 3}},
		{"wrong row count", "A", 2, [][]float64{{0.5}}, nil, nil},
		{"wrong col count", "A", 2, [][]float64{{0.5, 0.5}, {0.5, 0.5}}, []string{"B"}, []int{3}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewTabularCPD(tc.variable, tc.varCard, tc.values, tc.evidence, tc.evCard)
			if err == nil {
				t.Error("expected error")
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Validate
// ---------------------------------------------------------------------------

func TestValidate_Invalid(t *testing.T) {
	// Columns that do not sum to 1.
	cpd, err := NewTabularCPD("A", 2,
		[][]float64{
			{0.3, 0.4},
			{0.3, 0.4},
		},
		[]string{"B"}, []int{2},
	)
	if err != nil {
		t.Fatal(err)
	}
	if err := cpd.Validate(); err == nil {
		t.Error("expected validation error for columns not summing to 1")
	}
}

// ---------------------------------------------------------------------------
// ToFactor
// ---------------------------------------------------------------------------

func TestToFactor(t *testing.T) {
	cpd, _ := NewTabularCPD("A", 2,
		[][]float64{
			{0.4, 0.9},
			{0.6, 0.1},
		},
		[]string{"B"}, []int{2},
	)
	f := cpd.ToFactor()
	// Factor should have variables [A, B], card [2, 2].
	vars := f.Variables()
	if len(vars) != 2 || vars[0] != "A" || vars[1] != "B" {
		t.Errorf("ToFactor variables = %v", vars)
	}
	// P(A=0, B=0) = 0.4, P(A=0, B=1) = 0.9
	if !floatEq(f.GetValue(map[string]int{"A": 0, "B": 0}), 0.4) {
		t.Errorf("f(A=0,B=0) = %f, want 0.4", f.GetValue(map[string]int{"A": 0, "B": 0}))
	}
	if !floatEq(f.GetValue(map[string]int{"A": 0, "B": 1}), 0.9) {
		t.Errorf("f(A=0,B=1) = %f, want 0.9", f.GetValue(map[string]int{"A": 0, "B": 1}))
	}
	if !floatEq(f.GetValue(map[string]int{"A": 1, "B": 0}), 0.6) {
		t.Errorf("f(A=1,B=0) = %f, want 0.6", f.GetValue(map[string]int{"A": 1, "B": 0}))
	}
	if !floatEq(f.GetValue(map[string]int{"A": 1, "B": 1}), 0.1) {
		t.Errorf("f(A=1,B=1) = %f, want 0.1", f.GetValue(map[string]int{"A": 1, "B": 1}))
	}
}

// ---------------------------------------------------------------------------
// Copy
// ---------------------------------------------------------------------------

func TestTabularCPDCopy(t *testing.T) {
	cpd, _ := NewTabularCPD("A", 2,
		[][]float64{{0.3}, {0.7}},
		nil, nil,
	)
	c := cpd.Copy()
	if c.Variable() != "A" {
		t.Error("copy variable mismatch")
	}
	// Verify independence.
	f1 := cpd.ToFactor()
	f2 := c.ToFactor()
	f1.SetValue(map[string]int{"A": 0}, 99)
	if f2.GetValue(map[string]int{"A": 0}) == 99 {
		t.Error("copy was affected by modifying original factor")
	}
}

// ---------------------------------------------------------------------------
// Integration: CPD -> Factor -> Product -> Marginalize
// ---------------------------------------------------------------------------

func TestCPDIntegration(t *testing.T) {
	// Simple Bayesian network: B -> A
	// P(B=0)=0.4, P(B=1)=0.6
	// P(A=0|B=0)=0.2, P(A=0|B=1)=0.5
	// P(A=1|B=0)=0.8, P(A=1|B=1)=0.5
	cpdB, _ := NewTabularCPD("B", 2, [][]float64{{0.4}, {0.6}}, nil, nil)
	cpdA, _ := NewTabularCPD("A", 2,
		[][]float64{
			{0.2, 0.5},
			{0.8, 0.5},
		},
		[]string{"B"}, []int{2},
	)

	if err := cpdB.Validate(); err != nil {
		t.Fatal(err)
	}
	if err := cpdA.Validate(); err != nil {
		t.Fatal(err)
	}

	fB := cpdB.ToFactor()
	fA := cpdA.ToFactor()

	joint, err := FactorProduct(fA, fB)
	if err != nil {
		t.Fatal(err)
	}

	pA, err := joint.Marginalize([]string{"B"})
	if err != nil {
		t.Fatal(err)
	}

	// P(A=0) = 0.2*0.4 + 0.5*0.6 = 0.38
	// P(A=1) = 0.8*0.4 + 0.5*0.6 = 0.62
	if !floatEq(pA.GetValue(map[string]int{"A": 0}), 0.38) {
		t.Errorf("P(A=0) = %f, want 0.38", pA.GetValue(map[string]int{"A": 0}))
	}
	if !floatEq(pA.GetValue(map[string]int{"A": 1}), 0.62) {
		t.Errorf("P(A=1) = %f, want 0.62", pA.GetValue(map[string]int{"A": 1}))
	}
}
