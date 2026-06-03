//go:build unit

package models

import (
	"math"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

func mustNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

// ---------------------------------------------------------------------------
// BayesianNetwork: SetStates error
// ---------------------------------------------------------------------------

func TestBN_SetStates_UnknownVar(t *testing.T) {
	bn := NewBayesianNetwork()
	err := bn.SetStates("X", []string{"a", "b"})
	if err == nil {
		t.Error("expected error for unknown variable")
	}
}

// ---------------------------------------------------------------------------
// DiscreteMarkovNetwork: AddFactor / CheckModel
// ---------------------------------------------------------------------------

func TestDiscreteMarkovNetwork_AddFactor_Nil(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	err := dmn.AddFactor(nil)
	if err == nil {
		t.Error("expected error for nil factor")
	}
}

func TestDiscreteMarkovNetwork_AddFactor_NaN(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	f, err := factors.NewDiscreteFactor([]string{"A"}, []int{2}, []float64{math.NaN(), 0.5})
	if err != nil {
		t.Fatal(err)
	}
	err = dmn.AddFactor(f)
	if err == nil {
		t.Error("expected error for NaN values")
	}
}

func TestDiscreteMarkovNetwork_CheckModel_NoCPDs(t *testing.T) {
	dmn := NewDiscreteMarkovNetwork()
	dmn.AddNode("A")
	err := dmn.CheckModel()
	if err == nil {
		t.Error("expected error for no factors")
	}
}

// ---------------------------------------------------------------------------
// DiscreteBayesianNetwork: AddCPD / CheckModel errors
// ---------------------------------------------------------------------------

func TestDiscreteBN_AddCPD_Nil(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	err := dbn.AddCPD(nil)
	if err == nil {
		t.Error("expected error for nil CPD")
	}
}

func TestDiscreteBN_CheckModel_NoCPDs(t *testing.T) {
	dbn := NewDiscreteBayesianNetwork()
	mustNoErr(t, dbn.AddNode("A"))
	mustNoErr(t, dbn.SetStates("A", []string{"a0", "a1"}))
	err := dbn.CheckModel()
	if err == nil {
		t.Error("expected error for missing CPDs")
	}
}

// ---------------------------------------------------------------------------
// DynamicBayesianNetwork: AddNode / AddEdge error paths
// ---------------------------------------------------------------------------

func TestDynamicBN_AddNode_Duplicate(t *testing.T) {
	dbn := NewDynamicBayesianNetwork()
	mustNoErr(t, dbn.AddNode("A"))
	err := dbn.AddNode("A")
	if err == nil {
		t.Error("expected error for duplicate node")
	}
}

func TestDynamicBN_AddEdge_MissingNode(t *testing.T) {
	dbn := NewDynamicBayesianNetwork()
	mustNoErr(t, dbn.AddNode("A"))
	err := dbn.AddEdge("A", "Z")
	if err == nil {
		t.Error("expected error for missing node")
	}
}

// ---------------------------------------------------------------------------
// FunctionalBN: CheckModel
// ---------------------------------------------------------------------------

func TestFunctionalBN_CheckModel_NoFunctions(t *testing.T) {
	fbn := NewFunctionalBayesianNetwork()
	mustNoErr(t, fbn.AddNode("A"))
	err := fbn.CheckModel()
	if err == nil {
		t.Error("expected error for missing functions")
	}
}

// ---------------------------------------------------------------------------
// MarkovNetwork: CheckModel
// ---------------------------------------------------------------------------

func TestMarkovNetwork_CheckModel_Coverage(t *testing.T) {
	mn := NewMarkovNetwork()
	mn.AddNode("A")
	err := mn.CheckModel()
	if err == nil {
		t.Error("expected error for no factors")
	}
}

// ---------------------------------------------------------------------------
// NaiveBayes: AddEdgesFrom
// ---------------------------------------------------------------------------

func TestNaiveBayes_AddEdgesFrom_MissingFeature(t *testing.T) {
	nb, err := NewNaiveBayes("class", []string{"f1", "f2"})
	if err != nil {
		t.Fatal(err)
	}
	mustNoErr(t, nb.SetStates("class", []string{"c0", "c1"}))
	mustNoErr(t, nb.SetStates("f1", []string{"a", "b"}))
	mustNoErr(t, nb.SetStates("f2", []string{"x", "y"}))
	// AddEdgesFrom with a non-existent class variable
	err = nb.AddEdgesFrom("nonexistent", []string{"f1"})
	if err == nil {
		t.Error("expected error for non-existent class variable")
	}
}

// ---------------------------------------------------------------------------
// FactorGraph: CheckModel
// ---------------------------------------------------------------------------

func TestFactorGraph_CheckModel_NoFactors(t *testing.T) {
	fg := NewFactorGraph()
	mustNoErr(t, fg.AddVariable("A", 2))
	err := fg.CheckModel()
	if err == nil {
		t.Error("expected error for no factor nodes")
	}
}

// ---------------------------------------------------------------------------
// MarkovChain: basic error paths
// ---------------------------------------------------------------------------

func TestMarkovChain_Invalid(t *testing.T) {
	_, err := NewMarkovChain(nil, nil)
	if err == nil {
		t.Error("expected error for nil transition matrix")
	}
}

func TestMarkovChain_NonSquare(t *testing.T) {
	_, err := NewMarkovChain([][]float64{{0.5, 0.5}}, []string{"a", "b"})
	if err == nil {
		t.Error("expected error for non-square transition matrix")
	}
}
