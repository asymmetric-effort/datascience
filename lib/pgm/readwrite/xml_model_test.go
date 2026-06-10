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
// ReadXMLNative
// ---------------------------------------------------------------------------

func TestReadXMLNative_Basic(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<pgmgo-network name="example">
  <nodes>
    <node name="A" states="s0,s1"/>
    <node name="B" states="s0,s1"/>
  </nodes>
  <edges>
    <edge from="A" to="B"/>
  </edges>
  <cpds>
    <cpd variable="A" card="2">
      <values>0.6 0.4</values>
    </cpd>
    <cpd variable="B" card="2" evidence="A" evidence_card="2">
      <values>0.3 0.7 0.7 0.3</values>
    </cpd>
  </cpds>
</pgmgo-network>`

	bn, err := ReadXMLNative(strings.NewReader(input))
	if err != nil {
		t.Fatalf("ReadXMLNative failed: %v", err)
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
	// Values "0.3 0.7 0.7 0.3" => pc0: cs0=0.3 cs1=0.7, pc1: cs0=0.7 cs1=0.3
	// flat: [cs0*2+pc0, cs0*2+pc1, cs1*2+pc0, cs1*2+pc1] = [0.3, 0.7, 0.7, 0.3]
	assertFloatsClose(t, bData, []float64{0.3, 0.7, 0.7, 0.3}, "B CPD")
}

func TestReadXMLNative_InvalidXML(t *testing.T) {
	_, err := ReadXMLNative(strings.NewReader("<invalid"))
	if err == nil {
		t.Error("expected error for invalid XML")
	}
}

func TestReadXMLNative_BadValueCount(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<pgmgo-network name="test">
  <nodes>
    <node name="A" states="s0,s1"/>
  </nodes>
  <edges></edges>
  <cpds>
    <cpd variable="A" card="2">
      <values>0.6 0.3 0.1</values>
    </cpd>
  </cpds>
</pgmgo-network>`
	_, err := ReadXMLNative(strings.NewReader(input))
	if err == nil {
		t.Error("expected error for wrong value count")
	}
}

func TestReadXMLNative_EvidenceMismatch(t *testing.T) {
	input := `<?xml version="1.0" encoding="UTF-8"?>
<pgmgo-network name="test">
  <nodes>
    <node name="A" states="s0,s1"/>
    <node name="B" states="s0,s1"/>
  </nodes>
  <edges>
    <edge from="A" to="B"/>
  </edges>
  <cpds>
    <cpd variable="A" card="2">
      <values>0.6 0.4</values>
    </cpd>
    <cpd variable="B" card="2" evidence="A C" evidence_card="2">
      <values>0.3 0.7 0.7 0.3</values>
    </cpd>
  </cpds>
</pgmgo-network>`
	_, err := ReadXMLNative(strings.NewReader(input))
	if err == nil {
		t.Error("expected error for evidence/evidence_card count mismatch")
	}
}

// ---------------------------------------------------------------------------
// WriteXMLNative
// ---------------------------------------------------------------------------

func TestWriteXMLNative_Basic(t *testing.T) {
	bn := buildSampleBN(t)

	var buf bytes.Buffer
	if err := WriteXMLNative(&buf, bn); err != nil {
		t.Fatalf("WriteXMLNative failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "pgmgo-network") {
		t.Error("output missing pgmgo-network element")
	}
	if !strings.Contains(output, `name="Rain"`) {
		t.Error("output missing Rain node")
	}
	if !strings.Contains(output, `name="Sprinkler"`) {
		t.Error("output missing Sprinkler node")
	}
	if !strings.Contains(output, `from="Rain"`) {
		t.Error("output missing edge")
	}
}

// ---------------------------------------------------------------------------
// XML Native Round-trip
// ---------------------------------------------------------------------------

func TestXMLNative_RoundTrip(t *testing.T) {
	bn := buildSampleBN(t)

	var buf bytes.Buffer
	must(t, WriteXMLNative(&buf, bn))

	bn2, err := ReadXMLNative(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadXMLNative round-trip failed: %v\nXML:\n%s", err, buf.String())
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

func TestXMLNative_RoundTrip_MultipleParents(t *testing.T) {
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
	must(t, WriteXMLNative(&buf, bn))

	bn2, err := ReadXMLNative(strings.NewReader(buf.String()))
	must(t, err)

	gCPD2 := bn2.GetCPD("G")
	if gCPD2 == nil {
		t.Fatal("G CPD is nil after round-trip")
	}
	data1 := gCPD.ToFactor().Values().Data()
	data2 := gCPD2.ToFactor().Values().Data()
	assertFloatsClose(t, data2, data1, "G CPD round-trip")
}
