//go:build unit

package readwrite

import (
	"bytes"
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
)

// ---------------------------------------------------------------------------
// ReadCSVStructure
// ---------------------------------------------------------------------------

func TestReadCSVStructure_Basic(t *testing.T) {
	csv := "from,to\nA,B\nB,C\n"
	bn, err := ReadCSVStructure(strings.NewReader(csv))
	if err != nil {
		t.Fatalf("ReadCSVStructure failed: %v", err)
	}

	nodes := bn.Nodes()
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d: %v", len(nodes), nodes)
	}

	edges := bn.Edges()
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d: %v", len(edges), edges)
	}
	if edges[0] != [2]string{"A", "B"} {
		t.Errorf("edge 0 = %v, want [A B]", edges[0])
	}
	if edges[1] != [2]string{"B", "C"} {
		t.Errorf("edge 1 = %v, want [B C]", edges[1])
	}
}

func TestReadCSVStructure_Empty(t *testing.T) {
	_, err := ReadCSVStructure(strings.NewReader(""))
	if err == nil {
		t.Error("expected error for empty CSV")
	}
}

func TestReadCSVStructure_MissingHeader(t *testing.T) {
	_, err := ReadCSVStructure(strings.NewReader("x,y\nA,B\n"))
	if err == nil {
		t.Error("expected error for missing from/to header")
	}
}

func TestReadCSVStructure_HeaderOnly(t *testing.T) {
	bn, err := ReadCSVStructure(strings.NewReader("from,to\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bn.Nodes()) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(bn.Nodes()))
	}
}

// ---------------------------------------------------------------------------
// WriteCSVStructure
// ---------------------------------------------------------------------------

func TestWriteCSVStructure_Basic(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("A"))
	must(t, bn.AddNode("B"))
	must(t, bn.AddEdge("A", "B"))

	var buf bytes.Buffer
	if err := WriteCSVStructure(&buf, bn); err != nil {
		t.Fatalf("WriteCSVStructure failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "from,to") {
		t.Error("output missing header")
	}
	if !strings.Contains(output, "A,B") {
		t.Error("output missing edge A,B")
	}
}

// ---------------------------------------------------------------------------
// CSV Structure Round-trip
// ---------------------------------------------------------------------------

func TestCSVStructure_RoundTrip(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("A"))
	must(t, bn.AddNode("B"))
	must(t, bn.AddNode("C"))
	must(t, bn.AddEdge("A", "B"))
	must(t, bn.AddEdge("B", "C"))

	var buf bytes.Buffer
	must(t, WriteCSVStructure(&buf, bn))

	bn2, err := ReadCSVStructure(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadCSVStructure round-trip failed: %v\nCSV:\n%s", err, buf.String())
	}

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
}

// ---------------------------------------------------------------------------
// ReadCSVCPD / WriteCSVCPD
// ---------------------------------------------------------------------------

func TestCSVCPD_Unconditional(t *testing.T) {
	cpd, err := factors.NewTabularCPD("A", 2, [][]float64{{0.6}, {0.4}}, nil, nil)
	must(t, err)

	var buf bytes.Buffer
	must(t, WriteCSVCPD(&buf, cpd))

	cpd2, err := ReadCSVCPD(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadCSVCPD failed: %v\nCSV:\n%s", err, buf.String())
	}

	data1 := cpd.ToFactor().Values().Data()
	data2 := cpd2.ToFactor().Values().Data()
	assertFloatsClose(t, data2, data1, "unconditional CPD round-trip")
}

func TestCSVCPD_Conditional(t *testing.T) {
	cpd, err := factors.NewTabularCPD("B", 2,
		[][]float64{{0.3, 0.7}, {0.7, 0.3}},
		[]string{"A"}, []int{2})
	must(t, err)

	var buf bytes.Buffer
	must(t, WriteCSVCPD(&buf, cpd))

	cpd2, err := ReadCSVCPD(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadCSVCPD failed: %v\nCSV:\n%s", err, buf.String())
	}

	data1 := cpd.ToFactor().Values().Data()
	data2 := cpd2.ToFactor().Values().Data()
	assertFloatsClose(t, data2, data1, "conditional CPD round-trip")
}

func TestReadCSVCPD_ErrorTooFewRows(t *testing.T) {
	_, err := ReadCSVCPD(strings.NewReader("A,P\n"))
	if err == nil {
		t.Error("expected error for CSV CPD with only header")
	}
}

func TestReadCSVCPD_ErrorInvalidValue(t *testing.T) {
	_, err := ReadCSVCPD(strings.NewReader("A,P\ns0,abc\n"))
	if err == nil {
		t.Error("expected error for invalid probability value")
	}
}
