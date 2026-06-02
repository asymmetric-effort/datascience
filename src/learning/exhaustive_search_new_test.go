//go:build unit

package learning

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
)

func TestAllDAGs(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})

	es := NewExhaustiveSearch(data, simpleScoreFn)
	dags, err := es.AllDAGs()
	if err != nil {
		t.Fatalf("AllDAGs failed: %v", err)
	}

	// For 2 variables, there are 3 possible DAGs: no edge, A->B, B->A.
	if len(dags) != 3 {
		t.Errorf("expected 3 DAGs for 2 variables, got %d", len(dags))
	}
}

func TestAllDAGsThreeVars(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1}),
	})

	es := NewExhaustiveSearch(data, simpleScoreFn)
	dags, err := es.AllDAGs()
	if err != nil {
		t.Fatalf("AllDAGs failed: %v", err)
	}

	// For 3 variables, there are 25 possible DAGs.
	if len(dags) != 25 {
		t.Errorf("expected 25 DAGs for 3 variables, got %d", len(dags))
	}
}

func TestAllDAGsTooManyVars(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0}),
		"B": tabgo.NewSeries("B", []any{0}),
		"C": tabgo.NewSeries("C", []any{0}),
		"D": tabgo.NewSeries("D", []any{0}),
		"E": tabgo.NewSeries("E", []any{0}),
	})

	es := NewExhaustiveSearch(data, simpleScoreFn)
	_, err := es.AllDAGs()
	if err == nil {
		t.Error("expected error for too many variables")
	}
}
