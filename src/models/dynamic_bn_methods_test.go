//go:build unit

package models

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
)

func TestDBNAddNode(t *testing.T) {
	dbn := NewDynamicBayesianNetwork()
	if err := dbn.AddNode("A"); err != nil {
		t.Fatalf("AddNode: %v", err)
	}

	initNodes := dbn.Initial().Nodes()
	transNodes := dbn.Transition().Nodes()
	if len(initNodes) != 1 || initNodes[0] != "A" {
		t.Errorf("expected initial nodes [A], got %v", initNodes)
	}
	if len(transNodes) != 1 || transNodes[0] != "A" {
		t.Errorf("expected transition nodes [A], got %v", transNodes)
	}
}

func TestDBNAddNodeDuplicate(t *testing.T) {
	dbn := NewDynamicBayesianNetwork()
	_ = dbn.AddNode("A")
	if err := dbn.AddNode("A"); err == nil {
		t.Error("expected error for duplicate node")
	}
}

func TestDBNAddEdge(t *testing.T) {
	dbn := NewDynamicBayesianNetwork()
	_ = dbn.AddNode("A")
	_ = dbn.AddNode("B")
	if err := dbn.AddEdge("A", "B"); err != nil {
		t.Fatalf("AddEdge: %v", err)
	}

	initEdges := dbn.Initial().Edges()
	transEdges := dbn.Transition().Edges()
	if len(initEdges) != 1 {
		t.Errorf("expected 1 initial edge, got %d", len(initEdges))
	}
	if len(transEdges) != 1 {
		t.Errorf("expected 1 transition edge, got %d", len(transEdges))
	}
}

func TestDBNGetIntraEdges(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	intra := dbn.GetIntraEdges()
	if len(intra) != 1 {
		t.Errorf("expected 1 intra edge, got %d", len(intra))
	}
}

func TestDBNGetInterEdges(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	inter := dbn.GetInterEdges()
	if len(inter) != 1 {
		t.Errorf("expected 1 inter edge, got %d", len(inter))
	}
}

func TestDBNGetSliceNodes(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)

	nodes0, err := dbn.GetSliceNodes(0)
	if err != nil {
		t.Fatalf("GetSliceNodes(0): %v", err)
	}
	if len(nodes0) != 2 {
		t.Errorf("expected 2 nodes in slice 0, got %d", len(nodes0))
	}

	nodes1, err := dbn.GetSliceNodes(1)
	if err != nil {
		t.Fatalf("GetSliceNodes(1): %v", err)
	}
	if len(nodes1) != 2 {
		t.Errorf("expected 2 nodes in slice 1, got %d", len(nodes1))
	}

	_, err = dbn.GetSliceNodes(2)
	if err == nil {
		t.Error("expected error for invalid slice")
	}
}

func TestDBNGetCPDs(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	cpds := dbn.GetCPDs()
	if len(cpds) != 2 {
		t.Errorf("expected 2 CPDs, got %d", len(cpds))
	}
}

func TestDBNRemoveCPDs(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	dbn.RemoveCPDs("X")

	if dbn.Initial().GetCPD("X") != nil {
		t.Error("expected nil initial CPD for X after RemoveCPDs")
	}
	if dbn.Transition().GetCPD("X") != nil {
		t.Error("expected nil transition CPD for X after RemoveCPDs")
	}
}

func TestDBNInitializeInitialState(t *testing.T) {
	dbn := NewDynamicBayesianNetwork()
	_ = dbn.Initial().AddNode("X")
	_ = dbn.Transition().AddNode("X")

	// Add a placeholder CPD first so the node has a CPD.
	cpdX, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	_ = dbn.AddInitialCPD(cpdX)
	_ = dbn.AddTransitionCPD(cpdX)

	err := dbn.InitializeInitialState(map[string][]float64{
		"X": {0.3, 0.7},
	})
	if err != nil {
		t.Fatalf("InitializeInitialState: %v", err)
	}

	cpd := dbn.Initial().GetCPD("X")
	if cpd == nil {
		t.Fatal("expected CPD for X")
	}
	f := cpd.ToFactor()
	vals := f.Values().Data()
	if vals[0] != 0.3 || vals[1] != 0.7 {
		t.Errorf("expected [0.3, 0.7], got %v", vals)
	}
}

func TestDBNInitializeInitialStateEmpty(t *testing.T) {
	dbn := NewDynamicBayesianNetwork()
	_ = dbn.Initial().AddNode("X")

	err := dbn.InitializeInitialState(map[string][]float64{
		"X": {},
	})
	if err == nil {
		t.Error("expected error for empty distribution")
	}
}

func TestDBNMoralize(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	moral := dbn.Moralize()
	if moral == nil {
		t.Fatal("Moralize returned nil")
	}

	nodes := moral.Nodes()
	if len(nodes) != 2 {
		t.Errorf("expected 2 nodes in moral graph, got %d", len(nodes))
	}
}

func TestDBNGetMarkovBlanket(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	blanket, err := dbn.GetMarkovBlanket("X")
	if err != nil {
		t.Fatalf("GetMarkovBlanket: %v", err)
	}
	// X -> Y, so blanket of X is {Y}
	if len(blanket) != 1 || blanket[0] != "Y" {
		t.Errorf("expected blanket [Y], got %v", blanket)
	}
}

func TestDBNGetMarkovBlanketUnknown(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	_, err := dbn.GetMarkovBlanket("Z")
	if err == nil {
		t.Error("expected error for unknown node")
	}
}

func TestDBNGetConstantBN(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)

	bn0, err := dbn.GetConstantBN(0)
	if err != nil {
		t.Fatalf("GetConstantBN(0): %v", err)
	}
	if err := bn0.CheckModel(); err != nil {
		t.Fatalf("slice 0 CheckModel: %v", err)
	}

	bn1, err := dbn.GetConstantBN(1)
	if err != nil {
		t.Fatalf("GetConstantBN(1): %v", err)
	}
	if err := bn1.CheckModel(); err != nil {
		t.Fatalf("slice 1 CheckModel: %v", err)
	}

	// Modify the copy and ensure original is unaffected.
	bn0.RemoveCPD("X")
	if dbn.Initial().GetCPD("X") == nil {
		t.Error("original affected by copy modification")
	}

	_, err = dbn.GetConstantBN(5)
	if err == nil {
		t.Error("expected error for invalid slice")
	}
}

func TestDBNFit(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)

	// Create simple time series data.
	data := buildSimpleTimeSeriesData()

	err := dbn.Fit(data)
	if err != nil {
		t.Fatalf("Fit: %v", err)
	}
}

func TestDBNFitNilData(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	if err := dbn.Fit(nil); err == nil {
		t.Error("expected error for nil data")
	}
}

func TestDBNActiveTrailNodes(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	active := dbn.ActiveTrailNodes([]string{"X"}, nil)
	if !active["X"] {
		t.Error("X should be in active trail")
	}
	if !active["Y"] {
		t.Error("Y should be reachable from X")
	}
}

func TestDBNActiveTrailNodesObserved(t *testing.T) {
	// Build a network with X -> Y -> Z where observing Y blocks X-Z trail.
	dbn := NewDynamicBayesianNetwork()
	for _, n := range []string{"X", "Y", "Z"} {
		_ = dbn.Initial().AddNode(n)
		_ = dbn.Transition().AddNode(n)
	}
	_ = dbn.Initial().AddEdge("X", "Y")
	_ = dbn.Initial().AddEdge("Y", "Z")
	_ = dbn.Transition().AddEdge("X", "Y")
	_ = dbn.Transition().AddEdge("Y", "Z")

	// Add CPDs.
	cpdX, _ := factors.NewTabularCPD("X", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	cpdY, _ := factors.NewTabularCPD("Y", 2, [][]float64{{0.8, 0.2}, {0.2, 0.8}}, []string{"X"}, []int{2})
	cpdZ, _ := factors.NewTabularCPD("Z", 2, [][]float64{{0.9, 0.1}, {0.1, 0.9}}, []string{"Y"}, []int{2})
	_ = dbn.AddInitialCPD(cpdX)
	_ = dbn.AddInitialCPD(cpdY)
	_ = dbn.AddInitialCPD(cpdZ)
	_ = dbn.AddTransitionCPD(cpdX)
	_ = dbn.AddTransitionCPD(cpdY)
	_ = dbn.AddTransitionCPD(cpdZ)

	// Observing Y should block X -> Z trail.
	active := dbn.ActiveTrailNodes([]string{"X"}, []string{"Y"})
	if active["Z"] {
		t.Error("Z should NOT be reachable from X when Y is observed (chain)")
	}
}

func TestDBNSimulate(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	df, err := dbn.Simulate(10, 42)
	if err != nil {
		t.Fatalf("Simulate: %v", err)
	}

	if df.Len() != 10 {
		t.Errorf("expected 10 rows, got %d", df.Len())
	}

	xVals := df.Column("X").Values()
	yVals := df.Column("Y").Values()
	if len(xVals) != 10 || len(yVals) != 10 {
		t.Errorf("expected 10 values per column")
	}

	// All values should be valid states (0 or 1).
	for i := 0; i < 10; i++ {
		x := toInt(xVals[i])
		y := toInt(yVals[i])
		if x < 0 || x > 1 {
			t.Errorf("X[%d] = %d, out of range", i, x)
		}
		if y < 0 || y > 1 {
			t.Errorf("Y[%d] = %d, out of range", i, y)
		}
	}
}

func TestDBNSimulateDeterministic(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	df1, _ := dbn.Simulate(20, 42)
	df2, _ := dbn.Simulate(20, 42)

	x1 := df1.Column("X").Values()
	x2 := df2.Column("X").Values()
	for i := range x1 {
		if toInt(x1[i]) != toInt(x2[i]) {
			t.Errorf("X[%d] differs: %v vs %v", i, x1[i], x2[i])
		}
	}
}

func TestDBNSimulateInvalidModel(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	dbn.Initial().RemoveCPD("X")
	_, err := dbn.Simulate(10, 42)
	if err == nil {
		t.Error("expected error for invalid model")
	}
}

func TestDBNSimulateInvalidN(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)
	_, err := dbn.Simulate(0, 42)
	if err == nil {
		t.Error("expected error for n=0")
	}
	_, err = dbn.Simulate(-1, 42)
	if err == nil {
		t.Error("expected error for negative n")
	}
}

func TestDBNStates(t *testing.T) {
	dbn := buildSimpleDynamicBN(t)

	// Set states on initial network.
	_ = dbn.Initial().SetStates("X", []string{"off", "on"})
	_ = dbn.Initial().SetStates("Y", []string{"low", "high"})

	states := dbn.States()
	if len(states) != 2 {
		t.Errorf("expected 2 variables with states, got %d", len(states))
	}
	if states["X"][0] != "off" || states["X"][1] != "on" {
		t.Errorf("expected X states [off, on], got %v", states["X"])
	}
}

func TestDBNStatesEmpty(t *testing.T) {
	dbn := NewDynamicBayesianNetwork()
	states := dbn.States()
	if len(states) != 0 {
		t.Errorf("expected 0 states, got %d", len(states))
	}
}

// buildSimpleTimeSeriesData creates a simple DataFrame with columns X and Y.
func buildSimpleTimeSeriesData() *tabgo.DataFrame {
	return tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{0, 1, 0, 1, 0}),
		"Y": tabgo.NewSeries("Y", []any{0, 1, 1, 0, 0}),
	})
}
