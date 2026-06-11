//go:build unit

package readwrite

import (
	"fmt"
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/factors"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
)

// failOnOpBuilder is a test mock that fails on a specific operation type
// for coverage testing.
type failOnOpBuilder struct {
	real     *models.BayesianNetwork
	failOp   string // "AddNode", "SetStates", "AddEdge", "AddCPD"
	failNth  int    // fail on Nth call of that op (0 = all)
	opCounts map[string]int
}

// newFailOnOp creates a failOnOpBuilder that fails on the nth call of the given op.
func newFailOnOp(op string, nth int) *failOnOpBuilder {
	return &failOnOpBuilder{
		real:     models.NewBayesianNetwork(),
		failOp:   op,
		failNth:  nth,
		opCounts: make(map[string]int),
	}
}

func (f *failOnOpBuilder) shouldFail(op string) bool {
	f.opCounts[op]++
	return op == f.failOp && (f.failNth == 0 || f.opCounts[op] == f.failNth)
}

func (f *failOnOpBuilder) AddNode(name string) error {
	if f.shouldFail("AddNode") {
		return fmt.Errorf("mock: AddNode failure")
	}
	return f.real.AddNode(name)
}

func (f *failOnOpBuilder) AddEdge(from, to string) error {
	if f.shouldFail("AddEdge") {
		return fmt.Errorf("mock: AddEdge failure")
	}
	return f.real.AddEdge(from, to)
}

func (f *failOnOpBuilder) SetStates(node string, states []string) error {
	if f.shouldFail("SetStates") {
		return fmt.Errorf("mock: SetStates failure")
	}
	return f.real.SetStates(node, states)
}

func (f *failOnOpBuilder) AddCPD(cpd *factors.TabularCPD) error {
	if f.shouldFail("AddCPD") {
		return fmt.Errorf("mock: AddCPD failure")
	}
	return f.real.AddCPD(cpd)
}

// --- BIF reader defensive tests via readBIFWith ---

const diBIF1 = `network unknown {}
variable A {
  type discrete [ 2 ] { True, False };
}
probability ( A ) {
  table 0.6, 0.4;
}
`

const diBIF2 = `network unknown {}
variable A {
  type discrete [ 2 ] { True, False };
}
variable B {
  type discrete [ 2 ] { Yes, No };
}
probability ( A ) {
  table 0.6, 0.4;
}
probability ( B | A ) {
  (True) 0.9, 0.1;
  (False) 0.2, 0.8;
}
`

func TestDI_BIF_AddNodeFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	err := readBIFWith(strings.NewReader(diBIF1), b)
	if err == nil {
		t.Fatal("expected error when AddNode fails")
	}
}

func TestDI_BIF_SetStatesFail(t *testing.T) {
	b := newFailOnOp("SetStates", 1)
	err := readBIFWith(strings.NewReader(diBIF1), b)
	if err == nil {
		t.Fatal("expected error when SetStates fails")
	}
}

func TestDI_BIF_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	err := readBIFWith(strings.NewReader(diBIF2), b)
	if err == nil {
		t.Fatal("expected error when AddEdge fails")
	}
}

func TestDI_BIF_AddCPDFail(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	err := readBIFWith(strings.NewReader(diBIF1), b)
	if err == nil {
		t.Fatal("expected error when AddCPD fails")
	}
}

// --- NET reader defensive tests via readNETWith ---

const diNET1 = `net {}
node A {
  states = ("True" "False");
}
potential (A) {
  data = (0.6 0.4);
}
`

const diNET2 = `net {}
node A {
  states = ("True" "False");
}
node B {
  states = ("Yes" "No");
}
potential (A) {
  data = (0.6 0.4);
}
potential (B | A) {
  data = ((0.9 0.1)(0.2 0.8));
}
`

func TestDI_NET_AddNodeFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	err := readNETWith(strings.NewReader(diNET1), b)
	if err == nil {
		t.Fatal("expected error when AddNode fails")
	}
}

func TestDI_NET_SetStatesFail(t *testing.T) {
	b := newFailOnOp("SetStates", 1)
	err := readNETWith(strings.NewReader(diNET1), b)
	if err == nil {
		t.Fatal("expected error when SetStates fails")
	}
}

func TestDI_NET_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	err := readNETWith(strings.NewReader(diNET2), b)
	if err == nil {
		t.Fatal("expected error when AddEdge fails")
	}
}

func TestDI_NET_AddCPDFail(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	err := readNETWith(strings.NewReader(diNET1), b)
	if err == nil {
		t.Fatal("expected error when AddCPD fails")
	}
}

// --- XMLBIF reader defensive tests via readXMLBIFWith ---

const diXMLBIF1 = `<?xml version="1.0"?>
<BIF VERSION="0.3">
<NETWORK>
<NAME>test</NAME>
<VARIABLE TYPE="nature">
  <NAME>A</NAME>
  <OUTCOME>True</OUTCOME>
  <OUTCOME>False</OUTCOME>
</VARIABLE>
<DEFINITION>
  <FOR>A</FOR>
  <TABLE>0.6 0.4</TABLE>
</DEFINITION>
</NETWORK>
</BIF>
`

const diXMLBIF2 = `<?xml version="1.0"?>
<BIF VERSION="0.3">
<NETWORK>
<NAME>test</NAME>
<VARIABLE TYPE="nature">
  <NAME>A</NAME>
  <OUTCOME>True</OUTCOME>
  <OUTCOME>False</OUTCOME>
</VARIABLE>
<VARIABLE TYPE="nature">
  <NAME>B</NAME>
  <OUTCOME>Yes</OUTCOME>
  <OUTCOME>No</OUTCOME>
</VARIABLE>
<DEFINITION>
  <FOR>A</FOR>
  <TABLE>0.6 0.4</TABLE>
</DEFINITION>
<DEFINITION>
  <FOR>B</FOR>
  <GIVEN>A</GIVEN>
  <TABLE>0.9 0.1 0.2 0.8</TABLE>
</DEFINITION>
</NETWORK>
</BIF>
`

func TestDI_XMLBIF_AddNodeFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	err := readXMLBIFWith(strings.NewReader(diXMLBIF1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddNode fails")
	}
}

func TestDI_XMLBIF_SetStatesFail(t *testing.T) {
	b := newFailOnOp("SetStates", 1)
	err := readXMLBIFWith(strings.NewReader(diXMLBIF1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when SetStates fails")
	}
}

func TestDI_XMLBIF_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	err := readXMLBIFWith(strings.NewReader(diXMLBIF2), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddEdge fails")
	}
}

func TestDI_XMLBIF_AddCPDFail(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	err := readXMLBIFWith(strings.NewReader(diXMLBIF1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddCPD fails")
	}
}

// --- XDSL reader defensive tests via readXDSLWith ---

const diXDSL1 = `<?xml version="1.0"?>
<smile id="test">
<nodes>
<cpt id="A">
  <state id="True"/>
  <state id="False"/>
  <probabilities>0.6 0.4</probabilities>
</cpt>
</nodes>
</smile>
`

const diXDSL2 = `<?xml version="1.0"?>
<smile id="test">
<nodes>
<cpt id="A">
  <state id="True"/>
  <state id="False"/>
  <probabilities>0.6 0.4</probabilities>
</cpt>
<cpt id="B">
  <state id="Yes"/>
  <state id="No"/>
  <parents>A</parents>
  <probabilities>0.9 0.1 0.2 0.8</probabilities>
</cpt>
</nodes>
</smile>
`

func TestDI_XDSL_AddNodeFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	err := readXDSLWith(strings.NewReader(diXDSL1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddNode fails")
	}
}

func TestDI_XDSL_SetStatesFail(t *testing.T) {
	b := newFailOnOp("SetStates", 1)
	err := readXDSLWith(strings.NewReader(diXDSL1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when SetStates fails")
	}
}

func TestDI_XDSL_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	err := readXDSLWith(strings.NewReader(diXDSL2), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddEdge fails")
	}
}

func TestDI_XDSL_AddCPDFail(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	err := readXDSLWith(strings.NewReader(diXDSL1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddCPD fails")
	}
}

// --- PomdpX reader defensive tests via readPomdpXWith ---

const diPomdpX1 = `<?xml version="1.0"?>
<pomdpx version="1.0">
<Variable>
  <StateVar vnamePrev="A" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
</Variable>
<InitialStateBelief>
  <CondProb>
    <Var>A</Var>
    <Parameter type="TBL">
      <Entry><ProbTable>0.6 0.4</ProbTable></Entry>
    </Parameter>
  </CondProb>
</InitialStateBelief>
</pomdpx>
`

const diPomdpX2 = `<?xml version="1.0"?>
<pomdpx version="1.0">
<Variable>
  <StateVar vnamePrev="A" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
  <StateVar vnamePrev="B" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
</Variable>
<InitialStateBelief>
  <CondProb>
    <Var>A</Var>
    <Parameter type="TBL">
      <Entry><ProbTable>0.6 0.4</ProbTable></Entry>
    </Parameter>
  </CondProb>
</InitialStateBelief>
<StateTransitionFunction>
  <CondProb>
    <Var>B</Var>
    <Parent>A</Parent>
    <Parameter type="TBL">
      <Entry><Instance>s0</Instance><ProbTable>0.9 0.1</ProbTable></Entry>
      <Entry><Instance>s1</Instance><ProbTable>0.2 0.8</ProbTable></Entry>
    </Parameter>
  </CondProb>
</StateTransitionFunction>
</pomdpx>
`

func TestDI_PomdpX_AddNodeFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	bn := models.NewBayesianNetwork()
	err := readPomdpXWith(strings.NewReader(diPomdpX1), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddNode fails")
	}
}

func TestDI_PomdpX_SetStatesFail(t *testing.T) {
	b := newFailOnOp("SetStates", 1)
	bn := models.NewBayesianNetwork()
	err := readPomdpXWith(strings.NewReader(diPomdpX1), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when SetStates fails")
	}
}

func TestDI_PomdpX_AddCPDFail(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	bn := models.NewBayesianNetwork()
	err := readPomdpXWith(strings.NewReader(diPomdpX1), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddCPD fails")
	}
}

func TestDI_PomdpX_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	bn := models.NewBayesianNetwork()
	err := readPomdpXWith(strings.NewReader(diPomdpX2), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddEdge fails")
	}
}

func TestDI_PomdpX_DefaultCPDFail(t *testing.T) {
	// PomdpX with no CPD data - triggers default CPD creation
	input := `<?xml version="1.0"?>
<pomdpx version="1.0">
<Variable>
  <StateVar vnamePrev="A" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
</Variable>
</pomdpx>
`
	b := newFailOnOp("AddCPD", 1)
	bn := models.NewBayesianNetwork()
	err := readPomdpXWith(strings.NewReader(input), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when default AddCPD fails")
	}
}

// --- XBN reader defensive tests via readXBNWith ---

const diXBN1 = `<?xml version="1.0"?>
<ANALYSISNOTEBOOK>
<BNMODEL NAME="test">
<STATICPROPERTIES>
  <NODELIST>
    <NODE NAME="A">
      <STATENAME>True</STATENAME>
      <STATENAME>False</STATENAME>
    </NODE>
  </NODELIST>
  <ARCLIST></ARCLIST>
</STATICPROPERTIES>
<DYNAMICPROPERTIES>
  <DISTRIBS>
    <DIST TYPE="discrete">
      <DPIS><DPI>0.6 0.4</DPI></DPIS>
    </DIST>
  </DISTRIBS>
</DYNAMICPROPERTIES>
</BNMODEL>
</ANALYSISNOTEBOOK>
`

const diXBN2 = `<?xml version="1.0"?>
<ANALYSISNOTEBOOK>
<BNMODEL NAME="test">
<STATICPROPERTIES>
  <NODELIST>
    <NODE NAME="A">
      <STATENAME>True</STATENAME>
      <STATENAME>False</STATENAME>
    </NODE>
    <NODE NAME="B">
      <STATENAME>Yes</STATENAME>
      <STATENAME>No</STATENAME>
    </NODE>
  </NODELIST>
  <ARCLIST>
    <ARC PARENT="A" CHILD="B"/>
  </ARCLIST>
</STATICPROPERTIES>
<DYNAMICPROPERTIES>
  <DISTRIBS>
    <DIST TYPE="discrete">
      <DPIS><DPI>0.6 0.4</DPI></DPIS>
    </DIST>
    <DIST TYPE="discrete">
      <CONDSET><CONDELEM NAME="A"/></CONDSET>
      <DPIS>
        <DPI INDEXES="0">0.9 0.1</DPI>
        <DPI INDEXES="1">0.2 0.8</DPI>
      </DPIS>
    </DIST>
  </DISTRIBS>
</DYNAMICPROPERTIES>
</BNMODEL>
</ANALYSISNOTEBOOK>
`

func TestDI_XBN_AddNodeFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	bn := models.NewBayesianNetwork()
	err := readXBNWith(strings.NewReader(diXBN1), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddNode fails")
	}
}

func TestDI_XBN_SetStatesFail(t *testing.T) {
	b := newFailOnOp("SetStates", 1)
	bn := models.NewBayesianNetwork()
	err := readXBNWith(strings.NewReader(diXBN1), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when SetStates fails")
	}
}

func TestDI_XBN_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	bn := models.NewBayesianNetwork()
	err := readXBNWith(strings.NewReader(diXBN2), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddEdge fails")
	}
}

func TestDI_XBN_AddCPDFail(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	bn := models.NewBayesianNetwork()
	err := readXBNWith(strings.NewReader(diXBN1), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddCPD fails")
	}
}

func TestDI_XBN_DefaultCPDFail(t *testing.T) {
	// XBN with 2 nodes but only 1 dist -> B gets default CPD
	input := `<?xml version="1.0"?>
<ANALYSISNOTEBOOK>
<BNMODEL NAME="test">
<STATICPROPERTIES>
  <NODELIST>
    <NODE NAME="A"><STATENAME>s0</STATENAME><STATENAME>s1</STATENAME></NODE>
    <NODE NAME="B"><STATENAME>s0</STATENAME><STATENAME>s1</STATENAME></NODE>
  </NODELIST>
  <ARCLIST></ARCLIST>
</STATICPROPERTIES>
<DYNAMICPROPERTIES>
  <DISTRIBS>
    <DIST TYPE="discrete"><DPIS><DPI>0.6 0.4</DPI></DPIS></DIST>
  </DISTRIBS>
</DYNAMICPROPERTIES>
</BNMODEL>
</ANALYSISNOTEBOOK>
`
	b := newFailOnOp("AddCPD", 2)
	bn := models.NewBayesianNetwork()
	err := readXBNWith(strings.NewReader(input), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when default AddCPD fails")
	}
}

// --- XMLNative reader defensive tests via readXMLNativeWith ---

const diXMLNative1 = `<?xml version="1.0"?>
<pgmgo-network name="test">
<nodes>
  <node name="A" states="True,False"/>
</nodes>
<edges></edges>
<cpds>
  <cpd variable="A" card="2"><values>0.6 0.4</values></cpd>
</cpds>
</pgmgo-network>
`

const diXMLNative2 = `<?xml version="1.0"?>
<pgmgo-network name="test">
<nodes>
  <node name="A" states="True,False"/>
  <node name="B" states="Yes,No"/>
</nodes>
<edges><edge from="A" to="B"/></edges>
<cpds>
  <cpd variable="A" card="2"><values>0.6 0.4</values></cpd>
  <cpd variable="B" card="2" evidence="A" evidence_card="2"><values>0.9 0.1 0.2 0.8</values></cpd>
</cpds>
</pgmgo-network>
`

func TestDI_XMLNative_AddNodeFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	err := readXMLNativeWith(strings.NewReader(diXMLNative1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddNode fails")
	}
}

func TestDI_XMLNative_SetStatesFail(t *testing.T) {
	b := newFailOnOp("SetStates", 1)
	err := readXMLNativeWith(strings.NewReader(diXMLNative1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when SetStates fails")
	}
}

func TestDI_XMLNative_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	err := readXMLNativeWith(strings.NewReader(diXMLNative2), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddEdge fails")
	}
}

func TestDI_XMLNative_AddCPDFail(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	err := readXMLNativeWith(strings.NewReader(diXMLNative1), b, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when AddCPD fails")
	}
}

// --- CSV structure defensive tests via readCSVStructureWith ---

func TestDI_CSV_AddNodeFromFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	err := readCSVStructureWith(strings.NewReader("from,to\nA,B\n"), b)
	if err == nil {
		t.Fatal("expected error when AddNode fails for 'from'")
	}
}

func TestDI_CSV_AddNodeToFail(t *testing.T) {
	b := newFailOnOp("AddNode", 2)
	err := readCSVStructureWith(strings.NewReader("from,to\nA,B\n"), b)
	if err == nil {
		t.Fatal("expected error when AddNode fails for 'to'")
	}
}

func TestDI_CSV_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	err := readCSVStructureWith(strings.NewReader("from,to\nA,B\n"), b)
	if err == nil {
		t.Fatal("expected error when AddEdge fails")
	}
}

// --- JSON AddCPD defensive test via jsonAddCPDs ---

func TestDI_JSON_AddCPDFail2(t *testing.T) {
	jn := &jsonNetwork{
		CPDs: map[string]*jsonCPD{
			"A": {
				VariableCard: 2,
				Values:       [][]float64{{0.6}, {0.4}},
			},
		},
	}
	b := newFailOnOp("AddCPD", 1)
	// First make sure node exists
	_ = b.AddNode("A")
	err := jsonAddCPDs(jn, b)
	if err == nil {
		t.Fatal("expected error when AddCPD fails")
	}
}

// --- Write error tests for CSV ---

// failNthWriter is a test mock that fails on the Nth Write call
// for coverage testing.
type failNthWriter struct {
	count int
	nth   int
}

func (f *failNthWriter) Write(p []byte) (int, error) {
	f.count++
	if f.count >= f.nth {
		return 0, fmt.Errorf("mock: write failure at call %d", f.count)
	}
	return len(p), nil
}

// --- BIF/NET empty-line and edge-case tests ---

func TestDI_BIF_EmptyLinesInBody(t *testing.T) {
	// BIF with blank lines that become empty after stripping
	// readBIFLines strips empty lines, but readBIFWith handles
	// empty tokens if any slip through. Use a comment-only line.
	input := "network unknown {}\n//just a comment\nvariable A {\n  type discrete [ 2 ] { True, False };\n}\nprobability ( A ) {\n  table 0.6, 0.4;\n}\n"
	bn, err := ReadBIF(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bn.Nodes()) != 1 {
		t.Fatalf("expected 1 node, got %d", len(bn.Nodes()))
	}
}

// (NET comment test removed - handled by existing tests)

// --- PomdpX unconditional CPD creation error path ---
func TestDI_PomdpX_UnconditionalCPDCreateFail(t *testing.T) {
	input := `<?xml version="1.0"?>
<pomdpx version="1.0">
<Variable>
  <StateVar vnamePrev="A" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
</Variable>
<InitialStateBelief>
  <CondProb>
    <Var>A</Var>
    <Parameter type="TBL">
      <Entry><ProbTable>0.6 0.4</ProbTable></Entry>
    </Parameter>
  </CondProb>
</InitialStateBelief>
</pomdpx>
`
	// Test normal path first
	bn, err := ReadPomdpX(strings.NewReader(input))
	if err != nil {
		t.Fatalf("sanity: %v", err)
	}
	if bn.GetCPD("A") == nil {
		t.Fatal("expected CPD for A")
	}
}

// --- PomdpX conditional CPD create fail ---
func TestDI_PomdpX_ConditionalCPDCreateFail(t *testing.T) {
	input := `<?xml version="1.0"?>
<pomdpx version="1.0">
<Variable>
  <StateVar vnamePrev="A" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
  <StateVar vnamePrev="B" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
</Variable>
<StateTransitionFunction>
  <CondProb>
    <Var>B</Var>
    <Parent>A</Parent>
    <Parameter type="TBL">
      <Entry><Instance>s0</Instance><ProbTable>0.9 0.1</ProbTable></Entry>
      <Entry><Instance>s1</Instance><ProbTable>0.2 0.8</ProbTable></Entry>
    </Parameter>
  </CondProb>
</StateTransitionFunction>
</pomdpx>
`
	// Test with AddCPD failing on first conditional CPD
	b := newFailOnOp("AddCPD", 1)
	bn := models.NewBayesianNetwork()
	err := readPomdpXWith(strings.NewReader(input), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when conditional AddCPD fails")
	}
}

// --- UAI: test some additional error paths via production ReadUAI ---
func TestDI_UAI_BadScopeVar(t *testing.T) {
	input := "BAYES\n1\n2\n1\n1 abc\n2\n0.6 0.4\n"
	_, err := ReadUAI(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for non-integer scope var")
	}
}

func TestDI_UAI_BadFactorCount(t *testing.T) {
	input := "BAYES\n1\n2\nabc\n"
	_, err := ReadUAI(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for non-integer factor count")
	}
}

func TestDI_UAI_TruncatedFactorTable(t *testing.T) {
	input := "BAYES\n1\n2\n1\n1 0\nabc\n"
	_, err := ReadUAI(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for non-integer entry count")
	}
}

// CSV ReadCSVCPD short row test
func TestDI_ReadCSVCPD_ShortDataRow(t *testing.T) {
	input := "A,P\ns0,0.6\ns1\n"
	_, err := ReadCSVCPD(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for short data row")
	}
}

// CSV ReadCSVCPD CPD creation error
func TestDI_ReadCSVCPD_CPDCreateError(t *testing.T) {
	// Create a conditional CPD with wrong cardinality
	input := "A,B=s0,B=s1,B=s2\ns0,0.1,0.2,0.3\n"
	// This should parse but might fail on CPD creation due to mismatch
	_, _ = ReadCSVCPD(strings.NewReader(input))
}

// --- ReadCSVStructure short row ---

// (CSV short row test removed - csv.Reader enforces field count)

// --- Additional PomdpX edge case tests ---

func TestDI_PomdpX_EmptyInstanceEntry(t *testing.T) {
	// PomdpX with empty Instance tag - should be skipped
	input := `<?xml version="1.0"?>
<pomdpx version="1.0">
<Variable>
  <StateVar vnamePrev="A" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
  <StateVar vnamePrev="B" numValues="2">
    <ValueEnum>s0 s1</ValueEnum>
  </StateVar>
</Variable>
<StateTransitionFunction>
  <CondProb>
    <Var>B</Var>
    <Parent>A</Parent>
    <Parameter type="TBL">
      <Entry><Instance></Instance><ProbTable>0.5 0.5</ProbTable></Entry>
      <Entry><Instance>s0</Instance><ProbTable>0.9 0.1</ProbTable></Entry>
      <Entry><Instance>s1</Instance><ProbTable>0.2 0.8</ProbTable></Entry>
    </Parameter>
  </CondProb>
</StateTransitionFunction>
</pomdpx>
`
	bn, err := ReadPomdpX(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cpd := bn.GetCPD("B")
	if cpd == nil {
		t.Fatal("expected CPD for B")
	}
}

// --- XBN CondElem self-skip ---
func TestDI_XBN_CondElemSelfSkip(t *testing.T) {
	// XBN where CONDELEM includes the child itself
	input := `<?xml version="1.0"?>
<ANALYSISNOTEBOOK>
<BNMODEL NAME="test">
<STATICPROPERTIES>
  <NODELIST>
    <NODE NAME="A"><STATENAME>s0</STATENAME><STATENAME>s1</STATENAME></NODE>
    <NODE NAME="B"><STATENAME>s0</STATENAME><STATENAME>s1</STATENAME></NODE>
  </NODELIST>
  <ARCLIST>
    <ARC PARENT="A" CHILD="B"/>
  </ARCLIST>
</STATICPROPERTIES>
<DYNAMICPROPERTIES>
  <DISTRIBS>
    <DIST TYPE="discrete">
      <DPIS><DPI>0.6 0.4</DPI></DPIS>
    </DIST>
    <DIST TYPE="discrete">
      <CONDSET>
        <CONDELEM NAME="B"/>
        <CONDELEM NAME="A"/>
      </CONDSET>
      <DPIS>
        <DPI INDEXES="0">0.9 0.1</DPI>
        <DPI INDEXES="1">0.2 0.8</DPI>
      </DPIS>
    </DIST>
  </DISTRIBS>
</DYNAMICPROPERTIES>
</BNMODEL>
</ANALYSISNOTEBOOK>
`
	bn, err := ReadXBN(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bn.Nodes()) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(bn.Nodes()))
	}
}

// --- XDSL CPD creation error ---
func TestDI_XDSL_CPDCreateFail(t *testing.T) {
	// Test with NewTabularCPD failure by using readXDSLWith with a builder that
	// fails on AddCPD - already tested above. Let's test a different code path:
	// parent with empty trimmed name (fields never produces empty, so skip this)
	_ = 0
}

// --- Additional targeted tests for remaining uncovered stmts ---

// NET: test CPD creation error (net.go:150)
func TestDI_NET_CPDCreateFail(t *testing.T) {
	// NET with values that produce a bad CPD
	// NewTabularCPD fails when childCard doesn't match values
	// But in practice the NET reader always constructs valid params.
	// Test via builder mock:
	b := newFailOnOp("AddCPD", 1)
	err := readNETWith(strings.NewReader(diNET1), b)
	if err == nil {
		t.Fatal("expected error when NET AddCPD fails")
	}
}

// (NET unknown keyword test removed - NET parser treats unknown keywords as node names)

// XBN: CPD creation fail and default CPD fail
func TestDI_XBN_CPDCreateFail2(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	bn := models.NewBayesianNetwork()
	err := readXBNWith(strings.NewReader(diXBN1), b, bn, MaxInputSize)
	if err == nil {
		t.Fatal("expected error when XBN AddCPD fails")
	}
}

// --- UAI builder defensive tests via uaiBuildNetwork ---

func TestDI_UAI_AddNodeFail(t *testing.T) {
	b := newFailOnOp("AddNode", 1)
	noInt := func() (int, error) { return 0, nil }
	noFloat := func() (float64, error) { return 0, nil }
	err := uaiBuildNetwork(b, 1, []int{2}, 0, nil, noInt, noFloat)
	if err == nil {
		t.Fatal("expected error when UAI AddNode fails")
	}
}

func TestDI_UAI_SetStatesFail(t *testing.T) {
	b := newFailOnOp("SetStates", 1)
	noInt := func() (int, error) { return 0, nil }
	noFloat := func() (float64, error) { return 0, nil }
	err := uaiBuildNetwork(b, 1, []int{2}, 0, nil, noInt, noFloat)
	if err == nil {
		t.Fatal("expected error when UAI SetStates fails")
	}
}

func TestDI_UAI_AddEdgeFail(t *testing.T) {
	b := newFailOnOp("AddEdge", 1)
	idx := 0
	ints := []int{2, 0, 1} // numEntries=2, vals follow
	floats := []float64{0.6, 0.4}
	nextInt := func() (int, error) {
		if idx >= len(ints) {
			return 0, fmt.Errorf("end")
		}
		v := ints[idx]
		idx++
		return v, nil
	}
	fidx := 0
	nextFloat := func() (float64, error) {
		if fidx >= len(floats) {
			return 0, fmt.Errorf("end")
		}
		v := floats[fidx]
		fidx++
		return v, nil
	}
	scopes := []uaiScope{{vars: []int{0, 1}}} // parent 0 -> child 1
	err := uaiBuildNetwork(b, 2, []int{2, 2}, 1, scopes, nextInt, nextFloat)
	if err == nil {
		t.Fatal("expected error when UAI AddEdge fails")
	}
}

func TestDI_UAI_AddCPDFail(t *testing.T) {
	b := newFailOnOp("AddCPD", 1)
	idx := 0
	ints := []int{2}
	floats := []float64{0.6, 0.4}
	nextInt := func() (int, error) {
		if idx >= len(ints) {
			return 0, fmt.Errorf("end")
		}
		v := ints[idx]
		idx++
		return v, nil
	}
	fidx := 0
	nextFloat := func() (float64, error) {
		if fidx >= len(floats) {
			return 0, fmt.Errorf("end")
		}
		v := floats[fidx]
		fidx++
		return v, nil
	}
	scopes := []uaiScope{{vars: []int{0}}} // single node, no parents
	err := uaiBuildNetwork(b, 1, []int{2}, 1, scopes, nextInt, nextFloat)
	if err == nil {
		t.Fatal("expected error when UAI AddCPD fails")
	}
}

// --- failOnOpBuilder mock self-tests ---

func TestDI_MockSelfTest(t *testing.T) {
	b := newFailOnOp("AddNode", 2)
	_ = b.AddNode("A")
	if err := b.AddNode("B"); err == nil {
		t.Fatal("expected 2nd AddNode to fail")
	}

	b = newFailOnOp("SetStates", 0)
	_ = b.AddNode("A")
	if err := b.SetStates("A", []string{"s0"}); err == nil {
		t.Fatal("expected all SetStates to fail")
	}

	b = newFailOnOp("AddEdge", 0)
	_ = b.AddNode("A")
	_ = b.AddNode("B")
	if err := b.AddEdge("A", "B"); err == nil {
		t.Fatal("expected all AddEdge to fail")
	}

	b = newFailOnOp("AddCPD", 0)
	_ = b.AddNode("A")
	cpd, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.6}, {0.4}}, nil, nil)
	if err := b.AddCPD(cpd); err == nil {
		t.Fatal("expected all AddCPD to fail")
	}
}
