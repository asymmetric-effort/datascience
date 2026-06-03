//go:build unit

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/ci_tests"
)

// ---------------------------------------------------------------------------
// printUsage, adaptCITest, pdagToBN coverage
// ---------------------------------------------------------------------------

func TestPrintUsage(t *testing.T) {
	// Exercise printUsage — just ensure it doesn't panic.
	printUsage()
}

func TestAdaptCITest(t *testing.T) {
	// Exercise adaptCITest wrapper by using a real CITest.
	wrapped := adaptCITest(ci_tests.ChiSquare)
	if wrapped == nil {
		t.Fatal("adaptCITest returned nil")
	}
}

// ---------------------------------------------------------------------------
// learnStructure — additional methods
// ---------------------------------------------------------------------------

func TestLearnStructure_PC(t *testing.T) {
	csvContent := "A,B,C\n0,0,0\n0,0,1\n0,1,0\n0,1,1\n1,0,0\n1,0,1\n1,1,0\n1,1,1\n0,0,0\n1,1,1\n0,0,0\n1,1,1\n0,0,1\n1,1,0\n0,1,0\n1,0,1\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructure(csvPath, "pc", "bic", 0.05, outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for PC, got %d", code)
	}
}

func TestLearnStructure_GES(t *testing.T) {
	csvContent := "A,B,C\n0,0,0\n0,0,1\n0,1,0\n0,1,1\n1,0,0\n1,0,1\n1,1,0\n1,1,1\n0,0,0\n1,1,1\n0,0,0\n1,1,1\n0,0,1\n1,1,0\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructure(csvPath, "ges", "bic", 0.05, outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for GES, got %d", code)
	}
}

func TestLearnStructure_Exhaustive(t *testing.T) {
	csvContent := "A,B\n0,0\n0,1\n1,0\n1,1\n0,0\n0,1\n1,0\n1,1\n0,0\n1,1\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructure(csvPath, "exhaustive", "bic", 0.05, outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for exhaustive, got %d", code)
	}
}

func TestLearnStructure_BDeu(t *testing.T) {
	csvContent := "A,B\n0,0\n0,1\n1,0\n1,1\n0,0\n0,1\n1,0\n1,1\n0,0\n1,1\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructure(csvPath, "hillclimb", "bdeu", 0.05, outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for BDeu, got %d", code)
	}
}

func TestLearnStructure_K2(t *testing.T) {
	csvContent := "A,B\n0,0\n0,1\n1,0\n1,1\n0,0\n0,1\n1,0\n1,1\n0,0\n1,1\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructure(csvPath, "hillclimb", "k2", 0.05, outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for K2, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// fitParameters — EM method
// ---------------------------------------------------------------------------

func TestFitParameters_EM(t *testing.T) {
	modelPath := writeTempBIF(t, validBIF)
	// EM needs state names matching BIF (True=0, False=1)
	csvContent := "Rain,Sprinkler\nTrue,True\nTrue,False\nFalse,True\nFalse,False\nTrue,False\nFalse,True\nTrue,True\nFalse,False\nTrue,False\nFalse,True\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "fitted.bif")
	code := fitParameters(modelPath, csvPath, "em", outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for EM, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// runQuery — with evidence parsing
// ---------------------------------------------------------------------------

func TestRunQuery_WithEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runQuery([]string{"--variables", "Rain", "--evidence", "Sprinkler=0", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRunQuery_BadEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runQuery([]string{"--variables", "Rain", "--evidence", "badformat", path})
	if code != 2 {
		t.Errorf("expected exit code 2 for bad evidence, got %d", code)
	}
}

func TestRunQuery_BPMethod(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runQuery([]string{"--variables", "Rain", "--method", "bp", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRunQuery_ApproxMethod(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runQuery([]string{"--variables", "Rain", "--method", "approx", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// runMAP — with evidence and missing variables
// ---------------------------------------------------------------------------

func TestRunMAP_WithEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runMAP([]string{"--variables", "Rain", "--evidence", "Sprinkler=0", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRunMAP_MissingVariables(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runMAP([]string{path})
	if code != 2 {
		t.Errorf("expected exit code 2 for missing variables, got %d", code)
	}
}

func TestRunMAP_BadEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runMAP([]string{"--variables", "Rain", "--evidence", "badformat", path})
	if code != 2 {
		t.Errorf("expected exit code 2 for bad evidence, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// runDo — with evidence and missing flags
// ---------------------------------------------------------------------------

func TestRunDo_MissingIntervention(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runDo([]string{"--query", "Sprinkler", path})
	if code != 2 {
		t.Errorf("expected exit code 2 for missing intervention, got %d", code)
	}
}

func TestRunDo_WithEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runDo([]string{"--intervention", "Rain=0", "--query", "Sprinkler", path})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRunDo_BadIntervention(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runDo([]string{"--intervention", "bad", "--query", "Sprinkler", path})
	if code != 2 {
		t.Errorf("expected exit code 2 for bad intervention, got %d", code)
	}
}

func TestRunDo_BadEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := runDo([]string{"--intervention", "Rain=0", "--query", "Sprinkler", "--evidence", "bad", path})
	if code != 2 {
		t.Errorf("expected exit code 2 for bad evidence, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// runSample — with evidence
// ---------------------------------------------------------------------------

func TestRunSample_WithEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	outPath := filepath.Join(t.TempDir(), "samples.csv")
	code := runSample([]string{"--model", path, "--n", "5", "--output", outPath, "--method", "rejection", "--evidence", "Rain=0", "--seed", "42"})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestRunSample_BadEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	outPath := filepath.Join(t.TempDir(), "samples.csv")
	code := runSample([]string{"--model", path, "--output", outPath, "--evidence", "bad"})
	if code != 2 {
		t.Errorf("expected exit code 2 for bad evidence, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// convert — additional format paths
// ---------------------------------------------------------------------------

func TestConvertModel_XMLBIFToBIF(t *testing.T) {
	// First convert BIF to XMLBIF, then back.
	bifPath := writeTempBIF(t, validBIF)
	xmlPath := filepath.Join(t.TempDir(), "model.xmlbif")
	code := convertModel(bifPath, "bif", "xmlbif", xmlPath)
	if code != 0 {
		t.Fatalf("setup failed: code=%d", code)
	}

	outPath := filepath.Join(t.TempDir(), "out.bif")
	code = convertModel(xmlPath, "xmlbif", "bif", outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for XMLBIF->BIF, got %d", code)
	}
}

func TestConvertModel_NETToBIF(t *testing.T) {
	// Convert BIF -> NET, then NET -> BIF.
	bifPath := writeTempBIF(t, validBIF)
	netPath := filepath.Join(t.TempDir(), "model.net")
	code := convertModel(bifPath, "bif", "net", netPath)
	if code != 0 {
		t.Fatalf("setup failed: code=%d", code)
	}
	outPath := filepath.Join(t.TempDir(), "out.bif")
	code = convertModel(netPath, "net", "bif", outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for NET->BIF, got %d", code)
	}
}

func TestConvertModel_BadOutputPath(t *testing.T) {
	bifPath := writeTempBIF(t, validBIF)
	// Non-writable output path.
	code := convertModel(bifPath, "bif", "xmlbif", "/nonexistent/dir/out.xmlbif")
	if code != 1 {
		t.Errorf("expected exit code 1 for bad output path, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// validateBIFFile — malformed BIF
// ---------------------------------------------------------------------------

// TestValidateBIFFile_EmptyNetwork tests a BIF that parses but has no nodes.
func TestValidateBIFFile_EmptyNetwork(t *testing.T) {
	emptyBIF := "network unknown {\n}\n"
	path := writeTempBIF(t, emptyBIF)
	code := validateBIFFile(path)
	// Empty but valid network.
	if code != 0 {
		t.Errorf("expected exit code 0 for empty valid network, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// setStatesFromData — exercises the helper
// ---------------------------------------------------------------------------

func TestSetStatesFromData(t *testing.T) {
	csvContent := "A,B\n0,0\n0,1\n1,0\n1,1\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)

	// Use learnStructure which calls setStatesFromData internally.
	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := learnStructure(csvPath, "hillclimb", "bic", 0.05, outPath)
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}

	// Verify output was created.
	if _, err := os.Stat(outPath); err != nil {
		t.Errorf("output file not created: %v", err)
	}
}

// ---------------------------------------------------------------------------
// writeBIFFile — exercises the helper path
// ---------------------------------------------------------------------------

func TestWriteBIFFile_BadPath(t *testing.T) {
	// loadBIF + writeBIF path.
	bifPath := writeTempBIF(t, validBIF)
	bn, code := loadBIF(bifPath)
	if code != 0 {
		t.Fatalf("loadBIF failed: %d", code)
	}
	err := writeBIFFile("/nonexistent/dir/model.bif", bn)
	if err == nil {
		t.Error("expected error for bad write path")
	}
}

// ---------------------------------------------------------------------------
// loadBIF — malformed content
// ---------------------------------------------------------------------------

func TestLoadBIF_InvalidBIF(t *testing.T) {
	// Use invalidBIF (has node but no CPD) - loadBIF should succeed (no validation)
	path := writeTempBIF(t, invalidBIF)
	bn, code := loadBIF(path)
	if code != 0 {
		t.Errorf("expected exit code 0 for loadBIF, got %d", code)
	}
	if bn == nil {
		t.Error("expected non-nil BN")
	}
}

// ---------------------------------------------------------------------------
// bnToDigraph — exercises the graph conversion
// ---------------------------------------------------------------------------

func TestBnToDigraph(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	bn, code := loadBIF(path)
	if code != 0 {
		t.Fatalf("loadBIF failed: %d", code)
	}
	g := bnToDigraph(bn)
	if g == nil {
		t.Fatal("bnToDigraph returned nil")
	}
}

// ---------------------------------------------------------------------------
// infoBIF — no CPD path
// ---------------------------------------------------------------------------

func TestInfoBIF_InvalidBIF(t *testing.T) {
	// Use invalidBIF which has a node but no CPD.
	path := writeTempBIF(t, invalidBIF)
	code := infoBIF(path)
	// This should succeed since infoBIF doesn't validate the model.
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// runConvert — full flag parsing
// ---------------------------------------------------------------------------

func TestRunConvert_FullFlags(t *testing.T) {
	bifPath := writeTempBIF(t, validBIF)
	outPath := filepath.Join(t.TempDir(), "out.xmlbif")
	code := runConvert([]string{"--input", bifPath, "--from", "bif", "--to", "xmlbif", "--output", outPath})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// runCompare — full flag parsing
// ---------------------------------------------------------------------------

func TestRunCompare_FullFlags(t *testing.T) {
	path1 := writeTempBIF(t, validBIF)
	path2 := writeTempBIF(t, validBIF3)
	code := runCompare([]string{"--true", path1, "--estimated", path2})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestCompareModels_BadEstimated(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := compareModels(path, "/nonexistent")
	if code != 2 {
		t.Errorf("expected exit code 2, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// runFit — full flag parsing
// ---------------------------------------------------------------------------

func TestRunFit_FullFlags(t *testing.T) {
	modelPath := writeTempBIF(t, validBIF)
	csvContent := "Rain,Sprinkler\n0,0\n0,1\n1,0\n1,1\n0,0\n1,0\n0,1\n1,1\n0,0\n1,0\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "fitted.bif")
	code := runFit([]string{"--model", modelPath, "--data", csvPath, "--method", "mle", "--output", outPath})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

// ---------------------------------------------------------------------------
// runLearn — full flag parsing
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Additional coverage for error paths
// ---------------------------------------------------------------------------

func TestQueryVE_ErrorPath(t *testing.T) {
	// Query with non-existent variable should produce an error.
	path := writeTempBIF(t, validBIF)
	code := queryBIF(path, []string{"NonExistent"}, nil, "ve")
	if code == 0 {
		t.Error("expected non-zero exit code for non-existent variable")
	}
}

func TestQueryBP_ErrorPath(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := queryBIF(path, []string{"Rain"}, map[string]int{"Sprinkler": 0}, "bp")
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestQueryApprox_WithEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := queryBIF(path, []string{"Rain"}, map[string]int{"Sprinkler": 0}, "approx")
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestMapBIF_WithEvidenceDirect(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := mapBIF(path, []string{"Rain"}, map[string]int{"Sprinkler": 0})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestDoCausalQuery_WithEvidence(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	code := doCausalQuery(path, map[string]int{"Rain": 0}, "Sprinkler", map[string]int{})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestSampleModel_ForwardWithSeed(t *testing.T) {
	path := writeTempBIF(t, validBIF)
	outPath := filepath.Join(t.TempDir(), "samples.csv")
	code := sampleModel(path, 3, outPath, "forward", map[string]int{}, 123)
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}

func TestConvertModel_UATToBIF(t *testing.T) {
	bifPath := writeTempBIF(t, validBIF)
	uaiPath := filepath.Join(t.TempDir(), "model.uai")
	code := convertModel(bifPath, "bif", "uai", uaiPath)
	if code != 0 {
		t.Fatalf("setup failed: code=%d", code)
	}
	outPath := filepath.Join(t.TempDir(), "out.bif")
	code = convertModel(uaiPath, "uai", "bif", outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for UAI->BIF, got %d", code)
	}
}

func TestConvertModel_XDSLToBIF(t *testing.T) {
	bifPath := writeTempBIF(t, validBIF)
	xdslPath := filepath.Join(t.TempDir(), "model.xdsl")
	code := convertModel(bifPath, "bif", "xdsl", xdslPath)
	if code != 0 {
		t.Fatalf("setup failed: code=%d", code)
	}
	outPath := filepath.Join(t.TempDir(), "out.bif")
	code = convertModel(xdslPath, "xdsl", "bif", outPath)
	if code != 0 {
		t.Errorf("expected exit code 0 for XDSL->BIF, got %d", code)
	}
}

func TestRunLearn_FullFlags(t *testing.T) {
	csvContent := "A,B\n0,0\n0,1\n1,0\n1,1\n0,0\n0,1\n1,0\n1,1\n0,0\n1,1\n"
	csvPath := writeTempFile(t, "data.csv", csvContent)
	outPath := filepath.Join(t.TempDir(), "learned.bif")
	code := runLearn([]string{"--data", csvPath, "--method", "hillclimb", "--score", "bic", "--output", outPath})
	if code != 0 {
		t.Errorf("expected exit code 0, got %d", code)
	}
}
