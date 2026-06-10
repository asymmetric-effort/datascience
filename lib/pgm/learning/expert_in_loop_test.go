//go:build unit

package learning

import (
	"fmt"
	"strings"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/tabgo"
)

// ---------------------------------------------------------------------------
// Mock LLM client
// ---------------------------------------------------------------------------

// mockLLMClient implements LLMClient with canned responses keyed by substrings
// in the prompt. If no match is found, it returns a default response.
type mockLLMClient struct {
	// responses maps a substring to match in the prompt to the response to return.
	responses map[string]string
	// defaultResponse is returned when no key matches.
	defaultResponse string
	// calls tracks the number of Complete calls made.
	calls int
}

func newMockLLMClient() *mockLLMClient {
	return &mockLLMClient{
		responses:       make(map[string]string),
		defaultResponse: "UNKNOWN",
	}
}

func (m *mockLLMClient) Complete(prompt string, opts ...LLMOption) (string, error) {
	m.calls++
	for key, resp := range m.responses {
		if strings.Contains(prompt, key) {
			return resp, nil
		}
	}
	return m.defaultResponse, nil
}

func (m *mockLLMClient) ChatComplete(messages []Message, opts ...LLMOption) (string, error) {
	if len(messages) > 0 {
		return m.Complete(messages[len(messages)-1].Content, opts...)
	}
	return m.defaultResponse, nil
}

// errorLLMClient always returns an error, simulating an unavailable LLM.
type errorLLMClient struct{}

func (e *errorLLMClient) Complete(prompt string, opts ...LLMOption) (string, error) {
	return "", fmt.Errorf("llm unavailable")
}

func (e *errorLLMClient) ChatComplete(messages []Message, opts ...LLMOption) (string, error) {
	return "", fmt.Errorf("llm unavailable")
}

// ---------------------------------------------------------------------------
// Helper: create dummy data for tests
// ---------------------------------------------------------------------------

func expertDummyData(columns ...string) *tabgo.DataFrame {
	cols := make(map[string]*tabgo.Series)
	for _, c := range columns {
		cols[c] = tabgo.NewSeries(c, []any{0, 1})
	}
	return tabgo.NewDataFrame(cols)
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestExpertInLoop_VStructureWithLLMSupport tests the classic v-structure
// A -> B <- C where the LLM confirms the causal directions.
func TestExpertInLoop_VStructureWithLLMSupport(t *testing.T) {
	// True DAG: A -> B <- C (v-structure at B)
	// Independence structure:
	//   A _|_ C (marginally independent)
	//   A !_|_ B
	//   B !_|_ C
	//   A !_|_ C | B (explaining away)
	indeps := map[string]bool{
		ciKey("A", "B", nil):           false,
		ciKey("A", "C", nil):           true, // A and C are independent
		ciKey("B", "C", nil):           false,
		ciKey("A", "C", []string{"B"}): false, // not independent given B
		ciKey("A", "B", []string{"C"}): false,
		ciKey("B", "C", []string{"A"}): false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C")

	mock := newMockLLMClient()
	// LLM confirms: A causes B and C causes B.
	mock.responses["Does A cause B"] = "YES, confidence: 0.9"
	mock.responses["Does C cause B"] = "YES, confidence: 0.9"

	eil := NewExpertInLoop(data, mock, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("Estimate() error: %v", err)
	}

	// Expect v-structure: A -> B <- C.
	if !pdag.HasDirectedEdge("A", "B") {
		t.Errorf("expected directed edge A -> B")
	}
	if !pdag.HasDirectedEdge("C", "B") {
		t.Errorf("expected directed edge C -> B")
	}
	// A and C should NOT be adjacent.
	if pdag.Adjacent("A", "C") {
		t.Errorf("expected A and C to not be adjacent")
	}
	// LLM should have been consulted.
	if mock.calls == 0 {
		t.Errorf("expected LLM to be consulted, but calls = 0")
	}
}

// TestExpertInLoop_VStructureStatisticalOnly tests that v-structures are
// detected even when the LLM returns UNKNOWN (fallback to statistical).
func TestExpertInLoop_VStructureStatisticalOnly(t *testing.T) {
	// Same v-structure: A -> B <- C
	indeps := map[string]bool{
		ciKey("A", "B", nil):           false,
		ciKey("A", "C", nil):           true,
		ciKey("B", "C", nil):           false,
		ciKey("A", "C", []string{"B"}): false,
		ciKey("A", "B", []string{"C"}): false,
		ciKey("B", "C", []string{"A"}): false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C")

	// LLM returns UNKNOWN for everything.
	mock := newMockLLMClient()
	mock.defaultResponse = "UNKNOWN"

	eil := NewExpertInLoop(data, mock, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("Estimate() error: %v", err)
	}

	// Statistical evidence alone should orient the v-structure.
	if !pdag.HasDirectedEdge("A", "B") {
		t.Errorf("expected directed edge A -> B (statistical fallback)")
	}
	if !pdag.HasDirectedEdge("C", "B") {
		t.Errorf("expected directed edge C -> B (statistical fallback)")
	}
}

// TestExpertInLoop_LLMUnavailable tests that the algorithm gracefully handles
// an LLM that always returns errors.
func TestExpertInLoop_LLMUnavailable(t *testing.T) {
	indeps := map[string]bool{
		ciKey("A", "B", nil):           false,
		ciKey("A", "C", nil):           true,
		ciKey("B", "C", nil):           false,
		ciKey("A", "C", []string{"B"}): false,
		ciKey("A", "B", []string{"C"}): false,
		ciKey("B", "C", []string{"A"}): false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C")

	eil := NewExpertInLoop(data, &errorLLMClient{}, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("Estimate() error: %v", err)
	}

	// Should still orient the v-structure from statistical evidence.
	if !pdag.HasDirectedEdge("A", "B") {
		t.Errorf("expected directed edge A -> B (LLM unavailable fallback)")
	}
	if !pdag.HasDirectedEdge("C", "B") {
		t.Errorf("expected directed edge C -> B (LLM unavailable fallback)")
	}
}

// TestExpertInLoop_NilLLMClient tests that a nil LLM client is handled
// gracefully (pure statistical mode).
func TestExpertInLoop_NilLLMClient(t *testing.T) {
	indeps := map[string]bool{
		ciKey("A", "B", nil):           false,
		ciKey("A", "C", nil):           true,
		ciKey("B", "C", nil):           false,
		ciKey("A", "C", []string{"B"}): false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C")

	eil := NewExpertInLoop(data, nil, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("Estimate() error: %v", err)
	}

	// Statistical v-structure should still be detected.
	if !pdag.HasDirectedEdge("A", "B") {
		t.Errorf("expected directed edge A -> B")
	}
	if !pdag.HasDirectedEdge("C", "B") {
		t.Errorf("expected directed edge C -> B")
	}
}

// TestExpertInLoop_LLMInfluencesOrientation tests that the LLM can promote a
// non-v-structure to a v-structure when statistical evidence alone would not
// orient. This is the key integration: the LLM adds information beyond what
// the CI tests provide.
func TestExpertInLoop_LLMInfluencesOrientation(t *testing.T) {
	// Setup: A - B - C chain. Statistically, B IS in the separating set of
	// (A,C), so the PC algorithm would NOT orient a v-structure. But we
	// configure the LLM to say A->B and C->B (suggesting v-structure).
	indeps := map[string]bool{
		ciKey("A", "B", nil):           false,
		ciKey("A", "C", nil):           false,
		ciKey("B", "C", nil):           false,
		ciKey("A", "C", []string{"B"}): true, // B is in the sep set
		ciKey("A", "B", []string{"C"}): false,
		ciKey("B", "C", []string{"A"}): false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C")

	mock := newMockLLMClient()
	// LLM strongly believes A->B and C->B (v-structure).
	mock.responses["Does A cause B"] = "YES, confidence: 0.95"
	mock.responses["Does C cause B"] = "YES, confidence: 0.95"

	eil := NewExpertInLoop(data, mock, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("Estimate() error: %v", err)
	}

	// The LLM should have influenced the orientation. Since statistics say
	// B is in the separating set (no v-structure), but LLM says YES for both
	// directions, the combined signal should orient as a v-structure.
	if !pdag.HasDirectedEdge("A", "B") {
		t.Errorf("expected LLM to influence orientation: A -> B")
	}
	if !pdag.HasDirectedEdge("C", "B") {
		t.Errorf("expected LLM to influence orientation: C -> B")
	}
}

// TestExpertInLoop_LLMDoesNotOverrideWhenUncertain tests that when the LLM
// returns UNKNOWN, it does not override the statistical non-v-structure.
func TestExpertInLoop_LLMDoesNotOverrideWhenUncertain(t *testing.T) {
	// Chain A - B - C: B is in sepSet, so no v-structure statistically.
	indeps := map[string]bool{
		ciKey("A", "B", nil):           false,
		ciKey("A", "C", nil):           false,
		ciKey("B", "C", nil):           false,
		ciKey("A", "C", []string{"B"}): true,
		ciKey("A", "B", []string{"C"}): false,
		ciKey("B", "C", []string{"A"}): false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C")

	mock := newMockLLMClient()
	mock.defaultResponse = "UNKNOWN"

	eil := NewExpertInLoop(data, mock, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("Estimate() error: %v", err)
	}

	// With UNKNOWN LLM and no statistical v-structure, B should remain
	// connected to both A and C but edges should NOT be forcibly oriented
	// as a v-structure. The edges may be undirected or oriented by Meek rules.
	hasVStructure := pdag.HasDirectedEdge("A", "B") && pdag.HasDirectedEdge("C", "B")
	if hasVStructure {
		t.Errorf("did not expect v-structure when LLM is uncertain and stats say no v-structure")
	}
}

// TestExpertInLoop_ChainNoVStructure tests the A -> B -> C chain where no
// v-structure should be detected.
func TestExpertInLoop_ChainNoVStructure(t *testing.T) {
	// Chain: A -> B -> C
	// A _|_ C | B, and B IS in sepSet(A,C), so no v-structure.
	indeps := map[string]bool{
		ciKey("A", "B", nil):           false,
		ciKey("A", "C", nil):           false,
		ciKey("B", "C", nil):           false,
		ciKey("A", "C", []string{"B"}): true,
		ciKey("A", "B", []string{"C"}): false,
		ciKey("B", "C", []string{"A"}): false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C")

	mock := newMockLLMClient()
	mock.defaultResponse = "UNKNOWN"

	eil := NewExpertInLoop(data, mock, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("Estimate() error: %v", err)
	}

	// Skeleton should have edges A-B and B-C but not A-C.
	if pdag.Adjacent("A", "C") {
		t.Errorf("expected A and C to not be adjacent in chain structure")
	}
	// A-B and B-C should exist in some form (directed or undirected).
	if !pdag.Adjacent("A", "B") {
		t.Errorf("expected A and B to be adjacent")
	}
	if !pdag.Adjacent("B", "C") {
		t.Errorf("expected B and C to be adjacent")
	}
}

// TestExpertInLoop_EstimateBN tests the full pipeline including conversion to
// a BayesianNetwork.
func TestExpertInLoop_EstimateBN(t *testing.T) {
	indeps := map[string]bool{
		ciKey("A", "B", nil):           false,
		ciKey("A", "C", nil):           true,
		ciKey("B", "C", nil):           false,
		ciKey("A", "C", []string{"B"}): false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C")

	mock := newMockLLMClient()
	mock.responses["Does A cause B"] = "YES"
	mock.responses["Does C cause B"] = "YES"

	eil := NewExpertInLoop(data, mock, ciTest, 0.05)
	bn, err := eil.EstimateBN()
	if err != nil {
		t.Fatalf("EstimateBN() error: %v", err)
	}

	nodes := bn.Nodes()
	if len(nodes) != 3 {
		t.Errorf("expected 3 nodes, got %d", len(nodes))
	}

	// The BN should have edges; exact orientation depends on PDAG->DAG conversion.
	edges := bn.Edges()
	if len(edges) == 0 {
		t.Errorf("expected at least one edge in the BN")
	}
}

// TestExpertInLoop_TooFewVariables tests that Estimate returns an error when
// fewer than 2 variables are provided.
func TestExpertInLoop_TooFewVariables(t *testing.T) {
	data := expertDummyData("A")
	mock := newMockLLMClient()
	ciTest := mockCITestFromMap(nil)

	eil := NewExpertInLoop(data, mock, ciTest, 0.05)
	_, err := eil.Estimate()
	if err == nil {
		t.Error("expected error for fewer than 2 variables")
	}
}

// TestExpertInLoop_FourNodeDiamond tests a more complex structure with 4 nodes.
// True DAG: A -> B, A -> C, B -> D, C -> D (diamond with v-structure at D).
func TestExpertInLoop_FourNodeDiamond(t *testing.T) {
	indeps := map[string]bool{
		// Marginal independences.
		ciKey("A", "B", nil): false,
		ciKey("A", "C", nil): false,
		ciKey("A", "D", nil): false,
		ciKey("B", "C", nil): false,
		ciKey("B", "D", nil): false,
		ciKey("C", "D", nil): false,
		// Conditional independences.
		ciKey("B", "C", []string{"A"}):      true,  // B _|_ C | A
		ciKey("A", "D", []string{"B", "C"}): true,  // A _|_ D | {B,C}
		ciKey("A", "D", []string{"B"}):      false, // not independent given just B
		ciKey("A", "D", []string{"C"}):      false, // not independent given just C
		ciKey("B", "C", []string{"D"}):      false, // explaining away
		ciKey("B", "D", []string{"A"}):      false,
		ciKey("B", "D", []string{"C"}):      false,
		ciKey("C", "D", []string{"A"}):      false,
		ciKey("C", "D", []string{"B"}):      false,
		ciKey("A", "B", []string{"C"}):      false,
		ciKey("A", "B", []string{"D"}):      false,
		ciKey("A", "C", []string{"B"}):      false,
		ciKey("A", "C", []string{"D"}):      false,
	}
	ciTest := mockCITestFromMap(indeps)
	data := expertDummyData("A", "B", "C", "D")

	mock := newMockLLMClient()
	mock.responses["Does B cause D"] = "YES, confidence: 0.9"
	mock.responses["Does C cause D"] = "YES, confidence: 0.9"
	mock.defaultResponse = "UNKNOWN"

	eil := NewExpertInLoop(data, mock, ciTest, 0.05)
	pdag, err := eil.Estimate()
	if err != nil {
		t.Fatalf("Estimate() error: %v", err)
	}

	// B and C should not be adjacent (separated by A).
	if pdag.Adjacent("B", "C") {
		t.Errorf("expected B and C to not be adjacent (separated by A)")
	}

	// V-structure at D: B -> D <- C.
	// The sep set of (B,C) is {A}, and D is not in it, so statistical evidence
	// supports the v-structure at D.
	if !pdag.HasDirectedEdge("B", "D") {
		t.Errorf("expected directed edge B -> D (v-structure at D)")
	}
	if !pdag.HasDirectedEdge("C", "D") {
		t.Errorf("expected directed edge C -> D (v-structure at D)")
	}
}

// TestExpertInLoop_CombineSignals tests the combineSignals method directly.
func TestExpertInLoop_CombineSignals(t *testing.T) {
	eil := &ExpertInLoop{}

	tests := []struct {
		name        string
		statistical bool
		llm         llmOpinion
		want        bool
	}{
		{"stats_yes_llm_supports", true, llmSupports, true},
		{"stats_yes_llm_opposes", true, llmOpposes, true},
		{"stats_yes_llm_uncertain", true, llmUncertain, true},
		{"stats_no_llm_supports", false, llmSupports, true},
		{"stats_no_llm_opposes", false, llmOpposes, false},
		{"stats_no_llm_uncertain", false, llmUncertain, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := eil.combineSignals(tt.statistical, tt.llm)
			if got != tt.want {
				t.Errorf("combineSignals(%v, %v) = %v, want %v",
					tt.statistical, tt.llm, got, tt.want)
			}
		})
	}
}

// TestExpertInLoop_QueryLLMForVStructure tests the LLM query logic.
func TestExpertInLoop_QueryLLMForVStructure(t *testing.T) {
	tmpl := CausalPromptTemplate{}

	t.Run("both_yes", func(t *testing.T) {
		mock := newMockLLMClient()
		mock.responses["Does A cause B"] = "YES"
		mock.responses["Does C cause B"] = "YES"
		eil := &ExpertInLoop{llmClient: mock}
		got := eil.queryLLMForVStructure(tmpl, "A", "C", "B")
		if got != llmSupports {
			t.Errorf("expected llmSupports, got %d", got)
		}
	})

	t.Run("one_no", func(t *testing.T) {
		mock := newMockLLMClient()
		mock.responses["Does A cause B"] = "YES"
		mock.responses["Does C cause B"] = "NO"
		eil := &ExpertInLoop{llmClient: mock}
		got := eil.queryLLMForVStructure(tmpl, "A", "C", "B")
		if got != llmOpposes {
			t.Errorf("expected llmOpposes, got %d", got)
		}
	})

	t.Run("unknown", func(t *testing.T) {
		mock := newMockLLMClient()
		mock.defaultResponse = "UNKNOWN"
		eil := &ExpertInLoop{llmClient: mock}
		got := eil.queryLLMForVStructure(tmpl, "A", "C", "B")
		if got != llmUncertain {
			t.Errorf("expected llmUncertain, got %d", got)
		}
	})

	t.Run("nil_client", func(t *testing.T) {
		eil := &ExpertInLoop{llmClient: nil}
		got := eil.queryLLMForVStructure(tmpl, "A", "C", "B")
		if got != llmUncertain {
			t.Errorf("expected llmUncertain for nil client, got %d", got)
		}
	})

	t.Run("error_client", func(t *testing.T) {
		eil := &ExpertInLoop{llmClient: &errorLLMClient{}}
		got := eil.queryLLMForVStructure(tmpl, "A", "C", "B")
		if got != llmUncertain {
			t.Errorf("expected llmUncertain for error client, got %d", got)
		}
	})
}
