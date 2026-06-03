//go:build unit

package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/asymmetric-effort/pgmgo/lib/tabgo"
	"github.com/asymmetric-effort/pgmgo/src/models"
)

// ---------------------------------------------------------------------------
// run() tests — exercises the dispatch extracted from main()
// ---------------------------------------------------------------------------

func TestRun_NoArgs(t *testing.T) {
	code := run([]string{"pgmgo"})
	if code != 0 {
		t.Errorf("expected exit code 0 for no args, got %d", code)
	}
}

func TestRun_Version(t *testing.T) {
	for _, arg := range []string{"version", "--version", "-v"} {
		code := run([]string{"pgmgo", arg})
		if code != 0 {
			t.Errorf("expected exit code 0 for %q, got %d", arg, code)
		}
	}
}

func TestRun_Help(t *testing.T) {
	for _, arg := range []string{"help", "--help", "-h"} {
		code := run([]string{"pgmgo", arg})
		if code != 0 {
			t.Errorf("expected exit code 0 for %q, got %d", arg, code)
		}
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	code := run([]string{"pgmgo", "nonexistent_command"})
	if code != 1 {
		t.Errorf("expected exit code 1 for unknown command, got %d", code)
	}
}

func TestRun_Validate(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := run([]string{"pgmgo", "validate", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Query(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := run([]string{"pgmgo", "query", "--variables", "Rain", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Map(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := run([]string{"pgmgo", "map", "--variables", "Rain", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Info(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := run([]string{"pgmgo", "info", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Learn(t *testing.T) {
	csvContent := "A,B\n0,0\n0,1\n1,0\n1,1\n0,0\n0,1\n1,0\n1,1\n0,0\n1,1\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := run([]string{"pgmgo", "learn", "--data", csvPath, "--output", outPath})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Fit(t *testing.T) {
	modelPath := writeTempBIF(t, validBIF)
	csvContent := "Rain,Sprinkler\n0,0\n0,1\n1,0\n1,1\n0,0\n1,0\n0,1\n1,1\n0,0\n1,0\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "fitted.bif")
	code := run([]string{"pgmgo", "fit", "--model", modelPath, "--data", csvPath, "--output", outPath})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Sample(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	outPath := filepath.Join(t.TempDir(), "samples.csv")
	code := run([]string{"pgmgo", "sample", "--model", path, "--n", "5", "--output", outPath, "--seed", "42"})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Convert(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	outPath := filepath.Join(t.TempDir(), "out.xmlbif")
	code := run([]string{"pgmgo", "convert", "--input", path, "--from", "bif", "--to", "xmlbif", "--output", outPath})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Compare(t *testing.T) {
	path1 := writeTempBIF(t, validBIF)
	path2 := writeTempBIF(t, validBIF2)
	code := run([]string{"pgmgo", "compare", "--true", path1, "--estimated", path2})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRun_Do(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := run([]string{"pgmgo", "do", "--intervention", "Rain=0", "--query", "Sprinkler", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// failingBIFWriter is a test mock that simulates write failure for coverage testing.
// ---------------------------------------------------------------------------

type failingBIFWriter struct{}

func (f *failingBIFWriter) Write(w io.Writer, bn *models.BayesianNetwork) error {
	return errors.New("mock write error")
}

// ---------------------------------------------------------------------------
// convertModelImpl — write error via failing writer mock
// ---------------------------------------------------------------------------

// TestConvertModelImpl_WriteError exercises the write error path in convertModelImpl
// by injecting a failing bifWriter mock.
func TestConvertModelImpl_WriteError(t *testing.T) {
	bifPath := writeTempBIF(t, validBIF)
	outPath := filepath.Join(t.TempDir(), "out.bif")

	failWriters := map[string]bifWriter{
		"bif": &failingBIFWriter{},
	}

	code := convertModelImpl(bifPath, "bif", "bif", outPath, failWriters)
	if code != 1 {
		t.Errorf("expected exit code 1 for write error, got %d", code)
	}
}

// TestConvertModelImpl_UnknownOutputFormat exercises the unknown-output-format path
// with the new writer map approach.
func TestConvertModelImpl_UnknownOutputFormat(t *testing.T) {
	bifPath := writeTempBIF(t, validBIF)
	outPath := filepath.Join(t.TempDir(), "out.bif")

	code := convertModelImpl(bifPath, "bif", "badformat", outPath, formatWriterMap)
	if code != 2 {
		t.Errorf("expected exit code 2 for unknown output format, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// learnStructureFinalizeImpl — MLE failure fallback and write error
// ---------------------------------------------------------------------------

// TestLearnStructureFinalizeImpl_MLEFails_RandomCPDs exercises the MLE failure
// path where the fallback generates random CPDs, including the nStates > 2 branch.
func TestLearnStructureFinalizeImpl_MLEFails_RandomCPDs(t *testing.T) {
	// Create a BN with nodes that have >2 states and data that causes MLE to fail.
	bn := models.NewBayesianNetwork()
	bn.AddNode("X")
	bn.SetStates("X", []string{"a", "b", "c"}) // 3 states, triggers nStates > 2
	bn.AddNode("Y")
	bn.SetStates("Y", []string{"d", "e"})
	bn.AddEdge("X", "Y")

	// Data with columns that don't match the BN nodes, so MLE fails.
	m := map[string]*tabgo.Series{
		"Z": tabgo.NewSeries("Z", []any{0, 1, 0, 1}),
	}
	data := tabgo.NewDataFrame(m)

	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructureFinalizeImpl(bn, data, outPath, defaultBIFWriter)
	// Whether MLE fails depends on implementation, just exercise the path.
	_ = code
}

// TestLearnStructureFinalizeImpl_MLEFails_MoreStates exercises the nStates > 2
// branch with a node having many states.
func TestLearnStructureFinalizeImpl_MLEFails_MoreStates(t *testing.T) {
	bn := models.NewBayesianNetwork()
	bn.AddNode("A")
	bn.SetStates("A", []string{"s0", "s1", "s2", "s3", "s4"}) // 5 states

	// Provide mismatched data so MLE fails.
	m := map[string]*tabgo.Series{
		"Q": tabgo.NewSeries("Q", []any{0, 1}),
	}
	data := tabgo.NewDataFrame(m)

	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructureFinalizeImpl(bn, data, outPath, defaultBIFWriter)
	_ = code
}

// TestLearnStructureFinalizeImpl_WriteError exercises the write error path
// by injecting a failing bifWriter mock.
func TestLearnStructureFinalizeImpl_WriteError(t *testing.T) {
	// Create a minimal valid BN with CPDs so MLE succeeds.
	bn := models.NewBayesianNetwork()
	bn.AddNode("X")
	bn.SetStates("X", []string{"a", "b"})
	bn.GetRandomCPDs(2, 0)

	m := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{0, 1, 0, 1}),
	}
	data := tabgo.NewDataFrame(m)

	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructureFinalizeImpl(bn, data, outPath, &failingBIFWriter{})
	if code != 1 {
		t.Errorf("expected exit code 1 for write error, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// setStatesFromDataImpl — nil column path via mock columnLookup
// ---------------------------------------------------------------------------

// nilColumnLookup is a test mock that simulates missing columns for coverage testing.
type nilColumnLookup struct{}

func (n *nilColumnLookup) Column(name string) *tabgo.Series {
	return nil
}

// TestSetStatesFromDataImpl_NilColumn exercises the col == nil defensive path
// in setStatesFromDataImpl by injecting a nilColumnLookup mock.
func TestSetStatesFromDataImpl_NilColumn(t *testing.T) {
	bn := models.NewBayesianNetwork()
	bn.AddNode("X")
	// No states set for X, so the function should try to look up the column.

	m := map[string]*tabgo.Series{
		"X": tabgo.NewSeries("X", []any{0, 1}),
	}
	data := tabgo.NewDataFrame(m)

	// Inject a lookup that always returns nil.
	setStatesFromDataImpl(bn, data, &nilColumnLookup{})

	// X should still have no states (column was nil).
	states := bn.GetStates("X")
	if len(states) != 0 {
		t.Errorf("expected no states for X with nil column lookup, got %v", states)
	}
}

// TestSetStatesFromDataImpl_MixedPaths exercises both the existing-states-skip
// and nil-column paths in a single call.
func TestSetStatesFromDataImpl_MixedPaths(t *testing.T) {
	bn := models.NewBayesianNetwork()
	bn.AddNode("A")
	bn.SetStates("A", []string{"s0", "s1"}) // has existing states
	bn.AddNode("B")                         // no states, column will be nil

	m := map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
	}
	data := tabgo.NewDataFrame(m)

	setStatesFromDataImpl(bn, data, &nilColumnLookup{})

	statesA := bn.GetStates("A")
	if len(statesA) != 2 || statesA[0] != "s0" {
		t.Errorf("expected original states for A, got %v", statesA)
	}
	statesB := bn.GetStates("B")
	if len(statesB) != 0 {
		t.Errorf("expected no states for B, got %v", statesB)
	}
}

// ---------------------------------------------------------------------------
// dataFrameColumnLookup — test the wrapper's panic recovery
// ---------------------------------------------------------------------------

// TestDataFrameColumnLookup_MissingColumn exercises the panic recovery path
// in dataFrameColumnLookup for columns not present in the DataFrame.
func TestDataFrameColumnLookup_MissingColumn(t *testing.T) {
	m := map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
	}
	data := tabgo.NewDataFrame(m)
	lookup := &dataFrameColumnLookup{df: data}

	col := lookup.Column("NonExistent")
	if col != nil {
		t.Error("expected nil for missing column")
	}
}

// TestDataFrameColumnLookup_ExistingColumn verifies existing columns are returned.
func TestDataFrameColumnLookup_ExistingColumn(t *testing.T) {
	m := map[string]*tabgo.Series{
		"A": tabgo.NewSeries("A", []any{0, 1}),
	}
	data := tabgo.NewDataFrame(m)
	lookup := &dataFrameColumnLookup{df: data}

	col := lookup.Column("A")
	if col == nil {
		t.Error("expected non-nil for existing column")
	}
}

// ---------------------------------------------------------------------------
// writeBIFFileImpl — write error via failing writer mock
// ---------------------------------------------------------------------------

func TestWriteBIFFileImpl_WriteError(t *testing.T) {
	bn := models.NewBayesianNetwork()
	bn.AddNode("X")
	bn.SetStates("X", []string{"a", "b"})
	bn.GetRandomCPDs(2, 0)

	outPath := filepath.Join(t.TempDir(), "model.bif")
	err := writeBIFFileImpl(outPath, bn, &failingBIFWriter{})
	if err == nil {
		t.Error("expected error from failing writer")
	}
	// Clean up created file.
	os.Remove(outPath)
}

func TestWriteBIFFileImpl_BadPath(t *testing.T) {
	bn := models.NewBayesianNetwork()
	err := writeBIFFileImpl("/nonexistent/dir/model.bif", bn, defaultBIFWriter)
	if err == nil {
		t.Error("expected error for bad path")
	}
}
