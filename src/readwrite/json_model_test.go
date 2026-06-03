//go:build unit

package readwrite

import (
	"bytes"
	"strings"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// ---------------------------------------------------------------------------
// ReadJSON
// ---------------------------------------------------------------------------

func TestReadJSON_Basic(t *testing.T) {
	input := `{
  "name": "test",
  "nodes": ["A", "B"],
  "edges": [["A","B"]],
  "states": {"A": ["s0","s1"], "B": ["s0","s1"]},
  "cpds": {
    "A": {"variable_card": 2, "values": [[0.6],[0.4]]},
    "B": {"variable_card": 2, "values": [[0.3,0.7],[0.7,0.3]], "evidence": ["A"], "evidence_card": [2]}
  }
}`
	bn, err := ReadJSON(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadJSON failed: %v", err)
	}

	nodes := bn.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(nodes))
	}

	edges := bn.Edges()
	if len(edges) != 1 || edges[0] != [2]string{"A", "B"} {
		t.Errorf("edges = %v, want [[A B]]", edges)
	}

	aCPD := bn.GetCPD("A")
	if aCPD == nil {
		t.Fatal("A CPD is nil")
	}
	aData := aCPD.ToFactor().Values().Data()
	assertFloatsClose(t, aData, []float64{0.6, 0.4}, "A CPD")

	bCPD := bn.GetCPD("B")
	if bCPD == nil {
		t.Fatal("B CPD is nil")
	}
	bData := bCPD.ToFactor().Values().Data()
	assertFloatsClose(t, bData, []float64{0.3, 0.7, 0.7, 0.3}, "B CPD")
}

func TestReadJSON_InvalidJSON(t *testing.T) {
	_, err := ReadJSON(strings.NewReader("{invalid"))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestReadJSON_InvalidCPD(t *testing.T) {
	input := `{
  "name": "test",
  "nodes": ["A"],
  "edges": [],
  "states": {"A": ["s0","s1"]},
  "cpds": {
    "A": {"variable_card": 3, "values": [[0.6],[0.4]]}
  }
}`
	_, err := ReadJSON(strings.NewReader(input))
	if err == nil {
		t.Error("expected error for mismatched variable_card")
	}
}

// ---------------------------------------------------------------------------
// ReadJSONStructure
// ---------------------------------------------------------------------------

func TestReadJSONStructure_Basic(t *testing.T) {
	input := `{
  "name": "test",
  "nodes": ["X", "Y", "Z"],
  "edges": [["X","Y"], ["Y","Z"]],
  "states": {"X": ["a","b"]},
  "cpds": {"X": {"variable_card": 2, "values": [[0.5],[0.5]]}}
}`
	bn, err := ReadJSONStructure(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadJSONStructure failed: %v", err)
	}

	nodes := bn.Nodes()
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}

	edges := bn.Edges()
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(edges))
	}

	// CPDs should NOT be present (structure only).
	if bn.GetCPD("X") != nil {
		t.Error("expected no CPD for X in structure-only read")
	}
}

// ---------------------------------------------------------------------------
// WriteJSON
// ---------------------------------------------------------------------------

func TestWriteJSON_Basic(t *testing.T) {
	bn := buildSampleBN(t)

	var buf bytes.Buffer
	if err := WriteJSON(&buf, bn); err != nil {
		t.Fatalf("WriteJSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"name"`) {
		t.Error("output missing name field")
	}
	if !strings.Contains(output, `"Rain"`) {
		t.Error("output missing Rain")
	}
	if !strings.Contains(output, `"Sprinkler"`) {
		t.Error("output missing Sprinkler")
	}
}

// ---------------------------------------------------------------------------
// JSON Round-trip
// ---------------------------------------------------------------------------

func TestJSON_RoundTrip(t *testing.T) {
	bn := buildSampleBN(t)

	var buf bytes.Buffer
	must(t, WriteJSON(&buf, bn))

	bn2, err := ReadJSON(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadJSON round-trip failed: %v\nJSON:\n%s", err, buf.String())
	}

	// Verify nodes.
	nodes1 := bn.Nodes()
	nodes2 := bn2.Nodes()
	if len(nodes1) != len(nodes2) {
		t.Fatalf("node count mismatch: %d vs %d", len(nodes1), len(nodes2))
	}
	for i := range nodes1 {
		if nodes1[i] != nodes2[i] {
			t.Errorf("node %d: %q vs %q", i, nodes1[i], nodes2[i])
		}
	}

	// Verify edges.
	edges1 := bn.Edges()
	edges2 := bn2.Edges()
	if len(edges1) != len(edges2) {
		t.Fatalf("edge count mismatch: %d vs %d", len(edges1), len(edges2))
	}
	for i := range edges1 {
		if edges1[i] != edges2[i] {
			t.Errorf("edge %d: %v vs %v", i, edges1[i], edges2[i])
		}
	}

	// Verify states.
	for _, node := range nodes1 {
		s1 := bn.GetStates(node)
		s2 := bn2.GetStates(node)
		if len(s1) != len(s2) {
			t.Errorf("states for %q: %v vs %v", node, s1, s2)
			continue
		}
		for j := range s1 {
			if s1[j] != s2[j] {
				t.Errorf("state %q[%d]: %q vs %q", node, j, s1[j], s2[j])
			}
		}
	}

	// Verify CPD values.
	for _, node := range nodes1 {
		cpd1 := bn.GetCPD(node)
		cpd2 := bn2.GetCPD(node)
		if cpd1 == nil || cpd2 == nil {
			t.Errorf("CPD for %q: nil in one of the networks", node)
			continue
		}
		data1 := cpd1.ToFactor().Values().Data()
		data2 := cpd2.ToFactor().Values().Data()
		assertFloatsClose(t, data2, data1, node+" CPD round-trip")
	}
}

func TestJSON_RoundTrip_MultipleParents(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("D"))
	must(t, bn.SetStates("D", []string{"Easy", "Hard"}))
	must(t, bn.AddNode("I"))
	must(t, bn.SetStates("I", []string{"Low", "High"}))
	must(t, bn.AddNode("G"))
	must(t, bn.SetStates("G", []string{"A", "B", "C"}))

	must(t, bn.AddEdge("D", "G"))
	must(t, bn.AddEdge("I", "G"))

	dCPD, err := factors.NewTabularCPD("D", 2, [][]float64{{0.6}, {0.4}}, nil, nil)
	must(t, err)
	must(t, bn.AddCPD(dCPD))

	iCPD, err := factors.NewTabularCPD("I", 2, [][]float64{{0.7}, {0.3}}, nil, nil)
	must(t, err)
	must(t, bn.AddCPD(iCPD))

	gCPD, err := factors.NewTabularCPD("G", 3,
		[][]float64{
			{0.3, 0.05, 0.9, 0.5},
			{0.4, 0.25, 0.08, 0.3},
			{0.3, 0.7, 0.02, 0.2},
		},
		[]string{"D", "I"}, []int{2, 2},
	)
	must(t, err)
	must(t, bn.AddCPD(gCPD))

	var buf bytes.Buffer
	must(t, WriteJSON(&buf, bn))

	bn2, err := ReadJSON(strings.NewReader(buf.String()))
	must(t, err)

	gCPD2 := bn2.GetCPD("G")
	if gCPD2 == nil {
		t.Fatal("G CPD is nil after round-trip")
	}
	data1 := gCPD.ToFactor().Values().Data()
	data2 := gCPD2.ToFactor().Values().Data()
	assertFloatsClose(t, data2, data1, "G CPD round-trip")
}
