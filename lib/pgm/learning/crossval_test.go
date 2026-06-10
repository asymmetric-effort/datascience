//go:build unit

package learning_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/asymmetric-effort/datascience/lib/pgm/learning"
	"github.com/asymmetric-effort/datascience/lib/pgm/models"
	"github.com/asymmetric-effort/datascience/lib/tabgo"
	"github.com/asymmetric-effort/datascience/tests/testutil"
)

// cpdFixture holds expected CPD data from fixtures.
type cpdFixture struct {
	Variable     string      `json:"variable"`
	VariableCard int         `json:"variable_card"`
	Values       [][]float64 `json:"values"`
	Evidence     []string    `json:"evidence"`
	EvidenceCard []int       `json:"evidence_card"`
}

// mleInput holds the input section of an MLE/Bayesian test case.
type mleInput struct {
	Edges                [][]string     `json:"edges"`
	NodeCards            map[string]int `json:"node_cards"`
	DataColumns          []string       `json:"data_columns"`
	DataRows             [][]float64    `json:"data_rows"`
	EquivalentSampleSize float64        `json:"equivalent_sample_size"`
}

// mleExpected holds the expected section of an MLE/Bayesian test case.
type mleExpected struct {
	CPDA cpdFixture `json:"cpd_A"`
	CPDB cpdFixture `json:"cpd_B"`
}

// buildNetworkFromInput creates a BayesianNetwork from fixture input edges.
func buildNetworkFromInput(t *testing.T, input *mleInput) *models.BayesianNetwork {
	t.Helper()
	bn := models.NewBayesianNetwork()
	nodeSet := make(map[string]bool)
	for _, edge := range input.Edges {
		nodeSet[edge[0]] = true
		nodeSet[edge[1]] = true
	}
	for node := range nodeSet {
		if err := bn.AddNode(node); err != nil {
			t.Fatalf("AddNode(%q): %v", node, err)
		}
	}
	for _, edge := range input.Edges {
		if err := bn.AddEdge(edge[0], edge[1]); err != nil {
			t.Fatalf("AddEdge(%q, %q): %v", edge[0], edge[1], err)
		}
	}
	return bn
}

// buildDataFrameInt creates a DataFrame with integer-typed values from fixture rows.
func buildDataFrameInt(columns []string, rows [][]float64) *tabgo.DataFrame {
	anyRows := make([][]any, len(rows))
	for i, row := range rows {
		anyRow := make([]any, len(row))
		for j, v := range row {
			anyRow[j] = int(v)
		}
		anyRows[i] = anyRow
	}
	return tabgo.NewDataFrameFromRows(columns, anyRows)
}

// buildDataFrameStr creates a DataFrame with string-typed values from fixture rows.
// States are string representations of integers (e.g., "0", "1").
func buildDataFrameStr(columns []string, rows [][]float64) *tabgo.DataFrame {
	anyRows := make([][]any, len(rows))
	for i, row := range rows {
		anyRow := make([]any, len(row))
		for j, v := range row {
			anyRow[j] = strconv.Itoa(int(v))
		}
		anyRows[i] = anyRow
	}
	return tabgo.NewDataFrameFromRows(columns, anyRows)
}

// setStatesFromCards sets string state names ("0", "1", ...) based on cardinalities.
func setStatesFromCards(t *testing.T, bn *models.BayesianNetwork, nodeCards map[string]int) {
	t.Helper()
	for node, card := range nodeCards {
		states := make([]string, card)
		for i := 0; i < card; i++ {
			states[i] = strconv.Itoa(i)
		}
		if err := bn.SetStates(node, states); err != nil {
			t.Fatalf("SetStates(%q): %v", node, err)
		}
	}
}

// assertCPDClose compares a computed TabularCPD against expected fixture values
// with the given tolerance.
func assertCPDClose(t *testing.T, label string, cpd *cpdFixture, gotValues []float64, gotCard int, tol float64) {
	t.Helper()

	if gotCard != cpd.VariableCard {
		t.Errorf("%s: variable_card: expected %d, got %d", label, cpd.VariableCard, gotCard)
	}

	// Flatten expected values (rows=states, cols=parent configs) in the same
	// row-major order as TabularCPD: flat[state * numParentConfigs + parentConfig].
	numParentConfigs := 1
	for _, ec := range cpd.EvidenceCard {
		numParentConfigs *= ec
	}
	expectedFlat := make([]float64, 0, cpd.VariableCard*numParentConfigs)
	for _, row := range cpd.Values {
		expectedFlat = append(expectedFlat, row...)
	}

	if len(gotValues) != len(expectedFlat) {
		t.Fatalf("%s: value count mismatch: expected %d, got %d", label, len(expectedFlat), len(gotValues))
	}

	for i := range expectedFlat {
		if math.Abs(gotValues[i]-expectedFlat[i]) > tol {
			t.Errorf("%s: value[%d]: expected %.6f, got %.6f (tol=%.4f)", label, i, expectedFlat[i], gotValues[i], tol)
		}
	}
}

func TestCrossval_MLEParameterEstimation(t *testing.T) {
	ff := testutil.LoadFixtures(t, "learning/fixtures.json")
	tc := ff.FindTestCase(t, "mle_parameter_estimation")

	var input mleInput
	tc.UnmarshalInput(t, &input)

	var expected mleExpected
	tc.UnmarshalExpected(t, &expected)

	// Build the network structure.
	bn := buildNetworkFromInput(t, &input)

	// Build DataFrame with integer values (MLE uses .Int()).
	df := buildDataFrameInt(input.DataColumns, input.DataRows)

	mle := learning.NewMLE(bn, df)

	// Test node A.
	cpdA, err := mle.GetParameters("A")
	if err != nil {
		t.Fatalf("MLE.GetParameters(A): %v", err)
	}
	gotValsA := cpdA.ToFactor().Values().Data()
	assertCPDClose(t, "MLE cpd_A", &expected.CPDA, gotValsA, cpdA.VariableCard(), 0.02)

	// Test node B.
	cpdB, err := mle.GetParameters("B")
	if err != nil {
		t.Fatalf("MLE.GetParameters(B): %v", err)
	}
	gotValsB := cpdB.ToFactor().Values().Data()
	assertCPDClose(t, "MLE cpd_B", &expected.CPDB, gotValsB, cpdB.VariableCard(), 0.02)
}

func TestCrossval_BayesianBDeuParameterEstimation(t *testing.T) {
	ff := testutil.LoadFixtures(t, "learning/fixtures.json")
	tc := ff.FindTestCase(t, "bayesian_bdeu_parameter_estimation")

	var input mleInput
	tc.UnmarshalInput(t, &input)

	var expected mleExpected
	tc.UnmarshalExpected(t, &expected)

	// Build the network structure.
	bn := buildNetworkFromInput(t, &input)

	// BayesianEstimator requires states to be set on the network.
	setStatesFromCards(t, bn, input.NodeCards)

	// Build DataFrame with string values to match string state names.
	df := buildDataFrameStr(input.DataColumns, input.DataRows)

	be := learning.NewBayesianEstimator(bn, df, learning.BDeu, input.EquivalentSampleSize)

	// Run Estimate to populate CPDs in the network.
	if err := be.Estimate(); err != nil {
		t.Fatalf("BayesianEstimator.Estimate(): %v", err)
	}

	// Test node A.
	cpdA, err := be.GetParameters("A")
	if err != nil {
		t.Fatalf("BayesianEstimator.GetParameters(A): %v", err)
	}
	gotValsA := cpdA.ToFactor().Values().Data()
	assertCPDClose(t, "Bayesian cpd_A", &expected.CPDA, gotValsA, cpdA.VariableCard(), 0.02)

	// Test node B.
	cpdB, err := be.GetParameters("B")
	if err != nil {
		t.Fatalf("BayesianEstimator.GetParameters(B): %v", err)
	}
	gotValsB := cpdB.ToFactor().Values().Data()
	assertCPDClose(t, "Bayesian cpd_B", &expected.CPDB, gotValsB, cpdB.VariableCard(), 0.02)
}
