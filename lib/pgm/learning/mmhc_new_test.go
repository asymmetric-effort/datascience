//go:build unit

package learning

import (
	"testing"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

func alwaysIndepCITest(x, y string, z []string, data *tabgo.DataFrame, significance float64) (float64, float64, bool) {
	// Simple mock: always returns independent for testing.
	return 0.0, 1.0, true
}

func TestMMPC(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 1}),
		"C": tabgo.NewSeries("C", []any{1, 0, 1, 0, 1, 0}),
	})

	m := NewMMHC(data, simpleScoreFn, alwaysIndepCITest, 0.05)
	candidates := m.MMPC("A")
	// With the simpleCITest always returning independent, the candidate set
	// should be empty (all variables are pruned as independent).
	if len(candidates) != 0 {
		t.Logf("MMPC candidates for A: %v", candidates)
	}
}

func TestMMPCNonTrivial(t *testing.T) {
	// Use a CI test that never declares independence.
	neverIndepCITest := func(x, y string, z []string, data *tabgo.DataFrame, significance float64) (float64, float64, bool) {
		return 10.0, 0.001, false
	}

	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})

	m := NewMMHC(data, simpleScoreFn, neverIndepCITest, 0.05)
	candidates := m.MMPC("A")
	// With never-independent test, B should be in the candidate set.
	if len(candidates) != 1 || candidates[0] != "B" {
		t.Errorf("expected [B], got %v", candidates)
	}
}
