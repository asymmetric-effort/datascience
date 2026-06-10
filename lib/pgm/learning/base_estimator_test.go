//go:build unit

package learning

import (
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

func TestBaseEstimator_NewAndGetters(t *testing.T) {
	bn := models.NewBayesianNetwork()
	if err := bn.AddNode("A"); err != nil {
		t.Fatal(err)
	}

	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
	})

	base := NewBaseEstimator(bn, df)

	if base.GetModel() != bn {
		t.Error("GetModel returned wrong model")
	}
	if base.GetData() != df {
		t.Error("GetData returned wrong data")
	}
}

func TestBaseEstimator_NilModel(t *testing.T) {
	base := NewBaseEstimator(nil, nil)
	if base.GetModel() != nil {
		t.Error("expected nil model")
	}
	if base.GetData() != nil {
		t.Error("expected nil data")
	}
}

func TestStructureEstimator_NewAndGetters(t *testing.T) {
	df := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{0, 1, 2}),
		"Y": tabgo.NewSeries("Y", []any{1, 0, 1}),
	})

	se := NewStructureEstimator(df)
	if se.GetData() != df {
		t.Error("GetData returned wrong data")
	}

	vars := se.Variables()
	if len(vars) != 2 {
		t.Fatalf("expected 2 variables, got %d", len(vars))
	}

	// Variables should be a copy.
	vars[0] = "MODIFIED"
	origVars := se.Variables()
	if origVars[0] == "MODIFIED" {
		t.Error("Variables returned a reference, not a copy")
	}
}

func TestStructureEstimator_NilData(t *testing.T) {
	se := NewStructureEstimator(nil)
	if se.GetData() != nil {
		t.Error("expected nil data")
	}
	vars := se.Variables()
	if vars != nil {
		t.Error("expected nil variables for nil data")
	}
}
