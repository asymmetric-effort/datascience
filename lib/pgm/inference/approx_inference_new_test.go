//go:build unit

package inference

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

func TestApproxGetDistribution(t *testing.T) {
	// Two binary variables A and B with a simple joint factor.
	fAB, err := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2},
		[]float64{0.3, 0.1, 0.2, 0.4})
	if err != nil {
		t.Fatal(err)
	}

	ai := NewApproxInference([]*factors.DiscreteFactor{fAB}, 42)
	dist, err := ai.GetDistribution(100000)
	if err != nil {
		t.Fatalf("GetDistribution failed: %v", err)
	}

	vars := dist.Variables()
	if len(vars) != 2 {
		t.Errorf("expected 2 variables, got %d", len(vars))
	}

	// Values should roughly match the normalized input.
	total := 0.3 + 0.1 + 0.2 + 0.4
	expected := []float64{0.3 / total, 0.1 / total, 0.2 / total, 0.4 / total}
	data := dist.Values().Data()
	for i, exp := range expected {
		if math.Abs(data[i]-exp) > 0.05 {
			t.Errorf("value[%d]: expected ~%f, got %f", i, exp, data[i])
		}
	}
}

func TestApproxGetDistributionError(t *testing.T) {
	fA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{fA}, 42)

	_, err := ai.GetDistribution(0)
	if err == nil {
		t.Error("expected error for nSamples=0")
	}
}

func TestApproxMAPQuery(t *testing.T) {
	// Factor where A=1, B=1 is strongly preferred.
	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2},
		[]float64{0.01, 0.01, 0.01, 0.97})

	ai := NewApproxInference([]*factors.DiscreteFactor{fAB}, 42)
	assignment, err := ai.MAPQuery([]string{"A", "B"}, nil, 100000)
	if err != nil {
		t.Fatalf("MAPQuery failed: %v", err)
	}

	if assignment["A"] != 1 || assignment["B"] != 1 {
		t.Errorf("expected A=1, B=1, got A=%d, B=%d", assignment["A"], assignment["B"])
	}
}

func TestApproxMAPQueryErrors(t *testing.T) {
	fA, _ := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{0.5, 0.5})
	ai := NewApproxInference([]*factors.DiscreteFactor{fA}, 42)

	_, err := ai.MAPQuery(nil, nil, 100)
	if err == nil {
		t.Error("expected error for empty queryVars")
	}

	_, err = ai.MAPQuery([]string{"A"}, nil, 0)
	if err == nil {
		t.Error("expected error for nSamples=0")
	}
}
