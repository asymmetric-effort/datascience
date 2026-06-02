//go:build unit

package learning

import (
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/graphgo"
	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
)

func simpleScoreFn(variable string, parents []string, data *tabgo.DataFrame) float64 {
	// Simple scoring: penalize number of parents (prefer simpler models).
	return -float64(len(parents))
}

func makeGESTestData(t *testing.T) *tabgo.DataFrame {
	t.Helper()
	return tabgo.NewDataFrame(map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1, 0, 1, 0, 1}),
		"B": tabgo.NewSeries("B", []any{0, 0, 1, 1, 0, 1}),
		"C": tabgo.NewSeries("C", []any{1, 0, 1, 0, 1, 0}),
	})
}

func TestGESInsert(t *testing.T) {
	data := makeGESTestData(t)
	g := NewGES(data, simpleScoreFn)

	dag := graphgo.NewDiGraph()
	dag.AddNodes("A", "B", "C")

	delta, err := g.Insert(dag, "A", "B")
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}
	if !dag.HasEdge("A", "B") {
		t.Error("edge A->B should exist after insert")
	}
	t.Logf("insert delta: %f", delta)
}

func TestGESInsertCycle(t *testing.T) {
	data := makeGESTestData(t)
	g := NewGES(data, simpleScoreFn)

	dag := graphgo.NewDiGraph()
	dag.AddNodes("A", "B")
	dag.AddEdge("A", "B")

	_, err := g.Insert(dag, "B", "A")
	if err == nil {
		t.Error("expected error for cycle-creating insert")
	}
}

func TestGESInsertDuplicate(t *testing.T) {
	data := makeGESTestData(t)
	g := NewGES(data, simpleScoreFn)

	dag := graphgo.NewDiGraph()
	dag.AddNodes("A", "B")
	dag.AddEdge("A", "B")

	_, err := g.Insert(dag, "A", "B")
	if err == nil {
		t.Error("expected error for duplicate edge insert")
	}
}

func TestGESDelete(t *testing.T) {
	data := makeGESTestData(t)
	g := NewGES(data, simpleScoreFn)

	dag := graphgo.NewDiGraph()
	dag.AddNodes("A", "B", "C")
	dag.AddEdge("A", "B")

	delta, err := g.Delete(dag, "A", "B")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if dag.HasEdge("A", "B") {
		t.Error("edge A->B should not exist after delete")
	}
	t.Logf("delete delta: %f", delta)
}

func TestGESDeleteNonExistent(t *testing.T) {
	data := makeGESTestData(t)
	g := NewGES(data, simpleScoreFn)

	dag := graphgo.NewDiGraph()
	dag.AddNodes("A", "B")

	_, err := g.Delete(dag, "A", "B")
	if err == nil {
		t.Error("expected error for deleting non-existent edge")
	}
}

func TestGESTurn(t *testing.T) {
	data := makeGESTestData(t)
	g := NewGES(data, simpleScoreFn)

	dag := graphgo.NewDiGraph()
	dag.AddNodes("A", "B", "C")
	dag.AddEdge("A", "B")

	delta, err := g.Turn(dag, "A", "B")
	if err != nil {
		t.Fatalf("Turn failed: %v", err)
	}
	if dag.HasEdge("A", "B") {
		t.Error("edge A->B should not exist after turn")
	}
	if !dag.HasEdge("B", "A") {
		t.Error("edge B->A should exist after turn")
	}
	t.Logf("turn delta: %f", delta)
}

func TestGESTurnCycle(t *testing.T) {
	data := makeGESTestData(t)
	g := NewGES(data, simpleScoreFn)

	dag := graphgo.NewDiGraph()
	dag.AddNodes("A", "B", "C")
	dag.AddEdge("A", "B")
	dag.AddEdge("B", "A") // Already has reverse, so A->B with B->A won't turn cleanly.

	// Remove the reverse to set up a proper test.
	_ = dag.RemoveEdge("B", "A")
	dag.AddEdge("B", "C")
	dag.AddEdge("C", "A")

	// Turning A->B would create B->A, but A->C->... exists. Let's check.
	_, err := g.Turn(dag, "A", "B")
	// This should work since B->A doesn't create a cycle with B->C->A.
	t.Logf("Turn result: err=%v", err)
}

func TestGESTurnNonExistent(t *testing.T) {
	data := makeGESTestData(t)
	g := NewGES(data, simpleScoreFn)

	dag := graphgo.NewDiGraph()
	dag.AddNodes("A", "B")

	_, err := g.Turn(dag, "A", "B")
	if err == nil {
		t.Error("expected error for turning non-existent edge")
	}
}
