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
// PomdpX — Read tests
// ---------------------------------------------------------------------------

func TestReadPomdpX_FullConditional(t *testing.T) {
	// PomdpX with Variable definitions, InitialStateBelief, and
	// StateTransitionFunction containing conditional Entry elements.
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="Rain" numValues="2">
      <ValueEnum>True False</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="Sprinkler" numValues="2">
      <ValueEnum>On Off</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>Rain</Var>
      <Parameter type="TBL">
        <Entry>
          <ProbTable>0.2 0.8</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>Sprinkler</Var>
      <Parent>Rain</Parent>
      <Parameter type="TBL">
        <Entry>
          <Instance>True</Instance>
          <ProbTable>0.01 0.99</ProbTable>
        </Entry>
        <Entry>
          <Instance>False</Instance>
          <ProbTable>0.4 0.6</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`

	bn, err := ReadPomdpX(strings.NewReader(xmlData))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v", err)
	}

	// Check nodes.
	nodes := bn.Nodes()
	if len(nodes) != 2 {
		t.Fatalf("expected 2 nodes, got %d: %v", len(nodes), nodes)
	}

	// Check Rain CPD.
	rainCPD := bn.GetCPD("Rain")
	if rainCPD == nil {
		t.Fatal("Rain CPD is nil")
	}
	rainData := rainCPD.ToFactor().Values().Data()
	assertFloatsClose(t, rainData, []float64{0.2, 0.8}, "Rain CPD")

	// Check Sprinkler CPD.
	sprCPD := bn.GetCPD("Sprinkler")
	if sprCPD == nil {
		t.Fatal("Sprinkler CPD is nil")
	}
	sprEv := sprCPD.Evidence()
	if len(sprEv) != 1 || sprEv[0] != "Rain" {
		t.Fatalf("Sprinkler evidence = %v, want [Rain]", sprEv)
	}
	// data layout: [childState*numParentConfigs + parentConfig]
	// On: Rain=True=>0.01, Rain=False=>0.4
	// Off: Rain=True=>0.99, Rain=False=>0.6
	sprData := sprCPD.ToFactor().Values().Data()
	assertFloatsClose(t, sprData, []float64{0.01, 0.4, 0.99, 0.6}, "Sprinkler CPD")

	// Check edge.
	edges := bn.Edges()
	if len(edges) != 1 || edges[0] != [2]string{"Rain", "Sprinkler"} {
		t.Errorf("expected edge [Rain Sprinkler], got %v", edges)
	}
}

func TestReadPomdpX_MultiParentConditional(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="D" numValues="2">
      <ValueEnum>Easy Hard</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="I" numValues="2">
      <ValueEnum>Low High</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="G" numValues="3">
      <ValueEnum>A B C</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>D</Var>
      <Parameter type="TBL">
        <Entry>
          <ProbTable>0.6 0.4</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
    <CondProb>
      <Var>I</Var>
      <Parameter type="TBL">
        <Entry>
          <ProbTable>0.7 0.3</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>G</Var>
      <Parent>D I</Parent>
      <Parameter type="TBL">
        <Entry>
          <Instance>Easy Low</Instance>
          <ProbTable>0.3 0.4 0.3</ProbTable>
        </Entry>
        <Entry>
          <Instance>Easy High</Instance>
          <ProbTable>0.05 0.25 0.7</ProbTable>
        </Entry>
        <Entry>
          <Instance>Hard Low</Instance>
          <ProbTable>0.9 0.08 0.02</ProbTable>
        </Entry>
        <Entry>
          <Instance>Hard High</Instance>
          <ProbTable>0.5 0.3 0.2</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`

	bn, err := ReadPomdpX(strings.NewReader(xmlData))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v", err)
	}

	gCPD := bn.GetCPD("G")
	if gCPD == nil {
		t.Fatal("G CPD is nil")
	}
	gData := gCPD.ToFactor().Values().Data()
	// Parent configs (D x I): (Easy,Low)=0, (Easy,High)=1, (Hard,Low)=2, (Hard,High)=3
	expected := []float64{0.3, 0.05, 0.9, 0.5, 0.4, 0.25, 0.08, 0.3, 0.3, 0.7, 0.02, 0.2}
	assertFloatsClose(t, gData, expected, "G CPD")
}

func TestReadPomdpX_FlatTable(t *testing.T) {
	// Flat table without Instance tags (all probs in one ProbTable).
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="A" numValues="2">
      <ValueEnum>a0 a1</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="B" numValues="2">
      <ValueEnum>b0 b1</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>A</Var>
      <Parameter type="TBL">
        <Entry>
          <ProbTable>0.6 0.4</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>B</Var>
      <Parent>A</Parent>
      <Parameter type="TBL">
        <Entry>
          <ProbTable>0.2 0.8 0.75 0.25</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`

	bn, err := ReadPomdpX(strings.NewReader(xmlData))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v", err)
	}

	bCPD := bn.GetCPD("B")
	if bCPD == nil {
		t.Fatal("B CPD is nil")
	}
	bData := bCPD.ToFactor().Values().Data()
	// Flat: pc=0(a0): b0=0.2,b1=0.8; pc=1(a1): b0=0.75,b1=0.25
	assertFloatsClose(t, bData, []float64{0.2, 0.75, 0.8, 0.25}, "B CPD flat table")
}

func TestReadPomdpX_UnknownParent(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="X" numValues="2">
      <ValueEnum>a b</ValueEnum>
    </StateVar>
  </Variable>
  <StateTransitionFunction>
    <CondProb>
      <Var>X</Var>
      <Parent>Unknown</Parent>
      <Parameter type="TBL">
        <Entry>
          <Instance>s0</Instance>
          <ProbTable>0.5 0.5</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`
	_, err := ReadPomdpX(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for unknown parent")
	}
}

func TestReadPomdpX_WrongInstanceCount(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="A" numValues="2">
      <ValueEnum>a0 a1</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="B" numValues="2">
      <ValueEnum>b0 b1</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>A</Var>
      <Parameter><Entry><ProbTable>0.5 0.5</ProbTable></Entry></Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>B</Var>
      <Parent>A</Parent>
      <Parameter type="TBL">
        <Entry>
          <Instance>a0 extra</Instance>
          <ProbTable>0.5 0.5</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`
	_, err := ReadPomdpX(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for wrong instance part count")
	}
}

func TestReadPomdpX_WrongProbCount(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="A" numValues="2">
      <ValueEnum>a0 a1</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="B" numValues="2">
      <ValueEnum>b0 b1</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>A</Var>
      <Parameter><Entry><ProbTable>0.5 0.5</ProbTable></Entry></Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>B</Var>
      <Parent>A</Parent>
      <Parameter type="TBL">
        <Entry>
          <Instance>a0</Instance>
          <ProbTable>0.5 0.3 0.2</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`
	_, err := ReadPomdpX(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for wrong probability count in entry")
	}
}

func TestReadPomdpX_UnknownParentState(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="A" numValues="2">
      <ValueEnum>a0 a1</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="B" numValues="2">
      <ValueEnum>b0 b1</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>A</Var>
      <Parameter><Entry><ProbTable>0.5 0.5</ProbTable></Entry></Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>B</Var>
      <Parent>A</Parent>
      <Parameter type="TBL">
        <Entry>
          <Instance>unknown</Instance>
          <ProbTable>0.5 0.5</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`
	_, err := ReadPomdpX(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for unknown parent state in instance")
	}
}

func TestReadPomdpX_FlatTableWrongCount(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="A" numValues="2">
      <ValueEnum>a0 a1</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="B" numValues="2">
      <ValueEnum>b0 b1</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>A</Var>
      <Parameter><Entry><ProbTable>0.5 0.5</ProbTable></Entry></Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>B</Var>
      <Parent>A</Parent>
      <Parameter type="TBL">
        <Entry>
          <ProbTable>0.5 0.5 0.5</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`
	_, err := ReadPomdpX(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for wrong flat table size")
	}
}

func TestReadPomdpX_BadProbInConditional(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="A" numValues="2">
      <ValueEnum>a0 a1</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="B" numValues="2">
      <ValueEnum>b0 b1</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>A</Var>
      <Parameter><Entry><ProbTable>0.5 0.5</ProbTable></Entry></Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>B</Var>
      <Parent>A</Parent>
      <Parameter type="TBL">
        <Entry>
          <Instance>a0</Instance>
          <ProbTable>abc 0.5</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`
	_, err := ReadPomdpX(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for bad float in conditional entry")
	}
}

func TestReadPomdpX_BadProbInFlatTable(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="A" numValues="2">
      <ValueEnum>a0 a1</ValueEnum>
    </StateVar>
    <StateVar vnamePrev="B" numValues="2">
      <ValueEnum>b0 b1</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>A</Var>
      <Parameter><Entry><ProbTable>0.5 0.5</ProbTable></Entry></Parameter>
    </CondProb>
  </InitialStateBelief>
  <StateTransitionFunction>
    <CondProb>
      <Var>B</Var>
      <Parent>A</Parent>
      <Parameter type="TBL">
        <Entry>
          <ProbTable>abc 0.5 0.3 0.2</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </StateTransitionFunction>
</pomdpx>`
	_, err := ReadPomdpX(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for bad float in flat table")
	}
}

// ---------------------------------------------------------------------------
// PomdpX — Write tests
// ---------------------------------------------------------------------------

func TestWritePomdpX_ConditionalOutput(t *testing.T) {
	bn := buildThreeNodeBN(t)

	var buf bytes.Buffer
	if err := WritePomdpX(&buf, bn); err != nil {
		t.Fatalf("WritePomdpX failed: %v", err)
	}

	output := buf.String()

	// Should contain StateTransitionFunction for the conditional node.
	if !strings.Contains(output, "StateTransitionFunction") {
		t.Error("output missing StateTransitionFunction section")
	}
	if !strings.Contains(output, "InitialStateBelief") {
		t.Error("output missing InitialStateBelief section")
	}
	// Should have Instance tags.
	if !strings.Contains(output, "<Instance>") {
		t.Error("output missing Instance tags for conditional entries")
	}
	// Should have Parent tag.
	if !strings.Contains(output, "<Parent>") {
		t.Error("output missing Parent tags")
	}
}

// ---------------------------------------------------------------------------
// PomdpX — Full round-trip tests
// ---------------------------------------------------------------------------

func TestPomdpX_FullRoundTrip_ThreeNode(t *testing.T) {
	bn := buildThreeNodeBN(t)

	var buf bytes.Buffer
	if err := WritePomdpX(&buf, bn); err != nil {
		t.Fatalf("WritePomdpX failed: %v", err)
	}

	bn2, err := ReadPomdpX(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v\nOutput:\n%s", err, buf.String())
	}

	assertBNEqual(t, bn2, bn, "PomdpX full round-trip three-node")
}

func TestPomdpX_FullRoundTrip_MultiParent(t *testing.T) {
	bn := buildMultiParentBN(t)

	var buf bytes.Buffer
	if err := WritePomdpX(&buf, bn); err != nil {
		t.Fatalf("WritePomdpX failed: %v", err)
	}

	bn2, err := ReadPomdpX(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v\nOutput:\n%s", err, buf.String())
	}

	assertBNEqual(t, bn2, bn, "PomdpX full round-trip multi-parent")
}

func TestPomdpX_DoubleRoundTrip(t *testing.T) {
	bn := buildMultiParentBN(t)

	// First round-trip.
	var buf1 bytes.Buffer
	must(t, WritePomdpX(&buf1, bn))
	bn2, err := ReadPomdpX(strings.NewReader(buf1.String()))
	if err != nil {
		t.Fatalf("first ReadPomdpX failed: %v", err)
	}

	// Second round-trip.
	var buf2 bytes.Buffer
	must(t, WritePomdpX(&buf2, bn2))
	bn3, err := ReadPomdpX(strings.NewReader(buf2.String()))
	if err != nil {
		t.Fatalf("second ReadPomdpX failed: %v", err)
	}

	assertBNEqual(t, bn3, bn, "PomdpX double round-trip")
}

func TestPomdpX_RoundTrip_ThreeStates(t *testing.T) {
	// Test with variable cardinality > 2.
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("Color"))
	must(t, bn.SetStates("Color", []string{"Red", "Green", "Blue"}))
	must(t, bn.AddNode("Shade"))
	must(t, bn.SetStates("Shade", []string{"Light", "Dark"}))
	must(t, bn.AddEdge("Color", "Shade"))

	colorCPD, err := factors.NewTabularCPD("Color", 3,
		[][]float64{{0.5}, {0.3}, {0.2}}, nil, nil)
	must(t, err)
	must(t, bn.AddCPD(colorCPD))

	shadeCPD, err := factors.NewTabularCPD("Shade", 2,
		[][]float64{{0.7, 0.4, 0.1}, {0.3, 0.6, 0.9}},
		[]string{"Color"}, []int{3})
	must(t, err)
	must(t, bn.AddCPD(shadeCPD))

	var buf bytes.Buffer
	must(t, WritePomdpX(&buf, bn))

	bn2, err := ReadPomdpX(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v\nOutput:\n%s", err, buf.String())
	}

	assertBNEqual(t, bn2, bn, "PomdpX three-state round-trip")
}

// ---------------------------------------------------------------------------
// XBN — Read tests
// ---------------------------------------------------------------------------

func TestReadXBN_FullConditional(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<ANALYSISNOTEBOOK>
  <BNMODEL NAME="test">
    <STATICPROPERTIES>
      <NODELIST>
        <NODE NAME="Rain">
          <STATENAME>True</STATENAME>
          <STATENAME>False</STATENAME>
        </NODE>
        <NODE NAME="Sprinkler">
          <STATENAME>On</STATENAME>
          <STATENAME>Off</STATENAME>
        </NODE>
      </NODELIST>
      <ARCLIST>
        <ARC PARENT="Rain" CHILD="Sprinkler"/>
      </ARCLIST>
    </STATICPROPERTIES>
    <DYNAMICPROPERTIES>
      <DISTRIBS>
        <DIST TYPE="discrete">
          <DPIS>
            <DPI>0.2 0.8</DPI>
          </DPIS>
        </DIST>
        <DIST TYPE="discrete">
          <CONDSET>
            <CONDELEM NAME="Rain"/>
          </CONDSET>
          <DPIS>
            <DPI INDEXES="0">0.01 0.99</DPI>
            <DPI INDEXES="1">0.4 0.6</DPI>
          </DPIS>
        </DIST>
      </DISTRIBS>
    </DYNAMICPROPERTIES>
  </BNMODEL>
</ANALYSISNOTEBOOK>`

	bn, err := ReadXBN(strings.NewReader(xmlData))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v", err)
	}

	// Check Rain.
	rainCPD := bn.GetCPD("Rain")
	if rainCPD == nil {
		t.Fatal("Rain CPD is nil")
	}
	rainData := rainCPD.ToFactor().Values().Data()
	assertFloatsClose(t, rainData, []float64{0.2, 0.8}, "Rain CPD")

	// Check Sprinkler.
	sprCPD := bn.GetCPD("Sprinkler")
	if sprCPD == nil {
		t.Fatal("Sprinkler CPD is nil")
	}
	sprData := sprCPD.ToFactor().Values().Data()
	assertFloatsClose(t, sprData, []float64{0.01, 0.4, 0.99, 0.6}, "Sprinkler CPD")
}

func TestReadXBN_MultiParentIndexes(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<ANALYSISNOTEBOOK>
  <BNMODEL NAME="test">
    <STATICPROPERTIES>
      <NODELIST>
        <NODE NAME="D">
          <STATENAME>Easy</STATENAME>
          <STATENAME>Hard</STATENAME>
        </NODE>
        <NODE NAME="I">
          <STATENAME>Low</STATENAME>
          <STATENAME>High</STATENAME>
        </NODE>
        <NODE NAME="G">
          <STATENAME>A</STATENAME>
          <STATENAME>B</STATENAME>
          <STATENAME>C</STATENAME>
        </NODE>
      </NODELIST>
      <ARCLIST>
        <ARC PARENT="D" CHILD="G"/>
        <ARC PARENT="I" CHILD="G"/>
      </ARCLIST>
    </STATICPROPERTIES>
    <DYNAMICPROPERTIES>
      <DISTRIBS>
        <DIST TYPE="discrete">
          <DPIS>
            <DPI>0.6 0.4</DPI>
          </DPIS>
        </DIST>
        <DIST TYPE="discrete">
          <DPIS>
            <DPI>0.7 0.3</DPI>
          </DPIS>
        </DIST>
        <DIST TYPE="discrete">
          <CONDSET>
            <CONDELEM NAME="D"/>
            <CONDELEM NAME="I"/>
          </CONDSET>
          <DPIS>
            <DPI INDEXES="0 0">0.3 0.4 0.3</DPI>
            <DPI INDEXES="0 1">0.05 0.25 0.7</DPI>
            <DPI INDEXES="1 0">0.9 0.08 0.02</DPI>
            <DPI INDEXES="1 1">0.5 0.3 0.2</DPI>
          </DPIS>
        </DIST>
      </DISTRIBS>
    </DYNAMICPROPERTIES>
  </BNMODEL>
</ANALYSISNOTEBOOK>`

	bn, err := ReadXBN(strings.NewReader(xmlData))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v", err)
	}

	gCPD := bn.GetCPD("G")
	if gCPD == nil {
		t.Fatal("G CPD is nil")
	}
	gData := gCPD.ToFactor().Values().Data()
	expected := []float64{0.3, 0.05, 0.9, 0.5, 0.4, 0.25, 0.08, 0.3, 0.3, 0.7, 0.02, 0.2}
	assertFloatsClose(t, gData, expected, "G CPD multi-parent with INDEXES")
}

func TestReadXBN_NoIndexesFallback(t *testing.T) {
	// DPI without INDEXES attribute, values listed sequentially.
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<ANALYSISNOTEBOOK>
  <BNMODEL NAME="test">
    <STATICPROPERTIES>
      <NODELIST>
        <NODE NAME="A">
          <STATENAME>a0</STATENAME>
          <STATENAME>a1</STATENAME>
        </NODE>
        <NODE NAME="B">
          <STATENAME>b0</STATENAME>
          <STATENAME>b1</STATENAME>
        </NODE>
      </NODELIST>
      <ARCLIST>
        <ARC PARENT="A" CHILD="B"/>
      </ARCLIST>
    </STATICPROPERTIES>
    <DYNAMICPROPERTIES>
      <DISTRIBS>
        <DIST TYPE="discrete">
          <DPIS>
            <DPI>0.6 0.4</DPI>
          </DPIS>
        </DIST>
        <DIST TYPE="discrete">
          <CONDSET>
            <CONDELEM NAME="A"/>
          </CONDSET>
          <DPIS>
            <DPI>0.2 0.8</DPI>
            <DPI>0.75 0.25</DPI>
          </DPIS>
        </DIST>
      </DISTRIBS>
    </DYNAMICPROPERTIES>
  </BNMODEL>
</ANALYSISNOTEBOOK>`

	bn, err := ReadXBN(strings.NewReader(xmlData))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v", err)
	}

	bCPD := bn.GetCPD("B")
	if bCPD == nil {
		t.Fatal("B CPD is nil")
	}
	bData := bCPD.ToFactor().Values().Data()
	assertFloatsClose(t, bData, []float64{0.2, 0.75, 0.8, 0.25}, "B CPD no-indexes")
}

func TestReadXBN_BadDPIFloat(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<ANALYSISNOTEBOOK>
  <BNMODEL NAME="test">
    <STATICPROPERTIES>
      <NODELIST>
        <NODE NAME="X">
          <STATENAME>a</STATENAME>
          <STATENAME>b</STATENAME>
        </NODE>
      </NODELIST>
      <ARCLIST></ARCLIST>
    </STATICPROPERTIES>
    <DYNAMICPROPERTIES>
      <DISTRIBS>
        <DIST TYPE="discrete">
          <DPIS>
            <DPI>abc 0.5</DPI>
          </DPIS>
        </DIST>
      </DISTRIBS>
    </DYNAMICPROPERTIES>
  </BNMODEL>
</ANALYSISNOTEBOOK>`
	_, err := ReadXBN(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for bad float in DPI")
	}
}

func TestReadXBN_ConditionalBadFloat(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<ANALYSISNOTEBOOK>
  <BNMODEL NAME="test">
    <STATICPROPERTIES>
      <NODELIST>
        <NODE NAME="A"><STATENAME>a0</STATENAME><STATENAME>a1</STATENAME></NODE>
        <NODE NAME="B"><STATENAME>b0</STATENAME><STATENAME>b1</STATENAME></NODE>
      </NODELIST>
      <ARCLIST><ARC PARENT="A" CHILD="B"/></ARCLIST>
    </STATICPROPERTIES>
    <DYNAMICPROPERTIES>
      <DISTRIBS>
        <DIST TYPE="discrete">
          <DPIS><DPI>0.5 0.5</DPI></DPIS>
        </DIST>
        <DIST TYPE="discrete">
          <CONDSET><CONDELEM NAME="A"/></CONDSET>
          <DPIS>
            <DPI INDEXES="0">abc 0.5</DPI>
            <DPI INDEXES="1">0.3 0.7</DPI>
          </DPIS>
        </DIST>
      </DISTRIBS>
    </DYNAMICPROPERTIES>
  </BNMODEL>
</ANALYSISNOTEBOOK>`
	_, err := ReadXBN(strings.NewReader(xmlData))
	if err == nil {
		t.Error("expected error for bad float in conditional DPI")
	}
}

// ---------------------------------------------------------------------------
// XBN — Write tests
// ---------------------------------------------------------------------------

func TestWriteXBN_HasIndexes(t *testing.T) {
	bn := buildThreeNodeBN(t)

	var buf bytes.Buffer
	if err := WriteXBN(&buf, bn); err != nil {
		t.Fatalf("WriteXBN failed: %v", err)
	}

	output := buf.String()

	// Conditional distributions should have INDEXES attributes.
	if !strings.Contains(output, "INDEXES=") {
		t.Error("output missing INDEXES attribute in DPI elements")
	}
	// Should have CONDSET/CONDELEM for conditional nodes.
	if !strings.Contains(output, "CONDELEM") {
		t.Error("output missing CONDELEM elements")
	}
	// Should have STATENAME elements.
	if !strings.Contains(output, "STATENAME") {
		t.Error("output missing STATENAME elements")
	}
}

// ---------------------------------------------------------------------------
// XBN — Full round-trip tests
// ---------------------------------------------------------------------------

func TestXBN_FullRoundTrip_ThreeNode(t *testing.T) {
	bn := buildThreeNodeBN(t)

	var buf bytes.Buffer
	must(t, WriteXBN(&buf, bn))

	bn2, err := ReadXBN(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v\nOutput:\n%s", err, buf.String())
	}

	assertBNEqual(t, bn2, bn, "XBN full round-trip three-node")
}

func TestXBN_FullRoundTrip_MultiParent(t *testing.T) {
	bn := buildMultiParentBN(t)

	var buf bytes.Buffer
	must(t, WriteXBN(&buf, bn))

	bn2, err := ReadXBN(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v\nOutput:\n%s", err, buf.String())
	}

	assertBNEqual(t, bn2, bn, "XBN full round-trip multi-parent")
}

func TestXBN_DoubleRoundTrip(t *testing.T) {
	bn := buildMultiParentBN(t)

	var buf1 bytes.Buffer
	must(t, WriteXBN(&buf1, bn))
	bn2, err := ReadXBN(strings.NewReader(buf1.String()))
	if err != nil {
		t.Fatalf("first ReadXBN failed: %v", err)
	}

	var buf2 bytes.Buffer
	must(t, WriteXBN(&buf2, bn2))
	bn3, err := ReadXBN(strings.NewReader(buf2.String()))
	if err != nil {
		t.Fatalf("second ReadXBN failed: %v", err)
	}

	assertBNEqual(t, bn3, bn, "XBN double round-trip")
}

func TestXBN_RoundTrip_ThreeStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("Color"))
	must(t, bn.SetStates("Color", []string{"Red", "Green", "Blue"}))
	must(t, bn.AddNode("Shade"))
	must(t, bn.SetStates("Shade", []string{"Light", "Dark"}))
	must(t, bn.AddEdge("Color", "Shade"))

	colorCPD, err := factors.NewTabularCPD("Color", 3,
		[][]float64{{0.5}, {0.3}, {0.2}}, nil, nil)
	must(t, err)
	must(t, bn.AddCPD(colorCPD))

	shadeCPD, err := factors.NewTabularCPD("Shade", 2,
		[][]float64{{0.7, 0.4, 0.1}, {0.3, 0.6, 0.9}},
		[]string{"Color"}, []int{3})
	must(t, err)
	must(t, bn.AddCPD(shadeCPD))

	var buf bytes.Buffer
	must(t, WriteXBN(&buf, bn))

	bn2, err := ReadXBN(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v\nOutput:\n%s", err, buf.String())
	}

	assertBNEqual(t, bn2, bn, "XBN three-state round-trip")
}

// ---------------------------------------------------------------------------
// Cross-format tests — PomdpX and XBN agree on the same BN
// ---------------------------------------------------------------------------

func TestCrossFormat_PomdpX_XBN(t *testing.T) {
	bn := buildMultiParentBN(t)

	// Write to PomdpX, read back.
	var bufP bytes.Buffer
	must(t, WritePomdpX(&bufP, bn))
	bnP, err := ReadPomdpX(strings.NewReader(bufP.String()))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v", err)
	}

	// Write to XBN, read back.
	var bufX bytes.Buffer
	must(t, WriteXBN(&bufX, bn))
	bnX, err := ReadXBN(strings.NewReader(bufX.String()))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v", err)
	}

	// Both should match the original.
	assertBNEqual(t, bnP, bn, "PomdpX cross-format")
	assertBNEqual(t, bnX, bn, "XBN cross-format")
}

// ---------------------------------------------------------------------------
// Edge case: single node, no parents, no edges
// ---------------------------------------------------------------------------

func TestPomdpX_RoundTrip_SingleNode(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("Solo"))
	must(t, bn.SetStates("Solo", []string{"a", "b", "c"}))

	cpd, err := factors.NewTabularCPD("Solo", 3,
		[][]float64{{0.1}, {0.3}, {0.6}}, nil, nil)
	must(t, err)
	must(t, bn.AddCPD(cpd))

	var buf bytes.Buffer
	must(t, WritePomdpX(&buf, bn))

	bn2, err := ReadPomdpX(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v", err)
	}

	assertBNEqual(t, bn2, bn, "PomdpX single node")
}

func TestXBN_RoundTrip_SingleNode(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("Solo"))
	must(t, bn.SetStates("Solo", []string{"a", "b", "c"}))

	cpd, err := factors.NewTabularCPD("Solo", 3,
		[][]float64{{0.1}, {0.3}, {0.6}}, nil, nil)
	must(t, err)
	must(t, bn.AddCPD(cpd))

	var buf bytes.Buffer
	must(t, WriteXBN(&buf, bn))

	bn2, err := ReadXBN(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v", err)
	}

	assertBNEqual(t, bn2, bn, "XBN single node")
}
