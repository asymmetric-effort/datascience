//go:build unit

package base

import (
	"testing"
)

func TestNewDAG(t *testing.T) {
	d := NewDAG()
	if d == nil {
		t.Fatal("NewDAG returned nil")
	}
	if len(d.Nodes()) != 0 {
		t.Errorf("new DAG should have 0 nodes, got %d", len(d.Nodes()))
	}
	if len(d.Edges()) != 0 {
		t.Errorf("new DAG should have 0 edges, got %d", len(d.Edges()))
	}
}

func TestAddNode(t *testing.T) {
	d := NewDAG()
	if err := d.AddNode("A"); err != nil {
		t.Fatalf("AddNode failed: %v", err)
	}
	if !d.HasNode("A") {
		t.Error("HasNode returned false for added node")
	}
}

func TestAddNodeDuplicate(t *testing.T) {
	d := NewDAG()
	_ = d.AddNode("A")
	err := d.AddNode("A")
	if err == nil {
		t.Error("expected error when adding duplicate node")
	}
}

func TestAddNodes(t *testing.T) {
	d := NewDAG()
	if err := d.AddNodes("A", "B", "C"); err != nil {
		t.Fatalf("AddNodes failed: %v", err)
	}
	nodes := d.Nodes()
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}
	expected := []string{"A", "B", "C"}
	for i, n := range expected {
		if nodes[i] != n {
			t.Errorf("nodes[%d] = %q, want %q", i, nodes[i], n)
		}
	}
}

func TestAddNodesDuplicate(t *testing.T) {
	d := NewDAG()
	err := d.AddNodes("A", "B", "A")
	if err == nil {
		t.Error("expected error when adding duplicate node via AddNodes")
	}
	// A and B should still exist (added before the error).
	if !d.HasNode("A") || !d.HasNode("B") {
		t.Error("nodes added before error should be retained")
	}
}

func TestRemoveNode(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B")
	_ = d.AddEdge("A", "B")
	if err := d.RemoveNode("A"); err != nil {
		t.Fatalf("RemoveNode failed: %v", err)
	}
	if d.HasNode("A") {
		t.Error("removed node should not exist")
	}
	if d.HasEdge("A", "B") {
		t.Error("edges from removed node should be gone")
	}
}

func TestRemoveNodeNotFound(t *testing.T) {
	d := NewDAG()
	err := d.RemoveNode("X")
	if err == nil {
		t.Error("expected error when removing non-existent node")
	}
}

func TestAddEdge(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B")
	if err := d.AddEdge("A", "B"); err != nil {
		t.Fatalf("AddEdge failed: %v", err)
	}
	if !d.HasEdge("A", "B") {
		t.Error("HasEdge returned false for added edge")
	}
	if d.HasEdge("B", "A") {
		t.Error("reverse edge should not exist")
	}
}

func TestAddEdgeMissingNode(t *testing.T) {
	d := NewDAG()
	_ = d.AddNode("A")
	err := d.AddEdge("A", "B")
	if err == nil {
		t.Error("expected error when target node does not exist")
	}
	err = d.AddEdge("X", "A")
	if err == nil {
		t.Error("expected error when source node does not exist")
	}
}

func TestAddEdgeDuplicate(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B")
	_ = d.AddEdge("A", "B")
	err := d.AddEdge("A", "B")
	if err == nil {
		t.Error("expected error when adding duplicate edge")
	}
}

// --- Acyclicity validation tests ---

func TestAddEdgeRejectsSelfLoop(t *testing.T) {
	d := NewDAG()
	_ = d.AddNode("A")
	err := d.AddEdge("A", "A")
	if err == nil {
		t.Fatal("expected error for self-loop")
	}
	if d.HasEdge("A", "A") {
		t.Error("self-loop edge should not exist after rejection")
	}
}

func TestAddEdgeRejectsDirectCycle(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B")
	_ = d.AddEdge("A", "B")
	err := d.AddEdge("B", "A")
	if err == nil {
		t.Fatal("expected error for direct cycle B->A")
	}
	if d.HasEdge("B", "A") {
		t.Error("cycle-creating edge should not exist after rejection")
	}
	// Original edge should still be there.
	if !d.HasEdge("A", "B") {
		t.Error("original edge A->B should still exist")
	}
}

func TestAddEdgeRejectsIndirectCycle(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C")
	_ = d.AddEdge("A", "B")
	_ = d.AddEdge("B", "C")

	err := d.AddEdge("C", "A")
	if err == nil {
		t.Fatal("expected error for indirect cycle C->A")
	}
	if d.HasEdge("C", "A") {
		t.Error("cycle-creating edge should not exist after rejection")
	}
}

func TestAddEdgeRejectsLongCycle(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C", "D", "E")
	_ = d.AddEdge("A", "B")
	_ = d.AddEdge("B", "C")
	_ = d.AddEdge("C", "D")
	_ = d.AddEdge("D", "E")

	err := d.AddEdge("E", "A")
	if err == nil {
		t.Fatal("expected error for long cycle E->A")
	}
	if d.HasEdge("E", "A") {
		t.Error("cycle-creating edge should not exist")
	}
}

func TestAddEdgeAllowsNonCyclicEdge(t *testing.T) {
	// A->B, A->C, B->C should be valid (diamond, no cycle).
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C")
	_ = d.AddEdge("A", "B")
	_ = d.AddEdge("A", "C")
	err := d.AddEdge("B", "C")
	if err != nil {
		t.Fatalf("diamond pattern should be allowed: %v", err)
	}
}

func TestAddEdgeCycleRevertDoesNotCorrupt(t *testing.T) {
	// After a rejected cycle-creating edge, the graph should remain valid.
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C")
	_ = d.AddEdge("A", "B")
	_ = d.AddEdge("B", "C")
	_ = d.AddEdge("C", "A") // rejected

	// Graph should still be a valid DAG.
	order, err := d.TopologicalOrder()
	if err != nil {
		t.Fatalf("DAG should be valid after rejected edge: %v", err)
	}
	if len(order) != 3 {
		t.Errorf("expected 3 nodes in topo order, got %d", len(order))
	}

	// Should still be able to add valid edges.
	err = d.AddEdge("A", "C")
	if err != nil {
		t.Fatalf("should be able to add A->C after rejected C->A: %v", err)
	}
}

// --- RemoveEdge ---

func TestRemoveEdge(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B")
	_ = d.AddEdge("A", "B")
	if err := d.RemoveEdge("A", "B"); err != nil {
		t.Fatalf("RemoveEdge failed: %v", err)
	}
	if d.HasEdge("A", "B") {
		t.Error("edge should be removed")
	}
}

func TestRemoveEdgeNotFound(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B")
	err := d.RemoveEdge("A", "B")
	if err == nil {
		t.Error("expected error when removing non-existent edge")
	}
}

func TestRemoveEdgeEnablesNewEdge(t *testing.T) {
	// A->B, B->C. Edge C->A would be a cycle. But if we remove A->B, then C->A is fine.
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C")
	_ = d.AddEdge("A", "B")
	_ = d.AddEdge("B", "C")

	// C->A is a cycle right now.
	err := d.AddEdge("C", "A")
	if err == nil {
		t.Fatal("C->A should be rejected while A->B->C exists")
	}

	// Remove A->B, breaking the path.
	_ = d.RemoveEdge("A", "B")

	// Now C->A should be allowed (only edges: B->C, C->A, no cycle).
	err = d.AddEdge("C", "A")
	if err != nil {
		t.Fatalf("C->A should be allowed after removing A->B: %v", err)
	}
}

// --- Parents/Children ---

func TestParentsAndChildren(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C", "D")
	_ = d.AddEdge("A", "C")
	_ = d.AddEdge("B", "C")
	_ = d.AddEdge("C", "D")

	parents := d.Parents("C")
	if len(parents) != 2 || parents[0] != "A" || parents[1] != "B" {
		t.Errorf("Parents(C) = %v, want [A B]", parents)
	}

	children := d.Children("C")
	if len(children) != 1 || children[0] != "D" {
		t.Errorf("Children(C) = %v, want [D]", children)
	}

	// Root node has no parents.
	if len(d.Parents("A")) != 0 {
		t.Error("Parents(A) should be empty")
	}

	// Leaf node has no children.
	if len(d.Children("D")) != 0 {
		t.Error("Children(D) should be empty")
	}
}

// --- GetRoots / GetLeaves ---

func TestGetRoots(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C", "D")
	_ = d.AddEdge("A", "C")
	_ = d.AddEdge("B", "C")
	_ = d.AddEdge("C", "D")

	roots := d.GetRoots()
	if len(roots) != 2 || roots[0] != "A" || roots[1] != "B" {
		t.Errorf("GetRoots() = %v, want [A B]", roots)
	}
}

func TestGetLeaves(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C", "D")
	_ = d.AddEdge("A", "B")
	_ = d.AddEdge("A", "C")
	_ = d.AddEdge("B", "D")

	leaves := d.GetLeaves()
	if len(leaves) != 2 || leaves[0] != "C" || leaves[1] != "D" {
		t.Errorf("GetLeaves() = %v, want [C D]", leaves)
	}
}

func TestGetRootsAllRoots(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("X", "Y", "Z")
	// No edges: all nodes are roots and leaves.
	roots := d.GetRoots()
	if len(roots) != 3 {
		t.Errorf("expected 3 roots, got %d", len(roots))
	}
	leaves := d.GetLeaves()
	if len(leaves) != 3 {
		t.Errorf("expected 3 leaves, got %d", len(leaves))
	}
}

// --- TopologicalOrder ---

func TestTopologicalOrder(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C", "D")
	_ = d.AddEdge("A", "B")
	_ = d.AddEdge("A", "C")
	_ = d.AddEdge("B", "D")
	_ = d.AddEdge("C", "D")

	order, err := d.TopologicalOrder()
	if err != nil {
		t.Fatalf("TopologicalOrder failed: %v", err)
	}
	if len(order) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(order))
	}

	// Verify topological property: for every edge u->v, u appears before v.
	pos := make(map[string]int)
	for i, n := range order {
		pos[n] = i
	}
	edges := d.Edges()
	for _, e := range edges {
		if pos[e.Src] >= pos[e.Dst] {
			t.Errorf("topological violation: %s (pos %d) should come before %s (pos %d)",
				e.Src, pos[e.Src], e.Dst, pos[e.Dst])
		}
	}
}

func TestTopologicalOrderEmpty(t *testing.T) {
	d := NewDAG()
	order, err := d.TopologicalOrder()
	if err != nil {
		t.Fatalf("TopologicalOrder on empty DAG failed: %v", err)
	}
	if order != nil && len(order) != 0 {
		t.Errorf("expected nil or empty slice, got %v", order)
	}
}

func TestTopologicalOrderSingleNode(t *testing.T) {
	d := NewDAG()
	_ = d.AddNode("A")
	order, err := d.TopologicalOrder()
	if err != nil {
		t.Fatalf("TopologicalOrder failed: %v", err)
	}
	if len(order) != 1 || order[0] != "A" {
		t.Errorf("expected [A], got %v", order)
	}
}

// --- Nodes/Edges ---

func TestNodesSorted(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("C", "A", "B")
	nodes := d.Nodes()
	if len(nodes) != 3 || nodes[0] != "A" || nodes[1] != "B" || nodes[2] != "C" {
		t.Errorf("Nodes() should be sorted, got %v", nodes)
	}
}

func TestEdgesSorted(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C")
	_ = d.AddEdge("B", "C")
	_ = d.AddEdge("A", "C")
	_ = d.AddEdge("A", "B")

	edges := d.Edges()
	if len(edges) != 3 {
		t.Fatalf("expected 3 edges, got %d", len(edges))
	}
	// Should be sorted: (A,B), (A,C), (B,C).
	if edges[0].Src != "A" || edges[0].Dst != "B" {
		t.Errorf("edges[0] = (%s,%s), want (A,B)", edges[0].Src, edges[0].Dst)
	}
	if edges[1].Src != "A" || edges[1].Dst != "C" {
		t.Errorf("edges[1] = (%s,%s), want (A,C)", edges[1].Src, edges[1].Dst)
	}
	if edges[2].Src != "B" || edges[2].Dst != "C" {
		t.Errorf("edges[2] = (%s,%s), want (B,C)", edges[2].Src, edges[2].Dst)
	}
}

// --- Copy ---

func TestCopy(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B", "C")
	_ = d.AddEdge("A", "B")
	_ = d.AddEdge("B", "C")

	c := d.Copy()

	// Copy should have the same structure.
	if len(c.Nodes()) != 3 {
		t.Errorf("copy should have 3 nodes, got %d", len(c.Nodes()))
	}
	if !c.HasEdge("A", "B") || !c.HasEdge("B", "C") {
		t.Error("copy should have the same edges")
	}

	// Modifying the copy should not affect the original.
	_ = c.AddNode("D")
	_ = c.AddEdge("C", "D")
	if d.HasNode("D") {
		t.Error("original should not have node D after modifying copy")
	}
	if d.HasEdge("C", "D") {
		t.Error("original should not have edge C->D after modifying copy")
	}
}

func TestCopyAcyclicityPreserved(t *testing.T) {
	d := NewDAG()
	_ = d.AddNodes("A", "B")
	_ = d.AddEdge("A", "B")
	c := d.Copy()

	// The copy should still enforce acyclicity.
	err := c.AddEdge("B", "A")
	if err == nil {
		t.Error("copy should still enforce acyclicity")
	}
}

// --- Complex scenario ---

func TestComplexDAG(t *testing.T) {
	// Build the classic "asia" BN structure:
	// asia -> tub, smoke -> {lung, bronc}, tub -> either, lung -> either,
	// either -> {xray, dysp}, bronc -> dysp
	d := NewDAG()
	_ = d.AddNodes("asia", "tub", "smoke", "lung", "bronc", "either", "xray", "dysp")

	edges := [][2]string{
		{"asia", "tub"},
		{"smoke", "lung"},
		{"smoke", "bronc"},
		{"tub", "either"},
		{"lung", "either"},
		{"either", "xray"},
		{"either", "dysp"},
		{"bronc", "dysp"},
	}
	for _, e := range edges {
		if err := d.AddEdge(e[0], e[1]); err != nil {
			t.Fatalf("AddEdge(%s, %s) failed: %v", e[0], e[1], err)
		}
	}

	// Verify structure.
	if len(d.Nodes()) != 8 {
		t.Errorf("expected 8 nodes, got %d", len(d.Nodes()))
	}
	if len(d.Edges()) != 8 {
		t.Errorf("expected 8 edges, got %d", len(d.Edges()))
	}

	roots := d.GetRoots()
	if len(roots) != 2 || roots[0] != "asia" || roots[1] != "smoke" {
		t.Errorf("GetRoots() = %v, want [asia smoke]", roots)
	}

	leaves := d.GetLeaves()
	if len(leaves) != 2 || leaves[0] != "dysp" || leaves[1] != "xray" {
		t.Errorf("GetLeaves() = %v, want [dysp xray]", leaves)
	}

	parents := d.Parents("either")
	if len(parents) != 2 || parents[0] != "lung" || parents[1] != "tub" {
		t.Errorf("Parents(either) = %v, want [lung tub]", parents)
	}

	children := d.Children("either")
	if len(children) != 2 || children[0] != "dysp" || children[1] != "xray" {
		t.Errorf("Children(either) = %v, want [dysp xray]", children)
	}

	// Topological order should be valid.
	order, err := d.TopologicalOrder()
	if err != nil {
		t.Fatalf("TopologicalOrder failed: %v", err)
	}
	pos := make(map[string]int)
	for i, n := range order {
		pos[n] = i
	}
	for _, e := range edges {
		if pos[e[0]] >= pos[e[1]] {
			t.Errorf("topological violation: %s should come before %s", e[0], e[1])
		}
	}

	// Cycle-creating edge should be rejected.
	err = d.AddEdge("dysp", "asia")
	if err == nil {
		t.Error("dysp->asia should create a cycle")
	}
}
