//go:build unit

package models

import (
	"os"
	"strings"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/factors"
)

// ---------------------------------------------------------------------------
// LG BN: lgWriteImpl - exercise ALL write failure paths.
// ---------------------------------------------------------------------------
func TestLGWriteImpl_ExhaustiveFailure(t *testing.T) {
	lgbn := buildTestLGBN()
	var buf strings.Builder
	lgWriteImpl(&buf, lgbn)
	totalSize := buf.Len()

	for offset := 0; offset <= totalSize; offset++ {
		w := &lgFailingWriter{failAfter: offset}
		err := lgWriteImpl(w, lgbn)
		if offset < totalSize && err == nil {
			t.Fatalf("expected write failure at offset %d/%d", offset, totalSize)
		}
	}
}

// ---------------------------------------------------------------------------
// BIF writeBIF: exercise ALL write failure paths.
// ---------------------------------------------------------------------------
func TestWriteBIF_ExhaustiveFailure(t *testing.T) {
	bn := helperSimple2NodeBN(t)
	var buf strings.Builder
	bn.writeBIF(&buf)
	totalSize := buf.Len()

	for offset := 0; offset <= totalSize; offset++ {
		w := &failingWriter{failAfter: offset}
		_ = writeBIFImpl(w, bn)
	}
}

// ---------------------------------------------------------------------------
// LG BN: IsIMap - non-adjacent pair with a > b swap.
// ---------------------------------------------------------------------------
func TestLinearGaussianBN_IsIMap_NonAdjacentSwap(t *testing.T) {
	lgbn := NewLinearGaussianBayesianNetwork()
	lgbn.AddNode("A")
	lgbn.AddNode("B")
	lgbn.AddNode("C")
	lgbn.AddEdge("A", "C")
	lgbn.AddEdge("B", "C")
	cpdA, _ := factors.NewLinearGaussianCPD("A", 0, nil, 1, nil)
	cpdB, _ := factors.NewLinearGaussianCPD("B", 0, nil, 1, nil)
	cpdC, _ := factors.NewLinearGaussianCPD("C", 0, []float64{0.5, 0.3}, 1, []string{"A", "B"})
	lgbn.AddLinearGaussianCPD(cpdA)
	lgbn.AddLinearGaussianCPD(cpdB)
	lgbn.AddLinearGaussianCPD(cpdC)
	assertions := []IndependenceAssertion{
		{Event1: []string{"A"}, Event2: []string{"B"}, Given: nil},
		{Event1: []string{"B"}, Event2: []string{"A"}, Given: nil},
	}
	result, err := lgbn.IsIMap(assertions)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("IsIMap result: %v", result)
}

// ---------------------------------------------------------------------------
// BN: IsIMap - v-structure with JPD (exercises d-sep + JPD check path).
// ---------------------------------------------------------------------------
func TestBN_IsIMap_VStructureWithJPD(t *testing.T) {
	bn := NewBayesianNetwork()
	bn.AddNode("A")
	bn.AddNode("B")
	bn.AddNode("C")
	bn.AddEdge("A", "C")
	bn.AddEdge("B", "C")
	bn.SetStates("A", []string{"a0", "a1"})
	bn.SetStates("B", []string{"b0", "b1"})
	bn.SetStates("C", []string{"c0", "c1"})
	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	cpdC, _ := factors.NewTabularCPD("C", 2, [][]float64{{0.9, 0.6, 0.7, 0.1}, {0.1, 0.4, 0.3, 0.9}}, []string{"A", "B"}, []int{2, 2})
	bn.AddCPD(cpdA)
	bn.AddCPD(cpdB)
	bn.AddCPD(cpdC)

	pA := []float64{0.5, 0.5}
	pB := []float64{0.5, 0.5}
	pC := [][]float64{{0.9, 0.1}, {0.6, 0.4}, {0.7, 0.3}, {0.1, 0.9}}
	vals := make([]float64, 8)
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			for c := 0; c < 2; c++ {
				vals[a*4+b*2+c] = pA[a] * pB[b] * pC[a*2+b][c]
			}
		}
	}
	jpd, err := factors.NewJointProbabilityDistribution(
		[]string{"A", "B", "C"}, []int{2, 2, 2}, vals,
	)
	if err != nil {
		t.Fatalf("unexpected error creating JPD: %v", err)
	}
	result, err := bn.IsIMap(jpd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("IsIMap result: %v", result)
}

// ---------------------------------------------------------------------------
// BIF: loadBIF - conditional line not starting with "(".
// ---------------------------------------------------------------------------
func TestLoadBIF_ConditionalNonParenLine(t *testing.T) {
	input := `network unknown {
}
variable X {
  type discrete [ 2 ] { a, b };
}
variable Y {
  type discrete [ 2 ] { y0, y1 };
}
probability ( Y | X ) {
  something_else
  (a) 0.7, 0.3;
  (b) 0.4, 0.6;
}
`
	bn, err := loadBIF(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bn.Nodes()) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(bn.Nodes()))
	}
}

// ---------------------------------------------------------------------------
// MN: ToBayesianModel - 4 nodes to exercise more marginalization.
// ---------------------------------------------------------------------------
func TestMN_ToBayesianModel_FourNodes(t *testing.T) {
	mn := NewMarkovNetwork()
	mn.AddNode("A")
	mn.AddNode("B")
	mn.AddNode("C")
	mn.AddNode("D")
	mn.AddEdge("A", "B")
	mn.AddEdge("B", "C")
	mn.AddEdge("C", "D")
	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.3, 0.2, 0.4, 0.1})
	fBC, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{0.3, 0.2, 0.4, 0.1})
	fCD, _ := factors.NewDiscreteFactor([]string{"C", "D"}, []int{2, 2}, []float64{0.3, 0.2, 0.4, 0.1})
	mn.AddFactor(fAB)
	mn.AddFactor(fBC)
	mn.AddFactor(fCD)
	bn, err := mn.ToBayesianModel()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bn.Nodes()) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(bn.Nodes()))
	}
}

// ---------------------------------------------------------------------------
// SEM: FromLavaan - all line types.
// ---------------------------------------------------------------------------
func TestSEM_FromLavaan_AllLineTypesV2(t *testing.T) {
	syntax := "no_tilde\nY ~ X1 + X2\nZ ~ Y\nW ~\n"
	s, err := FromLavaan(syntax)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Variables()) < 4 {
		t.Fatalf("expected at least 4 variables, got %d", len(s.Variables()))
	}
}

// ---------------------------------------------------------------------------
// SEM: FromLisrel - empty rest.
// ---------------------------------------------------------------------------
func TestSEM_FromLisrel_EmptyRest(t *testing.T) {
	s, err := FromLisrel("X:")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eq := s.GetEquation("X")
	if eq == nil {
		t.Fatal("expected equation for X")
	}
}

// ---------------------------------------------------------------------------
// LG BN: LoadLinearGaussianBayesianNetwork - default case.
// ---------------------------------------------------------------------------
func TestLGBN_Load_DefaultCase(t *testing.T) {
	content := "network lg_bayesian_network {\n}\nunknown_keyword foo {\n  bar\n}\nvariable X {\n  type continuous;\n}\ndistribution X {\n  mean 0;\n  variance 1;\n}\n"
	tmpFile := "/tmp/test_lgbn_default_case.txt"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)
	bn, err := LoadLinearGaussianBayesianNetwork(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bn.Nodes()) != 1 {
		t.Fatalf("expected 1 node, got %d", len(bn.Nodes()))
	}
}

// ---------------------------------------------------------------------------
// LG BN: LoadLinearGaussianBayesianNetwork - empty tokens (line 244 skip).
// ---------------------------------------------------------------------------
func TestLGBN_Load_EmptyTokens(t *testing.T) {
	// A line that becomes empty after parsing should be skipped.
	content := "network lg_bayesian_network {\n}\n\nvariable X {\n  type continuous;\n}\ndistribution X {\n  mean 0;\n  variance 1;\n}\n"
	tmpFile := "/tmp/test_lgbn_empty_tokens.txt"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)
	bn, err := LoadLinearGaussianBayesianNetwork(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bn.Nodes()) != 1 {
		t.Fatalf("expected 1 node, got %d", len(bn.Nodes()))
	}
}

// ---------------------------------------------------------------------------
// SEM: FromGraph with edges.
// ---------------------------------------------------------------------------
func TestSEM_FromGraph_MultiEdgeV2(t *testing.T) {
	bn := NewBayesianNetwork()
	bn.AddNode("X")
	bn.AddNode("Y")
	bn.AddNode("Z")
	bn.AddEdge("X", "Y")
	bn.AddEdge("X", "Z")
	bn.AddEdge("Y", "Z")
	s, err := FromGraph(bn.dag)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eq := s.GetEquation("Z")
	if len(eq.Parents) != 2 {
		t.Fatalf("expected 2 parents for Z, got %d", len(eq.Parents))
	}
}

// ---------------------------------------------------------------------------
// BN: Simulate - exercise the rejection sampling "only N accepted" path.
// ---------------------------------------------------------------------------
func TestBN_Simulate_FewAccepted(t *testing.T) {
	bn := NewBayesianNetwork()
	bn.AddNode("A")
	bn.SetStates("A", []string{"a0", "a1"})
	// Very low probability for state 0.
	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.001}, {0.999}}, nil, nil)
	bn.AddCPD(cpdA)
	// Request many samples with rare evidence A=0.
	// maxAttempts = 100*1000 = 100000, P(A=0) = 0.001
	// Expected accepted ~100, but we need 100. Should barely work.
	_, err := bn.Simulate(100, map[string]int{"A": 0}, 42)
	if err != nil {
		// The "only N accepted" error path.
		if strings.Contains(err.Error(), "accepted") {
			t.Logf("hit rejection exhaustion path: %v", err)
			return
		}
		t.Fatalf("unexpected error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// BN: Simulate - TopologicalOrder error fallback.
// This requires a BN where TopologicalOrder fails, which needs a cycle.
// Can't happen with valid DAG.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// LG BN: Load with betas comma-separated.
// ---------------------------------------------------------------------------
func TestLGBN_Load_WithBetas(t *testing.T) {
	content := "network lg_bayesian_network {\n}\nvariable X {\n  type continuous;\n}\nvariable Y {\n  type continuous;\n}\ndistribution X {\n  mean 1;\n  variance 2;\n}\ndistribution Y | X {\n  mean 0.5;\n  variance 1;\n  betas 0.8;\n}\n"
	tmpFile := "/tmp/test_lgbn_betas.txt"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)
	bn, err := LoadLinearGaussianBayesianNetwork(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(bn.Nodes()) != 2 {
		t.Fatalf("expected 2 nodes, got %d", len(bn.Nodes()))
	}
}

// ---------------------------------------------------------------------------
// BN: IsIMap - 4-node chain to exercise d-sep with parents (lines 340-347).
// ---------------------------------------------------------------------------
func TestBN_IsIMap_ChainWithParents(t *testing.T) {
	// Chain: A -> B -> C -> D
	// A and C are non-adjacent. Parents(A)={}, so zSet is empty.
	// A and D are non-adjacent. Parents(A)={}, so zSet is empty.
	// B and D are non-adjacent. Parents(B)={A}, so zSet={A}.
	// DSep(B, D | {A}) should hold in a chain.
	bn := NewBayesianNetwork()
	bn.AddNode("A")
	bn.AddNode("B")
	bn.AddNode("C")
	bn.AddNode("D")
	bn.AddEdge("A", "B")
	bn.AddEdge("B", "C")
	bn.AddEdge("C", "D")
	bn.SetStates("A", []string{"a0", "a1"})
	bn.SetStates("B", []string{"b0", "b1"})
	bn.SetStates("C", []string{"c0", "c1"})
	bn.SetStates("D", []string{"d0", "d1"})
	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.6}, {0.4}}, nil, nil)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.9, 0.2}, {0.1, 0.8}}, []string{"A"}, []int{2})
	cpdC, _ := factors.NewTabularCPD("C", 2, [][]float64{{0.8, 0.3}, {0.2, 0.7}}, []string{"B"}, []int{2})
	cpdD, _ := factors.NewTabularCPD("D", 2, [][]float64{{0.7, 0.4}, {0.3, 0.6}}, []string{"C"}, []int{2})
	bn.AddCPD(cpdA)
	bn.AddCPD(cpdB)
	bn.AddCPD(cpdC)
	bn.AddCPD(cpdD)

	// Build JPD from the BN (compute exact joint).
	// P(A,B,C,D) = P(A)P(B|A)P(C|B)P(D|C)
	vals := make([]float64, 16)
	pA := []float64{0.6, 0.4}
	pBA := [][]float64{{0.9, 0.2}, {0.1, 0.8}}
	pCB := [][]float64{{0.8, 0.3}, {0.2, 0.7}}
	pDC := [][]float64{{0.7, 0.4}, {0.3, 0.6}}
	for a := 0; a < 2; a++ {
		for b := 0; b < 2; b++ {
			for c := 0; c < 2; c++ {
				for d := 0; d < 2; d++ {
					vals[a*8+b*4+c*2+d] = pA[a] * pBA[b][a] * pCB[c][b] * pDC[d][c]
				}
			}
		}
	}
	jpd, err := factors.NewJointProbabilityDistribution(
		[]string{"A", "B", "C", "D"}, []int{2, 2, 2, 2}, vals,
	)
	if err != nil {
		t.Fatalf("unexpected error creating JPD: %v", err)
	}
	result, err := bn.IsIMap(jpd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Since the JPD is derived from the BN, all implied independencies hold.
	if !result {
		t.Fatal("expected IsIMap=true for exact JPD")
	}
}

// BN: IsIMap - independence fails in JPD (line 345-347).
func TestBN_IsIMap_IndependenceFails(t *testing.T) {
	// A -> B -> C, but JPD has A and C correlated given {}
	// (which should be the case - they're not d-separated given {}).
	// Need a case where d-sep holds but JPD doesn't confirm it.
	// In chain A->B->C: A _|_ C | {parents(A)} = A _|_ C | {} (since A has no parents).
	// But in a chain, A and C are NOT d-separated given {} (since B is not observed).
	// So DSeparation returns false and this block is skipped.
	// For d-sep to hold: need parents(a) to block the path.
	// B and D in A->B->C->D: B _|_ D | {A} (since A is parent of B).
	// But B and D are NOT d-separated given {A} in a chain (path B->C->D exists).
	// Need a v-structure or fork.
	// Fork: A -> B, A -> C. B _|_ C | {parents(B)} = B _|_ C | {A}. YES d-separated.
	// But JPD where B _|_ C | {A} does NOT hold -> IsIMap returns false.
	bn := NewBayesianNetwork()
	bn.AddNode("A")
	bn.AddNode("B")
	bn.AddNode("C")
	bn.AddEdge("A", "B")
	bn.AddEdge("A", "C")
	bn.SetStates("A", []string{"a0", "a1"})
	bn.SetStates("B", []string{"b0", "b1"})
	bn.SetStates("C", []string{"c0", "c1"})
	cpdA, _ := factors.NewTabularCPD("A", 2, [][]float64{{0.5}, {0.5}}, nil, nil)
	cpdB, _ := factors.NewTabularCPD("B", 2, [][]float64{{0.9, 0.1}, {0.1, 0.9}}, []string{"A"}, []int{2})
	cpdC, _ := factors.NewTabularCPD("C", 2, [][]float64{{0.9, 0.1}, {0.1, 0.9}}, []string{"A"}, []int{2})
	bn.AddCPD(cpdA)
	bn.AddCPD(cpdB)
	bn.AddCPD(cpdC)

	// Create a JPD where B and C are NOT independent given A.
	// Just use a random joint that violates the independence.
	vals := []float64{
		0.20, 0.05, 0.15, 0.10, // A=0: (B=0,C=0), (B=0,C=1), (B=1,C=0), (B=1,C=1)
		0.05, 0.15, 0.10, 0.20, // A=1
	}
	jpd, err := factors.NewJointProbabilityDistribution(
		[]string{"A", "B", "C"}, []int{2, 2, 2}, vals,
	)
	if err != nil {
		t.Fatalf("unexpected error creating JPD: %v", err)
	}
	result, err := bn.IsIMap(jpd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// B _|_ C | A should fail in this JPD -> IsIMap returns false.
	if result {
		t.Fatal("expected IsIMap=false since B _|_ C | A fails in JPD")
	}
}

// ---------------------------------------------------------------------------
// MN: ToFactorGraph and ToBayesianModel - larger network coverage.
// ---------------------------------------------------------------------------
func TestMN_ToFactorGraph_ThreeNodeChain(t *testing.T) {
	mn := NewMarkovNetwork()
	mn.AddNode("A")
	mn.AddNode("B")
	mn.AddNode("C")
	mn.AddEdge("A", "B")
	mn.AddEdge("B", "C")
	fAB, _ := factors.NewDiscreteFactor([]string{"A", "B"}, []int{2, 2}, []float64{0.3, 0.2, 0.4, 0.1})
	fBC, _ := factors.NewDiscreteFactor([]string{"B", "C"}, []int{2, 2}, []float64{0.3, 0.2, 0.4, 0.1})
	mn.AddFactor(fAB)
	mn.AddFactor(fBC)
	fg, err := mn.ToFactorGraph()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fg.GetVariables()) != 3 {
		t.Fatalf("expected 3 variables, got %d", len(fg.GetVariables()))
	}
}
