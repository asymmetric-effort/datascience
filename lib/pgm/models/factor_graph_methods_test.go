//go:build unit

package models

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
)

func buildSimpleFactorGraph(t *testing.T) *FactorGraph {
	t.Helper()
	fg := NewFactorGraph()
	_ = fg.AddVariable("A", 2)
	_ = fg.AddVariable("B", 2)
	_ = fg.AddVariable("C", 2)

	f1, err := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.1, 0.2, 0.3, 0.4})
	if err != nil {
		t.Fatalf("NewDiscreteFactor f1: %v", err)
	}
	f2, err := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{0.4, 0.3, 0.2, 0.1})
	if err != nil {
		t.Fatalf("NewDiscreteFactor f2: %v", err)
	}

	if err := fg.AddFactor(f1); err != nil {
		t.Fatalf("AddFactor f1: %v", err)
	}
	if err := fg.AddFactor(f2); err != nil {
		t.Fatalf("AddFactor f2: %v", err)
	}

	return fg
}

func TestFactorGraphAddEdge(t *testing.T) {
	fg := buildSimpleFactorGraph(t)

	// A is already in factor 0's scope, so adding edge should succeed.
	if err := fg.AddEdge("A", 0); err != nil {
		t.Fatalf("AddEdge: %v", err)
	}
}

func TestFactorGraphAddEdgeUnknownVariable(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	if err := fg.AddEdge("Z", 0); err == nil {
		t.Error("expected error for unknown variable")
	}
}

func TestFactorGraphAddEdgeInvalidFactor(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	if err := fg.AddEdge("A", 10); err == nil {
		t.Error("expected error for out-of-range factor index")
	}
	if err := fg.AddEdge("A", -1); err == nil {
		t.Error("expected error for negative factor index")
	}
}

func TestFactorGraphAddEdgeNotInScope(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	// C is not in factor 0's scope (A,B).
	if err := fg.AddEdge("C", 0); err == nil {
		t.Error("expected error for variable not in factor scope")
	}
}

func TestFactorGraphRemoveFactors(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	if len(fg.GetFactors()) != 2 {
		t.Fatalf("expected 2 factors, got %d", len(fg.GetFactors()))
	}

	fg.RemoveFactors()

	if len(fg.GetFactors()) != 0 {
		t.Errorf("expected 0 factors after removal, got %d", len(fg.GetFactors()))
	}
	// Variables should still exist.
	if len(fg.GetVariables()) != 3 {
		t.Errorf("expected 3 variables, got %d", len(fg.GetVariables()))
	}
	// GetFactorsOf should return nil.
	if fg.GetFactorsOf("A") != nil {
		t.Error("expected nil factors for A after removal")
	}
}

func TestFactorGraphGetCardinality(t *testing.T) {
	fg := buildSimpleFactorGraph(t)

	card, err := fg.GetCardinality("A")
	if err != nil {
		t.Fatalf("GetCardinality: %v", err)
	}
	if card != 2 {
		t.Errorf("expected cardinality 2, got %d", card)
	}
}

func TestFactorGraphGetCardinalityUnknown(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	_, err := fg.GetCardinality("Z")
	if err == nil {
		t.Error("expected error for unknown variable")
	}
}

func TestFactorGraphGetFactorNodes(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	nodes := fg.GetFactorNodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 factor nodes, got %d", len(nodes))
	}
	if len(nodes[0]) != 2 {
		t.Errorf("expected 2 variables in factor 0, got %d", len(nodes[0]))
	}
}

func TestFactorGraphToJunctionTree(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	jt, err := fg.ToJunctionTree()
	if err != nil {
		t.Fatalf("ToJunctionTree: %v", err)
	}
	if jt == nil {
		t.Fatal("ToJunctionTree returned nil")
	}

	cliques := jt.Cliques()
	if len(cliques) == 0 {
		t.Error("expected at least one clique")
	}
}

func TestFactorGraphToJunctionTreeInvalid(t *testing.T) {
	fg := NewFactorGraph()
	_, err := fg.ToJunctionTree()
	if err == nil {
		t.Error("expected error for invalid factor graph")
	}
}

func TestFactorGraphGetPartitionFunction(t *testing.T) {
	fg := NewFactorGraph()
	_ = fg.AddVariable("X", 2)
	f, _ := factors.NewDiscreteFactor([]string{"X"}, []int{2}, []float64{0.3, 0.7})
	_ = fg.AddFactor(f)

	z, err := fg.GetPartitionFunction()
	if err != nil {
		t.Fatalf("GetPartitionFunction: %v", err)
	}
	if math.Abs(z-1.0) > 1e-6 {
		t.Errorf("expected partition function 1.0, got %f", z)
	}
}

func TestFactorGraphGetPartitionFunctionMultiple(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	z, err := fg.GetPartitionFunction()
	if err != nil {
		t.Fatalf("GetPartitionFunction: %v", err)
	}
	if z <= 0 {
		t.Errorf("expected positive partition function, got %f", z)
	}
}

func TestFactorGraphGetPartitionFunctionNoFactors(t *testing.T) {
	fg := NewFactorGraph()
	_, err := fg.GetPartitionFunction()
	if err == nil {
		t.Error("expected error for no factors")
	}
}

func TestFactorGraphCopy(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	cpy := fg.Copy()

	if len(cpy.GetVariables()) != len(fg.GetVariables()) {
		t.Errorf("copy has %d variables, expected %d", len(cpy.GetVariables()), len(fg.GetVariables()))
	}
	if len(cpy.GetFactors()) != len(fg.GetFactors()) {
		t.Errorf("copy has %d factors, expected %d", len(cpy.GetFactors()), len(fg.GetFactors()))
	}

	// Modify copy and ensure original is unaffected.
	cpy.RemoveFactors()
	if len(fg.GetFactors()) != 2 {
		t.Error("original was affected by copy modification")
	}
}

func TestFactorGraphCopyCheckModel(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	cpy := fg.Copy()
	if err := cpy.CheckModel(); err != nil {
		t.Fatalf("copied factor graph CheckModel: %v", err)
	}
}

func TestFactorGraphGetPointMassMessage(t *testing.T) {
	fg := buildSimpleFactorGraph(t)

	msg, err := fg.GetPointMassMessage("A", 0)
	if err != nil {
		t.Fatalf("GetPointMassMessage: %v", err)
	}

	vals := msg.Values().Data()
	if len(vals) != 2 {
		t.Fatalf("expected 2 values, got %d", len(vals))
	}
	if vals[0] != 1.0 || vals[1] != 0.0 {
		t.Errorf("expected [1.0, 0.0], got %v", vals)
	}
}

func TestFactorGraphGetPointMassMessageState1(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	msg, err := fg.GetPointMassMessage("A", 1)
	if err != nil {
		t.Fatalf("GetPointMassMessage: %v", err)
	}
	vals := msg.Values().Data()
	if vals[0] != 0.0 || vals[1] != 1.0 {
		t.Errorf("expected [0.0, 1.0], got %v", vals)
	}
}

func TestFactorGraphGetPointMassMessageUnknown(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	_, err := fg.GetPointMassMessage("Z", 0)
	if err == nil {
		t.Error("expected error for unknown variable")
	}
}

func TestFactorGraphGetPointMassMessageInvalidState(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	_, err := fg.GetPointMassMessage("A", 5)
	if err == nil {
		t.Error("expected error for invalid state")
	}
	_, err = fg.GetPointMassMessage("A", -1)
	if err == nil {
		t.Error("expected error for negative state")
	}
}

func TestFactorGraphGetUniformMessage(t *testing.T) {
	fg := buildSimpleFactorGraph(t)

	msg, err := fg.GetUniformMessage("A")
	if err != nil {
		t.Fatalf("GetUniformMessage: %v", err)
	}

	vals := msg.Values().Data()
	if len(vals) != 2 {
		t.Fatalf("expected 2 values, got %d", len(vals))
	}
	if math.Abs(vals[0]-0.5) > 1e-10 || math.Abs(vals[1]-0.5) > 1e-10 {
		t.Errorf("expected [0.5, 0.5], got %v", vals)
	}
}

func TestFactorGraphGetUniformMessageTernary(t *testing.T) {
	fg := NewFactorGraph()
	_ = fg.AddVariable("X", 3)

	msg, err := fg.GetUniformMessage("X")
	if err != nil {
		t.Fatalf("GetUniformMessage: %v", err)
	}

	vals := msg.Values().Data()
	for _, v := range vals {
		if math.Abs(v-1.0/3.0) > 1e-10 {
			t.Errorf("expected 1/3, got %f", v)
		}
	}
}

func TestFactorGraphGetUniformMessageUnknown(t *testing.T) {
	fg := buildSimpleFactorGraph(t)
	_, err := fg.GetUniformMessage("Z")
	if err == nil {
		t.Error("expected error for unknown variable")
	}
}
