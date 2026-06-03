//go:build unit

package example_models

import (
	"fmt"
	"testing"

	"github.com/asymmetric-effort/pgmgo/src/models"
)

func TestStudentCheckModel(t *testing.T) {
	bn := Student()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("Student CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 5 {
		t.Fatalf("Student should have 5 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 4 {
		t.Fatalf("Student should have 4 edges, got %d", len(bn.Edges()))
	}
}

func TestAsiaCheckModel(t *testing.T) {
	bn := Asia()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("Asia CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 8 {
		t.Fatalf("Asia should have 8 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 8 {
		t.Fatalf("Asia should have 8 edges, got %d", len(bn.Edges()))
	}
}

func TestAlarmCheckModel(t *testing.T) {
	bn := Alarm()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("Alarm CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 5 {
		t.Fatalf("Alarm should have 5 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 4 {
		t.Fatalf("Alarm should have 4 edges, got %d", len(bn.Edges()))
	}
}

func TestCancerCheckModel(t *testing.T) {
	bn := Cancer()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("Cancer CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 5 {
		t.Fatalf("Cancer should have 5 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 4 {
		t.Fatalf("Cancer should have 4 edges, got %d", len(bn.Edges()))
	}
}

func TestSachsStructure(t *testing.T) {
	bn := Sachs()
	if len(bn.Nodes()) != 11 {
		t.Fatalf("Sachs should have 11 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 17 {
		t.Fatalf("Sachs should have 17 edges, got %d", len(bn.Edges()))
	}
	// Sachs has no CPDs, so CheckModel should fail
	if err := bn.CheckModel(); err == nil {
		t.Fatal("Sachs CheckModel should fail (no CPDs)")
	}
}

func TestWaterSprinklerCheckModel(t *testing.T) {
	bn := WaterSprinkler()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("WaterSprinkler CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 4 {
		t.Fatalf("WaterSprinkler should have 4 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 4 {
		t.Fatalf("WaterSprinkler should have 4 edges, got %d", len(bn.Edges()))
	}
}

// --- New small models with full CPDs ---

func TestSurveyCheckModel(t *testing.T) {
	bn := Survey()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("Survey CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 6 {
		t.Fatalf("Survey should have 6 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 6 {
		t.Fatalf("Survey should have 6 edges, got %d", len(bn.Edges()))
	}
}

func TestMontyHallCheckModel(t *testing.T) {
	bn := MontyHall()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("MontyHall CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 3 {
		t.Fatalf("MontyHall should have 3 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 2 {
		t.Fatalf("MontyHall should have 2 edges, got %d", len(bn.Edges()))
	}
}

func TestDogProblemCheckModel(t *testing.T) {
	bn := DogProblem()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("DogProblem CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 5 {
		t.Fatalf("DogProblem should have 5 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 4 {
		t.Fatalf("DogProblem should have 4 edges, got %d", len(bn.Edges()))
	}
}

func TestFraudDetectionCheckModel(t *testing.T) {
	bn := FraudDetection()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("FraudDetection CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 6 {
		t.Fatalf("FraudDetection should have 6 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 6 {
		t.Fatalf("FraudDetection should have 6 edges, got %d", len(bn.Edges()))
	}
}

func TestMedicalDiagnosisCheckModel(t *testing.T) {
	bn := MedicalDiagnosis()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("MedicalDiagnosis CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 8 {
		t.Fatalf("MedicalDiagnosis should have 8 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 8 {
		t.Fatalf("MedicalDiagnosis should have 8 edges, got %d", len(bn.Edges()))
	}
}

func TestEarthquakeCheckModel(t *testing.T) {
	bn := Earthquake()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("Earthquake CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 5 {
		t.Fatalf("Earthquake should have 5 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 4 {
		t.Fatalf("Earthquake should have 4 edges, got %d", len(bn.Edges()))
	}
}

func TestCoinTossCheckModel(t *testing.T) {
	bn := CoinToss()
	if err := bn.CheckModel(); err != nil {
		t.Fatalf("CoinToss CheckModel failed: %v", err)
	}
	if len(bn.Nodes()) != 3 {
		t.Fatalf("CoinToss should have 3 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 2 {
		t.Fatalf("CoinToss should have 2 edges, got %d", len(bn.Edges()))
	}
}

// --- Large structure-only models ---

func TestChildStructure(t *testing.T) {
	bn := Child()
	if len(bn.Nodes()) != 20 {
		t.Fatalf("Child should have 20 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 25 {
		t.Fatalf("Child should have 25 edges, got %d", len(bn.Edges()))
	}
}

func TestInsuranceStructure(t *testing.T) {
	bn := Insurance()
	if len(bn.Nodes()) != 27 {
		t.Fatalf("Insurance should have 27 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 52 {
		t.Fatalf("Insurance should have 52 edges, got %d", len(bn.Edges()))
	}
}

func TestAlarmFullStructure(t *testing.T) {
	bn := AlarmFull()
	if len(bn.Nodes()) != 37 {
		t.Fatalf("AlarmFull should have 37 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 46 {
		t.Fatalf("AlarmFull should have 46 edges, got %d", len(bn.Edges()))
	}
}

func TestWaterStructure(t *testing.T) {
	bn := Water()
	if len(bn.Nodes()) != 32 {
		t.Fatalf("Water should have 32 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 66 {
		t.Fatalf("Water should have 66 edges, got %d", len(bn.Edges()))
	}
}

func TestMildewStructure(t *testing.T) {
	bn := Mildew()
	if len(bn.Nodes()) != 35 {
		t.Fatalf("Mildew should have 35 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 46 {
		t.Fatalf("Mildew should have 46 edges, got %d", len(bn.Edges()))
	}
}

func TestBarleyStructure(t *testing.T) {
	bn := Barley()
	if len(bn.Nodes()) != 48 {
		t.Fatalf("Barley should have 48 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 84 {
		t.Fatalf("Barley should have 84 edges, got %d", len(bn.Edges()))
	}
}

func TestHailfinderStructure(t *testing.T) {
	bn := Hailfinder()
	if len(bn.Nodes()) != 56 {
		t.Fatalf("Hailfinder should have 56 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 66 {
		t.Fatalf("Hailfinder should have 66 edges, got %d", len(bn.Edges()))
	}
}

func TestHepar2Structure(t *testing.T) {
	bn := Hepar2()
	if len(bn.Nodes()) != 70 {
		t.Fatalf("Hepar2 should have 70 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 123 {
		t.Fatalf("Hepar2 should have 123 edges, got %d", len(bn.Edges()))
	}
}

func TestWin95ptsStructure(t *testing.T) {
	bn := Win95pts()
	if len(bn.Nodes()) != 76 {
		t.Fatalf("Win95pts should have 76 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 112 {
		t.Fatalf("Win95pts should have 112 edges, got %d", len(bn.Edges()))
	}
}

func TestPathfinderStructure(t *testing.T) {
	bn := Pathfinder()
	if len(bn.Nodes()) != 109 {
		t.Fatalf("Pathfinder should have 109 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 195 {
		t.Fatalf("Pathfinder should have 195 edges, got %d", len(bn.Edges()))
	}
}

func TestPigsStructure(t *testing.T) {
	bn := Pigs()
	if len(bn.Nodes()) != 441 {
		t.Fatalf("Pigs should have 441 nodes, got %d", len(bn.Nodes()))
	}
	if len(bn.Edges()) != 592 {
		t.Fatalf("Pigs should have 592 edges, got %d", len(bn.Edges()))
	}
}

// --- Registry tests ---

func TestList(t *testing.T) {
	names := List()
	if len(names) < 25 {
		t.Fatalf("List() should return at least 25 models, got %d", len(names))
	}
	// Verify sorted
	for i := 1; i < len(names); i++ {
		if names[i] < names[i-1] {
			t.Fatalf("List() not sorted: %q before %q", names[i-1], names[i])
		}
	}
}

func TestGetKnown(t *testing.T) {
	for _, name := range List() {
		bn, err := Get(name)
		if err != nil {
			t.Fatalf("Get(%q) failed: %v", name, err)
		}
		if bn == nil {
			t.Fatalf("Get(%q) returned nil", name)
		}
		if len(bn.Nodes()) == 0 {
			t.Fatalf("Get(%q) returned empty network", name)
		}
	}
}

func TestGetCaseInsensitive(t *testing.T) {
	bn, err := Get("STUDENT")
	if err != nil {
		t.Fatalf("Get(STUDENT) should be case-insensitive: %v", err)
	}
	if len(bn.Nodes()) != 5 {
		t.Fatalf("Get(STUDENT) returned wrong model")
	}
}

func TestGetUnknown(t *testing.T) {
	_, err := Get("nonexistent_model")
	if err == nil {
		t.Fatal("Get(nonexistent_model) should return error")
	}
}

// --- Panic-path coverage for helper functions ---

func TestMustPanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("must(err) should panic on non-nil error")
		}
	}()
	must(fmt.Errorf("test error"))
}

func TestMustCPDPanicsOnBadInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("mustCPD should panic on invalid input")
		}
	}()
	// variableCard=0 is invalid and will cause NewTabularCPD to return an error
	mustCPD("X", 0, nil, nil, nil)
}

func TestMustCheckPanicsOnInvalidModel(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("mustCheck should panic when CheckModel fails")
		}
	}()
	// Build a network with a node but no CPD — CheckModel will fail
	bn := models.NewBayesianNetwork()
	_ = bn.AddNode("A")
	_ = bn.AddNode("B")
	_ = bn.AddEdge("A", "B")
	mustCheck(bn, "TestBroken")
}
