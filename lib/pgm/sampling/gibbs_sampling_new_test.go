//go:build unit

package sampling

import (
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
)

func buildSimpleBN(t *testing.T) *models.BayesianNetwork {
	t.Helper()
	bn := models.NewBayesianNetwork()
	for _, n := range []string{"A", "B"} {
		if err := bn.AddNode(n); err != nil {
			t.Fatal(err)
		}
	}
	_ = bn.AddEdge("A", "B")

	cpdA, err := factors.NewTabularCPD("A", 2, [][]float64{{0.6}, {0.4}}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = bn.AddCPD(cpdA)

	cpdB, err := factors.NewTabularCPD("B", 2, [][]float64{{0.9, 0.1}, {0.1, 0.9}}, []string{"A"}, []int{2})
	if err != nil {
		t.Fatal(err)
	}
	_ = bn.AddCPD(cpdB)

	return bn
}

func TestGenerateSample(t *testing.T) {
	bn := buildSimpleBN(t)
	gs, err := NewGibbsSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}

	sample, err := gs.GenerateSample(100, nil)
	if err != nil {
		t.Fatalf("GenerateSample failed: %v", err)
	}

	if len(sample) != 2 {
		t.Errorf("expected 2 variables in sample, got %d", len(sample))
	}

	// Check values are valid states.
	for v, val := range sample {
		if val < 0 || val > 1 {
			t.Errorf("variable %q has invalid state %d", v, val)
		}
	}
}

func TestGenerateSampleWithEvidence(t *testing.T) {
	bn := buildSimpleBN(t)
	gs, err := NewGibbsSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}

	evidence := map[string]int{"A": 1}
	sample, err := gs.GenerateSample(100, evidence)
	if err != nil {
		t.Fatalf("GenerateSample with evidence failed: %v", err)
	}

	if sample["A"] != 1 {
		t.Errorf("expected A=1 (evidence), got A=%d", sample["A"])
	}
}

func TestGenerateSampleErrors(t *testing.T) {
	bn := buildSimpleBN(t)
	gs, err := NewGibbsSampling(bn, 42)
	if err != nil {
		t.Fatal(err)
	}

	_, err = gs.GenerateSample(-1, nil)
	if err == nil {
		t.Error("expected error for negative burnIn")
	}

	_, err = gs.GenerateSample(0, map[string]int{"Z": 0})
	if err == nil {
		t.Error("expected error for unknown evidence variable")
	}

	_, err = gs.GenerateSample(0, map[string]int{"A": 5})
	if err == nil {
		t.Error("expected error for out-of-range evidence value")
	}
}
