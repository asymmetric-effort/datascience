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
// BIF Reader error paths
// ---------------------------------------------------------------------------

func TestReadBIF_MalformedVariableName(t *testing.T) {
	bif := "network test {\n}\nvariable {\n  type discrete [ 2 ] { A, B };\n}\n"
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for empty variable name")
	}
}

func TestReadBIF_NoTypeDecl(t *testing.T) {
	bif := "network test {\n}\nvariable X {\n  something else;\n}\n"
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for missing type declaration")
	}
}

func TestReadBIF_MalformedTypeDecl(t *testing.T) {
	bif := "network test {\n}\nvariable X {\n  type discrete [ 2 ] no_braces;\n}\n"
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for malformed type declaration")
	}
}

func TestReadBIF_EmptyStates(t *testing.T) {
	bif := "network test {\n}\nvariable X {\n  type discrete [ 0 ] { };\n}\n"
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for empty states")
	}
}

func TestReadBIF_UnknownVarInProb(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
probability ( Unknown ) {
  table 0.5, 0.5;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for unknown variable in probability")
	}
}

func TestReadBIF_UnknownParentInProb(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
probability ( X | Unknown ) {
  (a) 0.5, 0.5;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for unknown parent in probability")
	}
}

func TestReadBIF_MalformedProbHeader(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
probability no_parens {
  table 0.5, 0.5;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for malformed prob header")
	}
}

func TestReadBIF_EmptyVarInProbHeader(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
probability ( ) {
  table 0.5, 0.5;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for empty var in prob header")
	}
}

func TestReadBIF_TableWrongCount(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
probability ( X ) {
  table 0.3, 0.4, 0.3;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for wrong table value count")
	}
}

func TestReadBIF_CondTableWrongCount(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
variable Y {
  type discrete [ 2 ] { Y0, Y1 };
}
probability ( X ) {
  table 0.5, 0.5;
}
probability ( Y | X ) {
  table 0.1, 0.2, 0.3;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for wrong conditional table value count")
	}
}

func TestReadBIF_InvalidFloat(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
probability ( X ) {
  table abc, 0.5;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for invalid float")
	}
}

func TestReadBIF_MalformedConditionalLine(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
variable Y {
  type discrete [ 2 ] { Y0, Y1 };
}
probability ( X ) {
  table 0.5, 0.5;
}
probability ( Y | X ) {
  (A 0.3, 0.7;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for malformed conditional line")
	}
}

func TestReadBIF_WrongParentStateCount(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
variable Z {
  type discrete [ 2 ] { Z0, Z1 };
}
variable Y {
  type discrete [ 2 ] { Y0, Y1 };
}
probability ( X ) {
  table 0.5, 0.5;
}
probability ( Z ) {
  table 0.5, 0.5;
}
probability ( Y | X, Z ) {
  (A) 0.3, 0.7;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for wrong parent state count")
	}
}

func TestReadBIF_UnknownParentState(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
variable Y {
  type discrete [ 2 ] { Y0, Y1 };
}
probability ( X ) {
  table 0.5, 0.5;
}
probability ( Y | X ) {
  (Unknown) 0.3, 0.7;
  (B) 0.6, 0.4;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for unknown parent state")
	}
}

func TestReadBIF_CondValuesWrongCount(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
variable Y {
  type discrete [ 2 ] { Y0, Y1 };
}
probability ( X ) {
  table 0.5, 0.5;
}
probability ( Y | X ) {
  (A) 0.3, 0.4, 0.3;
  (B) 0.6, 0.4;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for wrong conditional value count")
	}
}

func TestReadBIF_CondInvalidFloat(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
variable Y {
  type discrete [ 2 ] { Y0, Y1 };
}
probability ( X ) {
  table 0.5, 0.5;
}
probability ( Y | X ) {
  (A) abc, 0.7;
  (B) 0.6, 0.4;
}
`
	_, err := ReadBIF(strings.NewReader(bif))
	if err == nil {
		t.Error("expected error for invalid float in conditional")
	}
}

// ---------------------------------------------------------------------------
// BIF Writer error paths
// ---------------------------------------------------------------------------

func TestWriteBIF_NoStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	// No states set.
	var buf bytes.Buffer
	err := WriteBIF(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no states")
	}
}

func TestWriteBIF_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A", "B"}))
	var buf bytes.Buffer
	err := WriteBIF(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no CPD")
	}
}

func TestWriteBIF_StateIndexOutOfRange(t *testing.T) {
	// Exercise the fallback path in bifDecomposeParentConfig where indices[i] >= len(states)
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A"}))
	// pc=2 with evidenceCard=[2] would produce indices[0]=0 for pc=0, indices[0]=1 for pc=1.
	// With only 1 state ("A"), index 1 >= len(states)=1 triggers fallback.
	names := bifDecomposeParentConfig(1, []string{"X"}, []int{2}, bn)
	if len(names) != 1 || names[0] != "state1" {
		t.Errorf("expected [state1], got %v", names)
	}
}

// ---------------------------------------------------------------------------
// NET error paths
// ---------------------------------------------------------------------------

func TestReadNET_MalformedNode(t *testing.T) {
	netData := "net\n{\n}\nnode {\n  states = (\"a\");\n}\n"
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for malformed node")
	}
}

func TestReadNET_NoStatesDecl(t *testing.T) {
	netData := "net\n{\n}\nnode X\n{\n  foo = bar;\n}\n"
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for missing states declaration")
	}
}

func TestReadNET_MalformedStates(t *testing.T) {
	netData := "net\n{\n}\nnode X\n{\n  states = no_parens;\n}\n"
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for malformed states")
	}
}

func TestReadNET_EmptyStates(t *testing.T) {
	netData := "net\n{\n}\nnode X\n{\n  states = ();\n}\n"
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for empty states")
	}
}

func TestReadNET_UnknownVarInPotential(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential (Unknown)
{
  data = (0.5 0.5);
}
`
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for unknown variable in potential")
	}
}

func TestReadNET_UnknownParent(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential (X | Unknown)
{
  data = ((0.5 0.5));
}
`
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for unknown parent")
	}
}

func TestReadNET_MalformedPotentialHeader(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential no_parens
{
  data = (0.5 0.5);
}
`
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for malformed potential header")
	}
}

func TestReadNET_EmptyVarInPotentialHeader(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential ( )
{
  data = (0.5 0.5);
}
`
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for empty var in potential header")
	}
}

func TestReadNET_NoDataDecl(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential (X)
{
  something = else;
}
`
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for missing data declaration")
	}
}

func TestReadNET_InvalidDataFloat(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential (X)
{
  data = (abc 0.5);
}
`
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for invalid float in data")
	}
}

func TestReadNET_WrongDataCount(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential (X)
{
  data = (0.3 0.4 0.3);
}
`
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for wrong data count")
	}
}

func TestReadNET_Comments(t *testing.T) {
	netData := `net
{
}
% This is a comment
node X
{
  states = ("a" "b"); % inline comment
}
potential (X)
{
  data = (0.3 0.7);
}
`
	bn, err := ReadNET(strings.NewReader(netData))
	if err != nil {
		t.Fatalf("ReadNET with comments failed: %v", err)
	}
	if len(bn.Nodes()) != 1 {
		t.Errorf("expected 1 node, got %d", len(bn.Nodes()))
	}
}

func TestWriteNET_NoStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	var buf bytes.Buffer
	err := WriteNET(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no states")
	}
}

func TestWriteNET_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A", "B"}))
	var buf bytes.Buffer
	err := WriteNET(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no CPD")
	}
}

// ---------------------------------------------------------------------------
// UAI error paths
// ---------------------------------------------------------------------------

func TestReadUAI_Empty(t *testing.T) {
	_, err := ReadUAI(strings.NewReader(""))
	if err == nil {
		t.Error("expected error for empty UAI file")
	}
}

func TestReadUAI_UnexpectedEnd(t *testing.T) {
	_, err := ReadUAI(strings.NewReader("BAYES\n"))
	if err == nil {
		t.Error("expected error for unexpected end of UAI file")
	}
}

func TestReadUAI_BadCardinality(t *testing.T) {
	_, err := ReadUAI(strings.NewReader("BAYES\n2\nabc 2\n"))
	if err == nil {
		t.Error("expected error for non-integer cardinality")
	}
}

func TestReadUAI_BadFloat(t *testing.T) {
	uai := "BAYES\n1\n2\n1\n1 0\n2\nabc 0.5\n"
	_, err := ReadUAI(strings.NewReader(uai))
	if err == nil {
		t.Error("expected error for non-float value")
	}
}

func TestReadUAI_WrongTableSize(t *testing.T) {
	uai := "BAYES\n1\n2\n1\n1 0\n3\n0.3 0.4 0.3\n"
	_, err := ReadUAI(strings.NewReader(uai))
	if err == nil {
		t.Error("expected error for wrong table size")
	}
}

func TestWriteUAI_NoStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	var buf bytes.Buffer
	err := WriteUAI(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no states")
	}
}

func TestWriteUAI_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A", "B"}))
	var buf bytes.Buffer
	err := WriteUAI(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no CPD")
	}
}

// ---------------------------------------------------------------------------
// XMLBIF error paths
// ---------------------------------------------------------------------------

func TestReadXMLBIF_UnknownVariable(t *testing.T) {
	xml := `<?xml version="1.0"?>
<BIF VERSION="0.3">
<NETWORK>
<NAME>test</NAME>
<VARIABLE TYPE="nature"><NAME>X</NAME><OUTCOME>a</OUTCOME><OUTCOME>b</OUTCOME></VARIABLE>
<DEFINITION><FOR>Unknown</FOR><TABLE>0.5 0.5</TABLE></DEFINITION>
</NETWORK>
</BIF>`
	_, err := ReadXMLBIF(strings.NewReader(xml))
	if err == nil {
		t.Error("expected error for unknown variable")
	}
}

func TestReadXMLBIF_UnknownParent(t *testing.T) {
	xml := `<?xml version="1.0"?>
<BIF VERSION="0.3">
<NETWORK>
<NAME>test</NAME>
<VARIABLE TYPE="nature"><NAME>X</NAME><OUTCOME>a</OUTCOME><OUTCOME>b</OUTCOME></VARIABLE>
<DEFINITION><FOR>X</FOR><GIVEN>Unknown</GIVEN><TABLE>0.5 0.5</TABLE></DEFINITION>
</NETWORK>
</BIF>`
	_, err := ReadXMLBIF(strings.NewReader(xml))
	if err == nil {
		t.Error("expected error for unknown parent")
	}
}

func TestReadXMLBIF_BadFloat(t *testing.T) {
	xml := `<?xml version="1.0"?>
<BIF VERSION="0.3">
<NETWORK>
<NAME>test</NAME>
<VARIABLE TYPE="nature"><NAME>X</NAME><OUTCOME>a</OUTCOME><OUTCOME>b</OUTCOME></VARIABLE>
<DEFINITION><FOR>X</FOR><TABLE>abc 0.5</TABLE></DEFINITION>
</NETWORK>
</BIF>`
	_, err := ReadXMLBIF(strings.NewReader(xml))
	if err == nil {
		t.Error("expected error for invalid float")
	}
}

func TestReadXMLBIF_WrongTableSize(t *testing.T) {
	xml := `<?xml version="1.0"?>
<BIF VERSION="0.3">
<NETWORK>
<NAME>test</NAME>
<VARIABLE TYPE="nature"><NAME>X</NAME><OUTCOME>a</OUTCOME><OUTCOME>b</OUTCOME></VARIABLE>
<DEFINITION><FOR>X</FOR><TABLE>0.3 0.4 0.3</TABLE></DEFINITION>
</NETWORK>
</BIF>`
	_, err := ReadXMLBIF(strings.NewReader(xml))
	if err == nil {
		t.Error("expected error for wrong table size")
	}
}

func TestWriteXMLBIF_NoStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	var buf bytes.Buffer
	err := WriteXMLBIF(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no states")
	}
}

func TestWriteXMLBIF_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A", "B"}))
	var buf bytes.Buffer
	err := WriteXMLBIF(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no CPD")
	}
}

// ---------------------------------------------------------------------------
// XDSL error paths
// ---------------------------------------------------------------------------

func TestReadXDSL_UnknownParent(t *testing.T) {
	xml := `<?xml version="1.0"?>
<smile id="test">
<nodes>
<cpt id="X"><state id="a"/><state id="b"/><parents>Unknown</parents><probabilities>0.5 0.5</probabilities></cpt>
</nodes>
</smile>`
	_, err := ReadXDSL(strings.NewReader(xml))
	if err == nil {
		t.Error("expected error for unknown parent")
	}
}

func TestReadXDSL_WrongProbCount(t *testing.T) {
	xml := `<?xml version="1.0"?>
<smile id="test">
<nodes>
<cpt id="X"><state id="a"/><state id="b"/><probabilities>0.3 0.4 0.3</probabilities></cpt>
</nodes>
</smile>`
	_, err := ReadXDSL(strings.NewReader(xml))
	if err == nil {
		t.Error("expected error for wrong prob count")
	}
}

func TestReadXDSL_BadFloat(t *testing.T) {
	xml := `<?xml version="1.0"?>
<smile id="test">
<nodes>
<cpt id="X"><state id="a"/><state id="b"/><probabilities>abc 0.5</probabilities></cpt>
</nodes>
</smile>`
	_, err := ReadXDSL(strings.NewReader(xml))
	if err == nil {
		t.Error("expected error for invalid float")
	}
}

func TestWriteXDSL_NoStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	var buf bytes.Buffer
	err := WriteXDSL(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no states")
	}
}

func TestWriteXDSL_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A", "B"}))
	var buf bytes.Buffer
	err := WriteXDSL(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no CPD")
	}
}

// ---------------------------------------------------------------------------
// XBN error paths
// ---------------------------------------------------------------------------

func TestReadXBN_InvalidXML(t *testing.T) {
	_, err := ReadXBN(strings.NewReader("<not valid xml"))
	if err == nil {
		t.Error("expected error for invalid XML")
	}
}

func TestWriteXBN_NoStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	var buf bytes.Buffer
	err := WriteXBN(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no states")
	}
}

func TestWriteXBN_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A", "B"}))
	var buf bytes.Buffer
	err := WriteXBN(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no CPD")
	}
}

// Exercise XBN with multi-parent conditional to cover dist parsing.
func TestXBN_RoundTrip_MultiParent(t *testing.T) {
	bn := buildMultiParentBN(t)

	var buf bytes.Buffer
	if err := WriteXBN(&buf, bn); err != nil {
		t.Fatalf("WriteXBN failed: %v", err)
	}

	bn2, err := ReadXBN(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v\nOutput:\n%s", err, buf.String())
	}

	nodes1 := bn.Nodes()
	nodes2 := bn2.Nodes()
	if len(nodes1) != len(nodes2) {
		t.Fatalf("XBN multi-parent round-trip: node count mismatch: %d vs %d", len(nodes1), len(nodes2))
	}
}

// Exercise XBN ReadXBN with default binary states.
func TestReadXBN_DefaultBinaryStates(t *testing.T) {
	xml := `<?xml version="1.0"?>
<ANALYSISNOTEBOOK>
  <BNMODEL NAME="test">
    <STATICPROPERTIES>
      <NODELIST>
        <NODE NAME="X"></NODE>
      </NODELIST>
      <ARCLIST></ARCLIST>
    </STATICPROPERTIES>
    <DYNAMICPROPERTIES></DYNAMICPROPERTIES>
  </BNMODEL>
</ANALYSISNOTEBOOK>`
	bn, err := ReadXBN(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v", err)
	}
	states := bn.GetStates("X")
	if len(states) != 2 || states[0] != "s0" {
		t.Errorf("expected default binary states, got %v", states)
	}
}

// ---------------------------------------------------------------------------
// PomdpX error paths
// ---------------------------------------------------------------------------

func TestReadPomdpX_InvalidXML(t *testing.T) {
	_, err := ReadPomdpX(strings.NewReader("<not valid xml"))
	if err == nil {
		t.Error("expected error for invalid XML")
	}
}

func TestWritePomdpX_NoStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	var buf bytes.Buffer
	err := WritePomdpX(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no states")
	}
}

func TestWritePomdpX_NoCPD(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A", "B"}))
	var buf bytes.Buffer
	err := WritePomdpX(&buf, bn)
	if err == nil {
		t.Error("expected error for variable with no CPD")
	}
}

// Exercise PomdpX with generated default state names.
func TestReadPomdpX_DefaultStateNames(t *testing.T) {
	xml := `<?xml version="1.0"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="X" numValues="3"></StateVar>
  </Variable>
</pomdpx>`
	bn, err := ReadPomdpX(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v", err)
	}
	states := bn.GetStates("X")
	if len(states) != 3 || states[0] != "s0" || states[2] != "s2" {
		t.Errorf("expected default states [s0 s1 s2], got %v", states)
	}
}

// Exercise PomdpX round-trip with multi-parent to cover conditional paths.
func TestPomdpX_WriteRead_MultiParent(t *testing.T) {
	bn := buildMultiParentBN(t)
	var buf bytes.Buffer
	if err := WritePomdpX(&buf, bn); err != nil {
		t.Fatalf("WritePomdpX failed: %v", err)
	}
	_, err := ReadPomdpX(strings.NewReader(buf.String()))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v\nOutput:\n%s", err, buf.String())
	}
}

// ---------------------------------------------------------------------------
// xdslFormatProbs (unreferenced but 0% coverage)
// ---------------------------------------------------------------------------

func TestXdslFormatProbs(t *testing.T) {
	result := xdslFormatProbs([]float64{0.1, 0.2, 0.7})
	if result != "0.1 0.2 0.7" {
		t.Errorf("xdslFormatProbs = %q, want %q", result, "0.1 0.2 0.7")
	}
}

func TestXdslFormatProbs_Empty(t *testing.T) {
	result := xdslFormatProbs(nil)
	if result != "" {
		t.Errorf("xdslFormatProbs(nil) = %q, want %q", result, "")
	}
}

// ---------------------------------------------------------------------------
// xmlbifParseFloats edge cases
// ---------------------------------------------------------------------------

func TestXmlbifParseFloats_Empty(t *testing.T) {
	vals, err := xmlbifParseFloats("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 0 {
		t.Errorf("expected 0 values, got %d", len(vals))
	}
}

func TestXmlbifParseFloats_Invalid(t *testing.T) {
	_, err := xmlbifParseFloats("abc")
	if err == nil {
		t.Error("expected error for invalid float")
	}
}

// ---------------------------------------------------------------------------
// bifParseFloats edge cases
// ---------------------------------------------------------------------------

func TestBifParseFloats_Empty(t *testing.T) {
	vals, err := bifParseFloats("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vals) != 0 {
		t.Errorf("expected 0 values, got %d", len(vals))
	}
}

// ---------------------------------------------------------------------------
// BIF with table keyword for conditional
// ---------------------------------------------------------------------------

func TestReadBIF_ConditionalTable(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
variable Y {
  type discrete [ 2 ] { Y0, Y1 };
}
probability ( X ) {
  table 0.5, 0.5;
}
probability ( Y | X ) {
  table 0.1, 0.9, 0.4, 0.6;
}
`
	bn, err := ReadBIF(strings.NewReader(bif))
	if err != nil {
		t.Fatalf("ReadBIF with conditional table failed: %v", err)
	}
	yCPD := bn.GetCPD("Y")
	if yCPD == nil {
		t.Fatal("Y CPD is nil")
	}
	yData := yCPD.ToFactor().Values().Data()
	// table ordering: parent configs outer, child states inner
	// pc=0 (A): Y0=0.1, Y1=0.9; pc=1 (B): Y0=0.4, Y1=0.6
	assertFloatsClose(t, yData, []float64{0.1, 0.4, 0.9, 0.6}, "Y CPD conditional table")
}

// ---------------------------------------------------------------------------
// NET without data= keyword
// ---------------------------------------------------------------------------

func TestReadNET_NoEquals(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential (X)
{
  data (0.3 0.7);
}
`
	_, err := ReadNET(strings.NewReader(netData))
	if err == nil {
		t.Error("expected error for missing = in data declaration")
	}
}

// ---------------------------------------------------------------------------
// FormatFloat coverage
// ---------------------------------------------------------------------------

func TestFormatFloat(t *testing.T) {
	result := formatFloat(0.123456789)
	if result != "0.123456789" {
		t.Errorf("formatFloat = %q", result)
	}
}

// ---------------------------------------------------------------------------
// bifCollectBlock edge case: no braces
// ---------------------------------------------------------------------------

func TestBifCollectBlock_NoBraces(t *testing.T) {
	lines := []string{"no braces here"}
	content, end := bifCollectBlock(lines, 0)
	if content != nil {
		t.Errorf("expected nil content, got %v", content)
	}
	if end != 1 {
		t.Errorf("expected end=1, got %d", end)
	}
}

// ---------------------------------------------------------------------------
// netCollectBlock edge case: no braces
// ---------------------------------------------------------------------------

func TestNetCollectBlock_NoBraces(t *testing.T) {
	lines := []string{"no braces here"}
	content, end := netCollectBlock(lines, 0)
	if content != nil {
		t.Errorf("expected nil content, got %v", content)
	}
	if end != 1 {
		t.Errorf("expected end=1, got %d", end)
	}
}

// ---------------------------------------------------------------------------
// Full round-trip coverage for all formats with multi-parent BN
// ---------------------------------------------------------------------------

func TestAllFormats_RoundTrip_MultiParent(t *testing.T) {
	bn := buildMultiParentBN(t)

	formats := []struct {
		name  string
		write func(*bytes.Buffer, *models.BayesianNetwork) error
		read  func(string) (*models.BayesianNetwork, error)
	}{
		{
			name: "BIF",
			write: func(buf *bytes.Buffer, bn *models.BayesianNetwork) error {
				return WriteBIF(buf, bn)
			},
			read: func(s string) (*models.BayesianNetwork, error) {
				return ReadBIF(strings.NewReader(s))
			},
		},
		{
			name: "XMLBIF",
			write: func(buf *bytes.Buffer, bn *models.BayesianNetwork) error {
				return WriteXMLBIF(buf, bn)
			},
			read: func(s string) (*models.BayesianNetwork, error) {
				return ReadXMLBIF(strings.NewReader(s))
			},
		},
		{
			name: "NET",
			write: func(buf *bytes.Buffer, bn *models.BayesianNetwork) error {
				return WriteNET(buf, bn)
			},
			read: func(s string) (*models.BayesianNetwork, error) {
				return ReadNET(strings.NewReader(s))
			},
		},
		{
			name: "XDSL",
			write: func(buf *bytes.Buffer, bn *models.BayesianNetwork) error {
				return WriteXDSL(buf, bn)
			},
			read: func(s string) (*models.BayesianNetwork, error) {
				return ReadXDSL(strings.NewReader(s))
			},
		},
	}

	for _, f := range formats {
		t.Run(f.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := f.write(&buf, bn); err != nil {
				t.Fatalf("Write%s failed: %v", f.name, err)
			}
			bn2, err := f.read(buf.String())
			if err != nil {
				t.Fatalf("Read%s failed: %v\nOutput:\n%s", f.name, err, buf.String())
			}
			nodes1 := bn.Nodes()
			nodes2 := bn2.Nodes()
			if len(nodes1) != len(nodes2) {
				t.Fatalf("node count mismatch: %d vs %d", len(nodes1), len(nodes2))
			}
			// Verify CPD data matches.
			for i, node := range nodes1 {
				cpd1 := bn.GetCPD(node)
				cpd2 := bn2.GetCPD(nodes2[i])
				if cpd1 == nil || cpd2 == nil {
					continue
				}
				d1 := cpd1.ToFactor().Values().Data()
				d2 := cpd2.ToFactor().Values().Data()
				assertFloatsClose(t, d2, d1, f.name+" "+node)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// UAI with multi-parent round trip & empty scope coverage
// ---------------------------------------------------------------------------

func TestReadUAI_EmptyScope(t *testing.T) {
	// scope with 0 vars, followed by 0 entries
	uai := "BAYES\n1\n2\n1\n0\n\n0\n"
	bn, err := ReadUAI(strings.NewReader(uai))
	if err != nil {
		t.Fatalf("ReadUAI failed: %v", err)
	}
	nodes := bn.Nodes()
	if len(nodes) != 1 {
		t.Errorf("expected 1 node, got %d", len(nodes))
	}
}

// Exercise read of PomdpX with ValueEnum
func TestReadPomdpX_ValueEnum(t *testing.T) {
	xml := `<?xml version="1.0"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="X" numValues="3">
      <ValueEnum>red green blue</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>X</Var>
      <Parameter>
        <Entry>
          <ProbTable>0.2 0.3 0.5</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </InitialStateBelief>
</pomdpx>`
	bn, err := ReadPomdpX(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v", err)
	}
	states := bn.GetStates("X")
	if len(states) != 3 || states[0] != "red" {
		t.Errorf("expected [red green blue], got %v", states)
	}
	cpd := bn.GetCPD("X")
	data := cpd.ToFactor().Values().Data()
	assertFloatsClose(t, data, []float64{0.2, 0.3, 0.5}, "PomdpX X CPD")
}

// Test ReadPomdpX with empty var name (should be skipped)
func TestReadPomdpX_EmptyVarName(t *testing.T) {
	xml := `<?xml version="1.0"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="" numValues="2"></StateVar>
    <StateVar vnamePrev="Y" numValues="2">
      <ValueEnum>y0 y1</ValueEnum>
    </StateVar>
  </Variable>
</pomdpx>`
	bn, err := ReadPomdpX(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("ReadPomdpX failed: %v", err)
	}
	nodes := bn.Nodes()
	if len(nodes) != 1 || nodes[0] != "Y" {
		t.Errorf("expected [Y], got %v", nodes)
	}
}

// Test ReadPomdpX with bad prob table float
func TestReadPomdpX_BadProbFloat(t *testing.T) {
	xml := `<?xml version="1.0"?>
<pomdpx version="1.0">
  <Variable>
    <StateVar vnamePrev="X" numValues="2">
      <ValueEnum>a b</ValueEnum>
    </StateVar>
  </Variable>
  <InitialStateBelief>
    <CondProb>
      <Var>X</Var>
      <Parameter>
        <Entry>
          <ProbTable>abc 0.5</ProbTable>
        </Entry>
      </Parameter>
    </CondProb>
  </InitialStateBelief>
</pomdpx>`
	_, err := ReadPomdpX(strings.NewReader(xml))
	if err == nil {
		t.Error("expected error for invalid float in ProbTable")
	}
}

// Test ReadBIF single-variable probability keyword
func TestReadBIF_SingleVarOnly(t *testing.T) {
	bif := `network test {
}
variable X {
  type discrete [ 2 ] { A, B };
}
probability ( X ) {
  table 0.5, 0.5;
}
`
	bn, err := ReadBIF(strings.NewReader(bif))
	if err != nil {
		t.Fatalf("ReadBIF failed: %v", err)
	}
	cpd := bn.GetCPD("X")
	if cpd == nil {
		t.Fatal("CPD is nil")
	}
}

// Test bifStripComment directly
func TestBifStripComment_NoComment(t *testing.T) {
	line := "table 0.5, 0.5;"
	result := bifStripComment(line)
	if result != line {
		t.Errorf("expected %q, got %q", line, result)
	}
}

func TestBifStripComment_WithComment(t *testing.T) {
	line := "table 0.5, 0.5; // comment"
	result := bifStripComment(line)
	if result != "table 0.5, 0.5; " {
		t.Errorf("expected %q, got %q", "table 0.5, 0.5; ", result)
	}
}

// Test net with no opening brace on same line as keyword
func TestReadNET_SeparateBrace(t *testing.T) {
	netData := `net
{
}
node X
{
  states = ("a" "b");
}
potential (X)
{
  data = (0.4 0.6);
}
`
	bn, err := ReadNET(strings.NewReader(netData))
	if err != nil {
		t.Fatalf("ReadNET failed: %v", err)
	}
	cpd := bn.GetCPD("X")
	data := cpd.ToFactor().Values().Data()
	assertFloatsClose(t, data, []float64{0.4, 0.6}, "X CPD")
}

// Exercise the XBN data mismatch fallback path
func TestReadXBN_DataMismatch(t *testing.T) {
	// Create XBN where DPI values count doesn't match expected
	xml := `<?xml version="1.0"?>
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
            <DPI>0.3</DPI>
          </DPIS>
        </DIST>
      </DISTRIBS>
    </DYNAMICPROPERTIES>
  </BNMODEL>
</ANALYSISNOTEBOOK>`
	bn, err := ReadXBN(strings.NewReader(xml))
	if err != nil {
		t.Fatalf("ReadXBN failed: %v", err)
	}
	cpd := bn.GetCPD("X")
	if cpd == nil {
		t.Fatal("X CPD is nil")
	}
	data := cpd.ToFactor().Values().Data()
	// Should have uniform distribution (fallback)
	assertFloatsClose(t, data, []float64{0.5, 0.5}, "X CPD fallback")
}

// WriteBIF with unconditional should write "table" keyword
func TestWriteBIF_Unconditional(t *testing.T) {
	bn := models.NewBayesianNetwork()
	must(t, bn.AddNode("X"))
	must(t, bn.SetStates("X", []string{"A", "B"}))
	cpd, err := factors.NewTabularCPD("X", 2, [][]float64{{0.3}, {0.7}}, nil, nil)
	must(t, err)
	must(t, bn.AddCPD(cpd))

	var buf bytes.Buffer
	must(t, WriteBIF(&buf, bn))
	if !strings.Contains(buf.String(), "table 0.3, 0.7;") {
		t.Errorf("expected table line in output, got:\n%s", buf.String())
	}
}
