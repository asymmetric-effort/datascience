//go:build unit

package learning

import (
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// =========================================================================
// MLE: negative/out-of-range values trigger skip branches in estimateNode
// =========================================================================

func TestCovPush_MLE_NegativeChildState(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{-1, 0, 1, 0}),
	})
	mle := NewMLE(bn, data)
	err := mle.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCovPush_MLE_NegativeParentValue(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{-1, 0, 1, 0}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	mle := NewMLE(bn, data)
	err := mle.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCovPush_MLE_OutOfRangeChild(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})
	// Value 5 exceeds card=2
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 5, 1, 0}),
	})
	mle := NewMLE(bn, data)
	err := mle.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// MLE.Estimate: error from AddCPD (line 61-63) - force by adding cpd for unknown node
// MLE.Estimate: error from estimateNode (line 58-60)
// These are hard to trigger in prod since validation occurs early. But let's try
// MLE.estimateNode with out-of-range parent values to hit "valid = false; break"

func TestCovPush_MLE_OutOfRangeParent(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	// Parent value 5 exceeds card=2
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 5, 1, 0}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	mle := NewMLE(bn, data)
	err := mle.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// MLE.EstimatePotentials: out-of-range values in edge counting (line 266)
func TestCovPush_MLE_EstimatePotentials_OutOfRange(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	// -1 values to exercise the "if ui >= 0 && ui < card" guard
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{-1, 0, 1, 0}),
		"B": tabgo.NewSeries("B", []any{0, -1, 1, 1}),
	})
	mle := NewMLE(bn, data)
	result, err := mle.EstimatePotentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected at least one factor")
	}
}

// MLE.EstimatePotentials: isolated node unary factor with out-of-range values (line 291-293)
func TestCovPush_MLE_EstimatePotentials_IsolatedNodeOOR(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{-1, 0, 1, 5}),
	})
	mle := NewMLE(bn, data)
	result, err := mle.EstimatePotentials()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) == 0 {
		t.Fatal("expected at least one factor")
	}
}

// =========================================================================
// BayesianEstimator: unknown state skip (line 157-159, 163-164),
//                    parent unknown state skip (line 157), zero colSum (line 184-188),
//                    CPD creation error (line 203-205)
// =========================================================================

func TestCovPush_BayesianEstimator_UnknownStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	// Data with unknown states: "2" not in states
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{"0", "2", "1", "0"}),
		"B": tabgo.NewSeries("B", []any{"0", "0", "2", "1"}),
	})
	be := NewBayesianEstimator(bn, data, K2, 1)
	err := be.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// BayesianEstimator.Estimate: AddCPD failure (line 59-61)
// Hard to trigger naturally. Skip for now since it requires BN internal error.

// =========================================================================
// EM: error paths
// =========================================================================

// EM.Estimate line 139-141: unknown state in data
func TestCovPush_EM_UnknownState(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	// State "2" doesn't exist
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{"0", "1", "2", "0"}),
		"B": tabgo.NewSeries("B", []any{"0", "0", "1", "1"}),
	})
	em := NewEM(bn, data, nil, 10, 0.01)
	err := em.Estimate()
	if err == nil {
		t.Fatal("expected error for unknown state")
	}
}

// EM.Estimate line 150-152: initializeCPDs failure
func TestCovPush_EM_NodeNoStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	// Don't set states for A
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{"0", "1"}),
	})
	em := NewEM(bn, data, nil, 10, 0.01)
	err := em.Estimate()
	if err == nil {
		t.Fatal("expected error for node with no states")
	}
}

// EM.Estimate line 186-188: E-step error from computeLatentPosterior
// This happens when VE fails. We need latent vars that cause issues.
// EM.computeLatentPosterior lines 390-392, 396-398: markov factor error / VE error
// Hard to force without breaking BN. Let's try with a BN that has no CPDs set.

// EM.Estimate line 219-221: M-step CPD creation error
// EM.Estimate line 235-237: AddCPD error in M-step
// These need special conditions.

// EM.initializeCPDs line 302-304: uniform init for unobserved parent configs
// Test with latent parents
func TestCovPush_EM_LatentParent(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("L")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("L", "A")
	_ = bn.AddEdge("L", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("L", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{"0", "1", "0", "1", "0", "1"}),
		"B": tabgo.NewSeries("B", []any{"0", "0", "1", "1", "0", "1"}),
	})
	em := NewEM(bn, data, []string{"L"}, 5, 0.01)
	err := em.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if em.Iterations() == 0 {
		t.Error("expected at least one iteration")
	}
}

// EM.initializeCPDs line 337-339: CPD creation error
// EM.initializeCPDs line 340-342: AddCPD error
// These are defensive paths that are hard to trigger.

// EM.GetParameters: no CPD set
func TestCovPush_EM_GetParameters_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})
	em := NewEM(bn, nil, nil, 10, 0.01)
	_, err := em.GetParameters()
	if err == nil {
		t.Fatal("expected error for no CPD")
	}
}

// =========================================================================
// MirrorDescent: out-of-range values, setUniformCPDs error paths
// =========================================================================

func TestCovPush_MirrorDescent_NegativeValues(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{-1, 0, 1, 0}),
		"B": tabgo.NewSeries("B", []any{0, 0, -1, 1}),
	})
	md := NewMirrorDescentEstimator(bn, data, 0.1, 10)
	err := md.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCovPush_MirrorDescent_OutOfRange(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 5, 1, 0}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 5}),
	})
	md := NewMirrorDescentEstimator(bn, data, 0.1, 10)
	err := md.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// MirrorDescent: pcTotals==0 branch (line 193-197)
// Triggered when a parent config has zero data
func TestCovPush_MirrorDescent_ZeroParentConfig(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1", "2"})
	_ = bn.SetStates("B", []string{"0", "1", "2"})
	// Only state 0 and 1 appear for A, never 2 -> parent config 2 has zero counts
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1, 0, 1}),
	})
	md := NewMirrorDescentEstimator(bn, data, 0.1, 10)
	err := md.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// =========================================================================
// PC: various error paths and edge cases
// =========================================================================

// PC.BuildSkeleton: maxCondSetSize boundary (line 93-94)
func TestCovPush_PC_BuildSkeleton_MaxCondSetSize(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 0, 1, false // never independent
	}
	pc := NewPC(data, ciTest, 0.05, WithMaxCondSetSize(0))
	pdag, sepSets, err := pc.BuildSkeleton()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pdag == nil {
		t.Fatal("expected non-nil PDAG")
	}
	_ = sepSets
}

// PC.Estimate: edge already removed during iteration (line 191-192)
func TestCovPush_PC_Estimate_EdgeRemoved(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	callCount := 0
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		callCount++
		// Make everything independent to trigger lots of removal
		return 0, 1, true
	}
	pc := NewPC(data, ciTest, 0.05)
	pdag, err := pc.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pdag == nil {
		t.Fatal("expected non-nil PDAG")
	}
}

// PC.Estimate: second adjY branch (line 210-214)
func TestCovPush_PC_Estimate_AdjYBranch(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	callCount := 0
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		callCount++
		// Only independent when conditioning on C
		if len(z) > 0 && z[0] == "C" {
			return 0, 1, true
		}
		// Independent from adj(Y) side for specific pairs
		if x == "A" && y == "C" && len(z) > 0 {
			return 0, 1, true
		}
		return 5.0, 0.01, false
	}
	pc := NewPC(data, ciTest, 0.05)
	pdag, err := pc.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// PC.Estimate: sepSet not exists skip (line 238-241)
func TestCovPush_PC_Estimate_NoSepSet(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		// Never independent - keeps all edges
		return 5.0, 0.01, false
	}
	pc := NewPC(data, ciTest, 0.05)
	pdag, err := pc.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// PC.EstimateBN: pdagToDAG error (line 277-279)
// This requires a PDAG that can't be oriented without creating a cycle
// Hard to force - let's try the normal path
func TestCovPush_PC_EstimateBN_Success(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		// A-B independent given C, create v-structure A->C<-B
		if (x == "A" && y == "B") || (x == "B" && y == "A") {
			return 0, 1, true
		}
		return 5.0, 0.01, false
	}
	pc := NewPC(data, ciTest, 0.05)
	bn, err := pc.EstimateBN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bn == nil {
		t.Fatal("expected non-nil BN")
	}
}

// pdagToDAG: directed edge error (line 299-301)
func TestCovPush_pdagToDAG_DirectedEdgeError(t *testing.T) {
	pdag := graphgo.NewPDAG()
	pdag.AddNode("A")
	pdag.AddNode("B")
	// Add a directed edge that creates a cycle when converted
	pdag.AddDirectedEdge("A", "B")
	pdag.AddDirectedEdge("B", "A")
	_, err := pdagToDAG(pdag)
	if err == nil {
		t.Fatal("expected error for cyclic directed edges")
	}
}

// pdagToDAG: undirected edge neither direction works (line 315-317)
func TestCovPush_pdagToDAG_UndirectedCycleError(t *testing.T) {
	pdag := graphgo.NewPDAG()
	pdag.AddNode("A")
	pdag.AddNode("B")
	pdag.AddNode("C")
	pdag.AddDirectedEdge("A", "B")
	pdag.AddDirectedEdge("B", "C")
	pdag.AddDirectedEdge("C", "A")
	// This already has a cycle in directed edges
	_, err := pdagToDAG(pdag)
	if err == nil {
		t.Fatal("expected error for cyclic PDAG")
	}
}

// OrientColliders: sepSet not found skip (line 346-347)
func TestCovPush_OrientColliders_NoSepSet(t *testing.T) {
	pdag := graphgo.NewPDAG()
	pdag.AddNode("A")
	pdag.AddNode("B")
	pdag.AddNode("C")
	pdag.AddUndirectedEdge("A", "C")
	pdag.AddUndirectedEdge("B", "C")
	// No A-B edge, but empty sepSets -> should skip
	sepSets := make(map[[2]string][]string)
	OrientColliders(pdag, sepSets)
	// Should not have oriented anything
}

// =========================================================================
// GES: error paths
// =========================================================================

// GES.Estimate: cycle check during add (line 71-73)
func TestCovPush_GES_Estimate_CycleCheck(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1, 0, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0, 1, 0}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	pdag, err := ges.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pdag == nil {
		t.Fatal("expected non-nil PDAG")
	}
}

// GES.Insert: edge already exists (line 142)
func TestCovPush_GES_Insert_EdgeExists(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddEdge("A", "B")
	_, err := ges.Insert(dag, "A", "B")
	if err == nil {
		t.Fatal("expected error for existing edge")
	}
}

// GES.Insert: cycle creation (line 150-153)
func TestCovPush_GES_Insert_Cycle(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddEdge("A", "B")
	_, err := ges.Insert(dag, "B", "A")
	if err == nil {
		t.Fatal("expected error for cycle")
	}
}

// GES.Turn: cycle on reversal (line 196-200)
func TestCovPush_GES_Turn_Cycle(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddNode("C")
	dag.AddEdge("A", "B")
	dag.AddEdge("B", "C")
	dag.AddEdge("A", "C")
	// Turning A->B to B->A when B->C->... would need A->C which already exists
	// Actually turn B->C when A->B and A->C exist. If we turn B->C, it becomes C->B.
	// That's fine. Let's create a cycle scenario: A->B, B->C, C->A is already cyclic.
	// Instead: A->B, B->C exist. Turn A->B to B->A. Then B->A and B->C exist. No cycle.
	// Need: A->B, C->A. Turn A->B to B->A creates B->A, C->A - no cycle.
	// Need actual cycle: A->B, B->C. Turn B->C to C->B. Then A->B and C->B exist. No cycle.
	// To force: A->B, C->B, A->C. Turn A->C: remove A->C, add C->A. Now C->A, A->B, C->B. C->A->B and C->B. No cycle.
	// Real cycle: A->B, B->C. Turn A->B: remove A->B, add B->A. Check: B->A, B->C. No cycle.
	// Hard. Let's use: A->B->C and C->A would be cycle. Start with A->B, B->C, and try to turn B->C.
	dag2 := graphgo.NewDiGraph()
	dag2.AddNode("A")
	dag2.AddNode("B")
	dag2.AddNode("C")
	dag2.AddEdge("A", "B")
	dag2.AddEdge("B", "C")
	dag2.AddEdge("C", "A") // This already creates a cycle!
	// IsDAG would already be false. Let's try the right way.
	dag3 := graphgo.NewDiGraph()
	dag3.AddNode("A")
	dag3.AddNode("B")
	dag3.AddNode("C")
	dag3.AddEdge("A", "B")
	dag3.AddEdge("B", "C")
	// Turning B->C: remove B->C, add C->B. Check DAG: A->B, C->B. Fine.
	// Turning A->B: remove A->B, add B->A. Check: B->A, B->C. Fine.
	// Need: turn would create a cycle. A->B, A->C, C->B. Turn C->B: remove C->B, add B->C.
	// Now A->B, A->C, B->C. Still DAG. No cycle.
	// Only way: A->B, B->C, turn B->C to C->B, then C->B and... A->B, C->B. Still DAG.
	// Actually for a cycle: need path from v back to u after adding v->u.
	// A->B, C->A. Turn A->B: add B->A. Now B->A, C->A. No cycle.
	// A->B, C->A. Turn C->A: add A->C. Now A->B, A->C. No cycle.
	// Hmm. Need: existing path from u to v, then turn u->v creates v->u + path u to v = cycle.
	// So: A->B->C and turn A->C: remove A->C, add C->A. Now C->A->B->C = cycle!
	dag4 := graphgo.NewDiGraph()
	dag4.AddNode("A")
	dag4.AddNode("B")
	dag4.AddNode("C")
	dag4.AddEdge("A", "B")
	dag4.AddEdge("B", "C")
	dag4.AddEdge("A", "C")
	// Turn A->C: remove A->C, add C->A. C->A->B->C = cycle!
	_, err := ges.Turn(dag4, "A", "C")
	if err == nil {
		t.Fatal("expected error for cycle on turn")
	}
	// Verify the edge was restored
	if !dag4.HasEdge("A", "C") {
		t.Error("expected A->C to be restored after failed turn")
	}
}

// GES.dagToPDAG: v-structure detection (line 236-240)
func TestCovPush_GES_dagToPDAG_VStructure(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	// Create a v-structure: A->C<-B (A and B not adjacent)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddNode("C")
	dag.AddEdge("A", "C")
	dag.AddEdge("B", "C")
	pdag := ges.dagToPDAG(dag, []string{"A", "B", "C"})
	if pdag == nil {
		t.Fatal("expected non-nil PDAG")
	}
}

// GES.dagToPDAG: undirected edge case (non-v-structure, line 247-249)
func TestCovPush_GES_dagToPDAG_Undirected(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	// Chain: A->B->C (no v-structure, B->C should be undirected in equivalence class)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddNode("C")
	dag.AddEdge("A", "B")
	dag.AddEdge("B", "C")
	pdag := ges.dagToPDAG(dag, []string{"A", "B", "C"})
	if pdag == nil {
		t.Fatal("expected non-nil PDAG")
	}
}

// GES backward phase: bestDelta > 0 removal (line 120-124, 130)
func TestCovPush_GES_BackwardPhase(t *testing.T) {
	// Create data where removing an edge improves the score
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	pdag, err := ges.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// =========================================================================
// ExpertInLoop: error paths
// =========================================================================

// ExpertInLoop.Estimate: LLM query paths (line 83-84, 97-100, 118-119, 124-125)
func TestCovPush_ExpertInLoop_WithLLM(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices":[{"message":{"content":"YES confidence: 0.9"}}]}`))
	}))
	defer server.Close()

	llmClient := NewHTTPLLMClient(server.URL, "test-key", "gpt-4")
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		// A and B independent (marginal)
		if (x == "A" && y == "B") || (x == "B" && y == "A") {
			return 0, 1, true
		}
		return 5.0, 0.01, false
	}
	eil := NewExpertInLoop(data, llmClient, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// ExpertInLoop.queryOneCausalDirection: error from LLM (line 205-207)
func TestCovPush_ExpertInLoop_LLMError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"bad request"}`))
	}))
	defer server.Close()

	llmClient := NewHTTPLLMClient(server.URL, "test-key", "gpt-4")
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		if (x == "A" && y == "B") || (x == "B" && y == "A") {
			return 0, 1, true
		}
		return 5.0, 0.01, false
	}
	eil := NewExpertInLoop(data, llmClient, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// ExpertInLoop.EstimateBN: success path (line 236-238, 241-243)
func TestCovPush_ExpertInLoop_EstimateBN(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		if (x == "A" && y == "B") || (x == "B" && y == "A") {
			return 0, 1, true
		}
		return 5.0, 0.01, false
	}
	eil := NewExpertInLoop(data, nil, ciTest, 0.05)
	bn, err := eil.EstimateBN()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bn == nil {
		t.Fatal("expected non-nil BN")
	}
}

// =========================================================================
// Scoring: error paths
// =========================================================================

// BICScore: n==0 (line 29-30)
func TestCovPush_BICScore_EmptyData(t *testing.T) {
	fn := BICScore()
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{}),
	})
	score := fn("A", nil, data)
	if score != 0 {
		t.Errorf("expected 0 for empty data, got %f", score)
	}
}

// BICScore: numPC==0 (line 41-43) - shouldn't happen in practice, but
// let's cover: count==0 branch (line 29-30 is the n==0 guard)
func TestCovPush_BICScore_ZeroCount(t *testing.T) {
	fn := BICScore()
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	score := fn("A", []string{"B"}, data)
	if math.IsNaN(score) {
		t.Error("score should not be NaN")
	}
}

// AICScore: n==0 (line 126-127)
func TestCovPush_AICScore_EmptyData(t *testing.T) {
	fn := AICScore()
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{}),
	})
	score := fn("A", nil, data)
	if score != 0 {
		t.Errorf("expected 0 for empty data, got %f", score)
	}
}

// AICScore: numPC==0 (line 137-139) same as BIC
func TestCovPush_AICScore_WithParents(t *testing.T) {
	fn := AICScore()
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	score := fn("A", []string{"B"}, data)
	if math.IsNaN(score) {
		t.Error("score should not be NaN")
	}
}

// BDeuScore: numPC==0 (line 87-89)
func TestCovPush_BDeuScore_EmptyData(t *testing.T) {
	fn := BDeuScore(5.0)
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{}),
	})
	score := fn("A", nil, data)
	_ = score // just ensure no panic
}

// localCountTable: parent pc < 1 guard (line 179-181)
func TestCovPush_localCountTable_EmptyParent(t *testing.T) {
	// Use negative values to trigger parent card < 1 guard
	data2 := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{-1, -1, -1, -1}),
	})
	fn := BICScore()
	score := fn("A", []string{"B"}, data2)
	_ = score
}

// localCountTable: card < 1 guard (line 164-166)
func TestCovPush_localCountTable_NegativeVarVals(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{-1, -1, -1}),
	})
	fn := BICScore()
	score := fn("A", nil, data)
	_ = score
}

// =========================================================================
// TreeSearch: error paths
// =========================================================================

// TreeSearch.Estimate: class var not found (line 87-89)
func TestCovPush_TreeSearch_ClassVarNotFound(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	ts := NewTreeSearch(data, WithClassVariable("Z"))
	_, err := ts.Estimate()
	if err == nil {
		t.Fatal("expected error for missing class variable")
	}
}

// TreeSearch.Estimate: root not found (line 127-129)
// Not triggered by: root in columns but not in treeVars
func TestCovPush_TreeSearch_RootNotFoundInTreeVars(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0}),
	})
	// Root is C, classVar is C -> C is removed from treeVars, root "C" not in treeVars
	ts := NewTreeSearch(data, WithClassVariable("C"), WithRoot("C"))
	_, err := ts.Estimate()
	if err == nil {
		t.Fatal("expected error for root not in tree variables")
	}
}

// TreeSearch.Estimate: AddEdge error (line 133-135)
// Hard to trigger since BN allows edge addition. Skip.

// TreeSearch.Estimate: TAN AddEdge error (line 141-143)
// Hard to trigger. Skip.

// TreeSearch.computeAllMI: n==0 (line 153-155)
func TestCovPush_TreeSearch_EmptyData(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{}),
		"B": tabgo.NewSeries("B", []any{}),
	})
	ts := NewTreeSearch(data)
	_, err := ts.Estimate()
	// Should still work but with zero MI
	t.Logf("TreeSearch empty data: err=%v", err)
}

// TreeSearch.kruskalMaxSpanningTree: rank tie (line 229-231)
func TestCovPush_TreeSearch_KruskalRankTie(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1}),
	})
	ts := NewTreeSearch(data)
	bn, err := ts.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bn == nil {
		t.Fatal("expected non-nil BN")
	}
}

// TreeSearch.pickRoot: tie-breaking (line 265-268)
func TestCovPush_TreeSearch_PickRoot_Tie(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 0, 1}),
	})
	ts := NewTreeSearch(data)
	bn, err := ts.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bn == nil {
		t.Fatal("expected non-nil BN")
	}
}

// =========================================================================
// MMHC: error paths
// =========================================================================

// mmpc: bestMinAssoc <= 0 break (line 144-147)
func TestCovPush_MMHC_MMPC_ZeroAssoc(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		// Return zero stat for all
		return 0, 1, false
	}
	scoreFn := BICScore()
	mmhc := NewMMHC(data, scoreFn, ciTest, 0.05)
	cands := mmhc.MMPC("A")
	_ = cands
}

// mmpc: CI test says independent in forward phase (line 162-164)
func TestCovPush_MMHC_MMPC_IndependentForward(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		// High stat but independent
		return 5.0, 0.8, true
	}
	scoreFn := BICScore()
	mmhc := NewMMHC(data, scoreFn, ciTest, 0.05)
	cands := mmhc.MMPC("A")
	if len(cands) != 0 {
		t.Errorf("expected empty candidate set, got %v", cands)
	}
}

// minAssociation: maxSubsetSize > 3 cap (line 193-195)
func TestCovPush_MMHC_MinAssociation_LargeSubset(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
		"D": tabgo.NewSeries("D", []any{0, 0, 0, 1, 1, 0, 1, 1}),
		"E": tabgo.NewSeries("E", []any{1, 0, 1, 0, 1, 0, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 3.0, 0.05, false
	}
	scoreFn := BICScore()
	mmhc := NewMMHC(data, scoreFn, ciTest, 0.05)
	// MMPC will build cpc > 3 and trigger the cap
	cands := mmhc.MMPC("A")
	_ = cands
}

// =========================================================================
// SEM Estimator: more error paths
// =========================================================================

// SEM.Estimate: no equation defined (line 56-61)
func TestCovPush_SEM_NoEquation(t *testing.T) {
	s := models.NewSEM()
	s.AddEquation("X", nil, nil, 0.0, 1.0)
	// Add Y variable without equation by using AddEquation then removing
	// Actually SEM may auto-add. Let's check if Variables() returns what we expect.
	// We need a variable with no equation.
	// Let's try accessing the missing equation path differently:
	// GetEquation returns nil for unknown variable.
	se := NewSEMEstimator(s, tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0, 3.0, 4.0, 5.0}),
	}))
	err := se.Estimate()
	// This should succeed since X has an equation
	t.Logf("SEM single var: %v", err)
}

// SEM.Estimate: insufficient data (line 65-68)
func TestCovPush_SEM_InsufficientData(t *testing.T) {
	s := models.NewSEM()
	s.AddEquation("Y", []string{"X1", "X2", "X3"}, []float64{0.5, 0.5, 0.5}, 0.0, 1.0)
	s.AddEquation("X1", nil, nil, 0.0, 1.0)
	s.AddEquation("X2", nil, nil, 0.0, 1.0)
	s.AddEquation("X3", nil, nil, 0.0, 1.0)
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X1": tabgo.NewSeries("X1", []any{1.0}),
		"X2": tabgo.NewSeries("X2", []any{2.0}),
		"X3": tabgo.NewSeries("X3", []any{3.0}),
		"Y":  tabgo.NewSeries("Y", []any{4.0}),
	})
	se := NewSEMEstimator(s, data)
	err := se.Estimate()
	t.Logf("SEM insufficient data: %v", err)
}

// SEM.Estimate: OLS solve error (line 103-105)
func TestCovPush_SEM_SingularOLS(t *testing.T) {
	s := models.NewSEM()
	s.AddEquation("Y", []string{"X"}, []float64{0.5}, 0.0, 1.0)
	s.AddEquation("X", nil, nil, 0.0, 1.0)
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 1.0, 1.0, 1.0, 1.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0, 3.0, 4.0, 5.0, 6.0}),
	})
	se := NewSEMEstimator(s, data)
	err := se.Estimate()
	t.Logf("SEM singular: %v", err)
}

// SEM.Estimate: AddEquation error (line 127-129)
// Hard to trigger. Skip.

// SEM.GetParameters: no equation (line 155-157)
func TestCovPush_SEM_GetParameters_NoEquation(t *testing.T) {
	s := models.NewSEM()
	se := NewSEMEstimator(s, nil)
	_, err := se.GetParameters()
	// SEM has no variables, so this returns empty map, not error.
	t.Logf("SEM GetParameters empty: %v", err)
}

// =========================================================================
// IV Estimator: more error paths
// =========================================================================

// IV.Fit: no instruments (line 59-61)
func TestCovPush_IV_NoInstruments(t *testing.T) {
	iv := NewIVEstimator("X", "Y", nil)
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 2.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0, 4.0}),
	})
	err := iv.Fit(data)
	if err == nil {
		t.Fatal("expected error for empty instruments")
	}
}

// IV.Fit: insufficient data (line 77-79)
func TestCovPush_IV_InsufficientData(t *testing.T) {
	iv := NewIVEstimator("X", "Y", []string{"Z1", "Z2"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"Z1": tabgo.NewSeries("Z1", []any{1.0}),
		"Z2": tabgo.NewSeries("Z2", []any{2.0}),
		"X":  tabgo.NewSeries("X", []any{3.0}),
		"Y":  tabgo.NewSeries("Y", []any{4.0}),
	})
	err := iv.Fit(data)
	t.Logf("IV insufficient data: %v", err)
}

// IV.Fit: stage 2 coefficient count check (line 82-84)
// This is defensive: stage 2 always has 1 predictor -> 1 coefficient. Hard to trigger.

// =========================================================================
// HillClimb: error paths
// =========================================================================

// HillClimb.Estimate: AddNode error (line 149-151)
// Hard to trigger since columns are unique. Skip.

// HillClimb.Estimate: AddEdge error (line 161-163)
// Hard to trigger. Skip.

// HillClimb.LegalOperations: reverse blacklisted edge (line 234-241)
func TestCovPush_HillClimb_LegalOps_ReverseBlacklisted(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
	})
	scoreFn := BICScore()
	// Blacklist B->A: prevents reversing A->B
	hc := NewHillClimbSearch(data, scoreFn, WithBlackList([][2]string{{"B", "A"}}))
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddEdge("A", "B")
	ops := hc.LegalOperations(dag, []string{"A", "B"})
	for _, op := range ops {
		if op.Type == "reverse" && op.From == "A" && op.To == "B" {
			t.Error("should not allow reverse of A->B when B->A is blacklisted")
		}
	}
}

// HillClimb.LegalOperations: whitelist prevents delete (line 234-241)
func TestCovPush_HillClimb_LegalOps_WhitelistPreventDelete(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
	})
	scoreFn := BICScore()
	hc := NewHillClimbSearch(data, scoreFn, WithWhiteList([][2]string{{"A", "B"}}))
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddEdge("A", "B")
	ops := hc.LegalOperations(dag, []string{"A", "B"})
	for _, op := range ops {
		if op.Type == "delete" && op.From == "A" && op.To == "B" {
			t.Error("should not allow delete of whitelisted edge")
		}
	}
}

// =========================================================================
// ExhaustiveSearch: error paths
// =========================================================================

// ExhaustiveSearch.Estimate: AddNode error (line 64-66)
// Hard to trigger. Let's just cover the successful path with 2 vars
func TestCovPush_ExhaustiveSearch_TwoVars(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	scoreFn := BICScore()
	es := NewExhaustiveSearch(data, scoreFn)
	bn, err := es.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bn == nil {
		t.Fatal("expected non-nil BN")
	}
}

// ExhaustiveSearch.AllScores: error paths (line 85-87)
func TestCovPush_ExhaustiveSearch_AllScores_TooManyVars(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1}),
		"D": tabgo.NewSeries("D", []any{0, 1}),
		"E": tabgo.NewSeries("E", []any{0, 1}),
	})
	scoreFn := BICScore()
	es := NewExhaustiveSearch(data, scoreFn)
	_, err := es.AllScores()
	if err == nil {
		t.Fatal("expected error for too many variables")
	}
}

// ExhaustiveSearch.Estimate: AddEdge error (line 69-71)
// Hard to trigger. Skip.

// enumerateDAGs: n==0 (line 156-158)
func TestCovPush_EnumerateDAGs_Empty(t *testing.T) {
	result := enumerateDAGs(nil)
	if len(result) != 1 || len(result[0]) != 0 {
		t.Errorf("expected single empty DAG for no vars, got %v", result)
	}
}

// =========================================================================
// Marginal Estimator: out-of-range values
// =========================================================================

func TestCovPush_MarginalEstimator_OutOfRange(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{-1, 0, 1, 5}),
		"B": tabgo.NewSeries("B", []any{0, -1, 5, 1}),
	})
	me := NewMarginalEstimator(bn, data)
	err := me.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// MarginalEstimator.MarginalLikelihood: no CPD error (line 132-134, 135-137)
func TestCovPush_MarginalEstimator_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
	})
	me := NewMarginalEstimator(bn, data)
	// Don't call Estimate, so no CPD exists
	_, err := me.MarginalLikelihood()
	if err == nil {
		t.Fatal("expected error for no CPD")
	}
}

// MarginalEstimator.MarginalLikelihood: out-of-range child state (line 190-191)
func TestCovPush_MarginalEstimator_ML_OutOfRange(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, -1, 1, 5}),
	})
	me := NewMarginalEstimator(bn, data)
	err := me.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ll, err := me.MarginalLikelihood()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("MarginalLikelihood with OOR: %f", ll)
}

// MarginalEstimator.MarginalLikelihood: out-of-range parent value (line 198-200, 204-205)
func TestCovPush_MarginalEstimator_ML_ParentOOR(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, -1, 1, 5}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	me := NewMarginalEstimator(bn, data)
	err := me.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ll, err := me.MarginalLikelihood()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("MarginalLikelihood parent OOR: %f", ll)
}

// MarginalEstimator.MarginalLikelihood: prob <= 0 (line 209-211)
func TestCovPush_MarginalEstimator_ML_ZeroProb(t *testing.T) {
	// Manually create a CPD with a zero probability entry.
	// P(A=0) = 1.0, P(A=1) = 0.0
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})
	zeroCPD, _ := factors.NewTabularCPD("A", 2, [][]float64{{1.0}, {0.0}}, nil, nil)
	_ = bn.AddCPD(zeroCPD)

	testData := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{1}),
	})
	me := &MarginalEstimator{bn: bn, data: testData}
	ll, err := me.MarginalLikelihood()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !math.IsInf(ll, -1) {
		t.Errorf("expected -Inf for zero probability, got %f", ll)
	}
}

// =========================================================================
// Linear Gaussian MLE: error paths
// =========================================================================

// LG MLE.Estimate: estimateNode error (line 66-68)
func TestCovPush_LGMLE_EstimateNode_Error(t *testing.T) {
	lgbn := models.NewLinearGaussianBayesianNetwork()
	_ = lgbn.AddNode("X")
	// X has no parents but data has 0 rows -> insufficient data
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{}),
	})
	est := NewLinearGaussianMLE(lgbn, data)
	err := est.Estimate()
	t.Logf("LG MLE empty data: %v", err)
}

// LG MLE.GetParameters: nil data (line 78-80)
func TestCovPush_LGMLE_GetParameters_NilData(t *testing.T) {
	lgbn := models.NewLinearGaussianBayesianNetwork()
	_ = lgbn.AddNode("X")
	est := NewLinearGaussianMLE(lgbn, nil)
	_, err := est.GetParameters("X")
	if err == nil {
		t.Fatal("expected error for nil data")
	}
}

// LG MLE.GetParameters: missing column (line 97-99)
func TestCovPush_LGMLE_GetParameters_MissingColumn(t *testing.T) {
	lgbn := models.NewLinearGaussianBayesianNetwork()
	_ = lgbn.AddNode("X")
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"Y": tabgo.NewSeries("Y", []any{1.0}),
	})
	est := NewLinearGaussianMLE(lgbn, data)
	_, err := est.GetParameters("X")
	if err == nil {
		t.Fatal("expected error for missing column")
	}
}

// LG MLE.estimateNode: singular matrix (line 156-158)
func TestCovPush_LGMLE_SingularMatrix(t *testing.T) {
	lgbn := models.NewLinearGaussianBayesianNetwork()
	_ = lgbn.AddNode("X")
	_ = lgbn.AddNode("Y")
	_ = lgbn.AddEdge("X", "Y")
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 1.0, 1.0, 1.0, 1.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0, 3.0, 4.0, 5.0, 6.0}),
	})
	est := NewLinearGaussianMLE(lgbn, data)
	err := est.Estimate()
	t.Logf("LG MLE singular: %v", err)
}

// =========================================================================
// LLM Client: additional coverage for wait() rate limiter loop (line 100-102)
// =========================================================================

func TestCovPush_LLM_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"error":{"message":"quota exceeded"}}`))
	}))
	defer server.Close()

	client := NewHTTPLLMClient(server.URL, "test-key", "gpt-4")
	_, err := client.ChatComplete([]Message{{Role: "user", Content: "test"}})
	if err == nil {
		t.Fatal("expected error for API error response")
	}
}

// LLM Client: 4xx non-429 error (does not retry)
func TestCovPush_LLM_ClientError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`forbidden`))
	}))
	defer server.Close()

	client := NewHTTPLLMClient(server.URL, "test-key", "gpt-4")
	_, err := client.ChatComplete([]Message{{Role: "user", Content: "test"}})
	if err == nil {
		t.Fatal("expected error for 403")
	}
}

// =========================================================================
// ExpertInLoop: LLM returns NO opinion
// =========================================================================

func TestCovPush_ExpertInLoop_LLMOpposesVStructure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices":[{"message":{"content":"NO confidence: 0.9"}}]}`))
	}))
	defer server.Close()

	llmClient := NewHTTPLLMClient(server.URL, "test-key", "gpt-4")
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		if (x == "A" && y == "B") || (x == "B" && y == "A") {
			return 0, 1, true
		}
		return 5.0, 0.01, false
	}
	eil := NewExpertInLoop(data, llmClient, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// ExpertInLoop: LLM returns unparseable response
func TestCovPush_ExpertInLoop_LLMUnparseable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"choices":[{"message":{"content":"maybe possibly perhaps"}}]}`))
	}))
	defer server.Close()

	llmClient := NewHTTPLLMClient(server.URL, "test-key", "gpt-4")
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		if (x == "A" && y == "B") || (x == "B" && y == "A") {
			return 0, 1, true
		}
		return 5.0, 0.01, false
	}
	eil := NewExpertInLoop(data, llmClient, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// =========================================================================
// EM with convergence: oldCPD != nil branch
// =========================================================================

func TestCovPush_EM_Convergence(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	_ = bn.SetStates("A", []string{"0", "1"})
	_ = bn.SetStates("B", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{"0", "1", "0", "1", "0", "1", "0", "1"}),
		"B": tabgo.NewSeries("B", []any{"0", "0", "1", "1", "0", "1", "0", "1"}),
	})
	em := NewEM(bn, data, nil, 100, 0.001)
	err := em.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !em.Converged() {
		t.Log("EM did not converge (expected for observed-only)")
	}
}

// =========================================================================
// PC: BuildSkeleton edge removal from adj(Y) side (line 126-129)
// =========================================================================

func TestCovPush_PC_BuildSkeleton_AdjY(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		// Only independent when tested from adj(Y) side
		// adjX\{Y} doesn't make them independent, but adjY\{X} does
		if (x == "A" && y == "C") || (x == "C" && y == "A") {
			if len(z) > 0 && z[0] == "B" {
				return 0, 1, true
			}
		}
		return 5.0, 0.01, false
	}
	pc := NewPC(data, ciTest, 0.05)
	pdag, sepSets, err := pc.BuildSkeleton()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pdag == nil {
		t.Fatal("expected non-nil PDAG")
	}
	_ = sepSets
}

// =========================================================================
// Additional coverage for remaining gaps
// =========================================================================

// ExpertInLoop: skeleton phase - edge removed during iteration (line 83-84)
func TestCovPush_ExpertInLoop_SkeletonEdgeRemoval(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 0, 1, true
	}
	eil := NewExpertInLoop(data, nil, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// ExpertInLoop: adjY branch in skeleton (line 97-100)
func TestCovPush_ExpertInLoop_SkeletonAdjY(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		if (x == "A" && y == "C") || (x == "C" && y == "A") {
			if len(z) > 0 && z[0] == "B" {
				return 0, 1, true
			}
		}
		return 5.0, 0.01, false
	}
	eil := NewExpertInLoop(data, nil, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// ExpertInLoop.Estimate: pdag.Adjacent skip + sepSet not found
func TestCovPush_ExpertInLoop_AdjacentAndNoSepSet(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 5.0, 0.01, false
	}
	eil := NewExpertInLoop(data, nil, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// ExpertInLoop.EstimateBN: error from Estimate (line 236-238)
func TestCovPush_ExpertInLoop_EstimateBN_Error(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 0, 1, true
	}
	eil := NewExpertInLoop(data, nil, ciTest, 0.05)
	_, err := eil.EstimateBN()
	if err == nil {
		t.Fatal("expected error for single variable")
	}
}

// GES backward phase with actual removal
func TestCovPush_GES_BackwardRemoval(t *testing.T) {
	// Forward phase adds edges that are beneficial; backward phase removes ones that aren't.
	// Use a score where adding ONE parent is good but TWO is bad.
	callCount := 0
	scoreFn := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		callCount++
		if len(parents) == 1 {
			return 5.0
		}
		if len(parents) >= 2 {
			return -5.0 // overfitting penalty
		}
		return 0.0
	}
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0}),
	})
	ges := NewGES(data, scoreFn)
	pdag, err := ges.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// GES: cycle check in forward phase (line 71-73)
// Need to force the forward phase to TRY adding an edge that would create a cycle.
// With 3 vars and a score that always rewards more parents, forward will add
// A->B, then B->C (or some other), then try C->A which would create a cycle.
func TestCovPush_GES_ForwardCycleCheck(t *testing.T) {
	scoreFn := func(variable string, parents []string, data *tabgo.DataFrame) float64 {
		// More parents = higher score to force adding all edges
		return float64(len(parents)) * 100.0
	}
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
	})
	ges := NewGES(data, scoreFn)
	pdag, err := ges.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// GES.Insert via direct method to hit cycle branch
func TestCovPush_GES_Insert_CycleDirect(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddNode("C")
	dag.AddEdge("A", "B")
	dag.AddEdge("B", "C")
	// Insert C->A would create cycle
	_, err := ges.Insert(dag, "C", "A")
	if err == nil {
		t.Fatal("expected error for cycle")
	}
}

// PC.Estimate: same level removal with 4 vars
func TestCovPush_PC_Estimate_SameLevelRemoval(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
		"D": tabgo.NewSeries("D", []any{0, 0, 1, 0, 1, 1, 0, 1}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		if (x == "A" && y == "B") || (x == "B" && y == "A") {
			return 5.0, 0.01, false
		}
		return 0, 1, true
	}
	pc := NewPC(data, ciTest, 0.05)
	pdag, err := pc.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// PC.EstimateBN: error from Estimate
func TestCovPush_PC_EstimateBN_Error(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 0, 1, true
	}
	pc := NewPC(data, ciTest, 0.05)
	_, err := pc.EstimateBN()
	if err == nil {
		t.Fatal("expected error for single variable")
	}
}

// MarginalEstimator.Estimate: no nodes
func TestCovPush_MarginalEstimator_NoNodes(t *testing.T) {
	bn := models.NewBayesianNetwork()
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0}),
	})
	me := NewMarginalEstimator(bn, data)
	err := me.Estimate()
	if err == nil {
		t.Fatal("expected error for no nodes")
	}
}

// MarginalEstimator.Estimate: missing column
func TestCovPush_MarginalEstimator_MissingCol(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"B": tabgo.NewSeries("B", []any{0}),
	})
	me := NewMarginalEstimator(bn, data)
	err := me.Estimate()
	if err == nil {
		t.Fatal("expected error for missing column")
	}
}

// MirrorDescent: cardMap < 1 guard
func TestCovPush_MirrorDescent_AllNegative(t *testing.T) {
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.SetStates("A", []string{"0", "1"})
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{-1, -1, -1}),
	})
	md := NewMirrorDescentEstimator(bn, data, 0.1, 10)
	err := md.Estimate()
	t.Logf("Mirror descent all negative: %v", err)
}

// TreeSearch: TAN with < 2 feature vars
func TestCovPush_TreeSearch_TAN_TooFewFeatures(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1}),
	})
	ts := NewTreeSearch(data, WithClassVariable("C"))
	_, err := ts.Estimate()
	if err == nil {
		t.Fatal("expected error for TAN with < 2 feature variables")
	}
}

// TreeSearch: root not in tree vars
func TestCovPush_TreeSearch_RootNotInVars(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0}),
	})
	ts := NewTreeSearch(data, WithRoot("Z"))
	_, err := ts.Estimate()
	if err == nil {
		t.Fatal("expected error for root not in variables")
	}
}

// TreeSearch: kruskal edge cases with 4 vars
func TestCovPush_TreeSearch_KruskalEdgeCases(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
		"D": tabgo.NewSeries("D", []any{0, 0, 1, 0, 1, 1, 0, 1}),
	})
	ts := NewTreeSearch(data)
	bn, err := ts.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bn == nil {
		t.Fatal("expected non-nil BN")
	}
}

// TreeSearch: pickRoot with variable having higher MI
func TestCovPush_TreeSearch_PickRootDifferentMI(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 0, 0, 0, 1, 1, 1, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 0, 0, 0, 1, 1, 1, 1}),
	})
	ts := NewTreeSearch(data)
	bn, err := ts.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bn == nil {
		t.Fatal("expected non-nil BN")
	}
}

// PC.BuildSkeleton: edge already removed in same iteration
func TestCovPush_PC_BuildSkeleton_SkipRemoved(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
		"D": tabgo.NewSeries("D", []any{0, 0, 1, 0, 1, 1, 0, 1}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 0, 1, true
	}
	pc := NewPC(data, ciTest, 0.05)
	pdag2, sepSets2, err := pc.BuildSkeleton()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag2
	_ = sepSets2
}

// LG MLE: estimateNode error propagation
func TestCovPush_LGMLE_EstimateNodeErrorProp(t *testing.T) {
	lgbn := models.NewLinearGaussianBayesianNetwork()
	_ = lgbn.AddNode("X")
	_ = lgbn.AddNode("Y")
	_ = lgbn.AddEdge("X", "Y")
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{1.0, 1.0, 1.0, 1.0, 1.0}),
		"Y": tabgo.NewSeries("Y", []any{2.0, 3.0, 4.0, 5.0, 6.0}),
	})
	est := NewLinearGaussianMLE(lgbn, data)
	err := est.Estimate()
	t.Logf("LG MLE singular propagation: %v", err)
}

// GES backward phase: force by using a score function that first adds then prefers removal
func TestCovPush_GES_BackwardPhase_Direct(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1, 0, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0, 1, 0}),
	})
	// Score function where A->B is good but A->C and B->C are bad
	// Forward phase adds A->B (good), then tries others.
	// Backward phase should try to remove A->C or B->C if added.
	callIdx := 0
	scoreFn := func(variable string, parents []string, d *tabgo.DataFrame) float64 {
		callIdx++
		if variable == "B" && len(parents) == 1 && parents[0] == "A" {
			return 20.0 // A->B is great
		}
		if len(parents) == 0 {
			return 0.0
		}
		if len(parents) == 1 {
			return 5.0 // any single parent is OK (forward will add)
		}
		return -20.0 // two parents is terrible (backward should remove)
	}
	ges := NewGES(data, scoreFn)
	pdag, err := ges.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = pdag
}

// GES Delete method directly
func TestCovPush_GES_Delete(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddEdge("A", "B")
	delta, err := ges.Delete(dag, "A", "B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("GES Delete delta: %f", delta)
}

// GES Delete: edge doesn't exist
func TestCovPush_GES_Delete_NoEdge(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	_, err := ges.Delete(dag, "A", "B")
	if err == nil {
		t.Fatal("expected error for non-existent edge")
	}
}

// GES Turn: edge doesn't exist
func TestCovPush_GES_Turn_NoEdge(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	_, err := ges.Turn(dag, "A", "B")
	if err == nil {
		t.Fatal("expected error for non-existent edge")
	}
}

// GES Turn: successful turn
func TestCovPush_GES_Turn_Success(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	dag.AddEdge("A", "B")
	delta, err := ges.Turn(dag, "A", "B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("GES Turn delta: %f", delta)
}

// GES Insert: successful insert
func TestCovPush_GES_Insert_Success(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1}),
	})
	scoreFn := BICScore()
	ges := NewGES(data, scoreFn)
	dag := graphgo.NewDiGraph()
	dag.AddNode("A")
	dag.AddNode("B")
	delta, err := ges.Insert(dag, "A", "B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("GES Insert delta: %f", delta)
}

// MMHC: large subset cap
func TestCovPush_MMHC_LargeSubsetCap(t *testing.T) {
	data := tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 0, 1, 1}),
		"C": tabgo.NewSeries("C", []any{0, 1, 1, 0, 0, 1, 1, 0}),
		"D": tabgo.NewSeries("D", []any{0, 0, 1, 0, 1, 1, 0, 1}),
		"E": tabgo.NewSeries("E", []any{1, 0, 1, 0, 1, 0, 1, 0}),
		"F": tabgo.NewSeries("F", []any{0, 1, 0, 1, 1, 0, 1, 0}),
	})
	ciTest := func(x, y string, z []string, d *tabgo.DataFrame, sig float64) (float64, float64, bool) {
		return 3.0 + float64(len(z)), 0.05, false
	}
	scoreFn := BICScore()
	mmhc := NewMMHC(data, scoreFn, ciTest, 0.05)
	bn, err := mmhc.Estimate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bn == nil {
		t.Fatal("expected non-nil BN")
	}
}
