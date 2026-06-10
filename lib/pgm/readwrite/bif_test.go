//go:build unit

package readwrite

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
)

const sampleBIF = `network unknown {
}

variable Rain {
  type discrete [ 2 ] { True, False };
}

variable Sprinkler {
  type discrete [ 2 ] { True, False };
}

probability ( Rain ) {
  table 0.2, 0.8;
}

probability ( Sprinkler | Rain ) {
  (True) 0.01, 0.99;
  (False) 0.4, 0.6;
}
`

// ---------------------------------------------------------------------------
// ReadBIF
// ---------------------------------------------------------------------------

func TestReadBIF_Basic(t *testing.T) {
	bn, err := ReadBIF(strings.NewReader(sampleBIF))
	if err != nil {
		t.Fatalf("ReadBIF failed: %v", err)
	}

	// Check nodes.
	nodes := bn.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d: %v", len(nodes), nodes)
	}

	// Check states.
	rainStates := bn.GetStates("Rain")
	if len(rainStates) != 2 || rainStates[0] != "True" || rainStates[1] != "False" {
		t.Errorf("Rain states = %v, want [True, False]", rainStates)
	}

	sprinklerStates := bn.GetStates("Sprinkler")
	if len(sprinklerStates) != 2 || sprinklerStates[0] != "True" || sprinklerStates[1] != "False" {
		t.Errorf("Sprinkler states = %v, want [True, False]", sprinklerStates)
	}

	// Check edges.
	edges := bn.Edges()
	if len(edges) != 1 {
		t.Fatalf("expected 1 edge, got %d: %v", len(edges), edges)
	}
	if edges[0] != [2]string{"Rain", "Sprinkler"} {
		t.Errorf("edge = %v, want [Rain, Sprinkler]", edges[0])
	}

	// Check Rain CPD (unconditional).
	rainCPD := bn.GetCPD("Rain")
	if rainCPD == nil {
		t.Fatal("Rain CPD is nil")
	}
	rainData := rainCPD.ToFactor().Values().Data()
	assertFloatsClose(t, rainData, []float64{0.2, 0.8}, "Rain CPD values")

	// Check Sprinkler CPD (conditional on Rain).
	sprCPD := bn.GetCPD("Sprinkler")
	if sprCPD == nil {
		t.Fatal("Sprinkler CPD is nil")
	}
	sprEvidence := sprCPD.Evidence()
	if len(sprEvidence) != 1 || sprEvidence[0] != "Rain" {
		t.Errorf("Sprinkler evidence = %v, want [Rain]", sprEvidence)
	}
	// Data layout: [childState * numParentConfigs + parentConfig]
	// child=True(0): Rain=True(0)=>0.01, Rain=False(1)=>0.4
	// child=False(1): Rain=True(0)=>0.99, Rain=False(1)=>0.6
	sprData := sprCPD.ToFactor().Values().Data()
	assertFloatsClose(t, sprData, []float64{0.01, 0.4, 0.99, 0.6}, "Sprinkler CPD values")
}

func TestReadBIF_Comments(t *testing.T) {
	bif := `// This is a comment
network unknown { // inline comment
}

variable X {
  type discrete [ 2 ] { A, B }; // state comment
}

// Another comment
probability ( X ) {
  table 0.3, 0.7; // prob comment
}
`
	bn, err := ReadBIF(strings.NewReader(bif))
	if err != nil {
		t.Fatalf("ReadBIF with comments failed: %v", err)
	}

	states := bn.GetStates("X")
	if len(states) != 2 || states[0] != "A" || states[1] != "B" {
		t.Errorf("X states = %v, want [A, B]", states)
	}

	cpd := bn.GetCPD("X")
	data := cpd.ToFactor().Values().Data()
	assertFloatsClose(t, data, []float64{0.3, 0.7}, "X CPD values")
}

func TestReadBIF_ThreeVarChain(t *testing.T) {
	bif := `network test {
}

variable A {
  type discrete [ 2 ] { a0, a1 };
}

variable B {
  type discrete [ 2 ] { b0, b1 };
}

variable C {
  type discrete [ 3 ] { c0, c1, c2 };
}

probability ( A ) {
  table 0.6, 0.4;
}

probability ( B | A ) {
  (a0) 0.2, 0.8;
  (a1) 0.75, 0.25;
}

probability ( C | B ) {
  (b0) 0.1, 0.2, 0.7;
  (b1) 0.5, 0.3, 0.2;
}
`
	bn, err := ReadBIF(strings.NewReader(bif))
	if err != nil {
		t.Fatalf("ReadBIF failed: %v", err)
	}

	nodes := bn.Nodes()
	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}

	edges := bn.Edges()
	if len(edges) != 2 {
		t.Fatalf("expected 2 edges, got %d", len(edges))
	}

	// C has 3 states.
	cStates := bn.GetStates("C")
	if len(cStates) != 3 {
		t.Errorf("C states length = %d, want 3", len(cStates))
	}

	// B CPD: B|A
	bCPD := bn.GetCPD("B")
	bData := bCPD.ToFactor().Values().Data()
	// child=b0(0): A=a0(0)=>0.2, A=a1(1)=>0.75
	// child=b1(1): A=a0(0)=>0.8, A=a1(1)=>0.25
	assertFloatsClose(t, bData, []float64{0.2, 0.75, 0.8, 0.25}, "B CPD values")

	// C CPD: C|B
	cCPD := bn.GetCPD("C")
	cData := cCPD.ToFactor().Values().Data()
	// child=c0(0): B=b0(0)=>0.1, B=b1(1)=>0.5
	// child=c1(1): B=b0(0)=>0.2, B=b1(1)=>0.3
	// child=c2(2): B=b0(0)=>0.7, B=b1(1)=>0.2
	assertFloatsClose(t, cData, []float64{0.1, 0.5, 0.2, 0.3, 0.7, 0.2}, "C CPD values")
}

func TestReadBIF_MultipleParents(t *testing.T) {
	bif := `network test {
}

variable D {
  type discrete [ 2 ] { d0, d1 };
}

variable I {
  type discrete [ 2 ] { i0, i1 };
}

variable G {
  type discrete [ 3 ] { g0, g1, g2 };
}

probability ( D ) {
  table 0.6, 0.4;
}

probability ( I ) {
  table 0.7, 0.3;
}

probability ( G | D, I ) {
  (d0, i0) 0.3, 0.4, 0.3;
  (d0, i1) 0.05, 0.25, 0.7;
  (d1, i0) 0.9, 0.08, 0.02;
  (d1, i1) 0.5, 0.3, 0.2;
}
`
	bn, err := ReadBIF(strings.NewReader(bif))
	if err != nil {
		t.Fatalf("ReadBIF failed: %v", err)
	}

	gCPD := bn.GetCPD("G")
	if gCPD == nil {
		t.Fatal("G CPD is nil")
	}
	ev := gCPD.Evidence()
	if len(ev) != 2 {
		t.Fatalf("G evidence = %v, want 2 parents", ev)
	}

	// Parent configs (D x I): (d0,i0)=0, (d0,i1)=1, (d1,i0)=2, (d1,i1)=3
	// Data: [g0: 0.3, 0.05, 0.9, 0.5, g1: 0.4, 0.25, 0.08, 0.3, g2: 0.3, 0.7, 0.02, 0.2]
	gData := gCPD.ToFactor().Values().Data()
	expected := []float64{0.3, 0.05, 0.9, 0.5, 0.4, 0.25, 0.08, 0.3, 0.3, 0.7, 0.02, 0.2}
	assertFloatsClose(t, gData, expected, "G CPD values")
}

// ---------------------------------------------------------------------------
// WriteBIF
// ---------------------------------------------------------------------------

func TestWriteBIF_Basic(t *testing.T) {
	bn := buildSampleBN(t)

	var buf bytes.Buffer
	if err := WriteBIF(&buf, bn); err != nil {
		t.Fatalf("WriteBIF failed: %v", err)
	}

	output := buf.String()

	// Check that output contains key sections.
	if !strings.Contains(output, "network unknown") {
		t.Error("output missing network header")
	}
	if !strings.Contains(output, "variable Rain") {
		t.Error("output missing variable Rain")
	}
	if !strings.Contains(output, "variable Sprinkler") {
		t.Error("output missing variable Sprinkler")
	}
	if !strings.Contains(output, "probability ( Rain )") {
		t.Error("output missing probability Rain")
	}
	if !strings.Contains(output, "probability ( Sprinkler | Rain )") {
		t.Error("output missing probability Sprinkler|Rain")
	}
	if !strings.Contains(output, "table 0.2, 0.8;") {
		t.Errorf("output missing Rain table values, got:\n%s", output)
	}
}

// ---------------------------------------------------------------------------
// Round-trip
// ---------------------------------------------------------------------------

func TestBIF_RoundTrip(t *testing.T) {
	bn := buildSampleBN(t)

	// Write.
	var buf bytes.Buffer
	if err := WriteBIF(&buf, bn); err != nil {
		t.Fatalf("WriteBIF failed: %v", err)
	}

	// Read back.
	bn2, err := ReadBIF(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadBIF failed on round-trip: %v\nBIF output:\n%s", err, buf.String())
	}

	// Verify structure.
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

func TestBIF_RoundTrip_MultipleParents(t *testing.T) {
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

	// G | D, I: parents in order [D, I], parent configs: (Easy,Low)=0, (Easy,High)=1, (Hard,Low)=2, (Hard,High)=3
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

	// Write.
	var buf bytes.Buffer
	must(t, WriteBIF(&buf, bn))

	// Read back.
	bn2, err := ReadBIF(strings.NewReader(buf.String()))
	must(t, err)

	// Verify G CPD.
	gCPD2 := bn2.GetCPD("G")
	if gCPD2 == nil {
		t.Fatal("G CPD is nil after round-trip")
	}
	data1 := gCPD.ToFactor().Values().Data()
	data2 := gCPD2.ToFactor().Values().Data()
	assertFloatsClose(t, data2, data1, "G CPD round-trip")
}

func TestReadBIF_Empty(t *testing.T) {
	bif := `network empty {
}
`
	bn, err := ReadBIF(strings.NewReader(bif))
	if err != nil {
		t.Fatalf("ReadBIF failed: %v", err)
	}
	if len(bn.Nodes()) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(bn.Nodes()))
	}
}

func TestReadBIF_MalformedVariable(t *testing.T) {
	bif := `network test {
}

variable {
  type discrete [ 2 ] { A, B };
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for malformed variable declaration")
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func buildSampleBN(t *testing.T) *models.BayesianNetwork {
	t.Helper()
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("Rain"))
	must(t, bn.SetStates("Rain", []string{"True", "False"}))
	must(t, bn.AddNode("Sprinkler"))
	must(t, bn.SetStates("Sprinkler", []string{"True", "False"}))
	must(t, bn.AddEdge("Rain", "Sprinkler"))

	rainCPD, err := factors.NewTabularCPD("Rain", 2,
		[][]float64{{0.2}, {0.8}}, nil, nil)
	must(t, err)
	must(t, bn.AddCPD(rainCPD))

	sprCPD, err := factors.NewTabularCPD("Sprinkler", 2,
		[][]float64{{0.01, 0.4}, {0.99, 0.6}},
		[]string{"Rain"}, []int{2})
	must(t, err)
	must(t, bn.AddCPD(sprCPD))

	return bn
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func assertFloatsClose(t *testing.T, got, want []float64, msg string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: length mismatch: got %d, want %d\ngot:  %v\nwant: %v",
			msg, len(got), len(want), got, want)
	}
	for i := range got {
		if math.Abs(got[i]-want[i]) > 1e-9 {
			t.Errorf("%s: index %d: got %f, want %f", msg, i, got[i], want[i])
		}
	}
}
